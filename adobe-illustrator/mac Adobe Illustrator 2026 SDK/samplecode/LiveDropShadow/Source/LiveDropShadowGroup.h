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

#ifndef __LiveDropShadowGroup_h__
#define __LiveDropShadowGroup_h__

class LiveDropShadowGroup
{
public:
	/**
	*/
	typedef struct{
		void init() {
			nothing = 0;
		}
		long nothing;	
	} ObjectData;

	/** Constructor
	*/
	LiveDropShadowGroup(){}

	/** Destructor
	*/
	virtual ~LiveDropShadowGroup(){}

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr PluginGroupUpdate(AIPluginGroupMessage* message);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr StartLensObject(AIArtHandle* object, ObjectData* data);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr FinishLensObject(AIArtHandle* object, ObjectData* data);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr appendPluginGroup(AIArtHandle* object, ObjectData* data, AIArtHandle path, AIPluginGroupHandle& pluginGroupHandle);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr ReleaseLensObject(AIArtHandle plugin);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr FirstLensObject(AIArtHandle* object, ObjectData* data, AIPluginGroupHandle& pluginGroupHandle);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr LastLensObject(AIArtHandle* object, ObjectData* data);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr ModifyLensObject(AIArtHandle* object, ObjectData* data, AIArtHandle plugin);

private:

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr EditLensObject(ObjectData* data);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr SetFilteredFill(AIBoolean* paint, AIFillStyle* fill);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr SetFilteredStroke(AIBoolean* paint, AIStrokeStyle* stroke);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr SetFilteredStyle(AIArtHandle object);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr SetBackingFill(AIBoolean* paint, AIFillStyle* fill);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr SetBackingStroke(AIBoolean* paint, AIStrokeStyle* stroke);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr SetBackingStyle(AIArtHandle object);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr SetBoundingFill(AIBoolean* paint, AIFillStyle* fill);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr SetBoundingStroke(AIBoolean* paint, AIStrokeStyle* stroke);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr SetBoundingStyle(AIArtHandle object);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr SetClippingStyle(AIArtHandle object);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr DuplicateCoveredArt(AIArtHandle source, AIArtHandle destination, AIRealRect* bounds, AIArtHandle start, AIBoolean* found);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr ReportError(AIErr result);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr setDropShadowStyle(AIArtHandle object);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr setDropShadowFill(AIFillStyle* fill);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr setDropShadowStroke(AIStrokeStyle* stroke);

	/** 
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr transformDropShadow(AIArtHandle art, AIReal tx, AIReal ty, AIReal theta, AIReal sx, AIReal sy);
};

#endif // End LiveDropShadowGroup.h
