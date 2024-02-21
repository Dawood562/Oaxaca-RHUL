window.onload = function() {
    let form = document.getElementById("loginform").addEventListener("submit", function(e) {
        e.preventDefault();

        var formData = new FormData(form);
    
        let un = document.getElementById("UN").value;
        let pw = document.getElementById("PW").value;
        // TODO: Authenticate the user using values un and pw

        document.getElementById('loginbox').style.visibility = "hidden";
        document.getElementById('kitchen-button').style.visibility = "visible";
        document.getElementById('waiter-button').style.visibility = "visible";
    })
}