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

#include "MarkedObjectsPreferences.hpp"
#include "MarkedObjectID.h"
#include "MarkedObjectsPanelController.h"

static const char* kStringPreference = "MarkedObjects";
static const char* kMarkedObjects_SubLocationXStr = "kMarkedObjects_SubLocationXStr";
static const char* kMarkedObjects_SubLocationYStr = "kMarkedObjects_SubLocationYStr";
static const char* kMarkedObjects_SubFontNameStr = "kMarkedObjects_SubFontNameStr";
static const char* kMarkedObjects_SubFontSizeStr = "kMarkedObjects_SubFontSizeStr";

static const char* kMarkedObjects_DefLocationXStr = "kMarkedObjects_DefLocationXStr";
static const char* kMarkedObjects_DefLocationYStr = "kMarkedObjects_DefLocationYStr";
static const char* kMarkedObjects_DefFontNameStr = "kMarkedObjects_DefFontNameStr";
static const char* kMarkedObjects_DefFontSizeStr = "kMarkedObjects_DefFontSizeStr";
static const char* kMarkedObjects_DefTextStr = "kMarkedObjects_DefTextStr";

static const char* kMarkedObjects_PrecisionStr = "kMarkedObjects_PrecisionStr";
static const char* kMarkedObjects_AutoSortStr = "kMarkedObjects_AutoSortStr";

MarkedObjectsPreferences* gPreferences = NULL;

MarkedObjectsPreferences::MarkedObjectsPreferences()
{
	subLoc.h = 6;
	subLoc.v = -6;
	subFontName = "ArialMT";
	subFontSize = 7;
	
	labelLoc.h = 20;
	labelLoc.v = -4;
	labelFontName = "ArialMT";
	labelFontSize = 10;
	labelDefaultText = "Default Text";

	precision = 0;
	autoSort = true;

	gPreferences = this;

	char tempChar[1024];

	// these *()^*&^ routines don't return a value if the preference is not there, very nice
	double value = subLoc.h;

	AIErr error = sAIPreference->GetRealPreference( kStringPreference, kMarkedObjects_SubLocationXStr, &value );
	if ( error == kNoErr ) subLoc.h = (AIReal)value;

	value = subLoc.v;
	error = sAIPreference->GetRealPreference( kStringPreference, kMarkedObjects_SubLocationYStr, &value );
	if ( error == kNoErr ) subLoc.v = (AIReal)value;

	memset(tempChar, 0, 1024);
	strcpy( tempChar, subFontName.c_str() );
	error = sAIPreference->GetStringPreference( kStringPreference, kMarkedObjects_SubFontNameStr, tempChar );
	if ( error == kNoErr ) subFontName.assign( tempChar, tempChar + strlen( tempChar ) );

	value = subFontSize;
	error = sAIPreference->GetRealPreference( kStringPreference, kMarkedObjects_SubFontSizeStr, &value );
	if ( error == kNoErr ) subFontSize = (AIReal)value;

	value = labelLoc.h;
	error = sAIPreference->GetRealPreference( kStringPreference, kMarkedObjects_DefLocationXStr, &value );
	if ( error == kNoErr ) labelLoc.h = (AIReal)value;

	value = labelLoc.v;
	error = sAIPreference->GetRealPreference( kStringPreference, kMarkedObjects_DefLocationYStr, &value );
	if ( error == kNoErr ) labelLoc.v = (AIReal)value;

	memset(tempChar, 0, 1024);
	strcpy( tempChar, labelFontName.c_str() );
	error = sAIPreference->GetStringPreference( kStringPreference, kMarkedObjects_DefFontNameStr, tempChar );
	if ( error == kNoErr ) labelFontName.assign( tempChar, tempChar + strlen( tempChar ) );

	value = labelFontSize;
	error = sAIPreference->GetRealPreference( kStringPreference, kMarkedObjects_DefFontSizeStr, &value );
	if ( error == kNoErr ) labelFontSize = (AIReal)value;

	memset(tempChar, 0, 1024);
	strcpy( tempChar, labelDefaultText.c_str() );
	error = sAIPreference->GetStringPreference( kStringPreference, kMarkedObjects_DefTextStr, tempChar );
	if ( error == kNoErr ) labelDefaultText.assign( tempChar, tempChar + strlen( tempChar ) );
	
	error = sAIPreference->GetIntegerPreference( kStringPreference, kMarkedObjects_PrecisionStr, &precision );

	error = sAIPreference->GetBooleanPreference( kStringPreference, kMarkedObjects_AutoSortStr, &autoSort );

}

MarkedObjectsPreferences::~MarkedObjectsPreferences()
{
	if ( moController != NULL && sAIPreference != NULL)
	{
	
		// for the marked object, font information, offsets, and size
		sAIPreference->PutRealPreference( kStringPreference, kMarkedObjects_SubLocationXStr, subLoc.h );
		sAIPreference->PutRealPreference( kStringPreference, kMarkedObjects_SubLocationYStr, subLoc.v );
		sAIPreference->PutStringPreference( kStringPreference, kMarkedObjects_SubFontNameStr, subFontName.c_str() );
		sAIPreference->PutRealPreference( kStringPreference, kMarkedObjects_SubFontSizeStr, subFontSize );

		sAIPreference->PutRealPreference( kStringPreference, kMarkedObjects_DefLocationXStr, labelLoc.h );
		sAIPreference->PutRealPreference( kStringPreference, kMarkedObjects_DefLocationYStr, labelLoc.v );
		sAIPreference->PutStringPreference( kStringPreference, kMarkedObjects_DefFontNameStr, labelFontName.c_str() );
		sAIPreference->PutRealPreference( kStringPreference, kMarkedObjects_DefFontSizeStr, labelFontSize );
		sAIPreference->PutStringPreference( kStringPreference, kMarkedObjects_DefTextStr, labelDefaultText.c_str() );

		sAIPreference->PutIntegerPreference( kStringPreference, kMarkedObjects_PrecisionStr, precision );
		sAIPreference->PutBooleanPreference( kStringPreference, kMarkedObjects_AutoSortStr, autoSort );
	}
	gPreferences = NULL;
}

int MarkedObjectsPreferences::DoModalPrefs( void )
{

	AIFontKey resultFont;
	ASErr result = sAIFont->FindFont( const_cast<char*>(moController->GetSubStringFontName().c_str()), 
	kAIAnyFontTechnology, kUnknownAIScript, false, &resultFont);

	if ( result != kNoErr ) 
	{
		sAIUser->MessageAlert(ai::UnicodeString("SubString font name is invalid"));
		return result;
	}

	result = sAIFont->FindFont( const_cast<char*>(moController->GetLabelStringFontName().c_str()), 
	kAIAnyFontTechnology, kUnknownAIScript, false, &resultFont);
	if ( result != kNoErr ) 
	{
		sAIUser->MessageAlert(ai::UnicodeString("Label font name is invalid"));
		return result;
	}

	SetSubStringLocation( moController->GetSubStringLocation() );
	SetSubStringFontName( moController->GetSubStringFontName() );
	SetSubStringFontSize( moController->GetSubStringFontSize() );

	SetLabelStringLocation( moController->GetLabelStringLocation() );
	SetLabelStringFontName( moController->GetLabelStringFontName() );
	SetLabelStringFontSize( moController->GetLabelStringFontSize() );
	SetLabelStringDefaultText( moController->GetLabelStringDefaultText() );

	SetPrecision( moController->GetPrecision() );
	SetAutoSort( moController->GetAutoSort() );
	return result;
}


// end MarkedObjectsPreferences.cpp
