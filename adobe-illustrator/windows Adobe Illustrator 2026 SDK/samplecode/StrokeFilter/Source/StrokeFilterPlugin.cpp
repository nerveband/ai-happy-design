//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/StrokeFilter/Source/StrokeFilterPlugin.cpp $
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
#include "AICSXS.h"
#include "StrokeFilterPlugin.h"
#include "StrokeFilterFlashController.h"
#include "StrokeFilters.h"

StrokeFilterPlugin *gPlugin = NULL;


Plugin *AllocatePlugin(SPPluginRef pluginRef)
{
	return new StrokeFilterPlugin(pluginRef);
}

/*
*/
void FixupReload(Plugin* plugin)
{
	StrokeFilterPlugin::FixupVTable((StrokeFilterPlugin*) plugin);
}

/*
*/
StrokeFilterPlugin::StrokeFilterPlugin(SPPluginRef pluginRef)
	: Plugin(pluginRef)
	, menuItem(NULL)
	, aboutPluginMenu(NULL)
	, fPathUpdateNotifier(NULL)
	, strokeFilterFlashController(NULL)
{
	strncpy(fPluginName, kStrokeFilterPluginName, kMaxStringLength);

	for (ASInt32 i = 0; i < kMaxFilters; ++i)
	{
		fFilters[i] = NULL;
	}
}

/*
*/
ASErr StrokeFilterPlugin::SetGlobal(Plugin *plugin)
{
	gPlugin = (StrokeFilterPlugin *) plugin;
	return kNoErr;
}

/*
*/
ASErr StrokeFilterPlugin::StartupPlugin(SPInterfaceMessage *message)
{
	ASErr error = Plugin::StartupPlugin( message );
	if (error) { return error; }

	error = AddMenus(message);
	if (error) { return error; }

	error = AddFilters(message);
	if (error) { return error; }

	error = AddNotifiers(message);
	if (error) { return error; }

	if (NULL == strokeFilterFlashController)
	{
		strokeFilterFlashController = new StrokeFilterFlashController();

		//Lock the plug-in as we register callbacks in PlugPlug Setup.
		error = Plugin::LockPlugin(true);
	}
	return error;
	// If StartupPlugin returns an error instead of kNoErr, Illustrator
	// will try to load this plug-in again later.
}

/*
*/
ASErr StrokeFilterPlugin::ShutdownPlugin(SPInterfaceMessage* message)
{
	message->d.globals = NULL;

	if (strokeFilterFlashController)
	{
		strokeFilterFlashController->RemoveEventListeners();
		delete strokeFilterFlashController;
		strokeFilterFlashController = NULL;
	}
	return Plugin::ShutdownPlugin(message);
}

ASErr StrokeFilterPlugin::ReloadPlugin(SPInterfaceMessage *message)
{
	return Plugin::ReloadPlugin(message);
}

ASErr StrokeFilterPlugin::UnloadPlugin(SPInterfaceMessage *message)
{
	return Plugin::UnloadPlugin(message);
}

/*
*/
ASErr StrokeFilterPlugin::AddMenus(SPInterfaceMessage* message) 
{
	ASErr error = kNoErr;

	// Add About Plugins menu item for this plug-in.
	SDKAboutPluginsHelper aboutPluginsHelper;
	error = aboutPluginsHelper.AddAboutPluginsMenuItem(message, 
				kSDKDefAboutSDKCompanyPluginsGroupName, 
				ai::UnicodeString(kSDKDefAboutSDKCompanyPluginsGroupNameString), 
				"StrokeFilter...", 
				&aboutPluginMenu);

	return error;
}

/*
*/
ASErr StrokeFilterPlugin::AddFilters(SPInterfaceMessage* message)
{
	ASErr error = kNoErr;
	AIMenuItemHandle menuItem = NULL;
	PlatformAddFilterData filterData;
	char categoryPStr[MAX_FILTER_NAME]; /* pascal string */
	char titlePStr[MAX_FILTER_NAME];    /* pascal string */
	char temp[] = kStokeFilterCategoryPStr;
	sprintf(categoryPStr, "%s", temp);
	CStrToPStr(categoryPStr, MAX_FILTER_NAME);
	for (short filterIndex = 0; 
		filterIndex < FILTER_STRING_INC; filterIndex++)
	{
		
		filterData.category = (unsigned char*) categoryPStr;
		sprintf(titlePStr, "%s", stokeFilterTitlePStr[filterIndex]);
		CStrToPStr(titlePStr, MAX_FILTER_NAME);
		filterData.title = (unsigned char*) titlePStr;
		
		error = sAIFilter->AddFilter(message->d.self, titlePStr, &filterData, 
					                 0, &this->fFilters[filterIndex]);
		if (error) { goto errorTag; }

		error = sAIFilter->GetFilterMenuItem(this->fFilters[filterIndex], &menuItem);
		if (error) { goto errorTag; }
		
		error = sAIMenu->UpdateMenuItemAutomatically(menuItem, kAutoEnableMenuItemAction, 0, 0, kIfPath, 0, 0, 0);
		if (error) { goto errorTag; }
	}

errorTag:
	return error;
}

ASErr StrokeFilterPlugin::AddNotifiers(SPInterfaceMessage* message)
{
	ASErr error = kNoErr;
	
	error = sAINotifier->AddNotifier( message->d.self, "StrokeFilterPlugin",
		kStrokeFilterUpdatePathNotifier, &fPathUpdateNotifier);

	if (error) { return error; }

	error = sAINotifier->AddNotifier(message->d.self, "StrokeFilterPlugin " kAICSXSPlugPlugSetupCompleteNotifier,
		kAICSXSPlugPlugSetupCompleteNotifier, &fCSXSPlugPlugSetupCompleteNotifier);

	return error;
}

ASErr StrokeFilterPlugin::GoMenuItem( AIMenuMessage *message )
{
	try 
	{
		if (message->menuItem == aboutPluginMenu)
		{
			SDKAboutPluginsHelper aboutPluginsHelper;
			aboutPluginsHelper.PopAboutBox(message, "About StrokeFilter", kSDKDefAboutSDKCompanyPluginsAlertString);
		}
	}
	catch(...)
	{
		return 'ESHD';
	}
	
	return kNoErr;
}

/*
*/
ASErr StrokeFilterPlugin::GetFilterParameters(AIFilterMessage* message)
{
	ASErr error;
	if (message->filter == this->fFilters[0])
	{
		error = strokeFilterFlashController->DoFlashExtension(message, StrokeFilterFlashController::kDashStroke);
	}
	else if (message->filter == this->fFilters[1])
	{
		error = strokeFilterFlashController->DoFlashExtension(message, StrokeFilterFlashController::kWaveStroke);
	}
	return error;
}

AIFilterHandle StrokeFilterPlugin::GetWaveFilterHandle()
{
	return this->fFilters[1];
}

AIFilterHandle StrokeFilterPlugin::GetDashFilterHandle()
{
	return this->fFilters[0];
}

/*
*/
ASErr StrokeFilterPlugin::GoFilter(AIFilterMessage* message)
{
	ASErr error = kNoErr;
	if(message->parameters != NULL)
	{
		StrokeFilters strokeFilters;
		if (message->filter == this->GetWaveFilterHandle())
		{
			StrokeFileterParamsWave* params = (StrokeFileterParamsWave*) message->parameters;
			error = strokeFilters.GoWavedStroke(*params);
		}
		else if( message->filter == this->GetDashFilterHandle())
		{
			StrokeFileterParamsDash* params = (StrokeFileterParamsDash*) message->parameters;
			error = strokeFilters.GoDashedStroke(*params);
		}
	}
	return error;
}
/*
*/
ASErr StrokeFilterPlugin::Notify(AINotifierMessage* message)
{
	ASErr error = kNoErr;
	
	if (strcmp(message->type, kStrokeFilterUpdatePathNotifier) == 0)
	{
		StrokeFilters strokeFilters;
		if (strokeFilterFlashController->GetFilterType() == StrokeFilterFlashController::kDashStroke) 
		{
			error = strokeFilters.GoDashedStroke(strokeFilterFlashController->GetDashParams());
		}
		else if (strokeFilterFlashController->GetFilterType() == StrokeFilterFlashController::kWaveStroke) 
		{
			error = strokeFilters.GoWavedStroke(strokeFilterFlashController->GetWaveParams());
		}
	}
	else if ( message->notifier == fCSXSPlugPlugSetupCompleteNotifier )
	{
		strokeFilterFlashController->RegisterCSXSEventListeners();
	}
	return error;
}

/*
*/
void StrokeFilterPlugin::CStrToPStr(char *s, ai::UnicodeString::size_type len)
{
	const ai::UnicodeString sAsUnicode((const char*)s);
	ai::PStr sAsPStr((unsigned char*) s);
	sAsUnicode.getToBuffer(sAsPStr, len, kAIUTF8CharacterEncoding );
}
