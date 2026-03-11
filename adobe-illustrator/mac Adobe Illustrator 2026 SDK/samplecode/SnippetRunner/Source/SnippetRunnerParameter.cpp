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
#include "IAIUnicodeString.h"
#include "SnippetRunnerParameter.h"
#include "SnippetRunnerTokenisor.h"
#include "SnippetRunnerPromptDialog.h"
#include "SnippetRunnerSuites.h"
#include "SnippetRunnerLog.h"
#include "SnippetRunnerPreferences.h"

SnippetRunnerParameter* SnippetRunnerParameter::fInstance = NULL;

/*
*/
SnippetRunnerParameter::SnippetRunnerParameter()
{
}

/*
*/
SnippetRunnerParameter::~SnippetRunnerParameter()
{
}


/*
*/
SnippetRunnerParameter* SnippetRunnerParameter::Instance()
{
	if (NULL == SnippetRunnerParameter::fInstance) {
		SnippetRunnerParameter::fInstance = new SnippetRunnerParameter();
	}
	return SnippetRunnerParameter::fInstance;
}

/*
*/
void SnippetRunnerParameter::DeleteInstance()
{
	if (NULL != SnippetRunnerParameter::fInstance) {
		delete SnippetRunnerParameter::fInstance;
		SnippetRunnerParameter::fInstance = NULL;
	}
}


/*
*/
void SnippetRunnerParameter::SetParameters(const ai::UnicodeString& params)
{
	fTokenisor.SetParameters(params);
}

/*
*/
ai::UnicodeString SnippetRunnerParameter::GetParameters() const
{
	return ai::UnicodeString();
}

/*
*/
ai::UnicodeString SnippetRunnerParameter::GetString(const ai::UnicodeString& prompt, const ai::UnicodeString& defaultValue) const
{
	ai::UnicodeString result = defaultValue;
	if (kCSVParameterInputMode == this->GetParameterInputMode()) {
		if (!fTokenisor.GetString(result)) {
			// No token available.
			throw ai::Error(kBadParameterErr);
		}
	}
	else if (kPromptParameterInputMode == this->GetParameterInputMode()) {
		SnippetRunnerPromptDialog parameterDialog;
		if (!parameterDialog.GetString(result, prompt)) {
			// Prompt cancelled.
			throw ai::Error(kCanceledErr);
		}
	}
	return result;
}

/*
*/
std::string	SnippetRunnerParameter::GetString(const std::string& prompt, const std::string& defaultValue) const
{
	return this->GetString(ai::UnicodeString(prompt), ai::UnicodeString(defaultValue)).getInStdString(kAIPlatformCharacterEncoding);
}

/*
*/
AIReal	SnippetRunnerParameter::GetReal(const ai::UnicodeString& prompt, const AIReal defaultValue, const AIReal lowerLimit, const AIReal upperLimit) const
{
	AIReal result = defaultValue;
	if (kCSVParameterInputMode == this->GetParameterInputMode()) {
		if (!fTokenisor.GetReal(result, lowerLimit, upperLimit)) {
			// No real token available.
			throw ai::Error(kBadParameterErr);
		}
	}
	else if (kPromptParameterInputMode == this->GetParameterInputMode()) {
		SnippetRunnerPromptDialog parameterDialog;
		if (!parameterDialog.GetReal(result, prompt, lowerLimit, upperLimit)) {
			// Prompt cancelled.
			throw ai::Error(kCanceledErr);
		}
	}
	return result;
}

/*
*/
AIReal	SnippetRunnerParameter::GetReal(const std::string& prompt, const AIReal defaultValue, const AIReal lowerLimit, const AIReal upperLimit) const
{
	return this->GetReal(ai::UnicodeString(prompt), defaultValue, lowerLimit, upperLimit);
}


/*
*/
ASInt32 SnippetRunnerParameter::GetInt(const ai::UnicodeString& prompt, const ASInt32 defaultValue, const ASInt32 lowerLimit, const ASInt32 upperLimit) const
{
	ASInt32 result = defaultValue;
	if (kCSVParameterInputMode == this->GetParameterInputMode()) {
		if (!fTokenisor.GetInt(result, lowerLimit, upperLimit)) {
			// No integer token available.
			throw ai::Error(kBadParameterErr);
		}
	}
	else if (kPromptParameterInputMode == this->GetParameterInputMode()) {
		SnippetRunnerPromptDialog parameterDialog;
		if (!parameterDialog.GetInt(result, prompt, lowerLimit, upperLimit)) {
			// Prompt cancelled.
			throw ai::Error(kCanceledErr);
		}
	}
	return result;
}

/*
*/
ASInt32 SnippetRunnerParameter::GetInt(const std::string& prompt, const ASInt32 defaultValue, const ASInt32 lowerLimit, const ASInt32 upperLimit) const
{
	return this->GetInt(ai::UnicodeString(prompt), defaultValue, lowerLimit, upperLimit);
}

/*
*/
ASBoolean SnippetRunnerParameter::GetBoolean(const ai::UnicodeString& prompt) const
{
	if (kPromptParameterInputMode == this->GetParameterInputMode()) {
		return (sAIUser->YesNoAlert(prompt) == kAIAnswer_Yes);
	}
	else
	{
		return false;
	}
}
/*
*/
ASBoolean SnippetRunnerParameter::GetBoolean(const std::string& prompt) const
{
	return this->GetBoolean(ai::UnicodeString(prompt));
}

/*
*/
ai::int32 SnippetRunnerParameter::GetChoice(const ai::UnicodeString& prompt, const std::vector<ai::UnicodeString>& choices, const ai::int32 defaultChoiceIndex) const
{
	ai::int32 result = defaultChoiceIndex;
	if (result >= choices.size() || choices.size() == 0) {
		result = 0;
	}
	// TODO could do with some asserts here
	if (kCSVParameterInputMode == this->GetParameterInputMode()) {
		if (!fTokenisor.GetChoice(result, choices)) {
			// No token available.
			throw ai::Error(kBadParameterErr);
		}
	}
	else if (kPromptParameterInputMode == this->GetParameterInputMode()) {
		SnippetRunnerPromptDialog parameterDialog;
		if (!parameterDialog.GetChoice(result, prompt, choices)) {
			// Prompt cancelled.
			throw ai::Error(kCanceledErr);
		}
	}
	return result;
}

/*
*/
ASUInt32 SnippetRunnerParameter::GetChoice(const std::string& prompt, const std::vector<std::string>& choices, const ASUInt32 defaultChoiceIndex) const
{
	std::vector<ai::UnicodeString> uChoices;
	for (std::vector<std::string>::const_iterator iter = choices.begin(); iter != choices.end(); iter++) {
		uChoices.push_back(ai::UnicodeString(*iter));
	}
	return this->GetChoice(ai::UnicodeString(prompt), uChoices, defaultChoiceIndex);
}

/*
*/
ai::FilePath SnippetRunnerParameter::GetFilePath(const ai::UnicodeString& prompt, ASBoolean isInput, const ai::FilePath defaultValue) const
{
	ai::FilePath result = defaultValue;
	if (kCSVParameterInputMode == this->GetParameterInputMode()) {
		ai::UnicodeString path;
		if (!fTokenisor.GetString(path)) {
			// No token available.
			throw ai::Error(kBadParameterErr);
		}
		result = ai::FilePath(path);
	}
	else if (kPromptParameterInputMode == this->GetParameterInputMode()) {
		AIFileDialogFilters* opts = NULL;
		ASErr status = kNoErr;
		if (isInput) {
			status = sAIUser->GetFileDialog(prompt, opts, result);
		}
		else {
			status = sAIUser->PutFileDialog(prompt, opts, result.GetFileName(), result);
		}
		if (status) {
			throw ai::Error(status);
		}
	}
	return result;
}

/*
*/
ai::FilePath SnippetRunnerParameter::GetFilePath(const std::string& prompt, ASBoolean isInput, const ai::FilePath defaultValue) const
{
	return this->GetFilePath(ai::UnicodeString(prompt), isInput, defaultValue);
}

/*
*/
ai::FilePath SnippetRunnerParameter::GetFolderPath(const ai::UnicodeString& prompt, const ai::FilePath defaultValue) const
{
	ai::FilePath result = defaultValue;
	if (kCSVParameterInputMode == this->GetParameterInputMode()) {
		ai::UnicodeString path;
		if (!fTokenisor.GetString(path)) {
			// No token available.
			throw ai::Error(kBadParameterErr);
		}
		result = ai::FilePath(path);
	}
	else if (kPromptParameterInputMode == this->GetParameterInputMode()) {
		ASErr status = sAIUser->GetDirectoryDialog(prompt, result);
		if (status) {
			throw ai::Error(status);
		}
	}
	return result;
}

/*
*/
ai::FilePath SnippetRunnerParameter::GetFolderPath(const std::string& prompt, const ai::FilePath defaultValue) const
{
	return this->GetFolderPath(ai::UnicodeString(prompt), defaultValue);
}

/*
*/
void SnippetRunnerParameter::Alert(const ai::UnicodeString& prompt) const
{
	if (kPromptParameterInputMode == this->GetParameterInputMode()) {
		SnippetRunnerPromptDialog parameterDialog;
		if (!parameterDialog.Alert(prompt)) {
			// Prompt cancelled.
			throw ai::Error(kCanceledErr);
		}
	}
	else {
		SnippetRunnerLog::Instance()->Write(prompt);
	}
}

/*
*/
void SnippetRunnerParameter::Alert(const std::string& prompt) const
{
	return this->Alert(ai::UnicodeString(prompt));
}

/*
*/
void SnippetRunnerParameter::SetParameterInputMode(ParameterInputMode mode)
{
	SnippetRunnerPreferences::Instance()->SetParameterInputMode(mode);
}

/*
*/
SnippetRunnerParameter::ParameterInputMode SnippetRunnerParameter::GetParameterInputMode() const
{
	return SnippetRunnerPreferences::Instance()->GetParameterInputMode();
}

// End SnippetRunnerParameter.cpp
