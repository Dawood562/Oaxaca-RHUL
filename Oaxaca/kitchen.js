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
        sock.send("KITCHEN");
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
    }else if(e.data == "CONFIRM"){
        // DO STUFF WHEN ORDER IS CONFIRMED
    }else if(e.data == "CANCEL"){
        // CLEANUP WHEN ORDER IS CANCELLED
    }
    else{
        console.log(e); // Display entire message if something went wrong for debugging
    }
}

function notifyService(){
    if(!sockInit){
        console.error("SOCKET NOT INITIALISED - CANNOT NOTIFY SERVICE")
    }
    sock.send("SERVICE")
}