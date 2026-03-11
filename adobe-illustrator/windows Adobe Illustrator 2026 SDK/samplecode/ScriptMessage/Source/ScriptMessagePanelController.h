//========================================================================================
//  
//  $File: //ai/mainline/devtech/sdk/public/samplecode/ScriptMessage/Source/ScriptMessagePanelController.h $
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
#include "FlashUIController.h"
#include "ScriptMessageSuites.h"

#ifndef __SCRIPTMESSAGEPANELCONTROLLER_H__
#define __SCRIPTMESSAGEPANELCONTROLLER_H__

class ScriptMessagePanelController : public FlashUIController
{
public:
	ScriptMessagePanelController();
	
	csxs::event::EventErrorCode RegisterCSXSEventListeners();
	
	csxs::event::EventErrorCode RemoveEventListeners();
	
	ASErr SendThemeToPanel();
	
	virtual ASErr SendData() {return kNoErr;}
	
	virtual void ParseData(const char* eventData) {;}
	
private:
	std::string AIRGBColorToXML(AIUIThemeColor color);
};

#endif
