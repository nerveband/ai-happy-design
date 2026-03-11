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

#ifndef __TRANSFORMBUTTONSID_H__
#define __TRANSFORMBUTTONSID_H__

#define kTransformButtonsPluginName			"TransformButtons"

#define kTransformButtonsDialogID			16100
#define kTransformButtonsDialogName			"Transform Buttons"

#define kRotateCCWID						1
#define kRotateCWID							2
#define kScaleUpID							3
#define kScaleDownID						4
#define kTranslateUpID						5
#define kTranslateDownID					6
#define kTranslateLeftID					7
#define kTranslateRightID					8

#define kTransformButtonsIcon				1009
#define kTransformButtonsIconRollover		1010

#define kTransformButtonsDialogMenuID		16001

#ifdef MAC_ENV
#define layerMinHeight						74		// 3 * layerLineHeight + bottomBarHeight-2 (see Resize)
#else
#define layerMinHeight						76		// 3 * layerLineHeight + bottomBarHeight
#endif
#define layerLineHeight						20

#define layerMinWidth						206
#define layerMaxWidth						618

#endif // End TransformButtonsID.h
