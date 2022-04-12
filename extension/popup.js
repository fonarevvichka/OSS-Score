function changeTimeFrame() {
    var inputResult = document.getElementById("newTimeFrame").value;
    if (inputResult != "") {
        chrome.storage.sync.set({key: inputResult}, function() {
            console.log('time frame is set to ' + inputResult);
        });    
        document.getElementById("newTimeFrame").value = '';
        updateTimeFrame(inputResult);
        sendMessage(inputResult);
    }
}

function sendMessage(timeFrame) {
    chrome.tabs.query({currentWindow: true, active: true}, function(tabs) {
        var activeTab = tabs[0];
        chrome.tabs.sendMessage(activeTab.id, {"timeFrame": timeFrame}, function(response) {
            // TODO handle/silence error in html when user on wrong tab
            try {
                console.log(response.message);
            } catch (error) {
                console.log(error);
                alert("Error: time frame could not be submitted. Make sure you are on the github tab.");
            }
        });
    });
}

function updateTimeFrame(timeFrame) {
    if (timeFrame == null) {
        document.getElementById('current').innerHTML = "Current Time Frame is infinite";
    }
    else if (timeFrame == "1") {
        document.getElementById('current').innerHTML = "Current Time Frame is " + timeFrame + " month";
    } else {
        document.getElementById('current').innerHTML = "Current Time Frame is " + timeFrame + " months";
    }
}

function resetTimeFrame() {
    updateTimeFrame(null);
    chrome.storage.sync.remove(['key'], function() {
        console.log("removed time frame");
    });
    sendMessage(null);
}

document.addEventListener('DOMContentLoaded', function() {
    var buttons = document.querySelectorAll('button');
    console.log(buttons.length);
    buttons[0].addEventListener('click', changeTimeFrame, false); // submit button
    buttons[1].addEventListener('click', resetTimeFrame, false); // reset button
    chrome.storage.sync.get(['key'], function(result) {
        if (result.key != null) { // there is a value stored
            updateTimeFrame(result.key);
            console.log('time frame currently is ' + result.key);
        } else { // no value stored
            console.log("time frame currently is infinite")
        }

    });
}, false)