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
#include "LiveDropShadowPlugin.h"

/*
*/
Plugin* AllocatePlugin(SPPluginRef pluginRef)
{
	return new LiveDropShadowPlugin(pluginRef);
}

/*
*/
void FixupReload(Plugin* plugin)
{
	LiveDropShadowPlugin::FixupVTable((LiveDropShadowPlugin*) plugin);
}

/*
*/
LiveDropShadowPlugin::LiveDropShadowPlugin(SPPluginRef pluginRef) 
	: Plugin(pluginRef),
	fAboutPluginMenu(NULL),
	fMakeMenuHandle(NULL),
	fReleaseMenuHandle(NULL),
	fAdjustMenuHandle(NULL),
	fPluginGroupHandle(NULL),
	fAppShutdownNotifier(NULL),
	fLiveDropShadowGroup(NULL)
{
	strncpy(fPluginName, kLiveDropShadowPluginName, kMaxStringLength);
	fObjectData.init();
}

/*
*/
ASErr LiveDropShadowPlugin::Message(char* caller, char* selector, void *message) 
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


/*
*/
ASErr LiveDropShadowPlugin::StartupPlugin(SPInterfaceMessage* message)
{
	AIErr result = kNoErr;
    
	result = Plugin::StartupPlugin(message);
	if (result) { return result;  }
	result = this->AddMenus(message);
    if (result) { return result;  }
	result = this->AddNotifiers(message);
	
	return result;
}

/*
*/
ASErr LiveDropShadowPlugin::ShutdownPlugin(SPInterfaceMessage* message)
{	
	message->d.globals = NULL;			
	return Plugin::ShutdownPlugin(message);
}

/*
*/
ASErr LiveDropShadowPlugin::PostStartupPlugin()
{
	ASErr result = kNoErr;

	// Create LiveDropShadowGroup
	if (this->fLiveDropShadowGroup == NULL) {
		this->fLiveDropShadowGroup = new LiveDropShadowGroup();
		// TODO error handling
	}
	return result;
}

/*
*/
ASErr LiveDropShadowPlugin::Notify(AINotifierMessage* message)
{
	ASErr result = kNoErr;
	if (message->notifier == this->fAppShutdownNotifier) {
		if (this->fLiveDropShadowGroup) {
			delete this->fLiveDropShadowGroup;
			this->fLiveDropShadowGroup = NULL;
			result = Plugin::LockPlugin(false);
		}
	}
	return result;
}

/*
*/
ASErr LiveDropShadowPlugin::PluginGroupUpdate(AIPluginGroupMessage* message)
{
	ASErr result = kNoErr;
	result = this->fLiveDropShadowGroup->PluginGroupUpdate(message);
	return result;
}

/*
*/
ASErr LiveDropShadowPlugin::AddNotifiers(SPInterfaceMessage* message)
{
    ASErr result = kNoErr;
	result = sAINotifier->AddNotifier(message->d.self, "LiveDropShadow " kAIApplicationShutdownNotifier, 
			kAIApplicationShutdownNotifier, &this->fAppShutdownNotifier);
	return result;
}

/*
*/
ASErr LiveDropShadowPlugin::AddMenus(SPInterfaceMessage* message)
{
    ASErr result = kNoErr;

	// Add a menu item to the About SDK Plug-ins menu group.
	SDKAboutPluginsHelper aboutPluginsHelper;
	aboutPluginsHelper.AddAboutPluginsMenuItem(message, 
				kSDKDefAboutSDKCompanyPluginsGroupName, 
				ai::UnicodeString(kSDKDefAboutSDKCompanyPluginsGroupNameString), 
				"LiveDropShadow...", 
				&this->fAboutPluginMenu);	

	const char* sdkText = "SDK";
	AIPlatformAddMenuItemDataUS sdkMenuData;
	sdkMenuData.groupName = kObjectPathsMenuGroup;
	sdkMenuData.itemText = ai::UnicodeString(sdkText);
	AIMenuItemHandle sdkMenuItemHandle = NULL;
	result = sAIMenu->AddMenuItem(message->d.self, "SDKObject", &sdkMenuData, 0, &sdkMenuItemHandle);
	if (result) return result;

	AIMenuGroup sdkGroupMenu = NULL;
	result = sAIMenu->AddMenuGroupAsSubMenu("SDKObject", kMenuGroupNoOptions, sdkMenuItemHandle, &sdkGroupMenu);
	if (result) return result;

	AIAddPluginGroupData pluginGroupData;
	pluginGroupData.minor = kMinorVersion;
	pluginGroupData.major = kMajorVersion;
    char description[] = "Live Drop Shadow Plug-in Group";
	pluginGroupData.desc  = description;
    char name[] = kNameSuffix "class";
	result = sAIPluginGroup->AddAIPluginGroup(message->d.self, name, &pluginGroupData, kPluginGroupWantsAutoTranslateOption, &this->fPluginGroupHandle);
		
	AIPlatformAddMenuItemDataUS menuItemData;
	AIMenuItemHandle throwawayMenu = NULL;
	if (result == kNoErr) {	
		menuItemData.groupName = "SDKObject";
		menuItemData.itemText = ai::UnicodeString(kGroupMenuString);
		result = sAIMenu->AddMenuItem(message->d.self, kNameSuffix "group", &menuItemData, 0, &throwawayMenu);	
	}
	
	if (result == kNoErr) {
		AIMenuGroup throwawayMenuGroup = NULL;
		sAIMenu->AddMenuGroupAsSubMenu(kNameSuffix "group", 0, throwawayMenu, &throwawayMenuGroup);
		if (throwawayMenuGroup == NULL)
			result = kTooManyMenuItemsErr;	
	}
		
	if (result == kNoErr) {
		menuItemData.groupName = kNameSuffix "group";
		menuItemData.itemText = ai::UnicodeString(kMakeMenuString);
		result = sAIMenu->AddMenuItem(message->d.self, kNameSuffix "make menu", &menuItemData, kMenuItemWantsUpdateOption, &this->fMakeMenuHandle);
	}
		
	if (result == kNoErr) {
		menuItemData.groupName = kNameSuffix "group";
		menuItemData.itemText = ai::UnicodeString(kReleaseMenuString);
		result = sAIMenu->AddMenuItem(message->d.self, kNameSuffix "release menu", &menuItemData, kMenuItemWantsUpdateOption, &this->fReleaseMenuHandle);
	}
		
	if (result == kNoErr) {
		menuItemData.groupName = kNameSuffix "group";
		menuItemData.itemText = ai::UnicodeString(kAdjustMenuString);
		result = sAIMenu->AddMenuItem(message->d.self, kNameSuffix "adjust menu", &menuItemData, kMenuItemWantsUpdateOption, &this->fAdjustMenuHandle);
	}	
	return result;
}

/*
*/
ASErr LiveDropShadowPlugin::UpdateMenuItem(AIMenuMessage* message) 
{
	ASErr result = kNoErr;
	
	AIArtHandle **store = NULL;
	ai::int32 count = 0;
	if (result == kNoErr)
	{	
		if (sAIMatchingArt->GetSelectedArt(&store, &count))
		{		
			store = NULL;
			count = 0;		
		}	
	}

	ASBoolean havePaths = false, havePlugins = false, haveOthers = false, enable = false;

	if (store != NULL)
	{	
		for (ai::int32 index = 0; (result == kNoErr) && (index < count); index++)
		{		
			AIArtHandle art = (*store)[index];
			short artType = kUnknownArt;
			
			if (result == kNoErr)		
				result = sAIArt->GetArtType(art, &artType);
			
			if (artType == kPluginArt)
			{				
				AIPluginGroupHandle artEntry;
				
				if (result == kNoErr)			
					result = sAIPluginGroup->GetPluginArtPluginGroup(art, &artEntry);
				
				if (result == kNoErr)
				{				
					if (artEntry != this->fPluginGroupHandle)
						artType = kUnknownArt;				
				}			
			}
			
			switch (artType) {			
				case kPathArt:
					havePaths = true;
				break;
				case kGroupArt:
				case kCompoundPathArt:
					// nothing
				break;
				case kTextFrameArt:
					havePaths = true;
				break;
				case kPluginArt:
					havePlugins = true;
				break;				
				default:
					haveOthers = true;
				break;			
			}		
		}	
	}
	
	if (result == kNoErr) {	
		if (message->menuItem == this->fMakeMenuHandle) {		
			enable = havePaths && !haveOthers;			
		} 
		else if (message->menuItem == this->fReleaseMenuHandle) {		
			enable = havePlugins;			
		} 
		else if (message->menuItem == this->fAdjustMenuHandle) {		
			enable = havePlugins;			
		} 
		else {		
			// unknown menu item			
		}	
	}

	if (result == kNoErr) {	
		if (enable) {		
			sAIMenu->EnableItem(message->menuItem);		
		} 
		else {		
			sAIMenu->DisableItem(message->menuItem);		
		}	
	}

	if (store != NULL) {	
		sAIMdMemory->MdMemoryDisposeHandle((AIMdMemoryHandle) store);	
	}
	
	return result;	
}

/*
*/
ASErr LiveDropShadowPlugin::GoMenuItem(AIMenuMessage* message) 
{
	ASErr result = kNoErr;

	if (message->menuItem == this->fAboutPluginMenu) {
		// Pop this plug-in's about box.
		SDKAboutPluginsHelper aboutPluginsHelper;
		aboutPluginsHelper.PopAboutBox(message, "About LiveDropShadow", kSDKDefAboutSDKCompanyPluginsAlertString);
		return result;
	}	

	ai::int32 count = 0;
	AIArtHandle **store = NULL;
	if (result == kNoErr)
	{
		if (sAIMatchingArt->GetSelectedArt(&store, &count))
		{
			store = NULL;
			count = 0;
		}
	}

	if (store != NULL) {
		AIArtHandle object = NULL;
		if (message->menuItem == this->fMakeMenuHandle) {		
			this->fLiveDropShadowGroup->StartLensObject(&object, &this->fObjectData);			
		} 
		else if (message->menuItem == this->fAdjustMenuHandle) {		
			this->fLiveDropShadowGroup->FirstLensObject(&object, &this->fObjectData, this->fPluginGroupHandle);
		}
		
		for (ai::int32 index = 0; (result == kNoErr) && (index < count); index++)
		{		
			AIArtHandle art = (*store)[index];
			short artType = kUnknownArt;			
			if (result == kNoErr)
				result = sAIArt->GetArtType(art, &artType);

			if (artType == kPluginArt)
			{
				AIPluginGroupHandle artEntry;

				if (result == kNoErr)
					result = sAIPluginGroup->GetPluginArtPluginGroup(art, &artEntry);
				
				if (result == kNoErr)
				{
					if (artEntry != this->fPluginGroupHandle)
						artType = kUnknownArt;
				}			
			}

			if (((message->menuItem == this->fMakeMenuHandle) && (artType == kPathArt)) || (artType == kTextFrameArt)) {
				if (result == kNoErr) {				
					result = this->fLiveDropShadowGroup->appendPluginGroup(&object, &this->fObjectData, art, this->fPluginGroupHandle);
				}			
			} 
			else if ((message->menuItem == this->fReleaseMenuHandle) && (artType == kPluginArt)) {			
				if (result == kNoErr) {				
					result = this->fLiveDropShadowGroup->ReleaseLensObject(art);				
				}			
			} 
			else if ((message->menuItem == this->fAdjustMenuHandle) && (artType == kPluginArt)) {			
				if (result == kNoErr) {				
					result = this->fLiveDropShadowGroup->ModifyLensObject(&object, &this->fObjectData, art);
				}			
			} 
			else {			
				// unknown menu item or unknown art type			
			}		
		}

		if (message->menuItem == this->fMakeMenuHandle) {		
			this->fLiveDropShadowGroup->FinishLensObject(&object, &this->fObjectData);			
		} 
		else if (message->menuItem == this->fAdjustMenuHandle) {		
			this->fLiveDropShadowGroup->LastLensObject(&object, &this->fObjectData);
		}				
	}

	if (store != NULL) {	
		sAIMdMemory->MdMemoryDisposeHandle((AIMdMemoryHandle) store);	
	}
	
	return result;	
}

// End LiveDropShadowPlugin.cpp
