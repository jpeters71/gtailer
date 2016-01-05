define(
	"main",
	[
		"MessageList"
	],
	function(MessageList) {
		var ws = new WebSocket("wss://bozo:bozopwd@localhost:7443/entry");
		var list = new MessageList(ws);
		ko.applyBindings(list);
	}
);
