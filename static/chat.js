function initGui(roomname) {
	$("#nojs").remove();
	
	var dialogoverlay = $("<div/>").prop("id", "dialogoverlay")
		.append($("<div/>").prop("id", "dialoginner")
			.append($("<h2/>"))
			.append($("<div/>")
				.append($("<p/>"))
				.append($("<input/>"))
				.append($("<div/>").addClass("buttons")))).hide();
	var h1 = $("<h1/>").text("Chatroom ").append($("<span/>").css("font-style", "italic").text(roomname));
	var chatlog = $("<div/>").prop("id", "chatlogwrap").append($("<ul/>").prop("id", "chatlog"));
	var buddiescontainer = $("<section/>").prop("id", "buddiescontainer")
		.append($("<h2/>").text("Buddies online"))
		.append($("<ul/>").prop("id", "buddies"));
	var chatinput = $("<section/>").prop("id", "chatinput")
		.append($("<input/>").prop("name", "chattext").prop("id", "chattext").prop("type", "text").prop("placeholder", "Type to chat..."))
		.append($("<button/>").prop("id", "sendbtn").text("Send"));
	$("body").append(dialogoverlay);
	$("#mainwrap")
		.append(h1)
		.append(chatlog)
		.append(buddiescontainer)
		.append(chatinput);
}

function mydialog(title, text, input, buttons) {
	$("#dialoginner h2").text(title);
	$("#dialoginner p").text(text);
	var callcallback;
	if(input === null) {
		$("#dialoginner input").hide();
		callcallback = function(cb) {cb();};
	} else {
		$("#dialoginner input").prop("placeholder", input).show();
		callcallback = function(cb) {cb($("#dialoginner input").val());};
	}
	
	$("#dialoginner .buttons > *").remove();
	var btncontainer = $("#dialoginner .buttons");
	
	for(var i in buttons) {
		var button = buttons[i];
		btncontainer.append($("<button/>").text(button.text).click(function() {
			$("#dialogoverlay").hide();
			callcallback(button.callback);
		}));
	}
	
	$("#dialogoverlay").show();
}

function randomColor() {
	var h = Math.floor(Math.random() * 360);
	var s = Math.floor((Math.random() * 0.5 + 0.5)*100);
	var l = Math.floor((Math.random() * 0.4 + 0.3)*100);
	return "hsl("+h+", "+s+"%, "+l+"%)";
}

function askTryAgain(reason, ws_url) {
	mydialog(
		"Could not join",
		"Could not Join chat (Reason: "+reason+")",
		null,
		[{
			"text": "Try again",
			"callback": (function() {
				pickUsername(ws_url);
			})
		}]
	);
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
	
	var elemNick = $("<span/>").addClass("nick").css("color", buddyColors[data.user]).text(data.user);
	var elemText = $("<span/>").addClass("msg").text(msgtext);
	var logentry = $("<li/>").addClass(data.type).append(elemNick).append(elemText);
	$("#chatlog").append(logentry);
	logentry.get()[0].scrollIntoView();
}

function initChatSender(ws) {
	var send = function() {
		var ct = $("#chattext")
		ws.send(ct.val());
		ct.val("");
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
				$("#chatinput").show();
				ws.onmessage = chatlogWriter;
				for(var i in data.buddies) {
					addBuddy(data.buddies[i]);
				}
				initChatSender(ws);
				ws.onclose = function(_) {
					mydialog("Connection lost", "Connection to server lost. Try refreshing the page", null, []);
				};
			} else {
				ws.close();
				askTryAgain(data.error, ws_url);
			}
		};
	};
}

function RunChat(ws_url, roomname) {
	initGui(roomname);
	pickUsername(ws_url);
}

function pickUsername(ws_url) {
	mydialog(
		"Pick a username",
		"Pick a username to join the chat",
		"username",
		[{
			"text": "OK",
			"callback": (function(nick) {Join(ws_url, nick);})
		}]
	);
}
