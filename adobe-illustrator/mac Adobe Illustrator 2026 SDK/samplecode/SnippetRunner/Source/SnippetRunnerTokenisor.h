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

#ifndef _SnippetRunnerTokenisor_
#define _SnippetRunnerTokenisor_

const ASUnicode kSnippetRunnerTokenisorSeparator = 0x002C; // Comma
const ASUnicode kSnippetRunnerTokenisorQuote = 0x0022; // QuotationMark
const ASUnicode kSnippetRunnerTokenisorSpace = 0x0020; // Space
const ASUnicode kSnippetRunnerTokenisorTab = 0x0009; // Tab

/** Tokenises a CSV string of parameters and provides helper methods 
	to convert parameters into various data types.

	@ingroup snippetrunner
*/
class SnippetRunnerTokenisor {
public:

	/** Constructor
		@param separator IN token delimiter.
		@param quote IN char representing a quote.
	*/
	SnippetRunnerTokenisor(ASUnicode separator = kSnippetRunnerTokenisorSeparator, ASUnicode quote = kSnippetRunnerTokenisorQuote);

	/** Destructor
	*/
	virtual ~SnippetRunnerTokenisor();

	/** Set the parameters.
		Tokenises params into a list of tokens using the separator specified on construction.
		Supports quotes to delimit tokens where white space is significant.

		@param params IN
	*/
	void SetParameters(const ai::UnicodeString& params);

	/** @return true if tokens are available, false otherwise.
	*/
	ASBoolean HasTokens() const;

	/** Get next token as a string.
		@param value OUT string token.
		@return true on success, false otherwise.
	*/
	ASBoolean GetString(ai::UnicodeString& value) const;

	/** Get next token as an integer.
		@param value OUT the value.
		@param lowerLimit IN (default 0, if both lowerLimit and upperlimit are 0 range checking is switched off)
		@param upperLimit IN (default 0)
		@return true on success, false otherwise.
	*/
	ASBoolean GetInt(ai::int32& value, const ai::int32& lowerLimit=0, const ai::int32& upperLimit=0) const;

	/** Get next token as a real number.
		@param value OUT the value.
		@param lowerLimit IN (default 0.0, if both lowerLimit and upperlimit are 0.0 range checking is switched off)
		@param upperLimit IN (default 0.0)
		@return true on success, false otherwise.
	*/
	ASBoolean GetReal(AIReal& value, const AIReal& lowerLimit=0.0, const AIReal& upperLimit=0.0) const;

	/** Get next token as the index of a selected choice in a list of valid choices. 
		@param value OUT the index of the selected choice
		@param choices IN vector of strings representing the choices available
		@return true on success, false otherwise.
	*/
	ASBoolean GetChoice(ai::int32& value, const std::vector<ai::UnicodeString>& choices) const;

private:
	/** Tokenises the parameters.
	*/
	void Tokenise();

	/** Checks if the character is white-space.
		@param c IN the character.
		@return true if whitespace, false otherwise.
	*/
	ASBoolean IsWhiteSpace(ai::UnicodeString::UTF32TextChar c);

	/** Parameter string passed by client.
	*/
	ai::UnicodeString fParameters;

	/** Vector of tokens parsed out of fParameters.
	*/
	std::vector<ai::UnicodeString> fTokens;

	/** Index of last token read from fTokens.
	*/
	ASUInt32 fTokenIndex;

	/** Character to be used to delimit tokens in parameter strings.
	*/
	const ASUnicode fSeparator;

	/** Character to be used to quote tokens where white space is significant.
	*/
	const ASUnicode fQuote;
};


#endif //_SnippetRunnerTokenisor_

// End SnippetRunnerTokenisor.h
