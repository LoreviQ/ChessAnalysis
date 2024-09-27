let moves = "";
let ready = false;
const baseURL = "http://127.0.0.1:5000/";

const intervalId = setInterval(() => {
    checkReadiness();
}, 1000);

chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    switch (request.action) {
        case "getTabUrl":
            chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
                sendResponse({ url: tabs[0].url });
            });
            return;
        case "updateMoveList":
            moves = request.moves;
            if (ready) {
                sendMovesToServer(moves);
            }
            break;
        case "getMoveList":
            sendResponse({ moves: moves });
            break;
        case "getReady":
            sendResponse({ ready: ready });
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
        body: JSON.stringify({ moves: moves }),
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
