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

#ifndef __TEXTFILEFORMATPLUGIN_H__
#define __TEXTFILEFORMATPLUGIN_H__

#include "IllustratorSDK.h"
#include "Plugin.hpp"
#include "TextFileFormatID.h"
#include "TextFileFormatSuites.h"
#include "SDKDef.h"
#include "SDKAboutPluginsHelper.h"

#define kMaxStringLength 256

/**	Creates a new TextFileFormatPlugin.
	@param pluginRef IN unique reference to this plugin.
	@return pointer to new TextFileFormatPlugin.
*/
Plugin* AllocatePlugin(SPPluginRef pluginRef);

/**	Reloads the TextFileFormatPlugin class state when the plugin is 
	reloaded by the application.
	@param plugin IN pointer to plugin being reloaded.
*/
void FixupReload(Plugin* plugin);

/**	Provides a plugin which demonstrates adding new file formats to open and save to.
*/
class TextFileFormatPlugin : public Plugin
{
public:
	/**	Constructor.
		@param pluginRef IN reference to this plugin.
	*/
	TextFileFormatPlugin(SPPluginRef pluginRef);

	/**	Destructor.
	*/
	virtual ~TextFileFormatPlugin(){}

	/**	Restores state of TextFileFormatPlugin during reload.
	*/
	FIXUP_VTABLE_EX(TextFileFormatPlugin, Plugin);

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

	/**	Performs actions required for menu item selected.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr GoMenuItem(AIMenuMessage* message);
	
	/**	Performs actions required for file format selected.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	virtual ASErr GoFileFormat(AIFileFormatMessage* message);

private:
	/**	File format handle for All Text as Text.
	*/
	AIFileFormatHandle fFileFormatAll;
	
	/**	File format handle for Selected Text as Text.
	*/
	AIFileFormatHandle fFileFormatSelected;
	
	/**	Menu item handle for this plugins About menu.
	*/
	AIMenuItemHandle fAboutPluginMenu;

	/**	Adds the menu items for this plugin to the application UI.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr AddMenus(SPInterfaceMessage* message);

	/**	Adds the file formats for this plugin to the application.
		@param message IN pointer to plugin and call information.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr AddFileFormats(SPInterfaceMessage* message);

	/**	Writes the text to the file.
		@param pathName IN path to file.
		@param selectedTextOnly IN true to write only selected text, false to write all.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr WriteText(const char* pathName, bool selectedTextOnly);
	
	/**	Reads the text from the file.
		@param pathName IN path to file.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr ReadText(const char * pathName);
	
	/**	Creates a new text art item from the text read in from a file.
		@param width IN of text art bounding box.
		@param height IN of text art bounding box.
		@return handle to new text art item.
	*/
	AIArtHandle CreateTextArt(AIReal width, AIReal height);
	
	/** Converts given C string into a Pascal string.
	*/
	void CStrToPStr(char *s, ai::UnicodeString::size_type len);
};
#endif // End TextFileFormatPlugin.h
