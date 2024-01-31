function requestMenu(){
    let data = requestMenuAll()
    
    data.then(r => {
        console.log(r)
        // UPDATE MENU BELOW HERE USING DATA FROM r
    })
}

// Asynchronously
async function requestMenu(scope){
    try{
        let response = await fetch("http://localhost:4444/menu", {
            method:"POST",
            headers:{
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                Scope: scope
            })
        })

        if(!response.ok){
            console.log("ERROR fetching menu")
        }

        let data = await response.json()
        return data
    }catch(error){
        console.error(error)
        return null
    }
}