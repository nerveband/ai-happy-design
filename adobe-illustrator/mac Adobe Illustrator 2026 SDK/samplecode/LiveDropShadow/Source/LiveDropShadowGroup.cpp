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

#include "IllustratorSDK.h"
#include "LiveDropShadowGroup.h"
#include "LiveDropShadowSuites.h"
#include "AITransformArt.h"

/*
*/
ASErr LiveDropShadowGroup::PluginGroupUpdate(AIPluginGroupMessage* message) 
{
	AIErr result = kNoErr;
	
	const AIReal DROP_SHADOW_OFFSET = 5;	
	
	AIArtHandle plugin = message->art;

	// Get the edit art.
	AIArtHandle editArt = NULL;
	result = sAIPluginGroup->GetPluginArtEditArt(plugin, &editArt);
	
	// Get the result art.
	AIArtHandle resultArt = NULL;
	result = sAIPluginGroup->GetPluginArtResultArt(plugin, &resultArt);

	size_t count = 0;
	result = sAIPluginGroup->GetPluginArtDataCount(plugin, &count);

	// Destroy all art that exists below the result art.
	AIArtHandle child = NULL;
	result = sAIArt->GetArtFirstChild(resultArt, &child);
	
	while ((result == kNoErr) && (child != NULL)) {	
		if (result == kNoErr) {		
			result = sAIArt->DisposeArt(child);
		}
		
		if (result == kNoErr) {
			result = sAIArt->GetArtFirstChild(resultArt, &child);
		}
	}

	// Build the drop shadow group.
	AIArtHandle dropShadowGroup = NULL;
	if (result == kNoErr) {	
		result = sAIArt->NewArt(kGroupArt, kPlaceInsideOnTop, resultArt, &dropShadowGroup);	
	}

	if (result == kNoErr) {	
		result = sAIArt->GetArtFirstChild(editArt, &child);	
	}

	AIArtHandle duplicate = NULL;
	while ((result == kNoErr) && (child != NULL)) {	
			result = sAIArt->DuplicateArt(child, kPlaceInsideOnTop, dropShadowGroup, &duplicate);
			result = sAIArt->GetArtSibling(child, &child);		
	}
	
	// Make a copy of the original edit art above the drop shadow group.
	if (result == kNoErr) {	
		result = sAIArt->GetArtFirstChild(editArt, &child);	
	}
	
	while ((result == kNoErr) && (child != NULL)) {	
			result = sAIArt->DuplicateArt(child, kPlaceAbove, dropShadowGroup, &duplicate);
			result = sAIArt->GetArtSibling(child, &child);		
	}
	
	// Set the drop shadow style.
	result = setDropShadowStyle(dropShadowGroup);
	
	// Transform the drop shadow.
	result = transformDropShadow(dropShadowGroup, DROP_SHADOW_OFFSET, -DROP_SHADOW_OFFSET, 0, 1, 1);	

	return result;		
}

/*
*/
ASErr LiveDropShadowGroup::SetFilteredFill(AIBoolean *paint, AIFillStyle *fill) 
{
	AIErr result = kNoErr;

	*paint = false;
	
	fill->color.kind = kGrayColor;
	fill->color.c.g.gray = kAIRealZero;
	
	fill->overprint = false;
	
	return result;
	
}

/*
*/
ASErr LiveDropShadowGroup::SetFilteredStroke(AIBoolean *paint, AIStrokeStyle *stroke) 
{
	AIErr result = kNoErr;

	*paint = true;

	stroke->color.kind = kGrayColor;
	stroke->color.c.g.gray = kAIRealOne;
	
	stroke->overprint = false;
	
	stroke->width = kAIRealOne;
	stroke->dash.length = 0;
	stroke->cap = kAIButtCap;
	stroke->join = kAIMiterJoin;
	stroke->miterLimit = _ShortToAIReal(4);
	
	return result;
	
}

/*
*/
ASErr LiveDropShadowGroup::SetFilteredStyle(AIArtHandle object) 
{
	AIErr result = kNoErr;

	short type;
	
	if (result == kNoErr) {	
		result = sAIArt->GetArtType(object, &type);	
	}
	AIPathStyle style;
	AIArtHandle child = NULL;
	switch (type) {	
		case kPathArt:
		case kMysteryPathArt:
		case kRasterArt:
			if (result == kNoErr) {			
				result = sAIPathStyle->GetPathStyle(object, &style, nullptr);			
			}
			if (result == kNoErr) {
				result = SetFilteredFill(&style.fillPaint, &style.fill);			
			}			
			if (result == kNoErr) {
				result = SetFilteredStroke(&style.strokePaint, &style.stroke);			
			}			
			if (result == kNoErr) {			
				result = sAIPathStyle->SetPathStyle(object, &style);			
			}
			break;			
		case kGroupArt:
		case kCompoundPathArt:
			if (result == kNoErr) {			
				result = sAIArt->GetArtFirstChild(object, &child);			
			}			
			while ((result == kNoErr) && (child != NULL)) {			
				if (result == kNoErr) {				
					result = SetFilteredStyle(child);					
				}				
				if (result == kNoErr) {				
					result = sAIArt->GetArtSibling(child, &child);				
				}				
			}
			break;			
		case kPlacedArt:		
			break;		
		case kPluginArt:
			break;
		default:		
			break;	
	}
	
	return result;	
}

/*
*/
ASErr LiveDropShadowGroup::SetBackingFill(AIBoolean *paint, AIFillStyle *fill) 
{
	AIErr result = kNoErr;

	*paint = true;
	
	fill->color.kind = kGrayColor;
	fill->color.c.g.gray = kAIRealZero;
	
	fill->overprint = false;
	
	return result;
	
}

/*
*/
ASErr LiveDropShadowGroup::SetBackingStroke(AIBoolean *paint, AIStrokeStyle *stroke) 
{
	AIErr result = kNoErr;

	*paint = false;

	stroke->color.kind = kGrayColor;
	stroke->color.c.g.gray = kAIRealOne;
	
	stroke->overprint = false;
	
	stroke->width = kAIRealOne;
	stroke->dash.length = 0;
	stroke->cap = kAIButtCap;
	stroke->join = kAIMiterJoin;
	stroke->miterLimit = _ShortToAIReal(4);
	
	return result;
	
}

/*
*/
ASErr LiveDropShadowGroup::SetBackingStyle(AIArtHandle object) 
{
	AIErr result = kNoErr;

	short type;
	if (result == kNoErr) {	
		result = sAIArt->GetArtType(object, &type);	
	}
	
	AIPathStyle style;
	AIArtHandle child = NULL;
	switch (type) {	
		case kPathArt:
		case kMysteryPathArt:
		case kRasterArt:
			if (result == kNoErr) {			
				result = sAIPathStyle->GetPathStyle(object, &style, nullptr);
			}
			if (result == kNoErr) {
				result = SetBackingFill(&style.fillPaint, &style.fill);			
			}			
			if (result == kNoErr) {
				result = SetBackingStroke(&style.strokePaint, &style.stroke);			
			}			
			if (result == kNoErr) {			
				result = sAIPathStyle->SetPathStyle(object, &style);			
			}
			break;			
		case kGroupArt:
		case kCompoundPathArt:
			if (result == kNoErr) {			
				result = sAIArt->GetArtFirstChild(object, &child);			
			}			
			while ((result == kNoErr) && (child != NULL)) {			
				if (result == kNoErr) {				
					result = SetBackingStyle( child );					
				}				
				if (result == kNoErr) {				
					result = sAIArt->GetArtSibling(child, &child);				
				}				
			}
			break;			
		case kPlacedArt:		
			break;		
		case kPluginArt:
			break;
		default:		
			break;	
	}	
	return result;	
}

/*
*/
ASErr LiveDropShadowGroup::SetBoundingFill(AIBoolean *paint, AIFillStyle *fill) 
{
	AIErr result = kNoErr;

	*paint = false;
	
	fill->color.kind = kGrayColor;
	fill->color.c.g.gray = kAIRealZero;
	
	fill->overprint = false;
	
	return result;	
}

/*
*/
ASErr LiveDropShadowGroup::SetBoundingStroke(AIBoolean *paint, AIStrokeStyle *stroke) 
{
	AIErr result = kNoErr;

	*paint = true;

	stroke->color.kind = kGrayColor;
	stroke->color.c.g.gray = kAIRealOne;
	
	stroke->overprint = false;
	
	stroke->width = kAIRealZero;
	stroke->dash.length = 0;
	stroke->cap = kAIButtCap;
	stroke->join = kAIMiterJoin;
	stroke->miterLimit = _ShortToAIReal(4);
	
	return result;	
}

/*
*/
ASErr LiveDropShadowGroup::SetBoundingStyle(AIArtHandle object) 
{
	AIErr result = kNoErr;

	short type;	
	if (result == kNoErr) {	
		result = sAIArt->GetArtType(object, &type);	
	}
	
	AIPathStyle style;
	AIArtHandle child = NULL;
	switch (type) {	
		case kPathArt:
		case kMysteryPathArt:
		case kRasterArt:		
			if (result == kNoErr) {			
				result = sAIPathStyle->GetPathStyle(object, &style, nullptr);
			}
			if (result == kNoErr) {
				result = SetBoundingFill(&style.fillPaint, &style.fill);			
			}			
			if (result == kNoErr) {
				result = SetBoundingStroke(&style.strokePaint, &style.stroke);			
			}			
			if (result == kNoErr) {			
				result = sAIPathStyle->SetPathStyle(object, &style);			
			}
			break;			
		case kGroupArt:
		case kCompoundPathArt:		
			if (result == kNoErr) {			
				result = sAIArt->GetArtFirstChild(object, &child);			
			}			
			while ((result == kNoErr) && (child != NULL)) {			
				if (result == kNoErr) {				
					result = SetBoundingStyle(child);					
				}				
				if (result == kNoErr) {				
					result = sAIArt->GetArtSibling(child, &child);				
				}				
			}
			break;			
		case kPlacedArt:		
			break;		
		case kPluginArt:
			break;
		default:		
			break;	
	}	
	return result;	
}

/*
*/
ASErr LiveDropShadowGroup::SetClippingStyle(AIArtHandle object) 
{
	AIErr result = kNoErr;

	short type;	
	if (object == NULL) {	
		return kNoErr;	
	}
	
	if (result == kNoErr) {	
		result = sAIArt->GetArtType(object, &type);	
	}
	
	AIPathStyle style;
	AIArtHandle child = NULL;
	switch (type) {	
		case kPathArt:
		case kMysteryPathArt:
		case kRasterArt:
			if (result == kNoErr) {			
				result = sAIPathStyle->GetPathStyle(object, &style, nullptr);
			}
			style.fillPaint = false;
			style.strokePaint = false;
			style.clip = true;			
			if (result == kNoErr) {			
				result = sAIPathStyle->SetPathStyle(object, &style);			
			}
			break;			
		case kGroupArt:
		case kCompoundPathArt:
			if (result == kNoErr) {			
				result = sAIArt->GetArtFirstChild(object, &child);			
			}			
			while ((result == kNoErr) && (child != NULL)) {			
				if (result == kNoErr) {				
					result = SetClippingStyle(child);
				}				
				if (result == kNoErr) {				
					result = sAIArt->GetArtSibling(child, &child);				
				}				
			}
			break;			
		case kPlacedArt:		
			break;		
		case kPluginArt:
			break;
		default:		
			break;		
	}
	
	return result;	
}

/*
*/
ASErr LiveDropShadowGroup::DuplicateCoveredArt(AIArtHandle source, AIArtHandle destination,
											   AIRealRect *bounds, AIArtHandle start, AIBoolean *found) 
{
	AIErr result = kNoErr;

	short type;
	if (result == kNoErr) {	
		result = sAIArt->GetArtType(source, &type);	
	}

	AIRealRect local = {0, 0, 0, 0};
	AIArtHandle child = NULL;
	AIArtHandle ahead = NULL;
	AIArtHandle clone = NULL;
	switch (type) {		
		case kPathArt:
		case kTextFrameArt:
		case kMysteryPathArt:
		case kPlacedArt:
		case kRasterArt:
		case kCompoundPathArt:			
			if (*found) {			
				if (result == kNoErr) {				
					result = sAIArt->GetArtBounds(source, &local);				
				}				
				if (sAIRealMath->AIRealRectOverlap(bounds, &local)) {				
					if (result == kNoErr) {				
						result = sAIArt->DuplicateArt(source, kPlaceInsideOnBottom, destination, &clone);					
					}				
				}					
			}
			break;			
		case kGroupArt:		
			if (result == kNoErr) {			
				result = sAIArt->GetArtBounds(source, &local);			
			}
			if ( sAIRealMath->AIRealRectOverlap(bounds, &local)) {			
				AIBoolean clipped = false;				
				if (result == kNoErr) {				
					result = sAIGroup->GetGroupClipped(source, &clipped);				
				}				
				if (clipped) {
					if (*found) {					
						if (result == kNoErr) {					
							result = sAIArt->DuplicateArt(source, kPlaceInsideOnBottom, destination, &clone);						
						}					
					}				
				} 
				else {					
					if (result == kNoErr) {				
						result = sAIArt->NewArt(kGroupArt, kPlaceInsideOnBottom, destination, &clone);					
					}					
					if (result == kNoErr) {					
						result = sAIArt->GetArtFirstChild(source, &child);					
					}					
					if (result == kNoErr) {					
						result = sAIArt->GetArtSibling(child, &ahead);						
					}
					while ((result == kNoErr) && (ahead != NULL)) {
						child = ahead;					
						if (result == kNoErr) {
							result = sAIArt->GetArtSibling(child, &ahead);						
						}						
					}					
					while ((result == kNoErr) && (child != NULL)) {					
						if (result == kNoErr) {						
							result = DuplicateCoveredArt(child, clone, bounds, start, found);							
						}					
						if (result == kNoErr) {						
							result = sAIArt->GetArtPriorSibling(child, &child);						
						}					
					}					
					if (result == kNoErr) {					
						result = sAIArt->GetArtFirstChild(clone, &child);					
					}					
					if (child == NULL) {					
						if (result == kNoErr) {						
							result = sAIArt->DisposeArt(clone);						
						}					
					}				
				}				
			}
			break;			
		case kPluginArt:		
			if (result == kNoErr) {			
				result = sAIArt->GetArtBounds(source, &local);			
			}
			if (sAIRealMath->AIRealRectOverlap(bounds, &local)) {			
				if (result == kNoErr) {			
					result = sAIPluginGroup->GetPluginArtEditArt(source, &child);				
				}				
				if (result == kNoErr) {				
					result = DuplicateCoveredArt(child, destination, bounds, start, found);					
				}			
			}
			break;			
		default:		
			break;		
	}
	if (source == start) {	
		*found = true;	
	}
	
	return result;	
}

/*
*/
ASErr LiveDropShadowGroup::StartLensObject(AIArtHandle *object, ObjectData *data) 
{
	AIErr result = kNoErr;
	*object = NULL;
	
	return result;	
}

/*
*/
ASErr LiveDropShadowGroup::FinishLensObject( AIArtHandle *object, ObjectData *data ) 
{
	AIErr result = kNoErr;		
	*object = NULL;
	
	return result;	
}

/*
*/
ASErr LiveDropShadowGroup::appendPluginGroup(AIArtHandle *object, ObjectData *data,
											 AIArtHandle path, AIPluginGroupHandle& pluginGroupHandle) 
{
	AIErr result = kNoErr;
	
	if (*object == NULL) {	
		if (result == kNoErr) {		
			result = EditLensObject(data);		
		}
		// Create a new kPluginArt above the selected art.
		if (result == kNoErr) {		
			result = sAIArt->NewArt(kPluginArt, kPlaceAbove, path, object);		
		}		
		if (result == kNoErr) {
			// register your pluginGroup
			result = sAIPluginGroup->UseAIPluginGroup(*object, pluginGroupHandle);		
		}		
		if (result == kNoErr) {		
			result = sAIPluginGroup->SetPluginArtDataCount(*object, sizeof(ObjectData));		
		}		
		if (result == kNoErr) {		
			result = sAIPluginGroup->SetPluginArtDataRange(*object, data, 0, sizeof(ObjectData));		
		}		
	}
	// unlock the path put it inside the edit art object
	if (result == kNoErr) {	
		result = sAIArt->SetArtUserAttr(path, kArtLocked, kArtLocked);	
	}	
	if (result == kNoErr) {	
		result = sAIArt->SetArtUserAttr(path, kArtLocked, 0);	
	}
	AIArtHandle editArt = NULL;
	if (result == kNoErr) {	
		result = sAIPluginGroup->GetPluginArtEditArt(*object, &editArt);	
	}
	if (result == kNoErr) {
		// put path inside editArt
		result = sAIArt->ReorderArt(path, kPlaceInsideOnBottom, editArt);	
	}	
	return result;	
}

/*
*/
ASErr LiveDropShadowGroup::ReleaseLensObject(AIArtHandle plugin) 
{
	AIErr result = kNoErr;
	
	if (result == kNoErr) {	
		result = sAIArt->SetArtUserAttr(plugin, kArtLocked, kArtLocked);	
	}
	
	if (result == kNoErr) {	
		result = sAIArt->SetArtUserAttr(plugin, kArtLocked, 0);	
	}

	AIArtHandle editArt = NULL;
	if (result == kNoErr) {	
		result = sAIPluginGroup->GetPluginArtEditArt(plugin, &editArt);	
	}

	AIArtHandle child = NULL;
	if (result == kNoErr) {	
		result = sAIArt->GetArtFirstChild(editArt, &child);	
	}
	
	while ((result == kNoErr) && (child != NULL)) {	
		if (result == kNoErr) {	
			AIArtHandle clone = NULL;
			result = sAIArt->DuplicateArt(child, kPlaceAbove, plugin, &clone);	
		}
		
		if (result == kNoErr) {	
			result = sAIArt->GetArtSibling(child, &child);	
		}	
	}
	
	if (result == kNoErr) {	
		result = sAIArt->DisposeArt(plugin);	
	}
	
	return result;	
}

/*
*/
ASErr LiveDropShadowGroup::FirstLensObject(AIArtHandle *object, ObjectData *data,
										   AIPluginGroupHandle& pluginGroupHandle) 
{
	AIErr result = kNoErr;

	ai::int32 count = 0;
	AIArtHandle **store = NULL;
	*object = NULL;	
	if (result == kNoErr)
	{	
		if (sAIMatchingArt->GetSelectedArt(&store, &count))
		{		
			store = NULL;
			count = 0;		
		}	
	}

	if (store != NULL)
	{
		for (ai::int32 index = 0; (result == kNoErr) && (index < count) && (*object == NULL); index++)
		{		
			AIArtHandle art = (*store)[index];
			short artType = kUnknownArt;
			
			if (result == kNoErr)	
				result = sAIArt->GetArtType(art, &artType);

			if (artType == kPluginArt)
			{			
				AIPluginGroupHandle artEntry = NULL;
				
				if (result == kNoErr)		
					result = sAIPluginGroup->GetPluginArtPluginGroup(art, &artEntry);				
				
				if (result == kNoErr)
				{				
					if (artEntry != pluginGroupHandle)					
						artType = kUnknownArt;	
				}			
			}

			if (artType == kPluginArt)
			{
				size_t length = 0;
				
				if (result == kNoErr)		  	
					result = sAIPluginGroup->GetPluginArtDataCount(art, &length);				
				
				if (result == kNoErr)
				{				
					if (length != sizeof(ObjectData))				
						artType = kUnknownArt;	
				}			
			}

			if (artType == kPluginArt)
			{
				if (result == kNoErr)
					result = sAIPluginGroup->GetPluginArtDataRange(art, data, 0, sizeof(ObjectData));

				if (result == kNoErr)
					*object = art;
			}			
		}		
	}

	if (store != NULL)	
		sAIMdMemory->MdMemoryDisposeHandle((AIMdMemoryHandle) store);	
	
	if (result == kNoErr)
	{  	
  		if (*object == NULL)		
  			result = kUnknownArtTypeErr;
	}

	if (result == kNoErr)
		result = EditLensObject(data);

	return result;	
}

/*
*/
ASErr LiveDropShadowGroup::LastLensObject(AIArtHandle *object, ObjectData *data) 
{
	AIErr result = kNoErr;	
	*object = NULL;
	
	return result;	
}

/*
*/
ASErr LiveDropShadowGroup::ModifyLensObject(AIArtHandle *object, ObjectData *data, AIArtHandle plugin) 
{
	AIErr result = kNoErr;
	if (result == kNoErr) {  	
		result = sAIPluginGroup->SetPluginArtDataCount(plugin, sizeof(ObjectData));	
	}

	if (result == kNoErr) {	
		result = sAIPluginGroup->SetPluginArtDataRange(plugin, &data, 0, sizeof(ObjectData));	
	}

	if (result == kNoErr) {
		sAIUser->MessageAlert(ai::UnicodeString("The data stores size and range were set."));
	}
	else {
		sAIUser->ErrorAlert(ai::UnicodeString("An error occurred, the data stores size and range were not set."));
	}
	
	return result;	
}

/*
*/
ASErr LiveDropShadowGroup::EditLensObject(ObjectData *data) 
{
	return kNoErr;	
}

/*
*/
ASErr LiveDropShadowGroup::setDropShadowStyle(AIArtHandle object) 
{
	AIErr result = kNoErr;

	short type;	
	if (result == kNoErr) {	
		result = sAIArt->GetArtType(object, &type);	
	}

	AIPathStyle style;
	AIArtHandle child = NULL;
	switch (type) {	
		case kPathArt:			
			if (result == kNoErr) {			
				result = sAIPathStyle->GetPathStyle(object, &style, nullptr);
			}
			if (result == kNoErr) {
				result = setDropShadowFill(&style.fill);			
			}			
			if (result == kNoErr) {
				result = setDropShadowStroke(&style.stroke);			
			}			
			if (result == kNoErr) {			
				result = sAIPathStyle->SetPathStyle(object, &style);			
			}
			break;
		case kMysteryPathArt:
		case kRasterArt:		
			if (result == kNoErr) {			
				result = sAIPathStyle->GetPathStyle(object, &style, nullptr);
			}
			if (result == kNoErr) {
				result = setDropShadowFill(&style.fill);			
			}			
			if (result == kNoErr) {
				result = setDropShadowStroke( &style.stroke );			
			}			
			if (result == kNoErr) {			
				result = sAIPathStyle->SetPathStyle(object, &style);			
			}
			break;			
		case kCompoundPathArt:
			if (result == kNoErr) {			
				result = sAIArt->GetArtFirstChild(object, &child);			
			}			
			while ((result == kNoErr) && (child != NULL)) {			
				if (result == kNoErr) {				
					result = setDropShadowStyle(child);					
				}				
				if (result == kNoErr) {				
					result = sAIArt->GetArtSibling(child, &child);				
				}				
			}
			break;			
		case kPlacedArt:		
			break;		
		case kPluginArt:
			break;
		default:		
			break;		
	}
	
	return result;	
}

/*
*/
ASErr LiveDropShadowGroup::setDropShadowFill(AIFillStyle *fill) 
{
	AIErr result = kNoErr;	
	fill->color.kind = kGrayColor;
	fill->color.c.g.gray = kAIRealOne * 0.5;
	fill->overprint = false;
	return result;
}

/*
*/
ASErr LiveDropShadowGroup::setDropShadowStroke(AIStrokeStyle *stroke) 
{
	AIErr result = kNoErr;
	stroke->color.kind = kGrayColor;
	stroke->color.c.g.gray = kAIRealOne * 0.5;
	stroke->overprint = false;	
	return result;	
}

/*
*/
ASErr LiveDropShadowGroup::transformDropShadow(AIArtHandle art, AIReal tx, AIReal ty, AIReal theta,
											   AIReal sx, AIReal sy) 
{
	AIErr result = kNoErr;	
	
	short type = kUnknownArt;
	sAIArt->GetArtType(art, &type);
	
	AIArtHandle child = NULL;
	AIRealPoint artCenter;
	AIRealMatrix artMatrix;		
	ai::int32 transformFlags = kTransformObjects | kScaleLines;	
	AIRealRect artBounds;	
	AIReal lineScale = (sAIRealMath->AIRealSqrt(sx)) * (sAIRealMath->AIRealSqrt(sy));
	switch (type) {	
		case kGroupArt:
			result = sAIArt->GetArtFirstChild(art, &child);		
			while ((result == kNoErr) && (child != NULL)) {		
				if (result == kNoErr) {
					result = transformDropShadow(child, tx, ty, theta, sx, sy);
				}			
				if (result == kNoErr) {
					result = sAIArt->GetArtSibling(child, &child);
				}
			}
		break;	
		case kCompoundPathArt:
			result = sAIArt->GetArtFirstChild(art, &child);		
			while ((result == kNoErr) && (child != NULL)) {		
				if (result == kNoErr) {
					result = transformDropShadow(child, tx, ty, theta, sx, sy);
				}			
				if (result == kNoErr) {
					result = sAIArt->GetArtSibling(child, &child);
				}			
			}
		break;	
		default:
			result = sAIArt->GetArtBounds(art, &artBounds);
			artCenter.h = artBounds.left + (artBounds.right-artBounds.left) / 2;
			artCenter.v = artBounds.bottom + (artBounds.top-artBounds.bottom) / 2;
		
			// Move object to center.
			sAIRealMath->AIRealMatrixSetTranslate(&artMatrix, -artCenter.h, -artCenter.v);
			// Translate object.
			sAIRealMath->AIRealMatrixConcatTranslate(&artMatrix, tx, ty);
			// Rotate object.
			sAIRealMath->AIRealMatrixConcatRotate(&artMatrix, theta);
			// Scale object.
			sAIRealMath->AIRealMatrixConcatScale(&artMatrix, sx, sy);
			// Move the object back to the center.
			sAIRealMath->AIRealMatrixConcatTranslate(&artMatrix, artCenter.h, artCenter.v);
		
			result = sAITransformArt->TransformArt(art, &artMatrix, lineScale, transformFlags);
		break;
	}
	return result;
}
// End LiveDropShadowGroup.cpp
