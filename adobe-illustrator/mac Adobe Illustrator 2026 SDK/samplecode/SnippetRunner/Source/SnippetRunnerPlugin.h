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

#ifndef __SnippetRunnerPlugin_h__
#define __SnippetRunnerPlugin_h__

#include "Plugin.hpp"

class SnippetRunnerPanelController;

/**	Creates a new SnippetRunnerPlugin.
	@param pluginRef IN unique reference to this plugin.
	@return pointer to new SnippetRunnerPlugin.
*/
Plugin* AllocatePlugin(SPPluginRef pluginRef);

/**	Reloads the SnippetRunnerPlugin class state when the plugin is 
	reloaded by the application.
	@param plugin IN pointer to plugin being reloaded.
*/
void FixupReload(Plugin* plugin);

/**	Hooks SnippetRunner up as an Illustrator plug-in.
	@ingroup snippetrunner
*/
class SnippetRunnerPlugin : public Plugin
{
private:
	/** Handle for menu item in Help > About SDK Plug-ins.
	*/
	AIMenuItemHandle			fAboutPluginMenu;	

	/** Handle for menu item in Window > SDK.
	*/
	AIMenuItemHandle			fShowHidePanelMenu;	

	/** Pointer to instance of SnippetRunnerPanelController
	*/
	SnippetRunnerPanelController*	fSnippetRunnerPanelController;
	
	/** Handle for app shutting down notifier.
	*/
	AINotifierHandle			fAppShutdownNotifier;

	/** Handle for art selection changed notifier.
	*/
	AINotifierHandle			fArtSelectionChangedNotifier;

	/** Handle for CSXS plugplug setup complete notifier.
	*/
	AINotifierHandle				fCSXSPlugPlugSetupCompleteNotifier;
	
	/** Add menu item associated with this plugin.
		@param message IN contains plugin information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr AddMenus(SPInterfaceMessage* message);

	/** Check if a menu group has already been created.
		@param targetGroupName IN name of the menu group to look for.
		@param groupAlreadyMade OUT true if group already exists, false if does not exist.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr MenuGroupExists(const char* targetGroupName, bool& groupAlreadyMade);

	/** Register notifiers this plugin wants to receive.
		@param message IN contains plugin information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr AddNotifiers(SPInterfaceMessage* message);

public:
	/** Constructor.
		@param pluginRef IN reference to this plugin.
	*/
	SnippetRunnerPlugin(SPPluginRef pluginRef);

	/** Destructor.
	*/
	virtual ~SnippetRunnerPlugin();

	/** Notify's the panel that the log has changed.
	*/
	void NotifyLogChanged();

	/** Notify's the panel that the unit testing has ended.
	*/
	void NotifyEndUnitTest();

	/**	Restores state of SnippetRunnerPlugin during reload.
	*/
	FIXUP_VTABLE_EX(SnippetRunnerPlugin, Plugin);

protected:
	/** Notify's the panel that the log has changed.
		@param caller IN sender of the message.
		@param selector IN action to take.
		@param message IN pointer to data required for action.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr Message(char *caller, char *selector, void *message);

	/**	Creates a global instance of the plugin.
		@param plugin IN plugin to create.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr SetGlobal(Plugin * plugin);

	/**	Handler called when a menu item is selected in the UI.
		@param message IN details of menu item.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr GoMenuItem( AIMenuMessage * message );

	/**	Handler called when a menu item should update itself.
		@param message IN details of the menu item.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr UpdateMenuItem( AIMenuMessage * message );
	
	/**	Initializes the plugin.
		@param message IN message sent by the plugin manager.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr StartupPlugin( SPInterfaceMessage * message );

	/**	Removes the plugin.
		@param message IN message sent by the plugin manager.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr ShutdownPlugin( SPInterfaceMessage * message );

	/** Performs plugin tasks that could not be performed until
		the application was started.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr PostStartupPlugin();

	/**	Handles notifiers this plugin is registered for when they 
		are received.
		@param message IN contains details of notifier.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr Notify( AINotifierMessage * message );

	/**	Acquires non-essential suites and suites that could not be 
		acquired until the application was started.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr AcquirePostStartupSuites();

	/**	Releases the suites that were aquired post startup.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr ReleasePostStartupSuites();

	/**	Start the unit test timer.
		@param message IN timer details.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr GoTimer(AITimerMessage* message);

};

#endif

// End SnippetRunnerPlugin.h
