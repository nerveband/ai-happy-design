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

#ifndef __MULTIARROWTOOLSUITES_H__
#define __MULTIARROWTOOLSUITES_H__

#include "IllustratorSDK.h"
#include "Suites.hpp"
#include "AIAnnotator.h"
#include "AIAnnotatorDrawer.h"
#include "AIStringFormatUtils.h"
//TODO: Remove below adm reference
//#include "ADMBasic11.h"

extern	"C"	AIMenuSuite*			sAIMenu;
extern	"C"	AIToolSuite*			sAITool;
extern	"C"	AIUnicodeStringSuite*	sAIUnicodeString;
extern	"C"	AIPluginSuite*			sAIPlugin;
extern	"C"	SPBlocksSuite*			sSPBlocks;
extern	"C"	AIUndoSuite*			sAIUndo;
extern	"C"	AIArtSuite*				sAIArt;
extern	"C"	AIRealMathSuite*		sAIRealMath;
extern	"C"	AIPathSuite*			sAIPath;
extern	"C"	AIPathStyleSuite*		sAIPathStyle;
extern  "C" AIAnnotatorSuite*		sAIAnnotator;
extern  "C" AIAnnotatorDrawerSuite* sAIAnnotatorDrawer;
extern  "C" AIDocumentViewSuite*	sAIDocumentView;
extern	"C" AIStringFormatUtilsSuite*	sAIStringFormatUtils;
//TODO: Remove below adm reference
//extern "C" ADMBasicSuite10 *sADMBasic;

#endif // End MultiArrowToolSuites.h
