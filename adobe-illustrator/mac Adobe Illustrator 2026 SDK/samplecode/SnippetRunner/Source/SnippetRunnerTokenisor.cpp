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
#include "IAIStringFormatUtils.h"

#include "SnippetRunnerTokenisor.h"


/*
*/
SnippetRunnerTokenisor::SnippetRunnerTokenisor(ASUnicode separator, ASUnicode quote) : fSeparator(separator), fQuote(quote)
{
}

/*
*/
SnippetRunnerTokenisor::~SnippetRunnerTokenisor()
{
}


/*
*/
void SnippetRunnerTokenisor::SetParameters(const ai::UnicodeString& params)
{
	this->fParameters = params;
	this->Tokenise();
	this->fTokenIndex = -1;
}

/*
*/
ASBoolean SnippetRunnerTokenisor::HasTokens() const
{
	return fTokens.size() > 0;
}

/*
*/
void SnippetRunnerTokenisor::Tokenise()
{
	// Parse the given string into a vector of tokens using
	// the separator identifed at construction time to trigger
	// when to delimit tokens.
	this->fTokens.clear();
	ai::UnicodeString token;
	ASBoolean inQuote = false;
	for (ASUInt32 i = 0; i < this->fParameters.size(); i++) {
		ai::UnicodeString::UTF32TextChar c = this->fParameters [i];
		if (!inQuote) {
			if (c == this->fQuote) {
				// Start of quoted token.
				inQuote = true;
			}
			else if (c == this->fSeparator) {
				// Buffer the token.
				this->fTokens.push_back(token);
				token.clear();
			}
			else if (!this->IsWhiteSpace(c)) {
				// Add character to token.
				token.push_back(c);
			}
		} 
		else {
			if (c == this->fQuote) {
				inQuote = false;
			}
			else {
				// Add character to token.
				token.push_back(c);
			}
		}
	}

	if (token.size() > 0) {
		this->fTokens.push_back(token);
	}
}

/*
*/
ASBoolean SnippetRunnerTokenisor::IsWhiteSpace(ai::UnicodeString::UTF32TextChar c)
{
	return c == kSnippetRunnerTokenisorSpace || c == kSnippetRunnerTokenisorTab ? true : false;
}

/*
*/
ASBoolean SnippetRunnerTokenisor::GetString(ai::UnicodeString& value) const
{
	ASBoolean result = false;
	const_cast<SnippetRunnerTokenisor*>(this)->fTokenIndex++;
	if (this->fTokenIndex < this->fTokens.size()) {
		value = fTokens[fTokenIndex];
		result = true;
	}
	return result;
}

/*
*/
ASBoolean SnippetRunnerTokenisor::GetInt(ASInt32& value, const ASInt32& lowerLimit, const ASInt32& upperLimit) const
{
	ASBoolean result = false;
	ai::UnicodeString token;
	if (this->GetString(token)) {
		// Convert to integer
		ASInt32 tmp = 0;
		if (ai::NumberFormat().parseString(token, tmp)) {
			if (lowerLimit == 0 && upperLimit == 0)  {
				value = tmp;
				result = true;
			}
			else if (tmp >= lowerLimit && tmp <= upperLimit) {
				value = tmp;
				result = true;
			}
		}
	}
	return result;
}

/*
*/
ASBoolean SnippetRunnerTokenisor::GetReal(AIReal& value, const AIReal& lowerLimit, const AIReal& upperLimit) const
{
	ASBoolean result = false;
	ai::UnicodeString token;
	if (this->GetString(token)) {
		// Convert to real
		AIReal tmp = 0.0;
		if (ai::NumberFormat().parseString(token, tmp)) {
			if (lowerLimit == 0.0 && upperLimit == 0.0)  {
				value = tmp;
				result = true;
			}
			else if (tmp >= lowerLimit && tmp <= upperLimit) {
				value = tmp;
				result = true;
			}
		}
	}
	return result;
}

/*
*/
ASBoolean SnippetRunnerTokenisor::GetChoice(ai::int32& value, const std::vector<ai::UnicodeString>& choices) const
{
	ASBoolean result = false;
	ai::UnicodeString token;
	if (this->GetString(token)) {
		// Lookup token in vector
		std::vector<ai::UnicodeString>::const_iterator theChoice = std::find(choices.begin(), choices.end(), token);
		if (theChoice != choices.end()) {
			value = theChoice - choices.begin();
			result = true;
		}
	}
	return result;
}

// End SnippetRunnerTokenisor.cpp
