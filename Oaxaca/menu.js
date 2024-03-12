// Will hold currently displayed version of menu
var currentMenu;
var editMode = false;
var currentEdit = -1;
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
    document.getElementById("MenuItemGridLayout").innerHTML = "";
    data.forEach(element => {
        document.getElementById("MenuItemGridLayout").innerHTML += createMenuItem(element.itemId ,element.itemName, element.imageURL, element.price, element.calories, element.allergens);
    });
}

function createMenuItem(id, itemName, imageURL, price, calories, allergens) {
    return `
    <div class='MenuItemDiv' id='item${id}'>
        <img class='MenuItemImg' src='http://localhost:4444/image/${imageURL}'>
        <div class='MenuItemDetails'>
            <label class='MenuItemName' id="itemName${id}">${itemName}</label><br>
            <input style="display: none" id='nameEditPrompt${id}' class='editMenuItemPrompt' type='text'>
            <label class='MenuItemPrice' id="itemPrice${id}">£${price.toFixed(2)}</label><br>
            <p id="priceContext${id}" class="editMenuContext">£</p><input style="display: none" id='priceEditPrompt${id}' class='editMenuItemPrompt' type='text'>
            <label class='MenuItemCalories' id="itemCalories${id}">${calories}kcal</label>
            <label class='MenuItemAllergens' id="itemAllergens${id}"><br><b>Allergens:</b><br>${renderAllergens(allergens)}</label>
            <input style="display: none" id='caloriesEditPrompt${id}' class='editMenuItemPrompt' type='text'><p id="caloriesContext${id}" class="editMenuContext">kcal</p>
        </div>
        <button id='addItem${id}' + class='addBasketButton' onclick='addToBasket(${id}, "${itemName}", ${price}, ${calories})'>Add to Basket</button>
        <button index="${id}" id="editItem${id}" style="display: none" class="editMenuItemButton">Edit</button>
        <button index="${id}" id="deleteItem${id}" style="display: none" class="deleteMenuItemButton">Delete</button>
        <button index="${id}" id="cancelEditItem${id}" style="display: none" class="cancelEditMenuItemButton" onclick="closeEdit(${id})">Cancel</button>
        <button index="${id}" id="submitEditItem${id}" style="display: none" class="submitEditMenuItemButton" onclick="submitEdit(${id})">Submit</button>
    </div>`;
}

function renderAllergens(allergens) {
    if(allergens.length === 0) {
        return "None";
    }
    let formatter = Intl.ListFormat("en", {
        style: "long",
        type: "conjunction"
    });
    return formatter.format(allergens);
}

async function refreshEditMenu() {
    await refreshMenu();
    document.getElementById("menuEditSection").style.display = "block";

    // Add edit button to each menu item
    document.querySelectorAll(".addBasketButton").forEach(item => {
        item.style.display = "none";
    })

    document.querySelectorAll(".editMenuItemButton").forEach(item => {
        item.style.display = "inline";
        item.addEventListener('click', function () {
            editMenuForItem(item.getAttribute("index"));
        });
    });

    document.querySelectorAll(".deleteMenuItemButton").forEach(item => {
        item.style.display = "inline";
        item.addEventListener("click", () => {
            deleteMenuItem(item.getAttribute("index"));
        });
    })
}

function closeEdit(id) {
    document.getElementById(`itemName${id}`).style.display = "inline";
    document.getElementById(`nameEditPrompt${id}`).style.display = "none";
    document.getElementById(`itemPrice${id}`).style.display = "inline";
    document.getElementById(`priceEditPrompt${id}`).style.display = "none";
    document.getElementById(`itemAllergens${id}`).style.display = "block";
    document.getElementById(`priceContext${id}`).style.display = "none";
    document.getElementById(`itemCalories${id}`).style.display = "inline";
    document.getElementById(`caloriesContext${id}`).style.display = "none";
    document.getElementById(`caloriesEditPrompt${id}`).style.display = "none";
    document.getElementById(`submitEditItem${id}`).style.display = "none";
    document.getElementById(`cancelEditItem${id}`).style.display = "none";
    document.getElementById(`editItem${id}`).style.display = "inline";
    document.getElementById(`deleteItem${id}`).style.display = "inline";
    currentEdit = -1;
}

// Toggle edit mode
function editMenu() {
    if (!editMode) {
        refreshEditMenu();
        editMode = true;
    } else {
        editMode = false;
        if(currentEdit !== -1) {
            closeEdit(currentEdit);
        }
        document.getElementById("menuEditSection").style.display = "none";
        document.querySelectorAll(".addBasketButton").forEach(item => {
            item.style.display = "block";
        })
        document.querySelectorAll('.editMenuItemButton').forEach(item => {
            item.style.display = "none";
        })
        document.querySelectorAll(".deleteMenuItemButton").forEach(item => {
            item.style.display = "none";
        })
        document.querySelectorAll(".submitEditMenuItemButton").forEach(item => {
            item.style.display = "none";
        })
        document.querySelectorAll(".cancelEditMenuItemButton").forEach(item => {
            item.style.display = "none";
        })
    }
}

function editMenuForItem(id) {
    let item = getMenuItemById(id);
    if(currentEdit === -1) {
        // Replace name with name tag and checkbox
        document.getElementById(`itemName${id}`).style.display = "none";
        document.getElementById(`nameEditPrompt${id}`).style.display = "inline";
        document.getElementById(`nameEditPrompt${id}`).value = item.itemName;

        document.getElementById(`itemAllergens${id}`).style.display = "none";

        document.getElementById(`itemPrice${id}`).style.display = "none";
        document.getElementById(`priceEditPrompt${id}`).style.display = "inline";
        document.getElementById(`priceEditPrompt${id}`).value = item.price;
        document.getElementById(`priceContext${id}`).style.display = "inline";

        document.getElementById(`itemCalories${id}`).style.display = "none";
        document.getElementById(`caloriesEditPrompt${id}`).style.display = "inline";
        document.getElementById(`caloriesEditPrompt${id}`).value = item.calories;
        document.getElementById(`caloriesContext${id}`).style.display = "inline";

        document.getElementById(`submitEditItem${id}`).style.display = "inline";
        document.getElementById(`cancelEditItem${id}`).style.display = "inline";

        document.getElementById(`editItem${id}`).style.display = "none";
        document.getElementById(`deleteItem${id}`).style.display = "none";
        currentEdit = id;
    }
}

function submitEdit(id) {
    submitMenuEdit(id);
    closeEdit(id);
}

function getMenuItemById(id) {
    return currentMenu.find(item => item.itemId == id);
}

// Needs to take in index so that it can get id
async function submitMenuEdit(id) {
    let name = document.getElementById(`nameEditPrompt${id}`).value;
    let itemPrice = parseFloat(document.getElementById(`priceEditPrompt${id}`).value);
    let itemCalories = parseInt(document.getElementById(`caloriesEditPrompt${id}`).value);

    try {
        let response = await fetch("http://localhost:4444/edit_item", {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                itemId: id,
                itemName: name,
                imageURL: getMenuItemById(id).imageURL,
                price: itemPrice,
                calories: itemCalories
            })
        });

    } catch (error) {
        alert("Failed to update item: " + error);
    }

    refreshEditMenu();
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

function deleteMenuItem(id) {
    removeItem(id).then(r => {
        if (r >= 0) {
            document.getElementById(`item${id}`).remove();
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

function applyFilter() {
    // TODO: implement this
    console.log("Applying filter");
}

function filterItems() {
    let searchTerm = document.getElementById('searchTerm').value.toLowerCase();
    let maxCalories = parseInt(document.getElementById('maxCalories').value) || 0;
    let maxPrice = parseFloat(document.getElementById('maxPrice').value) || 0;

    let data = requestMenu('', 0, 0); // Fetch all menu items
    data.then(r => {
        document.getElementById("MenuItemGridLayout").innerHTML = "";
        r.forEach(element => {
            // Check if the current item matches the filter criteria
            if (
                (searchTerm.length === 0 || element.itemName.toLowerCase().includes(searchTerm)) &&
                (maxCalories === 0 || element.calories <= maxCalories) &&
                (maxPrice === 0 || element.price <= maxPrice)
            ) {
                document.getElementById("MenuItemGridLayout").innerHTML += createMenuItem(element.itemName, element.imageURL, element.price, element.calories, element.allergens);
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
  
