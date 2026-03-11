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
#include "SnpDocumentHelper.h"
#include "SnpSelectionHelper.h"

using namespace ATE;

/** Provides functions that demonstrate how to create, modify, delete, 
	clear, apply and inspect character and paragraph styles.

	@see AIATECurrentTextFeaturesSuite, AIDocumentSuite, IDocumentTextResources,
	ICharStyle, ICharFeatures, IParaStyle, IParaFeatures, ITextRanges, ITextRange
  	@ingroup sdk_snippet
 */
class SnpTextStyles
{
	public:

		/** Constructor.
		*/
		SnpTextStyles() {}

		/** Destructor.
		*/
		virtual ~SnpTextStyles() {}

		/** Changes the font size and the character rotation of the normal 
			character style.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr ModifyNormalCharacterStyle();

		/** Iterates the set of character styles in the document reporting
			each style name to the SnippetRunner log.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr IterateCharacterStyles();

		/** Creates a new character style called 'SnpTextStyle'.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr CreateCharacterStyle();

		/** Creates 'SnpTextStyle' if not created and applies it to the
			selected text.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr ApplyCharacterStyle();

		/** Removes the association of the selected text range with its 
			character style	using ITextRanges::ClearNamedCharStyle.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr ClearCharacterStyle();

		/** Removes a character style from the current document	character 
			styles.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr DeleteCharacterStyle();

		/**	Writes the name of the character style that will be used on new text 
			items to the log.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr GetCurrentCharacterStyle();

		/** Changes the indent of the normal paragraph style.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr ModifyNormalParagraphStyle();

		/** Iterates the set of paragraph styles in the document reporting
			each style name to the SnippetRunner log.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr IterateParagraphStyles();

		/** Creates a new paragraph style called 'SnpTextStyle'.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr CreateParagraphStyle();

		/** Creates 'SnpTextStyle' if not created and applies it to the 
			selected paragraph.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr ApplyParagraphStyle();

		/** Removes the association of the selected text range with its 
			paragraph style	using ITextRanges::ClearNamedParaStyle.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr ClearParagraphStyle();

		/** Removes a paragraph style from the current document	paragraph
			styles.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr DeleteParagraphStyle();

		/**	Writes the name of the paragraph style that will be used on new text 
			items to the log.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr GetCurrentParagraphStyle();

};

/*
*/
ASErr SnpTextStyles::ModifyNormalCharacterStyle()
{
	ASErr result = kNoErr;
	try {
		DocumentTextResourcesRef docTextResourcesRef = NULL;
		result = sAIDocument->GetDocumentTextResources(&docTextResourcesRef);
		aisdk::check_ai_error(result);
		
		IDocumentTextResources resources(docTextResourcesRef);
		ICharStyle normalCharStyle = resources.GetNormalCharStyle();
		ICharFeatures oldFeatures = normalCharStyle.GetFeatures();
				
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		log->Write("Before modification:");
		log->BeginIndent();
		bool isAssigned = false;
		log->WritePrintf("Font size: %.1f", oldFeatures.GetFontSize(&isAssigned));
		log->WritePrintf("Character rotation value: %.1f", oldFeatures.GetCharacterRotation(&isAssigned));
		log->EndIndent();

		ICharFeatures features;
		features.SetFontSize(22.0);
		features.SetCharacterRotation(185.0);
		normalCharStyle.ReplaceOrAddFeatures(features);

		log->Write("After modification:");
		log->BeginIndent();
		log->WritePrintf("Font size: %.1f", features.GetFontSize(&isAssigned));
		log->WritePrintf("Character rotation value: %.1f", features.GetCharacterRotation(&isAssigned));
		log->EndIndent();
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch (ATE::Exception& ex) {
		result = ex.error;
	}
	return result;
}

/*
*/
ASErr SnpTextStyles::IterateCharacterStyles()
{
	ASErr result = kNoErr;
	try {
		// Get the document text resources.
		DocumentTextResourcesRef docTextResourcesRef = NULL;
		result = sAIDocument->GetDocumentTextResources(&docTextResourcesRef);
		aisdk::check_ai_error(result);		
		IDocumentTextResources resources(docTextResourcesRef);

		// Iterate the document character style set.
		ICharStyles charStyles = resources.GetCharStylesInDocument();
		ICharStylesIterator charStylesIter(charStyles);
		charStylesIter.MoveToFirst();
		while (charStylesIter.IsNotDone()) {
			// Get the name of the current style.
			ICharStyle charStyle = charStylesIter.Item();			
			ASUnicode* styleName = new ASUnicode[256];
			charStyle.GetName(styleName, 256);

			// Write the name to the log.
			ai::UnicodeString nameAsUnicode(styleName);
			SnippetRunnerLog::Instance()->Write(nameAsUnicode);
			charStylesIter.Next();
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch (ATE::Exception& ex) {
		result = ex.error;
	}
	return result;
}

/*
*/
ASErr SnpTextStyles::CreateCharacterStyle()
{
	ASErr result = kNoErr;
	try {		
		DocumentTextResourcesRef docTextResourcesRef = NULL;
		result = sAIDocument->GetDocumentTextResources(&docTextResourcesRef);
		aisdk::check_ai_error(result);

		IDocumentTextResources resources(docTextResourcesRef);
		const ASUnicode	*charStyleName = ai::UnicodeString("SnpTextStyle", kAIPlatformCharacterEncoding).as_ASUnicode().c_str();
		ICharStyle newCharStyle = resources.CreateCharStyle(charStyleName);
		if (!newCharStyle.IsNull()) {
			ICharFeatures features;
			features.SetUnderlinePosition(kUnderlineOn_RightInVertical);
			features.SetHorizontalScale(2.0);
			newCharStyle.SetFeatures(features);
		}
		else {
			SnippetRunnerLog::Instance()->Write("Character style already exists.");
			aisdk::check_ai_error(kBadParameterErr);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch (ATE::Exception& ex) {
		result = ex.error;
	}
	return result;
}

/*
*/
ASErr SnpTextStyles::ApplyCharacterStyle()
{
	ASErr result = kNoErr;
	try {
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);
		if (ranges.GetSize() > 0) {		
			DocumentTextResourcesRef docTextResourcesRef = NULL;
			result = sAIDocument->GetDocumentTextResources(&docTextResourcesRef);
			aisdk::check_ai_error(result);

			IDocumentTextResources resources(docTextResourcesRef);
			const ASUnicode	*charStyleName = ai::UnicodeString("SnpTextStyle", kAIPlatformCharacterEncoding).as_ASUnicode().c_str();
			ICharStyle newCharStyle = resources.GetCharStyle(charStyleName);
			if (newCharStyle.IsNull()) {
				// Create the character style.
				result = this->CreateCharacterStyle();
				aisdk::check_ai_error(result);
			}

			for (ASInt32 rangeIndex = 0; rangeIndex < ranges.GetSize(); rangeIndex++) {
				ITextRange range = ranges.Item(rangeIndex);
				if (range.GetSize() > 0) {
					const ASUnicode	*charStyleName = ai::UnicodeString("SnpTextStyle", kAIPlatformCharacterEncoding).as_ASUnicode().c_str();
					range.SetNamedCharStyle(charStyleName);
				}
				else {
					SnippetRunnerLog::Instance()->Write("No characters selected.");
					aisdk::check_ai_error(kBadParameterErr);
				}
			}
		}
		else {
			SnippetRunnerLog::Instance()->Write("Select a text range.");
			aisdk::check_ai_error(kBadParameterErr);
		}
	}
	catch(ai::Error& ex) {
		result = ex;
	}
	catch (ATE::Exception& ex) {
		result = ex.error;
	}
	return result;
}

/*
*/
ASErr SnpTextStyles::ClearCharacterStyle()
{
	ASErr result = kNoErr;
	try {
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);	
		if (ranges.GetSize() > 0) {
			ITextRange textRange = ranges.Item(0);
			if (textRange.GetSize() > 0) {
				ranges.ClearNamedCharStyle();
				ranges.ClearLocalCharFeatures();
			}
			else {
				SnippetRunnerLog::Instance()->Write("No characters selected.");
				aisdk::check_ai_error(result);
			}
		}
		else {
			SnippetRunnerLog::Instance()->Write("Select a text range.");
			aisdk::check_ai_error(result);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch (ATE::Exception& ex) {
		result = ex.error;
	}
	return result;
}

/*
*/
ASErr SnpTextStyles::DeleteCharacterStyle()
{
	ASErr result = kNoErr;
	try {
		// Get the current document text resources.
		DocumentTextResourcesRef docTextResourcesRef;
		result = sAIDocument->GetDocumentTextResources(&docTextResourcesRef);
		aisdk::check_ai_error(result);
		IDocumentTextResources resources(docTextResourcesRef);

		// Get the name of the character style to delete.
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ai::UnicodeString styleName = parameter->GetString(ai::UnicodeString("Style Name:"), ai::UnicodeString("SnpTextStyle"));

		// Delete the character style.
		bool deleted = resources.RemoveCharStyle((ASUnicode*)styleName.as_ASUnicode().c_str());
		if (!deleted) {
			SnippetRunnerLog::Instance()->Write("Style delete failed");
			aisdk::check_ai_error(kBadParameterErr);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch (ATE::Exception& ex) {
		result = ex.error;
	}
	return result;
}

/*
*/
ASErr SnpTextStyles::GetCurrentCharacterStyle()
{
	ASErr result = kNoErr;
	try {
		// Get the reference to the current char style.
		CharStyleRef charStyleRef;
		result = sAIATECurrentTextFeatures->GetCurrentCharStyle(&charStyleRef);

		// Create a new ICharStyle object using the CharStyleRef.
		ICharStyle charStyle(charStyleRef);

		// Get the name of the current style.
		ASUnicode* styleName = new ASUnicode[256];
		charStyle.GetName(styleName, 256);

		// Write the name to the log.
		ai::UnicodeString nameAsUnicode(styleName);
		SnippetRunnerLog::Instance()->Write(nameAsUnicode);
	}
	catch  (ai::Error& ex) {
		result = ex;
	}
	catch (ATE::Exception& ex) {
		result = ex.error;
	}
	return result;
}

/*
*/
ASErr SnpTextStyles::ModifyNormalParagraphStyle()
{
	ASErr result = kNoErr;
	try {
		DocumentTextResourcesRef docTextResourcesRef = NULL;
		result = sAIDocument->GetDocumentTextResources(&docTextResourcesRef);
		aisdk::check_ai_error(result);
		
		IDocumentTextResources resources(docTextResourcesRef);
		IParaStyle normalParaStyle = resources.GetNormalParaStyle();
		
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		log->Write("Before modification:");
		log->BeginIndent();
		bool isAssigned = false;
		log->WritePrintf("Start indent: %.1f", normalParaStyle.GetFeatures().GetStartIndent(&isAssigned));
		log->EndIndent();

		IParaFeatures features;
		features.SetStartIndent(10.0);
		normalParaStyle.ReplaceOrAddFeatures(features);

		log->Write("After modification:");
		log->BeginIndent();
		log->WritePrintf("Start indent: %.1f", normalParaStyle.GetFeatures().GetStartIndent(&isAssigned));
		log->EndIndent();
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch (ATE::Exception& ex) {
		result = ex.error;
	}
	return result;
}

/*
*/
ASErr SnpTextStyles::IterateParagraphStyles()
{
	ASErr result = kNoErr;
	try {
		// Get the document text resources.
		DocumentTextResourcesRef docTextResourcesRef = NULL;
		result = sAIDocument->GetDocumentTextResources(&docTextResourcesRef);
		aisdk::check_ai_error(result);		
		IDocumentTextResources resources(docTextResourcesRef);

		// Iterate the document paragraph style set.
		IParaStyles paraStyles = resources.GetParaStylesInDocument();
		IParaStylesIterator paraStylesIter(paraStyles);
		paraStylesIter.MoveToFirst();
		while (paraStylesIter.IsNotDone()) {
			// Get the name of the current style.
			IParaStyle paraStyle = paraStylesIter.Item();			
			ASUnicode* styleName = new ASUnicode[256];
			paraStyle.GetName(styleName, 256);

			// Write the name to the log.
			ai::UnicodeString nameAsUnicode(styleName);
			SnippetRunnerLog::Instance()->Write(nameAsUnicode);
			paraStylesIter.Next();
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch (ATE::Exception& ex) {
		result = ex.error;
	}
	return result;
}

/*
*/
ASErr SnpTextStyles::CreateParagraphStyle()
{
	ASErr result = kNoErr;
	try {
		DocumentTextResourcesRef docTextResourcesRef = NULL;
		result = sAIDocument->GetDocumentTextResources(&docTextResourcesRef);
		aisdk::check_ai_error(result);

		IDocumentTextResources resources(docTextResourcesRef);
		const ASUnicode	*paraStyleName = ai::UnicodeString("SnpTextStyle", kAIPlatformCharacterEncoding).as_ASUnicode().c_str();
		IParaStyle newParaStyle = resources.CreateParaStyle(paraStyleName);
		if (!newParaStyle.IsNull()) {
			IParaFeatures features;
			features.SetJustification(kCenterJustify);
			features.SetDesiredWordSpacing(2.0);
			newParaStyle.SetFeatures(features);
		}
		else {
			SnippetRunnerLog::Instance()->Write("Paragraph style already exists.");
			aisdk::check_ai_error(kBadParameterErr);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch (ATE::Exception& ex) {
		result = ex.error;
	}
	return result;
}

/*
*/
ASErr SnpTextStyles::ApplyParagraphStyle()
{
	ASErr result = kNoErr;
	try {
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);
		if (ranges.GetSize() > 0) {
			DocumentTextResourcesRef docTextResourcesRef = NULL;
			result = sAIDocument->GetDocumentTextResources(&docTextResourcesRef);
			aisdk::check_ai_error(result);

			IDocumentTextResources resources(docTextResourcesRef);
			const ASUnicode	*paraStyleName = ai::UnicodeString("SnpTextStyle", kAIPlatformCharacterEncoding).as_ASUnicode().c_str();
			IParaStyle newParaStyle = resources.GetParaStyle(paraStyleName);
			if (newParaStyle.IsNull()) {
				// Create the paragraph style.
				result = this->CreateParagraphStyle();
				aisdk::check_ai_error(result);
			}
			for (ASInt32 rangeIndex = 0; rangeIndex < ranges.GetSize(); rangeIndex++) {
				ITextRange range = ranges.Item(rangeIndex);
				const ASUnicode	*paraStyleName = ai::UnicodeString("SnpTextStyle", kAIPlatformCharacterEncoding).as_ASUnicode().c_str();
				range.SetNamedParaStyle(paraStyleName);
			}
		}
		else {
			SnippetRunnerLog::Instance()->Write("Select a text range.");
			aisdk::check_ai_error(kBadParameterErr);
		}
	}
	catch (ai::Error& ex) {
		result= ex;
	}
	catch (ATE::Exception& ex) {
		result = ex.error;
	}
	return result;
}

/*
*/
ASErr SnpTextStyles::ClearParagraphStyle()
{
	ASErr result = kNoErr;
	try {
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);	
		if (ranges.GetSize() > 0) {
			ITextRange textRange = ranges.Item(0);
			if (textRange.GetSize() > 0) {
				ranges.ClearNamedParaStyle();
				ranges.ClearLocalParaFeatures();
			}
			else {
				SnippetRunnerLog::Instance()->Write("No paragraphs selected.");
				aisdk::check_ai_error(kBadParameterErr);
			}
		}
		else {
			SnippetRunnerLog::Instance()->Write("Select a text range.");
			aisdk::check_ai_error(kBadParameterErr);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch (ATE::Exception& ex) {
		result = ex.error;
	}
	return result;
}

/*
*/
ASErr SnpTextStyles::DeleteParagraphStyle()
{
	ASErr result = kNoErr;
	try {
		// Get the current document text resources.
		DocumentTextResourcesRef docTextResourcesRef;
		result = sAIDocument->GetDocumentTextResources(&docTextResourcesRef);
		aisdk::check_ai_error(result);
		IDocumentTextResources resources(docTextResourcesRef);

		// Get the name of the paragraph style to delete.
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ai::UnicodeString styleName = parameter->GetString(ai::UnicodeString("Style Name:"), ai::UnicodeString("SnpTextStyle"));

		// Delete the style.
		bool deleted = resources.RemoveParaStyle((ASUnicode*)styleName.as_ASUnicode().c_str());
		if(!deleted) {
			SnippetRunnerLog::Instance()->Write("Style delete failed");
			aisdk::check_ai_error(kBadParameterErr);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch (ATE::Exception& ex) {
		result = ex.error;
	}
	return result;
}

/*
*/
ASErr SnpTextStyles::GetCurrentParagraphStyle()
{
	ASErr result = kNoErr;
	try {
		// Get the reference to the current para style.
		ParaStyleRef paraStyleRef;
		result = sAIATECurrentTextFeatures->GetCurrentParaStyle(&paraStyleRef);
		aisdk::check_ai_error(result);

		// Create a new IParaStyle object using the ParaStyleRef.
		IParaStyle paraStyle(paraStyleRef);

		// Get the name of the current style.
		ASUnicode* styleName = new ASUnicode[256];
		paraStyle.GetName(styleName, 256);

		// Write the name to the log.
		ai::UnicodeString nameAsUnicode(styleName);
		SnippetRunnerLog::Instance()->Write(nameAsUnicode);
	}
	catch  (ai::Error& ex) {
		result = ex;
	}
	catch (ATE::Exception& ex) {
		result = ex.error;
	}
	return result;
}


// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpTextStyles available to the SnippetRunner framework.
 */
class _SnpRunnableTextStyles : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableTextStyles () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableTextStyles () {}

		/* Returns name of snippet.
		*/
		std::string GetName() const;

		/* Returns a description of what the snippet does.
		*/
		std::string	GetDescription() const;

		/** Returns operations supported by this snippet.
		*/
		Operations GetOperations() const;

		/** Returns name of this snippet's default operation - must
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

		/** Sets the context manually for DeleteCharacterStyle and DeleteParagraphStyle.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr SetupContext(SnpRunnable::Context& runnableContext);
};

/*
*/
std::string _SnpRunnableTextStyles::GetName() const
{
	return "TextStyles";
}

/*
*/
std::string _SnpRunnableTextStyles::GetDescription() const 
{
	return "Demonstrates manipulation of text styles.";
}

/*
*/
SnpRunnable::Operations _SnpRunnableTextStyles::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("ModifyNormalCharacterStyle", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("IterateCharacterStyles", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("CreateCharacterStyle", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("ApplyCharacterStyle", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("ClearCharacterStyle", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("DeleteCharacterStyle", "selected text", kSnpRunCustomContext));
	operations.push_back(Operation("GetCurrentCharacterStyle", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("ModifyNormalParagraphStyle", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("IterateParagraphStyles", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("CreateParagraphStyle", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("ApplyParagraphStyle", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("ClearParagraphStyle", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("DeleteParagraphStyle", "selected text", kSnpRunCustomContext));
	operations.push_back(Operation("GetCurrentParagraphStyle", "document", kSnpRunNewDocumentContext));

	return operations;
}

/*
*/
std::string _SnpRunnableTextStyles::GetDefaultOperationName() const 
{
	return "ModifyNormalCharacterStyle";
	// Note: this must be one of the operation names returned by GetOperations.
}

/*
*/
std::vector<std::string> _SnpRunnableTextStyles::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Text Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableTextStyles::CanRun(SnpRunnable::Context& runnableContext)
{
	AIErr result = kNoErr;
	try {
		SnpSelectionHelper selectionHelper;
		if (!selectionHelper.IsDocumentSelected())
			return false;

		DocumentTextResourcesRef docTextResources = NULL;
		result = sAIDocument->GetDocumentTextResources(&docTextResources);
		aisdk::check_ai_error(result);
		IDocumentTextResources resources(docTextResources);

		if ("DeleteCharacterStyle" == runnableContext.GetOperation().GetName()) {
			ICharStyles styleSet = resources.GetCharStylesInDocument();
			return styleSet.GetSize() > 1;
		}
		else if ("DeleteParagraphStyle" == runnableContext.GetOperation().GetName()) {
			IParaStyles styleSet = resources.GetParaStylesInDocument();
			return styleSet.GetSize() > 1;
		}	
		else if ("ApplyCharacterStyle" == runnableContext.GetOperation().GetName() ||
				"ClearCharacterStyle" == runnableContext.GetOperation().GetName() ||
				"ApplyParagraphStyle" == runnableContext.GetOperation().GetName() ||
				"ClearParagraphStyle" == runnableContext.GetOperation().GetName()) {
			
			return selectionHelper.IsTextRangeSelected();
		}
		else
			return true;
	}
	catch (ai::Error& ex) {
		result = ex;
		return false;
	}
	catch (ATE::Exception& ex) {
		result = ex.error;
		return false;
	}
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableTextStyles::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpTextStyles instance;
	if ("ModifyNormalCharacterStyle" == runnableContext.GetOperation().GetName()) {
		result = instance.ModifyNormalCharacterStyle();
	}
	else if ("IterateCharacterStyles" == runnableContext.GetOperation().GetName()) {
		result = instance.IterateCharacterStyles();
	}
	else if ("CreateCharacterStyle" == runnableContext.GetOperation().GetName()) {
		result = instance.CreateCharacterStyle();
	}
	else if ("ApplyCharacterStyle" == runnableContext.GetOperation().GetName()) {
		result = instance.ApplyCharacterStyle();
	}
	else if ("ClearCharacterStyle" == runnableContext.GetOperation().GetName()) {
		result = instance.ClearCharacterStyle();
	}
	else if ("DeleteCharacterStyle" == runnableContext.GetOperation().GetName()) {
		result = instance.DeleteCharacterStyle();
	}
	else if ("GetCurrentCharacterStyle" == runnableContext.GetOperation().GetName()) {
		result = instance.GetCurrentCharacterStyle();
	}
	else if ("ModifyNormalParagraphStyle" == runnableContext.GetOperation().GetName()) {
		result = instance.ModifyNormalParagraphStyle();
	}
	else if ("IterateParagraphStyles" == runnableContext.GetOperation().GetName()) {
		result = instance.IterateParagraphStyles();
	}
	else if ("CreateParagraphStyle" == runnableContext.GetOperation().GetName()) {
		result = instance.CreateParagraphStyle();
	}
	else if ("ApplyParagraphStyle" == runnableContext.GetOperation().GetName()) {
		result = instance.ApplyParagraphStyle();
	}
	else if ("ClearParagraphStyle" == runnableContext.GetOperation().GetName()) {
		result = instance.ClearParagraphStyle();
	}
	else if ("DeleteParagraphStyle" == runnableContext.GetOperation().GetName()) {
		result = instance.DeleteParagraphStyle();
	}
	else if ("GetCurrentParagraphStyle" == runnableContext.GetOperation().GetName()) {
		result = instance.GetCurrentParagraphStyle();
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

ASErr _SnpRunnableTextStyles::SetupContext(SnpRunnable::Context &runnableContext)
{
	ASErr result = kNoErr;
	try {
		switch (runnableContext.GetOperation().GetUnitTestContext()) {
			case SnpRunnable::kSnpRunCustomContext:	{
					SnpTextStyles instance;
					SnpDocumentHelper docHelper;
					AIDocumentHandle document = NULL;
					SDK_ASSERT(!docHelper.NewDocument(ai::UnicodeString("Print"), document));
					SDK_ASSERT(document);

					if ("DeleteCharacterStyle" == runnableContext.GetOperation().GetName()) {
						// Unlink the text frames before continuing.
						instance.CreateCharacterStyle();
					}
					else if ("DeleteParagraphStyle" == runnableContext.GetOperation().GetName()){
						// Set the text focus before continuing.
						instance.CreateParagraphStyle();
					}
				}
				break;
			default: {
				result = SnpRunnable::SetupContext(runnableContext);
			}
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch (ATE::Exception& ex) {
		result = ex.error;
	}
	return result;
}


/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableTextStyles instance_SnpRunnableTextStyles;

// End SnpTextStyles.cpp

