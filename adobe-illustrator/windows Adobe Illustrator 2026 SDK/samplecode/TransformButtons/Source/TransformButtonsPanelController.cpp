//========================================================================================
//  
//  $File: //ai/ai15/devtech/sdk/public/samplecode/TransformButtons/Source/TransformButtonsPanelController.cpp $
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

#include "TransformButtonsSuites.h"
#include "Plugin.hpp"
#include "IllustratorSDK.h"
#include "AppContext.hpp"
#include "SDKErrors.h"
#include "TransformButtonsPanelController.h"

#define TRANSFORMBUTTONS_UI_EXTENSION	"com.adobe.illustrator.transformbuttonsui.panel"
#define EVENT_TYPE_OPERATION_CLICKED	"com.adobe.csxs.events.OperationClicked"

static void OperationClickedFunc(const csxs::event::Event* const eventParam, void* const context)
{
	TransformButtonsPanelController *transformButtonsPanelController = (TransformButtonsPanelController *)context;
	if(NULL == transformButtonsPanelController || eventParam == NULL)
		return;
	transformButtonsPanelController->ParseData(eventParam->data);

	do {
		AppContext appContext(gPlugin->GetPluginRef());
		ai::UnicodeString undoText("Undo Transform");
		ai::UnicodeString redoText("Redo Transform");
		
		sAIUndo->SetUndoTextUS(undoText, redoText);
		transformButtonsPanelController->Transform();
	} while(false);
}

/* Add event listeners
*/
csxs::event::EventErrorCode TransformButtonsPanelController::RegisterCSXSEventListeners()
{
	csxs::event::EventErrorCode result = csxs::event::kEventErrorCode_Success;
	do {
		result =  fPPLib.AddEventListener(EVENT_TYPE_OPERATION_CLICKED, OperationClickedFunc, this);
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
csxs::event::EventErrorCode TransformButtonsPanelController::RemoveEventListeners()
{
	csxs::event::EventErrorCode result = csxs::event::kEventErrorCode_Success;
	do {
		result =  fPPLib.RemoveEventListener(EVENT_TYPE_OPERATION_CLICKED, OperationClickedFunc, this);
		if (result != csxs::event::kEventErrorCode_Success) 
		{		
			break;
		}
	}
	while(false);
	return result;
}

/* SendDataToPanel
*/
ASErr TransformButtonsPanelController::SendData() 
{
	return kNoErr;
}

TransformButtonsPanelController::TransformButtonsPanelController()
: FlashUIController(TRANSFORMBUTTONS_UI_EXTENSION)
{}

void TransformButtonsPanelController::ParseData(const char* eventData)
{
	ai::UnicodeString dataStr(eventData);
	//// Get index of start of payload value string
	ai::UnicodeString openStr("<operation>");
	ai::UnicodeString::size_type openIndex = dataStr.find(openStr);
	openIndex = openIndex + openStr.length();

	//// Get index of end of payload value string
	ai::UnicodeString closeStr("</operation>");
	ai::UnicodeString::size_type closeIndex = dataStr.find(closeStr);	
	// Get new art value
	ai::UnicodeString::size_type strLength = closeIndex - openIndex;

	if (strLength > 0) 
	{
		ai::UnicodeString operationStr = dataStr.substr(openIndex, strLength);
		ai::NumberFormat numFormat;
		ASInt32 nOperation = 0;
		sAIStringFormatUtils->StringToInteger(numFormat, operationStr, nOperation);
		fTransformOperation = (TransformOperation)nOperation;
	}
}

void TransformButtonsPanelController::Transform()
{
	ai::int32 numMatches = 0;
	AIArtHandle** artList = NULL;
	TransformButtonsPanelController::GetSelection(&artList, &numMatches);
	for(ASInt32 i = 0; i < numMatches; ++i)
	{
		AIArtHandle art = (*artList)[i];
		
		switch(fTransformOperation)
		{
			case UP:
				TransformButtonsPanelController::transformWithParameters( art, 0, 72, 0, 1, 1 );
				break;

			case DOWN:
				TransformButtonsPanelController::transformWithParameters( art, 0, -72, 0, 1, 1 );
				break;
				
			case LEFT:
				TransformButtonsPanelController::transformWithParameters( art, -72, 0	, 0, 1, 1 );
				break;

			case RIGHT:
				TransformButtonsPanelController::transformWithParameters( art, 72, 0, 0, 1, 1 );
				break;
			case ROTATECCW:
				TransformButtonsPanelController::transformWithParameters( art, 0, 0, (AIReal) (kAIRealPi*0.5), 1, 1 );
				break;

			case ROTATECW:
				TransformButtonsPanelController::transformWithParameters( art, 0, 0, (AIReal) (kAIRealPi*1.5), 1, 1 );
				break;
				
			case SCALEUP:
				TransformButtonsPanelController::transformWithParameters( art, 0, 0, 0, 2.0, 2.0 );
				break;

			case SCALEDOWN:
				TransformButtonsPanelController::transformWithParameters( art, 0, 0, 0, 0.5, 0.5 );
				break;
			default:
				break;
		}	
	}
}

ASErr TransformButtonsPanelController::GetSelection(AIArtHandle ***matches, ai::int32 *numMatches)
{
    short types[]={kPathArt, kCompoundPathArt, kMysteryPathArt, kTextFrameArt, kPlacedArt, kSymbolArt, kForeignArt, kRasterArt };
	int size = sizeof(types) / sizeof(*types);
	AIMatchingArtSpec *spec = new AIMatchingArtSpec[size];	
	for (int nIndex = 0; nIndex < size; nIndex++)
	{
		spec[nIndex].type = types[nIndex];
		spec[nIndex].whichAttr = kArtFullySelected;
		spec[nIndex].attr = kArtFullySelected;
	}

	ASErr result = sAIMatchingArt->GetMatchingArt(spec, size, matches, numMatches);
	delete []spec;
	if (result)
		return result;
	if (!*numMatches)
		return kCanceledErr;
	return kNoErr;
}

/*
*/
ASErr TransformButtonsPanelController::transformWithParameters(AIArtHandle art, AIReal tx, AIReal ty,
															   AIReal theta, AIReal sx, AIReal sy ) 
{ 
	ASErr result = kNoErr;
	AIRealRect artBounds;
	AIRealPoint artCenter;
	AIRealMatrix artMatrix;	
	AIReal lineScale;
	ai::int32 transformFlags = kTransformObjects | kScaleLines;
	short type;
	
	lineScale = (sAIRealMath->AIRealSqrt(sx)) * (sAIRealMath->AIRealSqrt(sy));
	
	sAIArt->GetArtType( art, &type );
	
	
	switch(type) {
	
		case kGroupArt: return result;
	
		default:
		//result = sAIDocument->GetDocumentMaxArtboardBounds( &artboardBounds );
		result = sAIArt->GetArtBounds( art, &artBounds );
		artCenter.h = artBounds.left + (artBounds.right - artBounds.left) / 2;
		artCenter.v = artBounds.bottom + (artBounds.top - artBounds.bottom) / 2;
		
		// move object so that the centerpoint is at the origin
		sAIRealMath->AIRealMatrixSetTranslate( &artMatrix, -artCenter.h, -artCenter.v );
		// translate object by tx and ty
		sAIRealMath->AIRealMatrixConcatTranslate( &artMatrix, tx, ty );
		// rotate object by theta
		sAIRealMath->AIRealMatrixConcatRotate( &artMatrix, theta );
		// scale object by sx and sy
		sAIRealMath->AIRealMatrixConcatScale( &artMatrix, sx, sy );
		// move the object back to original position
		sAIRealMath->AIRealMatrixConcatTranslate( &artMatrix, artCenter.h, artCenter.v );
		
		result = sAITransformArt->TransformArt( art, &artMatrix, lineScale, transformFlags );
		
	}
	return result;
}

// End TransformButtonsPanelController.cpp
