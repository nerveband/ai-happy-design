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

#ifndef __TWIRLFILTERPLUGIN_H__
#define __TWIRLFILTERPLUGIN_H__

#include "IllustratorSDK.h"
#include "Plugin.hpp"
#include "SDKDef.h"
#include "SDKAboutPluginsHelper.h"
#include "TwirlFilterID.h"
#include "TwirlFilterSuites.h"
#include "TwirlFilter.h"

/** Maximum length of filter/effect name.
*/
const ASInt32 MAX_FILTER_NAME = 65;

/** Increment used when indexing through the filter/effect strings.
*/
const ASInt32 FILTER_STRING_INC = 2;

/**	Creates a new TwirlFilterPlugin.
	@param pluginRef IN unique reference to this plugin.
	@return pointer to new TwirlFilterPlugin.
*/
Plugin* AllocatePlugin(SPPluginRef pluginRef);

/**	Reloads the TwirlFilterPlugin class state when the plugin is 
	reloaded by the application.
	@param plugin IN pointer to plugin being reloaded.
*/
void FixupReload(Plugin* plugin);

/**	Provides a plugin which demonstrates adding a filter and a live effect.
*/
class TwirlFilterPlugin : public Plugin
{
public:
	/**	Constructor.
		@param pluginRef IN reference to this plugin.
	*/
	TwirlFilterPlugin(SPPluginRef pluginRef);

	/**	Destructor.
	*/
	virtual ~TwirlFilterPlugin(){}

	/**	Restores state of TwirlFilterPlugin during reload.
	*/
	FIXUP_VTABLE_EX(TwirlFilterPlugin, Plugin);

protected:
	/** Calls Plugin::Message and handles any errors returned.
		@param caller IN sender of the message.
		@param selector IN nature of the message.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr Message(char* caller, char* selector, void *message);

	/**	Calls Plugin::Startup and initialisation functions, such as 
		AddMenus and AddNotifiers.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr StartupPlugin(SPInterfaceMessage* message);	

	/**	Performs plugin tasks that could not be performed until
		the application was started.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr PostStartupPlugin();

	/** Deletes the TwirlFilter object after the application 
		shutdown notifier is received.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr Notify(AINotifierMessage* message);

	/**	Performs actions required for menu item selected.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr GoMenuItem(AIMenuMessage* message);
	
	/**	Applies the filter using the parameters provided by the user.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr GoFilter(AIFilterMessage* message);

	/**	Applies the live effect using the parameters provided by the user.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr GoLiveEffect(AILiveEffectGoMessage* message);

	/**	Gathers parameters from user via a dialog.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr GetFilterParameters(AIFilterMessage* message);

	/**	Gathers parameters from user via a dialog.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr EditLiveEffectParameters(AILiveEffectEditParamMessage* message);

private:
	/** Pointer to a TwirlFilter object.
	*/
	TwirlFilter* fTwirlFilter;

	/** Menu item handle for this plug-ins About menu.
	*/
	AIMenuItemHandle fAboutPluginMenu;

	/** Handle for the application shutdown notifier.
	*/
	AINotifierHandle fAppShutdownNotifier;

	/** Handle for CSXS plugplug setup complete notifier.
	*/
	AINotifierHandle fCSXSPlugPlugSetupCompleteNotifier;

	/** Handle for apply filter notifier
	*/
	AINotifierHandle fApplyFilterNotifier;

	/** Handle for apply live effect notifier
	*/
	AINotifierHandle fApplyLiveEffectNotifier;

	/** Array of AIFilterHandles to store references to this plugins filters.
	*/
	AIFilterHandle fFilters[kMaxFilters];

	/** Array of AILiveEffectHandles to store references to this plugins live effects.
	*/
	AILiveEffectHandle fEffects[kMaxEffects];

	/** Stores the number of filters added by this plugin.
	*/
	ASInt32 fNumFilters;

	/** Stores the number of effects added by this plugin.
	*/
	ASInt32 fNumEffects;

	/**	Adds the filters provided by this plugin.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr AddFilters(SPInterfaceMessage* message);

	/**	Adds the menu items for this plugin to the application UI.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr AddMenus(SPInterfaceMessage* message);

	/**	Adds the live effects provided by this plug-in.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr AddLiveEffects(SPInterfaceMessage* message);

	/**	Registers this plugin to receive the application shutdown notifier.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr AddNotifiers(SPInterfaceMessage* message);	
	
	/** Converts given C string into a Pascal string.
	*/
	void CStrToPStr(char *s, ai::UnicodeString::size_type len);
};

#endif // End TwirlFilterPlugin.h
