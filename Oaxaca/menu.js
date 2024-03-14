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

// Refreshes the menu with no filter
async function refreshMenu() {
    await fetchMenuWithFilter("", 0, 0, []);
}

async function fetchMenuWithFilter(searchTerm, maxPrice, maxCalories, excludedAllergens) {
    let allergenString = "";
    excludedAllergens.forEach(item => {
        allergenString += item + ",";
    })
    await fetch("http://localhost:4444/menu?" + new URLSearchParams({
        searchTerm: searchTerm,
        maxPrice: maxPrice,
        maxCalories: maxCalories,
        allergens: allergenString.substring(0, allergenString.length - 1)
    }), {
        method: "GET",
    })
    .then((res) => {
        if(res.ok) {
            return res.json();
        }
        throw new Error("Server returned an error");
    })
    .then((data) => {
        currentMenu = data;
        renderMenu();
    }).catch((error) => {
        alert("Failed to apply filter: " + error);
    });
}

function renderMenu() {
    document.getElementById("MenuItemGridLayout").innerHTML = "";
    currentMenu.forEach(element => {
        document.getElementById("MenuItemGridLayout").innerHTML += createMenuItem(element.itemId ,element.itemName, element.imageURL, element.price, element.calories, element.allergens);
    });
}

function createMenuItem(id, itemName, imageURL, price, calories, allergens) {
    return `
    <div class='MenuItemDiv' id='item${id}'>
        <img class='MenuItemImg' src='http://localhost:4444/image/${imageURL}'>
        <div class='MenuItemDetails'>
            <label class='MenuItemName' id="itemName${id}">${itemName}</label><br>
            <input style="display: none; max-width:90%" placeholder="Item Name" id='nameEditPrompt${id}' class='editMenuItemPrompt' type='text'>
            <br><label style="display: none" for='imageEditPrompt${id}'>Image</label>
            <input style="display: none" type='file' id='imageEditPrompt${id}' class='editMenuItemPrompt'>
            <label class='MenuItemPrice' id="itemPrice${id}">£${price.toFixed(2)}</label><br>
            <p id="priceContext${id}" class="editMenuContext">£</p><input style="display: none" id='priceEditPrompt${id}' placeholder="Price" class='editMenuItemPrompt' type='text'>
            <label class='MenuItemCalories' id="itemCalories${id}">${calories} kcal</label>
            <label class='MenuItemAllergens' id="itemAllergens${id}"><br><b>Allergens:</b><br>${renderAllergens(allergens)}</label>
            <input style="display: none" id='caloriesEditPrompt${id}' class='editMenuItemPrompt' type='text' placeholder="Calories"><p id="caloriesContext${id}" class="editMenuContext">kcal</p>
            <input style="display: none; max-width:90%" id='allergensEditPrompt${id}' class='editMenuItemPrompt' type='text' placeholder="Allergens">
        </div>
        <button id='addItem${id}' + class='addBasketButton' onclick='addToBasket(${id}, "${itemName}", ${price}, ${calories}, "${imageURL}")'>Add to Basket</button>
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
    let formatter = new Intl.ListFormat("en", {
        style: "long",
        type: "conjunction"
    });
    return formatter.format(allergens.map((item) => item.name));
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
    document.getElementById(`imageEditPrompt${id}`).style.display = "none";
    document.querySelector(`label[for='imageEditPrompt${id}']`).style.display = "none";
    document.getElementById(`itemPrice${id}`).style.display = "inline";
    document.getElementById(`priceEditPrompt${id}`).style.display = "none";
    document.getElementById(`itemAllergens${id}`).style.display = "block";
    document.getElementById(`allergensEditPrompt${id}`).style.display = "none";
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

        document.querySelector(`label[for='imageEditPrompt${id}']`).style.display = "inline";
        document.getElementById(`imageEditPrompt${id}`).style.display = "inline";

        document.getElementById(`itemAllergens${id}`).style.display = "none";
        document.getElementById(`allergensEditPrompt${id}`).style.display = "inline";
        let allergens = "";
        item.allergens.forEach(allergen => {
            allergens += allergen.name + ", ";
        });
        document.getElementById(`allergensEditPrompt${id}`).value = allergens.substring(0, allergens.length - 2);

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
    let allergens = allergenStringToArray(document.getElementById(`allergensEditPrompt${id}`).value)
    let newImage = document.getElementById(`imageEditPrompt${id}`).files[0];

    let imageURL;
    if (newImage) {
        imageURL = await uploadImage(newImage);
    } else {
        imageURL = getMenuItemById(id).imageURL;
    }

    try {
        let response = await fetch("http://localhost:4444/edit_item", {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                itemId: id, 
                itemName: name,
                imageURL: imageURL,
                price: itemPrice,
                calories: itemCalories,
                allergens: allergens
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
    let allergens = document.getElementById("newItemAllergensField").value;
    // Check that the user included an image
    if(image == null) {
        alert("Please upload an image for the item");
        return;
    }

    let imageURL = await uploadImage(image);
    let result = await addItemToDB(nameValue, imageURL, priceValue, caloriesValue, allergenStringToArray(allergens));

    if (result >= 0) {
        await refreshEditMenu();
    }
}

function getAllergenList() {
    let labels = [];
    let allergens = [];
    document.querySelectorAll('.excludedAllergens > label').forEach(item => {
        labels.push(item);
    });
    document.querySelectorAll('.excludedAllergens > input[type="checkbox"]').forEach(item => {
        if(item.checked) {
            // Find label matching this item
            let id = item.id;
            labels.forEach(label => {
                if(label.htmlFor === id) {
                    allergens.push(label.innerHTML);
                }
            });
        }
    });
    return allergens;
}

function allergenStringToArray(input) {
    return input.split(",").map(item => {
        return {
            name: item.trim()
        }
    }).filter((item) => {
        return item.name.length > 0
    });
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
async function addItemToDB(name, imageURL, _price, _calories, _allergens) {
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
        allergens: _allergens
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

// function to add menu item to basket
// Stores menu items in basket using cookies
// Cookie structure is CSV in form: id,itemName,price,calories,quantity
function addToBasket(itemId, itemName, price, calories, imageURL) {
  let quantity = 1;

  // This will work for now as we only store 1 type of cookie
  let previousCookieContent = document.cookie.split("basket=")[1];
  if (previousCookieContent == null) {
    document.cookie = "basket="
    previousCookieContent = document.cookie.split("basket=")[1];
  }

  let updated = false;
  let i = -1;
  let updatedIndex = -1;

  // Check that item is not already in basket
  previousCookieContent.split("#").forEach(element => {
    i++;
    let splitCookie = element.split(",")
    // Item found in basket
    if (splitCookie[0] == itemId) {
        updated = true;
        updatedIndex = i;
    }
  });

  let splitCookieContent = previousCookieContent.split("#");
  
  if(updated){  
    // Update quantity:
    let toModify = splitCookieContent[updatedIndex].split(",");
    toModify[4]++;
    
    let newCookie = "basket=";
    let count = 0;
    splitCookieContent.forEach(item => {
        if(item.length > 0){
            if(count == updatedIndex){
                // Add toModify
                newCookie+=toModify+"#";
            }else{
                // Add back original
                newCookie+=item+"#";
            }
        }
        count++;
    })
    document.cookie = newCookie;
  }
  else {
    document.cookie = "basket=" + itemId + "," + itemName + "," + price + "," + calories + "," + quantity + "," + imageURL + "#" + previousCookieContent;
  }
  updateBasketQuantity();
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
    let searchTerm = document.getElementById('searchTerm').value.toLowerCase();
    let maxCalories = parseInt(document.getElementById('maxCalories').value) || 0;
    let maxPrice = parseFloat(document.getElementById('maxPrice').value) || 0;
    let allergens = getAllergenList();
    fetchMenuWithFilter(searchTerm, maxPrice, maxCalories, allergens);
}

function getCookieQuantity(){
    // if basket cookie doesnt exist, create it
    if(document.cookie.indexOf("basket=") < 0){
        document.cookie = "basket="
    }

    let cookies = document.cookie.split(";");
    let basketCookie = "";
    // Find cookie with basket data
    for(let i = 0; i < cookies.length; i++){
        if(cookies[i].indexOf("basket=") >= 0){
            basketCookie = cookies[i].split("basket=")[1];
        }
    }
    if(basketCookie.length <= 0){
        return 0
    }

    // split basket
    let basketSplit = basketCookie.split("#")
    let basketQuantity = basketSplit.length;

    // if basket seperator causes extra empty basket then subtract 1
    if(basketSplit[basketSplit.length-1].length <= 0){
        basketQuantity--;
    }
    return basketQuantity;
}

// Updates basket quantity
function updateBasketQuantity(){
    let basketIcon = document.getElementById("basket-quantity");
    let currentBasketQuantity = getCookieQuantity();
    basketIcon.innerHTML = "🛒 "+currentBasketQuantity;
}

function goToOrderPage() {
    console.log("!");
    window.location.href = "order.html";
}