// $(function(){
//     $('#submitTimeFrame').click(function(){
//         var timeFrame = $('#newTimeFrame').val();
//         $('#newTimeFrame').val('');
//     })
// })

function changeTimeFrame() {
    var inputResult = document.getElementById("newTimeFrame").value;
    document.getElementById("newTimeFrame").value = '';
   }
   
document.addEventListener('DOMContentLoaded', function() {
document.querySelector('button').addEventListener('click', changeTimeFrame, false);
}, false)