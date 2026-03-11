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

/** Shows how to use the standard progress bar.

	@see AIUserSuite
  	@ingroup sdk_snippet
 */
class SnpProgressBar
{
	public:
		/** Displays a progress bar with a fixed textual description.
		*/
		ASErr BasicProgressBar();

		/** Displays a progress bar with a textual description that
			varies as the task progresses.
		*/
		ASErr ComplexProgressBar();

	private:
		/** Calls AIUserSuite::UpdateProgress
			@param current see AIUserSuite::UpdateProgress
			@param max see AIUserSuite::UpdateProgress
		*/
		void UpdateProgress(ai::int32 current, ai::int32 max);

		/** Calls AIUserSuite::SetProgressText
		*/
		void SetProgressText(const ai::UnicodeString& msg);

		/** Calls AIUserSuite::Cancel
		*/
		ASBoolean IsCancelled();

		/** Calls AUUserSuite::CloseProgress
		*/
		void CloseProgress();

		/** Dummy task.
			@param i
			@param j
			@return i * j
		*/
		ai::int32 MyTask(ai::int32 i, ai::int32 j);
};

/*
*/
ASErr SnpProgressBar::BasicProgressBar()
{
	ASErr result = kNoErr;

	// Set text of progress bar.
	this->SetProgressText(ai::UnicodeString("BasicProgressBar"));

	// Loop performing long running set of tasks.
	const ai::int32 max = 10000;
	for (ai::int32 taskNumber = 0; taskNumber < max; taskNumber++)
	{
		// Perform task
		for (ai::int32 j = 0; j < 100000; j++)
		{
			this->MyTask(taskNumber, j);
		}

		// Update progress.
		this->UpdateProgress(taskNumber + 1, max);

		// Check for cancellation.
		if (this->IsCancelled()) {
			result = kCanceledErr;
			break;
		}
	}

	// Close progress bar.
	this->CloseProgress();

	return result;
}

/*
*/
ASErr SnpProgressBar::ComplexProgressBar()
{
	ASErr result = kNoErr;
	const ai::int32 max = 10000;
	ai::UnicodeString maxString;
	ai::NumberFormat().toString(max, maxString);
	const ai::UnicodeString progressText1("ComplexProgressBar: Doing This");
	const ai::UnicodeString progressText2("ComplexProgressBar: Doing That");
	this->SetProgressText(progressText1);
	for (ai::int32 taskNumber = 0; taskNumber < max; taskNumber++) {

		// Perform task
		for (ai::int32 j = 0; j < 100000; j++)
		{
			this->MyTask(taskNumber, j);
		}

		// Update progress text message every so often
		if (taskNumber % 1200 == 0)
			this->SetProgressText(progressText2);
		else if (taskNumber % 2500 == 0)
			this->SetProgressText(progressText1);

		// Update progress.
		this->UpdateProgress(taskNumber + 1, max);

		// Check for cancellation.
		if (this->IsCancelled())
		{
			result = kCanceledErr;
			break;
		}
	}

	return result;
}

/*
*/
void SnpProgressBar::CloseProgress()
{
	if (sAIUser)
		sAIUser->CloseProgress();
}

/*
*/
void SnpProgressBar::UpdateProgress(ai::int32 current, ai::int32 max)
{
	if (sAIUser)
		sAIUser->UpdateProgress(current, max);
}

/*
*/
void SnpProgressBar::SetProgressText(const ai::UnicodeString& msg)
{
	if (sAIUser)
		sAIUser->SetProgressText(msg);
}

/*
*/
ASBoolean SnpProgressBar::IsCancelled()
{
	ASBoolean result = false;

	if (sAIUser)
		result = sAIUser->Cancel();

	return result;
}

/*
*/
ai::int32 SnpProgressBar::MyTask(ai::int32 i, ai::int32 j)
{
	return i * j;
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpProgressBar available to the SnippetRunner framework.
 */
class _SnpRunnableProgressBar : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableProgressBar () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableProgressBar () {}

		/* Returns name of snippet.
		*/
		std::string GetName() const;

		/* Returns a description of what the snippet does.
		*/
		std::string	GetDescription() const;

		/** Returns operations supported by this snippet.
		*/
		Operations GetOperations() const;

		/** Returns name of this snippet's default operation.
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
std::string _SnpRunnableProgressBar::GetName() const
{
	return "ProgressBar";
}

/*
*/
std::string _SnpRunnableProgressBar::GetDescription() const 
{
	return "Shows how to use the standard progress bar";
}

/*
*/
SnpRunnable::Operations _SnpRunnableProgressBar::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("BasicProgressBar", "", kSnpRunNotSupportedContext)); // skip unit test - interactive feature
	operations.push_back(Operation("ComplexProgressBar", "", kSnpRunNotSupportedContext)); // skip unit test - interactive feature
	return operations;
}

/*
*/
std::string _SnpRunnableProgressBar::GetDefaultOperationName() const 
{
	return "BasicProgressBar";
}

/*
*/
std::vector<std::string> _SnpRunnableProgressBar::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Helper Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableProgressBar::CanRun(SnpRunnable::Context& runnableContext)
{
	return true;
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableProgressBar::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpProgressBar instance;
	if ("BasicProgressBar" == runnableContext.GetOperation().GetName()) {
		result = instance.BasicProgressBar();
	}
	else if ("ComplexProgressBar" == runnableContext.GetOperation().GetName()) {
		result = instance.ComplexProgressBar();
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableProgressBar instance_SnpRunnableProgressBar;

// End SnpProgressBar.cpp

