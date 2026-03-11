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

#ifndef __MarkedObject_hpp__
#define __MarkedObject_hpp__

#include "MarkedObjectsSuites.h"
#include "MarkedObjectsPreferences.hpp"

typedef enum _MarkedObjectType {
	kPoint,
	kEllipse,
	kOval,
	kRectangle,
	kPolygon
} MarkedObjectType;

class MarkedObject 
{
private:
	AIArtHandle groupHandle;

	AIArtHandle pointHandle;
	AIArtHandle slashHandle;
	AIArtHandle backslashHandle;
	AIArtHandle IDHandle;
	AIArtHandle labelHandle;
	
	static AIRealRect mainArtBounds;
	static AIArtHandle insertNewArtHere;

	void CreateArt( AIRealPoint location, const string & id );
	AIArtHandle CreateGroupArt( AIRealPoint location, const string & id );
	AIArtHandle CreatePointArt( AIRealPoint location );
	AIArtHandle CreateLineArt( AIRealPoint start, AIRealPoint stop );
	AIArtHandle CreateTextArt( AIRealPoint location, const string & text, bool big );

	MarkedObject( const MarkedObject& );// don't create this yet
public:
	MarkedObject( AIRealPoint location, const string & id );
	MarkedObject( AIArtHandle newHandle );
	~MarkedObject();
	
	bool IsValid( void );
	bool IsParentValid( void );
	
	const string GetID( void ) const;
	void SetID( const string & id );
	
	const string GetLabel( void );
	void SetLabel( const string & newString );
	
	MarkedObjectType GetType( void );
	void SetType( MarkedObjectType newType );
	
	AIRealPoint GetLocation( void );
	void SetLocation( AIRealPoint newLocation );
	
	AIRealPoint GetLocationFromMainArt( void );
	static void SetMainArtBounds( AIRealRect newBounds );
	static AIRealRect GetMainArtBounds( void );
	static short GetBoundingBoxHeight( void );
	static short GetBoundingBoxWidth( void );
	
	static void SetInsertionArt( AIArtHandle insertionArt );
	
	const string GetFullString( void );
	
	void HideAndLock( void );
	
	void Select( void );
	
	bool IsSelected( void );
	
	bool operator>( const MarkedObject & rhs );
	
};

#endif	
// end MarkedObject.hpp
