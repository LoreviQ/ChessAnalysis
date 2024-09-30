let moves = "";
let ready = false;
const baseURL = "http://127.0.0.1:5000/";

checkReadiness();
setInterval(checkReadiness, 60000);

chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    console.log("Executing action: " + request.action);
    switch (request.action) {
        case "getTabUrl":
            chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
                console.log("Getting tab URL: " + tabs[0].url);
                sendResponse({ url: tabs[0].url });
            });
            return true;
        case "updateMoveList":
            moves = request.moves;
            if (ready) {
                console.log("Sending moves to server: " + moves);
                sendMovesToServer(moves);
            }
            break;
        case "getMoveList":
            console.log("Sending moves to popup: " + moves);
            sendResponse({ moves: moves });
            break;
        case "getReadiness":
            console.log("Sending readiness to popup: " + ready);
            sendResponse({ ready: ready });
            break;
        case "manualReadinessCheck":
            console.log("Manual readiness check");
            checkReadiness();
            break;
        default:
            break;
    }
});

// Function to send moves to the local server
function sendMovesToServer(moves) {
    fetch(baseURL + "update_moves", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({ moves: moves, gameid: "1" }),
    });
}

// Returns true if endpoint returns 200 OK
function checkReadiness() {
    fetch(baseURL + "readiness", {
        method: "GET",
    })
        .then((response) => {
            if (response.status === 200) {
                ready = true;
            } else {
                ready = false;
            }
        })
        .catch((error) => {
            ready = false;
        });
}
