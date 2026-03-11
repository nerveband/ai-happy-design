//========================================================================================
//  
//  $File: //ai/mainline/devtech/sdk/public/samplecode/ScriptMessage/Source/ScriptMessagePlugin.h $
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

#ifndef __ScriptMessagePlugin_h__
#define __ScriptMessagePlugin_h__

#include "Plugin.hpp"
#include "ScriptMessagePanelController.h"
#include "SDKDef.h"
#include "SDKAboutPluginsHelper.h"
#include "ScriptMessageSuites.h"
#include "ScriptMessageID.h"

/**	Creates a new ScriptMessagePlugin.
	@param pluginRef IN unique reference to this plugin.
	@return pointer to new ScriptMessagePlugin.
*/
Plugin* AllocatePlugin(SPPluginRef pluginRef);

/**	Reloads the ScriptMessagePlugin class state when the plugin is 
	reloaded by the application.
	@param plugin IN pointer to plugin being reloaded.
*/
void FixupReload(Plugin* plugin);

/**	Hooks ScriptMessage up as an Illustrator plug-in.
	@ingroup ScriptMessage
*/
class ScriptMessagePlugin : public Plugin
{
public:
	/** Constructor.
		@param pluginRef IN reference to this plugin.
	*/
	ScriptMessagePlugin(SPPluginRef pluginRef);
  
	/**	Restores state of ScriptMessagePlugin during reload.
	*/
	FIXUP_VTABLE_EX(ScriptMessagePlugin, Plugin);
	
	/** Override Message so we can get ScriptMessage messages
	*/
	virtual ASErr Message(char *caller, char *selector, void *message);
	
	virtual ASErr Notify(AINotifierMessage* message);
	
	/**	Calls Plugin::Startup and initialisation functions
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr StartupPlugin(SPInterfaceMessage *message);
	
	/**	Calls Plugin::Shutdown and cleanup functions
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr ShutdownPlugin(SPInterfaceMessage *message);

	/**	Performs actions required for menu item selected.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr GoMenuItem(AIMenuMessage* message);
	
private:
	AINotifierHandle fPlugPlugSetupCompleteNotifier;
	
	AINotifierHandle fBrightnessChangedNotifier;
	
	ScriptMessagePanelController* fPanelController;
	
	/**	Menu item handle for the about plugin dialog.
	*/
	AIMenuItemHandle fAboutPluginMenu;
};

#endif
