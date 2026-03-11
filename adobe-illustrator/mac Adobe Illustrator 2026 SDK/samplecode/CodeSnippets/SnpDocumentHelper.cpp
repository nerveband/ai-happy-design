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

// Project includes:
#include "SDKErrors.h"
#include "SnpRunnable.h"
#include "SnippetRunnerLog.h"
#include "SnippetRunnerParameter.h"
#include "SnippetRunnerSuites.h"
#include "SnippetRunnerPreferences.h"
#include "SnpArtSetHelper.h"
#include "SnpChooser.h"
#include "SnpDocumentHelper.h"
#include "SnpSelectionHelper.h"

/*
*/
ASErr SnpDocumentHelper::NewDocument(const ai::UnicodeString& preset, AIDocumentHandle& document)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIDocumentList);
		result = sAIDocumentList->New(const_cast<ai::UnicodeString&>(preset), NULL, kDialogNone, &document);
		aisdk::check_ai_error(result);
		SDK_ASSERT(document);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpDocumentHelper::NewCustomDocument(const ai::UnicodeString& preset, AIReal width, AIReal height, AIDocumentHandle& document)
{
	ASErr result = kNoErr;
	try {
		AINewDocumentPreset settings;
		SDK_ASSERT(sAIDocumentList);
		result = sAIDocumentList->GetPresetSettings(const_cast<ai::UnicodeString&>(preset), &settings);
		aisdk::check_ai_error(result);
		settings.docWidth = width;
		settings.docHeight = height;
		result = sAIDocumentList->New(const_cast<ai::UnicodeString&>(preset), &settings, kDialogNone, &document);
		aisdk::check_ai_error(result);
		SDK_ASSERT(document);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpDocumentHelper::NewCustomDocument(const ai::UnicodeString& preset, AIReal width, AIReal height, ASInt32 numArtboards, AIArtboardLayout artboardLayout, AIReal artboardSpacing, ASInt32 numRowsOrCols, AIDocumentHandle& document)
{
	ASErr result = kNoErr;
	try {
		AINewDocumentPreset settings;
		SDK_ASSERT(sAIDocumentList);
		result = sAIDocumentList->GetPresetSettings(const_cast<ai::UnicodeString&>(preset), &settings);
		aisdk::check_ai_error(result);
		settings.docWidth = width;
		settings.docHeight = height;
		settings.docNumArtboards = numArtboards;
		settings.docArtboardLayout = artboardLayout;
		settings.docArtboardSpacing = artboardSpacing;
		settings.docArtboardRowsOrCols = numRowsOrCols;
		result = sAIDocumentList->New(const_cast<ai::UnicodeString&>(preset), &settings, kDialogNone, &document);
		aisdk::check_ai_error(result);
		SDK_ASSERT(document);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpDocumentHelper::OpenDocument(const ai::UnicodeString& name, AIDocumentHandle& document)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIDocumentList);
		AIColorModel colorModel = kAIUnknownColorModel;		
		ai::FilePath path = this->GetAssetPath(name);
		ActionDialogStatus dialogStatus = path.Exists(false) ? kDialogOff : kDialogOn;
		AIBoolean forceCopy = false;
		result = sAIDocumentList->Open(path, colorModel, dialogStatus, forceCopy, &document);
		aisdk::check_ai_error(result);
		SDK_ASSERT(document);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpDocumentHelper::ActivateDocument(AIDocumentHandle document)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(document);
		SDK_ASSERT(sAIDocumentList);
		ASBoolean bSetFocus = true;
		result = sAIDocumentList->Activate(document, bSetFocus);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}


/*
*/
ASErr SnpDocumentHelper::PrintDocument(AIDocumentHandle document)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(document);
		SDK_ASSERT(sAIDocumentList);
		result = sAIDocumentList->Print(document, kDialogOn);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpDocumentHelper::SaveDocument(AIDocumentHandle document)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(document);
		SDK_ASSERT(sAIDocumentList);
		ASErr status = sAIDocumentList->Save(document);
		aisdk::check_ai_error(status);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpDocumentHelper::CloseDocument(AIDocumentHandle document)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(document);
		SDK_ASSERT(sAIDocumentList);
		result = sAIDocumentList->Close(document);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpDocumentHelper::CloseAllDocuments()
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIDocumentList);
		result = sAIDocumentList->CloseAll();
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpDocumentHelper::GetCurrentDocument(AIDocumentHandle& document)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIDocument);
		result = sAIDocument->GetDocument(&document);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpDocumentHelper::GetCurrentLayer(AILayerHandle& layer)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAILayer);
		result = sAILayer->GetCurrentLayer(&layer);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASBoolean SnpDocumentHelper::CanPlaceFile(void)
{
	ASBoolean result = true;
	try {
		AILayerHandle layer = NULL;
		ASErr status = this->GetCurrentLayer(layer);
		aisdk::check_ai_error(status);
		AIBoolean editable = false;
		status = sAILayer->GetLayerEditable(layer, &editable);
		aisdk::check_ai_error(status);
		AIBoolean visible = false;
		status = sAILayer->GetLayerVisible(layer, &visible);
		aisdk::check_ai_error(status);
		result = editable && visible;
	}
	catch (ai::Error) {
	}
	return result;
}

/*
*/
ASErr SnpDocumentHelper::PlaceFile(const ai::UnicodeString& name, AIArtHandle& art)
{
	ASErr result = kNoErr;

	try {
		ai::FilePath path = this->GetAssetPath(name);
		AIPlaceRequestData placeReqData;
		placeReqData.m_lPlaceMode = kVanillaPlace; // places the file and returns new art handle
		placeReqData.m_disableTemplate = true; // disable Template placement
		placeReqData.m_filemethod = 0; // place the file
		placeReqData.m_hNewArt = NULL; // new art handle
		placeReqData.m_hOldArt = NULL; // art to be replaced
		placeReqData.m_lParam = kPlacedArt; // type of art to create
		placeReqData.m_pFilePath = NULL; // if NULL, shows File Place dialog
		if (path.Exists(false) && name.size() != 0) {
			placeReqData.m_pFilePath = &path; // place the given file (prompt for file to be placed otherwise)
		}
		
		// Place the file.
		result = sAIPlaced->ExecPlaceRequest(placeReqData);
		aisdk::check_ai_error(result);

		// Verify that the handle of the placed art is returned.
		SDK_ASSERT(placeReqData.m_hNewArt);

		// Pass handle to placed art back to caller.
		art = placeReqData.m_hNewArt;
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}


/*
*/
ai::FilePath SnpDocumentHelper::GetAssetPath(const ai::UnicodeString& name)
{
	ai::FilePath path = SnippetRunnerPreferences::Instance()->GetAssetsFolderPath();
	path.AddComponent(ai::FilePath(name));
	return path;
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpDocumentHelper available to the SnippetRunner framework.
 */
class _SnpRunnableDocumentHelper : public SnpRunnable
{
	public:
		/* Constructor registers the snippet to the framework.
		 */
		_SnpRunnableDocumentHelper () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual	~_SnpRunnableDocumentHelper () {}

		/* Returns name of snippet.
		*/
		std::string GetName() const;

		/* Returns a description of what the snippet does.
		*/
		std::string GetDescription() const;

		/** Returns operations supported by this snippet.
		*/
		Operations GetOperations() const;

		/** Returns name of the snippet's default operation.
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
std::string _SnpRunnableDocumentHelper::GetName() const
{
	return "DocumentHelper";
}

/*
*/
std::string _SnpRunnableDocumentHelper::GetDescription() const 
{
	return "Creates, opens, saves and closes documents";
}

/*
*/
SnpRunnable::Operations _SnpRunnableDocumentHelper::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("NewDocument", "", kSnpRunAnyContext));
	operations.push_back(Operation("NewCustomDocument", "", kSnpRunAnyContext));
	operations.push_back(Operation("OpenDocument", "", kSnpRunAnyContext));
	operations.push_back(Operation("ActivateDocument", "documents", kSnpRunNotSupportedContext));	
	operations.push_back(Operation("PrintDocument", "document", kSnpRunNotSupportedContext));
	operations.push_back(Operation("SaveDocument", "document", kSnpRunNotSupportedContext));	
	operations.push_back(Operation("CloseDocument", "document", kSnpRunNewDocumentContext));	
	operations.push_back(Operation("PlaceFile", "document", kSnpRunNewDocumentContext));	
	return operations;
}

/*
*/
std::string _SnpRunnableDocumentHelper::GetDefaultOperationName() const 
{
	return "OpenDocument";
}

/*
*/
std::vector<std::string> _SnpRunnableDocumentHelper::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Helper Snippets");
	return categories;
}

/*
*/
ASBoolean _SnpRunnableDocumentHelper::CanRun(SnpRunnable::Context& runnableContext)
{
	SnpSelectionHelper selectionHelper;
	if ("ActivateDocument" == runnableContext.GetOperation().GetName() ||
		"PrintDocument" == runnableContext.GetOperation().GetName() ||
		"SaveDocument" == runnableContext.GetOperation().GetName() ||
		"CloseDocument" == runnableContext.GetOperation().GetName() ||
		"PlaceFile" == runnableContext.GetOperation().GetName()) {
				
			return selectionHelper.IsDocumentSelected();
	}
	else
		return true;
}

/*
*/
ASErr _SnpRunnableDocumentHelper::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr; 

	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnpDocumentHelper instance;
		if ("NewDocument" == runnableContext.GetOperation().GetName()) {
			AIDocumentHandle document = NULL;
			result = instance.NewDocument(ai::UnicodeString("Print"), document);
			aisdk::check_ai_error(result);
			SDK_ASSERT(document);
		}
		else if ("NewCustomDocument" == runnableContext.GetOperation().GetName()) {
			AIDocumentHandle document = NULL;
			result = instance.NewCustomDocument(ai::UnicodeString("Basic RGB"), 400.0, 300.0, 6, kAIArtboardLayoutGridByCol, 30, 3, document);
			aisdk::check_ai_error(result);
			SDK_ASSERT(document);
		}
		else if ("OpenDocument" == runnableContext.GetOperation().GetName()) {
			AIDocumentHandle document = NULL;
			result = instance.OpenDocument(ai::UnicodeString("sample-1.ai"), document);
			aisdk::check_ai_error(result);
			SDK_ASSERT(document);
		}
		else if ("ActivateDocument" == runnableContext.GetOperation().GetName()) {
			SnpChooser chooser;
			AIDocumentHandle document = chooser.PickDocument();
			result = instance.ActivateDocument(document);
			aisdk::check_ai_error(result);
		}
		else if ("PrintDocument" == runnableContext.GetOperation().GetName()) {
			AIDocumentHandle document = NULL;
			result = instance.GetCurrentDocument(document);
			aisdk::check_ai_error(result);
			result = instance.PrintDocument(document);
			aisdk::check_ai_error(result);
		}
		else if ("SaveDocument" == runnableContext.GetOperation().GetName()) {
			AIDocumentHandle document = NULL;
			result = instance.GetCurrentDocument(document);
			aisdk::check_ai_error(result);
			result = instance.SaveDocument(document);
			aisdk::check_ai_error(result);
		}
		else if ("CloseDocument" == runnableContext.GetOperation().GetName()) {
			AIDocumentHandle document = NULL;
			result = instance.GetCurrentDocument(document);
			aisdk::check_ai_error(result);
			result = instance.CloseDocument(document);
			aisdk::check_ai_error(result);
		}
		else if ("PlaceFile" == runnableContext.GetOperation().GetName()) {
			if (instance.CanPlaceFile()) {
				AIArtHandle art = NULL;
				result = instance.PlaceFile(ai::UnicodeString("thistle.psd"), art);
				aisdk::check_ai_error(result);
			} else {
				log->Write("Cannot place file");
			}
		}
		else {
			result = kBadParameterErr;
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*	Declare an instance to register the snippet hook with the framework.
*/
static _SnpRunnableDocumentHelper instance_SnpRunnableDocumentHelper;

// End SnpDocumentHelper.cpp

