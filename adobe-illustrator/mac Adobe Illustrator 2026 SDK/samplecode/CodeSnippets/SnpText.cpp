//========================================================================================
//  
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
#include "SnpArtSetHelper.h"
#include "SnpSelectionHelper.h"
#include "SnpArtHelper.h"

using namespace ATE;

/** Provides functions that demonstrate how to create, edit and delete 
	text items, perform spell checks, run find and replace and run hit tests.

	@see AITextFrameSuite, AILegacyTextConversionSuite, AIATETextUtilSuite, 
	AIDocumentSuite, AITextFrameHitSuite, AIHitTestSuite, ITextRange, ITextRanges,
	ITextFrame,	IStory, ISpell, IFind, IDocumentTextResources
	@ingroup sdk_snippet
 */
class SnpText
{
	public:

		/** Constructor.
		*/
		SnpText() {}

		/** Destructor.
		*/
		virtual ~SnpText() {}
		
		/** Creates a new point text item in the current document.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr CreatePointText();

		/** Creates a new in path text item in the current document.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr CreateInPathText();

		/** Creates a new on path text item in the current document.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr CreateOnPathText();

		/** Creates a new threaded in path text item in the current document.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr CreateThreadedInPathText();
		
		/** Adds a text string to the selected text range.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr InsertText();

		/** Clones the selected text and copies to a new text frame.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr CopyText();

		/** Moves the selected text to a new text frame.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr MoveText();

		/** Replaces the selected text with 'Text Replaced.'
			NOTE: The difference between this function and SnpText::
			FindAndReplace is this function replaces the selected text 
			with a text string and SnpText::FindAndReplace replaces all 
			occurances of a text string with another text string.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr ReplaceText();

		/** Selects all text in the selected text frame.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr SelectTextRange();

		/**	Deletes all selected text.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr DeleteTextRange();

		/** Converts all legacy text in the current document by calling
			AILegacyTextConversionSuite::ConvertAllToNative.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr ConvertLegacyText();

		/** Sets the story in the selected item to be in editing mode.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr SetTextFocus();

		/** Exits text editing mode using sAIDocument->LoseTextFocus.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr LoseTextFocus();

		/** Links the selected text frames as long as they are not already linked.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr LinkTextFrames();

		/** Removes any links from the selected text frames.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr UnlinkTextFrames();

		/**	Performs a spell check on the selected text items and writes 
			unknown words and the alternate suggestions from Illustrators
			dictionary to the log.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr RunSpellCheck();

		/**	Prompts the user for the text string to find and the text string 
			to replace the found string with, performs the find and replace 
			and writes the number of occurances replaced to the log.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr FindAndReplace();		

		/**	Prompt user for coordinates on the document then report whether a 
			text frame item was hit and if so what part.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr RunHitTest();

		/**	Iterates through every font available to the document and reports the
			name of each font to the SnippetRunner log.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr IterateAllFonts();

		/**	Iterates through each font used in the current document and reports the
			name of each font to the SnippetRunner log.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr IterateUsedFonts();

};

/*
*/
ASErr SnpText::CreatePointText()
{
	ASErr result = kNoErr;
	try {
		// Get the group art that contains all the art in the current layer.
		AIArtHandle artGroup = NULL;
		result = sAIArt->GetFirstArtOfLayer(NULL, &artGroup);
		aisdk::check_ai_error(result);

		// Add the new point text item to the layer.
		AITextOrientation orient = kHorizontalTextOrientation;
		AIRealPoint anchor = {10,10};
		AIArtHandle textFrame = NULL;
		result = sAITextFrame->NewPointText(kPlaceAboveAll, artGroup, orient, anchor, &textFrame);
		aisdk::check_ai_error(result);

		// Set the contents of the text range.
		TextRangeRef range = NULL;
		result = sAITextFrame->GetATETextRange(textFrame, &range);
		aisdk::check_ai_error(result);
		ITextRange crange(range);
		crange.InsertAfter( ai::UnicodeString("New point text item").as_ASUnicode().c_str() );
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
ASErr SnpText::CreateInPathText()
{
	ASErr result = kNoErr;
	try {
		// Get the group art that contains all the art in the current layer.
		AIArtHandle artGroup = NULL;
		result = sAIArt->GetFirstArtOfLayer(NULL, &artGroup);
		aisdk::check_ai_error(result);

		// Add the new in path text item to the layer.
		SnpArtHelper artHelper;
		AIArtHandle textFrame;
		result = artHelper.NewRectangle(textFrame);
		aisdk::check_ai_error(result);
		AITextOrientation orient = kHorizontalTextOrientation;
		result = sAITextFrame->NewInPathText(kPlaceAboveAll, artGroup, orient, textFrame, NULL, false, &textFrame);
		aisdk::check_ai_error(result);

		// Set the contents of the text range.
		TextRangeRef range = NULL;
		result = sAITextFrame->GetATETextRange(textFrame, &range);
		aisdk::check_ai_error(result);
		ITextRange iTextRange(range);
		iTextRange.InsertAfter(ai::UnicodeString("New in path text item").as_ASUnicode().c_str());
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
ASErr SnpText::CreateOnPathText()
{
	ASErr result = kNoErr;
	try {
		// Get the group art that contains all the art in the current layer.
		AIArtHandle artGroup = NULL;
		result = sAIArt->GetFirstArtOfLayer(NULL, &artGroup);
		aisdk::check_ai_error(result);

		// Add the new on path text item to the layer.
		AITextOrientation orient = kHorizontalTextOrientation;
		SnpArtHelper artHelper;
		AIArtHandle	textPath;
		result = artHelper.NewArc(textPath);
		aisdk::check_ai_error(result);
		result = sAITextFrame->NewOnPathText(kPlaceAboveAll, artGroup, orient, textPath, 0, 0, NULL, false, &textPath);
		aisdk::check_ai_error(result);

		// Set the contents of the text range.
		TextRangeRef range = NULL;
		result = sAITextFrame->GetATETextRange(textPath, &range);
		aisdk::check_ai_error(result);
		ITextRange iTextRange(range);
		iTextRange.InsertAfter(ai::UnicodeString("New on path text item").as_ASUnicode().c_str());
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
ASErr SnpText::CreateThreadedInPathText()
{
	ASErr result = kNoErr;
	try {
		// Get the group art that contains all the art in the current layer.
		AIArtHandle artGroup = NULL;
		result = sAIArt->GetFirstArtOfLayer(NULL, &artGroup);
		aisdk::check_ai_error(result);

		// New story to extend over several text frames.
		ai::UnicodeString newStory("This is a story. This story is intended to be long enough to extend over a number of text frames. Each text frame should continue displaying the story.");
		
		// Create first new text frame.
		AIRealRect rect1Position = {100, 700, 200, 650};
		AIArtHandle textPath1;
		SnpArtHelper artHelper;
		result = artHelper.NewRectangle(rect1Position, false, textPath1);
		aisdk::check_ai_error(result);
		// Create first text frame.
		result = sAITextFrame->NewInPathText(kPlaceAboveAll, 
										 artGroup, // Prep.
										 kHorizontalTextOrientation, 
										 textPath1, // Art path item.
										 NULL, // BaseFrame.
										 false, // Append.
										 &textPath1);
		aisdk::check_ai_error(result);

		// Insert the story into the first text frame.
		TextRangeRef textRange1 = NULL;
		result = sAITextFrame->GetATETextRange(textPath1, &textRange1);
		aisdk::check_ai_error(result);
		ITextRange iTextRange1(textRange1);
		iTextRange1.Remove();
		iTextRange1.InsertAfter(newStory.as_ASUnicode().c_str());
		
		// Create second text frame and link it to the first text frame.
		AIRealRect rect2Position = {200, 650, 300, 600};
		AIArtHandle textPath2;
		result = artHelper.NewRectangle(rect2Position, false, textPath2);
		aisdk::check_ai_error(result);
		result = sAITextFrame->NewInPathText(kPlaceBelow, 
										 textPath1, // Prep.
										 kHorizontalTextOrientation, 
										 textPath2, // Art path item.
										 textPath1, // BaseFrame.
										 true, // Append.
										 &textPath2);
		aisdk::check_ai_error(result);
		
		// Create third text frame and link it to the second text frame.
		AIRealRect rect3Position = {300, 600, 400, 550};
		AIArtHandle textPath3;
		result = artHelper.NewRectangle(rect3Position, false, textPath3);
		aisdk::check_ai_error(result);
		result = sAITextFrame->NewInPathText(kPlaceBelow, 
										 textPath2, // Prep.
										 kHorizontalTextOrientation, 
										 textPath3, // Art path item.
										 textPath2, // BaseFrame.
										 true, // Append.
										 &textPath3);
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
ASErr SnpText::InsertText()
{
	ASErr result = kNoErr;
	try {
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);
		if (ranges.GetSize() > 0) {
			ITextRange range = ranges.Item(0);
			range.InsertBefore(ai::UnicodeString("Text inserted.").as_ASUnicode().c_str());
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
ASErr SnpText::CopyText()
{
	ASErr result = kNoErr;
	try	{		
		// Get the selected text.
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);
		ITextRanges ranges(rangesRef);	
			
		if (ranges.GetSize() > 0) {
			ITextRange range = ranges.Item(0);
			if (range.GetSize() > 0)	{
				// Get the art handle for the text frame containing the selected text.
				ITextFrame frame = range.GetStory().GetFrame(0);
				TextFrameRef frameRef = frame.GetRef();
				AIArtHandle frameArt = NULL;
				result = sAITextFrame->GetAITextFrame(frameRef, &frameArt);
				aisdk::check_ai_error(result);

				// Create a new text frame item the same height and width 
				// as the text frame containing the selected text.
				AIRealRect frameBounds = {0, 0, 0, 0};
				result = sAIArt->GetArtBounds(frameArt, &frameBounds);
				aisdk::check_ai_error(result);
				AIReal height = frameBounds.top - frameBounds.bottom;
				AIRealRect rect = {0, 0, 0, 0};
				rect.top = frameBounds.top - (height + 10);
				rect.bottom = frameBounds.bottom - (height + 10);
				rect.left = frameBounds.left;
				rect.right = frameBounds.right;				
				AIArtHandle newFrame;
				SnpArtHelper artHelper;
				result = artHelper.NewRectangle(rect, false, newFrame);
				aisdk::check_ai_error(result);
				result = sAITextFrame->NewInPathText(kPlaceAboveAll, NULL, kHorizontalTextOrientation, 
					newFrame, NULL, false, &newFrame);
				aisdk::check_ai_error(result);

				// Copy the selected text to the new text frame.
				TextRangeRef rangeRef = NULL;		
				result = sAITextFrame->GetATETextRange(newFrame, &rangeRef);
				aisdk::check_ai_error(result);
				ITextRange textRange(rangeRef);	
				textRange.Remove();	
				ITextRange copiedRange = range.Clone();
				textRange.InsertAfter(copiedRange);
			}
			else {
				SnippetRunnerLog::Instance()->Write("Select some text.");
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
ASErr SnpText::MoveText()
{
	ASErr result = kNoErr;
	try {
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);
		ITextRanges ranges(rangesRef);	
			
		if (ranges.GetSize() > 0) {
			ITextRange range = ranges.Item(0);
			if (range.GetSize() > 0) {
				ITextFrame frame = range.GetStory().GetFrame(0);
				TextFrameRef frameRef = frame.GetRef();
				AIArtHandle frameArt = NULL;
				result = sAITextFrame->GetAITextFrame(frameRef, &frameArt);
				aisdk::check_ai_error(result);

				AIRealRect frameBounds = {0, 0, 0, 0};
				result = sAIArt->GetArtBounds(frameArt, &frameBounds);
				aisdk::check_ai_error(result);

				AIReal height = frameBounds.top - frameBounds.bottom;
				AIRealRect rect = {0, 0, 0, 0};
				rect.top = frameBounds.top - (height + 10);
				rect.bottom = frameBounds.bottom - (height + 10);
				rect.left = frameBounds.left;
				rect.right = frameBounds.right;

				AIArtHandle newFrame;
				SnpArtHelper artHelper;
				result = artHelper.NewRectangle(rect, false, newFrame);
				aisdk::check_ai_error(result);
				result = sAITextFrame->NewInPathText(kPlaceAboveAll, NULL, kHorizontalTextOrientation, newFrame, NULL, false, &newFrame);
				aisdk::check_ai_error(result);

				TextRangeRef rangeRef = NULL;		
				result = sAITextFrame->GetATETextRange(newFrame, &rangeRef);
				aisdk::check_ai_error(result);

				ITextRange textRange(rangeRef);	
				textRange.Remove();	
				ITextRange copiedRange = range.Clone();
				textRange.InsertAfter(copiedRange);
				range.Remove();
				range.DeSelect();
			}
			else {
				SnippetRunnerLog::Instance()->Write("Select some text.");
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
ASErr SnpText::ReplaceText()
{
	ASErr result = kNoErr;
	try {
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);	
			
		if (ranges.GetSize() > 0) {
			ITextRange range = ranges.Item(0);
			if (range.GetSize() > 0) {
				range.Remove();
				range.InsertAfter(ai::UnicodeString("Text replaced.").as_ASUnicode().c_str());
				range.DeSelect();
			}
			else {
				SnippetRunnerLog::Instance()->Write("Select some text.");
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
ASErr SnpText::SelectTextRange()
{
	ASErr result = kNoErr;
	try {
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);
		if (ranges.GetSize() > 0) {
			ITextRange textRange = ranges.Item(0);
			if (textRange.GetSize() > 0)	{
				ranges.Select(true);
			}
			else {
				SnippetRunnerLog::Instance()->Write("Text range is empty.");
				aisdk::check_ai_error(kBadParameterErr);
			}
		}
		else {
			SnippetRunnerLog::Instance()->Write("Select an item containing text.");
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
ASErr SnpText::DeleteTextRange()
{
	ASErr result = kNoErr;
	try {
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);
		if (ranges.GetSize() > 0) {
			// Delete each selected text range using ITextRange::Remove.
			for (ASInt32 i = 0; i < ranges.GetSize(); i++) {
				ITextRange range = ranges.Item(i);
				range.Remove();
				// DeSelect so text which moves into the position of the deleted text
				// is not selected.
				range.DeSelect();
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
ASErr SnpText::ConvertLegacyText()
{
	ASErr result = kNoErr;
	try {
		result = sAILegacyTextConversion->ConvertAllToNative(NULL);
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
ASErr SnpText::SetTextFocus()
{
	ASErr result = kNoErr;
	try {
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);
		if (ranges.GetSize() > 0) {
			// The following is a workaround to get the caret blinking at the beginning
			// of the selected text range.
			ITextRange textRange = ranges.Item(0);
			AIToolType toolNumber = NULL;
			result = sAITool->GetToolNumberFromName("Adobe Select Tool", &toolNumber);
			aisdk::check_ai_error(result);
			AIToolHandle toolHandle = NULL;
			result = sAITool->GetToolHandleFromNumber(toolNumber, &toolHandle);
			aisdk::check_ai_error(result);
			result = sAITool->SetSelectedTool(toolHandle);
			aisdk::check_ai_error(result);
			IStory story = textRange.GetStory();
			result = sAIDocument->SetTextFocus(story.GetRef());
			aisdk::check_ai_error(result);
			result = sAIDocument->RedrawDocument();
			aisdk::check_ai_error(result);
		}
		else {
			SnippetRunnerLog::Instance()->Write("Select an item containing a story.");
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
ASErr SnpText::LoseTextFocus()
{
	ASErr result = kNoErr;
	try {
		AIBoolean focus = NULL;
		result = sAIDocument->HasTextFocus(&focus);
		aisdk::check_ai_error(result);

		if (focus) {
			result = sAIDocument->LoseTextFocus();
			aisdk::check_ai_error(result);
		}
		else {
			SnippetRunnerLog::Instance()->Write("Document does not have text focus.");
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
ASErr SnpText::LinkTextFrames()
{
	ASErr result = kNoErr;
	try {		
		AIArtSet artSet = NULL;
		result = sAIArtSet->NewArtSet(&artSet);
		aisdk::check_ai_error(result);

		AIArtSpec findArtSpec = {kTextFrameArt, kArtSelected, kArtSelected};
		result = sAIArtSet->MatchingArtSet(&findArtSpec, 1, artSet);
		aisdk::check_ai_error(result);

		size_t count = 0;
		result = sAIArtSet->CountArtSet(artSet, &count);
		aisdk::check_ai_error(result);

		if (count > 1) {
			for (size_t i = 0; i < count; i++) {
				AIArtHandle baseFrame = NULL;
				result = sAIArtSet->IndexArtSet(artSet, i, &baseFrame);
				aisdk::check_ai_error(result);

				AIArtHandle nextFrame = NULL;
				result = sAIArtSet->NextInArtSet(artSet, baseFrame, &nextFrame);
				aisdk::check_ai_error(result);
				
				if (nextFrame != NULL) {
					AIBool8 linked = false;
					result = sAITextFrame->PartOfLinkedText(nextFrame, &linked);
					aisdk::check_ai_error(result);
							
					if (!linked) {
						result = sAITextFrame->Link(baseFrame, nextFrame);
						aisdk::check_ai_error(result);
					}
					else {
						SnippetRunnerLog::Instance()->Write("Will not link if a text frame is already linked.");
						aisdk::check_ai_error(kBadParameterErr);
					}
				}
			}
		}
		else {
			SnippetRunnerLog::Instance()->Write("Select at least 2 unlinked text frames.");
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
ASErr SnpText::UnlinkTextFrames()
{
	ASErr result = kNoErr;
	try {
		AIArtSet artSet = NULL;
		result = sAIArtSet->NewArtSet(&artSet);
		aisdk::check_ai_error(result);

		AIArtSpec findArtSpec = {kTextFrameArt, kArtSelected, kArtSelected};
		result = sAIArtSet->MatchingArtSet(&findArtSpec, 1, artSet);
		aisdk::check_ai_error(result);

		size_t count = 0;
		result = sAIArtSet->CountArtSet(artSet, &count);
		aisdk::check_ai_error(result);

		if (count > 0) {
			for (size_t i = 0; i < count; i++) {
				AIArtHandle selectedFrame = NULL;
				result = sAIArtSet->IndexArtSet(artSet, i, &selectedFrame);
				aisdk::check_ai_error(result);

				AIBool8 linked = false;
				result = sAITextFrame->PartOfLinkedText(selectedFrame, &linked);
				aisdk::check_ai_error(result);

				if (linked) {
					result = sAITextFrame->Unlink(selectedFrame, true, true);
					aisdk::check_ai_error(result);
				}
			}
		}
		else {
			SnippetRunnerLog::Instance()->Write("Select at least 1 linked text frame.");
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
ASErr SnpText::RunSpellCheck()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);	
		if (ranges.GetSize() > 0) {
			for (ASInt32 i = 0; i < ranges.GetSize(); i++) {
				// Get the path to the Illustrator dictionary.
				ITextRange range = ranges.Item(i);
				ai::FilePath fileSpec;
				result = sAIFolders->FindFolder(kAIDictionariesFolderType, false, fileSpec);
				aisdk::check_ai_error(result);

				// Create a new ISpell object.
				SpellRef spellRef = NULL;
				result = sAIATETextUtil->GetSpellFile(fileSpec, &spellRef);
				aisdk::check_ai_error(result);
				ISpell spellChecker = ISpell(spellRef);
				spellChecker.SetSearchRange(range);

				// Search the text range for misspelled words.
				IStory story = range.GetStory();
				ITextRange misspelledWord = story.GetTextRange(0, 0);
				SpellCheckingResult word = kNoProblems;
				Language language = kEnglishLanguage;
				if (spellChecker.IsInitialized()) {
					bool wordFound = true;
					while ((wordFound = spellChecker.FindOneMisspelledWord(&word, misspelledWord, &language))) {
						if (wordFound) {
							ai::AutoBuffer<ASUnicode> contents(misspelledWord.GetSize());
							misspelledWord.GetContents(contents.GetBuffer(), misspelledWord.GetSize());
							
							log->Write("Misspelled word:");
							log->BeginIndent();
							SnippetRunnerLog::Instance()->Write(ai::UnicodeString(contents.GetBuffer(), contents.GetCount()));
							log->EndIndent();
							log->Write("Spelling suggestions:");
							ASInt32 wordListSize = spellChecker.GetWordListSize();
							ASUnicode* wordBuffer = new ASUnicode[wordListSize];
							ASInt32 numberWords = 0;
							spellChecker.GetWordListContents(wordBuffer, wordListSize, &numberWords);

							log->BeginIndent();
							size_t offset = 0;
							for (ASInt32 word = 0; word < numberWords; ++word) {
								ai::UnicodeString contentsAsUnicode = ai::UnicodeString(&wordBuffer[offset]);
								log->Write(contentsAsUnicode);
								offset += contentsAsUnicode.length() + 1;
							}	
							log->EndIndent();
							delete []wordBuffer;
						}
						misspelledWord.SetStart(misspelledWord.GetEnd());
					}
				}
			}
		}
		else {
			log->Write("Select text range to spell check.");
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
ASErr SnpText::FindAndReplace()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);	
		if (ranges.GetSize() > 0) {
			SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();

			// Get current document text resources object.
			DocumentTextResourcesRef docTextResourcesRef = NULL;
			result = sAIDocument->GetDocumentTextResources(&docTextResourcesRef);
			aisdk::check_ai_error(result);
			IDocumentTextResources resources(docTextResourcesRef);

			// Get search and replace strings.
			ai::UnicodeString search = parameter->GetString(ai::UnicodeString("Find:"), ai::UnicodeString("find"));
			ai::UnicodeString replace = parameter->GetString(ai::UnicodeString("Replace with:"), ai::UnicodeString("replace"));

			// Convert search string to ASUnicode.
			const std::basic_string<ASUnicode>& searchStr = search.as_ASUnicode();
			const size_t searchSize = searchStr.length() + 1;
			ASUnicode* searchString = new ASUnicode[searchSize];
			memcpy(searchString, searchStr.data(), searchSize * sizeof(ASUnicode));

			// Convert replace string to ASUnicode.
			const std::basic_string<ASUnicode>& replaceStr = replace.as_ASUnicode();
			const size_t replaceSize = replaceStr.length() + 1;
			ASUnicode* replaceString = new ASUnicode[replaceSize];
			memcpy(replaceString, replaceStr.data(), replaceSize * sizeof(ASUnicode));

			// Create IFind object.
			IFind stringFinder = resources.GetFind();
			stringFinder.SetSearchChars(searchString);
			stringFinder.SetReplaceChars(replaceString);

			ASInt32 count = 0;

			for (ASInt32 i = 0; i < ranges.GetSize(); i++) {
				ITextRange range = ranges.Item(i);
				ITextRange searchRange = range;
				ASInt32 currentPoint = 0, startPoint = 0;
				stringFinder.SetSearchRange(range);
				stringFinder.GetPreReplaceAllSettings(&currentPoint, &startPoint);

				
				bool found = true;
				while (found)
				{
					found = stringFinder.FindMatch(searchRange);
					if (found)
					{
						found = stringFinder.ReplaceMatch(searchRange, false);
						
						// We now want to change the range we're searching in to avoid falling into an infinite loop
						// e.g. find: is, replace: isn't.
						// The next search should start from the end of the word we just changed
						ASInt32 nextRangeStart = searchRange.GetStart() + (replaceSize - 1);
						ASInt32 nextRangeEnd = range.GetEnd();
						
						searchRange.SetStart((nextRangeStart < nextRangeEnd) ? nextRangeStart : nextRangeEnd);
						searchRange.SetEnd(nextRangeEnd);
						count++;
					}
				}
				
				stringFinder.RestorePreReplaceAllSettings(currentPoint, startPoint);
			}
			log->WritePrintf("%d %s", count, "occurances replaced.");
			delete []searchString;
			delete []replaceString;
		}
		else {
			log->Write("Select text range to search.");
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
ASErr SnpText::RunHitTest()
{
	ASErr result = kNoErr;
	try {
		// Get the point to check for a hit.
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		AIReal xPos = parameter->GetReal(ai::UnicodeString("Enter x position:"), 800.0);
		AIReal yPos = parameter->GetReal(ai::UnicodeString("Enter y position:"), 750.0);

		AIRealPoint point;
		point.h = xPos;
		point.v = yPos;

		// Perform the hit test.
		AIHitRef hitRef = NULL;
		result = sAIHitTest->HitTest(NULL, &point, kAllHitRequest, &hitRef);
		aisdk::check_ai_error(result);

		// Check if an item was hit.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		if (sAIHitTest->IsHit(hitRef))
		{			
			log->Write("Item hit.");
			SnippetRunnerLog::Indent indent;
			AIArtHandle artHit = sAIHitTest->GetArt(hitRef);
			short artType = kUnknownArt;
			result = sAIArt->GetArtType(artHit, &artType);
			aisdk::check_ai_error(result);

			if (artType == kTextFrameArt)
			{
				log->Write("kTextFrameArt");
				ai::int32 hitPart = sAITextFrameHit->GetPart(hitRef);
				switch (hitPart) {
					case kAITextText:
						log->Write(ai::UnicodeString("Direct hit on text."));
						break;
					case kAITextInport:
						log->Write(ai::UnicodeString("Hit on 'in' port of either in-path or on-path text."));
						break;
					case kAITextOutPort:
						log->Write(ai::UnicodeString("Hit on 'out' port of either in-path or on-path text."));
						break;
					case kAITextStart:
						log->Write(ai::UnicodeString("Hit on start point of on-path text."));
						break;
					case kAITextMiddle:
						log->Write(ai::UnicodeString("Hit on middle point of on-path text."));
						break;
					case kAITextEnd:
						log->Write(ai::UnicodeString("Hit on end point of on-path text."));
						break;
					case kAITextNowhere:
						log->Write(ai::UnicodeString("Hit on empty area of text frame."));
						break;
				}		   
			}
		}
		else
			log->Write("No hit.");

		// Decrement reference count and free memory.
		sAIHitTest->Release(hitRef);
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
ASErr SnpText::IterateAllFonts()
{
	ASErr result = kNoErr;
	try {
		// Count the number of fonts available.
		ai::int32 fontCount = 0;
		result = sAIFont->CountFonts(&fontCount);
		aisdk::check_ai_error(result);

		for (ai::int32 i = 0; i < fontCount; i++) {
			// Report each font name to the log.
			AIFontKey fontKey = NULL;
			result = sAIFont->IndexFontList(i, &fontKey);
			aisdk::check_ai_error(result);
			ai::UnicodeString fontName;
			result = sAIFont->GetUserFontNameUnicode(fontKey, fontName, true);
			aisdk::check_ai_error(result);
			SnippetRunnerLog::Instance()->Write(ai::UnicodeString(fontName));
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpText::IterateUsedFonts()
{
	ASErr result = kNoErr;
	try {
		// Find the first text frame in the document.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		AIArtSpec specs[1] = {{kTextFrameArt,0,0}};
		SnpArtSetHelper textFrameArtSet(specs, 1);
		if (textFrameArtSet.GetCount() > 0) {
			AIArtHandle textFrameArt = textFrameArtSet [0];
			TextFrameRef textFrameRef = NULL;
			result = sAITextFrame->GetATETextFrame(textFrameArt, &textFrameRef);
			aisdk::check_ai_error(result);
			ITextFrame textFrame(textFrameRef);

			// Use the text frame to get all the text in the document.
			IStory story = textFrame.GetStory();
			IStories stories = story.GetStories();
			ITextRanges textRanges = stories.GetTextRanges();

			// Get the ICharInspector for the documents' text ranges.
			ICharInspector charInspector = textRanges.GetCharInspector();

			// Get the array of font references found in the text ranges.
			IArrayFontRef fontRefArray = charInspector.GetFont();

			// Iterate through the font references and report each 
			// associated font name to the log.
			ASInt32 count = 0;
			ASInt32 fontCount = fontRefArray.GetSize();
			if (fontCount > 0) {
				while (count < fontCount) {
					IFont currentFont = fontRefArray.Item(count);
					if (!currentFont.IsNull()) {
						FontRef fontRef = currentFont.GetRef();
						AIFontKey fontKey = NULL;
						result = sAIFont->FontKeyFromFont(fontRef, &fontKey);
						aisdk::check_ai_error(result);
						ai::UnicodeString fontName;
						result = sAIFont->GetUserFontNameUnicode(fontKey, fontName, true);
						aisdk::check_ai_error(result);
						log->Write(ai::UnicodeString(fontName));
					}
					count++;
				}
			}
			else {
				log->Write("No text in the current document.");
				aisdk::check_ai_error(kBadParameterErr);
			}
		}
		else {
			log->Write("No text frames in the current document");
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


// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpText available to the SnippetRunner framework.
 */
class _SnpRunnableText : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableText () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableText () {}

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

		/** Returns true if the snippet can run.
			@param runnableContext see ISnpRunnableContext for documentation.
			@return true if snippet can run, false otherwise
		 */
		ASBoolean			CanRun(SnpRunnable::Context& runnableContext);

		/** Runs the snippet.
			@param runnableContext see ISnpRunnableContext for documentation.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr		Run(SnpRunnable::Context& runnableContext);

		/** Sets the context manually for LinkTextFrames and LoseTextFocus.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr SetupContext(SnpRunnable::Context& runnableContext);
};

/*
*/
std::string _SnpRunnableText::GetName() const
{
	return "Text";
}

/*
*/
std::string _SnpRunnableText::GetDescription() const 
{
	return "Demonstrates manipulation of text items, spell checking, find and replace and hit tests.";
}

/*
*/
SnpRunnable::Operations _SnpRunnableText::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("CreatePointText", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("CreateInPathText", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("CreateOnPathText", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("CreateThreadedInPathText", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("InsertText", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("CopyText", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("MoveText", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("ReplaceText", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("SelectTextRange", "selected text art", kSnpRunCharactersContext));
	operations.push_back(Operation("DeleteTextRange", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("SetTextFocus", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("LoseTextFocus", "document in text focus", kSnpRunCustomContext));
	operations.push_back(Operation("LinkTextFrames", "selected text frames", kSnpRunCustomContext));
	operations.push_back(Operation("UnlinkTextFrames", "selected text frames", kSnpRunCharactersContext));
	operations.push_back(Operation("RunSpellCheck", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("FindAndReplace", "selected text", kSnpRunCharactersContext));	
	operations.push_back(Operation("RunHitTest", "document containing text art", kSnpRunUnitTestDocumentContext));	
	operations.push_back(Operation("IterateAllFonts", "document", kSnpRunNewDocumentContext));	
	operations.push_back(Operation("IterateUsedFonts", "document containing text art", kSnpRunUnitTestDocumentContext));

	return operations;
}

/*
*/
std::string _SnpRunnableText::GetDefaultOperationName() const 
{
	return "CreatePointText";
	// Note: this must be one of the operation names returned by GetOperations.
}

/*
*/
std::vector<std::string> _SnpRunnableText::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Text Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableText::CanRun(SnpRunnable::Context& runnableContext)
{
	AIErr result = kNoErr;
	try {
		SnpSelectionHelper selectionHelper;
		if (!selectionHelper.IsDocumentSelected())
			return false;

		if ("CopyText" == runnableContext.GetOperation().GetName() ||
			"MoveText" == runnableContext.GetOperation().GetName() ||
			"ReplaceText" == runnableContext.GetOperation().GetName() ||
			"SelectTextRange" == runnableContext.GetOperation().GetName() ||
			"DeleteTextRange" == runnableContext.GetOperation().GetName() ||
			"SetTextFocus" == runnableContext.GetOperation().GetName() ||
			"RunSpellCheck" == runnableContext.GetOperation().GetName() ||
			"FindAndReplace" == runnableContext.GetOperation().GetName()) {
			
				return selectionHelper.IsTextRangeSelected();
		}
		else if ("InsertText" == runnableContext.GetOperation().GetName()) {
			AIBoolean focus = false;
			result = sAIDocument->HasTextFocus(&focus);
			if (selectionHelper.IsTextRangeSelected() || focus)
				return true;
			else
				return false;
		}
		else if ("LoseTextFocus" == runnableContext.GetOperation().GetName()) {
			AIBoolean textFocus = false;
			result = sAIDocument->HasTextFocus(&textFocus);
			aisdk::check_ai_error(result);
			return textFocus;
		}
		else if ("LinkTextFrames" == runnableContext.GetOperation().GetName() ||
			"UnlinkTextFrames" == runnableContext.GetOperation().GetName()) {
			SnpArtSetHelper artSetHelper;
			artSetHelper.FindSelectedArt();

			SnpArtSetHelper textArtSet = artSetHelper.Filter(kTextFrameArt);
			return textArtSet.GetCount() > 1;
		}
		else if ("IterateUsedFonts" == runnableContext.GetOperation().GetName()) {
			AIArtSpec specs[1] = {{kTextFrameArt,0,0}};
			SnpArtSetHelper textFrameArtSet(specs, 1);
			return textFrameArtSet.GetCount() > 0;
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
ASErr _SnpRunnableText::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpText instance;
	if ("CreatePointText" == runnableContext.GetOperation().GetName()) {
		result = instance.CreatePointText();
	}
	else if ("CreateInPathText" == runnableContext.GetOperation().GetName()) {
		result = instance.CreateInPathText();
	}
	else if ("CreateOnPathText" == runnableContext.GetOperation().GetName()) {
		result = instance.CreateOnPathText();
	}
	else if ("CreateThreadedInPathText" == runnableContext.GetOperation().GetName()) {
		result = instance.CreateThreadedInPathText();
	}
	else if ("InsertText" == runnableContext.GetOperation().GetName()) {
		result = instance.InsertText();
	}
	else if ("CopyText" == runnableContext.GetOperation().GetName()) {
		result = instance.CopyText();
	}
	else if ("MoveText" == runnableContext.GetOperation().GetName()) {
		result = instance.MoveText();
	}
	else if ("ReplaceText" == runnableContext.GetOperation().GetName()) {
		result = instance.ReplaceText();
	}
	else if ("SelectTextRange" == runnableContext.GetOperation().GetName()) {
		result = instance.SelectTextRange();
	}
	else if ("DeleteTextRange" == runnableContext.GetOperation().GetName()) {
		result = instance.DeleteTextRange();
	}
	else if ("ConvertLegacyText" == runnableContext.GetOperation().GetName()) {
		result = instance.ConvertLegacyText();
	}
	else if ("SetTextFocus" == runnableContext.GetOperation().GetName()) {
		result = instance.SetTextFocus();
	}
	else if ("LoseTextFocus" == runnableContext.GetOperation().GetName()) {
		result = instance.LoseTextFocus();
	}
	else if ("LinkTextFrames" == runnableContext.GetOperation().GetName()) {
		result = instance.LinkTextFrames();
	}
	else if ("UnlinkTextFrames" == runnableContext.GetOperation().GetName()) {
		result = instance.UnlinkTextFrames();
	}
	else if ("RunSpellCheck" == runnableContext.GetOperation().GetName()) {
		result = instance.RunSpellCheck();
	}
	else if ("FindAndReplace" == runnableContext.GetOperation().GetName()) {
		result = instance.FindAndReplace();
	}
	else if ("RunHitTest" == runnableContext.GetOperation().GetName()) {
		result = instance.RunHitTest();
	}
	else if ("IterateAllFonts" == runnableContext.GetOperation().GetName()) {
		result = instance.IterateAllFonts();
	}
	else if ("IterateUsedFonts" == runnableContext.GetOperation().GetName()) {
		result = instance.IterateUsedFonts();
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

ASErr _SnpRunnableText::SetupContext(SnpRunnable::Context &runnableContext)
{
	ASErr result = kNoErr;
	try {
		switch (runnableContext.GetOperation().GetUnitTestContext()) {
			case SnpRunnable::kSnpRunCustomContext:	{
					SnpText instance;
					SnpDocumentHelper docHelper;
					AIDocumentHandle document = NULL;
					SDK_ASSERT(!docHelper.OpenDocument(ai::UnicodeString("unittest-1.ai"), document));
					SDK_ASSERT(document);

					// Select the text frames.
					SnpArtSetHelper artSet;
					artSet.FindLayerArt(ai::UnicodeString("layer.text"));
					artSet = artSet.Filter(kTextFrameArt);
					artSet.SelectAll();
					result = sAIDocument->RedrawDocument();
					aisdk::check_ai_error(result);

					if ("LinkTextFrames" == runnableContext.GetOperation().GetName()) {
						// Unlink the text frames before continuing.
						instance.UnlinkTextFrames();
					}
					else if ("LoseTextFocus" == runnableContext.GetOperation().GetName()){
						// Set the text focus before continuing.
						instance.SetTextFocus();
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
static _SnpRunnableText instance_SnpRunnableText;

// End SnpText.cpp

