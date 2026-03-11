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
#include "SnpLiveTrace.h"

/*
*/
ASErr SnpLiveTrace::CreateTracing()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		// Find selected objects which can be traced.
		SnpArtSetHelper artSet = this->FindCanTraceArtSet(true);
		if (artSet.GetCount() <= 0) {
			log->Write("No traceable art selected");
			throw ai::Error(kBadParameterErr);
		}
		else {
			// Get options to be used for tracings.
			SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();

			// Tracing Mode.
			std::vector<ai::UnicodeString> modeChoices;
			modeChoices.push_back(ai::UnicodeString("kAIVectorizeModeColor"));
			modeChoices.push_back(ai::UnicodeString("kAIVectorizeModeGray")); 
			modeChoices.push_back(ai::UnicodeString("kAIVectorizeModeBlackAndWhite"));
			VectorizeModes mode = static_cast<VectorizeModes>(parameter->GetChoice(ai::UnicodeString("Choose Tracing Mode"), modeChoices, 0));
			ai::UnicodeString modeStr("Mode: ");
			log->Write(modeStr.append(modeChoices[mode]));

			// Blur.
			AIReal blur = parameter->GetReal("Choose blur", 0.0, 0.0, 20.0);
			log->WritePrintf("Blur: %.1f px", blur);

			// Trace the art.
			for (size_t i = 0; i < artSet.GetCount(); i++) {
				AIArtHandle art = artSet[i];
				AIArtHandle tracedArt = NULL;
				result = this->CreateTracing(art, tracedArt, mode, blur);
			}
		}
	}
	catch (ai::Error& ex){
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveTrace::CreateTracing(const AIArtHandle& art, AIArtHandle& tracedArt, ai::int32 tracingModeKey,
								  AIReal blur, ai::int32 tracingMaxColorsKey, AIReal tracingResamplingKey)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(art);
		result = sAIVectorize->CreateTracing(kPlaceAbove, art, art, &tracedArt);
		aisdk::check_ai_error(result);

		// Optional tracing options, if no options are specified here the art will
		// be traced with default parameters.
		ai::Ref<AIDictionaryRef> options;
		result = sAIVectorize->AcquireTracingOptions(tracedArt, options << ai::Replace);
		aisdk::check_ai_error(result);

		result = sAIDictionary->SetIntegerEntry(options, sAIDictionary->Key(kTracingModeKey), tracingModeKey);
		aisdk::check_ai_error(result);

		//result = sAIDictionary->SetRealEntry(options,sAIDictionary->Key(kTracingBlurKey),blur);
		//aisdk::check_ai_error(result);

		result = sAIDictionary->SetIntegerEntry(options, sAIDictionary->Key(kTracingFullColorsFidelityKey), tracingMaxColorsKey);
		aisdk::check_ai_error(result);

		//result = sAIDictionary->SetRealEntry(options, sAIDictionary->Key(kTracingResampleResolutionKey), tracingResamplingKey);
		//aisdk::check_ai_error(result);

		// Embed swatches in the tracing result, should be done after creating a new
		// tracing or making any changes to an existing tracings' options.
		result = sAIVectorize->EmbedSwatches(tracedArt);
		aisdk::check_ai_error(result);

		// Mark the plugin art as dirty, so the application knows it has changed since
		// the last update.
		result = sAIPluginGroup->MarkPluginArtDirty(tracedArt);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex){
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveTrace::CreateTracingUsingOptionsFromTracingGroup()
{
	//	Use tracing options from an existing tracing group.
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();

		// Find selected objects that are tracings.
		SnpArtSetHelper tracingArtSet = this->FindIsTracingArtSet(true);
		if (tracingArtSet.GetCount() != 1) {
			log->Write("Incorrect number of tracing groups selected - select one tracing");
			throw ai::Error(kBadParameterErr);
		}

		// The tracingSrc object refers to the tracing whose options 
		// are to be applied when we create new tracings below.
		AIArtHandle tracingSrc = tracingArtSet[0];
		SDK_ASSERT(tracingSrc);

		// Find selected objects which can be traced.
		SnpArtSetHelper artSet = this->FindCanTraceArtSet(true);
		if (artSet.GetCount() <= 0) {
			log->Write("No traceable art selected");
			throw ai::Error(kBadParameterErr);
		}

		// Create the tracings.
		result = this->CreateTracingUsingOptionsFromTracingGroup(tracingSrc, artSet, tracingArtSet);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveTrace::CreateTracingUsingOptionsFromTracingGroup(AIArtHandle tracingSrc, const SnpArtSetHelper& artSet, SnpArtSetHelper& tracingArtSet)
{
	//	use tracing options from an existing tracing group
	ASErr result = kNoErr;
	try {
		for (size_t i = 0; i < artSet.GetCount(); i++) {	
			AIArtHandle art = artSet [i];
			AIArtHandle tracing = NULL;

			result = sAIVectorize->CreateTracing(kPlaceAbove, art, art, &tracing);
			aisdk::check_ai_error(result);

			ai::Ref<AIDictionaryRef> destOptions;
			result = sAIVectorize->AcquireTracingOptions(tracing, destOptions << ai::Replace);
			aisdk::check_ai_error(result);

			ai::Ref<AIDictionaryRef> srcOptions;
			result = sAIVectorize->AcquireTracingOptions(tracingSrc, srcOptions << ai::Replace);
			aisdk::check_ai_error(result);

			// Copy tracing options from source item to new tracing group.
			result = sAIDictionary->Copy(destOptions, srcOptions);
			aisdk::check_ai_error(result);

			// Embed swatches in the tracing result, should be done after creating a new
			// tracing or making any changes to an existing tracings' options.
			result = sAIVectorize->EmbedSwatches(tracing);
			aisdk::check_ai_error(result);

			// Mark the plugin art as dirty, so the application knows it has changed since
			// the last update.
			result = sAIPluginGroup->MarkPluginArtDirty(tracing);
			aisdk::check_ai_error(result);

			tracingArtSet.Add(tracing);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveTrace::ReleaseTracing()
{
	ASErr result = kNoErr;
	try {
		// Find selected objects that are tracings.
		SnpArtSetHelper tracingArtSet = this->FindIsTracingArtSet(true);
		if (tracingArtSet.GetCount() <= 0) {
			SnippetRunnerLog* log = SnippetRunnerLog::Instance();
			log->Write("No traced art selected");
			throw ai::Error(kBadParameterErr);
		}

		// Release the selected tracings.
		result = this->ReleaseTracing(tracingArtSet);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveTrace::ReleaseTracing(const SnpArtSetHelper& tracingArtSet)
{
	ASErr result = kNoErr;
	try {
		for (size_t i = 0; i < tracingArtSet.GetCount(); i++) {

			AIArtHandle tracing = tracingArtSet[i];
			SDK_ASSERT(tracing);

			// Find the soure art for the tracing.
			AIArtHandle art = NULL;
			result = sAIVectorize->GetSourceImage(tracing, &art);
			aisdk::check_ai_error(result);

			// Move the source art above the tracing.
			result = sAIArt->ReorderArt(art, kPlaceAbove, tracing);
			aisdk::check_ai_error(result);

			// Get rid of the tracing.
			result = sAIArt->DisposeArt(tracing);
			aisdk::check_ai_error(result);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveTrace::CopyTracing()
{
	ASErr result = kNoErr;
	try {	
		// Find selected objects that are tracings.
		SnpArtSetHelper tracingArtSet = this->FindIsTracingArtSet(true);
		if (tracingArtSet.GetCount() <= 0) {
			SnippetRunnerLog* log = SnippetRunnerLog::Instance();
			log->Write("No traced art selected");
			throw ai::Error(kBadParameterErr);
		}

		// Make a copy of each selected tracing.
		SnpArtSetHelper copyTracingArtSet;
		result = this->CopyTracing(tracingArtSet, copyTracingArtSet);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveTrace::CopyTracing(const SnpArtSetHelper& tracingArtSet, SnpArtSetHelper& copyTracingArtSet)
{
	ASErr result = kNoErr;
	try {
		// Make a copy of each tracing.
		for (size_t i = 0; i < tracingArtSet.GetCount(); i++) {
			AIArtHandle tracing = tracingArtSet[i];
			SDK_ASSERT(tracing);

			AIArtHandle art = NULL;
			result = sAIVectorize->CopyTracingArt(tracing, kPlaceAbove, tracing, &art, true);
			aisdk::check_ai_error(result);

			copyTracingArtSet.Add(art);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveTrace::PlaceAndCreateTracing()
{
	ASErr result = kNoErr;
	try {
		SnpDocumentHelper docHelper;
		if (!docHelper.CanPlaceFile()) {
			SnippetRunnerLog::Instance()->Write("Cannot place file");
			throw ai::Error(kBadParameterErr);
		}

		// Place a file.
		const ai::UnicodeString name; // Not set - so Place dialog is displayed.
		AIArtHandle placedArt = NULL;
		result = docHelper.PlaceFile(name, placedArt);
		aisdk::check_ai_error(result);

		// Trace the placed art, any non-traceable art placed will be caught in this->CreateTracing().
		result = this->CreateTracing();
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveTrace::InspectTracing()
{
	ASErr result = kNoErr;
	try {
		// Find selected objects that are tracings.
		SnpArtSetHelper tracingArtSet = this->FindIsTracingArtSet(true);
		if (tracingArtSet.GetCount() <= 0) {
			SnippetRunnerLog* log = SnippetRunnerLog::Instance();
			log->Write("No traced art selected");
			throw ai::Error(kBadParameterErr);
		}

		// Inspect their tracing info.
		result = this->InspectTracing(tracingArtSet);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpLiveTrace::InspectTracing(const SnpArtSetHelper& tracingArtSet)
{
	ASErr result = kNoErr;
	try {	
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();

		// Inspect each tracing and report information to the log.
		for (size_t i = 0; i < tracingArtSet.GetCount(); i++) {
			AIArtHandle tracing = tracingArtSet[i];
			
			// Get statistics from the tracing's dictionary.
			ai::Ref<AIDictionaryRef> stats;
			result = sAIVectorize->AcquireTracingStatistics(tracing,stats << ai::Replace);
			aisdk::check_ai_error(result);

			ai::int32 pathCountKey = 0;
			result = sAIDictionary->GetIntegerEntry(stats,sAIDictionary->Key(kTracingPathCountKey), &pathCountKey);
			aisdk::check_ai_error(result);

			ai::int32 anchorCountKey = 0;
			result = sAIDictionary->GetIntegerEntry(stats,sAIDictionary->Key(kTracingAnchorCountKey), &anchorCountKey);
			aisdk::check_ai_error(result);

			ai::int32 colorCountKey = 0;
			result = sAIDictionary->GetIntegerEntry(stats,sAIDictionary->Key(kTracingColorCountKey), &colorCountKey);
			aisdk::check_ai_error(result);

			ai::int32 areaCountKey = 0;
//			result = sAIDictionary->GetIntegerEntry(stats,sAIDictionary->Key(kTracingAreaCountKey), &areaCountKey);
			aisdk::check_ai_error(result);

			AIReal ppiKey = 0.0;
//			result = sAIDictionary->GetRealEntry(stats,sAIDictionary->Key(kTracingSourceImagePPIKey), &ppiKey);
			aisdk::check_ai_error(result);

			// Get native resolution.
			AIReal nativeDPI = 0.0;
			result = sAIVectorize->GetNativeResolution(tracing, &nativeDPI);
			aisdk::check_ai_error(result);

			// Get effective resolution.
			AIReal effectiveDPI = 0.0;
			result = sAIVectorize->GetEffectiveResolution(tracing, &effectiveDPI);
			aisdk::check_ai_error(result);

			ai::UnicodeString name;
			result = sAIArt->GetArtName(tracing, name, NULL);
			aisdk::check_ai_error(result);

			// Write information to the log.
			log->WritePrintf("Item %d:",i+1);
			SnippetRunnerLog::Indent indent;
			log->Write(name);
			log->WritePrintf("kTracingPathCountKey: %d", pathCountKey);
			log->WritePrintf("kTracingAnchorCountKey: %d", anchorCountKey);
			log->WritePrintf("kTracingColorCountKey: %d", colorCountKey);
			log->WritePrintf("kTracingAreaCountKey: %d", areaCountKey);
			log->WritePrintf("kTracingSourceImagePPIKey: %.1f px", ppiKey);
			log->WritePrintf("Native DPI: %.1f px", nativeDPI);
			log->WritePrintf("Effective DPI: %.1f px", effectiveDPI);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
SnpArtSetHelper SnpLiveTrace::FindCanTraceArtSet(ASBoolean selected)
{
	SnpArtSetHelper result;

	// Get a set of raster and placed art objects.
	ai::int32 attr = 0;
	if (selected) {
		attr = kArtSelected;
	}
	const short numSpecs = 2;
	AIArtSpec artSpec[numSpecs] = { 
		{kPlacedArt, attr, attr},
		{kRasterArt, attr, attr} 
	};
	SnpArtSetHelper artSet(artSpec, numSpecs);

	// The raster and placed art objects that are not already traced are the ones we want.
	for (size_t i = 0; i < artSet.GetCount(); i++) {
		AIArtHandle art = artSet[i];
		SDK_ASSERT(art);
		if (!this->IsArtInTracing(art)) {
			result.Add(art);
		}
	}
	return result;
}

/*
*/
SnpArtSetHelper SnpLiveTrace::FindIsTracingArtSet(ASBoolean selected)
{
	SnpArtSetHelper result;

	// Get the set of kPluginArt objects.
	ai::int32 attr = 0;
	if (selected) {
		attr = kArtSelected;
	}
	const short numSpecs = 1;
	AIArtSpec artSpec[numSpecs] = { 
		{kPluginArt, attr, attr}
	};
	SnpArtSetHelper artSet(artSpec, numSpecs);

	// The kPluginArt objects that are tracings are the ones that we want.
	for (size_t i = 0; i < artSet.GetCount(); i++) {
		AIArtHandle art = artSet[i];
		SDK_ASSERT(art);
		if (sAIVectorize->IsTracing(art)) {
			result.Add(art);
		}
	}
	return result;
}

/*
*/
ASBoolean SnpLiveTrace::IsArtInTracing(AIArtHandle art)
{
	ASBoolean result = false;
	ASErr status = kNoErr;
	while (art && status == kNoErr) {
		if (sAIVectorize->IsTracing(art)) {
			result = true;
			break;
		}
		status = sAIArt->GetArtParent(art, &art);
	}
	return result;
}

/*
*/
ASErr SnpLiveTrace::GetResultArt(AIArtHandle& tracedArt)
{
	ASErr result = kNoErr;
	try {
		// Get the result art from the tracing which AdjustObjectAIColors can access the colors of.
		AIArtHandle newArt = NULL;
		result = sAIVectorize->CopyTracingArt(tracedArt, kPlaceAbove, tracedArt, &newArt, true);
		aisdk::check_ai_error(result);
		if (newArt)
		{
			result = sAIArt->TransferAttributes(tracedArt, newArt, kTransferAll);
			aisdk::check_ai_error(result);
		}
		tracedArt = newArt;
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpLiveTrace available to the SnippetRunner framework.
 */
class _SnpRunnableLiveTrace : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableLiveTrace () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableLiveTrace () {}

		/* Returns name of snippet.
		*/
		std::string GetName() const;

		/* Returns a description of what the snippet does.
		*/
		std::string	GetDescription() const;
		
		/** Returns the operations this snippet supports.
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
		ASBoolean CanRun(SnpRunnable::Context& runnableContext);

		/* Runs the snippet.
			@param runnableContext see ISnpRunnableContext for documentation.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr Run(SnpRunnable::Context& runnableContext);

		/** Creates and selects tracings.
		*/
		ASErr SetupContext(SnpRunnable::Context& runnableContext);
};

/*
*/
std::string _SnpRunnableLiveTrace::GetName() const
{
	return "LiveTrace";
}

/*
*/
std::string _SnpRunnableLiveTrace::GetDescription() const 
{
	return "Performs basic Live Trace functions on selected items.";
}

/*
*/
SnpRunnable::Operations _SnpRunnableLiveTrace::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("CreateTracing", "image(s) to be traced selected", kSnpRunJPGContext));
	operations.push_back(Operation("CreateTracingUsingOptionsFromTracingGroup", "one tracing and one image to be traced selected", kSnpRunCustomContext));
	operations.push_back(Operation("ReleaseTracing", "tracing(s) selected", kSnpRunCustomContext));
	operations.push_back(Operation("CopyTracing", "tracing(s) selected", kSnpRunCustomContext));
	operations.push_back(Operation("PlaceAndCreateTracing", "document", kSnpRunNotSupportedContext)); // skip unit test because of interactive prompt
	operations.push_back(Operation("InspectTracing", "tracing(s) selected", kSnpRunCustomContext));
	return operations;
}

/*
*/
std::string _SnpRunnableLiveTrace::GetDefaultOperationName() const 
{
	return "CreateTracing";
}

/*
*/
std::vector<std::string> _SnpRunnableLiveTrace::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Tracing Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableLiveTrace::CanRun(SnpRunnable::Context& runnableContext)
{
	SnpSelectionHelper selectionHelper;
	if (!selectionHelper.IsDocumentSelected())
		return false;

	SnpLiveTrace instance;

	if ("ReleaseTracing" == runnableContext.GetOperation().GetName() ||
		"CopyTracing" == runnableContext.GetOperation().GetName() ||
		"InspectTracing" == runnableContext.GetOperation().GetName()) {
		
			SnpArtSetHelper tracings = instance.FindIsTracingArtSet(true);
			if (tracings.GetCount() > 0)
				return true;
			else 
				return false;
	}
	else if ("CreateTracing" == runnableContext.GetOperation().GetName()) {
		SnpArtSetHelper art = instance.FindCanTraceArtSet(true);
		if (art.GetCount() > 0 && (selectionHelper.IsArtSelected(kRasterArt) || 
			selectionHelper.IsArtSelected(kPlacedArt)))
				return true;
			else 
				return false;
	}
	else if ("CreateTracingUsingOptionsFromTracingGroup" == runnableContext.GetOperation().GetName()) {
		SnpArtSetHelper tracings = instance.FindIsTracingArtSet(true);
		SnpArtSetHelper art = instance.FindCanTraceArtSet(true);
		if (tracings.GetCount() == 1 && art.GetCount() >= 1)
			return true;
		else
			return false;
	}
	else
		return true;
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableLiveTrace::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpLiveTrace instance;
	if ("CreateTracing" == runnableContext.GetOperation().GetName()) {
		result = instance.CreateTracing();
	}
	else if ("CreateTracingUsingOptionsFromTracingGroup" == runnableContext.GetOperation().GetName()) {
		result = instance.CreateTracingUsingOptionsFromTracingGroup();
	}
	else if ("ReleaseTracing" == runnableContext.GetOperation().GetName()) {
		result = instance.ReleaseTracing();
	}
	else if ("CopyTracing" == runnableContext.GetOperation().GetName()) {
		result = instance.CopyTracing();
	}
	else if ("PlaceAndCreateTracing" == runnableContext.GetOperation().GetName()) {
		result = instance.PlaceAndCreateTracing();
	}
	else if ("InspectTracing" == runnableContext.GetOperation().GetName()) {
		result = instance.InspectTracing();
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

/*
*/
ASErr _SnpRunnableLiveTrace::SetupContext(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	try {
		switch (runnableContext.GetOperation().GetUnitTestContext()){
			case SnpRunnable::kSnpRunCustomContext: {
					SnpLiveTrace instance;
					SnpDocumentHelper docHelper;
					AIDocumentHandle document = NULL;
					SDK_ASSERT(!docHelper.OpenDocument(ai::UnicodeString("unittest-1.ai"), document));
					SDK_ASSERT(document);
					if ("CreateTracingUsingOptionsFromTracingGroup" == runnableContext.GetOperation().GetName()) {
						// Trace the first traceable object found in UnitTest doc.
						SnpArtSetHelper artSet = instance.FindCanTraceArtSet(false);
						artSet.SelectArt(artSet[0]);
						result = instance.CreateTracing();
						aisdk::check_ai_error(result);
						
						// Select the tracing.
						SnpSelectionHelper selectionHelper;
						selectionHelper.DeselectAll();
						SnpArtSetHelper tracingArtSet = instance.FindIsTracingArtSet(false);
						selectionHelper.SelectArt(tracingArtSet[0]);

						// Select another traceable object.
						selectionHelper.SelectArt(artSet[1]);
					}
					else {
						// Create and select some tracings.
						SnpArtSetHelper artSet = instance.FindCanTraceArtSet(false);
						artSet.SelectAll();
						result = instance.CreateTracing();
						aisdk::check_ai_error(result);
						SnpArtSetHelper tracingArtSet = instance.FindIsTracingArtSet(false);
						tracingArtSet.SelectAll();
					}
				}
				break;
			default: {
				result = SnpRunnable::SetupContext(runnableContext);
			}
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableLiveTrace instance_SnpRunnableLiveTrace;

// End SnpLiveTrace.cpp


