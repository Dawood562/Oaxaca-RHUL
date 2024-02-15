var sock
var sockInit = false

document.addEventListener('DOMContentLoaded', e=>{
    initSock()
})

// Called when page is initialised and initialises websocket connection
function initSock(){
    sock = new WebSocket("ws://localhost:4444/notifications")
    sock.onerror = function(event){
        // If unsuccessfully connected
        alert("Unsuccessfully to connect to backend websocket 🖥️🔥")
    }
    sock.addEventListener("open", e =>{
        // WE NEED TO ADD A WAY TO GET USERS TABLE NUMBER
        sock.send("CUSTOMER:7")
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
function notifyHelp(){
    sock.send("HELP")
}

// Notifies backend a customer has ordered something
function notifyNew(){
    sock.send("NEW")
}

   function showConfirmationSection() {
  document.getElementById('confirmationSection').style.display = 'block';
}


function submitOrder() {
     alert('you have submitted sucessfully');
} 
// havent tested yet
/* let order = JSON.parse(localStorage.getItem('order')) || [];

  if (order.length === 0) {
    alert('Your basket is empty. Add items before submitting your order.');
    return;
  }

 let orderType = document.getElementById('orderType').value;
 let tableNumber = document.getElementById('tableNumber').value;

 let orderData = {
   orderType: orderType,
  tableNumber: tableNumber,
    orderDetails: order,
  };

 //  Checking if the WebSocket connection is established
 if (sock.readyState !== WebSocket.OPEN) {
   alert('Error: WebSocket connection not established. Please refresh the page and try again.');
   return;
 }

  fetch('http://localhost:4444/add_order', {
   method: 'POST',
    headers: {
     'Content-Type': 'application/json',
   },
   body: JSON.stringify(orderData),
  })
    .then(response => response.json())
    .then(data => {
     console.log('Order submitted successfully:', data);

     localStorage.removeItem('order'); // Clears the basket after a successful order
      updateOrderDetails();
      updateBasketIcon();
      notifyNew(); // Notify backend about the new order
    })
    .catch(error => {
      console.error('Error submitting order:', error);
    });
}
*/