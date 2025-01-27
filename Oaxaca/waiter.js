var sock;
var sockInit = false;
var waiterID = -1
var waiterUsername = "";
// When a waiter confirms or cancels an order that order is added to this list
// and checked to make sure a waiter does not duplicate confirm an order
// var cancelConfirmBlacklist = [];
// var deliveredBlacklist = [];
var orderList = [];

document.addEventListener('DOMContentLoaded', e => {
    registerWaiter();
    initSock(); // This should be called within register waiter after registering waiter
});

function getUserNameFromCookies() {
    let cookieData = document.cookie;
    cookieData.split(";").forEach(cookie => {
        indexOfParam = cookie.indexOf("=");
        if (cookie.substring(0, indexOfParam).indexOf("username") != -1) {
            waiterUsername = cookie.substring(indexOfParam + 1, cookie.length);
        }
    })
}

// Returns waiter id from cookies if exists. If not then returns -1
function getWaiterIdFromCookies(){
    let storedID = sessionStorage.getItem("waiterID");
    return storedID;
}

async function registerWaiter() {
    getUserNameFromCookies();
    waiterID = getWaiterIdFromCookies();

    if(waiterUsername.length <= 0){
        console.log("WAITER NOT REGISTERED! INVALID USERNAME: {"+waiterUsername+"}");
        return;
    }

    if(waiterID > 0){
        console.error("Waiter ID already registered in cookies")
        refreshOrders();
        return;
    }

    const response = await fetch("http://localhost:4444/add_waiter", {
        method:"PUT",
        headers:{
            "Content-Type": "application/json"
        },
        body:JSON.stringify({
            "username": waiterUsername,
        })
    }).then(resp => resp.json()).then(data => {
        console.log(data)
        waiterID = Number(data.id);
        storeWaiterID();
        refreshOrders();
    })

}

function storeWaiterID(){
    if(waiterID < 0){
        sessionStorage.setItem("waiterID", -1);
        console.error("Cannot store waiter id - id is null")
        return;
    }

    sessionStorage.setItem("waiterID", waiterID);
}

// On leaving waiter page
//window.onbeforeunload = removeWaiter;

async function removeWaiter(){
    sessionStorage.setItem("waiterID", -1);
    const response = await fetch("http://localhost:4444/remove_waiter",{
        method:"POST",
        headers:{
            "Content-Type": "application/json"
        },
        body:JSON.stringify({
            "id": waiterID,
        })
    })
}

function logWaiterOut(){
    removeWaiter();
    alert("Waiter logged out!");
    document.cookie = "username=";
    document.location.href = "staff.html";
}

function initSock() {
    sock = new WebSocket("ws://localhost:4444/notifications");
    sock.onerror = function (event) {
        // If unsuccessfully connected
        alert("Unsuccessfully to connect to backend websocket 🖥️🔥");
    };
    sock.addEventListener("open", e => {
        // WE NEED TO ADD A WAY TO GET USERS TABLE NUMBER
        sock.send("WAITER:"+waiterUsername);
    });
    sock.addEventListener("message", e => handleMessages(e));
    sockInit = true;
}



function handleMessages(e) {
    console.log("Received WebSocket message:", e.data);

    if (e.data == "WELCOME") {
        console.log("Connected to backend websocket");
    } else if (e.data == "OK") {
        console.log("Notification successfully received");
    } else if (e.data == "SERVICE") {
        // NOTIFICATION SENT BY KITCHEN STAFF TO WAITERS - DO STUFF HERE
        alert("Kitchen has called service!")
    } else if (e.data.includes("HELP")) {
        // NOTIFICATION SENT BY CUSTOMER TO WAITERS - CUSTOMER IS AT TABLE 'tableNumber' - DO STUFF BELOW
        let tableNumber = e.data.split(":")[1]
        alert("Table " + tableNumber + " needs help!")
    } else if (e.data == "REFRESH") {
        refreshOrders();
    } else if (e.data == "NEW") {
        console.log("Notification: New order received");
        refreshOrders();
    } else if (e.data.startsWith("CANCELLED")) {
        let orderId = e.data.split(":")[1];
        console.log(`Notification: Order ${orderId} has been cancelled!`);
        alert(`Order ${orderId} has been cancelled!`);
        // Additional debugging or handling can be added here
    } else if (e.data == "CANCEL_CONFIRMATION") {
        console.log("Notification: Order cancellation confirmed");
        // ... (existing code)
    } else {
        console.log("Unknown message received:", e.data);
        // Display entire message if something unexpected is received
    }
}


async function refreshOrders() {
    if(waiterID < 0){
        console.log("Cannot refresh orders with no waiter id");
        return;
    }

    let table = document.getElementById("order_table");
    table.innerHTML = `<caption class="table-caption">Customer Orders</caption>
                        <tr>
                            <th>Table</th>
                            <th>Time</th>
                            <th>Items</th>
                            <th>Status</th>
                            <th>Payment Status</th>
                            <th></th>
                            <th></th>
                            <th></th>
                        </tr>`;
    let response = await fetch("http://localhost:4444/orders?");
    let data = await response.json();
    orderList = data;

    for (var order of data) {
        if (order.status !== 'Cancelled') {
            table.innerHTML += createOrder(order);
        }
    }
}


function createOrder(order) {
    let items = order.items.map((x) => x.itemId.itemName);
    let itemsStr = "";
    for (var item of items) {
        itemsStr += item + ", ";
    }
    return `<tr data-order-id="${order.orderId}">
        <td>${order.tableNumber}</td>
        <td>${new Date(order.orderTime).toLocaleTimeString()}</td>
        <td>${itemsStr.substring(0, itemsStr.length - 2)}</td>
        <td id="status`+order.orderId+`">${order.status}</td>
        <td>${order.paid ? "Paid" : "Unpaid"}</td>
        <td><button type="button" id="cancelOrderButton`+order.orderId+`" onclick="notifyCancellation(${order.orderId})">Cancel Order</button></td>
        <td><button type="button" id="confirmOrderButton`+order.orderId+`" onclick="notifyConfirmation(${order.orderId})">Confirm Order</button></td>
        <td><button type="button" id="confirmDeliveredButton`+order.orderId+`" onclick="notifyDelivered(${order.orderId})">Mark Delivered</button></td>
    </tr>`;
}

async function notifyCancellation(orderId) {
    if (!sockInit) {
        return console.error("SOCKET NOT INITIALIZED - CANNOT NOTIFY CANCELLATION");
    }

    console.log(`Sending cancellation request for order ${orderId}`);

    try {
        let response = await fetch(`http://localhost:4444/cancel/${orderId}`, {
            method: 'PATCH',
        });

        if (response.ok) {
            console.log(`Cancellation request for order ${orderId} successful`);

        } else if (response.status === 404) {
            console.log(`Order ${orderId} not found.`);
            alert(`Order ${orderId} not found.`);
        } else if (response.status === 409) {
            console.log(`Order ${orderId} has already been cancelled.`);
            alert(`Order ${orderId} has already been cancelled.`);
        } else {
            console.error(`Failed to cancel order ${orderId}. Status: ${response.status}`);
            alert(`Failed to cancel order ${orderId}. Please check console for details.`);
        }
    } catch (error) {
        console.error(`Error while cancelling order ${orderId}: ${error}`);
        alert(`Error while cancelling order ${orderId}. Please check console for details.`);
    }
}

async function notifyConfirmation(orderId) {
    if (!sockInit) {
        return console.error("SOCKET NOT INITIALIZED - CANNOT NOTIFY CONFIRMATION");
    }

    for(let i = 0; i < orderList.length; i++){
        if(orderList[i].orderId == orderId){
            if (orderList[i].status == "Preparing"){
                refreshOrders();
                return console.error("Order already confirmed/cancelled");
            }
        }
    }

    console.log(`Sending confirmation request for order ${orderId}`);

    try {
        let response = await fetch(`http://localhost:4444/confirm/${orderId}`, {
            method: 'PATCH',
        })

        if (response.ok) {
            console.log(`Confirmation request for order ${orderId} successful`);

        } else if (response.status === 404) {
            console.log(`Order ${orderId} not found.`);
            alert(`Order ${orderId} not found.`);
        } else if (response.status === 409) {
            console.log(`Order ${orderId} has already been confirmed.`);
            alert(`Order ${orderId} has already been confirmed.`);
        } else {
            console.error(`Failed to confirm order ${orderId}. Status: ${response.status}`);
            alert(`Failed to confirm order ${orderId}. Please check console for details.`);
        }
    } catch (error) {
        console.error(`Error while confirming order ${orderId}: ${error}`);
        alert(`Error while confirming order ${orderId}. Please check console for details.`);
    }
}

async function notifyDelivered(orderId){

    for(let i = 0; i < orderList.length; i++){
        if(orderList[i].orderId == orderId){
            if (orderList[i].status == "Awaiting Confirmation"){
                alert("Cannot mark order that has not been prepared as delivered!");
                return;
            }
            if (orderList[i].status == "Delivered"){
                refreshOrders();
                return console.error("Order already marked delivered");
            }
        }
    }

    try{
        let response = await fetch(`http://localhost:4444/delivered/${orderId}`, {
            method: 'PATCH',
        });
        if(response.ok){
            console.log(`Successfully marked order ${orderId} as delivered`);
        }else{
            console.error("Failed to mark as delivered:");
            console.log(response);
        }
    }catch (error){
        console.error(`Error while marking order ${orderId} as delivered`);
    }
    
}