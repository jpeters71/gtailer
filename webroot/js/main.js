define(
	"main",
	[
		"MessageList"
	],
	function(MessageList) {
		var ws = new WebSocket("wss://10.54.25.251:7443/entry");
		var list = new MessageList(ws);
		ko.applyBindings(list);
	}
);
