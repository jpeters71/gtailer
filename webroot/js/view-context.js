
var viewContext = (function() {
    
    var showList = function() {
        $.getJSON("/hosts", function(data) {
            var lst = [];
            $.each(data, function(key, val) {
                var strId = "vc-" + val.name.replace(".", "_");
                lst.push("<a href='#' id='" + strId + "'>" + val.name + "</a><br/>");
            });
            $("#host-list").html(lst.join(""));
            $.each(data, function(key, val) {
                var strId = "vc-" + val.name.replace(".", "_");
                $("#" + strId).click(function() {
                    if (this.id.indexOf("vc-") === 0) {
                        showContext(val.name);
                    }
                });
            });
            
        });
    }
    
    var showContext = function(strHost) {
        var strUrl = "/subscribe/" + strHost + "/view-context";
        
        if ($("#openInNewTab").is(":checked")) {
            // Open in a new tab
            window.open(strUrl, '_blank')
        }
        else {
            $.get(strUrl, function(data) {
                var strHtml = "<h3>" + strHost + "</h3>";
                
                strHtml += "<div class='codeView'>";
                var strEscaped = $("<div>").text(data).html();
                
                strHtml += "<code>" + strEscaped + "</code>";
                strHtml += "</div>";
                $("#contextView").html(strHtml);
            });
        }
    }
    
    return {
        init: function() {
            showList();
        }
    }
})();
 
$(document).ready(function(){ 
    viewContext.init() 

});
