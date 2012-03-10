$(function(){
  $("#command-input").keypress(function (e) {
    $("#debug").append(e);
    if ((e.which && e.which == 13) || (e.KeyCode && e.keyCode == 13)) {
    }
  });
});
