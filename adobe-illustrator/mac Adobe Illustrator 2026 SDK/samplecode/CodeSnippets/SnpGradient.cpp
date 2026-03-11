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
#include "SnpArtHelper.h"
#include "SnpArtSetHelper.h"

/** Demonstrates create, edit, display and delete of gradients.
	@ingroup sdk_snippet
 */
class SnpGradient
{
	public:
		/** Creates a new linear gradient in the current document then calls 
			SnpGradient::CreateGradient to set up the generic gradient values.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr CreateLinearGradient();

		/** Creates a new radial gradient in the current document then calls 
			SnpGradient::CreateGradient to set up the generic gradient values.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr CreateRadialGradient();

		/** Sets the generic gradient values such as the gradient stops, opacity etc.
			@param gradient IN/OUT buffer in which to return the updated gradient object.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr CreateGradient(AIGradientHandle& gradient);

		/** Edits the name of an existing gradient in the current document.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr EditGradientName();

		/** Edits a gradient stop of an existing gradient in the current document.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr EditGradientStop();

		/** Displays the chosen gradients details in the log.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr DisplayGradient();

		/** Displays the gradients details in the log.
			@param gradient IN the gradient to display.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr DisplayGradient(const AIGradientHandle gradient);

		/** Deletes the chosen gradient from the current document.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr DeleteGradient();

		//----------------------------- Utility functions ---------------------------------

		/** Collects the names of all gradients in a document.
			@param names IN/OUT buffer in which to return the gradient names.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr GetGradientNames(std::vector<ai::UnicodeString>& names);

		/** Prompts user for new gradient stop color and opacity values.
			@param gradientStop IN/OUT buffer in which to return new gradient stop values
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr SetGradientStop(const AIGradientHandle gradient, short index, AIGradientStop& gradientStop);

		/** Prompts user for red, green and blue color values.
			@param color IN/OUT buffer in which to return the color values.
			@param index IN the gradient stop index.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr GetRGBColor(AIColor& color, short index);

		/** Prompts user for cyan, magenta, yellow and black color values.
			@param color IN/OUT buffer in which to return the color values.
			@param index IN the gradient stop index.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr GetCMYKColor(AIColor& color, short index);

		/** Displays the name of the gradient in the log.
			@param gradient IN the gradient to display.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr DisplayGradientName(const AIGradientHandle gradient);

		/** Displays the color and opacity values of the gradient stop in the log.
			@param gradientStop IN the gradient stop to display.
			@param stopNumber IN gradient stop index.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr DisplayGradientStop(const AIGradientStop gradientStop, short stopNumber);

		/** Creates a new gradient swatch.
			@param gradient IN the gradient to apply to the new swatch.
			@param swatch IN/OUT buffer in which to return the swatch ref.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr CreateGradientSwatch(const AIGradientHandle gradient, AISwatchRef& swatch);

		/** Applies a gradient to a new art item or selected art item.
			@param swatch IN the gradient swatch to apply to the new or selected art.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr ApplyGradientSwatchToArt(const AISwatchRef swatch);

		/** Sets the fill color of the path style of an art item.
			@param art IN the art item to apply the path style to.
			@param swatch IN the gradient swatch to apply to the new or selected art.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr SetArtStyle(const AIArtHandle art, const AISwatchRef swatch);

};

/*
*/
ASErr SnpGradient::CreateLinearGradient()
{
	ASErr result = kNoErr;
	try {
		// Create new gradient.
		AIGradientHandle gradient = NULL;
		result = sAIGradient->NewGradient(&gradient);
		aisdk::check_ai_error(result);

		// Get gradient name.
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ai::UnicodeString gradientName = parameter->GetString(ai::UnicodeString("Enter new gradient name"), ai::UnicodeString("My Linear Gradient"));
		
		// Create unique gradient name using string input, will add a number on the end if the 
		// name is already in use.
		result = sAIGradient->NewGradientName(gradientName);
		aisdk::check_ai_error(result);

		// Set the gradient name.
		result = sAIGradient->SetGradientName(gradient, gradientName);
		aisdk::check_ai_error(result);

		// Set the gradient type.
		result = sAIGradient->SetGradientType(gradient, kLinearGradient);
		aisdk::check_ai_error(result);

		// Get the generic values.
		result = this->CreateGradient(gradient);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGradient::CreateRadialGradient()
{
	ASErr result = kNoErr;
	try {
		// Create new gradient.
		AIGradientHandle gradient = NULL;
		result = sAIGradient->NewGradient(&gradient);
		aisdk::check_ai_error(result);

		// Get gradient name.
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ai::UnicodeString gradientName = parameter->GetString(ai::UnicodeString("Enter new gradient name"), ai::UnicodeString("My Radial Gradient"));
		
		// Create unique gradient name using string input, will add a number on the end if the 
		// name is already in use.
		result = sAIGradient->NewGradientName(gradientName);
		aisdk::check_ai_error(result);

		// Set the gradient name.
		result = sAIGradient->SetGradientName(gradient, gradientName);
		aisdk::check_ai_error(result);

		// Set the gradient type.
		result = sAIGradient->SetGradientType(gradient, kRadialGradient);
		aisdk::check_ai_error(result);

		// Get the generic values.
		result = this->CreateGradient(gradient);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGradient::CreateGradient(AIGradientHandle& gradient)
{
	ASErr result = kNoErr;
	try {
		// Get the number of gradient stops required for the gradient.
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		short gradientStopCount = short(parameter->GetInt(ai::UnicodeString("Enter number of gradient stops: "), 2, 2, 50));

		// Define first and last stop index.
		short firstStop = 0;
		short lastStop = gradientStopCount - 1;

		// Get values for all gradient stops from user.
		for (short i = 0; i < gradientStopCount; i++) {
			ai::NumberFormat numFormat;
			ai::UnicodeString indexString;
			indexString = numFormat.toString(i+1, indexString);
			parameter->Alert(ai::UnicodeString("Set gradient stop ").append(indexString));
			AIGradientStop gradientStop;
			if (i != firstStop && i != lastStop) {
				// Get new gradient stop values from user.
				result = this->SetGradientStop(gradient, i, gradientStop);
				aisdk::check_ai_error(result);
				// Add a new stop with specified values.
				result = sAIGradient->InsertGradientStop(gradient, i, &gradientStop);
				aisdk::check_ai_error(result);
			}
			else {
				// Get the existing gradient stop values.
				result = sAIGradient->GetNthGradientStop(gradient, i, &gradientStop);
				aisdk::check_ai_error(result);
				// Get the new gradient stop values from the user.
				result = this->SetGradientStop(gradient, i, gradientStop);
				aisdk::check_ai_error(result);
				// Set the gradient stop values.
				result = sAIGradient->SetNthGradientStop(gradient, i, &gradientStop);
				aisdk::check_ai_error(result);
			}
		}

		// Create gradient swatch.
		AISwatchRef swatch;
		result = this->CreateGradientSwatch(gradient, swatch);
		aisdk::check_ai_error(result);

		// Apply gradient swatch to new art.
		result = this->ApplyGradientSwatchToArt(swatch);
		aisdk::check_ai_error(result);

		// Display gradient details in the log.
		result = this->DisplayGradient(gradient);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGradient::EditGradientName()
{
	ASErr result = kNoErr;
	try {
		// Create a vector containing the names of all gradient in the document.
		std::vector<ai::UnicodeString> names;
		result = this->GetGradientNames(names);
		aisdk::check_ai_error(result);

		ASInt32 count = 0;
		result = sAIGradient->CountGradients(&count);
		aisdk::check_ai_error(result);

		AIGradientHandle chosenGradient = NULL;

		if (count == 0) {
			SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
			parameter->Alert("There are no gradients in this document.");
		}
		else {
			// Get gradient to edit.
			SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
			ASInt32 choice = parameter->GetChoice(ai::UnicodeString("Select gradient to edit:"), names);
			ai::UnicodeString chosenGradientName = names[choice];

			result = sAIGradient->GetGradientByName(chosenGradientName, &chosenGradient);
			aisdk::check_ai_error(result);

			SDK_ASSERT(chosenGradient);

			ai::UnicodeString newName = parameter->GetString(ai::UnicodeString("Enter new name:"), ai::UnicodeString("My New Gradient"));

			// Create unique gradient name using string input.
			result = sAIGradient->NewGradientName(newName);
			aisdk::check_ai_error(result);

			result = sAIGradient->SetGradientName(chosenGradient, newName);
			aisdk::check_ai_error(result);

			// Display the new name in the log.
			result = this->DisplayGradientName(chosenGradient);
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
ASErr SnpGradient::EditGradientStop()
{
	ASErr result = kNoErr;
	try {
		// Create a vector containing the names of all gradient in the document.
		std::vector<ai::UnicodeString> names;
		result = this->GetGradientNames(names);
		aisdk::check_ai_error(result);

		ASInt32 count = 0;
		result = sAIGradient->CountGradients(&count);
		aisdk::check_ai_error(result);

		AIGradientHandle chosenGradient = NULL;

		if (count == 0) {
			SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
			parameter->Alert("There are no gradients in this document.");
		}
		else {
			// Get gradient stop to edit.
			SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
			ASInt32 choice = parameter->GetChoice(ai::UnicodeString("Select gradient to edit:"), names);
			ai::UnicodeString chosenGradientName = names[choice];

			result = sAIGradient->GetGradientByName(chosenGradientName, &chosenGradient);
			aisdk::check_ai_error(result);

			SDK_ASSERT(chosenGradient);

			short gradientStopCount = 0;
			result = sAIGradient->GetGradientStopCount(chosenGradient, &gradientStopCount);
			aisdk::check_ai_error(result);

			// Select gradient stop to edit.
			ASInt32 stop = parameter->GetInt(ai::UnicodeString("Select gradient stop to edit:"), 0, 0, gradientStopCount - 1);
			
			// Get the chosen gradient stop.
			AIGradientStop gradientStop;
			result = sAIGradient->GetNthGradientStop(chosenGradient, short(stop), &gradientStop);
			aisdk::check_ai_error(result);

			result = this->SetGradientStop(chosenGradient, short(stop), gradientStop);
			aisdk::check_ai_error(result);

			// Set the new gradient stop values.
			result = sAIGradient->SetNthGradientStop(chosenGradient, short(stop), &gradientStop);
			aisdk::check_ai_error(result);

			// Display the new values in the log.
			result = this->DisplayGradientStop(gradientStop, short(stop));
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
ASErr SnpGradient::DisplayGradient()
{
	ASErr result = kNoErr;
	try {
		// Create a vector containing the names of all gradient in the document.
		std::vector<ai::UnicodeString> names;
		result = this->GetGradientNames(names);
		aisdk::check_ai_error(result);

		ASInt32 count = 0;
		result = sAIGradient->CountGradients(&count);
		aisdk::check_ai_error(result);

		AIGradientHandle chosenGradient = NULL;

		if (count == 0) {
			SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
			parameter->Alert("There are no gradients in this document.");
		}
		else {
			// Get gradient to display.
			SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
			ASInt32 choice = parameter->GetChoice(ai::UnicodeString("Select gradient to display:"), names);
			ai::UnicodeString chosenGradientName = names[choice];

			result = sAIGradient->GetGradientByName(chosenGradientName, &chosenGradient);
			aisdk::check_ai_error(result);

			result = this->DisplayGradient(chosenGradient);
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
ASErr SnpGradient::DisplayGradient(AIGradientHandle gradient)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(gradient);

		result = this->DisplayGradientName(gradient);
		aisdk::check_ai_error(result);

		SnippetRunnerLog::Indent indent;

		short gradientStopCount = 0;
		result = sAIGradient->GetGradientStopCount(gradient, &gradientStopCount);
		aisdk::check_ai_error(result);
		for (short i = 0; i < gradientStopCount; i++) {
			AIGradientStop gradientStop;
			result = sAIGradient->GetNthGradientStop(gradient, i, &gradientStop);
			aisdk::check_ai_error(result);

			result = this->DisplayGradientStop(gradientStop, i);
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
ASErr SnpGradient::DeleteGradient()
{
	ASErr result = kNoErr;
	try {
		// Create a vector containing the names of all gradient in the document.
		std::vector<ai::UnicodeString> names;
		result = this->GetGradientNames(names);
		aisdk::check_ai_error(result);

		ASInt32 count = 0;
		result = sAIGradient->CountGradients(&count);
		aisdk::check_ai_error(result);

		AIGradientHandle chosenGradient = NULL;

		if (count == 0) {
			SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
			parameter->Alert("There are no gradients in this document.");
		}
		else {
			// Get gradient to edit.
			SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
			ASInt32 choice = parameter->GetChoice(ai::UnicodeString("Select gradient to delete:"), names);
			ai::UnicodeString chosenGradientName = names[choice];

			result = sAIGradient->GetGradientByName(chosenGradientName, &chosenGradient);
			aisdk::check_ai_error(result);

			SDK_ASSERT(chosenGradient);

			// Delete the gradient.
			result = sAIGradient->DeleteGradient(chosenGradient);
			aisdk::check_ai_error(result);
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
ASErr SnpGradient::GetGradientNames(std::vector<ai::UnicodeString>& names)
{
	ASErr result = kNoErr;
	try {
		ai::int32 count = 0;
		result = sAIGradient->CountGradients(&count);
		aisdk::check_ai_error(result);
		for (ai::int32 i = 0; i < count; i++)
		{
			AIGradientHandle gradient = NULL;
			result = sAIGradient->GetNthGradient(i, &gradient);
			aisdk::check_ai_error(result);
			SDK_ASSERT(gradient);
			
			ai::UnicodeString gradientName;
			result = sAIGradient->GetGradientName(gradient, gradientName);
			aisdk::check_ai_error(result);
			names.push_back(gradientName);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGradient::SetGradientStop(const AIGradientHandle gradient, short index, AIGradientStop& gradientStop)
{
	ASErr result = kNoErr;
	try {
		short docColorModel = 0;
		result = sAIDocument->GetDocumentColorModel(&docColorModel);
		aisdk::check_ai_error(result);

		// Get the new color values for the gradient stop.
		AIColor color;
		if (docColorModel == kDocRGBColor) {
			result = this->GetRGBColor(color, index);
			aisdk::check_ai_error(result);
		}
		else {
			result = this->GetCMYKColor(color, index);
			aisdk::check_ai_error(result);
		}
		
		// Get valid range the gradient ramp point can be set to by looking at the next and previous stop ramp points.
		short stopCount = 0;
		result = sAIGradient->GetGradientStopCount(gradient, &stopCount);
		aisdk::check_ai_error(result);
		AIGradientStop previousStop, nextStop;
		AIReal previousRampPoint = 1, nextRampPoint = 100, defaultRampPoint = 1;
		if ((index > 0) && (index != (stopCount - 1))) {
			// Get previous stop ramp point.
			result = sAIGradient->GetNthGradientStop(gradient, index - 1, &previousStop);
			aisdk::check_ai_error(result);
			previousRampPoint = previousStop.rampPoint;
			// Set the default ramp point.
			defaultRampPoint = (previousRampPoint + ((nextRampPoint - previousRampPoint)/2));
		}
		else if (index == stopCount - 1) {
			// Set the default ramp point.
			defaultRampPoint = 100;
		}

		if ((index < stopCount - 1) && index != 0) {
			// Get next stop ramp point
			result = sAIGradient->GetNthGradientStop(gradient, index + 1, &nextStop);
			aisdk::check_ai_error(result);
			nextRampPoint = nextStop.rampPoint;
			// Set the default ramp point.
			defaultRampPoint = (previousRampPoint + ((nextRampPoint - previousRampPoint)/2));
		}
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		AIReal opacity = parameter->GetReal(ai::UnicodeString("Enter gradient stop opacity value:"), 100, 0, 100);
		AIReal rampPoint = parameter->GetReal(ai::UnicodeString("Enter gradient stop ramp point:"), defaultRampPoint, previousRampPoint, nextRampPoint);
		AIReal midPoint = parameter->GetReal(ai::UnicodeString("Enter gradient stop mid point:"), 50, 13, 87);
		gradientStop.opacity = opacity / 100;
		gradientStop.color = color;
		gradientStop.rampPoint = rampPoint;
		gradientStop.midPoint = midPoint;
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGradient::GetRGBColor(AIColor& color, short index)
{
	ASErr result = kNoErr;
	try {
		color.kind = kThreeColor;
		// Set the default color value to 0 if first gradient stop, 255 otherwise.
		AIReal defaultColorVal = 255;
		if (index > 0) {
			defaultColorVal = 0;
		}
		// Get the values from the user.
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		AIReal red = parameter->GetReal(ai::UnicodeString("Enter red value:"), defaultColorVal, 0, 255);
		AIReal green = parameter->GetReal(ai::UnicodeString("Enter green value:"), defaultColorVal, 0, 255);
		AIReal blue = parameter->GetReal(ai::UnicodeString("Enter blue value:"), defaultColorVal, 0, 255);
		color.c.rgb.red = red / 255;
		color.c.rgb.green = green / 255;
		color.c.rgb.blue = blue / 255;
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGradient::GetCMYKColor(AIColor& color, short index)
{
	ASErr result = kNoErr;
	try {
		color.kind = kFourColor;
		// Set the default color value to 0 if first gradient stop, 100 otherwise.
		AIReal defaultColorVal = 0;
		if (index > 0) {
			defaultColorVal = 100;
		}
		// Get the values from the user.
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		AIReal cyan = parameter->GetReal(ai::UnicodeString("Enter cyan value:"), defaultColorVal, 0, 100);
		AIReal magenta = parameter->GetReal(ai::UnicodeString("Enter magenta value:"), defaultColorVal, 0, 100);
		AIReal yellow = parameter->GetReal(ai::UnicodeString("Enter yellow value:"), defaultColorVal, 0, 100);
		AIReal black = parameter->GetReal(ai::UnicodeString("Enter black value:"), defaultColorVal, 0, 100);
		color.c.f.cyan = cyan / 100;
		color.c.f.magenta = magenta / 100;
		color.c.f.yellow = yellow / 100;
		color.c.f.black = black / 100;
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGradient::DisplayGradientName(const AIGradientHandle gradient)
{
	ASErr result = kNoErr;
	try {
		ai::UnicodeString gradientName;
		result = sAIGradient->GetGradientName(gradient, gradientName);
		aisdk::check_ai_error(result);

		SnippetRunnerLog::Instance()->Write(ai::UnicodeString("Gradient Name: ").append(gradientName));
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGradient::DisplayGradientStop(const AIGradientStop gradientStop, short stopNumber)
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog::Instance()->WritePrintf("Gradient stop %d:", stopNumber + 1);
		SnippetRunnerLog::Indent indent;
		AIColor color = gradientStop.color;
		if (color.kind == kThreeColor) {
			SnippetRunnerLog::Instance()->WritePrintf("Red: %.1f, Green: %.1f, Blue: %.1f", color.c.rgb.red * 255, color.c.rgb.green * 255, color.c.rgb.blue * 255);
		}
		else if (color.kind == kFourColor) {
			SnippetRunnerLog::Instance()->WritePrintf("Cyan: %.1f, Magenta: %.1f, Yellow: %.1f, Black: %.1f", color.c.f.cyan * 100, color.c.f.magenta * 100, color.c.f.yellow * 100, color.c.f.black * 100);
		}
		else if (color.kind == kGrayColor) {
			SnippetRunnerLog::Instance()->WritePrintf("Gray: %.1f", color.c.g.gray * 100);
		}
		else {
			SnippetRunnerLog::Instance()->Write(ai::UnicodeString("Unsupported color type in use."));
		}
		SnippetRunnerLog::Instance()->WritePrintf("Opacity: %.1f", gradientStop.opacity * 100);
		SnippetRunnerLog::Instance()->WritePrintf("Ramp point: %.1f", gradientStop.rampPoint);
		SnippetRunnerLog::Instance()->WritePrintf("Mid point: %.1f", gradientStop.midPoint);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGradient::CreateGradientSwatch(const AIGradientHandle gradient, AISwatchRef& swatch)
{
	ASErr result = kNoErr;
	try {
		// Insert new swatch at end of general swatch group.
		swatch = sAISwatchList->InsertNthSwatch(NULL, -1);

		// Get gradient name.
		ai::UnicodeString gradientName;
		result = sAIGradient->GetGradientName(gradient, gradientName);
		aisdk::check_ai_error(result);

		// Set swatch name.
		result = sAISwatchList->SetSwatchName(swatch, gradientName);
		aisdk::check_ai_error(result);

		// Create the swatch color using the gradient.
		AIColor swatchColor;
		swatchColor.kind = kGradient;
		swatchColor.c.b.gradient = gradient;
		swatchColor.c.b.gradientAngle = 0;
		swatchColor.c.b.gradientLength = 1;
		AIRealPoint origin = {0,0};
		swatchColor.c.b.gradientOrigin = origin;
		swatchColor.c.b.hiliteAngle = 0;
		swatchColor.c.b.hiliteLength = 0;
		
		// Set the swatch color.
		result = sAISwatchList->SetAIColor(swatch, &swatchColor);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGradient::ApplyGradientSwatchToArt(const AISwatchRef swatch)
{
	ASErr result = kNoErr;
	try {
		SnpSelectionHelper selectionHelper;
		if (selectionHelper.IsArtSelected(kPathArt)) {
			// Find all selected path art.
			const short numSpecs = 1;
			AIArtSpec artSpec[numSpecs] = { 
				{kPathArt, kArtSelected, kArtSelected},
			};
			SnpArtSetHelper artSet(artSpec, numSpecs);
			// Iterate through selected path art applying gradient swatch.
			for (size_t i = 0; i < artSet.GetCount(); i++) {
				// Set the path style of the art.
				AIArtHandle art = artSet[i];
				result = this->SetArtStyle(art, swatch);
				aisdk::check_ai_error(result);
			}
		}
		else {
			// Create new art and apply gradient swatch.
			AIRealRect rect = {300, 200, 100, 400};
			AIArtHandle art = NULL;
			SnpArtHelper artHelper;
			result = artHelper.NewRectangle(rect, false, art);
			aisdk::check_ai_error(result);

			// Set the style of the art.
			result = this->SetArtStyle(art, swatch);
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
ASErr SnpGradient::SetArtStyle(const AIArtHandle art, const AISwatchRef swatch)
{
	ASErr result = kNoErr;
	try {
		AIPathStyle pathStyle;
		result = sAIPathStyle->GetPathStyle(art, &pathStyle, nullptr);
		aisdk::check_ai_error(result);

		// Get the color from the swatch.
		AIColor color;
		result = sAISwatchList->GetAIColor(swatch, &color);
		aisdk::check_ai_error(result);

		// Set the fill of the art path style to the swatch color.
		pathStyle.fill.color = color;
        pathStyle.stroke.color = color;

		// Update the path style of the art.
		result = sAIPathStyle->SetPathStyle(art, &pathStyle);
		aisdk::check_ai_error(result);

		// Get the aspect ratio.
		AIReal fillAspectRatio = SnippetRunnerParameter::Instance()->GetReal(ai::UnicodeString("Enter fill aspect ratio (applicable only for radial gradients):"), 100, 0, 100);
		AIReal strokeAspectRatio = SnippetRunnerParameter::Instance()->GetReal(ai::UnicodeString("Enter stroke aspect ratio (applicable only for radial gradients):"), 100, 0, 100);

		// Create the new art style.
		AIArtStyleHandle artStyle;
		AIArtStylePaintData artStylePaintData;

		result = sAIArtStyle->GetArtStyle(art, &artStyle);
		aisdk::check_ai_error(result);
		AIPathStyle usedStyle;
		result = sAIArtStyle->GetPaintAttributes(artStyle, &usedStyle, &artStylePaintData);
		aisdk::check_ai_error(result);
		artStylePaintData.fillGradientAspectRatio = fillAspectRatio / 100;
		artStylePaintData.strokeGradientAspectRatio = strokeAspectRatio / 100;
		result = sAIArtStyle->NewStyle(&pathStyle, &artStylePaintData, NULL, &artStyle);
		aisdk::check_ai_error(result);

		// Apply the art style to the art item.
		result = sAIArtStyle->SetArtStyle(art, artStyle);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpGradient available to the SnippetRunner framework.
 */
class _SnpRunnableGradient : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableGradient () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableGradient () {}

		/* Returns name of snippet.
		*/
		std::string GetName() const;

		/* Returns a description of what the snippet does.
		*/
		std::string	GetDescription() const;

		/* Returns operations supported by this snippet.
		*/
		Operations GetOperations() const;

		/* Returns name of this snippet's default operation - must
			be one of the operation names returned by GetOperations.
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
		ASBoolean			CanRun(SnpRunnable::Context& runnableContext);

		/* Runs the snippet.
			@param runnableContext see ISnpRunnableContext for documentation.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr		Run(SnpRunnable::Context& runnableContext);
};

/*
*/
std::string _SnpRunnableGradient::GetName() const
{
	// TODO - edit as required.
	return "Gradient";
}

/*
*/
std::string _SnpRunnableGradient::GetDescription() const 
{
	// TODO - edit as required.
	return "Demonstrates manipulation of gradients";
}

/*
*/
SnpRunnable::Operations _SnpRunnableGradient::GetOperations() const
{
	// TODO - edit as required.
	SnpRunnable::Operations operations;
	operations.push_back(Operation("CreateLinearGradient", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("CreateRadialGradient", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("EditGradientName", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("EditGradientStop", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("DisplayGradient", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("DeleteGradient", "document", kSnpRunNewDocumentContext));
	// Add a line as above for each operation your snippet wants to advertise to the framework.
	// Operation(const std::string& name, const std::string& preconditions, eSnpRunContextType unitTestContext)
	// name - name of the operation 
	// preconditions - description of preconditions to be displayed in SnippetRunner panel
	// unitTestContext - context you want the framework to automatically set for your snippet during unit test.
	return operations;
}

/*
*/
std::string _SnpRunnableGradient::GetDefaultOperationName() const 
{
	// TODO - edit as required.
	return "CreateLinearGradient";
	// Note: this must be one of the operation names returned by GetOperations.
}

/*
*/
std::vector<std::string> _SnpRunnableGradient::GetCategories() const
{
	// TODO - Modify the category below to the most relevant category for this snippet
	// or specify a new category.
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Color Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableGradient::CanRun(SnpRunnable::Context& runnableContext)
{
	SnpSelectionHelper selectionHelper;
	if ("CreateLinearGradient" == runnableContext.GetOperation().GetName() || 
		"CreateRadialGradient" == runnableContext.GetOperation().GetName() ||
		"EditGradientName" == runnableContext.GetOperation().GetName() || 
		"EditGradientStop" == runnableContext.GetOperation().GetName() ||
		"DisplayGradient" == runnableContext.GetOperation().GetName() ||  
		"DeleteGradient" == runnableContext.GetOperation().GetName())
		return selectionHelper.IsDocumentSelected();
	else
		return true;
	// Note: this will control trhe enabling and disabling of the Run button in the SnippetRunner panel
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableGradient::Run(SnpRunnable::Context& runnableContext)
{
	// TODO edit as required by your snippet - you don't have to have a Do method in your class.
	ASErr result = kNoErr;
	SnpGradient instance;
	if ("CreateLinearGradient" == runnableContext.GetOperation().GetName()) {
		result = instance.CreateLinearGradient();
	}
	else if ("CreateRadialGradient" == runnableContext.GetOperation().GetName()) {
		result = instance.CreateRadialGradient();
	}
	else if ("EditGradientName" == runnableContext.GetOperation().GetName()) {
		result = instance.EditGradientName();
	}
	else if ("EditGradientStop" == runnableContext.GetOperation().GetName()) {
		result = instance.EditGradientStop();
	}
	else if ("DisplayGradient" == runnableContext.GetOperation().GetName()) {
		result = instance.DisplayGradient();
	}
	else if ("DeleteGradient" == runnableContext.GetOperation().GetName()) {
		result = instance.DeleteGradient();
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableGradient instance_SnpRunnableGradient;

// End SnpGradient.cpp

