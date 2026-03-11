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
#include "SnpPatternHelper.h"
#include "SnpSelectionHelper.h"
#include "SnpArtHelper.h"
#include "SnpArtSetHelper.h"

/**
	Demonstrates the creation of a new pattern programmatically, inspecting the patterns avaialable with any document,
	filling a particular art object with a patten, and deleting patterns from the Swatch palette.
*/

class SnpPattern
{
	public:
		/**
			It shows how patterns available in Ai can be used to fill an art programmatically. 
			A rectangle art object is used in which the pattern gets inserted on selecting this option.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr ApplyPattern(void);

		/**
			This creates a sample pattern on a blank document. 
			This could be one of the  ways of creating patterns programmatically. 
			A new pattern is created out of that art, by selecting a rectangle path out of it( hidden), 
			and a new name is given to the pattern . The user can crop it from anywhere, and add to his swatch manually, 
			or he could add the pattern name to the swatch list through his code by using the pattern name. 
			The new pattern is added to the swatch palette after running this function.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr CreatePattern(void);

		/**
			Displays the patterns available in a document. 
			The user is prompted for the name of the pattern to be displayed.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr GetPattern(void);

		/**
			Prompt for the number of the pattern to edit (defualts to the first pattern)
			Then opens the given pattern in edit mode if edit mode isn't running
		*/
		ASErr EditPattern(void);

		/**
			If we are in pattern edit mode, then this quits is without saving changes made
		*/
		ASErr CancelPatternEditMode(void);
	
		/** 
			Save pattern being edited
		*/
		ASErr SavePatternEdits(void);
	
		/**
			Gives the total number of patterns, and number and name of each pattern available in an open document .
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr InspectPatterns(void);

		/**
			Changes the name of a pattern by name. 
			By default the first pattern name is changed to "My Pattern". 
			User is prompted for a new name and may change this for his selection .
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr ChangePatternName(void);
		
		/**
			Delete a pattern by Name.
			This by default deletes the first pattern available with Ai document. 
			The user is prompted for the pattern number to be deleted.
			Since these are 0-indexed, the choices available are 0 : number of patterns-1.
			If any specific pattern is to be deleted, its name can be found out through InspectPatterns() to know the index,
			and that can be passed as a parameter.
			So the user can call InspectPatterns first, and then can give the required pattern 
			name to be deleted through this function.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr DeletePatternByName(void);

		/**
			Delete All the Patterns in the document.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr DeleteAllPatterns(void);		
};

/*
*/
ASErr SnpPattern::ApplyPattern(void)
{
	
	ASErr result = kNoErr;
	try {
		AIPatternHandle pattern = NULL;
	
		//Draw a rectangle to retrieve the pattern in.
		AIRealRect rectbounds;
		rectbounds.left = 200;
		rectbounds.top =  500;
		rectbounds.right = 300;
		rectbounds.bottom = 600;

		SnpArtHelper artHelper;
		AIArtHandle art;
		result = artHelper.NewRectangle(rectbounds,false,art);
		aisdk::check_ai_error(result);			
			
		//Count the total number of swatch groups in the document.
		int numSwatchGroups = sAISwatchGroup->CountSwatchGroups(NULL);

		AISwatchRef swatchRef;
		//For each swatch group, get the number of swatches.
		for(int i=0; i<numSwatchGroups; i++){

			AISwatchGroupRef groupRef = sAISwatchGroup->GetNthSwatchGroup(NULL, i);
			int numSwatches = sAISwatchGroup->CountSwatches(groupRef);

				//Check for each swatch if it represents a pattern.
				for (int swatchIndex = 0; swatchIndex < numSwatches; swatchIndex++){
					if(( swatchRef = sAISwatchGroup->GetNthSwatch(groupRef,swatchIndex)) != NULL){
						AIColor color;
						result = sAISwatchList->GetAIColor(swatchRef, &color);
						aisdk::check_ai_error(result);

						if(color.kind == kPattern){
							ai::UnicodeString swatchName;
							ai::int32 count;
							result = sAIPattern->CountPatterns(&count);
							aisdk::check_ai_error(result);
							result = sAIPattern->GetNthPattern(count-1,&pattern);
							aisdk::check_ai_error(result);

							sAIPattern->ValidatePattern(pattern);
							
							SnpArtSetHelper ArtSetHelper;
							//Select the art before filling with pattern.
							ArtSetHelper.SelectArt(art);
							short type;
							result = sAIArt->GetArtType(art, &type);
							aisdk::check_ai_error(result);

							//Apply pattern to the selected art.
							if(type == kPathArt){
								AIPathStyle style;
								result = sAIPathStyle->GetPathStyle(art,&style, nullptr);
								aisdk::check_ai_error(result);
								AIPathStyleMap pathStyle;
								result = sAIPathStyle->SetCurrentPathStyle(&style, &pathStyle, NULL);
								aisdk::check_ai_error(result);
								result = sAIPathStyle->ApplyAIColorToPathStyle(&style,NULL,&color,NULL,false,false,art,NULL);
								aisdk::check_ai_error(result);
								result = sAIPathStyle->SetPathStyle(art,&style);
								aisdk::check_ai_error(result);
							}										
						}
					}
				}
			}
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpPattern :: CreatePattern(void)
{
	
	ASErr result = kNoErr;
	try {
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
				
		AIArtHandle basePath = NULL;
		result = sAIArt->NewArt(kGroupArt, kPlaceAboveAll, basePath, &basePath);
		aisdk::check_ai_error(result);

		//Base rectangle to create the pattern. 		
		AIRealRect rectbounds;
		rectbounds.left = 200;
		rectbounds.top =  300;
		rectbounds.right = 300;
		rectbounds.bottom = 200;

		SnpArtHelper artHelper;
		AIArtHandle rectart;
		result = artHelper.NewRectangle(rectbounds,false,rectart);
		aisdk::check_ai_error(result);

		result = sAIArt->ReorderArt(rectart,kPlaceInsideOnTop,basePath);
		aisdk::check_ai_error(result);

		//Create Oval art for creating own sample pattern.
		AIRealPoint realpoint,source;
		source.h = rectbounds.left;
		source.v = rectbounds.bottom;
		realpoint.h = 25;
		realpoint.v = 25;
		AIArtHandle ovalArt;
				
		for(AIRealPoint xyVertSource = source;xyVertSource.v <= rectbounds.top;	xyVertSource.v += 25){

				for(AIRealPoint xysource = xyVertSource; xysource.h <= rectbounds.right; xysource.h += 25){

					result = sAIShapeConstruction->NewPointSizeOval(xysource,realpoint,true,false,false,&ovalArt);
					aisdk::check_ai_error(result);

					result = sAIArt->ReorderArt(ovalArt,kPlaceInsideOnTop,basePath);
					aisdk::check_ai_error(result);
					
					// Set style with no fill and stroke.
					AIPathStyle style;
					result = sAIPathStyle->GetInitialPathStyle(&style);
					aisdk::check_ai_error(result);

					style.fillPaint = true;
					style.fill.color.kind = kThreeColor;
					style.fill.color.c.rgb.red = 1.0;
					result = sAIPathStyle->SetPathStyle(ovalArt,&style);
					aisdk::check_ai_error(result);		
				}
		}
		//Dispose the helping rectange used for creating the pattern.
		result = sAIArt->DisposeArt(rectart);
		aisdk::check_ai_error(result);
		
		AIPatternHandle pattern;
		result = sAIPattern->NewPattern(&pattern);
		aisdk::check_ai_error(result);

		// Set style with no fill and stroke.
		AIPathStyle style;
		result = sAIPathStyle->GetPathStyle(basePath,&style, nullptr);
		aisdk::check_ai_error(result);
		style.fillPaint = false;
		style.strokePaint = false;
		style.fill.color.kind = kPattern;
		style.fill.color.c.p.pattern = pattern;
		result = sAIPathStyle->SetPathStyle(basePath,&style);
		aisdk::check_ai_error(result);
	
    	result = sAIPattern->SetPatternArt(pattern,basePath );
		aisdk::check_ai_error(result);

		//Set the name of the new Pattern.
		ai::UnicodeString name("My Pattern");
		result = sAIPattern->SetPatternName(pattern, name);
		aisdk::check_ai_error(result);
		log->Write(ai::UnicodeString("Pattern name: ").append(name));

		//Set AIColor for the new Pattern.
		AISwatchRef swatchRef = sAISwatchList->InsertNthSwatch(NULL,-1);
		AIColor color;
		color.kind = kPattern;
		color.c.p.pattern = pattern;
		color.c.p.shiftDist = kAIRealZero;
		color.c.p.shiftAngle = kAIRealZero;
		color.c.p.scale.h = 100;
		color.c.p.scale.v = 100;
		color.c.p.rotate = kAIRealZero;
		color.c.p.reflect = false;
		color.c.p.reflectAngle = kAIRealZero;
		color.c.p.shearAngle = kAIRealZero;
		color.c.p.shearAxis = kAIRealZero;
		
		sAIRealMath->AIRealMatrixSetIdentity(&color.c.p.transform);
		result = sAISwatchList->SetAIColor(swatchRef,&color);
		aisdk::check_ai_error(result);
			
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpPattern::EditPattern(void)
{
	ASErr result = kNoErr;
	
	// Get the log.
	SnippetRunnerLog* log = SnippetRunnerLog::Instance();					
	
	try
	{		
		// Prompt for the pattern the user wants to edit
		ai::int32 count;
		result = sAIPattern->CountPatterns(&count);
		aisdk::check_ai_error(result);
		
		ASInt32 choice = SnippetRunnerParameter::Instance()->GetInt("Select Pattern Number to get:", 0,0,count-1);
		AIPatternHandle pattern;
		result = sAIPattern->GetNthPattern(choice,&pattern);
		aisdk::check_ai_error(result);
		
		log->Write("Entering pattern edit mode");
		
		result = sAIPattern->EnterPatternEditMode(pattern);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex)
	{
		result = ex;
	}
	
	return result;
}

/*
*/
ASErr SnpPattern::CancelPatternEditMode(void)
{
	// Get the log.
	SnippetRunnerLog* log = SnippetRunnerLog::Instance();					
	
	ASErr result = kNoErr;
	AIPatternHandle patternBeingEdited = NULL;
	
	// This ought to be true, as it is part of the preconditions
	AIBoolean beingEdited = sAIPattern->GetPatternEditMode(&patternBeingEdited, NULL);
	
	if (beingEdited)
	{
		log->Write("Canceling pattern edit mode");
		result = sAIPattern->EndPatternEditMode(patternBeingEdited, false);
	}
	else
	{
		result = kCantHappenErr;
	}
	
	return result;
}

/*
*/
ASErr SnpPattern::SavePatternEdits(void)
{
	// Get the log.
	SnippetRunnerLog* log = SnippetRunnerLog::Instance();					
	
	ASErr result = kNoErr;
	AIPatternHandle patternBeingEdited = NULL;
	
	// This ought to be true, as it is part of the preconditions
	AIBoolean beingEdited = sAIPattern->GetPatternEditMode(&patternBeingEdited, NULL);
	
	if (beingEdited)
	{
		log->Write("Saving changes and exiting pattern edit mode");
		result = sAIPattern->EndPatternEditMode(patternBeingEdited, true);
	}
	else
	{
		result = kCantHappenErr;
	}
	
	return result;
}

/*
*/
ASErr SnpPattern::GetPattern(void)
{	
	ASErr result = kNoErr;
	try {
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();					
		
		ai::int32 count;
		result = sAIPattern->CountPatterns(&count);
		aisdk::check_ai_error(result);
		
		ASInt32 choice = SnippetRunnerParameter::Instance()->GetInt("Select Pattern Number to get:", 0,0,count-1);
		AIPatternHandle pattern;
		result = sAIPattern->GetNthPattern(choice,&pattern);
		aisdk::check_ai_error(result);
		result = sAIPattern->GetPattern(pattern);
		aisdk::check_ai_error(result);
		
		ai::UnicodeString name;	
		result = sAIPattern->GetPatternName(pattern, name);
		aisdk::check_ai_error(result);
		log->Write(ai::UnicodeString("Pattern name: ").append(name));
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/

ASErr SnpPattern :: InspectPatterns()
{
	ASErr result = kNoErr;
	try {		
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
			
		ai::int32 count;
		result = sAIPattern->CountPatterns(&count);
		aisdk::check_ai_error(result);

		log->WritePrintf("TotalPatternCount %d", count);
		for(ai::int32 i=0; i < count; i++)
		{
			log->WritePrintf("PatternCount %d", i);
			AIPatternHandle pattern;
			result = sAIPattern->GetNthPattern(i, &pattern);
			aisdk::check_ai_error(result);

			ai::UnicodeString name;
			result = sAIPattern->GetPatternName(pattern, name);
			aisdk::check_ai_error(result);
			log->Write(ai::UnicodeString("Pattern name: ").append(name));
		}
	}
	catch (ai::Error& ex)
	{
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpPattern::ChangePatternName(void)
{
	
	ASErr result = kNoErr;
	try {
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		log->Write("Changing Name of the first Pattern");
		SnippetRunnerLog::Indent indent;
			
		ai::int32 count;
		result = sAIPattern->CountPatterns(&count);
		aisdk::check_ai_error(result);
		
		if(count > 0){
			AIPatternHandle pattern;
			result = sAIPattern->GetNthPattern(0,&pattern);
			aisdk::check_ai_error(result);
			ai::UnicodeString name;
			result = sAIPattern->GetPatternName(pattern, name);
			aisdk::check_ai_error(result);

			log->Write(ai::UnicodeString("Pattern name before change: ").append(name));

			ai::UnicodeString newName = SnippetRunnerParameter::Instance()->GetString(ai::UnicodeString("Enter new name for the pattern"),ai::UnicodeString("Changed Pattern"));
			result = sAIPattern->SetPatternName(pattern,newName);
			aisdk::check_ai_error(result);
			log->Write(ai::UnicodeString("Pattern name after change: ").append(newName));						
		}		
	}	
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpPattern::DeletePatternByName(void)
{
	
	ASErr result = kNoErr;
	try {
		// Get input parameters one at a time.
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();

		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		
		ai::UnicodeString name;	
		std::vector<ai::UnicodeString> names;

		ai::int32 count;
		result = sAIPattern->CountPatterns(&count);
		aisdk::check_ai_error(result);

		if(count > 0){
			AIPatternHandle pattern;
			for(ai::int32 i = count-1; i >= 0; i--)
			{
				result = sAIPattern->GetNthPattern(i, &pattern);
				aisdk::check_ai_error(result);

				result = sAIPattern->GetPatternName(pattern, name);
				aisdk::check_ai_error(result);
			
				names.push_back(name);
				log->Write(ai::UnicodeString("Pattern name: ").append(name));
				
			}
			ASUInt32 choice = parameter->GetChoice(ai::UnicodeString("Select Pattern to delete:"), names,0);
			ai::UnicodeString chosenPatternName = names[choice];

			log->Write(ai::UnicodeString("Deleting Pattern name: ").append(chosenPatternName));
			result = sAIPattern->DeletePattern(pattern);
			aisdk::check_ai_error(result);
		}
		else
			log->Write("No pattern to be deleted");
	}
	catch (ai::Error& ex)
	{
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpPattern::DeleteAllPatterns(void)
{
	ASErr result = kNoErr;
	try {
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
	
		ai::int32 count;
		result = sAIPattern->CountPatterns(&count);
		aisdk::check_ai_error(result);

		log->WritePrintf("TotalPatternCount %d", count);
		for(ai::int32 i = count; i > 0; i--)
		{
			log->WritePrintf("PatternCount %d", i);

			AIPatternHandle pattern;
			result = sAIPattern->GetNthPattern(i-1, &pattern);
			aisdk::check_ai_error(result);

			ai::UnicodeString name;
			result = sAIPattern->GetPatternName(pattern, name);
			aisdk::check_ai_error(result);

			result = sAIPattern->DeletePattern(pattern);
			aisdk::check_ai_error(result);
			log->Write(ai::UnicodeString("Deleting Pattern name: ").append(name));
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpPattern available to the SnippetRunner framework.
 */
class _SnpRunnablePattern : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnablePattern () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnablePattern () {}

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
std::string _SnpRunnablePattern::GetName() const
{
	return "Pattern";
}

/*
*/
std::string _SnpRunnablePattern::GetDescription() const 
{
	return "Provides Pattern Operations";
}

/*
*/
SnpRunnable::Operations _SnpRunnablePattern::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("ApplyPattern", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("CreatePattern", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("GetPattern", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("InspectPatterns", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("ChangePatternName", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("DeletePatternByName", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("DeleteAllPatterns", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("EditPattern", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("CancelPatternEditMode", "document, pattern edit mode", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("SavePatternEdits", "document, pattern edit mode", kSnpRunUnitTestDocumentContext));
	return operations;
}

/*
*/
std::string _SnpRunnablePattern::GetDefaultOperationName() const 
{
	return "ApplyPattern";
}

/*
*/
std::vector<std::string> _SnpRunnablePattern::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Color Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnablePattern::CanRun(SnpRunnable::Context& runnableContext)
{
	SnpSelectionHelper selectionHelper;
	if (!selectionHelper.IsDocumentSelected())
		return false;
	
	SnpPatternHelper patternHelper;	
	if ("EditPattern" == runnableContext.GetOperation().GetName())
		return !patternHelper.IsInPatternEditMode();
	
	if ("CancelPatternEditMode" == runnableContext.GetOperation().GetName()
		|| "SavePatternEdits" == runnableContext.GetOperation().GetName())
		return patternHelper.IsInPatternEditMode();
	
	return true;
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnablePattern::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpPattern instance;
	if ("ApplyPattern" == runnableContext.GetOperation().GetName()) {
		result = instance.ApplyPattern();
	}
	else if ("CreatePattern" == runnableContext.GetOperation().GetName()) {
		result = instance.CreatePattern();
	}
	else if ("GetPattern" == runnableContext.GetOperation().GetName()) {
		result = instance.GetPattern();
	}
	else if ("InspectPatterns" == runnableContext.GetOperation().GetName()) {
		result = instance.InspectPatterns();
	}
	else if ("ChangePatternName" == runnableContext.GetOperation().GetName()) {
		result = instance.ChangePatternName();
	}
	else if ("DeletePatternByName" == runnableContext.GetOperation().GetName()) {
		result = instance.DeletePatternByName();
	}
	else if ("DeleteAllPatterns" == runnableContext.GetOperation().GetName()) {
		result = instance.DeleteAllPatterns();
	}
	else if ("EditPattern" == runnableContext.GetOperation().GetName()) {
		result = instance.EditPattern();
	}
	else if ("CancelPatternEditMode" == runnableContext.GetOperation().GetName()) {
		result = instance.CancelPatternEditMode();
	}
	else if ("SavePatternEdits" == runnableContext.GetOperation().GetName()) {
		result = instance.SavePatternEdits();
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnablePattern instance_SnpRunnablePattern;

// End SnpPattern.cpp

