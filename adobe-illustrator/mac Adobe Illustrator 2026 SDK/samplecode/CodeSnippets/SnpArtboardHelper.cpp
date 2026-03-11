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

#include "IllustratorSDK.h"

#include "SDKErrors.h"

// Framework includes:
#include "SnpRunnable.h"
#include "SnippetRunnerSuites.h" 
#include "SnippetRunnerLog.h"
#include "SnippetRunnerParameter.h"
#include "SnpSelectionHelper.h"

#include "SnpArtboardHelper.h"
#include <list>

/*
*/
ASErr SnpArtboardHelper::AddArtboard()
{
	ASErr result = kNoErr;
	try {
		// Get values from user.
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		AIRealRect artboardRect = {0, 0, 0, 0}; // left, top, right, bottom
		artboardRect.left = parameter->GetReal(ai::UnicodeString("Enter left boundary of new crop area in points:"), 100.0);
		artboardRect.top = parameter->GetReal(ai::UnicodeString("Enter top boundary of new crop area in points:"), 200.0);
		artboardRect.right = parameter->GetReal(ai::UnicodeString("Enter right boundary of new crop area in points:"), 200.0);
		artboardRect.bottom = parameter->GetReal(ai::UnicodeString("Enter bottom boundary of new crop area in points:"), 100.0);
		ASReal par = parameter->GetReal(ai::UnicodeString("Enter Pixel Aspect Ratio (used if units are pixels):"), 0.0);
		ASBoolean showCenter = parameter->GetBoolean(ai::UnicodeString("Show center mark?"));
		ASBoolean showCrossHairs = parameter->GetBoolean(ai::UnicodeString("Show cross hairs?"));
		ASBoolean showSafeAreas = parameter->GetBoolean(ai::UnicodeString("Show safe areas?"));
		// Call AddArtboard with user input or default values.
		AddArtboard(artboardRect, par, showCenter, showCrossHairs, showSafeAreas);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpArtboardHelper::AddArtboard(AIRealRect artboardRect, AIReal par, ASBoolean showCenter, ASBoolean showCrossHairs, ASBoolean showSafeAreas)
{
	ASErr result = kNoErr;
	try {
		// Create ArtboardProperties using input values.
		ai::ArtboardProperties artboardProperties;
		ai::ArtboardList artboardList;
		sAIArtboard->GetArtboardList(artboardList);
		sAIArtboard->SetPosition(artboardProperties, artboardRect);
		sAIArtboard->SetPAR(artboardProperties, par);
		sAIArtboard->SetShowDisplayMark(artboardProperties, ai::ArtboardProperties::kCenter, showCenter);
		sAIArtboard->SetShowDisplayMark(artboardProperties, ai::ArtboardProperties::kCrossHair, showCrossHairs);
		sAIArtboard->SetShowDisplayMark(artboardProperties, ai::ArtboardProperties::kSafeAreas, showSafeAreas);
		
		// Call AddArtboard with AIArtboard.
		result = AddArtboard(artboardProperties);
		aisdk::check_ai_error(result);
	} catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpArtboardHelper::AddArtboard(ai::ArtboardProperties artboardProperties)
{
	ASErr result = kNoErr;
	try {
		// Add the new artboard to the document.
		ASInt32 artboardIndex = 0;
		ai::ArtboardList artboardList;
		sAIArtboard->GetArtboardList(artboardList);
		result = sAIArtboard->AddNew(artboardList, artboardProperties, artboardIndex);
		aisdk::check_ai_error(result);
		SnippetRunnerLog::Instance()->WritePrintf("Artboard number %d added.", artboardIndex + 1);
		SnippetRunnerLog::Instance()->WritePrintf("");
	} catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}


/*
*/
ASErr SnpArtboardHelper::DeleteArtboard()
{
	ASErr result = kNoErr;
	try {
		ASInt32 artboardCount = this->CountArtboards();
		if (artboardCount > 1) {
			ASInt32 artboardIndex = SnippetRunnerParameter::Instance()->GetInt("Enter artboard number to delete:", 1, 1, artboardCount);
			result = this->DeleteArtboard(artboardIndex-1); // 0-based collection.
			aisdk::check_ai_error(result);
		}
		else
			SnippetRunnerLog::Instance()->Write(ai::UnicodeString("Cannot delete the only artboard in the document."));
	} catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpArtboardHelper::DeleteArtboard(ASInt32 artboardIndex)
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		ai::ArtboardList artboardList;
		sAIArtboard->GetArtboardList(artboardList);
		result = sAIArtboard->Delete(artboardList, artboardIndex);
		if (result == kAICantDeleteLastArtboardErr)
			log->Write(ai::UnicodeString("Cannot delete the last artboard in the document."));
		aisdk::check_ai_error(result);
		log->WritePrintf("Artboard number %d deleted.", artboardIndex + 1);
	} catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpArtboardHelper::DeleteArtboards()
{
	ASErr result = kNoErr;
	AIArtboardRangeHandle rangeHandle = NULL;
	try {
		ASInt32 artboardCount = this->CountArtboards();
		if (artboardCount > 1) {
			ai::UnicodeString artboardRange = SnippetRunnerParameter::Instance()->GetString(ai::UnicodeString("Enter artboard range to delete:"), ai::UnicodeString("1"));
			result = sAIArtboardRange->ValidateString(&artboardRange);
			aisdk::check_ai_error(result);	
			result = sAIArtboardRange->Create(false, &artboardRange, &rangeHandle);
			aisdk::check_ai_error(result);
			result = this->DeleteArtboards(rangeHandle); 
			aisdk::check_ai_error(result);
		}
		else
			SnippetRunnerLog::Instance()->Write(ai::UnicodeString("Cannot delete the only artboard in the document."));
	} catch (ai::Error& ex) {
		result = ex;
	}
	if (rangeHandle != NULL ) {
		result = sAIArtboardRange->Dispose(rangeHandle);
		aisdk::check_ai_error(result);
	}
	return result;
}

/*
*/
ASErr SnpArtboardHelper::DeleteArtboards(AIArtboardRangeHandle artboardRange)
{
	ASErr result = kNoErr;
	AIArtboardRangeIterator rangeIter = NULL;
	try {
		// Create list of artboard indexes to delete.
		result = sAIArtboardRange->RemoveDuplicate(artboardRange);
		aisdk::check_ai_error(result);
		result = sAIArtboardRange->Begin(artboardRange, &rangeIter);
		aisdk::check_ai_error(result);						
		ASInt32 artboardIndex = 0;
		// After AIArtboardRange::Begin, first call to AIArtboardRange::Next gets 
		// the first artboard index in the range.
		result = sAIArtboardRange->Next(rangeIter, &artboardIndex);	

		// Create list to hold artboard indexes.
		list<ASInt32> artboardIndexes;
		list<ASInt32>::iterator iter;
		while (result != kEndOfRangeErr) {
			aisdk::check_ai_error(result);
			artboardIndexes.push_back(artboardIndex);
			result = sAIArtboardRange->Next(rangeIter, &artboardIndex);
		}
		// Sort the indexes into descending order.
		artboardIndexes.sort();
		artboardIndexes.reverse();
		
		// Delete each artboard in descending order, indexing is updated after a deletion.
		for(iter = artboardIndexes.begin(); iter != artboardIndexes.end(); ++iter) {
			result = this->DeleteArtboard(*iter);
			aisdk::check_ai_error(result);
		}
	} catch (ai::Error& ex) {
		result = ex;
	}
	if (rangeIter != NULL) {
		result = sAIArtboardRange->DisposeIterator(rangeIter);
		aisdk::check_ai_error(result);
	}
	return result;
}

/*
*/
ASErr SnpArtboardHelper::EditArtboard()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ASInt32 artboardIndex = parameter->GetInt("Enter artboard number to edit:", 1, 1, this->CountArtboards());
		AIRealRect artboardRect = {0, 0, 0, 0};
		artboardRect.left = parameter->GetReal(ai::UnicodeString("Enter new value for artboard left boundary (points):"), 200.0);
		artboardRect.top = parameter->GetReal(ai::UnicodeString("Enter new value for artboard top boundary (points):"), 300.0);
		artboardRect.right = parameter->GetReal(ai::UnicodeString("Enter new value for artboard right boundary (points):"), 300.0);
		artboardRect.bottom = parameter->GetReal(ai::UnicodeString("Enter new value for artboard bottom boundary (points):"), 200.0);
		ASReal par = parameter->GetReal(ai::UnicodeString("Enter new pixel aspect ratio (used if units are pixels):"), 0.0);
		ASBoolean showCenter = parameter->GetBoolean(ai::UnicodeString("Show center mark?"));
		ASBoolean showCrossHairs = parameter->GetBoolean(ai::UnicodeString("Show cross hairs?"));
		ASBoolean showSafeAreas = parameter->GetBoolean(ai::UnicodeString("Show safe areas?"));
		result = this->EditArtboard(artboardIndex-1 /* 0-based index*/, artboardRect, par, showCenter, showCrossHairs, showSafeAreas);
		aisdk::check_ai_error(result);
	} catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpArtboardHelper::EditArtboard(ASInt32 artboardIndex, AIRealRect artboardRect, AIReal par, ASBoolean showCenter, ASBoolean showCrossHairs, ASBoolean showSafeAreas)
{
	ASErr result = kNoErr;
	try {
		ai::ArtboardList artboardList;
		ai::ArtboardProperties artboardProperties; 
		result = sAIArtboard->GetArtboardProperties(artboardList, artboardIndex, artboardProperties);
		aisdk::check_ai_error(result);
		artboardProperties.SetPosition(artboardRect);
		artboardProperties.SetPAR(par);
		artboardProperties.SetShowDisplayMark(ai::ArtboardProperties::kCenter, showCenter);
		artboardProperties.SetShowDisplayMark(ai::ArtboardProperties::kCrossHair, showCrossHairs);
		artboardProperties.SetShowDisplayMark(ai::ArtboardProperties::kSafeAreas, showSafeAreas);
		result = this->EditArtboard(artboardIndex, artboardProperties);
		aisdk::check_ai_error(result);
	} catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpArtboardHelper::EditArtboard(ASInt32 artboardIndex, ai::ArtboardProperties artboardProperties)
{
	ASErr result = kNoErr;
	try {
		ai::ArtboardList artboardList;
		result = sAIArtboard->GetArtboardList(artboardList);
		aisdk::check_ai_error(result);
		result = sAIArtboard->Update(artboardList, artboardIndex, artboardProperties);
		aisdk::check_ai_error(result);
		result = this->DisplayArtboard(artboardIndex);
		aisdk::check_ai_error(result);
		// Workaround for bug 1789334, have to force redraw of artboards by
		// selecting the crop tool.
		AIToolType toolNumber = NULL;
		result = sAITool->GetToolNumberFromName("Adobe Crop Tool", &toolNumber);
		aisdk::check_ai_error(result);
		AIToolHandle toolHandle = NULL;
		result = sAITool->GetToolHandleFromNumber(toolNumber, &toolHandle);
		aisdk::check_ai_error(result);
		result = sAITool->SetSelectedTool(toolHandle);
		aisdk::check_ai_error(result);
		result = sAITool->GetToolNumberFromName("Adobe Select Tool", &toolNumber);
		aisdk::check_ai_error(result);
		result = sAITool->GetToolHandleFromNumber(toolNumber, &toolHandle);
		aisdk::check_ai_error(result);
		result = sAITool->SetSelectedTool(toolHandle);
		aisdk::check_ai_error(result);
	} catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpArtboardHelper::DisplayArtboard()
{
	ASErr result = kNoErr;
	try {
		// Get the number of the artboard to display.
		ASInt32 artboardIndex = SnippetRunnerParameter::Instance()->GetInt(ai::UnicodeString("Enter artboard number to display:"), 1, 1, this->CountArtboards());
		result = this->DisplayArtboard(artboardIndex-1); // 0-based collection.
		aisdk::check_ai_error(result);
	} catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpArtboardHelper::DisplayArtboard(ASInt32 artboardIndex)
{
	ASErr result = kNoErr;
	try {
		// Display the properties of the specified artboard in the log.
		ai::ArtboardList artboardList;
		ai::ArtboardProperties artboardProperties; 
		result = sAIArtboard->GetArtboardProperties(artboardList, artboardIndex, artboardProperties);
		aisdk::check_ai_error(result);
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		log->WritePrintf("Property values for artboard %d:", artboardIndex+1);
		AIRealRect artboardRect;
		result = artboardProperties.GetPosition(artboardRect);
		SnippetRunnerLog::Indent indent;
		log->WritePrintf("Left boundary = %.2f pts", artboardRect.left);
		log->WritePrintf("Top boundary = %.2f pts", artboardRect.top);
		log->WritePrintf("Right boundary = %.2f pts", artboardRect.right);
		log->WritePrintf("Bottom boundary = %.2f pts", artboardRect.bottom);
		AIReal par;
		artboardProperties.GetPAR(par);
		log->WritePrintf("Pixel Aspect Ratio = %.2f", par);
		bool showCenter;
		artboardProperties.GetShowDisplayMark(ai::ArtboardProperties::kCenter,showCenter);
		log->WritePrintf("Show center mark = %s", showCenter ? "True" : "False");
		bool showCrossHairs;
		artboardProperties.GetShowDisplayMark(ai::ArtboardProperties::kCrossHair,showCrossHairs);
		log->WritePrintf("Show cross hairs = %s", showCrossHairs ? "True" : "False");
		bool showSafeAreas;
		artboardProperties.GetShowDisplayMark(ai::ArtboardProperties::kSafeAreas,showSafeAreas);
		log->WritePrintf("Show safe areas = %s", showSafeAreas ? "True" : "False");
	} catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpArtboardHelper::DisplayActiveArtboard()
{
	ASErr result = kNoErr;
	try {
		ASInt32 artboardIndex = 0;
		ai::ArtboardList artboardList;
		result = sAIArtboard->GetArtboardList(artboardList);
		aisdk::check_ai_error(result);
		result = sAIArtboard->GetActive(artboardList,artboardIndex);
		aisdk::check_ai_error(result);
		result = this->DisplayArtboard(artboardIndex);
		aisdk::check_ai_error(result);
	} catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpArtboardHelper::DisplayAllArtboards()
{
	ASErr result = kNoErr;
	try {
		for(ASInt32 i = 0; i < this->CountArtboards(); i++) {
			result = this->DisplayArtboard(i);
			aisdk::check_ai_error(result);
		}		
	} catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpArtboardHelper::SetActiveArtboard()
{
	ASErr result = kNoErr;
	try {
		// Get the number of the artboard to activate.
		ASInt32 artboardIndex = SnippetRunnerParameter::Instance()->GetInt(ai::UnicodeString("Enter artboard number to set active:"), 1, 1, this->CountArtboards());
		result = this->SetActiveArtboard(artboardIndex-1); // 0-based collection.
		aisdk::check_ai_error(result);
	} catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpArtboardHelper::SetActiveArtboard(ASInt32 artboardIndex)
{
	ASErr result = kNoErr;
	try {
		ASInt32 artboardIndex = 0;
		ai::ArtboardList artboardList;
		result = sAIArtboard->GetArtboardList(artboardList);
		aisdk::check_ai_error(result);
		result = sAIArtboard->SetActive(artboardList, artboardIndex);
		aisdk::check_ai_error(result);
	} catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASInt32 SnpArtboardHelper::CountArtboards()
{
	ASErr result = kNoErr;
	ASInt32 countArtboards = 0;
	try {
		ai::ArtboardList artboardList;
		result = sAIArtboard->GetArtboardList(artboardList);
		aisdk::check_ai_error(result);
		result = sAIArtboard->GetCount(artboardList, countArtboards);
		aisdk::check_ai_error(result);
	} catch (ai::Error& ex) {
		result = ex;
	}
	return countArtboards;
}

ASErr SnpArtboardHelper::RenameActiveArtboard()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ai::UnicodeString name = parameter->GetString(ai::UnicodeString("Enter ative artboard new name:"), ai::UnicodeString("ative artboard"));
		result = this->RenameActiveArtboard(name); 
		aisdk::check_ai_error(result);
	} catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

ASErr SnpArtboardHelper::RenameActiveArtboard(const ai::UnicodeString& name)
{
	ASErr result = kNoErr;
	try {
		ASInt32 artboardIndex = 0;
		ai::ArtboardList artboardList;
		result = sAIArtboard->GetArtboardList(artboardList);
		aisdk::check_ai_error(result);
		result = sAIArtboard->GetActive(artboardList, artboardIndex);
		aisdk::check_ai_error(result);
		ai::ArtboardProperties artboardProperties; 
		result = sAIArtboard->GetArtboardProperties(artboardList, artboardIndex, artboardProperties);
		aisdk::check_ai_error(result);
		ASErr result = artboardProperties.SetName(name);

		sAIArtboard->Update(artboardList, artboardIndex, artboardProperties);
		aisdk::check_ai_error(result);

	} catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpArtboardHelper available to the SnippetRunner framework.
 */
class _SnpRunnableArtboardHelper : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableArtboardHelper () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableArtboardHelper () {}

		/* Returns name of snippet.
		*/
		std::string GetName() const;

		/* Returns a description of what the snippet does.
		*/
		std::string	GetDescription() const;

		/* Returns operations supported by this snippet.
		*/
		Operations GetOperations() const;

		/* Returns name of this snippet's default operation - must
			be one of the operation names returned by GetOperations.
		*/
		std::string GetDefaultOperationName() const;

		/** Returns the categories a snippet belongs to.
			@return default categories.
		*/
		std::vector<std::string> GetCategories() const;

		/* Returns true if the snippet can run.
			@param runnableContext see ISnpRunnableContext for documentation.
			@return true if snippet can run, false otherwise
		 */
		ASBoolean			CanRun(SnpRunnable::Context& runnableContext);

		/* Runs the snippet.
			@param runnableContext see ISnpRunnableContext for documentation.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr		Run(SnpRunnable::Context& runnableContext);
};

/*
*/
std::string _SnpRunnableArtboardHelper::GetName() const
{
	return "ArtboardHelper";
}

/*
*/
std::string _SnpRunnableArtboardHelper::GetDescription() const 
{
	return "Demonstrates how to manipulate artboards and display artboard properties.";
}

/*
*/
SnpRunnable::Operations _SnpRunnableArtboardHelper::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("AddArtboard", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("DeleteArtboard", "document with at least 2 artboards", kSnpRunNewDocumentContext));
	operations.push_back(Operation("DeleteArtboards", "document with at least 2 artboards", kSnpRunNewDocumentContext));
	operations.push_back(Operation("EditArtboard", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("DisplayArtboard", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("DisplayActiveArtboard", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("DisplayAllArtboards", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("SetActiveArtboard", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("RenameActiveArtboard", "document", kSnpRunNewDocumentContext));
	return operations;
}

/*
*/
std::string _SnpRunnableArtboardHelper::GetDefaultOperationName() const 
{
	return "AddArtboard";
	// Note: this must be one of the operation names returned by GetOperations.
}

/*
*/
std::vector<std::string> _SnpRunnableArtboardHelper::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Helper Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableArtboardHelper::CanRun(SnpRunnable::Context& runnableContext)
{
	SnpSelectionHelper selectionHelper;
	SnpArtboardHelper artboardHelper;
	if ("AddArtboard" == runnableContext.GetOperation().GetName() || 
		"EditArtboard" == runnableContext.GetOperation().GetName() || 
		"DisplayArtboard" == runnableContext.GetOperation().GetName() ||  
		"DisplayActiveArtboard" == runnableContext.GetOperation().GetName() || 
		"DisplayAllArtboards" == runnableContext.GetOperation().GetName() || 
		"SetActiveArtboard" == runnableContext.GetOperation().GetName() ||
		"RenameActiveArtboard" == runnableContext.GetOperation().GetName())
		return selectionHelper.IsDocumentSelected();
	else if ("DeleteArtboard" == runnableContext.GetOperation().GetName() || 
			 "DeleteArtboards" == runnableContext.GetOperation().GetName()) {
		if (selectionHelper.IsDocumentSelected() && (artboardHelper.CountArtboards() > 1))
			return true;
		else
			return false;
	}
	else
		return true;
	// Note: this will control the enabling and disabling of the Run button in the SnippetRunner panel
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableArtboardHelper::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpArtboardHelper instance;
	if ("AddArtboard" == runnableContext.GetOperation().GetName()) {
		result = instance.AddArtboard();
	}
	else if ("DeleteArtboard" == runnableContext.GetOperation().GetName()) {
		result = instance.DeleteArtboard();
	}
	else if ("DeleteArtboards" == runnableContext.GetOperation().GetName()) {
		result = instance.DeleteArtboards();
	}
	else if ("EditArtboard" == runnableContext.GetOperation().GetName()) {
		result = instance.EditArtboard();
	}
	else if ("DisplayArtboard" == runnableContext.GetOperation().GetName()) {
		result = instance.DisplayArtboard();
	}
	else if ("DisplayActiveArtboard" == runnableContext.GetOperation().GetName()) {
		result = instance.DisplayActiveArtboard();
	}
	else if ("DisplayAllArtboards" == runnableContext.GetOperation().GetName()) {
		result = instance.DisplayAllArtboards();
	}
	else if ("SetActiveArtboard" == runnableContext.GetOperation().GetName()) {
		result = instance.SetActiveArtboard();
	}
	else if ("RenameActiveArtboard" == runnableContext.GetOperation().GetName()) {
		result = instance.RenameActiveArtboard();
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableArtboardHelper instance_SnpRunnableArtboardHelper;

// End SnpArtboardHelper.cpp

