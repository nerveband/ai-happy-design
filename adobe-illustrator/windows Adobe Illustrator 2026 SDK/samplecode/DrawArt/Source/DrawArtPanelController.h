//========================================================================================
//  
//  $File: //ai/mainline/devtech/sdk/public/samplecode/DrawArt/Source/DrawArtPanelController.h $
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

#ifndef __DrawArtPANELCONTROLLER_H__
#define __DrawArtPANELCONTROLLER_H__

#include "IllustratorSDK.h"
#include "HtmlUIController.h"
#include "DrawArtSuites.h"
#include "DrawArtPlugin.h"

extern DrawArtPlugin *gPlugin;

class DrawArtPanelController : public HtmlUIController
{
public:
	DrawArtPanelController();
	
	csxs::event::EventErrorCode RegisterCSXSEventListeners();
	
	csxs::event::EventErrorCode RemoveEventListeners();
	
	/*
	 Sends the bitmap string inside <data> tags to the panel
	 */
	ASErr SendBitmapToPanel(std::string& bitmap);
	
	virtual ASErr SendData() {return kNoErr;}
	
	virtual void ParseData(const char* eventData) {;}
	
};

#endif
