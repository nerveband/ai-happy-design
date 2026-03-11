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

// Framework includes:
#include "SDKErrors.h"
#include "SnpRunnable.h"
#include "SnippetRunnerLog.h"
#include "SnippetRunnerParameter.h"
#include "SnippetRunnerTokenisor.h"
#include "SnpUnitTestHelper.h"

/** Tests SnippetRunner internals.

  	@ingroup sdk_snippet
 */
class SnpTestSnippetRunner
{
	public:
		/** 
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr Do();
		
		/** 
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr DoParameterTest(void);

		/** 
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr DoTokenisorTest(void);

};

/*
*/
ASErr SnpTestSnippetRunner::Do(void)
{
	ASErr result = kNoErr;
	try {
		{
			SnpUnitTestHelper::Logger ut("SnpTestSnippetRunner::DoParameterTest");
			ut.SetStatus(this->DoParameterTest());
		}

		{
			SnpUnitTestHelper::Logger ut("SnpTestSnippetRunner::DoTokenisorTest");
			ut.SetStatus(this->DoTokenisorTest());
		}
	}
	catch (...) {
		result = kCantHappenErr;
	}
	return result;
}

/*
*/
ASErr SnpTestSnippetRunner::DoParameterTest(void)
{
	ASErr result = kNoErr;

	try {
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ai::UnicodeString s = parameter->GetString(ai::UnicodeString("s"), ai::UnicodeString("s-value"));
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		log->WritePrintf("s: %s", s.getInStdString(kAIPlatformCharacterEncoding).c_str());

		std::string s1 = parameter->GetString("s1", "s1-value");
		log->Write(std::string("s1: ") + s1);

		AIReal r = parameter->GetReal(ai::UnicodeString("r (range checked)"), 1.0, 0.0, 10.0);
		log->WritePrintf("r: %f", r);
		r = parameter->GetReal("r", 2.0);
		log->WritePrintf("r: %f", r);

		ASInt32 i = parameter->GetInt(ai::UnicodeString("i (range checked)"), 1, 0, 10);
		log->WritePrintf("i: %d", i);
		i = parameter->GetInt("i", 2);
		log->WritePrintf("i: %d", i);

		std::vector<ai::UnicodeString> choices;
		choices.push_back(ai::UnicodeString("Do this"));
		choices.push_back(ai::UnicodeString("Do that"));
		choices.push_back(ai::UnicodeString("Do the other"));
		ASUInt32 c = parameter->GetChoice(ai::UnicodeString("c"), choices, 1);
		log->WritePrintf("c: %d", c);

		std::vector<std::string> choices2;
		choices2.push_back("Do this again");
		choices2.push_back("Do that again");
		choices2.push_back("Do the other again");
		ASUInt32 c2 = parameter->GetChoice("c2", choices2, 2);
		log->WritePrintf("c2: %d", c2);

		ai::FilePath path = parameter->GetFilePath(ai::UnicodeString("Choose input file"));
		log->WritePrintf("input path: %s", path.GetFullPath().getInStdString(kAIPlatformCharacterEncoding).c_str());

		path = parameter->GetFilePath(ai::UnicodeString("Choose output file"), false);
		log->WritePrintf("output path: %s", path.GetFullPath().getInStdString(kAIPlatformCharacterEncoding).c_str());

		path = parameter->GetFolderPath(ai::UnicodeString("Choose folder"));
		log->WritePrintf("folder path: %s", path.GetFullPath().getInStdString(kAIPlatformCharacterEncoding).c_str());

	}
	catch (...) {
		result = kCantHappenErr;
	}
	
	return result;
}

/*
*/
ASErr SnpTestSnippetRunner::DoTokenisorTest(void)
{
	ASErr result = kNoErr;

	try {
		SnippetRunnerTokenisor instance;
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		log->Write("GetString test");
		instance.SetParameters(ai::UnicodeString("0, 1 , 2 "));
		{ 
			ai::UnicodeString v;
			instance.GetString(v);
			aisdk::check_assertion(v == ai::UnicodeString("0"));
		}
		{ 
			ai::UnicodeString v;
			instance.GetString(v);
			aisdk::check_assertion(v == ai::UnicodeString("1"));
		}
		{ 
			ai::UnicodeString v;
			instance.GetString(v);
			aisdk::check_assertion(v == ai::UnicodeString("2"));
		}
		log->Write("GetString test passed");
		
		log->Write("GetString quoting test");
		instance.SetParameters(ai::UnicodeString("0,\"1  1\",2"));
		{ 
			ai::UnicodeString v;
			instance.GetString(v);
			aisdk::check_assertion(v == ai::UnicodeString("0"));
		}
		{ 
			ai::UnicodeString v;
			instance.GetString(v);
			aisdk::check_assertion(v == ai::UnicodeString("1  1"));
		}
		{ 
			ai::UnicodeString v;
			instance.GetString(v);
			aisdk::check_assertion(v == ai::UnicodeString("2"));
		}
		log->Write("GetString quoting test passed");		
		
		log->Write("GetString blank parameter test");
		instance.SetParameters(ai::UnicodeString("0,,2"));
		{ 
			ai::UnicodeString v;
			instance.GetString(v);
			aisdk::check_assertion(v == ai::UnicodeString("0"));
		}
		{ 
			ai::UnicodeString v("foo");
			instance.GetString(v);
			aisdk::check_assertion(v.size() == 0);
		}
		{ 
			ai::UnicodeString v;
			instance.GetString(v);
			aisdk::check_assertion(v == ai::UnicodeString("2"));
		}
		log->Write("GetString blank parameter test passed");

		log->Write("Mixed parameter type test");
		instance.SetParameters(ai::UnicodeString("hello,1,2.0,foof"));
		{
			ai::UnicodeString sValue;
			ai::int32 iValue = 0;
			AIReal rValue = 0.0;
			ai::int32 cValue = 0;
			std::vector<ai::UnicodeString> choices;
			choices.push_back(ai::UnicodeString("Hello"));
			choices.push_back(ai::UnicodeString("foof"));
			choices.push_back(ai::UnicodeString("there"));

			instance.GetString(sValue);
			aisdk::check_assertion(sValue == ai::UnicodeString("hello"));
			instance.GetInt(iValue);
			aisdk::check_assertion(iValue == 1);
			instance.GetReal(rValue);
			aisdk::check_assertion(rValue == 2.0);
			instance.GetChoice(cValue, choices);
			aisdk::check_assertion(cValue == 1);
		}
		log->Write("Mixed parameter type test passed");

	}
	catch (...) {
		result = kCantHappenErr;
	}

	return result;
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpTestSnippetRunner available to the SnippetRunner framework.
 */
class _SnpTestSnippetRunner : public SnpRunnable
{
	public:
		/* This constructor registers your snippet with the framework.
		 */
		_SnpTestSnippetRunner () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpTestSnippetRunner () {}

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

		/* Run the snippet.
			@param runnableContext see ISnpRunnableContext for documentation.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr		Run(SnpRunnable::Context& runnableContext);
};

/*
*/
std::string _SnpTestSnippetRunner::GetName() const
{
	return "TestSnippetRunner";
}

/*
*/
std::string _SnpTestSnippetRunner::GetDescription() const 
{
	return "Provides unit tests for internal SnippetRunner classes";
}

/*
*/
SnpRunnable::Operations _SnpTestSnippetRunner::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("Do", "", kSnpRunAnyContext));
	return operations;
}

/*
*/
std::string _SnpTestSnippetRunner::GetDefaultOperationName() const 
{
	return "Do";
}

/*
*/
std::vector<std::string> _SnpTestSnippetRunner::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Utility Snippets");
	return categories;
}

/*
*/
ASBoolean	_SnpTestSnippetRunner::CanRun(SnpRunnable::Context& runnableContext)
{
	ASBoolean result = true;
	return result;
}

/*
*/
ASErr _SnpTestSnippetRunner::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpTestSnippetRunner instance;
	result = instance.Do();
	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpTestSnippetRunner instance_SnpTestSnippetRunner;

// End, SnpTestSnippetRunner.cpp

