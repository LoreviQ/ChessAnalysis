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
    // Displays the move list in popup
    const moveListDiv = document.getElementById("move-list");
    chrome.runtime.sendMessage({ action: "getMoveList" }, function (response) {
        console.log(response);
        if (response && response.moves) {
            moves = response.moves;
            moveList = moves.split(" ");
            moveListDiv.innerHTML = generateMoveListHTML(moveList);
        }
    });
});

// Function to generate HTML for move list
function generateMoveListHTML(moveList) {
    let html = "<ol>";
    for (let i = 0; i < moveList.length; i += 3) {
        const player1move = moveList[i + 1];
        const player2move = moveList[i + 2];
        html += `<li>${player1move} : ${player2move}</li>`;
    }
    html += "</ol>";
    return html;
}
