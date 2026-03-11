var snippets = [];

$(function() {
    var csInterface = new CSInterface();
	    
    updateThemeWithAppSkinInfo(csInterface.hostEnvironment.appSkinInfo);
	
    // Update the color of the panel when the theme color of the product changed.
    csInterface.addEventListener(CSInterface.THEME_COLOR_CHANGED_EVENT, onAppThemeColorChanged);
	
	// add all the event listeners we need..
	csInterface.addEventListener('com.adobe.csxs.events.SnippetRunnerUpdatePanel', parseTreeXML);
	csInterface.addEventListener('com.adobe.csxs.events.SnippetRunnerCanRun', canRun);
	csInterface.addEventListener('com.adobe.csxs.events.SnippetRunnerLog', getNewLogInfo);
	
	$('#run').click(function(event) {		
		event.preventDefault();
		
		var selectedOperationName = $('#operations option:selected').val();
		var selectedSnippetName = $('#snippets option:selected').val();
		
		var event = new CSEvent('com.adobe.csxs.events.SnippetRunnerRun', 'APPLICATION', 'ILST', 'SnippetRunnerUI');
		event.data = '<payload>';
		event.data += '<snippet>' + selectedSnippetName + '</snippet>';
		event.data += '<op>' + selectedOperationName + '</op>';
		event.data += '</payload>';
		csInterface.dispatchEvent(event);
	});
	
	$('#clear').click(function(event) {		
		event.preventDefault();
		$('#log').val('');
		csInterface.dispatchEvent(new CSEvent('com.adobe.csxs.events.SnippetRunnerClearLog', 'APPLICATION', 'ILST', 'SnippetRunnerUI'));
	});
	
	$('#save').click(function(event) {		
		event.preventDefault();
		csInterface.dispatchEvent(new CSEvent('com.adobe.csxs.events.SnippetRunnerSaveLog', 'APPLICATION', 'ILST', 'SnippetRunnerUI'));
	});
	
	$('#snippets').change(function() {
		// disable the run button whenever a snippet is selected
		$('#run').attr("disabled", "disabled");
		
		// enable the operations drop down..
		$('#operations').removeAttr('disabled');
		
		// blank the precondition being shown..
		$('#precondition').html('');
		
		var selectedSnippetName = $('#snippets option:selected').val();
		
		// update the operations dropdown to show the operations for the selected snippet..
		$('select[name="operations"]').empty();
		
		// add a prompt option..
		var operationSelect = document.getElementById('operations');
	    var promptOpt = document.createElement('option');
	    promptOpt.innerHTML = 'Select';
	    promptOpt.value = 'Select';
		operationSelect.appendChild(promptOpt);

		var snippetObj = snippets[selectedSnippetName];
		snippetObj.operations.forEach(function(op) {
			var operationName = op.name;
			
		    var opt = document.createElement('option');
		    opt.innerHTML = operationName;
		    opt.value = operationName;
		    operationSelect.appendChild(opt);
		});
	});
	
	$('#operations').change(function() {
		var selectedOperationName = $('#operations option:selected').val();
		var selectedSnippetName = $('#snippets option:selected').val();
		
		// set the precondition text
		var snippetObj = snippets[selectedSnippetName];
		snippetObj.operations.forEach(function(op) {
			if (op.name == selectedOperationName) {
				$('#precondition').html(op.precondition);
			}
		});
		
		// dispatch an event asking the native plug-in whether this operation can be run..
		// this will callback the 'canRun' function which enables/disables the run button.
		var event = new CSEvent('com.adobe.csxs.events.SnippetRunnerSelectionChange', 'APPLICATION', 'ILST', 'SnippetRunnerUI');
		event.data = '<payload>';
		event.data += '<snippet>' + selectedSnippetName + '</snippet>';
		event.data += '<op>' + selectedOperationName + '</op>';
		event.data += '</payload>';
		csInterface.dispatchEvent(event);
	});
	
	// let the native plug-in part of this sample know that we are ready to receive events now..	
	var panelReadyEvent = new CSEvent('com.adobe.csxs.events.SnippetRunnerPanelReady', 'APPLICATION', 'ILST', 'SnippetRunnerUI');
	csInterface.dispatchEvent(panelReadyEvent);
});

function parseTreeXML(event) {
    var xmlData = $.parseXML(event.data);
    var $xml = $(xmlData);
	
	// empty the drop down menu..
	$('#select[name="snippets"]').empty();
	$('#select[name="operations"]').empty();
	
	var snippetSelect = document.getElementById('snippets');
	
	// add a prompt option to drop downs
    var promptOpt = document.createElement('option');
    promptOpt.innerHTML = 'Select';
    promptOpt.value = 'Select';
	snippetSelect.appendChild(promptOpt);
	
	// populate the drop down menu and our snippets array..
	$xml.find('snippet').each(function() {
		var snippetName = $(this).attr('label');
		var description = $(this).find('description').attr('label');
		var $operations = $(this).find('operation');
		
	    var opt = document.createElement('option');
	    opt.innerHTML = snippetName;
	    opt.value = snippetName;
		
		// add to snippets drop down menu..
	    snippetSelect.appendChild(opt);
		
		// store a snippet object in our array of snippets..
		var snippet = {	 
			name: snippetName, 
			description: description,
			operations:[]
		};
		
		// add the snippets operations to the snippet object..
		$(this).find('operation').each(function() {
			var operationName = $(this).attr('label');
			var operation = { name: operationName };
			
			// add pre/post conditions to operation object..
			$(this).find('precondition').each(function() {
				operation.precondition = $(this).attr('label');
			});
			
			snippet.operations.push(operation);
		});
		
		snippets[snippetName] = snippet;
	});
}

function canRun(event) {
    var xmlData = $.parseXML(event.data);
    var $xml = $(xmlData);
	
	var snippetName = $xml.find('snippet').text();
	var operationName = $xml.find('op').text();
	var canRun = $xml.find('canrun').text();
	
	// if the currently selected snippet and operation match
	// the one this received data is about, enable/disable the
	// run button accordingly..
	if (canRun == 'false') {
		$('#run').attr("disabled", "disabled");
	} else {
		$('#run').removeAttr("disabled"); 
	}
}

function getNewLogInfo(event) {
    var xmlData = $.parseXML(event.data);
    var $xml = $(xmlData);
	
	// set log
	var logTxt = '';
	$xml.find('logline').each(function() {
		logTxt += $(this).text();
		logTxt += '\n';
		$('#log').val(logTxt);
	});
}

/**
 * Update the theme with the AppSkinInfo retrieved from the host product.
 */
function updateThemeWithAppSkinInfo(appSkinInfo) {
    //Update the background color of the panel
    var panelBackgroundColor = appSkinInfo.panelBackgroundColor.color;
    document.body.bgColor = toHex(panelBackgroundColor);
        
    var styleId = "ppstyle";
	addRule(styleId, ".default", "font-size:" + appSkinInfo.baseFontSize + "px" + "; color:" + reverseColor(panelBackgroundColor) + "; background-color:" + toHex(panelBackgroundColor, 20));
    addRule(styleId, "button", "border-color: " + toHex(panelBackgroundColor, -50));
}

function addRule(stylesheetId, selector, rule) {
    var stylesheet = document.getElementById(stylesheetId);
    
    if (stylesheet) {
        stylesheet = stylesheet.sheet;
           if( stylesheet.addRule ){
               stylesheet.addRule(selector, rule);
           } else if( stylesheet.insertRule ){
               stylesheet.insertRule(selector + ' { ' + rule + ' }', stylesheet.cssRules.length);
           }
    }
}

function reverseColor(color, delta) {
    return toHex({red:Math.abs(255-color.red), green:Math.abs(255-color.green), blue:Math.abs(255-color.blue)}, delta);
}

/**
 * Convert the Color object to string in hexadecimal format;
 */
function toHex(color, delta) {
    function computeValue(value, delta) {
        var computedValue = !isNaN(delta) ? value + delta : value;
        if (computedValue < 0) {
            computedValue = 0;
        } else if (computedValue > 255) {
            computedValue = 255;
        }

        computedValue = computedValue.toString(16);
        return computedValue.length == 1 ? "0" + computedValue : computedValue;
    }

    var hex = "";
    if (color) {
        with (color) {
             hex = computeValue(red, delta) + computeValue(green, delta) + computeValue(blue, delta);
        };
    }
    return "#" + hex;
}

function onAppThemeColorChanged(event) {
    // Should get a latest HostEnvironment object from application.
    var skinInfo = JSON.parse(window.__adobe_cep__.getHostEnvironment()).appSkinInfo;
    // Gets the style information such as color info from the skinInfo, 
    // and redraw all UI controls of your extension according to the style info.
    updateThemeWithAppSkinInfo(skinInfo);
} 

function evalScript(script, callback) {
    new CSInterface().evalScript(script, callback);
}