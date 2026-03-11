//========================================================================================
//  
//  $File: //ai/mainline/devtech/sdk/public/samplecode/ScriptMessage/Source/ScriptMessageSuites.cpp $
//
//  $Revision: #1 $
//
//  Copyright 2010 Adobe Systems Incorporated. All rights reserved.
//  
//  NOTICE:  Adobe permits you to use, modify, and distribute this file in accordance 
//  with the terms of the Adobe license agreement accompanying it.  If you have received
//  this file from a source other than Adobe, then your use, modification, or 
//  distribution of it requires the prior written permission of Adobe.
//  
//========================================================================================

#include "IllustratorSDK.h"
#include "ScriptMessageSuites.h"

// Suite externs
extern "C"
{
	SPBlocksSuite*			sSPBlocks = NULL;
	AIUnicodeStringSuite*	sAIUnicodeString = NULL;
	AIArtSuite*				sAIArt = NULL;
	AIUIThemeSuite*			sAIUITheme = NULL;
}

// Import suites
ImportSuite gImportSuites[] = 
{
	kSPBlocksSuite, kSPBlocksSuiteVersion, &sSPBlocks,
	kAIUnicodeStringSuite, kAIUnicodeStringVersion, &sAIUnicodeString,
	kAIArtSuite, kAIArtSuiteVersion, &sAIArt,
	kAIUIThemeSuite, kAIUIThemeSuiteVersion, &sAIUITheme,
	nullptr, 0, nullptr
};
