// Will hold currently displayed version of menu
var currentMenu;
var editMode = false;
// Called when menu page is initially loaded
function initMenuAll(){
    if(editMode){
        editMenu()
    }

    let data = requestMenu(0,0,0); // Zero value = none specified
    data.then(r => {
        console.log(r)
        currentMenu = r
        let index = 0;
        document.getElementById("MenuItemGridLayout").innerHTML = ""
        r.forEach(element => {
            document.getElementById("MenuItemGridLayout").innerHTML+= createMenuItem(index, element.itemName, element.price, element.calories);
            index++;
        });
    })
}

// Function that takes in data and turns into menu item to be displayed
function createMenuItem(index, itemName, price, calories){
    let comp = "<div class='MenuItemDiv' id='"+index+"'> <img class='MenuItemImg' src='image/foodimg.jpg'><br> <div class='MenuItemDetails'><label class='MenuItemName'>"+itemName+"</label><br><label class='MenuItemPrice'>£"+price.toFixed(2)+"</label><label class='MenuItemCalories'>"+calories+"kcal</label></div></div>";
    return comp;
}

// Toggle edit mode
function editMenu(){

    if(!editMode){
        document.getElementById("menuSectionAll").innerHTML += "<div id='newItemDiv'><h3>Add new menu item:</h3><p><label>Name:</label><input type='text' id='newItemNameField'>  <label>Price:</label><input type='text' id='newItemPriceField'>  <label>Calories:</label><input type='text' id='newItemCaloriesField'>    <button onclick='addMenuItem()'>+</button></p></div>";
        document.getElementById("menuSectionAll").innerHTML += "<div id='removeItemDiv'><h3>Delete menu item:</h3><p><label>Enter the name of the item you want to delete from the menu:</label><input type='text' id='deleteItemNameField'> <button onclick='deleteMenuItem()'>-</button></p></div>";
        
        // Add edit button to each menu item
        for(let i = 0; i < currentMenu.length; i++){
            document.getElementById("item"+i).innerHTML += "<button index='"+i+"' class='editMenuItem'>EDIT</button>"
        }

        // Retrieve all edit buttons into iterable
        const editButtons = document.querySelectorAll(".editMenuItem")

        // For each edit button add click listener
        editButtons.forEach(item => {
            item.addEventListener('click', function(){
                toggleEditMode(item.getAttribute("index"))
            })
        })

        editMode = true;
    }else{
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

// From here a textfield should appear where menu item was with the previous content that was there inside it
// Also submit button should apear instead of 
function toggleEditMode(index){
    let previousText = document.getElementById("item"+index).textContent;

    previousText = previousText.substring(0,previousText.length-4); // Gets rid of edit

    document.getElementById("item"+index).innerHTML = "<input id='menuEditBox' type='text' value='"+previousText+"'> <button id='menuEditSubmitButton'>submit</button>"
    document.getElementById("menuEditSubmitButton").addEventListener('click', function (){
        submitMenuEdit(index);
    });

    // Must now replace new item with item from server
}

// Needs to take in index so that it can get id
async function submitMenuEdit(index){
    let userInput = document.getElementById("menuEditBox").value
    userInput = userInput.split("--");
    let id = currentMenu[index].itemId;
    let name = userInput[0].substring(0,userInput[0].length-1);
    let itemPrice = parseFloat(userInput[1].substring(2, userInput[1].length-1));
    let itemCalories = parseInt(userInput[2].substring(1, userInput[2].length-4));
    console.log(id+","+name+","+itemPrice+","+itemCalories);

    try{
        let response = await fetch("http://localhost:4444/edit_item", {
            method:'PATCH',
            headers:{
                'Content-Type':'application/json'
            },
            body: JSON.stringify({
                itemId: id,
                itemName: name,
                price: itemPrice,
                calories: itemCalories
            })
        })
    }catch(error){
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
async function addItemToDB(name, _price, _calories){
    try{
        let response = await fetch("http://localhost:4444/add_item", {
            method:'POST',
            headers:{
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                itemName: name,
                price: _price,
                calories: _calories
            }) 
        })
    }catch(error){
        console.error(error)
    }
}

// Takes in name and returns id if found
function getIdFromName(name){
    for(let i = 0; i < currentMenu.length;i++){
        if(currentMenu[i].itemName == name){
            return currentMenu[i].itemId;
        }
    }
    console.error("COULD NOT FIND ID FROM ITEM NAME: "+name);
    return null;
}

// Contacts backend to remove menu item
async function removeItem(name){
    console.log(currentMenu)
    nameId = String(getIdFromName(name))
    try{
        let response = await fetch("http://localhost:4444/remove_item?" + new URLSearchParams({
            itemId: nameId
        }).toString(), {
            method: "DELETE"
        })

    }catch(error){
        console.error(error)
    }
}

// Fetches data from backend or throws error to console and provides example menu data
async function requestMenu(userSearchTerm, userMaxPrice, userMaxCalories){
    try{
        let response = await fetch("http://localhost:4444/menu?" + new URLSearchParams({
            ItemName: userSearchTerm,
            Price: userMaxPrice,
            Calories: userMaxCalories 
        }).toString())

        if(!response.ok){
            console.log("ERROR fetching menu")
        }

        let data = await response.json()
        return data
    }catch(error){
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
