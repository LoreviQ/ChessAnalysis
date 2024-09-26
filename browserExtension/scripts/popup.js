document.addEventListener("DOMContentLoaded", function () {
    chrome.runtime.sendMessage({ action: "getTabUrl" }, function (response) {
        const url = response.url;
        const moveListDiv = document.getElementById("move-list");

        if (url.includes("chess.com")) {
            moveListDiv.textContent = "Connected to chess.com";
        } else if (url.includes("lichess.org")) {
            moveListDiv.textContent = "Lichess not yet supported";
        } else {
            moveListDiv.textContent = "Unsupported website";
        }
    });
});
