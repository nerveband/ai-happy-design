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

// Project includes:
#include "SDKErrors.h"
#include "SnpRunnable.h"
#include "SnippetRunnerLog.h"
#include "SnippetRunnerParameter.h"
#include "SnippetRunnerSuites.h"
#include "SnpChooser.h"

/*
*/
AIDocumentHandle SnpChooser::PickDocument()
{
	AIDocumentHandle document = NULL;
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIDocumentList);
		SDK_ASSERT(sAIDocument);
		std::vector<AIDocumentHandle> documents;
		std::vector<ai::UnicodeString> names;
		ai::int32 count = 0;
		result = sAIDocumentList->Count(&count);
		aisdk::check_ai_error(result);
		for (ai::int32 i = 0; i < count; i++) {
			AIDocumentHandle doc = NULL;
			result = sAIDocumentList->GetNthDocument(&doc, i);
			aisdk::check_ai_error(result);
			SDK_ASSERT(doc);
			result = sAIDocumentList->Activate(doc, true);
			aisdk::check_ai_error(result);
			ai::FilePath path;
			result = sAIDocument->GetDocumentFileSpecification(path);
			aisdk::check_ai_error(result);
			documents.push_back(doc);
			names.push_back(path.GetFileName());
		}
		if (count == 0) {
			result = kNoDocumentErr;
			SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
			parameter->Alert("Open some documents and try again");
		}
		else if (count == 1) {
			document = documents [0];
		}
		else {
			SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
			ASUInt32 c = parameter->GetChoice(ai::UnicodeString("Pick a document"), names, 0);
			document = documents [c];
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return document;
}

/*
*/
AILayerHandle SnpChooser::PickLayer()
{
	AILayerHandle layer = NULL;
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAILayer);
		std::vector<AILayerHandle> layers;
		std::vector<ai::UnicodeString> names;
		ai::int32 count = 0;
		result = sAILayer->CountLayers(&count);
		aisdk::check_ai_error(result);
		for (ai::int32 i = 0; i < count; i++) {
			AILayerHandle lyr = NULL;
			result = sAILayer->GetNthLayer(i, &lyr);
			aisdk::check_ai_error(result);
			SDK_ASSERT(lyr);
			ai::UnicodeString name;
			result = sAILayer->GetLayerTitle(lyr, name);
			aisdk::check_ai_error(result);
			layers.push_back(lyr);
			names.push_back(name);
		}
		if (count == 0) {
			result = kNoDocumentErr;
			SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
			parameter->Alert("Create some layers and try again");
		}
		else if (count == 1) {
			layer = layers [0];
		}
		else {
			SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
			ASUInt32 c = parameter->GetChoice(ai::UnicodeString("Pick a layer"), names, 0);
			layer = layers [c];
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return layer;
}

/*
*/
AIArtType SnpChooser::PickArtType()
{
	AIArtType value = kGroupArt;
	ASErr result = kNoErr;
	try {
		std::vector<AIArtType> values;
		values.push_back(kAnyArt);
		values.push_back(kGroupArt);
		values.push_back(kPathArt);
		values.push_back(kCompoundPathArt);
		values.push_back(kPlacedArt);
		values.push_back(kMysteryPathArt);
		values.push_back(kRasterArt);
		values.push_back(kPluginArt);
		values.push_back(kMeshArt);
		values.push_back(kTextFrameArt);
		values.push_back(kSymbolArt);
		values.push_back(kForeignArt);
		values.push_back(kLegacyTextArt);
		std::vector<ai::UnicodeString> names;
		names.push_back(ai::UnicodeString("kAnyArt"));
		names.push_back(ai::UnicodeString("kGroupArt"));
		names.push_back(ai::UnicodeString("kPathArt"));
		names.push_back(ai::UnicodeString("kCompoundPathArt"));
		names.push_back(ai::UnicodeString("kPlacedArt"));
		names.push_back(ai::UnicodeString("kMysteryPathArt"));
		names.push_back(ai::UnicodeString("kRasterArt"));
		names.push_back(ai::UnicodeString("kPluginArt"));
		names.push_back(ai::UnicodeString("kMeshArt"));
		names.push_back(ai::UnicodeString("kTextFrameArt"));
		names.push_back(ai::UnicodeString("kSymbolArt"));
		names.push_back(ai::UnicodeString("kForeignArt"));
		names.push_back(ai::UnicodeString("kLegacyTextArt"));
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ASUInt32 c = parameter->GetChoice(ai::UnicodeString("Pick an AIArtType"), names, 0);
		value = values [c];
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return value;
}

/*
*/
AIDocumentColorModelValue SnpChooser::PickDocumentColorModel()
{
	AIDocumentColorModelValue result = kDocCMYKColor;
	std::vector<std::string> choices;
	choices.push_back("kDocRGBColor");
	choices.push_back("kDocCMYKColor");
	ASUInt32 choice = SnippetRunnerParameter::Instance()->GetChoice("Pick Document Color Model", choices, 1);
	result = choice ? kDocCMYKColor : kDocRGBColor;
	return result;
}

/*
*/
AIDocumentRulerUnitValue SnpChooser::PickDocumentRulerUnit()
{
	AIDocumentRulerUnitValue result = kInchesUnits;
	std::vector<std::string> choices;
	choices.push_back("kUnknownUnits");
	choices.push_back("kInchesUnits");
	choices.push_back("kCentimetersUnits");
	choices.push_back("kPointsUnits");
	choices.push_back("kPicasUnits");
	choices.push_back("kMillimetersUnits");
	choices.push_back("kPixelsUnits");
	choices.push_back("kQUnits");
	ASUInt32 choice = SnippetRunnerParameter::Instance()->GetChoice("Pick Document Ruler Units", choices, 1);
	result = choice ? static_cast<AIDocumentRulerUnitValue>(choice) : result;
	return result;
}

/*
*/
AIRasterResolution SnpChooser::PickDocumentRasterResolution()
{
	AIRasterResolution result = kAIRasterResolutionScreen;
	std::vector<std::string> choices;
	choices.push_back("kAIRasterResolutionScreen (72 PPI)");
	choices.push_back("kAIRasterResolutionMedium (150 PPI)");
	choices.push_back("kAIRasterResolutionHigh (300 PPI)");
	result = static_cast<AIRasterResolution>(SnippetRunnerParameter::Instance()->GetChoice("Pick Raster Resolution", choices, 0));
	return result;
}

/*
*/
AIPreviewMode SnpChooser::PickPreviewMode()
{
	AIPreviewMode result = kAIPreviewModeDefault;
	std::vector<std::string> choices;
	choices.push_back("kAIPreviewModeDefault");
	choices.push_back("kAIPreviewModePixelPreview");
	choices.push_back("kAIPreviewModeOverprintPreview");
	result = static_cast<AIPreviewMode>(SnippetRunnerParameter::Instance()->GetChoice("Pick Preview Mode", choices, 0));
	return result;
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpChooser available to the SnippetRunner framework.
 */
class _SnpRunnableTypeChooser : public SnpRunnable
{
	public:
		/* Constructor registers the snippet to the framework.
		 */
		_SnpRunnableTypeChooser () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual	~_SnpRunnableTypeChooser () {}

		/* Returns name of snippet.
		*/
		std::string GetName() const;

		/* Returns a description of what the snippet does.
		*/
		std::string GetDescription() const;

		/** Returns operations supported by this snippet.
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
};

/*
*/
std::string _SnpRunnableTypeChooser::GetName() const
{
	return "Chooser";
}

/*
*/
std::string _SnpRunnableTypeChooser::GetDescription() const 
{
	return "Prompts for various Illustrator type values.";
}

/*
*/
SnpRunnable::Operations _SnpRunnableTypeChooser::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("PickDocument", "", kSnpRunNewDocumentContext));
	operations.push_back(Operation("PickLayer", "", kSnpRunNewDocumentContext));
	operations.push_back(Operation("PickArtType", "", kSnpRunAnyContext));
	operations.push_back(Operation("PickDocumentColorModel", "", kSnpRunAnyContext));
	operations.push_back(Operation("PickDocumentRulerUnit", "", kSnpRunAnyContext));
	operations.push_back(Operation("PickDocumentRasterResolution", "", kSnpRunAnyContext));
	operations.push_back(Operation("PickPreviewMode", "", kSnpRunAnyContext));
	return operations;
}

/*
*/
std::string _SnpRunnableTypeChooser::GetDefaultOperationName() const 
{
	return "PickDocument";
}

/*
*/
std::vector<std::string> _SnpRunnableTypeChooser::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Helper Snippets");
	return categories;
}

/*
*/
ASBoolean _SnpRunnableTypeChooser::CanRun(SnpRunnable::Context& runnableContext)
{
	return true;
}

/*
*/
ASErr _SnpRunnableTypeChooser::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr; 
	SnpChooser instance;
	if ("PickDocument" == runnableContext.GetOperation().GetName()) {
		instance.PickDocument();
	}
	else if ("PickLayer" == runnableContext.GetOperation().GetName()) {
		instance.PickLayer();
	}
	else if ("PickArtType" == runnableContext.GetOperation().GetName()) {
		instance.PickArtType();
	}
	else if ("PickDocumentColorModel" == runnableContext.GetOperation().GetName()) {
		instance.PickDocumentColorModel();
	}
	else if ("PickDocumentRulerUnit" == runnableContext.GetOperation().GetName()) {
		instance.PickDocumentRulerUnit();
	}
	else if ("PickDocumentRasterResolution" == runnableContext.GetOperation().GetName()) {
		instance.PickDocumentRasterResolution();
	}
	else if ("PickPreviewMode" == runnableContext.GetOperation().GetName()) {
		instance.PickPreviewMode();
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

/*	Declare an instance to register the snippet hook with the framework.
*/
static _SnpRunnableTypeChooser instance_SnpRunnableTypeChooser;

// End SnpChooser.cpp

