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

#ifndef _SnippetRunnerManager_
#define _SnippetRunnerManager_

// General includes:
#include <vector>

// Forward declarations:
#include "SnpRunnable.h"

/** Registers and provides access to code snippets (implemented as a Meyers singleton).

	@see SnippetRunnerUnitTestManager
	@ingroup snippetrunner
*/
class SnippetRunnerManager
{
public:
	/** Returns the global instance of this singleton, call like this SnippetRunnerManager::Instance().
		@return the global instance of this singleton class.
	*/
	static SnippetRunnerManager& Instance();

	/** Adds snippet to list of registered snippets. If a snippet is
		already registered under this name the request will be ignored.
		@param runnable snippet to be registered.
	*/
	void AddSnippet(SnpRunnable* runnable);

	/** Removes snippet from list of registered snippets.
		@param runnable snippet to be removed.
	*/
	void RemoveSnippet(SnpRunnable* runnable);

	/** Looks up the named snippet in the list of registered snippets.
		@param snippetName IN name of the snippet.
		@return the snippet.
	*/
	SnpRunnable* FindSnippet(const std::string& snippetName) const;

	/** Collection of registered snippets.
	*/
	typedef	vector<SnpRunnable*> Snippets;

	/** Returns a collection of registered snippets.
		@param category IN the category of snippets to return.
		@return a collection of registered snippets.
	*/
	Snippets GetSnippets(const std::string category = "All Snippets") const;

	/** Returns context for a snippet run.
		@param snippetName IN name of the snippet.
		@param operationName IN name of the snippets operation.
		@return the required context.
	*/
	SnpRunnable::Context GetContext(const std::string& snippetName, const std::string& operationName) const;

	/** Returns true if the snippet can run, false otherwise.
		@param snippetName IN name of the snippet.
		@param operationName IN name of the snippets operation.
		@return true if snippet can run, false otherwise.
	*/
	ASBoolean CanRun(const std::string& snippetName, const std::string& operationName) const;

	/** Runs the given snippet.
		@param snippetName IN name of the snippet.
		@param operationName IN name of the operation.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr Run(const std::string& snippetName, const std::string& operationName) const;

	/** Checks if unit test is running.
		@return true if unit test is running, false otherwise.
	*/
	ASBoolean IsUnitTestRunning() const;

	/** Queues unit test of snippet operations using SnippetRunnerUnitTestManager.
		@param snippetNames IN vector of snippet names to run.
		@param operationName IN tests all operations of all named snippets if string is empty, tests given operation otherwise.
		@return kNoErr on success, other ASErr otherwise.
		@see SnippetRunnerUnitTestManager
	*/
	ASErr RunUnitTest(const std::vector<std::string>& snippetNames, const std::string& operationName);

	/** Writes header to the log, calls SnpRunnable to
		run the snippet and writes footer to the log to indicate the result.
		@param snippetName IN name of the snippet.
		@param operationName IN name of the operation.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr Execute(const std::string& snippetName, const std::string& operationName) const;

private:
	/** Constructor is hidden from public because this is a singleton.
	*/
	SnippetRunnerManager();

	/** Copy constructor is hidden.
		@param rhs
	*/
	SnippetRunnerManager(SnippetRunnerManager const&);

	/** Assignment operator is hidden.
		@param IN SnippetRunnerManager object to assign to lhs.
		@return copy of SnippetRunnerManager passed in.
	*/
	SnippetRunnerManager& operator=(SnippetRunnerManager const&);

	/** Destructor is hidden.
	*/
	~SnippetRunnerManager(); // destructor hidden
};

#endif // _SnippetRunnerManager_

// End SnippetRunnerManager.h
