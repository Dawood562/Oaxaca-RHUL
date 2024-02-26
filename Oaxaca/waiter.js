var sock
var sockInit = false

document.addEventListener('DOMContentLoaded', e=>{
    registerWaiter();
    initSock(); // This should be called within register waiter after registering waiter
    refreshOrders();
})

var waiterUsername = "";
function getUserNameFromCookies() {
    let cookieData = document.cookie;
    cookieData.split(";").forEach(cookie => {
        indexOfParam = cookie.indexOf("=");
        if(cookie.substring(0, indexOfParam).indexOf("username") != -1){
            waiterUsername = cookie.substring(indexOfParam+1, cookie.length);
        }
    })
}

async function registerWaiter(){
    getUserNameFromCookies();
    let randId = Math.floor(Math.random()*Number.MAX_SAFE_INTEGER);
    console.log(randId);

    if(waiterUsername.length <= 0){
        console.log("WAITER NOT REGISTERED! INVALID USERNAME: {"+waiterUsername+"}");
        return;
    }

    const response = await fetch("http://localhost:4444/add_waiter", {
        method:"PUT",
        headers:{
            "Content-Type": "application/json"
        },
        body:JSON.stringify({
            "id": randId,
            "waiterUsername": waiterUsername,
            "tableNumber": []
        })
    })
    console.log(response);
}

function initSock(){
    sock = new WebSocket("ws://localhost:4444/notifications")
    sock.onerror = function(event){
        // If unsuccessfully connected
        alert("Unsuccessfully to connect to backend websocket 🖥️🔥");
    }
    sock.addEventListener("open", e =>{
        // WE NEED TO ADD A WAY TO GET USERS TABLE NUMBER
        sock.send("WAITER:Ben");
    })
    sock.addEventListener("message", e => handleMessages(e));
    sockInit = true;
}

function handleMessages(e){
    console.log(e)
    if (e.data == "WELCOME"){
        console.log("Connected to backend websocket");
    }else if (e.data == "OK"){
        console.log("Notification successfully received");
    }else if (e.data == "SERVICE"){
        // NOTIFICATION SENT BY KITCHEN STAFF TO WAITERS - DO STUFF HERE
        alert("Kitchen has called service!")
    }else if (e.data.includes("HELP")){
        // NOTIFICATION SENT BY CUSTOMER TO WAITERS - CUSTOMER IS AT TABLE 'tableNumber' - DO STUFF BELOW
        let tableNumber = e.data.split(":")[1]
        alert("Table " + tableNumber + " needs help!")
    }else if(e.data == "NEW"){
        refreshOrders();
    }else{
        console.log(e); // Display entire message if something went wrong for debugging
    }
}

async function refreshOrders(){
    let table = document.getElementById("order_table");
    table.innerHTML = `<caption>Customer Orders</caption>
                        <tr>
                            <th>Table</th>
                            <th>Time</th>
                            <th>Items</th>
                            <th>Status</th>
                            <th></th>
                            <th></th>
                        </tr>`;
    let response = await fetch("http://localhost:4444/orders")
    let data = await response.json()
    for(var order of data) {
        table.innerHTML += createOrder(order)
    }
}

function createOrder(order) {
    let items = order.items.map((x) => x.itemId.itemName);
    let itemsStr = "";
    for(var item of items) {
        itemsStr += item + ", "
    }
    return `<tr>
        <td>${order.tableNumber}</td>
        <td>${new Date(order.orderTime).toLocaleTimeString()}</td>
        <td>${itemsStr.substring(0, itemsStr.length - 2)}</td>
        <td>${order.status}</td>
        <td><button type="submit">Cancel Order</button></td>
        <td><button type="submit">Complete</button></td>
    </tr>`;
}

function notifyConfirmation(){
    if(!sockInit){
        return console.error("SOCKET NOT INITIALISED - CANNOT NOTIFY CONFIRMATION")
    }
    sock.send("CONFIRM");
}

function notifyCancellation(){
    if(!sockInit){
        return console.error("SOCKET NOT INITIALISED - CANNOT NOTIFY CANCELLATION")
    }
    sock.send("CANCEL");
}