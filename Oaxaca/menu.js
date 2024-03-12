// Will hold currently displayed version of menu
var currentMenu;
var editMode = false;

// Called when menu page is initially loaded
function initMenuAll() {
    if (editMode) {
        editMenu();
    }

    refreshMenu();
}

async function refreshMenu() {
    let data = await requestMenu(0, 0, 0); // Zero value = none specified
    currentMenu = data;
    let index = 0;
    document.getElementById("MenuItemGridLayout").innerHTML = "";
    data.forEach(element => {
        document.getElementById("MenuItemGridLayout").innerHTML += createMenuItem(index, element.itemId ,element.itemName, element.imageURL, element.price, element.calories);
        index++;
    });
}

function createMenuItem(index, id, itemName, imageURL, price, calories) {
    let comp = `<div class='MenuItemDiv' itemid="${id}" id='item` + index + `'> <img class='MenuItemImg' src='http://localhost:4444/image/${imageURL}'><div class='MenuItemDetails'><label class='MenuItemName'>` + itemName + `</label><br><label class='MenuItemPrice'>£` + price.toFixed(2) + `</label><br><label class='MenuItemCalories'>` + calories + ` kcal</label></div>`;
    comp += "<button class='addBasketButton' onclick='addToBasket(" + index + "," + id + ", \"" + itemName + "\", " + price + ", " + calories + ")'>Add to Basket</button>";
    return comp;
}

async function refreshEditMenu() {
    await refreshMenu();
    document.getElementById("MenuEditSection").innerHTML = "";
    document.getElementById("MenuEditSection").innerHTML += `
    <div id='newItemDiv'>
        <h3>Add new menu item:</h3>
        <p>
            <label>Name:</label>
            <input type='text' id='newItemNameField'><br><br>
            <label>Upload Image</label><br><br>
            <input type='file' id='newItemFileUpload'><br><br>
            <label>Price:</label>
            <input type='text' id='newItemPriceField'><br><br>
            <label>Calories:</label><input type='text' id='newItemCaloriesField'><br><br>
            <button onclick='addMenuItem()'>Add Item</button>
        </p>
    </div>`;

    // Add edit button to each menu item
    for (let i = 0; i < currentMenu.length; i++) {
        document.getElementById("item" + i).innerHTML += `
        <button index="${i}" class="editMenuItemButton">Edit</button>
        <button index="${i}" class="deleteMenuItemButton">Delete</button>`;
    }

    const editButtons = document.querySelectorAll(".editMenuItemButton");
    const deleteButtons = document.querySelectorAll(".deleteMenuItemButton");

    editButtons.forEach(item => {
        item.addEventListener('click', function () {
            editMenuForItem(item.getAttribute("index"));
        });
    });

    deleteButtons.forEach(item => {
        item.addEventListener("click", () => {
            deleteMenuItem(item.getAttribute("index"));
        });
    })
}

// Toggle edit mode
function editMenu() {
    if (!editMode) {
        refreshEditMenu();
        editMode = true;
    } else {
        editMode = false;
        document.getElementById("newItemDiv").remove();
        const editButtons = document.querySelectorAll('.editMenuItemButton');
        const deleteButtons = document.querySelectorAll('.deleteMenuItemButton');
        editButtons.forEach(item => {
            item.remove();
        });
        deleteButtons.forEach(item => {
            item.remove();
        })
    }
}


var editingText = false;
// Holds location of item being edited
var currentlyEditing = -1;
var oldText;

// Should replace name with label displaying name: and a textfield
// Should replace price and calories of price label and textfield too
function editMenuForItem(index) {
    // Index 0 = item name
    // Index 2 = item price
    // Index 3 = item calories
    let itemToEdit = document.getElementById("item" + index).childNodes[4].childNodes;

    // Replace name with name tag and checkbox
    itemToEdit[0].innerHTML = "<label class='editMenuItemPrompt'>New name:</label><input id='nameEditPrompt' class='editMenuItemPrompt' type='text'>";
    itemToEdit[2].innerHTML = "<label class='editMenuItemPrompt'>New price:</label><input id='priceEditPrompt' class='editMenuItemPrompt' type='text'>";
    itemToEdit[3].innerHTML = "<label class='editMenuItemPrompt'>New calories:</label><input id='caloriesEditPrompt' class='editMenuItemPrompt' type='text'>";

    document.getElementById("item" + index).innerHTML += "<button index='" + index + "' class='editMenuItemButton'>Submit</button>";
    document.querySelector(".editMenuItemButton").addEventListener('click', function () {
        submitMenuEdit(index);
        let toRemove = document.getElementsByClassName("editMenuItemPrompt");
        toRemove[0].innerHTML = "<label>Submitted!</label>";
        for (let i = 1; i < toRemove.length; i++) {
            toRemove[i].innerHTML = "";
        }
        document.getElementById("nameEditPrompt").remove();
        document.getElementById("priceEditPrompt").remove();
        document.getElementById("caloriesEditPrompt").remove();
        document.querySelector(".editMenuItemButton").remove();
    });
}

// Needs to take in index so that it can get id
async function submitMenuEdit(index) {
    let id = currentMenu[index].itemId;
    let name = document.getElementById("nameEditPrompt").value;
    let itemPrice = parseFloat(document.getElementById("priceEditPrompt").value);
    let itemCalories = parseInt(document.getElementById("caloriesEditPrompt").value);

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
        });

    } catch (error) {
        console.error(error);
    }
}

// Functions to add and delete menu items
async function addMenuItem() {
    let nameValue = document.getElementById("newItemNameField").value;
    let image = document.getElementById("newItemFileUpload").files[0];
    let priceValue = parseFloat(document.getElementById("newItemPriceField").value);
    let caloriesValue = parseInt(document.getElementById("newItemCaloriesField").value);

    // Check that the user included an image
    if(image == null) {
        alert("Please upload an image for the item");
        return;
    }

    let imageURL = await uploadImage(image);
    let result = await addItemToDB(nameValue, imageURL, priceValue, caloriesValue);

    if (result >= 0) {
        await refreshEditMenu();
    }
}

async function uploadImage(image) {
    let data = new FormData();
    var resultFilename = "";
    data.append("file", image);

    await fetch("http://localhost:4444/upload", {
        method: "POST",
        body: data,
    })
    .then((res) => res.text())
    .then((filename) => {
        resultFilename = filename;
    }).catch(() => alert("Failed to upload image"));

    return resultFilename;
}

function deleteMenuItem(index) {
    let id = document.getElementById(`item${index}`).getAttribute("itemid");
    removeItem(id).then(r => {
        if (r >= 0) {
            document.getElementById(`item${index}`).remove();
        }
    });

}

// Add menu item
async function addItemToDB(name, imageURL, _price, _calories) {
  try {
    let response = await fetch("http://localhost:4444/add_item", {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        itemName: name,
        imageURL: imageURL,
        price: _price,
        calories: _calories,
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
async function removeItem(id) {
    try {
        let response = await fetch("http://localhost:4444/remove_item?" + new URLSearchParams({
            itemId: id
        }).toString(), {
            method: "DELETE"
        });
        return 0;
    } catch (error) {
        console.error(error);
        return -1;
    }
}

// Fetches data from backend or throws error to console and provides example menu data
async function requestMenu(userSearchTerm, userMaxPrice, userMaxCalories) {
  try {
    let response = await fetch("http://localhost:4444/menu?" + new URLSearchParams({
      searchTerm: userSearchTerm,
      Price: userMaxPrice,
      Calories: userMaxCalories,
    }).toString())

        if (!response.ok) {
            console.log("ERROR fetching menu");
        }

        let data = await response.json();
        return data;
    } catch (error) {
        console.error(error);
        // Return example error if unable to connect to backend
        return ([{
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
        }]);
    }
}

// function to add menu item to basket
// Stores menu items in basket using cookies
// Cookie structure is CSV in form: id,itemName,price,calories,quantity
function addToBasket(index, itemId, itemName, price, calories) {
  let quantity = 1;

  // This will work for now as we only store 1 type of cookie
  let previousCookieContent = document.cookie.split("basket=")[1];
  if (previousCookieContent == null) {
    document.cookie = "basket="
    previousCookieContent = document.cookie.split("basket=")[1];
  }

  let updated = false;

  // Check that item is not already in basket
  previousCookieContent.split("#").forEach(element => {
    let splitCookie = element.split(",")
    // Item found in basket
    if (splitCookie[0] == itemId) {
      // Then update quantity instead
      let newQuantity = Number(splitCookie[4]) + Number(quantity);

      let indexOfItem = previousCookieContent.indexOf(element)
      let indexOfEndOfItem = previousCookieContent.indexOf("#", indexOfItem)
      let updatedCookieSegment = previousCookieContent.substring(indexOfItem, indexOfEndOfItem - 1) + newQuantity
      let updatedCookie = previousCookieContent.substring(0, indexOfItem) + updatedCookieSegment + previousCookieContent.substring(indexOfEndOfItem, previousCookieContent.length);
      document.cookie = "basket=" + updatedCookie;

      updated = true;
    }
  });

  if (!updated) {
    document.cookie = "basket=" + itemId + "," + itemName + "," + price + "," + calories + "," + quantity + "#" + previousCookieContent;
  }

  console.log("Current basket:" + document.cookie);

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
document.addEventListener('DOMContentLoaded', function () {
    updateOrderDetails();
});

function filterItems() {
    let searchTerm = document.getElementById('searchTerm').value.toLowerCase();
    let maxCalories = parseInt(document.getElementById('maxCalories').value) || 0;
    let maxPrice = parseFloat(document.getElementById('maxPrice').value) || 0;

    let data = requestMenu('', 0, 0); // Fetch all menu items
    data.then(r => {
        let index = 0;
        document.getElementById("MenuItemGridLayout").innerHTML = "";
        r.forEach(element => {
            // Check if the current item matches the filter criteria
            if (
                (searchTerm.length === 0 || element.itemName.toLowerCase().includes(searchTerm)) &&
                (maxCalories === 0 || element.calories <= maxCalories) &&
                (maxPrice === 0 || element.price <= maxPrice)
            ) {
                document.getElementById("MenuItemGridLayout").innerHTML += createMenuItem(index, element.itemName, element.imageURL, element.price, element.calories);
                index++;
            }
        });
    });
}

// Initialises basket quantity on page load to number of items in basket cookies
function initBasketQuantity(){
    if(document.cookie.length <= 0){
      document.cookie="basket="
    }
    let cookieList = document.cookie.split(";");
    let basketCookie = "";
    cookieList.forEach(cookie => {
      if(cookie.indexOf("basket=")!=-1){
        cookie = cookie.substring(cookie.indexOf("basket=")+"basket=".length, cookie.length)
        console.log(cookie)
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
  
