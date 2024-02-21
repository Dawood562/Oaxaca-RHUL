var sock
var sockInit = false
var basketData = []

document.addEventListener('DOMContentLoaded', e=>{
    initSock();
    initBasketData()
    showOrdersToPage();
})

function initBasketData(){
    let basketCookies = document.cookie.split("basket=")[1].split("#");
    basketCookies.forEach(item => {
        if(item.length > 0){
            // For each item in basket
            let splitData = item.split(",")
            let itemData = {
                id:splitData[0],
                name:splitData[1],
                price:splitData[2],
                calories:splitData[3],
                quantity:splitData[4]
            }
            basketData.push(itemData)
        }     
    });
    console.log(basketData)

}

// Called when page is initialised and initialises websocket connection
function initSock(){
    sock = new WebSocket("ws://localhost:4444/notifications")
    sock.onerror = function(event){
        // If unsuccessfully connected
        alert("Unsuccessfully to connect to backend websocket 🖥️🔥")
    }
    sock.addEventListener("open", e =>{
        // WE NEED TO ADD A WAY TO GET USERS TABLE NUMBER
        sock.send("CUSTOMER:1")
    })

    sock.addEventListener("message", e => handleMessages(e))
}

// Displays connection and communication status to console
function handleMessages(e){
    if (e.data == "WELCOME"){
        console.log("Connected to backend websocket")
    }else if (e.data == "OK"){
        console.log("Notification successfully received")
    }else{
        console.log(e) // Display entire message if something went wrong for debugging
    }
}

// Notifies backend a customer needs help
function sendHelp(){
    sock.send("HELP")
}

// Notifies backend a customer has ordered something
function notifyNew(){
    sock.send("NEW")
}

   function showConfirmationSection() {
  document.getElementById('confirmationSection').style.display = 'block';
}

function showOrdersToPage(){
    let orderStore = document.getElementById("orderHeading");
    basketData.forEach(itemData => {
        let item = document.createElement('li');
        item.className = "orderPageItem"
        item.innerHTML = `
            <label class='orderPageItemData'>${itemData.name}</label>
            <label class='orderPageItemData'>quantity: ${itemData.quantity}</label>
            <label class='orderPageItemData'>Calories: ${Number(itemData.calories) * Number(itemData.quantity)} kcal</label>
            <label class='orderPageItemData'>Price: £${(Number(itemData.price) * Number(itemData.quantity)).toFixed(2)}</label>
            <button class="removeButton"><i class = "fa fa-trash"></i></button>`
        orderStore.appendChild(item)
    });
}

function submitOrder() {
    let tableNum = Number(document.getElementById('tableNumber').value);

    if (basketData.length == 0) {
        alert('Your basket is empty. Add items before submitting your order.');
        return;
    }

    // Fill in body of item array to send to backend
    let totalBill = 0.0;
    let itemObjects = [];
    basketData.forEach(element => {
        itemObjects.push({item: Number(element.id), notes:"Stuff"});
        totalBill += Number(element.price);
    });

    //  Checking if the WebSocket connection is established
    if (sock.readyState !== WebSocket.OPEN) {
        alert('Error: WebSocket connection not established. Please refresh the page and try again.');
        return;
    }

    
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
    .then(() => {
        console.log('Order submitted successfully');

        localStorage.removeItem('order'); // Clears the basket after a successful order
        updateOrderDetails();
        notifyNew(); // Notify backend about the new order
    });
}
