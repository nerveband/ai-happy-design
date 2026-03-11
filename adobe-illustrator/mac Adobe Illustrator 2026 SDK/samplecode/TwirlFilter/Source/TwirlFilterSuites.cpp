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
#include "TwirlFilterSuites.h"

extern "C"
{
	AIMenuSuite*			sAIMenu = NULL;
	AIUnicodeStringSuite*	sAIUnicodeString = NULL;
	SPBlocksSuite*			sSPBlocks = NULL;
	AIFilterSuite*			sAIFilter = NULL;
	AILiveEffectSuite*		sAILiveEffect = NULL;
	AIMdMemorySuite*		sAIMdMemory = NULL;
	AIDictionarySuite*		sAIDictionary = NULL;
	AIMatchingArtSuite*		sAIMatchingArt = NULL;
	AIPathSuite*			sAIPath = NULL;
	AIRealMathSuite*		sAIRealMath = NULL;
	AIArtSuite*				sAIArt = NULL;
	//AIUIUtilsSuite*         sAIUIUtils = NULL;
	AIStringFormatUtilsSuite*	sAIStringFormatUtils = NULL;
};

ImportSuite gImportSuites[] = 
{
	kAIMenuSuite, kAIMenuSuiteVersion, &sAIMenu,
	kAIUnicodeStringSuite, kAIUnicodeStringSuiteVersion, &sAIUnicodeString,
	kSPBlocksSuite, kSPBlocksSuiteVersion, &sSPBlocks,
	kAIFilterSuite, kAIFilterVersion, &sAIFilter,
	kAILiveEffectSuite, kAILiveEffectVersion, &sAILiveEffect,
	kAIMdMemorySuite, kAIMdMemorySuiteVersion, &sAIMdMemory,
	kAIDictionarySuite, kAIDictionaryVersion, &sAIDictionary,
	kAIMatchingArtSuite, kAIMatchingArtSuiteVersion, &sAIMatchingArt,
	kAIPathSuite, kAIPathSuiteVersion,	&sAIPath,
	kAIRealMathSuite, kAIRealMathVersion, &sAIRealMath,
	kAIArtSuite, kAIArtVersion, &sAIArt,
	//kAIUIUtilsSuite, kAIUIUtilsSuiteVersion, &sAIUIUtils,
	kAIStringFormatUtilsSuite, kAIStringFormatUtilsSuiteVersion, &sAIStringFormatUtils,
	nullptr, 0, nullptr
};

// End TwirlFilterSuites.cpp
