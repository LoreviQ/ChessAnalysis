document.addEventListener("DOMContentLoaded", function () {
    chrome.runtime.sendMessage({ action: "getTabUrl" }, function (response) {
        const url = response.url;
        const connectionStatusDiv = document.getElementById("connection-status");

        if (url.includes("chess.com")) {
            connectionStatusDiv.textContent = "Connected to chess.com";
        } else if (url.includes("lichess.org")) {
            connectionStatusDiv.textContent = "Lichess not yet supported";
        } else {
            connectionStatusDiv.textContent = "Unsupported website";
        }
    });
});
