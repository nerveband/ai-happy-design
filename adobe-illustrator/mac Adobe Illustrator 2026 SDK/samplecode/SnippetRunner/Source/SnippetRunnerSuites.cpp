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
#include "SnippetRunnerSuites.h"
#include "ATETextSuitesImportHelper.h"

// Defines pointers for suites used by code snippets.
extern "C" 
{
	AIDocumentListSuite*			sAIDocumentList = NULL;
	AIDocumentSuite*				sAIDocument = NULL;
	AIVectorizeSuite*				sAIVectorize = NULL;
	AIArtSuite*						sAIArt = NULL;
	AICountedObjectSuite*			sAICountedObject = NULL;
	AIDictionarySuite*				sAIDictionary = NULL;
	AIArtSetSuite*					sAIArtSet = NULL;
	AIPlacedSuite*					sAIPlaced = NULL;
	AISymbolSuite*					sAISymbol = NULL;
	AIRasterizeSuite*				sAIRasterize = NULL;
	AILayerSuite*					sAILayer = NULL;
	AIPluginGroupSuite*				sAIPluginGroup = NULL;
	AISwatchListSuite*				sAISwatchList = NULL;
	AISwatchGroupSuite*				sAISwatchGroup = NULL;
	AICustomColorSuite*				sAICustomColor = NULL;
	AIColorHarmonySuite*			sAIColorHarmony = NULL;
	AIPathStyleSuite*				sAIPathStyle = NULL;
	AIMatchingArtSuite*				sAIMatchingArt = NULL;
	ASUserInteractionSuite*			sASUserInteraction = NULL;
	AIActionManagerSuite*			sAIActionManager = NULL;
	AIRealMathSuite*				sAIRealMath = NULL;
	AIColorConversionSuite*			sAIColorConversion = NULL;
	AITextFrameSuite*				sAITextFrame = NULL;
	AILegacyTextConversionSuite*	sAILegacyTextConversion = NULL;
	AIToolSuite*					sAITool = NULL;
	AIATEPaintSuite*				sAIATEPaint = NULL;
	AIFontSuite*					sAIFont = NULL;
	AIHardSoftSuite*				sAIHardSoft = NULL;
	AIPathSuite*					sAIPath = NULL;
	AIATETextUtilSuite*				sAIATETextUtil = NULL;
	AIATECurrentTextFeaturesSuite*	sAIATECurrentTextFeatures = NULL;
	AITextFrameHitSuite*			sAITextFrameHit = NULL;
	AIHitTestSuite*					sAIHitTest = NULL;
	AIShapeConstructionSuite*		sAIShapeConstruction = NULL;
	//AICropAreaSuite*				sAICropArea = NULL;
	//AICropAreaRangeSuite*			sAICropAreaRange = NULL;
	AIArtboardSuite*				sAIArtboard = NULL;
	AIArtboardRangeSuite*			sAIArtboardRange = NULL;
	AIGradientSuite*				sAIGradient = NULL;
	AIArtStyleSuite*				sAIArtStyle = NULL;
	AIFXGFileFormatSuite*			sAIFXGFileFormat = NULL;
	AIXMLDocumentSuite*				sAIXMLDocument = NULL;
	AIDictionaryIteratorSuite*		sAIDictionaryIterator = NULL;
	AIPatternSuite*					sAIPattern = NULL;
	AIArtStyleParserSuite*			sAIArtStyleParser = NULL;
	AIEntrySuite*					sAIEntry = NULL;
	AIXMLNodeSuite*					sAIXMLNode = NULL;
	AIXMLNodeListSuite*				sAIXMLNodeList = NULL;
	AIXMLNamedNodeMapSuite*			sAIXMLNamedNodeMap = NULL;
	AIXMLElementSuite*				sAIXMLElement = NULL;
	AIDrawArtSuite*					sAIDrawArt = NULL;
	AIAssertionSuite*				sAIAssertion = NULL;
	EXTERN_TEXT_SUITES
	// Add a line above for each additional suite your snippet needs.
	// NOTE: Also add a line to gPostStartupSuites and edit SnippetRunnerSuites.h to declare your suite pointer.
};

// Identifies suites used by code snippets.
// Acquired on application startup (kAIApplicationStartupNotifier) and released on plug-in shutdown.
ImportSuite		gPostStartupSuites[] = {

	kAIDocumentListSuite, kAIDocumentListSuiteVersion, &sAIDocumentList,
	kAIDocumentSuite, kAIDocumentSuiteVersion, &sAIDocument,
	kAIVectorizeSuite, kAIVectorizeVersion, &sAIVectorize,
	kAIArtSuite, kAIArtSuiteVersion, &sAIArt,
	kAICountedObjectSuite, kAICountedObjectSuiteVersion, &sAICountedObject,
	kAIDictionarySuite, kAIDictionarySuiteVersion, &sAIDictionary,
	kAIArtSetSuite, kAIArtSetSuiteVersion, &sAIArtSet,
	kAIPlacedSuite, kAIPlacedSuiteVersion, &sAIPlaced,
	kAISymbolSuite, kAISymbolSuiteVersion, &sAISymbol,
	kAIRasterizeSuite, kAIRasterizeSuiteVersion, &sAIRasterize,
	kAILayerSuite,kAILayerSuiteVersion, &sAILayer,
	kAIPluginGroupSuite, kAIPluginGroupSuiteVersion, &sAIPluginGroup,
	kAISwatchListSuite, kAISwatchListSuiteVersion, &sAISwatchList,
	kAISwatchGroupSuite, kAISwatchGroupSuiteVersion, &sAISwatchGroup,
	kAICustomColorSuite, kAICustomColorSuiteVersion, &sAICustomColor,
	kAIColorHarmonySuite, kAIColorHarmonySuiteVersion, &sAIColorHarmony,
	kAIPathStyleSuite, kAIPathStyleSuiteVersion, &sAIPathStyle,
	kAIMatchingArtSuite, kAIMatchingArtSuiteVersion, &sAIMatchingArt,
	kASUserInteractionSuite, kASUserInteractionSuiteVersion, &sASUserInteraction,
	kAIActionManagerSuite, kAIActionManagerVersion, &sAIActionManager,
	kAIRealMathSuite, kAIRealMathSuiteVersion, &sAIRealMath,
	kAIColorConversionSuite, kAIColorConversionSuiteVersion, &sAIColorConversion,
	kAITextFrameSuite, kAITextFrameSuiteVersion, &sAITextFrame,
	kAILegacyTextConversionSuite, kAILegacyTextConversionSuiteVersion, &sAILegacyTextConversion,
	kAIToolSuite, kAIToolSuiteVersion, &sAITool,
	kAIATEPaintSuite, kAIATEPaintSuiteVersion, &sAIATEPaint,
	kAIFontSuite, kAIFontSuiteVersion, &sAIFont,
	kAIHardSoftSuite, kAIHardSoftSuiteVersion, &sAIHardSoft,
	kAIPathSuite, kAIPathSuiteVersion, &sAIPath,
	kAIATETextUtilSuite, kAIATETextUtilSuiteVersion, &sAIATETextUtil,
	kAIATECurrentTextFeaturesSuite, kAIATECurrentTextFeaturesSuiteVersion, &sAIATECurrentTextFeatures,
	kAITextFrameHitSuite, kAITextFrameHitSuiteVersion, &sAITextFrameHit,
	kAIHitTestSuite, kAIHitTestSuiteVersion, &sAIHitTest,
	kAIShapeConstructionSuite, kAIShapeConstructionSuiteVersion, &sAIShapeConstruction,
	//kAICropAreaSuite, kAICropAreaSuiteVersion, &sAICropArea,
	//kAICropAreaRangeSuite, kAICropAreaRangeSuiteVersion, &sAICropAreaRange,
	kAIArtboardSuite, kAIArtboardSuiteVersion, &sAIArtboard,
	kAIArtboardRangeSuite, kAIArtboardRangeSuiteVersion, &sAIArtboardRange,
	kAIGradientSuite, kAIGradientSuiteVersion, &sAIGradient,
	kAIArtStyleSuite, kAIArtStyleSuiteVersion, &sAIArtStyle,
	//TODO Illustrator suite removed for build 439.1 due once it is fixed this needs to be uncommented
	//kAIFXGFileFormatSuite, kAIFXGFileFormatSuiteVersion, &sAIFXGFileFormat,
	kAIXMLDocumentSuite, kAIXMLDocumentSuiteVersion, &sAIXMLDocument,
	kAIDictionaryIteratorSuite, kAIDictionaryIteratorSuiteVersion, &sAIDictionaryIterator,
	kAIPatternSuite, kAIPatternSuiteVersion, &sAIPattern,
	kAIArtStyleParserSuite, kAIArtStyleParserSuiteVersion, &sAIArtStyleParser,
	kAIEntrySuite, kAIEntrySuiteVersion, &sAIEntry,
	kAIXMLNodeSuite, kAIXMLNodeSuiteVersion, &sAIXMLNode,
	kAIXMLNodeListSuite, kAIXMLNodeListSuiteVersion, &sAIXMLNodeList,
	kAIXMLNamedNodeMapSuite, kAIXMLNamedNodeMapSuiteVersion, &sAIXMLNamedNodeMap,
	kAIXMLElementSuite, kAIXMLElementSuiteVersion, &sAIXMLElement,
	kAICSXSExtensionSuite, kAICSXSExtensionSuiteVersion, &sAICSXSExtension,
	kAIDrawArtSuite, kAIDrawArtSuiteVersion, &sAIDrawArt,
	// Add a line above for each additional suite your snippet needs.
	IMPORT_TEXT_SUITES
	nullptr, 0, nullptr
};

// Defines pointers for suites used by the SnippetRunner framework.
extern "C" 
{
	AIUnicodeStringSuite*		sAIUnicodeString = NULL;
	SPBlocksSuite*				sSPBlocks = NULL;
	AIMenuSuite*				sAIMenu = NULL;
	AIStringFormatUtilsSuite*	sAIStringFormatUtils = NULL;
	AIPreferenceSuite*			sAIPreference = NULL;
	AITimerSuite*				sAITimer = NULL;
	AIDataFilterSuite*			sAIDataFilter = NULL;
	AICSXSExtensionSuite*		sAICSXSExtension = NULL;
};

// Identifies the suites acquired on plug-in startup (Plugin::StartupPlugin) and released on plug-in shutdown.
ImportSuite gImportSuites[] = 
{
	kAIUnicodeStringSuite, kAIUnicodeStringSuiteVersion, &sAIUnicodeString,
	kAIArtboardSuite, kAIArtboardSuiteVersion, &sAIArtboard,
	kSPBlocksSuite, kSPBlocksSuiteVersion, &sSPBlocks,
	kAIMenuSuite, kAIMenuSuiteVersion, &sAIMenu,
	kAIStringFormatUtilsSuite, kAIStringFormatUtilsVersion, &sAIStringFormatUtils,
	kAIPreferenceSuite, kAIPreferenceSuiteVersion, &sAIPreference,
	kAITimerSuite, kAITimerSuiteVersion, &sAITimer,
	kAIDataFilterSuite, kAIDataFilterSuiteVersion, &sAIDataFilter,
	kAIAssertionSuite, kAIAssertionSuiteVersion, &sAIAssertion,
	nullptr, 0, nullptr
};

// End SnippetRunnerSuites.cpp
