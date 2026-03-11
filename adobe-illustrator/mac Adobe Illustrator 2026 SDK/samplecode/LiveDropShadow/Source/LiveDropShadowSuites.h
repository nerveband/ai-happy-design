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

#ifndef __LiveDropShadowSuites_h__
#define __LiveDropShadowSuites_h__

#include "IllustratorSDK.h"
#include "Suites.hpp"

#include "AIPluginGroup.h"
#include "AIMatchingArt.h"
#include "AIRealMath.h"
#include "AIGroup.h"
#include "AITransformArt.h"

extern	"C" AIMenuSuite*			sAIMenu;
extern	"C" AIPluginGroupSuite*		sAIPluginGroup;
extern	"C" AIMatchingArtSuite*		sAIMatchingArt;
extern	"C" AIArtSuite*				sAIArt;
extern	"C" AIMdMemorySuite*		sAIMdMemory;
extern	"C" AINotifierSuite*		sAINotifier;
extern	"C" AIPathStyleSuite*		sAIPathStyle;
extern	"C" AIRealMathSuite*		sAIRealMath;
extern	"C" AIGroupSuite*			sAIGroup;
extern	"C" AITransformArtSuite*	sAITransformArt;
extern	"C" AIUnicodeStringSuite*	sAIUnicodeString;
extern	"C" SPBlocksSuite*			sSPBlocks;

#endif // End LiveDropShadowSuites.h
