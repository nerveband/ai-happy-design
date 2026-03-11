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

#ifndef __SnpArtHelper_h__
#define __SnpArtHelper_h__

/** Provides helper functions to create new art items in a document.
  	@ingroup sdk_snippet
	@see AIArtSuite, AIShapeConstructionSuite, AIPathSuite, AIPathSegment
 */
class SnpArtHelper
{
	public:

		/** Adds a new rectangle with default parameters to a document.
			@param artHandle IN/OUT reference to the art item.
			@return kNoErr on success, other error code otherwise.
			@see AIShapeConstructionSuite::NewRect
		*/
		ASErr NewRectangle(AIArtHandle& artHandle);

		/** Adds a new rectangle with specified parameters to a document.
			@param rect IN top, left, bottom, right coordinates of new rectangle.
			@param reversed IN true to reverse the object.
			@param artHandle IN/OUT reference to the art item.
			@return kNoErr on success, other error code otherwise.
			@see AIShapeConstructionSuite::NewRect
		*/
		ASErr NewRectangle(const AIRealRect rect, const AIBoolean reversed, AIArtHandle& artHandle);

		/** Adds a new rounded rectangle with default parameters to a document.
			@param artHandle IN/OUT reference to the art item.
			@return kNoErr on success, other error code otherwise.
			@see AIShapeConstructionSuite::NewRoundedRect
		*/
		ASErr NewRoundedRectangle(AIArtHandle& artHandle);

		/** Adds a new rounded rectangle with specified parameters to a document.
			@param rect IN top, left, bottom, right coordinates of new rounded rectangle.
			@param horizRadius IN horizontal radius of rounded corners.
			@param vertRadius IN vertical radius of rounded corners.
			@param reversed IN true to reverse the object.
			@param artHandle IN/OUT reference to the art item.
			@return kNoErr on success, other error code otherwise.
			@see AIShapeConstructionSuite::NewRoundedRect
		*/
		ASErr NewRoundedRectangle(const AIRealRect rect, const AIReal horizRadius, const AIReal vertRadius, const AIBoolean reversed, AIArtHandle& artHandle);

		/** Adds a new oval circumscribed around a rectangle with default parameters to a document.
			@param artHandle IN/OUT reference to the art item.
			@return kNoErr on success, other error code otherwise.
			@see AIShapeConstructionSuite::NewCircumscribedOval
		*/
		ASErr NewOval(AIArtHandle& artHandle);

		/** Adds a new oval circumscribed around a rectangle with specified parameters to a document.
			@param rect IN top, left, bottom, right coordinates of new oval.
			@param reversed IN true to reverse the object.
			@param artHandle IN/OUT reference to the art item.
			@return kNoErr on success, other error code otherwise.
			@see AIShapeConstructionSuite::NewCircumscribedOval
		*/
		ASErr NewOval(const AIRealRect rect, const AIBoolean reversed, AIArtHandle& artHandle);

		/** Adds a new polygon with default parameters to a document.
			@param artHandle IN/OUT reference to the art item.
			@return kNoErr on success, other error code otherwise.
			@see AIShapeConstructionSuite::NewRegularPolygon
		*/
		ASErr NewPolygon(AIArtHandle& artHandle);

		/** Adds a new polygon with specified parameters to a document.
			@param numSides IN number of sides.
			@param centerX IN center x coordinate.
			@param centerY IN center y coordinate.
			@param radius IN radius of polygon.
			@param reversed IN true to reverse the object.
			@param artHandle IN/OUT reference to the art item.
			@return kNoErr on success, other error code otherwise.
			@see AIShapeConstructionSuite::NewRegularPolygon
		*/
		ASErr NewPolygon(const short numSides, const AIReal centerX, const AIReal centerY, const AIReal radius, const AIBoolean reversed, AIArtHandle& artHandle);

		/** Adds a new star with default parameters to a document.
			@param artHandle IN/OUT reference to the art item.
			@return kNoErr on success, other error code otherwise.
			@see AIShapeConstructionSuite::NewStar
		*/
		ASErr NewStar(AIArtHandle& artHandle);

		/** Adds a new star with specified parameters to a document.
			@param numPoints IN number of points.
			@param centerX IN center x coordinate.
			@param centerY IN center y coordinate.
			@param innerRadius IN inner radius of the star.
			@param outerRadius IN outer radius of the star.
			@param reversed IN true to reverse the object.
			@param artHandle IN/OUT reference to the art item.
			@return kNoErr on success, other error code otherwise.
			@see AIShapeConstructionSuite::NewStar
		*/
		ASErr NewStar(const short numPoints, const AIReal centerX, const AIReal centerY, const AIReal innerRadius, const AIReal outerRadius, const AIBoolean reversed, AIArtHandle& artHandle);

		/** Adds a new spiral with default parameters to a document.
			@param artHandle IN/OUT reference to the art item.
			@return kNoErr on success, other error code otherwise.
			@see AIShapeConstructionSuite::NewStar
		*/
		ASErr NewSpiral(AIArtHandle& artHandle);

		/** Adds a new spiral with specified parameters to a document.
			@param firstArcCenter IN center point of first arc.
			@param start IN start point.
			@param decayPercent IN percent object should decay.
			@param numQuarterTurns IN number of quarter turns.
			@param clockwiseFromOutside IN true to spiral clockwise.
			@param artHandle IN/OUT reference to the art item.
			@return kNoErr on success, other error code otherwise.
			@see AIShapeConstructionSuite::NewStar
		*/
		ASErr NewSpiral(const AIRealPoint firstArcCenter, const AIRealPoint start, const AIReal decayPercent, const short numQuarterTurns, const AIBoolean clockwiseFromOutside, AIArtHandle& artHandle);

		/** Adds a new arc shaped path with default parameters to a document.
			@param artHandle IN/OUT reference to the art item.
			@return kNoErr on success, other error code otherwise.
			@see AIArtSuite::NewArt, AIPathSuite, AIPathSegment
		*/
		ASErr NewArc(AIArtHandle& artHandle);

};

#endif // __SnpArtHelper_h__

// End SnpArtHelper.h
