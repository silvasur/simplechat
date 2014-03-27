function randomColor() {
	var h = Math.floor(Math.random() * 360);
	var s = Math.floor((Math.random() * 0.5 + 0.5)*100);
	var l = Math.floor((Math.random() * 0.4 + 0.3)*100);
	return "hsl("+h+", "+s+"%, "+l+"%)";
}

function askTryAgain(reason, ws_url) {
	if(confirm("Could not join chat (Reason: "+ reason +"). Try again?")) {
		window.setTimeout(RunChat, 1, ws_url); // We use a timeout, so we don't accidentally fill up the call stack.
	}
}

var buddyColors = {};

function addBuddy(nick) {
	var found = false;
	$("#buddies li").each(function(index) {
		if($(this).text() == nick) {
			found = true;
		}
	});
	if(!found) {
		var col = randomColor();
		buddyColors[nick] = col;
		console.log(col);
		$("#buddies").append($("<li/>").css("color", col).text(nick));
	}
}

function removeBuddy(nick) {
	$("#buddies li").each(function(index) {
		var self = $(this);
		if(self.text() == nick) {
			self.remove();
		}
	});
}

function chatlogWriter(event) {
	var data = JSON.parse(event.data);
	var msgtext = "";
	switch(data.type) {
	case "chat":
		msgtext = data.text;
		break;
	case "join":
		msgtext = "joined the room";
		addBuddy(data.user);
		break;
	case "leave":
		msgtext = "left the room";
		removeBuddy(data.user);
		break;
	}
	
	var elemNick = $("<span/>").addClass("nick").prop("style", "color: " + buddyColors[data.user]).text(data.user);
	var elemText = $("<span/>").addClass("msg").text(msgtext);
	var logentry = $("<li/>").addClass(data.type).append(elemNick).append(elemText);
	$("#chatlog").append(logentry);
	window.scrollTo(0, logentry.offset().top);
}

function initChatSender(ws) {
	var send = function() {
		var ct = $("#chattext")
		ws.send(ct.prop("value"));
		ct.prop("value","");
		ct.focus();
	};
	
	$("#sendbtn").click(send);
	$("#chattext").keyup(function(event) {
		if(event.keyCode==13) {
			send();
		}
	});
}

function Join(ws_url, nick) {
	var ws = new WebSocket(ws_url, "chat");
	ws.onopen = function(_) {
		ws.send(nick); 
		ws.onmessage = function(event) {
			var data = JSON.parse(event.data);
			if(data.ok) {
				ws.onmessage = chatlogWriter;
				for(i in data.buddies) {
					addBuddy(data.buddies[i]);
				}
				initChatSender(ws);
				ws.onclose = function(_) {
					alert("Connection lost. Try refreshing the page.");
				};
			} else {
				ws.close();
				askTryAgain(data.error, ws_url);
			}
		};
	};
}

function RunChat(ws_url) {
	var nick = "";
	while(nick == "") {
		nick = prompt("Choose a nickname");
		if(nick === null) {
			return;
		}
	}
	Join(ws_url, nick);
}
