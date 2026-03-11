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

#ifndef _SnippetRunnerUnitTestManager_
#define _SnippetRunnerUnitTestManager_

#include "SnippetRunnerParameter.h"

/** Manages unit testing - snippet operations are unit tested
	via this class - the tests are run by a timer.

	@ingroup snippetrunner
*/
class SnippetRunnerUnitTestManager 
{

public:
	/** Returns the global instance of this singleton.
		@return the global instance of this singleton class.
	*/
	static SnippetRunnerUnitTestManager* Instance();

	/**	Deletes the singleton instance of this class.
	*/
	static void DeleteInstance();

	/** Queues unit test of snippet operations.
		@param snippetNames IN all snippet names to run in unit test.
		@param operationName IN tests all operations of all named snippets if string is empty, tests given operation otherwise.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr RunUnitTest(const std::vector<std::string>& snippetNames, const std::string& operationName);

	/** Unit tests one snippet operation per timer call.
		@param message contents of the timer message sent by kSelectorAIGoTimer.
	*/
	void GoTimer(AITimerMessage* message);

	/** Unit test manager states.
	*/
	typedef enum {
		kSnpRunUnitTestIdle,
		kSnpRunUnitTestQueued,
		kSnpRunUnitTestRunning
	} eSnpRunUnitTestState;

	/** Gets the state of unit test manager.
		@return unit test manager state
	*/
	SnippetRunnerUnitTestManager::eSnpRunUnitTestState GetUnitTestState() {return fUnitTestState;}

	/** Sends message that unit test has ended.
	*/
	void NotifyEndUnitTest();

private:
	/** Constructor is hidden - instantiate using SnippetRunnerUnitTestManager::Instance.
	*/
	SnippetRunnerUnitTestManager();

	/** Destructor is hidden - destroy using SnippetRunnerUnitTestManager::DeleteInstance.
	*/
	~SnippetRunnerUnitTestManager();

	/** Copy constructor is hidden.
		@param rhs
	*/
	SnippetRunnerUnitTestManager(SnippetRunnerUnitTestManager const&);

	/** Assignment operator is hidden.
		@param IN SnippetRunnerUnitTestManager object to assign to lhs.
		@return copy of SnippetRunnerUnitTestManager passed in.
	*/
	SnippetRunnerUnitTestManager& operator=(SnippetRunnerUnitTestManager const&);  // assign op is hidden

private:
	/** Points to the singleton instance of this class.
	*/
	static SnippetRunnerUnitTestManager* fInstance;

	/** Stores current state of unit testing.
	*/
	eSnpRunUnitTestState fUnitTestState;

	/** Describes a unit test, the first member of the pair is the name of a snippet
		and the second member is the name of the operation to be tested.	
	*/
	typedef std::pair<std::string, std::string> UnitTest;

	/** Container for unit tests.
	*/
	typedef std::vector<UnitTest> UnitTests;

	/** Stores current queue of operations to be tested.
	*/
	UnitTests fUnitTests;

	/** Stores handle of timer used to run unit tests.
	*/
	AITimerHandle fUnitTestTimer;

	/** Stores parameter input mode to be restored at end of unit test.
	*/
	SnippetRunnerParameter::ParameterInputMode fParameterInputMode;

	/** Stores operations that failed during a unit test run.
	*/
	std::vector<std::string> fUnitTestFailedOperations;
};

#endif //_SnippetRunnerUnitTestManager_

// End SnippetRunnerUnitTestManager.h
