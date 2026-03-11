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
#include "MultiArrowToolSuites.h"

extern "C" {
	AIMenuSuite*			sAIMenu = NULL;
	AIToolSuite*			sAITool = NULL;
	AIUnicodeStringSuite*	sAIUnicodeString = NULL;
	AIPluginSuite*			sAIPlugin = NULL;
	SPBlocksSuite*			sSPBlocks = NULL;
	AIUndoSuite*			sAIUndo = NULL;
	AIArtSuite*				sAIArt = NULL;
	AIRealMathSuite*		sAIRealMath = NULL;
	AIPathSuite*			sAIPath = NULL;
	AIPathStyleSuite*		sAIPathStyle = NULL;
	AIAnnotatorSuite*		sAIAnnotator = NULL;
	AIAnnotatorDrawerSuite* sAIAnnotatorDrawer = NULL;
	AIDocumentViewSuite*	sAIDocumentView = NULL;
	AIStringFormatUtilsSuite*	sAIStringFormatUtils = NULL;
	//TODO: Remove below adm reference
	//ADMBasicSuite10 *sADMBasic = NULL;
}

ImportSuite gImportSuites[] = {
	kAIMenuSuite, kAIMenuSuiteVersion, &sAIMenu,
	kAIToolSuite, kAIToolVersion, &sAITool,
	kAIUnicodeStringSuite, kAIUnicodeStringVersion, &sAIUnicodeString,
	kAIPluginSuite, kAIPluginSuiteVersion, &sAIPlugin,
	kSPBlocksSuite, kSPBlocksSuiteVersion, &sSPBlocks,
	kAIUndoSuite, kAIUndoSuiteVersion, &sAIUndo,
	kAIArtSuite, kAIArtVersion, &sAIArt,
	kAIRealMathSuite, kAIRealMathVersion, &sAIRealMath,	
	kAIPathSuite, kAIPathSuiteVersion, &sAIPath,
	kAIPathStyleSuite, kAIPathStyleSuiteVersion, &sAIPathStyle,
	kAIAnnotatorSuite, kAIAnnotatorSuiteVersion, &sAIAnnotator,
	kAIAnnotatorDrawerSuite, kAIAnnotatorDrawerSuiteVersion, &sAIAnnotatorDrawer,
	kAIDocumentViewSuite, kAIDocumentViewSuiteVersion, &sAIDocumentView,
	kAIStringFormatUtilsSuite, kAIStringFormatUtilsSuiteVersion, &sAIStringFormatUtils,
	//TODO: Remove below adm reference
	//kADMBasicSuite, kADMBasicSuiteVersion10, &sADMBasic,
	nullptr, 0, nullptr
};

// End MultiArrowToolSuites.cpp
