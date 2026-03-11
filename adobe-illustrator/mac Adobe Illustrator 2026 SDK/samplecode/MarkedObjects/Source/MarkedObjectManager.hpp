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

#ifndef __MarkedObjectManager_hpp__
#define __MarkedObjectManager_hpp__

#include "MarkedObject.hpp"
#include "MarkedObjectsPreferences.hpp"

class MarkedObject;

class MarkedObjectManager 
{
private:
	vector<MarkedObject *> markedObjects;
	AIDictionaryRef documentDict;
	AILayerHandle layerHandle;
	AILayerHandle currentLayerHandle;
	
	void FindLayer( void );

public:
	MarkedObjectManager();
	~MarkedObjectManager();
	
	ASInt32 GetNextMarkedObjectID( void );
	void SetNextMarkedObjectID( ASInt32 newID );
	void AddMarkedObject( MarkedObject * mo );
	
	bool Hide( AIBoolean hide );
	bool Edit( AIBoolean edit );
	void MakeLayerTemplate( AIBoolean layerTemplate );
	void WriteAll( ofstream & fileOut );
	
	void SelectLayer( void );
	void CreateLayer( void );
	void SetInsertionPoint( void );
	void ResetInsertionPoint( void );
	
	void ResetIDs( void );

	bool HasMarkedObjects( void );

	MarkedObject *GetMarkedObject(int index)
	{
		if(markedObjects.size()>0)
		return markedObjects[index];
		return NULL;
	}

	ASInt32 GetMarkedObjectSize( void )
	{
		return markedObjects.size();
	}

};

#endif	
// end MarkedObjectManager.hpp
