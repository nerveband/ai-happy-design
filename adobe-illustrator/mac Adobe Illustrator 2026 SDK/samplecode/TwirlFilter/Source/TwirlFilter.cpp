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

//	Twirl anchor points around the path's center.

#include "IllustratorSDK.h"
#include "TwirlFilter.h"
#include "TwirlFilterSuites.h"

#define kPrecisionVal 2
#define maxNewSetInputChars 16

static AIBoolean IsValidInput(AIValidationStage validationStage, const ai::UnicodeString &strValue, ai::UnicodeString &errorStr, void *userData) 
{	
    ai::UnicodeString errorMessage("Please enter a value between -360 and 360");
	switch (validationStage) {
		case kAIValidationStage_interactive:
		case kAIValidationStage_onOk:
			if (strValue.empty()) 
			{
				return false;
			}
			else
			{
				ai::NumberFormat numFormat;
				float fValue;
				ASErr err = sAIStringFormatUtils->StringToReal(numFormat, strValue, fValue);
				if(err == kNoErr && fValue > -360.0 && fValue < 360.0)
				{
					*(float*)userData = fValue;
					return true;
				}
				else
                {
                    errorStr = errorMessage;
					return false;
                }
			}
		default:
            errorStr = errorMessage;
			return false;
	}
}

ASErr TwirlFilter::TwirlGetParms(AILiveEffectEditParamMessage* pb)
{
	parms_t parms;
	ASBoolean canceled = false;
	parms_t orig_parms;
	AIErr error;
	AIArtHandle **matches;
	ai::int32 count, total_segments=4;
	AILiveEffectParameters liveEffectDictionary = pb->parameters;
	AIDictKey dictKey;
    
	error = GetSelection( &matches, &count, &total_segments );
	if ( error ) goto errorTag;
    
	sAIMdMemory->MdMemoryDisposeHandle((AIMdMemoryHandle)matches);
    
	dictKey = sAIDictionary->Key("angle");
	if ( sAIDictionary->IsKnown( liveEffectDictionary, dictKey ) )
	{
		sAIDictionary->GetRealEntry( liveEffectDictionary, dictKey, &(parms.angle) );
		dictKey = sAIDictionary->Key("position");
	}
	else
	{
		parms.angle = DEFAULT_ANGLE;
	}
	
	orig_parms = parms;
	{
		ai::UnicodeString label("Angle:");
		ai::UnicodeString title(kTwirlDialogName);
		ai::UnicodeString angleUString;
		sAIUser->IUAIRealToString(DEFAULT_ANGLE, kPrecisionVal, angleUString);
		IsValidInputProc proc = IsValidInput;
		float fValue =0;
		error = sAIUser->GetInputFromUser(title, NULL, label, angleUString, proc, &fValue, maxNewSetInputChars);
        
		if (!error)
		{
			dictKey = sAIDictionary->Key("angle");
			sAIDictionary->SetRealEntry( liveEffectDictionary, dictKey, fValue );
			dictKey = sAIDictionary->Key("position");
			sAIDictionary->SetBooleanEntry( liveEffectDictionary, dictKey, TRUE );
			sAILiveEffect->UpdateParameters( pb->context );
			canceled = false;
		}
		else
		{
			parms = orig_parms;
			canceled = true;
		}
	}
    
	return (canceled ? kCanceledErr : kNoErr);
    
errorTag:
	return (error);
}			 

ASErr TwirlFilter::TwirlGetParms(AIFilterMessage* pb)
{
	parms_t **parms;
	ASBoolean canceled = false;
	parms_t orig_parms;
	AIErr error;
	AIArtHandle **matches;
	ai::int32 count, total_segments=4;
    
	error = GetSelection( &matches, &count, &total_segments );
	if ( error ) 
		goto errorTag;
    
	sAIMdMemory->MdMemoryDisposeHandle((AIMdMemoryHandle)matches);
    
	parms = (parms_t **) pb->parameters;
	if ( !parms ) {
        
		error = sAIMdMemory->MdMemoryNewHandle( sizeof(parms_t), (AIMdMemoryHandle *)&pb->parameters);
		if ( error )
			goto errorTag;
        
		if ( !pb->parameters ) {
			error = kOutOfMemoryErr;
			goto errorTag;
        }
        
		parms = (parms_t **) pb->parameters;
        
		(**parms).angle = DEFAULT_ANGLE;
	}		
	
	orig_parms = **parms;
	
	{
		ai::UnicodeString title(kTwirlDialogName);
		ai::UnicodeString label("Angle:");
		ai::UnicodeString angleUString;
		sAIUser->IUAIRealToString(orig_parms.angle, kPrecisionVal, angleUString);
		IsValidInputProc proc = IsValidInput;
		float angleValue = 0.0;
        
        ai::UnicodeString message("TwirlFilter works better with many points");
        
		error = sAIUser->GetInputFromUser(title, &message, label, angleUString, proc, &angleValue, 16);
        
		if (!error)
		{
			(**parms).angle = angleValue;
			canceled = false;
		}
		else
		{
			**parms = orig_parms;
			canceled = true;
		}
	}
	return (canceled ? kCanceledErr : kNoErr);
    
errorTag:
	return (error);
}

ASErr TwirlFilter::GetSelection(AIArtHandle*** matches, ai::int32* numMatches, ai::int32* totalSegments)
{
	AIMatchingArtSpec spec[2];
	ai::int32 i;
	ASErr error = kNoErr;

	*matches = nullptr;

	spec[0].type = kPathArt;
	spec[0].whichAttr = kArtSelected;
	spec[0].attr = kArtSelected;

	spec[1].type = kMysteryPathArt;
	spec[1].whichAttr = kArtSelected;
	spec[1].attr = kArtSelected;

	error = sAIMatchingArt->GetMatchingArt( spec, 2, matches, numMatches );
	if ( error )
	{
		ReleaseSelection(*matches);
		return error;
	}

	if ( !*numMatches )
	{
		error = kBadSelectionErr;
		goto errorTag;
	}

	*totalSegments = 0;
	for ( i = 0; i < *numMatches; ++i )
	{
		ai::int16 count;
		error = sAIPath->GetPathSegmentCount( (**matches)[i], &count );

		if ( error )
			goto errorTag;

		*totalSegments += count;
	}

	return kNoErr;

errorTag:
	ReleaseSelection(*matches);

	return error;
}


void TwirlFilter::ReleaseSelection(AIArtHandle** matches)
{
	if (matches) 
		sAIMdMemory->MdMemoryDisposeHandle((AIMdMemoryHandle)matches);
}

ASErr TwirlFilter::GoTwirl(AIFilterMessage* pb)
{
	parms_t** parms = (parms_t **) pb->parameters;
	AIArtHandle **matches = nullptr;
	ai::int32 i, count, progress_current, progress_total;
	ASBoolean bounded;
	AIRealRect bounds;
	AIErr error = kNoErr;
    
	error = GetSelection(&matches, &count, &progress_total);
	if ( error ) goto errorTag;
    
	(**parms).twirl = sAIRealMath->DegreeToRadian((**parms).angle);
    
	bounded = false;
	for ( i = 0; i < count; ++i ) {
		AIRealRect one_bounds;
        
		error = sAIArt->GetArtBounds( (*matches)[i], &one_bounds ); // not (*matches[i]). njn
		if ( error ) goto errorTag;
		if ( !bounded ) {
			bounds = one_bounds;
			bounded = true;
		}
		else sAIRealMath->AIRealRectUnion( &bounds, &one_bounds, &bounds );
	}
    
    
	(**parms).center.h = sAIRealMath->AIRealDiv( bounds.right + bounds.left, _ShortToAIReal(2) );
	(**parms).center.v = sAIRealMath->AIRealDiv( bounds.bottom + bounds.top, _ShortToAIReal(2) );
    
	(**parms).diagonal = sAIRealMath->AIRealPointLength( (AIRealPoint *) &bounds, &(**parms).center );
    
    
	progress_current = 0;
	for ( i = 0; i < count; ++i ) {
        
		short progress_increment;
		AIArtHandle path = (*matches)[i];
        
		error = sAIPath->GetPathSegmentCount( path, &progress_increment );
		if ( error ) goto errorTag;
        
		error = TwirlPath( pb, path, progress_current, progress_total );
		if ( error ) goto errorTag;
        
		progress_current += progress_increment;
        
		sAIUser->UpdateProgress( progress_current, progress_total );
		if ( sAIUser->Cancel() ) {
			error = kCanceledErr;
			goto errorTag;
		}
	}
    
	sAIMdMemory->MdMemoryDisposeHandle( (AIMdMemoryHandle)matches );
	return kNoErr;
    
errorTag:
	if ( matches )
		sAIMdMemory->MdMemoryDisposeHandle((AIMdMemoryHandle) matches );
    
	return error;
}

ASErr TwirlFilter::TwirlPath(AIFilterMessage *pb, AIArtHandle path, ai::int32 progress_current,	ai::int32 progress_total)
{
	parms_t **parms = (parms_t **) pb->parameters;
	ai::int16 count, i;
	AIRealPoint center;
	AIReal diagonal, twirl;
	ASErr error = kNoErr;
    
	center = (**parms).center;
	diagonal = (**parms).diagonal;
	twirl = (**parms).twirl;
    
	error = sAIPath->GetPathSegmentCount( path, &count );
	if ( error ) goto errorTag;
    
	for ( i = 0; i < count; ++i ) {
        
		AIPathSegment segment;
		AIRealMatrix transform;
		AIReal angle;
		AIReal length;
        
		error = sAIPath->GetPathSegments( path, i, 1, &segment );
		if ( error ) goto errorTag;
        
		length = sAIRealMath->AIRealPointLength( &center, &segment.p );
		if ( length < 0 ) length = -length;
		angle = sAIRealMath->AIRealMul( twirl, sAIRealMath->AIRealDiv( length, diagonal ) );
        
		sAIRealMath->AIRealMatrixSetTranslate( &transform, -center.h, -center.v );
		sAIRealMath->AIRealMatrixConcatRotate( &transform, angle );
		sAIRealMath->AIRealMatrixConcatRotate( &transform, -twirl );
		sAIRealMath->AIRealMatrixConcatTranslate( &transform, center.h, center.v );
        
		sAIRealMath->AIRealMatrixXformPoint( &transform, &segment.p,   &segment.p );
		sAIRealMath->AIRealMatrixXformPoint( &transform, &segment.in,  &segment.in );
		sAIRealMath->AIRealMatrixXformPoint( &transform, &segment.out, &segment.out );
        
		error = sAIPath->SetPathSegments( path, i, 1, &segment );
		if ( error ) goto errorTag;
        
		sAIUser->UpdateProgress( progress_current + i, progress_total );
        
	}
    
	return kNoErr;
    
errorTag:
	return error;
}

ASErr TwirlFilter::DoTwirl(parms_t parms, AIArtHandle art, ai::int16 &progress_increment)
{
	ASErr result = kNoErr;
	AIArtHandle child;
	short artType= kUnknownArt;	

	result = sAIArt->GetArtFirstChild(art, &child);
	while (child && !result )
	{
		result = sAIArt->GetArtType( child, &artType );
		if ( result )
			goto errorTag;

		if ( artType == kPathArt )
		{
			result = sAIPath->GetPathSegmentCount( child, &progress_increment );
			if ( result )
				goto errorTag;

			result = TwirlPath( parms, child, 1, 1 );
			if ( result )
				goto errorTag;
		}
		else if( artType == kGroupArt )
		{
			result = DoTwirl(parms, child, progress_increment); 
			if ( result )
				goto errorTag;
		}
		if ( !result )
			sAIArt->GetArtSibling(child, &child);
	}
errorTag:
	return result;
}

ASErr TwirlFilter::GoTwirl(AILiveEffectGoMessage* pb)
{
	AILiveEffectParameters liveEffectDictionary = pb->parameters;
	parms_t parms;
	ai::int16 progress_increment;
	ai::int32 progress_current;
	ASErr error = kNoErr;
	AIArtHandle path = pb->art;
	AIArtHandle anotherart = pb->art;
    
	short artType;
	AIDictKey dictKey;
    
	dictKey = sAIDictionary->Key("angle");
	if (sAIDictionary->IsKnown(liveEffectDictionary, dictKey)){
		sAIDictionary->GetRealEntry(liveEffectDictionary, dictKey, &(parms.angle));
		dictKey = sAIDictionary->Key("position");
	}
	else {
		return kNoErr;
	}
    
	parms.twirl = sAIRealMath->DegreeToRadian( parms.angle );
    
	AIRealRect one_bounds;
    
	error = sAIArt->GetArtBounds( pb->art, &one_bounds );
	if ( error ) goto errorTag;
	error = sAIArt->GetArtType( anotherart, &artType );
	if ( error ) goto errorTag;
    
	parms.center.h = sAIRealMath->AIRealDiv( one_bounds.right + one_bounds.left, _ShortToAIReal(2) );
	parms.center.v = sAIRealMath->AIRealDiv( one_bounds.bottom + one_bounds.top, _ShortToAIReal(2) );
    
	parms.diagonal = sAIRealMath->AIRealPointLength( (AIRealPoint *) &one_bounds, &(parms.center) );
    
	DoTwirl(parms, anotherart, progress_increment);
    
	path = anotherart ? anotherart : path;
    
	progress_current = 0;
    
    
	progress_current += progress_increment;
    
	if ( sAIUser->Cancel() ) 
	{
		error = kCanceledErr;
	}
    
errorTag:
    
	return error;
}

ASErr TwirlFilter::TwirlPath(parms_t &pb, AIArtHandle path, ai::int32 progress_current, ai::int32 progress_total)
{
	ai::int16 count, i;
	AIRealPoint center;
	AIReal diagonal, twirl;
	ASErr error = kNoErr;

	center = pb.center;
	diagonal = pb.diagonal;
	twirl = pb.twirl;

	error = sAIPath->GetPathSegmentCount( path, &count );
	if ( error )
		goto errorTag;

	for ( i = 0; i < count; ++i )
	{
		AIPathSegment segment;
		AIRealMatrix transform;
		AIReal angle;
		AIReal length;

		error = sAIPath->GetPathSegments( path, i, 1, &segment );
		if ( error )
			goto errorTag;

		length = sAIRealMath->AIRealPointLength( &center, &segment.p );
		if ( length < 0 )
			length = -length;

		angle = sAIRealMath->AIRealMul( twirl, sAIRealMath->AIRealDiv( length, diagonal ) );

		sAIRealMath->AIRealMatrixSetTranslate( &transform, -center.h, -center.v );
		sAIRealMath->AIRealMatrixConcatRotate( &transform, angle );
		sAIRealMath->AIRealMatrixConcatRotate( &transform, -twirl );
		sAIRealMath->AIRealMatrixConcatTranslate( &transform, center.h, center.v );

		sAIRealMath->AIRealMatrixXformPoint( &transform, &segment.p,   &segment.p );
		sAIRealMath->AIRealMatrixXformPoint( &transform, &segment.in,  &segment.in );
		sAIRealMath->AIRealMatrixXformPoint( &transform, &segment.out, &segment.out );

		error = sAIPath->SetPathSegments( path, i, 1, &segment );
		if ( error )
			goto errorTag;

		sAIUser->UpdateProgress( progress_current + i, progress_total );
	}

	return kNoErr;

errorTag:
	return error;
}

// End TwirlFilter.cpp
