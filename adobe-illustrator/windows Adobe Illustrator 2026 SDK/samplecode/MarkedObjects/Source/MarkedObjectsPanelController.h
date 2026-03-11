//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/MarkedObjects/Source/MarkedObjectsPanelController.h $
//
//  $Revision: #1 $
//
//  Copyright 2010 Adobe Systems Incorporated. All rights reserved.
//  
//  NOTICE:  Adobe permits you to use, modify, and distribute this file in accordance 
//  with the terms of the Adobe license agreement accompanying it.  If you have received
//  this file from a source other than Adobe, then your use, modification, or 
//  distribution of it requires the prior written permission of Adobe.
//  
//========================================================================================
#ifndef __MARKEDOBJECTSPANELCONTROLLER_H__
#define __MARKEDOBJECTSPANELCONTROLLER_H__

#include "SDKPlugPlug.h"
#include "IllustratorSDK.h"
#include "FlashUIController.h"
#include "MarkedObjectsPlugin.hpp"

extern MarkedObjectsPlugin *gPlugin;

class MarkedObjectsPanelController : public FlashUIController
{
public:

	MarkedObjectsPanelController();

	csxs::event::EventErrorCode RegisterCSXSEventListeners();

	csxs::event::EventErrorCode RemoveEventListeners();
	
	AIErr Initialise(AIFilterMessage* pb);

	ASErr SendData();

	ASErr SendTheme();

	std::string AIRGBColorToXML(AIUIThemeColor color);

	void ParseData(const char* eventData);

	void ParsePreferenceData(const char* eventData);

	void ClearMarkedObjectList(void);

	void AddMarkedObject( MarkedObject * markedObject );
	
	void UpdateMarkedObjectSelected(void);
	
	void Clear( void );

	void UpdateBoundingBoxInfo( void );

	void UpdatePanelFontNameList();

	MarkedObject * GetMarkedObjectbyIndex(int index);

	void SetMarkedObjectManager(MarkedObjectManager *markedObjectManager)
	{
		fMarkedObjectManager = markedObjectManager;
	}

	void NewGlobalObjectSelected(void);

	void UpdateMarkedObject(void);

	bool GetHideButtonStatus()
	{
		return hideButtonStatus;
	}

	void SetHideButtonStatus(bool hideStatus)
	{
		hideButtonStatus = hideStatus;
	}

	ai::FilePath GetFilePath()
	{
		return ioFilePath;
	}

	void SetFilePath(ai::FilePath filePath)
	{
		ioFilePath = filePath;
	}

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

	AIErr IsPrimaryStageVisible(AIBoolean& isVisible);

private:
	AIRealPoint		subLoc;
	string			subFontName;
	float			subFontSize;
	AIRealPoint		labelLoc;
	string			labelFontName;
	float			labelFontSize;
	string			labelDefaultText;
	ai::int32		precision;
	AIBoolean		autoSort;
	
	ASInt32			objectIndex;
	AIRect			bounds;
	MarkedObjectManager *fMarkedObjectManager;
	
	float			locX;
	float			locY;
	string			objectID;
	string			objectLabel;
	bool			hideButtonStatus;

	ai::FilePath	ioFilePath;

private:
	void UpdatePanelByMarkedObject(MarkedObject *markedObject, int index);

public:
	static ai::UnicodeString GetElement(ai::UnicodeString srcData, const char* value);

	static ai::UnicodeString UnescapeXML(ai::UnicodeString srcValue);
};
#endif
