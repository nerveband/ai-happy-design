//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/StrokeFilter/Source/StrokeFileterParams.h $
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

#ifndef __StrokeFileterParams_H__
#define __StrokeFileterParams_H__

const ASInt32 DASH_ARRAY_LEN = 6;

class StrokeFileterParamsDash
{
public:
	StrokeFileterParamsDash();
	StrokeFileterParamsDash(const StrokeFileterParamsDash &rhs);
	StrokeFileterParamsDash& operator=(const StrokeFileterParamsDash &rhs);

	AIBoolean dashed;		// true - dashed stroke; false - solid stroke
	AIBoolean adjust;		// true - adjust dash
	AIReal dashes[DASH_ARRAY_LEN];		// dash and gap length in pt
};

class StrokeFileterParamsWave
{
public:
	StrokeFileterParamsWave();
	StrokeFileterParamsWave(const StrokeFileterParamsWave &rhs);
	StrokeFileterParamsWave& operator=(const StrokeFileterParamsWave &rhs);

	AIBoolean exchange;		// true - exchange width1 and width2
	AIBoolean replace;		// true - replace existing width parameters; false - append
	int number;             // number of wave sections
	AIReal width1;			// the left width of a wave section
	AIReal width2;			// the right width of a wave section
};

#endif
