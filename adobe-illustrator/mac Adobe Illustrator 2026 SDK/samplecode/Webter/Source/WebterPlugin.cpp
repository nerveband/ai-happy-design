//========================================================================================
//  
//  $File$
//
//  $Revision$
//
//  Copyright 1987 Adobe Systems Incorporated. All rights reserved.
//  
//  NOTICE:  Adobe permits you to use, modify, and distribute this file in accordance 
//  with the terms of the Adobe license agreement accompanying it.  If you have received
//  this file from a source other than Adobe, then your use, modification, or 
//  distribution of it requires the prior written permission of Adobe.
//  
//========================================================================================

#include "IllustratorSDK.h"

#include <math.h>
#include <stdio.h>
#include <string.h>
#include <stddef.h>

#include "WebterPlugin.hpp"
#include "WebterPanelController.h"
#include "WebterSuites.h"

#include "SDKDef.h"
#include "SDKAboutPluginsHelper.h"
#include "WebterID.h"
#include "AICSXS.h"
#include "IAIArtMatcher.hpp"

WebterPlugin *gPlugin = NULL;

Plugin *AllocatePlugin(SPPluginRef pluginRef)
{
	return new WebterPlugin(pluginRef);
}

/*
*/
void FixupReload(Plugin* plugin)
{
	WebterPlugin::FixupVTable((WebterPlugin*) plugin);
}

/*
*/
WebterPlugin::WebterPlugin(SPPluginRef pluginRef) :
	Plugin(pluginRef),
	fRegisterEventNotifierHandle(NULL),
	fWebterPanelController(NULL),
	fMenuItem(NULL),
	fAboutPluginMenu(NULL),
	fArtSelectionChangedNotifier(NULL)
{
	fTool[0] = fTool[1] = fTool[2] = nullptr;
	strncpy(fPluginName, "Webter", kMaxStringLength);
}

/*
*/
ASErr WebterPlugin::Message(char* caller, char* selector, void *message) 
{
	ASErr error = kNoErr;

	try {
		error = Plugin::Message(caller, selector, message);
	}
	catch (ai::Error& ex) {
		error = ex;
	}
	catch (...) {
		error = kCantHappenErr;
	}
	if (error) {
		if (error == kUnhandledMsgErr) {
			// Defined by Plugin.hpp and used in Plugin::Message - ignore.
			error = kNoErr;
		}
		else {
			Plugin::ReportError(error, caller, selector, message);
		}
	}	
	return error;
}


ASErr WebterPlugin::SetGlobal(Plugin *plugin)
{
	gPlugin = (WebterPlugin *) plugin;
	return kNoErr;
}

/*
*/
ASErr WebterPlugin::GoMenuItem( AIMenuMessage *message )
{
	if ( message->menuItem == fMenuItem ) 
	{
		fWebterPanelController->LoadExtension();
	}
	else if (message->menuItem == fAboutPluginMenu) 
	{
		SDKAboutPluginsHelper aboutPluginsHelper;
		aboutPluginsHelper.PopAboutBox(message, "About Webter", kSDKDefAboutSDKCompanyPluginsAlertString);
	}
	
	return kNoErr;
}

ASErr WebterPlugin::UpdateMenuItem( AIMenuMessage *message )
{
	//TODO - Use the AICSXSExtension suite.
	ASBoolean visible = false;
	if ( visible )
	{
		sAIMenu->SetItemText( message->menuItem, ai::UnicodeString("Hide Webter Dialog") );
	}
	else
	{
		sAIMenu->SetItemText( message->menuItem, ai::UnicodeString("Show Webter Dialog" ));
	}
	return kNoErr;
}

ASErr WebterPlugin::StartupPlugin( SPInterfaceMessage *message )
{
	ASErr error = Plugin::StartupPlugin( message );

	if (!error && this->fWebterPanelController == NULL) 
	{
		this->fWebterPanelController = new WebterPanelController();
		if (error) { return error; }
	}
	
	if (!error) 
	{
		// Add About Plugins menu item for this plug-in.
		SDKAboutPluginsHelper aboutPluginsHelper;
		error = aboutPluginsHelper.AddAboutPluginsMenuItem(message, 
					kSDKDefAboutSDKCompanyPluginsGroupName, 
					ai::UnicodeString(kSDKDefAboutSDKCompanyPluginsGroupNameString), 
					"Webter...", 
					&fAboutPluginMenu);
	}

	// Add other menu items.
	if ( !error )
	{
		// search through the menus and
		// see if another SDK plug in has already set up the group
		AIMenuGroup menuGroup, pluginMenuGroup,dummyGroup;
		AIMenuItemHandle dummyItem;
		AIPlatformAddMenuItemDataUS menuData;
		bool groupAlreadyMade = false;
		const char *sdkMenuStr = "SDK";
		const char *menuGroupStr = "SDKWindowGroup";
		const char *pluginMenuStr = "Webter";

		menuData.groupName = kToolPalettesMenuGroup;//kOthersPalettesMenuGroup;
		menuData.itemText = ai::UnicodeString( sdkMenuStr );
		
		ai::int32 count;
		error = sAIMenu->CountMenuGroups( &count );
		if ( error ) return error;
		while (count > 0)
		{
			error = sAIMenu->GetNthMenuGroup( count-1, &dummyGroup );
			if ( error ) return error;
			char *name;
			error = sAIMenu->GetMenuGroupName( dummyGroup, (const char**)&name );
			if ( error ) return error;
			if ( std::strcmp(name, menuGroupStr ) == 0 )
			{
				groupAlreadyMade = true;
				count = -1;
			}
			count--;
		}

		if ( !groupAlreadyMade && !error)
		{
			error = sAIMenu->AddMenuItem( message->d.self,
				                          menuGroupStr,
										  &menuData,
										  0,
										  &dummyItem );
			if (!error)
				error = sAIMenu->AddMenuGroupAsSubMenu(menuGroupStr, kMenuGroupSortedAlphabeticallyOption, dummyItem, &menuGroup );
		}

		menuData.groupName = menuGroupStr;
		menuData.itemText = ai::UnicodeString(pluginMenuStr);
		error = sAIMenu->AddMenuItem( message->d.self, 
								      pluginMenuStr, 
								      &menuData, 
								      0,
  								      &dummyItem );
		if ( error ) return error;

		error = sAIMenu->AddMenuGroupAsSubMenu( pluginMenuStr, 0, dummyItem, &pluginMenuGroup );
		if ( error ) return error;

		menuData.groupName = pluginMenuStr;
		char menuStr[] = "Show Webter Dialog";
		menuData.itemText = ai::UnicodeString( menuStr );

		if (!error)
			error = sAIMenu->AddMenuItem( message->d.self, 
								  	      "Webter Menu Item", 
										  &menuData, 
										  kMenuItemWantsUpdateOption,
										  &fMenuItem );
	}
	
	// add a notifier	
	if (!error)
	{
		error = sAINotifier->AddNotifier(fPluginRef, "WebterPlugin",
									 kAIArtSelectionChangedNotifier, &fArtSelectionChangedNotifier);
		if (!error)
		{
			error = sAINotifier->AddNotifier(fPluginRef, "WebterPlugin",
									 kAIDocumentChangedNotifier, &fDocumentChangedNotifier);

			if (!error) 
			{
				error = sAINotifier->AddNotifier(fPluginRef, "WebterPlugin", 
									 kAIDocumentAboutToCloseNotifier, &fDocumentClosedNotifier);
			}
		}

		error = sAINotifier->AddNotifier( fPluginRef, "Register Event Notify",
		kAICSXSPlugPlugSetupCompleteNotifier, &fRegisterEventNotifierHandle);
	}
										

	return error;
}

ASErr WebterPlugin::ShutdownPlugin( SPInterfaceMessage *message )
{
	if (fWebterPanelController)
	{
		fWebterPanelController->RemoveEventListeners();
		fWebterPanelController->SetArtDictionary(NULL, ai::UnicodeString("No art selected"));
		fWebterPanelController->SendDocDictionaryToPanel();
		delete fWebterPanelController;
		fWebterPanelController = NULL;
		Plugin::LockPlugin(false);
	}

	message->d.globals = NULL;	
	ASErr error = Plugin::ShutdownPlugin(message);

	return error;

}

ASErr WebterPlugin::Notify(AINotifierMessage *message)
{
	ASErr error = kNoErr;
	
	if (message->notifier == fArtSelectionChangedNotifier)
	{
		AIDictionaryRef theDictionary = NULL;
		ai::int32 theArtCount = 0;
		ai::ArtMatcher theArtStore(AIMatchingArtSpec{ kAnyArt, kArtSelected, kArtSelected });
		theArtStore.GetMatchingArt();

		if (theArtStore)
		{
			AIArtHandle theArt = theArtStore.GetMatchedArt(0);
			AIArtHandle child = theArt;
			while (child)
			{
				theArt = child;
				sAIArt->GetArtFirstChild(theArt, &child);
			}
			
			ASBoolean isDefaultName (false);
			ai::UnicodeString artName;
			sAIArt->GetArtName(theArt, artName, &isDefaultName);
			fWebterPanelController->SetArtDictionary(theArt, artName);
		}
		else
		{
			fWebterPanelController->SetArtDictionary(NULL, ai::UnicodeString("No art selected"));
		}
	}
	else if (message->notifier == fDocumentClosedNotifier)
	{
		// No document, so no artwork ==> clear artwork dictionary too
		fWebterPanelController->SetArtDictionary(NULL, ai::UnicodeString("No art selected"));
		
		// Release handle on document
		fWebterPanelController->SendDocDictionaryToPanel();
	}
	else if (message->notifier == fDocumentChangedNotifier)
	{
		AIDocumentHandle doc(nullptr);
		sAIDocument->GetDocument(&doc);
		
		if (doc == nullptr)
		{
			// No document, so no artwork ==> clear artwork dictionary too
			fWebterPanelController->SetArtDictionary(nullptr);
		}
		
		fWebterPanelController->SendDocDictionaryToPanel();
	}
	else if ( message->notifier == fRegisterEventNotifierHandle)
	{
		fWebterPanelController->RegisterCSXSEventListeners();
	}

	return error;
}


