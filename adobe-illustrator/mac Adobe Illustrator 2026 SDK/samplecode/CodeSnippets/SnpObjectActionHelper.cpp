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
#include "AIObjectAction.h"

#include "SDKErrors.h"

// Framework includes:
#include "SnpRunnable.h"
#include "SnippetRunnerSuites.h"
#include "SnippetRunnerLog.h"
#include "SnippetRunnerParameter.h"
#include "SnpObjectActionHelper.h"
#include "SnpSelectionHelper.h"

/*
*/
SnpObjectActionHelper::VPB::VPB() : fActionParamValueRef(NULL)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIActionManager);
		result = sAIActionManager->AINewActionParamValue(&fActionParamValueRef);
		SDK_ASSERT(!result);
		SDK_ASSERT(fActionParamValueRef);
	}
	catch (ai::Error) {
	}
}

/*
*/
SnpObjectActionHelper::VPB::~VPB()
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIActionManager);
		result = sAIActionManager->AIDeleteActionParamValue(fActionParamValueRef);
		SDK_ASSERT(!result);
		fActionParamValueRef = NULL;
	}
	catch (ai::Error) {
	}
}



/*
*/
ASErr SnpObjectActionHelper::CopySelection()
{
	ASErr result = kNoErr;
	try {
		SnpObjectActionHelper::VPB vpb;

		result = this->CopySelection(kDialogOff, vpb);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;

}

/*
*/
ASErr SnpObjectActionHelper::CopySelection(ActionDialogStatus dialogStatus, AIActionParamValueRef parameters)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIActionManager);
		result = sAIActionManager->PlayActionEvent(kAICopySelectionAction, dialogStatus, parameters);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;

}

/*
*/
ASErr SnpObjectActionHelper::PasteClipboard()
{
	ASErr result = kNoErr;
	try {
		SnpObjectActionHelper::VPB vpb;

		result = this->PasteClipboard(kDialogOff, vpb);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;

}

/*
*/
ASErr SnpObjectActionHelper::PasteClipboard(ActionDialogStatus dialogStatus, AIActionParamValueRef parameters)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIActionManager);
		result = sAIActionManager->PlayActionEvent(kAIPasteClipboardAction, dialogStatus, parameters);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;

}


// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpObjectActionHelper available to the SnippetRunner framework.
 */
class _SnpRunnableObjectActionHelper : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableObjectActionHelper () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableObjectActionHelper () {}

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
std::string _SnpRunnableObjectActionHelper::GetName() const
{
	return "ObjectActionHelper";
}

/*
*/
std::string _SnpRunnableObjectActionHelper::GetDescription() const 
{
	return "Helper class for programming object actions";
}

/*
*/
SnpRunnable::Operations _SnpRunnableObjectActionHelper::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("CopySelection", "document containing art", kSnpRunPathContext));
	operations.push_back(Operation("PasteClipboard", "document", kSnpRunNewDocumentContext));
	return operations;
}

/*
*/
std::string _SnpRunnableObjectActionHelper::GetDefaultOperationName() const 
{
	return "CopySelection";
}

/*
*/
std::vector<std::string> _SnpRunnableObjectActionHelper::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Helper Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableObjectActionHelper::CanRun(SnpRunnable::Context& runnableContext)
{
	SnpSelectionHelper selectionHelper;
	if (!selectionHelper.IsDocumentSelected())
		return false;

	if ("CopySelection" == runnableContext.GetOperation().GetName()) {
		return selectionHelper.IsArtSelected();
	}
	else
		return true;
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableObjectActionHelper::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpObjectActionHelper instance;
	if ("CopySelection" == runnableContext.GetOperation().GetName()) {
		result = instance.CopySelection();
	}
	else if ("PasteClipboard" == runnableContext.GetOperation().GetName()) {
		result = instance.PasteClipboard();
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableObjectActionHelper instance_SnpRunnableObjectActionHelper;

// End SnpObjectActionHelper.cpp

