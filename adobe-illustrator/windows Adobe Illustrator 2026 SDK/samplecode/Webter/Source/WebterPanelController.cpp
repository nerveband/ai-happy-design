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

#include "WebterPanelController.h"
#include "AppContext.hpp"
#include "SDKErrors.h"
#include "Plugin.hpp"
#include "AIDictionary.h"
#include "AIStringFormatUtils.h"
#include "DictionaryToXmlConverter.h"
#include "WebterSuites.h"

#define WEBTER_UI_EXTENSION				"com.adobe.illustrator.webterui.panel"
#define ILST_APPID						"ILST"

// Events we listen out for
#define EVENT_TYPE_ADDENTRY_CLICKED		"com.adobe.csxs.events.WebterAddEntry"
#define EVENT_TYPE_ITEM_SELECTED		"com.adobe.csxs.events.WebterItemSelect"
#define EVENT_TYPE_PANEL_READY			"com.adobe.csxs.events.WebterPanelReady"

// Events we dispatch
#define EVENT_UPDATE_ENTRY				"com.adobe.csxs.events.WebterUpdateEntry"
#define EVENT_UPDATE_DEFN				"com.adobe.csxs.events.WebterUpdateDefinitionAndType"

static void AddEntryToDictionaryClickedFunc(const csxs::event::Event* const eventParam, void* const context)
{
	AppContext appContext(gPlugin->GetPluginRef());

	WebterPanelController *webterPanelController = (WebterPanelController *)context;
	if ( NULL == webterPanelController || eventParam == NULL )
		return;
	webterPanelController->ParseData(eventParam->data);
	webterPanelController->AddEntryToDictionary();
}

static void ItemSelectedFunc(const csxs::event::Event* const eventParam, void* const context)
{
	AppContext appContext(gPlugin->GetPluginRef());

	WebterPanelController *webterPanelController = (WebterPanelController *)context;
	if ( NULL == webterPanelController || eventParam == NULL )
		return;
	webterPanelController->ParseData(eventParam->data);
	webterPanelController->UpdateDefinitionAndType();
}

static void PanelReadyFunc(const csxs::event::Event* const eventParam, void* const context)
{
	AppContext appContext(gPlugin->GetPluginRef());

	WebterPanelController *webterPanelController = (WebterPanelController *)context;
	if ( NULL == webterPanelController || eventParam == NULL )
		return;
	
	webterPanelController->SendDocDictionaryToPanel();
	webterPanelController->SendArtDictionaryToPanel();
}

WebterPanelController::WebterPanelController(void)
	: FlashUIController(WEBTER_UI_EXTENSION)
{
}

WebterPanelController::~WebterPanelController(void)
{
}

/* Add event listeners
*/
csxs::event::EventErrorCode WebterPanelController::RegisterCSXSEventListeners()
{
	csxs::event::EventErrorCode result = csxs::event::kEventErrorCode_Success;
	do
	{
		result = fPPLib.AddEventListener(EVENT_TYPE_PANEL_READY, PanelReadyFunc, this);
		if (result != csxs::event::kEventErrorCode_Success)
		{
			break;
		}
	
		result =  fPPLib.AddEventListener(EVENT_TYPE_ADDENTRY_CLICKED, AddEntryToDictionaryClickedFunc, this);
		if (result != csxs::event::kEventErrorCode_Success) 
		{
			break;
		}

		result =  fPPLib.AddEventListener(EVENT_TYPE_ITEM_SELECTED, ItemSelectedFunc, this);
		if (result != csxs::event::kEventErrorCode_Success) 
		{
			break;
		}
	}
	while(false);
	return result;
}

/* Remove event listeners
*/
csxs::event::EventErrorCode WebterPanelController::RemoveEventListeners()
{
	csxs::event::EventErrorCode result = csxs::event::kEventErrorCode_Success;
	do
	{
		result = fPPLib.RemoveEventListener(EVENT_TYPE_PANEL_READY, PanelReadyFunc, this);
		if (result != csxs::event::kEventErrorCode_Success)
		{
			break;
		}
		
		result =  fPPLib.RemoveEventListener(EVENT_TYPE_ADDENTRY_CLICKED, AddEntryToDictionaryClickedFunc, this);
		if (result != csxs::event::kEventErrorCode_Success) 
		{
			break;
		}

		result =  fPPLib.RemoveEventListener(EVENT_TYPE_ITEM_SELECTED, ItemSelectedFunc, this);
		if (result != csxs::event::kEventErrorCode_Success) 
		{
			break;
		}
	}
	while(false);
	return result;
}

void WebterPanelController::ParseData(const char* eventData)
{
	string srcData(eventData);
	
	ai::NumberFormat numFormat;
	string itemStr = WebterPanelController::GetElement(srcData, "entry");
	if ( !itemStr.empty() )
	{
		ASInt32 temp = atoi(itemStr.c_str());
		itemType = (ItemType)temp;
	}

	string entryStr = WebterPanelController::GetElement(srcData, "entrystring");
	if ( !entryStr.empty() )
	{
		entryString = entryStr;
	}

	string defintionStr = WebterPanelController::GetElement(srcData, "value");
	if ( !defintionStr.empty() )
	{
		stringValue = defintionStr;
	}

	string entryTypeStr = WebterPanelController::GetElement(srcData, "entrytype");
	if ( !entryTypeStr.empty() )
	{
		entryType= atoi(entryTypeStr.c_str());
	}

}

void WebterPanelController::SendDocDictionaryToPanel()
{
	AIDocumentHandle doc(nullptr);
	sAIDocument->GetDocument(&doc);
	string payload;
	if (doc)
	{
		DocDictionaryToXmlConverter converter(ai::Dictionary (ai::Dictionary::kCurrentDocument), entryString);
		payload = converter.GetPayload();
	}
	else
	{
		payload = "<payload><dict>doc</dict><null>true</null></payload>";
	}

	dispatchDictionaryUpdateEvent(payload);
}

csxs::event::EventErrorCode WebterPanelController::dispatchDictionaryUpdateEvent(const string& payload)
{
	csxs::event::Event event = {EVENT_UPDATE_ENTRY,
		csxs::event::kEventScope_Application,
		ILST_APPID,
		NULL,
		payload.c_str()};
	
	return fPPLib.DispatchEvent(&event);
}

void WebterPanelController::SetArtDictionary(AIArtHandle art, const ai::UnicodeString& artName)
{
	if (art == nullptr)
	{
		safeArtHandle = nullptr;
		displayNameOfSelectedArt = ai::UnicodeString("No art selected").as_Roman();
		return;
	}
		
	sAIArt->GetSafeArtHandle(art, &safeArtHandle);

	ai::UnicodeString::size_type strLength = artName.length();
	if (strLength > 0) 
	{
		displayNameOfSelectedArt = artName.as_Roman();
		
		// Strip '<' & '>' if present
		if(displayNameOfSelectedArt.at(0) == '<')
			displayNameOfSelectedArt = displayNameOfSelectedArt.substr(1,
				displayNameOfSelectedArt.length()-2);
	}
	else
	{
		displayNameOfSelectedArt = "No art selected:";
	}
	
	SendArtDictionaryToPanel();
}

void WebterPanelController::SendArtDictionaryToPanel()
{
	if (safeArtHandle == nullptr)
		return;

	AIArtHandle art(nullptr);
	sAIArt->GetArtHandle(safeArtHandle, &art);
	if (art == nullptr)
		return;

	ArtDictionaryToXmlConverter converter(ai::Dictionary(art), displayNameOfSelectedArt);
	std::string payload = converter.GetPayload();
	dispatchDictionaryUpdateEvent(payload);
}

/* Called when an item is selected. Determines which dictionary to look up,
 * then looks up the selected item, gets a string representation of the entry,
 * if the entry is an integer/real/boolean/string, otherwise gets the type of
 * the entry and converts it to a string.
 * Then it dispatches this information to the flash panel.
 */
void WebterPanelController::UpdateDefinitionAndType()
{
	auto dict = GetCurrentDict();
	
	auto iter = dict.find(sAIDictionary->Key(entryString.c_str()));
	if (iter == dict.end())
		return;

	auto entry = (*iter).second;

	auto typeOfEntry = entry.GetType();
	
	string value;
	
	// Now get the value, if we can
	switch (typeOfEntry)
	{
		case IntegerType:
		{
			char IntBuffer[33];
			ai::int32 intValue = entry;
			sprintf(IntBuffer,"%d",intValue);
			value = IntBuffer;
			break;
		}
			
		case RealType:
		{
			char Realbuffer[33];
			AIReal realValue = entry;
			sprintf(Realbuffer,"%f",realValue);
			value = Realbuffer;			
			break;
		}
			
		case StringType:
		{
			value = entry.StringVal();
			break;
		}
			
		case BooleanType:
		{
			bool boolEntry = entry;
			value = (boolEntry) ? "True" : "False";
			break;
		}
			
		default:
		{
			value = "---";
			break;
		}
	}
	
	
	stringstream payloadStream;
	payloadStream << "<payload><dict>" << ((itemType == DOCUMENT) ? "doc" : "art")
	<< "</dict><key>" << entryString << "</key><value>" << value
	<< "</value><type>" << typeOfEntry << "</type></payload>";
	
	std::string payloadString = payloadStream.str();
	// Create and dispatch event
	csxs::event::Event event = {EVENT_UPDATE_DEFN,
								 csxs::event::kEventScope_Application,
								 ILST_APPID,
								 NULL,
								 payloadString.c_str()};
	fPPLib.DispatchEvent(&event);
}

void WebterPanelController::AddEntryToDictionary()
{
	auto dict = GetCurrentDict();
	
	auto key = sAIDictionary->Key(entryString.c_str());

	AIBoolean success = true;
	switch (entryType)
	{
		case IntegerType:
		{
			ai::int32 intVal = atoi(stringValue.c_str());
			dict.Set(key, ai::Entry(intVal));
			break;
		}
		case RealType:
		{	
			AIReal floatVal = (float)atof(stringValue.c_str());
			dict.Set(key, ai::Entry(floatVal));
			break;
		}
		case StringType:
		{
			dict.Set(key, ai::Entry(stringValue));
			break;
		}
		case BooleanType:
		{
			// TODO: Revisit this once we change the xml sent from the panel
			AIBoolean boolVal = atoi(stringValue.c_str());
			dict.Set(key, ai::Entry(boolVal != false));
			break;
		}
		default:
		{
			success = false;
			sAIUser->ErrorAlert(ai::UnicodeString("I don't support that type!"));
			break;
		}
	}

	if (success) // Let's send the updated dictionary to the panel
	{
		if (itemType == DOCUMENT)
			SendDocDictionaryToPanel();
		else
			SendArtDictionaryToPanel();
	}
}

string WebterPanelController::GetElement(const string& srcData, const char* value)
{
	string openStr("<");
	string contentStr(value);
	string endStr(">");
	openStr.append(contentStr);
	openStr.append(endStr);
	size_t openIndex = srcData.find(openStr);
	openIndex = openIndex + openStr.length();
	string strValue("");

    string colseStr("</");
	colseStr.append(contentStr);
	colseStr.append(endStr);
	size_t closeIndex = srcData.find(colseStr);
	size_t strLength = closeIndex - openIndex;
	if ( strLength > 0 && strLength < srcData.length())
	{
		strValue = srcData.substr(openIndex, strLength);
	}
	
	return strValue;
}

ai::Dictionary WebterPanelController::GetCurrentDict()
{
	ai::Dictionary dict(itemType == DOCUMENT ? ai::Dictionary(ai::Dictionary::kCurrentDocument) :
		[&] { AIArtHandle art(nullptr); sAIArt->GetArtHandle(safeArtHandle, &art); return ai::Dictionary(art); } ()
	);

	return dict;
}
