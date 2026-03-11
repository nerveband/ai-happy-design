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
#include "SnpArtSetHelper.h"
#include "SnpRunnable.h"
#include "SnippetRunnerSuites.h" 
#include "SnippetRunnerLog.h"
#include "SnpUnitTestHelper.h"
#include "SnpDocumentHelper.h"
#include "SnpSelectionHelper.h"
#include "SnpChooser.h"

/* Creates empty art set.
*/
SnpArtSetHelper::SnpArtSetHelper() : fArtSet(NULL)
{
	ASErr status = sAIArtSet->NewArtSet(&this->fArtSet);
	SDK_ASSERT(!status);
}

/* Creates art set containing specified art.
*/
SnpArtSetHelper::SnpArtSetHelper(AIArtSpec* specs, short numSpecs) : fArtSet(NULL)
{
	ASErr status = kNoErr;
	try {
		status = sAIArtSet->NewArtSet(&this->fArtSet);
		aisdk::check_ai_error(status);
		
		status = sAIArtSet->MatchingArtSet(specs, numSpecs, this->fArtSet);
		aisdk::check_ai_error(status);
	}
	catch (ai::Error& ex) {
		status = ex;
		if (this->fArtSet) {
			sAIArtSet->DisposeArtSet(&this->fArtSet);
		}
	}
	SDK_ASSERT(!status);
}

/* Copy constructor
*/
SnpArtSetHelper::SnpArtSetHelper(const SnpArtSetHelper& rhs) : fArtSet(NULL)
{
	ASErr status = kNoErr;
	try {
		status = sAIArtSet->NewArtSet(&this->fArtSet);
		aisdk::check_ai_error(status);
		
		for (size_t i = 0; i < rhs.GetCount(); i++) {
			AIArtHandle art = rhs[i];
			status = sAIArtSet->AddArtToArtSet(this->fArtSet, art);
			aisdk::check_ai_error(status);
		}
	}
	catch (ai::Error& ex) {
		status = ex;
		if (this->fArtSet) {
			sAIArtSet->DisposeArtSet(&this->fArtSet);
		}
	}
	SDK_ASSERT(!status);
} // end copy constructor

/* Destructor.
*/
SnpArtSetHelper::~SnpArtSetHelper()
{
	if(this->fArtSet) {
		ASErr status = sAIArtSet->DisposeArtSet(&this->fArtSet);
		this->fArtSet = NULL;
		SDK_ASSERT(!status);
	}
}

/* = operator (assignment)
*/
SnpArtSetHelper& SnpArtSetHelper::operator=(const SnpArtSetHelper& rhs)
{
	ASErr status = kNoErr;
    if (this != &rhs) // make sure not same object
	{  
		try
		{
			status = sAIArtSet->DisposeArtSet(&this->fArtSet);
			aisdk::check_ai_error(status);

			this->fArtSet = NULL;

			status = sAIArtSet->NewArtSet(&this->fArtSet);
			aisdk::check_ai_error(status);
			
			for (size_t i = 0; i < rhs.GetCount(); i++)
			{
				AIArtHandle art = rhs[i];
				status = sAIArtSet->AddArtToArtSet(this->fArtSet, art);
				aisdk::check_ai_error(status);
			}
		}
		catch (ai::Error& ex)
		{
			status = ex;
			if (this->fArtSet != NULL)
				sAIArtSet->DisposeArtSet(&this->fArtSet);
		}
    }
	SDK_ASSERT(!status);
    return *this;    // Return ref for multiple assignment
}//end operator=

/* index operator
*/
AIArtHandle SnpArtSetHelper::operator [] (size_t index) const
{
	AIArtHandle result = NULL;
	if (index < this->GetCount()) {
		ASErr status = sAIArtSet->IndexArtSet(this->fArtSet, index, &result);
		SDK_ASSERT(!status);
	}
	else {
		SDK_ASSERT_MSG(false, "Invalid index");
	}
	return result;
}

/*
*/
size_t SnpArtSetHelper::GetCount() const
{
	size_t result = 0;
	ASErr status = sAIArtSet->CountArtSet(this->fArtSet, &result);
	SDK_ASSERT(!status);
	return result;
}

/* Finds set of selected art.
*/
void SnpArtSetHelper::FindSelectedArt(void)
{
	ASErr status = sAIArtSet->SelectedArtSet(this->fArtSet);
	SDK_ASSERT(!status);
}


/* Finds art on given layer.
*/
void SnpArtSetHelper::FindLayerArt(const ai::UnicodeString& layerName)
{
	ASErr status = kNoErr;
	AILayerHandle layer = NULL;
	status = sAILayer->GetLayerByTitle(&layer, layerName);
	aisdk::check_ai_error(status);

	status = sAIArtSet->LayerArtSet(layer, this->fArtSet);
	aisdk::check_ai_error(status);
}

/*	Returns a new art set containing art in this set that 
	matches a given type.
*/
SnpArtSetHelper SnpArtSetHelper::Filter(AIArtType include)
{
	SnpArtSetHelper result;
	ASErr status = kNoErr;
	for (size_t i = 0; i < this->GetCount(); i++) {
		AIArtHandle art = this->operator [](i);
		short type = kUnknownArt;
		status = sAIArt->GetArtType(art, &type);
		aisdk::check_ai_error(status);
		if (type == include) {
			result.Add(art);
		}
	}
	SDK_ASSERT(!status);
	return result;
}


/*
*/
void SnpArtSetHelper::Add(AIArtHandle art)
{
	ASErr status = sAIArtSet->AddArtToArtSet(this->fArtSet,art);
	SDK_ASSERT(!status);
}

/*
*/
void SnpArtSetHelper::SelectArt(AIArtHandle art)
{
	SnpSelectionHelper selectionHelper;
	if (art && selectionHelper.CanSelectArt(art)) {
		selectionHelper.SelectArt(art);
	}
}

/*
*/
void SnpArtSetHelper::SelectAll()
{
	SnpSelectionHelper selectionHelper;
	selectionHelper.DeselectAll();
	for (size_t i = 0; i < this->GetCount(); i++) {
		AIArtHandle art = this->operator [](i);
		if (art && selectionHelper.CanSelectArt(art)) {
			selectionHelper.SelectArt(art);
		}
	}
}

/*
*/
void SnpArtSetHelper::Dump()
{
	SnippetRunnerLog* log = SnippetRunnerLog::Instance();
	for (size_t i = 0; i < this->GetCount(); i++) {
		AIArtHandle art = this->operator [](i);
		if (art) {
			short type;
			sAIArt->GetArtType(art, &type);
			ai::UnicodeString name;
			sAIArt->GetArtName(art, name, NULL);
			log->WritePrintf("i: %d, type: %d", i, type);
			ai::UnicodeString msg("name: ");
			msg.append(name);
			log->Write(msg);
		}
	}
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpArtSetHelper available to the SnippetRunner framework.
 */
class _SnpRunnableArtSetHelper : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableArtSetHelper () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableArtSetHelper () {}

		/* Returns name of snippet.
		*/
		std::string GetName() const;

		/* Returns a description of what the snippet does.
		*/
		std::string	GetDescription() const;

		/** Returns the operations this snippet supports.
		*/
		Operations GetOperations() const;

		/** Returns name of the snippet's default operation.
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

		/** Runs custom unit test.
		*/
		ASErr RunUnitTest(SnpRunnable::Context& runnableContext);
};

/*
*/
std::string _SnpRunnableArtSetHelper::GetName() const
{
	return "ArtSetHelper";
}

/*
*/
std::string _SnpRunnableArtSetHelper::GetDescription() const 
{
	return "Helper class that finds sets of art in a document";
}

/*
*/
SnpRunnable::Operations _SnpRunnableArtSetHelper::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("FindSelectedArt", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("FindLayerArt", "document", kSnpRunNotSupportedContext));
	operations.push_back(Operation("FindArtType", "document", kSnpRunNotSupportedContext));
	return operations;
}

/*
*/
std::string _SnpRunnableArtSetHelper::GetDefaultOperationName() const 
{
	return "FindSelectedArt";
}

/*
*/
std::vector<std::string> _SnpRunnableArtSetHelper::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Helper Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableArtSetHelper::CanRun(SnpRunnable::Context& runnableContext)
{
	SnpSelectionHelper selectionHelper;
	return selectionHelper.IsDocumentSelected();
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableArtSetHelper::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	try {
		if ("FindSelectedArt" == runnableContext.GetOperation().GetName()) {
			SnpArtSetHelper instance;
			instance.FindSelectedArt();
			instance.Dump();
		}
		else if ("FindLayerArt" == runnableContext.GetOperation().GetName()) {
			SnpChooser chooser;
			AILayerHandle layer = chooser.PickLayer();
			if (layer) {
				ai::UnicodeString title;
				result = sAILayer->GetLayerTitle(layer, title);
				if (!result) {
					SnpArtSetHelper instance;
					instance.FindLayerArt(title);
					instance.Dump();
				}
			}
		}
		else if ("FindArtType" == runnableContext.GetOperation().GetName()) {

			SnpChooser chooser;
			AIArtType artType = chooser.PickArtType();
            AIArtSpec specs[1] = {{static_cast<ai::int16>(artType),0,0}};
			SnpArtSetHelper instance(specs, 1);
			instance.Dump();
		}
		else {
			result = kBadParameterErr;
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr _SnpRunnableArtSetHelper::RunUnitTest(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;

	if ("FindSelectedArt" == runnableContext.GetOperation().GetName()) {
		
		// Run all unit tests when this operation is called.
		
		result = this->SetupContext(runnableContext);

		try {
			SnpUnitTestHelper::Logger ut("All art in unittest-1.ai");
			AIArtSpec specs[1] = {{kAnyArt,0,0}};
			SnpArtSetHelper artSet(specs, 1);
			artSet.Dump();
			SDK_ASSERT(artSet.GetCount() > 0);
			ut.SetStatus(kNoErr);
		}
		catch (ai::Error& ex) {
			result = ex;
		}

		try {
			SnpUnitTestHelper::Logger ut("All art in layer layer.eps in unittest-1.ai");
			SnpArtSetHelper artSet;
			artSet.FindLayerArt(ai::UnicodeString("layer.eps"));
			artSet.Dump();
			SDK_ASSERT(artSet.GetCount() == 2);
			ut.SetStatus(kNoErr);
		}
		catch (ai::Error& ex) {
			result = ex;
		}

		try {
			SnpUnitTestHelper::Logger ut("Raster art in layer layer.jpg in unittest-1.ai");
			SnpArtSetHelper artSet;
			artSet.FindLayerArt(ai::UnicodeString("layer.jpg"));
			artSet = artSet.Filter(kRasterArt);
			artSet.Dump();
			SDK_ASSERT(artSet.GetCount() == 1);
			ut.SetStatus(kNoErr);
		}
		catch (ai::Error& ex) {
			result = ex;
		}

		try {
			SnpUnitTestHelper::Logger ut("Selected art in unittest-1.ai");
			SnpArtSetHelper jpgArtSet;
			jpgArtSet.FindLayerArt(ai::UnicodeString("layer.jpg"));
			jpgArtSet = jpgArtSet.Filter(kRasterArt);
			SDK_ASSERT(jpgArtSet.GetCount() == 1);
			SnpSelectionHelper selectionHelper;
			SDK_ASSERT(!selectionHelper.SelectArt(jpgArtSet [0]));
			SnpArtSetHelper artSet;
			artSet.FindSelectedArt();
			artSet.Dump();
			SDK_ASSERT(artSet.GetCount() == 1);
			ut.SetStatus(kNoErr);
		}
		catch (ai::Error& ex) {
			result = ex;
		}

		ASErr tmp = this->TeardownContext(runnableContext);
		if (tmp && !result) {
			result = tmp;
		}
	}
	else {
		SnippetRunnerLog::Instance()->Write("SKIPPED - already tested by FindSelectedArt");
	}

	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableArtSetHelper instance_SnpRunnableArtSetHelper;


// End SnpArtSetHelper.cpp

