
var viewContext = (function() {
    var op = "";
     
    var showList = function() {
        $.getJSON("/hosts?operation=" + op, function(data) {
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
        var strUrl = "/subscribe/" + strHost + "/" + op;
        
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
                $("#fileView").html(strHtml);
            });
        }
    }
    var getUrlVars = function() {
        var vars = [], hash;
        var hashes = window.location.href.slice(window.location.href.indexOf('?') + 1).split('&');
        for(var i = 0; i < hashes.length; i++)
        {
            hash = hashes[i].split('=');
            vars.push(hash[0]);
            vars[hash[0]] = hash[1];
        }
        return vars;
    }    
    
    return {
        init: function() {
            op = getUrlVars()["op"];
            showList();
        }
    }
    
})();
 
$(document).ready(function(){ 
    viewContext.init() 

});
