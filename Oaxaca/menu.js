function requestMenuAll(){
    let data = requestMenu()
    
    data.then(r => {
        console.log(r)
    })
}

// Fetches data from backend or throws error to console and provides example menu data
async function requestMenu(MenuScope){
    try{
        let response = await fetch("http://localhost:4444/menu", {
            method:"POST",
            headers:{
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                Scope: MenuScope
            })
        })

        if(!response.ok){
            console.log("ERROR fetching menu")
        }

        let data = await response.json()
        return data
    }catch(error){
        console.error(error)
        // Return example error if unable to connect to backend
        return JSON.stringify(
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