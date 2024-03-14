window.onload = function () {
    let form = document.getElementById("loginform").addEventListener("submit", function (e) {
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

function closeLoginBox() {
    document.getElementById('loginbox').style.display = 'none';
}

function storeAccDetails() {
    document.cookie = "username="; // Clear username entry
    let enteredUsername = document.getElementById("UN").value;
    document.cookie = "username=" + enteredUsername;
    console.log("Cookies: " + document.cookie);
}

function goToKitchen() {
    window.location.href = "kitchen.html";
}

function goToWaiter() {
    window.location.href = "waiter.html";
}