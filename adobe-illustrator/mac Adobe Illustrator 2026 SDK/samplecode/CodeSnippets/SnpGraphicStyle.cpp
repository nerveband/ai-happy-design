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

#include "IllustratorSDK.h"

#include "SDKErrors.h"

// Framework includes:
#include "SnpRunnable.h"
#include "SnippetRunnerSuites.h" 
#include "SnippetRunnerLog.h"
#include "SnippetRunnerParameter.h"
#include "SnpSelectionHelper.h"

/** 
	This class provides functions for the Appearance and Graphic Style Palette of Ai.
	Provides functionalities for displaying the art styles available with a document;
	inspecting and editing the fill, stroke and blend properties of a selected art object;
	inspecting and editing the Live Effects applied on an art object;
	and for removing the paint fields and the effects of an art object.
	
 */
class SnpGraphicStyle
{
	public:
		/**
			Displays the art styles available with a document.
			@return kNoErr on success, other ASErr otherwise.
		*/		
		ASErr DisplayArtStyle(void);

		/**
			Creates and adds new Graphic Style in the Graphics Style Palette, 
			based on an existing style of a selected art object.
			User is prmpted for nameing the new Graphic Style.
			@return kNoErr on success, other ASErr otherwise.
		*/	
		ASErr NewStyleFromExisting(void);

		/**
			Inspects and displays the Fill,Stroke, and Blend properties of the selected art object.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr InspectFillStrokeAndBlend(void);

		/**
			Inspects and displays the various live effects that have been applied on an artstyle.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr InspectLiveEffects(void);

		/**
			Modifies various effects like Fill,Stroke, and Blend, Pre/Post effects etc. of the selected art object.
			User is prompted for changing the properties.
			Various operations are done on the effects as an illustration of how modifications can be done on the styles.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr ModifyEffects(void);

		/**
			Counts and modifies the various live effects that have been applied on an artstyle.
			No editing happens in case of 0 Live Effects on the art object.
			LiveEffects if not already there,may be manually applied on selected artwork. Then various properties may then
			be modified through this function.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr ModifyLiveEffects(void);

		/**
			Removes a Graphic Style from the Graphic Styles Palette by choice.
			User is prompted for the name of the style to be deleted.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr RemoveArtStyleByChoice(void);

		/**
			Removes the effects of the selected art object by choice.
			User is prompted for the number of paintfield to be deleted.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr RemoveEffectByChoice(void);

		/**
			Removes all the effects of the selected art object.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr RemoveAllEffects(void);

		//------------------------------- Utility Functions -------------------------------

		/**
			Creates new Fill Style for the Fill Field of the selected art object.
			@param parserPaintField IN/OUT paintfield for which fill is being modified.
			@param fill IN/OUT that is being modified for parserPaintField.
			@return kNoErr on success, other ASErr otherwise.
		*/	
		ASErr CreateFillStyle(AIParserPaintField* parserPaintField, AIFillStyle fill);

		/**
			Creates new Stroke Style for the Stroke Field of the selected art object.
			@param parserPaintField IN/OUT paintfield for which fill is being modified.
			@param stroke  IN/OUT that is being modified for parserPaintField.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr CreateStrokeStyle(AIParserPaintField* parserPaintField, AIStrokeStyle stroke);

		/**
			Give new color properties to the AIFillStyle object to create a new Fill Style.
			@param fill IN/OUT whose color properties are changed, and then returned back to the calling function, 
			where the new style is being created.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr FillNewColor(AIFillStyle* fill);

		/**
			Give new color properties to the AIStrokeStyle object to create a new Stroke Style.
			@param stroke IN/OUT whose color properties are changed, and then returned back to the calling function, 
			where the new style is being created.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr StrokeNewColor(AIStrokeStyle* stroke);

		/**
			Gets the Fill Style for the Fill Field of the selected art object.
			Expected values should be kNoneColor, kThreeColor, or kFourColor.
			@param fill IN whose fill style is being parsed.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr GetFillColorKind(const AIFillStyle fill);

		/**
			Gets the Stroke Style for the Stroke Field of the selected art object.
			Expected values should be kNoneColor, kThreeColor, or kFourColor.
			@param stroke IN whose stroke style is being parsed.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr GetStrokeColorKind(const AIStrokeStyle stroke);
};

/*
*/
ASErr SnpGraphicStyle :: StrokeNewColor(AIStrokeStyle* stroke)
{
	ASErr result = kNoErr;
	SnippetRunnerLog* log = SnippetRunnerLog::Instance();
	log->Write("");
	log->Write("New Stroke color");
	SnippetRunnerLog::Indent indent;
	try {
		//Give new values to Stroke.
		stroke->color.kind = kThreeColor;
		stroke->color.c.rgb.red = 1.0;
		stroke->overprint = true;
		stroke->width = 2.5;			
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGraphicStyle :: CreateStrokeStyle(AIParserPaintField* parserPaintField, AIStrokeStyle stroke)
{
	ASErr result = kNoErr;
	SnippetRunnerLog* log = SnippetRunnerLog::Instance();
	SnippetRunnerLog::Indent indent;
	
	try {		
			//Get the current fill style.
			result = this->GetStrokeColorKind(stroke);
			aisdk::check_ai_error(result);

			//Give new color properties to the stroke.
			result = this->StrokeNewColor(&stroke);
			aisdk::check_ai_error(result);						

			//Set the new stroke.
			result = sAIArtStyleParser->SetStroke(*parserPaintField,&stroke,NULL);
			aisdk::check_ai_error(result);
			
			log->Write("The fill style after changing is as follows :");
			//Get the new fill style.
			result = this->GetStrokeColorKind(stroke);
			aisdk::check_ai_error(result);
		
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/

ASErr SnpGraphicStyle :: CreateFillStyle(AIParserPaintField* parserPaintField, AIFillStyle fill)
{
	ASErr result = kNoErr;
	SnippetRunnerLog* log = SnippetRunnerLog::Instance();
	SnippetRunnerLog::Indent indent;
	
	try {		
			//Get the current fill style.
			result = this->GetFillColorKind(fill);
			aisdk::check_ai_error(result);

			//Give new color properties to the fill.
			result = this->FillNewColor(&fill);
			aisdk::check_ai_error(result);						

			//Set the new fill.
			result = sAIArtStyleParser->SetFill(*parserPaintField,&fill,NULL);
			aisdk::check_ai_error(result);
			
			log->Write("The fill style after changing is as follows");
			//Get the new fill style.
			result = this->GetFillColorKind(fill);
			aisdk::check_ai_error(result);		
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGraphicStyle :: FillNewColor(AIFillStyle* fill)
{
	ASErr result = kNoErr;
	SnippetRunnerLog* log = SnippetRunnerLog::Instance();
	log->Write("");
	log->Write("New Fill color");
	SnippetRunnerLog::Indent indent;
	try {	
		//Give new values to Fill.
		fill->color.c.f.black = (AIReal)0.4;
		fill->color.c.f.cyan =(AIReal) 0.8;
		fill->color.c.f.magenta =(AIReal) 0.6;
		fill->color.c.f.yellow =(AIReal) 0.1;	
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGraphicStyle::GetFillColorKind(const AIFillStyle fill)
{
	ASErr result = kNoErr;
	SnippetRunnerLog* log = SnippetRunnerLog::Instance();
	log->Write("Displaying Fill color");
	SnippetRunnerLog::Indent indent;
	try {
		switch(fill.color.kind){
			case kGrayColor:
				log->Write("Color kind : kGrayColor");
				log->WritePrintf("Gray = %f",fill.color.c.g.gray);
				break;
			case kFourColor:
				log->Write("Color kind : kFourColor");
				log->WritePrintf("Cyan = %f",fill.color.c.f.cyan);
				log->WritePrintf("Magenta = %f",fill.color.c.f.magenta);
				log->WritePrintf("Yellow= %f",fill.color.c.f.yellow);
				log->WritePrintf("Black = %f",fill.color.c.f.black);
				break;
			case kThreeColor:
				log->Write("Color kind : kThreeColor");
				log->WritePrintf("Red = %f",fill.color.c.rgb.red);
				log->WritePrintf("Green = %f",fill.color.c.rgb.green);
				log->WritePrintf("Blue = %f",fill.color.c.rgb.blue);
				break;
			case kNoneColor:
				log->Write("Color kind : kNoneColor");
				break;
			default:
				log->Write("None of the known colors");
				break;
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGraphicStyle::GetStrokeColorKind(const AIStrokeStyle stroke)
{
	ASErr result = kNoErr;
	SnippetRunnerLog* log = SnippetRunnerLog::Instance();
	SnippetRunnerLog::Indent indent;
	log->Write("Displaying Stroke color");
	log->WritePrintf("Width = %f",stroke.width);
	log->WritePrintf("MiterLimit= %f",stroke.miterLimit);
	if(stroke.overprint)
		log->Write("Overprint true");
	else
		log->Write("Overprint false");

	try {
		switch(stroke.color.kind){
			case kGrayColor:
				log->Write("Color kind : kGrayColor");
				log->WritePrintf("Gray = %f",stroke.color.c.g.gray);
				break;
			case kFourColor:
				log->Write("Color kind : kFourColor");
				log->WritePrintf("Cyan = %f",stroke.color.c.f.cyan);
				log->WritePrintf("Magenta = %f",stroke.color.c.f.magenta);
				log->WritePrintf("Yellow= %f",stroke.color.c.f.yellow);
				log->WritePrintf("Black = %f",stroke.color.c.f.black);
				break;
			case kThreeColor:
				log->Write("Color kind : kThreeColor");
				log->WritePrintf("Red = %f",stroke.color.c.rgb.red);
				log->WritePrintf("Green = %f",stroke.color.c.rgb.green);
				log->WritePrintf("Blue = %f",stroke.color.c.rgb.blue);
				break;
			case kNoneColor:
				log->Write("Color kind : kNoneColor");
				break;
			default:
				log->Write("None of the known colors");
				break;
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGraphicStyle::NewStyleFromExisting(void)
{
	ASErr result = kNoErr;
	try {
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerLog::Indent indent;
		
		AIArtHandle **matchingArt = NULL;
		ai::int32 numMatches;
		result = sAIMatchingArt->GetSelectedArt(&matchingArt,&numMatches);
		aisdk::check_ai_error(result);

		for( ai::int32 i = 0; i < numMatches; i++){
			AIArtHandle art;
			art = (*matchingArt)[i];

			if( i == 1 ){

				//Get the artsyle of selected art.
				AIArtStyleHandle style;
				result = sAIArtStyle->GetArtStyle(art,&style);
				aisdk::check_ai_error(result);

				//Create a new parser.
				AIStyleParser parser = NULL;
				result = sAIArtStyleParser->NewParser(&parser);
				aisdk::check_ai_error(result);

				//Parse the selected artStyle.
				result = sAIArtStyleParser->ParseStyle(parser,style);
				aisdk::check_ai_error(result);

				ASBoolean parsable;
				parsable = sAIArtStyleParser->IsStyleParseable(parser);

				if(parsable){
					//Get the blend field.
					AIParserBlendField blendField;
					result = sAIArtStyleParser->GetStyleBlendField(parser,&blendField);

					//Check whether the selected item contains paint.
					if(sAIArtStyleParser->ContainsPaint(parser)){
						log->Write("Contains paint");

						ai::int32 numPaintFields = sAIArtStyleParser->CountPaintFields(parser);
						log->WritePrintf("Number of paint fields in art number %d is %d ",i+1,numPaintFields);

						//Get paintfields on selected artwork.
						for(ai::int32 n = 0; n < numPaintFields ;n++){
							AIParserPaintField parserPaintField;
							result = sAIArtStyleParser->GetNthPaintField(parser,n,&parserPaintField);
							aisdk::check_ai_error(result);

							//Check whether the selected paintfield is a fill.
							if(sAIArtStyleParser->IsFill(parserPaintField)){
								log->Write("");
								log->WritePrintf("Selected item Fill");

								AIFillStyle fill;
								result = sAIArtStyleParser->GetFill(parserPaintField,&fill,NULL);
								aisdk::check_ai_error(result);

								//Create a new Fill.
								result = this->CreateFillStyle(&parserPaintField,fill);
								aisdk::check_ai_error(result);

							}
							//Check whether the selected paintfield is a stroke.
							if(sAIArtStyleParser->IsStroke(parserPaintField)){
								log->Write("");
								log->WritePrintf("Selected item Stroke");
								AIStrokeStyle stroke;
								result = sAIArtStyleParser->GetStroke(parserPaintField,&stroke,NULL);
								aisdk::check_ai_error(result);

								//Create a new Stroke.
								result = this->CreateStrokeStyle(&parserPaintField,stroke);
								aisdk::check_ai_error(result);
							}
						}
					}
					//Create new artstyle based on the modifications in parser.s
					AIArtStyleHandle newArtStyle;
					result = sAIArtStyleParser->CreateNewStyle(parser,&newArtStyle);
					aisdk::check_ai_error(result);

					//Give a name to the new Graphic Style.
					ai::UnicodeString newName = SnippetRunnerParameter::Instance()->GetString(ai::UnicodeString("Enter name for the new ArtStyle"),ai::UnicodeString("New Graphic Style"));
					AIArtStyleHandle namedStyle;

					//Add the new art style to the Graphic Style Palette.
					result = sAIArtStyle->AddNamedStyle(newArtStyle,newName,true,&namedStyle);
					aisdk::check_ai_error(result);
					log->Write(ai::UnicodeString("New name for ArtStyle is: ").append(newName));

					//Apply the new style to the selected art object.
					result = sAIArtStyle->SetCurrentArtStyle(newArtStyle);
					aisdk::check_ai_error(result);
				}
				//Free the memory associated with the parser as it is not required now.
				result = sAIArtStyleParser->DisposeParser(parser);
				aisdk::check_ai_error(result);
			}				
		}		
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGraphicStyle::InspectFillStrokeAndBlend(void)
{
	ASErr result = kNoErr;
	try {
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerLog::Indent indent;
		
		AIArtHandle **matchingArt = NULL;
		ai::int32 numMatches;
		result = sAIMatchingArt->GetSelectedArt(&matchingArt,&numMatches);
		aisdk::check_ai_error(result);

		for( ai::int32 i = 0; i < numMatches; i++){
			AIArtHandle art;
			art = (*matchingArt)[i];
			
			//Get the artsyle of selected art.
			AIArtStyleHandle style;
			result = sAIArtStyle->GetArtStyle(art,&style);
			aisdk::check_ai_error(result);

			//Create a new parser.
			AIStyleParser parser = NULL;
			result = sAIArtStyleParser->NewParser(&parser);
			aisdk::check_ai_error(result);
			
			//Parse the selected artstyle.
			result = sAIArtStyleParser->ParseStyle(parser,style);
			aisdk::check_ai_error(result);

			ASBoolean boolean;
			boolean = sAIArtStyleParser->IsStyleParseable(parser);

			//Get the blendfield and check its properties.
			AIParserBlendField blendField;
			result = sAIArtStyleParser->GetStyleBlendField(parser,&blendField);

			AIBoolean blend;
			blend = sAIArtStyleParser->IsNormalOpaqueBlend(blendField);
			if(blend)
				log->WritePrintf("Blend field has normal opacity");

			AIBoolean blendVisible;
			result = sAIArtStyleParser->GetBlendFieldVisible(blendField,&blendVisible);
			aisdk::check_ai_error(result);
			if(blendVisible)
				log->Write("Blend field is visible");
			else
			{
				log->Write("Blend field is not visible, so setting it visible");
				result = sAIArtStyleParser->SetBlendFieldVisible(style,blendField,true);
				aisdk::check_ai_error(result);
				result = sAIArtStyleParser->GetBlendFieldVisible(blendField,&blendVisible);
				aisdk::check_ai_error(result);
				if(blendVisible)
					log->Write("Now the blend style is visible");
				else
					log->Write("Blend Style is still not visible");
			}

			
			boolean = sAIArtStyleParser->ContainsPaint(parser);
			if(boolean)
				log->Write("Parser for this art contains paint");
			else
				log->Write("Parser for this art does not contain paint");

			boolean = sAIArtStyleParser->ContainsEffects(parser);
			if(boolean)
				log->Write("Parser for this art contains effects");
			else
				log->Write("Parser for this art does not contain effects");

			boolean = sAIArtStyleParser->ContainsTransparency(parser);
			if(boolean)
				log->Write("Parser for this art contains Transparency");
			else
				log->Write("Parser for this art does not contain Transparency");


			ai::int32 numPaintFields = sAIArtStyleParser->CountPaintFields(parser);
			log->Write("");
			log->WritePrintf("Number of paint fields in art number %d is %d ",i+1,numPaintFields);
			
			//Get each paintfield, and check its fill and stroke.
			for(ai::int32 n = 0; n < numPaintFields ;n++){
				AIParserPaintField parserPaintField;
				result = sAIArtStyleParser->GetNthPaintField(parser,n,&parserPaintField);
				aisdk::check_ai_error(result);

				log->Write("");
				log->WritePrintf("Paint field number : %d",n);		

				if(sAIArtStyleParser->IsFill(parserPaintField)){
					AIFillStyle fill;
					result = sAIArtStyleParser->GetFill(parserPaintField,&fill,NULL);
					aisdk::check_ai_error(result);

					AIBoolean evenOdd;
					result = sAIArtStyleParser->GetEvenOdd(parserPaintField,&evenOdd);
					if(evenOdd)
						log->Write("EvenOdd rule has been applied on this paintfield");
					else
						log->Write("EvenOdd rule has not been applied on this paintfield");

					result = this->GetFillColorKind(fill);
					aisdk::check_ai_error(result);
				}
				if(sAIArtStyleParser->IsStroke(parserPaintField)){
					AIStrokeStyle stroke;
					result = sAIArtStyleParser->GetStroke(parserPaintField,&stroke,NULL);
					aisdk::check_ai_error(result);

					result = this->GetStrokeColorKind(stroke);
					aisdk::check_ai_error(result);
				}
			}
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGraphicStyle::InspectLiveEffects(void)
{
	ASErr result = kNoErr;
	try {	
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerLog::Indent indent;

		AIArtHandle **matchingArt = NULL;
		ai::int32 numMatches;
		result = sAIMatchingArt->GetSelectedArt(&matchingArt,&numMatches);
		aisdk::check_ai_error(result);

		for( ai::int32 i = 0; i < numMatches; i++){
			AIArtHandle art;
			art = (*matchingArt)[i];

			//Get the artsyle of selected art.
			AIArtStyleHandle style;
			result = sAIArtStyle->GetArtStyle(art,&style);
			aisdk::check_ai_error(result);

			//Create a new parser.
			AIStyleParser parser = NULL;
			result = sAIArtStyleParser->NewParser(&parser);
			aisdk::check_ai_error(result);

			result = sAIArtStyleParser->ParseStyle(parser,style);
			aisdk::check_ai_error(result);

			ai::int32 numPaintFields = sAIArtStyleParser->CountPaintFields(parser);
			log->WritePrintf("Number of paint fields in art number %d is %d ",i+1,numPaintFields);

			//For each paintfield, inspect the Live Effects applied.
			for( ai::int32 n = 0; n < numPaintFields; n++){
				log->Write("");
				log->WritePrintf(" For paint field number : %d", n+1);
				AIParserPaintField paintfield;
				result = sAIArtStyleParser->GetNthPaintField(parser,n,&paintfield);
				aisdk::check_ai_error(result);

				ASInt32 effectsCount;
				effectsCount = sAIArtStyleParser->CountEffectsOfPaintField(paintfield);
				log->WritePrintf("Number of Live Effects : %d",effectsCount);

				for(int k = 0; k < effectsCount; k++){
					log->Write("");
					log->WritePrintf("Live Effect #%d : ",k);
					AIParserLiveEffect liveEffect;
					result = sAIArtStyleParser->GetNthEffectOfPaintField(paintfield,k,&liveEffect);
					aisdk::check_ai_error(result);

					AIBoolean visible;
					result = sAIArtStyleParser->GetEffectVisible(liveEffect,&visible);
					aisdk::check_ai_error(result);
					if(visible)
						log->WritePrintf("LiveEffect %d is visible",k);
					else
						log->WritePrintf("LiveEffect %d is not visible",k);

					const char* effectName = NULL;
					ai::int32 minor,major;
					result = sAIArtStyleParser->GetLiveEffectNameAndVersion(liveEffect,&effectName,&minor,&major);
					aisdk::check_ai_error(result);
					log->WritePrintf("LiveEffect for paintfield %d is %s",n+1,effectName);
					
				}
			}
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}


/*
*/
ASErr SnpGraphicStyle::DisplayArtStyle(void)
{
	ASErr result = kNoErr;
	SnippetRunnerLog::Indent indent;
	try {
		//Get the default art style of the document.
		AIArtStyleHandle defaultStyle;
		result = sAIArtStyle->GetDefaultArtStyle(&defaultStyle);
		aisdk::check_ai_error(result);
		
		//Get the default art style name of the document.
		ai::UnicodeString name;
		result = sAIArtStyle->GetDefaultArtStyleName(name);
		aisdk::check_ai_error(result);

		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		log->Write(ai::UnicodeString("Default ArtStyle  = ").append(name));

		ai::int32 count;
		result = sAIArtStyle->CountNamedArtStyles(&count);
		aisdk::check_ai_error(result);
		log->WritePrintf("Total number of ArtStyles available : %d",count);
		
		//Dsiplay the art styles available in a document.
		for(ai::int32 i = 0; i < count; i++){
			AIArtStyleHandle style;
			result = sAIArtStyle->GetNthNamedArtStyle(i,&style);
			aisdk::check_ai_error(result);

			ai::UnicodeString styleName;
			ASBoolean boolean;
			result = sAIArtStyle->GetArtStyleName(style,styleName,&boolean);
			aisdk::check_ai_error(result);
			log->Write("");
			log->WritePrintf(" ArtStyle Number %d : ",i+1);
			log->Write(ai::UnicodeString("is  ").append(styleName));

			AIPathStyle pathstyle;
			AIArtStylePaintData paintData;
			result = sAIArtStyle->GetPaintAttributes(style,&pathstyle,&paintData);
			aisdk::check_ai_error(result);

			AIDocumentHandle doc;
			result = sAIDocument->GetDocument(&doc);
			aisdk::check_ai_error(result);

			ai::int32 namedDocCount;
			result = sAIArtStyle->CountNamedArtStylesFromDocument(&namedDocCount,doc);
			aisdk::check_ai_error(result);

			for( ai::int32 i = 0; i < namedDocCount ; i++){
				AIArtStyleHandle docStyle;
				result = sAIArtStyle->GetNthNamedArtStyleFromDocument(i,&docStyle,doc);
				aisdk::check_ai_error(result);
			}
		}
		/* Sort the named styles alphabetically and display them once. */
		result = sAIArtStyle->SortNamedStyles();
		aisdk::check_ai_error(result);

		//Flag to control the display only once.
		static int displaycount = 0;
		displaycount++; 
		
		if( displaycount == 1){
			log->Write(" Now displaying the art styles after Sorting alphabetically ");
			result = this->DisplayArtStyle();
			aisdk::check_ai_error(result);
			displaycount = 0;
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}


/*
*/
ASErr SnpGraphicStyle::ModifyEffects(void)
{
	ASErr result = kNoErr;
	
	try {
		SnippetRunnerLog::Indent indent;
		AIArtHandle **matchingArt = NULL;
		ai::int32 numMatches;
		result = sAIMatchingArt->GetSelectedArt(&matchingArt,&numMatches);
		aisdk::check_ai_error(result);

		for( ai::int32 i = 0; i < numMatches; i++){
			AIArtHandle art;
			art = (*matchingArt)[i];

			//Get the art style of the selected art.
			AIArtStyleHandle style;
			result = sAIArtStyle->GetArtStyle(art,&style);
			aisdk::check_ai_error(result);

			//Create a new parser.
			AIStyleParser parser = NULL;
			result = sAIArtStyleParser->NewParser(&parser);
			aisdk::check_ai_error(result);

			result = sAIArtStyleParser->ParseStyle(parser,style);
			aisdk::check_ai_error(result);

			ai::int32 preEffectsCount = sAIArtStyleParser->CountPreEffects(parser);

			// Write output to the log.
			SnippetRunnerLog* log = SnippetRunnerLog::Instance();
			log->WritePrintf("Number of Pre-Effects in new parser : %d ",preEffectsCount);

			//Remove any pre-effects applied to the artstyle.
			for ( ASInt32 j = preEffectsCount-1; j >= 0 ; j--)
			{
				AIParserLiveEffect effect;
				result = sAIArtStyleParser->GetNthPreEffect( parser, j, &effect );
				aisdk::check_ai_error(result);
				
				result = sAIArtStyleParser->RemovePreEffect( parser, effect, true );
				aisdk::check_ai_error(result);			
			}
			ai::int32 postEffectsCount = sAIArtStyleParser->CountPostEffects(parser) ;

			if (postEffectsCount > 0 )
				log->WritePrintf("Number of Post-Effects in new parser : %d ",postEffectsCount);

			//Check whether the selected item contains paint.
			if(sAIArtStyleParser->ContainsPaint(parser)){
				log->Write("Contains paint");

				ai::int32 numPaintFields = sAIArtStyleParser->CountPaintFields(parser);
				log->WritePrintf("Number of paint fields in art number %d is %d ",i+1,numPaintFields);
				for(ai::int32 n = 0; n < numPaintFields ;n++){
					AIParserPaintField parserPaintField;
					result = sAIArtStyleParser->GetNthPaintField(parser,n,&parserPaintField);
					aisdk::check_ai_error(result);

					//Check whether the selected paintfield is a fill.
					if(sAIArtStyleParser->IsFill(parserPaintField)){
						log->WritePrintf("Selected item Fill");

						AIFillStyle fill;
						result = sAIArtStyleParser->GetFill(parserPaintField,&fill,NULL);
						aisdk::check_ai_error(result);

						AILiveEffectHandle liveEffect = NULL;						
						sAIArtStyleParser->GetPaintLiveEffectInfo(parserPaintField, &liveEffect, nullptr);
						if (liveEffect){
							//Prompt the user to modify the applied Effects.
							result = sAIArtStyleParser->EditPaintEffectParameters(style, parserPaintField);						
							aisdk::check_ai_error(result);
						}
						else
							log->Write("No Live Effect applied to fill");

					}
					//Check whether the selected paintfield is a stroke.
					if(sAIArtStyleParser->IsStroke(parserPaintField)){
						log->WritePrintf("Selected item Stroke");
						AIStrokeStyle stroke;
						result = sAIArtStyleParser->GetStroke(parserPaintField,&stroke,NULL);
						aisdk::check_ai_error(result);

						AILiveEffectHandle liveEffect = NULL;						
						sAIArtStyleParser->GetPaintLiveEffectInfo(parserPaintField, &liveEffect, nullptr);
						if (liveEffect){
							//Prompt the user to modify the applied LiveEffects.
							result = sAIArtStyleParser->EditPaintEffectParameters(style, parserPaintField);						
							aisdk::check_ai_error(result);
						}
						else
							log->Write("No Live Effect applied to stroke");
					}
				}
			}				
			// Get the Focus Stroke.
			AIParserPaintField focalStroke;
			result = sAIArtStyleParser->GetFocusStroke( parser, &focalStroke );
			aisdk::check_ai_error(result);
			
			//Get the Focus Fill.
			AIParserPaintField focalFill;
			result = sAIArtStyleParser->GetFocusFill( parser, &focalFill );
			aisdk::check_ai_error(result);
			

			//Get and modify Stroke and Fill values.
			AIStrokeStyle strokeStyle;
			result = sAIArtStyleParser->GetStroke( focalStroke, &strokeStyle, NULL );
			aisdk::check_ai_error(result);
			
			AIFillStyle fillStyle;
			fillStyle.color = strokeStyle.color;
			fillStyle.overprint = strokeStyle.overprint;

			result = sAIArtStyleParser->SetFill( focalFill, &fillStyle, NULL );
			aisdk::check_ai_error(result);

			ai::int32 strokeEffectsCount = sAIArtStyleParser->CountEffectsOfPaintField( focalStroke );
			log->WritePrintf("Number of Effects of Focal Stroke : %d ",strokeEffectsCount);

			ai::int32 fillEffectsCount = sAIArtStyleParser->CountEffectsOfPaintField( focalFill );
			log->WritePrintf("Number of Effects of Focal Fill : %d ",fillEffectsCount);

		}				
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}


/*
*/
ASErr SnpGraphicStyle::ModifyLiveEffects(void)
{
	ASErr result = kNoErr;
	try {	
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerLog::Indent indent;

		AIArtHandle **matchingArt = NULL;
		ai::int32 numMatches;
		result = sAIMatchingArt->GetSelectedArt(&matchingArt,&numMatches);
		aisdk::check_ai_error(result);

		for( ai::int32 i = 0; i < numMatches; i++){
			AIArtHandle art;
			art = (*matchingArt)[i];

			//Get the artsyle of selected art.
			AIArtStyleHandle style;
			result = sAIArtStyle->GetArtStyle(art,&style);
			aisdk::check_ai_error(result);

			//Create a new parser.
			AIStyleParser parser = NULL;
			result = sAIArtStyleParser->NewParser(&parser);
			aisdk::check_ai_error(result);

			result = sAIArtStyleParser->ParseStyle(parser,style);
			aisdk::check_ai_error(result);

			ai::int32 numPaintFields = sAIArtStyleParser->CountPaintFields(parser);
			log->WritePrintf("Number of paint fields in art number %d is %d ",i+1,numPaintFields);

			//For each paintfield, inspect the Live Effects applied.
			for( int n = 0; n < numPaintFields; n++){
				log->WritePrintf(" For paint field number : %d", n+1);
				AIParserPaintField paintfield;
				result = sAIArtStyleParser->GetNthPaintField(parser,n,&paintfield);
				aisdk::check_ai_error(result);

				ASInt32 effectsCount;
				effectsCount = sAIArtStyleParser->CountEffectsOfPaintField(paintfield);
				log->WritePrintf("Number of Live Effects : %d",effectsCount);

				for(int k = 0; k < effectsCount; k++){
					AIParserLiveEffect liveEffect;
					result = sAIArtStyleParser->GetNthEffectOfPaintField(paintfield,k,&liveEffect);
					aisdk::check_ai_error(result);

					AIBoolean visible;
					result = sAIArtStyleParser->GetEffectVisible(liveEffect,&visible);
					aisdk::check_ai_error(result);
					if(visible)
						log->WritePrintf("LiveEffect %d is visible",k);
					else
						log->WritePrintf("LiveEffect %d is not visible",k);

					const char* effectName = NULL;
					ai::int32 minor,major;
					result = sAIArtStyleParser->GetLiveEffectNameAndVersion(liveEffect,&effectName,&minor,&major);
					aisdk::check_ai_error(result);
					log->WritePrintf("LiveEffect for paintfield %d is %s",n+1,effectName);

					//For each liveEffect applied, prompt the user with an editing dialog to edit the parameters of each liveEffect.
					result = sAIArtStyleParser->EditEffectParameters(style,liveEffect);
					aisdk::check_ai_error(result);
				}
			}
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGraphicStyle::RemoveArtStyleByChoice(void)
{
	ASErr result = kNoErr;
	SnippetRunnerLog::Indent indent;
	try {
		//Get the default art style of the document.
		AIArtStyleHandle defaultStyle;
		result = sAIArtStyle->GetDefaultArtStyle(&defaultStyle);
		aisdk::check_ai_error(result);
		
		ai::int32 count;
		result = sAIArtStyle->CountNamedArtStyles(&count);
		aisdk::check_ai_error(result);

		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		log->WritePrintf("Total number of ArtStyles available : %d",count);

		std::vector<ai::UnicodeString> names;
		AIArtStyleHandle style;
		
		//Display the art styles available in a document.
		for(ai::int32 i = 0; i < count; i++){
			
			result = sAIArtStyle->GetNthNamedArtStyle(i,&style);
			aisdk::check_ai_error(result);

			ai::UnicodeString styleName;
			ASBoolean boolean;
			result = sAIArtStyle->GetArtStyleName(style,styleName,&boolean);
			aisdk::check_ai_error(result);

			names.push_back(styleName);				
							
		}
		
		//Prompt the user for the art style to be deleted.
		SnippetRunnerParameter* parameter = SnippetRunnerParameter::Instance();
		ASInt32 choice = parameter->GetChoice(ai::UnicodeString("Select Graphic Style to delete:"), names,1);
		ai::UnicodeString chosenPatternName = names[choice];

		if(choice){
			//Get the name of the selected art style.
			ai::UnicodeString styleNametoDelete;
			result = sAIArtStyle->GetArtStyleByName(&style,chosenPatternName,false);
			aisdk::check_ai_error(result);

			//Remove the selected art style.
			AIArtStyleHandle anonStyle;
			result = sAIArtStyle->RemoveNamedStyle(style,&anonStyle);
			aisdk::check_ai_error(result);
			
			result = this->DisplayArtStyle();
			aisdk::check_ai_error(result);
		}
		else
			log->Write("Default art style cannot be deleted");
		
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}


/*
*/
ASErr SnpGraphicStyle::RemoveEffectByChoice(void)
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog::Indent indent;
		AIArtHandle **matchingArt = NULL;
		ai::int32 numMatches;
		result = sAIMatchingArt->GetSelectedArt(&matchingArt,&numMatches);
		aisdk::check_ai_error(result);

		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();

		for( ai::int32 i = 0; i < numMatches; i++){
			AIArtHandle art;
			art = (*matchingArt)[i];

			if( i == 1 )
			{
				//Get the art style of the selected art object.
				AIArtStyleHandle style;
				result = sAIArtStyle->GetArtStyle(art,&style);
				aisdk::check_ai_error(result);

				//Create a new parser.
				AIStyleParser parser = NULL;
				result = sAIArtStyleParser->NewParser(&parser);
				aisdk::check_ai_error(result);

				result = sAIArtStyleParser->ParseStyle(parser,style);
				aisdk::check_ai_error(result);

				ai::int32 numPaintFields = sAIArtStyleParser->CountPaintFields(parser);
				log->WritePrintf("Number of paint fields in art number %d is %d ",i+1,numPaintFields);
				ASInt32 choice = SnippetRunnerParameter::Instance()->GetInt("Select Paint Number to delete:", 0,0,numPaintFields-1);

				AIParserPaintField paintfield;
				result = sAIArtStyleParser->GetNthPaintField(parser,choice,&paintfield);
				aisdk::check_ai_error(result);

				//Remove the effects of the selected paintfield.
				ASInt32 numEffects;
				numEffects = sAIArtStyleParser->CountEffectsOfPaintField(paintfield);
				if(numEffects){
					ASInt32 choiceEffects = SnippetRunnerParameter::Instance()->GetInt("Select Effect to delete:", 0,0,numEffects-1);
					AIParserLiveEffect liveeffect ;
					result = sAIArtStyleParser->GetNthEffectOfPaintField(paintfield,choiceEffects,&liveeffect);
					aisdk::check_ai_error(result);

					const char* name;
					ai::int32 minor,major;
					result = sAIArtStyleParser->GetLiveEffectNameAndVersion(liveeffect,&name,&minor,&major);
					aisdk::check_ai_error(result);
					log->WritePrintf("Selected %s live effect to be deleted",name);
					result = sAIArtStyleParser->RemoveEffectOfPaintField(parser,paintfield,liveeffect,true);
					aisdk::check_ai_error(result);
				}	
				else
					log->Write("Selected paint contains no Live Effect to be deleted");

				//Apply the new style to the art.
				AIArtStyleHandle newStyle;
				result = sAIArtStyleParser->CreateNewStyle(parser,&newStyle);
				aisdk::check_ai_error(result);

				result = sAIArtStyle->SetArtStyle(art,newStyle);
				aisdk::check_ai_error(result);
			}
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpGraphicStyle::RemoveAllEffects(void)
{
	ASErr result = kNoErr;
	try {
		SnippetRunnerLog::Indent indent;
		AIArtHandle **matchingArt = NULL;
		ai::int32 numMatches;
		result = sAIMatchingArt->GetSelectedArt(&matchingArt,&numMatches);
		aisdk::check_ai_error(result);

		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();

		for( ai::int32 i = 0; i < numMatches; i++){
			AIArtHandle art;
			art = (*matchingArt)[i];

			//Get the art style of the selected art object.
			AIArtStyleHandle style;
			result = sAIArtStyle->GetArtStyle(art,&style);
			aisdk::check_ai_error(result);

			//Create a new parser.
			AIStyleParser parser = NULL;
			result = sAIArtStyleParser->NewParser(&parser);
			aisdk::check_ai_error(result);

			result = sAIArtStyleParser->ParseStyle(parser,style);
			aisdk::check_ai_error(result);

			long numPaintFields;
			numPaintFields = sAIArtStyleParser->CountPaintFields(parser);
			log->WritePrintf("Number of paint fields in art number %d is %d ",i+1,numPaintFields);
			
			for(int j=numPaintFields-1; j >= 0 ; j--){
				AIParserPaintField paintfield;
				result = sAIArtStyleParser->GetNthPaintField(parser,j,&paintfield);
				aisdk::check_ai_error(result);

				//Remove all the effects of the selected paintfield.
				ASInt32 numEffects;
				numEffects = sAIArtStyleParser->CountEffectsOfPaintField(paintfield);
				if(numEffects){
					result = sAIArtStyleParser->RemoveAllEffects(parser);
					aisdk::check_ai_error(result);
				}	
				result = sAIArtStyleParser->RemovePaintField(parser,paintfield,true);
				aisdk::check_ai_error(result);

				//Apply the new style to the art.
				AIArtStyleHandle newStyle;
				result = sAIArtStyleParser->CreateNewStyle(parser,&newStyle);
				aisdk::check_ai_error(result);

				result = sAIArtStyle->SetArtStyle(art,newStyle);
				aisdk::check_ai_error(result);
			}
		
			numPaintFields = sAIArtStyleParser->CountPaintFields(parser);
			log->WritePrintf("Number of paint fields after deleting paintfield %d is %d ",i,numPaintFields);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}


// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpGraphicStyle available to the SnippetRunner framework.
 */
class _SnpRunnableGraphicStyle : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableGraphicStyle () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableGraphicStyle () {}

		/* Returns name of snippet.
		*/
		std::string GetName() const;

		/* Returns a description of what the snippet does.
		*/
		std::string	GetDescription() const;

		/* Returns operations supported by this snippet.
		*/
		Operations GetOperations() const;

		/* Returns name of this snippet's default operation - must
			be one of the operation names returned by GetOperations.
		*/
		std::string GetDefaultOperationName() const;

		/** Returns the categories a snippet belongs to.
			@return default categories.
		*/
		std::vector<std::string> GetCategories() const;

		/* Returns true if the snippet can run.
			@param runnableContext see ISnpRunnableContext for documentation.
			@return true if snippet can run, false otherwise
		 */
		ASBoolean			CanRun(SnpRunnable::Context& runnableContext);

		/* Runs the snippet.
			@param runnableContext see ISnpRunnableContext for documentation.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr		Run(SnpRunnable::Context& runnableContext);
};

/*
*/
std::string _SnpRunnableGraphicStyle::GetName() const
{
	return "Graphic Styles";
}

/*
*/
std::string _SnpRunnableGraphicStyle::GetDescription() const 
{
	return "Provides GraphicStyle functions";
}

/*
*/
SnpRunnable::Operations _SnpRunnableGraphicStyle::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("NewStyleFromExisting", "Selected art object", kSnpRunAnyContext));
	operations.push_back(Operation("InspectFillStrokeAndBlend", "Selected art object", kSnpRunAnyContext));
	operations.push_back(Operation("InspectLiveEffects", "Selected art object", kSnpRunAnyContext));
	operations.push_back(Operation("DisplayArtStyle", "document", kSnpRunAnyContext));
	operations.push_back(Operation("ModifyEffects", "Selected art object", kSnpRunAnyContext));
	operations.push_back(Operation("ModifyLiveEffects", "Selected art object", kSnpRunAnyContext));
	operations.push_back(Operation("RemoveArtStyleByChoice", "document" , kSnpRunAnyContext));
	operations.push_back(Operation("RemoveEffectByChoice", "Selected art object" , kSnpRunAnyContext));
	operations.push_back(Operation("RemoveAllEffects", "Selected art object", kSnpRunAnyContext));
	return operations;
}

/*
*/
std::string _SnpRunnableGraphicStyle::GetDefaultOperationName() const 
{
	return "DisplayArtStyle";
}

/*
*/
std::vector<std::string> _SnpRunnableGraphicStyle::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Color Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableGraphicStyle::CanRun(SnpRunnable::Context& runnableContext)
{	
	SnpSelectionHelper selectionHelper;

	if (("DisplayArtStyle" == runnableContext.GetOperation().GetName()) ||
		("RemoveArtStyleByChoice" == runnableContext.GetOperation().GetName()) )
		return selectionHelper.IsDocumentSelected();

	if (("NewStyleFromExisting" == runnableContext.GetOperation().GetName()) ||
		("InspectFillStrokeAndBlend" == runnableContext.GetOperation().GetName()) ||
		("InspectLiveEffects" == runnableContext.GetOperation().GetName()) ||
		("ModifyEffects" == runnableContext.GetOperation().GetName()) ||
		("ModifyLiveEffects" == runnableContext.GetOperation().GetName())||
		("RemoveEffectByChoice" == runnableContext.GetOperation().GetName())||
		("RemoveAllEffects" == runnableContext.GetOperation().GetName())){

		if (selectionHelper.IsArtSelected(kPathArt)||
			selectionHelper.IsArtSelected(kPlacedArt) ||
			selectionHelper.IsArtSelected(kRasterArt) ||
			selectionHelper.IsArtSelected(kGroupArt))
					return true;
				else 
					return false;
	}
	else
		return false;
	
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableGraphicStyle::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpGraphicStyle instance;
	if ("NewStyleFromExisting" == runnableContext.GetOperation().GetName()) {
		result = instance.NewStyleFromExisting();
	}
	else if ("InspectFillStrokeAndBlend" == runnableContext.GetOperation().GetName()) {
		result = instance.InspectFillStrokeAndBlend();
	}
	else if ("InspectLiveEffects" == runnableContext.GetOperation().GetName()) {
		result = instance.InspectLiveEffects();
	}
	else if ("DisplayArtStyle" == runnableContext.GetOperation().GetName()) {
		result = instance.DisplayArtStyle();
	}	
	else if ("ModifyEffects" == runnableContext.GetOperation().GetName()) {
		result = instance.ModifyEffects();
	}
	else if ("ModifyLiveEffects" == runnableContext.GetOperation().GetName()) {
		result = instance.ModifyLiveEffects();
	}
	else if ("RemoveArtStyleByChoice" == runnableContext.GetOperation().GetName()) {
		result = instance.RemoveArtStyleByChoice();
	}
	else if ("RemoveEffectByChoice" == runnableContext.GetOperation().GetName()) {
		result = instance.RemoveEffectByChoice();
	}
	else if ("RemoveAllEffects" == runnableContext.GetOperation().GetName()) {
		result = instance.RemoveAllEffects();
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableGraphicStyle instance_SnpRunnableGraphicStyle;

// End SnpGraphicStyle.cpp

