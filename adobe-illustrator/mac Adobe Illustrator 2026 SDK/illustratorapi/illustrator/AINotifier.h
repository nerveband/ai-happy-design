#ifndef __AINotifier__
#define __AINotifier__

/*******************************************************************************
* ADOBE CONFIDENTIAL
*
* Copyright 1986 Adobe
* All Rights Reserved.
*
* NOTICE: Adobe permits you to use, modify, and distribute this file in
* accordance with the terms of the Adobe license agreement accompanying
* it. If you have received this file from a source other than Adobe,
* then your use, modification, or distribution of it requires the prior
* written permission of Adobe.
*******************************************************************************/


/*******************************************************************************
 **
 **	Imports
 **
 **/

#ifndef __AITypes__
#include "AITypes.h"
#endif

#ifndef __AIPlugin__
#include "AIPlugin.h"
#endif


#include "AIHeaderBegin.h"

/** @file AINotifier.h */


/*******************************************************************************
 **
 ** Constants
 **
 **/

#define kAINotifierSuite			"AI Notifier Suite"
#define kAINotifierSuiteVersion		AIAPI_VERSION(5)
#define kAINotifierVersion			kAINotifierSuiteVersion

/** @ingroup Callers
	See \c #AINotifierSuite. */
#define kCallerAINotify 		"AI Notifier"
/** @ingroup Selectors
	See \c #AINotifierSuite. */
#define kSelectorAINotify		"AI Notify"


/*******************************************************************************
 **
 ** Types
 **
 **/

/** Opaque reference to a notifier, never dereferenced. See \c #AINotifierSuite. */
typedef struct _t_AINotifierOpaque *AINotifierHandle;

/** The contents of a notifier message. */
struct AINotifierMessage {
	/** The message data */
	SPMessageData d;
	/** The notifier handle returned from \c #AINotifierSuite::AddNotifier(). If more than
		one notifier is installed, can be used to determine the instance that should
		handle the message. */
	AINotifierHandle notifier;
	/** The type of event for which the plug-in is being notified.
		If more than one notifier type is installed, can be used alone or
		with the notifier reference to determine the best way to handle
		the call. */
	const char *type;
	/** Data whose type depends on the type of notification being sent.
		For instance, notifiers related to plug-in tools receive
		the \c #AIToolHandle to which the notifier message refers. */
	void *notifyData;
	/** Indicates that the notification is being sent parallely to multiple plugins.
		Refer \c #kNotifierOptionSupportsParallelDispatch. */
	AIBoolean isParallelDispatch;
};

/** The contents of a notifier message. */
enum AINotifierOptions : ai::uint32
{
	/** The message data */
	kNotifierOptionNone = 0,
	/** For internal use only.
		Indicates that the plugin can handle this notification in parallel with other plugins.
		During a parallel dispatch, the plugin may update its cached data but must not change anyhing
		in artwork DOM or do any UI updates.
		The parallel notification may be received by the plugin either on background or main thread.
		The plugin should use \c #isParallelDispatch in \c #AINotifierMessage to check
		which case it is.
		After parallel dispatch, the plugin will receive the same notification in traditional serial way also
		in the main thread. The plugin should do UI updates at that time. 
		Use of this flag during adding a notifier does not guarentee that the parallel dispatch notification will
		be sent to the plugin, so the plugin code should write fallback in the main thread to update cache. */
	kNotifierOptionSupportsParallelDispatch = 1 << 0,
};


/*******************************************************************************
 **
 **	Suite Record
 **
 **/

/**	@ingroup Suites
	This suite provides functions that allow your plug-in to use Illustrator's
	event notification system.

	Illustrator sends a notifier to plug-in to inform it of an events for which
	it has registered an interest. These functions allow you to request particular
	notifications, turn them on and off, and find out which plug-ins are listening
	to notifications.

	Notifiers can be used by themselves as a background process or with other plug-in
	types, such as a menus or windows, to learn when an	update is needed.
	Specific notifier type definitions are not a part of this suite,
	but are found in the suites to which they are related. For instance, the \c #AIArtSuite
	defines the \c #kAIArtSelectionChangedNotifier and \c #kAIArtPropertiesChangedNotifier.

	Illustrator notifies a plug-in of an event by sending a message to its main entry
	point with the caller \c #kCallerAINotify and selector \c #kSelectorAINotify. The
	message data is defined by \c #AINotifierMessage.

	Notifications are sometimes sent during an idle loop, so a plug-in should not rely
	on receiving them synchronously with a document state change. Some notifications
	are sent when something might have changed, but do not guarantee that something has
	changed.

 	\li Acquire this suite using \c #SPBasicSuite::AcquireSuite() with the constants
		\c #kAINotifierSuite and \c #kAINotifierVersion.
  */
struct AINotifierSuite {

	/**  Registers interest in a notification. Use at startup.
		 	@param self This plug-in.
			@param name The unique identifying name of this plug-in.
			@param type The notification type, as defined in the related suite.
				See @ref Notifiers.
			@param notifier [out] A buffer in which to return the notifier reference.
				 If your plug-in installs multiple notifications, store this in '
				 \c globals to compare when receiving a notification.
		*/
	AIAPI AIErr (*AddNotifier) ( SPPluginRef self, const char *name, const char *type,
				AINotifierHandle *notifier );

	/**  Registers interest in a notification. Use at startup.
			@param self This plug-in.
			@param name The unique identifying name of this plug-in.
			@param type The notification type, as defined in the related suite.
				See @ref Notifiers.
			@param options A set of option flags.
				 A logical OR of \c #AINotifierOptions.
			@param notifier [out] A buffer in which to return the notifier reference.
				 If your plug-in installs multiple notifications, store this in '
				 \c globals to compare when receiving a notification.
		*/
	AIAPI AIErr (*AddNotifierEx) ( SPPluginRef self, const char *name, const char *type,
				AINotifierOptions options, AINotifierHandle *notifier );

	/** Retrieves the unique identifying name of a notifier.
			@param notifier The notifier reference.
			@param name [out] A buffer in which to return the name. Do not modify this string.
		*/
	AIAPI AIErr (*GetNotifierName) ( AINotifierHandle notifier, const char **name );

	/** Retrieves the type of a notifier.
			@param notifier The notifier reference.
			@param type [out] A buffer in which to return the type. Do not modify this string.
		*/
	AIAPI AIErr (*GetNotifierType) ( AINotifierHandle notifier, const char **type );

	/** Retrieves the options of a notifier.
			@param notifier The notifier reference.
			@param options [out] The notifier options.
		*/
	AIAPI AIErr (*GetNotifierOptions) ( AINotifierHandle notifier,
				AINotifierOptions &options );

	/** Reports whether a notifier is active. An active notifier can receive events.
			@param notifier The notifier reference.
			@param active [out] A buffer in which to return true if the notifier is active.
		*/
	AIAPI AIErr (*GetNotifierActive) ( AINotifierHandle notifier, AIBoolean *active );

	/** Turns a notifier on or off. When a notifier is active (on) it receives
		notification of events from Illustrator.
			@param notifier The notifier reference.
			@param active True to turn the notifier on, false to turn it off.
		*/
	AIAPI AIErr (*SetNotifierActive) ( AINotifierHandle notifier, AIBoolean active );

	/** Retrieves the plug-in that installed a notifier. You can pass this
		reference to functions in the \c #AIPluginSuite.
			@param notifier The notifier reference.
			@param plugin [out] A buffer in which to return the plug-in reference.
		*/
	AIAPI AIErr (*GetNotifierPlugin) ( AINotifierHandle notifier, SPPluginRef *plugin );

	/** Gets the number of installed notifiers. Use this with \c #GetNthNotifier()
		to iterate through notifiers.
			@param count [out] A buffer in which to return the number of notification registrations.
		*/
	AIAPI AIErr (*CountNotifiers) ( ai::int32 *count );

	/** Retrieves a notifier by index position. Use this with \c #CountNotifiers()
		to iterate through notifiers.
			@param n The 0-based position index of the notifier.
			@param notifier [out] A buffer in which to return the notifier reference.
		*/
	AIAPI AIErr (*GetNthNotifier) ( ai::int32 n, AINotifierHandle *notifier );

	/** Broadcasts a notification to all the plug-ins that have subscribed to the
		specified type. Use this mechanism for communicating with other plug-ins.
	 		@param type	The notification type, as defined in the related suite.
				See @ref Notifiers.
			@param notifyData Type-specific data to pass to the notification.
	 	*/
	AIAPI AIErr (*Notify) ( const char *type, void *notifyData );

};


#include "AIHeaderEnd.h"

#endif
