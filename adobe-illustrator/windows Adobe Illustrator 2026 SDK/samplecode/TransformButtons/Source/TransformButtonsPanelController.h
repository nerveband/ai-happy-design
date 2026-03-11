//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/TransformButtons/Source/TransformButtonsPanelController.h $
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
#ifndef __TRANSFORMBUTTONSPANELCONTROLLER_H__
#define __TRANSFORMBUTTONSPANELCONTROLLER_H__

#include "SDKPlugPlug.h"
#include "IllustratorSDK.h"
#include "TransformButtonsID.h"
#include "FlashUIController.h"
#include "TransformButtonsPlugin.h"

extern TransformButtonsPlugin *gPlugin;

enum TransformOperation
{
		UP = 0 ,   
		DOWN,  
		LEFT, 
		RIGHT,
		ROTATECCW,
		ROTATECW,
		SCALEUP,
		SCALEDOWN
};


class TransformButtonsPanelController : public FlashUIController
{
public:
	
	TransformButtonsPanelController();

	/** Registers the events this plug-in will listen for with PlugPlug.
		@return error code if error found, or success otherwise.
	*/
	csxs::event::EventErrorCode RegisterCSXSEventListeners();

	/** Removes the previously added event listeners from PlugPlug.
		@return error code if error found, or success otherwise.
	*/
	csxs::event::EventErrorCode RemoveEventListeners();
	
	AIErr Initialise(AIFilterMessage* pb);

	ASErr SendData();

	void ParseData(const char* eventData);

	void Transform();

private:

	TransformOperation fTransformOperation; 

	static ASErr transformWithParameters( AIArtHandle art, AIReal tx, AIReal ty, AIReal theta, AIReal sx, AIReal sy );

	static ASErr GetSelection(AIArtHandle ***matches, ai::int32 *numMatches);
};

#endif
