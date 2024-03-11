var sock;
var sockInit = false;
var waiterID = -1
var waiterUsername = "";

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

async function registerWaiter() {
    getUserNameFromCookies();
    let randId = Math.floor(Math.random() * Number.MAX_SAFE_INTEGER);
    console.log(randId);

    if (waiterUsername.length <= 0) {
        console.log("WAITER NOT REGISTERED! INVALID USERNAME: {" + waiterUsername + "}");
        return;
    }

    const response = await fetch("http://localhost:4444/add_waiter", {
        method: "PUT",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            "id": randId,
            "waiterUsername": waiterUsername,
            "tableNumber": []
        })
    }).then(resp => resp.json()).then(data => {
        console.log(data)
        waiterID = Number(data.id);
        refreshOrders();
    })
}

// On leaving waiter page
document.addEventListener("beforeunload", (e) => {

    // Unregister waiter
    removeWaiter();
})

function removeWaiter() {

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
                        </tr>`;
    let response = await fetch("http://localhost:4444/orders?" + new URLSearchParams({
        "waiterId": waiterID
    }));
    let data = await response.json();

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
        <td>${order.status}</td>
        <td><button type="button" onclick="notifyCancellation(${order.orderId})">Cancel Order</button></td>
        <td><button type="button" onclick="notifyConfirmation(${order.orderId})">Confirm Order</button></td>
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

            // Update the order status to "Cancelled" directly on the client side
            let row = document.querySelector(`[data-order-id="${orderId}"]`);
            if (row) {
                row.querySelector('td:nth-child(4)').textContent = 'Cancelled';
                row.querySelector('td:nth-child(5) button').disabled = true;
                row.querySelector('td:nth-child(6) button').disabled = true;
                alert(`Order ${orderId} has been cancelled.`);
            } else {
                console.error(`Row with order ID ${orderId} not found.`);
            }
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

    console.log(`Sending confirmation request for order ${orderId}`);

    try {
        let response = await fetch(`http://localhost:4444/confirm/${orderId}`, {
            method: 'PATCH',
        });

        if (response.ok) {
            console.log(`Confirmation request for order ${orderId} successful`);

            // Update the order status to "Confirmed" directly on the client side
            let row = document.querySelector(`[data-order-id="${orderId}"]`);
            if (row) {
                row.querySelector('td:nth-child(4)').textContent = 'Confirmed';
                row.querySelector('td:nth-child(5) button').disabled = true;
                row.querySelector('td:nth-child(6) button').disabled = true;
                alert(`Order ${orderId} has been confirmed.`);
            } else {
                console.error(`Row with order ID ${orderId} not found.`);
            }
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
