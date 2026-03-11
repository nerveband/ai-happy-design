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

#ifndef __SnpUnitTestHelper_h__
#define __SnpUnitTestHelper_h__

/** Helper class for snippet unit testing.

	A number of the pre-defined contexts depend on objects in the "unittest-1.ai" Illustrator
	document. This document contains objects which are grouped by type into different layers.
	You can add new layers with new types of object to this document and update the 
	SnpUnitTestHelper::SetupContext implementation to find and return those objects. This can
	be done for objects that more than one snippet will need for testing.

  	@ingroup sdk_snippet
 */
class SnpUnitTestHelper
{
	public:
		/** Called to set up the given context.
			@param runnableContext IN context the snippet requires.
			@param runnable IN the snippet.
			@return kNoErr on success, other ASErr otherwise.
			@see SnpRunnable::SetupContext
		 */
		ASErr SetupContext(SnpRunnable::Context& runnableContext, SnpRunnable* runnable);

		/** Called to tear down the given context - implementation saves and closes all open documents.
			@param runnableContext IN context the snippet requires.
			@param runnable IN the snippet.
			@return kNoErr on success, other ASErr otherwise.
			@see SnpRunnable::SetupContext
		 */
		ASErr TeardownContext(SnpRunnable::Context& runnableContext, SnpRunnable* runnable);

		/** Creates a new document.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr NewDocument();

		/** Opens a document, if a document of the given name can be found
			in SnippetRunner's asset folder it is opened, otherwise the user
			is prompted to locate the file.
			@param name IN name of asset file.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr OpenDocument(const ai::UnicodeString& name);

		/** Saves a document to SnippetRunner's asset folder then closes the document.
			@param name IN name to save document to.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr SaveAndCloseDocument(const ai::UnicodeString& name);

		/** Closes all open documents.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr CloseAllDocuments();

		/** Returns the categories a snippet belongs to.
			@return default categories.
		*/
		std::vector<std::string> GetCategories() const;

		/** Stack based helper class for logging a customised unit test - if
			you define SnpRunnable::RunUnitTest for your snippet this class
			can be useful for formatting your output to the log.

			@ingroup snippetrunner
		*/
		class Logger {
		public:
			/** Constructor - logs beginning of unit test.
			*/
			Logger(const char* name);

			/** Sets pass or fail status for the test and updates the overall result
				@param status IN the status of the unit test, kNoErr for pass, other ASErr for fail
				@return the error code given by status
			*/
			ASErr SetStatus(ASErr status);

			/** Destructor - logs end of test with pass / fail result.
			*/
			~Logger();

		private:
			/** Name of unit test.
			*/
			const char* fName;

			/**	Contains pass or fail status of the unit test.
			*/
			ASErr fStatus;
		};

};

#endif // __SnpUnitTestContext_h__

// End SnpUnitTestContext.cpp

