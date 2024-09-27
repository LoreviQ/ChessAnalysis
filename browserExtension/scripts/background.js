let moves = "";

chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    if (request.action === "getTabUrl") {
        chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
            sendResponse({ url: tabs[0].url });
        });
        return true;
    } else if (request.action === "updateMoveList") {
        moves = request.moves;
        sendMovesToServer(moves);
    } else if (request.action === "getMoveList") {
        sendResponse({ moves: moves });
    }
});

// Function to send moves to the local server
function sendMovesToServer(moves) {
    fetch("http://127.0.0.1:5000/update_moves", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ moves: moves }),
    })
        .then((response) => response.json())
        .then((data) => console.log("Success:", data))
        .catch((error) => console.error("Error:", error));
}
