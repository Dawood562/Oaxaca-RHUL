// Will hold currently displayed version of menu
var currentMenu;
var editMode = false;

// Called when menu page is initially loaded
function initMenuAll() {
    if (editMode) {
        editMenu();
    }

    let data = requestMenu(0, 0, 0); // Zero value = none specified
    data.then(r => {
        currentMenu = r;
        let index = 0;
        document.getElementById("MenuItemGridLayout").innerHTML = "";
        r.forEach(element => {
            document.getElementById("MenuItemGridLayout").innerHTML += createMenuItem(index, element.itemName, element.price, element.calories);
            index++;
        });
    });
}

function createMenuItem(index, itemName, price, calories) {
    let comp = "<div class='MenuItemDiv' id='item" + index + "'> <img class='MenuItemImg' src='image/foodimg.jpg'><br> <div class='MenuItemDetails'><label class='MenuItemName'>" + itemName + "</label><br><label class='MenuItemPrice'>£" + price.toFixed(2) + "</label><label class='MenuItemCalories'>" + calories + "kcal</label></div>";
    comp += "<button class='addBasketButton' onclick='addToBasket(" + index + ", \"" + itemName + "\", " + price + ", " + calories + ")'>Add to Basket</button>";
    comp += "<button class='editMenuItemButton' onclick='editMenuForItem(" + index + ")'>Edit</button></a></div>"; // Add Edit button
    return comp;
}

// Toggle edit mode
function editMenu() {
    console.log("Edit menu button clicked");
    if (!editMode) {
        document.getElementById("MenuEditSection").innerHTML += "<div id='newItemDiv'><h3>Add new menu item:</h3><p><label>Name:</label><input type='text' id='newItemNameField'>  <label>Price:</label><input type='text' id='newItemPriceField'>  <label>Calories:</label><input type='text' id='newItemCaloriesField'>    <button onclick='addMenuItem()'>+</button></p></div>";
        document.getElementById("MenuEditSection").innerHTML += "<div id='removeItemDiv'><h3>Delete menu item:</h3><p><label>Enter the name of the item you want to delete from the menu:</label><input type='text' id='deleteItemNameField'> <button onclick='deleteMenuItem()'>-</button></p></div>";

        // Add edit button to each menu item
        for (let i = 0; i < currentMenu.length; i++) {
            document.getElementById("item" + i).innerHTML += "<button index='" + i + "' class='editMenuItemButton'>Edit</button>";
        }

        const editButtons = document.querySelectorAll(".editMenuItemButton");

        editButtons.forEach(item => {
            item.addEventListener('click', function () {
                editMenuForItem(item.getAttribute("index"));
            });
        });

        editMode = true;
    } else {
        document.getElementById("newItemDiv").remove();
        document.getElementById("removeItemDiv").remove();
        const editButtons = document.querySelectorAll('.editMenuItemButton');
        editButtons.forEach(item => {
            item.remove();
        });
        editMode = false;
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
    let priceValue = parseFloat(document.getElementById("newItemPriceField").value);
    let caloriesValue = parseInt(document.getElementById("newItemCaloriesField").value);
    let result = await addItemToDB(nameValue, priceValue, caloriesValue);

    if (result >= 0) {
        document.getElementById("MenuItemGridLayout").innerHTML += createMenuItem(currentMenu.length, nameValue, priceValue, caloriesValue);
    }
}

function deleteMenuItem() {
    let nameValue = document.getElementById("deleteItemNameField").value;
    removeItem(nameValue).then(r => {
        if (r >= 0) {
            // First childNodes holds each item on menu
            // Second childNodes holds the details of the first childnodes menu item
            // Third childNodes holds the actual details. index 0 is name
            for (let i = 0; i < document.getElementById("MenuItemGridLayout").childNodes.length; i++) {
                if (document.getElementById("MenuItemGridLayout").childNodes[i].childNodes[4].childNodes[0].innerHTML == nameValue) {
                    let gridLayout = document.getElementById("MenuItemGridLayout");
                    let itemToRemove = document.getElementById("MenuItemGridLayout").childNodes[0];
                    gridLayout.removeChild(itemToRemove);
                }
            }
        }
    });

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
        });
        return 0;
    } catch (error) {
        console.error(error);
        return -1;
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
    nameId = String(getIdFromName(name));
    try {
        let response = await fetch("http://localhost:4444/remove_item?" + new URLSearchParams({
            itemId: nameId
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
            maxPrice: userMaxPrice,
            maxCalories: userMaxCalories
        }).toString());

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

function addToBasket(index, itemName, price, calories) {
    let order = JSON.parse(localStorage.getItem('order')) || [];

    const item = {
        index: index,
        itemName: itemName,
        price: price,
        calories: calories
    };
    order.push(item);
    localStorage.setItem('order', JSON.stringify(order));
    const orderDetailsDiv = document.getElementById('orderDetails');
    const li = document.createElement('li');
    li.innerHTML = `
        <h3>${item.itemName}</h3>
        <p>Price: £${item.price.toFixed(2)}</p>
        <p>Calories: ${item.calories} kcal</p>
    `;
    orderDetailsDiv.appendChild(li);
}

function updateOrderDetails() {
    const order = JSON.parse(localStorage.getItem('order'));
    const orderDetailsDiv = document.getElementById('orderDetails');

    if (order && order.length > 0) {
        order.forEach(item => {
            const li = document.createElement('li');
            li.innerHTML = `
                <h3>${item.itemName}</h3>
                <p>Price: £${item.price.toFixed(2)}</p>
                <p>Calories: ${item.calories} kcal</p>
            `;
            orderDetailsDiv.appendChild(li);
        });
    } else {
        orderDetailsDiv.innerHTML = `<p>No items selected.</p>`;
    }
}

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
                document.getElementById("MenuItemGridLayout").innerHTML += createMenuItem(index, element.itemName, element.price, element.calories);
                index++;
            }
        });
    });
}
