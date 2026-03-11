//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/CodeSnippets/SnpSymbolHelper.cpp $
//
//  $Revision: #1 $
//
//  Copyright 2010 Adobe Systems Incorporated. All rights reserved.
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
#include "SnpSelectionHelper.h"

#include "SnpSymbolHelper.h"
#include <list>

/*
*/
ai::int32 SnpSymbolHelper::CountSymbols()
{
	ai::int32 nCount = 0;
	try
	{
		ASErr result = sAISymbol->CountSymbolPatterns(&nCount, false);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error&)
	{
		nCount = 0;
	}

	return nCount;
}

/*
*/
ASErr SnpSymbolHelper::AddSymbol()
{
	ASErr result = kNoErr;
	try {
		// Get values from user.
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		AIPatternHandle symbolPattern;
		result = sAISymbol->NewSymbolPatternFromSel(&symbolPattern, kSymbolCenterPoint);
		aisdk::check_ai_error(result);
		ai::UnicodeString name = parameter->GetString(ai::UnicodeString("Enter Symbol new name:"), ai::UnicodeString("Symbol"));
		result = sAISymbol->SetSymbolPatternName(symbolPattern, name);

		// Define export type choices.
		std::vector<ai::UnicodeString> exportTypes;
		exportTypes.push_back(ai::UnicodeString("Graphic"));
		exportTypes.push_back(ai::UnicodeString("Movie Clip"));

		ASUInt32 exportType = parameter->GetChoice(ai::UnicodeString("Export Types:"), exportTypes, 0);

		if (exportType == 0) 
		{ // Graphic
			result = sAISymbol->SetSymbolExportType(symbolPattern, kGraphicType);
			aisdk::check_ai_error(result);
		}
		else if(exportType == 1)
		{ 
			// kMovieClipType.
			result = sAISymbol->SetSymbolExportType(symbolPattern, kMovieClipType);
			aisdk::check_ai_error(result);
		}
		else 
		{
			aisdk::check_ai_error(kBadParameterErr);
		}
	
		ASBoolean bEnableGuides = parameter->GetBoolean(ai::UnicodeString("Enable Guides?"));
		result = sAISymbol->SetSymbolEnableGuides(symbolPattern, bEnableGuides);
		aisdk::check_ai_error(result);
		ASBoolean bPixelPerfect = parameter->GetBoolean(ai::UnicodeString("Pixel Perfect?"));
		result = sAISymbol->SetSymbolPixelPerfect(symbolPattern, bPixelPerfect);
		aisdk::check_ai_error(result);

		AIRealRect slice9Grid = {0, 0, 0, 0}; // left, top, right, bottom
		slice9Grid.left = parameter->GetReal(ai::UnicodeString("Enter left boundary of new symbol:"), 100.0);
		slice9Grid.top = parameter->GetReal(ai::UnicodeString("Enter top boundary of new symbol:"), 200.0);
		slice9Grid.right = parameter->GetReal(ai::UnicodeString("Enter right boundary of new symbol:"), 200.0);
		slice9Grid.bottom = parameter->GetReal(ai::UnicodeString("Enter bottom boundary of new symbol:"), 100.0);
		result = sAISymbol->SetSymbol9SliceGrid(symbolPattern, slice9Grid);
		aisdk::check_ai_error(result);

		SnippetRunnerLog::Instance()->WritePrintf("A new symbol is added.");
		SnippetRunnerLog::Instance()->WritePrintf("");

	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpSymbolHelper::DeleteSymbol()
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		
		AIPatternHandle symbolPattern;
		ai::UnicodeString chosenPatternName;
		
		result = this->ChooseSymbol(&chosenPatternName, NULL);
		aisdk::check_ai_error(result);

		log->Write(ai::UnicodeString("Delete the symbol named: ").append(chosenPatternName));

		sAISymbol->GetSymbolPatternByName(chosenPatternName, &symbolPattern);
		aisdk::check_ai_error(result);

		result = sAISymbol->DeleteSymbolPattern(symbolPattern);
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
ASErr SnpSymbolHelper::RenameSymbol()
{
	ASErr result = kNoErr;
	try {
		// Get values from user.
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		ai::UnicodeString chosenPatternName;
		result = ChooseSymbol(&chosenPatternName, NULL);

		AIPatternHandle symbolPattern;
		sAISymbol->GetSymbolPatternByName(chosenPatternName, &symbolPattern);
		aisdk::check_ai_error(result);

		ai::UnicodeString name = parameter->GetString(ai::UnicodeString("Enter symbol new name:"), ai::UnicodeString("Symbol"));
		result = sAISymbol->SetSymbolPatternName(symbolPattern, name);
		log->Write(ai::UnicodeString("Rename the symbol named").append(chosenPatternName).append(ai::UnicodeString(" as ")).append(name));
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
		SnippetRunnerLog::Instance()->Write("Can not rename. The name had been used.");
	}	
	return result;
}

/*
*/
ASErr SnpSymbolHelper::PlaceSymbolInDocument()
{
	ASErr result = kNoErr;
	try 
	{
		ai::int32 symbolIndex = 0;
		result = ChooseSymbol(NULL, &symbolIndex);
		AIPatternHandle symbolPattern;
		result = sAISymbol->GetNthSymbolPattern(symbolIndex, &symbolPattern, false);
		aisdk::check_ai_error(result);
		AIArtHandle art;
		result = sAISymbol->NewInstanceCenteredInView(symbolPattern, kPlaceBelowAll, NULL, &art);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpSymbolHelper::DisplaySymbolProperties()
{
	ASErr result = kNoErr;
	try {
		ai::int32 symbolIndex = 0;
		result = ChooseSymbol(NULL, &symbolIndex);
		aisdk::check_ai_error(result);
		result = this->DisplaySymbolProperties(symbolIndex);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpSymbolHelper::DisplaySymbolProperties(ai::int32 symbolIndex)
{
	ASErr result = kNoErr;
	SnippetRunnerLog* log = SnippetRunnerLog::Instance();
	try {
		// Get values from user.
		AIPatternHandle symbolPattern;
		result = sAISymbol->GetNthSymbolPattern(symbolIndex, &symbolPattern, false);
		aisdk::check_ai_error(result);
		ai::UnicodeString name;
		result = sAISymbol->GetSymbolPatternName(symbolPattern, name);
		aisdk::check_ai_error(result);
		log->Write(ai::UnicodeString("Symbol name = ").append(name));

		AISymbolFlashExportType exportType;
		result = sAISymbol->GetSymbolExportType(symbolPattern, &exportType);
		aisdk::check_ai_error(result);
		(kGraphicType == exportType) ? log->Write("The symbol export type = Graphic.") : log->Write("The symbol export type = Movie Clip");
		
		AIBoolean isEnableGuides;
		result = sAISymbol->GetSymbolEnableGuides(symbolPattern, &isEnableGuides);
		aisdk::check_ai_error(result);
		isEnableGuides ? log->Write("The symbol is enable guides for 9-slice scaling.") : log->Write("The symbol is not enable guides for 9-slice scaling.");
		
		AIBoolean isEnablePixelPerfect;
		result = sAISymbol->GetSymbolEnablePixelPerfect(symbolPattern, &isEnablePixelPerfect);
		aisdk::check_ai_error(result);
		isEnablePixelPerfect ? log->Write("The symbol is pixel perfect.") : log->Write("The symbol is not pixel perfect.");
				
		AIBoolean isListed = sAISymbol->IsSymbolPatternListed(symbolPattern);
		isListed ? log->Write("The symbol is listed in current document.") : log->Write("The symbol is not listed in current document.");
		
		AIBoolean isEdited = sAISymbol->EditingSymbolPattern(symbolPattern);
		isEdited ? log->Write("The symbol is being edited.") : log->Write("The symbol is not being edited.");

	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}



/*
*/
ASErr SnpSymbolHelper::Display9SliceScaledInfo()
{
	ASErr result = kNoErr;
	try {
		ai::int32 symbolIndex = 0;
		result = ChooseSymbol(NULL, &symbolIndex);
		aisdk::check_ai_error(result);
		result = this->Display9SliceScaledInfo(symbolIndex);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpSymbolHelper::Display9SliceScaledInfo(ai::int32 symbolIndex)
{
	ASErr result = kNoErr;
	SnippetRunnerLog* log = SnippetRunnerLog::Instance();
	try {
		AIPatternHandle symbolPattern;
		result = sAISymbol->GetNthSymbolPattern(symbolIndex, &symbolPattern, false);
		aisdk::check_ai_error(result);
		AIRealRect symbolRect; 
		result = sAISymbol->GetSymbol9SliceGrid(symbolPattern, &symbolRect);
		aisdk::check_ai_error(result);
		if(!sAIRealMath->AIRealRectEmpty(&symbolRect) )
		{
			log->WritePrintf("The symbol is 9-Slice Scaled");
			log->WritePrintf("Left boundary = %.2f pts", symbolRect.left);
			log->WritePrintf("Top boundary = %.2f pts", symbolRect.top);
			log->WritePrintf("Right boundary = %.2f pts", symbolRect.right);
			log->WritePrintf("Bottom boundary = %.2f pts", symbolRect.bottom);
			aisdk::check_ai_error(result);
		}
		else
		{
			log->WritePrintf("The symbol is not 9-Slice Scaled");
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

ASErr SnpSymbolHelper::ChooseSymbol(ai::UnicodeString *name, ai::int32 *symbolIndex)
{
	if(name == NULL && symbolIndex == NULL)
		return kBadParameterErr;

	ASErr result = kNoErr;
	ai::UnicodeString symbolName;
	std::vector<ai::UnicodeString> names;
	ai::int32 symbolCount = this->CountSymbols();
	AIPatternHandle symbolPattern;
	for(ai::int32 i = 0; i < symbolCount; i++)
	{
		result = sAISymbol->GetNthSymbolPattern(i, &symbolPattern, false);
		aisdk::check_ai_error(result);
		
		result = sAISymbol->GetSymbolPatternName(symbolPattern, symbolName);
		aisdk::check_ai_error(result);
		names.push_back(symbolName);
	}
	SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
	ASUInt32 choice = parameter->GetChoice(ai::UnicodeString("Select symbol:"), names, 0);
	if(name) { *name = names[choice]; }
	if(symbolIndex) { *symbolIndex = choice; }

	return result;
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpSymbolHelper available to the SnippetRunner framework.
 */
class _SnpRunnableSymbolHelper : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableSymbolHelper () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableSymbolHelper () {}

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
std::string _SnpRunnableSymbolHelper::GetName() const
{
	return "SymbolHelper";
}

/*
*/
std::string _SnpRunnableSymbolHelper::GetDescription() const 
{
	return "Demonstrates how to manipulate Symbols";
}

/*
*/
SnpRunnable::Operations _SnpRunnableSymbolHelper::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("AddSymbol", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("DeleteSymbol", "document with at least 2 Symbols", kSnpRunNewDocumentContext));
	operations.push_back(Operation("RenameSymbol", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("PlaceSymbolInDocument", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("DisplaySymbolProperties", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("Display9SliceScaledInfo", "document", kSnpRunNewDocumentContext));
	return operations;
}

/*
*/
std::string _SnpRunnableSymbolHelper::GetDefaultOperationName() const 
{
	return "AddSymbol";
}

/*
*/
std::vector<std::string> _SnpRunnableSymbolHelper::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Helper Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableSymbolHelper::CanRun(SnpRunnable::Context& runnableContext)
{
	SnpSelectionHelper selectionHelper;
	SnpSymbolHelper SymbolHelper;
	if("AddSymbol" == runnableContext.GetOperation().GetName())
	{
		return selectionHelper.IsDocumentSelected();
	}
	else if("DeleteSymbol" == runnableContext.GetOperation().GetName() ||
		"RenameSymbol" == runnableContext.GetOperation().GetName() || 
		"PlaceSymbolInDocument" == runnableContext.GetOperation().GetName() || 
		"DisplaySymbolProperties" == runnableContext.GetOperation().GetName() ||
		"Display9SliceScaledInfo" == runnableContext.GetOperation().GetName()) 
	{
		if(selectionHelper.IsDocumentSelected() && (SymbolHelper.CountSymbols() > 0))
		{
			return true;
		}
	}
	return false;
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableSymbolHelper::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpSymbolHelper instance;
	if ("AddSymbol" == runnableContext.GetOperation().GetName()) {
		result = instance.AddSymbol();
	}
	else if ("DeleteSymbol" == runnableContext.GetOperation().GetName()) {
		result = instance.DeleteSymbol();
	}
	else if ("RenameSymbol" == runnableContext.GetOperation().GetName()) {
		result = instance.RenameSymbol();
	}
	else if ("PlaceSymbolInDocument" == runnableContext.GetOperation().GetName()) {
		result = instance.PlaceSymbolInDocument();
	}
	else if ("DisplaySymbolProperties" == runnableContext.GetOperation().GetName()) {
		result = instance.DisplaySymbolProperties();
	}
	else if ("Display9SliceScaledInfo" == runnableContext.GetOperation().GetName()) {
		result = instance.Display9SliceScaledInfo();
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableSymbolHelper instance_SnpRunnableSymbolHelper;

// End SnpSymbolHelper.cpp

