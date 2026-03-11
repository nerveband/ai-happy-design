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
#include "TwirlFilterPlugin.h"
#include "AICSXS.h"

#define CHECK_RESULT if(result) { return result;}

/*
*/
Plugin* AllocatePlugin(SPPluginRef pluginRef)
{
	return new TwirlFilterPlugin(pluginRef);
}

/*
*/
void FixupReload(Plugin* plugin)
{
	TwirlFilterPlugin::FixupVTable((TwirlFilterPlugin*) plugin);
}

/*
*/
TwirlFilterPlugin::TwirlFilterPlugin(SPPluginRef pluginRef) 
	: Plugin(pluginRef),
	fAboutPluginMenu(NULL),
	fAppShutdownNotifier(NULL),
	fCSXSPlugPlugSetupCompleteNotifier(NULL),
	fNumFilters(0),
	fTwirlFilter(NULL)
{
	strncpy(fPluginName, kTwirlFilterPluginName, kMaxStringLength);
}

/*
*/
ASErr TwirlFilterPlugin::Message(char* caller, char* selector, void *message) 
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
ASErr TwirlFilterPlugin::StartupPlugin(SPInterfaceMessage* message)
{
	ASErr result = kNoErr;
	result = Plugin::StartupPlugin(message);
	CHECK_RESULT
	result = this->AddMenus(message);
    CHECK_RESULT
	result = this->AddFilters(message);
    CHECK_RESULT
	result = this->AddLiveEffects(message);
    CHECK_RESULT
	result = this->AddNotifiers(message);
	
	return result;
}

/*
*/
ASErr TwirlFilterPlugin::PostStartupPlugin()
{
	ASErr result = kNoErr;

	// Create LiveDropShadowGroup
	if (this->fTwirlFilter == NULL)
	{
		this->fTwirlFilter = new TwirlFilter();
		// TODO error handling
	}

	return result;
}

/*
*/
ASErr TwirlFilterPlugin::Notify(AINotifierMessage* message)
{
	ASErr result = kNoErr;
	if (message->notifier == this->fAppShutdownNotifier)
	{
		if (this->fTwirlFilter)
		{
			delete this->fTwirlFilter;
		}
        Plugin::LockPlugin(false);
	}
	else if ( message->notifier == fApplyFilterNotifier )
	{
        GoFilter((AIFilterMessage *)message);
	}
	
	else if ( message->notifier == fApplyLiveEffectNotifier )
	{
        GoLiveEffect((AILiveEffectGoMessage *) message);
	}
	
	return result;
}

/*
*/
ASErr TwirlFilterPlugin::AddNotifiers(SPInterfaceMessage* message)
{
    ASErr result = kNoErr;

	result = sAINotifier->AddNotifier(message->d.self, "TwirlFilter " kAIApplicationShutdownNotifier, 
		kAIApplicationShutdownNotifier, &this->fAppShutdownNotifier);

	result = sAINotifier->AddNotifier(message->d.self, "TwirlFilter " kAICSXSPlugPlugSetupCompleteNotifier,
		kAICSXSPlugPlugSetupCompleteNotifier, &fCSXSPlugPlugSetupCompleteNotifier);

	result = sAINotifier->AddNotifier(message->d.self, kTwirlFilterApplyCaller, kTwirlApplyFilterNotifier,
		&fApplyFilterNotifier);
	
	
	result = sAINotifier->AddNotifier(message->d.self, kTwirlFilterApplyCaller, kTwirlApplyEffectNotifier,
		&fApplyLiveEffectNotifier);
	

	return result;
}

/*
*/
ASErr TwirlFilterPlugin::AddMenus(SPInterfaceMessage* message) 
{
	ASErr error = kNoErr;

	// Add a menu item to the About SDK Plug-ins menu group.
	SDKAboutPluginsHelper aboutPluginsHelper;
	error = aboutPluginsHelper.AddAboutPluginsMenuItem(message, 
				kSDKDefAboutSDKCompanyPluginsGroupName, 
				ai::UnicodeString(kSDKDefAboutSDKCompanyPluginsGroupNameString), 
				"TwirlFilter...", 
				&this->fAboutPluginMenu);

	return error;
}

/*
*/
ASErr TwirlFilterPlugin::AddFilters(SPInterfaceMessage* message)
{
	ASErr error = kNoErr;

	AIMenuItemHandle menuItem;
	PlatformAddFilterData filterData;
	short	filterIndex = 0;

	filterData.category = (unsigned char*) " SDK";
	filterData.title = (unsigned char*) " Twirl...";
	char filterName[] = "Twirl Filter";
	error = sAIFilter->AddFilter(message->d.self, filterName, &filterData, 
				0, &this->fFilters[filterIndex]);
	if (error) 
		goto errorTag;

	error = sAIFilter->GetFilterMenuItem(this->fFilters[filterIndex], &menuItem);
	if (error) 
		goto errorTag;
	
	error = sAIMenu->UpdateMenuItemAutomatically(menuItem, kAutoEnableMenuItemAction, 0, 0,
		kIfPath, 0, 0, 0);
	if (error) 
		goto errorTag;
	

	this->fNumFilters = filterIndex;

errorTag:
	return error;
}

/*
*/

ASErr TwirlFilterPlugin::AddLiveEffects(SPInterfaceMessage* message)
{
	ASErr error = kNoErr;
	AILiveEffectData effectData;
	AddLiveEffectMenuData menuData;
	AIMenuItemHandle menuHandle = nullptr;
	short filterIndex = 0;

    char categoryStr[] = " SDK";
	char nameStr[] = " Twirl...";

	effectData.self = message->d.self;
	effectData.name = nameStr;				// Localized name of the effect
	effectData.title = nameStr;
	effectData.majorVersion = 1;
	effectData.minorVersion = 0;
	effectData.prefersAsInput = (kGroupInputArt | kPathInputArt | kCompoundPathInputArt);//kAnyInputArt;	// AIStyleFilterPreferedInputArtType
	effectData.styleFilterFlags = kPostEffectFilter;	// AIStyleFilterFlags

	error = sAILiveEffect->AddLiveEffect(&effectData, &this->fEffects[filterIndex]);
	if (error) 
		goto errorTag;

	menuData.category = categoryStr;
	menuData.title = nameStr;
	menuData.options = 0;

	error = sAILiveEffect->AddLiveEffectMenuItem(this->fEffects[filterIndex],
												nameStr,
												&menuData,
												&menuHandle,
												NULL);

	if (error) 
		goto errorTag;
	error = sAIMenu->UpdateMenuItemAutomatically(menuHandle, kAutoEnableMenuItemAction, 0, 0,
		kIfPath, 0, 0, 0);
	if (error) 
		goto errorTag;
	if (error) 
		goto errorTag;
    

	this->fNumEffects = filterIndex;

errorTag:

	return error;
}

/*
*/
ASErr TwirlFilterPlugin::GoMenuItem(AIMenuMessage* message) 
{
	ASErr error = kNoErr;
	if (message->menuItem == this->fAboutPluginMenu) {
		// Pop this plug-in's about box.
		SDKAboutPluginsHelper aboutPluginsHelper;
		aboutPluginsHelper.PopAboutBox(message, "About TwirlFilter", kSDKDefAboutSDKCompanyPluginsAlertString);
	}	
	return error;
}

ASErr TwirlFilterPlugin::GoLiveEffect(AILiveEffectGoMessage* message)
{
	return fTwirlFilter->GoTwirl(message);
}


ASErr TwirlFilterPlugin::GoFilter(AIFilterMessage* message)
{
	return fTwirlFilter->GoTwirl(message);
}

/*
*/
ASErr TwirlFilterPlugin::GetFilterParameters(AIFilterMessage* message)
{
	ASErr error = kNoErr;

	if (message->filter == this->fFilters[0])
	{
		// Let's get parameters that we can (e.g. centre position, diagonal, etc.)
		error = fTwirlFilter->TwirlGetParms(message);
	}
	
	return error;
}


ASErr TwirlFilterPlugin::EditLiveEffectParameters(AILiveEffectEditParamMessage* message)
{
	ASErr error = kNoErr;
		
	if (message->effect == this->fEffects[0])
	{
		// Let's get parameters that we can (e.g. centre position, diagonal, etc.)
		error = this->fTwirlFilter->TwirlGetParms(message);
	}
	
	return error;
}


/////////////////////////////////////////
//
// Utility functions
//
/////////////////////////////////////////

/*
*/
void TwirlFilterPlugin::CStrToPStr(char *s, ai::UnicodeString::size_type len)
{
	const ai::UnicodeString sAsUnicode((const char*)s);
	ai::PStr sAsPStr((unsigned char*) s);
	sAsUnicode.getToBuffer(sAsPStr, len, kAIUTF8CharacterEncoding );
}

// End TwirlFilterPlugin.cpp
