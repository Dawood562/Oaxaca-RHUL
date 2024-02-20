// Will hold currently displayed version of menu
var currentMenu;
var editMode = false;
// Called when menu page is initially loaded
function initMenuAll() {
  if (editMode) {
    editMenu()
  }

  let data = requestMenu(0, 0, 0); // Zero value = none specified
  data.then(r => {
    console.log(r)
    currentMenu = r
    let index = 0;
    document.getElementById("MenuItemGridLayout").innerHTML = ""
    r.forEach(element => {
      document.getElementById("MenuItemGridLayout").innerHTML += createMenuItem(index, element.itemId, element.itemName, element.price, element.calories);
      index++;
    });
  })
}

// Function that takes in data and turns into menu item to be displayed
function createMenuItem(index, id, itemName, price, calories) {
  let comp = "<div class='MenuItemDiv' id='item" + index + "'> <img class='MenuItemImg' src='image/foodimg.jpg'><br> <div class='MenuItemDetails'><label class='MenuItemName'>" + itemName + "</label><br><label class='MenuItemPrice'>£" + price.toFixed(2) + "</label><label class='MenuItemCalories'>" + calories + "kcal</label></div>";
  comp += "<input type='number' id='itemQuantityInput" + index + "' min='1' value='1' class='itemQuantityInput'>";
  comp += "<button id='addToBasketButton" + index + "' onclick='addToBasket(" + index + ", " + id + ", \"" + itemName + "\", " + price + ", " + calories + ")'>Add to Basket</button></a></div>";

  return comp;
}


// Toggle edit mode
function editMenu() {

  if (!editMode) {
    document.getElementById("MenuEditSection").innerHTML += "<div id='newItemDiv'><h3>Add new menu item:</h3><p><label>Name:</label><input type='text' id='newItemNameField'>  <label>Price:</label><input type='text' id='newItemPriceField'>  <label>Calories:</label><input type='text' id='newItemCaloriesField'>    <button onclick='addMenuItem()'>+</button></p></div>";
    document.getElementById("MenuEditSection").innerHTML += "<div id='removeItemDiv'><h3>Delete menu item:</h3><p><label>Enter the name of the item you want to delete from the menu:</label><input type='text' id='deleteItemNameField'> <button onclick='deleteMenuItem()'>-</button></p></div>";

    // Add edit button to each menu item
    for (let i = 0; i < currentMenu.length; i++) {
      document.getElementById("item" + i).innerHTML += "<button index='" + i + "' id='itemEditButton' class='editMenuItem'>EDIT</button>"
      let thing = document.getElementById("item"+i).childNodes[5].remove()
    }

    // Retrieve all edit buttons into iterable
    const editButtons = document.querySelectorAll(".editMenuItem")

    // For each edit button add click listener
    editButtons.forEach(item => {
      item.addEventListener('click', function() {
        EditMenuField(item.getAttribute("index"))
      })
    })

    editMode = true;
  } else {
    document.getElementById("newItemDiv").remove();
    document.getElementById("removeItemDiv").remove();
    const editButtons = document.querySelectorAll('.editMenuItem');
    editButtons.forEach(item => {
      item.remove();
    })
    editMode = false;
  }
}


// Holds global variable on if a field is currently being edited
var editingText = false;
// Holds location of item being edited
var currentlyEditing = -1;
var oldText

// Should replace name with label displaying name: and a textfield
// Should replace price and calories of price label and textfield too
function EditMenuField(index) {
  // Index 0 = item name
  // Index 2 = item price
  // Index 3 = item calories
  let itemToEdit = document.getElementById("item" + index).childNodes[4].childNodes;


  // Replace name with name tag and checkbox
  itemToEdit[0].innerHTML = "<label class='editMenuItemPrompt'>New name:</label><input id='nameEditPrompt' class='editMenuItemPrompt' type='text'>";
  itemToEdit[2].innerHTML = "<label class='editMenuItemPrompt'>New price:</label><input id='priceEditPrompt' class='editMenuItemPrompt' type='text'>";
  itemToEdit[3].innerHTML = "<label class='editMenuItemPrompt'>New calories:</label><input id='caloriesEditPrompt' class='editMenuItemPrompt' type='text'>";

  document.getElementById("itemEditButton").remove();
  document.getElementById("item" + index).innerHTML += "<button index='" + index + "' id='itemEditButton32' class='editMenuItem'>Submit</button>"
  document.getElementById("itemEditButton32").addEventListener('click', function() {
    submitMenuEdit(index);
      let toRemove = document.getElementsByClassName("editMenuItemPrompt")
        toRemove[0].innerHTML = "<label>Submitted!</label>"
        for(let i=1; i<toRemove.length; i++){
            toRemove[i].innerHTML = ""
        }
        document.getElementById("nameEditPrompt").remove()
        document.getElementById("priceEditPrompt").remove()
        document.getElementById("caloriesEditPrompt").remove()
        document.getElementById("itemEditButton32").remove()

    })


  // MUST REPLACE EDIT BUTTON TO BE CALLED SUBMIT AND MAKE IT CALL submitMenuEdit() FUNCTION
  // SHOULD THEN REPLACE ALL THE BOXES WITH A SIMPLE "SUCCESSFUL" LABEL TO INDICATE CHANGE
}

// Needs to take in index so that it can get id
async function submitMenuEdit(index) {
  let id = currentMenu[index].itemId;
  let name = document.getElementById("nameEditPrompt").value;
  let itemPrice = parseFloat(document.getElementById("priceEditPrompt").value);
  let itemCalories = parseInt(document.getElementById("caloriesEditPrompt").value);
  console.log(id + "," + name + "," + itemPrice + "," + itemCalories);
  console.log(document.getElementById("nameEditPrompt"));

  try {
    let response = await fetch("http://localhost:4444/edit_item", {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        itemId: id,
        itemName: name,
        price: itemPrice,
        calories: itemCalories
      })
    })

  } catch (error) {
    console.error(error);
  }
}

// Functions to add and delete menu items
async function addMenuItem(){
    let nameValue = document.getElementById("newItemNameField").value;
    let priceValue = parseFloat(document.getElementById("newItemPriceField").value);
    let caloriesValue = parseInt(document.getElementById("newItemCaloriesField").value);
    let result = await addItemToDB(nameValue, priceValue, caloriesValue);

    if(result >= 0){
        document.getElementById("MenuItemGridLayout").innerHTML+= createMenuItem(currentMenu.length, nameValue, priceValue, caloriesValue);
    }
    
}

function deleteMenuItem(){
    let nameValue = document.getElementById("deleteItemNameField").value;
    removeItem(nameValue).then(r => {
        if (r >= 0){
            // First childNodes holds each item on menu
            // Second childNodes holds the details of the first childnodes menu item
            // Third childNodes holds the actual details. index 0 is name
            for (let i = 0; i < document.getElementById("MenuItemGridLayout").childNodes.length; i++){
                if (document.getElementById("MenuItemGridLayout").childNodes[i].childNodes[4].childNodes[0].innerHTML == nameValue){
                    let gridLayout = document.getElementById("MenuItemGridLayout");
                    let itemToRemove = document.getElementById("MenuItemGridLayout").childNodes[0];
                    gridLayout.removeChild(itemToRemove);
                }
            }
        }
    })
    
}

// Add menu item
async function addItemToDB(name, _price, _calories) {
  try {
    let response = await fetch("http://localhost:4444/add_item", {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        itemName: name,
        price: _price,
        calories: _calories
      })
    })
     return 0
  } catch (error) {
    console.error(error)
     return -1
  }
}

// Takes in name and returns id if found
function getIdFromName(name) {
  for (let i = 0; i < currentMenu.length; i++) {
    if (currentMenu[i].itemName == name) {
      return currentMenu[i].itemId;
    }
  }
  console.error("COULD NOT FIND ID FROM ITEM NAME: " + name);
  return null;
}

// Contacts backend to remove menu item
async function removeItem(name) {
  console.log(currentMenu)
  nameId = String(getIdFromName(name))
  try {
    let response = await fetch("http://localhost:4444/remove_item?" + new URLSearchParams({
      itemId: nameId
    }).toString(), {
      method: "DELETE"
    })
return 0
  } catch (error) {
    console.error(error)
    return -1
  }
}

// Fetches data from backend or throws error to console and provides example menu data
async function requestMenu(userSearchTerm, userMaxPrice, userMaxCalories) {
  try {
    let response = await fetch("http://localhost:4444/menu?" + new URLSearchParams({
      searchTerm: userSearchTerm,
      Price: userMaxPrice,
      Calories: userMaxCalories
    }).toString())

    if (!response.ok) {
      console.log("ERROR fetching menu")
    }

    let data = await response.json()
    return data
  } catch (error) {
    console.error(error)
    // Return example error if unable to connect to backend
    return (
      [{
        itemName: "Tequila",
        price: 6.90,
        calories: 12,
        Type: "DRINK"
      }, {
        itemName: "Olives",
        price: 7.99,
        calories: 165,
        Type: "APPETIZER"
      }, {
        itemName: "Mozzarella Sticks",
        price: 8.99,
        calories: 349,
        Type: "ENTREES"
      }, {
        itemName: "Ice Cream",
        price: 7.42,
        calories: 632,
        Type: "DESSERT"
      }])
  }
}
// function to add menu item to basket
// Stores menu items in basket using cookies
// Cookie structure is CSV in form: id,itemName,price,calories,quantity
function addToBasket(index, itemId, itemName, price, calories) {
  let quantity = document.getElementById("itemQuantityInput"+index).value

  // This will work for now as we only store 1 type of cookie
  let previousCookieContent = document.cookie.split("basket=")[1]


  let updated = false;


  // Check that item is not already in basket
  previousCookieContent.split("#").forEach(element => {
    let splitCookie = element.split(",")
    // Item found in basket
    if(splitCookie[0] == itemId){
      // Then update quantity instead
      let newQuantity = Number(splitCookie[4])+Number(quantity);
      
      let indexOfItem = previousCookieContent.indexOf(element)
      let indexOfEndOfItem = previousCookieContent.indexOf("#", indexOfItem)
      let updatedCookieSegment = previousCookieContent.substring(indexOfItem, indexOfEndOfItem-1)+newQuantity
      let updatedCookie = previousCookieContent.substring(0,indexOfItem)+updatedCookieSegment+previousCookieContent.substring(indexOfEndOfItem, previousCookieContent.length);
      document.cookie = "basket="+updatedCookie;

      updated = true;
    }
  });

  if(!updated){
    document.cookie="basket="+itemId+","+itemName+","+price+","+calories+","+quantity+"#"+previousCookieContent;
  }

  console.log("Current basket:"+document.cookie);
  
}

//function to update basket icon with the item quantity in order
function updateBasketIcon() {
 let order = JSON.parse(localStorage.getItem('order')) || [];
 let basketIcon = document.getElementById('basketIcon');
 let totalQuantity = order.reduce((total, item) => total + item.quantity, 0);
 basketIcon.textContent = `🛒 ${totalQuantity}`;
}
//function to update order details with items in the order
function updateOrderDetails() {
 let order = JSON.parse(localStorage.getItem('order'));
 let orderDetailsDiv = document.getElementById('orderDetails');
 let totalDiv = document.getElementById('orderTotal');
 orderDetailsDiv.innerHTML = '';
 let orderTotal = 0;
 if (order && order.length > 0) {
   order.forEach(item => {
     orderTotal += item.price * item.quantity;
     let li = document.createElement('li');
     li.innerHTML = `
              <h3>${item.itemName}</h3>
              <p> quantity: ${item.quantity}</p>
              <p>Calories: ${item.calories * item.quantity} kcal</p>
              <p>Price: £${(item.price * item.quantity).toFixed(2)}</p>
              <button class="removeButton"><i class = "fa fa-trash"></i></button>
          `;
     orderDetailsDiv.appendChild(li);
   });
   totalDiv.textContent = `Total: £${orderTotal.toFixed(2)}`;
   let removeButtons = document.querySelectorAll('.removeButton');
   removeButtons.forEach(button => {
     button.addEventListener('click', () => {
       removeFromOrder(button.parentElement.querySelector('h3').textContent);
     });
   });
 } else {
   orderDetailsDiv.innerHTML = `basket empty`;
   totalDiv.textContent = '';
 }
}
//function ro remove menu item from the order
function removeFromOrder(itemName) {
 let order = JSON.parse(localStorage.getItem('order')) || [];
 const itemIndex = order.findIndex(item => item.itemName === itemName);
 if (itemIndex >= 0) {
   order.splice(itemIndex, 1);
   localStorage.setItem('order', JSON.stringify(order));
   updateOrderDetails();
 }
}
//updates order details and basket icon when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
  initializeMenu();
  initializeCategoryFilter();
  updateOrderDetails();
  updateBasketIcon();
});

function refreshCookies(){
  // Only refresh cookies if loading menu page again.
  // Should keep cookies when redirected to order page
  if (document.title.indexOf("Menu") != -1){
    document.cookie="basket=" // Reset basket on refresh or new session
    console.log("Refreshed cookies!")
  }
}

function filterItems(){
   let searchTerm = document.getElementById('searchTerm').value;
   let maxCalories = parseInt(document.getElementById('maxCalories').value) || 0;
   let maxPrice = parseFloat(document.getElementById('maxPrice').value) || 0;
   if (searchTerm.length <= 0){
       searchTerm = 0;
   }
   let data = requestMenu(searchTerm, maxPrice, maxCalories);
   data.then(r =>{
       let index = 0;
       document.getElementById("MenuItemGridLayout").innerHTML = "";
       r.forEach(element => {
           document.getElementById("MenuItemGridLayout").innerHTML+= createMenuItem(index, element.itemName, element.price, element.calories);
           index++;
       });
   })
}


