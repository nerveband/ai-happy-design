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

#ifndef  __SnpSelectionHelper_h__
#define  __SnpSelectionHelper_h__

/** Shows how to select and deselect objects.

	@see AIArtSuite, AIMatchingArtSuite, SnpArtSetHelper
  	@ingroup sdk_snippet
 */
class SnpSelectionHelper
{
	public:
		/** Checks if the user is currently working with a document.
			@return true if there is a current document, false otherwise.
			@see SnpDocumentHelper::GetCurrentDocument
		*/
		ASBoolean IsDocumentSelected(void);

		/** Checks if art objects are selected.
			@param artType IN specifies type of selected art to match (default kAnyArt)
			@return true if art objects are selected, false otherwise.
			@see SnpArtSetHelper
		*/
		ASBoolean IsArtSelected(AIArtType artType = kAnyArt);

		/** Checks if a range of text characters is selected.
			@return true if a range of text characters is selected, false otherwise.
		*/
		ASBoolean IsTextRangeSelected(void);

		/** Checks whether given art object can be selected - objects that are on layers 
			that are locked or hidden should not be selected.
			@param art IN object to select.
			@return true if given art object can be selected, false otherwise.
		*/
		ASBoolean CanSelectArt(AIArtHandle art);

		/** Select given art object.
			@param art IN object to select.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr SelectArt(AIArtHandle art);

		/** Deselects all art in the document the user is currently working with.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr DeselectAll();
};

#endif // __SnpSelectionHelper_h__

// End SnpSelectionHelper.h

