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

#ifndef __SnpObjectActionHelper_h__
#define __SnpObjectActionHelper_h__

/** Helper class that shows how to program objects with actions.

	@see AIActionManagerSuite
  	@ingroup sdk_snippet
 */
class SnpObjectActionHelper
{
public:

		/** Value parameter block (VPB) - a container for action parameters.
			@see AIActionManagerSuite
		*/
		class VPB
		{
		public:
			/**	Constructor
			*/
			VPB();

			/**	Destructor
			*/
			~VPB();

			/**	Returns the value parameter block structure.
			*/
			inline operator AIActionParamValueRef(void) const {return fActionParamValueRef;}

			

		private:
			/** The value parameter block.
			*/
			AIActionParamValueRef fActionParamValueRef;
		}; // end class VPB

public:

		/** Gathers parameters and copies the selection.
			@return kNoErr on success, other error code otherwise.
			@see kAICopySelectionAction
		*/
		ASErr CopySelection();

		/** Copies selection to the clipboard.
			@param dialogStatus IN whether parameters dialog appears by default.
			@param parameters IN the value parameter block.
			@return kNoErr on success, other error code otherwise.
			@see kAICopySelectionAction
		*/
		ASErr CopySelection(ActionDialogStatus dialogStatus, AIActionParamValueRef parameters);

		/** Gathers parameters and pastes the contents of the clipboard to the document.
			@return kNoErr on success, other error code otherwise.
			@see kAIPasteClipboardAction
		*/
		ASErr PasteClipboard();

		/** Pastes clipboard to document.
			@param dialogStatus IN whether parameters dialog appears by default.
			@param parameters IN the value parameter block.
			@return kNoErr on success, other error code otherwise.
			@see kAIPasteClipboardAction
		*/
		ASErr PasteClipboard(ActionDialogStatus dialogStatus, AIActionParamValueRef parameters);

	};

#endif // __SnpObjectActionHelper_h__

// End SnpObjectActionHelper.h
