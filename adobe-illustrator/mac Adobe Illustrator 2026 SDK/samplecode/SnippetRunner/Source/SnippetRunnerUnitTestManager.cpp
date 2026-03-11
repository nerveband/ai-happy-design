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

// General includes:
#include <vector>

#include "SDKErrors.h"
#include "SnippetRunnerID.h"
#include "SnippetRunnerSuites.h"
#include "SnippetRunnerParameter.h"
#include "SnippetRunnerLog.h"
#include "SnippetRunnerManager.h"
#include "SnpDocumentHelper.h"
#include "SnpRunnable.h"
#include "SnippetRunnerPlugin.h"
#include "SnippetRunnerUnitTestManager.h"

extern SnippetRunnerPlugin*	gPlugin;

SnippetRunnerUnitTestManager* SnippetRunnerUnitTestManager::fInstance = NULL;

/*
*/
SnippetRunnerUnitTestManager::SnippetRunnerUnitTestManager() : fUnitTestTimer(NULL), fUnitTestState(kSnpRunUnitTestIdle)
{
}

/*
*/
SnippetRunnerUnitTestManager::~SnippetRunnerUnitTestManager()
{
}

/*
*/
SnippetRunnerUnitTestManager* SnippetRunnerUnitTestManager::Instance()
{
	if (NULL == SnippetRunnerUnitTestManager::fInstance) {
		SnippetRunnerUnitTestManager::fInstance = new SnippetRunnerUnitTestManager();
	}
	return SnippetRunnerUnitTestManager::fInstance;
}

/*
*/
void SnippetRunnerUnitTestManager::DeleteInstance()
{
	if (NULL != SnippetRunnerUnitTestManager::fInstance) {
		delete SnippetRunnerUnitTestManager::fInstance;
		SnippetRunnerUnitTestManager::fInstance = NULL;
	}
}

/*
*/
ASErr SnippetRunnerUnitTestManager::RunUnitTest(const std::vector<std::string>& snippetNames, const std::string& operationName)
{
	ASErr result = kNoErr;

	try {
		SDK_ASSERT_MSG(this->fUnitTestState == kSnpRunUnitTestIdle, "Unit test already running");
		SDK_ASSERT(snippetNames.size());

		// Get the current parameter input mode for subsequent restoration.
		this->fParameterInputMode = SnippetRunnerParameter::Instance()->GetParameterInputMode();

		// Get user to confirm that unit test is to be run.
		SnippetRunnerParameter::Instance()->SetParameterInputMode(SnippetRunnerParameter::kPromptParameterInputMode);
		SnippetRunnerParameter::Instance()->Alert("Begin Unit Test (hold Esc key to cancel the test run)?");
					
		// Use defaults for unit testing.
		SnippetRunnerParameter::Instance()->SetParameterInputMode(SnippetRunnerParameter::kDefaultParameterInputMode);

		// Close all open documents before starting the test.
		SnpDocumentHelper docHelper;
		result = docHelper.CloseAllDocuments();
		aisdk::check_ai_error(result);

		// Initialise unit test queue.
		this->fUnitTestState = kSnpRunUnitTestIdle;
		this->fUnitTests.clear();

		// Are we testing one operation or many?
		if (operationName.size()) {
			// Queue up one operation.
			this->fUnitTests.push_back(UnitTest(snippetNames[0], operationName));
		}
		else {
			// Queue up all operations in the given set of snippets.
			for (std::vector<std::string>::const_iterator iter = snippetNames.begin(); iter != snippetNames.end(); iter++) {
				SnpRunnable* snp = SnippetRunnerManager::Instance().FindSnippet(*iter);
				if (!snp) {
					SDK_ASSERT_MSG_NOTHROW(false, "Couldn't find snippet");
					continue;
				}

				// Get the current snippet's operations.
				SnpRunnable::Operations operations = snp->GetOperations();
				for (SnpRunnable::Operations::iterator opsIter = operations.begin(); opsIter != operations.end(); opsIter++) {
					this->fUnitTests.push_back(std::pair<std::string, std::string>(*iter, opsIter->GetName()));
				}
			} // end for
		}

		// Create the timer if necessary.
		if (!this->fUnitTestTimer) {
            char name[] = "SnippetRunnerUnitTestManager";
			result = sAITimer->AddTimer(gPlugin->GetPluginRef(), name, 10, &this->fUnitTestTimer);
			SDK_ASSERT_MSG_NOTHROW(!result && this->fUnitTestTimer, "Couldn't start unit test AddTimer call failed");
		}

		// Activate the timer.
		if (this->fUnitTestTimer) {
			result = sAITimer->SetTimerActive(this->fUnitTestTimer, true);
			if (result) {
				SDK_ASSERT_MSG_NOTHROW(!result, "Couldn't start unit test SetTimerActive call failed");
			}
		}
		if (!result) {
			this->fUnitTestState = kSnpRunUnitTestQueued;
		}
		else {
			this->fUnitTestState = kSnpRunUnitTestIdle;
			this->fUnitTests.clear();
		}

	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
void SnippetRunnerUnitTestManager::GoTimer(AITimerMessage* message)
{
	ASErr result = kNoErr;
	SnippetRunnerLog* log = SnippetRunnerLog::Instance();
	try {
		// Check that we expected to be called.
		if (this->fUnitTestState == kSnpRunUnitTestIdle) {
			if (this->fUnitTestTimer) {
				sAITimer->SetTimerActive(this->fUnitTestTimer, false);
			}
			SDK_ASSERT_MSG(false, "Unexpected timer call."); // will throw		
		}

		// TODO Verify this timer message is ours.

		// Handle starting of unit test run.
		if (this->fUnitTestState == kSnpRunUnitTestQueued) {
			this->fUnitTestFailedOperations.clear();
			log->Write("Begin Unit Test");
			log->BeginIndent();
			this->fUnitTestState = kSnpRunUnitTestRunning;
		}

		// Check for cancellation.
		ASBoolean isCancelled = false;
/*		ADMVirtualKeyArray keys;
		if (sADMTracker->GetKeysState(&keys)) {
			if (keys [kADMEscapeKey]) {
				// Cancel 
				log->Write("UNIT TEST CANCELLED");
				isCancelled = true;
			}
		}*/

		// Unit test the next queued snippet operation.
		if (!isCancelled && this->fUnitTests.size()) {
			// Get the next operation from the queue.
			UnitTests::iterator iter = this->fUnitTests.begin();

			// Execute the operation.
			ASErr status = kNoErr;
			try {
				status = SnippetRunnerManager::Instance().Execute(iter->first, iter->second);
			}
			catch (ai::Error& ex) {
				status = ex;
			}
			catch (...) {
				status = kCantHappenErr; // unhandled exception
			}
			if (status) {
				this->fUnitTestFailedOperations.push_back(iter->first + "::" + iter->second);
			}

			// Remove the operation from the queue.
			this->fUnitTests.erase(iter);
		}

		// Handle end of the unit test.
		if (isCancelled || !this->fUnitTests.size()) {

			// End of the unit test.
			this->fUnitTestState = kSnpRunUnitTestIdle;

			// Turn off the timer.
			if (this->fUnitTestTimer) {
				sAITimer->SetTimerActive(this->fUnitTestTimer, false);
			}

			// Restore the parameter input mode.
			SnippetRunnerParameter::Instance()->SetParameterInputMode(this->fParameterInputMode);

			// Update the log.
			if (this->fUnitTestFailedOperations.size()) {
				log->Write("UNIT TEST FAILURE SUMMARY:");
				{
					SnippetRunnerLog::Indent indent;
					for (std::vector<std::string>::const_iterator iter = this->fUnitTestFailedOperations.begin(); iter != this->fUnitTestFailedOperations.end(); iter++) {
						log->Write(*iter);
					}
				}
			}
			log->EndIndent();
			log->Write("End Unit Test");
			log->NotifyLogChanged();
			this->NotifyEndUnitTest();
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
}

/*
*/
void SnippetRunnerUnitTestManager::NotifyEndUnitTest()
{
	if (gPlugin) {
		gPlugin->NotifyEndUnitTest();
	}
}

// End SnippetRunnerUnitTestManager.cpp
