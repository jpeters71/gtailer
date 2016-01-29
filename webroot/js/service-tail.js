
var gtailerInit = (function() {
    var ws;
    
    var startListening = function() {
        ws = new WebSocket("wss://localhost:7443/entry", "protocolOne");
        
        ws.onmessage = function (event) {
            var evj = $.parseJSON(event.data);
            $("#logMessages").append("<div class='rTableRow'>" +
                "<div class='rTableCell-l'>" + evj.name + "</div>" +
                "<div class='rTableCell-r'>" + $("<div/>").text(evj.text).html() + "</div>" +
                "</div>"
            );
            if (!$("#holdBtn").is(":checked")) {
                var elemMsg = $("#logMessagesPar");
                var scrollHeight = elemMsg[0].scrollHeight;
                elemMsg.scrollTop(scrollHeight);
            }
        }
    }
    
    var subscribe = function(strHost) {
        var strCleanedHost = strHost.substring("st-cb-".length);
        $.ajax({
            url: "/subscribe/" + strCleanedHost + "/service-tail"
        });
    }
    
    var showList = function() {
        $.getJSON("/hosts", function(data) {
            var lst = [];
            $.each(data, function(key, val) {
                lst.push("<input type='checkbox' id='st-cb-" + val.name + "'>" + val.name + "</input><br/>");
            });
            $("#host-list").html(lst.join(""));
            $(":checkbox").change(function() {
                if (this.id.indexOf("st-cb-") === 0) {
                    if(this.checked) {
                        if (!ws) {
                            startListening();
                        }
                        subscribe(this.id);
                    }
                    else {
                        alert(this.id + " UNCHECKED");
                    }
                }
            });

        });
    }
    
    return {
        init: function() {
            showList();
        }
    }
})();
 
$(document).ready(function(){ 
    gtailerInit.init() 

});
