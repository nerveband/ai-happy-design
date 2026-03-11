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
#include "LiveDropShadowSuites.h"

extern "C"
{
	AIMenuSuite*			sAIMenu = NULL;
	AIPluginGroupSuite*		sAIPluginGroup = NULL;
	AIMatchingArtSuite*		sAIMatchingArt = NULL;
	AIArtSuite*				sAIArt = NULL;
	AIMdMemorySuite*		sAIMdMemory = NULL;
	AIPathStyleSuite*		sAIPathStyle = NULL;
	AIRealMathSuite*		sAIRealMath = NULL;
	AIGroupSuite*			sAIGroup = NULL;
	AITransformArtSuite*	sAITransformArt = NULL;
	AIUnicodeStringSuite*	sAIUnicodeString = NULL;
	SPBlocksSuite*			sSPBlocks = NULL;
};

ImportSuite gImportSuites[] = 
{
	kAIMenuSuite, kAIMenuSuiteVersion, &sAIMenu,
	kAIPluginGroupSuite, kAIPluginGroupSuiteVersion, &sAIPluginGroup,
	kAIMatchingArtSuite, kAIMatchingArtSuiteVersion, &sAIMatchingArt,
	kAIArtSuite, kAIArtSuiteVersion, &sAIArt,
	kAIMdMemorySuite, kAIMdMemorySuiteVersion, &sAIMdMemory,
	kAIPathStyleSuite, kAIPathStyleSuiteVersion, &sAIPathStyle,
	kAIRealMathSuite, kAIRealMathSuiteVersion, &sAIRealMath,
	kAIGroupSuite, kAIGroupSuiteVersion, &sAIGroup,
	kAITransformArtSuite, kAITransformArtSuiteVersion, &sAITransformArt,
	kAIUnicodeStringSuite, kAIUnicodeStringSuiteVersion, &sAIUnicodeString,
	kSPBlocksSuite, kSPBlocksSuiteVersion, &sSPBlocks,
	nullptr, 0, nullptr
};

// End LiveDropShadowSuites.cpp
