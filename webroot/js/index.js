
var indexInit = (function() {
    
    return {
        init: function() {
            $.getJSON("/operations", function(data) {
                var lst = [];
                $.each(data, function(key, val) {
                    lst.push("<li id='li_" + val.name + "'><a href='" + val['link-url'] + "'>" + val.title + "</a></li>");
                });
                $("#operation-list").html(lst.join(""));
            });
        }
    }
})();
 
$(document).ready(function(){ indexInit.init() });
