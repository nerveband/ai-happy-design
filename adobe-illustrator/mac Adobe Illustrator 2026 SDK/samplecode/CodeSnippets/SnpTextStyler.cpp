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

using namespace ATE;

/** Provides functions that demonstrate how to apply, clear and inspect 
	character and paragraph features.

	@see AIATECurrentTextFeaturesSuite, AIATEPaintSuite, AIDocumentSuite, ITextRanges,
	ITextRange, ICharFeatures, IApplicationPaint, IParaFeatures, ICharInspector, IFont,
	IParaInspector, IParagraphsIterator, ITabStops, ITabStop, IParagraph, IStory
  	@ingroup sdk_snippet
 */
class SnpTextStyler
{
	public:

		/** Constructor.
		*/
		SnpTextStyler() {}

		/** Destructor.
		*/
		virtual ~SnpTextStyler() {}

		/** Gets the fill color type and stroke width of the selected characters,
			gets the overriding features if they exist in the selection, the details 
			are	written to the SnippetRunner log window.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr GetCharacterFeatures();

		/** Sets the font size and fill color local character features of the 
			selected characters.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr ApplyLocalCharacterFeatures();

		/** Removes the local character features of the selected characters
			using ITextRanges::ClearLocalCharFeatures.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr ClearLocalCharacterFeatures();

		/**	Sets individual character features that will be applied to new text items.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr SetCurrentCharacterOverrides();

		/** Sets the justification local paragraph feature of the selected paragraphs.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr ApplyLocalParagraphFeatures();

		/** Removes the local paragraph features of the selected paragraphs
			using ITextRanges::ClearLocalParaFeatures.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr ClearLocalParagraphFeatures();

		/**	Sets individual paragraph features that will be used in new text items.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr SetCurrentParagraphOverrides();

		/** Gets all font names and sizes found in the selected characters
			using an ICharInspector and writes them to the SnippetRunner log window.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr InspectSelectedCharacterFeatures();

		/** Gets all the justifications and start indent values in the selected
			paragraphs using an IParaInspector and writes them to the SnippetRunner 
			log window.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr InspectSelectedParagraphFeatures();

		/** Adds new tab stops to the selected paragraphs.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr AddTabStops();

		/** Gets all tab stop positions of the selected text range and 
			writes them to the SnippetRunner log window.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr InspectSelectedParagraphTabStops();

		/** Removes all tab stops from the selected paragraphs.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr RemoveTabStops();

		/** Sets the kern type of all selected characters to kOpticalKern.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr SetKernType();

};

/*
*/
ASErr SnpTextStyler::GetCharacterFeatures()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();	
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);		

		ITextRanges ranges(rangesRef);
		if (ranges.GetSize() > 0) {
			ITextRange textRange = ranges.Item(0);
			if (textRange.GetSize() > 0) {
				ICharFeatures features;
				if (ranges.HasLocalCharFeatures()) {
					features = textRange.GetUniqueLocalCharFeatures();
					log->Write("Overriding features");
				}
				else {
					features = textRange.GetUniqueCharFeatures();
				}
				bool isValid = false;
				IApplicationPaint fillColor = features.GetFillColor(&isValid);
				AIPathStyle style;
				AIPathStyleMap styleMap;
				SnippetRunnerLog::Indent indent;
				if (isValid) {
					result = sAIATEPaint->GetAIPathStyleAndMap(features.GetRef(), &style, &styleMap);
					aisdk::check_ai_error(result);

					log->Write("Fill color type: ");
					SnippetRunnerLog::Indent indent2;
					AIColor color = style.fill.color;
					AIColorTag  fillColorTag = color.kind;
					switch (fillColorTag) {
					case kGrayColor:
						log->Write("kGrayColor");
						break;
					case kFourColor:
						log->Write("kFourColor");
						break;
					case kPattern:
						log->Write("kPattern");
						break;
					case kCustomColor:
						log->Write("kCustomColor");
						break;
					case kGradient:
						log->Write("kGradient");
						break;
					case kThreeColor:
						log->Write("kThreeColor");
						break;
					case kNoneColor:
						log->Write("kNoneColor");
						break;
					}
				}
				else {
					log->Write("More than one color type in selection.");
				}
				isValid = false;
				bool strokeApplied = features.GetStroke(&isValid);
				if (isValid && strokeApplied) {
					result = sAIATEPaint->GetAIPathStyleAndMap(features.GetRef(), &style, &styleMap);
					aisdk::check_ai_error(result);

					log->WritePrintf("Stroke width: %.1f pt", style.stroke.width);
				}
				else {
					log->Write("Stroke width: No stroke applied.");
				}
			}
			else {
				log->Write("Select one or more characters.");
				aisdk::check_ai_error(kBadParameterErr);
			}
		}
		else {
			log->Write("Select a text range.");
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
ASErr SnpTextStyler::ApplyLocalCharacterFeatures()
{
	ASErr result = kNoErr;
	try {
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);
		if (ranges.GetSize() > 0) {
			ICharFeatures features;
			features.SetFontSize(23.0);
			AIColor newColor;			
			newColor.kind = kFourColor;
			newColor.c.f.cyan = 1.0;
			newColor.c.f.magenta = 0.0;
			newColor.c.f.yellow = 0.0;
			newColor.c.f.black = 0.0;
			ApplicationPaintRef paintRef = NULL;
			result = sAIATEPaint->CreateATEApplicationPaint(&newColor, &paintRef);
			aisdk::check_ai_error(result);
			IApplicationPaint paint(paintRef);
			features.SetFillColor(paint);
			ranges.SetLocalCharFeatures(features);
		}
		else {
			SnippetRunnerLog::Instance()->Write("Select some characters.");
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
ASErr SnpTextStyler::ClearLocalCharacterFeatures()
{
	ASErr result = kNoErr;
	try {
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);	
		if (ranges.GetSize() > 0) {
			ranges.ClearLocalCharFeatures();
		}
		else {
			SnippetRunnerLog::Instance()->Write("Select a range of text.");
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
ASErr SnpTextStyler::SetCurrentCharacterOverrides()
{
	ASErr result = kNoErr;
	try {
		// Create new ICharFeatures object.
		ICharFeatures features;
		features.SetHorizontalScale(5.0);
		features.SetVerticalScale(5.0);
		features.SetAutoKernType(kOpticalKern);

		// Get the CharFeaturesRef for the ICharFeatures object.
		CharFeaturesRef featuresRef = features.GetRef();

		// Apply the overrides.
		result = sAIATECurrentTextFeatures->SetCurrentCharOverrides(featuresRef);
		aisdk::check_ai_error(result);
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
ASErr SnpTextStyler::ApplyLocalParagraphFeatures()
{
	ASErr result = kNoErr;
	try {
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);
		if (ranges.GetSize() > 0) {
			IParaFeatures features;
			features.SetJustification(kCenterJustify);
			ranges.SetLocalParaFeatures(features);
		}
		else {
			SnippetRunnerLog::Instance()->Write("Select a paragraph.");
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
ASErr SnpTextStyler::ClearLocalParagraphFeatures()
{
	ASErr result = kNoErr;
	try {
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);
		if (ranges.GetSize() > 0) {
			ranges.ClearLocalParaFeatures();
		}
		else {
			SnippetRunnerLog::Instance()->Write("Select a paragraph.");
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
ASErr SnpTextStyler::SetCurrentParagraphOverrides()
{
	ASErr result = kNoErr;
	try {
		// Create the new IParaFeatures object.
		IParaFeatures features;
		features.SetMinWordSpacing(10.0);
		features.SetJustification(kCenterJustify);

		// Get the ParaFeaturesRef for the IParaFeatures object.
		ParaFeaturesRef featuresRef = features.GetRef();

		// Apply the overrides.
		result = sAIATECurrentTextFeatures->SetCurrentParaOverrides(featuresRef);
		aisdk::check_ai_error(result);
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
ASErr SnpTextStyler::InspectSelectedCharacterFeatures()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		// Get each font and font size found in selection and list in the log file.
		ITextRanges ranges(rangesRef);	
		IArrayReal fontSizeArray;
		IArrayFontRef fontRefArray;
		ICharInspector inspector = ranges.GetCharInspector();		
		fontSizeArray = inspector.GetFontSize();
		fontRefArray = inspector.GetFont();
		ASInt32 countFontRefs = fontRefArray.GetSize();
		ASInt32 countFontSizes = fontSizeArray.GetSize();
		if (countFontRefs > 0 && countFontSizes > 0) {
			log->Write("Fonts found in selection:");
			for (ASInt32 i = 0; i < countFontRefs; i++) {
				IFont font = fontRefArray.Item(i);
				FontRef fontRef = font.GetRef();
				AIFontKey fontKey = NULL;
				result = sAIFont->FontKeyFromFont(fontRef, &fontKey);
				aisdk::check_ai_error(result);

				ai::UnicodeString fontName;
				result = sAIFont->GetUserFontNameUnicode(fontKey, fontName, true);
				aisdk::check_ai_error(result);
				SnippetRunnerLog::Indent indent;
				log->Write(fontName);
			}

			log->Write("Font sizes found in selection:");
			for (ASInt32 i = 0; i < countFontSizes; i++) {
				SnippetRunnerLog::Indent indent;
				log->WritePrintf("%.1f%s",fontSizeArray.Item(i)," pt");
			}
		}
		else {
			log->Write("Select characters.");
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
ASErr SnpTextStyler::InspectSelectedParagraphFeatures()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);	
		ASInt32 selectedRanges = ranges.GetSize();
		if (selectedRanges > 0) {
			IParaInspector inspector = ranges.GetParaInspector();
			IArrayParagraphJustification paraJustifications = inspector.GetJustification();
			log->Write("Paragraph justifications found:");
			for (ASInt32 i = 0; i < paraJustifications.GetSize(); i++) {
				SnippetRunnerLog::Indent indent;
				ParagraphJustification justification = paraJustifications.Item(i);
				switch (justification) {
				case kLeftJustify:
					log->Write("kLeftJustify");
					break;
				case kRightJustify:
					log->Write("kRightJustify");
					break;
				case kCenterJustify:
					log->Write("kCenterJustify");
					break;
				case kFullJustifyLastLineLeft:
					log->Write("kFullJustifyLastLineLeft");
					break;
				case kFullJustifyLastLineRight:
					log->Write("kFullJustifyLastLineRight");
					break;
				case kFullJustifyLastLineCenter:
					log->Write("kFullJustifyLastLineCenter");
					break;
				case kFullJustifyLastLineFull:
					log->Write("kFullJustifyLastLineFull");
					break;
				}
			}
			IArrayReal firstLineIndents = inspector.GetFirstLineIndent();
			log->Write("First line indents found:");
			for (ASInt32 i = 0; i < firstLineIndents.GetSize(); i++) {
				SnippetRunnerLog::Indent indent;
				log->WritePrintf("%.1f%s", firstLineIndents.Item(i), " pt");
			}
		}
		else {
			log->Write("Select paragraphs.");
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
ASErr SnpTextStyler::AddTabStops(void)
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);	
		if (ranges.GetSize() > 0) {
			// Create tab stops.
			ITabStop tabStop1, tabStop2;
			tabStop1.SetPosition(10.0);
			tabStop2.SetPosition(20.0);
			ITabStops tabStopsSet;
			tabStopsSet.ReplaceOrAdd(tabStop1);
			tabStopsSet.ReplaceOrAdd(tabStop2);

			IParagraphsIterator paraIter = ranges.GetParagraphsIterator();
			paraIter.MoveToFirst();
			while (paraIter.IsNotDone()) {
				IParagraph para = paraIter.Item();
				ITextRange textRange = para.GetTextRange();
				IParaFeatures features = textRange.GetUniqueParaFeatures();
				features.SetTabStops(tabStopsSet);
				textRange.SetLocalParaFeatures(features);
				paraIter.Next();
			}
		}
		else {
			log->Write("Select at least 1 paragraph.");
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
ASErr SnpTextStyler::InspectSelectedParagraphTabStops()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);	
		if (ranges.GetSize() > 0) {
			IParagraphsIterator paraIter = ranges.GetParagraphsIterator();
			paraIter.MoveToFirst();
			ASInt32 paraCount = 0;
			while (paraIter.IsNotDone()) {
				IParagraph currentPara = paraIter.Item();
				ITextRange textRange = currentPara.GetTextRange();
				IParaInspector inspector = textRange.GetParaInspector();
				IArrayTabStopsRef tabStopsRefArray = inspector.GetTabStops();
				ASInt32 refArraySize = tabStopsRefArray.GetSize();
				if (refArraySize > 0) {
					for (ASInt32 i = 0; i < refArraySize; i++) {
						SnippetRunnerLog::Indent indent;
						ITabStops tabStopsArray = tabStopsRefArray.Item(i);
						ASInt32 tabArraySize = tabStopsArray.GetSize();
						if (tabArraySize > 0) {
							for (ASInt32 j = 0; j < tabArraySize; j++) {
								ITabStop tabStop = tabStopsArray.Item(j);
								AIReal tabPosition = tabStop.GetPosition();
								log->WritePrintf("Paragraph %d%s%d%s%.1f", paraCount, 
									", tab stop ", j, " position: ", tabPosition);
							}
						}
						else {
							log->WritePrintf("Paragraph %d%s", paraCount, ", no tab stops.");
						}
					}
				}
				paraCount++;
				paraIter.Next();
			}
		}
		else {
			log->Write("Select paragraphs.");
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
ASErr SnpTextStyler::RemoveTabStops(void)
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);	
		if (ranges.GetSize() > 0) {
			IParagraphsIterator paraIter = ranges.GetParagraphsIterator();
			paraIter.MoveToFirst();
			while (paraIter.IsNotDone()) {
				IParagraph para = paraIter.Item();
				ITextRange textRange = para.GetTextRange();
				IParaFeatures features = textRange.GetUniqueParaFeatures();
				bool isAssigned = false;
				ITabStops tabStopsSet = features.GetTabStops(&isAssigned);
				if (isAssigned) {
					tabStopsSet.RemoveAll();
					features.SetTabStops(tabStopsSet);
					textRange.SetLocalParaFeatures(features);
				}
				paraIter.Next();
			}
		}
		else {
			log->Write("Select the paragraph(s) to remove the tab stops from.");
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
ASErr SnpTextStyler::SetKernType(void)
{
	ASErr result = kNoErr;
	try {
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);	
		if (ranges.GetSize() > 0) {
			for (ASInt32 i = 0; i < ranges.GetSize(); i++) {
				ITextRange range = ranges.Item(i);
				IStory story = range.GetStory();
				story.SetKernForSelection(range, kOpticalKern);
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

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpTextStyler available to the SnippetRunner framework.
 */
class _SnpRunnableTextStyler : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableTextStyler () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableTextStyler () {}

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
};

/*
*/
std::string _SnpRunnableTextStyler::GetName() const
{
	return "TextStyler";
}

/*
*/
std::string _SnpRunnableTextStyler::GetDescription() const 
{
	return "Demonstrates manipulation of text features.";
}

/*
*/
SnpRunnable::Operations _SnpRunnableTextStyler::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("GetCharacterFeatures", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("ApplyLocalCharacterFeatures", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("ClearLocalCharacterFeatures", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("SetCurrentCharacterOverrides", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("ApplyLocalParagraphFeatures", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("ClearLocalParagraphFeatures", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("SetCurrentParagraphOverrides", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("InspectSelectedCharacterFeatures", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("InspectSelectedParagraphFeatures", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("AddTabStops", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("InspectSelectedParagraphTabStops", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("RemoveTabStops", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("SetKernType", "selected text", kSnpRunCharactersContext));

	return operations;
}

/*
*/
std::string _SnpRunnableTextStyler::GetDefaultOperationName() const 
{
	return "GetCharacterFeatures";
	// Note: this must be one of the operation names returned by GetOperations.
}

/*
*/
std::vector<std::string> _SnpRunnableTextStyler::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Text Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableTextStyler::CanRun(SnpRunnable::Context& runnableContext)
{
	SnpSelectionHelper selectionHelper;
	if (!selectionHelper.IsDocumentSelected())
		return false;

	if ("GetCharacterFeatures" == runnableContext.GetOperation().GetName() ||
		"ApplyLocalCharacterFeatures" == runnableContext.GetOperation().GetName() ||
		"ClearLocalCharacterFeatures" == runnableContext.GetOperation().GetName() ||
		"ApplyLocalParagraphFeatures" == runnableContext.GetOperation().GetName() ||
		"ClearLocalParagraphFeatures" == runnableContext.GetOperation().GetName() ||
		"InspectSelectedCharacterFeatures" == runnableContext.GetOperation().GetName() ||
		"InspectSelectedParagraphFeatures" == runnableContext.GetOperation().GetName() ||
		"AddTabStops" == runnableContext.GetOperation().GetName() ||
		"InspectSelectedParagraphTabStops" == runnableContext.GetOperation().GetName() ||
		"RemoveTabStops" == runnableContext.GetOperation().GetName() ||
		"SetKernType" == runnableContext.GetOperation().GetName()) {
		
			return selectionHelper.IsTextRangeSelected();
	}
	else
		return true;
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableTextStyler::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpTextStyler instance;
	if ("GetCharacterFeatures" == runnableContext.GetOperation().GetName()) {
		result = instance.GetCharacterFeatures();
	}
	else if ("ApplyLocalCharacterFeatures" == runnableContext.GetOperation().GetName()) {
		result = instance.ApplyLocalCharacterFeatures();
	}
	else if ("ClearLocalCharacterFeatures" == runnableContext.GetOperation().GetName()) {
		result = instance.ClearLocalCharacterFeatures();
	}
	else if ("SetCurrentCharacterOverrides" == runnableContext.GetOperation().GetName()) {
		result = instance.SetCurrentCharacterOverrides();
	}
	else if ("ApplyLocalParagraphFeatures" == runnableContext.GetOperation().GetName()) {
		result = instance.ApplyLocalParagraphFeatures();
	}
	else if ("ClearLocalParagraphFeatures" == runnableContext.GetOperation().GetName()) {
		result = instance.ClearLocalParagraphFeatures();
	}
	else if ("SetCurrentParagraphOverrides" == runnableContext.GetOperation().GetName()) {
		result = instance.SetCurrentParagraphOverrides();
	}
	else if ("InspectSelectedCharacterFeatures" == runnableContext.GetOperation().GetName()) {
		result = instance.InspectSelectedCharacterFeatures();
	}
	else if ("InspectSelectedParagraphFeatures" == runnableContext.GetOperation().GetName()) {
		result = instance.InspectSelectedParagraphFeatures();
	}
	else if ("AddTabStops" == runnableContext.GetOperation().GetName()) {
		result = instance.AddTabStops();
	}
	else if ("InspectSelectedParagraphTabStops" == runnableContext.GetOperation().GetName()) {
		result = instance.InspectSelectedParagraphTabStops();
	}
	else if ("RemoveTabStops" == runnableContext.GetOperation().GetName()) {
		result = instance.RemoveTabStops();
	}
	else if ("SetKernType" == runnableContext.GetOperation().GetName()) {
		result = instance.SetKernType();
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableTextStyler instance_SnpRunnableTextStyler;

// End SnpTextStyler.cpp

