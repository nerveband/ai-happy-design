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

// Project includes:
#include "SDKErrors.h"
#include "SnippetRunnerID.h"
#include "SnpRunnable.h"
#include "SnippetRunnerLog.h"
#include "SnippetRunnerPlugin.h"
#include "SnippetRunnerSuites.h"
#include "SnippetRunnerUnitTestManager.h"
#include "SnippetRunnerManager.h"

// Can't use dynamic structures (e.g. std::vector) to 
// register snippets because snippets load statically and Illustrator's
// new/delete operators are not available until the application has started.
// So we use a fixed size static array instead.
const	ASUInt32		kSnippetRunnerRegistrySize = 200;
static	SnpRunnable*	sSnippetRegistry[kSnippetRunnerRegistrySize];
static	ASUInt32		sSnippetCount = 0;

/*
*/
SnippetRunnerManager& SnippetRunnerManager::Instance()
{
	static SnippetRunnerManager gSnippetRunnerManager;
	return gSnippetRunnerManager;
}

/*
*/
SnippetRunnerManager::SnippetRunnerManager()
{
}

/*
*/
SnippetRunnerManager::~SnippetRunnerManager()
{
}

/*
*/
void SnippetRunnerManager::AddSnippet(SnpRunnable* runnable)
{
	if (sSnippetCount < kSnippetRunnerRegistrySize && runnable != NULL) {
		sSnippetRegistry [sSnippetCount] = runnable;
		sSnippetCount++;
	}
}

/*
*/
void SnippetRunnerManager::RemoveSnippet(SnpRunnable* runnable)
{
	for (ASUInt32 i = 0; i < sSnippetCount; i++) {
		if (sSnippetRegistry [i] != NULL && sSnippetRegistry [i] == runnable) {
			sSnippetRegistry [i] = NULL;
		}
	}
}

/* 
*/
SnpRunnable* SnippetRunnerManager::FindSnippet(const std::string& snippetName) const
{
	SnpRunnable* result = NULL;
	for (ASUInt32 i = 0; i < sSnippetCount; i++) {
		if (sSnippetRegistry [i] != NULL && snippetName == sSnippetRegistry [i]->GetName()) {
			result = sSnippetRegistry [i];
		}
	}
	return result;
}

/* 
*/
SnpRunnable::Context SnippetRunnerManager::GetContext(const std::string& snippetName, const std::string& operationName) const
{
	SnpRunnable::Operation operation("", "", SnpRunnable::kSnpRunNotSupportedContext);
	SnpRunnable* snp = this->FindSnippet(snippetName);
	if (snp) {
		SnpRunnable::Operations operations = snp->GetOperations();
		for (SnpRunnable::Operations::const_iterator iter = operations.begin(); iter != operations.end(); iter++) {
			if (operationName == iter->GetName()) {
				operation = *(iter);
				break;
			}
		}
	}
	return SnpRunnable::Context(operation);
}

/*
*/
ASBoolean SnippetRunnerManager::CanRun(const std::string& snippetName, const std::string& operationName) const
{
	ASBoolean result = false;
	//individual snippet or category
	if(!(snippetName == "category"))
	{
		// Ask the snippet if it can run.
		SnpRunnable* snp = this->FindSnippet(snippetName);
		if (snp) {
			std::string snippetOperation = operationName;
			if (snippetOperation.size() == 0) {
					snippetOperation = snp->GetDefaultOperationName();
			}
			SnpRunnable::Context context = this->GetContext(snippetName, snippetOperation);
			result = snp->CanRun(context);

		}
	}
	else
	{
		SnippetRunnerManager::Snippets snippets = GetSnippets(operationName);
		for (SnippetRunnerManager::Snippets::const_iterator iter = snippets.begin(); iter != snippets.end(); iter++) {
				SnpRunnable* snp = *iter;
				SnpRunnable::Operations operations = snp->GetOperations();
				for(SnpRunnable::Operations::const_iterator operations_iter = operations.begin(); operations_iter != operations.end(); operations_iter++)
				{
					SnpRunnable::Operation operation = *operations_iter;
					if(!result)
					{
						SnpRunnable::Context context = this->GetContext(snippetName, operation.GetName());
						result = snp->CanRun(context);
						if(result)
							break;
					}
				}
				if(result)
					break;
		}
	}
	return result;
}

/*
*/
ASErr SnippetRunnerManager::Run(const std::string& snippetName, const std::string& operationName) const
{
	ASErr result = kNoErr;
	if (operationName.size() != 0) {
		// Run given operation.
		result = this->Execute(snippetName, operationName);
	}
	else {
		// Run default operation.
		SnpRunnable* snp = SnippetRunnerManager::Instance().FindSnippet(snippetName);
		std::string snippetOperation;
		if (snp) {
			snippetOperation = snp->GetDefaultOperationName();
		}
		result = this->Execute(snippetName, snippetOperation);
	}
	return result;
}

/*
*/
ASErr SnippetRunnerManager::RunUnitTest(const std::vector<std::string>& snippetNames, const std::string& operationName)
{
	return SnippetRunnerUnitTestManager::Instance()->RunUnitTest(snippetNames, operationName);
}

/*
*/
ASErr SnippetRunnerManager::Execute(const std::string& snippetName, const std::string& snippetOperation) const
{
	ASErr result = kNoErr;
	SnpRunnable* snp = SnippetRunnerManager::Instance().FindSnippet(snippetName);
	SnippetRunnerLog* log = SnippetRunnerLog::Instance();
	if (snp) {
		log->Write("Begin " + snippetName + "::" + snippetOperation);
		log->BeginIndent();
		try {
			SnpRunnable::Context context = this->GetContext(snippetName, snippetOperation);
			if (this->IsUnitTestRunning()) {
				result = snp->RunUnitTest(context);
			} 
			else {
				result = snp->Run(context);
			}
			if (result == kNoErr) {
				log->Write("SUCCESS");
			}
			else if (result == kCanceledErr) {
				log->Write("CANCELLED");
			}
			else if (result == kBadParameterErr) {
				log->Write("FAILED: bad parameter");
			}
			else {
				log->WritePrintf("FAILED: error=%x", result);
			}
		}
		catch (ai::Error& ex) {
			ASErr result = ex;
			log->WritePrintf("FAILED: Unhandled ai::Error exception=%x", result);
		}
		catch (...) {
			log->WritePrintf("FAILED: Unhandled exception");
		}				
		log->EndIndent();
		log->Write("End " + snippetName + "::" + snippetOperation);
	}
	else {
		log->Write(std::string("ERROR: cannot find snippet") + snippetName);
	}
	log->NotifyLogChanged();
	return result;
}

/*
*/
static bool SnippetSortAscending(const SnpRunnable* lhs, const SnpRunnable* rhs)
{
	return lhs->GetName() < rhs->GetName();
}

/*
*/
SnippetRunnerManager::Snippets SnippetRunnerManager::GetSnippets(const std::string category) const
{
	Snippets snippets;
	for (ASUInt32 i = 0; i < sSnippetCount; i++) {
		if (sSnippetRegistry [i] != NULL) {
			SnpRunnable* snippet = sSnippetRegistry [i];
			std::vector<std::string> categories = snippet->GetCategories();
			std::vector<std::string>::iterator iter;
			ASBoolean found = false;
			for (iter = categories.begin(); iter != categories.end(); iter++) {
				if (category.compare(0, iter->length(), iter->c_str()) == 0) {
					found = true;
					break;
				}
			}
			if (found)
				snippets.push_back(sSnippetRegistry [i]);
		}
	}
	std::sort(snippets.begin(), snippets.end(), SnippetSortAscending);
	return snippets;
}

/*
*/
ASBoolean SnippetRunnerManager::IsUnitTestRunning() const
{
	return SnippetRunnerUnitTestManager::Instance()->GetUnitTestState() != SnippetRunnerUnitTestManager::kSnpRunUnitTestIdle;
}
// End SnippetRunnerManager.cpp
