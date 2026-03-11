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

/** Demonstrates throwing and catching an ATE Exception.
	@see ATE::Exception
  	@ingroup sdk_snippet
 */
class SnpTextException
{
	public:
		/** Throws an ATE Exception then reports its details to the 
			SnippetRunner log.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr ThrowATEException(void);
};

/*
*/
ASErr SnpTextException::ThrowATEException(void)
{
	AIErr result = kNoErr;
	try {
		ATE::ATEErr error = ATE::kMissingSpellingDictionary;
		throw ATE::Exception(error);
	}	
	catch (ATE::Exception& ex) {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		result = ex.error;
		if (result == ATE::kMissingSpellingDictionary) {
			log->WritePrintf("ATE error code: %d - kMissingSpellingDictionary", result);
			// Correct exception caught so test passed.
			result = kNoErr;
		}
	}
	return result;
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpTextException available to the SnippetRunner framework.
 */
class _SnpRunnableTextException : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableTextException () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableTextException () {}

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
std::string _SnpRunnableTextException::GetName() const
{
	return "TextException";
}

/*
*/
std::string _SnpRunnableTextException::GetDescription() const 
{
	return "Demonstrates throwing and catching an ATE Exception";
}

/*
*/
SnpRunnable::Operations _SnpRunnableTextException::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("ThrowATEException", "", kSnpRunAnyContext));
	return operations;
}

/*
*/
std::string _SnpRunnableTextException::GetDefaultOperationName() const 
{
	return "ThrowATEException";
	// Note: this must be one of the operation names returned by GetOperations.
}

/*
*/
std::vector<std::string> _SnpRunnableTextException::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Text Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableTextException::CanRun(SnpRunnable::Context& runnableContext)
{
	return true;
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableTextException::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpTextException instance;
	if ("ThrowATEException" == runnableContext.GetOperation().GetName()) {
		result = instance.ThrowATEException();
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableTextException instance_SnpRunnableTextException;

// End SnpTextException.cpp

