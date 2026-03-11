//========================================================================================
//  
//  $File$
//
//  $Revision$
//
//  Copyright 2008 Adobe Systems Incorporated. All rights reserved.
//  
//  NOTICE:  Adobe permits you to use, modify, and distribute this file in accordance 
//  with the terms of the Adobe license agreement accompanying it.  If you have received
//  this file from a source other than Adobe, then your use, modification, or 
//  distribution of it requires the prior written permission of Adobe.
//  
//========================================================================================

#ifndef __ANNOTATOR_H__
#define __ANNOTATOR_H__

#include "AIAnnotator.h"
#include "SDKErrors.h"
#include "AnnotatorSuites.h"

/** This class provides a very basic demonstration of the use of the 
	AIAnnotatorSuite and AIAnnotatorDrawerSuite to add your own annotations 
	to artwork.

	@see AIAnnotatorSuite, AIAnnotatorDrawerSuite, AIHitTestSuite, AIDocumentViewSuite
*/
class Annotator
{
private:
	/** Handle for art object the cursor was last over.
	*/
	AIArtHandle fArtLastHit;

	/** Defines whether the art annotation is drawn during next
		annotator draw call.
	*/
	ASBoolean fArtUpdate;

	/** Defines whether the cursor annotation is drawn during next
		annotator draw call.
	*/
	ASBoolean fCursorUpdate;

	/** Bounds of the art bounds annotation.
	*/
	AIRect fArtRect;

	/** Bounds of the cursor location annotation.
	*/
	AIRect fCursorRect;

	/** The art bounds as artwork coordinates.
	*/
	AIRealRect fArtArtworkBounds;

	/** The cursor location as an artwork coordinate.
	*/
	AIRealPoint fCursorArtPoint;

	/** The cursor location as a view coordinate.
	*/
	AIPoint fCursorViewPoint;

public:

	/**	Constructor.
		@param pluginRef IN reference to this plugin.
	*/
	Annotator();

	/**	Destructor.
	*/
	virtual ~Annotator(){}

	/** Initialise both coordinates in an AIPoint struct to 0.
		@param point IN the point to initialise.
	*/
	void Init(AIPoint& point);

	/** Initialise both coordinates in an AIRealPoint struct to 0.
		@param point IN the point to initialise.
	*/
	void Init(AIRealPoint& point);

	/** Initialise the 4 values in an AIRect struct to 0.
		@param rect IN the rectangle to initialise.
	*/
	void Init(AIRect& rect);

	/** Initialise the 4 values in an AIRealRect struct to 0.
		@param rect IN the rectangle to initialise.
	*/
	void Init(AIRealRect& rect);

	/** Checks if cursor is over art and updates art selection, if 
		necessary, to reflect a change in cursor location.
		@param message IN message data.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr CheckForArtHit(AIToolMessage* message);

	/**	Draws the art bounds annotation.
		@param message IN message data.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr DrawArtAnnotation(AIAnnotatorMessage* message);

	/**	Draws the cursor location annotation.
		@param message IN message data.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr DrawCursorAnnotation(AIAnnotatorMessage* message);

	/**	Invalidates an area on the artboard so affected annotations will be redrawn.
		@param invalRealRect IN area to invalidate as an AIRealRect.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr InvalidateRect(const AIRealRect& invalRealRect);

	/**	Invalidates an area on the artboard so affected annotations will be redrawn
		@param invalRectIN area to invalidate as an AIRect.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr InvalidateRect(const AIRect& invalRect);

	/**	Converts an AIRealRect in artwork coordinates to an AIRect in view coordinates.
		@param artworkBounds IN the rectangle bounds in artwork coordinates.
		@param viewBounds OUT the rectangle bounds in view coordinates.
		@return kNoErr on success, other ASErr otherwise. This API considers rotation
	*/
	ASErr ArtworkBoundsToViewBounds(const AIRealRect& artworkBounds, AIRect& viewBounds);

	/**	Converts an AIRealRect in artwork coordinates to an AIRect in view coordinates.
		@param artworkBounds IN the rectangle bounds in artwork coordinates.
		@param viewBounds OUT the rectangle bounds in view coordinates.
		@return kNoErr on success, other ASErr otherwise. This API ignores rotation
	*/
	ASErr ArtworkBoundsToViewBoundsUnrotated(const AIRealRect& artworkBounds, AIRect& viewBounds);

	/**	Converts the bounds provided in an AIRealRect to unicode strings.
		@param rect IN the rectangle containing the bounds to convert.
		@param topString OUT the unicode string representing the top boundary of the rectangle.
		@param leftString OUT the unicode string representing the left boundary of the rectangle.
		@param bottomString OUT the unicode string representing the bottom boundary of the rectangle.
		@param rightString OUT the unicode string representing the right boundary of the rectangle.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr GetRectString(const AIRealRect& rect, ai::UnicodeString& topString, ai::UnicodeString& leftString, ai::UnicodeString& bottomString, ai::UnicodeString& rightString);

	/**	Converts a point coordinate to a unicode string.
		@param point IN the point to convert.
		@param pointStr OUT the unicode string representing the point coordinate.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr GetPointString(const AIRealPoint& point, ai::UnicodeString& pointStr);
};

#endif // __ANNOTATOR_H__

// End Annotator.h
