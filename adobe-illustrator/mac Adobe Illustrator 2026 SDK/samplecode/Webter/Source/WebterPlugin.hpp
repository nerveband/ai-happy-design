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

#ifndef __WebterPlugin_hpp__
#define __WebterPlugin_hpp__

#include "Plugin.hpp"
#include "WebterID.h"

class WebterPanelController;

static const int kToolStrings 	= 16100;
static const int kToolIcons 	= 16050;
static const int kCursorID		= 16000;
static const int kNoIconErr		= '!ico';

/**	Creates a new WebterPlugin.
	@param pluginRef IN unique reference to this plugin.
	@return pointer to new WebterPlugin.
*/
Plugin* AllocatePlugin(SPPluginRef pluginRef);

/**	Reloads the WebterPlugin class state when the plugin is 
	reloaded by the application.
	@param plugin IN pointer to plugin being reloaded.
*/
void FixupReload(Plugin* plugin);

class WebterPlugin : public Plugin
{
public:
	/**	Constructor.
		@param pluginRef IN reference to this plugin.
	*/
	WebterPlugin(SPPluginRef pluginRef);

	/**	Destructor.
	*/
	virtual ~WebterPlugin(){}

	/**	Restores state of WebterPlugin during reload.
	*/
	FIXUP_VTABLE_EX(WebterPlugin, Plugin);

protected:

	/** Calls Plugin::Message and handles any errors returned.
		@param caller IN sender of the message.
		@param selector IN nature of the message.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr Message(char* caller, char* selector, void *message);

	virtual ASErr SetGlobal(Plugin* plugin);

	/**	Calls Plugin::Startup and initialisation functions, such as 
		AddMenus and AddNotifiers.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr StartupPlugin( SPInterfaceMessage *message );

	/**	Sets global data to NULL and calls Plugin::Shutdown.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr ShutdownPlugin( SPInterfaceMessage *message );

	/**	Performs actions required for menu item selected.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr GoMenuItem( AIMenuMessage *message );

	/**	Updates state of menu items.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr UpdateMenuItem( AIMenuMessage *message );
	
	/** Deletes the WebterDialog object after the application 
		shutdown notifier is received.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr Notify( AINotifierMessage *message );

private:	
	AIMenuItemHandle		fMenuItem;
	AIToolHandle			fTool[3];
	AIRealPoint				fStartingPoint, fEndPoint;
	AIArtHandle				fArtHandle;
	AINotifierHandle		fArtSelectionChangedNotifier;
	AINotifierHandle		fDocumentChangedNotifier;
	AINotifierHandle		fDocumentClosedNotifier;
	AIMenuItemHandle		fAboutPluginMenu;
	WebterPanelController	*fWebterPanelController;
	AINotifierHandle		fRegisterEventNotifierHandle;
};

#endif
