//========================================================================================
//  
//  $File$
//
//  $Revision$
//
//  Copyright 2008 Adobe Systems Incorporated. All rights reserved.
//  
//  NOTICE:  Adobe permits you to use, modify, and distribute this file in accordance 
//  with the terms of the Adobe license agreement accompanying it.  If you have received
//  this file from a source other than Adobe, then your use, modification, or 
//  distribution of it requires the prior written permission of Adobe.
//  
//========================================================================================

#ifndef __ANNOTATORPLUGIN_H__
#define __ANNOTATORPLUGIN_H__

#include "AnnotatorID.h"
#include "SDKDef.h"
#include "SDKAboutPluginsHelper.h"
#include "AIAnnotator.h"
#include "AnnotatorSuites.h"
#include "Annotator.h"
#include "SDKErrors.h"


/**	Creates a new AnnotatorPlugin.
	@param pluginRef IN unique reference to this plugin.
	@return pointer to new AnnotatorPlugin.
*/
Plugin* AllocatePlugin(SPPluginRef pluginRef);

/**	Reloads the AnnotatorPlugin class state when the plugin is 
	reloaded by the application.
	@param plugin IN pointer to plugin being reloaded.
*/
void FixupReload(Plugin* plugin);

/** Provides a plug-in which adds a tool to the toolbar, when this tool is
	selected and the cursor is over artwork the artwork is annotated with 
	information displaying the art bounds and the cursor location within the art.
*/
class AnnotatorPlugin : public Plugin
{
private:
	/** Handle for the Annotator tool.
	*/
	AIToolHandle		fToolHandle;

	/** Handle for the About SDK Plug-ins menu item.
	*/
	AIMenuItemHandle	fAboutPluginMenu;

	/** Handle for the annotator added by this plug-in.
	*/
	AIAnnotatorHandle	fAnnotatorHandle;

	/** Handle for the selection changed notifier.
	*/
	AINotifierHandle	fNotifySelectionChanged;

	/** Handle for illustrator shutdown notifier.
	*/
	AINotifierHandle	fShutdownApplicationNotifier;

	/** Pointer to Annotator object.
	*/
	Annotator*			fAnnotator;	
    
    /** Handle for the resource manager used for tool cursor
     */
    AIResourceManagerHandle fResourceManagerHandle;

public:
	/**	Constructor.
		@param pluginRef IN reference to this plugin.
	*/
	AnnotatorPlugin(SPPluginRef pluginRef);

	/**	Destructor.
	*/
	virtual ~AnnotatorPlugin(){}

	/**	Restores state of AnnotatorPlugin during reload.
	*/
	FIXUP_VTABLE_EX(AnnotatorPlugin, Plugin);
	
protected:

	/**	Creates a global instance of the plugin.
		@param plugin IN plugin to create.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr SetGlobal(Plugin* plugin);

	/**	Initializes the plugin.
		@param message IN message sent by the plugin manager.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr StartupPlugin(SPInterfaceMessage* message);

	/**	Creates a new Annotator object.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr PostStartupPlugin();

	/**	Removes the plugin.
		@param message IN message sent by the plugin manager.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr ShutdownPlugin(SPInterfaceMessage* message);

	/** Notify's the plug-in of any operations to perform.
		@param caller IN sender of the message.
		@param selector IN action to take.
		@param message IN pointer to message data.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr Message(char* caller, char* selector, void* message);

	/**	Handler called when a menu item is selected in the UI.
		@param message IN message data.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr GoMenuItem(AIMenuMessage* message);

	/**	Handler called when a menu item should update itself.
		@param message IN message data.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr UpdateMenuItem(AIMenuMessage* message);
	
	/**	Handles notifiers this plugin is registered for when they 
		are received.
		@param message IN message data.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr Notify(AINotifierMessage* message);

	/**	As the cursor moves checks whether any art is 'hit'.
		@param message IN message data.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr TrackToolCursor(AIToolMessage* message);

	/**	Activates this plug-ins annotator.
		@param message IN message data.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr SelectTool(AIToolMessage* message);
	
	/**	Deactivates this plug-ins annotator.
		@param message IN message data.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr DeselectTool(AIToolMessage* message);

	/**	Adds this plug-ins tool to the toolbar.
		@param message IN message data.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr AddTool(SPInterfaceMessage* message);

	/**	Adds this plug-ins annotator to the applications list of
		annotators.
		@param message IN message data.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr AddAnnotator(SPInterfaceMessage* message);

	/**	Registers this plug-in to receive the selection changed
		notifier.
		@param message IN message data.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr AddNotifier(SPInterfaceMessage* message);

	/**	Draws/redraws the annotations when prompted by the application.
		@param message IN message data.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr DrawAnnotation(AIAnnotatorMessage* message);

	/**	Invalidates the annotations within an area of the artwork when
		prompted by the application.
		@param message IN message data.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr InvalAnnotation(AIAnnotatorMessage* message);
};

#endif // __ANNOTATORPLUGIN_H__
