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
#include "TransformButtonsSuites.h"

extern "C"
{
	AIMenuSuite*				sAIMenu = NULL;
	AIPreferenceSuite*			sAIPreference = NULL;
	AIMatchingArtSuite*			sAIMatchingArt = NULL;
	AIRealMathSuite*			sAIRealMath = NULL;
	AIArtSuite*					sAIArt = NULL;
	AITransformArtSuite*		sAITransformArt = NULL;
	AIUnicodeStringSuite*		sAIUnicodeString = NULL;
	AIUndoSuite*				sAIUndo = NULL;
	AIStringFormatUtilsSuite*	sAIStringFormatUtils = NULL;
	SPBlocksSuite*			sSPBlocks = NULL;

};

ImportSuite gImportSuites[] = 
{
	kAIMenuSuite, kAIMenuSuiteVersion, &sAIMenu,
	kAIPreferenceSuite, kAIPreferenceSuiteVersion, &sAIPreference,
	kAIMatchingArtSuite, kAIMatchingArtSuiteVersion, &sAIMatchingArt,
	kAIRealMathSuite, kAIRealMathSuiteVersion, &sAIRealMath,
	kAIArtSuite, kAIArtVersion, &sAIArt,
	kAITransformArtSuite, kAITransformArtSuiteVersion, &sAITransformArt,
	kAIUnicodeStringSuite, kAIUnicodeStringVersion, &sAIUnicodeString,
	kAIUndoSuite, kAIUndoSuiteVersion, &sAIUndo,
	kAIStringFormatUtilsSuite, kAIStringFormatUtilsSuiteVersion, &sAIStringFormatUtils,
	kSPBlocksSuite, kSPBlocksSuiteVersion, &sSPBlocks,
	nullptr, 0, nullptr
};

// End TransformButtonsSuites.cpp
