//========================================================================================
//  
//  $File: //ai/mainline/devtech/sandbox/lomorgan/public/samplecode/Webter/Source/WebterSuites.h $
//
//  $Revision: #1 $
//
//  Copyright 1987 Adobe Systems Incorporated. All rights reserved.
//  
//  NOTICE:  Adobe permits you to use, modify, and distribute this file in accordance 
//  with the terms of the Adobe license agreement accompanying it.  If you have received
//  this file from a source other than Adobe, then your use, modification, or 
//  distribution of it requires the prior written permission of Adobe.
//  
//========================================================================================

#ifndef __WEBTERSUITES_H__
#define __WEBTERSUITES_H__

#include "IllustratorSDK.h"
#include "Suites.hpp"
#include "AIStringFormatUtils.h"
#include "AIArtSetGenerator.h"
#include "AICountedObject.h"

extern	"C" SPAccessSuite				*sSPAccess;
extern	"C" SPBlocksSuite				*sSPBlocks;
extern	"C" AIArtSuite					*sAIArt;
extern	"C" AIMatchingArtSuite			*sAIMatchingArt;
extern	"C" AIMdMemorySuite				*sAIMdMemory;
extern	"C" AIMenuSuite					*sAIMenu;
extern	"C" AINotifierSuite				*sAINotifier;
extern	"C" AIDocumentSuite				*sAIDocument;
extern	"C" AIDictionarySuite			*sAIDictionary;
extern	"C" AIUnicodeStringSuite		*sAIUnicodeString;
extern	"C"	AIStringFormatUtilsSuite	*sAIStringFormatUtils;
extern  "C" AIArtSet					*sAIArtSet;
extern  "C"	AIArtSetGeneratorSuite		*sAIArtSetGenerator;
extern  "C" AICountedObjectSuite		*sAICountedObject;
extern  "C" AIArraySuite				*sAIArray;
extern  "C" AIEntrySuite				*sAIEntry;

#endif
