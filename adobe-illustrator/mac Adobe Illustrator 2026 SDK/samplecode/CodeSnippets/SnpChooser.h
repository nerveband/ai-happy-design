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

#ifndef  __SnpTypeChooser_h__
#define  __SnpTypeChooser_h__

/** Helper class that prompts for commonly used Illustrator types.

  	@ingroup sdk_snippet
 */
class SnpChooser
{
	public:
		/** Lists the names of open documents in a prompt 
			and returns handle to chosen document.
			@return handle of chosen document.
		*/
		AIDocumentHandle PickDocument();

		/** Lists the names of layers in the current documents 
			in a prompt and returns handle to chosen layer -
			if the document only has one layer no prompt is
			made and this layer is returned.
			@return handle of chosen layer.
		*/
		AILayerHandle PickLayer();

		/** Prompts for a type of art.
			@return chosen art type.
		*/
		AIArtType PickArtType();

		/** Prompts for a document color model.
			@return chosen color model.
		*/
		AIDocumentColorModelValue PickDocumentColorModel();

		/** Prompts for a document ruler unit.
			@return chosen document ruler unit.
		*/
		AIDocumentRulerUnitValue PickDocumentRulerUnit();

		/** Prompts for a raster resolution.
			@return chosen raster resolution.
		*/
		AIRasterResolution PickDocumentRasterResolution();

		/** Prompts for a preview mode.
			@return chosen preview mode.
		*/
		AIPreviewMode PickPreviewMode();
};

#endif // __SnpTypeChooser_h__

// End SnpChooser.h

