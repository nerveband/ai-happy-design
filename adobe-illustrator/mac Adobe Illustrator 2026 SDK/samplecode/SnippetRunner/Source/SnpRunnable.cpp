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

// Interface includes:

// General includes:
#include "SDKErrors.h"

// Project includes:
#include "SnpRunnable.h"
#include "SnippetRunnerManager.h"
#include "SnippetRunnerLog.h"
#include "SnpUnitTestHelper.h"


const char* kEmptyString = "";

/* Constructor register's snippet with the framework.
*/
SnpRunnable::SnpRunnable()
{
	SnippetRunnerManager& snipRunnerManager = SnippetRunnerManager::Instance();
	snipRunnerManager.AddSnippet(this);
}

/* Destructor removes snippet from the framework.
*/
SnpRunnable::~SnpRunnable()
{
	SnippetRunnerManager& snipRunnerManager = SnippetRunnerManager::Instance();
	snipRunnerManager.RemoveSnippet(this);
}

/* Default unit test.
*/
ASErr SnpRunnable::RunUnitTest(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	if (runnableContext.GetOperation().GetUnitTestContext() != SnpRunnable::kSnpRunNotSupportedContext) {
		result = this->SetupContext(runnableContext);
		if (!result) {
			result = this->Run(runnableContext);
			ASErr tmp = this->TeardownContext(runnableContext);
			if (tmp) {
				SnippetRunnerLog::Instance()->WritePrintf("FAILED: SnpRunnable::TeardownContext error=%x", tmp);
				if (!result) {
					// Return the teardown error to the caller if no errors have occurred yet.
					result = tmp;
				}
			}
		}
		else {
			SnippetRunnerLog::Instance()->WritePrintf("FAILED: SnpRunnable::SetupContext error=%x", result);
		}
	}
	else {
		SnippetRunnerLog::Instance()->Write("SKIPPED (kSnpRunNotSupportedContext)");
	}
	return result;
}

/* Default categories.
*/
std::vector<std::string>	SnpRunnable::GetCategories() const
{
	std::vector<std::string> categories;
	categories.push_back("All Snippets");
	return categories;
}

/*
*/
ASErr SnpRunnable::SetupContext(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	try {
		SnpUnitTestHelper unitTestHelper;
		result = unitTestHelper.SetupContext(runnableContext, this);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpRunnable::TeardownContext(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	try {
		SnpUnitTestHelper unitTestHelper;
		result = unitTestHelper.TeardownContext(runnableContext, this);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

// End SnpRunnable.cpp

