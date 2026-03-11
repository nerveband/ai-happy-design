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

#include "IllustratorSDK.h"
#include "Annotator.h"

/** A light yellow color that we will use as the background for our annotations
*/
const AIRGBColor BKGND_COLOR = {65000, 65000, 40000};

/** Black, the color we will use for the borders of our annotations.
*/
const AIRGBColor BORDER_COLOR = {0, 0, 0};

/*
*/
Annotator::Annotator(): fArtLastHit(NULL), fArtUpdate(false), fCursorUpdate(false)
{
	this->Init(fCursorViewPoint);
	this->Init(fCursorArtPoint);
	this->Init(fArtRect);
	this->Init(fCursorRect);
	this->Init(fArtArtworkBounds);
}

/*
*/
void Annotator::Init(AIPoint& point)
{
	point.h = 0, point.v = 0;
}

/*
*/
void Annotator::Init(AIRealPoint& point)
{
	point.h = 0, point.v = 0;
}

/*
*/
void Annotator::Init(AIRect& rect)
{
	rect.left = 0; rect.top = 0; rect.right = 0; rect.bottom = 0;
}

/*
*/
void Annotator::Init(AIRealRect& rect)
{
	rect.left = 0; rect.top = 0; rect.right = 0; rect.bottom = 0;
}

/*
*/
ASErr Annotator::CheckForArtHit(AIToolMessage* message)
{
	ASErr result = kNoErr;
	try {
		// Set cursor location.
		this->fCursorArtPoint = message->cursor;
		SDK_ASSERT(sAIDocumentView);
		result = sAIDocumentView->ArtworkPointToViewPoint(NULL, &fCursorArtPoint, &fCursorViewPoint);
		aisdk::check_ai_error(result);

		// Check if cursor is over any art.
		AIHitRef hitRef = NULL;
		AIToolHitData toolHitData;
		SDK_ASSERT(sAIHitTest);
		result = sAIHitTest->HitTest(NULL, &fCursorArtPoint, kAllHitRequest, &hitRef);
		aisdk::check_ai_error(result);
		result = sAIHitTest->GetHitData(hitRef, &toolHitData);
		aisdk::check_ai_error(result);
		result = sAIHitTest->Release(hitRef);
		aisdk::check_ai_error(result);
		
		if (toolHitData.hit && toolHitData.object != NULL && toolHitData.object != fArtLastHit) {
			// New art item hit, set last art object hit in class to art item hit.
			fArtLastHit = toolHitData.object;
			fArtUpdate = true;
			fCursorUpdate = true;

			// Get the art bounds in artwork coordinates.
			SDK_ASSERT(sAIArt);
			result = sAIArt->GetArtTransformBounds(fArtLastHit, NULL, kNoStrokeBounds, &fArtArtworkBounds);
			aisdk::check_ai_error(result);

			AIPoint viewTopRight{ 0, 0 };
			AIRealPoint artworkTopRight{ this->fArtArtworkBounds.right, this->fArtArtworkBounds.top };
			SDK_ASSERT(sAIDocumentView);
			result = sAIDocumentView->ArtworkPointToViewPoint(NULL, &artworkTopRight, &viewTopRight);
			aisdk::check_ai_error(result);

			// Set the annotation location to the right of the art.
			fArtRect.left = viewTopRight.h + 10;
			fArtRect.top = viewTopRight.v;
			fArtRect.right = fArtRect.left + 80;
			fArtRect.bottom = fArtRect.top + 70;
			
			// If other art is selected, deselect all art.
			SDK_ASSERT(sAIMatchingArt);
			if (sAIMatchingArt->IsSomeArtSelected) {
				result = sAIMatchingArt->DeselectAll();
				aisdk::check_ai_error(result); 
			}

			// Select art hit.
			result = sAIArt->SetArtUserAttr(fArtLastHit, kArtSelected, kArtSelected);
			aisdk::check_ai_error(result);
		}
		else if (fArtLastHit != NULL && toolHitData.object == fArtLastHit) {
			fCursorUpdate = true;
			result = this->InvalidateRect(fCursorRect);
			aisdk::check_ai_error(result);
		}
		else {
			// No art hit, deselect all art.
			SDK_ASSERT(sAIMatchingArt);
			result = sAIMatchingArt->DeselectAll();
			aisdk::check_ai_error(result);
			fArtLastHit = NULL;
			fArtUpdate = false;
			fCursorUpdate = false;
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr Annotator::DrawArtAnnotation(AIAnnotatorMessage* message)
{
	ASErr result = kNoErr;
	try {
		if (fArtUpdate) {

			// Draw a light yellow filled rectangle
			SDK_ASSERT(sAIAnnotatorDrawer);
			sAIAnnotatorDrawer->SetColor(message->drawer, BKGND_COLOR);
			result = sAIAnnotatorDrawer->DrawRect(message->drawer, fArtRect, true);
			aisdk::check_ai_error(result);

			// Draw a black border
			sAIAnnotatorDrawer->SetColor(message->drawer, BORDER_COLOR);
			sAIAnnotatorDrawer->SetLineWidth(message->drawer, 0.5);
			result = sAIAnnotatorDrawer->DrawRect(message->drawer, fArtRect, false);
			aisdk::check_ai_error(result);

			// Draw art bounds text content.
			result = sAIAnnotatorDrawer->SetFontPreset(message->drawer, kAIAFSmall);
			aisdk::check_ai_error(result);

			ai::UnicodeString topString, leftString, bottomString, rightString;
			result = this->GetRectString(this->fArtArtworkBounds, topString, leftString, bottomString, rightString);
			aisdk::check_ai_error(result);
			
			AIReal textPadding = sAIAnnotatorDrawer->GetFontSize(message->drawer) + 5;
			AIPoint textBottomLeft;
			textBottomLeft.h = fArtRect.left + 5;
			textBottomLeft.v = fArtRect.top + textPadding;

			// Add art top bounds text.
			result = sAIAnnotatorDrawer->DrawText(message->drawer, topString, textBottomLeft, false);
			aisdk::check_ai_error(result);

			// Add art left bounds text.
			textBottomLeft.v = textBottomLeft.v + textPadding;
			result = sAIAnnotatorDrawer->DrawText(message->drawer, leftString, textBottomLeft, false);
			aisdk::check_ai_error(result);

			// Add art bottom bounds text.
			textBottomLeft.v = textBottomLeft.v + textPadding;
			result = sAIAnnotatorDrawer->DrawText(message->drawer, bottomString, textBottomLeft, false);
			aisdk::check_ai_error(result);

			// Add art right bounds text.
			textBottomLeft.v = textBottomLeft.v + textPadding;
			result = sAIAnnotatorDrawer->DrawText(message->drawer, rightString, textBottomLeft, false);
			aisdk::check_ai_error(result);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr Annotator::DrawCursorAnnotation(AIAnnotatorMessage* message)
{
	ASErr result = kNoErr;
	try {	
		if (fCursorUpdate) {
			// Get cursor location as string.
			ai::UnicodeString pointStr;
			result = this->GetPointString(fCursorArtPoint, pointStr);
			aisdk::check_ai_error(result);

			// Get a point a little above and to the right of the cursor location
			ai::int32 offset = 5;
			AIPoint offsetCursorPoint;
			offsetCursorPoint.h = fCursorViewPoint.h + offset;
			offsetCursorPoint.v = fCursorViewPoint.v - offset;

			// Get a rectangle near the cursor representing the dimensions of our text
			AIRect textBounds;
			result = sAIAnnotatorDrawer->GetTextBounds(message->drawer, pointStr, &offsetCursorPoint, false, textBounds, false);
			aisdk::check_ai_error(result);

			// Now construct the rectangle we wish to draw, adding space above and to the right for padding
			fCursorRect = textBounds;
			fCursorRect.top -= 5;
			fCursorRect.right += 20;
			
			// Draw a filled rectangle with a light yellow background color
			sAIAnnotatorDrawer->SetColor(message->drawer, BKGND_COLOR);
			result = sAIAnnotatorDrawer->DrawRect(message->drawer, fCursorRect, true);
			aisdk::check_ai_error(result);

			// Draw a black border around the rectangle
			sAIAnnotatorDrawer->SetColor(message->drawer, BORDER_COLOR);
			sAIAnnotatorDrawer->SetLineWidth(message->drawer, 0.5);
			result = sAIAnnotatorDrawer->DrawRect(message->drawer, fCursorRect, false);
			aisdk::check_ai_error(result);

			// Add the text to the rectangle
			result = sAIAnnotatorDrawer->SetFontPreset(message->drawer, kAIAFSmall);
			aisdk::check_ai_error(result);
			result = sAIAnnotatorDrawer->DrawTextAligned(message->drawer, pointStr, kAICenter, kAIMiddle, fCursorRect, false);
			aisdk::check_ai_error(result);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr Annotator::InvalidateRect(const AIRealRect& invalRealRect)
{
	ASErr result = kNoErr;
	try {		
		// invalRealRect is in artwork coordinates, convert to view 
		// coordinates for AIAnnotatorSuite::InvalAnnotationRect.
		AIRect invalRect;
		result = this->ArtworkBoundsToViewBounds(invalRealRect, invalRect);
		aisdk::check_ai_error(result);

		// Invalidate the rect bounds so it is redrawn.
		SDK_ASSERT(sAIAnnotator);
		result = sAIAnnotator->InvalAnnotationRect(NULL, &invalRect);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr Annotator::InvalidateRect(const AIRect& invalRect)
{
	ASErr result = kNoErr;
	try {		
		// Invalidate the rect bounds so it is redrawn.
		SDK_ASSERT(sAIAnnotator);
		result = sAIAnnotator->InvalAnnotationRect(NULL, &invalRect);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr Annotator::ArtworkBoundsToViewBounds(const AIRealRect& artworkBounds, AIRect& viewBounds)
{
	ASErr result = kNoErr;
	try {		
		SDK_ASSERT(sAIDocumentView);
		result = sAIDocumentView->ArtworkRectToViewRect(NULL, &artworkBounds, &viewBounds);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

ASErr Annotator::ArtworkBoundsToViewBoundsUnrotated(const AIRealRect& artworkBounds, AIRect& viewBounds)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIDocumentView);
		result = sAIDocumentView->ArtworkRectToViewRectUnrotated(NULL, &artworkBounds, &viewBounds);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr Annotator::GetRectString(const AIRealRect& rect, ai::UnicodeString& topString,
							   ai::UnicodeString& leftString, ai::UnicodeString& bottomString,
							   ai::UnicodeString& rightString)
{
	ASErr result = kNoErr;
	try
	{		
		ai::int32 precision = 2;
		ai::UnicodeString top, left, bottom, right;
		ai::NumberFormat numFormat;
		top = numFormat.toString(-rect.top, precision, top, true);
		left = numFormat.toString(rect.left, precision, left, true);
		bottom = numFormat.toString(-rect.bottom, precision, bottom, true);
		right = numFormat.toString(rect.right, precision, right, true);
		topString = ai::UnicodeString("Top: ").append(top);
		leftString = ai::UnicodeString("Left: ").append(left);
		bottomString = ai::UnicodeString("Bottom: ").append(bottom);
		rightString = ai::UnicodeString("Right: ").append(right);
	}
	catch (ai::Error& ex)
	{
		result = ex;
	}

	return result;
}

/*
*/
ASErr Annotator::GetPointString(const AIRealPoint& point, ai::UnicodeString& pointStr)
{
	ASErr result = kNoErr;
	try
	{		
		ai::int32 precision = 2;
		ai::UnicodeString horiz, vert;
		ai::NumberFormat numFormat;
		horiz = numFormat.toString(point.h, precision, horiz);
		vert = numFormat.toString(-point.v, precision, vert);
		pointStr = pointStr.append(ai::UnicodeString("x: ").append(horiz).append(ai::UnicodeString(", y: ").append(vert)));
	}
	catch (ai::Error& ex)
	{
		result = ex;
	}

	return result;
}

