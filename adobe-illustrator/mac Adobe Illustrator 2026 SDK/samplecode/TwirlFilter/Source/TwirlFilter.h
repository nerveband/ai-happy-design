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

#ifndef __TWIRLFILTER_H__
#define __TWIRLFILTER_H__

#include "IllustratorSDK.h"
#include "TwirlFilterID.h"

#define DEFAULT_ANGLE	_ShortToAIReal(50)
#define kBadSelectionErr  	'!sel'

// This structure will be passed back to Illustrator, so you can store 
// information here that you will use next time the filter is called.  
// This is different from the global variables block of the plug-in. 
// Since a plug-in might have multiple filters, each one can have parameter
// block with different values, and the correct parameter block will be
// passed to the filter when it is called.  This is easier than storing all
// possible parameters in the globals block.  The structure is allocated using
// the MDMemory suite functions. 
typedef struct
{
	AIReal angle;
	AIReal twirl;
	AIRealPoint center;
	AIReal diagonal;
}
parms_t;

class TwirlFilter 
{
public:
	ASErr TwirlGetParms(AIFilterMessage* pb);
	ASErr TwirlGetParms(AILiveEffectEditParamMessage* pb);
    ASErr GoTwirl(AIFilterMessage* pb);
    ASErr GoTwirl(AILiveEffectGoMessage* pb);

	ASErr GetSelection(AIArtHandle*** matches, ai::int32* numMatches, ai::int32* totalSegments);
	void ReleaseSelection(AIArtHandle **matches);

	ASErr TwirlPath(parms_t &pb, AIArtHandle art, ai::int32 progress_current, ai::int32 progress_total);
    ASErr TwirlPath(AIFilterMessage *pb, AIArtHandle path, ai::int32 progress_current, ai::int32 progress_total);

	ASErr DoTwirl(parms_t parms, AIArtHandle art, ai::int16 &progress_increment);
};
#endif // End TwirlFilter.h
