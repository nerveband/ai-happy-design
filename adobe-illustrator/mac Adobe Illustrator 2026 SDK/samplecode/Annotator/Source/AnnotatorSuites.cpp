//========================================================================================
//  
//  $File$
//
//  $Revision$
//
//  Copyright 2008 Adobe Systems Incorporated. All rights reserved.
//  
//  NOTICE:  Adobe permits you to use, modify, and distribute this file in accordance 
//  with the terms of the Adobe license agreement accompanying it.  If you have received
//  this file from a source other than Adobe, then your use, modification, or 
//  distribution of it requires the prior written permission of Adobe.
//  
//========================================================================================

#include "IllustratorSDK.h"

#include "AnnotatorSuites.h"

extern "C" 
{ 
	AIUnicodeStringSuite*		sAIUnicodeString = NULL;
	SPBlocksSuite*				sSPBlocks = NULL;
	AIAnnotatorSuite*			sAIAnnotator = NULL;
	AIAnnotatorDrawerSuite*		sAIAnnotatorDrawer = NULL;
	AIToolSuite*				sAITool = NULL;
	AIArtSetSuite*				sAIArtSet = NULL;
	AIArtSuite*					sAIArt = NULL;
	AIHitTestSuite*				sAIHitTest = NULL;
	AIDocumentViewSuite*		sAIDocumentView = NULL;
	AIDocumentSuite*			sAIDocument = NULL;
	AIMatchingArtSuite*			sAIMatchingArt = NULL;
	AIStringFormatUtilsSuite*	sAIStringFormatUtils = NULL;
}

ImportSuite gImportSuites[] = 
{
	kAIUnicodeStringSuite, kAIUnicodeStringSuiteVersion, &sAIUnicodeString,
	kSPBlocksSuite, kSPBlocksSuiteVersion, &sSPBlocks,
	kAIAnnotatorSuite, kAIAnnotatorSuiteVersion, &sAIAnnotator,
	kAIAnnotatorDrawerSuite, kAIAnnotatorDrawerSuiteVersion, &sAIAnnotatorDrawer,
	kAIToolSuite, kAIToolSuiteVersion, &sAITool,
	kAIArtSetSuite, kAIArtSetSuiteVersion, &sAIArtSet,
	kAIArtSuite, kAIArtSuiteVersion, &sAIArt,
	kAIHitTestSuite, kAIHitTestSuiteVersion, &sAIHitTest,
	kAIDocumentViewSuite, kAIDocumentViewSuiteVersion, &sAIDocumentView,
	kAIDocumentSuite, kAIDocumentSuiteVersion, &sAIDocument,
	kAIMatchingArtSuite, kAIMatchingArtSuiteVersion, &sAIMatchingArt,
	kAIStringFormatUtilsSuite, kAIStringFormatUtilsSuiteVersion, &sAIStringFormatUtils,
	nullptr, 0, nullptr
};
// End AnnotatorSuites.cpp
