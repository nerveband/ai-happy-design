function sendMessage(textToSend) {
	var resp = '';
	
	// if any art is selected, send the message on it. 
	// Otherwise send on the app object
	if (app.selection && app.selection.length != 0) {
		for (var i = 0; i < app.selection.length; i++) {
			resp = app.selection[i].sendScriptMessage('ScriptMessage', 'Send Text', textToSend);
		}
	} else {
		resp = app.sendScriptMessage('ScriptMessage', 'Send Text', textToSend);
	}
	
	return resp;
}