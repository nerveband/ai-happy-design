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

#ifndef  __SnpDocumentHelper_h__
#define  __SnpDocumentHelper_h__

/** Helper class that shows how to create, open, save and close Illustrator documents.

	@see AIDocumentListSuite, AINewDocumentPreset, AIDocumentSuite
  	@ingroup sdk_snippet
 */
class SnpDocumentHelper
{
	public:
		/** Creates a new Illustrator document from the given preset.
			@param preset IN name of the document preset to be used
			@param document OUT handle of the created document
			@return kNoErr on success, other error code otherwise.
			@see AIDocumentListSuite::New
		*/
		ASErr NewDocument(const ai::UnicodeString& preset, AIDocumentHandle& document);

		/** Creates a new Illustrator document from the given preset
			with a custom width and height.
			@param preset IN name of the document preset to be used.
			@param width IN custom document width points.
			@param height IN custom document height in points.
			@param document OUT handle of the created document.
			@return kNoErr on success, other error code otherwise.
			@see AIDocumentListSuite::New
		*/
		ASErr NewCustomDocument(const ai::UnicodeString& preset, AIReal width, AIReal height, AIDocumentHandle& document);

		/** Creates a new Illustrator document from the given preset
			with custom settings for width, height and artboards.
			@param preset IN name of the document preset to be used.
			@param width IN custom document width points.
			@param height IN custom document height in points.
			@param numArtboards IN number of artboards in the new custom document.
			@param artboardLayout IN the layout of the artboards in the new document e.g. a grid organised by row.
			@param artboardSpacing IN the spacing between the artboards.
			@param numRowsOrCols IN the number of either rows or cols to span the grid of artboards over.
			@param document OUT handle of the created document.
			@return kNoErr on success, other error code otherwise.
			@see AIDocumentListSuite::New
		*/
		ASErr NewCustomDocument(const ai::UnicodeString& preset, AIReal width, AIReal height, ASInt32 numArtboards, AIArtboardLayout artboardLayout, AIReal artboardSpacing, ASInt32 numRowsOrCols, AIDocumentHandle& document);

		/** Opens a document of the given file name in SnippetRunner's assets folder.
			@param name IN file name of the document to be opened.
			@param document OUT handle of the created document.
			@return kNoErr on success, other error code otherwise.
			@see AIDocumentListSuite::Open
		*/
		ASErr OpenDocument(const ai::UnicodeString& name, AIDocumentHandle& document);

		/** Make the given document the one the user is working with.
			@param document IN handle of the document to be activated.
			@return kNoErr on success, other error code otherwise.
			@see AIDocumentListSuite::Activate
		*/
		ASErr ActivateDocument(AIDocumentHandle document);

		/** Print the given document.
			@param document IN handle of the document to be printed.
			@return kNoErr on success, other error code otherwise.
			@see AIDocumentListSuite::Print
		*/
		ASErr PrintDocument(AIDocumentHandle document);

		/** Save the given document.
			@param document IN handle of the document to be saved.
			@return kNoErr on success, other error code otherwise.
			@see AIDocumentListSuite::Save
		*/
		ASErr SaveDocument(AIDocumentHandle document);

		/** Close the given document.
			@param document IN handle of the document to be closed.
			@return kNoErr on success, other error code otherwise.
			@see AIDocumentListSuite::Close
		*/
		ASErr CloseDocument(AIDocumentHandle document);

		/** Close all open documents.
			@return kNoErr on success, other error code otherwise.
			@see AIDocumentListSuite::CloseAll
		*/
		ASErr CloseAllDocuments();

		/** Gets the handle of the document the user is currently working with.
			@param document OUT handle of document user is working with.
			@return kNoErr on success, kNoDocumentErr if no documents open, other error otherwise.
			@see AIDocumentSuite::GetDocument
		*/
		ASErr GetCurrentDocument(AIDocumentHandle& document);

		/** Gets the handle of the layer the user is currently working with.
			@param layer OUT handle of the layer the user is currently working with.
			@return kNoErr on success, kNoDocumentErr if no documents open, other error otherwise.
			@see AILayerSuite::GetCurrentLayer
		*/
		ASErr GetCurrentLayer(AILayerHandle& layer);

		/** Checks if a file can be placed in the current layer.
			@return true if the file can be placed, false otherwise.
		*/
		ASBoolean CanPlaceFile(void);

		/** Places a file in the current document.
			@param name IN name of the file from SnippetRunner's assets folder to be placed, if the file
				does not exist the user will be prompted to locate the file to be placed.
			@param art OUT handle to the placed art object
			@return kNoErr on success, other error otherwise.
			@see AIPlacedSuite::ExecPlaceRequest
		*/
		ASErr PlaceFile(const ai::UnicodeString& name, AIArtHandle& art);

		/** Gets the path to the given file in SnippetRunner's assets folder.
			@param name IN name of the file.
			@return path to the given file in SnippetRunner's assets folder
		*/
		ai::FilePath GetAssetPath(const ai::UnicodeString& name);
};

#endif // __SnpDocumentHelper_h__

// End SnpDocumentHelper.h

