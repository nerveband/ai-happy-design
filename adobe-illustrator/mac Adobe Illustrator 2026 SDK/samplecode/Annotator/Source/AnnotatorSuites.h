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

#ifndef __ANNOTATORSUITES_H__
#define __ANNOTATORSUITES_H__


#include "IllustratorSDK.h"
#include "AIAnnotator.h"
#include "AIAnnotatorDrawer.h"
#include "AIStringFormatUtils.h"
#include "AIArt.h"

extern  "C" AIUnicodeStringSuite*		sAIUnicodeString;
extern  "C" SPBlocksSuite*				sSPBlocks;
extern	"C" AIAnnotatorSuite*			sAIAnnotator;
extern	"C"	AIAnnotatorDrawerSuite*		sAIAnnotatorDrawer;
extern	"C" AIToolSuite*				sAITool;
extern	"C" AIArtSetSuite*				sAIArtSet;
extern	"C" AIArtSuite*					sAIArt;
extern	"C" AIHitTestSuite*				sAIHitTest;
extern	"C" AIDocumentViewSuite*		sAIDocumentView;
extern	"C" AIDocumentSuite*			sAIDocument;
extern	"C" AIMatchingArtSuite*			sAIMatchingArt;
extern	"C" AIStringFormatUtilsSuite*	sAIStringFormatUtils;

#endif // __ANNOTATORSUITES_H__

