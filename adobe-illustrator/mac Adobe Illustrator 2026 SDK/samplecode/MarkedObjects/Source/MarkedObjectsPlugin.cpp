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
#include "MarkedObjectsPlugin.hpp"
#include "MarkedObjectsPanelController.h"
#include "AICSXS.h"
#include "SDKErrors.h"

#define CHECK_ERROR if ( error ) {return error;}

extern ImportSuite gPostStartupSuites[];

MarkedObjectsPlugin *gPlugin = NULL;

Plugin *AllocatePlugin(SPPluginRef pluginRef)
{
	return new MarkedObjectsPlugin(pluginRef);
}

/*
*/
void FixupReload(Plugin* plugin)
{
	MarkedObjectsPlugin::FixupVTable((MarkedObjectsPlugin*) plugin);
}

/*
*/
ASErr MarkedObjectsPlugin::SetGlobal(Plugin *plugin)
{
	gPlugin = (MarkedObjectsPlugin *) plugin;
	return kNoErr;
}

/*
*/
MarkedObjectsPlugin::MarkedObjectsPlugin(SPPluginRef pluginRef) :
	Plugin(pluginRef),
	menuItem(NULL),
	fMarkedObjectsPanelController(NULL),
	artSelectionChangedNotifier(NULL),
	toolHandle(NULL),
	objectManager(NULL),
	preferences(NULL),
	aboutPluginMenu(NULL),
	fHideLayerNotifier(NULL),
	fResetNumberIDNotifier(NULL),
	fSaveObjectNotifier(NULL),
	fPreferencesSetNotifier(NULL),
	fObjectSelectedNotifier(NULL),
    fResourceManagerHandle(NULL),
	fShutdownApplicationNotifier(NULL)
{	
	strncpy(fPluginName, kMarkedObjectPluginName, kMaxStringLength);	
}

/*
*/
ASErr MarkedObjectsPlugin::StartupPlugin( SPInterfaceMessage *message )
{
	ASErr error = Plugin::StartupPlugin( message );
	CHECK_ERROR

	// Add About Plugins menu item for this plug-in.
	SDKAboutPluginsHelper aboutPluginsHelper;
	error = aboutPluginsHelper.AddAboutPluginsMenuItem(message, 
				kSDKDefAboutSDKCompanyPluginsGroupName, 
				ai::UnicodeString(kSDKDefAboutSDKCompanyPluginsGroupNameString), 
				"MarkedObjects...", 
				&aboutPluginMenu);
	CHECK_ERROR

	if (NULL == fMarkedObjectsPanelController)
	{
		fMarkedObjectsPanelController = new MarkedObjectsPanelController();
		CHECK_ERROR
	}

	if (NULL == preferences)
	{
		preferences = new MarkedObjectsPreferences();
		preferences->SetMarkedObjectsPanelController(fMarkedObjectsPanelController);
	}

	// Add other menu items.
	AIPlatformAddMenuItemDataUS menuData;
	AIPlatformAddMenuItemDataUS pluginMenuData;
	AIMenuGroup	menuGroup;
	AIMenuGroup	pluginMenuGroup;
	AIMenuGroup dummyGroup;
	AIMenuItemHandle dummyItem;
	AIMenuItemHandle pluginMenuItemHandle;
	bool groupAlreadyMade = false;
		
	const char * menuCStr = "SDKWindowGroup";
	const char * pluginCStr = "Marked Objects";

	ai::int32 count;
	error = sAIMenu->CountMenuGroups( &count );
	CHECK_ERROR
	while (count > 0)
	{
		error = sAIMenu->GetNthMenuGroup( count-1, &dummyGroup );
		CHECK_ERROR
		char *name;
		error = sAIMenu->GetMenuGroupName( dummyGroup, (const char**)&name );
		CHECK_ERROR
		if ( std::strcmp(name, menuCStr ) == 0 )
		{
			groupAlreadyMade = true;
			count = -1;
		}
		count--;
	}

	menuData.groupName = kOtherPalettesMenuGroup;
	ai::UnicodeString menuStr( "SDK" );
	menuData.itemText = menuStr;
	
	if ( !groupAlreadyMade )
	{
		error = sAIMenu->AddMenuItem( message->d.self, menuCStr, &menuData, 0, &dummyItem );
		CHECK_ERROR
		
		error = sAIMenu->AddMenuGroupAsSubMenu( menuCStr, kMenuGroupSortedAlphabeticallyOption, dummyItem, &menuGroup );
		CHECK_ERROR
	}
	

	pluginMenuData.groupName = menuCStr;
	pluginMenuData.itemText = ai::UnicodeString(pluginCStr);

	error = sAIMenu->AddMenuItem( message->d.self, pluginCStr, &pluginMenuData, 0,&pluginMenuItemHandle );
	CHECK_ERROR

	error = sAIMenu->AddMenuGroupAsSubMenu( pluginCStr, 0, pluginMenuItemHandle, &pluginMenuGroup );
	CHECK_ERROR
		
	pluginMenuData.groupName = pluginCStr;
	error = sAIMenu->AddMenuItem( message->d.self, 
								  "Marked Objects Menu", 
								  &pluginMenuData, 
								  kMenuItemWantsUpdateOption,
								  &menuItem );
	CHECK_ERROR
	
	// add some notifiers
	error = sAINotifier->AddNotifier( fPluginRef, "MarkedObjectsPlugin",
								      kAIUIBrightnessChangedNotifier, &brightnessChangedNotifier);
	CHECK_ERROR
	error = sAINotifier->AddNotifier( fPluginRef, "MarkedObjectsPlugin",
								      kAIArtSelectionChangedNotifier, &artSelectionChangedNotifier);
	CHECK_ERROR
	error = sAINotifier->AddNotifier( fPluginRef, "MarkedObjectsPlugin",
									 kAIDocumentChangedNotifier , &documentChangedNotifier);
	CHECK_ERROR
	error = sAINotifier->AddNotifier( fPluginRef, "MarkedObjectsPlugin",
									 "AI Command Notifier: Before Save", &documentSaveNotifierBefore );
	CHECK_ERROR
	error = sAINotifier->AddNotifier( fPluginRef, "MarkedObjectsPlugin",
									 "AI Command Notifier: Before Save As", &documentSaveAsNotifierBefore );
	CHECK_ERROR
	error = sAINotifier->AddNotifier( fPluginRef, "MarkedObjectsPlugin",
									 "AI Command Notifier: After Save", &documentSaveNotifierAfter );
	CHECK_ERROR
	error = sAINotifier->AddNotifier( fPluginRef, "MarkedObjectsPlugin",
									 "AI Command Notifier: After Save As", &documentSaveAsNotifierAfter );
	CHECK_ERROR
	error = sAINotifier->AddNotifier( fPluginRef, "Register Event Notify",
		kAICSXSPlugPlugSetupCompleteNotifier, &fRegisterEventNotifierHandle);
	CHECK_ERROR

	error = sAINotifier->AddNotifier( fPluginRef, "MarkedObjectsPlugin",
		kHideLayerNotifier, &fHideLayerNotifier);
	CHECK_ERROR

	error = sAINotifier->AddNotifier( fPluginRef, "MarkedObjectsPlugin",
		kResetNumberIDNotifier, &fResetNumberIDNotifier);
    CHECK_ERROR

	error = sAINotifier->AddNotifier( fPluginRef, "MarkedObjectsPlugin",
		kSaveObjectNotifier, &fSaveObjectNotifier);
	CHECK_ERROR

	error = sAINotifier->AddNotifier( fPluginRef, "MarkedObjectsPlugin",
		kPreferencesSetNotifier, &fPreferencesSetNotifier);
	CHECK_ERROR

	error = sAINotifier->AddNotifier( fPluginRef, "MarkedObjectsPlugin",
		kObjectSelectedNotifier, &fObjectSelectedNotifier);
	CHECK_ERROR

    error = sAINotifier->AddNotifier( fPluginRef, "MarkedObjectsPlugin",
                                     kAIApplicationShutdownNotifier, &fShutdownApplicationNotifier);
	CHECK_ERROR

	// add a tool
	AIAddToolData toolData;
		
	toolData.title = ai::UnicodeString::FromRoman("Marked Object Tool");
	toolData.tooltip = ai::UnicodeString::FromRoman("Marked Object Tool");
	toolData.sameGroupAs = kNoTool;

	toolData.darkIconResID = kMarkedObjectsToolIconResID;
	toolData.normalIconResID = kMarkedObjectsToolIconResID;
    
    sAITool->GetToolNumberFromName( kAddAnchorPointTool, &toolData.sameGroupAs );	
	toolData.sameToolsetAs = toolData.sameGroupAs;
	toolData.iconType = ai::IconType::kSVG;

	ai::int32 options = kToolWantsToTrackCursorOption;
	

		
	error = sAITool->AddTool( message->d.self, kMarkedObjectsTool, 
			                  toolData, options, &toolHandle );
	CHECK_ERROR
    
	return error;
}

ASErr MarkedObjectsPlugin::PostStartupPlugin()
{
	ASErr error = kNoErr;
    error = sAIUser->CreateCursorResourceMgr(fPluginRef,&fResourceManagerHandle);
	error = AcquirePostStartupSuites();
	return error;
}

/*
*/
ASErr MarkedObjectsPlugin::ShutdownPlugin( SPInterfaceMessage *message )
{
	delete preferences;
	
	// Destroy the flash UI
	if (fMarkedObjectsPanelController != NULL)
	{
		fMarkedObjectsPanelController->RemoveEventListeners();
		delete fMarkedObjectsPanelController;
		fMarkedObjectsPanelController = NULL;
		Plugin::LockPlugin(false);
	}

	ReleasePostStartupSuites();
			
	delete objectManager;

	ASErr error = Plugin::ShutdownPlugin( message );

	return error;
}

/* TODO: This should be moved to the plugin framework */
ASErr MarkedObjectsPlugin::AcquirePostStartupSuites()
{
	ASErr error = kNoErr;

	for ( int i = 0; gPostStartupSuites[i].name != nullptr; ++i )
	{
		if ( gPostStartupSuites[i].suite != nullptr )
		{
			ASErr tmperr = sSPBasic->AcquireSuite(	gPostStartupSuites[i].name, 
					gPostStartupSuites[i].version, 
					(const void **) gPostStartupSuites[i].suite );
			
			SDK_ASSERT_MSG_NOTHROW(!tmperr, 
					aisdk::format_args("AcquireSuite failed for suite=%s version=%d", 
					gPostStartupSuites[i].name, 
					gPostStartupSuites[i].version));

			if (tmperr && !error)
			{
				// A suite could not be acquired - note first error encountered for later return to caller then carry on.
				error = tmperr;
			}
		}
	}

	return error;
}

ASErr MarkedObjectsPlugin::ReleasePostStartupSuites()
{
	ASErr error = kNoErr;

	for ( int i = 0; gPostStartupSuites[i].name != nullptr; ++i )
	{
		if ( gPostStartupSuites[i].suite != nullptr )
		{
			void **s = (void **) gPostStartupSuites[i].suite;
			if (*s != nullptr)
			{
				ASErr tmperr = sSPBasic->ReleaseSuite( gPostStartupSuites[i].name, gPostStartupSuites[i].version );
				*s = nullptr;
				SDK_ASSERT_MSG_NOTHROW(!tmperr, 
						aisdk::format_args("ReleaseSuite failed for suite=%s version=%d", 
						gPostStartupSuites[i].name, 
						gPostStartupSuites[i].version));

				if (tmperr && !error)
				{
					// A suite could not be released - note first error encountered for later return to caller then carry on.
					error = tmperr;
				}
			}
		}
	}

	return error;
}

ASErr MarkedObjectsPlugin::GoMenuItem( AIMenuMessage *message )
{
	try 
	{
		if ( message->menuItem == menuItem ) 
		{
			if (fMarkedObjectsPanelController)
			{
				AIBoolean visible = false;

				fMarkedObjectsPanelController->IsPrimaryStageVisible(visible);

				if(visible)
				{
					fMarkedObjectsPanelController->UnloadExtension();
				}
				else
				{
					fMarkedObjectsPanelController->LoadExtension();
				}
			}
		} 
		else if (message->menuItem == aboutPluginMenu)
		{
			SDKAboutPluginsHelper aboutPluginsHelper;
			aboutPluginsHelper.PopAboutBox(message, "About MarkedObjects",
				kSDKDefAboutSDKCompanyPluginsAlertString);
		}
	}
	catch(...)
	{
		return 'ESHD';
	}
	
	return kNoErr;
}

ASErr MarkedObjectsPlugin::UpdateMenuItem( AIMenuMessage *message )
{
	AIBoolean visible = false;

	if ( fMarkedObjectsPanelController != NULL )
		fMarkedObjectsPanelController->IsPrimaryStageVisible(visible);

	sAIMenu->CheckItem(menuItem, visible);

	return kNoErr;
}

ASErr MarkedObjectsPlugin::Notify( AINotifierMessage *message )
{
	ASErr error = kNoErr;
	
	try 
	{

		if ( message->notifier == artSelectionChangedNotifier )
		{
			GetMarkedObjects();

			fMarkedObjectsPanelController->UpdateMarkedObjectSelected();

			if ( !IsMarkedObjectSelected() )
				fMarkedObjectsPanelController->Clear();
				
			fMarkedObjectsPanelController->UpdateBoundingBoxInfo();
		}

		if ( message->notifier == documentChangedNotifier )
		{
			fMarkedObjectsPanelController->Clear();

			fMarkedObjectsPanelController->UpdateMarkedObjectSelected();

			GetMarkedObjects();
		}
		
		if ( ( message->notifier == documentSaveNotifierBefore ) || ( message->notifier == documentSaveAsNotifierBefore ) )
		{
            if (objectManager != NULL)
            {
                objectManager->MakeLayerTemplate( true );
            }
		}

		if ( ( message->notifier == documentSaveNotifierAfter ) || ( message->notifier == documentSaveAsNotifierAfter ) )
		{
            if (objectManager != NULL)
            {
                if (objectManager->HasMarkedObjects())
                {
                    ai::FilePath fileSpec;
		
                    sAIDocument->GetDocumentFileSpecification( fileSpec );
			
                    WriteMarkedPoints( fileSpec );
                }
            }
		}

		if ( message->notifier == fRegisterEventNotifierHandle)
		{
			fMarkedObjectsPanelController->RegisterCSXSEventListeners();
		}

		
		if (message->notifier == brightnessChangedNotifier)
		{
			error = fMarkedObjectsPanelController->SendTheme();
		}

		if (strcmp(message->type, kHideLayerNotifier) == 0)
		{
			this->HideMarkedObjectsLayer(fMarkedObjectsPanelController->GetHideButtonStatus());
		}

		if (strcmp(message->type, kResetNumberIDNotifier) == 0)
		{
			ResetIDs();
			GetMarkedObjects();
			UnSelectAllArt();
		}

		if (strcmp(message->type, kSaveObjectNotifier) == 0)
		{
			WriteMarkedPoints( fMarkedObjectsPanelController->GetFilePath());
		}

		if (strcmp(message->type, kPreferencesSetNotifier) == 0)
		{
			DoModalPrefs();
			GetMarkedObjects();
		}

		if (strcmp(message->type, kObjectSelectedNotifier) == 0)
		{
			fMarkedObjectsPanelController->NewGlobalObjectSelected();
		}

		if (message->notifier == fShutdownApplicationNotifier)
		{
			if(fResourceManagerHandle != NULL)
			{
				error = sAIUser->DisposeCursorResourceMgr(fResourceManagerHandle);
				fResourceManagerHandle = NULL;
			}
		}


	}
	
	catch(...)
	{
		return 'EPN ';
	}

	return error;
}

ASErr MarkedObjectsPlugin::TrackToolCursor( AIToolMessage* message )
{
	AIErr error = kNoErr;
	if(sAIUser != NULL)
		error = sAIUser->SetSVGCursor(kMarkedObjectsToolIconResID, fResourceManagerHandle);
	return error;
}

ASErr MarkedObjectsPlugin::ToolMouseUp( AIToolMessage* message )
{
	ASErr error = kNoErr;
	
	try 
	{
		objectManager->CreateLayer();
		
		AIBoolean editState = objectManager->Edit( true );

		AIBoolean hiddenState = objectManager->Hide( false );

		objectManager->SetInsertionPoint();
		
		ostringstream idAsString;
		
		idAsString << objectManager->GetNextMarkedObjectID();
		
		MarkedObject* newMarkedObject = new MarkedObject( message->cursor, idAsString.str() );
		fMarkedObjectsPanelController->AddMarkedObject( newMarkedObject );
		objectManager->AddMarkedObject( newMarkedObject );
		
		UnSelectAllArt();
		
		newMarkedObject->Select();

		objectManager->Edit( editState );

		objectManager->Hide( hiddenState );

	}
	
	catch (...)
	{
		return 'EDTM';
	}
	
	return error;
}

ASErr MarkedObjectsPlugin::DeselectTool( AIToolMessage* message )
{
	try 
	{
		UnSelectAllArt();
		// If the user deselects using the shortcut keys the ResetInsertionPoint will
		// start a timer and try to ResetInsertionPoint later.
		objectManager->ResetInsertionPoint();
	}

	catch(...)
	{
		AITimerHandle timer;
        char name[] = "ResetInsertionPointTimer";
		sAITimer->AddTimer( message->d.self, name, 10, &timer );
	}

	return kNoErr;
}

ASErr MarkedObjectsPlugin::GoTimer( AITimerMessage* message )
{
	try
	{
		UnSelectAllArt();
		
		objectManager->ResetInsertionPoint();

		sAITimer->SetTimerActive( message->timer, false );
	}
	
	catch(...)
	{
		sAITimer->SetTimerActive( message->timer, true );
	}
	
	return kNoErr;
}

ASErr MarkedObjectsPlugin::HideMarkedObjectsLayer( AIBoolean hide )
{
	objectManager->Hide(hide);
	return kNoErr;
}

void MarkedObjectsPlugin::WriteMarkedPoints( const ai::FilePath & fileSpec )
{
	ai::FilePath fullPath( fileSpec );
	
	fullPath.RemoveExtension();
	
	fullPath.AddExtension( ".bb" );
	
	ofstream fileOut( fullPath.GetFullPath().as_Platform().c_str(), ios_base::trunc|ios_base::out );

	int documentVersion = sAIDocument->GetAIVersion();
	
	fileOut << "%%Creator: Adobe Illustrator " << documentVersion << endl;

	fileOut << "%%Title: ";

	string fileName = fileSpec.GetFileName().as_Roman();
	
	fileOut << fileName << endl;

	AIFileFormatHandle fileformat = NULL;
	const char *fileformatname;

	AIErr error = sAIDocument->GetDocumentFileFormat( &fileformat );
	if( error ) throw( error );

	if(fileformat != NULL)
	{
		error = sAIFileFormat->GetFileFormatName( fileformat, &fileformatname );
		if( error ) throw( error );
	}
	else
	{
        char name[] = "Document not saved to a format";
		fileformatname = name;
	}
	
	fileOut << "%%File Format: " << fileformatname << endl;

	fileOut << "%%Date: ";

	AIUserDateTime dateTime;    	
    error = sAIUser->GetDateAndTime(&dateTime);
	if(error) throw( error );
	
	ai::UnicodeString timeString, dateString;	
	error = sAIUser->GetDateString(&dateTime,0,dateString);
	error = sAIUser->GetTimeString(&dateTime,0,timeString);
	if(!error)
	{
		fileOut << dateString.as_Platform().c_str() << " ";
		fileOut << timeString.as_Platform().c_str() << endl;
	}
	else
	{
		fileOut << "Error getting date/time";
	}
				
	fileOut << "%%BoundingBox: ";
	
	AIRealRect mainArtBounds = MarkedObject::GetMainArtBounds();

	short left = 0;
	short bottom = 0;
	short right = 0;
	short top = 0;

	if ( documentVersion >= 9 )
	{
		right = MarkedObject::GetBoundingBoxWidth();
		top = MarkedObject::GetBoundingBoxHeight();
	}
	else
	{
		left = static_cast<short>( _AIRealTrunc( mainArtBounds.left ));
		bottom = static_cast<short>( _AIRealTrunc( mainArtBounds.bottom ));
		AIReal temp = (AIReal)( mainArtBounds.right ) - static_cast<short>( mainArtBounds.right );
		right = static_cast<short>( mainArtBounds.right );
		if ( temp > 0.0 )
			right++;
		temp = (AIReal)( mainArtBounds.top ) - static_cast<short>( mainArtBounds.top );
		top = static_cast<short>( mainArtBounds.top );
		if ( temp > 0.0 )
			top++;
	}

	fileOut << left << " " << bottom << " " << right << " " << top << endl;
	
	fileOut << "%%HiResBoundingBox: ";

	AIReal leftHiRes;
	AIReal bottomHiRes;
	AIReal rightHiRes;
	AIReal topHiRes;

	if ( documentVersion >= 9 )
	{
		leftHiRes = 0;
		bottomHiRes = 0;
		rightHiRes = mainArtBounds.right - mainArtBounds.left;
		topHiRes = mainArtBounds.top - mainArtBounds.bottom;
	}
	else
	{
		leftHiRes = mainArtBounds.left;
		bottomHiRes = mainArtBounds.bottom;
		rightHiRes = mainArtBounds.right;
		topHiRes = mainArtBounds.top;
	}

	ios::fmtflags flags_old = fileOut.flags();
	fileOut.flags( flags_old | ios::fixed );

	std::streamsize prec_old = fileOut.precision();
	
	fileOut.precision( 6 );

	fileOut << leftHiRes << " " << bottomHiRes << " " << rightHiRes << " " << topHiRes << endl;

	fileOut.precision( prec_old );
	fileOut.flags( flags_old );


	fileOut << "%%ArtBox: ";
	left = static_cast<short>( _AIRealTrunc( mainArtBounds.left ));
	bottom = static_cast<short>( _AIRealTrunc( mainArtBounds.bottom ));
	AIReal temp = (AIReal)( mainArtBounds.right ) - static_cast<short>( mainArtBounds.right );
	right = static_cast<short>( mainArtBounds.right );
	if ( temp > 0.0 )
		right++;
	temp = (AIReal)( mainArtBounds.top ) - static_cast<short>( mainArtBounds.top );
	top = static_cast<short>( mainArtBounds.top );
	if ( temp > 0.0 )
		top++;

	fileOut << left << " " << bottom << " " << right << " " << top << endl;

	fileOut << "%%HiResArtBox: ";
	leftHiRes = mainArtBounds.left;
	bottomHiRes = mainArtBounds.bottom;
	rightHiRes = mainArtBounds.right;
	topHiRes = mainArtBounds.top;

	flags_old = fileOut.flags();
	fileOut.flags( flags_old | ios::fixed );

	prec_old = fileOut.precision();
	fileOut.precision( 6 );

	fileOut << leftHiRes << " " << bottomHiRes << " " << rightHiRes << " " << topHiRes << endl;

	fileOut.precision( prec_old );
	fileOut.flags( flags_old );
	
	fileOut << "%%Coordinates: LL" << endl;

	fileOut << "%%StartMarkedPoints" << endl;

	objectManager->WriteAll( fileOut );

	fileOut << "%%EndMarkedPoints" << endl;
}

void MarkedObjectsPlugin::GetMarkedObjects( void )
{
	if ( fMarkedObjectsPanelController != NULL )
	{
		fMarkedObjectsPanelController->ClearMarkedObjectList();
		delete objectManager;
		objectManager = new MarkedObjectManager;
		fMarkedObjectsPanelController->SetMarkedObjectManager(objectManager);	
		AIArtHandle **artStore = NULL;
		ai::int32 artCount = 0;
				
		AIMatchingArtSpec spec;
		spec.type = kAnyArt;
		spec.whichAttr = 0;
		spec.attr = 0;

		// Ignore the error, we could be searching when there is no document
		// and we don't want that error reported, just let the null artStore
		// take over and silently quit the routine.
		sAIMatchingArt->GetMatchingArt( &spec, 1, &artStore, &artCount );
				
		if ( artStore != NULL )
		{
			// Go through all the art, deal with the non-MOs now the MOs will 
			// be dealt with after calculating the bounding box of the real art.
			
			vector<MarkedObject*> MOs;

			AIRealRect boundsOfMainArt = { 0, 0, 0, 0 };
			
			MOs.reserve( artCount );
			
			bool firstArtBounds = true;
			
			for ( ai::int32 x = 0; x < artCount; x++ )
			{
				AIArtHandle artHandle = (*artStore)[x];
				
				MarkedObject* markedObject = new MarkedObject( artHandle );
				
				if ( markedObject == NULL ) return;
					
				if ( markedObject->IsValid() )
				{
					MOs.push_back( markedObject );
				}
				else 
				{
				
					if ( ! ( markedObject->IsParentValid() ) )
					{
						AIRealRect nonMOBounds;
						short artType;
						
						AIErr error = sAIArt->GetArtType( artHandle, &artType );
						if ( error ) throw( error );
						
						ai::int32 attr = kArtHidden;
						error = sAIArt->GetArtUserAttr( artHandle, kArtHidden, &attr );
						if ( error ) throw( error );
						
						if ( artType != kGroupArt && !attr )
						{
							error = sAIArt->GetArtBounds( artHandle, &nonMOBounds );
							if ( error ) throw( error );
							
							if ( firstArtBounds )
							{
								boundsOfMainArt.top = nonMOBounds.top;
								boundsOfMainArt.left = nonMOBounds.left;
								boundsOfMainArt.bottom = nonMOBounds.bottom;
								boundsOfMainArt.right = nonMOBounds.right;
								firstArtBounds = false;
							}
							else
							{
								if ( nonMOBounds.top > boundsOfMainArt.top )
									boundsOfMainArt.top = nonMOBounds.top;
								if ( nonMOBounds.left < boundsOfMainArt.left )
									boundsOfMainArt.left = nonMOBounds.left;
								if ( nonMOBounds.bottom < boundsOfMainArt.bottom )
									boundsOfMainArt.bottom = nonMOBounds.bottom;
								if ( nonMOBounds.right > boundsOfMainArt.right )
									boundsOfMainArt.right = nonMOBounds.right;
							}
						}
					}
					
					delete markedObject;
				}
			}
			
			MarkedObject::SetMainArtBounds( boundsOfMainArt );

			// Go through all the MOs and add them to the manager and the dialog.
			vector<MarkedObject*>::iterator mo_iter = MOs.begin();
			vector<MarkedObject*>::iterator mo_iter_end = MOs.end();
			
			while ( mo_iter != mo_iter_end )
			{
				objectManager->AddMarkedObject( *mo_iter );
				fMarkedObjectsPanelController->AddMarkedObject( *mo_iter );
				mo_iter++;
			}
			
			AIErr error = sAIMdMemory->MdMemoryDisposeHandle( (void**) artStore );
			if ( error ) throw( error );
		}
	}
}

void MarkedObjectsPlugin::ResetIDs( void )
{
	objectManager->MakeLayerTemplate( false );
	objectManager->ResetIDs();
}

void MarkedObjectsPlugin::SelectOurTool( void )
{
	sAITool->SetSelectedTool( toolHandle );
}

ASBoolean MarkedObjectsPlugin::IsMarkedObjectSelected( void )
{
	ASBoolean result = false;
	
	AIArtHandle **artStore = NULL;
	ai::int32 artCount = 0;
				
	AIMatchingArtSpec spec;
	spec.type = kAnyArt;
	spec.whichAttr = kArtSelected;
	spec.attr = kArtSelected;

	// Ignore the error, we could be searching when there is no document
	// and we don't want that error reported, just let the null artStore
	// take over and silently quit the routine.		
	sAIMatchingArt->GetMatchingArt( &spec, 1, &artStore, &artCount );
				
	if ( artStore != NULL )
	{
		for ( ai::int32 x = 0; x < artCount && !result; x++ )
		{
			AIArtHandle artHandle = (*artStore)[x];
	
			MarkedObject* markedObject = new MarkedObject( artHandle );
			
			if ( markedObject->IsValid() )
			{
				result = true;
			}
			
			delete markedObject;
		}
	
		AIErr error = sAIMdMemory->MdMemoryDisposeHandle( (void**) artStore );
		if ( error ) throw( error );
	}

	return result;
}

void MarkedObjectsPlugin::UnSelectAllArt( void )
{
	AIArtHandle **artStore = NULL;
	ai::int32 artCount = 0;
				
	AIMatchingArtSpec spec;
	spec.type = kAnyArt;
	spec.whichAttr = kArtSelected;
	spec.attr = kArtSelected;

	sAIMatchingArt->GetMatchingArt( &spec, 1, &artStore, &artCount );
				
	if ( artStore != NULL )
	{
		for ( ai::int32 x = 0; x < artCount; x++ )
		{
			AIArtHandle artHandle = (*artStore)[x];
			
			short artType = kUnknownArt;
			
			sAIArt->GetArtType( artHandle, &artType );
			
			sAIArt->SetArtUserAttr( artHandle, kArtSelected | kArtFullySelected, 0 );
		}
		
		sAIMdMemory->MdMemoryDisposeHandle( (void**) artStore );
	}
}

void MarkedObjectsPlugin::DoModalPrefs( void )
{
	preferences->DoModalPrefs();
}

void MarkedObjectsPlugin::ResetInsertionPoint( void )
{
	objectManager->ResetInsertionPoint();
}
// end MarkedObjectsPlugin.cpp
