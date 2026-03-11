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

#include "MarkedObjectsSuites.h"

extern "C" 
{ 
	AIArtSuite					*sAIArt;
	AIDocumentSuite				*sAIDocument;
	AIFileFormatSuite			*sAIFileFormat;
	AIMatchingArtSuite			*sAIMatchingArt;
	AIMdMemorySuite				*sAIMdMemory;
	AIMenuSuite					*sAIMenu;
	AIDictionarySuite			*sAIDictionary;
	AIDictionaryIteratorSuite	*sAIDictionaryIterator;
	AIPathSuite					*sAIPath;
	AIPathStyleSuite			*sAIPathStyle;
	AIPreferenceSuite			*sAIPreference;
	AIToolSuite					*sAITool;
	AILayerSuite				*sAILayer;
	AITimerSuite				*sAITimer;
	AIRuntimeSuite				*sAIRuntime;
	AIUndoSuite					*sAIUndo;
	AIDataFilterSuite			*sAIDataFilter;
	AIStringFormatUtilsSuite	*sAIStringFormatUtils;
	AITextFrameSuite			*sAITextFrame;
	AIFontSuite					*sAIFont;
	AIUnicodeStringSuite		*sAIUnicodeString;
	SPBlocksSuite				*sSPBlocks;
	AICSXSExtensionSuite		*sAICSXSExtension;
	AIUIThemeSuite				*sAIUITheme;
	//TODO: remove below adm reference
//	ADMBasicSuite10				*sADMBasic;
	EXTERN_TEXT_SUITES 
}

ImportSuite gImportSuites[] = 
{
	kAIArtSuite, kAIArtSuiteVersion, &sAIArt,
	kAIDocumentSuite, kAIDocumentVersion, &sAIDocument, 
	kAIFileFormatSuite, kAIFileFormatSuiteVersion, &sAIFileFormat,
	kAIMatchingArtSuite, kAIMatchingArtSuiteVersion, &sAIMatchingArt,
	kAIMdMemorySuite, kAIMdMemorySuiteVersion, &sAIMdMemory,
	kAIMenuSuite, kAIMenuSuiteVersion, &sAIMenu,
	kAIDictionarySuite, kAIDictionarySuiteVersion, &sAIDictionary,
	kAIDictionaryIteratorSuite, kAIDictionaryIteratorSuiteVersion, &sAIDictionaryIterator,
	kAIPathSuite, kAIPathSuiteVersion, &sAIPath,
	kAIPathStyleSuite, kAIPathStyleSuiteVersion, &sAIPathStyle,
	kAIPreferenceSuite, kAIPreferenceSuiteVersion, &sAIPreference,
	kAIToolSuite, kAIToolSuiteVersion, &sAITool,
	kAILayerSuite, kAILayerSuiteVersion, &sAILayer,
	kAITimerSuite, kAITimerSuiteVersion, &sAITimer,
	kAIRuntimeSuite, kAIRuntimeSuiteVersion, &sAIRuntime,
	kAIUndoSuite, kAIUndoSuiteVersion, &sAIUndo,
	kAIUnicodeStringSuite, kAIUnicodeStringSuiteVersion, &sAIUnicodeString,
	kSPBlocksSuite, kSPBlocksSuiteVersion, &sSPBlocks,
	kAITextFrameSuite, kAITextFrameSuiteVersion, &sAITextFrame,
	kAIFontSuite, kAIFontSuiteVersion, &sAIFont,
	kAIDataFilterSuite, kAIDataFilterSuiteVersion, &sAIDataFilter,
	kAIStringFormatUtilsSuite, kAIStringFormatUtilsSuiteVersion, &sAIStringFormatUtils,
	kAIUIThemeSuite, kAIUIThemeSuiteVersion, &sAIUITheme,
	//remove below adm reference
//	kADMBasicSuite, kADMBasicSuiteVersion10, &sADMBasic,
	IMPORT_TEXT_SUITES

	nullptr, 0, nullptr
};

ImportSuite gPostStartupSuites[] = 
{
	kAICSXSExtensionSuite, kAICSXSExtensionSuiteVersion, &sAICSXSExtension,
	nullptr, 0, nullptr
};
// end MarkedObjectsSuites.cpp
