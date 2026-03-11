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

#ifndef __SnpExportFXG_h__
#define __SnpExportFXG_h__

#include "IllustratorSDK.h"

#include "SDKErrors.h"

// Framework includes:
#include "SnpRunnable.h"
#include "SnippetRunnerSuites.h"
#include "SnippetRunnerLog.h"
#include "SnippetRunnerParameter.h"
#include "SnpSelectionHelper.h"
#include "SnpArtSetHelper.h"
#include "SnpArtHelper.h"
#include "SnpDocumentHelper.h"
#include "SnpFileHelper.h"
#include "SnippetRunnerUnitTestManager.h"


/** Implements a basic stream interface which is passed to AIFXGFileFormatSuite::ExportFXGToStream.
*/
struct SnpExportFXGStreamInterface : public AIFXGWriteStreamInterface
{
	/** Constructor.
	*/
	SnpExportFXGStreamInterface();

	/** Destructor.
	*/
	~SnpExportFXGStreamInterface();

	/** Prompts the user for a file save location and opens a data filter
		for the file.
		@param imageName IN the name of the image.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr OpenFile();

	/** Closes the current data filter if open.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr CloseFile();
	
	/** Writes the stream bytes to the file using the current data filter.
		@param buffer IN buffer to write to the stream.
		@param maxToWrite IN max number of bytes to write.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr WriteBytesToFile(const ASUInt8* buffer, const size_t  maxToWrite);

	/** Returns the full path to the data filter output file.
		@return the full file path.
	*/
	const ai::FilePath& GetCurrentFilePath() const {return fFilePath;}

	/** Sets the file path of the file to write to using the current data filter.
		@param filePath IN path to stream output file.
		@return kNoErr on success, other ASErr otherwise.
	*/
	ASErr SetFilePath(const ai::FilePath& filePath) {fFilePath = filePath; return kNoErr;};

private:

	/** File path used by current data filter.
	*/
	ai::FilePath fFilePath;

	/** Pointer to current data filter.
	*/
	AIDataFilter* fDataFilter;

	/** 
	*/
	ASUInt32 fCount;
};


/** Demonstrates the use of the FXG File Format suite.
  	@ingroup sdk_snippet
*/
class SnpExportFXG
{
	public:

		/** Constructor.
		 */
		SnpExportFXG(): fStreamInterface(NULL){}

		/** Destructor.
		 */
		~SnpExportFXG(){if(fStreamInterface) delete fStreamInterface;}

		/** Exports the FXG data of the entire document to a file.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr ExportDocumentToFXGFile();

		/** Exports the FXG data of the selected art to files.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr ExportArtToFXGFile();

		/** Exports the FXG data of the entire document to a file and the log
			using the stream interface implemented in this snippet.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr ExportDocumentToFXGStream();

		/** Exports the FXG data of the selected art to files and the log using
			the stream interface implemented in this snippet.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr ExportArtToFXGStream();

		//----------------------------- Utility functions ---------------------------------

		/** Prompts the user for the save parameters.
			@param params IN/OUT buffer in which to return the save parameters.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr GetParams(AIFXGFileFormatSaveParams& params);

		/** Prompts user for the location to save the FXG file.
			@param filePath IN/OUT buffer in which to return the file path.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr GetFilePath(ai::FilePath& filePath);

		/** Implements and returns a stream interface.
			@return pointer to a stream interface.
		 */
		AIFXGWriteStreamInterface* GetStreamInterface() ;

		/** Sets the file path of the output file in the stream interface.
			@param filePath IN the output file path.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr SetFilePath(const ai::FilePath& filePath);

	private:

		/** Pointer to the stream interface.
		*/
		SnpExportFXGStreamInterface* fStreamInterface;

};

/** Provides an interface between the snippet and the FXGFileFormat suite for exporting fxg streams,
	implements some of the functions from the _t_AIFXGWriteStreamInterface struct.
*/
namespace FXGStreamInterface
{
	/** Calls SnpExportFXGStreamInterface::OpenFile.
		@param stream IN stream pointer as passed in through AIFXGFileFormatSuite::ExportFXGToStream.
		@param assetType IN the asset type to open.
		@param assetName IN/OUT name of the asset to open.
		@param mimeType IN the mime type of the asset to open.
		@return kNoErr on success, other ASErr otherwise.
	*/
	AIAPI AIErr	BeginAsset(AIFXGWriteStreamInterface* stream, ai::int16 assetType, const ai::UnicodeString& assetName, const ai::UnicodeString& mimeType);

	/** Calls SnpExportFXGStreamInterface::CloseFile.
		@param stream IN stream pointer as passed in through AIFXGFileFormatSuite::ExportFXGToStream.
		@param assetType IN asset type to close.
		@return kNoErr on success, other ASErr otherwise.
	*/
	AIAPI AIErr	EndAsset(AIFXGWriteStreamInterface* stream, ai::int16 assetType);

	/** Calls SnpExportFXGStreamInterface::GetCurrentFilePath.
		@param stream IN stream pointer as passed in through AIFXGFileFormatSuite::ExportFXGToStream.
		@param publicID IN/OUT buffer to return the name of the current stream.
		@return kNoErr on success, other ASErr otherwise.
	*/
	AIAPI AIErr	GetPublicID(AIFXGWriteStreamInterface* stream, ai::UnicodeString& publicID);

	/** Calls SnpExportFXGStreamInterface::WriteBytesToFile.
		@param stream IN stream pointer as passed in through AIFXGFileFormatSuite::ExportFXGToStream.
		@param buffer IN buffer to write to stream.
		@param maxToWrite IN the max number of bytes to write.
		@return number of bytes written.
	*/
	AIAPI size_t WriteBytes(AIFXGWriteStreamInterface* stream, const ASUInt8* buffer, const size_t maxToWrite);

	template<typename T>
	/** Changes an object of one type to be a pointer to another type, in the context of this snippet it
		changes an object of type SnpExportFXGStreamInterface to a pointer to an AIFXGWriteStreamInterface.
	*/
	class TPtr
	{
	public:
		/** Constructor
		*/
		TPtr(void* ptr)
		{
			// Change the object passed in to a pointer to an object of the type specified at instantiation.
			fPtr = reinterpret_cast<T*>(ptr);
		}
		/** Overload -> function to return the pointer.
		*/
		T* operator ->() {return fPtr;}
	private:
		/** Pointer to type specified at instantiation.
		*/
		T* fPtr;
	};
}

#endif
