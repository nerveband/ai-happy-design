//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/CodeSnippets/SnpSymbolHelper.h $
//
//  $Revision: #1 $
//
//  Copyright 2010 Adobe Systems Incorporated. All rights reserved.
//  
//  NOTICE:  Adobe permits you to use, modify, and distribute this file in accordance 
//  with the terms of the Adobe license agreement accompanying it.  If you have received
//  this file from a source other than Adobe, then your use, modification, or 
//  distribution of it requires the prior written permission of Adobe.
//  
//========================================================================================

#ifndef __SnpSymbolHelper_h__
#define __SnpSymbolHelper_h__

/** Provides helper functions to create, edit, display and delete symbol.
  	@ingroup sdk_snippet
	@see AISymbolSuite
 */
class SnpSymbolHelper
{
	public:

		/** Gets the number of symbols in the current document.
			@return the number of symbols in the current document.
			@see AISymbolSuite::CountSymbolPatterns
		 */
		ai::int32 CountSymbols();

		/** Creates a new listed symbol definition from selected art
			@return kNoErr on success, other ASErr otherwise.
			@see sAISymbolSuite::NewSymbolPatternFromSel
		 */
		ASErr AddSymbol();

		/** Deletes a symbol definition,gets symbol name from user or uses the default value.
			@return kNoErr on success, other ASErr otherwise.
			@see AISymbolSuite::DeleteSymbolPattern
		 */

		ASErr DeleteSymbol();
	
		/** rename the symbol in the current document by user or default.
			@return kNoErr on success, other ASErr otherwise.
			@see AISymbolSuite::SetSymbolPatternName
		 */
		ASErr RenameSymbol();

		/** Adds a new symbol to current document.
			@return kNoErr on success, other ASErr otherwise.
			@see AISymbolSuite::NewInstanceCenteredInView
		 */
		ASErr PlaceSymbolInDocument();

		/** Displays the details of a symbol in the log, uses user specified symbol number or default.
			@return kNoErr on success, other ASErr otherwise.
			@see AISymbolSuite::Get
		 */
		ASErr DisplaySymbolProperties();

		/** Displays the details of the symbol at the specified index, in the log.
			@param symbol IN the index of the symbol to display.
			@return kNoErr on success, other ASErr otherwise.
			@see AISymbolSuite::Get
		 */
		ASErr DisplaySymbolProperties(ai::int32 symbolIndex);

		/** Retrieves the symbol scaling grid for 9-slice scaling for a symbol pattern in the current document by user or default.
			@return kNoErr on success, other ASErr otherwise.
			@see AISymbolSuite::GetSymbol9SliceGrid
		 */
		ASErr Display9SliceScaledInfo();

		/** Retrieves the symbol scaling grid for 9-slice scaling for a symbol patternat the specified index, in the log.
			@param symbol IN the index of the symbol to display.
			@return kNoErr on success, other ASErr otherwise.
			@see AISymbolSuite::GetSymbol9SliceGrid
		 */
		ASErr Display9SliceScaledInfo(ai::int32 symbolIndex);

private:
		
		ASErr ChooseSymbol(ai::UnicodeString *name, ai::int32 *symbolIndex);
};

#endif // __SnpSymbolHelper_h__

// End SnpSymbolHelper.h
