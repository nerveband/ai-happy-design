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

#ifndef _SnpRunnable_
#define _SnpRunnable_

// Project includes:
#include "SnippetRunnerID.h"

/** Hooks a snippet into the SnippetRunner framework. To register the
	snippet with the framework a static instance of the snippet must
	be declared. The constructor registers the snippet via SnippetRunnerManager. 

	This registration happens each time the plug-in is loaded into memory. 
	This approach that is ok for a development tool type plug-in like SnippetRunner 
	but is NOT one that should be used for functional plug-ins.

	@ingroup snippetrunner
*/
class SnpRunnable
{
public:
		/** Context type - SnippetRunner can automatically set up contexts of these types
			for unit testing.
			@see SnpRunnable::RunUnitTest, SnpRunnable::SetupContext, SnpRunnable::TeardownContext
		*/
		typedef enum {	
				kSnpRunNotSupportedContext, // Context type features are not supported.
				kSnpRunAnyContext, // Requires no context 
				kSnpRunNoDocumentContext, // Requires no documents open
				kSnpRunNewDocumentContext, // Requires document with nothing selected
				kSnpRunUnitTestDocumentContext, // Requires unit test document with nothing selected
				kSnpRunPathContext, // Requires path selected
				kSnpRunJPGContext, // Requires JPEG image selected
				kSnpRunPSDContext, // Requires PSD image selected
				kSnpRunEPSContext, // Requires linked EPS selected
				kSnpRunPDFContext, // Requires linked PDF selected
				kSnpRunTypeContext, // Requires text frame selected
				kSnpRunCharactersContext, // Requires text range selected
				/** selecting custom requires overriding SetupContext in your SnpRunnable derived class */ 
				kSnpRunCustomContext
		} eSnpRunContextType;

		/** Describes one operation supported by the snippet.
		*/
		class Operation {
		public:
			/** Constructor.
			*/
			Operation(const std::string& name, const std::string& preconditions, eSnpRunContextType unitTestContext) : fName(name), fPreconditions(preconditions), fUnitTestContext(unitTestContext) {}

			/** Destructor.
			*/
			virtual ~Operation() {}

			/** Gets the name of the snippet operation.
				@return snippet operation name.
			*/
			std::string GetName() const {return fName;}

			/** Get the string describing the preconditions for this operation for display in the framework's UI.
				@return precondition string.
			*/
			std::string GetPreconditions() const {return fPreconditions;}

			/** Gets the context this operation requires for unit test-
				a number of contexts can be automatically supplied by the framework which
				saves you writing custom unit test code.
				@return the context.
			*/
			eSnpRunContextType GetUnitTestContext() const {return fUnitTestContext;}

		private:

			/** Name of the operation.
			*/
			std::string fName;

			/**	Describes the preconditions necessary for 
				successful execution of the operation.
			*/
			std::string fPreconditions;

			/** Identifies the context this operation requires for unit test.
				@see SnpRunnable::RunUnitTest, SnpRunnable::SetupContext, SnpRunnable::TeardownContext
			*/
			eSnpRunContextType fUnitTestContext;
		};

		/** Describes all operations supported by the snippet.
		*/
		typedef std::vector<Operation> Operations;


		/** Context for execution of a snippet.

			@ingroup snippetrunner
		 */
		class Context
		{
			public:
				/** Constructor
				*/
				Context(const SnpRunnable::Operation& operation) : fOperation(operation) {}

				/** Gets the snippet operation.
					@return the snippet operation.
				*/
				SnpRunnable::Operation GetOperation(void) const {return fOperation;}
			private:
				/** The snippet operation.
				*/
				SnpRunnable::Operation fOperation;
		};

	public:

		/**	Registers the snippet with SnippetRunnerManager.
		*/
		SnpRunnable();

		/**	Deregisters the snippet with SnippetRunnerManager.
		*/
		virtual ~SnpRunnable();

		/** Gets the name of snippet.
			@return the name of the snippet.
		*/
		virtual std::string GetName() const = 0;

		/** Gets a description of what the snippet does.
			@return the snippet description.
		*/
		virtual std::string	GetDescription() const = 0;

		/** Returns true if preconditions for running are met, false otherwise.
			@param runnableContext IN context the snippet requires.
			@return true if preconditions for running are met, false otherwise.
		*/
		virtual ASBoolean CanRun(Context& runnableContext) = 0;

		/** Runs the snippet.
			@param runnableContext IN context the snippet requires.
			@return kNoErr on success, other ASErr otherwise
		*/
		virtual ASErr Run(Context& runnableContext) = 0;

		/** Gets the operations a snippet supports - a snippet must have at least one operation.
			@return the snippets operations.
		*/
		virtual Operations GetOperations() const = 0;

		/** Gets the name of the snippet's default operation.
			@return name of default snippet operation.
		*/
		virtual std::string GetDefaultOperationName() const = 0;

		/** Gets the default categories a snippet belongs to.
			@return default categories.
		*/
		virtual std::vector<std::string> GetCategories() const;

		/** Default behavior for snippet unit test.
			@param runnableContext IN context the snippet requires.
			@return kNoErr on success, other ASErr otherwise
		*/
		virtual ASErr RunUnitTest(SnpRunnable::Context& runnableContext);

		/** Finds out whether the snippet can be loaded.
			@return true (default) if the snippet should be presented in the framework's UI, false otherwise.
		*/
		virtual ASBoolean CanLoad() const {return true;}

		/** Sets up the default context for the snippet so that it
			can be run by SnpRunnable::RunUnitTest.
			@param runnableContext IN context the snippet requires.
			@return kNoErr on success, other ASErr otherwise
		*/
		virtual ASErr SetupContext(SnpRunnable::Context& runnableContext);

		/** Tears down the default context for the snippet after it
			has been run by SnpRunnable::RunUnitTest.
			@param runnableContext IN context the snippet requires.
			@return kNoErr on success, other ASErr otherwise
		*/
		virtual ASErr TeardownContext(SnpRunnable::Context& runnableContext);
};

#endif // SnpRunnable

// End SnpRunnable.h
