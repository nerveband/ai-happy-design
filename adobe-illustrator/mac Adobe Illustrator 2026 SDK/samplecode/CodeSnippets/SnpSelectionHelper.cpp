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
#include "SnippetRunnerSuites.h"
#include "SnpArtSetHelper.h"
#include "SnpDocumentHelper.h"
#include "SnpChooser.h"
#include "SnpSelectionHelper.h"

/*
*/
ASBoolean SnpSelectionHelper::IsDocumentSelected(void)
{
	ASBoolean result = false;
	AIDocumentHandle document = NULL;
	SnpDocumentHelper docHelper;
	ASErr status = docHelper.GetCurrentDocument(document);
	result = !status && document;
	return result;
}

/*
*/
ASBoolean SnpSelectionHelper::IsArtSelected(AIArtType artType)
{
	ASBoolean result = false;
	if (this->IsDocumentSelected()) {
		SnpArtSetHelper selectedArt;
		selectedArt.FindSelectedArt();
		if (artType == kAnyArt) {
			result = selectedArt.GetCount() > 0;
		}
		else {
			selectedArt = selectedArt.Filter(artType);
			result = selectedArt.GetCount() > 0;
		}
	}
	return result;
}

/*
*/
ASBoolean SnpSelectionHelper::IsTextRangeSelected(void)
{
	ASBoolean result = false;
	if (this->IsDocumentSelected()) {
		ATE::TextRangesRef rangesRef = NULL;
		AIErr error = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(error);

		ATE::ITextRanges ranges(rangesRef); 
		if (ranges.GetSize() > 0) {
			AIBoolean textFocus = false;
			error = sAIDocument->HasTextFocus(&textFocus);
			aisdk::check_ai_error(error);

			if (textFocus) {
				ATE::ITextRange textRange = ranges.Item(0);
				if (textRange.GetSize() > 0)
					result = true;
			}
			else
				result = true;
		}
	}
	return result;
}


/*
*/
ASBoolean SnpSelectionHelper::CanSelectArt(AIArtHandle art)
{
	ASBoolean result = false;
	try {
		SDK_ASSERT(sAIArt);
		ai::int32 attr = 0;
		ASErr status = sAIArt->GetArtUserAttr(art, kArtLocked | kArtHidden, &attr);
		aisdk::check_ai_error(status);
		result = !(attr & kArtLocked) && !(attr & kArtHidden);
	}
	catch (ai::Error) {
	}
	return result;
}

/*
*/
ASErr SnpSelectionHelper::SelectArt(AIArtHandle art)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(this->CanSelectArt(art));
		result = sAIArt->SetArtUserAttr(art, kArtSelected, kArtSelected);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpSelectionHelper::DeselectAll()
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIMatchingArt);
		result = sAIMatchingArt->DeselectAll();
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpSelectionHelper available to the SnippetRunner framework.
 */
class _SnpRunnableSelectionHelper : public SnpRunnable
{
	public:
		/* Constructor registers the snippet to the framework.
		 */
		_SnpRunnableSelectionHelper () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual	~_SnpRunnableSelectionHelper () {}

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
std::string _SnpRunnableSelectionHelper::GetName() const
{
	return "SelectionHelper";
}

/*
*/
std::string _SnpRunnableSelectionHelper::GetDescription() const 
{
	return "Manipulates selection";
}

/*
*/
SnpRunnable::Operations _SnpRunnableSelectionHelper::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("IsDocumentSelected", "", kSnpRunNewDocumentContext));	
	operations.push_back(Operation("IsArtSelected", "", kSnpRunPathContext));	
	operations.push_back(Operation("IsTextRangeSelected", "", kSnpRunCharactersContext));	
	operations.push_back(Operation("SelectArtOnLayer", "document", kSnpRunUnitTestDocumentContext));	
	return operations;
}

/*
*/
std::string _SnpRunnableSelectionHelper::GetDefaultOperationName() const 
{
	return "IsDocumentSelected";
}

/*
*/
std::vector<std::string> _SnpRunnableSelectionHelper::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Helper Snippets");
	return categories;
}

/*
*/
ASBoolean _SnpRunnableSelectionHelper::CanRun(SnpRunnable::Context& runnableContext)
{
	SnpSelectionHelper selectionHelper;
	if ("SelectArtOnLayer" == runnableContext.GetOperation().GetName())
		return selectionHelper.IsDocumentSelected();
	else
		return true;
}

/*
*/
ASErr _SnpRunnableSelectionHelper::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr; 

	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnpSelectionHelper instance;
		if ("IsDocumentSelected" == runnableContext.GetOperation().GetName()) {
			if (instance.IsDocumentSelected()) {
				log->Write("yes");
			} else {
				log->Write("no");
			}
		}
		else if ("IsArtSelected" == runnableContext.GetOperation().GetName()) {
			if (instance.IsArtSelected()) {
				log->Write("yes");
			} else {
				log->Write("no");
			}
		}
		else if ("IsTextRangeSelected" == runnableContext.GetOperation().GetName()) {
			if (instance.IsTextRangeSelected()) {
				log->Write("yes");
			} else {
				log->Write("no");
			}
		}
		else if ("SelectArtOnLayer" == runnableContext.GetOperation().GetName()) {
			SnpChooser chooser;
			AILayerHandle layer = chooser.PickLayer();
			ai::UnicodeString title;
			result = sAILayer->GetLayerTitle(layer, title);
			if (!result) {
				SnpArtSetHelper instance;
				instance.FindLayerArt(title);
				instance.SelectAll();
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
static _SnpRunnableSelectionHelper instance_SnpRunnableSelectionHelper;

// End SnpSelectionHelper.cpp

