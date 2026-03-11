//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/StrokeFilter/Source/StrokeFilterSuites.cpp $
//
//  $Revision: #2 $
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
#include "StrokeFilterSuites.h"

extern "C" 
{ 
	
	AIArtSuite *sAIArt;
	AIMatchingArtSuite *sAIMatchingArt;
	AIMdMemorySuite *sAIMdMemory;
	AIMenuSuite *sAIMenu;
	AIDictionarySuite *sAIDictionary;
	AIPathSuite *sAIPath;
	AIPathStyleSuite *sAIPathStyle;
	AILayerSuite *sAILayer;
	AIUndoSuite *sAIUndo;
	AIDataFilterSuite *sAIDataFilter;
	AIBeautifulStrokesSuite *sAIBeautifulStrokes;
	AIArraySuite *sAIArray;
	AIEntrySuite *sAIEntry;
	AIFilterSuite *sAIFilter;
	AILiveEffectSuite *sAILiveEffect;
	AIRealMathSuite *sAIRealMath;
	AIArtStyleSuite *sAIArtStyle;
	AIArtStyleParserSuite *sAIArtStyleParser;
	AITextFrameSuite *sAITextFrame;
	AIFontSuite *sAIFont;
	AIUnicodeStringSuite *sAIUnicodeString;
	AIStringFormatUtilsSuite *sAIStringFormatUtils;
	SPBlocksSuite *sSPBlocks;

	EXTERN_TEXT_SUITES 
}

ImportSuite gImportSuites[] = 
{

	kAIArtSuite, kAIArtSuiteVersion, &sAIArt,
	kAIMatchingArtSuite, kAIMatchingArtSuiteVersion, &sAIMatchingArt,
	kAIMdMemorySuite, kAIMdMemorySuiteVersion, &sAIMdMemory,
	kAIMenuSuite, kAIMenuSuiteVersion, &sAIMenu,
	kAIDictionarySuite, kAIDictionarySuiteVersion, &sAIDictionary,
	kAIPathSuite, kAIPathSuiteVersion, &sAIPath,
	kAIPathStyleSuite, kAIPathStyleSuiteVersion, &sAIPathStyle,
	kAILayerSuite, kAILayerSuiteVersion, &sAILayer,
	kAIUndoSuite, kAIUndoSuiteVersion, &sAIUndo,
	kAIDataFilterSuite, kAIDataFilterSuiteVersion, &sAIDataFilter,
	kAIBeautifulStrokesSuite, kAIBeautifulStrokesSuiteVersion, &sAIBeautifulStrokes,
	kAIArraySuite, kAIArraySuiteVersion, &sAIArray,
	kAIEntrySuite, kAIEntrySuiteVersion, &sAIEntry,
	kAIFilterSuite, kAIFilterSuiteVersion, &sAIFilter,
	kAILiveEffectSuite, kAILiveEffectVersion, &sAILiveEffect,
	kAIRealMathSuite, kAIRealMathVersion, &sAIRealMath,
	kAIArtStyleSuite, kAIArtStyleSuiteVersion, &sAIArtStyle,
	kAIArtStyleParserSuite,	kAIArtStyleParserSuiteVersion, &sAIArtStyleParser,
	kAITextFrameSuite, kAITextFrameSuiteVersion, &sAITextFrame,
	kAIFontSuite, kAIFontSuiteVersion, &sAIFont,
	kAIUnicodeStringSuite, kAIUnicodeStringSuiteVersion, &sAIUnicodeString,
	kSPBlocksSuite, kSPBlocksSuiteVersion, &sSPBlocks,
	kAIStringFormatUtilsSuite, kAIStringFormatUtilsSuiteVersion, &sAIStringFormatUtils,

	IMPORT_TEXT_SUITES

	nullptr, 0, nullptr
};
// end StrokeFilterSuites.cpp
