function changeTimeFrame() {
    var inputResult = document.getElementById("newTimeFrame").value;
    if(isInt(inputResult)){
        if (inputResult != "") {
            chrome.storage.sync.set({key: inputResult}, function() {
                console.log('time frame is set to ' + inputResult);
            });    
            document.getElementById("newTimeFrame").value = '';
            sendMessage(inputResult);
        }      
    } else {
        alert("Error: time frame must be a positive integer");
    }

}

function isInt(str) {
    return /^[1-9]\d*$/.test(str);
}

function sendMessage(timeFrame) {
    chrome.tabs.query({currentWindow: true, active: true}, function(tabs) {
        var activeTab = tabs[0];
        chrome.tabs.sendMessage(activeTab.id, {"timeFrame": timeFrame}, function(response) {
            if (!window.chrome.runtime.lastError) {
                console.log(response.message);
                updateTimeFrame(timeFrame);
            } else {
                console.log(chrome.runtime.lastError.message);
                alert("Error: time frame could not be submitted. Make sure you are on the github tab.");
            }

        });
    });
}

function updateTimeFrame(timeFrame) {
    if (timeFrame == null) {
        document.getElementById('current').innerHTML = "Current Time Frame: 12 months";
    }
    else if (timeFrame == "1") {
        document.getElementById('current').innerHTML = "Current Time Frame: " + timeFrame + " month";
    } else {
        document.getElementById('current').innerHTML = "Current Time Frame: " + timeFrame + " months";
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
    var submit = document.getElementById('submit');
    var reset = document.getElementById('reset');
    submit.addEventListener('click', changeTimeFrame, false); // submit button
    reset.addEventListener('click', resetTimeFrame, false); // reset button
    chrome.storage.sync.get(['key'], function(result) {
        if (result.key != null) { // there is a value stored
            updateTimeFrame(result.key);
            console.log('time frame currently is ' + result.key);
        } else { // no value stored
            console.log("time frame currently is 12 months")
        }

    });
}, false)