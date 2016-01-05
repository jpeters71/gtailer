define(
	"Message",
	[],
	function() {

		function Message(model) {
			if (model !== undefined) {
				this.host = ko.observable(model.host);
				this.name = ko.observable(model.name);
				this.text = ko.observable(model.text);
                this.printout = ko.observable(this.name() + "::" + this.text());
			} else {
				this.host = ko.observable("Anonymous");
				this.name = ko.observable("Anonymous");
				this.text = ko.observable("");
                this.printout = ko.observable(this.name + "::" + this.text);
			}

			this.toModel = function() {
				return {
					host: this.host(),
					name: this.name(),
					text: this.text(),
                    printout: this.name() + "::" + this.text()
				}; 
			}
		}

		return Message;
	}
);
