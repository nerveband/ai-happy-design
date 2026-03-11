//========================================================================================
//  
//  $File: //ai/mainline/devtech/sdk/public/samplecode/MarkedObjects/Resources/MarkedObjects.r $
//
//  $Revision: #6 $
//
//  Copyright 1987 Adobe Systems Incorporated. All rights reserved.
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
 * PluginName.r -> eg. MarkedObjects.r - for standard image definitions
 * PluginName2x.r -> eg. MarkedObjects2x.r - for HiDPI image definitions
 *
 *
 * Changes in resource handling:
 * In 10.8 SDK, Apple deprecated all .r/.rsrc resource reading APIs. Now all resources files
 * will need to be copied to <YourPlugin.aip>/Contents/Resources/<extension>/ as post build step.
 *
 * The resources are added in the Run Script build phase of the target. Each resource is added 
 * to a folder named after its extension. 
 * e.g All png folders should be copied to <YourPlugin.aip>/Contents/Resources/png/.
 *
 * 
 * IDToFile.txt
 * --------------
 * You will see in this sample that IDToFile.txt has 2 entries:
 * 
 *      16050.png		MarkedObjectsTool.png
 *      16051.png		MarkedObjectsIconRollover.png
 * this points to the standard DPI images.
 *
 * HiDPI images - the name of the resource must be the standard images name 
 * with @2x or @3to2x added on. You do not need to define these explicitly in IDToFile.txt as 
 * they will be deduced from the images you add in Run Script build phase.
 *
 */

#include "MarkedObjectID.h"

#include <Carbon/Carbon.r>
