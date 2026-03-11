//========================================================================================
//  
//  $File$
//
//  $Revision$
//
//  Copyright 2008 Adobe Systems Incorporated. All rights reserved.
//  
//  NOTICE:  Adobe permits you to use, modify, and distribute this file in accordance 
//  with the terms of the Adobe license agreement accompanying it.  If you have received
//  this file from a source other than Adobe, then your use, modification, or 
//  distribution of it requires the prior written permission of Adobe.
//  
//========================================================================================

#include "IllustratorSDK.h"
#include "AnnotatorPlugin.h"

AnnotatorPlugin *gPlugin = NULL;

Plugin *AllocatePlugin(SPPluginRef pluginRef)
{
	return new AnnotatorPlugin(pluginRef);
}

/*
*/
void FixupReload(Plugin* plugin)
{
	AnnotatorPlugin::FixupVTable((AnnotatorPlugin*) plugin);
}

/*
*/
ASErr AnnotatorPlugin::SetGlobal(Plugin* plugin)
{
	gPlugin = (AnnotatorPlugin*) plugin;
	return kNoErr;
}

/*
*/
AnnotatorPlugin::AnnotatorPlugin(SPPluginRef pluginRef) :
	Plugin(pluginRef), fToolHandle(NULL), fAboutPluginMenu(NULL), 
	fAnnotatorHandle(NULL), fNotifySelectionChanged(NULL),
	fAnnotator(NULL),
    fResourceManagerHandle(NULL),
	fShutdownApplicationNotifier(NULL)
{	
	strncpy(fPluginName, kAnnotatorPluginName, kMaxStringLength);
}

/*
*/
ASErr AnnotatorPlugin::StartupPlugin(SPInterfaceMessage* message)
{
	ASErr result = kNoErr;
	try 
	{
		result = Plugin::StartupPlugin(message);
		aisdk::check_ai_error(result);
		// Add About Plugins menu item for this plug-in.
		SDKAboutPluginsHelper aboutPluginsHelper;
		result = aboutPluginsHelper.AddAboutPluginsMenuItem(message, 
					kSDKDefAboutSDKCompanyPluginsGroupName, 
					ai::UnicodeString(kSDKDefAboutSDKCompanyPluginsGroupNameString), 
					"Annotator...", 
					&fAboutPluginMenu);
		aisdk::check_ai_error(result);

		// Add tool.
		result = this->AddTool(message);
		aisdk::check_ai_error(result);

		// Add annotator.
		result = this->AddAnnotator(message);
		aisdk::check_ai_error(result);

		// Add notifier.
		result = this->AddNotifier(message);
		aisdk::check_ai_error(result);

	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch(...)
	{
		result = kCantHappenErr;
	}
	return result;
}

/*
*/
ASErr AnnotatorPlugin::PostStartupPlugin()
{
	ASErr result = kNoErr;
	try {
		if (fAnnotator == NULL) {
			fAnnotator = new Annotator();
			SDK_ASSERT(fAnnotator);
		}
        result = sAIUser->CreateCursorResourceMgr(fPluginRef,&fResourceManagerHandle);
        aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch(...)
	{
		result = kCantHappenErr;
	}
	return result;
}

/*
*/
ASErr AnnotatorPlugin::ShutdownPlugin(SPInterfaceMessage* message)
{
	ASErr result = kNoErr;
	try {
		result = Plugin::ShutdownPlugin( message );
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch (...) {
		result = kCantHappenErr;
	}
	return result;
}

/*
*/
ASErr AnnotatorPlugin::Message(char* caller, char* selector, void* message)
{
	ASErr result = kNoErr;
	try {
		result = Plugin::Message(caller, selector, message);
		
		if (result == kUnhandledMsgErr) {
			if (strcmp(caller, kCallerAIAnnotation) == 0) {
				if (strcmp(selector, kSelectorAIDrawAnnotation) == 0) {
					result = this->DrawAnnotation((AIAnnotatorMessage*)message);
					aisdk::check_ai_error(result);
				}
				else if (strcmp(selector, kSelectorAIInvalAnnotation) == 0) {
					result = this->InvalAnnotation((AIAnnotatorMessage*)message);
					aisdk::check_ai_error(result);
				}
			}
		}
		else
			aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch (...) {
		result = kCantHappenErr;
	}
	return result;
}

/*
*/
ASErr AnnotatorPlugin::GoMenuItem(AIMenuMessage* message)
{
	ASErr result = kNoErr;
	try 
	{
		if (message->menuItem == fAboutPluginMenu) {
			SDKAboutPluginsHelper aboutPluginsHelper;
			aboutPluginsHelper.PopAboutBox(message, "About Annotator", kSDKDefAboutSDKCompanyPluginsAlertString);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch(...)
	{
		result = kCantHappenErr;
	}
	
	return result;
}

/*
*/
ASErr AnnotatorPlugin::UpdateMenuItem(AIMenuMessage* message)
{
	ASErr result = kNoErr;
	try 
	{
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch(...)
	{
		result = kCantHappenErr;
	}
	return result;
}

/*
*/
ASErr AnnotatorPlugin::Notify(AINotifierMessage* message)
{
	ASErr result = kNoErr;
	
	try 
	{
		if (message->notifier == fNotifySelectionChanged) {
			// Get the bounds of the current document view.
			AIRealRect viewBounds = {0, 0, 0, 0};
			result = sAIDocumentView->GetDocumentViewBounds(NULL, &viewBounds);
			aisdk::check_ai_error(result);
			// Invalidate the entire document view bounds.
			result = fAnnotator->InvalidateRect(viewBounds);
			aisdk::check_ai_error(result);
		}
		if (message->notifier == fShutdownApplicationNotifier)
		{
			if(fResourceManagerHandle != NULL)
			{
				sAIUser->DisposeCursorResourceMgr(fResourceManagerHandle);
				fResourceManagerHandle = NULL;
			}
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch(...)
	{
		result = kCantHappenErr;
	}

	return result;
}

/*
*/
ASErr AnnotatorPlugin::TrackToolCursor(AIToolMessage* message)
{
	ASErr result = kNoErr;
	try {
		if (this->fAnnotator) {
			result = this->fAnnotator->CheckForArtHit(message);
			aisdk::check_ai_error(result);
			if(sAIUser != NULL)
				result = sAIUser->SetSVGCursor(kAnnotatorToolIconResourceID, fResourceManagerHandle);
            aisdk::check_ai_error(result);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch(...)
	{
		result = kCantHappenErr;
	}
	return result;
}

/*
*/
ASErr AnnotatorPlugin::SelectTool(AIToolMessage* message)
{
	ASErr result = kNoErr;
	try {
		result = sAIAnnotator->SetAnnotatorActive(fAnnotatorHandle, true);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch(...)
	{
		result = kCantHappenErr;
	}
	return result;
}

/*
*/
ASErr AnnotatorPlugin::DeselectTool(AIToolMessage* message)
{
	ASErr result = kNoErr;
	try {
		result = sAIAnnotator->SetAnnotatorActive(fAnnotatorHandle, false);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch(...)
	{
		result = kCantHappenErr;
	}
	return result;
}

/*
*/
ASErr AnnotatorPlugin::AddTool(SPInterfaceMessage *message)
{
	ASErr result = kNoErr;
	try {
		ASErr			error = kNoErr;
		AIAddToolData	data;
        data.title = ai::UnicodeString::FromRoman("SDK Annotator Tool");
		data.tooltip = ai::UnicodeString::FromRoman("SDK Annotator Tool");
		data.sameGroupAs = kNoTool;
		data.sameToolsetAs = kNoTool;
		data.normalIconResID = kAnnotatorToolIconResourceID;
		data.darkIconResID = kAnnotatorToolIconResourceID;
		data.iconType = ai::IconType::kSVG;

		ai::int32 options = kToolWantsToTrackCursorOption;

		error = sAITool->GetToolNumberFromName(kAnnotatorTool, &data.sameGroupAs); // the Crop Area tool is now the first in its group
		if (error) return error;

		error = sAITool->GetToolNumberFromName(kAnnotatorTool, &data.sameToolsetAs); // and at the top of its tool stack
		if (error) return error;

		if (!result) {
			result = sAITool->AddTool(message->d.self, kAnnotatorTool, data, options, &fToolHandle);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch(...)
	{
		result = kCantHappenErr;
	}
	return result;
}

/*
*/
ASErr AnnotatorPlugin::AddAnnotator(SPInterfaceMessage *message)
{
	ASErr result = kNoErr;
	try {
		result = sAIAnnotator->AddAnnotator(message->d.self, "SDK Annotator", &fAnnotatorHandle);
		aisdk::check_ai_error(result);

		result = sAIAnnotator->SetAnnotatorActive(fAnnotatorHandle, false);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch(...)
	{
		result = kCantHappenErr;
	}
	return result;
}

/*
*/
ASErr AnnotatorPlugin::AddNotifier(SPInterfaceMessage *message)
{
	ASErr result = kNoErr;
	try {
		result = sAINotifier->AddNotifier(fPluginRef, "AnnotatorPlugin", kAIArtSelectionChangedNotifier, &fNotifySelectionChanged);
		aisdk::check_ai_error(result);
		result = sAINotifier->AddNotifier(fPluginRef, "AnnotatorPlugin", kAIApplicationShutdownNotifier, &fShutdownApplicationNotifier);
		aisdk::check_ai_error(result);

	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch(...)
	{
		result = kCantHappenErr;
	}
	return result;
}

/*
*/
ASErr AnnotatorPlugin::DrawAnnotation(AIAnnotatorMessage* message)
{
	ASErr result = kNoErr;
	try {
		if (this->fAnnotator) {
			result = this->fAnnotator->DrawArtAnnotation(message);
			aisdk::check_ai_error(result);

			result = this->fAnnotator->DrawCursorAnnotation(message);
			aisdk::check_ai_error(result);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch(...)
	{
		result = kCantHappenErr;
	}
	return result;
}

/*
*/
ASErr AnnotatorPlugin::InvalAnnotation(AIAnnotatorMessage* message)
{
	ASErr result = kNoErr;
	try {
		result = fAnnotator->InvalidateRect(*message->invalidationRects);
		aisdk::check_ai_error(result);
	}
	catch (ai::Error& ex) {
		result = ex;
	}
	catch(...)
	{
		result = kCantHappenErr;
	}
	return result;
}

// end MarkedObjectsPlugin.cpp
