$(function() {
    var csInterface = new CSInterface();  
	
    updateThemeWithAppSkinInfo(csInterface.hostEnvironment.appSkinInfo);
	
    // Update the color of the panel when the theme color of the product changed.
    csInterface.addEventListener(CSInterface.THEME_COLOR_CHANGED_EVENT, onAppThemeColorChanged);
	
	$('#sendBtn').click(function(event) {		
		event.preventDefault();
		
		// check that a document is open (using extendscript)..
		evalScript('app.documents.length', function(result) {
			if (result == 0) {
				alert('An open document is required');
			} else {
				// call the 'sendScriptMessage(..) function defined in main.jsx
				evalScript('sendMessage("'+$('#txt').val()+'")', function(result) {
					if ($('#txtReceived').val() != '') {
						$('#txtReceived').val($('#txtReceived').val() + '\n==============\n'); 
					}
					$('#txtReceived').val($('#txtReceived').val() + result)
				});
			}
		});
	});
	
	$('#clearBtn').click(function(event) {		
		event.preventDefault();
		$('#txtReceived').val('');
		$('#txt').val('');
	});
	
	// let the native plug-in part of this sample know that we are ready to receive events now..	
	var panelReadyEvent = new CSEvent('com.adobe.aisdk.sm.panelready', 'APPLICATION', 'ILST', 'ScriptMessageUI');
	csInterface.dispatchEvent(panelReadyEvent);
});

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
           } else if( stylesheet.insertRule ) {
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