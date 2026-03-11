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

#ifndef _SnippetRunnerPreferences_
#define _SnippetRunnerPreferences_

#include "SnippetRunnerParameter.h"

/** Provides access to this plug-in's preferences.

	@ingroup snippetrunner
*/
class SnippetRunnerPreferences {

public:
	/** Returns the global instance of this singleton.
		@return the global instance of this singleton class.
	*/
	static SnippetRunnerPreferences* Instance();

	/**	Deletes the singleton instance of this class.
	*/
	static void DeleteInstance();

	/** Gets the path to the assets folder 
		(the user is prompted to locate the folder if the preference is not set).
		@return the path to the assets folder.
	*/
	ai::FilePath GetAssetsFolderPath();

	/** Sets the path to the assets folder.
		@param assetsFolder IN the path to the assets folder.
	*/
	void SetAssetsFolderPath(const ai::FilePath& assetsFolder) {fAssetsFolderPath = assetsFolder;}

	/** Prompts for and sets the assets folder path.
	*/
	void PromptAssetsFolderPath();

	/** Sets the input mode.
		@param mode IN the parameter input mode.
	*/
	virtual void SetParameterInputMode(SnippetRunnerParameter::ParameterInputMode mode) {fParameterInputMode = mode;}

	/** Gets the input mode.
		@return the parameter input mode.
	*/
	virtual SnippetRunnerParameter::ParameterInputMode GetParameterInputMode() const {return static_cast<SnippetRunnerParameter::ParameterInputMode>(fParameterInputMode);}

	/** Sets the unit test flag.
		@param unitTest IN unit test flag value.
	*/
	virtual void SetUnitTest(ASBoolean unitTest) {fUnitTest = unitTest;}

	/** Gets the unit test flag.
		@return true if unit test flag set, false if not.
	*/
	virtual ASBoolean GetUnitTest() const {return fUnitTest;}


private:
	/** Constructor is hidden - instantiate using SnippetRunnerPreferences::Instance.
	*/
	SnippetRunnerPreferences();

	/** Destructor is hidden - destroy using SnippetRunnerPreferences::DeleteInstance.
	*/
	~SnippetRunnerPreferences();

	/** Copy constructor is hidden.
		@param rhs
	*/
	SnippetRunnerPreferences(SnippetRunnerPreferences const&);

	/** Assignment operator is hidden.
		@param IN SnippetRunnerPreferences object to assign to lhs.
		@return copy of SnippetRunnerPreferences passed in.
	*/
	SnippetRunnerPreferences& operator=(SnippetRunnerPreferences const&);  // assign op is hidden

private:
	/** Points to the singleton instance of this class.
	*/
	static SnippetRunnerPreferences* fInstance;

	/**	Contains the path to the assets folder.
	*/
	ai::FilePath fAssetsFolderPath;

	/** The parameter input mode - default, prompt, CSV.
	*/
	ai::int32 fParameterInputMode;

	/**	Unit test flag.
	*/
	ASBoolean fUnitTest;
};

#endif //_SnippetRunnerPreferences_

// End SnippetRunnerPreferences.h
