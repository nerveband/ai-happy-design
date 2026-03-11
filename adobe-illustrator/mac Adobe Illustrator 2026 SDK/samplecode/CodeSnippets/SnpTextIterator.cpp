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
#include "SnpArtSetHelper.h"

using namespace ATE;

/** Provides functions that demonstrate how to iterate text characters, 
	words, paragraphs, stories, ranges,	frames, lines, glyph runs, text 
	runs and kern types.

	@see AIHardSoftSuite, AIDocumentSuite, ITextRanges, ITextRange, IStory,
	ITextFrame, ITextLine, IGlyphRun, ITextFramesIterator, ITextLinesIterator,
	IGlyphRunsIterator, IWordsIterator, IParagraphsIterator
	@ingroup sdk_snippet
 */
class SnpTextIterator
{
	public:

		/** Constructor.
		*/
		SnpTextIterator() {}

		/** Destructor.
		*/
		virtual ~SnpTextIterator() {}
		
		/** Reports all selected characters to the SnippetRunner log.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr IterateSelectedCharacters();

		/** Reports all selected words to the SnippetRunner log.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr IterateSelectedWords();

		/** Reports all selected paragraphs to the SnippetRunner log.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr IterateSelectedParagraphs();

		/** Reports the content of each story in the current document to the 
			SnippetRunner log.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr IterateStories();

		/** Reports all selected stories to the SnippetRunner log.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr IterateSelectedStories();

		/** Reports all selected text ranges to the SnippetRunner log.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr IterateSelectedTextRanges();

		/** Reports the content of each text frame in the current document
			to the SnippetRunner log.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr IterateTextFrames();

		/** Reports all selected text frames to the SnippetRunner log.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr IterateSelectedTextFrames();

		/** Reports all selected lines to the SnippetRunner log.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr IterateLinesInSelectedFrames();

		/** Reports all selected glyph runs to the SnippetRunner log.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr IterateSelectedGlyphRuns();

		/** Reports all selected text runs to the SnippetRunner log.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr IterateSelectedTextRuns();

		/** Reports all selected character kern types to the SnippetRunner log.
			@return kNoErr if successful, otherwise return error code
		*/
		ASErr IterateSelectedCharacterKernTypes();

};

/*
*/
ASErr SnpTextIterator::IterateSelectedCharacters()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);
		if (ranges.GetSize() > 0) {
			ASInt32 rangeCount = 0;
			for (ASInt32 i = 0; i < ranges.GetSize(); i++) {
				log->WritePrintf("Text range %d: ", rangeCount);
				rangeCount++;
				ITextRange range = ranges.Item(i);
				if (range.GetSize() > 0) {
					SnippetRunnerLog::Indent indent;
					ASInt32 rangeSize = range.GetSize();
					ai::AutoBuffer<ASUnicode> contents(rangeSize);
					range.GetContents(contents.GetBuffer(), rangeSize);
					if (contents.IsValid()) {
						ai::UnicodeString contentsAsUnicode(contents.GetBuffer(), contents.GetCount());
						for (ASInt32 charIndex = 0; charIndex < rangeSize; charIndex++) {
							log->WritePrintf("Character %d:", charIndex);
							SnippetRunnerLog::Indent indent2;
							log->Write(ai::UnicodeString(1, contentsAsUnicode.at(charIndex)));
						}
					}
				}
				else {
					log->Write("Select some characters.");
					aisdk::check_ai_error(kBadParameterErr);
				}
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
ASErr SnpTextIterator::IterateSelectedWords()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);	
		if (ranges.GetSize() > 0) {
			SnippetRunnerLog::Indent indent;
			ASInt32 wordCount = 0;
			IWordsIterator iWords(ranges);
			for (iWords.MoveToFirst(); iWords.IsNotDone(); iWords.Next()) {
				ITextRange word = iWords.Item();
				ASInt32 rangeSize = word.GetSize();
				ai::AutoBuffer<ASUnicode> contents(rangeSize);
				word.GetContents(contents.GetBuffer(), rangeSize);
				if (contents.IsValid()) {
					log->WritePrintf("Word %d:", wordCount);
					SnippetRunnerLog::Indent indent2;
					log->Write(ai::UnicodeString(contents.GetBuffer(), contents.GetCount()));
				}
				wordCount++;
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
ASErr SnpTextIterator::IterateSelectedParagraphs()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);
		if (ranges.GetSize() > 0) {
			IParagraphsIterator paragraphsIter = ranges.GetParagraphsIterator();
			if (!paragraphsIter.IsEmpty()) {
				SnippetRunnerLog::Indent indent;
				ASInt32 paraCount = 0;
				while (paragraphsIter.IsNotDone()) {
					IParagraph paragraph = paragraphsIter.Item();
					ITextRange paraTextRange = paragraph.GetTextRange();
					ai::AutoBuffer<ASUnicode> contents(paraTextRange.GetSize());
					paraTextRange.GetContents(contents.GetBuffer(), paraTextRange.GetSize());
					log->WritePrintf("Paragraph %d:", paraCount);
					SnippetRunnerLog::Indent indent2;
					log->Write(ai::UnicodeString(contents.GetBuffer(), contents.GetCount()));
					paraCount++;
					paragraphsIter.Next();
				}
			}
			else {
				log->Write("No paragraphs found.");
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
ASErr SnpTextIterator::IterateStories()
{
	ASErr result = kNoErr;
	try {
		// Find text frame art.
		SDK_ASSERT(sAITextFrame);
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		AIArtSpec specs[1] = {{kTextFrameArt,0,0}};
		SnpArtSetHelper textFrameArtSet(specs, 1);
		if (textFrameArtSet.GetCount() > 0) {
			AIArtHandle textFrameArt = textFrameArtSet [0];
			// Create ITextFrame object.
			TextFrameRef textFrameRef = NULL;
			result = sAITextFrame->GetATETextFrame(textFrameArt, &textFrameRef);
			aisdk::check_ai_error(result);
			ITextFrame textFrame(textFrameRef);
			// Get all the stories in the current document.
			ITextRange textRange = textFrame.GetTextRange();
			IStory story = textRange.GetStory();
			IStories stories = story.GetStories();
			if (stories.GetSize() > 0) {
				for (ASInt32 index = 0; index < stories.GetSize(); index++) {
					IStory currentStory = stories.Item(index);
					ITextRange storyRange = currentStory.GetTextRange();
					ai::AutoBuffer<ASUnicode> contents(storyRange.GetSize());
					storyRange.GetContents(contents.GetBuffer(), storyRange.GetSize());
					log->WritePrintf("Story %d:", index);
					SnippetRunnerLog::Indent indent;
					log->Write(ai::UnicodeString(contents.GetBuffer(), contents.GetCount()));
				}
			}
		}
		else {
			log->Write("Create some text art.");
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
ASErr SnpTextIterator::IterateSelectedStories()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);
		if (ranges.GetSize() > 0) {
			ASInt32 storyCount = 0;
			for (ASInt32 rangeIndex = 0; rangeIndex < ranges.GetSize(); rangeIndex++) {
				SnippetRunnerLog::Indent indent;
				ITextRange textRange = ranges.Item(rangeIndex);
				if (textRange.GetSize() > 0) {
					IStory story = textRange.GetStory();
					ITextRange storyRange = story.GetTextRange();
					ai::AutoBuffer<ASUnicode> contents(storyRange.GetSize());
					storyRange.GetContents(contents.GetBuffer(), storyRange.GetSize());
					log->WritePrintf("Story %d:", storyCount);
					SnippetRunnerLog::Indent indent2;
					log->Write(ai::UnicodeString(contents.GetBuffer(), contents.GetCount()));
					storyCount++;
				}
				else {
					log->Write("Select some text.");
					aisdk::check_ai_error(kBadParameterErr);
				}
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
ASErr SnpTextIterator::IterateSelectedTextRanges()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);
		if (ranges.GetSize() > 0) {
			ASInt32 rangeCount = 0;
			for (ASInt32 rangeIndex = 0; rangeIndex < ranges.GetSize(); rangeIndex++) {
				SnippetRunnerLog::Indent indent;
				ITextRange textRange = ranges.Item(rangeIndex);
				if (textRange.GetSize() > 0) {
					ai::AutoBuffer<ASUnicode> contents(textRange.GetSize());
					textRange.GetContents(contents.GetBuffer(), textRange.GetSize());
					log->WritePrintf("Text range %d:", rangeCount);
					SnippetRunnerLog::Indent indent2;
					log->Write(ai::UnicodeString(contents.GetBuffer(), contents.GetCount()));
					rangeCount++;
				}
				else {
					log->Write("Select some text.");
					aisdk::check_ai_error(kBadParameterErr);
				}
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
ASErr SnpTextIterator::IterateTextFrames()
{
	ASErr result = kNoErr;
	try
	{
		SDK_ASSERT(sAITextFrame);
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		AIArtSpec specs[1] = {{kTextFrameArt,0,0}};
		SnpArtSetHelper textFrameArtSet(specs, 1);
		if (textFrameArtSet.GetCount() > 0)
		{
			for (size_t i = 0; i < textFrameArtSet.GetCount(); i++)
			{
				AIArtHandle textFrameArt = textFrameArtSet [i];
				TextFrameRef textFrameRef = NULL;
				result = sAITextFrame->GetATETextFrame(textFrameArt, &textFrameRef);
				aisdk::check_ai_error(result);
				ITextFrame textFrame(textFrameRef);
				log->WritePrintf("Text frame %d", i);
				SnippetRunnerLog::Indent indent;
				ITextRange textRange = textFrame.GetTextRange();
				ai::AutoBuffer<ASUnicode> contents(textRange.GetSize());
				textRange.GetContents(contents.GetBuffer(), textRange.GetSize());
				log->Write(ai::UnicodeString(contents.GetBuffer(), contents.GetCount()));
			}
		}
		else
		{
			log->Write("Create some text art.");
			aisdk::check_ai_error(kBadParameterErr);
		}
	}
	catch (ai::Error& ex)
	{
		result = ex;
	}
	catch (ATE::Exception& ex)
	{
		result = ex.error;
	}
	return result;
}

/*
*/
ASErr SnpTextIterator::IterateSelectedTextFrames()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);
		if (ranges.GetSize() > 0) {
			SnippetRunnerLog::Indent indent;
			for (ASInt32 rangeIndex = 0; rangeIndex < ranges.GetSize(); rangeIndex++) {
				log->WritePrintf("Text range %d: ", rangeIndex);
				ITextRange range = ranges.Item(rangeIndex);
				ITextFramesIterator framesIter = range.GetTextFramesIterator();
				if (!framesIter.IsEmpty()) {
					SnippetRunnerLog::Indent indent2;
					ASInt32 frameCount = 0;
					while (framesIter.IsNotDone()) {
						ITextFrame textFrame = framesIter.Item();
						ITextRange rangeInFrame = textFrame.GetTextRange();
						ai::AutoBuffer<ASUnicode> contents(range.GetSize());
						rangeInFrame.GetContents(contents.GetBuffer(), rangeInFrame.GetSize());
						log->WritePrintf("Text frame %d:", frameCount);
						SnippetRunnerLog::Indent indent3;
						log->Write(ai::UnicodeString(contents.GetBuffer(), contents.GetCount()));
						frameCount++;
						framesIter.Next();
					}
				}
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
ASErr SnpTextIterator::IterateLinesInSelectedFrames()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);
		if (ranges.GetSize() > 0) {
			SnippetRunnerLog::Indent indent;
			for (ASInt32 rangeIndex = 0; rangeIndex < ranges.GetSize(); rangeIndex++) {
				ITextRange range = ranges.Item(rangeIndex);
				ITextFramesIterator framesIter = range.GetTextFramesIterator();
				log->WritePrintf("Text range %d:", rangeIndex);
				if (!framesIter.IsEmpty()) {
					SnippetRunnerLog::Indent indent2;
					ASInt32 frameCount = 0;
					while (framesIter.IsNotDone()) {
						log->WritePrintf("Text frame %d:", frameCount);
						frameCount++;
						ITextFrame textFrame = framesIter.Item();
						ITextLinesIterator linesIter = textFrame.GetTextLinesIterator();
						if (!linesIter.IsEmpty()) {
							SnippetRunnerLog::Indent indent3;
							ASInt32 lineCount = 0;
							while (linesIter.IsNotDone()) {
								ITextLine textLine = linesIter.Item();
								ITextRange textLineRange = textLine.GetTextRange();
								ai::AutoBuffer<ASUnicode> contents(range.GetSize());
								textLineRange.GetContents(contents.GetBuffer(), textLineRange.GetSize());
								log->WritePrintf("Text line %d:", lineCount);
								SnippetRunnerLog::Indent indent4;
								log->Write(ai::UnicodeString(contents.GetBuffer(), contents.GetCount()));
								lineCount++;
								linesIter.Next();
							}
						}
						else {
							log->Write("No text lines found.");
							aisdk::check_ai_error(kBadParameterErr);
						}
						framesIter.Next();
					}
				}
				else {
					log->Write("No text frames found.");
					aisdk::check_ai_error(kBadParameterErr);
				}
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
ASErr SnpTextIterator::IterateSelectedGlyphRuns()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);	
		if (ranges.GetSize() > 0) {
			for (ASInt32 nRange = 0; nRange < ranges.GetSize(); nRange++) {
				SnippetRunnerLog::Indent indent;
				// For each range get the frame list.
				ITextRange	range = ranges.Item(nRange);
				ITextFramesIterator frames = range.GetTextFramesIterator();
				log->WritePrintf("%s%d","Text Range ", nRange);
				while (frames.IsNotDone()) {
					// For each frame get the textlines.
					ITextFrame frame = frames.Item();
					ITextLinesIterator lines = frame.GetTextLinesIterator();
					while (lines.IsNotDone()) {
						ITextLine line = lines.Item();
						IGlyphRunsIterator glyphRuns = line.GetGlyphRunsIterator();
						ASInt32 glyphCount = 0;
						while (glyphRuns.IsNotDone()) {
							SnippetRunnerLog::Indent indent;
							IGlyphRun glyphRun = glyphRuns.Item();
							if (glyphRun.GetCharacterCount() > 0) {
								// Write the glyph run contents to the log.
								ai::AutoBuffer<ASUnicode> contents(glyphRun.GetCharacterCount());
								glyphRun.GetContents(contents.GetBuffer(), glyphRun.GetCharacterCount());
								log->WritePrintf("Glyph %d:", glyphCount);
								SnippetRunnerLog::Indent indent2;
								log->Write(ai::UnicodeString(contents.GetBuffer(), contents.GetCount()));
							
								// Get the position of each glyph run.
								AIRealMatrix GlyphMatrix = glyphRun.GetMatrix();
								// Concat by [1 0 0 -1 0 0] as text coordinates are going positive the other direction.
								AIRealMatrix flipY={1,0,0,-1,0,0};
								sAIRealMath->AIRealMatrixConcat(&flipY, &GlyphMatrix, &GlyphMatrix );
								// ATE space is application independent, and doesn't know about this Illustrator soft/hard coordinate thingy
								// take care of it here.
								result = sAIHardSoft->AIRealMatrixRealSoft(&GlyphMatrix);
								aisdk::check_ai_error(result);
							}
							glyphCount++;
							glyphRuns.Next();
						}
						lines.Next();
					}
					frames.Next();
				}
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
ASErr SnpTextIterator::IterateSelectedTextRuns()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);	
		if (ranges.GetSize() > 0) {
			ASInt32 count = 0;
			ITextRunsIterator iRuns(ranges.GetTextRunsIterator());
			for (iRuns.MoveToFirst(); iRuns.IsNotDone(); iRuns.Next()) {
				ITextRange runRange = iRuns.Item();
				ai::AutoBuffer<ASUnicode> contents(runRange.GetSize());
				runRange.GetContents(contents.GetBuffer(), runRange.GetSize());
				log->WritePrintf("Text run %d:", count);
				SnippetRunnerLog::Indent indent;
				log->Write(ai::UnicodeString(contents.GetBuffer(), contents.GetCount()));
				count++;
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
ASErr SnpTextIterator::IterateSelectedCharacterKernTypes()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		TextRangesRef rangesRef = NULL;
		result = sAIDocument->GetTextSelection(&rangesRef);
		aisdk::check_ai_error(result);

		ITextRanges ranges(rangesRef);	
		if (ranges.GetSize() > 0) {
			SnippetRunnerLog::Indent indent;
			for (ASInt32 i = 0; i < ranges.GetSize(); i++) {
				ITextRange range = ranges.Item(i);
				log->WritePrintf("Text Range %d characters:",i);
				IStory story = range.GetStory();
				SnippetRunnerLog::Indent indent2;
				for (ASInt32 charIndex = 0; charIndex < story.GetSize(); charIndex++) {
					ASInt32 kernRef = 0;
					AutoKernType kernType = story.GetModelKernAtChar(charIndex, &kernRef);
					switch (kernType) {
						case kNoAutoKern:
							log->WritePrintf("%d. %s", charIndex, "kNoAutoKern");
							break;
						case kMetricKern:
							log->WritePrintf("%d. %s", charIndex, "kMetricKern");
							break;
						case kOpticalKern:
							log->WritePrintf("%d. %s", charIndex, "kOpticalKern");
							break;
                        case kMetricRomanOnlyKern:
                            log->WritePrintf("%d. %s", charIndex, "kMetricRomanOnlyKern");
							break;
					}
				}
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

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpTextIterator available to the SnippetRunner framework.
 */
class _SnpRunnableTextIterator : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableTextIterator () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableTextIterator () {}

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

};

/*
*/
std::string _SnpRunnableTextIterator::GetName() const
{
	return "TextIterator";
}

/*
*/
std::string _SnpRunnableTextIterator::GetDescription() const 
{
	return "Demonstrates iteration of text items.";
}

/*
*/
SnpRunnable::Operations _SnpRunnableTextIterator::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("IterateSelectedCharacters", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("IterateSelectedWords", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("IterateSelectedParagraphs", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("IterateStories", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("IterateSelectedStories", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("IterateSelectedTextRanges", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("IterateTextFrames", "selected text frames", kSnpRunCharactersContext));
	operations.push_back(Operation("IterateSelectedTextFrames", "selected text frames", kSnpRunCharactersContext));
	operations.push_back(Operation("IterateLinesInSelectedFrames", "selected text frames", kSnpRunCharactersContext));
	operations.push_back(Operation("IterateSelectedGlyphRuns", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("IterateSelectedTextRuns", "selected text", kSnpRunCharactersContext));
	operations.push_back(Operation("IterateSelectedCharacterKernTypes", "selected text", kSnpRunCharactersContext));

	return operations;
}

/*
*/
std::string _SnpRunnableTextIterator::GetDefaultOperationName() const 
{
	return "IterateSelectedCharacters";
	// Note: this must be one of the operation names returned by GetOperations.
}

/*
*/
std::vector<std::string> _SnpRunnableTextIterator::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Text Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableTextIterator::CanRun(SnpRunnable::Context& runnableContext)
{
	SnpSelectionHelper selectionHelper;
	if ("IterateStories" == runnableContext.GetOperation().GetName() ||
		"IterateTextFrames" == runnableContext.GetOperation().GetName())
		return selectionHelper.IsDocumentSelected();
	else
		return selectionHelper.IsTextRangeSelected();
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableTextIterator::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpTextIterator instance;
	if ("IterateSelectedCharacters" == runnableContext.GetOperation().GetName()) {
		result = instance.IterateSelectedCharacters();
	}
	else if ("IterateSelectedWords" == runnableContext.GetOperation().GetName()) {
		result = instance.IterateSelectedWords();
	}
	else if ("IterateSelectedParagraphs" == runnableContext.GetOperation().GetName()) {
		result = instance.IterateSelectedParagraphs();
	}
	else if ("IterateStories" == runnableContext.GetOperation().GetName()) {
		result = instance.IterateStories();
	}
	else if ("IterateSelectedStories" == runnableContext.GetOperation().GetName()) {
		result = instance.IterateSelectedStories();
	}
	else if ("IterateSelectedTextRanges" == runnableContext.GetOperation().GetName()) {
		result = instance.IterateSelectedTextRanges();
	}
	else if ("IterateTextFrames" == runnableContext.GetOperation().GetName()) {
		result = instance.IterateTextFrames();
	}
	else if ("IterateSelectedTextFrames" == runnableContext.GetOperation().GetName()) {
		result = instance.IterateSelectedTextFrames();
	}
	else if ("IterateLinesInSelectedFrames" == runnableContext.GetOperation().GetName()) {
		result = instance.IterateLinesInSelectedFrames();
	}
	else if ("IterateSelectedGlyphRuns" == runnableContext.GetOperation().GetName()) {
		result = instance.IterateSelectedGlyphRuns();
	}
	else if ("IterateSelectedTextRuns" == runnableContext.GetOperation().GetName()) {
		result = instance.IterateSelectedTextRuns();
	}
	else if ("IterateSelectedCharacterKernTypes" == runnableContext.GetOperation().GetName()) {
		result = instance.IterateSelectedCharacterKernTypes();
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableTextIterator instance_SnpRunnableTextIterator;

// End SnpTextIterator.cpp

