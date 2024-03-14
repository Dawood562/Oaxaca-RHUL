var sock
var sockInit = false
var basketData = []

document.addEventListener('DOMContentLoaded', e => {
    initTableNo();
    initBasketData()
    showOrdersToPage();
})

function initTableNo() {
    // Retrieve table number from cookies
    tableNo = null;
    let cookieData = document.cookie;
    cookieData.split(";").forEach(cookie => {
        indexOfParam = cookie.indexOf("=");
        if (cookie.substring(0, indexOfParam).indexOf("username") != -1 && cookie.substring(indexOfParam, cookie.length).indexOf("table") != -1) {
            tableNo = cookie.substring(indexOfParam + 6, cookie.length);
        }
    })

    if (tableNo != null) {
        tableNoField = document.getElementById("tableNumber");
        tableNoField.value = tableNo;
        tableNoField.setAttribute('readonly', true);
        initSock(tableNo);
    }
}

function initBasketData() {
    let splitCookies = document.cookie.split(";");
    let basketCookies = "";
    let basketEmpty = true;
    splitCookies.forEach(cookie => {
        if (cookie.indexOf("basket=") >= 0) {
            basketCookies = cookie.split("basket=")[1].split("#");
            if (basketCookies.length > 1) {
                basketEmpty = false;
            }
        }
    })

    if (basketEmpty) {
        console.log("Basket is empty!");
        return;
    }

    basketCookies.forEach(item => {
        if (item.length > 0) {
            // For each item in basket
            let splitData = item.split(",")
            let itemData = {
                id: splitData[0],
                name: splitData[1],
                price: splitData[2],
                calories: splitData[3],
                quantity: splitData[4],
                imageURL: splitData[5] // Add the imageURL property
            }
            basketData.push(itemData)
        }
    });
}

// Called when page is initialised and initialises websocket connection
function initSock(tableNo) {
    sock = new WebSocket("ws://localhost:4444/notifications")
    sock.onerror = function (event) {
        // If unsuccessfully connected
        alert("Unsuccessfully to connect to backend websocket 🖥️🔥")
    }
    sock.addEventListener("open", e => {
        // WE NEED TO ADD A WAY TO GET USERS TABLE NUMBER
        sock.send("CUSTOMER:".concat(tableNo))
    })

    sock.addEventListener("message", e => handleMessages(e))
    document.cookie = "username=table".concat(tableno); // Set "username" entry to table number
}

function closeSock() {
    sock.close();
}

// Displays connection and communication status to console
function handleMessages(e) {
    if (e.data == "WELCOME") {
        console.log("Connected to backend websocket")
    } else if (e.data == "OK") {
        console.log("Notification successfully received")
    } else if (e.data == "REFRESH"){
        refreshOrderStatus()
    } else {
        console.log(e) // Display entire message if something went wrong for debugging
    }
}

// Notifies backend a customer needs help
function sendHelp() {
    sock.send("HELP")
}

function showConfirmationSection() {
    document.getElementById('confirmationSection').style.display = 'block';
}

function showOrdersToPage() {
    let orderStore = document.getElementById("orderHeading");

    if (basketData.length === 0) {
        let emptyMessage = document.createElement('div');
        emptyMessage.innerHTML = `
            <p>Nothing here.</p>
            <button onclick="location.href='menu.html'">Start Shopping</button>
        `;
        emptyMessage.style.textAlign = 'center';
        emptyMessage.style.marginTop = '2rem';
        emptyMessage.style.fontWeight = '500';
        orderStore.appendChild(emptyMessage);
    } else {
        basketData.forEach(itemData => {
            let item = document.createElement('div'); // Change to 'div' instead of 'li'
            item.className = "orderPageItem";
            item.innerHTML = `
                <div class='orderItemEntry'>
                    <img src="http://localhost:4444/image/${itemData.imageURL}" alt="${itemData.name}" class="orderItemImage">
                    <div class="orderItemDetails">
                        <label class='orderPageItemDataN'>${itemData.name}</label>
                        <label class='orderPageItemDataP'>£${(Number(itemData.price) * Number(itemData.quantity)).toFixed(2)}</label>
                        <label class='orderPageItemDataC'>${Number(itemData.calories) * Number(itemData.quantity)} kcal</label>
                        <label class='orderPageItemDataQ'>Quantity: ${itemData.quantity}</label>
                    </div>
                    <button class="removeButton" onclick='removeOrderFromList(${itemData.id})'>Remove<i class="fa fa-trash"></i></button>
                </div>`;
            orderStore.appendChild(item);
        });
    }
}

function removeOrderFromList(id) {
    let removedItem = false;
    for (let i = 0; i < basketData.length; i++) {
        if (basketData[i].id == id) {
            basketData.splice(i, 1);
            i = basketData.length;
            removedItem = true;
            console.log("New size of basketData: " + basketData.length)
            removeOrderFromCookie(id);
        }
    }

    if (removedItem) {
        // Clear orders from page:
        removeAllOrders();
        showOrdersToPage();
    }

}

function removeAllOrders() {
    let orderList = document.getElementById("orderHeading").childNodes;
    // Backwards for loop because when element 1 removed, element 2 takes 1's place
    for (let i = orderList.length - 1; i > 0; i--) {
        orderList[1].remove();
    }
}

function removeOrderFromCookie(id) {
    let cookies = document.cookie.split("basket=")[1].split("#");
    let newCookieData = "basket=";
    cookies.forEach(cookie => {
        if (cookie.length > 0) {
            // If cookie doesnt match expected cookie to remove then add it back to cookie 
            if (cookie.split(",")[0] != id) {
                newCookieData += cookie + "#";
            }
        }
    })
    // Update cookies
    document.cookie = newCookieData;
}

var orderID = -1;
function submitOrder() {

    try {
        x = sock.readyState;
    } catch (e) {
        initSock(tableNum);
    }

    let tableNum = Number(document.getElementById('tableNumber').value);

    if (basketData.length == 0) {
        alert('Your basket is empty. Add items before submitting your order.');
        return;
    }

    // Fill in body of item array to send to backend
    let totalBill = 0.0;
    let itemObjects = [];
    basketData.forEach(element => {
        itemObjects.push({ item: Number(element.id), notes: "Stuff" });
        totalBill += Number(element.price);
    });


    fetch('http://localhost:4444/add_order?', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            bill: totalBill,
            tableNumber: tableNum,
            items: itemObjects
        })
    })
        .then((res) => {
            if (res.ok) {
                return res.text();
            }
            throw new Error("Could not place order");
        })
        .then((data) => {
            console.log(`Order of id ${data} submitted successfully`);
            orderID = Number(data);

            localStorage.removeItem('order'); // Clears the basket after a successful order
            localStorage.setItem("orderID", data);

            document.cookie = "basket=";
            removeAllOrders();
            replaceWithOrderStatus();
        })
        .catch((err) => {
            alert(err);
        });

    // Change table number field so it cannot be changed
    tableNoField = document.getElementById("tableNumber");
    tableNoField.setAttribute('readonly', true);

}

function replaceWithOrderStatus() {
    let orderList = document.getElementById("orderHeading");
    let submissionNotification = document.createElement('div');
    submissionNotification.innerHTML = `
    <label class='orderPageItem'>Submitted order!</label> <br>
    <label>Status: </label> <label id=orderStatus>Awaiting confirmation...</label> <button onclick="refreshOrderStatus()" id="refreshStatusButton">Refresh</button>
    `;
    orderList.appendChild(submissionNotification);
}

async function refreshOrderStatus() {
    if (orderID < 0) {
        console.error("Cannot refresh order with no order id!")
        return;
    }
    let response = await fetch(`http://localhost:4444/status/${orderID}`, {
        method: 'GET',
    }).then((res) => res.text()).then((data) => {
        console.log("Retrieved status of: " + data);
        document.getElementById("orderStatus").innerHTML = data;
    });
}

function sendPayment() {
    let id = localStorage.getItem("orderID");
    fetch("http://localhost:4444/pay/" + id, {
        method: "PATCH"
    })
        .then((res) => {
            if (res.ok) {
                alert("Payment received. Thank you for dining with us :)");
            } else {
                alert("There was a problem with payment. Please try again later.");
            }
        })
}
