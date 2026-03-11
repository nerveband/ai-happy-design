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

#include "SDKErrors.h"

// Framework includes:
#include "SnpRunnable.h"
#include "SnippetRunnerSuites.h"
#include "SnippetRunnerLog.h"
#include "SnippetRunnerParameter.h"
#include "SnpSelectionHelper.h"
#include "SnpArtSetHelper.h"
#include "SnippetRunnerManager.h"
#include "IAIRef.h"
#include "SnpLiveTrace.h"

/**	Set the AITracingModeValue, can be 
		kAITracingModeColor
		kAITracingModeGray
		kAITracingModeBlackAndWhite.
*/
const VectorizeModes kSnpLiveColorTracingMode = kAIVectorizeModeColor;

/** Set the tracing option for the blur value to use.
*/
const AIReal kSnpLiveColorBlur = 0;

/** Set the tracing option key for maximum number of colors.
*/
const ai::int32 kSnpLiveColorTracingMaxColors = 30;

/** Set the tracing option key for the resample resolution.
*/
const AIReal kSnpLiveColorTracingResampleResolution = 300;

/** Set the name of the swatch group created by 
	CreateTintedSpotColors function.
*/
const char* kCreateTintedSpotColorsSGName = "SnpLiveColorCreateTintedSpotColors";

/** Set the name of the swatch group created by 
	GroupSpotColors function.
*/
const char* kGroupSpotColorsSGName = "SnpLiveColorGroupSpotColors";

/** Set the name of the swatch group created by 
	ViewInColorHarmonyDialog function.
*/
const char* kViewInColorHarmonyDialogSGName = "SnpLiveColorViewInColorHarmonyDialog";

/** Demonstrates how to convert colors in a selection to black and white or spot 
	colors, and manipulate color groups in the Swatches palette, Color Guide and the 
	Color Harmony modal dialog.

	@see AIColorHarmonySuite, AIPathStyleSuite, AIMatchingArtSuite, AIColorConversionSuite,
	AISwatchListSuite, AISwatchGroupSuite, AICustomColorSuite
	@ingroup sdk_snippet
 */
class SnpLiveColor
{
	private:
		/**	Structure for associating an AIColor with an int to keep 
			track of the number of times each color is found in a selection.
		*/
        struct CountedColor {
			/**	The color object contained within the CountedColor.
			*/
			AIColor fColor;
			
			/**	The number of instances of fColor found in the selected art.
			*/
			ASInt32 fColorCount;

			/** Constructor.
			*/
			CountedColor(AIColor color, ASInt32 colorCount) {
				fColor = color;
				fColorCount = colorCount;
			}
		};

		/** Typedef for a container of CountedColor objects.
		*/
		typedef std::vector<CountedColor> CountedColors;

		/** Structure used to pass user defined data to AdjustColorCallback().
		*/
        struct CallbackData{
			/**	Container of CountedColor objects.
			*/
			CountedColors fCountedColors;

			/**	Store for total number of colors in the selection.
			*/
			ASInt32 fTotalColors;

			/**	Store for index of current color in the selection.
			*/
			ASInt32 fCurrentColor;
			
			/** Constructor.
			*/
			CallbackData(ASInt32 totalColors = 0, ASInt32 currentColor = 0) {
				fTotalColors = totalColors;
				fCurrentColor = currentColor;
			}
		};

	public:

		/**	Creates a new spot color with a 50% tint for each color found in the selected art,
			using AIPathStyleSuite::AdjustObjectAIColors to visit each color, the colors used
			in the selected art are not modified. The new spot swatches are created in a new 
			swatch group called SnpLiveColorCreateTintedSpotColors in the swatches palette.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr CreateTintedSpotColors();
		
		/** Finds all spot colors in the current document swatch list and groups them together
			in a new swatch group called SnpLiveColorGroupSpotColors in the swatches palette. 
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr GroupSpotColors();
		
		/** Converts the selected art to black and white - all light colors become white and 
			all dark colors become black.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr ConvertToBlackAndWhite();

		/** Finds all colors used in the selected art, creates a new swatch group containing all the
			unique colors found called SnpLiveColorViewInColorHarmonyDialog in the swatches palette 
			and opens the color group in the Color Harmony dialog.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr ViewInColorHarmonyDialog();

		/** Finds the most frequently occuring color in the selected art then opens the Color
			Guide with the selection's most frequent color set as the base color. The base color is 
			the first color displayed in the harmony rule display of the color guide and the color
			within the large bold circle in the color wheel, the base color is the only color attribute 
			of the color wheel that can be set through the public api.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr SetColorGuideBaseColor();
		
		//------------------------------ Callback Functions -------------------------------
		/*  The following are the callback functions for AdjustObjectAIColors*/

		/** Callback function of type AIAdjustColorFunc for CreateTintedSpotColors, for each
			color found in the selected art initiates the process to convert the color to a
			new spot color swatch with a 50% tint.
			@param color IN/OUT the color to modify and return.
			@param userData IN pointer to a CallbackData structure.
			@param reportResult OUT returns an error if encountered.
			@param modifiedColor OUT returns true if the color was modified, false if not.
		*/
		static void CreateTintedSpotAdjustColorFunc(AIColor* color, void* userData, AIErr* reportResult, ASBoolean* modifiedColor);

		/** Callback function of type AIAdjustColorFunc for ConvertToBlackAndWhite, for each 
			color found	in the selected art initiates the process to convert the color to either
			black or white depending on its grayscale value.
			@param color IN/OUT the color to modify and return.
			@param userData IN pointer to a CallbackData structure.
			@param reportResult OUT returns an error if encountered.
			@param modifiedColor OUT returns true if the color was modified, false if not.
		*/
		static void ConvertToBlackAndWhiteAdjustColorFunc(AIColor* color, void* userData, AIErr* reportResult, ASBoolean* modifiedColor);

		/** Callback function of type AIAdjustColorFunc for ViewInColorHarmonyDialog, for each 
			color found	in the selected art initiates the process of adding the color to a new 
			swatch group which will be opened in the Color Harmony dialog at the end.
			@param color IN/OUT the color to modify and return.
			@param userData IN pointer to a CallbackData structure.
			@param reportResult OUT returns an error if encountered.
			@param modifiedColor OUT returns true if the color was modified, false if not.
		*/
		static void ViewInColorHarmonyAdjustColorFunc(AIColor* color, void* userData, AIErr* reportResult, ASBoolean* modifiedColor);

		/** Callback function of type AIAdjustColorFunc for SetColorGuideBaseColor, for each 
			color found	in the selected art initiates the process of parsing the color in order
			to count the number of times the color occurs in the art.
			@param color IN/OUT the color to modify and return.
			@param userData IN pointer to a CallbackData structure.
			@param reportResult OUT returns an error if encountered.
			@param modifiedColor OUT returns true if the color was modified, false if not.
		*/
		static void SetColorGuideBaseColorAdjustColorFunc(AIColor* color, void* userData, AIErr* reportResult, ASBoolean* modifiedColor);

		/** Callback function that counts the number of colors in the selection so a progress bar can be implemented.
			@param color IN/OUT the color to modify and return.
			@param userData IN pointer to a CallbackData structure.
			@param reportResult OUT returns an error if encountered.
			@param modifiedColor OUT returns true if the color was modified, false if not.
		*/
		static void CountColorFunc(AIColor* color, void* userData, AIErr* reportResult, ASBoolean* modifiedColor);

		//------------------------------- Utility Functions -------------------------------


		/**	Changes newColor to a spot color using its original color values and a unique
			custom color name.
			@param newColor IN/OUT object in which to return the new spot color.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr ConvertToSpot(AIColor& newColor);

		/**	Converts color to its corresponding grayscale value, then converts the gray value
			to either black or white depending on the depth of gray - < 50% becomes white, >= 50%
			becomes black.
			@param color IN/OUT object in which to return the new gray color.
			@param modified OUT ASBoolean returns true if color was modified, false if not.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr ConvertToBlackAndWhite(AIColor& color, ASBoolean& modifiedColor);

		/**	Creates a new swatch group if it does not already exist or returns the swatch group if it 
			does, if deleteIfExists is true, deletes the swatch group if it exists then creates and 
			returns the new swatch group.
			@param swatchGroupName IN the name of the swatch group to look for, create and return.
			@param swatchGroup OUT reference to the new/existing swatch group.
			@param deleteIfExists IN if true deletes the swatch group if found then creates a new swatch
			group with the same name, if false returns the existing swatch group if found.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr CreateSwatchGroup(const ai::UnicodeString& swatchGroupName, AISwatchGroupRef& swatchGroup, ASBoolean deleteIfExists);

		/**	Finds the swatch group with the name matching groupName and adds newColor to the swatch
			group if it is not already in the group.
			@param newColor IN object to be added to the swatch group.
			@param groupName IN the name of the swatch group the new color should be added to.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr AddColorToGroup(AIColor& newColor, const ai::UnicodeString& groupName);

		/** Iterates through the colors in the swatch group to check whether the current
			color has already been added to the group.
			@param newColor IN object to look for in the group.
			@param swatchGroupRef IN reference to the swatch group to search.
			@param exists OUT true if color is found false if not.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr ColorInGroup(const AIColor& newColor, const AISwatchGroupRef& swatchGroupRef, ASBoolean& exists);

		/** Compares two AICustomColors and reports whether they are equal.
			@param newColor IN first custom color.
			@param existingColor IN second custom color to compare with first.
			@return ASBoolean true if the custom colors are equal, false if not
		*/
		ASBoolean Equal(const AICustomColorHandle& newColor, const AICustomColorHandle& existingColor) const;

		
		/** Compares two AIColors and reports whether they are equal.
			@param newColor IN first color.
			@param existingColor IN second color to compare with first.
			@return ASBoolean true if the colors are equal, false if not
		*/
		ASBoolean Equal(const AIColor& newColor,const AIColor& existingColor) const;

		/** Uses SnpArtSetHelper to find all the raster and placed items in the artset and
			converts all found to tracings. The vector result art of each tracing is acquired 
			and the tracings disposed of using this->GetResultArt().
			@param selectedArt IN/OUT art set in which to pass and return the art set.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr ConvertToVectors(SnpArtSetHelper& selectedArt);
		
		/** Iterates through each CountedColor in the vector comparing its AIColor object with 
			newColor, adds a new CountedColor, containing newColor, to the vector if the color 
			is not found, if the color is found increments the matching CountedColors' colorCount
			attribute.
			@param newColor IN object in which to pass in the new color.
			@param colors IN/OUT vector of CountedColors to search and return.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr ParseColor(AIColor& newColor,std::vector<CountedColor>& colors);

		/** Searches through the vector of CountedColors, looking for the most frequently
			occuring color.
			@param colors IN vector of CountedColors to search.
			@param commonColor IN/OUT object in which to return the most common color.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr GetMostCommonColor(std::vector<CountedColor>& colors, AIColor& commonColor);
};


/*
*/
ASErr SnpLiveColor::CreateTintedSpotColors()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerLog::Indent indent;
		
		// Find all selected art.
		const short numSpecs = 1;
		AIArtSpec artSpec[numSpecs] = { 
			{kAnyArt, kArtSelected, kArtSelected},
		};
		SnpArtSetHelper artSet(artSpec, numSpecs);
		
		if (artSet.GetCount() > 0) {
			AISwatchGroupRef swatchGroup = NULL;
			result = this->CreateSwatchGroup(ai::UnicodeString(kCreateTintedSpotColorsSGName), swatchGroup, true);
			aisdk::check_ai_error(result);

			// Convert raster, traced and placed items into vector art items 
			// so AdjustObjectAIColors can access the colors.		
			result = this->ConvertToVectors(artSet);
			aisdk::check_ai_error(result);

			// Visit all colors in the selection which meet visitFlags.
			VisitAIColorFlags visitFlags = (kVisitColorsSolidOnly | kVisitGlobalObjectsOnceOnly | kVisitColorsSelectedOnly);
			ASBoolean modified = false;
			// Data to pass to the callback function.
			CallbackData data = CallbackData();
			// First count the number of colors in the selection.
			result = sAIPathStyle->AdjustObjectAIColors(NULL, this->CountColorFunc, &data, visitFlags, &modified);
			// Do not call aisdk::check_ai_error(result) as there are values returned by AdjustObjectAIColors
			// which are not errors, such as kIterationCanQuitReply which means the process completed successfully.
			if (result == kNoErr || result == kIterationCanQuitReply) {
				result = sAIPathStyle->AdjustObjectAIColors(NULL, this->CreateTintedSpotAdjustColorFunc, &data, visitFlags, &modified);
				if (result == kIterationCanQuitReply)
					result = kNoErr;
			}
			aisdk::check_ai_error(result);

			// Release any tracing from selected traced art to return art to
			// original state. TODO: a kUntouchableArtObjectErr is returned if i try to dispose
			// of the result art here so for the moment it remains in the document
			// on top of the original art item.
			SnpLiveTrace trc;
			SnpArtSetHelper tracingArtSet = trc.FindIsTracingArtSet(true);
			result = trc.ReleaseTracing(tracingArtSet);
			aisdk::check_ai_error(result);
		}
		else {
			log->Write(ai::UnicodeString("No art selected"));
			aisdk::check_ai_error(kBadParameterErr);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveColor::GroupSpotColors()
{
	ASErr result = kNoErr;
	try {
		int swatchCount = sAISwatchList->CountSwatches(NULL);
		// Iterate through document swatch list, looking for spot color swatches.
		for (int i = swatchCount-1; i >= 0; i--) {
			// Get the swatches color object.
			AISwatchRef swatchRef = sAISwatchList->GetNthSwatch(NULL, i);
			AIColor swatchColor;
			result = sAISwatchList->GetAIColor(swatchRef,&swatchColor);
			aisdk::check_ai_error(result);			
			if (swatchColor.kind == kCustomColor) {

				// Color is a custom color, get its AICustomColor object.
				AICustomColorHandle colorHandle = swatchColor.c.c.color;
				SDK_ASSERT(colorHandle);
				AICustomColor customColor;
				result = sAICustomColor->GetCustomColor(colorHandle, &customColor);
				aisdk::check_ai_error(result);
				if (customColor.flag == kCustomSpotColor) {
					// Get swatch's current location.
					AISwatchGroupRef currentGroup = NULL;
					ai::int32 groupIndex = 0, swatchIndex = 0;
					result = sAISwatchGroup->LocateSwatch(NULL, swatchRef, &currentGroup, &groupIndex, &swatchIndex);
					aisdk::check_ai_error(result);

					// Get the swatch group to contain the spot colors, create if necessary.
					AISwatchGroupRef swatchGroup = NULL;
					result = this->CreateSwatchGroup(ai::UnicodeString(kGroupSpotColorsSGName), swatchGroup, false);
					aisdk::check_ai_error(result);

					// Check whether swatch is already in the swatch group.
					if (currentGroup != swatchGroup) {
						// Move swatch to its new location in the new swatch group.
						result = sAISwatchGroup->ReorderSwatchBetweenGroups(currentGroup, swatchIndex, swatchGroup, -1);
						aisdk::check_ai_error(result);
						// Have to increment i here, as the reorder moves another swatch into the current
						// position, if i is not incremented to nullify the decrement in the for loop the 
						// swatch will be missed.
						i++;
					}
				}
			}
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveColor::ConvertToBlackAndWhite()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerLog::Indent indent;

		// Find all selected art.
		const short numSpecs = 1;
		AIArtSpec artSpec[numSpecs] = { 
			{kAnyArt, kArtSelected, kArtSelected},
		};
		SnpArtSetHelper artSet(artSpec, numSpecs);
		
		if (artSet.GetCount() > 0) {
			// Convert raster, traced and placed items into vector art items which AdjustObjectAIColors
			// can access.
			result = this->ConvertToVectors(artSet);
			aisdk::check_ai_error(result);

			// Visit all colors in the selection which meet visitFlags.
			VisitAIColorFlags visitFlags = (kVisitColorsSolidOnly | kVisitGlobalObjectsOnceOnly | kVisitColorsSelectedOnly);
			ASBoolean modified = false;
			// Data to pass to the callback function.
			CallbackData data = CallbackData();
			// First count the number of colors in the selection.
			result = sAIPathStyle->AdjustObjectAIColors(NULL, this->CountColorFunc, &data, visitFlags, &modified);
			// Do not call aisdk::check_ai_error(result) as there are values returned by AdjustObjectAIColors
			// which are not errors, such as kIterationCanQuitReply which means the process completed successfully.
			if (result == kNoErr || result == kIterationCanQuitReply) {
				if (result == kIterationCanQuitReply)
					result = kNoErr;
				result = sAIPathStyle->AdjustObjectAIColors(NULL, this->ConvertToBlackAndWhiteAdjustColorFunc, &data, visitFlags, &modified);
			}
			else {
				aisdk::check_ai_error(result);
			}

			// Release any tracing from selected traced art to return art to
			// original state. TODO: a kUntouchableArtObjectErr is returned if i try to dispose
			// of the result art here so for the moment it remains in the document
			// on top of the original art item.
			SnpLiveTrace trc;
			SnpArtSetHelper tracingArtSet = trc.FindIsTracingArtSet(true);
			result = trc.ReleaseTracing(tracingArtSet);
			aisdk::check_ai_error(result);
		}
		else {
			log->Write(ai::UnicodeString("No art selected"));
			aisdk::check_ai_error(kBadParameterErr);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveColor::ViewInColorHarmonyDialog()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerLog::Indent indent;

		// Find all selected art.
		const short numSpecs = 1;
		AIArtSpec artSpec[numSpecs] = { 
			{kAnyArt, kArtSelected, kArtSelected},
		};
		SnpArtSetHelper artSet(artSpec, numSpecs);

		if (artSet.GetCount() > 0) {		
			AISwatchGroupRef swatchGroup = NULL;
			result = this->CreateSwatchGroup(ai::UnicodeString(kViewInColorHarmonyDialogSGName), swatchGroup, true);
			aisdk::check_ai_error(result);

			// Convert raster, traced and placed items into vector art items so AdjustObjectAIColors
			// can access the colors.
			result = this->ConvertToVectors(artSet);
			aisdk::check_ai_error(result);

			// Visit all colors in the selection which meet visitFlags.
			VisitAIColorFlags visitFlags = (kVisitColorsSolidOnly | kVisitGlobalObjectsOnceOnly | kVisitColorsSelectedOnly);
			ASBoolean modified = false;
			// Data to pass to the callback function.
			CallbackData data = CallbackData();
			// First count the number of colors in the selection.
			result = sAIPathStyle->AdjustObjectAIColors(NULL, this->CountColorFunc, &data, visitFlags, &modified);
			// Do not call aisdk::check_ai_error(result) as there are values returned by AdjustObjectAIColors
			// which are not errors, such as kIterationCanQuitReply which means the process completed successfully.
			if (result == kNoErr || result == kIterationCanQuitReply) {
				if (result == kIterationCanQuitReply)
					result = kNoErr;
				result = sAIPathStyle->AdjustObjectAIColors(NULL, this->ViewInColorHarmonyAdjustColorFunc, &data, visitFlags, &modified);
			}
			else {
				aisdk::check_ai_error(result);
			}

			// Release any tracing from selected traced art to return art to
			// original state. TODO: a kUntouchableArtObjectErr is returned if i try to dispose
			// of the result art here so for the moment it remains in the document
			// on top of the original art item.
			SnpLiveTrace trc;
			SnpArtSetHelper tracingArtSet = trc.FindIsTracingArtSet(true);
			result = trc.ReleaseTracing(tracingArtSet);
			aisdk::check_ai_error(result);

			// Check whether running in unit test mode.
			SnippetRunnerManager& snpManager = SnippetRunnerManager::Instance();
			if (!snpManager.IsUnitTestRunning()) {
				// Not running unit test so open the Color Harmony modal dialog with 
				// the new swatch group loaded.
				result = sAIColorHarmony->ShowWithColorGroup(swatchGroup, false);
				aisdk::check_ai_error(result);
			}
		}
		else {
			log->Write(ai::UnicodeString("No art selected"));
			aisdk::check_ai_error(kBadParameterErr);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveColor::SetColorGuideBaseColor()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerLog::Indent indent;

		// Find all selected art.
		const short numSpecs = 1;
		AIArtSpec artSpec[numSpecs] = { 
			{kAnyArt, kArtSelected, kArtSelected},
		};
		SnpArtSetHelper artSet(artSpec, numSpecs);
		
		if (artSet.GetCount() > 0 ) {
			// Convert raster, traced and placed items into vector art items so AdjustObjectAIColors
			// can access the colors.
			result = this->ConvertToVectors(artSet);
			aisdk::check_ai_error(result);
			
			// Visit all colors in the selection which meet visitFlags.
			VisitAIColorFlags visitFlags = (kVisitColorsSolidOnly | kVisitGlobalObjectsOnceOnly | kVisitColorsSelectedOnly);
			ASBoolean modified = false;
			// Data to pass to the callback function.
			CallbackData data = CallbackData();
			// First count the number of colors in the selection.
			result = sAIPathStyle->AdjustObjectAIColors(NULL, this->CountColorFunc, &data, visitFlags, &modified);
			// Do not call aisdk::check_ai_error(result) as there are values returned by AdjustObjectAIColors
			// which are not errors, such as kIterationCanQuitReply which means the process completed successfully.
			if (result == kNoErr || result == kIterationCanQuitReply) {
				if (result == kIterationCanQuitReply)
					result = kNoErr;
				result = sAIPathStyle->AdjustObjectAIColors(NULL, this->SetColorGuideBaseColorAdjustColorFunc, &data, visitFlags, &modified);
			}
			else {
				aisdk::check_ai_error(result);
			}

			// Release any tracing from selected traced art to return art to
			// original state. TODO: a kUntouchableArtObjectErr is returned if i try to dispose
			// of the result art here so for the moment it remains in the document
			// on top of the original art item.
			SnpLiveTrace trc;
			SnpArtSetHelper tracingArtSet = trc.FindIsTracingArtSet(true);
			result = trc.ReleaseTracing(tracingArtSet);
			aisdk::check_ai_error(result);
			
			AIColor commonColor;
			result = this->GetMostCommonColor(data.fCountedColors, commonColor);
			aisdk::check_ai_error(result);

			// Now set the base color to the most common color and show the color guide panel.
			result = sAIColorHarmony->SetColorGuideBaseColor(&commonColor);
			aisdk::check_ai_error(result);

			result = sAIColorHarmony->ShowColorGuide();
			aisdk::check_ai_error(result);
		}
		else {
			log->Write(ai::UnicodeString("No art selected"));
			aisdk::check_ai_error(kBadParameterErr);
		}

	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}


// ---------------------------------------- Callback Functions -----------------------------------------------------

/*
*/
void SnpLiveColor::CreateTintedSpotAdjustColorFunc(AIColor* color, void* userData, AIErr* reportResult, ASBoolean* modifiedColor)
{
	ASErr result = kNoErr;
	try {
		SnpLiveColor slc;		
		*reportResult = kNoErr;
		*modifiedColor = false;
		CallbackData *data = (CallbackData*)userData;

		data->fCurrentColor++;

		// Update the progress bar.
		sAIUser->SetProgressText(ai::UnicodeString("Creating tinted spot colors"));
		sAIUser->UpdateProgress(data->fCurrentColor, data->fTotalColors);

		AIColor newColor = *color;	// Don't modify original color.
		result = slc.ConvertToSpot(newColor);
		aisdk::check_ai_error(result);

		// Add the new spot color to the swatch group.
		result = slc.AddColorToGroup(newColor, ai::UnicodeString(kCreateTintedSpotColorsSGName));
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		*reportResult = ex;
	}
}
/*
*/
void SnpLiveColor::ConvertToBlackAndWhiteAdjustColorFunc(AIColor* color, void* userData, AIErr* reportResult, ASBoolean* modifiedColor)
{
	ASErr result = kNoErr;
	try {
		SnpLiveColor slc;
		*reportResult = kNoErr;	
		*modifiedColor = false;

		CallbackData *data = (CallbackData*)userData;

		data->fCurrentColor++;

		// Update the progress bar.
		sAIUser->SetProgressText(ai::UnicodeString("Converting to black and white"));
		sAIUser->UpdateProgress(data->fCurrentColor, data->fTotalColors);

		result = slc.ConvertToBlackAndWhite(*color, *modifiedColor);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		*reportResult = ex;
	}
}
/*
*/
void SnpLiveColor::ViewInColorHarmonyAdjustColorFunc(AIColor* color, void* userData, AIErr* reportResult, ASBoolean* modifiedColor)
{
	ASErr result = kNoErr;
	try {
		SnpLiveColor slc;
		*reportResult = kNoErr;	
		*modifiedColor = false;

		CallbackData *data = (CallbackData*)userData;

		data->fCurrentColor++;

		// Update the progress bar.
		sAIUser->SetProgressText(ai::UnicodeString("Collecting colors in selection"));
		sAIUser->UpdateProgress(data->fCurrentColor, data->fTotalColors);

		AIColor newColor = *color; // Don't modify original color.
		result = slc.AddColorToGroup(newColor, ai::UnicodeString(kViewInColorHarmonyDialogSGName));
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		*reportResult = ex;
	}
}
/*
*/
void SnpLiveColor::SetColorGuideBaseColorAdjustColorFunc(AIColor* color, void* userData, AIErr* reportResult, ASBoolean* modifiedColor)
{
	ASErr result = kNoErr;
	try {
		SnpLiveColor slc;
		*reportResult = kNoErr;	
		*modifiedColor = false;

		CallbackData *data = (CallbackData*)userData;

		data->fCurrentColor++;

		// Update the progress bar.
		sAIUser->SetProgressText(ai::UnicodeString("Finding most common color"));
		sAIUser->UpdateProgress(data->fCurrentColor, data->fTotalColors);

		AIColor newColor = *color; // Don't modify original color.
		result = slc.ParseColor(newColor, data->fCountedColors);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		*reportResult = ex;
	}
}

/*
*/
void SnpLiveColor::CountColorFunc(AIColor* color, void* userData, AIErr* reportResult, ASBoolean* modifiedColor)
{
	CallbackData *data = (CallbackData*)userData;
	*reportResult = kNoErr;	
	*modifiedColor = false;

	data->fTotalColors++;
}


//------------------------------------------------ Utility Functions ------------------------------------------------

/*
*/
ASErr SnpLiveColor::ConvertToSpot(AIColor& newColor)
{
	ASErr result = kNoErr;
	try {
		// Create the new spot color using the values from newColor.
		AICustomColor customColor;
		customColor.flag = kCustomSpotColor;
		short kind = newColor.kind;
		switch (kind) {
			default:
				break;
			case kGrayColor: {
				// kCustomFourColor
				customColor.kind = kCustomFourColor;
				customColor.c.f.cyan = 0.0;
				customColor.c.f.magenta = 0.0;
				customColor.c.f.yellow = 0.0;
				customColor.c.f.black = newColor.c.g.gray;
				break;
			}
			case kThreeColor: {
				// kCustomThreeColor
				customColor.kind = kCustomThreeColor;
				customColor.c.rgb.red = newColor.c.rgb.red;
				customColor.c.rgb.green = newColor.c.rgb.green;
				customColor.c.rgb.blue = newColor.c.rgb.blue;
				break;
			}
			case kFourColor: {
				// kCustomFourColor
				customColor.kind = kCustomFourColor;
				customColor.c.f.cyan = newColor.c.f.cyan;
				customColor.c.f.magenta = newColor.c.f.magenta;
				customColor.c.f.yellow = newColor.c.f.yellow;
				customColor.c.f.black = newColor.c.f.black;
				break;
			}
		}
		// Get a new unique name for the spot color.
		ai::UnicodeString name;
		result = sAICustomColor->NewCustomColorName(name);
		aisdk::check_ai_error(result);

		// Create the new custom color handle using the spot color and name set above.
		AICustomColorHandle customColorHandle = NULL;
		result = sAICustomColor->NewCustomColor(&customColor, name, &customColorHandle);
		aisdk::check_ai_error(result);

		// Set new colors kind to kCustomColor and its color to the spot color.
		newColor.kind = kCustomColor;
		newColor.c.c.color = customColorHandle;
		// Set the tint to 50%.
		newColor.c.c.tint = 0.5;
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveColor::ConvertToBlackAndWhite(AIColor& color, ASBoolean& modifiedColor)
{
	ASErr result = kNoErr;
	modifiedColor = false;
	try {
		AIColorConversionSpaceValue toSpace = kAIGrayColorSpace;
		AIColorConversionSpaceValue srcSpace = kAICMYKColorSpace;
		
		AIFloatSampleComponent src[4];
		AIFloatSampleComponent dst[4];
		
		if (color.kind == kCustomColor) {
			/* Convert color from a custom color to:
				- If a global process color returns the process equivalent.
				- If a spot color returns the process approximation of the 
					appearance of the spot color.
				- If a registration color returns the process color used to
					represent the registration color.
			*/
			result = sAICustomColor->ConvertToNonGlobal(color.c.c.color, color.c.c.tint, &color);
			aisdk::check_ai_error(result);
		}
		// Get the individual color values.
		switch (color.kind) {
			case kFourColor: {
				src[0] = color.c.f.cyan;
				src[1] = color.c.f.magenta;
				src[2] = color.c.f.yellow;
				src[3] = color.c.f.black;
				break;				 
			}
			case kThreeColor: {
				srcSpace = kAIRGBColorSpace;
				src[0] = color.c.rgb.red;
				src[1] = color.c.rgb.green;
				src[2] = color.c.rgb.blue;
				break;				 
			}
			case kGrayColor: {
				// Already a grayscale color, so only need to convert to black or white.
				break;				 
			}
			default: {
				aisdk::check_ai_error(kColorConversionErr);
			}
		}
		// If RGB or CMYK convert to Grayscale.
		if (color.kind == kFourColor || color.kind == kThreeColor) {
			result = sAIColorConversion->ConvertSampleColor(srcSpace, &src[0], toSpace, &dst[0],
				AIColorConvertOptions::kForExport, NULL);
			aisdk::check_ai_error(result);
		
			color.kind = kGrayColor;
			color.c.g.gray = ((AIReal)1.0) - dst[0];
		}

		// Convert the gray color to either black or white.
		if (color.c.g.gray >= 0.0 && color.c.g.gray < 0.5) {
			// Less than 50% Gray, convert to white.
			color.c.g.gray = 0.0;
			modifiedColor = true;
		}
		else if (color.c.g.gray >= 0.5 && color.c.g.gray <= 1.0) {
			// More than 50% Gray convert to black.
			color.c.g.gray = 1.0;
			modifiedColor = true;
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveColor::CreateSwatchGroup(const ai::UnicodeString& swatchGroupName, AISwatchGroupRef& swatchGroup, ASBoolean deleteIfExists)
{
	ASErr result = kNoErr;
	try {
		// String swatchGroupName contains the name to be given to the new swatch group, check if a swatch group 
		// with that name already exists. If the swatch group does exist and deleteIfExists is true delete the 
		// swatch group and create and return a new empty swatch group with the name contained in swatchGroupName. 
		// If the swatch group exists and deleteIfExists is false, return the existing swatch group. If the swatch 
		// group does not exists create and return a new empty swatch group with the name contained in swatchGroupName.

		AISwatchGroupRef swatchGroupRef = sAISwatchGroup->GetSwatchGroupByName(NULL, swatchGroupName);
		if (swatchGroupRef && deleteIfExists) {
			// Delete swatches one by one so any global colors will also be removed. If just the
			// swatch group is removed, custom colors will be moved to the general swatch group
			// even if deleteSwatches is true, we don't want that safety catch here.
			int swatchCount = sAISwatchGroup->CountSwatches(swatchGroupRef);
			for (int i = 0; i < swatchCount; i++) {
				// Delete the swatch, even if a global color.
				result = sAISwatchGroup->RemoveNthSwatch(swatchGroupRef, 0, true, false);
				aisdk::check_ai_error(result);
			}
			// Now remove the swatch group.
			result = sAISwatchGroup->RemoveSwatchGroup(NULL, swatchGroupRef, true);
			aisdk::check_ai_error(result);
		}
		else if (swatchGroupRef && !deleteIfExists) {
			swatchGroup = swatchGroupRef;
			return result;
		}
		result = sAISwatchGroup->NewSwatchGroup(NULL, kAISGKindGeneric, -1, &swatchGroup);
		aisdk::check_ai_error(result);
		result = sAISwatchGroup->SetSwatchGroupName(swatchGroup, swatchGroupName);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveColor::AddColorToGroup(AIColor& newColor, const ai::UnicodeString& groupName)
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog::Indent indent;
		// Find the swatch group.
		AISwatchGroupRef swatchGroupRef = sAISwatchGroup->GetSwatchGroupByName(NULL, groupName);
		SDK_ASSERT(swatchGroupRef);
		// Check if color is already in the group.
		ASBoolean colorInGroup = false;
		result = this->ColorInGroup(newColor, swatchGroupRef, colorInGroup);
		aisdk::check_ai_error(result);
		if (!colorInGroup) {
			// Color not already in the group so add.
			AISwatchRef swatch = sAISwatchGroup->InsertNthSwatch(swatchGroupRef, &newColor, -1);
			SDK_ASSERT(swatch);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveColor::ColorInGroup(const AIColor& newColor, const AISwatchGroupRef& swatchGroupRef, ASBoolean& exists)
{
	ASErr result = kNoErr;
	exists = false;
	int swatchCount = sAISwatchGroup->CountSwatches(swatchGroupRef);
	// Iterate through all the colors already present in the swatch group.
	for (int i = 0; i < swatchCount; i++) {
		AISwatchRef swatchRef = sAISwatchGroup->GetNthSwatch(swatchGroupRef, i);
		SDK_ASSERT(swatchRef);
		AIColor existingColor;
		try {
			result = sAISwatchList->GetAIColor(swatchRef, &existingColor);
			aisdk::check_ai_error(result);
			
			if (newColor.kind == kCustomColor) {
				// Check if the custom color has already been added to the swatch group.
				exists = this->Equal(newColor.c.c.color, existingColor.c.c.color);
				if (exists) {
					return result;
				}
			}
			else {
				// Check if the color has already been added to the swatch group.
				exists = this->Equal(newColor,existingColor);
				if(exists) {
					return result;
				}
			}
		}
		catch (ai::Error& ex) {
			result = ex;
		}
	}
	return result;
}

/*
*/
ASBoolean SnpLiveColor::Equal(const AICustomColorHandle& newColorHandle, const AICustomColorHandle& existingColorHandle) const
{
	// Initial check that neither handles are NULL.
	SDK_ASSERT(newColorHandle);
	SDK_ASSERT(existingColorHandle);

	// Get the AICustomColors from the AICustomColorHandles.
	AICustomColor newCustomColor, existingCustomColor;
	ASErr result = sAICustomColor->GetCustomColor(newColorHandle, &newCustomColor);
	aisdk::check_ai_error(result);
	result = sAICustomColor->GetCustomColor(existingColorHandle, &existingCustomColor);
	aisdk::check_ai_error(result);

	if (newCustomColor.kind == existingCustomColor.kind) {
		switch (newCustomColor.kind) {
			case kCustomFourColor: {
				if (newCustomColor.c.f.cyan == existingCustomColor.c.f.cyan &&
					newCustomColor.c.f.magenta == existingCustomColor.c.f.magenta &&
					newCustomColor.c.f.yellow == existingCustomColor.c.f.yellow &&
					newCustomColor.c.f.black == existingCustomColor.c.f.black) {
						// Found a match.
						return true;
				}
				break;
			}
			case kCustomThreeColor: {
				if (newCustomColor.c.rgb.red == existingCustomColor.c.rgb.red &&
					newCustomColor.c.rgb.green == existingCustomColor.c.rgb.green &&
					newCustomColor.c.rgb.blue == existingCustomColor.c.rgb.blue) {
						// Found a match.
						return true;
				}
				break;
			}
		    case kCustomLabColor: {
				if (newCustomColor.c.lab.l == existingCustomColor.c.lab.l &&
					newCustomColor.c.lab.a == existingCustomColor.c.lab.a &&
					newCustomColor.c.lab.b == existingCustomColor.c.lab.b) {
						// Found a match.
						return true;
				}
				break;
			}
		}
	}
	// No match found.
	return false;
}

/*
*/
ASBoolean SnpLiveColor::Equal(const AIColor& newColor, const AIColor& existingColor) const
{

	if (newColor.kind == existingColor.kind) {
		switch(newColor.kind) {
			case kFourColor: {
				if (newColor.c.f.cyan == existingColor.c.f.cyan &&
					newColor.c.f.magenta == existingColor.c.f.magenta &&
					newColor.c.f.yellow == existingColor.c.f.yellow &&
					newColor.c.f.black == existingColor.c.f.black) {
						// Found a match.
						return true;
				}
				break;				 
			}
			case kThreeColor: {
				if (newColor.c.rgb.red == existingColor.c.rgb.red &&
					newColor.c.rgb.green == existingColor.c.rgb.green &&
					newColor.c.rgb.blue == existingColor.c.rgb.blue) {
						// Found a match.
						return true;
				}
				break;				 
			}
			case kGrayColor: {
				if(newColor.c.g.gray == existingColor.c.g.gray) {
					// Found a match.
					return true;
				}
				break;				 
			}
			case kCustomColor: {
				if (newColor.c.c.color == existingColor.c.c.color &&
					newColor.c.c.tint == existingColor.c.c.tint) {
						// Found a match.
						return true;
				}
				break;				 
			}
			case kPattern: {
				if (newColor.c.p.pattern == existingColor.c.p.pattern &&
					newColor.c.p.reflect == existingColor.c.p.reflect &&
					newColor.c.p.reflectAngle == existingColor.c.p.reflectAngle &&
					newColor.c.p.rotate == existingColor.c.p.rotate &&
					newColor.c.p.scale.h == existingColor.c.p.scale.h &&
					newColor.c.p.scale.v == existingColor.c.p.scale.v &&
					newColor.c.p.shearAngle == existingColor.c.p.shearAngle &&
					newColor.c.p.shearAxis == existingColor.c.p.shearAxis &&
					newColor.c.p.shiftAngle == existingColor.c.p.shiftAngle &&
					newColor.c.p.shiftDist == existingColor.c.p.shiftDist &&
					sAIRealMath->AIRealMatrixEqual((AIRealMatrix *)&newColor.c.p.transform,
					(AIRealMatrix *)&existingColor.c.p.transform)) {
						// Found a match.
						return true;
				}
				break;				 
			}
			case kGradient: {
				if (newColor.c.b.gradient == existingColor.c.b.gradient &&
					newColor.c.b.gradientAngle == existingColor.c.b.gradientAngle &&
					newColor.c.b.gradientLength == existingColor.c.b.gradientLength &&
					newColor.c.b.gradientOrigin.h == existingColor.c.b.gradientOrigin.h &&
					newColor.c.b.gradientOrigin.v == existingColor.c.b.gradientOrigin.v &&
					newColor.c.b.hiliteAngle == existingColor.c.b.hiliteAngle &&
					newColor.c.b.hiliteLength == existingColor.c.b.hiliteLength &&
					sAIRealMath->AIRealMatrixEqual((AIRealMatrix *)&newColor.c.b.matrix,
					(AIRealMatrix *)&existingColor.c.b.matrix)) {
						// Found a match.
						return true;
				}
				break;				 
			}
			case kNoneColor: {
				// Don't add if there is no color.
				return true;
				break;				 
			}
			default: {
				// Unknown kind, don't add.
				return true;
				break;		 
			}
		}
	}
	// Did not find a match.
	return false;
}

/*
*/
ASErr SnpLiveColor::ConvertToVectors(SnpArtSetHelper& artSet)
{
	ASErr result = kNoErr;
	try {
		for (size_t i = 0; i < artSet.GetCount(); i++) {
			AIArtHandle art = artSet[i];
			SDK_ASSERT(art);
			short type = kAnyArt;
			result = sAIArt->GetArtType(art, &type);
			aisdk::check_ai_error(result);
			if (type == kPlacedArt || type == kRasterArt) {
				AIArtHandle tracing = NULL;
				SnpLiveTrace trace;
				result = trace.CreateTracing(art, tracing, kSnpLiveColorTracingMode, kSnpLiveColorBlur,
					kSnpLiveColorTracingMaxColors, kSnpLiveColorTracingResampleResolution);
				aisdk::check_ai_error(result);

				// Redraw to update the selection.
				result = sAIDocument->RedrawDocument();
				aisdk::check_ai_error(result);

				result = trace.GetResultArt(tracing);
				aisdk::check_ai_error(result);

				artSet.Add(tracing);
			}
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveColor::ParseColor(AIColor& newColor, std::vector<CountedColor>& colors)
{
	ASErr result = kNoErr;
	try {
		if (colors.size() == 0) {
			// Color vector empty so add new color.
			CountedColor newCountedColor = CountedColor(newColor, 1);
			colors.push_back(newCountedColor);
		}
		else {
			// Iterate through contents checking whether the color has already been added.
			std::vector<CountedColor>::iterator iter;
			ASBoolean found = false;
			for (iter = colors.begin(); iter != colors.end(); iter++) {
				AIColor existingColor = iter->fColor;		
				// Check for equality of each objects color.
				if (newColor.kind == kCustomColor) {
					// Compare custom colors.
					found = this->Equal(newColor.c.c.color, existingColor.c.c.color);
					if (found)
						break;
				}
				else {
					// Compare colors.
					found = this->Equal(newColor, existingColor);		
					if (found)
						break;
				}
			}
			if (found) {
				// Colors are equal so increment color count.
				iter->fColorCount++;
			}
			else {
				// Colors not equal so add the new color to the vector.
				CountedColor newCountedColor = CountedColor(newColor,1);
				colors.push_back(newCountedColor);
			}
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveColor::GetMostCommonColor(std::vector<CountedColor>& colors, AIColor& commonColor)
{
	ASErr result = kNoErr;
	try {
		ASUInt32 highestCount = 0;
		std::vector<CountedColor>::iterator iter;
		for (iter = colors.begin(); iter != colors.end(); iter++) {
			ASUInt32 count = iter->fColorCount;
			if (count > highestCount) {
				highestCount = count;
				commonColor = iter->fColor;
			}
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpLiveColor available to the SnippetRunner framework.
 */
class _SnpRunnableLiveColor : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableLiveColor () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableLiveColor () {}

		/* Returns name of snippet.
		*/
		std::string GetName() const;

		/* Returns a description of what the snippet does.
		*/
		std::string	GetDescription() const;

		/* Returns the operations this snippet supports.
		*/
		Operations GetOperations() const;

		/* Returns name of this snippet's default operation.
		*/
		std::string GetDefaultOperationName() const;

		/** Returns the categories a snippet belongs to.
			@return default categories.
		*/
		std::vector<std::string> GetCategories() const;

		/* Returns true if the snippet can run.
			@param runnableContext see ISnpRunnableContext for documentation.
			@return true if snippet can run, false otherwise
		 */
		ASBoolean CanRun(SnpRunnable::Context& runnableContext);

		/* Runs the snippet.
			@param runnableContext see ISnpRunnableContext for documentation.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr Run(SnpRunnable::Context& runnableContext);

};

/*
*/
std::string _SnpRunnableLiveColor::GetName() const
{
	return "LiveColor";
}

/*
*/
std::string _SnpRunnableLiveColor::GetDescription() const 
{
	return "Performs Live Color functions on selected items.";
}

/*
*/
SnpRunnable::Operations _SnpRunnableLiveColor::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("CreateTintedSpotColors", "selected item(s)", kSnpRunEPSContext));
	operations.push_back(Operation("GroupSpotColors", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("ConvertToBlackAndWhite", "selected item(s)", kSnpRunPSDContext));
	operations.push_back(Operation("ViewInColorHarmonyDialog", "selected item(s)", kSnpRunPDFContext));
	operations.push_back(Operation("SetColorGuideBaseColor", "selected item(s)", kSnpRunJPGContext));
	return operations;
}

/*
*/
std::string _SnpRunnableLiveColor::GetDefaultOperationName() const 
{
	return "CreateTintedSpotColors";
}

/*
*/
std::vector<std::string> _SnpRunnableLiveColor::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Color Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableLiveColor::CanRun(SnpRunnable::Context& runnableContext)
{
	SnpSelectionHelper selectionHelper;
	if (!selectionHelper.IsDocumentSelected())
		return false;

	if ("CreateTintedSpotColors" == runnableContext.GetOperation().GetName() ||
		"ConvertToBlackAndWhite" == runnableContext.GetOperation().GetName() ||
		"ViewInColorHarmonyDialog" == runnableContext.GetOperation().GetName() ||
		"SetColorGuideBaseColor" == runnableContext.GetOperation().GetName()) {
		
			if (selectionHelper.IsArtSelected(kPathArt)||
				selectionHelper.IsArtSelected(kPlacedArt) ||
				selectionHelper.IsArtSelected(kRasterArt) ||
				selectionHelper.IsArtSelected(kGroupArt))
				return true;
			else 
				return false;
	}
	else
		return true;
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableLiveColor::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpLiveColor instance;
	if ("CreateTintedSpotColors" == runnableContext.GetOperation().GetName()) {
		result = instance.CreateTintedSpotColors();
	}
	else if ("GroupSpotColors" == runnableContext.GetOperation().GetName()) {
		result = instance.GroupSpotColors();
	}
	else if ("ConvertToBlackAndWhite" == runnableContext.GetOperation().GetName()) {
		result = instance.ConvertToBlackAndWhite();
	}
	else if ("ViewInColorHarmonyDialog" == runnableContext.GetOperation().GetName()) {
		result = instance.ViewInColorHarmonyDialog();
	}
	else if ("SetColorGuideBaseColor" == runnableContext.GetOperation().GetName()) {
		result = instance.SetColorGuideBaseColor();
	}
	else  {
		result = kBadParameterErr;
	}
	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableLiveColor instance_SnpRunnableLiveColor;

// End SnpLiveColor.cpp

