//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/SnippetRunner/Source/SnippetRunnerPanelController.h $
//
//  $Revision: #2 $
//
//  Copyright 1987 Adobe Systems Incorporated. All rights reserved.
//  
//  NOTICE:  Adobe permits you to use, modify, and distribute this file in accordance 
//  with the terms of the Adobe license agreement accompanying it.  If you have received
//  this file from a source other than Adobe, then your use, modification, or 
//  distribution of it requires the prior written permission of Adobe.
//  
//========================================================================================

#ifndef __SNIPPETRUNNERPANELCONTROLLER_H__
#define __SNIPPETRUNNERPANELCONTROLLER_H__

#include "FlashUIController.h"
#include "SDKPlugPlug.h"
#include "SnippetRunnerPlugin.h"
#include "SnippetRunnerParameter.h"
#include <map>

extern SnippetRunnerPlugin *gPlugin;

class SnippetRunnerPanelController : public FlashUIController
{
public:
	SnippetRunnerPanelController();
		
	csxs::event::EventErrorCode RegisterCSXSEventListeners();
	
	csxs::event::EventErrorCode RemoveEventListeners();
	
	ASErr SendData() {return kNoErr;};
	
	/** Handles refreshing the view of the log - see SnippetRunnerLog.
	*/
	void HandleLogChanged(void);
	
	/** Handles updating the "can run" status of a snippet when something changes
		in Illustrator.
	*/
	void HandleModelChanged(void);
	
	void ParseData(const char* eventData) {return ;};
	
	/** Send a list of snippets to the flash panel
	*/
	void UpdatePanelSnippetTree(void);
	
	/** Send the preferences to the flash panel
	*/
	void UpdatePanelPrefs(void);
	
	/** Handle the selected snippet in the flash panel changing, send
		whether or not we can run the snippet back to the panel
	*/
	void SelectionChanged(const char* eventData);
	
	/** Handle a run snippet event from the flash panel
	*/
	void RunSnippet(const char* eventData);
	
	/** Handle preferences being sent from the flash panel
	*/
	void ReceivePrefsChange(const char* eventData);
	
	/** Clear the log
	*/
	void ClearLog(void);
	
	/** Save the log
	*/
	void SaveLog(void);

	AIErr IsPrimaryStageVisible(AIBoolean& visible); 

	static ai::UnicodeString EscapeForXML(ai::UnicodeString srcString);

	static ai::UnicodeString FindAndReplaceString(ai::UnicodeString srcString,ai::UnicodeString findString,ai::UnicodeString replace);
	
private:
	
	/** The currently selected snippet.
	*/
	std::string fSelectedSnippet;
	
	/** The currently selected operation.
	*/
	std::string fSelectedOp;
	
	/** The cached "can run" status of the selected snippet
	*/
	ASBoolean fCanRunSelected;
	
	/** A map between parameter input mode identifiers and their values
	*/
	std::map<std::string, SnippetRunnerParameter::ParameterInputMode> fParamMap;
	
	/** Determine whether or not the given snippet and operation can be run
	*/
	ASBoolean CanRun(std::string snippetName, std::string opName);
	
	/** Dispatch a can run event for the given snippet and operation
	*/
	void DispatchCanRunEvent(std::string snippetName, std::string opName, ASBoolean canRun);
	
	/** The last line of the log sent to the panel
	*/
	int fLastLogLineSent;
};

#endif
