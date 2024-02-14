var sock

document.addEventListener("DOMContentLoaded", function(){
    sock = new WebSocket("ws://localhost:4444/notifications")
    sock.addEventListener("open", e =>{
        console.log("Connected to backend websocket")
    })
})