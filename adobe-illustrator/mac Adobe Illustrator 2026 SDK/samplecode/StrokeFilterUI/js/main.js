$(function() {
    var csInterface = new CSInterface();

    updateThemeWithAppSkinInfo(csInterface.hostEnvironment.appSkinInfo);
	
    // Update the color of the panel when the theme color of the product changed.
    csInterface.addEventListener(CSInterface.THEME_COLOR_CHANGED_EVENT, onAppThemeColorChanged);
	
	// add all the event listeners we need..
	csInterface.addEventListener('com.adobe.csxs.events.StrokeUpdateDashPanel', updateDashPanel);
	csInterface.addEventListener('com.adobe.csxs.events.StrokeUpdateWavePanel', updateWavePanel);
	
	$('#dashOkBtn').click(function(event) {		
		event.preventDefault();

		var event = new CSEvent('com.adobe.csxs.events.StrokeOkClicked', 'APPLICATION', 'ILST', 'StrokeFilterUI');
		event.data = '<payload>';
		event.data += '<state>DashPanel</state>';
		event.data += '<dashed>' + $('#dashedStrokeCB').is(':checked') + '</dashed>'; // value is true or false
		event.data += '<adjust>' + $('#adjustDashCB').is(':checked') + '</adjust>'; // value is true or false
		event.data += '<dash1>' + $('#dash1').val() + '</dash1>';
		event.data += '<gap1>' + $('#gap1').val() + '</gap1>';
		event.data += '<dash2>' + $('#dash2').val() + '</dash2>';
		event.data += '<gap2>' + $('#gap2').val() + '</gap2>';
		event.data += '<dash3>' + $('#dash3').val() + '</dash3>';
		event.data += '<gap3>' + $('#gap3').val() + '</gap3>';
		event.data += '</payload>';
		csInterface.dispatchEvent(event);
		
		csInterface.closeExtension();
	});
	
	$('#waveOkBtn').click(function(event) {		
		event.preventDefault();
		
		var event = new CSEvent('com.adobe.csxs.events.StrokeOkClicked', 'APPLICATION', 'ILST', 'StrokeFilterUI');
		event.data = '<payload>';
		event.data += '<state>WavePanel</state>';
		event.data += '<exchange>' + $('#exchangeCB').is(':checked') + '</exchange>'; // value is true or false
		event.data += '<replace>' + $('#replaceCB').is(':checked') + '</replace>'; // value is true or false
		event.data += '<number>' + $('#number').val() + '</number>';
		event.data += '<width1>' + $('#width1').val() + '</width1>';
		event.data += '<width2>' + $('#width2').val() + '</width2>';
		event.data += '</payload>';
		csInterface.dispatchEvent(event);
		
		csInterface.closeExtension();
	});
	
	$('#dashCancelBtn, #waveCancelBtn').click(function(event) {		
		event.preventDefault();
		csInterface.dispatchEvent(new CSEvent('com.adobe.csxs.events.StrokeCancelClicked', 'APPLICATION', 'ILST', 'StrokeFilterUI'));
		csInterface.closeExtension();
	});
		
	// let the native plug-in part of this sample know that we are ready to receive events now..	
	var panelReadyEvent = new CSEvent('com.adobe.csxs.events.StrokePanelReady', 'APPLICATION', 'ILST', 'StrokeFilterUI');
	csInterface.dispatchEvent(panelReadyEvent);
});

function updateDashPanel(event) {
    var xmlData = $.parseXML(event.data);
    var $xml = $(xmlData);
	
	var dashedStroke = ($xml.find('dashed').text() == 'true') ? true : false;
	var adjustDash = ($xml.find('adjust').text() == 'true') ? true : false;
	var dash1 = $xml.find('dash1').text();
	var dash2 = $xml.find('dash2').text();
	var dash3 = $xml.find('dash3').text();
	var gap1 = $xml.find('gap1').text();
	var gap2 = $xml.find('gap2').text();
	var gap3 = $xml.find('gap3').text();
	
	$('#waveForm').hide();
	$('#dashForm').show();
	
	if (dashedStroke) {
		$('#dashedStrokeCB').prop('checked', true);
	}
	
	if (adjustDash) {
		$('#adjustDashCB').prop('checked', true);
	}
	
	$('#dash1').val(dash1);
	$('#dash2').val(dash2);
	$('#dash3').val(dash3);
	$('#gap1').val(gap1);
	$('#gap2').val(gap2);
	$('#gap3').val(gap3);
}

function updateWavePanel(event) {
    var xmlData = $.parseXML(event.data);
    var $xml = $(xmlData);
	
	var exchangeWidth = ($xml.find('exchange').text() == 'true') ? true : false;
	var replace = ($xml.find('replace').text() == 'true') ? true : false;
	var number = $xml.find('number').text();
	var width1 = $xml.find('width1').text();
	var width2 = $xml.find('width2').text();
	
	$('#dashForm').hide();
	$('#waveForm').show();
	
	if (exchangeWidth) {
		$('#exchangeCB').prop('checked', true);
	}
	
	if (replace) {
		$('#replaceCB').prop('checked', true);
	}
	
	$('#number').val(number);
	$('#width1').val(width1);
	$('#width2').val(width2);
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