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
#include "SnpArtHelper.h"

// Framework includes:
#include "SnpRunnable.h"
#include "SnippetRunnerSuites.h"
#include "SnippetRunnerLog.h"
#include "SnippetRunnerParameter.h"
#include "SnpSelectionHelper.h"

/*
*/
ASErr SnpArtHelper::NewRectangle(AIArtHandle& artHandle)
{
	ASErr result = kNoErr;
	try {
		// Set the default position and add the rectangle.
		const AIRealRect rect = {100, 100, 200, 200};
		result = this->NewRectangle(rect, false, artHandle);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpArtHelper::NewRectangle(const AIRealRect rect, const AIBoolean reversed, AIArtHandle &artHandle)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIShapeConstruction);
		result = sAIShapeConstruction->NewRect(rect.top, rect.left, rect.bottom, rect.right, reversed, &artHandle);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpArtHelper::NewRoundedRectangle(AIArtHandle& artHandle)
{
	ASErr result = kNoErr;
	try {
		// Set the default position and add the rectangle.
		const AIRealRect rect = {100, 100, 200, 200};
		result = this->NewRoundedRectangle(rect, 5.0, 5.0, false, artHandle);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpArtHelper::NewRoundedRectangle(const AIRealRect rect, const AIReal horizRadius, const AIReal vertRadius, const AIBoolean reversed, AIArtHandle &artHandle)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIShapeConstruction);
		result = sAIShapeConstruction->NewRoundedRect(rect.top, rect.left, rect.bottom, rect.right, horizRadius, vertRadius, reversed, &artHandle);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpArtHelper::NewOval(AIArtHandle& artHandle)
{
	ASErr result = kNoErr;
	try {
		// Set the default position and add the oval.
		const AIRealRect rect = {100, 100, 200, 200};
		result = this->NewOval(rect, false, artHandle);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpArtHelper::NewOval(const AIRealRect rect, const AIBoolean reversed, AIArtHandle &artHandle)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIShapeConstruction);
		result = sAIShapeConstruction->NewCircumscribedOval(rect.top, rect.left, rect.bottom, rect.right, reversed, &artHandle);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpArtHelper::NewPolygon(AIArtHandle& artHandle)
{
	ASErr result = kNoErr;
	try {
		// Set the default parameters and add the polygon.
		result = this->NewPolygon(5, 100, 100, 50, false, artHandle);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpArtHelper::NewPolygon(const short numSides, const AIReal centerX, const AIReal centerY, const AIReal radius, const AIBoolean reversed, AIArtHandle &artHandle)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIShapeConstruction);
		result = sAIShapeConstruction->NewRegularPolygon(numSides, centerX, centerY, radius, reversed, &artHandle);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpArtHelper::NewStar(AIArtHandle& artHandle)
{
	ASErr result = kNoErr;
	try {
		// Set the default parameters and add the polygon.
		result = this->NewStar(5, 100, 100, 50, 100, false, artHandle);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpArtHelper::NewStar(const short numPoints, const AIReal centerX, const AIReal centerY, const AIReal innerRadius, const AIReal outerRadius, const AIBoolean reversed, AIArtHandle &artHandle)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIShapeConstruction);
		result = sAIShapeConstruction->NewStar(numPoints, centerX, centerY, innerRadius, outerRadius, reversed, &artHandle);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpArtHelper::NewSpiral(AIArtHandle& artHandle)
{
	ASErr result = kNoErr;
	try {
		// Set the default parameters and add the spiral.
		AIRealPoint firstArcCenter = {300, 300};
		AIRealPoint start = {300, 350};
		result = this->NewSpiral(firstArcCenter, start, 80, 20, true, artHandle);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpArtHelper::NewSpiral(const AIRealPoint firstArcCenter, const AIRealPoint start, const AIReal decayPercent, const short numQuarterTurns, const AIBoolean clockwiseFromOutside, AIArtHandle &artHandle)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIShapeConstruction);
		result = sAIShapeConstruction->NewSpiral(firstArcCenter, start, decayPercent, numQuarterTurns, clockwiseFromOutside, &artHandle);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpArtHelper::NewArc(AIArtHandle& artHandle)
{
	ASErr result = kNoErr;
	try {
		// Create new art item in the document.
		result = sAIArt->NewArt(kPathArt, kPlaceAboveAll, NULL, &artHandle);
		aisdk::check_ai_error(result);

		// Define path segments.
		AIPathSegment segment = {{0,0},{0,0},{0,0},false};	
		segment.corner = false;
		short i = 0;
		
		segment.p.h = 300;
		segment.p.v = 300;
		segment.in.v = 200;
		segment.in.h = segment.out.h = segment.p.h;
		segment.out.v = 400;
		
		result = sAIPath->SetPathSegments(artHandle, i++, 1, &segment);
		aisdk::check_ai_error(result);

		segment.p.h = 400;
		segment.p.v = 400;
		segment.in = segment.out = segment.p;
		
		result = sAIPath->SetPathSegments(artHandle, i++, 1, &segment);
		aisdk::check_ai_error(result);

		segment.p.h = 500;
		segment.p.v = 300;
		segment.in.v = 400;
		segment.in.h = segment.out.h = segment.p.h;
		segment.out.v = 200;
		
		result = sAIPath->SetPathSegments(artHandle, i++, 1, &segment);
		aisdk::check_ai_error(result);

		// Keep path open.
		result = sAIPath->SetPathClosed(artHandle, false);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpArtHelper available to the SnippetRunner framework.
 */
class _SnpRunnableArtHelper : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableArtHelper () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableArtHelper () {}

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
std::string _SnpRunnableArtHelper::GetName() const
{
	return "ArtHelper";
}

/*
*/
std::string _SnpRunnableArtHelper::GetDescription() const 
{
	return "Helper class that adds art to a document";
}

/*
*/
SnpRunnable::Operations _SnpRunnableArtHelper::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("NewRectangle", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("NewRoundedRectangle", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("NewOval", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("NewPolygon", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("NewStar", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("NewSpiral", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("NewArc", "document", kSnpRunNewDocumentContext));
	return operations;
}

/*
*/
std::string _SnpRunnableArtHelper::GetDefaultOperationName() const 
{
	return "NewRectangle";
}

/*
*/
std::vector<std::string> _SnpRunnableArtHelper::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Helper Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableArtHelper::CanRun(SnpRunnable::Context& runnableContext)
{
	SnpSelectionHelper selectionHelper;
	return selectionHelper.IsDocumentSelected();
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableArtHelper::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpArtHelper instance;
	AIArtHandle artHandle;
	if ("NewRectangle" == runnableContext.GetOperation().GetName()) {
		result = instance.NewRectangle(artHandle);
		aisdk::check_ai_error(result);
		SDK_ASSERT(artHandle);
	}
	else if ("NewRoundedRectangle" == runnableContext.GetOperation().GetName()) {
		result = instance.NewRoundedRectangle(artHandle);
		aisdk::check_ai_error(result);
		SDK_ASSERT(artHandle);
	}
	else if ("NewOval" == runnableContext.GetOperation().GetName()) {
		result = instance.NewOval(artHandle);
		aisdk::check_ai_error(result);
		SDK_ASSERT(artHandle);
	}
	else if ("NewPolygon" == runnableContext.GetOperation().GetName()) {
		result = instance.NewPolygon(artHandle);
		aisdk::check_ai_error(result);
		SDK_ASSERT(artHandle);
	}
	else if ("NewStar" == runnableContext.GetOperation().GetName()) {
		result = instance.NewStar(artHandle);
		aisdk::check_ai_error(result);
		SDK_ASSERT(artHandle);
	}
	else if ("NewSpiral" == runnableContext.GetOperation().GetName()) {
		result = instance.NewSpiral(artHandle);
		aisdk::check_ai_error(result);
		SDK_ASSERT(artHandle);
	}
	else if ("NewArc" == runnableContext.GetOperation().GetName()) {
		result = instance.NewArc(artHandle);
		aisdk::check_ai_error(result);
		SDK_ASSERT(artHandle);
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableArtHelper instance_SnpRunnableArtHelper;

// End SnpArtHelper.cpp

