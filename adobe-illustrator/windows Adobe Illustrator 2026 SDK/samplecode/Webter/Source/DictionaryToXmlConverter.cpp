//========================================================================================
//  
//  $File: //ai/mainline/devtech/sandbox/lomorgan/public/samplecode/Webter/Source/DictionaryToXmlConverter.cpp $
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

#include "DictionaryToXmlConverter.h"
#include "WebterSuites.h"

namespace
{
	std::string ToCommaSeparatedString(const ai::Dictionary& dict)
	{
		stringstream entryNameList;

		size_t entryCount(0);
		for (auto entry : dict)
		{
			std::string keyString = sAIDictionary->GetKeyString(entry.first);
			if (keyString[0] == '-')
				keyString[0] = ' ';

			entryNameList << "," << keyString << "#id#=" << entryCount;
			++entryCount;
		}

		printf("Number of entries is: %d", entryCount);

		return entryNameList.str();
	}
}
//--------------------------------------------------------------------------- DictionaryToXmlConverter
DictionaryToXmlConverter::DictionaryToXmlConverter(ai::Dictionary dictionary)
: dict(dictionary)
{
}

std::string DictionaryToXmlConverter::GetPayload()
{
	std::stringstream xml;
	
	xml << "<payload>"
	<< "<dict>"
	<< GetDictLabel()
	<< "</dict>";
	
	if (dict.get() == nullptr)
	{
		xml << "<null>true</null>";
	}
	else
	{
		xml << "<null>false</null>"
		<< "<entries>"
		<< ToCommaSeparatedString(GetDictionary())
		<< "</entries>"
		<< GetAdditionalTags();
	}
	
	xml	<< "</payload>";
	
	return xml.str();
}

//--------------------------------------------------------------------------- DocDictionaryToXmlConverter
DocDictionaryToXmlConverter::DocDictionaryToXmlConverter(ai::Dictionary dict, std::string dictName)
: DictionaryToXmlConverter(dict),
fDictName(std::move(dictName))
{}

std::string DocDictionaryToXmlConverter::GetAdditionalTags()
{
	std::stringstream xml;
	xml << "<dictname>"
		<< fDictName
		<< "</dictname>";
	return xml.str();
}

//--------------------------------------------------------------------------- ArtDictionaryToXmlConverter
ArtDictionaryToXmlConverter::ArtDictionaryToXmlConverter(ai::Dictionary dict, std::string artName)
: DictionaryToXmlConverter(dict),
fArtName(std::move(artName))
{}

std::string ArtDictionaryToXmlConverter::GetAdditionalTags()
{
	std::stringstream xml;
	xml << "<artname>"
	<< fArtName
	<< "</artname>";
	
	return xml.str();
}
