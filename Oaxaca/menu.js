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
    let comp = "<div class='MenuItemDiv' id='item"+index+"'> <img class='MenuItemImg' src='image/foodimg.jpg'><br> <div class='MenuItemDetails'><label class='MenuItemName'>"+itemName+"</label><br><label class='MenuItemPrice'>£"+price.toFixed(2)+"</label><label class='MenuItemCalories'>"+calories+"kcal</label></div></div>";
    return comp;
}

// Toggle edit mode
function editMenu(){

    if(!editMode){
        document.getElementById("MenuEditSection").innerHTML += "<div id='newItemDiv'><h3>Add new menu item:</h3><p><label>Name:</label><input type='text' id='newItemNameField'>  <label>Price:</label><input type='text' id='newItemPriceField'>  <label>Calories:</label><input type='text' id='newItemCaloriesField'>    <button onclick='addMenuItem()'>+</button></p></div>";
        document.getElementById("MenuEditSection").innerHTML += "<div id='removeItemDiv'><h3>Delete menu item:</h3><p><label>Enter the name of the item you want to delete from the menu:</label><input type='text' id='deleteItemNameField'> <button onclick='deleteMenuItem()'>-</button></p></div>";
        
        // Add edit button to each menu item
        for(let i = 0; i < currentMenu.length; i++){
            document.getElementById("item"+i).innerHTML += "<button index='"+i+"' id='itemEditButton' class='editMenuItem'>EDIT</button>"
        }

        // Retrieve all edit buttons into iterable
        const editButtons = document.querySelectorAll(".editMenuItem")

        // For each edit button add click listener
        editButtons.forEach(item => {
            item.addEventListener('click', function(){
                EditMenuField(item.getAttribute("index"))
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

// Should replace name with label displaying name: and a textfield
// Should replace price and calories of price label and textfield too
function EditMenuField(index){
    // Index 0 = item name
    // Index 2 = item price
    // Index 3 = item calories
    let itemToEdit = document.getElementById("item"+index).childNodes[4].childNodes;
    

    // Replace name with name tag and checkbox
    itemToEdit[0].innerHTML = "<label class='editMenuItemPrompt'>New name:</label><input id='nameEditPrompt' class='editMenuItemPrompt' type='text'>";
    itemToEdit[2].innerHTML = "<label class='editMenuItemPrompt'>New price:</label><input id='priceEditPrompt' class='editMenuItemPrompt' type='text'>";
    itemToEdit[3].innerHTML = "<label class='editMenuItemPrompt'>New calories:</label><input id='caloriesEditPrompt' class='editMenuItemPrompt' type='text'>";

    document.getElementById("itemEditButton").remove();
    document.getElementById("item"+index).innerHTML += "<button index='"+index+"' id='itemEditButton32' class='editMenuItem'>Submit</button>"
    document.getElementById("itemEditButton32").addEventListener('click', function (){
        submitMenuEdit(index);
    })

    // MUST REPLACE EDIT BUTTON TO BE CALLED SUBMIT AND MAKE IT CALL submitMenuEdit() FUNCTION
    // SHOULD THEN REPLACE ALL THE BOXES WITH A SIMPLE "SUCCESSFUL" LABEL TO INDICATE CHANGE
}

// Needs to take in index so that it can get id
async function submitMenuEdit(index){
    let id = currentMenu[index].itemId;
    let name = document.getElementById("nameEditPrompt").value;
    let itemPrice = parseFloat(document.getElementById("priceEditPrompt").value);
    let itemCalories = parseInt(document.getElementById("caloriesEditPrompt").value);
    console.log(id+","+name+","+itemPrice+","+itemCalories);
    console.log(document.getElementById("nameEditPrompt"));

    try{
        let response = await fetch("http://localhost:4444/edit_item", {
            method:'PUT',
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