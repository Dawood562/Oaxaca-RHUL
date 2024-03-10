var sock
var sockInit = false

document.addEventListener('DOMContentLoaded', e => {
    initSock(); // This should be called within register waiter after registering waiter
    refreshOrders();
})

function initSock() {
    sock = new WebSocket("ws://localhost:4444/notifications")
    sock.onerror = function (event) {
        // If unsuccessfully connected
        alert("Unsuccessfully to connect to backend websocket 🖥️🔥");
    }
    sock.addEventListener("open", e => {
        // WE NEED TO ADD A WAY TO GET USERS TABLE NUMBER
        sock.send("KITCHEN");
    })
    sock.addEventListener("message", e => handleMessages(e));
    sockInit = true;
}

function notifyService() {
    sock.send("SERVICE");
}

function handleMessages(e) {
    console.log(e)
    if (e.data == "WELCOME") {
        console.log("Connected to backend websocket");
    } else if (e.data == "OK") {
        console.log("Notification successfully received");
    } else if (e.data == "CONFIRM") {
        refreshOrders();
    } else if (e.data == "CANCEL") {
        refreshOrders();
    } else if (e.data == "REFRESH") {
        refreshOrders();
    } else {
        console.log(e); // Display entire message if something went wrong for debugging
    }
}

async function refreshOrders() {
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
    let response = await fetch("http://localhost:4444/orders?confirmed=true")
    let data = await response.json()
    for (var order of data) {
        table.innerHTML += createOrder(order)
    }
}

function createOrder(order) {
    let items = order.items.map((x) => x.itemId.itemName);
    let itemsStr = "";
    for (var item of items) {
        itemsStr += item + ", "
    }
    return `<tr>
        <td>${order.tableNumber}</td>
        <td>${new Date(order.orderTime).toLocaleTimeString()}</td>
        <td>${itemsStr.substring(0, itemsStr.length - 2)}</td>
        <td>${order.status}</td>
        <td><button onclick="completeOrder(${order.id})">Complete Order</td>
    </tr>`;
}

function confirmOrder(id) {
    fetch("http://localhost:4444/confirm/" + id, {
        method: "PATCH"
    })
        .then((res) => {
            if (!res.ok) {
                alert("Failed to confirm order!");
            }
        })
}

function cancelOrder(id) {
    fetch("http://localhost:4444/cancel/" + id, {
        method: "PATCH"
    })
        .then((res) => {
            if (!res.ok) {
                alert("Failed to cancel order!");
            }
        })
}