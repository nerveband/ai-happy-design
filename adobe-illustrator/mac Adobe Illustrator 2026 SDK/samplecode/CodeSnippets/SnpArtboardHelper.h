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

#ifndef __SnpArtboardHelper_h__
#define __SnpArtboardHelper_h__

/** Provides helper functions to create, edit, display and delete artboards in a document.
  	@ingroup sdk_snippet
	@see AIArtboardSuite, AIArtboardRangeSuite
 */
class SnpArtboardHelper
{
	public:

		/** Adds a new artboard in a document using input from user or default
			values.
			@return kNoErr on success, other ASErr otherwise.
			@see AIArtboardSuite::AddNew
		 */
		ASErr AddArtboard();

		/** Adds a new artboard to a document using parameters provided.
			@param ArtboardRect IN the left, top, right, bottom values for the artboard.
			@param par IN the pixel aspect ratio (used if units are pixels).
			@param showCenter IN true to show center mark, false to not.
			@param showCrossHairs IN true to show cross hairs, false to not.
			@param showSafeAreas IN true to show safe areas, false to not.
			@return kNoErr on success, other ASErr otherwise.
			@see AIArtboardSuite::AddNew
		 */
		ASErr AddArtboard(AIRealRect ArtboardRect, AIReal par, ASBoolean showCenter, ASBoolean showCrossHairs, ASBoolean showSafeAreas);

		/** Adds a new artboard using values provided in AIArtboard object.
			@param Artboard IN properties of new artboard in AIArtboard object.
			@return kNoErr on success, other ASErr otherwise.
			@see AIArtboardSuite::AddNew
		 */
		ASErr AddArtboard(ai::ArtboardProperties artboardProperties);

		/** Deletes an artboard if there is more than 1 in the document,
			gets artboard number from user or uses the default value.
			@return kNoErr on success, other ASErr otherwise.
			@see AIArtboardSuite::Delete
		 */
		ASErr DeleteArtboard();

		/** Deletes the artboard at the index provided.
			@param artboardIndex IN the index of the artboard to delete.
			@return kNoErr on success, other ASErr otherwise.
			@see AIArtboardSuite::Delete
		 */
		ASErr DeleteArtboard(ASInt32 ardboardIndex);

		/** Deletes a range of artboards specified by user or default range.
			@return kNoErr on success, other ASErr otherwise.
			@see AIArtboardRangeHandle, AIArtboardSuite::Delete
		 */
		ASErr DeleteArtboards();

		/** Deletes a range of artboards.
			@param artboardRange IN the range of artboards to delete.
			@return kNoErr on success, other ASErr otherwise.
			@see AIArtboardRangeHandle, AIArtboardSuite::Delete
		 */
		ASErr DeleteArtboards(AIArtboardRangeHandle ardboardRange);

		/** Edits an existing artboard, gets artboard index from the user or uses the default.
			@return kNoErr on success, other ASErr otherwise.
			@see AIArtboardSuite::Update
		*/
		ASErr EditArtboard();

		/** Edits an existing artboard using parameters provided.
			@param artboardIndex IN the index of the artboard to edit.
			@param ArtboardRect IN the left, top, right, bottom values for the artboard.
			@param par IN the pixel aspect ratio of the artboard (used if units are pixels).
			@param showCenter IN true to show center mark false to not.
			@param showCrossHairs IN true to show cross hairs, false to not.
			@param showSafeAreas IN true to show safe areas, false to not.
			@return kNoErr on success, other ASErr otherwise.
			@see AIArtboardSuite::Update
		 */
		ASErr EditArtboard(ASInt32 artboardIndex, AIRealRect ArtboardRect, AIReal par, ASBoolean showCenter, ASBoolean showCrossHairs, ASBoolean showSafeAreas);

		/** Edits an exiting artboard using values in AIArtboardPtr.
			@param artboardIndex IN the index of the artboard to edit.
			@param ArtboardPtr IN the AIArtboardPtr object containing the new values for the artboard.
			@return kNoErr on success, other ASErr otherwise.
			@see AIArtboardSuite::Update
		*/
		ASErr EditArtboard(ASInt32 artboardIndex, ai::ArtboardProperties artboardProperties);

		/** Displays the details of an artboard in the log, uses user specified artboard number or default.
			@return kNoErr on success, other ASErr otherwise.
			@see AIArtboardSuite::Get
		 */
		ASErr DisplayArtboard();

		/** Displays the details of the artboard at the specified index, in the log.
			@param artboardIndex IN the index of the artboard to display.
			@return kNoErr on success, other ASErr otherwise.
			@see AIArtboardSuite::Get
		 */
		ASErr DisplayArtboard(ASInt32 artboardIndex);

		/** Displays the details of the active artboard in the log.
			@return kNoErr on success, other ASErr otherwise.
			@see AIArtboardSuite::GetActive
		 */
		ASErr DisplayActiveArtboard();

		/** Displays the details of all artboards in the log.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr DisplayAllArtboards();

		/** Sets the active artboard using index specified by user or default.
			@return kNoErr on success, other ASErr otherwise.
			@see AIArtboardSuite::SetActive
		 */
		ASErr SetActiveArtboard();

		/** Sets the active artboard using the index parameter provided.
			@return kNoErr on success, other ASErr otherwise.
			@see AIArtboardSuite::SetActive
		 */
		ASErr SetActiveArtboard(ASInt32 artboardIndex);

		/** Gets the number of artboards in the current document.
			@return the number of artboards in the current document.
			@see AIArtboardSuite::GetCount
		 */
		ASInt32 CountArtboards();

		/** rename the active artboard in the current document by user or default.
			@return the number of artboards in the current document.
			@see AIArtboardSuite::GetCount
		 */
		ASErr RenameActiveArtboard();

		
		/** rename the active artboard in the current document using the name parameter provided.
			@return the number of artboards in the current document.
			@see AIArtboardSuite::GetCount
		 */
		ASErr RenameActiveArtboard(const ai::UnicodeString& name);
		

};

#endif // __SnpArtboardHelper_h__

// End SnpArtboardHelper.h
