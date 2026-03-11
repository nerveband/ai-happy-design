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

#include "ATESuites.h"
#include "ATETypes.h"
#include "IText.h"
#include "ATETextSuitesImportHelper.h"

#include "AITextFrame.h"
#include "AIFont.h"
#include "AIFolders.h"
#include "AIStringFormatUtils.h"
#include "AIUITheme.h"
//TODO: Remove below adm references
//#include "ADMBasic11.h"

extern	"C"	SPBasicSuite* sSPBasic;
extern "C" AIDictionarySuite *sAIDictionary;
extern "C" AIFileFormatSuite *sAIFileFormat;
extern "C" AIMatchingArtSuite *sAIMatchingArt;
extern "C" AIMdMemorySuite *sAIMdMemory;
extern "C" AIPathSuite *sAIPath;
extern "C" AIPathStyleSuite *sAIPathStyle;
extern "C" AIPreferenceSuite* sAIPreference;
extern "C" AIToolSuite *sAITool;
extern "C" AITimerSuite *sAITimer;
extern "C" AIRuntimeSuite *sAIRuntime;
extern "C" AIUndoSuite *sAIUndo;
extern "C" AIDocumentSuite *sAIDocument;
extern "C" AIArtSuite *sAIArt;
extern "C" AILayerSuite* sAILayer;
extern "C" AIMenuSuite *sAIMenu;
extern "C" AINotifierSuite *sAINotifier;
extern "C" AIDataFilterSuite *sAIDataFilter;
extern "C" AIStringFormatUtilsSuite *sAIStringFormatUtils;
extern "C" AICSXSExtensionSuite *sAICSXSExtension;
// globally defined suites that are part of the framework
extern "C" SPAccessSuite	*sSPAccess;
extern "C" AITextFrameSuite *sAITextFrame;
extern "C" AIFontSuite *sAIFont;
extern "C" AIUnicodeStringSuite *sAIUnicodeString;
extern "C" SPBlocksSuite *sSPBlocks;
extern "C" 	AIUIThemeSuite	*sAIUITheme;
//TODO: Remove below adm reference
//extern "C" ADMBasicSuite10 *sADMBasic;

//end MarkedObjectsSuites.h
