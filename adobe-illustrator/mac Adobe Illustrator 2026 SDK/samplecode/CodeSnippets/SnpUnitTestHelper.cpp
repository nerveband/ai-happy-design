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
#include "SnippetRunnerSuites.h"
#include "SnpRunnable.h"
#include "SnpDocumentHelper.h"
#include "SnpDocumentActionHelper.h"
#include "SnpArtSetHelper.h"
#include "SnpUnitTestHelper.h"
#include "SnippetRunnerParameter.h"
#include "SnippetRunnerLog.h"

/*
*/
ASErr SnpUnitTestHelper::SetupContext(SnpRunnable::Context& runnableContext, SnpRunnable* runnable)
{
	ASErr result = kNoErr;
	try {
		switch (runnableContext.GetOperation().GetUnitTestContext())
		{
			case SnpRunnable::kSnpRunNotSupportedContext:
				{
					result = kCantHappenErr;
				}
				break;
			case SnpRunnable::kSnpRunAnyContext: 
				{
					result = kNoErr;
				}
				break;
			case SnpRunnable::kSnpRunNoDocumentContext: 
				{
					result = kNoErr;
				}
				break;
			case SnpRunnable::kSnpRunNewDocumentContext: 
				{	
					result = this->NewDocument();
				}
				break;
			case SnpRunnable::kSnpRunUnitTestDocumentContext: 
				{
					result = this->OpenDocument(ai::UnicodeString("unittest-1.ai"));
				}
				break;
			case SnpRunnable::kSnpRunPathContext: 
				{
					result = this->OpenDocument(ai::UnicodeString("unittest-1.ai"));
					SnpArtSetHelper artSet;
					artSet.FindLayerArt(ai::UnicodeString("layer.path"));
					artSet = artSet.Filter(kPathArt);
					artSet.SelectAll();
					result = sAIDocument->RedrawDocument();
					aisdk::check_ai_error(result);
				}
				break;
			case SnpRunnable::kSnpRunJPGContext: 
				{
					result = this->OpenDocument(ai::UnicodeString("unittest-1.ai"));
					SnpArtSetHelper artSet;
					artSet.FindLayerArt(ai::UnicodeString("layer.jpg"));
					artSet = artSet.Filter(kRasterArt);
					artSet.SelectAll();
					result = sAIDocument->RedrawDocument();
					aisdk::check_ai_error(result);
				}
				break;
			case SnpRunnable::kSnpRunPSDContext: 
				{
					result = this->OpenDocument(ai::UnicodeString("unittest-1.ai"));
					SnpArtSetHelper artSet;
					artSet.FindLayerArt(ai::UnicodeString("layer.psd"));
					artSet = artSet.Filter(kRasterArt);
					artSet.SelectAll();
					result = sAIDocument->RedrawDocument();
					aisdk::check_ai_error(result);
				}
				break;
			case SnpRunnable::kSnpRunEPSContext: 
				{
					result = this->OpenDocument(ai::UnicodeString("unittest-1.ai"));
					SnpArtSetHelper artSet;
					artSet.FindLayerArt(ai::UnicodeString("layer.eps"));
					artSet = artSet.Filter(kPlacedArt);
					artSet.SelectAll();
					result = sAIDocument->RedrawDocument();
					aisdk::check_ai_error(result);
				}
				break;
			case SnpRunnable::kSnpRunPDFContext: 
				{
					result = this->OpenDocument(ai::UnicodeString("unittest-1.ai"));
					SnpArtSetHelper artSet;
					artSet.FindLayerArt(ai::UnicodeString("layer.pdf"));
					artSet = artSet.Filter(kPlacedArt);
					artSet.SelectAll();
					result = sAIDocument->RedrawDocument();
					aisdk::check_ai_error(result);
				}
				break;
			case SnpRunnable::kSnpRunCharactersContext: 
				{
					result = this->OpenDocument(ai::UnicodeString("unittest-1.ai"));
					SnpArtSetHelper artSet;
					artSet.FindLayerArt(ai::UnicodeString("layer.text"));
					artSet = artSet.Filter(kTextFrameArt);
					artSet.SelectAll();
					result = sAIDocument->RedrawDocument();
					aisdk::check_ai_error(result);
				}
				break;
			case SnpRunnable::kSnpRunCustomContext:
				SDK_ASSERT_MSG(false, "SnpRunnable::kSnpRunCustomContext requires custom implementation of SetupContext(SnpRunnable::Context& runnableContext)");
				break;
			default:
				SDK_ASSERT_MSG(false, "Unknown ISnpRunnableContext type");
		}

	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpUnitTestHelper::TeardownContext(SnpRunnable::Context& runnableContext, SnpRunnable* runnable)
{
	ASErr result = kNoErr;
	try {
		std::string fileName = "out-" + runnable->GetName() + "-" + runnableContext.GetOperation().GetName() + ".ai";
		result = this->SaveAndCloseDocument(ai::UnicodeString(fileName));
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpUnitTestHelper::NewDocument()
{
	SnpDocumentHelper docHelper;
	AIDocumentHandle document = NULL;
	return docHelper.NewDocument(ai::UnicodeString("Print"), document);
}

/*
*/
ASErr SnpUnitTestHelper::OpenDocument(const ai::UnicodeString& name)
{
	SnpDocumentHelper docHelper;
	AIDocumentHandle document = NULL;
	return docHelper.OpenDocument(name, document);
}

/*
*/
ASErr SnpUnitTestHelper::SaveAndCloseDocument(const ai::UnicodeString& name)
{
	ASErr result = kNoErr;
	try {
		SnpDocumentHelper docHelper;
		AIDocumentHandle document = NULL;
		docHelper.GetCurrentDocument(document);
		if (document) {
			// Save document via an action to avoid any dialogs.
			SnpDocumentActionHelper docActionHelper;
			result = docActionHelper.SaveDocumentAsAI(name);
			SDK_ASSERT(!result);

			// Close the document.
			result = docHelper.CloseDocument(document);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpUnitTestHelper::CloseAllDocuments()
{
	SnpDocumentHelper docHelper;
	return docHelper.CloseAllDocuments();
}

/*
*/
std::vector<std::string> SnpUnitTestHelper::GetCategories() const
{
	std::vector<std::string> categories;
	categories.push_back("All Snippets");
	categories.push_back("Helper Snippets");
	return categories;
}

/*
*/
SnpUnitTestHelper::Logger::Logger(const char* name) : fName(name), fStatus(kCantHappenErr)
{
	SnippetRunnerLog* log = SnippetRunnerLog::Instance();
	log->Write(std::string(fName) + std::string(" Begin"));
	log->BeginIndent();
}

/*
*/
ASErr SnpUnitTestHelper::Logger::SetStatus(ASErr status)
{
	fStatus = status;
	return status;
}

/*
*/
SnpUnitTestHelper::Logger::~Logger()
{
	SnippetRunnerLog* log = SnippetRunnerLog::Instance();
	if (fStatus) {
		log->WritePrintf("FAILED: error=%x", fStatus);
	}
	else {
		log->Write("PASSED");
	}
	log->EndIndent();
	log->Write(std::string(fName) + std::string(" End"));
}

// End SnpUnitTestHelper.cpp

