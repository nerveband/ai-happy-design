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

#ifndef __LiveDropShadowPlugin_h__
#define __LiveDropShadowPlugin_h__

#include "LiveDropShadowGroup.h"
#include "Plugin.hpp"
#include "LiveDropShadowSuites.h"
#include "LiveDropShadowID.h"
#include "AIMenuGroups.h"
#include "SDKDef.h"
#include "SDKAboutPluginsHelper.h"

/**	Creates a new LiveDropShadowPlugin.
	@param pluginRef IN unique reference to this plugin.
	@return pointer to new LiveDropShadowPlugin.
*/
Plugin* AllocatePlugin(SPPluginRef pluginRef);

/**	Reloads the LiveDropShadowPlugin class state when the plugin is 
	reloaded by the application.
	@param plugin IN pointer to plugin being reloaded.
*/
void FixupReload(Plugin* plugin);

/**	Provides a plugin which demonstrates adding a drop shadow to path items.
*/
class LiveDropShadowPlugin : public Plugin
{
public:
	/**	Constructor.
		@param pluginRef IN reference to this plugin.
	*/
	LiveDropShadowPlugin(SPPluginRef pluginRef);

	/**	Destructor.
	*/
	virtual ~LiveDropShadowPlugin(){}

	/**	Restores state of LiveDropShadowPlugin during reload.
	*/
	FIXUP_VTABLE_EX(LiveDropShadowPlugin, Plugin);

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

	/**	Sets global data to NULL and calls Plugin::Shutdown.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr ShutdownPlugin(SPInterfaceMessage* message);

	/**	Performs plugin tasks that could not be performed until
		the application was started.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr PostStartupPlugin();

	/** Deletes the LiveDropShadowGroup object after the application 
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

	/**	Updates state of menu items.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr UpdateMenuItem(AIMenuMessage* message);

	/**	Updates the plugin group item.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr PluginGroupUpdate(AIPluginGroupMessage* message);

private:
	/** Menu item handle for the about plugin dialog.
	*/
	AIMenuItemHandle fAboutPluginMenu;

	/**	Menu item handle for adding a drop shadow.
	*/
	AIMenuItemHandle fMakeMenuHandle;

	/**	Menu item handle for removing a drop shadow.
	*/
	AIMenuItemHandle fReleaseMenuHandle;

	/** Menu item handle for setting the data store count and range.
	*/
	AIMenuItemHandle fAdjustMenuHandle;

	/** Plugin group handle.
	*/
	AIPluginGroupHandle	fPluginGroupHandle;

	/** LiveDropShadowGroup object.
	*/
	LiveDropShadowGroup* fLiveDropShadowGroup;

	/** ObjectData object used by LiveDropShadowGroup.
	*/
	LiveDropShadowGroup::ObjectData fObjectData;

	/** Application shutdown notifier handler.
	*/
	AINotifierHandle fAppShutdownNotifier;
	
	/**	Adds the menu items for this plugin to the application UI.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr AddMenus(SPInterfaceMessage* message);

	/**	Registers this plugin to receive the application shutdown notifier.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr AddNotifiers(SPInterfaceMessage* message);

};

#endif // End LiveDropShadowPlugin.h
