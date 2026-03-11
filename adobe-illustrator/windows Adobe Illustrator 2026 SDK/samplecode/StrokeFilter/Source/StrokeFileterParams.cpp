//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/StrokeFilter/Source/StrokeFileterParams.cpp $
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
#include "IllustratorSDK.h"
#include "StrokeFileterParams.h"


StrokeFileterParamsDash::StrokeFileterParamsDash():
	dashed(FALSE),
	adjust(FALSE)
							
{
	dashes[0] = 12.0;

	for (ASInt32 i = 1; i < DASH_ARRAY_LEN; i++)
	{
		dashes[i] = 0.0;
	}
}

StrokeFileterParamsDash::StrokeFileterParamsDash(const StrokeFileterParamsDash &rhs)
{
	dashed = rhs.dashed;
	adjust = rhs.adjust;

	for (ASInt32 i = 0; i < DASH_ARRAY_LEN; i++)
	{
		dashes[i] = rhs.dashes[i];
	}
}

StrokeFileterParamsDash& StrokeFileterParamsDash::operator=(const StrokeFileterParamsDash &rhs)
{
	if(this == &rhs) 
		return *this;

	this->dashed = rhs.dashed;
	this->adjust = rhs.adjust;

	for (ASInt32 i = 0; i < DASH_ARRAY_LEN; i++)
	{
		this->dashes[i] = rhs.dashes[i];
	}

	return *this; 
}

StrokeFileterParamsWave::StrokeFileterParamsWave():
	exchange(FALSE),
	replace(TRUE),
	number(4),
	width1(1.0),
	width2(2.0)
{
}

StrokeFileterParamsWave::StrokeFileterParamsWave(const StrokeFileterParamsWave &rhs)
{ 
	exchange = rhs.exchange;
	replace = rhs.replace;
	number = rhs.number;
	width1 = rhs.width1;
	width2 = rhs.width2;
}

StrokeFileterParamsWave& StrokeFileterParamsWave::operator=(const StrokeFileterParamsWave &rhs)
{
	if(this == &rhs)
		return *this;

	this->exchange = rhs.exchange;
	this->replace = rhs.replace;
	this->number = rhs.number;
	this->width1 = rhs.width1;
	this->width2 = rhs.width2;

	return *this; 
}
