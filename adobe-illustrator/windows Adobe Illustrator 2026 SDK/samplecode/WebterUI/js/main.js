var docEntries = [];
var artEntries = [];

$(function() {
    var csInterface = new CSInterface();  
	
    updateThemeWithAppSkinInfo(csInterface.hostEnvironment.appSkinInfo);
	
    // Update the color of the panel when the theme color of the product changed.
    csInterface.addEventListener(CSInterface.THEME_COLOR_CHANGED_EVENT, onAppThemeColorChanged);
	
	// add all the event listeners we need..
	csInterface.addEventListener('com.adobe.csxs.events.WebterUpdateEntry', handleDictionary);
	csInterface.addEventListener('com.adobe.csxs.events.WebterUpdateDefinitionAndType', handleItem);
	
	$('#docAddBtn').click(function(event) {		
		event.preventDefault();
		
		var event = new CSEvent('com.adobe.csxs.events.WebterAddEntry', 'APPLICATION', 'ILST', 'WebterUI');
		event.data = '<payload>';
		event.data += '<entry>0</entry>';
		event.data += '<entrystring>' + $('#docEntryDropDown').val() + '</entrystring>';
		event.data += '<value>' + $('#docDefinitionTxt').val() + '</value>';
		event.data += '<entrytype>' + $('#docTypeDropDown option:selected').val() + '</entrytype>';
		event.data += '</payload>';
		csInterface.dispatchEvent(event);
	});
	
	$('#artAddBtn').click(function(event) {		
		event.preventDefault();
		var event = new CSEvent('com.adobe.csxs.events.WebterAddEntry', 'APPLICATION', 'ILST', 'WebterUI');
		event.data = '<payload>';
		event.data += '<entry>1</entry>';
		event.data += '<entrystring>' + $('#artEntryDropDown').val() + '</entrystring>';
		event.data += '<value>' + $('#artDefinitionTxt').val() + '</value>';
		event.data += '<entrytype>' + $('#artTypeDropDown option:selected').val() + '</entrytype>';
		event.data += '</payload>';
		csInterface.dispatchEvent(event);
	});
	
	// let the native plug-in part of this sample know that we are ready to receive events now..	
	var panelReadyEvent = new CSEvent('com.adobe.csxs.events.WebterPanelReady', 'APPLICATION', 'ILST', 'WebterUI');
	csInterface.dispatchEvent(panelReadyEvent);
});

function handleDictionary(event) {
    var xmlData = $.parseXML(event.data);
    var $xml = $(xmlData);
		
	var isDoc = ($xml.find('dict').text() == 'doc') ? true : false;
	var isNull = ($xml.find('null').text() == 'true') ? true : false;
	
	// set selected art label..
	var artName = (!isNull && !isDoc) ? $xml.find('artname').text() : 'No art selected';
	$('#artLabel').text(artName);
	
	if (isDoc) {		
		if (!isNull) {			
			$('#docEntryDropDown').empty();
			docEntries = [];
			
			// set our data model..
			docEntries = parseEntries($xml.find('entries'));
			
			// add all the entries to dropdown..
			var docEntryDropDown = document.getElementById('docEntryDropDown');
			for (var e in docEntries) {
				var entry = docEntries[e];

				// add to dropdown
			    var opt = document.createElement('option');
			    opt.innerHTML = entry;
			    opt.value = e;
			    docEntryDropDown.appendChild(opt);
			}
			
			// enable controls
			$('#docEntryDropDown').removeAttr('disabled');
			$('#docDefinitionTxt').removeAttr('disabled');
			$('#docTypeDropDown').removeAttr('disabled');
			$('#docAddBtn').removeAttr('disabled');
			
		} else {
			// disable controls
			$('#docEntryDropDown').attr("disabled", "disabled");
			$('#docDefinitionTxt').attr("disabled", "disabled");
			$('#docTypeDropDown').attr("disabled", "disabled");
			$('#docAddBtn').attr("disabled", "disabled");
		}
	} else {		
		if (!isNull) {
			$('#artEntryDropDown').empty();
			artEntries = [];
			
			// set our data model..
			artEntries = parseEntries($xml.find('entries'));
			
			// add all the entries to dropdown..
			var artEntryDropDown = document.getElementById('artEntryDropDown');
			for (var e in artEntries) {
				var entry = artEntries[e];

				// add to dropdown
			    var opt = document.createElement('option');
			    opt.innerHTML = entry;
			    opt.value = e;
			    artEntryDropDown.appendChild(opt);
			}
			
			// enable controls
			$('#artEntryDropDown').removeAttr('disabled');
			$('#artDefinitionTxt').removeAttr('disabled');
			$('#artTypeDropDown').removeAttr('disabled');
			$('#artAddBtn').removeAttr('disabled');
			
		} else {
			// disable controls
			$('#artEntryDropDown').attr("disabled", "disabled");
			$('#artDefinitionTxt').attr("disabled", "disabled");
			$('#artTypeDropDown').attr("disabled", "disabled");
			$('#artAddBtn').attr("disabled", "disabled");
		}
	}
}

// utility function for parsing the received entries string into 
// a more usable associative array of id => label
function parseEntries($entriesXml) {
	var entriesToReturn = [];
	
	var entries = $entriesXml.text().split(',');
	for (var e in entries) {
		var entry = entries[e];
		if (entry.length != 0) {
			var splitEntry = entry.split('#id#=');
			if (splitEntry.length == 2) {
				var label = splitEntry[0];
				var id = splitEntry[1];
				entriesToReturn[id] = label;
			}
		}
	}
	
	return entriesToReturn;
}

function handleItem(event) {
    var xmlData = $.parseXML(event.data);
    var $xml = $(xmlData);
	
	var dict = $xml.find('dict').text();
	var key = $xml.find('key').text();
	var value = $xml.find('value').text();
	var type = $xml.find('type').text();
	
	if (dict == 'doc') {
		$("#docTypeDropDown option[value='"+type+"']").attr("selected", "selected");
		$('#docDefinitionTxt').val(value);
		
	} else {
		$("#artTypeDropDown option[value='"+type+"']").attr("selected", "selected");
		$('#artDefinitionTxt').val(value);
	}	
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