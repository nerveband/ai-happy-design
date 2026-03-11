//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/StrokeFilter/Source/StrokeFilters.cpp $
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
#include "StrokeFilters.h"
#include "StrokeFilterSuites.h"
#include "SDKErrors.h"

// Launch "Dashed Stroke" filter.
ASErr StrokeFilters::GoDashedStroke(StrokeFileterParamsDash parms_dash)
{
	ASErr error = kNoErr;
	int arrayLen = 0;
	// If Dash stroke dialog
	for (int i = DASH_ARRAY_LEN - 1; i >= 0; --i)
	{
		if (!sAIRealMath->EqualWithinTol(parms_dash.dashes[i], (AIReal) 0.0, (AIReal) 0.001))
		{
			arrayLen = i + 1;
			break;
		}
	}
	
	error = beautifyStrokeDash(parms_dash.dashed, parms_dash.adjust, parms_dash.dashes, arrayLen);
	return error;
}

// Launch "Waved Stroke" filter.
ASErr StrokeFilters::GoWavedStroke(StrokeFileterParamsWave parms_wave)
{
	ASErr error = kNoErr;
	if (parms_wave.exchange)
	{
		AIReal temp = parms_wave.width1;
		parms_wave.width1 = parms_wave.width2;
		parms_wave.width2 = temp;
	}

	if (parms_wave.replace)
	{
		error = beautifyStrokeWaveReplace(parms_wave.number, parms_wave.width1, parms_wave.width2);
		if (error) {goto errorTag;}
	}
	else
	{
		error = beautifyStrokeWaveAppend(parms_wave.number, parms_wave.width1, parms_wave.width2);
		if (error) {goto errorTag;}
	}
errorTag:
	return error;
}

/*
*/
ASErr StrokeFilters::beautifyStrokeDash(AIBoolean dashed, AIBoolean dashAdjust, AIReal dashArray[], int dashArrayLen)
{
	ASErr error = kNoErr;

	// get current path style
	AIPathStyle pathStyle;
	AIPathStyleMap pathStyleMap;
	AIDictionaryRef advParams = NULL;
	error = sAIPathStyle->GetCurrentPathStyle(&pathStyle, &pathStyleMap, &advParams, nullptr);
	if (error) goto errorTag;

	pathStyleMap.stroke.dash.length = true;

	if (dashed)
	{
		// make the line dashed
		// set dash and gap length
		dashArrayLen = (dashArrayLen < kMaxDashComponents) ? dashArrayLen : kMaxDashComponents;
		pathStyle.stroke.dash.length = dashArrayLen;
		for (int i = 0; i < dashArrayLen; i++)
		{
			pathStyle.stroke.dash.array[i] = dashArray[i];		
			pathStyleMap.stroke.dash.array[i] = true;
		}

		// set dash adjustment
		if (!advParams)
		{
			// The existing selection had no advanced stroke parameters
			error = sAIBeautifulStrokes->CreateParamsDict(NULL, advParams);
			if (error) goto errorTag;
		}
		error = sAIBeautifulStrokes->SetDashAdjustment(advParams, dashAdjust);
		if (error) goto errorTag;
	}
	else
	{
		// make the line solid
		pathStyle.stroke.dash.length = 0;
	}

	// set current path style
	error = sAIPathStyle->SetCurrentPathStyle(&pathStyle, &pathStyleMap, advParams);

	// fall through errorTag
errorTag:
	if (advParams)
	{
		sAIDictionary->Release(advParams);
	}
	return error;
}

/*
*/
ASErr StrokeFilters::beautifyStrokeWaveReplace(int number, AIReal width1, AIReal width2)
{
	ASErr error = kNoErr;
	
	AIDictionaryRef advParams = NULL;	// beautiful-stroke parameters
	AIArrayRef tVals = NULL;			// Array of fractional distances
	AIArrayRef widths = NULL;			// Array of stroke widths
	
	try
	{
		// get current path style
		AIPathStyle style;
		AIPathStyleMap map;
		error = sAIPathStyle->GetCurrentPathStyle(&style, &map, &advParams, nullptr);
		aisdk::check_ai_error(error);
		
		// Change stroke width (unit pt) to the larger one of width1 and width2.
		// This number is showed on Stroke panel.
		AIReal strokeWidth = width1 > width2 ? width1 : width2;
		style.stroke.width = strokeWidth;
		
		//------------------------------
		
		// Beautiful-stroke parameters
		if (!advParams)
		{
			error = sAIBeautifulStrokes->CreateParamsDict(NULL, advParams);
			aisdk::check_ai_error(error);
		}
		
		// width source
		error = sAIBeautifulStrokes->SetWidthSource(advParams, kUserProfile);
		aisdk::check_ai_error(error);
		
		// Array of fractional distances
		error = sAIArray->CreateArray(&tVals);
		AIReal section = sAIRealMath->AIRealDiv((AIReal) 1.0, (AIReal) number);
		AIReal fraction = (AIReal) 0.0;
		for (int i = 0; i < number; ++i)
		{
			error = sAIArray->AppendEntry(tVals, sAIEntry->FromReal(fraction));
			fraction = sAIRealMath->AIRealAdd(fraction, section);
		}
		error = sAIArray->AppendEntry(tVals, sAIEntry->FromReal((AIReal) 1.0));
		
		// Array of stroke widths
		error = sAIArray->CreateArray(&widths);
		AIReal widthArray[] = {width1/strokeWidth, width2/strokeWidth}; // relative widths
		for (int i = 0; i <= number; ++i)
		{
			error = sAIArray->AppendEntry(widths, sAIEntry->FromReal(widthArray[i%2]));
		}
		
		// width profile
		error = sAIBeautifulStrokes->SetCenteredWidthProfile(advParams, tVals, widths);
		aisdk::check_ai_error(error);
		
		// set current path style
		error = sAIPathStyle->SetCurrentPathStyle(&style, &map, advParams);
		aisdk::check_ai_error(error);
		
		// release memory
		if (advParams)
		{
			sAIDictionary->Release(advParams);
			advParams = NULL;
		}
		
		if (tVals)
		{
			sAIArray->Release(tVals);
			tVals = NULL;
		}
		
		if (widths)
		{
			sAIArray->Release(widths);
			widths = NULL;
		}
	
	}
	catch(...)
	{
		// release memory
		if (advParams)
		{
			sAIDictionary->Release(advParams);
			advParams = NULL;
		}
		
		if (tVals)
		{
			sAIArray->Release(tVals);
			tVals = NULL;
		}
		
		if (widths)
		{
			sAIArray->Release(widths);
			widths = NULL;
		}
	}

	return error;
}

/*
*/
ASErr StrokeFilters::beautifyStrokeWaveAppend(int number, AIReal width1, AIReal width2)
{
	ASErr error = kNoErr;
	ai::int32 progress_total, progress_current, count;
	AIArtHandle** matches = NULL;

	error = GetSelection(&matches, &count);
	if (error) goto errorTag;
	
	progress_total = count;

	for (ai::int32 i = 0; i < count; ++i )
	{
		AIArtHandle artHandle = (*matches)[i];

		short artType;
		error = sAIArt->GetArtType(artHandle, &artType);
		if ( error ) goto errorTag;

		if (artType == kPathArt)
		{
			beautifyStrokeWaveAppend(artHandle, number, width1, width2);
		}

		progress_current += 1;
		sAIUser->UpdateProgress(progress_current, progress_total);

		if ( sAIUser->Cancel() ) {
			error = kCanceledErr;
			goto errorTag;
		}
	}

	sAIMdMemory->MdMemoryDisposeHandle((AIMdMemoryHandle) matches);
	return kNoErr;

errorTag:
	if (matches)
	{
		sAIMdMemory->MdMemoryDisposeHandle((AIMdMemoryHandle) matches);
	}
	return error;
}

/*
*/
ASErr StrokeFilters::beautifyStrokeWaveAppend(AIArtHandle art, int number, AIReal width1, AIReal width2)
{
	AIErr error = kNoErr;

	AIReal segment = sAIRealMath->AIRealDiv((AIReal) 1.0, (AIReal) number);
	AIReal widths[] = {width1, width2};
	WidthMarkerFlags options = kNoWidthMarkerFlags;

	AIReal fraction = (AIReal) 0.0;
	for (int i = 0; i <= number; ++i)
	{
		error = sAIBeautifulStrokes->SetCenteredWidthAtDistance(art, fraction, widths[i%2], options);
		fraction = sAIRealMath->AIRealAdd(fraction, segment);
	}

	return error;
}

/*
*/
ASErr StrokeFilters::GetSelection(AIArtHandle*** matches, ai::int32* numMatches)
{
	ASErr error = kNoErr;

	*matches = nullptr;

	AIMatchingArtSpec spec[2];
	spec[0].type = kPathArt;
	spec[0].whichAttr = kArtSelected;
	spec[0].attr = kArtSelected;

	spec[1].type = kMysteryPathArt;
	spec[1].whichAttr = kArtSelected;
	spec[1].attr = kArtSelected;

	error = sAIMatchingArt->GetMatchingArt(spec, 2, matches, numMatches);
	if (error) { goto errorTag; }

	return error;

errorTag:
	if (*matches)
	{
		sAIMdMemory->MdMemoryDisposeHandle((AIMdMemoryHandle)*matches);
	}

	return error;
}


// End StrokeFilters.cpp
