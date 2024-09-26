// Function to handle changes to move list
function handleChanges(mutationsList, observer) {
    const moveListElement = document.querySelector("wc-simple-move-list");
    if (moveListElement) {
        console.log("moveListElement content:", moveListElement.textContent.trim());
    }
}

// Create an observer instance linked to the callback function
const observer = new MutationObserver(handleChanges);

// Function to start observing the move list element
function startObserving() {
    const moveListElement = document.querySelector("wc-simple-move-list");
    if (moveListElement) {
        observer.observe(moveListElement, {
            attributes: true, // Observe attribute changes
            childList: true, // Observe direct children changes
            subtree: true, // Observe all descendants
        });
        console.log("Started observing moveListElement");
    } else {
        console.log("moveListElement not found");
    }
}

// Start observing the move list element when it is available
const intervalId = setInterval(() => {
    if (document.querySelector("wc-simple-move-list")) {
        startObserving();
        clearInterval(intervalId);
    }
}, 1000);
