window.onload = function () {
    let form = document.getElementById("loginform").addEventListener("submit", function (e) {
        e.preventDefault();

        var formData = new FormData(form);

        let un = document.getElementById("UN").value;
        let pw = document.getElementById("PW").value;
        // TODO: Authenticate the user using values un and pw
        if (authUser(un, pw)) {
            closeLoginBox();
            showNavButtons();
        }
    })
    
    checkIfLoggedIn();
}
// Skip login dialogue if already logged in
function checkIfLoggedIn() {
    // Grab username
    let username = "";
    let cookieData = document.cookie;
    cookieData.split(";").forEach(cookie => {
        indexOfParam = cookie.indexOf("=");
        if (cookie.substring(0, indexOfParam).indexOf("username") != -1) {
            username = cookie.substring(indexOfParam + 1, cookie.length);
        }
    })

    if (username.length >= 1) {
        if (sessionStorage.getItem("waiterID") != -1) {
            goToWaiter();
        } else {
            closeLoginBox();
            showNavButtons();
        }
    }
}


function closeLoginBox() {
    document.getElementById('loginbox').style.display = 'none';
}

function showNavButtons() {
    document.getElementById('kitchen-button').style.visibility = "visible";
    document.getElementById('waiter-button').style.visibility = "visible";
}

function storeAccDetails() {
    document.cookie = "username="; // Clear username entry
    let enteredUsername = document.getElementById("UN").value;
    let enteredPassword = document.getElementById("PW").value;
    if (!authUser(enteredUsername, enteredPassword)) {
        return;
    }
    document.cookie = "username=" + enteredUsername;
    console.log("Cookies: " + document.cookie);
}

function authUser(un, pw) {
    if (un.length >= 1) {
        return true;
    }
    // Add authentication (not in spec)
    return false;
}

function goToKitchen() {
    window.location.href = "kitchen.html";
}

function goToWaiter() {
    window.location.href = "waiter.html";
}