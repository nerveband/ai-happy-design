//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/StrokeFilter/Source/StrokeFilterFlashController.h $
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

#ifndef __StrokeFilterFlashController_H__
#define __StrokeFilterFlashController_H__

#include "SDKPlugPlug.h"
#include "StrokeFilters.h"
#include "StrokeFileterParams.h"
#include "FlashUIController.h"
#include "StrokeFilterPlugin.h"

/** StrokeFilterFlashController
	Controls data passing between the flash panel and StrokeFilter plug-in.	
*/

extern StrokeFilterPlugin *gPlugin;

class StrokeFilterFlashController : public FlashUIController
{
	public:
		enum FilterType
		{
			kDashStroke = 0,
			kWaveStroke = 1
		};
	public:
		StrokeFilterFlashController();

		/** Load the extension which opens the stroke filter flash panel dialog.
		*/
		AIErr DoFlashExtension(AIFilterMessage* pb, FilterType filter);

		/** Registers the events this plug-in will listen for with PlugPlug.
			@return error code if error found, or success otherwise.
		*/
		csxs::event::EventErrorCode RegisterCSXSEventListeners();

		/** Removes the previously added event listeners from PlugPlug.
			@return error code if error found, or success otherwise.
		*/
		csxs::event::EventErrorCode RemoveEventListeners();

		/** Send the current frame label values to the flash panel -
			values will either initialize the flash panel or populate it 
			with values from current item, used by InitializeDialogFields.
			Data is constructed as an XML string before being sent.
		*/
		ASErr SendWaveData();

		/** Send the current frame label values to the flash panel -
			values will either initialize the flash panel or populate it 
			with values from current item, used by InitializeDialogFields.
			Data is constructed as an XML string before being sent.
		*/
		ASErr SendDashData();

		void ParseDashData(const char* eventData);

		void ParseWaveData(const char* eventData);

		FilterType GetFilterType();

		StrokeFileterParamsDash GetDashParams();

		StrokeFileterParamsWave GetWaveParams();
		
		ASErr SendData(){return kNoErr;};

		void ParseData(const char* eventData){return ;};

	private:

		FilterType filterType;

		// Dash variables
		StrokeFileterParamsDash dashParms;

		// Wave variables
		StrokeFileterParamsWave waveParms;

		StrokeFilters strokeFilters;
};

#endif
