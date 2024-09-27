document.addEventListener("DOMContentLoaded", function () {
    console.log("loading popup.js");

    // Displays the connection status based on the URL of the active tab
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

    // displays the server status
    const serverStatusDiv = document.getElementById("server-status");
    chrome.runtime.sendMessage({ action: "getReady" }, function (response) {
        if (response.ready) {
            serverStatusDiv.textContent = "Server is ready";
        } else {
            serverStatusDiv.textContent = "Server is not ready";
        }
    });

    // Displays the move list in popup
    const moveListDiv = document.getElementById("move-list");
    chrome.runtime.sendMessage({ action: "getMoveList" }, function (response) {
        console.log(response);
        if (response && response.moves) {
            moves = response.moves;
            moveListDiv.innerHTML = generateMoveListHTML(moves);
        }
    });
});

// Function to generate HTML for move list
function generateMoveListHTML(moves) {
    let html = "<ol>";
    for (let i = 0; i < moves.length; i += 3) {
        const player1move = moves[i + 1];
        const player2move = moves[i + 2];
        html += `<li>${player1move} : ${player2move}</li>`;
    }
    html += "</ol>";
    return html;
}
