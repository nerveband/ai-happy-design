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

#ifndef __TEXTFILEFORMATSUITES_H__
#define __TEXTFILEFORMATSUITES_H__

#include "IllustratorSDK.h"
#include "Suites.hpp"
#include "ATETextSuitesImportHelper.h"

extern	"C"	AIUnicodeStringSuite*	sAIUnicodeString;
extern  "C" SPBlocksSuite*			sSPBlocks;
extern	"C" AIFileFormatSuite*		sAIFileFormat;
extern	"C" AIDocumentSuite*		sAIDocument;
extern	"C" AITextFrameSuite*		sAITextFrame;
extern	"C" AIArtSuite*				sAIArt;
extern	"C" AIPathSuite*			sAIPath;
extern	"C" AIMatchingArtSuite*		sAIMatchingArt;
extern	"C" AIMdMemorySuite*		sAIMdMemory;

#endif // End TextFileFormatSuites.h
