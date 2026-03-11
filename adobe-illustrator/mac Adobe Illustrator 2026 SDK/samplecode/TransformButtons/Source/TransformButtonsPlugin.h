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

#ifndef __TRANSFORMBUTTONSPLUGIN_H__
#define __TRANSFORMBUTTONSPLUGIN_H__

#include "IllustratorSDK.h"
#include "Plugin.hpp"
#include "SDKDef.h"
#include "SDKAboutPluginsHelper.h"
#include "TransformButtonsSuites.h"
class TransformButtonsPanelController;

/**	Creates a new TransformButtonsPlugin.
	@param pluginRef IN unique reference to this plugin.
	@return pointer to new TransformButtonsPlugin.
*/
Plugin* AllocatePlugin(SPPluginRef pluginRef);

/**	Reloads the TransformButtonsPlugin class state when the plugin is 
	reloaded by the application.
	@param plugin IN pointer to plugin being reloaded.
*/
void FixupReload(Plugin* plugin);

/**	Provides a plugin which demonstrates applying one-click transforms to art.
*/
class TransformButtonsPlugin : public Plugin
{
public:
	/**	Constructor.
		@param pluginRef IN reference to this plugin.
	*/
	TransformButtonsPlugin(SPPluginRef pluginRef);

	/**	Destructor.
	*/
	virtual ~TransformButtonsPlugin(){}

	/**	Restores state of TransformButtonsPlugin during reload.
	*/
	FIXUP_VTABLE_EX(TransformButtonsPlugin, Plugin);

protected:
	/** Calls Plugin::Message and handles any errors returned.
		@param caller IN sender of the message.
		@param selector IN nature of the message.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr Message(char* caller, char* selector, void *message);

	/**	Calls Plugin::Startup and initialisation functions, such as 
		AddMenus.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr SetGlobal(Plugin* plugin);

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

	/** Performs plugin tasks that could not be performed until
		the application was started.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr PostStartupPlugin();

	/** Deletes the TransformButtonsDialog object after the application 
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

private:
	/**	Pointer to instance of TransformButtonsDialog.
	*/
	TransformButtonsPanelController* fTransformButtonsPanelController;

	/**	Menu item handle to show/hide the TransformButtons dialog.
	*/
	AIMenuItemHandle fShowTransformButtonsDialog;

	/** Menu item handle for the TransformButtons About menu item.
	*/
	AIMenuItemHandle fAboutPluginMenu;

	AINotifierHandle			fRegisterEventNotifierHandle;

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

#endif // End TransformButtonsPlugin.h
