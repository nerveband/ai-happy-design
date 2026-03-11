//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/StrokeFilter/Source/StrokeFilterPlugin.h $
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

#ifndef __StrokeFilterPlugin_hpp__
#define __StrokeFilterPlugin_hpp__

#include "StrokeFilterID.h"
#include "SDKDef.h"
#include "SDKAboutPluginsHelper.h"
#include "StrokeFilterSuites.h"
#include "StrokeFileterParams.h"

class StrokeFilterFlashController;

/** Maximum length of filter/effect name.
*/
const ASInt32 MAX_FILTER_NAME = 65;

/** Increment used when indexing through the filter/effect strings.
*/
const ASInt32 FILTER_STRING_INC = 2;

const char stokeFilterTitlePStr[2][MAX_FILTER_NAME] = {"Dash" , "Wave"};

/**	Creates a new StrokeFilterPlugin.
	@param pluginRef IN unique reference to this plugin.
	@return pointer to new StrokeFilterPlugin.
*/
Plugin* AllocatePlugin(SPPluginRef pluginRef);

/**	Reloads the StrokeFilterPlugin class state when the plugin is 
	reloaded by the application.
	@param plugin IN pointer to plugin being reloaded.
*/
void FixupReload(Plugin* plugin);


class StrokeFilterPlugin : public Plugin
{
public:
	/**	Constructor.
		@param pluginRef IN reference to this plugin.
	*/
	StrokeFilterPlugin(SPPluginRef pluginRef);

	/**	Destructor.
	*/
	~StrokeFilterPlugin(){}

	//StrokeFilterFlashController* GetStrokeFilterFlashController() const { return strokeFilterFlashController; }

	/**	Restores state of StrokeFilterPlugin during reload.
	*/
	FIXUP_VTABLE_EX(StrokeFilterPlugin, Plugin);

public:	AIFilterHandle GetWaveFilterHandle(void);

public:	AIFilterHandle GetDashFilterHandle(void);
	
protected:

	ASErr SetGlobal(Plugin* plugin);

	ASErr StartupPlugin(SPInterfaceMessage* message);

	/**	Sets global data to NULL and calls Plugin::Shutdown.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr ShutdownPlugin(SPInterfaceMessage* message);

	ASErr ReloadPlugin(SPInterfaceMessage* message);

	ASErr UnloadPlugin(SPInterfaceMessage* message);

	/**	Adds the menu items for this plugin to the application UI.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr AddMenus(SPInterfaceMessage* message);

	/**	Adds the filters provided by this plugin.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr AddFilters(SPInterfaceMessage* message);

	ASErr AddNotifiers(SPInterfaceMessage* message);

	ASErr GoMenuItem(AIMenuMessage* message);

	/**	Gathers parameters from user via a dialog.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr GetFilterParameters(AIFilterMessage* message);

	ASErr GoFilter(AIFilterMessage* message);

	ASErr Notify(AINotifierMessage* message);

	/** Converts given C string into a Pascal string.
	*/
	void CStrToPStr(char* s, ai::UnicodeString::size_type len);

private:
	AIMenuItemHandle menuItem;

	AIMenuItemHandle aboutPluginMenu;

	AIFilterHandle fFilters[kMaxFilters];

	AINotifierHandle fPathUpdateNotifier;

	AINotifierHandle fCSXSPlugPlugSetupCompleteNotifier;

	StrokeFilterFlashController* strokeFilterFlashController;
};

#endif
