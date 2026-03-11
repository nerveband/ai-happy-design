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

#ifndef __SnpDocumentActionHelper_h__
#define __SnpDocumentActionHelper_h__

#include "AIPDFFormatAction.h"
#include "AISVGAction.h"

/** Helper class that shows how to program with actions.

	@see AIActionManagerSuite
  	@ingroup sdk_snippet
 */
class SnpDocumentActionHelper
{
public:
		/** Value parameter block (VPB) - a container for action parameters.
			@see AIActionManagerSuite
		*/
		class VPB
		{
		public:
			/**	Constructor
			*/
			VPB();

			/**	Destructor
			*/
			~VPB();

			/**	Returns the value parameter block structure.
			*/
			inline operator AIActionParamValueRef(void) const {return fActionParamValueRef;}

			// kAINewDocumentAction parameters

			/**	Sets the new document name - kAINewDocumentNameKey, in the value parameter
				block.
				@param name IN new document name.
			*/
			void SetNewDocumentName(const ai::UnicodeString& name);

			/**	Sets the document settings file key - kAINewDocumentSettingsFileKey, 
				in the value parameter block.
				@param settingsFile IN path to document settings file.
			*/
			void SetNewDocumentSettingsFile(const ai::FilePath& settingsFile);

			/**	Sets the document color model key - kAINewDocumentColorModelKey, 
				in the value parameter block.
				@param colorModel IN document color model.
			*/
			void SetNewDocumentColorModel(AIDocumentColorModelValue colorModel);

			/**	Sets the width key of the new document - kAINewDocumentWidthKey, 
				in the value parameter block.
				@param width IN new document width.
			*/
			void SetNewDocumentWidth(AIReal width);

			/**	Sets the height key of the new document - kAINewDocumentHeightKey, 
				in the value parameter block.
				@param height IN new document height.
			*/
			void SetNewDocumentHeight(AIReal height);

			/**	Sets the ruler units key of the new document - kAINewDocumentRulerUnitsKey, 
				in the value parameter block.
				@param rulerUnit IN new document ruler unit value.
			*/
			void SetNewDocumentRulerUnits(AIDocumentRulerUnitValue rulerUnit);

			/**	Sets the raster resolution key of the new document - kAINewDocumentRasterResolutionKey, 
				in the value parameter block.
				@param rasterResolution IN new document raster resolution value.
			*/
			void SetNewDocumentRasterResolution(AIRasterResolution rasterResolution);

			/**	Sets the document preview mode key of the new document - kAINewDocumentPreviewModeKey, 
				in the value parameter block.
				@param previewMode IN new document preview mode.
			*/
			void SetNewDocumentPreviewMode(AIPreviewMode previewMode);

			// kAIOpenDocumentAction parameters

			/**	Sets the open document name key - kAIOpenDocumentNameKey, 
				in the value parameter block.
				@param name IN name of document to open.
			*/
			void SetOpenDocumentName(const ai::FilePath& name);

			/**	Sets the open document color model key - kAIOpenDocumentColorModelKey, 
				in the value parameter block.
				@param colorModel IN color model of document to open.
			*/
			void SetOpenDocumentColorModel(AIDocumentColorModelValue colorModel);

			// kAISaveDocumentAsAction parameters

			/**	Sets the name key of the document to save as - kAISaveDocumentAsNameKey, 
				in the value parameter block.
				@param name IN name to save document as.
			*/
			void SetSaveDocumentAsName(const ai::FilePath& name);

			/**	Sets the format key of the document to save as - kAISaveDocumentAsFormatKey, 
				in the value parameter block.
				@param format IN format to save document as.
			*/
			void SetSaveDocumentAsFormat(const char* format);  

			/**	Sets the get params key of the document to save as - kAISaveDocumentAsGetParamsKey, 
				in the value parameter block.
				@param getParams IN true to get params, false if not.
			*/
			void SetSaveDocumentAsGetParams(ASBoolean getParams);

			// kAINativeFileFormat parameters

			/**	Sets the native version key of the document to save - kAINativeVersionKey, 
				in the value parameter block, default is 13.
				@param version IN the AIVersion to save as.
				@see AINativeAction.h
			*/
			void SetSaveDocumentAsAIVersionKey(AIVersion version);

			/**	Sets the pdf compatibility key of the document to save - kAINativePDFCompatibilityKey, 
				in the value parameter block.
				@param compatibility IN true creates a PDF compatible file.
				@see AINativeAction.h
			*/
			void SetSaveDocumentAsAIPDFCompatibilityKey(ASBoolean compatibility);

			/**	Sets the include images key of the document to save - kAINativeIncludeImagesKey, 
				in the value parameter block. If 0 does not embed linked images, if not 0 embeds 
				linked images.
				@param embedImages IN the include images key.
				@see AINativeAction.h
			*/
			void SetSaveDocumentAsAIIncludeImagesKey(ASInt32 embedImages);

			/**	Sets the embed fonts key of the document to save - kAINativeEmbedAllFontsKey, 
				in the value parameter block. If 0 does not embed fonts, if not 0 embeds fonts.
				@param embedFonts IN the embed fonts key.
				@see AINativeAction.h
			*/
			void SetSaveDocumentAsAIEmbedFontsKey(ASInt32 embedFonts);

			/**	Sets the subset fonts ratio of the document to save - kAINativeSubsetFontsRatioKey, 
				in the value parameter block. Threshold value used when subsetting fonts.
				@param subsetFontsRatio IN the subset fonts ratio key.
				@see SnpDocumentActionHelper::VPB::SetSaveDocumentAsAISubsetFontsBelowKey
				@see AINativeAction.h
			*/
			void SetSaveDocumentAsAISubsetFontsRatioKey(ASInt32 subsetFontsRatio);

			/**	Sets the subset fonts below key of the document to save - kAINativeSubsetFontsBelowKey, 
				in the value parameter block. 
				@param subsetFonts IN True to subset the fonts below the threshold, False does not
				subset the fonts.
				@see SnpDocumentActionHelper::VPB::SetSaveDocumentAsAISubsetFontsRatioKey
				@see AINativeAction.h
			*/
			void SetSaveDocumentAsAISubsetFontsBelowKey(ASBoolean subsetFonts);

			// kAIEPSFileFormat parameters

			/**	Sets the version key of the document to save - kAIEPSVersionKey, 
				in the value parameter block.
				@param version IN the AIVersion to save as.
				@see AINativeAction.h
			*/
			void SetSaveDocumentAsEPSVersionKey(ASInt32 version);

			/**	Sets the include thumbnail key of the document to save - kAIEPSIncludeThumbnailKey, 
				in the value parameter block.
				@param includeThumb IN the include thumbnail key.
				@see AINativeAction.h
			*/
			void SetSaveDocumentAsEPSIncludeThumbnailKey(ASBoolean includeThumb);

			/**	Sets the compatible gradient print key of the document to save - kAIEPSCompatibleGradientPrintKey, 
				in the value parameter block.
				@param gradientPrint IN true to make compatible with gradient print and mesh.
				@see AINativeAction.h
			*/
			void SetSaveDocumentAsEPSCompatibleGradientPrintKey(ASBoolean gradientPrint);

			/**	Sets the cmyk postscript key of the document to save - kAIEPSCMYKPostScriptKey, 
				in the value parameter block.
				@param cmykPostScript IN true to convert RGB to CMYK.
				@see AINativeAction.h
			*/
			void SetSaveDocumentAsEPSCMYKPostScriptKey(ASBoolean cmykPostScript);

			/**	Sets the include placed items key of the document to save - kAIEPSIncludePlacedKey, 
				in the value parameter block. True embeds linked images, False does not embed.
				@param includePlaced IN the include placed items key.
				@see AINativeAction.h
			*/
			void SetSaveDocumentAsEPSIncludePlacedKey(ASBoolean includePlaced);

			// kAIPDFFileFormat parameters

			/**	Sets the pdf compatibility key of the document to save - kAIPDFCompatibilityKey, 
				in the value parameter block.
				@param compatibility IN the pdf compatibility key.
				@see AIPDFFormatAction.h
			*/
			void SetSaveDocumentAsPDFCompatibilityKey(AIPDFCompatibility compatibility);

			/**	Sets the round trip key of the document to save - kAIPDFRoundTripKey, 
				in the value parameter block.
				@param roundTrip IN the round trip key.
				@see AIPDFFormatAction.h
			*/
			void SetSaveDocumentAsPDFRoundTripKey(ASBoolean roundTrip);

			/**	Sets the generate thumbnails key of the document to save - kAIPDFGenerateThumbnailsKey, 
				in the value parameter block.
				@param thumbnails IN the generate thumbnails key.
				@see AIPDFFormatAction.h
			*/
			void SetSaveDocumentAsPDFGenerateThumbnailsKey(ASBoolean thumbnails);

			/**	Sets the view pdf file key of the document to save - kAIPDFViewPDFFileKey, 
				in the value parameter block.
				@param viewPDFFile IN the view pdf file key.
				@see AIPDFFormatAction.h
			*/
			void SetSaveDocumentAsPDFViewPDFFileKey(ASBoolean viewPDFFile);

			/**	Sets the fast web view key of the document to save - kAIPDFFastWebViewKey, 
				in the value parameter block.
				@param fastWebView IN the fast web view key.
				@see AIPDFFormatAction.h
			*/
			void SetSaveDocumentAsPDFFastWebViewKey(ASBoolean fastWebView);

			// kAISVGFileFormat parameters

			/**	Sets the font subsetting key of the document to save - kAISVGFontSubsettingKey, 
				in the value parameter block.
				@param subsetting IN the font subsetting key.
				@see AISVGAction.h
			*/
			void SetSaveDocumentAsSVGFontSubsettingKey(AISVGFontSubsetting subsetting);

			/**	Sets the embed raster location key of the document to save - kAISVGEmbedRasterLocationKey, 
				in the value parameter block.
				@param embedRaster IN the embed raster location key.
				@see AISVGAction.h
			*/
			void SetSaveDocumentAsSVGEmbedRasterLocationKey(AISVGRasterLocation embedRaster);

			/**	Sets the SVG XML document type definition key of the document to save - kAISVGDTDKey, 
				in the value parameter block.
				@param dtd IN the document type definition key.
				@see AISVGAction.h
			*/
			void SetSaveDocumentAsSVGDTDKey(AISVGDTD dtd);

			/**	Sets the font format key of the document to save - kAISVGFontFormatKey, 
				in the value parameter block.
				@param fontFormat IN the font format key.
				@see AISVGAction.h
			*/
			void SetSaveDocumentAsSVGFontFormatKey(AISVGFontType fontFormat);
			
			// kAICloseDocumentAction parameters

			/**	Sets the close and save key - kAICloseAndSaveDocumentKey, 
				in the value parameter block.
				@param saveDocument IN true to get save closing document, false if not.
			*/
			void SetCloseAndSaveDocument(ASBoolean saveDocument);

			/**	Sets the close all documents key - kAICloseAllDocumentsKey, 
				in the value parameter block.
				@param closeAll IN true to close all open documents, false if not.
			*/
			void SetCloseAllDocuments(ASBoolean closeAll);

			// kAIPlaceDocumentAction parameters

			/**	Sets the name key of the document to place - kAIPlaceDocumentActionNameKey, 
				in the value parameter block.
				@param filePath IN path of file to place.
			*/
			void SetPlaceDocumentActionName(const ai::FilePath& filePath);

			/**	Sets the link key of the document to place - kAIPlaceDocumentActionLinkKey, 
				in the value parameter block.
				@param link IN if true the file is linked, false the file is embedded.
			*/
			void SetPlaceDocumentActionLink(ASBoolean link); // embedded if false

			/**	Sets the replace key of the document to place - kAIPlaceDocumentActionReplaceKey, 
				in the value parameter block.
				@param replace IN if true the file replaces the current selection.
			*/
			void SetPlaceDocumentActionReplace(ASBoolean replace);

			/**	Sets the template layer key of the document to place - kAIPlaceDocumentActionTemplateKey, 
				in the value parameter block.
				@param createTemplateLayer IN if true places the file on a template layer.
			*/
			void SetPlaceDocumentActionTemplate(ASBoolean createTemplateLayer);

			// kAIExportDocumentAction parameters

			/**	Sets the name key of the document to export - kAIExportDocumentNameKey, 
				in the value parameter block.
				@param name IN name to save document as.
			*/
			void SetExportDocumentName(const ai::FilePath& name);

			/**	Sets the extension key of the document to export - kAIExportDocumentExtensionKey, 
				in the value parameter block.
				@param name IN name to save document as.
			*/
			void SetExportDocumentExtensionKey(const char* extn);

			/**	Sets the format key of the document to export to - kAIExportDocumentFormatKey, 
				in the value parameter block.
				@param format IN format to save document as.
			*/
			void SetExportDocumentFormat(const char* format);  

		private:
			/** The value parameter block.
			*/
			AIActionParamValueRef fActionParamValueRef;
		}; // end class VPB

	public:

		/** Gathers action parameters and creates a new document.
			@return kNoErr on success, other error code otherwise.
			@see kAINewDocumentAction
		*/
		ASErr NewDocument();

		/** Creates a new document using the given parameters.
			@param dialogStatus IN whether parameters dialog appears by default.
			@param parameters IN the value parameter block.
			@return kNoErr on success, other error code otherwise.
			@see kAINewDocumentAction
			@see SnpDocumentActionHelper::VPB
		*/
		ASErr NewDocument(ActionDialogStatus dialogStatus, AIActionParamValueRef parameters);

		/** Gathers action parameters and opens a document.
			@return kNoErr on success, other error code otherwise.
			@see kAIOpenDocumentAction
		*/
		ASErr OpenDocument(const ai::UnicodeString& name);

		/** Opens a new document using the given parameters.
			@param dialogStatus IN whether parameters dialog appears by default.
			@param parameters IN the value parameter block.
			@return kNoErr on success, other error code otherwise.
			@see kAIOpenDocumentAction
			@see SnpDocumentActionHelper::VPB
		*/
		ASErr OpenDocument(ActionDialogStatus dialogStatus, AIActionParamValueRef parameters);

		/** Gathers action parameters and saves a document.
			@return kNoErr on success, other error code otherwise.
			@see kAIOpenDocumentAction
		*/
		ASErr SaveDocument();

		/** Saves a document using the given parameters.
			@param dialogStatus IN whether parameters dialog appears by default.
			@param parameters IN the value parameter block.
			@return kNoErr on success, other error code otherwise.
			@see kAISaveDocumentAction
			@see SnpDocumentActionHelper::VPB
		*/
		ASErr SaveDocument(ActionDialogStatus dialogStatus, AIActionParamValueRef parameters);

		/** Gathers action parameters and saves document as native AI.
			@param name IN name to save the document as.
			@return kNoErr on success, other error code otherwise.
			@see kAISaveDocumentAsAction
		*/
		ASErr SaveDocumentAsAI(const ai::UnicodeString& name);

		/** Gathers action parameters and saves document as EPS.
			@param name IN name to save the document as.
			@return kNoErr on success, other error code otherwise.
			@see kAISaveDocumentAsAction
		*/
		ASErr SaveDocumentAsEPS(const ai::UnicodeString& name);

		/** Gathers action parameters and saves document as PDF.
			@param name IN name to save the document as.
			@return kNoErr on success, other error code otherwise.
			@see kAISaveDocumentAsAction
		*/
		ASErr SaveDocumentAsPDF(const ai::UnicodeString& name);

		/** Gathers action parameters and saves document as SVG.
			@param name IN name to save the document as.
			@return kNoErr on success, other error code otherwise.
			@see kAISaveDocumentAsAction
		*/
		ASErr SaveDocumentAsSVG(const ai::UnicodeString& name);

		/** Gathers action parameters and saves document as compressed SVG.
			@param name IN name to save the document as.
			@return kNoErr on success, other error code otherwise.
			@see kAISaveDocumentAsAction
		*/
		ASErr SaveDocumentAsSVGZ(const ai::UnicodeString& name);

		/** Gathers common parameters and saves a document.
			@param name IN name to save the document as.
			@param vpb IN the value parameter block.
			@return kNoErr on success, other error code otherwise.
			@see kAISaveDocumentAsAction
			@see SnpDocumentActionHelper::VPB
		*/
		ASErr SaveDocumentAs(const ai::UnicodeString& name, SnpDocumentActionHelper::VPB& vpb);

		/** Saves a document as using the given parameters.
			@param dialogStatus IN whether parameters dialog appears by default.
			@param parameters IN the value parameter block.
			@return kNoErr on success, other error code otherwise.
			@see kAISaveDocumentAsAction
			@see SnpDocumentActionHelper::VPB
		*/
		ASErr SaveDocumentAs(ActionDialogStatus dialogStatus, AIActionParamValueRef parameters);

		/** Gathers action parameters and closes document.
			@return kNoErr on success, other error code otherwise.
			@see kAICloseDocumentAction
		*/
		ASErr CloseDocument();

		/** Closes a document using the given parameters.
			@param dialogStatus IN whether parameters dialog appears by default.
			@param parameters IN the value parameter block.
			@return kNoErr on success, other error code otherwise.
			@see kAICloseDocumentAction
			@see SnpDocumentActionHelper::VPB
		*/
		ASErr CloseDocument(ActionDialogStatus dialogStatus, AIActionParamValueRef parameters);

		/** Gathers action parameters and places a file.
			@return kNoErr on success, other error code otherwise.
			@see kAIPlaceDocumentAction
		*/
		ASErr PlaceFile(const ai::UnicodeString& name);

		/** Places a file using the given parameters.
			@param dialogStatus IN whether parameters dialog appears by default.
			@param parameters IN the value parameter block.
			@return kNoErr on success, other error code otherwise.
			@see kAIPlaceDocumentAction
			@see SnpDocumentActionHelper::VPB
		*/
		ASErr PlaceFile(ActionDialogStatus dialogStatus, AIActionParamValueRef parameters);

		/** Gathers action parameters and exports the selection to an AI file.
			@param name IN the name to save the document as.
			@return kNoErr on success, other error code otherwise.
			@see kAIExportDocumentAction
		*/
		ASErr ExportSelectionToAI(const ai::UnicodeString& name);

		/** Gathers action parameters and exports the selection to an EPS file.
			@param name IN the name to save the document as.
			@return kNoErr on success, other error code otherwise.
			@see kAIExportDocumentAction
		*/
		ASErr ExportSelectionToEPS(const ai::UnicodeString& name);

		/** Exports the selection.
			@param name IN the name to save the document as.
			@return kNoErr on success, other error code otherwise.
			@see kAIExportDocumentAction
		*/
		ASErr ExportSelection(const ai::UnicodeString& name, const char* format, const char* extn);

		/** Gathers action parameters and exports a document to AutoCAD DWG.
			@param name IN the name to save the document as.
			@return kNoErr on success, other error code otherwise.
			@see kAIExportDocumentAction
		*/
		ASErr ExportDocumentToDWG(const ai::UnicodeString& name);

		/** Gathers action parameters and exports a document to AutoCAD DXF.
			@param name IN the name to save the document as.
			@return kNoErr on success, other error code otherwise.
			@see kAIExportDocumentAction
		*/
		ASErr ExportDocumentToDXF(const ai::UnicodeString& name);

		/** Gathers action parameters and exports a document to PSD.
			@param name IN the name to save the document as.
			@return kNoErr on success, other error code otherwise.
			@see kAIExportDocumentAction
		*/
		ASErr ExportDocumentToPSD(const ai::UnicodeString& name);

		/** Exports a document using the given parameters.
			@param name IN the name to save the document as.
			@param format IN the format to export the document to.
			@param extn IN the extension for the new file.
			@return kNoErr on success, other error code otherwise.
			@see kAIExportDocumentAction
			@see SnpDocumentActionHelper::VPB
		*/
		ASErr ExportDocument(const ai::UnicodeString& name, const char* format, const char* extn);
		
		/** Exports a document using the given parameters.
			@param dialogStatus IN whether parameters dialog appears by default.
			@param parameters IN the value parameter block.
			@return kNoErr on success, other error code otherwise.
			@see kAIExportDocumentAction
			@see SnpDocumentActionHelper::VPB
		*/
		ASErr ExportDocument(ActionDialogStatus dialogStatus, AIActionParamValueRef parameters);

};

#endif // __SnpDocumentActionHelper_h__

// End SnpDocumentActionHelper.h

