//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/StrokeFilter/Source/StrokeFilterFlashController.cpp $
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
#include "IllustratorSDK.h"
#include "StrokeFilterFlashController.h"
#include "SDKErrors.h"
#include "StrokeFilterSuites.h"

#define STROKE_FILTER_UI_EXTENSION	"com.adobe.illustrator.strokefilterui.dialog"

#define EVENT_TYPE_OK_CLICKED		"com.adobe.csxs.events.StrokeOkClicked"
#define	EVENT_TYPE_CANCEL_CLICKED	"com.adobe.csxs.events.StrokeCancelClicked"
#define EVENT_TYPE_PANEL_READY		"com.adobe.csxs.events.StrokePanelReady"


static void OkClickedFunc(const csxs::event::Event* const eventParam, void* const context)
{
	StrokeFilterFlashController *strokeFilterFlashController = (StrokeFilterFlashController *)context;

	if (strokeFilterFlashController->GetFilterType() == StrokeFilterFlashController::kDashStroke)
	{
		// Dash stroke dialog
		strokeFilterFlashController->ParseDashData(eventParam->data);
		sAINotifier->Notify(kStrokeFilterUpdatePathNotifier, nullptr);
	}
	else if (strokeFilterFlashController->GetFilterType() == StrokeFilterFlashController::kWaveStroke)
	{
		// Wave stroke dialog
		strokeFilterFlashController->ParseWaveData(eventParam->data);
		sAINotifier->Notify(kStrokeFilterUpdatePathNotifier, nullptr);
	}
}

static void PanelReadyFunc(const csxs::event::Event* const eventParam, void* const context)
{
	StrokeFilterFlashController *strokeFilterFlashController = (StrokeFilterFlashController *)context;

	if (strokeFilterFlashController->GetFilterType() == StrokeFilterFlashController::kDashStroke)
	{
		strokeFilterFlashController->SendDashData();
	}
	else if (strokeFilterFlashController->GetFilterType() == StrokeFilterFlashController::kWaveStroke)
	{
		strokeFilterFlashController->SendWaveData();
	}
}

StrokeFilterFlashController::StrokeFilterFlashController(void)
: FlashUIController(STROKE_FILTER_UI_EXTENSION),
filterType(kDashStroke)
{
}

/* Start the extension
*/
ASErr StrokeFilterFlashController::DoFlashExtension(AIFilterMessage* pb, FilterType filter)
{
	filterType = filter;

	// Load the StrokeFilter UI extension
	PlugPlugErrorCode result2 = LoadExtension();
	if (result2 != PlugPlugErrorCode_success)
	{
		SDK_ASSERT(false);
		return kCantHappenErr;
	}

	return kNoErr;
}

/* Add event listeners
*/
csxs::event::EventErrorCode StrokeFilterFlashController::RegisterCSXSEventListeners()
{

	csxs::event::EventErrorCode result = csxs::event::kEventErrorCode_Success;
	do {
		result = fPPLib.AddEventListener(EVENT_TYPE_OK_CLICKED, OkClickedFunc, this);
		if (result != csxs::event::kEventErrorCode_Success)
		{
			break;
		}

		result = fPPLib.AddEventListener(EVENT_TYPE_PANEL_READY, PanelReadyFunc, this);
		if (result != csxs::event::kEventErrorCode_Success)
		{
			break;
		}
	}
	while(false);

	return result;
}

/* Remove event listeners
*/
csxs::event::EventErrorCode StrokeFilterFlashController::RemoveEventListeners()
{
	csxs::event::EventErrorCode result = csxs::event::kEventErrorCode_Success;
	do {
		result = fPPLib.RemoveEventListener(EVENT_TYPE_OK_CLICKED, OkClickedFunc, this);
		if (result != csxs::event::kEventErrorCode_Success)
		{
			break;
		}

		result = fPPLib.RemoveEventListener(EVENT_TYPE_PANEL_READY, PanelReadyFunc, this);
		if (result != csxs::event::kEventErrorCode_Success)
		{
			break;
		}
	}
	while(false);

	return result;
}

/* SendDashDataToPanel
*/
ASErr StrokeFilterFlashController::SendDashData() 
{
	ASErr error = kNoErr;

	// Get dashed value
	std::string dashedCStr = "false";
	if (dashParms.dashed == 1) {
		dashedCStr = "true";
	}

	// Get adjust value
	std::string adjustCStr = "false";
	if (dashParms.adjust == 1) {
		adjustCStr = "true";
	}

	ai::NumberFormat numFormat;
	// Get dash1 value
	ai::UnicodeString valueUString;
	error = sAIStringFormatUtils->RealToString(numFormat, dashParms.dashes[0], 1, true, valueUString);
	if (error)
		return error;
	// Convert to std::string
	std::string dash1CStr = valueUString.as_Platform();

	// Get gap1 value
	error = sAIStringFormatUtils->RealToString(numFormat, dashParms.dashes[1], 1, true, valueUString);
	if (error)
		return error;
	// Convert to std::string
	std::string gap1CStr = valueUString.as_Platform();

	// Get dash2 value
	error = sAIStringFormatUtils->RealToString(numFormat, dashParms.dashes[2], 1, true, valueUString);
	if (error)
		return error;
	// Convert to std::string
	std::string dash2CStr = valueUString.as_Platform();

	// Get gap2 value
	error = sAIStringFormatUtils->RealToString(numFormat, dashParms.dashes[3], 1, true, valueUString);
	if (error)
		return error;
	// Convert to std::string
	std::string gap2CStr = valueUString.as_Platform();

	// Get dash3 value
	error = sAIStringFormatUtils->RealToString(numFormat, dashParms.dashes[4], 1, true, valueUString);
	if (error)
		return error;
	// Convert to std::string
	std::string dash3CStr = valueUString.as_Platform();

	// Get gap3 value
	error = sAIStringFormatUtils->RealToString(numFormat, dashParms.dashes[5], 1, true, valueUString);
	if (error)
		return error;
	// Convert to std::string
	std::string gap3CStr = valueUString.as_Platform();

	// Construct XML string
	std::string xmlString = "<payload><dashed>" + dashedCStr + "</dashed><adjust>" + adjustCStr + "</adjust><dash1>" + dash1CStr + "</dash1>" +
		"<gap1>" + gap1CStr + "</gap1><dash2>" + dash2CStr + "</dash2><gap2>" + gap2CStr + "</gap2>" +
		"<dash3>" + dash3CStr + "</dash3><gap3>" + gap3CStr + "</gap3></payload>";
	
	// Create and dispatch event
	csxs::event::Event event = {"com.adobe.csxs.events.StrokeUpdateDashPanel",
                                      csxs::event::kEventScope_Application,
                                      "StrokeFilter",
                                      NULL,
                                      xmlString.c_str()};
	fPPLib.DispatchEvent(&event);

	return error;
}

/* SendWaveDataToPanel
*/
ASErr StrokeFilterFlashController::SendWaveData() 
{
	ASErr error = kNoErr;

	// Get dashed value
	std::string exchangeCStr = "false";
	if (waveParms.exchange == 1) {
		exchangeCStr = "true";
	}

	// Get adjust value
	std::string replaceCStr = "false";
	if (waveParms.replace == 1) {
		replaceCStr = "true";
	}

	ai::NumberFormat numFormat;
	// Get number value
	ai::UnicodeString valueUString;
	error = sAIStringFormatUtils->IntegerToString(numFormat, waveParms.number, valueUString);
	if (error)
		return error;
	// Convert to std::string
	std::string numberCStr = valueUString.as_Platform();

	// Get width1 value
	error = sAIStringFormatUtils->RealToString(numFormat, waveParms.width1, 1, true, valueUString);
	if (error)
		return error;
	// Convert to std::string
	std::string width1CStr = valueUString.as_Platform();

	// Get width2 value
	error = sAIStringFormatUtils->RealToString(numFormat, waveParms.width2, 1, true, valueUString);
	if (error)
		return error;
	// Convert to std::string
	std::string width2CStr = valueUString.as_Platform();

	// Construct XML string
	std::string xmlString = "<payload><exchange>" + exchangeCStr + "</exchange><replace>" + replaceCStr + "</replace><number>" + numberCStr + "</number>" +
		"<width1>" + width1CStr + "</width1><width2>" + width2CStr + "</width2></payload>";

	// Create and dispatch event
	csxs::event::Event event = {"com.adobe.csxs.events.StrokeUpdateWavePanel",
                                      csxs::event::kEventScope_Application,
                                      "test",
                                      NULL,
                                      xmlString.c_str()};
	fPPLib.DispatchEvent(&event);

	return error;
}

/* Parse received data from flash panel.
*/
void StrokeFilterFlashController::ParseDashData(const char* eventData)
{
	ASErr error = kNoErr;
	ai::UnicodeString dataStr(eventData);

	// Get index of start of dashed value string
	ai::UnicodeString openDashedStr("<dashed>");
	ai::UnicodeString::size_type openIndex = dataStr.find(openDashedStr);
	openIndex = openIndex + openDashedStr.length();

	// Get index of end of dashed value string
	ai::UnicodeString closeDashedStr("</dashed>");
	ai::UnicodeString::size_type closeIndex = dataStr.find(closeDashedStr);	

	// Get dashed value
	ai::UnicodeString::size_type strLength = closeIndex - openIndex;
	if (strLength > 0) {
		ai::UnicodeString dashedStr = dataStr.substr(openIndex, strLength);
		if (dashedStr == ai::UnicodeString("true")) {
			dashParms.dashed = TRUE;
		}
		else {
			dashParms.dashed = FALSE;
		}
	}

	// Get index of start of adjust dash value string
	ai::UnicodeString openAdjustStr("<adjust>");
	openIndex = dataStr.find(openAdjustStr);
	openIndex = openIndex + openAdjustStr.length();

	// Get index of end of adjust dash value string
	ai::UnicodeString closeAdjustStr("</adjust>");
	closeIndex = dataStr.find(closeAdjustStr);	

	// Get adjust dash value
	strLength = closeIndex - openIndex;
	if (strLength > 0) {
		ai::UnicodeString adjustStr = dataStr.substr(openIndex, strLength);
		if (adjustStr == ai::UnicodeString("true")) {
			dashParms.adjust = TRUE;
		}
		else {
			dashParms.adjust = FALSE;
		}
	}

	ai::NumberFormat numFormat;

	// Get index of start of first dash value
	ai::UnicodeString openDash1Str("<dash1>");
	openIndex = dataStr.find(openDash1Str);
	openIndex = openIndex + openDash1Str.length();

	// Get index of end of first dash value
	ai::UnicodeString closeDash1Str("</dash1>");
	closeIndex = dataStr.find(closeDash1Str);	

	// Get dash value
	strLength = closeIndex - openIndex;
	if (strLength > 0) {
		ai::UnicodeString dash1Str = dataStr.substr(openIndex, strLength);
		// convert to AIReal
		float dash1 = 0.0;
		error = sAIStringFormatUtils->StringToReal(numFormat, dash1Str, dash1);
		// Add to array
		dashParms.dashes[0] = dash1;
	}

	// Get index of start of first gap value
	ai::UnicodeString openGap1Str("<gap1>");
	openIndex = dataStr.find(openGap1Str);
	openIndex = openIndex + openGap1Str.length();

	// Get index of end of first gap value
	ai::UnicodeString closeGap1Str("</gap1>");
	closeIndex = dataStr.find(closeGap1Str);	

	// Get gap value
	strLength = closeIndex - openIndex;
	if (strLength > 0) {
		ai::UnicodeString gap1Str = dataStr.substr(openIndex, strLength);
		// convert to AIReal
		float gap1 = 0.0;
		error = sAIStringFormatUtils->StringToReal(numFormat, gap1Str, gap1);
		// Add to array
		dashParms.dashes[1] = gap1;
	}

	// Get index of start of second dash value
	ai::UnicodeString openDash2Str("<dash2>");
	openIndex = dataStr.find(openDash2Str);
	openIndex = openIndex + openDash2Str.length();

	// Get index of end of second dash value
	ai::UnicodeString closeDash2Str("</dash2>");
	closeIndex = dataStr.find(closeDash2Str);	

	// Get dash value
	strLength = closeIndex - openIndex;
	if (strLength > 0) {
		ai::UnicodeString dash2Str = dataStr.substr(openIndex, strLength);
		// convert to AIReal
		float dash2 = 0.0;
		error = sAIStringFormatUtils->StringToReal(numFormat, dash2Str, dash2);
		// Add to array
		dashParms.dashes[2] = dash2;
	}

	// Get index of start of second gap value
	ai::UnicodeString openGap2Str("<gap2>");
	openIndex = dataStr.find(openGap2Str);
	openIndex = openIndex + openGap2Str.length();

	// Get index of end of second gap value
	ai::UnicodeString closeGap2Str("</gap2>");
	closeIndex = dataStr.find(closeGap2Str);	

	// Get gap value
	strLength = closeIndex - openIndex;
	if (strLength > 0) {
		ai::UnicodeString gap2Str = dataStr.substr(openIndex, strLength);
		// convert to AIReal
		float gap2 = 0.0;
		error = sAIStringFormatUtils->StringToReal(numFormat, gap2Str, gap2);
		// Add to array
		dashParms.dashes[3] = gap2;
	}	

	// Get index of start of third dash value
	ai::UnicodeString openDash3Str("<dash3>");
	openIndex = dataStr.find(openDash3Str);
	openIndex = openIndex + openDash3Str.length();

	// Get index of end of third dash value
	ai::UnicodeString closeDash3Str("</dash3>");
	closeIndex = dataStr.find(closeDash3Str);	

	// Get dash value
	strLength = closeIndex - openIndex;
	if (strLength > 0) {
		ai::UnicodeString dash3Str = dataStr.substr(openIndex, strLength);
		// convert to AIReal
		float dash3 = 0.0;
		error = sAIStringFormatUtils->StringToReal(numFormat, dash3Str, dash3);
		// Add to array
		dashParms.dashes[4] = dash3;
	}

	// Get index of start of third gap value
	ai::UnicodeString openGap3Str("<gap3>");
	openIndex = dataStr.find(openGap3Str);
	openIndex = openIndex + openGap3Str.length();

	// Get index of end of third gap value
	ai::UnicodeString closeGap3Str("</gap3>");
	closeIndex = dataStr.find(closeGap3Str);	

	// Get gap value
	strLength = closeIndex - openIndex;
	if (strLength > 0) {
		ai::UnicodeString gap3Str = dataStr.substr(openIndex, strLength);
		// convert to AIReal
		float gap3 = 0.0;
		error = sAIStringFormatUtils->StringToReal(numFormat, gap3Str, gap3);
		// Add to array
		dashParms.dashes[5] = gap3;
	}		
	sAIFilter->SetFilterParameters(gPlugin->GetDashFilterHandle(),&dashParms);
}

/* Parse received data from flash panel.
*/
void StrokeFilterFlashController::ParseWaveData(const char* eventData)
{
	ASErr error = kNoErr;
	ai::UnicodeString dataStr(eventData);

	// Get index of start of exchange value string
	ai::UnicodeString openExchangeStr("<exchange>");
	ai::UnicodeString::size_type openIndex = dataStr.find(openExchangeStr);
	openIndex = openIndex + openExchangeStr.length();

	// Get index of end of exchange value string
	ai::UnicodeString closeExchangeStr("</exchange>");
	ai::UnicodeString::size_type closeIndex = dataStr.find(closeExchangeStr);	

	// Get exchange value
	ai::UnicodeString::size_type strLength = closeIndex - openIndex;
	if (strLength > 0) {
		ai::UnicodeString exchangeStr = dataStr.substr(openIndex, strLength);
		if (exchangeStr == ai::UnicodeString("true")) {
			waveParms.exchange = TRUE;
		}
		else {
			waveParms.exchange = FALSE;
		}
	}

	// Get index of start of replace value string
	ai::UnicodeString openReplaceStr("<replace>");
	openIndex = dataStr.find(openReplaceStr);
	openIndex = openIndex + openReplaceStr.length();

	// Get index of end of replace value string
	ai::UnicodeString closeReplaceStr("</replace>");
	closeIndex = dataStr.find(closeReplaceStr);	

	// Get replace value
	strLength = closeIndex - openIndex;
	if (strLength > 0) {
		ai::UnicodeString replaceStr = dataStr.substr(openIndex, strLength);
		if (replaceStr == ai::UnicodeString("true")) {
			waveParms.replace = TRUE;
		}
		else {
			waveParms.replace = FALSE;
		}
	}

	ai::NumberFormat numFormat;

	// Get index of start of number value
	ai::UnicodeString openNumberStr("<number>");
	openIndex = dataStr.find(openNumberStr);
	openIndex = openIndex + openNumberStr.length();

	// Get index of end of number value
	ai::UnicodeString closeNumberStr("</number>");
	closeIndex = dataStr.find(closeNumberStr);	

	// Get number value
	strLength = closeIndex - openIndex;
	if (strLength > 0) {
		ai::UnicodeString numberStr = dataStr.substr(openIndex, strLength);
		// convert to int
		ASInt32 number = 0;
		error = sAIStringFormatUtils->StringToInteger(numFormat, numberStr, number);
		// Add to array
		waveParms.number = number;
	}

	// Get index of start of width1 value
	ai::UnicodeString openWidth1Str("<width1>");
	openIndex = dataStr.find(openWidth1Str);
	openIndex = openIndex + openWidth1Str.length();

	// Get index of end of width1 value
	ai::UnicodeString closeWidth1Str("</width1>");
	closeIndex = dataStr.find(closeWidth1Str);	

	// Get number value
	strLength = closeIndex - openIndex;
	if (strLength > 0) {
		ai::UnicodeString width1Str = dataStr.substr(openIndex, strLength);
		// convert to float
		float width1 = 0.0;
		error = sAIStringFormatUtils->StringToReal(numFormat, width1Str, width1);
		// Add to array
		waveParms.width1 = width1;
	}

	// Get index of start of width2 value
	ai::UnicodeString openWidth2Str("<width2>");
	openIndex = dataStr.find(openWidth2Str);
	openIndex = openIndex + openWidth2Str.length();

	// Get index of end of width2 value
	ai::UnicodeString closeWidth2Str("</width2>");
	closeIndex = dataStr.find(closeWidth2Str);	

	// Get number value
	strLength = closeIndex - openIndex;
	if (strLength > 0) {
		ai::UnicodeString width2Str = dataStr.substr(openIndex, strLength);
		// convert to float
		float width2 = 0.0;
		error = sAIStringFormatUtils->StringToReal(numFormat, width2Str, width2);
		// Add to array
		waveParms.width2 = width2;
	}
	sAIFilter->SetFilterParameters(gPlugin->GetWaveFilterHandle(),&waveParms);
}

StrokeFilterFlashController::FilterType StrokeFilterFlashController::GetFilterType()
{
	return filterType;
}

StrokeFileterParamsDash StrokeFilterFlashController::GetDashParams()
{
	return dashParms;
}

StrokeFileterParamsWave StrokeFilterFlashController::GetWaveParams()
{
	return waveParms;
}
