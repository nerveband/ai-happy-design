//========================================================================================
//  
//  $File: //ai/mainline/devtech/sdk/public/samplecode/Annotator/Resources/Annotator.r $
//
//  $Revision: #5 $
//
//  Copyright 2008 Adobe Systems Incorporated. All rights reserved.
//  
//  NOTICE:  Adobe permits you to use, modify, and distribute this file in accordance 
//  with the terms of the Adobe license agreement accompanying it.  If you have received
//  this file from a source other than Adobe, then your use, modification, or 
//  distribution of it requires the prior written permission of Adobe.
//  
//========================================================================================


/*
 * Two Resource files are required for Mac one for the standard definition resources and one
 * for HiDPI resources. The resource entries for both are required regardless of whether
 * you will be providing HiDPI icons for your tools. If you do not provide the two resource 
 * files you will receive [CANT] Errors on startup.
 * 
 * The resource files should follow the naming convention:
 * PluginName.r -> eg. Annotator.r - for standard image definitions
 * PluginName2x.r -> eg. Annotator2x.r - for HiDPI image definitions
 *
 *
 * Changes in resource handling:
 * In 10.8 SDK, Apple deprecated all .r/.rsrc resource reading APIs. Now all resources files
 * will need to be copied to <YourPlugin.aip>/Contents/Resources/<extension>/ as post build step.
 *
 * In this sample there are 4 resources used:
 * - AnnotatorTool.png
 * - AnnotatorTool@2x.png
 * - AnnotatorTool@3to2x.png
 * - AnnotatorTool.xml
 *
 * The resources are added in the Run Script build phase of the target. Each resource is added 
 * to a folder named after its extension. 
 * e.g All png folders should be copied to <YourPlugin.aip>/Contents/Resources/png/.
 * 
 * The script for Annotator plugin looks like this: 
 *
 *    -- This creates directories inside the plugin package. There should be a separate directory 
 *    for each type of resource. For Annotator plugin these directories will be png and xml.
 *    mkdir -p "${CONFIGURATION_BUILD_DIR}/${PRODUCT_NAME}.aip/Contents/Resources/png"
 *    mkdir -p "${CONFIGURATION_BUILD_DIR}/${PRODUCT_NAME}.aip/Contents/Resources/xml"
 * 
 *    -- This creates a txt directory for IDToFile.txt which is used to map IDs to resource file names.
 *    mkdir -p "${CONFIGURATION_BUILD_DIR}/${PRODUCT_NAME}.aip/Contents/Resources/txt"
 *
 *    -- Copy resources to Annotator.aip/Contents/Resources/<extension>/
 *   cp -f "./Resources/raw/AnnotatorTool.png" "${CONFIGURATION_BUILD_DIR}/${PRODUCT_NAME}.aip/Contents/Resources/png/AnnotatorTool.png"
 *   cp -f "./Resources/raw/AnnotatorTool@2x.png" "${CONFIGURATION_BUILD_DIR}/${PRODUCT_NAME}.aip/Contents/Resources/png/AnnotatorTool@2x.png"
 *   cp -f "./Resources/raw/AnnotatorTool@3to2x.png" "${CONFIGURATION_BUILD_DIR}/${PRODUCT_NAME}.aip/Contents/Resources/png/AnnotatorTool@3to2x.png"
 *   cp -f "./Resources/raw/AnnotatorTool.xml" "${CONFIGURATION_BUILD_DIR}/${PRODUCT_NAME}.aip/Contents/Resources/xml/AnnotatorTool.xml"
 *
 *    -- Make sure IDToFile.txt is copied over as well.
 *    cp -f "./Resources/raw/IDToFile.txt" "${CONFIGURATION_BUILD_DIR}/${PRODUCT_NAME}.aip/Contents/Resources/txt/IDToFile.txt"
 *
 *    -- Exit the script
 *    exit
 *
 * 
 * IDToFile.txt
 * --------------
 * You will see in this sample that IDToFile.txt has 2 entries:
 * 
 * 1. 16151.png		AnnotatorTool.png
 * this points to the standard DPI image.
 * 
 * 2. 16151.xml		AnnotatorTool.xml
 * this points to an xml file that defines the cursor hotspot for the tools image
 *
 * HiDPI images - the name of the resource must be the standard images name 
 * with @2x or @3to2x added on. You do not need to define these explicitly in IDToFile.txt as 
 * they will be deduced from the images you add in Run Script build phase.
 *
 */

