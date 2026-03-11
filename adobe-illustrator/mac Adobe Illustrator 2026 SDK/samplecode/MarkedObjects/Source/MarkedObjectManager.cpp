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

#include "MarkedObjectManager.hpp"

static const char* kMarkedObjects_DKNextMOID = "kMarkedObjects_DKNextMOID";
static const char* kMarkedObjects_LayerTitle = "Marked Objects DO NOT MODIFY";

MarkedObjectManager::MarkedObjectManager()
{
	try 
	{
		markedObjects.clear();

		layerHandle = NULL;
		currentLayerHandle = NULL;
				
		FindLayer();
		
		SetInsertionPoint();
		
		if ( sAIDocument->GetDictionary( &documentDict ) )
			documentDict = NULL;
	}
	
	catch(...)
	{
		return;
	}
}

MarkedObjectManager::~MarkedObjectManager()
{
	try 
	{
		vector<MarkedObject*>::iterator mo_iter = markedObjects.begin();
		vector<MarkedObject*>::iterator mo_iter_end = markedObjects.end();
		
		while ( mo_iter != mo_iter_end )
		{
			delete *mo_iter;
			mo_iter++;
		}
		
		sAIDictionary->Release( documentDict );
	}
	
	catch(...)
	{
		return;
	}
}

ASInt32 MarkedObjectManager::GetNextMarkedObjectID( void )
{
	AIDictKey dictKey = sAIDictionary->Key( kMarkedObjects_DKNextMOID );

	ASInt32 nextMarkedObjectID = 1;		

	sAIDictionary->GetIntegerEntry( documentDict, dictKey, &nextMarkedObjectID );	
		
	SetNextMarkedObjectID( nextMarkedObjectID + 1 );
	
	return nextMarkedObjectID;
}

void MarkedObjectManager::SetNextMarkedObjectID( ASInt32 newID )
{
	AIDictKey dictKey = sAIDictionary->Key( kMarkedObjects_DKNextMOID );

	sAIDictionary->SetIntegerEntry( documentDict, dictKey, newID );		
}

void MarkedObjectManager::AddMarkedObject( MarkedObject* mo )
{
	try 
	{
		if ( layerHandle == NULL )
		{
			CreateLayer();
		}
			
		markedObjects.push_back( mo );
	}
	
	catch(...)
	{
		return;
	}
}

bool MarkedObjectManager::Hide( AIBoolean hide )
{
	AIBoolean visible = true;

	try 
	{
		AIErr error = kNoErr;
		
		if ( layerHandle != NULL )
		{
			error = sAILayer->GetLayerVisible( layerHandle, &visible );
			if ( error ) throw( error );
			
			if ( visible != !hide )
				error = sAILayer->SetLayerVisible( layerHandle, !hide );
			if ( error ) throw( error );
		}
		
		return !visible;
	}
	
	catch(...)
	{
		return !visible;
	}
}	

bool MarkedObjectManager::Edit( AIBoolean edit )
{
	AIBoolean editable = true;

	try
	{
		AIErr error = kNoErr;

		if ( layerHandle != NULL )
		{
			error = sAILayer->GetLayerEditable( layerHandle, &editable );
			if ( error ) throw( error );

			if ( edit != editable )
				error = sAILayer->SetLayerEditable( layerHandle, edit );
			if ( error ) throw( error );
		}

		return editable != 0;
	}

	catch(...)
	{
		return editable != 0;
	}
}

void MarkedObjectManager::MakeLayerTemplate( AIBoolean layerTemplate )
{
	try 
	{
		AIErr error = kNoErr;
		
		if ( layerHandle != NULL )
		{
			error = sAILayer->SetLayerIsTemplate( layerHandle, layerTemplate );
			if ( error ) throw( error );
			AIBoolean lT;
			error = sAILayer->GetLayerIsTemplate( layerHandle, &lT );
			if ( error ) throw( error );
			if ( lT != layerTemplate )
			{
				return;
			}
		}
	}
	
	catch(...)
	{
		return;
	}
}	

class IDLess {
public:
	bool operator()(MarkedObject* x, MarkedObject* y)
	{
		if ( x == NULL ) return 0;
		if ( y == NULL ) return 0;

		return ( *y > *x );
	}
};

void MarkedObjectManager::WriteAll( ofstream& fileOut )
{
	try 
	{
	
		if (gPreferences->GetAutoSort())
			sort( markedObjects.begin(), markedObjects.end(), IDLess() );

		vector<MarkedObject*>::iterator mo_iter = markedObjects.begin();
		vector<MarkedObject*>::iterator mo_end = markedObjects.end();
		
		while ( mo_iter != mo_end )
		{
			const string st = (*mo_iter)->GetFullString();

			fileOut << st;

			fileOut << endl;
			
			mo_iter++;
		}

	}
	
	catch(...)
	{
		return;
	}
}

void MarkedObjectManager::FindLayer( void )
{
	AILayerHandle tempHandle = NULL;
	ai::UnicodeString layerTitle;
	
	sAILayer->GetFirstLayer( &tempHandle );
	while ( tempHandle != NULL )
	{
		ASErr error = sAILayer->GetLayerTitle( tempHandle, layerTitle );
		if ( error ) return;

		if (strcmp( layerTitle.as_Platform().c_str(), kMarkedObjects_LayerTitle ) == 0 )
		{
			layerHandle = tempHandle;
			tempHandle = NULL;
		} 
		else 
		{
			error = sAILayer->GetNextLayer( tempHandle, &tempHandle );
			if ( error ) return;
		}
	}
}

void MarkedObjectManager::CreateLayer( void )
{
	if ( layerHandle == NULL )
	{
		ASErr error = sAILayer->InsertLayer( NULL, kPlaceAboveAll, &layerHandle );
		if ( error ) return;
	
		ai::UnicodeString layerTitle( kMarkedObjects_LayerTitle );
	
		sAILayer->SetLayerTitle( layerHandle, layerTitle );
	}
}

void MarkedObjectManager::SetInsertionPoint( void )
{

	
	if ( layerHandle != NULL )
	{
		MakeLayerTemplate( false );

		AIArtHandle layerArt = NULL;

		ASErr error = sAIArt->GetFirstArtOfLayer( layerHandle, &layerArt );	
		if ( error == kNoErr && layerArt != NULL )
			MarkedObject::SetInsertionArt( layerArt );
	}
	else
	{
		MarkedObject::SetInsertionArt( NULL );
	}
}

void MarkedObjectManager::ResetInsertionPoint( void )
{
	AILayerHandle aLayerHandle = NULL;
	
	AIErr error = sAILayer->GetFirstLayer( &aLayerHandle );
	if ( error || aLayerHandle == NULL ) throw( error ? error : 'NULL' );

	do { 
		ai::UnicodeString layerTitle;
		error = sAILayer->GetLayerTitle( aLayerHandle, layerTitle );
		if ( error ) throw( error );
		
		if (strcmp( layerTitle.as_Platform().c_str(), kMarkedObjects_LayerTitle ) == 0 )
			error = sAILayer->GetNextLayer( aLayerHandle, &aLayerHandle );
		else
			error = 1;
			
	} while ( aLayerHandle != NULL  && !error );
	 
	AIArtHandle firstLayerArt = NULL;
	
	error = sAIArt->GetFirstArtOfLayer( aLayerHandle, &firstLayerArt );
	if ( error || firstLayerArt == NULL ) throw( error ? error : 'NULL' );
	
	error = sAIArt->SetInsertionPoint( firstLayerArt );
	if ( error ) throw( 0 );
}

void MarkedObjectManager::ResetIDs( void )
{
	ASInt32 id = 1;
	
	vector<MarkedObject*>::iterator mo_iter = markedObjects.begin();
	vector<MarkedObject*>::iterator mo_iter_end = markedObjects.end();
		
	while ( mo_iter != mo_iter_end )
	{
		ostringstream idAsString;
		idAsString << id++;
		(*mo_iter)->SetID( idAsString.str().c_str() );
		mo_iter++;
	}

	SetNextMarkedObjectID( id );
}

bool MarkedObjectManager::HasMarkedObjects( void )
{
	return markedObjects.size() != 0;
}
// end MarkedObjectManager.cpp
