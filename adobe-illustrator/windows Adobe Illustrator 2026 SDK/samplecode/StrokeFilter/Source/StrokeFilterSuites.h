//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/StrokeFilter/Source/StrokeFilterSuites.h $
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
#ifndef __STROKEFILTERSUITES_H__
#define __STROKEFILTERSUITES_H__

#include "IllustratorSDK.h"

#include "ATESuites.h"
#include "ATETypes.h"
#include "IText.h"
#include "ATETextSuitesImportHelper.h"

#include "AITextFrame.h"
#include "AIFont.h"
#include "AIFolders.h"
#include "AIStringFormatUtils.h"
#include "AIBeautifulStrokes.h"
#include "AIArtStyleParser.h"

extern "C" AIArtSuite *sAIArt;
extern "C" AIMatchingArtSuite *sAIMatchingArt;
extern "C" AIMdMemorySuite *sAIMdMemory;
extern "C" AIMenuSuite *sAIMenu;
extern "C" AIDictionarySuite *sAIDictionary;
extern "C" AIPathSuite *sAIPath;
extern "C" AIPathStyleSuite *sAIPathStyle;
extern "C" AILayerSuite *sAILayer;
extern "C" AIUndoSuite *sAIUndo;
extern "C" AIDataFilterSuite *sAIDataFilter;

extern "C" AIBeautifulStrokesSuite *sAIBeautifulStrokes;
extern "C" AIArraySuite *sAIArray;
extern "C" AIEntrySuite *sAIEntry;
extern "C" AIFilterSuite *sAIFilter;
extern "C" AILiveEffectSuite *sAILiveEffect;
extern "C" AIRealMathSuite *sAIRealMath;
extern "C" AIArtStyleSuite *sAIArtStyle;
extern "C" AIArtStyleParserSuite *sAIArtStyleParser;

extern "C" AITextFrameSuite *sAITextFrame;
extern "C" AIFontSuite *sAIFont;
extern "C" AIUnicodeStringSuite *sAIUnicodeString;
extern "C" AIStringFormatUtilsSuite *sAIStringFormatUtils;
extern "C" SPBlocksSuite *sSPBlocks;

// globally defined suites that are part of the framework
extern "C" SPAccessSuite *sSPAccess;

#endif
//end StrokeFilterSuites.h
