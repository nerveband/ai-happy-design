  //========================================================================================
//  
//  $File: //ai/mainline/devtech/sdk/public/samplecode/DrawArt/Source/DrawArtPanelController.cpp $
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
#include "DrawArtPanelController.h"
#include "DrawArtPlugin.h"
#include "AppContext.hpp"

// Our extension's ID
#define DrawArt_UI_EXTENSION "DrawArt"

// Events we listen for
static const char* BUTTON_PRESSED  = "com.adobe.aisdk.sm.html.buttonpressed";

// Events we dispatch
static const char* EVENT_THEME_CHANGE = "com.adobe.aisdk.sm.html.themechange";

static const char* BITMAP_DRAWN = "com.adobe.aisdk.sm.html.bitmapdrawn";

// Illustrator's CSXS app id
static const char* ILST_APP = "ILST";


static void ButtonPressedFunc(const csxs::event::Event* const eventParam, void* const context)
{
	AppContext appContext(gPlugin->GetPluginRef());
	DrawArtPanelController* panelController = (DrawArtPanelController*) context;
	std::string bitmap;
	gPlugin->DrawAGMPort(bitmap);
	panelController->SendBitmapToPanel(bitmap);
}

DrawArtPanelController::DrawArtPanelController()
:HtmlUIController(DrawArt_UI_EXTENSION)
{}

csxs::event::EventErrorCode DrawArtPanelController::RegisterCSXSEventListeners()
{
	csxs::event::EventErrorCode result = csxs::event::kEventErrorCode_Success;
	do
	{
		result = htmlPPLib.AddEventListener(BUTTON_PRESSED, ButtonPressedFunc, this);
		if (result != csxs::event::kEventErrorCode_Success) 
		{
			break;
		}
	}
	while (false);
	
	return result;
}

csxs::event::EventErrorCode DrawArtPanelController::RemoveEventListeners()
{
	csxs::event::EventErrorCode result = csxs::event::kEventErrorCode_Success;
	do
	{
		result = htmlPPLib.RemoveEventListener(BUTTON_PRESSED, ButtonPressedFunc, this);
		if (result != csxs::event::kEventErrorCode_Success) 
		{
			break;
		}
	}
	while (false);
	
	return result;
}

ASErr DrawArtPanelController::SendBitmapToPanel(std::string& bitmap)
{
	AIErr error = kNoErr;
	stringstream myString;
	myString << "<data>"
		<< bitmap
		<< "</data>";
	
	std::string xmlString = myString.str();
	csxs::event::Event event = {BITMAP_DRAWN,
		csxs::event::kEventScope_Application,
		ILST_APP,
		NULL,
		xmlString.c_str()};
	htmlPPLib.DispatchEvent(&event);

	return error;
}
