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

#ifndef _SnippetRunnerLog_
#define _SnippetRunnerLog_

/** Logs messages from snippets (implemented as Gamma singleton). 

	To write to the log:
	<code>
		SnippetRunnerLog* log = SnippetRunnerLog::Instance();
		log->Write("Hello");
		log->Write(ai::UnicodeString("MyUnicodeData"));
		log->WritePrintf("%s %d", "number=", 125);
	</code>

	@ingroup snippetrunner
*/
class SnippetRunnerLog {
public:
	/** Returns the global instance of this singleton.
		@return the global instance of this singleton class.
	*/
	static SnippetRunnerLog* Instance();

	/**	Deletes the singleton instance of this class.
	*/
	static void DeleteInstance();

	/** Type used to store one log message.
	*/
	typedef ai::UnicodeString LogLine;

	/** Collection used to store the log messages.
	*/
	typedef std::vector<LogLine> Log;

	/** Writes given line to the log.
		@param logline IN the line to write to the log.
	*/
	void Write(const LogLine& logline);

	/** Writes given line to the log.
		@param logline IN the line to write to the log.
	*/
	void Write(const std::string& logline);

	/** Writes given line to the log.
		@param logline IN the line to write to the log.
	*/
	void Write(char const* logline);

	/** Writes given line to the log using pringf format with variable set of arguments.
		@param format IN various formats to use.
		@param ... IN undefined number of input parameters to write to the log.
	*/
	void WritePrintf(const char *format, ... );

	/** Indents the log.
	*/
	void BeginIndent();

	/** Un-indents the log.
	*/
	void EndIndent();

	/** Saves the log to a file.
	*/
	void Save();

	/** Returns the log data.
		@return the log.
	*/
	const Log& GetLog();

	/** Clears the log.
	*/
	void Clear();

	/** Broadcast notification to the view that the log has changed.
	*/
	void NotifyLogChanged();

public:
	/** Stack based helper class that indents the log on construction and
		un-indents on destruction.
		@ingroup snippetrunner
	*/
	class Indent {
	public:
		/** Constructor.
		*/
		Indent() {SnippetRunnerLog::Instance()->BeginIndent();}

		/** Destructor.
		*/
		~Indent() {SnippetRunnerLog::Instance()->EndIndent();}
	};

private:
	/** Constructor is hidden - instantiate using SnippetRunnerLog::Instance.
	*/
	SnippetRunnerLog();

	/** Destructor is hidden - destroy using SnippetRunnerLog::DeleteInstance.
	*/
	~SnippetRunnerLog();

	/** Copy constructor is hidden.
		@param rhs
	*/
	SnippetRunnerLog(SnippetRunnerLog const&);

	/** Assignment operator is hidden.
		@param IN SnippetRunnerLog object to assign to lhs.
		@return copy of SnippetRunnerLog passed in.
	*/
	SnippetRunnerLog& operator=(SnippetRunnerLog const&);  // assign op is hidden

private:
	/** Points to the singleton instance of this class.
	*/
	static SnippetRunnerLog* fInstance;

	/** Stores the log data.
	*/
	Log fLog;

	/** Stores how indented the log is at any one time.
	*/
	ASInt32 fCurrentIndentLevel;
};

#endif //_SnippetRunnerLog_

// End SnippetRunnerLog.h
