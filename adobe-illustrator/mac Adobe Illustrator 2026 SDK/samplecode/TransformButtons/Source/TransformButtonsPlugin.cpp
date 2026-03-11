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
#include "TransformButtonsPlugin.h"
#include "TransformButtonsPanelController.h"
#include "AICSXS.h"

TransformButtonsPlugin *gPlugin = NULL;

/*
*/
Plugin* AllocatePlugin(SPPluginRef pluginRef)
{
	return new TransformButtonsPlugin(pluginRef);
}

/*
*/
void FixupReload(Plugin* plugin)
{
	TransformButtonsPlugin::FixupVTable((TransformButtonsPlugin*) plugin);
}

/*
*/
TransformButtonsPlugin::TransformButtonsPlugin(SPPluginRef pluginRef) 
	: Plugin(pluginRef),
	fShowTransformButtonsDialog(NULL),
	fAboutPluginMenu(NULL),
	fTransformButtonsPanelController(NULL),
	fRegisterEventNotifierHandle(NULL)
{
	strncpy(fPluginName, kTransformButtonsPluginName, kMaxStringLength);
}

/*
*/
ASErr TransformButtonsPlugin::Message(char* caller, char* selector, void *message) 
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


ASErr TransformButtonsPlugin::SetGlobal(Plugin *plugin)
{
	gPlugin = (TransformButtonsPlugin *) plugin;
	return kNoErr;
}

/*
*/
ASErr TransformButtonsPlugin::StartupPlugin(SPInterfaceMessage* message)
{
	AIErr error = kNoErr;
	error = Plugin::StartupPlugin(message);
    if (error) { return error; }
	if (this->fTransformButtonsPanelController == NULL) {
		this->fTransformButtonsPanelController = new TransformButtonsPanelController();
		error = Plugin::LockPlugin(true);
		if (error) { return error; }
	}
	error = this->AddMenus(message);
    if (error) { return error; }
	error = this->AddNotifiers(message);
	return error;
}

/*
*/
ASErr TransformButtonsPlugin::ShutdownPlugin(SPInterfaceMessage* message)
{	
	if (fTransformButtonsPanelController)
	{
		fTransformButtonsPanelController->RemoveEventListeners();
		delete fTransformButtonsPanelController;
		fTransformButtonsPanelController = NULL;
		Plugin::LockPlugin(false);
	}
	message->d.globals = NULL;	
	return Plugin::ShutdownPlugin(message);
}


/*
*/
ASErr TransformButtonsPlugin::PostStartupPlugin()
{
	return kNoErr;
}

/*
*/
ASErr TransformButtonsPlugin::GoMenuItem(AIMenuMessage *message) 
{
	ASErr error = kNoErr;
	
	if (message->menuItem == this->fAboutPluginMenu) {
		// Pop this plug-in's about box.
		SDKAboutPluginsHelper aboutPluginsHelper;
		aboutPluginsHelper.PopAboutBox(message, "About TransformButtons", kSDKDefAboutSDKCompanyPluginsAlertString);
	}	
	else if (message->menuItem == this->fShowTransformButtonsDialog) 
	{
		if (fTransformButtonsPanelController)
		{
			ASBoolean visible = false;//fTransformButtonsPanelController->GetPanelReady();	
			if(visible)
			{
				fTransformButtonsPanelController->UnloadExtension();
				//fTransformButtonsPanelController->SetPanelReady(false);
			}
			else
			{
				fTransformButtonsPanelController->LoadExtension();
			}
		}
	}
	
	if (error)
		goto error;

error:
	return error;
}

/*
*/
ASErr TransformButtonsPlugin::UpdateMenuItem(AIMenuMessage* message) 
{
	ASErr error = kNoErr;

	ASBoolean visible = false;
	/*if (fTransformButtonsPanelController)
	{
		;
		//visible = freeGridDialogController->GetPanelReady();
	}*/
	if (visible) {
		sAIMenu->SetItemText(message->menuItem, ai::UnicodeString("Hide Transform Buttons Dialog"));
	}
	else {
		sAIMenu->SetItemText(message->menuItem, ai::UnicodeString("Show Transform Buttons Dialog"));
	}

	if (error)
		goto error;

error:
	return error;
}

/*
*/
ASErr TransformButtonsPlugin::AddNotifiers(SPInterfaceMessage* message) 
{
	ASErr error = kNoErr;

	error = sAINotifier->AddNotifier( fPluginRef, "Register Event Notify",
		kAICSXSPlugPlugSetupCompleteNotifier, &fRegisterEventNotifierHandle);

	if (error) goto error;

error:
	return error;
}

/*
*/
ASErr TransformButtonsPlugin::AddMenus(SPInterfaceMessage* message)
{
	ASErr error = kNoErr;

	// Add a menu item to the About SDK Plug-ins menu group.
	SDKAboutPluginsHelper aboutPluginsHelper;
	aboutPluginsHelper.AddAboutPluginsMenuItem(message, 
				kSDKDefAboutSDKCompanyPluginsGroupName, 
				ai::UnicodeString(kSDKDefAboutSDKCompanyPluginsGroupNameString), 
				"TransformButtons...", 
				&this->fAboutPluginMenu);

	// Add other menus.
	AIPlatformAddMenuItemDataUS menuData;
	const char *sdkMenuStr = "SDK";
	const char *menuGroupStr = "SDKWindowGroup";
	const char *pluginMenuStr = "Transform Buttons";
	AIMenuGroup	menuGroup,pluginMenuGroup,dummyGroup;
	AIMenuItemHandle dummyItem;
	bool groupAlreadyMade = false;

	menuData.groupName = kOtherPalettesMenuGroup;
	menuData.itemText = ai::UnicodeString( sdkMenuStr );
	
	// search through the menus and 
	// see if another SDK plug in has already set up the group
	ai::int32 count;
	error = sAIMenu->CountMenuGroups( &count );
	if ( error ) goto error;
	while (count > 0)
	{
		error = sAIMenu->GetNthMenuGroup( count-1, &dummyGroup );
		if ( error ) goto error;
		char *name;
		error = sAIMenu->GetMenuGroupName( dummyGroup, (const char**)&name );
		if ( error ) goto error;
		if ( std::strcmp(name, menuGroupStr ) == 0 )
		{
			groupAlreadyMade = true;
			count = -1;
		}
		count--;
	}

	if ( !groupAlreadyMade )
	{
		error = sAIMenu->AddMenuItem( message->d.self, 
								      menuGroupStr, 
								      &menuData, 
								      0,
  								      &dummyItem );
		if ( error ) goto error;

		error = sAIMenu->AddMenuGroupAsSubMenu( menuGroupStr, kMenuGroupSortedAlphabeticallyOption, dummyItem, &menuGroup );
		if ( error ) goto error;
	}
	
	menuData.groupName = menuGroupStr;
	menuData.itemText = ai::UnicodeString(pluginMenuStr);
	error = sAIMenu->AddMenuItem( message->d.self, 
								      menuGroupStr, 
								      &menuData, 
								      0,
  								      &dummyItem );
	if ( error ) goto error;

	error = sAIMenu->AddMenuGroupAsSubMenu( pluginMenuStr, 0, dummyItem, &pluginMenuGroup );
	if ( error ) goto error;
	menuData.groupName = pluginMenuStr;
	error = sAIMenu->AddMenuItem( message->d.self, 
							      "Transform Buttons Menu Item", 
							      &menuData, 
							      kMenuItemWantsUpdateOption,
							      &this->fShowTransformButtonsDialog );
	if ( error )
		goto error;
		
error:
	return error;
}

/*
*/
ASErr TransformButtonsPlugin::Notify(AINotifierMessage* message) 
{
	ASErr error = kNoErr;

	if ( message->notifier == fRegisterEventNotifierHandle)
	{
		fTransformButtonsPanelController->RegisterCSXSEventListeners();
	}
	return error;
}
// End TransformButtonsPlugin.cpp
