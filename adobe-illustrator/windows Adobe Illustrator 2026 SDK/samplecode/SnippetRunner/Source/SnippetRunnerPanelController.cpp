//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/SnippetRunner/Source/SnippetRunnerPanelController.cpp $
//
//  $Revision: #2 $
//
//  Copyright 2010 Adobe Systems Incorporated. All rights reserved.
//  
//  NOTICE:  Adobe permits you to use, modify, and distribute this file in accordance 
//  with the terms of the Adobe license agreement accompanying it.  If you have received
//  this file from a source other than Adobe, then your use, modification, or 
//  distribution of it requires the prior written permission of Adobe.
//  
//========================================================================================

#include "SnippetRunnerPanelController.h"
#include "AppContext.hpp"
#include "SnippetRunnerManager.h"
#include "SnippetRunnerLog.h"
#include "SnippetRunnerPreferences.h"
#include "SnippetRunnerSuites.h"
#include <map>

#define SNIPPETRUNNER_UI_EXTENSION "com.adobe.illustrator.snippetrunnerui.panel"

// Events we listen for
#define EVENT_TYPE_PANEL_READY					"com.adobe.csxs.events.SnippetRunnerPanelReady"
#define EVENT_TYPE_SELECTION_CHANGE				"com.adobe.csxs.events.SnippetRunnerSelectionChange"
#define EVENT_TYPE_RUN							"com.adobe.csxs.events.SnippetRunnerRun"
#define EVENT_TYPE_CLEAR_LOG					"com.adobe.csxs.events.SnippetRunnerClearLog"
#define EVENT_TYPE_SAVE_LOG						"com.adobe.csxs.events.SnippetRunnerSaveLog"
#define EVENT_TYPE_PREFS_FROM_UI				"com.adobe.csxs.events.SnippetRunnerSetPrefs"

// Events we dispatch
#define EVENT_TYPE_UPDATE_PANEL_SNIPPET_TREE	"com.adobe.csxs.events.SnippetRunnerUpdatePanel"
#define EVENT_TYPE_UPDATE_PANEL_PREFS			"com.adobe.csxs.events.SnippetRunnerPrefsToPanel"
#define EVENT_TYPE_CAN_RUN						"com.adobe.csxs.events.SnippetRunnerCanRun"
#define EVENT_TYPE_NEW_LOG_DATA					"com.adobe.csxs.events.SnippetRunnerLog"

#define APP_ID									"ILST"

typedef std::map<std::string, SnippetRunnerParameter::ParameterInputMode> param_map;

param_map init_paramMap()
{
	std::map<std::string, SnippetRunnerParameter::ParameterInputMode> ret;
	ret["default"] = SnippetRunnerParameter::kDefaultParameterInputMode;
	ret["csv"] = SnippetRunnerParameter::kCSVParameterInputMode;
	ret["prompt"] = SnippetRunnerParameter::kPromptParameterInputMode;
	ret["assets"] = SnippetRunnerParameter::kAssetFolderInputMode;
	return ret;
}

static void PanelReadyFunc(const csxs::event::Event* const eventParam, void* const context)
{
	SnippetRunnerPanelController *controller = (SnippetRunnerPanelController *) context;
	
	if ( controller )
	{
		controller->UpdatePanelSnippetTree();
		controller->UpdatePanelPrefs();
	}
}

static void SelectionChangedFunc(const csxs::event::Event* const eventParam, void* const context)
{
	AppContext appContext(gPlugin->GetPluginRef());
	SnippetRunnerPanelController *controller = (SnippetRunnerPanelController *) context;
	
	if ( controller && eventParam )
		controller->SelectionChanged(eventParam->data);
	
}

static void RunSnippetHandler(const csxs::event::Event* const eventParam, void* const context)
{
	AppContext appContext(gPlugin->GetPluginRef());
	SnippetRunnerPanelController *controller = (SnippetRunnerPanelController *) context;
	
	if ( controller && eventParam )
		controller->RunSnippet(eventParam->data);

}

static void ClearLogHandler(const csxs::event::Event* const eventParam, void* const context)
{
	AppContext appContext(gPlugin->GetPluginRef());
	SnippetRunnerPanelController *controller = (SnippetRunnerPanelController *) context;
	
	if ( controller )
		controller->ClearLog();
}

static void SaveLogHandler(const csxs::event::Event* const eventParam, void* const context)
{
	AppContext appContext(gPlugin->GetPluginRef());
	SnippetRunnerPanelController *controller = (SnippetRunnerPanelController *) context;
	
	if ( controller )
		controller->SaveLog();
}

static void GetPrefsHandler(const csxs::event::Event* const eventParam, void* const context)
{
	AppContext appContext(gPlugin->GetPluginRef());
	SnippetRunnerPanelController *controller = (SnippetRunnerPanelController *) context;
	
	if ( controller )
		controller->ReceivePrefsChange(eventParam->data);
}

// Get the contents of the first instance of the "tag" tag in the xmlString,
// assuming the tag has no attributes
static std::string getTagContents(const std::string xmlString, const std::string tag)
{
	// Get the index of the char beyond the end of the opening tag
	std::string openTag("<" + tag + ">");
	std::string::size_type openIndex = xmlString.find(openTag);
	openIndex += openTag.length();
	
	// Get the index of the first char of the ending tag
	std::string closingTag("</" + tag + ">");
	std::string::size_type closingIndex = xmlString.find(closingTag);
	
	// Now get the tag contents
	std::string tagContents("");
	std::string::size_type strLength = closingIndex - openIndex;
	if (strLength > 0)
		tagContents = xmlString.substr(openIndex, strLength);
	
	return tagContents;
}

csxs::event::EventErrorCode SnippetRunnerPanelController::RegisterCSXSEventListeners()
{
	csxs::event::EventErrorCode result = csxs::event::kEventErrorCode_Success;
	
	do
	{
		// Listener for panel ready
		result = fPPLib.AddEventListener(EVENT_TYPE_PANEL_READY, PanelReadyFunc, this);
		if ( result != csxs::event::kEventErrorCode_Success )
		{
			break;
		}
		
		// Listener for snippet selection events
		result = fPPLib.AddEventListener(EVENT_TYPE_SELECTION_CHANGE, SelectionChangedFunc, this);
		if ( result != csxs::event::kEventErrorCode_Success )
		{
			break;
		}
		
		// Listener for run snippet events
		result = fPPLib.AddEventListener(EVENT_TYPE_RUN, RunSnippetHandler, this);
		if ( result != csxs::event::kEventErrorCode_Success )
		{
			break;
		}
		
		// Listener for clear log events
		result = fPPLib.AddEventListener(EVENT_TYPE_CLEAR_LOG, ClearLogHandler, this);
		if ( result != csxs::event::kEventErrorCode_Success )
		{
			break;
		}
		
		// Listener for log save events
		result = fPPLib.AddEventListener(EVENT_TYPE_SAVE_LOG, SaveLogHandler, this);
		if ( result != csxs::event::kEventErrorCode_Success )
		{
			break;
		}
		
		result = fPPLib.AddEventListener(EVENT_TYPE_PREFS_FROM_UI, GetPrefsHandler, this);
		if ( result != csxs::event::kEventErrorCode_Success )
		{
			break;
		}
	}
	while ( false );
		
	return result;
}

csxs::event::EventErrorCode SnippetRunnerPanelController::RemoveEventListeners()
{
	csxs::event::EventErrorCode result = csxs::event::kEventErrorCode_Success;

	do
	{
		result = fPPLib.RemoveEventListener(EVENT_TYPE_PANEL_READY, PanelReadyFunc, this);
		if ( result != csxs::event::kEventErrorCode_Success )
		{
			break;
		}

		result = fPPLib.RemoveEventListener(EVENT_TYPE_SELECTION_CHANGE, SelectionChangedFunc, this);
		if ( result != csxs::event::kEventErrorCode_Success )
		{
			break;
		}

		result = fPPLib.RemoveEventListener(EVENT_TYPE_RUN, RunSnippetHandler, this);
		if ( result != csxs::event::kEventErrorCode_Success )
		{
			break;
		}
		
		result = fPPLib.RemoveEventListener(EVENT_TYPE_CLEAR_LOG, ClearLogHandler, this);
		if ( result != csxs::event::kEventErrorCode_Success )
		{
			break;
		}
		
		result = fPPLib.RemoveEventListener(EVENT_TYPE_SAVE_LOG, SaveLogHandler, this);
		if ( result != csxs::event::kEventErrorCode_Success )
		{
			break;
		}
		
		result = fPPLib.RemoveEventListener(EVENT_TYPE_PREFS_FROM_UI, GetPrefsHandler, this);
		if ( result != csxs::event::kEventErrorCode_Success )
		{
			break;
		}
	}
	while ( false );
	
	return result;
}

void SnippetRunnerPanelController::UpdatePanelSnippetTree()
{
	// Get a hold of the snippets
	SnippetRunnerManager& snipRunnerManager = SnippetRunnerManager::Instance();
	SnippetRunnerManager::Snippets snippets = snipRunnerManager.GetSnippets();
	
	std::map<std::string, std::string> catBlobs;
	
	// Construct the XML
	for (SnippetRunnerManager::Snippets::const_iterator iter = snippets.begin();
		 iter != snippets.end(); iter++)
	{
		SnpRunnable* snippet = *iter;
		
		// First build up a string containing the XML for the snippet.
		stringstream snipXml;
		snipXml << "<snippet label=\"" + snippet->GetName() + "\">" 
		<< "<description label=\"Description: " + snippet->GetDescription() + "\" />";
		
		// Now populate the operations tag
		SnpRunnable::Operations operations = snippet->GetOperations();
		if ( operations.size() )
		{
			snipXml << "<operations label=\"Operations:\">";
			
			for (SnpRunnable::Operations::const_iterator iter = operations.begin();
				 iter != operations.end(); iter++)
			{
				snipXml << "<operation label=\"" + iter->GetName() + "\">";
				
				// Now add any preconditions
				std::string preconditions = iter->GetPreconditions();
				if ( preconditions.size() )
				{
					snipXml << "<precondition label=\"Preconditions: " + preconditions + "\" />";
				}
				snipXml << "</operation>";
			}
			snipXml << "</operations>";
		}
		snipXml << "</snippet>";
		
		// Now that we have the snippet XML blob we need to add it to each of the
		// categories that the snippet belongs to.
		std::vector<std::string> categories = snippet->GetCategories();
		for (std::vector<std::string>::const_iterator catIter = categories.begin();
			 catIter != categories.end(); catIter++)
		{
			catBlobs[*catIter] += snipXml.str();
		}
	}
	
	// Now construct the full CSXS payload
	stringstream fullXml;
	fullXml << "<payload>";
	for (std::map<std::string, std::string>::const_iterator mapIter = catBlobs.begin();
		 mapIter != catBlobs.end(); mapIter++)
	{
		fullXml << "<category label=\"" + mapIter->first + "\" >"
		<< mapIter->second << "</category>";
	}
	fullXml << "</payload>";
	
	const std::string temp = fullXml.str();
	const char* tempcstr = temp.c_str();

	// Create the event
	csxs::event::Event event = {EVENT_TYPE_UPDATE_PANEL_SNIPPET_TREE,
		csxs::event::kEventScope_Application,
		APP_ID,
		NULL,
		tempcstr};
	
	fPPLib.DispatchEvent(&event);
}

void SnippetRunnerPanelController::SelectionChanged(const char* eventData)
{
	/* Recieved XML will look like this:
	 * <payload>
	 *   <snippet>snippet name</snippet>
	 *   <op>op name</op>
	 * </payload>
	 */
	
	// Get the snippet and operation names
	std::string dataStr(eventData);
	fSelectedSnippet = getTagContents(dataStr, "snippet");
	fSelectedOp = getTagContents(dataStr, "op");
	
	if ( fSelectedSnippet == "")
		return;
	
	fCanRunSelected = CanRun(fSelectedSnippet, fSelectedOp);
	
	DispatchCanRunEvent(fSelectedSnippet, fSelectedOp, fCanRunSelected);
}

void SnippetRunnerPanelController::RunSnippet(const char* eventData)
{
	// Get the snippet and operation names
	std::string dataStr(eventData);
	std::string snippetName = getTagContents(dataStr, "snippet");
	std::string opName = getTagContents(dataStr, "op");
	
	SnippetRunnerManager& snipRunnerManager = SnippetRunnerManager::Instance();
	
	if (!SnippetRunnerPreferences::Instance()->GetUnitTest())
	{
		// Not in unit test mode, run the snippet
		
		// TODO: Check for CSV parameter mode and prompt for params

			snipRunnerManager.Run(snippetName, opName);
	}
	else
	{
		// Unit test mode
		// if a snippet is selected, test it
		std::vector<std::string> snippetNames;

		if(snippetName == "category")
		{
			SnippetRunnerManager::Snippets snippets = snipRunnerManager.GetSnippets(opName);
			for (SnippetRunnerManager::Snippets::const_iterator iter = snippets.begin(); iter != snippets.end(); iter++) 
			{
				SnpRunnable* snp = *iter;
				SnpRunnable::Operations operations = snp->GetOperations();
				for(SnpRunnable::Operations::const_iterator operations_iter = operations.begin(); operations_iter != operations.end(); operations_iter++)
				{
					SnpRunnable::Operation operation = *operations_iter;
					snippetNames.push_back(snp->GetName());
				}
			}
			opName = "";
		}
		else
		{
			snippetNames.push_back(snippetName);
		}
		snipRunnerManager.RunUnitTest(snippetNames, opName);
		
		// TODO: If no snippet is selected, test every snippet in the category
	}
}

void SnippetRunnerPanelController::ReceivePrefsChange(const char* eventData)
{
	std::string dataStr(eventData);
	
	// Get the unit test mode
	ASBoolean unitTest = (getTagContents(dataStr, "unittest").compare("true") == 0);
	
	// Get the parameter input mode as a string
	SnippetRunnerParameter::ParameterInputMode inputMode;
	std::string tempParamMode = getTagContents(dataStr, "parammode");
	
	// Get the enum value for the parameter mode
	param_map::const_iterator it = fParamMap.find(tempParamMode);
	if (it == fParamMap.end())
	{
		inputMode = SnippetRunnerParameter::kDefaultParameterInputMode;
	}
	else
	{
		inputMode = it->second;
	}
	if (inputMode == SnippetRunnerParameter::kAssetFolderInputMode)
		SnippetRunnerPreferences::Instance()->PromptAssetsFolderPath();
	// set the preferences
	SnippetRunnerPreferences* prefs = SnippetRunnerPreferences::Instance();
	prefs->SetParameterInputMode(inputMode);
	prefs->SetUnitTest(unitTest);	
}

ASBoolean SnippetRunnerPanelController::CanRun(std::string snippetName, std::string opName)
{
	// Determine whether the snippet can be run
	ASBoolean canRun = false;
	const ASBoolean isUnitTestRunning = SnippetRunnerManager::Instance().IsUnitTestRunning();
	
	if ( SnippetRunnerPreferences::Instance()->GetUnitTest() )
	{
		// Unit test preference is checked but no tests are running - enable
		// the run button, instead of running the operation it will run the snippet
		// test suite.
		if (!isUnitTestRunning)
			canRun = true;
	}
	else if (!isUnitTestRunning)
	{
		// Not in test mode and no tests are still running, check if the operation can
		// be run
		canRun = SnippetRunnerManager::Instance().CanRun(snippetName, opName);
	}
	
	return canRun;
}

void SnippetRunnerPanelController::DispatchCanRunEvent(std::string snippetName,
													   std::string opName,
													   ASBoolean canRun)
{
	/* Dispatch XML of the following form:
	 * <payload>
	 *   <snippet>snippet name</snippet>
	 *   <op>op name</op>
	 *   <canrun>true/false</canrun>
	 * </payload>
	 */
	
	stringstream dispatchData;
	dispatchData << "<payload><snippet>" << snippetName
	<< "</snippet><op>" << opName
	<< "</op><canrun>" 
	<< (canRun ? "true" : "false")
	<< "</canrun></payload>";
	
	const std::string temp = dispatchData.str();
	const char* tempcstr = temp.c_str();

	csxs::event::Event event = {EVENT_TYPE_CAN_RUN,
		csxs::event::kEventScope_Application,
		APP_ID,
		NULL,
		tempcstr};
	
	fPPLib.DispatchEvent(&event);
}

void SnippetRunnerPanelController::UpdatePanelPrefs()
{
	/* Dispatch xml of the following form:
	 * <payload>
	 *   <parammode>default/prompt/csv</parammode>
	 *   <unittest>true/false</unittest>
	 * </payload>
	 */
	
	SnippetRunnerPreferences* prefs = SnippetRunnerPreferences::Instance();
	ASBoolean unitTest = prefs->GetUnitTest();
	SnippetRunnerParameter::ParameterInputMode paramMode = prefs->GetParameterInputMode();
	
	string paramModeStr;
	
	for (param_map::const_iterator it = fParamMap.begin(); it != fParamMap.end(); it++)
	{
		if (it->second == paramMode)
		{
			paramModeStr = it->first;
			break;
		}
	}
			
	stringstream dispatchData;
	dispatchData << "<payload><parammode>"
	<< paramModeStr
	<< "</parammode><unittest>"
	<< (unitTest ? "true" : "false")
	<< "</unittest></payload>";

	const std::string temp = dispatchData.str();
	const char* tempcstr = temp.c_str();

	csxs::event::Event event = {EVENT_TYPE_UPDATE_PANEL_PREFS,
		csxs::event::kEventScope_Application,
		APP_ID,
		NULL,
		tempcstr};
	
	fPPLib.DispatchEvent(&event);
}

void SnippetRunnerPanelController::HandleLogChanged()
{
	const SnippetRunnerLog::Log& logData = SnippetRunnerLog::Instance()->GetLog();
	stringstream dispatchData;
	
	if ( !logData.size() )
		return;
	
	dispatchData << "<payload>";
	
	while ( fLastLogLineSent < logData.size() )
	{
		dispatchData << "<logline>"
			<< SnippetRunnerPanelController::EscapeForXML(logData[fLastLogLineSent]).as_Platform()
			<< "</logline>";
		++fLastLogLineSent;
	}
	
	dispatchData << "</payload>";

	const std::string temp = dispatchData.str();
	const char* tempcstr = temp.c_str();
	
	csxs::event::Event event = {EVENT_TYPE_NEW_LOG_DATA,
		csxs::event::kEventScope_Application,
		APP_ID,
		NULL,
		tempcstr};
	
	fPPLib.DispatchEvent(&event);
}

ai::UnicodeString SnippetRunnerPanelController::EscapeForXML(ai::UnicodeString srcString)
{
	ai::UnicodeString findCharacters [] = {
		ai::UnicodeString("&"),
		ai::UnicodeString ("<"),
		ai::UnicodeString(">"),
		ai::UnicodeString("\""),
		ai::UnicodeString("'") };
	ai::UnicodeString replaceCharacters [] = {
		ai::UnicodeString("&amp;"),
		ai::UnicodeString ("&lt;"),
		ai::UnicodeString("&gt;"),
		ai::UnicodeString("&quot;"),
		ai::UnicodeString("&apos")
	};

	for(int i = 0; i < 5;i++)
	{
		srcString = FindAndReplaceString(srcString,findCharacters[i],replaceCharacters[i]);
	}
	return srcString;
}

ai::UnicodeString SnippetRunnerPanelController::FindAndReplaceString(ai::UnicodeString srcString,ai::UnicodeString findString,ai::UnicodeString replace)
{
	ai::UnicodeString::size_type start = srcString.find(findString);
	if(start < srcString.length())
	{
		srcString.replace(start,findString.length(),replace);
	}
	return srcString;
}

AIErr SnippetRunnerPanelController::IsPrimaryStageVisible(AIBoolean& isVisible)
{
	AIErr error = kNoErr;
	
	if(sAICSXSExtension != NULL)
		error = sAICSXSExtension->IsPrimaryStageVisible(SNIPPETRUNNER_UI_EXTENSION, isVisible);

	return error;
}

void SnippetRunnerPanelController::ClearLog()
{
	SnippetRunnerLog::Instance()->Clear();
	fLastLogLineSent = 0;
}

void SnippetRunnerPanelController::SaveLog()
{
	SnippetRunnerLog::Instance()->Save();
}

void SnippetRunnerPanelController::HandleModelChanged()
{
	// Get the can run status of the selected snippet
	ASBoolean canRunNow = CanRun(fSelectedSnippet, fSelectedOp);
	
	// If the status has changed, dispatch an event & store the new status
	if ( fCanRunSelected != canRunNow )
	{
		DispatchCanRunEvent(fSelectedSnippet, fSelectedOp, canRunNow);
		fCanRunSelected = canRunNow;
	}
}

SnippetRunnerPanelController::SnippetRunnerPanelController()
: FlashUIController(SNIPPETRUNNER_UI_EXTENSION),
fSelectedSnippet(""),
fSelectedOp(""),
fCanRunSelected(false),
fLastLogLineSent(0)
{
	fParamMap = init_paramMap();
}
