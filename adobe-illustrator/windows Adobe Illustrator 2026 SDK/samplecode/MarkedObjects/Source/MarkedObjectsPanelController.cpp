//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/MarkedObjects/Source/MarkedObjectsPanelController.cpp $
//
//  $Revision: #2 $
//
//  Copyright 2010 Adobe Systems Incorporated. All rights reserved.
//  
//  NOTICE:  Adobe permits you to use, modify, and distribute this file in accordance 
//  with the terms of the Adobe license agreement accompanying it.  If you have received
//  this file from a source other than Adobe, then your use, modification, or 
//  distribution of it requires the prior written permission of Adobe.
//  
//========================================================================================

#include "IllustratorSDK.h"
#include "MarkedObjectsPanelController.h"
#include "AppContext.hpp"
#include "MarkedObjectsSuites.h"
#include "Plugin.hpp"

#define MARKEDOBJECTS_UI_EXTENSION	"com.adobe.illustrator.markedobjectsui.panel"

#define EVENT_TYPE_UPDATE_CLICKED			"com.adobe.csxs.events.MarkedObjectsUpdateClicked"
#define EVENT_TYPE_HIDEALL_CLICKED			"com.adobe.csxs.events.MarkedObjectsHideAllClicked"
#define EVENT_TYPE_SAVE_CLICKED				"com.adobe.csxs.events.MarkedObjectsSave"
#define EVENT_TYPE_RENMUBERIDS_CLICKED		"com.adobe.csxs.events.MarkedObjectsRenumberIDs"
#define EVENT_TYPE_PREFREENCES_CLICKED		"com.adobe.csxs.events.MarkedObjectsPreferences"
#define EVENT_TYPE_ITEM_SELECTED			"com.adobe.csxs.events.MarkedObjectsItemSelect"
#define EVENT_TYPE_PANEL_READY				"com.adobe.csxs.events.MarkedObjectsPanelReady"
#define EVENT_TYPE_PREFREENCESSTACK_READY	"com.adobe.csxs.events.MarkedObjectsPreferencesStackReady"
#define EVENT_TYPE_BRIGHTNESS_CHANGE		"com.adobe.csxs.events.BrightnessChange"

static void UpdateClickedFunc(const csxs::event::Event* const eventParam, void* const context)
{
	MarkedObjectsPanelController *markedObjectsPanelController = (MarkedObjectsPanelController *)context;
	if ( NULL == markedObjectsPanelController || eventParam == NULL )
		return;
	markedObjectsPanelController->ParseData(eventParam->data);

	do {
		AppContext appContext(gPlugin->GetPluginRef());
		ai::UnicodeString undoText("Undo UpdateObjects");
		ai::UnicodeString redoText("Redo UpdateObjects");
		
		sAIUndo->SetUndoTextUS(undoText, redoText);
		markedObjectsPanelController->UpdateMarkedObject();
	} while(false);
}

static void HideAllClickedFunc(const csxs::event::Event* const eventParam, void* const context)
{
	MarkedObjectsPanelController *markedObjectsPanelController = (MarkedObjectsPanelController *)context;
	if ( NULL == markedObjectsPanelController || eventParam == NULL )
		return;

	ai::UnicodeString dataStr(eventParam->data);
	ai::UnicodeString hideStr = MarkedObjectsPanelController::GetElement(dataStr, "hide");
	if ( !hideStr.empty())
	{
		if ( hideStr.compare(ai::UnicodeString("true")) == 0 )
		{
			markedObjectsPanelController->SetHideButtonStatus(true);
			sAINotifier->Notify(kHideLayerNotifier, nullptr);
		}
		else if ( hideStr.compare(ai::UnicodeString("false")) == 0 )
		{
			markedObjectsPanelController->SetHideButtonStatus(false);
			sAINotifier->Notify(kHideLayerNotifier, nullptr);
		}
	}
}

static void SaveClickedFunc(const csxs::event::Event* const eventParam, void* const context)
{
	MarkedObjectsPanelController *markedObjectsPanelController = (MarkedObjectsPanelController *)context;
	if ( NULL == markedObjectsPanelController || eventParam == NULL )
		return;

	const ai::UnicodeString title( "Where do you want to save the Marked Objects?" );
	const AIFileDialogFilters * dlgOpts = NULL;
	const ai::UnicodeString defaultName;
	ai::FilePath ioFilePath;
	AIErr e = sAIUser->PutFileDialog( title, dlgOpts, defaultName, ioFilePath );

	AIBoolean doSave = true;
	AIBoolean doOverwrite = false;

	// Check if the file exists
	while ( (ioFilePath.Exists(true) && !doOverwrite) && doSave )
	{
		const ai::UnicodeString existsMsg("This file already exists, would you like to overwrite it?");
		AIAnswer answer = sAIUser->QuestionAlert(existsMsg);
		if ( answer == kAIAnswer_No )
		{
			// Prompt again
			e = sAIUser->PutFileDialog( title, dlgOpts, defaultName, ioFilePath );
		}
		else if (answer == kAIAnswer_Yes)
		{
			// Overwrite file
			doOverwrite = true;
		}
		else // Cancel (or error)
		{
			// Don't save the file
			doSave = false;
		}
	}

	if ( e == kNoErr && doSave )
	{
		markedObjectsPanelController->SetFilePath(ioFilePath);
		sAINotifier->Notify(kSaveObjectNotifier, nullptr);
	}
}

static void ReNumberIDClickedFunc(const csxs::event::Event* const eventParam, void* const context)
{
	sAINotifier->Notify(kResetNumberIDNotifier, nullptr);
}

static void PreferencesStackReady(const csxs::event::Event* const eventParam, void* const context)
{
	MarkedObjectsPanelController *markedObjectsPanelController = (MarkedObjectsPanelController *)context;
	markedObjectsPanelController->SendData();
}

static void ItemSelectedFunc(const csxs::event::Event* const eventParam, void* const context)
{
	MarkedObjectsPanelController *markedObjectsPanelController = (MarkedObjectsPanelController *)context;
	if ( NULL == markedObjectsPanelController || eventParam == NULL )
		return;
	markedObjectsPanelController->ParseData(eventParam->data);
	sAINotifier->Notify(kObjectSelectedNotifier, nullptr);
}

static void PreferencesClickedFunc(const csxs::event::Event* const eventParam, void* const context)
{
	MarkedObjectsPanelController *markedObjectsPanelController = (MarkedObjectsPanelController *)context;
	if ( NULL == markedObjectsPanelController || eventParam == NULL )
		return;

	markedObjectsPanelController->ParsePreferenceData(eventParam->data);
	sAINotifier->Notify(kPreferencesSetNotifier, nullptr);
}

static void PanelReadyFunc(const csxs::event::Event* const eventParam, void* const context)
{
	MarkedObjectsPanelController *markedObjectsPanelController = (MarkedObjectsPanelController *)context;
	markedObjectsPanelController->UpdatePanelFontNameList();
	markedObjectsPanelController->SendTheme();
}

MarkedObjectsPanelController::MarkedObjectsPanelController()
: FlashUIController(MARKEDOBJECTS_UI_EXTENSION),
objectID(""),
objectLabel(""),
locX(0),
locY(0),
objectIndex(0),
hideButtonStatus(false)
{}

/* Add event listeners
*/
csxs::event::EventErrorCode MarkedObjectsPanelController::RegisterCSXSEventListeners()
{
	csxs::event::EventErrorCode result = csxs::event::kEventErrorCode_Success;
	do {
		result = fPPLib.AddEventListener(EVENT_TYPE_UPDATE_CLICKED, UpdateClickedFunc, this);
		if ( result != csxs::event::kEventErrorCode_Success ) 
		{
			break;
		}

		result = fPPLib.AddEventListener(EVENT_TYPE_HIDEALL_CLICKED, HideAllClickedFunc, this);
		if ( result != csxs::event::kEventErrorCode_Success ) 
		{
			break;
		}

		result = fPPLib.AddEventListener(EVENT_TYPE_SAVE_CLICKED, SaveClickedFunc, this);
		if ( result != csxs::event::kEventErrorCode_Success ) 
		{
			break;
		}

		result =  fPPLib.AddEventListener(EVENT_TYPE_RENMUBERIDS_CLICKED, ReNumberIDClickedFunc, this);
		if ( result != csxs::event::kEventErrorCode_Success ) 
		{
			break;
		}

		result = fPPLib.AddEventListener(EVENT_TYPE_PREFREENCES_CLICKED, PreferencesClickedFunc, this);
		if ( result != csxs::event::kEventErrorCode_Success ) 
		{
			break;
		}

		result = fPPLib.AddEventListener(EVENT_TYPE_PANEL_READY, PanelReadyFunc, this);
		if ( result != csxs::event::kEventErrorCode_Success ) {
			break;
		}

		result = fPPLib.AddEventListener(EVENT_TYPE_ITEM_SELECTED, ItemSelectedFunc, this);
		if ( result != csxs::event::kEventErrorCode_Success ) {
			break;
		}

		result = fPPLib.AddEventListener(EVENT_TYPE_PREFREENCESSTACK_READY, PreferencesStackReady, this);
		if ( result != csxs::event::kEventErrorCode_Success ) {
			break;
		}	
	}
	while(false);
	return result;
}

/* Remove event listeners
*/
csxs::event::EventErrorCode MarkedObjectsPanelController::RemoveEventListeners()
{
	csxs::event::EventErrorCode result = csxs::event::kEventErrorCode_Success;
	do {
		result = fPPLib.RemoveEventListener(EVENT_TYPE_UPDATE_CLICKED, UpdateClickedFunc, this);
		if ( result != csxs::event::kEventErrorCode_Success ) 
		{		
			break;
		}

		result = fPPLib.RemoveEventListener(EVENT_TYPE_HIDEALL_CLICKED, HideAllClickedFunc, this);
		if ( result != csxs::event::kEventErrorCode_Success ) 
		{		
			break;
		}

		result = fPPLib.RemoveEventListener(EVENT_TYPE_SAVE_CLICKED, SaveClickedFunc, this);
		if ( result != csxs::event::kEventErrorCode_Success ) 
		{		
			break;
		}

		result = fPPLib.RemoveEventListener(EVENT_TYPE_RENMUBERIDS_CLICKED, ReNumberIDClickedFunc, this);
		if ( result != csxs::event::kEventErrorCode_Success ) 
		{		
			break;
		}

		result = fPPLib.RemoveEventListener(EVENT_TYPE_PREFREENCES_CLICKED, PreferencesClickedFunc, this);
		if ( result != csxs::event::kEventErrorCode_Success ) 
		{		
			break;
		}
		result = fPPLib.RemoveEventListener(EVENT_TYPE_PANEL_READY, PanelReadyFunc, this);
		if ( result != csxs::event::kEventErrorCode_Success ) 
		{
			break;
		}

		result = fPPLib.RemoveEventListener(EVENT_TYPE_ITEM_SELECTED, ItemSelectedFunc, this);
		if ( result != csxs::event::kEventErrorCode_Success ) {
			break;
		}

		result = fPPLib.RemoveEventListener(EVENT_TYPE_PREFREENCESSTACK_READY, PreferencesStackReady, this);
		if ( result != csxs::event::kEventErrorCode_Success ) {
			break;
		}
	}
	while(false);
	return result;
}

/* SendDataToPanel
*/
ASErr MarkedObjectsPanelController::SendData() 
{
	ASErr error = kNoErr;
	
	this->SetSubStringLocation( gPreferences->GetSubStringLocation() );
	this->SetSubStringFontName( gPreferences->GetSubStringFontName() );
	this->SetSubStringFontSize( gPreferences->GetSubStringFontSize() );

	this->SetLabelStringLocation( gPreferences->GetLabelStringLocation() );
	this->SetLabelStringFontName( gPreferences->GetLabelStringFontName() );
	this->SetLabelStringFontSize( gPreferences->GetLabelStringFontSize() );
	this->SetLabelStringDefaultText( gPreferences->GetLabelStringDefaultText() );

	this->SetPrecision( gPreferences->GetPrecision() );
	this->SetAutoSort( gPreferences->GetAutoSort() );

	stringstream myString;
	myString << "<payload>"  << "<subfontsize>" << 
		subFontSize << "</subfontsize><subfontname>" << subFontName	<< "</subfontname><sublocx>" << 
		subLoc.h << "</sublocx>" << "<sublocy>"  << subLoc.v
		<< "</sublocy>" << "<labelfontsize>" << labelFontSize << "</labelfontsize>" << 
		"<labelfontname>" << labelFontName << "</labelfontname><labellocx>" << 
		labelLoc.h << "</labellocx>" << "<labellocy>"  << labelLoc.v
		<< "</labellocy>" << "<defaulttext>" << labelDefaultText << "</defaulttext>" << "<precision>" << precision << "</precision>" <<
    "<autosort>" << (autoSort == true ? "1" : "0") << "</autosort>" << "</payload>";

	std::string xmlString = myString.str();
	// Create and dispatch event
	csxs::event::Event event = {"com.adobe.csxs.events.MarkedObjectsUpdatePreferencesPanel",
								 csxs::event::kEventScope_Application,
								 "MarkedObjectsUpdatePreferencesPanel",
								 NULL,
								 xmlString.c_str()};
	fPPLib.DispatchEvent(&event);
	return error;
}

ASErr MarkedObjectsPanelController::SendTheme()
{
	AIErr error = kNoErr;
	
	AIUIThemeColor backGroundcolor;	
	AIUIThemeColor textcolor;
	AIUIThemeColor focusColor;
	AIUIThemeColor inputBackgroundColor;
	AIUIThemeColor inputTextColor;
	error = sAIUITheme->GetUIThemeColor(kAIUIThemeSelectorPanel,kAIUIComponentColorBackground,backGroundcolor);
	error = sAIUITheme->GetUIThemeColor(kAIUIThemeSelectorPanel,kAIUIComponentColorText,textcolor);
	error = sAIUITheme->GetUIThemeColor(kAIUIThemeSelectorPanel,kAIUIComponentColorFocusRing,focusColor);
	error = sAIUITheme->GetUIThemeColor(kAIUIThemeSelectorPanel,kAIUIComponentColorEditText,inputTextColor);
	error = sAIUITheme->GetUIThemeColor(kAIUIThemeSelectorPanel,kAIUIComponentColorEditTextBackground,inputBackgroundColor);

	stringstream myString;
	myString << "<payload><theme>" 
		<< "<backgroundColor>" 
			<< AIRGBColorToXML(backGroundcolor)
		<< "</backgroundColor>"
		<< "<textColor>" 
			<< AIRGBColorToXML(textcolor)
		<< "</textColor>"
		<< "<focusColor>" 
			<< AIRGBColorToXML(focusColor)
		<< "</focusColor>" 
		<< "<inputBackgroundColor>" 
			<< AIRGBColorToXML(inputBackgroundColor)
		<< "</inputBackgroundColor>" 
		<< "<inputTextColor>" 
			<< AIRGBColorToXML(inputTextColor)
		<< "</inputTextColor>" 
		<< "</theme></payload>";
	std::string xmlString = myString.str();
	csxs::event::Event event = {EVENT_TYPE_BRIGHTNESS_CHANGE,
								 csxs::event::kEventScope_Application,
								 "MarkedObjectsUpdatePreferencesPanel",
								 NULL,
								 xmlString.c_str()};
	fPPLib.DispatchEvent(&event);
	return error;
}

std::string MarkedObjectsPanelController::AIRGBColorToXML(AIUIThemeColor color)
{
	int colorDepth = 255;
	int red = (int) (color.red * colorDepth);
	int green = (int) (color.green * colorDepth);
	int blue = (int) (color.blue * colorDepth);
	stringstream myString;
	myString 
		<< "<color>"
			<< "<alpha>" << color.alpha << "</alpha>"
			<< "<red>" << red << "</red>" 
			<< "<green>" << green << "</green>" 
			<< "<blue>" << blue << "</blue>" 
		<< "</color>";
	return myString.str();
}

void MarkedObjectsPanelController::ParsePreferenceData(const char* eventData)
{
	ASErr error = kNoErr;
	ai::UnicodeString dataStr(eventData);
	
	ai::NumberFormat numFormat;
	ai::UnicodeString subFontSizeStr = MarkedObjectsPanelController::GetElement(dataStr, "subfontsize");
	if ( !subFontSizeStr.empty())
	{
		error = sAIStringFormatUtils->StringToReal(numFormat, subFontSizeStr, subFontSize);
	}

	ai::UnicodeString subFontNameStr = MarkedObjectsPanelController::GetElement(dataStr, "subfontname");
	if ( !subFontNameStr.empty())
	{
		subFontName = subFontNameStr.as_Roman();
	}

	ai::UnicodeString subLocXStr = MarkedObjectsPanelController::GetElement(dataStr, "sublocx");
	if ( !subLocXStr.empty())
	{
		float x = 0.;
		error = sAIStringFormatUtils->StringToReal(numFormat, subLocXStr, x);
		subLoc.h = x;
	}

	ai::UnicodeString subLocYStr = MarkedObjectsPanelController::GetElement(dataStr, "sublocy");
	if ( !subLocYStr.empty())
	{
		float y = 0.;
		error = sAIStringFormatUtils->StringToReal(numFormat, subLocYStr, y);
		subLoc.v = y;
	}

	ai::UnicodeString labelFontSizeStr = MarkedObjectsPanelController::GetElement(dataStr, "labelfontsize");
	if ( !labelFontSizeStr.empty())
	{
		error = sAIStringFormatUtils->StringToReal(numFormat, labelFontSizeStr, labelFontSize);
	}

	ai::UnicodeString labelFontNameStr = MarkedObjectsPanelController::GetElement(dataStr, "labelfontname");
	if ( !labelFontNameStr.empty())
	{
		labelFontName = labelFontNameStr.as_Roman();
	}

	ai::UnicodeString labelLocXStr = MarkedObjectsPanelController::GetElement(dataStr, "labellocx");
	if ( !labelLocXStr.empty())
	{
		float x = 0.;
		error = sAIStringFormatUtils->StringToReal(numFormat, labelLocXStr, x);
		labelLoc.h = x;
	}

	ai::UnicodeString labelLocYStr = MarkedObjectsPanelController::GetElement(dataStr, "labellocy");
	if ( !labelLocYStr.empty())
	{
		float y = 0.;
		error = sAIStringFormatUtils->StringToReal(numFormat, labelLocYStr, y);
		labelLoc.v = y;
	}

	ai::UnicodeString defaultTextStr = MarkedObjectsPanelController::GetElement(dataStr, "defaulttext");
	if ( !defaultTextStr.empty())
	{
		labelDefaultText = defaultTextStr.as_Roman();
	}

	ai::UnicodeString precisionStr = MarkedObjectsPanelController::GetElement(dataStr, "precision");
	if ( !precisionStr.empty())
	{
		error = sAIStringFormatUtils->StringToInteger(numFormat, precisionStr, precision);
	}

	ai::UnicodeString autoSortStr = MarkedObjectsPanelController::GetElement(dataStr, "autosort");
	if ( !autoSortStr.empty())
	{
		if ( autoSortStr.compare(ai::UnicodeString("true")) == 0 || autoSortStr.compare(ai::UnicodeString("1"))  == 0 )
		{
			autoSort = true;
		}
		else
		{
			autoSort = false;
		}
	}
}

void MarkedObjectsPanelController::ParseData(const char* eventData)
{
	ASErr error = kNoErr;
	ai::UnicodeString srcData(eventData);
	
	ai::NumberFormat numFormat;
	ai::UnicodeString itemStr = MarkedObjectsPanelController::GetElement(srcData, "item");
	if ( !itemStr.empty() )
	{
		error = sAIStringFormatUtils->StringToInteger(numFormat, itemStr, objectIndex);
	}
	
	ai::UnicodeString labelStr = MarkedObjectsPanelController::GetElement(srcData, "label");
	if ( !labelStr.empty() )
	{
		objectLabel = labelStr.as_Roman();
	}

	ai::UnicodeString IDStr = MarkedObjectsPanelController::GetElement(srcData, "ID");
	if ( !IDStr.empty() )
	{
		objectID = IDStr.as_Roman();
	}

	ai::UnicodeString locXStr  = MarkedObjectsPanelController::GetElement(srcData, "locx");
	if ( !locXStr.empty() )
	{
		error = sAIStringFormatUtils->StringToReal(numFormat, locXStr, locX);
	}
	
	ai::UnicodeString locYStr  = MarkedObjectsPanelController::GetElement(srcData, "locy");
	if ( !locYStr.empty() )
	{
		error = sAIStringFormatUtils->StringToReal(numFormat, locYStr, locY);
	}
}

void MarkedObjectsPanelController::AddMarkedObject( MarkedObject* markedObject )
{
		// Create and dispatch event
	//const string fullString = markedObject->GetFullString();
	string xmlString = "<payload><id>"+markedObject->GetID()+"</id><markedobject>" + markedObject->GetFullString() + "</markedobject></payload>";
	csxs::event::Event event = {"com.adobe.csxs.events.MarkedObjectsAddMarkedObject",
								 csxs::event::kEventScope_Application,
								 "MarkedObject",
								 NULL,
								 xmlString.c_str()};
	fPPLib.DispatchEvent(&event);

}

void MarkedObjectsPanelController::NewGlobalObjectSelected(void)
{
	ASErr error = kNoErr;
	AILayerHandle layerHandle;
	ai::UnicodeString layerTitle("Marked Objects DO NOT MODIFY"); // set in MarkedObjectManager.cpp
	sAILayer->GetLayerByTitle(&layerHandle,layerTitle);
	AIBoolean visible = true;
	error = sAILayer->GetLayerVisible(layerHandle,&visible);
		if ( error || visible == false ) return;

	MarkedObject *markedObject = this->GetMarkedObjectbyIndex(objectIndex);
	if ( markedObject == NULL )
	{
		gPlugin->UnSelectAllArt();
		gPlugin->ResetInsertionPoint();
		return;
	}

	UpdatePanelByMarkedObject(markedObject, objectIndex);
	gPlugin->UnSelectAllArt();
	
	markedObject->Select();
}

void MarkedObjectsPanelController::UpdateMarkedObjectSelected( void )
{
	// Construct XML string
	if ( fMarkedObjectManager == NULL)
		return;
	int count = fMarkedObjectManager->GetMarkedObjectSize();
	
	if ( count <= 0 )
		return;

	for(int i = count -1; i >= 0; i--)
	{
		MarkedObject* markedObject = fMarkedObjectManager->GetMarkedObject(i);
		if ( markedObject == NULL ) return;
		if ( markedObject->IsSelected() )
		{
			
			UpdatePanelByMarkedObject(markedObject, i);
			break;
		}
	}
}

void MarkedObjectsPanelController::UpdatePanelByMarkedObject(MarkedObject *markedObject, int index)
{
	if ( markedObject == NULL ) return;
	AIRealPoint location = markedObject->GetLocation();

	stringstream myString;
	myString << "<payload><item>" << index << "</item><label>" << markedObject->GetLabel()
		<< "</label><ID>" << markedObject->GetID() << "</ID>" << "<locx>"  << location.h
		<< "</locx>" << "<locy>" << location.v << "</locy></payload>";

	std::string xmlstring = myString.str();

	csxs::event::Event event = {"com.adobe.csxs.events.MarkedObjectsUpdatePanel",
								 csxs::event::kEventScope_Application,
								 "MarkedObjectsUpdatePanel",
								 NULL,
								 xmlstring.c_str()};
	fPPLib.DispatchEvent(&event);
}

void MarkedObjectsPanelController::UpdatePanelFontNameList()
{
	// Provide up to 40 fonts.
	ai::int32 count = 0;
	ASErr error = kNoErr;
	error = sAIFont->CountFonts( &count );
	
	// Provide up to 40 fonts.
	string fontNameList;
	string lastFontName;
	int entries = 0;
	for (int i = 0; i < count && !error && entries < 40; i++ )
	{
		AIFontKey fontKey;
		const int tempLength = 1024;
		char temp[tempLength];
		memset(temp, 0, tempLength);
		error = sAIFont->IndexFontList( i, &fontKey );
		if ( !error )
			error = sAIFont->GetPostScriptFontName( fontKey, temp, tempLength );
		string fontName;
		fontName.assign( temp, temp + strlen( temp ) );
        AIFontStyle* fontStyle = new AIFontStyle();
        if(sAIFont != NULL){
            if ( !error ){
                error = sAIFont->GetFontInfo(fontKey,fontStyle);
            }
        }
		if ( fontName.compare( lastFontName ) != 0 && fontName.find( "-" ) == basic_string<char>::npos && (fontStyle != NULL && fontStyle->format != kUnknownFormat) )
		{
			lastFontName = fontName;
			fontNameList.append(fontName);
			if ( entries != 39 )
				fontNameList.append(",");
			entries++;
		}
	}

	stringstream myString;
	myString << "<payload><fontlist>" << fontNameList << "</fontlist></payload>" ;
	std::string xmlstring = myString.str();

	csxs::event::Event event = {"com.adobe.csxs.events.MarkedObjectsUpdatePanelFontNameList",
								 csxs::event::kEventScope_Application,
								 "MarkedObjectsUpdatePanelFontNameList",
								 NULL,
								 xmlstring.c_str()};
	fPPLib.DispatchEvent(&event);
}

void MarkedObjectsPanelController::UpdateMarkedObject(void)
{
	// update art by the panel set
	try 
	{
		AIRealPoint location1;
		MarkedObject* markedObject = this->GetMarkedObjectbyIndex(objectIndex);
		
		if ( markedObject == NULL ) return;
		location1 = markedObject->GetLocation();
		string lable = markedObject->GetLabel();

		markedObject->SetLabel( objectLabel );
		markedObject->SetID( objectID );
	
		AIRealPoint location;
		location.h = locX;
		location.v = locY;
		
		markedObject->SetLocation( location );
		markedObject->Select();
		gPlugin->GetMarkedObjects();
	}
	
	catch(...)
	{
		return;
	}
}

void MarkedObjectsPanelController::Clear( void )
{
	stringstream myString;
	myString << "<payload><item>" << -1 << "</item><label>" << ""
		<< "</label><ID>" << "" << "</ID>" << "<locx>"  << 0
		<< "</locx>" << "<locy>" << 0 << "</locy></payload>";

	std::string xmlString = myString.str();
	csxs::event::Event event = {"com.adobe.csxs.events.MarkedObjectsUpdatePanel",
								 csxs::event::kEventScope_Application,
								 "MarkedObjectsUpdatePanel",
								 NULL,
								 xmlString.c_str()};
	fPPLib.DispatchEvent(&event);

}

void MarkedObjectsPanelController::UpdateBoundingBoxInfo( void )
{
	stringstream myString;
	myString << "<payload><width>" << MarkedObject::GetBoundingBoxWidth() << "</width><height>" << 
		MarkedObject::GetBoundingBoxHeight()<< "</height></payload>";

	std::string xmlString = myString.str();
	csxs::event::Event event = {"com.adobe.csxs.events.MarkedObjectsUpdateBoundingBoxInfo",
								 csxs::event::kEventScope_Application,
								 "MarkedObjectsUpdateBoundingBoxInfo",
								 NULL,
								 xmlString.c_str()};
	fPPLib.DispatchEvent(&event);
}

void MarkedObjectsPanelController::ClearMarkedObjectList( void)
{
	
	// Create and dispatch event
	csxs::event::Event event = {"com.adobe.csxs.events.MarkedObjectsClearMarkedObject",
								 csxs::event::kEventScope_Application,
								 "MarkedObjects",
								 NULL,
								 NULL};
	fPPLib.DispatchEvent(&event);
}

AIErr MarkedObjectsPanelController::IsPrimaryStageVisible(AIBoolean& isVisible)
{
	AIErr error = kNoErr;

	if ( sAICSXSExtension != NULL )
		error = sAICSXSExtension->IsPrimaryStageVisible(MARKEDOBJECTS_UI_EXTENSION, isVisible);

	return error;
}

MarkedObject * MarkedObjectsPanelController::GetMarkedObjectbyIndex(int index)
{
	return fMarkedObjectManager->GetMarkedObject(index);
}

ai::UnicodeString MarkedObjectsPanelController::GetElement(ai::UnicodeString srcData, const char* value)
{
	ai::UnicodeString openStr("<");
	ai::UnicodeString contentStr(value);
	ai::UnicodeString endStr(">");
	openStr.append(contentStr);
	openStr.append(endStr);
	ai::UnicodeString::size_type openIndex = srcData.find(openStr);
	openIndex = openIndex + openStr.length();
	ai::UnicodeString strValue("");
    
	ai::UnicodeString colseStr("</");
	colseStr.append(contentStr);
	colseStr.append(endStr);
	ai::UnicodeString::size_type closeIndex = srcData.find(colseStr);
	ai::UnicodeString::size_type strLength = closeIndex - openIndex;
	if ( strLength > 0 )
	{
		strValue = srcData.substr(openIndex, strLength);
		strValue = MarkedObjectsPanelController::UnescapeXML(strValue);
    }
	
	return strValue;
}

ai::UnicodeString MarkedObjectsPanelController::UnescapeXML(ai::UnicodeString srcValue)
{
	//Only unescaping ampersand currently could be expanded to other charaters
	ai::UnicodeString ampersandStr("&amp;");
	ai::UnicodeString::size_type start = srcValue.find(ampersandStr);
	if(start <= srcValue.length())
	{
		srcValue.replace(start,ampersandStr.length(),ai::UnicodeString("&"));
	}
	return srcValue;
}
