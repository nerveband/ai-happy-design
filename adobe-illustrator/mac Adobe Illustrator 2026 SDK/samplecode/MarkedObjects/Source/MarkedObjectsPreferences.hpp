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

#ifndef __MarkedObjectsPreferences_hpp__
#define __MarkedObjectsPreferences_hpp__

#include "MarkedObjectsPlugin.hpp"

class MarkedObjectsPanelController;

class MarkedObjectsPreferences
{
private:
	MarkedObjectsPanelController* 		moController;
	
	AIRealPoint					subLoc;
	string						subFontName;
	AIReal						subFontSize;
	
	AIRealPoint					labelLoc;
	string						labelFontName;
	AIReal						labelFontSize;
	string						labelDefaultText;
	ai::int32					precision;
	AIBoolean					autoSort;
		
public:
	MarkedObjectsPreferences();
	~MarkedObjectsPreferences();
	
	int DoModalPrefs( void );
	
	AIRealPoint	GetSubStringLocation( void ) { return subLoc; }
	const string GetSubStringFontName( void )	{ return subFontName; }
	AIReal GetSubStringFontSize( void )	{ return subFontSize; }

	AIRealPoint	GetLabelStringLocation( void ) { return labelLoc; }
	const string GetLabelStringFontName( void ) { return labelFontName; }
	AIReal GetLabelStringFontSize( void ) { return labelFontSize; }
	const string GetLabelStringDefaultText( void ) { return labelDefaultText; }
	
	void SetSubStringLocation( AIRealPoint inPoint ) { subLoc = inPoint; }
	void SetSubStringFontName( const string fontName ) { subFontName = fontName; }
	void SetSubStringFontSize( AIReal fontSize ) { subFontSize = fontSize; }

	void SetLabelStringLocation( AIRealPoint inPoint ) { labelLoc = inPoint; }
	void SetLabelStringFontName( const string fontName ) { labelFontName = fontName; }
	void SetLabelStringFontSize( AIReal fontSize ) { labelFontSize = fontSize; }
	void SetLabelStringDefaultText( const string newText ) { labelDefaultText = newText; }	

	void SetPrecision( ai::int32 p ) { precision = p; }
	ai::int32 GetPrecision( void ) { return precision; }

	void SetAutoSort( AIBoolean as ) { autoSort = as; }
	AIBoolean GetAutoSort( void ) { return autoSort; }

	void SetMarkedObjectsPanelController(MarkedObjectsPanelController *markedObjectsPanelController)
	{
		moController = markedObjectsPanelController;
	}
};

extern MarkedObjectsPreferences* gPreferences;

#endif
// end MarkedObjectsPreferences.hpp
