//========================================================================================
//  
//  $File: //ai/mainline/devtech/sdk/public/samplecode/ScriptMessage/Source/ScriptMessagePlugin.cpp $
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
#include "AIScriptMessage.h"
#include "ScriptMessagePlugin.h"
#include "AICSXS.h"

/** Get a string representation of an AIArtType
 */
static std::string GetArtTypeName(short artType)
{
	switch (artType)
	{
		case(kAnyArt): return "AnyArt";
		case(kUnknownArt): return "UnknownArt";
		case(kGroupArt): return "GroupArt";
		case(kPathArt): return "PathArt";
		case(kCompoundPathArt): return "CompoundPathArt";
		case(kTextArtUnsupported): return "TextArtUnsupported";
		case(kTextPathArtUnsupported): return "TextPathArtUnsupported";
		case(kTextRunArtUnsupported): return "TextRunArtUnsupported";
		case(kPlacedArt): return "PlacedArt";
		case(kMysteryPathArt): return "MysteryPathArt";
		case(kRasterArt): return "RasterArt";
		case(kPluginArt): return "PluginArt";
		case(kMeshArt): return "MeshArt";
		case(kTextFrameArt): return "TextFrameArt";
		case(kSymbolArt): return "SymbolArt";
		case(kForeignArt): return "ForeignArt";
		case(kLegacyTextArt): return "LegacyTextArt";
		default: return "Error";
	}
}

Plugin* AllocatePlugin(SPPluginRef pluginRef)
{
	return new ScriptMessagePlugin(pluginRef);
}

void FixupReload(Plugin* plugin)
{
	ScriptMessagePlugin::FixupVTable((ScriptMessagePlugin*) plugin);
}

ScriptMessagePlugin::ScriptMessagePlugin(SPPluginRef pluginRef) :
	Plugin(pluginRef),
	fPanelController(NULL),
	fPlugPlugSetupCompleteNotifier(NULL),
	fAboutPluginMenu(NULL),
	fBrightnessChangedNotifier(NULL)
{
	strncpy(fPluginName, kScriptMessagePluginName, kMaxStringLength);
}

ASErr ScriptMessagePlugin::StartupPlugin(SPInterfaceMessage *message)
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
				"ScriptMessage...", 
				&fAboutPluginMenu);
	
	// Create panel controller
	if (fPanelController == NULL)
	{
		fPanelController = new ScriptMessagePanelController();
		
		error = Plugin::LockPlugin(true);
		
		if (error)
			return error;
	}
		
	// Add notifier for adding CSXS event listeners
	error = sAINotifier->AddNotifier(fPluginRef, kScriptMessagePluginName,
									 kAICSXSPlugPlugSetupCompleteNotifier, &fPlugPlugSetupCompleteNotifier);
	
	if (error)
		return error;
	
	// Add notifier for theme changes
	error = sAINotifier->AddNotifier(fPluginRef, kScriptMessagePluginName,
									 kAIUIBrightnessChangedNotifier, &fBrightnessChangedNotifier);
	
	return error;
}

ASErr ScriptMessagePlugin::ShutdownPlugin(SPInterfaceMessage *message)
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

ASErr ScriptMessagePlugin::Notify(AINotifierMessage* message)
{
	ASErr error = kNoErr;
	
	if (message->notifier == fBrightnessChangedNotifier)
	{
		if (fPanelController != NULL)
			fPanelController->SendThemeToPanel();
	}
	else if (message->notifier == fPlugPlugSetupCompleteNotifier)
	{
		if (fPanelController != NULL)
			fPanelController->RegisterCSXSEventListeners();
	}
	
	return error;
}

ASErr ScriptMessagePlugin::Message(char *caller, char *selector, void *message)
{
	ASErr error = kNoErr;
	
	/** If there is artwork selected, then display information about it,
		otherwise just display a basic message.
	*/
	if (strcmp(caller, kCallerAIScriptMessage ) == 0)
	{
		error = kNoErr;
		AIScriptMessage *msg = (AIScriptMessage*) message;

		ai::UnicodeString msgString("Message recieved!");
		msgString.append(ai::UnicodeString("\rSelector: "));
		msgString.append(ai::UnicodeString(selector));
		msgString.append(ai::UnicodeString("\rMessage: "));
		msgString.append(msg->inParam);
		
		if (msg->art != NULL) // artwork selected
		{
			msgString.append(ai::UnicodeString("\rArtwork: "));
			
			// Get the name of the art object
			ai::UnicodeString nameOfArt;
			ASBoolean isDefaultName;
			sAIArt->GetArtName(msg->art, nameOfArt, &isDefaultName);
			msgString.append(nameOfArt);
			if (isDefaultName)
				msgString.append(ai::UnicodeString(" (default name)"));
			
			// Get the note attached to the art, if there is one
			msgString.append(ai::UnicodeString("\rNotes: "));
			if (sAIArt->HasNote(msg->art))
			{
				ai::UnicodeString artNote;
				sAIArt->GetNote(msg->art, artNote);
				msgString.append(artNote);
			}
			else
			{
				msgString.append(ai::UnicodeString("None"));
			}
			
			// Get the type of the Art object
			msgString.append(ai::UnicodeString("\rArtType: "));
			short artType;
			sAIArt->GetArtType(msg->art, &artType);
			msgString.append(ai::UnicodeString(GetArtTypeName(artType)));
			
		}
		else // No artwork selected
		{
			msgString.append(ai::UnicodeString("\rArtwork: None"));
		}
		
		// Display an alert to the user
		sAIUser->MessageAlert(ai::UnicodeString(msgString));
								 
		//Return the message to the script
		msg->outParam = msgString;

	}
	else
	{
		error = Plugin::Message(caller, selector, message);
	}

	return error;
}

/*
*/
ASErr ScriptMessagePlugin::GoMenuItem(AIMenuMessage* message) 
{
	// Compare the menuItem selected (in the message) with our stored values to see if we do anything
	if (message->menuItem == this->fAboutPluginMenu) {
		// Pop this plug-in's about box.
		SDKAboutPluginsHelper aboutPluginsHelper;
		aboutPluginsHelper.PopAboutBox(message, "About ScriptMessage", kSDKDefAboutSDKCompanyPluginsAlertString);
	}
	return kNoErr;
}