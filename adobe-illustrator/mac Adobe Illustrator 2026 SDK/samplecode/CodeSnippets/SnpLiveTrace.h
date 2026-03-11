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
#include "IAIRef.h"
#include "SDKErrors.h"

// Framework includes:
#include "SnpRunnable.h"
#include "SnippetRunnerSuites.h"
#include "SnippetRunnerLog.h"
#include "SnippetRunnerParameter.h"
#include "SnpArtSetHelper.h"
#include "SnpDocumentHelper.h"
#include "SnpSelectionHelper.h"

/** Shows how to perform live tracings using default and specified parameters, remove
	live trace effects, copy out the result art from a tracing group, create a new 
	live trace art item from a file and display the statistics of a tracing group.

	The API which contains the tracing presets is not exposed currently so they cannot be 
	accessed, instead of using the presets this sample uses either the default tracing options,
	options input by the user or options from a previously traced art item.

	@see AITracingSuite, AIDictionarySuite, AIArtSuite, AIArtSetSuite, AICountedObjectSuite
  	@ingroup sdk_snippet
 */
class SnpLiveTrace
{
	public:
		/** If running in default parameter mode or unit test mode - traces the selected source art item(s) 
			using default tracing options. If running in prompt parameter mode prompts user for values of 
			Tracing Mode and Blur options, then traces selected source art item(s) using the selected values 
			along with default values for the rest of the tracing options.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr CreateTracing();

		/** Traces the given art handle.
			@param art IN the art to be traced.
			@param tracing OUT the traced art.
			@param tracingModeKey IN tracing mode - color, grayscale or black and white, default color.
			@param blur IN amount of blur to be applied, default 0.
			@param tracingMaxColorsKey IN max colors allowed in tracing, default 10.
			@param tracingResamplingKey IN resampling resoulution, default 150.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr CreateTracing(const AIArtHandle& art, AIArtHandle& tracing,
			ai::int32 tracingModeKey = kAIVectorizeModeColor, AIReal blur = 0,
			ai::int32 tracingMaxColorsKey = 10, AIReal tracingResamplingKey = 150);

		/** Traces the selected source art item(s) using the options applied to the selected 
			tracing group, if tracing art from a file select the tracing group before
			running the snippet
			NB. This function requires only 1 tracing group and at least 1 source art item 
			selected.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr CreateTracingUsingOptionsFromTracingGroup();

		/** Traces art using the tracing options from an existing tracing.
			object speficied by the tracingS
			@param tracingSrc IN the thacing whose options are to be used.
			@param artSet IN the art to be traced.
			@param tracingArtSet OUT the traced art.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr CreateTracingUsingOptionsFromTracingGroup(AIArtHandle tracingSrc, const SnpArtSetHelper& artSet, SnpArtSetHelper& tracingArtSet);

		/** Removes tracing effects from all selected tracing groups, returns error if no 
			tracing groups are selected
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr ReleaseTracing();

		/** Removes the tracings of the given objects.
			@param tracingArtSet IN the tracings to be removed.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr ReleaseTracing(const SnpArtSetHelper& tracingArtSet);

		/** Copies the result art group from the selected tracing groups to the current
			document, returns error if no tracing groups are selected
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr CopyTracing();

		/** Copies the given tracings.
			@param tracingArtSet IN the tracing art to copy.
			@param copyTracingArtSet OUT the result art from the copied tracings.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr CopyTracing(const SnpArtSetHelper& tracingArtSet, SnpArtSetHelper& copyTracingArtSet);

		/** Places a file and then traces it.
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr PlaceAndCreateTracing();

		/** Reports tracing related properties of selected tracings to the SnippetRunnerLog, 
			returns error if no tracings are selected
			@return kNoErr on success, other ASErr otherwise.
		 */
		ASErr InspectTracing();

		/** Reports tracing related properties to the SnippetRunnerLog.
			@param tracingArtSet IN the tracings to be inspected.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr InspectTracing(const SnpArtSetHelper& tracingArtSet);

		/** Returns set of objects that can be traced.
			@param selected IN if true search only selected art, if false search all
			art in the document.
			@return set of art which can be traced.
		*/
		SnpArtSetHelper FindCanTraceArtSet(ASBoolean selected);

		/** Returns set of objects that are tracings.
			@param selected IN if true search only selected art, if false search all
			art in the document.
			@return set of art which are tracing group items.
		*/
		SnpArtSetHelper FindIsTracingArtSet(ASBoolean selected);

		/** Checks whether an art item is part of a tracing group and returns the result.
			@param art IN the art to check.
			@return true if art object or any ancestor is part of a tracing, false otherwise.
		*/
		ASBoolean IsArtInTracing(AIArtHandle art);
		
		/** Gets the result art of the traced art passed in.
			@param tracedArt IN the traced art.
			@return kNoErr on success, other ASErr otherwise.
		*/
		ASErr GetResultArt(AIArtHandle& tracedArt);
};
