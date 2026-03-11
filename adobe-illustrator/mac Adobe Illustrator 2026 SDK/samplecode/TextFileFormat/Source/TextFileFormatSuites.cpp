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
#include "TextFileFormatSuites.h"

extern "C"
{
	AIUnicodeStringSuite*	sAIUnicodeString = NULL;
	SPBlocksSuite*			sSPBlocks = NULL;
	AIFileFormatSuite*		sAIFileFormat = NULL;
	AIDocumentSuite*		sAIDocument = NULL;
	AITextFrameSuite*		sAITextFrame = NULL;
	AIArtSuite*				sAIArt = NULL;
	AIPathSuite*			sAIPath = NULL;
	AIMatchingArtSuite*		sAIMatchingArt = NULL;
	AIMdMemorySuite*		sAIMdMemory = NULL;
	EXTERN_TEXT_SUITES
};

ImportSuite gImportSuites[] = 
{
	kAIUnicodeStringSuite, kAIUnicodeStringSuiteVersion, &sAIUnicodeString,
	kSPBlocksSuite, kSPBlocksSuiteVersion, &sSPBlocks,
	kAIFileFormatSuite, kAIFileFormatVersion, &sAIFileFormat,
	kAIDocumentSuite, kAIDocumentVersion, &sAIDocument,
	kAITextFrameSuite, kAITextFrameVersion, &sAITextFrame,
	kAIArtSuite, kAIArtSuiteVersion, &sAIArt,
	kAIPathSuite, kAIPathVersion, &sAIPath,
	kAIMatchingArtSuite, kAIMatchingArtVersion, &sAIMatchingArt,
	kAIMdMemorySuite, kAIMdMemorySuiteVersion, &sAIMdMemory,
	IMPORT_TEXT_SUITES
	nullptr, 0, nullptr
};


// End TextFileFormatSuites.cpp
