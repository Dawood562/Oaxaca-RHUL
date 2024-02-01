// Called when menu page is initially loaded
function initMenuAll(){
    let data = requestMenu(0,0,0); // Zero value = none specified
    
    data.then(r => {
        console.log(r)
        let index = 0;
        document.getElementById("menuSectionAll").innerHTML = ""
        r.forEach(element => {
            let item = '<li id="item'+index+'" class="genericMenuItem" data-calories="'+element.calories+'" data-price="'+element.price+'">'+element.itemName+'</li>';
            document.getElementById("menuSectionAll").innerHTML+= item
            index++;
        });
    })
}

var editMode = false;
// Toggle edit mode
function editMenu(){
    // Add textfield to bottom of menu list
    // Add plus button to the right of text field

    // Enable edit mode
    if(!editMode){
        document.getElementById("menuSectionAll").innerHTML += "<div id='newItemDiv'>   <label>Name:</label><input type='text' id='newItemNameField'>  <label>Price:</label><input type='text' id='newItemPriceField'>  <label>Calories:</label><input type='text' id='newItemCaloriesField'>    <button onclick='addMenuItem()'>+</button></div>";
        editMode = true;
    }else{
        document.getElementById("newItemDiv").remove();
        editMode = false;
    }
}

function addMenuItem(){
    let nameValue = document.getElementById("newItemNameField").value;
    let priceValue = parseFloat(document.getElementById("newItemPriceField").value);
    let caloriesValue = parseInt(document.getElementById("newItemCaloriesField").value);
    addItemToDB(nameValue, priceValue, caloriesValue);
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

// Fetches data from backend or throws error to console and provides example menu data
async function requestMenu(userSearchTerm, userMaxPrice, userMaxCalories){
    try{
        let response = await fetch("http://localhost:4444/menu?" + new URLSearchParams({
            ItemName: userSearchTerm,
            Price: userMaxPrice,
            userMaxCalories 
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
            ItemName: "Tequila",
            Price: 6.90,
            Calories: 12,
            Type: "DRINK"
        }, {
            ItemName: "Olives",
            Price: 7.99,
            Calories: 165,
            Type: "APPETIZER"
        }, {
            ItemName: "Mozzarella Sticks",
            Price: 8.99,
            Calories: 349,
            Type: "ENTREES"
        }, {
            ItemName: "Ice Cream",
            Price: 7.42,
            Calories: 632,
            Type: "DESSERT"
        }])
    }
}