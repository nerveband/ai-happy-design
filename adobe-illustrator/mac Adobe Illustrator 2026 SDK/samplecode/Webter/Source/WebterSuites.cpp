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
#include "WebterSuites.h"

extern "C" {
	SPBlocksSuite				*sSPBlocks = nullptr;
	AIArtSuite					*sAIArt = nullptr;
	AIDocumentSuite				*sAIDocument = nullptr;
	AIMatchingArtSuite			*sAIMatchingArt = nullptr;
	AIMdMemorySuite				*sAIMdMemory = nullptr;
	AIPreferenceSuite			*sAIPreference = nullptr;
	AIMenuSuite					*sAIMenu = nullptr;
	AIDictionarySuite			*sAIDictionary = nullptr;
	AIDictionaryIteratorSuite	*sAIDictionaryIterator = nullptr;
	AIUnicodeStringSuite		*sAIUnicodeString = nullptr;
	AIStringFormatUtilsSuite	*sAIStringFormatUtils = nullptr;
	AIArtSet					*sAIArtSet = nullptr;
	AIArtSetGeneratorSuite		*sAIArtSetGenerator = nullptr;
	AICountedObjectSuite		*sAICountedObject = nullptr;
	AIArraySuite				*sAIArray = nullptr;
	AIEntrySuite				*sAIEntry = nullptr;
}

ImportSuite gImportSuites[] = 
{
	kAIArtSuite, kAIArtSuiteVersion, &sAIArt,
	kAIDocumentSuite, kAIDocumentVersion, &sAIDocument,
	kAIMatchingArtSuite, kAIMatchingArtSuiteVersion, &sAIMatchingArt,
	kAIMdMemorySuite, kAIMdMemorySuiteVersion, &sAIMdMemory,
	kAIPreferenceSuite, kAIPreferenceSuiteVersion, &sAIPreference,
	kAIMenuSuite, kAIMenuSuiteVersion, &sAIMenu,
	kAIDictionarySuite, kAIDictionarySuiteVersion, &sAIDictionary,
	kAIDictionaryIteratorSuite, kAIDictionaryIteratorSuiteVersion, &sAIDictionaryIterator,
	kAIUnicodeStringSuite, kAIUnicodeStringVersion,&sAIUnicodeString,
	kAIStringFormatUtilsSuite, kAIStringFormatUtilsSuiteVersion, &sAIStringFormatUtils,
	kSPBlocksSuite, kSPBlocksSuiteVersion, &sSPBlocks,
	kAIArtSetSuite, kAIArtSetSuiteVersion, &sAIArtSet,
	kAIArtSetGeneratorSuite, kAIArtSetGeneratorSuiteVersion, &sAIArtSetGenerator,
	kAICountedObjectSuite, kAICountedObjectSuiteVersion, &sAICountedObject,
	kAIArraySuite, kAIArraySuiteVersion, &sAIArray,
	kAIEntrySuite, kAIEntrySuiteVersion, &sAIEntry,
	nullptr, 0, nullptr
};

