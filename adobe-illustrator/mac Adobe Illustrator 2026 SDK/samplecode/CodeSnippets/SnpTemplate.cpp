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
#include "SnippetRunnerSuites.h" // TODO declare additional suites here if you need them
#include "SnippetRunnerLog.h"
#include "SnippetRunnerParameter.h"

/** 
	\li	Provides the template code snippet to use as a start point for new snippets.

	To use the template take the following steps.

	Make a copy of SnpTemplate.cpp and rename it according to the
	function of your snippet, for example SnpCreateFrame.cpp. 

	Add your snippet (SnpCreateFrame.cpp) to the SnippetRunner projects
	for Windows and Macintosh.

	Replace string "Template" with "CreateFrame". Note don't
	do the seemingly similar replace of "SnpTemplate" with 
	"SnpCreateFrame" or the name of your snippet in the user
	interface won't get set properly.

    There are normally two classes in a code snippet:
	\li SnpTemplate in which you define your functional code, and,
	\li _SnpRunnableTemplate using which hooks your snippet into
		the SnippetRunner framework.

	Edit the _SnpRunnableTemplate code to describe your snippet to the framework:
	\li GetName returns the name of your snippet
	\li GetDescription returns a description of what your snippet does
	\li GetOperations returns operations supported by your snippet
	\li GetDefaultOperationName returns the name of your snippet's default operation

	Edit the _SnpRunnableTemplate::CanRun method to check the preconditions
	your snippet requires are met.

	Edit your functional methods(s) into your SnpTemplate class.

	Edit the _SnpRunnableTemplate::Run method to call your snippet's
	function methods.

	Use SnippetRunnerParameter to get input parameters.

	Use SnippetRunnerLog to log output messages.

	If your snippet requires suites that are not already acquired and released by the
	framework perform the following steps:
	\li Declare the suite pointers you need in SnippetRunnerSuites.h.
	\li Define the suite pointers you need and add the suite and version IDs 
	in SnippetRunnerSuites.cpp so that the suites are available for you to use.

  	@ingroup sdk_snippet
 */
class SnpTemplate
{
	public:
		/** TODO: Replace with functional method(s).
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr Do(void);
};

/*
*/
ASErr SnpTemplate::Do(void)
{
	// TODO replace with functional code
	ASErr result = kNoErr;
	try {
		// Get input parameters one at a time.
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ai::UnicodeString param = parameter->GetString(ai::UnicodeString("Enter parameter"), ai::UnicodeString("default"));
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		log->Write(ai::UnicodeString("Template parameter = ").append(param));
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpTemplate available to the SnippetRunner framework.
 */
class _SnpRunnableTemplate : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableTemplate () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableTemplate () {}

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
std::string _SnpRunnableTemplate::GetName() const
{
	// TODO - edit as required.
	return "Template";
}

/*
*/
std::string _SnpRunnableTemplate::GetDescription() const 
{
	// TODO - edit as required.
	return "Provides a template";
}

/*
*/
SnpRunnable::Operations _SnpRunnableTemplate::GetOperations() const
{
	// TODO - edit as required.
	SnpRunnable::Operations operations;
	operations.push_back(Operation("Do", "myPreconditions", kSnpRunAnyContext));
	// Add a line as above for each operation your snippet wants to advertise to the framework.
	// Operation(const std::string& name, const std::string& preconditions, eSnpRunContextType unitTestContext)
	// name - name of the operation 
	// preconditions - description of preconditions to be displayed in SnippetRunner panel
	// unitTestContext - context you want the framework to automatically set for your snippet during unit test.
	return operations;
}

/*
*/
std::string _SnpRunnableTemplate::GetDefaultOperationName() const 
{
	// TODO - edit as required.
	return "Do";
	// Note: this must be one of the operation names returned by GetOperations.
}

/*
*/
std::vector<std::string> _SnpRunnableTemplate::GetCategories() const
{
	// TODO - Modify the category below to the most relevant category for this snippet
	// or specify a new category.
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Utility Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableTemplate::CanRun(SnpRunnable::Context& runnableContext)
{
	// TODO - edit as required.
	return true;
	// Note: this will control trhe enabling and disabling of the Run button in the SnippetRunner panel
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableTemplate::Run(SnpRunnable::Context& runnableContext)
{
	// TODO edit as required by your snippet - you don't have to have a Do method in your class.
	ASErr result = kNoErr;
	SnpTemplate instance;
	if ("Do" == runnableContext.GetOperation().GetName()) {
		result = instance.Do();
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableTemplate instance_SnpRunnableTemplate;

// End SnpTemplate.cpp

