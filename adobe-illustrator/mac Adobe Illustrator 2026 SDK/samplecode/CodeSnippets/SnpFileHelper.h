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

#ifndef __SnpFileHelper_h__
#define __SnpFileHelper_h__

/** Shows how to browse for files and folders using AIUserSuite.

	@see AIUserSuite
  	@ingroup sdk_snippet
 */
class SnpFileHelper
{
	public:
		/** Prompts for a path to a folder.
			@param title IN the prompt dialog title.
			@param ioFilePath OUT path to the chosen folder.
			@return kNoErr on success, other error otherwise.
			@see AIUserSuite::GetDirectoryDialog
		*/
		ASErr GetDirectoryDialog(const ai::UnicodeString &title, ai::FilePath& ioFilePath);

		/** Prompts for a path to an input file.
			@param title IN the prompt dialog title.
			@param ioFilePath OUT path to the chosen folder.
			@return kNoErr on success, other error otherwise.
			@see AIUserSuite::GetFileDialog
		*/
		ASErr GetFileDialog(const ai::UnicodeString &title, ai::FilePath& ioFilePath);

		/** Prompts for a path to an output file.
			@param title IN the prompt dialog title.
			@param defaultName IN default file name to save into.
			@param ioFilePath IN/OUT path to the chosen folder.
			@return kNoErr on success, other error otherwise.
			@see AIUserSuite::PutFileDialog
		*/
		ASErr PutFileDialog(const ai::UnicodeString &title, const ai::UnicodeString& defaultName, ai::FilePath& ioFilePath);

		/** Prompts for a path to an Illustrator file for input -
			shows how to apply a filter in the dialog so only Illustrator
			documents are presented.
			@param title IN the prompt dialog title.
			@param ioFilePath OUT path to the chosen file.
			@return kNoErr on success, other error otherwise.
			@see AIUserSuite::GetFileDialog
		*/
		ASErr GetIllustratorFileDialog(const ai::UnicodeString &title, ai::FilePath& ioFilePath);

		/** Prompts for a path to an Illustrator file for output -
			shows how to apply a filter in the dialog so only Illustrator
			documents are presented.
			@param title IN the prompt dialog title.
			@param defaultName IN default file name to save into.
			@param ioFilePath IN/OUT path to the chosen file.
			@return kNoErr on success, other error otherwise.
			@see AIUserSuite::PutFileDialog
		*/
		ASErr PutIllustratorFileDialog(const ai::UnicodeString &title, const ai::UnicodeString& defaultName, ai::FilePath& ioFilePath);

		/** Prompts for a path to an input text file -
			shows how to apply a filter in the dialog so only text
			files are presented.
			@param title IN the prompt dialog title.
			@param ioFilePath OUT path to the chosen file.
			@return kNoErr on success, other error otherwise.
			@see AIUserSuite::GetFileDialog
		*/
		ASErr GetTextFileDialog(const ai::UnicodeString &title, ai::FilePath& ioFilePath);	

		/** Prompts for a path to an output text file -
			shows how to apply a filter in the dialog so only text
			files are presented.
			@param title IN the prompt dialog title.
			@param defaultName IN default file name to save into.
			@param ioFilePath IN/OUT path to the chosen file.
			@return kNoErr on success, other error otherwise.
			@see AIUserSuite::PutFileDialog
		*/
		ASErr PutTextFileDialog(const ai::UnicodeString &title, const ai::UnicodeString& defaultName, ai::FilePath& ioFilePath);
};

#endif // SnpFileHelper.h

// End SnpFileHelper.h

