let moveList = "";

chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    if (request.action === "getTabUrl") {
        chrome.tabs.query({ active: true, currentWindow: true }, (tabs) => {
            sendResponse({ url: tabs[0].url });
        });
        return true;
    } else if (request.action === "updateMoveList") {
        moveList = request.moves;
    } else if (request.action === "getMoveList") {
        sendResponse({ moves: moveList });
    }
});
