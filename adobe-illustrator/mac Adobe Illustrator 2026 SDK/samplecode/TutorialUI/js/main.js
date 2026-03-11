
$(function() {
    var csInterface = new CSInterface();
	
    updateThemeWithAppSkinInfo(csInterface.hostEnvironment.appSkinInfo);
	
    // Update the color of the panel when the theme color of the product changed.
    csInterface.addEventListener(CSInterface.THEME_COLOR_CHANGED_EVENT, onAppThemeColorChanged);
	
	csInterface.addEventListener('com.adobe.aisdk.TutorialGetParamsReq', respondToParams);
	
	// add an event listener for user hitting keyboard return in the dialog..
	document.onkeypress = function(e) {
		if (e.keyCode == 13) {
			dispatchParams(false);	
		}
	};
	
	$('#okBtn').click(function(event) {		
		event.preventDefault();
		dispatchParams(false);
	});
	
	$('#cancelBtn').click(function(event) {		
		event.preventDefault();
		dispatchParams(true);
	});
	
	// let the native plug-in part of this sample know that we are ready to receive events now..	
	var panelReadyEvent = new CSEvent('com.adobe.aisdk.TutorialDialogReady', 'APPLICATION', 'ILST', 'TutorialUI');
	csInterface.dispatchEvent(panelReadyEvent);
});

function respondToParams (event) {	
    var xmlData = $.parseXML(event.data);
    var $xml = $(xmlData);
	var $defaults = $xml.find('default');
	var corners = $defaults.find('corners').text();
	var size = $defaults.find('size').text();
	
	$('#corners').val(corners);
	$('#size').val(size);
}

function dispatchParams (cancelled) {								
	var event = new CSEvent('com.adobe.aisdk.TutorialGetParamsResp', 'APPLICATION', 'ILST', 'TutorialUI');
	event.data = '<payload>';
	event.data +=	'<cancelled>' + cancelled + '</cancelled>';
	event.data +=	'<corners>' + $('#corners').val() + '</corners>';
	event.data +=	'<size>' + $('#size').val() + '</size>';
	event.data +=	'</payload>';
	new CSInterface().dispatchEvent(event);	
	new CSInterface().closeExtension();
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