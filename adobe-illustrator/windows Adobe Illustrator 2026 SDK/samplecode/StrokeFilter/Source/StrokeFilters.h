//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/StrokeFilter/Source/TwirlFilters.h $
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

#ifndef __StrokeFilters_H__
#define __StrokeFilters_H__

#include "IllustratorSDK.h"
#include "StrokeFilterID.h"
#include "StrokeFileterParams.h"

class StrokeFilters 
{
public:
	// Launch "Dashed Stroke" filter.
	ASErr GoDashedStroke(StrokeFileterParamsDash parms_dash);
	// Launch "Waved Stroke" filter.
	ASErr GoWavedStroke(StrokeFileterParamsWave parms_wave);

	// Make stroke(s) dashed.
	ASErr beautifyStrokeDash(AIBoolean dashed, AIBoolean dashAdjust, AIReal dashes[], int num);

	// Make stroke(s) waved (with variable widths) by replacing existing width parameters.
	ASErr beautifyStrokeWaveReplace(int number, AIReal width1, AIReal width2);

	// Make stroke(s) waved (with variable widths) by appending new width parameters.
	ASErr beautifyStrokeWaveAppend(int number, AIReal width1, AIReal width2);
	// Make a stroke waved (with variable widths) by appending new width parameters.
	ASErr beautifyStrokeWaveAppend(AIArtHandle art, int number, AIReal width1, AIReal width2);

	// Get selected paths
	ASErr GetSelection(AIArtHandle*** matches, ai::int32* numMatches);
};
#endif // End TwirlFilter.h
