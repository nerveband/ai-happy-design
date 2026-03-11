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
#include "AIDocumentAction.h"
#include "AINativeAction.h"
#include "AIDxfDwgPrefs.h"
#include "AIPhotoshopPrefs.h"

#include "SDKErrors.h"

// Framework includes:
#include "SnpRunnable.h"
#include "SnippetRunnerSuites.h"
#include "SnippetRunnerLog.h"
#include "SnippetRunnerParameter.h"
#include "SnpDocumentHelper.h"
#include "SnpChooser.h"
#include "SnpDocumentActionHelper.h"
#include "SnpSelectionHelper.h"
#include "SnpArtSetHelper.h"
#include "SnpArtHelper.h"
#include "SnpObjectActionHelper.h"

/*
*/
SnpDocumentActionHelper::VPB::VPB() : fActionParamValueRef(NULL)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIActionManager);
		result = sAIActionManager->AINewActionParamValue(&fActionParamValueRef);
		SDK_ASSERT(!result);
		SDK_ASSERT(fActionParamValueRef);
	}
	catch (ai::Error) {
	}
}

/*
*/
SnpDocumentActionHelper::VPB::~VPB()
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIActionManager);
		result = sAIActionManager->AIDeleteActionParamValue(fActionParamValueRef);
		SDK_ASSERT(!result);
		fActionParamValueRef = NULL;
	}
	catch (ai::Error) {
	}
}

/*
*/
void SnpDocumentActionHelper::VPB::SetNewDocumentName(const ai::UnicodeString& name)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetStringUS(this->fActionParamValueRef, kAINewDocumentNameKey, name);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetNewDocumentSettingsFile(const ai::FilePath& settingsFile)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetStringUS(this->fActionParamValueRef, kAINewDocumentSettingsFileKey, settingsFile.GetFullPath());
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetNewDocumentColorModel(AIDocumentColorModelValue colorModel)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetInteger(this->fActionParamValueRef, kAINewDocumentColorModelKey, colorModel);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetNewDocumentWidth(AIReal width)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetReal(this->fActionParamValueRef, kAINewDocumentWidthKey, width);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetNewDocumentHeight(AIReal height)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetReal(this->fActionParamValueRef, kAINewDocumentHeightKey, height);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetNewDocumentRulerUnits(AIDocumentRulerUnitValue rulerUnit)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetInteger(this->fActionParamValueRef, kAINewDocumentRulerUnitsKey, rulerUnit);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetNewDocumentRasterResolution(AIRasterResolution rasterResolution)
{
	SDK_ASSERT(sAIActionManager);
	ASUInt32 rasterResolutionPPI = 72;
	switch (rasterResolution)
	{	
		default:
		case kAIRasterResolutionScreen:
		{
			rasterResolutionPPI = 72;
		}
		break;
		case kAIRasterResolutionMedium:
		{
			rasterResolutionPPI = 150;
		}
		break;
		case kAIRasterResolutionHigh:
		{
			rasterResolutionPPI = 300;
		}
		break;
	}
	ASErr result = sAIActionManager->AIActionSetInteger(this->fActionParamValueRef, kAINewDocumentRasterResolutionKey, rasterResolutionPPI); 
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetNewDocumentPreviewMode(AIPreviewMode previewMode)
{
	SDK_ASSERT(sAIActionManager);
	ASUInt32 previewModeMask = 0x020;
	switch (previewMode)
	{	
		default:
		case kAIPreviewModeDefault:
		{
			previewModeMask = 0x020;
		}
		break;
		case kAIPreviewModePixelPreview:
		{
			previewModeMask = 0x040;
		}
		break;
		case kAIPreviewModeOverprintPreview:
		{
			previewModeMask = 0x100;
		}
		break;
	}
	ASErr result = sAIActionManager->AIActionSetInteger(this->fActionParamValueRef, kAINewDocumentPreviewModeKey, previewModeMask);
	aisdk::check_ai_error(result);
}

// kAIOpenDocumentAction parameters
/*
*/
void SnpDocumentActionHelper::VPB::SetOpenDocumentName(const ai::FilePath& name)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetStringUS(this->fActionParamValueRef, kAIOpenDocumentNameKey, name.GetFullPath());
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetOpenDocumentColorModel(AIDocumentColorModelValue colorModel)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetInteger(this->fActionParamValueRef, kAIOpenDocumentColorModelKey, colorModel);
	aisdk::check_ai_error(result);
}

// kAISaveDocumentAsAction parameters
/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsName(const ai::FilePath& name)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetStringUS(this->fActionParamValueRef, kAISaveDocumentAsNameKey, name.GetFullPath());
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsFormat(const char* format)  
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetString(this->fActionParamValueRef, kAISaveDocumentAsFormatKey, format);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsGetParams(ASBoolean getParams)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetBoolean(this->fActionParamValueRef, kAISaveDocumentAsGetParamsKey, getParams);
	aisdk::check_ai_error(result);
}

// kAINativeFileFormat parameters

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsAIVersionKey(AIVersion version)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetInteger(this->fActionParamValueRef, kAINativeVersionKey, version);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsAIPDFCompatibilityKey(ASBoolean compatibility)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetBoolean(this->fActionParamValueRef, kAINativePDFCompatibilityKey, compatibility);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsAIIncludeImagesKey(ASInt32 embedImages)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetInteger(this->fActionParamValueRef, kAINativeIncludeImagesKey, embedImages);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsAIEmbedFontsKey(ASInt32 embedFonts)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetInteger(this->fActionParamValueRef, kAINativeEmbedAllFontsKey, embedFonts);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsAISubsetFontsRatioKey(ASInt32 subsetFontsRatio)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetInteger(this->fActionParamValueRef, kAINativeSubsetFontsRatioKey, subsetFontsRatio);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsAISubsetFontsBelowKey(ASBoolean subsetFonts)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetBoolean(this->fActionParamValueRef, kAINativeSubsetFontsBelowKey, subsetFonts);
	aisdk::check_ai_error(result);
}

// kAIEPSFileFormat parameters

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsEPSVersionKey(ASInt32 version)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetInteger(this->fActionParamValueRef, kAIEPSVersionKey, version);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsEPSIncludeThumbnailKey(ASBoolean includeThumb)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetBoolean(this->fActionParamValueRef, kAIEPSIncludeThumbnailKey, includeThumb);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsEPSCompatibleGradientPrintKey(ASBoolean gradientPrint)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetBoolean(this->fActionParamValueRef, kAIEPSCompatibleGradientPrintKey, gradientPrint);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsEPSCMYKPostScriptKey(ASBoolean cmykPostScript)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetBoolean(this->fActionParamValueRef, kAIEPSCMYKPostScriptKey, cmykPostScript);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsEPSIncludePlacedKey(ASBoolean includePlaced)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetBoolean(this->fActionParamValueRef, kAIEPSIncludePlacedKey, includePlaced);
	aisdk::check_ai_error(result);
}

// kAIPDFFileFormat parameters

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsPDFCompatibilityKey(AIPDFCompatibility compatibility)
{
	SDK_ASSERT(sAIActionManager);
	ASUInt32 pdfCompatibilityMask = 0x0017;
	switch (compatibility)
	{	
		default:
		case kAIPDFCompatibility13:
		{
			pdfCompatibilityMask = 0x0013;
		}
		break;
		case kAIPDFCompatibility14:
		{
			pdfCompatibilityMask = 0x0014;
		}
		break;
		case kAIPDFCompatibility15:
		{
			pdfCompatibilityMask = 0x0015;
		}
		break;
		case kAIPDFCompatibility16:
		{
			pdfCompatibilityMask = 0x0016;
		}
		break;
		case kAIPDFCompatibility17:
		{
			pdfCompatibilityMask = 0x0017;
		}
		break;
	}
	ASErr result = sAIActionManager->AIActionSetInteger(this->fActionParamValueRef, kAIPDFCompatibilityKey, pdfCompatibilityMask);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsPDFRoundTripKey(ASBoolean roundTrip)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetBoolean(this->fActionParamValueRef, kAIPDFRoundTripKey, roundTrip);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsPDFGenerateThumbnailsKey(ASBoolean thumbnails)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetBoolean(this->fActionParamValueRef, kAIPDFGenerateThumbnailsKey, thumbnails);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsPDFViewPDFFileKey(ASBoolean viewPDFFile)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetBoolean(this->fActionParamValueRef, kAIPDFViewPDFFileKey, viewPDFFile);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsPDFFastWebViewKey(ASBoolean fastWebView)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetBoolean(this->fActionParamValueRef, kAIPDFFastWebViewKey, fastWebView);
	aisdk::check_ai_error(result);
}

// kAISVGFileFormat parameters

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsSVGFontSubsettingKey(AISVGFontSubsetting subsetting)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetInteger(this->fActionParamValueRef, kAISVGFontSubsettingKey, subsetting);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsSVGEmbedRasterLocationKey(AISVGRasterLocation embedRaster)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetInteger(this->fActionParamValueRef, kAISVGEmbedRasterLocationKey, embedRaster);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsSVGDTDKey(AISVGDTD dtd)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetInteger(this->fActionParamValueRef, kAISVGDTDKey, dtd);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetSaveDocumentAsSVGFontFormatKey(AISVGFontType fontFormat)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetInteger(this->fActionParamValueRef, kAISVGFontFormatKey, fontFormat);
	aisdk::check_ai_error(result);
}

// kAICloseDocumentAction parameters
/*
*/
void SnpDocumentActionHelper::VPB::SetCloseAndSaveDocument(ASBoolean saveDocument)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetBoolean(this->fActionParamValueRef, kAICloseAndSaveDocumentKey, saveDocument);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetCloseAllDocuments(ASBoolean closeAll)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetBoolean(this->fActionParamValueRef, kAICloseAllDocumentsKey, closeAll);
	aisdk::check_ai_error(result);
}


// kAIPlaceDocumentAction parameters
/*
*/
void SnpDocumentActionHelper::VPB::SetPlaceDocumentActionName(const ai::FilePath& filePath)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetStringUS(this->fActionParamValueRef, kAIPlaceDocumentActionNameKey, filePath.GetFullPath());
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetPlaceDocumentActionLink(ASBoolean link)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetBoolean(this->fActionParamValueRef, kAIPlaceDocumentActionLinkKey, link);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetPlaceDocumentActionReplace(ASBoolean replace)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetBoolean(this->fActionParamValueRef, kAIPlaceDocumentActionReplaceKey, replace);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetPlaceDocumentActionTemplate(ASBoolean createTemplateLayer)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetBoolean(this->fActionParamValueRef, kAIPlaceDocumentActionTemplateKey, createTemplateLayer);
	aisdk::check_ai_error(result);
}

// kAIExportDocumentAction parameters
/*
*/
void SnpDocumentActionHelper::VPB::SetExportDocumentName(const ai::FilePath& name)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetStringUS(this->fActionParamValueRef, kAIExportDocumentNameKey, name.GetFullPath());
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetExportDocumentExtensionKey(const char* extn)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetString(this->fActionParamValueRef, kAIExportDocumentExtensionKey, extn);
	aisdk::check_ai_error(result);
}

/*
*/
void SnpDocumentActionHelper::VPB::SetExportDocumentFormat(const char* format)
{
	SDK_ASSERT(sAIActionManager);
	ASErr result = sAIActionManager->AIActionSetString(this->fActionParamValueRef, kAIExportDocumentFormatKey, format);
	aisdk::check_ai_error(result);
}

/*
*/
ASErr SnpDocumentActionHelper::NewDocument()
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIActionManager);
		SnpDocumentActionHelper::VPB vpb;
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();

		ai::UnicodeString documentName = parameter->GetString(ai::UnicodeString("Document name"), ai::UnicodeString("MyDocument"));
		vpb.SetNewDocumentName(documentName);

		ai::FilePath profilesFolder;
		result = sAIFolders->FindFolder(kAIUserWritableStartupFileFolderType, false, profilesFolder);
		aisdk::check_ai_error(result);
		ai::FilePath defaultPresetFilePath = profilesFolder;
		defaultPresetFilePath.AddComponent(ai::FilePath(ai::UnicodeString("Print.ai")));
		ai::FilePath presetFilePath = parameter->GetFilePath("New Document Profile", true, defaultPresetFilePath);
		vpb.SetNewDocumentSettingsFile(presetFilePath);

		SnpChooser chooser;
		AIDocumentColorModelValue documentColorModel = chooser.PickDocumentColorModel();
		vpb.SetNewDocumentColorModel(documentColorModel);

		AIReal documentWidth = parameter->GetReal("Document width", 612.0);
		vpb.SetNewDocumentWidth(documentWidth);

		AIReal documentHeight = parameter->GetReal("Document height", 792.0);
		vpb.SetNewDocumentHeight(documentHeight);

		AIDocumentRulerUnitValue documentRulerUnit = chooser.PickDocumentRulerUnit();
		vpb.SetNewDocumentRulerUnits(documentRulerUnit);

		AIRasterResolution rasterResolution = chooser.PickDocumentRasterResolution();
		vpb.SetNewDocumentRasterResolution(rasterResolution);

		AIPreviewMode previewMode = chooser.PickPreviewMode();
		vpb.SetNewDocumentPreviewMode(previewMode);

		result = this->NewDocument(kDialogOff, vpb);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDocumentActionHelper::NewDocument(ActionDialogStatus dialogStatus, AIActionParamValueRef parameters)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIActionManager);
		result = sAIActionManager->PlayActionEvent(kAINewDocumentAction, dialogStatus, parameters);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDocumentActionHelper::OpenDocument(const ai::UnicodeString& name)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIActionManager);

		// Stores action parameters.
		SnpDocumentActionHelper::VPB vpb;

		// Name parameter
		SnpDocumentHelper docHelper;
		ai::FilePath path = docHelper.GetAssetPath(name);
		SnippetRunnerLog::Instance()->Write(path.GetFullPath());
		vpb.SetOpenDocumentName(path);

		// Color model parameter.
		SnpChooser chooser;
		AIDocumentColorModelValue colorModel = chooser.PickDocumentColorModel();
		vpb.SetOpenDocumentColorModel(colorModel);

		// Play the action.
		result = this->OpenDocument(kDialogOff, vpb);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDocumentActionHelper::OpenDocument(ActionDialogStatus dialogStatus, AIActionParamValueRef parameters)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIActionManager);
		result = sAIActionManager->PlayActionEvent(kAIOpenDocumentAction, dialogStatus, parameters);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;

}

/*
*/
ASErr SnpDocumentActionHelper::SaveDocument()
{
	return this->SaveDocument(kDialogOff, NULL);
}

/*
*/
ASErr SnpDocumentActionHelper::SaveDocument(ActionDialogStatus dialogStatus, AIActionParamValueRef parameters)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIActionManager);
		result = sAIActionManager->PlayActionEvent(kAISaveDocumentAction, dialogStatus, parameters);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;

}

/*
*/
ASErr SnpDocumentActionHelper::SaveDocumentAsAI(const ai::UnicodeString& name)
{
	ASErr result = kNoErr;
	try {
		SnpDocumentActionHelper::VPB vpb;

		// Format.
		vpb.SetSaveDocumentAsFormat(kAINativeFileFormat);		

		// Version number.
		AIVersion version = kAIVersion14;
		std::vector<std::string> versionChoices;
		versionChoices.push_back("kAIVersion1");
		versionChoices.push_back("kAIVersion88");
		versionChoices.push_back("kAIVersion3");
		versionChoices.push_back("kAIVersion4");
		versionChoices.push_back("kAIVersion5");
		versionChoices.push_back("kAIVersion6");
		versionChoices.push_back("kAIVersion7");
		versionChoices.push_back("kAIVersion8");
		versionChoices.push_back("kAIVersion9");
		versionChoices.push_back("kAIVersion10");
		versionChoices.push_back("kAIVersion11");
		versionChoices.push_back("kAIVersion12");
		versionChoices.push_back("kAIVersion13");
		versionChoices.push_back("kAIVersion14");
		version = static_cast<AIVersion>(SnippetRunnerParameter::Instance()->GetChoice("Version", versionChoices, 13));
		vpb.SetSaveDocumentAsAIVersionKey(version);

		// Pdf compatibility.
		ASBoolean compatibility = SnippetRunnerParameter::Instance()->GetBoolean("Generate PDF compatible document");
		vpb.SetSaveDocumentAsAIPDFCompatibilityKey(compatibility);

		// Include images.
		ASInt32 embedImages = (ASInt32)SnippetRunnerParameter::Instance()->GetBoolean("Embed images?");
		vpb.SetSaveDocumentAsAIIncludeImagesKey(embedImages);
		
		// Embed fonts.
		ASInt32 embedFonts = (ASInt32)SnippetRunnerParameter::Instance()->GetBoolean("Embed fonts?");
		vpb.SetSaveDocumentAsAIEmbedFontsKey(embedFonts);

		// Subset fonts ratio.
		ASInt32 subsetFontsRatio = SnippetRunnerParameter::Instance()->GetInt("Subset fonts ratio", 100, 0, 100);
		vpb.SetSaveDocumentAsAISubsetFontsRatioKey(subsetFontsRatio);
		
		// Subset fonts below.
		ASBoolean subsetFonts = SnippetRunnerParameter::Instance()->GetBoolean("Subset fonts below ratio");
		vpb.SetSaveDocumentAsAISubsetFontsBelowKey(subsetFonts);

		// Gather common parameters then save.
		result = this->SaveDocumentAs(name, vpb);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDocumentActionHelper::SaveDocumentAsEPS(const ai::UnicodeString& name)
{
	ASErr result = kNoErr;
	try {
		SnpDocumentActionHelper::VPB vpb;		

		// Format parameter.
		vpb.SetSaveDocumentAsFormat(kAIEPSFileFormat);		

		// Include thumbnail.
		ASBoolean includeThumb = SnippetRunnerParameter::Instance()->GetBoolean("Include thumbnail");
		vpb.SetSaveDocumentAsEPSIncludeThumbnailKey(includeThumb);

		// Compatible gradient print.
		ASBoolean compatible = SnippetRunnerParameter::Instance()->GetBoolean("Make compatible with gradient print and mesh");
		vpb.SetSaveDocumentAsEPSCompatibleGradientPrintKey(compatible);

		// CMYK postscript.
		ASBoolean cmykPostScript = SnippetRunnerParameter::Instance()->GetBoolean("Convert RGB to CMYK for PostScript?");
		vpb.SetSaveDocumentAsEPSCMYKPostScriptKey(cmykPostScript);

		// Include placed items.
		ASBoolean includePlaced = SnippetRunnerParameter::Instance()->GetBoolean(ai::UnicodeString("Include placed items?"));
		vpb.SetSaveDocumentAsEPSIncludePlacedKey(includePlaced);

		// Gather common parameters then save.
		result = this->SaveDocumentAs(name, vpb);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDocumentActionHelper::SaveDocumentAsPDF(const ai::UnicodeString& name)
{
	ASErr result = kNoErr;
	try {
		SnpDocumentActionHelper::VPB vpb;

		// Format parameter.
		vpb.SetSaveDocumentAsFormat(kAIPDFFileFormat);	

		// Compatibility.
		AIPDFCompatibility compatibility = kAIPDFCompatibility17;
		std::vector<std::string> pdfCompatibilities;
		pdfCompatibilities.push_back("kAIPDFCompatibility13 (Acrobat 4)");
		pdfCompatibilities.push_back("kAIPDFCompatibility14 (Acrobat 5)");
		pdfCompatibilities.push_back("kAIPDFCompatibility15 (Acrobat 6)");
		pdfCompatibilities.push_back("kAIPDFCompatibility16 (Acrobat 7)");
		pdfCompatibilities.push_back("kAIPDFCompatibility17 (Acrobat 8)");
		compatibility = static_cast<AIPDFCompatibility>(SnippetRunnerParameter::Instance()->GetChoice("PDF compatibility", pdfCompatibilities, 4));
		vpb.SetSaveDocumentAsPDFCompatibilityKey(compatibility);

		// Round trip.
		ASBoolean roundTrip = SnippetRunnerParameter::Instance()->GetBoolean("Round trip information?");
		vpb.SetSaveDocumentAsPDFRoundTripKey(roundTrip);

		// Generate thumbnails.
		ASBoolean thumbnails = SnippetRunnerParameter::Instance()->GetBoolean("Generate thumbnails?");
		vpb.SetSaveDocumentAsPDFGenerateThumbnailsKey(thumbnails);

		// View PDF file.
		ASBoolean viewPDFFile = SnippetRunnerParameter::Instance()->GetBoolean("View PDF file?");
		vpb.SetSaveDocumentAsPDFViewPDFFileKey(viewPDFFile);

		// Fast web view.
		ASBoolean fastWebView = SnippetRunnerParameter::Instance()->GetBoolean("Fast web view?");
		vpb.SetSaveDocumentAsPDFFastWebViewKey(fastWebView);

		// Gather common parameters then save.
		result = this->SaveDocumentAs(name, vpb);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDocumentActionHelper::SaveDocumentAsSVG(const ai::UnicodeString& name)
{
	ASErr result = kNoErr;
	try {
		SnpDocumentActionHelper::VPB vpb;

		// Format parameter.
		vpb.SetSaveDocumentAsFormat(kAISVGFileFormat);	

		// Font subsetting.
		AISVGFontSubsetting subsetting = kAISVGSubsettingNone;
		std::vector<std::string> subsettings;
		subsettings.push_back("kAISVGSubsettingNone");
		subsettings.push_back("kAISVGSubsettingGlyphsOnly");
		subsettings.push_back("kAISVGSubsettingCommonEnglishOnly");
		subsettings.push_back("kAISVGSubsettingGlyphsPlusEnglish");
		subsettings.push_back("kAISVGSubsettingCommonRomanOnly");
		subsettings.push_back("kAISVGSubsettingGlyphsPlusRoman");
		subsettings.push_back("kAISVGSubsettingOutputEntireFont");
		subsetting = static_cast<AISVGFontSubsetting>(SnippetRunnerParameter::Instance()->GetChoice("Font subsetting", subsettings));
		vpb.SetSaveDocumentAsSVGFontSubsettingKey(subsetting);

		// Embed raster location key.
		AISVGRasterLocation embedRaster = kAISVGImageEmbed;
		std::vector<std::string> embedLocations;
		embedLocations.push_back("kAISVGImageEmbed");
		embedLocations.push_back("kAISVGImageLink");
		embedRaster = static_cast<AISVGRasterLocation>(SnippetRunnerParameter::Instance()->GetChoice("Embed raster location", embedLocations));
		vpb.SetSaveDocumentAsSVGEmbedRasterLocationKey(embedRaster);

		// Document type definition.
		AISVGDTD dtd = kAISVGDTD11;
		std::vector<std::string> dtdChoices;
		dtdChoices.push_back("kAISVGDTD10");
		dtdChoices.push_back("kAISVGDTD11");
		dtdChoices.push_back("kAISVGDTDTiny11");
		dtdChoices.push_back("kAISVGDTDTiny11Plus");
		dtdChoices.push_back("kAISVGDTDBasic11");
		dtdChoices.push_back("kAISVGDTDTiny12");
		dtd = static_cast<AISVGDTD>(SnippetRunnerParameter::Instance()->GetChoice("Document type definition", dtdChoices));
		vpb.SetSaveDocumentAsSVGDTDKey(dtd);

		// Font format.
		AISVGFontType fontFormat = kAISVGFontCEF;
		std::vector<std::string> fontFormats;
		fontFormats.push_back("kAISVGFontCEF");
		fontFormats.push_back("kAISVGFontSVG");
		fontFormats.push_back("kAISVGFontOutline");
		fontFormat = static_cast<AISVGFontType>(SnippetRunnerParameter::Instance()->GetChoice("Font format", fontFormats));
		vpb.SetSaveDocumentAsSVGFontFormatKey(fontFormat);

		// Gather common parameters then save.
		result = this->SaveDocumentAs(name, vpb);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDocumentActionHelper::SaveDocumentAsSVGZ(const ai::UnicodeString& name)
{
	ASErr result = kNoErr;
	try {
		SnpDocumentActionHelper::VPB vpb;

		// Format parameter.
		vpb.SetSaveDocumentAsFormat(kAISVGCompressedFileFormat);		

		// Gather common parameters then save.
		result = this->SaveDocumentAs(name, vpb);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDocumentActionHelper::SaveDocumentAs(const ai::UnicodeString& name, SnpDocumentActionHelper::VPB& vpb)
{
	ASErr result = kNoErr;
	try {
		// Name parameter.
		SnpDocumentHelper docHelper;
		ai::FilePath path = docHelper.GetAssetPath(name);
		SnippetRunnerLog::Instance()->Write(path.GetFullPath());
		vpb.SetSaveDocumentAsName(path);

		// Enable dialogs parameter.
		ASBoolean getParams = SnippetRunnerParameter::Instance()->GetBoolean("Enable dialogs during save");
		vpb.SetSaveDocumentAsGetParams(getParams);

		// Play the action.
		result = this->SaveDocumentAs(kDialogOff, vpb);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;

}

/*
*/
ASErr SnpDocumentActionHelper::SaveDocumentAs(ActionDialogStatus dialogStatus, AIActionParamValueRef parameters)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIActionManager);
		result = sAIActionManager->PlayActionEvent(kAISaveDocumentAsAction, dialogStatus, parameters);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;

}

/*
*/
ASErr SnpDocumentActionHelper::CloseDocument()
{
	ASErr result = kNoErr;
	try {
		SnpDocumentActionHelper::VPB vpb;

		// Save parameter.
		ASBoolean saveDocument = SnippetRunnerParameter::Instance()->GetBoolean("Save document before close");
		vpb.SetCloseAndSaveDocument(saveDocument);

		// Close all documents parameter.
		ASBoolean closeAll = SnippetRunnerParameter::Instance()->GetBoolean("Close all documents");
		vpb.SetCloseAllDocuments(closeAll);

		// Play the action.
		result = this->CloseDocument(kDialogOff, vpb);
		aisdk::check_ai_error(result);

	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDocumentActionHelper::CloseDocument(ActionDialogStatus dialogStatus, AIActionParamValueRef parameters)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIActionManager);
		result = sAIActionManager->PlayActionEvent(kAICloseDocumentAction, dialogStatus, parameters);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;

}

/*
*/
ASErr SnpDocumentActionHelper::PlaceFile(const ai::UnicodeString& name)
{
	ASErr result = kNoErr;
	try {
		SnpDocumentActionHelper::VPB vpb;

		// Name parameter.
		SnpDocumentHelper docHelper;
		ai::FilePath filePath = docHelper.GetAssetPath(name);
		SnippetRunnerLog::Instance()->Write(filePath.GetFullPath());
		vpb.SetPlaceDocumentActionName(filePath);

		// Link parameter.
		ASBoolean link = SnippetRunnerParameter::Instance()->GetBoolean("Linked (embedded if false)");
		vpb.SetPlaceDocumentActionLink(link);

		// Replace parameter.
		ASBoolean replace = SnippetRunnerParameter::Instance()->GetBoolean("Replace current selection");
		vpb.SetPlaceDocumentActionReplace(replace);

		// Template layer parameter.
		ASBoolean createTemplateLayer = SnippetRunnerParameter::Instance()->GetBoolean("Create template layer");
		vpb.SetPlaceDocumentActionTemplate(createTemplateLayer);

		result = this->PlaceFile(kDialogOff, vpb);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDocumentActionHelper::PlaceFile(ActionDialogStatus dialogStatus, AIActionParamValueRef parameters)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIActionManager);
		result = sAIActionManager->PlayActionEvent(kAIPlaceDocumentAction, dialogStatus, parameters);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDocumentActionHelper::ExportSelectionToAI(const ai::UnicodeString& name)
{
	ASErr result = kNoErr;
	try {
		result = this->ExportSelection(name, kAINativeFileFormat, kAINativeFileFormatExtension);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDocumentActionHelper::ExportSelectionToEPS(const ai::UnicodeString& name)
{
	ASErr result = kNoErr;
	try {
		result = this->ExportSelection(name, kAIEPSFileFormat, kAIEPSFileFormatExtension);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDocumentActionHelper::ExportSelection(const ai::UnicodeString& name, const char* format, const char* extn)
{
	ASErr result = kNoErr;
	try {
		// Copy selection to clipboard.
		SnpObjectActionHelper objActionHelper;
		result = objActionHelper.CopySelection();
		aisdk::check_ai_error(result);

		// Open new document.
		result = this->NewDocument();
		aisdk::check_ai_error(result);

		// Workaround for 1348729, see Readme.txt for full details.
		AIDocumentHandle document = NULL;
		result = sAIDocumentList->GetNthDocument(&document, 0);
		aisdk::check_ai_error(result);
		result = sAIDocumentList->Activate(document, true);
		aisdk::check_ai_error(result);

		// Paste clipboard.
		result = objActionHelper.PasteClipboard();
		aisdk::check_ai_error(result);

		// Create new value parameter block.
		SnpDocumentActionHelper::VPB vpb;

		// Name.
		SnpDocumentHelper docHelper;
		ai::FilePath path = docHelper.GetAssetPath(name);
		SnippetRunnerLog::Instance()->Write(path.GetFullPath());
		vpb.SetExportDocumentName(path);

		// Format.
		vpb.SetExportDocumentFormat(format);
		
		// Extension.
		vpb.SetExportDocumentExtensionKey(extn);

		// Enable dialogs parameter.
		ASBoolean showDialog = SnippetRunnerParameter::Instance()->GetBoolean("Display dialog?");
		ActionDialogStatus dialogStatus = kDialogOff;
		if (showDialog)
			dialogStatus = kDialogOn;

		result = this->ExportDocument(dialogStatus, vpb);
		aisdk::check_ai_error(result);

		// Use AIDocumentListSuite to close the document due to application
		// context update issue using action, see Readme.txt for full details.
		result = sAIDocumentList->Close(document);
		aisdk::check_ai_error(result);

	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDocumentActionHelper::ExportDocumentToDWG(const ai::UnicodeString& name)
{
	ASErr result = kNoErr;
	try {
		result = this->ExportDocument(name, kDwgExportFileFormatName, kFormatDwgExtn);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDocumentActionHelper::ExportDocumentToDXF(const ai::UnicodeString& name)
{
	ASErr result = kNoErr;
	try {
		result = this->ExportDocument(name, kDxfExportFileFormatName, kFormatDxfExtn);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDocumentActionHelper::ExportDocumentToPSD(const ai::UnicodeString& name)
{
	ASErr result = kNoErr;
	try {
		result = this->ExportDocument(name, kPSDExportFormatName, kPSDFormatExtn);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDocumentActionHelper::ExportDocument(const ai::UnicodeString& name, const char* format, const char* extn)
{
	ASErr result = kNoErr;
	try {
		SnpDocumentActionHelper::VPB vpb;

		// Name.
		SnpDocumentHelper docHelper;
		ai::FilePath path = docHelper.GetAssetPath(name);
		SnippetRunnerLog::Instance()->Write(path.GetFullPath());
		vpb.SetExportDocumentName(path);

		// Format.
		vpb.SetExportDocumentFormat(format);
		
		// Extension.
		vpb.SetExportDocumentExtensionKey(extn);

		// Enable dialogs parameter.
		ASBoolean showDialog = SnippetRunnerParameter::Instance()->GetBoolean("Display dialog?");
		ActionDialogStatus dialogStatus = kDialogOff;
		if (showDialog)
			dialogStatus = kDialogOn;

		// Call export action.
		result = this->ExportDocument(dialogStatus, vpb);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}


/*
*/
ASErr SnpDocumentActionHelper::ExportDocument(ActionDialogStatus dialogStatus, AIActionParamValueRef parameters)
{
	ASErr result = kNoErr;
	try {
		SDK_ASSERT(sAIActionManager);
		result = sAIActionManager->PlayActionEvent(kAIExportDocumentAction, dialogStatus, parameters);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpDocumentActionHelper available to the SnippetRunner framework.
 */
class _SnpRunnableDocumentActionHelper : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableDocumentActionHelper () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableDocumentActionHelper () {}

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

		/* Runs the snippet.
			@param runnableContext see ISnpRunnableContext for documentation.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr		Run(SnpRunnable::Context& runnableContext);

};

/*
*/
std::string _SnpRunnableDocumentActionHelper::GetName() const
{
	return "DocumentActionHelper";
}

/*
*/
std::string _SnpRunnableDocumentActionHelper::GetDescription() const 
{
	return "Helper class for programming document actions";
}

/*
*/
SnpRunnable::Operations _SnpRunnableDocumentActionHelper::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("NewDocument", "", kSnpRunNotSupportedContext)); // skip it - see 1348729
	operations.push_back(Operation("OpenDocument", "", kSnpRunNotSupportedContext));  // skip it - see 1348729
	operations.push_back(Operation("SaveDocument", "document", kSnpRunNotSupportedContext)); // skip it - interactive
	operations.push_back(Operation("SaveDocumentAsAI", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("SaveDocumentAsEPS", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("SaveDocumentAsPDF", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("SaveDocumentAsSVG", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("SaveDocumentAsSVGZ", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("CloseDocument", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("PlaceFile", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("ExportSelectionToAI", "document", kSnpRunPathContext));
	operations.push_back(Operation("ExportSelectionToEPS", "document", kSnpRunPathContext));
	operations.push_back(Operation("ExportDocumentToDWG", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("ExportDocumentToDXF", "document", kSnpRunUnitTestDocumentContext));
	operations.push_back(Operation("ExportDocumentToPSD", "document containing art", kSnpRunUnitTestDocumentContext));
	return operations;
}

/*
*/
std::string _SnpRunnableDocumentActionHelper::GetDefaultOperationName() const 
{
	return "NewDocument";
}

/*
*/
std::vector<std::string> _SnpRunnableDocumentActionHelper::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Helper Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableDocumentActionHelper::CanRun(SnpRunnable::Context& runnableContext)
{
	if ("NewDocument" == runnableContext.GetOperation().GetName() ||
		"OpenDocument" == runnableContext.GetOperation().GetName()) {
				
			return true;
	}
		
	SnpSelectionHelper selectionHelper;
	if (!selectionHelper.IsDocumentSelected())
		return false;


	if ("ExportSelectionToAI" == runnableContext.GetOperation().GetName() ||
		"ExportSelectionToEPS" == runnableContext.GetOperation().GetName()) {
		
			return selectionHelper.IsArtSelected();
	}
	else if ("ExportDocumentToPSD" == runnableContext.GetOperation().GetName()) {

		ai::int32 attr = 0;
		const short numSpecs = 4;
		AIArtSpec artSpec[numSpecs] = { 
			{kPathArt, attr, attr},
			{kTextFrameArt, attr, attr},
			{kPlacedArt, attr, attr},
			{kRasterArt, attr, attr} 
		};
		SnpArtSetHelper artSet(artSpec, numSpecs);
		return artSet.GetCount() > 0;
	}
	else
		return true;
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableDocumentActionHelper::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	try {
		SnpDocumentActionHelper instance;
		if ("NewDocument" == runnableContext.GetOperation().GetName()) {
			result = instance.NewDocument();
		}
		else if ("OpenDocument" == runnableContext.GetOperation().GetName()) {
			result = instance.OpenDocument(ai::UnicodeString("sample-1.ai"));
		}
		else if ("SaveDocument" == runnableContext.GetOperation().GetName()) {
			result = instance.SaveDocument();
		}
		else if ("SaveDocumentAsAI" == runnableContext.GetOperation().GetName()) {
			result = instance.SaveDocumentAsAI(ai::UnicodeString("SnpDocumentActionHelper.ai"));
		}
		else if ("SaveDocumentAsEPS" == runnableContext.GetOperation().GetName()) {
			result = instance.SaveDocumentAsEPS(ai::UnicodeString("SnpDocumentActionHelper.eps"));
		}
		else if ("SaveDocumentAsPDF" == runnableContext.GetOperation().GetName()) {
			result = instance.SaveDocumentAsPDF(ai::UnicodeString("SnpDocumentActionHelper.pdf"));
		}
		else if ("SaveDocumentAsSVG" == runnableContext.GetOperation().GetName()) {
			result = instance.SaveDocumentAsSVG(ai::UnicodeString("SnpDocumentActionHelper.svg"));
		}
		else if ("SaveDocumentAsSVGZ" == runnableContext.GetOperation().GetName()) {
			result = instance.SaveDocumentAsSVGZ(ai::UnicodeString("SnpDocumentActionHelper.svgz"));
		}
		else if ("CloseDocument" == runnableContext.GetOperation().GetName()) {
			result = instance.CloseDocument();
		}
		else if ("PlaceFile" == runnableContext.GetOperation().GetName()) {
			result = instance.PlaceFile(ai::UnicodeString("pavilion.jpg"));
		}
		else if ("ExportSelectionToAI" == runnableContext.GetOperation().GetName()) {
			result = instance.ExportSelectionToAI(ai::UnicodeString("SnpDocumentActionHelperExportSelection.ai"));
		}
		else if ("ExportSelectionToEPS" == runnableContext.GetOperation().GetName()) {
			result = instance.ExportSelectionToEPS(ai::UnicodeString("SnpDocumentActionHelperExportSelection.eps"));
		}
		else if ("ExportDocumentToDWG" == runnableContext.GetOperation().GetName()) {
			result = instance.ExportDocumentToDWG(ai::UnicodeString("SnpDocumentActionHelperExport.dwg"));
		}
		else if ("ExportDocumentToDXF" == runnableContext.GetOperation().GetName()) {
			result = instance.ExportDocumentToDXF(ai::UnicodeString("SnpDocumentActionHelperExport.dxf"));
		}
		else if ("ExportDocumentToPSD" == runnableContext.GetOperation().GetName()) {
			result = instance.ExportDocumentToPSD(ai::UnicodeString("SnpDocumentActionHelperExport.psd"));
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

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableDocumentActionHelper instance_SnpRunnableDocumentActionHelper;

// End SnpDocumentActionHelper.cpp

