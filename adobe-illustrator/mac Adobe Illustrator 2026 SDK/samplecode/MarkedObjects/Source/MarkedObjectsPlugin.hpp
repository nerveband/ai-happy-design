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

#ifndef __MarkedObjectsPlugin_hpp__
#define __MarkedObjectsPlugin_hpp__

#include "MarkedObjectManager.hpp"
#include "MarkedObjectsPreferences.hpp"
#include "MarkedObjectID.h"
#include "SDKDef.h"
#include "SDKAboutPluginsHelper.h"

/**	Creates a new MarkedObjectsPlugin.
	@param pluginRef IN unique reference to this plugin.
	@return pointer to new MarkedObjectsPlugin.
*/
Plugin* AllocatePlugin(SPPluginRef pluginRef);

/**	Reloads the MarkedObjectsPlugin class state when the plugin is 
	reloaded by the application.
	@param plugin IN pointer to plugin being reloaded.
*/
void FixupReload(Plugin* plugin);

class MarkedObjectsPreferences;
class MarkedObjectManager;
class MarkedObjectsPanelController;

class MarkedObjectsPlugin : public Plugin
{
public:
	/**	Constructor.
		@param pluginRef IN reference to this plugin.
	*/
	MarkedObjectsPlugin(SPPluginRef pluginRef);

	/**	Destructor.
	*/
	virtual ~MarkedObjectsPlugin(){}

	ASErr HideMarkedObjectsLayer( AIBoolean hide );

	ASErr SelectMarkedObject(int index);

	void GetMarkedObjects(void);
	
	void WriteMarkedPoints( const ai::FilePath & fileSpec );

	void UpdateMarkedObjectString( int itemSelected, const string & newString );

	void ResetIDs( void );
	
	void SelectOurTool( void );
	
	void UnSelectAllArt( void );

	void ResetInsertionPoint( void );
	
	void SetDialogPreferences( void );
	
	void DoModalPrefs( void );

	ASBoolean IsMarkedObjectSelected( void );

	/**	Restores state of MarkedObjectsPlugin during reload.
	*/
	FIXUP_VTABLE_EX(MarkedObjectsPlugin, Plugin);
	
protected:

	virtual ASErr SetGlobal(Plugin * plugin);

	virtual ASErr StartupPlugin( SPInterfaceMessage * message );

	virtual ASErr ShutdownPlugin( SPInterfaceMessage * message );

	virtual ASErr PostStartupPlugin( void );

	virtual ASErr GoMenuItem( AIMenuMessage * message );

	virtual ASErr UpdateMenuItem( AIMenuMessage * message );
	
	virtual ASErr Notify( AINotifierMessage * message );
	
	// add tool handlers
	virtual ASErr TrackToolCursor( AIToolMessage * message );

	virtual ASErr ToolMouseUp( AIToolMessage * message );

	virtual ASErr DeselectTool( AIToolMessage * message );

	virtual ASErr GoTimer(AITimerMessage * message );

private:
	AIMenuItemHandle			menuItem;
	/**	Pointer to instance of TransformButtonsDialog.
	*/
	MarkedObjectsPanelController* fMarkedObjectsPanelController;

	AINotifierHandle			artSelectionChangedNotifier;
	AINotifierHandle			brightnessChangedNotifier;
	AINotifierHandle			documentChangedNotifier;
	AINotifierHandle			documentSaveNotifierBefore;
	AINotifierHandle			documentSaveAsNotifierBefore;
	AINotifierHandle			documentSaveNotifierAfter;
	AINotifierHandle			documentSaveAsNotifierAfter;
	AIToolHandle				toolHandle;
	MarkedObjectManager *		objectManager;
	MarkedObjectsPreferences *	preferences;
	AIMenuItemHandle			aboutPluginMenu;

	AINotifierHandle			fRegisterEventNotifierHandle;
	AINotifierHandle			fHideLayerNotifier;
	AINotifierHandle			fResetNumberIDNotifier;
	AINotifierHandle			fSaveObjectNotifier;
	AINotifierHandle			fPreferencesSetNotifier;
	AINotifierHandle			fObjectSelectedNotifier;
	AINotifierHandle			fShutdownApplicationNotifier;
    AIResourceManagerHandle     fResourceManagerHandle;

	ASErr AcquirePostStartupSuites( void );

	ASErr ReleasePostStartupSuites( void );
};

#endif
