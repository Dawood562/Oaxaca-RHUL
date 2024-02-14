var sock
var sockInit = false

document.addEventListener('DOMContentLoaded', e=>{
    initSock()
})

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
    if (e.data == "WELCOME"){
        console.log("Connected to backend websocket");
    }else if (e.data == "OK"){
        console.log("Notification successfully received");
    }else if (e.data == "SERVICE"){
        // NOTIFICATION SENT BY KITCHEN STAFF TO WAITERS - DO STUFF HERE
    }
    else{
        console.log(e); // Display entire message if something went wrong for debugging
    }
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
