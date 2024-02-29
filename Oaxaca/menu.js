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
function addMenuItem(){
    let nameValue = document.getElementById("newItemNameField").value;
    let priceValue = parseFloat(document.getElementById("newItemPriceField").value);
    let caloriesValue = parseInt(document.getElementById("newItemCaloriesField").value);
    addItemToDB(nameValue, priceValue, caloriesValue);
}

function deleteMenuItem(){
    let nameValue = document.getElementById("deleteItemNameField").value;
    removeItem(nameValue);
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
  let previousCookieContent = document.cookie.split("basket=")[1];
  if(previousCookieContent == null){
    document.cookie = "basket="
    previousCookieContent = document.cookie.split("basket=")[1];
  }

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

    let previousBasket = document.getElementById("basketIcon").innerHTML;
    let previousBasketQuantity = previousBasket.substring(previousBasket.length-1,previousBasket.length);
    let newQuantity = Number(previousBasketQuantity)+1;
    document.getElementById("basketIcon").innerHTML = "🛒 "+newQuantity;
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

//function ro remove menu item from the order
function removeFromOrder(itemName) {
 let order = JSON.parse(localStorage.getItem('order')) || [];
 const itemIndex = order.findIndex(item => item.itemName === itemName);
 if (itemIndex >= 0) {
   order.splice(itemIndex, 1);
   localStorage.setItem('order', JSON.stringify(order));
 }
}
//updates order details and basket icon when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
  if(document.title.indexOf("Menu") != -1){
    initializeMenu();
    initializeCategoryFilter();
    updateBasketIcon();
  }
});

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

function initializeMenu() {
  const activeButton = document.querySelector('#menuFilter button.active');
  if (activeButton) {
    const category = activeButton.getAttribute('data-category');
    const filter = {
      searchTerm: '',
      maxCalories: 0,
      maxPrice: 0,
      category: category
    };
    filterMenu(filter, activeButton);
  } else {
    filterItems();
  }
}

// Initialises basket quantity on page load to number of items in basket cookies
function initBasketQuantity(){
  let cookieList = document.cookie.split(";");
  let basketCookie = "";
  cookieList.forEach(cookie => {
    if(cookie.indexOf("basket=")!=-1){
        
      // we found basket cookie
      let splitCookie = cookie.split("#")
      let basketCount = splitCookie.length
      if(splitCookie[splitCookie.length-1].length <= 0){
        basketCount--;
      }
      console.log(basketCount);
      document.getElementById("basketIcon").innerHTML="🛒 "+basketCount;
    }
  })
}