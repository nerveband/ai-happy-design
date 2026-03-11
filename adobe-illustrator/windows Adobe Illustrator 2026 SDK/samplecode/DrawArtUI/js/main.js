(function () {
    'use strict';
    
    // Event we dispatch (to send to native Illustrator plug-in)
    var BUTTON_PRESSED = "com.adobe.aisdk.sm.html.buttonpressed",

    // Event we listen for (from native Illustrator plug-in)
    BITMAP_DRAWN = "com.adobe.aisdk.sm.html.bitmapdrawn";

    var csInterface = new CSInterface();
    
    // Reloads extension (e.g. reset the extension)
    function reloadPanel() {
        location.reload();
    }
    
    function decodeImage(event) {
        console.log("Event \"" + event.type + " \" received.");

        // Use $('#img')[0] to get the native DOM element rather than jQuery extended object
        // to obtain canvas context.
        var context = $('#img')[0].getContext("2d");
        
        // Parses <data> object sent from the plugin
        var xmlData = $.parseXML(event.data);
        var $xml = $(xmlData);
        var bitmapData = $xml.find("data").text();
       
        // Decodes base64 bitmap data
        var byteArray =  atob(bitmapData);

        // Creates an ImageData object
        var imgData = context.createImageData(50, 50);
        
        // Populates ImageData object with bitmap data
        for (var i = 0; i < byteArray.length; i++) {
            imgData.data[i] = byteArray.charCodeAt(i);
        }       
        context.putImageData(imgData, 0, 0);
	}
    
    function buttonPressed() {
        // If no document opened, open one. See /jsx/hostscript.jsx.
        csInterface.evalScript("checkDocumentIsOpen()");
        
        // Create CSEvent and dispatch it to the plugin's event listener.
        var panelReadyEvent = new CSEvent(BUTTON_PRESSED, "APPLICATION", "ILST", "DrawArtHtml");
        csInterface.dispatchEvent(panelReadyEvent);
        console.log("Event \"" + BUTTON_PRESSED + " \" dispatched.");
    }
    
    function appComplete() {
        // Adds a listener for the BITMAP_DRAWN event.
        csInterface.addEventListener(BITMAP_DRAWN, decodeImage);
    }
    
    function init() {   
        themeManager.init();
        $("#btn_reload").click(reloadPanel);
        $("#AGMPortButton").click(buttonPressed);
    }
        
    init();
    appComplete();
}());