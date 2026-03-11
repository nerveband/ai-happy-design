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
#include "TextFileFormatPlugin.h"

/*
*/
Plugin* AllocatePlugin(SPPluginRef pluginRef)
{
	return new TextFileFormatPlugin(pluginRef);
}

/*
*/
void FixupReload(Plugin* plugin)
{
	TextFileFormatPlugin::FixupVTable((TextFileFormatPlugin*) plugin);
}

/*
*/
TextFileFormatPlugin::TextFileFormatPlugin(SPPluginRef pluginRef) 
	: Plugin(pluginRef)
{
	strncpy(fPluginName, kTextFileFormatPluginName, kMaxStringLength);
}

/*
*/
ASErr TextFileFormatPlugin::Message(char* caller, char* selector, void *message) 
{
	ASErr error = kNoErr;

	try {
		error = Plugin::Message(caller, selector, message);
	}
	catch (ai::Error& ex) {
		error = ex;
	}
	catch (...) {
		error = kCantHappenErr;
	}
	if (error) {
		if (error == kUnhandledMsgErr) {
			// Defined by Plugin.hpp and used in Plugin::Message - ignore.
			error = kNoErr;
		}
		else {
			Plugin::ReportError(error, caller, selector, message);
		}
	}	
	return error;
}


/*
*/
ASErr TextFileFormatPlugin::StartupPlugin(SPInterfaceMessage* message)
{
	ASErr error = kNoErr;
	error = Plugin::StartupPlugin(message);
	if (error) { return error;  }
	error = this->AddMenus(message);
    if (error) { return error;  }
	error = this->AddFileFormats(message);
	
	return error;
}


ASErr TextFileFormatPlugin::GoMenuItem(AIMenuMessage* message)
{
	ASErr error = kNoErr;
	if (message->menuItem == this->fAboutPluginMenu) {
		// Pop this plug-in's about box.
		SDKAboutPluginsHelper aboutPluginsHelper;
		aboutPluginsHelper.PopAboutBox(message, "About TextFileFormat", kSDKDefAboutSDKCompanyPluginsAlertString);
	}	
	return error;
}


ASErr TextFileFormatPlugin::AddMenus(SPInterfaceMessage* message) {
	ASErr error = kNoErr;
	// Add a menu item to the About SDK Plug-ins menu group.
	SDKAboutPluginsHelper aboutPluginsHelper;
	error = aboutPluginsHelper.AddAboutPluginsMenuItem(message, 
				kSDKDefAboutSDKCompanyPluginsGroupName, 
				ai::UnicodeString(kSDKDefAboutSDKCompanyPluginsGroupNameString), 
				"TextFileFormat...", 
				&this->fAboutPluginMenu);
	return error;
}

ASErr TextFileFormatPlugin::AddFileFormats(SPInterfaceMessage* message) 
{
	ASErr error = kNoErr;

	PlatformAddFileFormatData affd;
	char pstrSelected[kMaxStringLength] = "Selected Text as TEXT";
	char pstrAll[kMaxStringLength] = "All Text as TEXT";
	
    affd.title = ai::UnicodeString::FromRoman(pstrSelected);
	affd.titleOrder = 0;
	affd.extension = ai::UnicodeString::FromRoman("txt");
	
	error = sAIFileFormat->AddFileFormat( message->d.self, "Selected text format",
				                          &affd, kFileFormatWrite,
										  &this->fFileFormatSelected, kNoExtendedOptions );
	if (error)
		return error;

	
	affd.title = ai::UnicodeString::FromRoman(pstrAll);
	error = sAIFileFormat->AddFileFormat( message->d.self, "All text format",
			                              &affd, kFileFormatWrite | kFileFormatRead,
									      &this->fFileFormatAll, kNoExtendedOptions );
	return error;
}

ASErr TextFileFormatPlugin::GoFileFormat(AIFileFormatMessage* message) 
{
	ASErr error = kNoErr;
	char pathName[300];

	message->GetFilePath().GetFullPath().as_Roman( pathName, 300);
	
	if ( message->option & kFileFormatWrite ) 
	{
		// Get whatever art is supposed to be saved.  You might walk the tree writing out 
		// information on each object as you come to it, or as in this case
		// just get a bunch of art and work in batch mode.  It depends on what your doing.					

		error = WriteText( pathName, message->fileFormat == this->fFileFormatSelected );
	}

	else if ( message->option & kFileFormatRead ) 
	{
		error = ReadText( pathName );
	} 
	
	return error;
}

ASErr TextFileFormatPlugin::WriteText(const char* pathName, bool selectedTextOnly) 
{
	ASErr error = kNoErr;
	
#ifdef MAC_ENV
	FILE * fp = fopen( pathName, "w" );
#else
	FILE *fp;
	errno_t err;
	err = fopen(&fp, pathName, "w" );
#endif
	
	if ( fp != NULL )
	{

		ATE::ITextRanges iTextRanges;

		if ( selectedTextOnly )
		{
			TextRangesRef textSelectionRef;
			error = sAIDocument->GetTextSelection( &textSelectionRef );
			if ( !error )
			{
				iTextRanges = ATE::ITextRanges( textSelectionRef );
			}
		}
		else
		{
			AIArtHandle **artStore = NULL;
			ai::int32 artCount = 0;
				
			AIMatchingArtSpec spec;
			spec.type = kTextFrameArt;
			spec.whichAttr = 0;
			spec.attr = 0;
			
			(void) sAIMatchingArt->GetMatchingArt( &spec, 1, &artStore, &artCount );

			for( ai::int32 i = 0; i < artCount && !error; i++ )
			{
				AIArtHandle artHandle = (*artStore)[i];
				TextRangeRef textRange;
				error = sAITextFrame->GetATETextRange( artHandle, &textRange );
				if ( error ) continue;
				ATE::ITextRange iTextRange( textRange );
				iTextRanges.Add( iTextRange );
			}

			(void) sAIMdMemory->MdMemoryDisposeHandle( (void**) artStore );

		}

		if ( !error )
		{

			ASInt32 textRangeSize = iTextRanges.GetSize();

			for ( int i = 0; i < textRangeSize; i++ ) 
			{
				ATE::ITextRange iTextRange = iTextRanges.Item( i );
				ASInt32 strLength = iTextRange.GetSize();
				if ( strLength > 0 )
				{
					std::vector<char> vc(strLength);
					ASInt32 conLength = iTextRange.GetContents( &vc[0], strLength );
					if ( conLength == strLength )
					{
						std::string contents;
						contents.assign( vc.begin(), vc.begin() + strLength );
						fputs( contents.c_str(), fp );
					}
				}
			}
		}
	
		fclose( fp );
	
	}

	return error;
}

ASErr TextFileFormatPlugin::ReadText(const char * pathName)
{
	ASErr error = kNoErr;

#ifdef MAC_ENV
	FILE * fp = fopen( pathName, "r" );
#else
	FILE *fp;
	errno_t err;
	err = fopen(&fp, pathName, "r" );
#endif

	if ( fp != NULL ) 
	{
		AIDocumentSetup docSetup;
		
		error = sAIDocument->GetDocumentSetup( &docSetup );
		
		if ( !error )
		{
			AIArtHandle textHandle = CreateTextArt( docSetup.width, docSetup.height );
			
			if ( textHandle != NULL )
			{
				TextRangeRef textRange;

				error = sAITextFrame->GetATETextRange( textHandle, &textRange );
		
				if ( !error )
				{
					ATE::ITextRange iTextRange( textRange );

					iTextRange.Remove();

					const unsigned int textSize = 1024 * 100;
					std::unique_ptr<char> text( new char[textSize] ); 
					size_t bytesRead = 0;

					do 
					{
						bytesRead = fread( text.get(), 1, textSize, fp );
						if ( ! ferror( fp ) && bytesRead )
						{
                            iTextRange.InsertAfter( ai::UnicodeString(text.get()).as_ASUnicode().c_str(), bytesRead);
						}
					} while ( ! ferror( fp ) && ! feof( fp ) && bytesRead );
				}
			}
		}

		fclose( fp );
	
	}
						
	return error;
}		

AIArtHandle TextFileFormatPlugin::CreateTextArt(AIReal width, AIReal height)
{
	AIArtHandle artHandle = NULL;
	AIArtHandle textArtHandle = NULL;
	ASErr error = sAIArt->NewArt(kPathArt, kPlaceAboveAll, NULL, &artHandle);
	if ( !error )
	{
		error = sAIPath->SetPathSegmentCount( artHandle, 4);
		if ( !error )
		{
			AIPathSegment segments[4];

			const int shrink = 10;
				
			AIReal left = 0 + shrink;
			AIReal top = 0 + shrink;
				
			AIReal right = width - shrink;
			AIReal bottom = height - shrink;
				
			segments[0].p.h = left; 	
			segments[0].p.v = top; 		
			segments[1].p.h = right; 	
			segments[1].p.v = top; 		
			segments[2].p.h = right; 	
			segments[2].p.v = bottom; 	
			segments[3].p.h = left; 	
			segments[3].p.v = bottom;
								
			segments[0].in = segments[0].p;
			segments[0].out = segments[0].p;
				
			segments[1].in = segments[1].p;
			segments[1].out = segments[1].p;
				
			segments[2].in = segments[2].p;
			segments[2].out = segments[2].p;
				
			segments[3].in = segments[3].p;
			segments[3].out = segments[3].p;
				
			error = sAIPath->SetPathSegments( artHandle, 0, 4, segments );
			if ( !error )
				error = sAIPath->SetPathClosed( artHandle, true );		
			if ( !error )
				error = sAITextFrame->NewInPathText(kPlaceAboveAll, NULL, 
													kHorizontalTextOrientation, 
													artHandle, NULL, false, 
				                                    &textArtHandle );
			else
				error = sAIArt->DisposeArt( artHandle );
		}
	}

	return textArtHandle;
}

/////////////////////////////////////////
//
// Utility functions
//
/////////////////////////////////////////

/*
*/
void TextFileFormatPlugin::CStrToPStr(char *s, ai::UnicodeString::size_type len)
{
	const ai::UnicodeString sAsUnicode((const char*)s);
	ai::PStr sAsPStr((unsigned char*) s);
	sAsUnicode.getToBuffer(sAsPStr, len, kAIUTF8CharacterEncoding );
}

// End TextFileFormatPlugin.cpp
