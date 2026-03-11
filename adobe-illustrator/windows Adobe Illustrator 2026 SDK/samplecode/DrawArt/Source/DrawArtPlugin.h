//========================================================================================
//  
//  $File: //ai/mainline/devtech/sdk/public/samplecode/DrawArtHtml/Source/DrawArtHtmlPlugin.h $
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

#ifndef __DrawArtHtmlPlugin_h__
#define __DrawArtHtmlPlugin_h__

#include "Plugin.hpp"
#include "SDKDef.h"
#include "SDKAboutPluginsHelper.h"
#include "DrawArtSuites.h"
#include "DrawArtID.h"

class DrawArtPanelController;

/**	Creates a new DrawArtHtmlPlugin.
	@param pluginRef IN unique reference to this plugin.
	@return pointer to new DrawArtHtmlPlugin.
*/
Plugin* AllocatePlugin(SPPluginRef pluginRef);

/**	Reloads the DrawArtHtmlPlugin class state when the plugin is 
	reloaded by the application.
	@param plugin IN pointer to plugin being reloaded.
*/
void FixupReload(Plugin* plugin);

/**	Hooks DrawArtHtml up as an Illustrator plug-in.
	@ingroup DrawArtHtml
*/
class DrawArtPlugin : public Plugin
{
public:
	/** Constructor.
		@param pluginRef IN reference to this plugin.
	*/
	DrawArtPlugin(SPPluginRef pluginRef);
  
	/**	Restores state of DrawArtHtmlPlugin during reload.
	*/
	FIXUP_VTABLE_EX(DrawArtPlugin, Plugin);
	
	/** Override Message so we can get DrawArtHtml messages
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
	
	public:	ASErr DrawAGMPort(std::string& encodeBitmap);
	
	ASErr SetGlobal(Plugin *plugin);
	
private:
	AINotifierHandle fPlugPlugSetupCompleteNotifier;
	
	DrawArtPanelController* fPanelController;
	
	/**	Menu item handle for the about plugin dialog.
	*/
	AIMenuItemHandle fAboutPluginMenu;
};

#endif
