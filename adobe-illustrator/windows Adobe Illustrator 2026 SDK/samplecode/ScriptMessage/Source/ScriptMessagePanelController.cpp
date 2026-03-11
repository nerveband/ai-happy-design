//========================================================================================
//  
//  $File: //ai/mainline/devtech/sdk/public/samplecode/ScriptMessage/Source/ScriptMessagePanelController.cpp $
//
//  $Revision: #1 $
//
//  Copyright 2010 Adobe Systems Incorporated. All rights reserved.
//  
//  NOTICE:  Adobe permits you to use, modify, and distribute this file in accordance 
//  with the terms of the Adobe license agreement accompanying it.  If you have received
//  this file from a source other than Adobe, then your use, modification, or 
//  distribution of it requires the prior written permission of Adobe.
//  
//========================================================================================

#include "IllustratorSDK.h"
#include "ScriptMessagePanelController.h"

// Our extension's ID
#define SCRIPTMESSAGE_UI_EXTENSION "com.adobe.illustrator.scriptmessageui.panel"

// Events we listen for
static const char* EVENT_PANEL_READY = "com.adobe.aisdk.sm.panelready";

// Events we dispatch
static const char* EVENT_THEME_CHANGE = "com.adobe.aisdk.sm.themechange";

// Illustrator's CSXS app id
static const char* ILST_APP = "ILST";

static void PanelReadyFunc(const csxs::event::Event* const eventParam, void* const context)
{
	ScriptMessagePanelController* panelController = (ScriptMessagePanelController*) context;
	panelController->SendThemeToPanel();
}

ScriptMessagePanelController::ScriptMessagePanelController()
:FlashUIController(SCRIPTMESSAGE_UI_EXTENSION)
{}

csxs::event::EventErrorCode ScriptMessagePanelController::RegisterCSXSEventListeners()
{
	csxs::event::EventErrorCode result = csxs::event::kEventErrorCode_Success;
	do
	{
		result = fPPLib.AddEventListener(EVENT_PANEL_READY, PanelReadyFunc, this);
		if (result != csxs::event::kEventErrorCode_Success) 
		{
			break;
		}
	}
	while (false);
	
	return result;
}

csxs::event::EventErrorCode ScriptMessagePanelController::RemoveEventListeners()
{
	csxs::event::EventErrorCode result = csxs::event::kEventErrorCode_Success;
	do
	{
		result = fPPLib.RemoveEventListener(EVENT_PANEL_READY, PanelReadyFunc, this);
		if (result != csxs::event::kEventErrorCode_Success) 
		{
			break;
		}
	}
	while (false);
	
	return result;
}

ASErr ScriptMessagePanelController::SendThemeToPanel()
{
	AIErr error = kNoErr;
	
	AIUIThemeColor backgroundColor;			// Panel background coluor
	AIUIThemeColor textColor;				// Panel text colour
	AIUIThemeColor focusColor;				// Colour around controls with focus
	AIUIThemeColor inputBackgroundColor;	// Background colour of text boxes, etc.
	AIUIThemeColor inputTextColor;			// Text colour of text boxes, etc.
	
	// Get the colours
	error = sAIUITheme->GetUIThemeColor(kAIUIThemeSelectorPanel,kAIUIComponentColorBackground, backgroundColor);
	error = sAIUITheme->GetUIThemeColor(kAIUIThemeSelectorPanel,kAIUIComponentColorText, textColor);
	error = sAIUITheme->GetUIThemeColor(kAIUIThemeSelectorPanel,kAIUIComponentColorFocusRing, focusColor);
	error = sAIUITheme->GetUIThemeColor(kAIUIThemeSelectorPanel,kAIUIComponentColorEditText, inputTextColor);
	error = sAIUITheme->GetUIThemeColor(kAIUIThemeSelectorPanel,kAIUIComponentColorEditTextBackground, inputBackgroundColor);
	
	// Construct XML payload
	stringstream myString;
	myString << "<payload><theme>" 
		<< "<backgroundColor>" 
			<< AIRGBColorToXML(backgroundColor)
		<< "</backgroundColor>"
		<< "<textColor>" 
			<< AIRGBColorToXML(textColor)
		<< "</textColor>"
		<< "<focusColor>" 
			<< AIRGBColorToXML(focusColor)
		<< "</focusColor>" 
		<< "<inputTextColor>" 
			<< AIRGBColorToXML(inputTextColor)
		<< "</inputTextColor>" 
		<< "<inputBackgroundColor>" 
			<< AIRGBColorToXML(inputBackgroundColor)
		<< "</inputBackgroundColor>"
	<< "</theme></payload>";
	
	// Dispatch event
	std::string xmlString = myString.str();
	csxs::event::Event event = {EVENT_THEME_CHANGE,
		csxs::event::kEventScope_Application,
		ILST_APP,
		NULL,
		xmlString.c_str()};
	fPPLib.DispatchEvent(&event);
	return error;
}

std::string ScriptMessagePanelController::AIRGBColorToXML(AIUIThemeColor color)
{
	int colorDepth = 255;
	int red = (int) (color.red * colorDepth);
	int green = (int) (color.green * colorDepth);
	int blue = (int) (color.blue * colorDepth);
	stringstream myString;
	myString 
		<< "<color>"
			<< "<a>" << color.alpha << "</a>"
			<< "<r>" << red << "</r>" 
			<< "<g>" << green << "</g>" 
			<< "<b>" << blue << "</b>" 
		<< "</color>";
	return myString.str();
}
