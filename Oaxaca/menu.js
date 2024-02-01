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
        document.getElementById("menuSectionAll").innerHTML = ""
        r.forEach(element => {
            let item = '<li id="item'+index+'" class="genericMenuItem" data-calories="'+element.calories+'" data-price="'+element.price+'">'+element.itemName+' -- £'+element.price.toFixed(2)+' -- '+element.calories+'kcal</li>';
            document.getElementById("menuSectionAll").innerHTML+= item
            index++;
        });
    })
}

// Toggle edit mode
function editMenu(){

    if(!editMode){
        document.getElementById("menuSectionAll").innerHTML += "<div id='newItemDiv'><h3>Add new menu item:</h3><p><label>Name:</label><input type='text' id='newItemNameField'>  <label>Price:</label><input type='text' id='newItemPriceField'>  <label>Calories:</label><input type='text' id='newItemCaloriesField'>    <button onclick='addMenuItem()'>+</button></p></div>";
        document.getElementById("menuSectionAll").innerHTML += "<div id='newItemDiv'><h3>Delete menu item:</h3><p><label>Enter the name of the item you want to delete from the menu:</label><input type='text' id='deleteItemNameField'> <button onclick='deleteMenuItem()'>-</button></p></div>";
        editMode = true;
    }else{
        document.getElementById("newItemDiv").remove();
        editMode = false;
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
    deleteItemFromDB(nameValue);
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

function getIdFromName(name){
    for(let i = 0; i < currentMenu.length;i++){
        console.log(currentMenu[i].itemName+":"+name)
        if(currentMenu[i].itemName == name){
            console.log(currentMenu[i].itemId)
            return currentMenu[i].itemId;
        }
    }
    console.log("COULD NOT FIND ID FROM ITEM NAME: "+name);
    return null;
}

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
