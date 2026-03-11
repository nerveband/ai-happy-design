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

#ifndef _SnippetRunnerParameter_
#define _SnippetRunnerParameter_

#include "SnippetRunnerTokenisor.h"

/**	Provides input parameters to snippets (implemented as Gamma singleton)

	To get a string parameter:
	<code>
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		std::string param = parameter->GetString("Enter parameter", "default");
	</code>

	Methods are provided that allow string, integer, real and choice parameters to
	be obtained.

	Supports three modes of parameter input:
		/li default: use default parameter values supplied by the snippet
		/li prompt: prompt the user interactively for parameter values
		/li csv: get parameter values from the CSV string set by SnippetRunnerParameter::SetParameters

	@ingroup snippetrunner
*/
class SnippetRunnerParameter {
public:
	/** Returns the global instance of this singleton.
		@return the global instance of this singleton class.
	*/
	static SnippetRunnerParameter* Instance();

	/**	Deletes the singleton instance of this class.
	*/
	static void DeleteInstance();
	
	/**	Returns an ai::UnicodeString.
		<code>
		ai::UnicodeString swatchName = SnippetRunnerParameter::Instance()->GetString(ai::UnicodeString("Swatch name"), ai::UnicodeString("Black"));
		</code>

		@param prompt IN describing what the string represents
		@param defaultValue IN (default "")
		@return a string value
	*/
	virtual ai::UnicodeString	GetString(const ai::UnicodeString& prompt, const ai::UnicodeString& defaultValue = ai::UnicodeString("")) const;

	/**	Returns a std::string.
		<code>
		std::string swatchName = SnippetRunnerParameter::Instance()->GetString("Swatch name", "Black");
		</code>

		@param prompt IN describing what the string represents
		@param defaultValue IN (default "")
		@return a string value
	*/
	virtual std::string	GetString(const std::string& prompt, const std::string& defaultValue = std::string("")) const;

	/**	Returns a real.
		<code>
		AIReal strokeWeight = SnippetRunnerParameter::Instance()->GetReal(ai::UnicodeString("Stroke weight"), 1.0);
		</code>

		@param prompt IN describing what the real number represents
		@param defaultValue IN (default 0.0)
		@param lowerLimit IN (default 0.0, if both lowerLimit and upperlimit are 0.0 range checking is switched off)
		@param upperLimit IN (default 0.0)
		@return a real value
	*/
	virtual AIReal GetReal(const ai::UnicodeString& prompt, const AIReal defaultValue=0.0, const AIReal lowerLimit=0.0, const AIReal upperLimit=0.0) const;

	/**	Returns a real.
		<code>
		AIReal strokeWeight = SnippetRunnerParameter::Instance()->GetReal("Stroke weight", 1.0);
		</code>

		@param prompt IN describing what the PMReal represents
		@param defaultValue IN (default 0.0)
		@param lowerLimit IN (default 0.0, if both lowerLimit and upperlimit are 0.0 range checking is switched off)
		@param upperLimit IN (default 0.0)
		@return a real value
	*/
	virtual AIReal GetReal(const std::string& prompt, const AIReal defaultValue=0.0, const AIReal lowerLimit=0.0, const AIReal upperLimit=0.0) const;

	/**	Returns an integer.
		<code>
		ASInt32 frameCount = SnippetRunnerParameter::Instance()->GetInt32(ai::UnicodeString("Number of frames"), 1);
		</code>

		@param prompt IN describing what the integer represents
		@param defaultValue IN (default 0)
		@param lowerLimit IN (default 0, if both lowerLimit and upperlimit are 0 range checking is switched off)
		@param upperLimit IN (default 0)
		@return an integer value
	*/
	virtual ASInt32 GetInt(const ai::UnicodeString& prompt, const ASInt32 defaultValue=0, const ASInt32 lowerLimit=0, const ASInt32 upperLimit=0) const;

	/**	Returns an integer.
		<code>
		ASInt32 frameCount = SnippetRunnerParameter::Instance()->GetInt32("Number of frames", 1);
		</code>

		@param prompt IN describing what the integer represents
		@param defaultValue IN (default 0)
		@param lowerLimit IN (default 0, if both lowerLimit and upperlimit are 0 range checking is switched off)
		@param upperLimit IN (default 0)
		@return an integer value
	*/
	virtual ASInt32 GetInt(const std::string& prompt, const ASInt32 defaultValue=0, const ASInt32 lowerLimit=0, const ASInt32 upperLimit=0) const;

	/**	Returns a boolean.
		<code>
		ASBoolean isEnabled = SnippetRunnerParameter::Instance()->GetBoolean(ai::UnicodeString("Enable interaction"));
		</code>

		@param prompt IN describing what the boolean represents
		@return a boolean value
	*/
	virtual ASBoolean GetBoolean(const ai::UnicodeString& prompt) const;

		/**	Returns a boolean.

		@param prompt IN describing what the boolean represents
		@return a boolean value
	*/
	virtual ASBoolean GetBoolean(const std::string& prompt) const;

	/**	Return the index of a selected choice from a list of choices described 
		by a vector of strings. The code below requests a choice between 
		"DoThis" and "DoThat" to be made.
		<code>
		std::vector<ai::UnicodeString> choices;
		choices.push_back(ai::UnicodeString("DoThis"));
		choices.push_back(ai::UnicodeString("DoThat"));
		enum {kDoThis, kDoThat};
		ai::int32 choice  = SnippetRunnerParameter::Instance()->GetChoice("What do you want to do.", choices);
		switch (choice) 
		{
			case 0:
				break;
			case 1:
				break;
		}
		</code>

		@param prompt IN describing what the choice represents
		@param choices IN vector of strings representing the choices available
		@param defaultChoiceIndex IN (default 0)
		@return the index of the selected choice
	*/
	virtual ai::int32 GetChoice(const ai::UnicodeString& prompt, const std::vector<ai::UnicodeString>& choices, const ai::int32 defaultChoiceIndex=0) const;

	/**	Return the index of a selected choice from a list of choices described 
		by a vector of strings. The code below requests a choice between 
		"DoThis" and "DoThat" to be made.
		<code>
		std::vector<std::string> choices;
		choices.push_back(std::string("DoThis"));
		choices.push_back(std::string("DoThat"));
		enum {kDoThis, kDoThat};
		ASUInt32 choice  = SnippetRunnerParameter::Instance()->GetChoice("What do you want to do.", choices);
		switch (choice) 
		{
			case 0:
				break;
			case 1:
				break;
		}
		</code>

		@param prompt IN describing what the choice represents
		@param choices IN vector of strings representing the choices available
		@param defaultChoiceIndex IN (default 0)
		@return the index of the selected choice
	*/
	virtual ASUInt32 GetChoice(const std::string& prompt, const std::vector<std::string>& choices, const ASUInt32 defaultChoiceIndex=0) const;

	/**	Returns path to a file to be used for input.
		<code>
		ai::FilePath path = SnippetRunnerParameter::Instance()->GetFilePath(ai::UnicodeString("Choose input file"));
		</code>

		@param prompt IN describing the path required
		@param isInput true (default) for an input file (pops a GetFileDialog), false for an output file (pops a PutFileDialog)
		@param defaultValue IN (default empty path)
		@return path to chosen file
	*/
	virtual ai::FilePath GetFilePath(const ai::UnicodeString& prompt, ASBoolean isInput = true, const ai::FilePath defaultValue=ai::FilePath()) const;

	/**	Returns path to a file to be used for input.

		@param prompt IN describing the path required
		@param isInput true (default) for an input file (pops a GetFileDialog), false for an output file (pops a PutFileDialog)
		@param defaultValue IN (default empty path)
		@return path to chosen file
	*/
	virtual ai::FilePath GetFilePath(const std::string& prompt, ASBoolean isInput = true, const ai::FilePath defaultValue=ai::FilePath()) const;

	/**	Returns path to a folder to be used for input.
		<code>
		ai::FilePath path = SnippetRunnerParameter::Instance()->GetFolderPath(ai::UnicodeString("Choose input folder"));
		</code>

		@param prompt IN describing the path required
		@param defaultValue IN (default empty path)
		@return path to chosen folder
	*/
	virtual ai::FilePath GetFolderPath(const ai::UnicodeString& prompt, const ai::FilePath defaultValue=ai::FilePath()) const;

	/**	Returns path to a folder to be used for input.

		@param prompt IN describing the path required
		@param defaultValue IN (default empty path)
		@return path to chosen folder
	*/
	virtual ai::FilePath GetFolderPath(const std::string& prompt, const ai::FilePath defaultValue=ai::FilePath()) const;

	/** Pops an interactive alert with given message if mode is kPromptParameterInputMode,
		otherwise writes message to SnippetRunnerLog.
	*/
	virtual void Alert(const ai::UnicodeString& prompt) const;

	/** Pops an interactive alert with given message if mode is kPromptParameterInputMode,
		otherwise writes message to SnippetRunnerLog.
	*/
	virtual void Alert(const std::string& prompt) const;

	/** Set the parameters the snippet should run with.
		@param params a CSV string containing parameter values or empty string.
	*/
	virtual void		SetParameters(const ai::UnicodeString& params);

	/** Get the parameters the snippet was run with.
		@return a CSV string containing parameter values.
	 */
	virtual ai::UnicodeString		GetParameters() const;

	/** Indicates parameter input mode to be used.
	*/
	typedef enum {
		/** Use default parameter values supplied by the snippet. */
		kDefaultParameterInputMode, 
		/** Prompt the user for parameter values required by the snippet. */
		kPromptParameterInputMode, 
		/** Use the CSV string supplied by SnippetRunnerParameter::SetParameters to source parameter values required by the snippet. */
		kCSVParameterInputMode,
		kAssetFolderInputMode
	} ParameterInputMode;

	/** Sets the input mode.
	*/
	virtual void SetParameterInputMode(ParameterInputMode mode);

	/** Gets the input mode.
	*/
	virtual ParameterInputMode GetParameterInputMode() const;

	/** Stack based helper class that sets the parameter input mode on construction and restores previous mode on destruction.
		@ingroup snippetrunner
	*/
	class ModePreserver {
	public:
		ModePreserver(ParameterInputMode mode) {fMode = SnippetRunnerParameter::Instance()->GetParameterInputMode();SnippetRunnerParameter::Instance()->SetParameterInputMode(mode);}
		~ModePreserver() {SnippetRunnerParameter::Instance()->SetParameterInputMode(fMode);}
	private:
		SnippetRunnerParameter::ParameterInputMode fMode;
	};

private:
	/** Constructor is hidden - instantiate using SnippetRunnerLog::Instance.
	*/
	SnippetRunnerParameter();

	/** Destructor is hidden - destroy using SnippetRunnerLog::DeleteInstance.
	*/
	~SnippetRunnerParameter();

	/** Copy constructor is hidden.
		@param rhs
	*/
	SnippetRunnerParameter(SnippetRunnerParameter const&);

	/** Assignment operator is hidden.
		@param IN SnippetRunnerParameter object to assign to lhs.
		@return copy of SnippetRunnerParameter passed in.
	*/
	SnippetRunnerParameter& operator=(SnippetRunnerParameter const&);

private:
	/** Points to the singleton instance of this class.
	*/
	static SnippetRunnerParameter* fInstance;

	/** Tokenises CSV parameters.
	*/
	SnippetRunnerTokenisor fTokenisor;
};

#endif //_SnippetRunnerParameter_

// End SnippetRunnerParameter.h
