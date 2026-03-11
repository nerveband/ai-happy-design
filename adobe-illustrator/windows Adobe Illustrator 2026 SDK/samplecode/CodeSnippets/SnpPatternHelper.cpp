//========================================================================================
//  
//  $File: //ai/mainline/devtech/sandbox/lomorgan/public/samplecode/CodeSnippets/SnpPatternHelper.cpp $
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

#include "IllustratorSDK.h"
#include "SnpPatternHelper.h"
#include "SnippetRunnerSuites.h"

/*
 */
ASBoolean SnpPatternHelper::IsInPatternEditMode()
{
	return sAIPattern->GetPatternEditMode(NULL, NULL);
}
