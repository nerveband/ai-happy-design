//========================================================================================
//  
//  $File: //ai/mainline/devtech/sandbox/lomorgan/public/samplecode/Webter/Source/DictionaryToXmlConverter.h $
//
//  $Revision: #1 $
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
#include "IAIDictionary.hpp"

#ifndef __DICTIONARYTOXMLCONVERTER_H__
#define __DICTIONARYTOXMLCONVERTER_H__

/** An abstract class for converting a dictionary to xml
	Class follows the strategy pattern - subclasses should
	implement GetDictLabel() to set the contents of the dict
	tag in the produced xml, and can supplement the xml payload
	with custom xml by overriding GetAdditionalTags(), which by
	default returns an empty string.
*/
class DictionaryToXmlConverter
{
public:
	/** Constructor.
		@param dictionary IN the dictionary to convert to XML
	*/
	DictionaryToXmlConverter(ai::Dictionary dictionary);
	
	/** Get the xml payload
		@returns XML payload of the following form:
		<payload>
			<dict>return value of GetDictLabel()</dict>
			<null>true if the underlying AI dictionary is null, false o/w</null>
			<entries>
				Comma separated list of dictionary entries. Tag is not present if
				the underlying AI dictionary is null.
			</entries>
			<!-- Additional tags returned by GetAdditionalTags() go here -->
		</payload>
	*/
	virtual std::string GetPayload();

	/** Get the label of the dictionary, used by GetPayload().
		Subclasses must implement this.
		@returns a textual description of the type of dictionary.
	*/
	virtual std::string GetDictLabel() = 0;
	
	/** Get any additional tags to add to the payload.
		Subclasses may choose to implement this if they wish to supplement the XML
		produced by GetPayload()
		@returns additional tags, empty string by default
	*/
	virtual std::string GetAdditionalTags() { return ""; };
	
	/** Get the dictionary that will be converted to XML
		@return the Dictionary passed in to the constructor
	*/
	ai::Dictionary GetDictionary() { return dict; };

protected:
	ai::Dictionary dict;
};

/** A specialisation of the DictionaryToXmlConverter that produces
 the xml payload for a document dictionary
 */
class DocDictionaryToXmlConverter : public DictionaryToXmlConverter
{
public:
	DocDictionaryToXmlConverter(ai::Dictionary dict, std::string dictName);
	
	virtual std::string GetDictLabel() { return "doc"; };

	virtual std::string GetAdditionalTags();

private:
	std::string fDictName;
};

/** A specialisation of the DictionaryToXmlConverter that produces
 the xml payload for an art dictionary
 */
class ArtDictionaryToXmlConverter : public DictionaryToXmlConverter
{
public:
	ArtDictionaryToXmlConverter(ai::Dictionary dict, std::string artName);
	
	virtual std::string GetDictLabel() { return "art"; };
	
	virtual std::string GetAdditionalTags();
	
private:
	std::string fArtName;
	
};

#endif
