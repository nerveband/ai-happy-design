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

#include "MarkedObject.hpp"

const char * const kNotMarkedObject = NULL;

using namespace ATE;

static const char* kMarkedObjects_IDStr = "kMarkedObjects_IDStr";
static const char* kMarkedObjects_TypeStr = "kMarkedObjects_TypeStr";

AIRealRect MarkedObject::mainArtBounds = { 0, 0, 0, 0 };
AIArtHandle MarkedObject::insertNewArtHere = NULL;

static AIErr ConvertToInt(const string & s, int& i);

MarkedObject::MarkedObject( AIRealPoint location, const string & id )
{
	groupHandle = NULL;
	
	pointHandle = NULL;
	slashHandle = NULL;
	backslashHandle = NULL;
	IDHandle = NULL;
	labelHandle = NULL;
	
	CreateArt( location, id );	
}

MarkedObject::MarkedObject( AIArtHandle newHandle )
{
	groupHandle = newHandle;

	pointHandle = NULL;
	slashHandle = NULL;
	backslashHandle = NULL;
	IDHandle = NULL;
	labelHandle = NULL;
	
	sAIArt->GetArtFirstChild( groupHandle, &pointHandle );
	sAIArt->GetArtSibling( pointHandle, &slashHandle );
	sAIArt->GetArtSibling( slashHandle, &backslashHandle );
	sAIArt->GetArtSibling( backslashHandle, &IDHandle );
	sAIArt->GetArtSibling( IDHandle, &labelHandle );
}

MarkedObject::~MarkedObject()
{
}

void MarkedObject::CreateArt( AIRealPoint location, const string & id )
{	
	try 
	{
		groupHandle = CreateGroupArt( location, id );
		
		pointHandle = CreatePointArt( location );
		
		AIRealPoint start = location;
		start.h -= 4;
		start.v -= 4;
		
		AIRealPoint stop = location;
		stop.h += 4;
		stop.v += 4;
		
		slashHandle = CreateLineArt( start, stop );
		
		start = location;
		start.h -= 4;
		start.v += 4;
		
		stop = location;
		stop.h += 4;
		stop.v -= 4;
		
		backslashHandle = CreateLineArt( start, stop );
		
		AIRealPoint subLocation = location;
		AIRealPoint subOffset = gPreferences->GetSubStringLocation();
		subLocation.h += subOffset.h;
		subLocation.v += subOffset.v;
		
		IDHandle = CreateTextArt( subLocation, id, false );
		
		AIRealPoint labelLocation = location;
		AIRealPoint labelOffset = gPreferences->GetLabelStringLocation();
		labelLocation.h += labelOffset.h;
		labelLocation.v += labelOffset.v;

		labelHandle = CreateTextArt( labelLocation, gPreferences->GetLabelStringDefaultText(), true );
		
		SetID( id );
		
		SetLabel( gPreferences->GetLabelStringDefaultText() );
	}

	catch(...)
	{
		return;
	}
}

AIArtHandle MarkedObject::CreateGroupArt( AIRealPoint location, const string & id )
{
	AIArtHandle group = NULL;
	
	AIErr error;
	
	if ( MarkedObject::insertNewArtHere == NULL )
		error = sAIArt->NewArt( kGroupArt, kPlaceAboveAll, NULL, &group );
	else
		error = sAIArt->NewArt( kGroupArt, kPlaceInsideOnTop, MarkedObject::insertNewArtHere, &group );
	if ( error ) throw( error );
		
	ai::UnicodeString idu( id );
	error = sAIArt->SetArtName( group, idu );
	if ( error ) throw( error );
	
	return ( group );
}

AIArtHandle MarkedObject::CreatePointArt( AIRealPoint location )
{
	AIArtHandle pointHandle = NULL;
	
	AIErr error = sAIArt->NewArt( kPathArt, kPlaceInsideOnBottom, groupHandle, &pointHandle );
	if ( error ) throw( error );
		
	error = sAIPath->SetPathSegmentCount( pointHandle, 1 );
	if ( error ) throw( error );
		
	AIPathSegment segment;
	segment.p.h = location.h;
	segment.p.v = location.v;
	segment.in = segment.out = segment.p;
	segment.corner = true;
		
	error = sAIPath->SetPathSegments( pointHandle, 0, 1, &segment );
	if ( error ) throw( error );

	error = sAIPath->SetPathClosed( pointHandle, true );
	if ( error ) throw( error );

	return( pointHandle );
}

AIArtHandle MarkedObject::CreateLineArt( AIRealPoint start, AIRealPoint stop )
{
	AIArtHandle lineHandle = NULL;
	
	AIErr error = sAIArt->NewArt( kPathArt, kPlaceInsideOnBottom, groupHandle, &lineHandle );
	if ( error ) throw( error );
		
	error = sAIPath->SetPathSegmentCount( lineHandle, 2 );
	if ( error ) throw( error );
		
	AIPathSegment segment[2];
	segment[0].p.h = start.h;
	segment[0].p.v = start.v;
	segment[0].in = segment[0].out = segment[0].p;
	segment[0].corner = true;
		
	segment[1].p.h = stop.h;
	segment[1].p.v = stop.v;
	segment[1].in = segment[1].out = segment[1].p;
	segment[1].corner = true;

	error = sAIPath->SetPathSegments( lineHandle, 0, 2, segment );
	if ( error ) throw( error );

	error = sAIPath->SetPathClosed( lineHandle, true );
	if ( error ) throw( error );
	
	AIPathStyle style;
	error = sAIPathStyle->GetPathStyle( lineHandle, &style, nullptr);
	if ( error ) throw( error );
	
	style.strokePaint = true;
	style.stroke.color.kind = kFourColor;
	style.stroke.color.c.f.cyan = 0;
	style.stroke.color.c.f.magenta = 0;
	style.stroke.color.c.f.yellow = 0;
	style.stroke.color.c.f.black = 100;
	style.stroke.width = 1;
	style.stroke.dash.length = 0;

	error = sAIPathStyle->SetPathStyle( lineHandle, &style );
	if ( error ) throw( error );

	return( lineHandle );
}

AIArtHandle MarkedObject::CreateTextArt( AIRealPoint location, const string & text, bool big )
{
	AIErr error = 0;

	AIArtHandle textHandle = NULL;

	error = sAITextFrame->NewPointText( kPlaceInsideOnBottom, 
	                                    groupHandle, 
											  kHorizontalTextOrientation, 
											  location, &textHandle );

	if ( error ) throw( error );

	TextRangeRef textRange;

	error = sAITextFrame->GetATETextRange( textHandle, &textRange );
	if ( error ) throw( error );

	ITextRange iTextRange( textRange );
	
	iTextRange.Remove();

	iTextRange.InsertAfter( ai::UnicodeString(text).as_ASUnicode().c_str() );

	ICharFeatures pFeatures;

	string fontName;
	if (big)
		fontName = gPreferences->GetLabelStringFontName();
	else
		fontName = gPreferences->GetSubStringFontName();

	AIFontKey fontKey;
	error = sAIFont->FindFont( (char*)fontName.c_str(), kAITrueTypeFontTechnology, kRomanAIScript, true, &fontKey );
	if ( error ) throw( error );
		                
	FontRef fontRef;
	error = sAIFont->FontFromFontKey ( fontKey, &fontRef );

	IFont iFont( fontRef );

	pFeatures.SetFont( iFont );

	pFeatures.SetFontSize( big ? gPreferences->GetLabelStringFontSize() 
						 : gPreferences->GetSubStringFontSize());

	iTextRange.SetLocalCharFeatures(pFeatures);

	return ( textHandle );

}

bool MarkedObject::IsValid( void )
{
	return ( !GetID().empty() );
}

bool MarkedObject::IsParentValid( void )
{
	AIArtHandle parentHandle;

	bool returnValue = true;
	
	AIErr error = sAIArt->GetArtParent( groupHandle, &parentHandle );
	if ( error ) return( false );
	
	const char* id = kNotMarkedObject;
	
	AIDictKey dictKey = sAIDictionary->Key( kMarkedObjects_IDStr );
	
	AIDictionaryRef dict;
	
	if ( error == kNoErr ) error = sAIArt->GetDictionary( parentHandle, &dict );
	if ( error ) returnValue = false ;
	
	if ( error == kNoErr ) error = sAIDictionary->GetStringEntry( dict, dictKey, &id );
	if ( error ) returnValue = false ;

	sAIDictionary->Release( dict );
		
	return( returnValue );
}

const string MarkedObject::GetID( void ) const
{
	string sid;

	const char* id = kNotMarkedObject;
	
	AIDictKey dictKey = sAIDictionary->Key( kMarkedObjects_IDStr );
	
	AIDictionaryRef groupDict = NULL;
	
	AIErr error = sAIArt->GetDictionary( groupHandle, &groupDict );

	if ( !error )
		error = sAIDictionary->GetStringEntry( groupDict, dictKey, &id );

	if ( !error )
		sid.assign( id, id + strlen(id) );

	sAIDictionary->Release( groupDict );
	
	return sid;
}

void MarkedObject::SetID( const string & id )
{
	AIDictKey dictKey = sAIDictionary->Key( kMarkedObjects_IDStr );
	
	AIDictionaryRef groupDict = NULL;
	
	sAIArt->GetDictionary( groupHandle, &groupDict );
	
	sAIDictionary->SetStringEntry( groupDict, dictKey, id.c_str() );
	
	sAIDictionary->Release( groupDict );
	
	TextRangeRef textRange;

	AIErr error = sAITextFrame->GetATETextRange( IDHandle, &textRange );
	if ( error ) throw( error );

	ITextRange iTextRange( textRange );
		
	iTextRange.Remove();

	iTextRange.InsertAfter( ai::UnicodeString(id).as_ASUnicode().c_str() );

	ai::UnicodeString idu( id );

	error = sAIArt->SetArtName( groupHandle, idu );

	if ( error ) throw( error );
}

const string MarkedObject::GetLabel( void )
{
	string label;
		
	TextRangeRef textRange;

	AIErr error = sAITextFrame->GetATETextRange( labelHandle, &textRange );
	if ( error ) throw( error );

	ITextRange iTextRange( textRange );

	ASInt32 strLength = iTextRange.GetSize();
	if ( strLength <= 0 ) return label;

	std::vector<char> vc( strLength );

	ASInt32 getContentsLength = iTextRange.GetContents( &vc[0], strLength );

	if ( getContentsLength != strLength ) return label;

	while ( vc[strLength - 1] == '\r' )
	{
		vc[strLength - 1] = 0;
		strLength--;
	}

	label.assign( vc.begin(), vc.begin() + strLength );

	return label;
	
}
 
void MarkedObject::SetLabel( const string & newString )
{

	TextRangeRef textRange;

	AIErr error = sAITextFrame->GetATETextRange( labelHandle, &textRange );
	if ( error ) throw( error );

	ITextRange iTextRange( textRange );
		
	iTextRange.Remove();

	iTextRange.InsertAfter( ai::UnicodeString(newString).as_ASUnicode().c_str() );
}

MarkedObjectType MarkedObject::GetType( void )
{
	AIDictKey dictKey = sAIDictionary->Key( kMarkedObjects_TypeStr );
	
	ASInt32 typeAsInt = -1;

	AIDictionaryRef groupDict = NULL;
	
	sAIArt->GetDictionary( groupHandle, &groupDict );
	
	sAIDictionary->GetIntegerEntry( groupDict, dictKey, &typeAsInt );
		
	sAIDictionary->Release( groupDict );
	
	return (MarkedObjectType)typeAsInt;
}

void MarkedObject::SetType( MarkedObjectType newType )
{
	AIDictKey dictKey = sAIDictionary->Key( kMarkedObjects_TypeStr );
	
	AIDictionaryRef groupDict = NULL;
	
	sAIArt->GetDictionary( groupHandle, &groupDict );
	
	sAIDictionary->SetIntegerEntry( groupDict, dictKey, newType );

	sAIDictionary->Release( groupDict );
}

AIRealPoint MarkedObject::GetLocation( void )
{
	AIRealRect bounds;
	AIRealPoint location = { 0, 0};
	
	sAIArt->GetArtBounds( pointHandle, &bounds );
	
	location.h = bounds.left;
	location.v = bounds.top;
	
	return location;
}

void MarkedObject::SetLocation( AIRealPoint newLocation )
{
	AIRealPoint oldLocation = GetLocation();
	
	if ( ( newLocation.h != oldLocation.h ) || ( newLocation.v != oldLocation.v ) )
	{
		const string label = GetLabel();

		MarkedObjectType type = GetType();
		
		const string id = GetID();
		
		sAIArt->DisposeArt( groupHandle );
		
		CreateArt( newLocation, id );
		
		SetLabel( label );
		
		SetType( type );
	}
}

const string MarkedObject::GetFullString( void )
{
	ostringstream outStream;

	outStream << "%%MarkedPoint:(";

	ai::int32 precision = gPreferences->GetPrecision();
	
	ios::fmtflags flags_old = outStream.flags();


	outStream.flags( flags_old | ios::fixed );
	outStream.precision( precision );

	AIRealPoint location = GetLocationFromMainArt();

	if (precision > 0)
	{
		outStream << location.h;
		outStream << ",";
		outStream << location.v;
	}
	else
	{
		outStream << (short)location.h;
		outStream << ",";
		outStream << (short)location.v;
	}
	
	outStream << ") : ";
	
	outStream << GetID();

	const string label = GetLabel();

	outStream << " %" << label;
	
	return outStream.str();
}

void MarkedObject::HideAndLock( void )
{
	sAIArt->SetArtUserAttr( groupHandle, kArtHidden | kArtLocked, kArtHidden | kArtLocked );
}

void MarkedObject::Select( void )
{
	sAIArt->SetArtUserAttr( pointHandle, kArtSelected, kArtSelected );
	sAIArt->SetArtUserAttr( slashHandle, kArtSelected, kArtSelected );
	sAIArt->SetArtUserAttr( backslashHandle, kArtSelected, kArtSelected );
	sAIArt->SetArtUserAttr( IDHandle, kArtSelected, kArtSelected );
	sAIArt->SetArtUserAttr( labelHandle, kArtSelected, kArtSelected );
}

AIRealPoint MarkedObject::GetLocationFromMainArt( void )
{
	AIRealPoint MOLocation = GetLocation();
	
	MOLocation.h -= MarkedObject::mainArtBounds.left;
	MOLocation.v -= MarkedObject::mainArtBounds.bottom;
	
	return ( MOLocation );
}

void MarkedObject::SetMainArtBounds( AIRealRect newBounds )
{
	MarkedObject::mainArtBounds.top = newBounds.top;
	MarkedObject::mainArtBounds.left = newBounds.left;
	MarkedObject::mainArtBounds.bottom = newBounds.bottom;
	MarkedObject::mainArtBounds.right = newBounds.right;
}

AIRealRect MarkedObject::GetMainArtBounds( void )
{
	return ( MarkedObject::mainArtBounds );
}

short MarkedObject::GetBoundingBoxWidth( void )
{
	AIReal width = MarkedObject::mainArtBounds.right - MarkedObject::mainArtBounds.left;

	AIReal roundUp = (AIReal)(width) - short(width);

	short sWidth = (short)width;

	if ( roundUp > 0.0 )
			sWidth++;

	return ( sWidth );
}

short MarkedObject::GetBoundingBoxHeight( void )
{
	AIReal height = MarkedObject::mainArtBounds.top - MarkedObject::mainArtBounds.bottom;

	AIReal roundUp = (AIReal)(height) - short(height);

	short sHeight = (short)height;

	if ( roundUp > 0.0 )
			sHeight++;

	return ( sHeight );
}

bool MarkedObject::IsSelected( void )
{
	ai::int32 attrib = 0;

	AIErr error = sAIArt->GetArtUserAttr( groupHandle, kArtSelected, &attrib );

	if ( error ) return ( false );

	if ( attrib ) return ( true );

	return ( false );
}

void MarkedObject::SetInsertionArt( AIArtHandle insertionArt )
{
	MarkedObject::insertNewArtHere = insertionArt;
}

bool MarkedObject::operator>( const MarkedObject& rhs )
{
	const string rhs_ID = rhs.GetID();
	int rhs_I;
	AIErr rhs_error = ConvertToInt( rhs_ID, rhs_I );
	
	int this_I;
	AIErr this_error = ConvertToInt( GetID(), this_I );
	
	if ( rhs_error == 0 && this_error == 0 )
		return this_I > rhs_I;
	else if ( rhs_error && this_error )
		return ( rhs_ID.compare( GetID() ) > 0 );
	else if ( rhs_error == 0 )
		return true;
	else 
		return false;
	
	return true;
}

static AIErr ConvertToInt(const string & s, int& i)
{
	unsigned int index = 0;

	bool error = false;
	
	if (s.length() > 0)
	{
		if (s.at(index) == '+')
			index++;
		else if (s.at(index) == '-')
			index++;
		
		while ( index < s.length() && !error )
		{
			if (s.at(index) > '9' || s.at(index) < '0')
				error = true;
			index++;
		}
	}
	
	if (!error)
		i = atoi( s.c_str() );
	
	return error;
}
// end MarkedObject.cpp
