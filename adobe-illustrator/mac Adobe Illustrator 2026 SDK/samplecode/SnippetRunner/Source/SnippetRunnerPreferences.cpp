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

#include "SnippetRunnerID.h"
#include "SnippetRunnerSuites.h"
#include "SnippetRunnerPreferences.h"
#include "SnippetRunnerParameter.h"
#include "SDKErrors.h"

const char* kSnippetRunnerPreferencePrefix = "sdk/" kSnippetRunnerPluginName "/";
const char* kSnippetRunnerAssetsFolderPath = "kSnippetRunnerAssetsFolderPath";
const char* kSnippetRunnerParameterInputMode = "kSnippetRunnerParameterInputMode";
const char* kSnippetRunnerUnitTest = "kSnippetRunnerUnitTest";
SnippetRunnerPreferences* SnippetRunnerPreferences::fInstance = NULL;

/*
*/
SnippetRunnerPreferences::SnippetRunnerPreferences() : fParameterInputMode(SnippetRunnerParameter::kDefaultParameterInputMode), fUnitTest(false)
{
	ASErr status = sAIPreference->GetFilePathSpecificationPreference(kSnippetRunnerPreferencePrefix, 
		kSnippetRunnerAssetsFolderPath, 
		fAssetsFolderPath);
	SDK_ASSERT_NOTHROW(!status);
	status = sAIPreference->GetIntegerPreference(kSnippetRunnerPreferencePrefix, 
		kSnippetRunnerParameterInputMode, 
		&fParameterInputMode);
	SDK_ASSERT_NOTHROW(!status);
	status = sAIPreference->GetBooleanPreference(kSnippetRunnerPreferencePrefix, 
		kSnippetRunnerUnitTest, 
		&fUnitTest);
	SDK_ASSERT_NOTHROW(!status);
}

/*
*/
SnippetRunnerPreferences::~SnippetRunnerPreferences()
{
	ASErr status = kNoErr;
	if (!fAssetsFolderPath.IsEmpty()) {
		status = sAIPreference->PutFilePathSpecificationPreference(kSnippetRunnerPreferencePrefix, 
			kSnippetRunnerAssetsFolderPath,
			fAssetsFolderPath);
		SDK_ASSERT_NOTHROW(!status);
	}
	status = sAIPreference->PutIntegerPreference(kSnippetRunnerPreferencePrefix, 
		kSnippetRunnerParameterInputMode, 
		fParameterInputMode);
	SDK_ASSERT_NOTHROW(!status);
	status = sAIPreference->PutBooleanPreference(kSnippetRunnerPreferencePrefix, 
		kSnippetRunnerUnitTest, 
		fUnitTest);
	SDK_ASSERT_NOTHROW(!status);
}

/*
*/
SnippetRunnerPreferences* SnippetRunnerPreferences::Instance()
{
	if (NULL == SnippetRunnerPreferences::fInstance) {
		SnippetRunnerPreferences::fInstance = new SnippetRunnerPreferences();
	}
	return SnippetRunnerPreferences::fInstance;
}

/*
*/
void SnippetRunnerPreferences::DeleteInstance()
{
	if (NULL != SnippetRunnerPreferences::fInstance) {
		delete SnippetRunnerPreferences::fInstance;
		SnippetRunnerPreferences::fInstance = NULL;
	}
}

/*
*/
ai::FilePath SnippetRunnerPreferences::GetAssetsFolderPath()
{
	bool resolveLinks = false;
	if (!this->fAssetsFolderPath.Exists(resolveLinks)) {
		this->PromptAssetsFolderPath();
	}
	return this->fAssetsFolderPath;
}

/*
*/
void SnippetRunnerPreferences::PromptAssetsFolderPath()
{
	SnippetRunnerParameter::ModePreserver mode(SnippetRunnerParameter::kPromptParameterInputMode);
	this->fAssetsFolderPath = SnippetRunnerParameter::Instance()->GetFolderPath("Locate the Assets folder (<SDK>/samplecode/CodeSnippets/examplefiles)");
}

// End SnippetRunnerPreferences.cpp
