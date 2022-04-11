function changeTimeFrame() {
    var inputResult = document.getElementById("newTimeFrame").value;
    document.getElementById("newTimeFrame").value = '';
    chrome.tabs.query({currentWindow: true, active: true}, function(tabs) {
        var activeTab = tabs[0];
        chrome.tabs.sendMessage(activeTab.id, {"timeFrame": inputResult}, function(response) {
            try {
                console.log(response.message);
            } catch (error) {
                console.log(error);
                alert("Error: time frame could not be submitted. Make sure you are on the github tab.");
            }
        });
    });
}
   
document.addEventListener('DOMContentLoaded', function() {
    document.querySelector('button').addEventListener('click', changeTimeFrame, false);
}, false)