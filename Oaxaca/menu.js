// Called when menu page is initially loaded
function initMenuAll(){
    let data = requestMenu(0,0,0); // Zero value = none specified
    
    data.then(r => {
        let index = 0;
        document.getElementById("menuSectionAll").innerHTML = ""
        r.forEach(element => {
            
            let item = '<li id="item'+index+'" class="genericMenuItem" data-calories="'+element.Calories+'" data-price="'+element.Price+'">'+element.ItemName+'</li>';
            document.getElementById("menuSectionAll").innerHTML+= item
            index++;
        });
    })
}

// ONLY FOR USE IN TESTING ONLY
// TO ADD ITEM TO YOUR LOCAL TESTING DB RUN THIS IN BROWSER CONSOLE
async function addTestItem(){
    try{
        let response = await fetch("http://localhost:4444/add_item", {
            method:'POST',
            headers:{
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                ItemName: "Tequila",
                Price: 2.50,
                Calories: 12
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