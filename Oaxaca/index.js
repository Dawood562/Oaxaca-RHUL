document.addEventListener('DOMContentLoaded', e => {
    initCookies();
});

// Called when page is loaded. This both initialises and resets cookies.
function initCookies() {
    document.cookie = "tableNo=;";
    document.cookie =  "basket=;";
}