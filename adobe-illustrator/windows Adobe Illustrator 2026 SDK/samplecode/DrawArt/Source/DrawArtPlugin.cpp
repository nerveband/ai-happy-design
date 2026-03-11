//========================================================================================
//  
//  $File: //ai/mainline/devtech/sdk/public/samplecode/DrawArtHtml/Source/DrawArtPlugin.cpp $
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
#include "AIDrawArt.h"
#include "DrawArtPlugin.h"
#include "DrawArtPanelController.h"
#include "AICSXS.h"

#include "SDKErrors.h"

// Definitions for Base64 encoding
#include <string>
#include <vector>

//#ifdef MAC_ENV
#define TCHAR wchar_t //Unicode
#define TCHAR char    //Not unicode
#define TEXT(x) x     //Not unicode
#define DWORD long
#define BYTE unsigned char
//#endif

DrawArtPlugin* gPlugin = NULL;

Plugin* AllocatePlugin(SPPluginRef pluginRef)
{
	return new DrawArtPlugin(pluginRef);
}

void FixupReload(Plugin* plugin)
{
	DrawArtPlugin::FixupVTable((DrawArtPlugin*) plugin);
}

DrawArtPlugin::DrawArtPlugin(SPPluginRef pluginRef) :
	Plugin(pluginRef),
	fPanelController(NULL),
	fPlugPlugSetupCompleteNotifier(NULL),
	fAboutPluginMenu(NULL)
{
	strncpy(fPluginName, kDrawArtPluginName, kMaxStringLength);
}

ASErr DrawArtPlugin::StartupPlugin(SPInterfaceMessage *message)
{
	ASErr error;
	
	error = Plugin::StartupPlugin(message);
	
	if (error)
		return error;

	// Add a menu item to the About SDK Plug-ins menu group.
	SDKAboutPluginsHelper aboutPluginsHelper;
	aboutPluginsHelper.AddAboutPluginsMenuItem(message, 
				kSDKDefAboutSDKCompanyPluginsGroupName, 
				ai::UnicodeString(kSDKDefAboutSDKCompanyPluginsGroupNameString), 
				"DrawArtHtml...", 
				&fAboutPluginMenu);
	
	// Create panel controller
	if (fPanelController == NULL)
	{
		fPanelController = new DrawArtPanelController();
		error = Plugin::LockPlugin(true);
		
		if (error)
			return error;
	}
		
	// Add notifier for adding CSXS event listeners
	error = sAINotifier->AddNotifier(fPluginRef, kDrawArtPluginName,
									 kAICSXSPlugPlugSetupCompleteNotifier, &fPlugPlugSetupCompleteNotifier);
	
	if (error)
		return error;
	
	return error;
}

/*
 */
ASErr DrawArtPlugin::SetGlobal(Plugin *plugin)
{
	gPlugin = (DrawArtPlugin *) plugin;
	return kNoErr;
}

ASErr DrawArtPlugin::ShutdownPlugin(SPInterfaceMessage *message)
{
	// Destroy the flash UI
	if (fPanelController != NULL)
	{
		fPanelController->RemoveEventListeners();
		delete fPanelController;
		fPanelController = NULL;
		Plugin::LockPlugin(false);
	}
	
	return kNoErr;
}

ASErr DrawArtPlugin::Notify(AINotifierMessage* message)
{
	ASErr error = kNoErr;
	
	if (message->notifier == fPlugPlugSetupCompleteNotifier)
	{
		if (fPanelController != NULL)
			fPanelController->RegisterCSXSEventListeners();
	}
	
	return error;
}

ASErr DrawArtPlugin::Message(char *caller, char *selector, void *message)
{
	ASErr error = kNoErr;
	

	error = Plugin::Message(caller, selector, message);


	return error;
}

/*
*/
ASErr DrawArtPlugin::GoMenuItem(AIMenuMessage* message) 
{
	// Compare the menuItem selected (in the message) with our stored values to see if we do anything
	if (message->menuItem == this->fAboutPluginMenu) {
		// Pop this plug-in's about box.
		SDKAboutPluginsHelper aboutPluginsHelper;
		aboutPluginsHelper.PopAboutBox(message, "About DrawArtHtml", kSDKDefAboutSDKCompanyPluginsAlertString);
	}
	return kNoErr;
}

//Lookup table for encoding
//If you want to use an alternate alphabet, change the characters here
const static TCHAR encodeLookup[] = TEXT("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/");
const static TCHAR padCharacter = TEXT('=');
std::basic_string<TCHAR> base64Encode(std::vector<BYTE> inputBuffer)
{
	std::basic_string<TCHAR> encodedString;
	encodedString.reserve(((inputBuffer.size()/3) + (inputBuffer.size() % 3 > 0)) * 4);
	DWORD temp;
	std::vector<BYTE>::iterator cursor = inputBuffer.begin();
	for(size_t idx = 0; idx < inputBuffer.size()/3; idx++)
	{
		temp  = (*cursor++) << 16; //Convert to big endian
		temp += (*cursor++) << 8;
		temp += (*cursor++);
		encodedString.append(1,encodeLookup[(temp & 0x00FC0000) >> 18]);
		encodedString.append(1,encodeLookup[(temp & 0x0003F000) >> 12]);
		encodedString.append(1,encodeLookup[(temp & 0x00000FC0) >> 6 ]);
		encodedString.append(1,encodeLookup[(temp & 0x0000003F)      ]);
	}
	switch(inputBuffer.size() % 3)
	{
        case 1:
			temp  = (*cursor++) << 16; //Convert to big endian
			encodedString.append(1,encodeLookup[(temp & 0x00FC0000) >> 18]);
			encodedString.append(1,encodeLookup[(temp & 0x0003F000) >> 12]);
			encodedString.append(2,padCharacter);
			break;
        case 2:
			temp  = (*cursor++) << 16; //Convert to big endian
			temp += (*cursor++) << 8;
			encodedString.append(1,encodeLookup[(temp & 0x00FC0000) >> 18]);
			encodedString.append(1,encodeLookup[(temp & 0x0003F000) >> 12]);
			encodedString.append(1,encodeLookup[(temp & 0x00000FC0) >> 6 ]);
			encodedString.append(1,padCharacter);
			break;
	}
	return encodedString;
}

std::vector<BYTE> convertARGBtoRGBA(std::vector<BYTE> argb )
{
    std::vector<BYTE> rgba(argb.size());
    for(int i =0; i < argb.size(); i+=4)
    {
        rgba[i] = argb[i+1];    //R
        rgba[i+1] = argb[i+2];  //G
        rgba[i+2] = argb[i+3];  //B
        rgba[i+3] = argb[i];    //A
    }
    return rgba;
}

ASErr DrawArtPlugin::DrawAGMPort(std::string& encodedBitmap)
{
	
	ASErr error = kNoErr;
	
	AISwatchRef swatchRef = sAISwatchList->GetSwatchByName(NULL, ai::UnicodeString("CMYK Red"));
	
	AIDrawArtAGMPort port;
	error = sAIDrawArtSuite->CreateImagePort(50, 50, &port.port);
	aisdk::check_ai_error(error);
	
	AIDrawColorData colorData;
	AIColor color;
	error = sAISwatchList->GetAIColor(swatchRef,&color);
	
	colorData.color = color;
	colorData.output.port = port;
	colorData.type = 2;
	colorData.style = 0;
	colorData.rect.bottom = 50;
	colorData.rect.right = 50;
	colorData.rect.left =0;
	colorData.rect.top = 0;
	colorData.options = 0;
	colorData.width = 0;
	
	aisdk::check_ai_error(error);
	colorData.color = color;
	colorData.output.port = port;
	
	error = sAIDrawArtSuite->DrawColorSwatch(&colorData);
	
	aisdk::check_ai_error(error);
	size_t length = 0;
	ai::uint32* width = 0;
	ai::uint32* height = 0;
	const ai::uint8* pixelData = NULL;
	error = sAIDrawArtSuite->GetImagePixelData(port.port, &pixelData, &length, width, height);
	std::vector<BYTE> buffer;
	//length = length / 4; //When using individual values
	for(size_t i = 0; i < length ;i++)
	{
		/*
		 How to get the individual values
		 const ai::uint8* a = pixelData++;
		 const ai::uint8* r = pixelData++;	
		 const ai::uint8* g = pixelData++;
		 const ai::uint8* b = pixelData++;
		 ai::uint32 hex = NULL;
		 hex = (*a << 24) | (*r << 16) | (*g << 8) | *b;
		 */
		
		buffer.push_back(*pixelData);
		pixelData++;
	}
    std::vector<BYTE> rgbaBuffer = convertARGBtoRGBA(buffer);
    std::string encoded = base64Encode(rgbaBuffer);
	encodedBitmap.assign(encoded);
	
	aisdk::check_ai_error(error);
	error = sAIDrawArtSuite->DestroyImagePort(port.port);
	aisdk::check_ai_error(error);
	return error;
}
