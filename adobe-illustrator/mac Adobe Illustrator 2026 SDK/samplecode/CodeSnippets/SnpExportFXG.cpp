//========================================================================================
//  
//  $File$
//
//  $Revision$
//
//  Copyright 2008 Adobe Systems Incorporated. All rights reserved.
//  
//  NOTICE:  Adobe permits you to use, modify, and distribute this file in accordance 
//  with the terms of the Adobe license agreement accompanying it.  If you have received
//  this file from a source other than Adobe, then your use, modification, or 
//  distribution of it requires the prior written permission of Adobe.
//  
//========================================================================================

#include "IllustratorSDK.h"

#include "SnpExportFXG.h"


/*
*/
ASErr SnpExportFXG::ExportDocumentToFXGFile()
{
	ASErr result = kNoErr;
	if(sAIFXGFileFormat != NULL)
	{
		try {
			// Initialise AIFXGFileFormatSaveParams object to default params.
			AIFXGFileFormatSaveParams params;
			result = sAIFXGFileFormat->CreateFXGSaveParams(&params);
			aisdk::check_ai_error(result);

			// Get AIFXGFileFormatSaveParams values from user.
			result = this->GetParams(params);
			aisdk::check_ai_error(result);

			// Get file path.
			ai::FilePath filePath;
			result = this->GetFilePath(filePath);
			aisdk::check_ai_error(result);

			// Export entire document to file.
			result = sAIFXGFileFormat->ExportFXGToFile(NULL, filePath, params, kAIFXGExportFlagDefault,
				ai::UnicodeString(""), ai::UnicodeString("Writing document to FXG file"), NULL);
			aisdk::check_ai_error(result);
		}
		catch (ai::Error& ex) {
			result = ex;
		}	
	}
	else
	{
		result = kSPSuiteNotFoundError;
	}
	return result;
}

/*
*/
ASErr SnpExportFXG::ExportArtToFXGFile()
{
	ASErr result = kNoErr;
	if(sAIFXGFileFormat != NULL)
	{
		try {
			// Initialise AIFXGFileFormatSaveParams object to default params.
			AIFXGFileFormatSaveParams params;
			result = sAIFXGFileFormat->CreateFXGSaveParams(&params);
			aisdk::check_ai_error(result);

			// Get AIFXGFileFormatSaveParams values from user.
			result = this->GetParams(params);
			aisdk::check_ai_error(result);

			// Get file path.
			ai::FilePath filePath;
			result = this->GetFilePath(filePath);
			aisdk::check_ai_error(result);

			// Get selected art.
			const short numSpecs = 1;
			AIArtSpec artSpec[numSpecs] = { 
				{kAnyArt, kArtSelected, kArtSelected},
			};
			SnpArtSetHelper artSet(artSpec, numSpecs);

			ai::UnicodeString extension(filePath.GetFileExtension());
		
			// Export art to file.
			size_t artCount = artSet.GetCount();
			if (artCount > 1) {
				// Iterate through each art item selected and export to separate FXG files.
				for (size_t i = 0; i < artCount; i++) {
					ai::NumberFormat numFormat;
					ai::UnicodeString indexString;
					indexString = numFormat.toString(i+1, indexString);
				
					ai::FilePath artFilePath(filePath);
					ai::UnicodeString nameOnly(artFilePath.GetFileNameNoExt());
					nameOnly.append(indexString);
					artFilePath.RemoveExtension();
					artFilePath.RemoveComponent();
					artFilePath.AddComponent(nameOnly);
					artFilePath.AddExtension(extension);

					AIArtHandle art = artSet[i];
					result = sAIFXGFileFormat->ExportFXGToFile(art, artFilePath, params, kAIFXGExportFlagDefault,
						ai::UnicodeString(""), ai::UnicodeString("Writing art to FXG files"), NULL);
					aisdk::check_ai_error(result);
				}
			}
			else if (artCount == 1) {
				// Export the FXG of the selected art.
				AIArtHandle art = artSet[0];
				result = sAIFXGFileFormat->ExportFXGToFile(art, filePath, params, kAIFXGExportFlagDefault,
					ai::UnicodeString(""), ai::UnicodeString("Writing art to FXG file"), NULL);
				aisdk::check_ai_error(result);
			}
			else {
				SnippetRunnerParameter::Instance()->Alert(ai::UnicodeString("Please select art in a document."));
				result = kBadParameterErr;
			}
		}
		catch (ai::Error& ex) {
			result = ex;
		}	
	}
	else
	{
		result = kSPSuiteNotFoundError;
	}
	return result;
}

/*
*/
ASErr SnpExportFXG::ExportDocumentToFXGStream()
{
	ASErr result = kNoErr;
	if(sAIFXGFileFormat != NULL)
	{
		try {
			// Initialise AIFXGFileFormatSaveParams object to default params.
			AIFXGFileFormatSaveParams params;
			result = sAIFXGFileFormat->CreateFXGSaveParams(&params);
			aisdk::check_ai_error(result);

			// Get AIFXGFileFormatSaveParams values from user.
			result = this->GetParams(params);
			aisdk::check_ai_error(result);

			// Get file path.
			ai::FilePath filePath;
			result = this->GetFilePath(filePath);
			aisdk::check_ai_error(result);

			// Set file path.
			result = this->SetFilePath(filePath);
			aisdk::check_ai_error(result);

			// Get stream.
			AIFXGWriteStreamInterface* stream = this->GetStreamInterface();

			// Export entire document to stream.
			result = sAIFXGFileFormat->ExportFXGToStream(NULL, stream, params, kAIFXGExportFlagDefault,
				ai::UnicodeString(""), ai::UnicodeString("Writing document to FXG stream"), NULL);
			aisdk::check_ai_error(result);
		}
		catch (ai::Error& ex) {
			result = ex;
		}
	}
	else
	{
		result = kSPSuiteNotFoundError;
	}
	return result;
}

/*
*/
ASErr SnpExportFXG::ExportArtToFXGStream()
{
	ASErr result = kNoErr;
	if(sAIFXGFileFormat != NULL)
	{
		try {
			// Initialise AIFXGFileFormatSaveParams object to default params.
			AIFXGFileFormatSaveParams params;
			result = sAIFXGFileFormat->CreateFXGSaveParams(&params);
			aisdk::check_ai_error(result);

			// Get AIFXGFileFormatSaveParams values from user.
			result = this->GetParams(params);
			aisdk::check_ai_error(result);

			// Get file path.
			ai::FilePath filePath;
			result = this->GetFilePath(filePath);
			aisdk::check_ai_error(result);

			// Set file path.
			result = this->SetFilePath(filePath);
			aisdk::check_ai_error(result);

			// Get stream.
			AIFXGWriteStreamInterface* stream = this->GetStreamInterface();

			// Get selected art.
			const short numSpecs = 1;
			AIArtSpec artSpec[numSpecs] = { 
				{kAnyArt, kArtSelected, kArtSelected},
			};
			SnpArtSetHelper artSet(artSpec, numSpecs);

			// Export art to file.
			size_t artCount = artSet.GetCount();
			if (artCount > 1) {
				// Iterate through each art item selected and export to separate FXG streams.
				for (size_t i = 0; i < artCount; i++) {
					ai::NumberFormat numFormat;
					ai::UnicodeString indexString;
					indexString = numFormat.toString(i+1, indexString);

					AIArtHandle art = artSet[i];
					result = sAIFXGFileFormat->ExportFXGToStream(art, stream, params, kAIFXGExportFlagDefault,
						ai::UnicodeString(""), ai::UnicodeString("Writing art to FXG stream"), NULL);
					aisdk::check_ai_error(result);
				}
			}
			else if (artCount == 1) {
				// Export the FXG of the selected art to a stream.
				AIArtHandle art = artSet[0];
				result = sAIFXGFileFormat->ExportFXGToStream(art, stream, params, kAIFXGExportFlagDefault,
					ai::UnicodeString(""),ai::UnicodeString("Writing art to FXG stream"), NULL);
				aisdk::check_ai_error(result);
			}
			else {
				SnippetRunnerParameter::Instance()->Alert(ai::UnicodeString("Please select art in a document."));
				result = kBadParameterErr;
			}
		}
		catch (ai::Error& ex) {
			result = ex;
		}	
	}
	else
	{
		result = kSPSuiteNotFoundError;
	}
	return result;
}

//--------------------------------------- Utility Functions ---------------------------------------------------

/*
*/
ASErr SnpExportFXG::GetParams(AIFXGFileFormatSaveParams& params)
{
	ASErr result = kNoErr;
	if(sAIFXGFileFormat != NULL)
	{
		try {
			SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();

			// Set save version.
			AIFXGVersion version = kAIFXGDefaultVersion;
			std::vector<std::string> versionChoices;
			versionChoices.push_back("kAIFXGVersion1pt0");
			versionChoices.push_back("kAIFXGDefaultVersion");
			version = static_cast<AIFXGVersion>((parameter->GetChoice("FXG Version:", versionChoices, 0)+1));
			result = sAIFXGFileFormat->SetFXGSaveVersion(params, version);
			aisdk::check_ai_error(result);

			// Set Preserve AI Editing Capabilities flag.
			ASBoolean preserveEditCap = parameter->GetBoolean(ai::UnicodeString("Preserve AI Editing Capabilities?"));
			result = sAIFXGFileFormat->SetFXGSavePreserveAIEditCapability(params, preserveEditCap);
			aisdk::check_ai_error(result);

			// Create vector of preserve policies.
			std::vector<std::string> policyChoices;
			policyChoices.push_back("kAIFXGPreserveAppearanceExpand");
			policyChoices.push_back("kAIFXGPreserveAppearanceRasterize");
			policyChoices.push_back("kAIFXGPreserveEditibility");

			// Set Text Policy.
			AIFXGPreservePolicy textPreservePolicy = kAIFXGPreserveEditibility;
			textPreservePolicy = static_cast<AIFXGPreservePolicy>((parameter->GetChoice("Text Preserve Policy:", policyChoices, 0)+1));
			result = sAIFXGFileFormat->SetFXGSaveTextPolicy(params, textPreservePolicy);
			aisdk::check_ai_error(result);

			// Set Filters Policy.
			AIFXGPreservePolicy filtersPreservePolicy = kAIFXGPreserveEditibility;
			filtersPreservePolicy = static_cast<AIFXGPreservePolicy>((parameter->GetChoice("Filters Preserve Policy:", policyChoices, 0)+1));
			result = sAIFXGFileFormat->SetFXGSaveFilterPolicy(params, filtersPreservePolicy);
			aisdk::check_ai_error(result);

			// Set Gradients Policy.
			AIFXGPreservePolicy gradientsPreservePolicy = kAIFXGPreserveEditibility;
			gradientsPreservePolicy = static_cast<AIFXGPreservePolicy>((parameter->GetChoice("Gradients Preserve Policy:", policyChoices, 0)+1));
			result = sAIFXGFileFormat->SetFXGSaveGradientPolicy(params, gradientsPreservePolicy);
			aisdk::check_ai_error(result);

			// Set Include Symbols flag.
			ASBoolean includeSymbols = parameter->GetBoolean(ai::UnicodeString("Include Symbols?"));
			result = sAIFXGFileFormat->SetFXGSaveIncludeSymbol(params, includeSymbols);
			aisdk::check_ai_error(result);

			// Set Include XMP flag.
			ASBoolean includeXMP = parameter->GetBoolean(ai::UnicodeString("Include XMP?"));
			result = sAIFXGFileFormat->SetFXGSaveIncludeXMP(params, includeXMP);
			aisdk::check_ai_error(result);

			// Create vector of Rasterize Resolution values.
			std::vector<std::string> rasterResChoices;
			rasterResChoices.push_back("kAIFXGResolutionScreen (72 ppi)");
			rasterResChoices.push_back("kAIFXGResolutionMedium (150 ppi)");
			rasterResChoices.push_back("kAIFXGResolutionHigh (300 ppi)");
			rasterResChoices.push_back("kAIFXGResolutionCustom");

			// Set Rasterize Resolution value.
			AIFXGRasterizeResolution rasterResolution = kAIFXGResolutionScreen;
			ASInt32 rasterizeResolution = parameter->GetChoice("Select Rasterize Resolution:", rasterResChoices, 0);
			switch (rasterizeResolution) {
				case 0: {
					rasterResolution = kAIFXGResolutionScreen;
					break;	
				}
				case 1: {
					rasterResolution = kAIFXGResolutionMedium;
					break;
				}
				case 2: {
					rasterResolution = kAIFXGResolutionHigh;
					break;					
				}
				case 3: {
					rasterResolution = kAIFXGResolutionCustom;
					break;
				}
				default:
					break;
			}
			ASInt32 customResolution = 0;
			if (rasterResolution == kAIFXGResolutionCustom) {
				// Prompt user for custom resolution value.
				customResolution = parameter->GetInt(ai::UnicodeString("Enter custom resolution value:"), 72);
			}

			result = sAIFXGFileFormat->SetFXGSaveRasterizeResolution(params, rasterResolution, customResolution);
			aisdk::check_ai_error(result);

			// Set Downsample Linked Images flag.
			ASBoolean downsample = parameter->GetBoolean(ai::UnicodeString("Downsample Linked Images?"));
			result = sAIFXGFileFormat->SetFXGSaveDownsampleLinkedImages(params, downsample);
			aisdk::check_ai_error(result);

			// Set Write Images flag.
			ASBoolean writeImages = parameter->GetBoolean(ai::UnicodeString("Write Images?"));
			result = sAIFXGFileFormat->SetFXGSaveWriteImages(params, writeImages);
			aisdk::check_ai_error(result);

		}
		catch (ai::Error& ex) {
			result = ex;
		}	
	}
	else
	{
		result = kSPSuiteNotFoundError;
	}
	return result;
}

/*
*/
ASErr SnpExportFXG::GetFilePath(ai::FilePath& filePath)
{
	ASErr result = kNoErr;
	try {
		ai::UnicodeString fileName(ai::UnicodeString("FXGData.fxg"));
		if (SnippetRunnerUnitTestManager::Instance()->GetUnitTestState() == 2/* kSnpRunUnitTestRunning */) {
			// Unit tests running so set the file path to the assets folder.
			SnpDocumentHelper docHelper;
			filePath = docHelper.GetAssetPath(fileName);
		}
		else {
			SnpFileHelper fileHelper;
			result = fileHelper.PutFileDialog(ai::UnicodeString("Select location to save the FXG data: "), fileName, filePath);
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
AIFXGWriteStreamInterface* SnpExportFXG::GetStreamInterface() 
{
	if(fStreamInterface == NULL)
	{
		fStreamInterface = new SnpExportFXGStreamInterface();
		fStreamInterface->BeginAsset		= FXGStreamInterface::BeginAsset;
		fStreamInterface->EndAsset			= FXGStreamInterface::EndAsset;
		fStreamInterface->WriteBytes		= FXGStreamInterface::WriteBytes;
		fStreamInterface->GetPublicID		= FXGStreamInterface::GetPublicID;
	}										

	return fStreamInterface;
}

/*
*/
ASErr SnpExportFXG::SetFilePath(const ai::FilePath& filePath)
{
	ASErr result = kNoErr;
	try {
		if (fStreamInterface == NULL)
			this->GetStreamInterface();

		result = fStreamInterface->SetFilePath(filePath);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

//------------------------------------ SnpExportFXGStreamInterface functions --------------------------------------------------

/*
*/
SnpExportFXGStreamInterface::SnpExportFXGStreamInterface(): fDataFilter(NULL), fCount(0) {}

/*
*/
SnpExportFXGStreamInterface::~SnpExportFXGStreamInterface(){}

/*
*/
ASErr SnpExportFXGStreamInterface::OpenFile()
{
	ASErr result = kNoErr;
	try {
		ai::UnicodeString extension = fFilePath.GetFileExtension();
		ai::NumberFormat numFormat;
		ai::UnicodeString indexString;
		indexString = numFormat.toString(++fCount, indexString);
		
		ai::FilePath artFilePath(fFilePath);
		ai::UnicodeString nameOnly(artFilePath.GetFileNameNoExt());
		nameOnly.append(indexString);
		artFilePath.RemoveExtension();
		artFilePath.RemoveComponent();
		artFilePath.AddComponent(nameOnly);
		artFilePath.AddExtension(extension);

		// Create a new data filter.
        result = sAIDataFilter->NewFileDataFilter(artFilePath, "write", 'ART5', 'TEXT', &fDataFilter);
		aisdk::check_ai_error(result);
        result = sAIDataFilter->LinkDataFilter(NULL, fDataFilter);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}

/*
*/
ASErr SnpExportFXGStreamInterface::CloseFile()
{
	ASErr result = kNoErr;
	try {
		if (fDataFilter) {
			result = sAIDataFilter->UnlinkDataFilter(fDataFilter, NULL);
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
ASErr SnpExportFXGStreamInterface::WriteBytesToFile(const ASUInt8* buffer, const size_t  maxToWrite)
{
	ASErr result = kNoErr;
	try {
		if (fDataFilter) {
			const char* newBuffer = reinterpret_cast<const char*>(buffer);

			// Write the FXG stream to the file via the filter.
			size_t localLen = maxToWrite;
			result = sAIDataFilter->WriteDataFilter(fDataFilter, const_cast<char*>(newBuffer), &localLen);
			aisdk::check_ai_error(result);

			// Also write the FXG stream to the log.
			SnippetRunnerLog::Instance()->Write(newBuffer);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	return result;
}


//------------------------------------ FXGStreamInterface namespace --------------------------------------------------

/**
*/
namespace FXGStreamInterface
{
	/*
	*/
	AIAPI AIErr		BeginAsset(AIFXGWriteStreamInterface* stream, ai::int16 assetType, const ai::UnicodeString& assetName, const ai::UnicodeString& mimeType)
	{
		ASErr result = kNoErr;
		SDK_ASSERT(stream);

		TPtr<SnpExportFXGStreamInterface> streamInterface(stream);
		result = streamInterface->OpenFile();
		return result;

	}

	/*
	*/
	AIAPI AIErr		EndAsset(AIFXGWriteStreamInterface* stream, ai::int16 assetType)
	{
		ASErr result = kNoErr;
		SDK_ASSERT(stream);

		TPtr<SnpExportFXGStreamInterface> streamInterface(stream);
		result = streamInterface->CloseFile();
		return result;
	}

	/*
	*/
	AIAPI AIErr		GetPublicID(AIFXGWriteStreamInterface* stream, ai::UnicodeString& publicID)
	{
		ASErr result = kNoErr;
		TPtr<SnpExportFXGStreamInterface> streamInterface(stream);
		publicID = streamInterface->GetCurrentFilePath().GetFullPath();
		return result;
	}

	/*
	*/
	AIAPI size_t	WriteBytes(AIFXGWriteStreamInterface* stream, const ASUInt8* buffer, const size_t  maxToWrite)
	{
		ASErr result = kNoErr;
		TPtr<SnpExportFXGStreamInterface> streamInterface(stream);
		result = streamInterface->WriteBytesToFile(buffer, maxToWrite);
		if(result == kNoErr)
			return maxToWrite;
		
		return 0;
	}	
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpExportFXG available to the SnippetRunner framework.
 */
class _SnpRunnableExportFXG : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableExportFXG () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableExportFXG () {}

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

		/** Selects art in unit test document.
		*/
		ASErr SetupContext(SnpRunnable::Context& runnableContext);
};

/*
*/
std::string _SnpRunnableExportFXG::GetName() const
{
	return "ExportFXG";
}

/*
*/
std::string _SnpRunnableExportFXG::GetDescription() const 
{
	return "Demonstrates export of documents and art items to FXG files either directly or via a stream.";
}

/*
*/
SnpRunnable::Operations _SnpRunnableExportFXG::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("ExportDocumentToFXGFile", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("ExportArtToFXGFile", "document with selected art", kSnpRunCustomContext));
	operations.push_back(Operation("ExportDocumentToFXGStream", "document", kSnpRunNewDocumentContext));
	operations.push_back(Operation("ExportArtToFXGStream", "document with selected art", kSnpRunCustomContext));
	return operations;
}

/*
*/
std::string _SnpRunnableExportFXG::GetDefaultOperationName() const 
{
	return "ExportDocumentToFXGFile";
	// Note: this must be one of the operation names returned by GetOperations.
}

/*
*/
std::vector<std::string> _SnpRunnableExportFXG::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("FXG Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableExportFXG::CanRun(SnpRunnable::Context& runnableContext)
{
	SnpSelectionHelper selectionHelper;
	if (!selectionHelper.IsDocumentSelected())
		return false;

	if (("ExportArtToFXGFile" == runnableContext.GetOperation().GetName() || 
		"ExportArtToFXGStream" == runnableContext.GetOperation().GetName()) &&
		selectionHelper.IsArtSelected())
			return true;
	else 
			return true;
	return true;
	// Note: this will control the enabling and disabling of the Run button in the SnippetRunner panel
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableExportFXG::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpExportFXG instance;
	if ("ExportDocumentToFXGFile" == runnableContext.GetOperation().GetName()) {
		result = instance.ExportDocumentToFXGFile();
	}
	else if ("ExportArtToFXGFile" == runnableContext.GetOperation().GetName()) {
		result = instance.ExportArtToFXGFile();
	}
	else if ("ExportDocumentToFXGStream" == runnableContext.GetOperation().GetName()) {
		result = instance.ExportDocumentToFXGStream();
	}
	else if ("ExportArtToFXGStream" == runnableContext.GetOperation().GetName()) {
		result = instance.ExportArtToFXGStream();
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

/*
*/
ASErr _SnpRunnableExportFXG::SetupContext(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	try {
		switch (runnableContext.GetOperation().GetUnitTestContext()){
			case SnpRunnable::kSnpRunCustomContext: {
					SnpDocumentHelper docHelper;
					SnpSelectionHelper selectionHelper;
					AIDocumentHandle document = NULL;
					SDK_ASSERT(!docHelper.OpenDocument(ai::UnicodeString("unittest-1.ai"), document));
					SDK_ASSERT(document);
					if ("ExportArtToFXGFile" == runnableContext.GetOperation().GetName() ||
						"ExportArtToFXGStream" == runnableContext.GetOperation().GetName()) {
						const short numSpecs = 1;
						AIArtSpec artSpec[numSpecs] = { 
							{kAnyArt, NULL, NULL},
						};
						SnpArtSetHelper artSet(artSpec, numSpecs);
						
						if (artSet.GetCount() > 0) {
							artSet.SelectAll();
						}
						else {
							SnpArtHelper artHelper;
							AIArtHandle art = NULL;
							result = artHelper.NewRectangle(art);
							aisdk::check_ai_error(result);
							
							selectionHelper.SelectArt(art);
							aisdk::check_ai_error(result);
						}
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
static _SnpRunnableExportFXG instance_SnpRunnableExportFXG;

// End SnpExportFXG.cpp
