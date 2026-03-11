var markedObjects = new Array();
var selectedObjectId = null;
var selectedIndex = null;

$(function() {
    var csInterface = new CSInterface();
	
    var appName = csInterface.hostEnvironment.appName;
	
    updateThemeWithAppSkinInfo(csInterface.hostEnvironment.appSkinInfo);
	
    // Update the color of the panel when the theme color of the product changed.
    csInterface.addEventListener(CSInterface.THEME_COLOR_CHANGED_EVENT, onAppThemeColorChanged);
	
	csInterface.addEventListener("com.adobe.csxs.events.MarkedObjectsUpdatePanel", updatePanel);
	csInterface.addEventListener("com.adobe.csxs.events.MarkedObjectsUpdateBoundingBoxInfo", updateBoundingBoxInfo);
	csInterface.addEventListener("com.adobe.csxs.events.MarkedObjectsUpdatePanelFontNameList", updatePanelFontNameList);
	csInterface.addEventListener("com.adobe.csxs.events.MarkedObjectsUpdatePreferencesPanel", updatePreferencesPanel);
	csInterface.addEventListener("com.adobe.csxs.events.MarkedObjectsAddMarkedObject", addMarkedObject);
	csInterface.addEventListener("com.adobe.csxs.events.MarkedObjectsClearMarkedObject", clearMarkedObject);
	
	$('#save_menu').on('click', function(e) {
		csInterface.dispatchEvent(new CSEvent('com.adobe.csxs.events.MarkedObjectsSave', 'APPLICATION', 'ILST', 'MarkedObjectsUI'));
	});
	
	$('#renumber_menu').on('click', function(e) {
		csInterface.dispatchEvent(new CSEvent('com.adobe.csxs.events.MarkedObjectsRenumberIDs', 'APPLICATION', 'ILST', 'MarkedObjectsUI'));
	});
	
	$('#preferences_menu').on('click', function(e) {
		$('#home').hide();
		$('#preferences').show();
		csInterface.dispatchEvent(new CSEvent('com.adobe.csxs.events.MarkedObjectsPreferencesStackReady', 'APPLICATION', 'ILST', 'MarkedObjectsUI'));
	});
	
	// make a jqueryui selectable list..
	$('#list').selectable({
		selected: function( event, ui ) {
			selectedObjectId = ui.selected.id;
			
			// find the index of the selected object in the markedObjects array..
			markedObjects.filter(function(markedObject, index) { 
				if (markedObject.id == selectedObjectId) {
					selectedIndex = index;
				} 
			});
	  
			var event = new CSEvent('com.adobe.csxs.events.MarkedObjectsItemSelect', 'APPLICATION', 'ILST', 'MarkedObjectsUI');
			event.data = '<payload>';
			event.data += '<item>' + selectedIndex + '</item>';
			event.data += '<label>' + $('#label').val() + '</label>';
			event.data += '<ID>' + $('#id').val() + '</ID>';
			event.data += '<locx>' + $('#homeX').val() + '</locx>';
			event.data += '<locy>' + $('#homeY').val() + '</locy>';
			event.data += '</payload>';
			csInterface.dispatchEvent(event);
		}
	});
	
	$('#preferencesForm').submit(function(event) {		
		event.preventDefault();
		
		var event = new CSEvent('com.adobe.csxs.events.MarkedObjectsPreferences', 'APPLICATION', 'ILST', 'MarkedObjectsUI');
		event.data = '<payload>';
		event.data += '<subfontsize>' + $('#fontSize').val() + '</subfontsize>';
		event.data += '<subfontname>' + $('#fontName').val() + '</subfontname>';
		event.data += '<sublocx>' + $('#x').val() + '</sublocx>';
		event.data += '<sublocy>' + $('#y').val() + '</sublocy>';
		event.data += '<defaulttext>' + $('#defaultText').val() + '</defaulttext>';
		event.data += '<precision>' + $('#pointPrecision').val() + '</precision>';
		event.data += '<autosort>' + $('#autoSort').prop('checked') + '</autosort>';
		event.data += '</payload>';
		csInterface.dispatchEvent(event);
		
		$('#home').show();
		$('#preferences').hide();
	});
	
	$('#updateBtn').click(function(event) {		
		event.preventDefault();
		
		var event = new CSEvent('com.adobe.csxs.events.MarkedObjectsUpdateClicked', 'APPLICATION', 'ILST', 'MarkedObjectsUI');
		event.data = '<payload>';
		event.data += '<item>' + selectedIndex + '</item>';
		event.data += '<label>' + $('#label').val() + '</label>';
		event.data += '<ID>' + $('#id').val() + '</ID>';
		event.data += '<locx>' + $('#homeX').val() + '</locx>';
		event.data += '<locy>' + $('#homeY').val() + '</locy>';
		event.data += '</payload>';
		csInterface.dispatchEvent(event);
	});
	
	$('#hideShowBtn').click(function(event) {		
		event.preventDefault();
		
		var event = new CSEvent('com.adobe.csxs.events.MarkedObjectsHideAllClicked', 'APPLICATION', 'ILST', 'MarkedObjectsUI');
				
		var buttonLabel = $(this).attr('value');
		if (buttonLabel == 'Hide All') {
			event.data = '<payload>';
			event.data += '<hide>true</hide>';
			event.data += '</payload>';
			$(this).attr('value', 'Show All');
			
		} else {
			event.data = '<payload>';
			event.data += '<hide>false</hide>';
			event.data += '</payload>';	
			$(this).attr('value', 'Hide All');
		}
		
		csInterface.dispatchEvent(event);
	});

	// let the native plug-in part of this sample know that we are ready to receive events now..	
	var panelReadyEvent = new CSEvent('com.adobe.csxs.events.MarkedObjectsPanelReady', 'APPLICATION', 'ILST', 'MarkedObjectsUI');
	csInterface.dispatchEvent(panelReadyEvent);
});

function updatePanel(event) {
    var xmlData = $.parseXML(event.data);
    var $xml = $(xmlData);
	
	if ($xml.find('item').text() != -1) {	
		var x = $xml.find('locx').text();
		var y = $xml.find('locy').text();
		var label = $xml.find('label').text();
		var id = $xml.find('ID').text();
		var itemIndex = $xml.find('item').text();
		
		$('#homeX').val(x);
		$('#homeY').val(y);
		$('#id').val(id);
		$('#label').val(label);
		$('input[name="submit"]').removeAttr('disabled');
	} 
}

function updateBoundingBoxInfo(event) {
    var xmlData = $.parseXML(event.data);
    var $xml = $(xmlData);
	
	var width = $xml.find('width').text();
	var height = $xml.find('height').text();
	
	$('#width').val(width);
	$('#height').val(height);
}

function updatePanelFontNameList(event) {
    var xmlData = $.parseXML(event.data);
    var $xml = $(xmlData);
	
	// parse fontlist out of event XML..
	var fonts = $xml.find('fontlist').text();
	var fontsArr = fonts.split(',');
	
	// empty the drop down menus..
	$('select[name="fontName"]').empty();
	$('select[name="defaultFontName"]').empty();
	
	// populate the drop down menus..
	var fontNameSelect = document.getElementById('fontName');
	var defaultFontNameSelect = document.getElementById('defaultFontName');
	for (var i = 0; i < fontsArr.length; i++) {
	    var opt = document.createElement('option');
	    opt.innerHTML = fontsArr[i];
	    opt.value = fontsArr[i];
		
	    fontNameSelect.appendChild(opt);
		// you can't use the same element twice, so it must be cloned to re-use it
		defaultFontNameSelect.appendChild(opt.cloneNode(true));
	}
}

function updatePreferencesPanel(event) {
    var xmlData = $.parseXML(event.data);
    var $xml = $(xmlData);
	
	var fontSize = $xml.find('subfontsize').text();
	var fontName = $xml.find('subfontname').text();
	var x = $xml.find('sublocx').text();
	var y = $xml.find('sublocy').text();
	
	var defaultFontSize = $xml.find('labelfontsize').text();
	var defaultFontName = $xml.find('labelfontname').text();
	var defaultX = $xml.find('labellocx').text();
	var defaultY = $xml.find('labellocy').text();
	var defaultText = $xml.find('defaulttext').text();
	
	var precision = $xml.find('precision').text();
	var autoSort = $xml.find('autosort').text();

	// set form field values..
	$('#fontSize').val(fontSize);
	$('#fontName option[value="' + fontName + '"]').prop('selected', true);
	$('#x').val(x);
	$('#y').val(y);
	
	$('#defaultFontSize').val(defaultFontSize);
	$('#defaultFontName option[value="' + defaultFontName + '"]').prop('selected', true);
	$('#defaultX').val(defaultX);
	$('#defaultY').val(defaultY);
	$('#defaultText').val(defaultText);
	
	$('#pointPrecision').val(precision);
	$('#autoSort').prop('checked', (autoSort == '1'));
}

function addMarkedObject(event) {
    var xmlData = $.parseXML(event.data);
    var $xml = $(xmlData);
	
	var id = $xml.find('id').text();
	var object = $xml.find('markedobject').text();
	
	// add marked object to array..
	markedObject = {id:id, object: object};
	markedObjects.push(markedObject);
	
	// add object to UI list..
	$('#list').append('<li id="'+id+'" class="ui-widget-content">'+object+'</li>');
	
	$('input[name="hide"]').removeAttr('disabled');
}

function clearMarkedObject(event) { 
	markedObjects = [];
	$('#list').empty();
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