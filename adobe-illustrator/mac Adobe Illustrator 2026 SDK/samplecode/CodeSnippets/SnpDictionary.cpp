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

#include "SDKErrors.h"

// Framework includes:
#include "SnpRunnable.h"
#include "SnippetRunnerSuites.h" 
#include "SnippetRunnerLog.h"
#include "SnippetRunnerParameter.h"
#include "SnpArtHelper.h"
#include "SnpSelectionHelper.h"

/* Different string values of keys that would be used in this snippet */
#define kFillVisibilityKey		"FillVisibility"
#define kStrokeVisibilityKey	"StrokeVisibility"
#define kFillOpacityKey			"FillOpacity"
#define kStrokeOpacityKey		"StrokeOpacity"
#define kChildDictKey			"ChildDictionary"
#define kUnicodeStringKey		"UnicodeStringValue"
#define kStringKey				"StringValue"
#define kArtKey					"ArtValue"
#define kIntKey					"IntegerValue"
#define kSwapKey				"SwapEntry"
#define kNewArtKey				"NewArtKey"
#define kRectArtKey				"RectArtKey"
#define kTempKey				"TemporaryKey"

#define INTEGER_VALUE			8
#define TEMP_INT_VALUE			10

/* This is the gloabl variable for the Parent Dictionary used throughout this snippet*/
AIDictionaryRef ParentDict = NULL;

/* Flag to know whether a type is Dictionary or art object */
ASInt32 kArtFlag = 0; 



/** Provides functions to create a dictionary, add different types of entries,
	modify dictionary entries, and delete dictionaries.
	There is one Parent Dictionary which acts as a base dictionary for all other functions in this class.
	All other operations are taken on this Parent Dictionary.
*/
class SnpDictionary
{
	public:

		/*Various keys used in the snippet */
		/* Keys for the parent Dictionary */
		/* Key for Child Dictionary Entry to be added in the Parent Dictionary */
		AIDictKey fchildDictKey ;
		/* Key for UnicodeString Entry to be added in the Parent Dictionary */
		AIDictKey fUnicodeStringKey;
		/* Key for String Entry to be added in the Parent Dictionary */
		AIDictKey fStringKey ;
		/* Key for Art Entry to be added in the Parent Dictionary */
		AIDictKey fArtKey ;
		/* Key for Integer Entry to be added in the Parent Dictionary */
		AIDictKey fIntKey ;

		/* Keys for Child Dictionary */
		/* Dummy Key for boolean Entry to be added in the Child Dictionary */
		AIDictKey fFillVisibleKey ;
		/* Dummy Key for boolean Entry to be added in the Child Dictionary */
		AIDictKey fStrokeVisibleKey;
		/* Dummy Key for Real Entry to be added in the Child Dictionary */
		AIDictKey fFillOpacityKey ;
		/* Dummy Key for Real Entry to be added in the Child Dictionary */
		AIDictKey fStrokeOpacityKey; 

		/* Keys for Modifying Entries in the Parent Dictionary */
		/* Dummy Key to be used with SwapEntries API*/
		AIDictKey fSwapKey;
		/* Dummy Art Key for new art entry to be inserted in the Parent Dictionary */
		AIDictKey fNewArtKey;
		/* Dummy Art Key for new art entry to be inserted in the Parent Dictionary */
		AIDictKey fRectArtKey;
		/* Dummy Temporary Key for new entry to be inserted in the Parent Dictionary */
		AIDictKey fTempKey;

		
		/** Struct for the child dictionary entries.
			This structure acts as a dummy for setting values for the artstyle of an artwork.
			However, the purpose of this structure here is to display the creation and insertion of a new dictionary to an existing dictionary.
		*/
		typedef struct{
			/* Corresponding Boolean value for fFillVisibleKey */
			AIBoolean fFillVisible;
			/* Corresponding Boolean value for fStrokeVisibleKey */
			AIBoolean fStrokeVisible;
			/* Corresponding Boolean value for fFillOpacityKey */
			AIReal fFillOpacity;
			/* Corresponding Boolean value for fStrokeOpacityKey */
			AIReal fStrokeOpacity;
		}DictArtStyle;

		/** Constructor.
			Creates a parent dictionary with many other entries, if one does not exist. 
			It creates a parent dictionary, and sets its varios entires.
			@return kNoErr on success, other ASErr otherwise.
		 */
		SnpDictionary(void);		

		/**
			Creates and Sets entries for the child dictionary.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr CreateChildDictionary(void);

		/**
			Display dictionary entries.
			@param dictRef IN whose entries are to be displayed.
			@return kNoErr on success, other ASErr otherwise.
			*/
		ASErr DisplayDictionary(const ConstAIDictionaryRef dictRef);
		
		/**
			Make different modifications in the dictionary.
			@param dictRef IN whose entries are being modified.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr ModifyDictionary(const AIDictionaryRef dictRef);

		/**
			Delete the created dictionary.
			@param dictRef IN which is being deleted.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr DeleteDictionary(const AIDictionaryRef dictRef);
		
		
		//------------------------------- Utility Functions -------------------------------
		/**
			Sets entries for the parent dictionary.
			@param  dictRef IN whose entries are to be created.This dictionary reference contains the entries after
			running this function.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr AddEntries(const AIDictionaryRef dictRef);

		/**
			Sets struct values for the child dictionary entries.
			@param styleParam IN/OUT  whose individual structure entries are to set.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr SetDictStructValues(DictArtStyle* styleParam);

		/**
			Find a specific art type for an art entry.
			@param art IN object whose type is to be found through this function.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr FindArtType( AIArtHandle art);

		/**
			Delete the entries and keys of the dictionary.
			@param dictRef IN the dictionary containing the entries to be deleted.
			@param  key IN corresponding to the entry to be deleted.
			@return The size of the dictionary whose entries are deleted.
		*/
		ASInt32 DeleteEntries(AIDictionaryRef dictRef,AIDictKey dictkey);
};

/*
*/
SnpDictionary :: SnpDictionary(void)
{
	
	ASErr result = kNoErr;
	try {	
		/* Create the parent dictionary if one does not exist already. */
		if(!ParentDict){

			/* Create and Set the Keys defined in SnpDictionary class for the Parent Dictionary. */		
			fchildDictKey = sAIDictionary->Key(kChildDictKey);
			fUnicodeStringKey = sAIDictionary->Key(kUnicodeStringKey);
			fStringKey = sAIDictionary->Key(kStringKey);
			fArtKey = sAIDictionary->Key(kArtKey);
			fIntKey = sAIDictionary->Key(kIntKey);
			fSwapKey = sAIDictionary->Key(kSwapKey);
			fTempKey = sAIDictionary->Key(kTempKey);
				
			/* Create and Set the Keys for the Child Dictionary. */
			fFillVisibleKey = sAIDictionary->Key(kFillVisibilityKey);
			fStrokeVisibleKey = sAIDictionary->Key(kStrokeVisibilityKey);
			fFillOpacityKey = sAIDictionary->Key(kFillOpacityKey);
			fStrokeOpacityKey = sAIDictionary->Key(kStrokeOpacityKey);
			
			/* Create Parent Dictionary. */
			result = sAIDictionary->CreateDictionary(&ParentDict);
			aisdk::check_ai_error(result);
			/* Fill the parent dictionary with Keys*/
			result = this->AddEntries(ParentDict);
			aisdk::check_ai_error(result);
			result = sAIDictionary->SetUnicodeStringEntry(ParentDict,fUnicodeStringKey, ai::UnicodeString("Parent Dictionary"));
			aisdk::check_ai_error(result);

			/* Get the size of the Parent Dictionary. */
			ASInt32 size = sAIDictionary->Size(ParentDict);
			// Write size to the log.
			SnippetRunnerLog* log = SnippetRunnerLog::Instance();
			log->WritePrintf("Entry Count in Parent Dictionary is %d",size);
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
}

/*
*/
ASErr SnpDictionary::AddEntries(AIDictionaryRef dictRef)
{	
	ASErr result = kNoErr;

	try {			
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		log->Write("");
		log->Write("Creating entries for Parent Dictionary");
	
		const char* stringvalue = " This is SnpDictionary Snippet ";
		result = sAIDictionary->SetStringEntry(dictRef,fStringKey,stringvalue);
		aisdk::check_ai_error(result);

		int IntEntry = INTEGER_VALUE;
		result = sAIDictionary->SetIntegerEntry(dictRef,fIntKey,IntEntry);

		//Create a new art object to be inserted as an entry to the parent dictionary.
		SnpArtHelper artHelper;
		AIArtHandle art;
		result = artHelper.NewRectangle(art);
		aisdk::check_ai_error(result);
		
		ai::UnicodeString name;
		ASBoolean boolname;
		result = sAIArt->GetArtName(art,name,&boolname);
		aisdk::check_ai_error(result);
		log->Write(ai::UnicodeString("Name of the art object in the Layers Palette: ").append(name));

		//Get Art Type for the art object inserted
		short type = 0;
		result = sAIArt->GetArtType(art,&type);
		aisdk::check_ai_error(result);

		if(type){
			/*If this is an art object , then set this flag to 1. 
			It would be used to identify art entry as there is no specific entry type for art objects,
			and Ai treats them as DictType.*/
			kArtFlag = 1;
		}

		/* Add this entry to the dictionary */
		result = sAIDictionary->NewArtEntry(dictRef,fArtKey,type);
		aisdk::check_ai_error(result);
		
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDictionary::SetDictStructValues(DictArtStyle* styleParam)
{	
	ASErr result = kNoErr;
	try {
		/* Set the values for the Child Dictionary entries.*/
		styleParam->fFillVisible = true;
		styleParam->fStrokeVisible = true;
		styleParam->fFillOpacity = 2.0;
		styleParam->fStrokeOpacity = 4.0;
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDictionary :: FindArtType(AIArtHandle art)
{	
	ASErr result = kNoErr;
	try {
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();

		AIArtHandle newart = art;
		short arttype;
		result = sAIArt->GetArtType(newart,&arttype);
		aisdk::check_ai_error(result);

		switch(arttype){
				case kGroupArt:
					log->Write("GroupArt Art object");
					break;
				case kRasterArt:
					log->Write("RasterArt Art object");
					break;
				case kPathArt:
					log->Write("PathArt Art object");
					break;
				default:
					log->Write("Unknown type");
					break;
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;

}

/*
*/
ASInt32 SnpDictionary :: DeleteEntries(AIDictionaryRef dictRef,AIDictKey dictkey)
{
	ASErr result = kNoErr;
	try {
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();

		AIDictionaryRef CurrentDictRef = dictRef;
		AIDictKey DelKey = dictkey;

		/* Get the entry type, and display the type accordingly */
			AIEntryType entryType;
			result = sAIDictionary->GetEntryType(CurrentDictRef, DelKey, &entryType);
			aisdk::check_ai_error(result);

			switch(entryType){
				case BooleanType:
					log->Write("Deleting Boolean type ");
					break;

				case RealType:
					log->Write("Deleting Real type ");
					break;

				case IntegerType:
					log->Write("Deleting Integer type ");
					break;
			
				case StringType:
					log->Write("Deleting String type ");
					break;

				case UnicodeStringType:
					log->Write("Deleting UnicodeString type ");
					break;

				case DictType:
						if((sAIDictionary->Size(CurrentDictRef)) == 0)
							kArtFlag =1;
						/* if kArtFlag ==1, it means it is an art object type */
						if(kArtFlag){
							log->Write("Deleting Art type ");
							kArtFlag = 0;
							break;
						}
					
					/* if kArtFlag =0, it means its another dictionary entry in the present dictionary */
					else{
						log->Write("Deleting DictionaryType type ");
						break;
					}

				default:
					log->Write("None of known entry Type");
					break;
			}

		ASInt32 size;
		result = sAIDictionary->DeleteEntry(CurrentDictRef,DelKey);
		aisdk::check_ai_error(result);		
		return (size = sAIDictionary->Size(CurrentDictRef));
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return 0;

}

/*
*/
ASErr SnpDictionary::CreateChildDictionary(void)
{
	ASErr result = kNoErr;

	try {
				
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();	
		SnippetRunnerLog::Indent indent;

		/* Create Child Dictionary */
		AIDictionaryRef dictRef = NULL;
		static ASBoolean childDictflag = true;

		//Create a child dictionary only if one does not exist.
		if(childDictflag){
			log->Write("");
			log->Write("Creating entries for Child Dictionary");
			result = sAIDictionary->CreateDictionary(&dictRef);
			aisdk::check_ai_error(result);
			childDictflag = false;
			
			DictArtStyle artStyle;	
			/* Set individual components of the dictionary. */
			result = this->SetDictStructValues(&artStyle);
			aisdk::check_ai_error(result);

			/* Set the value for a Boolean Entry in the Child Dictionary .*/
			result = sAIDictionary->SetBooleanEntry(dictRef,fFillVisibleKey,artStyle.fFillVisible);
			aisdk::check_ai_error(result);

			/* Set the value for a Boolean Entry in the Child Dictionary. */
			result = sAIDictionary->SetBooleanEntry(dictRef,fStrokeVisibleKey,artStyle.fStrokeVisible);
			aisdk::check_ai_error(result);

			/* Set the value for a Real Entry in the Child Dictionary. */
			result = sAIDictionary->SetRealEntry(dictRef,fFillOpacityKey,artStyle.fFillOpacity);
			aisdk::check_ai_error(result);

			/* Set the value for a Real Entry in the Child Dictionary. */
			result = sAIDictionary->SetRealEntry(dictRef,fStrokeOpacityKey,artStyle.fStrokeOpacity);
			aisdk::check_ai_error(result);

			/* Add the Child Dictionary to the Parent Dictionary. */
			result = sAIDictionary->SetDictEntry(ParentDict,fchildDictKey,dictRef);
			aisdk::check_ai_error(result);

			/* Get the new size of the Parent Dictionary. */
			ASInt32 size = sAIDictionary->Size(ParentDict);
			log->WritePrintf("Entry Count in Parent Dictionary after inserting Child Dictionary is %d",size);
	
			/*Display the resulting dictionaries now. */
			result = this->DisplayDictionary(ParentDict);
			aisdk::check_ai_error(result);

		}
		else{				
			log->Write("This operation cannot be performed");
			log->Write("Either Child Dictionary already exists or the Parent Dictionary has been deleted");			
		}

	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDictionary::DisplayDictionary(ConstAIDictionaryRef dictRef)
{	
	ASErr result = kNoErr;
	try {
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerLog::Indent indent;
			
		AIDictionaryIterator iterator;
		/* Create a new iterator for traversing the dictionary */
		result = sAIDictionary->Begin(dictRef,&iterator);
		aisdk::check_ai_error(result);
			
		for( int i = 1; !sAIDictionaryIterator->AtEnd(iterator) ;sAIDictionaryIterator->Next(iterator),i++){
			log->Write("");			

			/* For every entry, get the corresponding key */
			AIDictKey dictKey = sAIDictionaryIterator->GetKey(iterator);
			log->WritePrintf("Dictionary Entry number : %d ",i);

			/* Get the entry type, and display the type accordingly */
			AIEntryType entryType;
			result = sAIDictionary->GetEntryType(dictRef, dictKey, &entryType);
			aisdk::check_ai_error(result);

			ai::UnicodeString unicode;

			switch(entryType){
				case BooleanType:
					ASBoolean boolean;
					result = sAIDictionary->GetBooleanEntry(dictRef,dictKey,&boolean);
					aisdk::check_ai_error(result);
					(boolean == false) ? log->Write("Boolean value : False") : log->Write("Boolean value : True");
					break;

				case RealType:
					AIReal real;
					result = sAIDictionary->GetRealEntry(dictRef,dictKey,&real);
					aisdk::check_ai_error(result);
					log->WritePrintf("Real value : %f",real);
					break;

				case IntegerType:
					ASInt32 integer;
					result = sAIDictionary->GetIntegerEntry(dictRef,dictKey,&integer);
					aisdk::check_ai_error(result);
					log->WritePrintf("Integer value : %d",integer);
					break;
			
				case StringType:
					const char* string;
					result = sAIDictionary->GetStringEntry(dictRef,dictKey,&string);
					aisdk::check_ai_error(result);
					log->WritePrintf("String value : %s",string);
					break;

				case UnicodeStringType:
					result = sAIDictionary->GetUnicodeStringEntry(dictRef,dictKey,unicode);
					log->Write(ai::UnicodeString("UnicodeString name: ").append(unicode));
					break;

				case DictType:
					/* if kArtFlag ==1, it means it is an art object type */
					if(kArtFlag){
						AIArtHandle art;
						result = sAIDictionary->GetArtEntry(dictRef,dictKey,&art);
						aisdk::check_ai_error(result);
						log->Write(" Art Type");
						kArtFlag = 0;

						result = this->FindArtType(art);
						aisdk::check_ai_error(result);
						break;
					}
					/* if kArtFlag = 0, it means its another dictionary entry in the present dictionary */
					else{
						AIDictionaryRef dict;
						result = sAIDictionary->GetDictEntry(dictRef,dictKey,&dict);
						aisdk::check_ai_error(result);
						
						log->Write(" Dict Type");
						result = this->DisplayDictionary(dict);
						aisdk::check_ai_error(result);
						break;
					}

				default:
					log->Write("None of known entry Type");
					break;
			}
		}
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/

ASErr SnpDictionary :: ModifyDictionary(AIDictionaryRef ParentDict)
{
	ASErr result = kNoErr;
	
	try {
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		
		/* Make a clone of the current dictionary */
		AIDictionaryRef cloneDict ;
		result = sAIDictionary->Clone(ParentDict,&cloneDict);
		aisdk::check_ai_error(result);

		/*Get the size of the dictionary */
		ASInt32 size = sAIDictionary->Size(cloneDict);
		log->WritePrintf("Total initial size of the Dictionary is : %d",size);

		/* If there are entries in the dictionary, then make the modifications */
		if( size != 0){

			fTempKey = sAIDictionary->Key(kTempKey);
			int TempEntry = TEMP_INT_VALUE;
			result = sAIDictionary->SetIntegerEntry(cloneDict,fTempKey,TempEntry);
			aisdk::check_ai_error(result);
			size = sAIDictionary->Size(cloneDict);
			log->WritePrintf("The size of the dictionary setting an entry is : %d",size);
		
			/* Illustration of copying one entry to another */
			if(sAIDictionary->IsKnown(cloneDict,fTempKey)){
				result = sAIDictionary->CopyEntry(cloneDict,cloneDict,fTempKey,sAIDictionary->Key(kSwapKey));
				aisdk::check_ai_error(result);
				size = sAIDictionary->Size(cloneDict);
				log->WritePrintf("The size of the dictionary after copying an entry is : %d",size);
			}
			
			/* Illustration of swapping one entry with another */
			if(sAIDictionary->IsKnown(cloneDict,fTempKey)){
				result = sAIDictionary->SwapEntries(cloneDict,cloneDict,fTempKey,sAIDictionary->Key(kUnicodeStringKey));
				aisdk::check_ai_error(result);
			}

			/* Illustration of getting an entry */
			AIEntryRef entryRef;
			if(sAIDictionary->IsKnown(cloneDict,fTempKey))
			entryRef = sAIDictionary->Get(cloneDict,fTempKey);
			size = sAIDictionary->Size(cloneDict);
			log->WritePrintf("The size of the dictionary after getting an element is : %d",size);

			/* Illustration of setting an entry */
			result = sAIDictionary->Set(cloneDict,fTempKey,entryRef);
			aisdk::check_ai_error(result);
			size = sAIDictionary->Size(cloneDict);
			log->WritePrintf("The size of the dictionary after setting an element is : %d",size);

			
			/* Illustration of setting a new art entry */
			AIArtHandle newArt;
			SnpArtHelper artHelper;
			result = artHelper.NewArc(newArt);

			short arttype;
			result = sAIArt->GetArtType(newArt,&arttype);
			aisdk::check_ai_error(result);
			fNewArtKey = sAIDictionary->Key(kNewArtKey);
			result = sAIDictionary->NewArtEntry(cloneDict,sAIDictionary->Key(kNewArtKey),arttype);
			aisdk::check_ai_error(result);

			/*If this is an art object , then set this flag to 1. 
			It would be used to identify art entry as there is no specific entry type for art objects,
			and Ai treats them as DictType.*/
			kArtFlag = 1;
			

			SnpDictionary instance;
			result = instance.FindArtType(newArt);
			aisdk::check_ai_error(result);

			ai::UnicodeString name;
			ASBoolean boolname;
			result = sAIArt->GetArtName(newArt,name,&boolname);
			aisdk::check_ai_error(result);
			
			AIArtHandle getArt;
			result = sAIDictionary->GetArtEntry(cloneDict,sAIDictionary->Key(kNewArtKey),&getArt);
			aisdk::check_ai_error(result);
			size = sAIDictionary->Size(cloneDict);
			log->WritePrintf("The size of the dictionary after inserting art elements is : %d",size);

			/* Illustration of copying an art entry from artwork to entry */
			AIArtHandle rectArt;
			result = artHelper.NewRectangle(rectArt);
			fRectArtKey = sAIDictionary->Key(kRectArtKey);
			result = sAIDictionary->CopyArtToEntry(cloneDict,fRectArtKey,rectArt);
			aisdk::check_ai_error(result);
			result = instance.FindArtType( rectArt);
			aisdk::check_ai_error(result);
			
			result = sAIArt->GetArtName(rectArt,name,&boolname);
			aisdk::check_ai_error(result);
			log->Write(ai::UnicodeString("Name of the art object in the Layers Palette: ").append(name));

			size = sAIDictionary->Size(cloneDict);
			log->WritePrintf("The size of the dictionary after Copying art elements is : %d",size);

			/* Illustration of moving an art entry from entry to artwork */
			AIArtHandle movedArt;
			result = sAIDictionary->MoveEntryToArt(cloneDict,fRectArtKey,kPlaceBelowAll,NULL,&movedArt);
			aisdk::check_ai_error(result);
			result = instance.FindArtType(movedArt);
			aisdk::check_ai_error(result);
			size = sAIDictionary->Size(cloneDict);
			log->WritePrintf("Now the size of the dictionary after moving art element is : %d",size);
		}

		/* Display the dictionary */
		result = this->DisplayDictionary(cloneDict);
		aisdk::check_ai_error(result);

	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

/*
*/
ASErr SnpDictionary::DeleteDictionary(AIDictionaryRef dictRef)
{
	
	ASErr result = kNoErr;
	try {
		ASInt32 size = sAIDictionary->Size(dictRef);
		// Write output to the log.
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		SnippetRunnerLog::Indent indent;
		log->WritePrintf("Entry Count in present Dictionary is %d",size);
		log->Write("");
		
		AIDictionaryIterator iterator;
		result = sAIDictionary->Begin(dictRef,&iterator);
		aisdk::check_ai_error(result);

		/* Delete the entries till the last entry is found. */
		while(size>0){
			AIDictKey dictKey = sAIDictionaryIterator->GetKey(iterator);
			size = this->DeleteEntries(dictRef,dictKey);
			result = sAIDictionary->Begin(dictRef,&iterator);
			aisdk::check_ai_error(result);
		}

		/* Display the dictionary. */
		result = this->DisplayDictionary(dictRef);
		aisdk::check_ai_error(result);
		
		/* Release the dictionary object. */
		sAIDictionary->Release(dictRef);
	}
	catch (ai::Error& ex) {
		result = ex;
	}	
	return result;
}

// --------------------------------- SnippetRunner framework hook ---------------------------------------------------

/* Makes the snippet SnpTemplate available to the SnippetRunner framework.
 */
class _SnpRunnableDictionary : public SnpRunnable
{
	public:
		/* Constructor registers the snippet with the framework.
		 */
		_SnpRunnableDictionary () : SnpRunnable() {}

		/* Destructor.
		 */
		virtual ~_SnpRunnableDictionary () {}

		/* Returns name of snippet.
		*/
		std::string GetName() const;

		/* Returns a description of what the snippet does.
		*/
		std::string	GetDescription() const;

		/* Returns operations supported by this snippet.
		*/
		Operations GetOperations() const;

		/* Returns name of this snippet's default operation - must
			be one of the operation names returned by GetOperations.
		*/
		std::string GetDefaultOperationName() const;

		/** Returns the categories a snippet belongs to.
			@return default categories.
		*/
		std::vector<std::string> GetCategories() const;

		/* Returns true if the snippet can run.
			@param runnableContext see ISnpRunnableContext for documentation.
			@return true if snippet can run, false otherwise
		 */
		ASBoolean			CanRun(SnpRunnable::Context& runnableContext);

		/* Runs the snippet.
			@param runnableContext see ISnpRunnableContext for documentation.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr		Run(SnpRunnable::Context& runnableContext);
};

/*
*/
std::string _SnpRunnableDictionary::GetName() const
{
	return "Dictionary";
}

/*
*/
std::string _SnpRunnableDictionary::GetDescription() const 
{
	return "Provides Dictionary Operations";
}

/*
*/
SnpRunnable::Operations _SnpRunnableDictionary::GetOperations() const
{
	SnpRunnable::Operations operations;
	operations.push_back(Operation("CreateChildDictionary", "document,Can be performed just once", kSnpRunNotSupportedContext));
	operations.push_back(Operation("DisplayDictionary", "document", kSnpRunNotSupportedContext));
	operations.push_back(Operation("ModifyDictionary", "document", kSnpRunNotSupportedContext));
	operations.push_back(Operation("DeleteDictionary", "document", kSnpRunNotSupportedContext));
	return operations;
}

/*
*/
std::string _SnpRunnableDictionary::GetDefaultOperationName() const 
{	
	return "CreateChildDictionary";
	
}

/*
*/
std::vector<std::string> _SnpRunnableDictionary::GetCategories() const
{
	std::vector<std::string> categories = SnpRunnable::GetCategories();
	categories.push_back("Data Utility Snippets");
	return categories;
}

/* Checks if preconditions are met.
*/
ASBoolean _SnpRunnableDictionary::CanRun(SnpRunnable::Context& runnableContext)
{
	SnpSelectionHelper selectionHelper;
	return selectionHelper.IsDocumentSelected();
	
	
}

/* Instantiates and calls your snippet class.
*/
ASErr _SnpRunnableDictionary::Run(SnpRunnable::Context& runnableContext)
{
	ASErr result = kNoErr;
	SnpDictionary instance;
	if ("CreateChildDictionary" == runnableContext.GetOperation().GetName()) {
		result = instance.CreateChildDictionary();
	}
	else if ("DisplayDictionary" == runnableContext.GetOperation().GetName()) {
		result = instance.DisplayDictionary(ParentDict);
	}
	else if ("ModifyDictionary" == runnableContext.GetOperation().GetName()) {
		result = instance.ModifyDictionary(ParentDict);
	}
	else if ("DeleteDictionary" == runnableContext.GetOperation().GetName()) {
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		log->Write("Deleting Parent Dictionary");
		result = instance.DeleteDictionary(ParentDict);
	}
	else {
		result = kBadParameterErr;
	}
	return result;
}

/*	Declares an instance to register the snippet hook with the framework.
*/
static _SnpRunnableDictionary instance_SnpRunnableDictionary;

// End SnpDictionary.cpp

