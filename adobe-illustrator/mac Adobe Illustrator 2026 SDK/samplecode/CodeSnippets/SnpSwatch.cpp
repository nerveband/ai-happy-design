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

/**	Provides functions that demonstrate how to create, display, edit and delete CMYK, 
	RGB and Grayscale process and spot colors and swatch groups in the current document 
	swatch list.

	@see AISwatchListSuite, AISwatchGroupSuite, AICustomColorSuite
  	@ingroup sdk_snippet
 */
class SnpSwatch
{
	public:
		
		/** Creates a new AISwatch object using SnpSwatch::GetColorTypeAndValue to get the new color
			values and prompts the user for the new swatch name.
			@param swatchGroupRef IN/OUT the group that the swatch is added to if not NULL.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr CreateSwatch(AISwatchGroupRef& swatchGroupRef);

		/** Prompts the user for the swatch name then, if found calls InspectSwatch(AISwatchRef).
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr InspectSwatch() const;

		/** Sends the name and color details to the log.
			@param swatch IN the swatch to inspect.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr InspectSwatch(const AISwatchRef& swatch) const;

		/** Calls InspectSwatch(AISwatchRef) for each swatch in the current document swatch list
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr InspectAllSwatches()const;

		/** Prompts the user for the swatch name then, if found and swatch is not a pattern 
			or gradient, calls SnpSwatch::GetColorTypeAndValue to get the new color values.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr EditSwatchColor();

		/** Prompts the user for a swatch name then if found, calls SnpSwatch::GetSwatchName.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr EditSwatchName();

		/** Prompts the user for a swatch name then deletes the swatch if found.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr DeleteSwatch() const;

		/** Creates a new AISwatchGroup object, prompting the user for the group name and 
			the kind of swatches it allows.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr CreateSwatchGroup();

		/**	Prompts the user for the name of the swatch group then lists the swatches
			contained in the group.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr InspectSwatchGroup() const;

		/**	Lists the name of all swatch groups in the current document.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr InspectAllSwatchGroups() const;

		/** Prompts the user for the name of a swatch group to edit, then gets the new 
			name for the swatch group from the user.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr EditSwatchGroupName();

		/** Prompts the user for a swatch name and swatch group name then moves the 
			swatch to the group, if the group name is empty the swatch is moved to 
			the general swatch group.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr MoveSwatchToGroup() const;

		/**	Prompts the user for the name of the swatch group to delete and whether or 
			not to delete the contained swatches, then deletes the chosen swatch group.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr DeleteSwatchGroup() const;

		/** Orders all swatches in the current documents swatch list into groups
			by kind. E.g. all pattern swatches are moved to PatternColorGroup.
			The "Sort by Kind" option on the Swatches panel sorts all the swatches 
			by kind and alphabetically, this function sorts by kind alone and filters 
			each swatch into either a Solid, Pattern or Gradient swatch group depending
			on whether the swatch contains a solid, pattern or gradient color.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr ReorderAllSwatchesByKind() const;

		//----------------------------- Utility functions ---------------------------------

		/**	Creates a new AIColor object from user specified parameters.
			@param color IN/OUT AIColor object which returns the values specified.
			@param customColor IN/OUT AICustomColor object which returns the values specified.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr GetColorTypeAndValue(AIColor& color, AICustomColor& customColor);

		/**	Displays the color values associated with the color kind.
			@param color IN AIColor object whose values are displayed.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr DisplayColor(const AIColor& color) const;

		/** Checks the swatch name the user entered against the swatch names already in use.
			@param name IN swatch name to look for in current document swatch list.
			@return true if name is already in use, false otherwise.
		 */
		ASBoolean CheckSwatchName(const ai::UnicodeString& name);

		/**	Writes the swatch name to the log.
			@param swatch IN the swatch containing the name to be displayed.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr InspectSwatchName(const AISwatchRef& swatch) const;

		/** Checks the swatch group name the user entered against the swatch group names already in use.
			@param name IN swatch group name to look for in current document swatch list.
			@return true if name is already in use, false otherwise.
		 */
		ASBoolean CheckSwatchGroupName(const ai::UnicodeString& name);

		/**	Writes the swatch group name to the log..
			@param swatchGroup IN the swatch group containing the name to be displayed.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr InspectSwatchGroupName(const AISwatchGroupRef& swatchGroup) const;
};


/*
*/
ASErr SnpSwatch::CreateSwatch(AISwatchGroupRef& swatchGroupRef)
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog::Indent indent;
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();		

		AISwatchRef newSwatch = NULL;
		AIColor color;
		if (swatchGroupRef) {
			AISwatchGroupKind groupKind = sAISwatchGroup->GetSwatchGroupKind(swatchGroupRef);
			if (groupKind == kAISGKindGeneric || groupKind == kAISGKindSolidColorsOnly) {
				// Adds the new swatch to the end of the specified swatch group.
				newSwatch = sAISwatchGroup->InsertNthSwatch(swatchGroupRef, &color, -1);
			}
			else {
				log->Write(ai::UnicodeString("Cannot create a new swatch in this type of group."));
				aisdk::check_ai_error(kBadParameterErr);
			}
		}
		else {
			// Add to the general swatch group of the current document swatch list, 
			// -1 adds the swatch to the end of the general swatch group.
			newSwatch = sAISwatchList->InsertNthSwatch(NULL, -1);
		}

		if (newSwatch) {
			// Get the name for the new swatch from the user.
			ai::UnicodeString swatchName;
			ASBoolean validName = false;
			while (!validName) {				
				swatchName = parameter->GetString(ai::UnicodeString("Swatch name:"), 
					ai::UnicodeString("SnpSwatch::CreateSwatch"));
				validName = this->CheckSwatchName(swatchName);
				if (!validName) {
					sAIUser->ErrorAlert(ai::UnicodeString("Invalid swatch name, please enter another."));
				}
			}

			AICustomColor customColor;
			result = this->GetColorTypeAndValue(color, customColor);
			aisdk::check_ai_error(result);

			if (color.kind == kCustomColor) {
				// Get a new unique name for the custom color from the swatch name.
				result = sAICustomColor->NewCustomColorName(swatchName);
				aisdk::check_ai_error(result);
				// Create the new kCustomColor using the custom color created in 
				// SnpSwatch::GetColorTypeAndValue().
				AICustomColorHandle customColorHandle = NULL;
				result = sAICustomColor->NewCustomColor(&customColor, swatchName, &customColorHandle);
				aisdk::check_ai_error(result);
				color.c.c.color = customColorHandle;
			}	
			else {
				result = sAISwatchList->SetSwatchName(newSwatch, swatchName);
				aisdk::check_ai_error(result);
			}

			result = sAISwatchList->SetAIColor(newSwatch, &color);
			aisdk::check_ai_error(result);

			// Write the swatch details to the log.
			result = this->InspectSwatchName(newSwatch);
			aisdk::check_ai_error(result);

			result = this->DisplayColor(color);
			aisdk::check_ai_error(result);
		}
		else {
			log->Write(ai::UnicodeString("Swatch creation failed"));
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
ASErr SnpSwatch::InspectSwatch() const
{
	ASErr result = kNoErr;
	try {	
		SnippetRunnerLog::Indent indent;
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ai::UnicodeString swatchName = parameter->GetString(ai::UnicodeString("Swatch name:"), 
			ai::UnicodeString("K=100"));

		// Find the swatch specified by the user.
		AISwatchRef swatchRef = sAISwatchList->GetSwatchByName(NULL, swatchName);
		if (swatchRef) {
			// Inspect the swatch.
			result = this->InspectSwatch(swatchRef);
			aisdk::check_ai_error(result);
		}
		else {
			log->Write(ai::UnicodeString("Swatch not found."));
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
ASErr SnpSwatch::InspectSwatch(const AISwatchRef& swatch) const
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		// Display the swatch's name.
		ai::UnicodeString swatchName;
		result = sAISwatchList->GetSwatchName(swatch, swatchName);
		aisdk::check_ai_error(result);
		log->Write(ai::UnicodeString("Swatch name: ").append(swatchName));

		// Display the swatch's color.
		AIColor color;
		result = sAISwatchList->GetAIColor(swatch,&color);
		aisdk::check_ai_error(result);
		result = this->DisplayColor(color);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpSwatch::InspectAllSwatches() const
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		// Get count of all current document swatches.
		int swatchCount = sAISwatchList->CountSwatches(NULL);
		// Inspect each swatch.
		for (int i = 0; i < swatchCount; i++) {
			log->WritePrintf("Swatch %d",i);
			SnippetRunnerLog::Indent indent;

			AISwatchRef swatchRef = sAISwatchList->GetNthSwatch(NULL, i);
			result = this->InspectSwatch(swatchRef);
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
ASErr SnpSwatch::EditSwatchColor()
{
	ASErr result = kNoErr;
	try {	
		SnippetRunnerLog::Indent indent;
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ai::UnicodeString swatchName = parameter->GetString(ai::UnicodeString("Swatch name:"), 
			ai::UnicodeString("K=100")); // 'K=100' swatch is in all document presets.

		// Get the AISwatchRef from the name.
		AISwatchRef swatchRef = sAISwatchList->GetSwatchByName(NULL, swatchName);
		if (swatchRef) {
			AIColor color;
			AICustomColor customColor;
			result = sAISwatchList->GetAIColor(swatchRef, &color);
			aisdk::check_ai_error(result);

			if (color.kind == kGrayColor || color.kind == kFourColor || color.kind == kThreeColor
				|| color.kind == kCustomColor) {
			
				result = this->GetColorTypeAndValue(color, customColor);
				aisdk::check_ai_error(result);

				if (color.kind == kCustomColor) {
					// Reset the swatch name to "", so NewCustomColorName does not have to put a 1
					// on the end to make the name unique.
					result = sAISwatchList->SetSwatchName(swatchRef, ai::UnicodeString(""));
					aisdk::check_ai_error(result);
					// Get a new unique name for the custom color from the swatch name.
					result = sAICustomColor->NewCustomColorName(swatchName);
					aisdk::check_ai_error(result);
					// Create the new kCustomColor using the custom color created in 
					// SnpSwatch::GetColorTypeAndValue().
					AICustomColorHandle customColorHandle = NULL;
					result = sAICustomColor->NewCustomColor(&customColor, swatchName, 
						&customColorHandle);
					aisdk::check_ai_error(result);
					color.c.c.color = customColorHandle;
				}

				result = sAISwatchList->SetAIColor(swatchRef, &color);
				aisdk::check_ai_error(result);

				// Write the swatch name to the log.
				result = this->InspectSwatchName(swatchRef);
				aisdk::check_ai_error(result);

				result = this->DisplayColor(color);
				aisdk::check_ai_error(result);
			}
			else {
				log->Write("Able to edit RGB, CMYK and Grayscale, process and spot colors only.");	
			}
		}
		else {
			log->Write(ai::UnicodeString("Swatch not found."));
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
ASErr SnpSwatch::EditSwatchName()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog::Indent indent;
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ai::UnicodeString swatchName = parameter->GetString(ai::UnicodeString("Swatch name to edit:"), 
			ai::UnicodeString("K=100"));
		// Find the swatch specified by the user.
		AISwatchRef swatchRef = sAISwatchList->GetSwatchByName(NULL, swatchName);
		if (swatchRef) {
			// Get the new swatch name.
			ai::UnicodeString newSwatchName;
			ASBoolean validName = false;
			while (!validName) {				
				newSwatchName = parameter->GetString(ai::UnicodeString("New swatch name:"), 
					ai::UnicodeString("SnpSwatch::EditSwatchName"));
				validName = this->CheckSwatchName(newSwatchName);
				if (!validName) {
					sAIUser->ErrorAlert(ai::UnicodeString("Invalid swatch name, please enter another."));
				}
			}

			result = sAISwatchList->SetSwatchName(swatchRef, newSwatchName);
			aisdk::check_ai_error(result);

			// Display the new swatch name.
			result = this->InspectSwatchName(swatchRef);
			aisdk::check_ai_error(result);
		}
		else {
			log->Write(ai::UnicodeString("Swatch not found."));
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
ASErr SnpSwatch::DeleteSwatch() const
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog::Indent indent;
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();		
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ai::UnicodeString swatchName = parameter->GetString(ai::UnicodeString("Swatch name:"), 
			ai::UnicodeString("K=100"));
		// Find the swatch specified by the user.
		AISwatchRef swatchRef = sAISwatchList->GetSwatchByName(NULL, swatchName);
		if (swatchRef) {
			// Remove the swatch from the current document swatch list, deleting the
			// custom color as well if necessary
			result = sAISwatchList->RemoveSwatch(NULL, swatchRef, true);
			aisdk::check_ai_error(result);
		}
		else {
			log->Write(ai::UnicodeString("Swatch not found."));
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
ASErr SnpSwatch::CreateSwatchGroup()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerLog::Indent indent;
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();

		// Get swatch group name.
		ai::UnicodeString swatchGroupName;
		ASBoolean validName = false;
		while (!validName) {				
			swatchGroupName = parameter->GetString(ai::UnicodeString("Swatch group name:"), 
				ai::UnicodeString("SnpSwatch::CreateSwatchGroup"));
			validName = this->CheckSwatchGroupName(swatchGroupName);
			if (!validName) {
				sAIUser->ErrorAlert(ai::UnicodeString("Invalid swatch group name, please enter another."));
			}
		}

		// Create kind choices.
		std::vector<ai::UnicodeString> choices;
		choices.push_back(ai::UnicodeString("kAISGKindGeneric"));
		choices.push_back(ai::UnicodeString("kAISGKindSolidColorsOnly"));
		choices.push_back(ai::UnicodeString("kAISGKindPatternsOnly"));
		choices.push_back(ai::UnicodeString("kAISGKindGradientsOnly"));
		AISwatchGroupKind swatchGroupKind = static_cast<AISwatchGroupKind>(parameter->GetChoice(
			ai::UnicodeString("Swatches allowed in group:"), choices, 0));
		
		int index = sAISwatchGroup->CountSwatchGroups(NULL); // NULL for current document swatch list.
		AISwatchGroupRef swatchGroupRef = NULL;

		result = sAISwatchGroup->NewSwatchGroup(NULL, swatchGroupKind, index, &swatchGroupRef);
		aisdk::check_ai_error(result);

		// Set swatch group name.
		result = sAISwatchGroup->SetSwatchGroupName(swatchGroupRef, swatchGroupName);
		aisdk::check_ai_error(result);

		result = this->InspectSwatchGroupName(swatchGroupRef);
		aisdk::check_ai_error(result);

		log->Write(choices[swatchGroupKind]);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpSwatch::InspectSwatchGroup() const
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog::Indent indent;
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();		
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ai::UnicodeString swatchGroupName = parameter->GetString(ai::UnicodeString("Swatch group name:"), 
			ai::UnicodeString("Grayscale"));

		AISwatchGroupRef swatchGroupRef = sAISwatchGroup->GetSwatchGroupByName(NULL,swatchGroupName);
		if (swatchGroupRef) {
			// Write details to the log.
			log->Write(ai::UnicodeString("Swatch group: ").append(swatchGroupName));
			int swatchCount = sAISwatchGroup->CountSwatches(swatchGroupRef);

			for (int i = 0; i < swatchCount; i++) {
				// Displays the details of the swatches contained in the chosen swatch group.
				AISwatchRef swatchRef = sAISwatchGroup->GetNthSwatch(swatchGroupRef, i);
				if (swatchRef) {
					log->WritePrintf("Swatch %d:",i);
					SnippetRunnerLog::Indent indent2;

					AIColor swatchColor;
					result = sAISwatchList->GetAIColor(swatchRef, &swatchColor);
					aisdk::check_ai_error(result);

					result = this->InspectSwatchName(swatchRef);
					aisdk::check_ai_error(result);

					result = this->DisplayColor(swatchColor);
					aisdk::check_ai_error(result);
				}
			}			
		}
		else {
			log->Write(ai::UnicodeString("Swatch group not found."));
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
ASErr SnpSwatch::InspectAllSwatchGroups() const
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerLog::Indent indent;
		// Count the number of swatch groups in the current document swatch list.
		int countGroups = sAISwatchGroup->CountSwatchGroups(NULL);
		for (int i = 0; i < countGroups; i++) {
			// Write each swatch group name to the log.
			// NOTE: The general swatch group at index 0 does not have a name.
			log->WritePrintf("Swatch Group %d", i);
			AISwatchGroupRef swatchGroupRef = sAISwatchGroup->GetNthSwatchGroup(NULL, i);
			result = this->InspectSwatchGroupName(swatchGroupRef);
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
ASErr SnpSwatch::EditSwatchGroupName()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog::Indent indent;
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ai::UnicodeString swatchGroupName = parameter->GetString(ai::UnicodeString("Swatch group name to edit:"), 
			ai::UnicodeString("Grayscale"));
		// Find the swatch group specified by the user.
		AISwatchGroupRef swatchGroupRef = sAISwatchGroup->GetSwatchGroupByName(NULL, swatchGroupName);
		if (swatchGroupRef) {
			// Get the new swatch group name.
			ai::UnicodeString newSwatchGroupName;
			ASBoolean validName = false;
			while (!validName) {				
				newSwatchGroupName = parameter->GetString(ai::UnicodeString("New swatch group name:"), 
					ai::UnicodeString("SnpSwatch::EditSwatchGroupName"));
				validName = this->CheckSwatchGroupName(newSwatchGroupName);
				if (!validName) {
					sAIUser->ErrorAlert(ai::UnicodeString("Invalid swatch group name, please enter another."));
				}
			}
			result = sAISwatchGroup->SetSwatchGroupName(swatchGroupRef, newSwatchGroupName);
			aisdk::check_ai_error(result);

			// Write the new swatch group name to the log.
			result = this->InspectSwatchGroupName(swatchGroupRef);
			aisdk::check_ai_error(result);
		}
		else {
			log->Write(ai::UnicodeString("Swatch group not found."));
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
ASErr SnpSwatch::MoveSwatchToGroup() const
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog::Indent indent;
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();		
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ai::UnicodeString swatchName = parameter->GetString(ai::UnicodeString("Swatch name:"), 
			ai::UnicodeString("Black"));
		
		AISwatchRef swatchRef = sAISwatchList->GetSwatchByName(NULL, swatchName);
		if (swatchRef) {
			AISwatchGroupRef fromSwatchGroup = NULL;
			ai::int32 groupIndex = 0, swatchIndex = 0;

			// Finds the exact location of the chosen swatch.
			result = sAISwatchGroup->LocateSwatch(NULL, swatchRef, &fromSwatchGroup, 
				&groupIndex, &swatchIndex);
			aisdk::check_ai_error(result);

			// Gets the swatch group the swatch is to be moved to.
			ai::UnicodeString toGroupName = parameter->GetString(
				ai::UnicodeString("To group name: [Leave blank for general swatch group]"), 
				ai::UnicodeString("Grayscale"));
			AISwatchGroupRef toSwatchGroup = sAISwatchGroup->GetSwatchGroupByName(NULL, toGroupName);
			if (toSwatchGroup) {
				if (toSwatchGroup != fromSwatchGroup) {
					int toIndex = sAISwatchGroup->CountSwatches(toSwatchGroup);

					// Moves the swatch.
					result = sAISwatchGroup->ReorderSwatchBetweenGroups(fromSwatchGroup, swatchIndex, 
						toSwatchGroup, toIndex);
					if (result == kInvalidSwatchTypeForDest)
						log->Write(ai::UnicodeString("Incorrect swatch kind for group."));
					aisdk::check_ai_error(result);

					log->WritePrintf("Moved \"%s\" to \"%s\"",
						swatchName.getInStdString(kAIPlatformCharacterEncoding).c_str(),
						toGroupName.getInStdString(kAIPlatformCharacterEncoding).c_str());
				}
				else {
					log->Write(ai::UnicodeString("Swatch already in that swatch group."));
					aisdk::check_ai_error(result);
				}
			}
			else {
				log->Write(ai::UnicodeString("Swatch group not found."));
				aisdk::check_ai_error(kBadParameterErr);
			}
		}
		else {
			log->Write(ai::UnicodeString("Swatch not found."));
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
ASErr SnpSwatch::DeleteSwatchGroup() const
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog::Indent indent;
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ai::UnicodeString swatchGroupName = parameter->GetString(ai::UnicodeString("Swatch group name:"), 
			ai::UnicodeString("Grayscale"));

		AISwatchGroupRef swatchGroupRef = sAISwatchGroup->GetSwatchGroupByName(NULL, swatchGroupName);
		if (swatchGroupRef) {
			ASBoolean deleteSwatches = parameter->GetBoolean(ai::UnicodeString("Delete swatches?"));
			// If swatches are not deleted, they are moved to the general swatch group.
			result = sAISwatchGroup->RemoveSwatchGroup(NULL, swatchGroupRef, deleteSwatches);
			aisdk::check_ai_error(result);
		}
		else {
			log->Write(ai::UnicodeString("Swatch group not found."));
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
ASErr SnpSwatch::ReorderAllSwatchesByKind() const
{
	ASErr result = kNoErr;
	try {
		// Get a reference to the general swatch group, index 0 in the current document swatch list.
		AISwatchGroupRef generalSwatchGroup = sAISwatchGroup->GetNthSwatchGroup(NULL, 0);
		// Delete all existing deletable swatch groups and move their swatches into the general swatch group.
		int totalSwatchGroups = sAISwatchGroup->CountSwatchGroups(NULL);
		for (int i = 1; i < totalSwatchGroups; i++) {
			result = sAISwatchGroup->RemoveNthSwatchGroup(NULL, 1, false);
			aisdk::check_ai_error(result);
		}

		// Create the new swatch groups.
		AISwatchGroupRef solidGroup = NULL, patternGroup = NULL, gradientGroup = NULL;
		result = sAISwatchGroup->NewSwatchGroup(NULL, kAISGKindSolidColorsOnly, 1, &solidGroup);
		aisdk::check_ai_error(result);
		result = sAISwatchGroup->SetSwatchGroupName(solidGroup, ai::UnicodeString("SolidColorGroup"));
		aisdk::check_ai_error(result);
		result = sAISwatchGroup->NewSwatchGroup(NULL, kAISGKindPatternsOnly, 2, &patternGroup);
		aisdk::check_ai_error(result);
		result = sAISwatchGroup->SetSwatchGroupName(patternGroup, ai::UnicodeString("PatternColorGroup"));
		aisdk::check_ai_error(result);
		result = sAISwatchGroup->NewSwatchGroup(NULL, kAISGKindGradientsOnly, 3, &gradientGroup);
		aisdk::check_ai_error(result);
		result = sAISwatchGroup->SetSwatchGroupName(gradientGroup, ai::UnicodeString("GradientColorGroup"));
		aisdk::check_ai_error(result);

		// Iterate through all swatches in the current document swatch list and sort by type into
		// one of the 4 swatch groups created above.
		int totalSwatches = sAISwatchList->CountSwatches(NULL);
		for (int i = 0; i < totalSwatches; i++) {
			// Always get the 0th element as the swatches are renumbered after a swatch is moved.
			AISwatchRef swatchRef = sAISwatchList->GetNthSwatch(NULL, 0);
			AIColor swatchColor;
			result = sAISwatchList->GetAIColor(swatchRef, &swatchColor);
			aisdk::check_ai_error(result);
			AIColorTag colorKind = swatchColor.kind;
			if (sAISwatchGroup->AllowsColorKind(solidGroup, colorKind)) {
				// Swatch is a solid color.
				int index = sAISwatchGroup->CountSwatches(solidGroup);
				// Using index puts the swatch to the end of the group.
				result = sAISwatchGroup->ReorderSwatchBetweenGroups(generalSwatchGroup, 0, solidGroup, index);
				aisdk::check_ai_error(result);
			}
			else if (sAISwatchGroup->AllowsColorKind(patternGroup, colorKind)) {
				// Swatch is a pattern.
				int index = sAISwatchGroup->CountSwatches(patternGroup);
				// Using index puts the swatch to the end of the group.
				result = sAISwatchGroup->ReorderSwatchBetweenGroups(generalSwatchGroup, 0, patternGroup, index);
				aisdk::check_ai_error(result);
			}
			else if (sAISwatchGroup->AllowsColorKind(gradientGroup, colorKind)) {
				// Swatch is a gradient.
				int index = sAISwatchGroup->CountSwatches(gradientGroup);
				// Using index puts the swatch to the end of the group.
				result = sAISwatchGroup->ReorderSwatchBetweenGroups(generalSwatchGroup, 0, gradientGroup, index);
				aisdk::check_ai_error(result);
			}
			else {
				// Cannot go into any of the above swatch groups so leave it where it is, should never happen
				// as all color kinds should fall into one of the above three categories.
			}
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

//--------------------------------------- Utility Functions ---------------------------------------------------

/*
*/
ASErr SnpSwatch::GetColorTypeAndValue(AIColor& color, AICustomColor& customColor)
{
	ASErr result = kNoErr;
	try {		
		// Define color mode choices.
		std::vector<ai::UnicodeString> colorModes;
		colorModes.push_back(ai::UnicodeString("CMYK"));
		colorModes.push_back(ai::UnicodeString("RGB"));
		colorModes.push_back(ai::UnicodeString("Grayscale")); 
		
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ASInt32 colorMode = parameter->GetChoice(ai::UnicodeString("Color Mode:"), colorModes, 0);

		// Get the current document color model.
		short documentColorModel = 0;
		result = sAIDocument->GetDocumentColorModel(&documentColorModel);
		aisdk::check_ai_error(result);

		// Define color type choices.
		std::vector<ai::UnicodeString> colorTypes;
		colorTypes.push_back(ai::UnicodeString("Process Color"));
		colorTypes.push_back(ai::UnicodeString("Spot Color"));

		// Get new color values for the swatch.
		switch (colorMode) {
			default:
			case 0: { // CMYK.
				// Color selected is CMYK, now get whether process color or spot color.
				ASInt32 colorType = parameter->GetChoice(ai::UnicodeString("Color Mode:"), colorTypes, 0);

				AIReal cyan = parameter->GetReal(ai::UnicodeString("Cyan: [%]"), 0.0, 0.0, 100.0) / 100;
				AIReal magenta = parameter->GetReal(ai::UnicodeString("Magenta: [%]"), 0.0, 0.0, 100.0) / 100;
				AIReal yellow = parameter->GetReal(ai::UnicodeString("Yellow: [%]"), 0.0, 0.0, 100.0) / 100;
				AIReal black = parameter->GetReal(ai::UnicodeString("Black: [%]"), 0.0, 0.0, 100.0) / 100;

				if (colorType == 0) { // Process Color.
					// kFourColor
					color.kind = kFourColor;
					color.c.f.cyan = cyan;
					color.c.f.magenta = magenta;
					color.c.f.yellow = yellow;
					color.c.f.black = black;
					if(documentColorModel == kDocRGBColor)
						sAIUser->ErrorAlert(ai::UnicodeString("The CMYK color values will be converted to RGB values to align with the current document color model."));
				}
				else if(colorType == 1) { // Spot Color.
					// kCustomFourColor.
					color.kind = kCustomColor;			
					customColor.flag = kCustomSpotColor;
					customColor.kind = kCustomFourColor;
					customColor.c.f.cyan = cyan;
					customColor.c.f.magenta = magenta;
					customColor.c.f.yellow = yellow;
					customColor.c.f.black = black;
					
				}
				else {
					aisdk::check_ai_error(kBadParameterErr);
				}
				break;
			} // End CMYK.
			case 1: { // RGB.
				// Color selected is RGB, now get whether process color or spot color.
				ASInt32 colorType = parameter->GetChoice(ai::UnicodeString("Color Mode:"), colorTypes, 0);

				AIReal red = parameter->GetReal(ai::UnicodeString("Red: [%]"), 0.0, 0.0, 100.0) / 100;
				AIReal green = parameter->GetReal(ai::UnicodeString("Green: [%]"), 0.0, 0.0, 100.0) / 100;
				AIReal blue = parameter->GetReal(ai::UnicodeString("Blue: [%]"), 0.0, 0.0, 100.0) / 100;

				if (colorType == 0) { // Process Color.
					// kThreeColor
					color.kind = kThreeColor;
					color.c.rgb.red = red;
					color.c.rgb.green = green;
					color.c.rgb.blue = blue;
					if(documentColorModel == kDocCMYKColor)
						sAIUser->ErrorAlert(ai::UnicodeString("The RGB color values will be converted to CMYK values to align with the current document color model."));
				}
				else if (colorType == 1) { // Spot Color.
					// kCustomThreeColor.
					color.kind = kCustomColor;			
					customColor.flag = kCustomSpotColor;
					customColor.kind = kCustomThreeColor;
					customColor.c.rgb.red = red;
					customColor.c.rgb.green = green;
					customColor.c.rgb.blue = blue;
				}
				else {
					aisdk::check_ai_error(kBadParameterErr);
				}
				break;
			} // End RGB.
			case 2: { // Grayscale.
				// Color selected is Grayscale, now get whether process color or spot color.
				ASInt32 colorType = parameter->GetChoice(ai::UnicodeString("Color Mode:"), colorTypes, 0);

				AIReal gray = parameter->GetReal(ai::UnicodeString("Gray: [%]"), 0.0, 0.0, 100.0) / 100;

				if (colorType == 0) { // Process Color.
					// kGrayColor.
					color.kind = kGrayColor;
					color.c.g.gray = gray;
				}
				else if (colorType == 1) { // Spot Color.
					// Treat as kCustomFourColor.
					color.kind = kCustomColor;			
					customColor.flag = kCustomSpotColor;
					customColor.kind = kCustomFourColor;
					customColor.c.f.cyan = 0;
					customColor.c.f.magenta = 0;
					customColor.c.f.yellow = 0;
					customColor.c.f.black = gray;
				}
				else {
					aisdk::check_ai_error(kBadParameterErr);
				}
				break;	
			} // End Grayscale.
		}
		if(color.kind == kCustomColor) {
			// Get the tint value, 0 is full tint, 100 is white.
			color.c.c.tint = (AIReal)1.0 - (parameter->GetReal(ai::UnicodeString("Tint:"), 0.0, 0.0, 100.0) / 100);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpSwatch::DisplayColor(const AIColor& color) const
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		log->Write(ai::UnicodeString("Swatch Color:"));
		SnippetRunnerLog::Indent indent;
		// Get the value of each color and convert to a percentage.
		// CMYK values range from 0 - 100, RGB values range from 0 - 255, in the log they are displayed
		// as percent.
		AIColorTag colorKind = color.kind;
		switch (colorKind) {
			default:
			case kGrayColor: {
				log->Write(ai::UnicodeString("kGrayColor"));
				log->WritePrintf("Gray: %.1f %%", color.c.g.gray * 100);
				break;		
			}
			case kFourColor: {
				log->Write(ai::UnicodeString("kFourColor"));
				log->WritePrintf("Cyan: %.1f %%", color.c.f.cyan * 100);
				log->WritePrintf("Magenta: %.1f %%", color.c.f.magenta * 100);
				log->WritePrintf("Yellow: %.1f %%", color.c.f.yellow * 100);
				log->WritePrintf("Black: %.1f %%", color.c.f.black * 100);
				break;		
			}
			case kPattern: {
				log->Write(ai::UnicodeString("kPattern"));
				break;		
			}
			case kCustomColor: {
				log->Write(ai::UnicodeString("kCustomColor"));
				SnippetRunnerLog::Indent indent2;

				AICustomColor customColor;
				AICustomColorHandle customColorHandle = color.c.c.color;
				result = sAICustomColor->GetCustomColor(customColorHandle, &customColor);
				aisdk::check_ai_error(result);

				AICustomColorTag customColorKind = customColor.kind;
				switch (customColorKind) {
					default:
					case kCustomFourColor: {
						log->Write(ai::UnicodeString("kCustomFourColor"));
						log->WritePrintf("Cyan: %.1f %%", customColor.c.f.cyan * 100);
						log->WritePrintf("Magenta: %.1f %%", customColor.c.f.magenta * 100);
						log->WritePrintf("Yellow: %.1f %%", customColor.c.f.yellow * 100);
						log->WritePrintf("Black: %.1f %%", customColor.c.f.black * 100);
						break;
					}
					case kCustomThreeColor: {
						log->Write(ai::UnicodeString("kCustomThreeColor"));
						log->WritePrintf("Red: %.1f %%", customColor.c.rgb.red * 100);
						log->WritePrintf("Green: %.1f %%", customColor.c.rgb.green * 100);
						log->WritePrintf("Blue: %.1f %%", customColor.c.rgb.blue * 100);
						break;
					}
					case kCustomLabColor: {
						log->Write(ai::UnicodeString("kCustomLabColor"));
						log->WritePrintf("l: %.1f %%", customColor.c.lab.l * 100);
						log->WritePrintf("a: %.1f %%",customColor.c.lab.a * 100);
						log->WritePrintf("b: %.1f %%",customColor.c.lab.b * 100);
						break;
					}
				}
				log->WritePrintf("Tint: %.1f %%", ((AIReal)1.0 - (color.c.c.tint)) * 100);
				break;		
			}
			case kGradient: {
				log->Write(ai::UnicodeString("kGradient"));
				break;		
			}
			case kThreeColor: {
				log->Write(ai::UnicodeString("kThreeColor"));
				log->WritePrintf("Red: %.1f %%", color.c.rgb.red * 100);
				log->WritePrintf("Green: %.1f %%", color.c.rgb.green * 100);
				log->WritePrintf("Blue: %.1f %%", color.c.rgb.blue * 100);
				break;		
			}
			case kNoneColor: {
				log->Write(ai::UnicodeString("No color"));
				break;		
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
ASBoolean SnpSwatch::CheckSwatchName(const ai::UnicodeString& name)
{
	int nameEntered = name.compare(ai::UnicodeString(""));
	if (nameEntered == 0) {
		return false;
	}
	else {
		AISwatchRef swatchRef = sAISwatchList->GetSwatchByName(NULL, name);
		if (swatchRef) {
			return false;
		}
		return true;
	}
	return true;
}

/*
*/
ASErr SnpSwatch::InspectSwatchName(const AISwatchRef& swatch) const
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		ai::UnicodeString swatchName;
		result = sAISwatchList->GetSwatchName(swatch, swatchName);
		aisdk::check_ai_error(result);
		log->Write(ai::UnicodeString("Swatch Name: ").append(swatchName));
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASBoolean SnpSwatch::CheckSwatchGroupName(const ai::UnicodeString& name)
{
	int nameEntered = name.compare(ai::UnicodeString(""));
	if (nameEntered == 0) {
		return false;
	}
	else {
		AISwatchGroupRef swatchGroupRef = sAISwatchGroup->GetSwatchGroupByName(NULL, name);
		if (swatchGroupRef) {
			return false;
		}
		return true;
	}
	return true;
}

/*
*/
ASErr SnpSwatch::InspectSwatchGroupName(const AISwatchGroupRef& swatchGroup) const
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerLog::Indent indent;
		ai::UnicodeString swatchGroupName;
		result = sAISwatchGroup->GetSwatchGroupName(swatchGroup, swatchGroupName);
		aisdk::check_ai_error(result);
		log->Write(ai::UnicodeString("Swatch Group Name: ").append(swatchGroupName));
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpSwatch available to the SnippetRunner framework.
 */
class _SnpRunnableSwatch : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableSwatch () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableSwatch () {}

		/* Returns name of snippet.
		*/
		std::string GetName() const;

		/* Returns a description of what the snippet does.
		*/
		std::string	GetDescription() const;

		/* Returns operations supported by this snippet.
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
		ASBoolean	CanRun(SnpRunnable::Context& runnableContext);

		/* Runs the snippet.
			@param runnableContext see ISnpRunnableContext for documentation.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr	Run(SnpRunnable::Context& runnableContext);
};

/*
*/
std::string _SnpRunnableSwatch::GetName() const
{
	return "Swatch";
}

/*
*/
std::string _SnpRunnableSwatch::GetDescription() const 
{
	return "Manipulates colors, swatches and swatch groups";
}

/*
*/
SnpRunnable::Operations _SnpRunnableSwatch::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("CreateSwatch", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("InspectSwatch", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("InspectAllSwatches", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("EditSwatchColor", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("EditSwatchName", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("DeleteSwatch", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("CreateSwatchGroup", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("InspectSwatchGroup", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("InspectAllSwatchGroups", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("EditSwatchGroupName", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("MoveSwatchToGroup", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("DeleteSwatchGroup", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("ReorderAllSwatchesByKind", "document", kSnpRunUnitTestDocumentContext));
	return operations;
}

/*
*/
std::string _SnpRunnableSwatch::GetDefaultOperationName() const 
{
	return "CreateSwatch";
}

/*
*/
std::vector<std::string> _SnpRunnableSwatch::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Color Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableSwatch::CanRun(SnpRunnable::Context& runnableContext)
{
	SnpSelectionHelper selectionHelper;
	return selectionHelper.IsDocumentSelected();
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableSwatch::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpSwatch instance;
	if ("CreateSwatch" == runnableContext.GetOperation().GetName()) {
		AISwatchGroupRef swatchGroupRef = NULL;
		result = instance.CreateSwatch(swatchGroupRef);
	}
	else if ("InspectSwatch" == runnableContext.GetOperation().GetName()) {
		result = instance.InspectSwatch();
	}
	else if ("InspectAllSwatches" == runnableContext.GetOperation().GetName()) {
		result = instance.InspectAllSwatches();
	}
	else if ("EditSwatchColor" == runnableContext.GetOperation().GetName()) {
		result = instance.EditSwatchColor();
	}
	else if ("EditSwatchName" == runnableContext.GetOperation().GetName()) {
		result = instance.EditSwatchName();
	}
	else if ("DeleteSwatch" == runnableContext.GetOperation().GetName()) {
		result = instance.DeleteSwatch();
	}
	else if ("CreateSwatchGroup" == runnableContext.GetOperation().GetName()) {
		result = instance.CreateSwatchGroup();
	}
	else if ("InspectSwatchGroup" == runnableContext.GetOperation().GetName()) {
		result = instance.InspectSwatchGroup();
	}
	else if ("InspectAllSwatchGroups" == runnableContext.GetOperation().GetName()) {
		result = instance.InspectAllSwatchGroups();
	}
	else if ("EditSwatchGroupName" == runnableContext.GetOperation().GetName()) {
		result = instance.EditSwatchGroupName();
	}
	else if ("MoveSwatchToGroup" == runnableContext.GetOperation().GetName()) {
		result = instance.MoveSwatchToGroup();
	}
	else if ("DeleteSwatchGroup" == runnableContext.GetOperation().GetName()) {
		result = instance.DeleteSwatchGroup();
	}
	else if ("ReorderAllSwatchesByKind" == runnableContext.GetOperation().GetName()) {
		result = instance.ReorderAllSwatchesByKind();
	}
	else  {
		result = kBadParameterErr;
	}
	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableSwatch instance_SnpRunnableSwatch;

// End SnpSwatch.cpp

