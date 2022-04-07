// $(function(){
//     $('#submitTimeFrame').click(function(){
//         var timeFrame = $('#newTimeFrame').val();
//         $('#newTimeFrame').val('');
//     })
// })

function changeTimeFrame() {
    var inputResult = document.getElementById("newTimeFrame").value;
    document.getElementById("newTimeFrame").value = '';
    chrome.tabs.query({currentWindow: true, active: true}, function(tabs) {
        var activeTab = tabs[0];
        chrome.tabs.sendMessage(activeTab.id, {"timeFrame": inputResult}, function(response) {
            console.log(response.message);
        });
       });
}
   
document.addEventListener('DOMContentLoaded', function() {
    document.querySelector('button').addEventListener('click', changeTimeFrame, false);
}, false)