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

#include "IAIStringFormatUtils.h"

// Project includes
#include "SnippetRunnerID.h"
#include "SnippetRunnerPromptDialog.h"
#include "SnippetRunnerPlugin.h"
#include "SnippetRunnerSuites.h"

extern SnippetRunnerPlugin*	gPlugin;

AIReal SnippetRunnerPromptDialog::fUpperLimit;
AIReal SnippetRunnerPromptDialog::fLowerLimit;

/*
 * Called by GetInputFromUser will check that the input is between the specified ranges.
*/
AIBoolean SnippetRunnerPromptDialog::ValidNumber(AIValidationStage validationStage, const ai::UnicodeString &input, ai::UnicodeString &errorStf, void * userdata) 
{ 
    AIBoolean result = false;
    AIReal tmp = 0.0;

    if (ai::NumberFormat().parseString(input, tmp)) 
    {
		//let 0, 0 be any limit.
        if (fLowerLimit == 0.0 && fUpperLimit == 0.0)  
        {
            result = true;
        }
        else if (tmp >= fLowerLimit && tmp <= fUpperLimit) 
        {
            result = true;
        }
    }

	if(!result){
		stringstream errorMessageStream;
		errorMessageStream	<< "Please enter a number between "
							<< fLowerLimit << " and " << fUpperLimit;
		ai::UnicodeString errorMessage = ai::UnicodeString(errorMessageStream.str());
		errorStf = errorMessage;
	}

    return result;
}

SnippetRunnerPromptDialog::SnippetRunnerPromptDialog()
{
}

SnippetRunnerPromptDialog::~SnippetRunnerPromptDialog( )
{
}

/** Prompts the user for a string.
	@param value OUT the string entered by the user.
	@param prompt IN string prompt to display in the dialog.
	@return true if user pressed OK button, false otherwise.
*/
ASBoolean SnippetRunnerPromptDialog::GetString(ai::UnicodeString& value, const ai::UnicodeString& prompt)
{
	ASBoolean result = false;
    AIErr error = sAIUser->GetInputFromUser(ai::UnicodeString(kSnippetRunnerPromptDialogName),NULL, prompt, value, &AnyString, NULL, -1);
	if(!error) 
    {
        result = true;
	}
	return result;
}

/** Prompts the user for a real value.
	@param value OUT the value entered by the user.
	@param prompt IN string prompt to display in the dialog.
	@param lowerLimit IN the smallest value that can be entered.
	@param upperLimit IN the largest value that can be entered.
	@return true if user pressed OK button, false otherwise.
*/
ASBoolean SnippetRunnerPromptDialog::GetReal(AIReal& value, const ai::UnicodeString& prompt, const AIReal lowerLimit, const AIReal upperLimit)
{
	ASBoolean result = false;
    SnippetRunnerPromptDialog::fLowerLimit = lowerLimit;
    SnippetRunnerPromptDialog::fUpperLimit = upperLimit;
    ai::UnicodeString stringValue;
    ai::NumberFormat().toString(value,5,stringValue);
    AIErr error = sAIUser->GetInputFromUser(ai::UnicodeString(kSnippetRunnerPromptDialogName),NULL, prompt, stringValue, &ValidNumber, NULL, -1);
	if (!error) {
        result = true;
        ai::NumberFormat().parseString(stringValue, value);
	}
	return result;
}

/*
	Prompts the user for an integer value between fLowerLimit and fUpperLimit
	@param value OUT the value entered by the user.
	@message the message to show above the input box
	@param prompt IN string prompt to display in the dialog.
*/
AIBoolean SnippetRunnerPromptDialog::PromptForInt(ai::UnicodeString* message, ai::UnicodeString prompt, ASInt32& value)
{
    AIBoolean result = false;
    ai::UnicodeString stringValue;
    ai::NumberFormat().toString(value,stringValue);
    AIErr error = sAIUser->GetInputFromUser(ai::UnicodeString(kSnippetRunnerPromptDialogName),message, prompt, stringValue, &ValidNumber, NULL, -1);
	if (!error) {
        result = true;
        ai::NumberFormat().parseString(stringValue, value);
	}
    return result;
}

/** Prompts the user for an integer value.
	@param value OUT the value entered by the user.
	@param prompt IN string prompt to display in the dialog.
	@param lowerLimit IN the smallest value that can be entered.
	@param upperLimit IN the largest value that can be entered.
	@return true if user pressed OK button, false otherwise.
*/
ASBoolean SnippetRunnerPromptDialog::GetInt(ASInt32& value, const ai::UnicodeString& prompt, const ASInt32 lowerLimit, const ASInt32 upperLimit)
{
    SnippetRunnerPromptDialog::fLowerLimit = static_cast<AIReal>(lowerLimit);
    SnippetRunnerPromptDialog::fUpperLimit = static_cast<AIReal>(upperLimit);
   	return PromptForInt(NULL,prompt,value);
}


/** Prompts the user to choose from a list of string options.
	@param value OUT the users choice.
	@param prompt IN string prompt to display in the dialog.
	@param choices IN vector of strings to represent choices.
	@return true if user pressed OK button, false otherwise.
*/
AIBoolean SnippetRunnerPromptDialog::GetChoice( ai::int32& value, const ai::UnicodeString& prompt, const std::vector<ai::UnicodeString>& choices)
{
    SnippetRunnerPromptDialog::fLowerLimit = static_cast<AIReal>(0);
    SnippetRunnerPromptDialog::fUpperLimit = static_cast<AIReal>((choices.size()-1));
    int i = 0;
    stringstream choiceList;
    for( std::vector<ai::UnicodeString>::const_iterator n = choices.begin(); n != choices.end(); ++n)
    {
        choiceList << i << ": " << n->as_UTF8() << "\t";
        i++;
    }
    
    ai::UnicodeString stringValue;
    ai::UnicodeString message = ai::UnicodeString(choiceList.str());
    
	return PromptForInt(&message, prompt,value);
}

/** Displays a message to the user.
	@param prompt IN string message.
	@return true if user pressed OK button, false otherwise.
*/
ASBoolean SnippetRunnerPromptDialog::Alert(const ai::UnicodeString& prompt)
{
	return sAIUser->OKCancelAlert(prompt,true,NULL);
}
