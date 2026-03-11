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

#ifndef __SnippetRunnerPromptDialog_h__
#define __SnippetRunnerPromptDialog_h__

/**	Provides dialog which prompts for parameter input.

	@ingroup snippetrunner
*/
#include "IAIStringFormatUtils.h"


class SnippetRunnerPromptDialog
{
public:
	/**	Constructor.
	*/
	SnippetRunnerPromptDialog();

	/**	Destructor.
	*/
	virtual ~SnippetRunnerPromptDialog();
    
	/** Prompts the user for a string.
		@param value OUT the string entered by the user.
		@param prompt IN string prompt to display in the dialog.
		@return true if user pressed OK button, false otherwise.
	*/
	ASBoolean GetString(ai::UnicodeString& value, const ai::UnicodeString& prompt);

	/** Prompts the user for a real value.
		@param value OUT the value entered by the user.
		@param prompt IN string prompt to display in the dialog.
		@param lowerLimit IN the smallest value that can be entered.
		@param upperLimit IN the largest value that can be entered.
		@return true if user pressed OK button, false otherwise.
	*/
	ASBoolean GetReal(AIReal& value, const ai::UnicodeString& prompt, const AIReal lowerLimit=0.0, const AIReal upperLimit=0.0);

	/** Prompts the user for an integer value.
		@param value OUT the value entered by the user.
		@param prompt IN string prompt to display in the dialog.
		@param lowerLimit IN the smallest value that can be entered.
		@param upperLimit IN the largest value that can be entered.
		@return true if user pressed OK button, false otherwise.
	*/
	ASBoolean GetInt(ASInt32& value, const ai::UnicodeString& prompt, const ASInt32 lowerLimit=0, const ASInt32 upperLimit=0);

	/** Prompts the user to choose from a list of string options.
		@param value OUT the users choice.
		@param prompt IN string prompt to display in the dialog.
		@param choices IN vector of strings to represent choices.
		@return true if user pressed OK button, false otherwise.
	*/
	AIBoolean GetChoice(ai::int32& value, const ai::UnicodeString& prompt, const std::vector<ai::UnicodeString>& choices);

	/** Displays a message to the user.
		@param prompt IN string message.
		@return true if user pressed OK button, false otherwise.
	*/
	ASBoolean Alert(const ai::UnicodeString& prompt);
    
private:
    
	/*
	* validation callback for GetInputFromUser will make sure that the number input is between fLowerLimit and fUpperLimit.
	*/
    static AIBoolean ValidNumber(AIValidationStage validationStage, const ai::UnicodeString &input, ai::UnicodeString &errorStf, void * userdata); 
  
	/*
	* validation callback for GetInputFromUser will let any string be entered.
	*/
	static AIBoolean AnyString (AIValidationStage validationStage, const ai::UnicodeString &input, ai::UnicodeString &errorStf, void * userdata){ return true; };

	/*
	* Brings up a modal dialog that accepts int values
	*/
    AIBoolean PromptForInt(ai::UnicodeString* message, ai::UnicodeString prompt, ASInt32& value);
    
    static AIReal fLowerLimit;
    
    static AIReal fUpperLimit;
};

#endif // __SnippetRunnerPromptDialog_h__

// End SnippetRunnerPromptDialog.h
