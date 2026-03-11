//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/Webter/Source/WebterPanelController.h $
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
#ifndef __WEBTERPANELCONTROLLER_H__
#define __WEBTERPANELCONTROLLER_H__

#include "SDKPlugPlug.h"
#include "IllustratorSDK.h"
#include "WebterID.h"
#include "FlashUIController.h"
#include "WebterPlugin.hpp"
#include "IAIDictionary.hpp"
#include <memory>

extern WebterPlugin *gPlugin;

enum ItemType
{
	DOCUMENT,
	ART
};

class WebterPanelController : public FlashUIController
{
public:
	WebterPanelController(void);

	~WebterPanelController(void);

	csxs::event::EventErrorCode RegisterCSXSEventListeners();

	csxs::event::EventErrorCode RemoveEventListeners();

	ASErr SendData() { return kNoErr;};

	void ParseData(const char* eventData);

	void SendDocDictionaryToPanel();
	
	void SetArtDictionary(AIArtHandle art, const ai::UnicodeString& artName = ai::UnicodeString());

	void SendArtDictionaryToPanel();
	
	void UpdateDefinitionAndType();

	void AddEntryToDictionary();

	static string GetElement(const string& srcData, const char* value);

	ai::Dictionary GetCurrentDict();

private:
	
	//std::string buildStringOfEntriesFromDictionary(AIDictionaryRef dict);

	std::string getDocDictionaryPayloadFromStringOfEntries(const string& entryList);
	std::string getArtDictionaryPayloadFromStringOfEntries(const string& entryList);

	csxs::event::EventErrorCode dispatchDictionaryUpdateEvent(const string& payload);
	
	AISafeArtHandle safeArtHandle = nullptr;

	std::string displayNameOfSelectedArt;
	
	ItemType itemType  = DOCUMENT;
	AIEntryType entryType = UnknownType;
	string stringValue;
	string entryString;	
};

#endif
