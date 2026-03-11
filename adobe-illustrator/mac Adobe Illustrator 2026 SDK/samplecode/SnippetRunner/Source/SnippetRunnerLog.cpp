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

// General includes:
#include <vector>
#include <string>
#include "SDKErrors.h"
#include "SnippetRunnerPlugin.h"
#include "SnippetRunnerSuites.h"
#include "SnpFileHelper.h"
#include "AIDataFilter.h"
#include "SnippetRunnerPreferences.h"
#include "SnippetRunnerLog.h"

extern SnippetRunnerPlugin*	gPlugin;

SnippetRunnerLog* SnippetRunnerLog::fInstance = NULL;

/*
*/
SnippetRunnerLog::SnippetRunnerLog() :  fCurrentIndentLevel(0)
{
}

/*
*/
SnippetRunnerLog::~SnippetRunnerLog()
{
}

/*
*/
SnippetRunnerLog* SnippetRunnerLog::Instance()
{
	if (NULL == SnippetRunnerLog::fInstance) {
		SnippetRunnerLog::fInstance = new SnippetRunnerLog();
	}
	return SnippetRunnerLog::fInstance;
}

/*
*/
void SnippetRunnerLog::DeleteInstance()
{
	if (NULL != SnippetRunnerLog::fInstance) {
		delete SnippetRunnerLog::fInstance;
		SnippetRunnerLog::fInstance = NULL;
	}
}

/*
*/
void SnippetRunnerLog::Write(const LogLine& logline)
{
	ai::UnicodeString message;
	const ai::UnicodeString kIndentString("  ");
	for (ASInt32 i = 0; i < this->fCurrentIndentLevel; i++) {
		message += kIndentString;
	}
	message += logline;
	this->fLog.push_back(message);
}

/*
*/
void SnippetRunnerLog::Write(const std::string& logline)
{
	this->Write(ai::UnicodeString(logline));
}

/*
*/
void SnippetRunnerLog::Write(char const* logline)
{
	this->Write(ai::UnicodeString(logline));
}

/*
*/
void SnippetRunnerLog::WritePrintf(const char *format, ... )
{
	static char buf[4096];
	va_list		argptr;
	va_start ( argptr, format );
#ifdef MAC_ENV
	int len = vsnprintf ( buf, sizeof(buf), format, argptr );
#endif
#ifdef WIN_ENV
	int len = _vsnprintf ( buf, sizeof(buf), format, argptr );
#endif
	va_end ( argptr );
	if(len > 0) {
		this->Write(buf);
	}
}

/*
*/
void SnippetRunnerLog::BeginIndent()
{
	fCurrentIndentLevel++;
}

/*
*/
void SnippetRunnerLog::EndIndent()
{
	if (fCurrentIndentLevel > 0) {
		fCurrentIndentLevel--;
	}
}

/*
*/
void SnippetRunnerLog::Save()
{
	ASErr status = kNoErr;
	AIDataFilter* filter = NULL;
	try {
		// Get path to wherever the user wants to save the log.
		ai::FilePath ioFilePath = SnippetRunnerPreferences::Instance()->GetAssetsFolderPath();
		ai::UnicodeString defaultName("SnippetRunnerLog.txt");
		SnpFileHelper fileHelper;
		status = fileHelper.PutTextFileDialog(ai::UnicodeString("Choose log file location"), defaultName, ioFilePath);
		aisdk::check_ai_error(status);

		// Create a new data filter.
        status = sAIDataFilter->NewFileDataFilter(ioFilePath, "write", 'ART5', 'TEXT', &filter);
		aisdk::check_ai_error(status);
        status = sAIDataFilter->LinkDataFilter(NULL, filter);
		aisdk::check_ai_error(status);

#if WIN_ENV
		const char* lineEnding = "\r\n";
		size_t lineEndingCount = 2;
#else
		const char* lineEnding = "\n";
		size_t lineEndingCount = 1;
#endif // WIN_ENV

		// Write the log data to the file via the filter.
		for (SnippetRunnerLog::Log::const_iterator iter = this->fLog.begin(); iter != this->fLog.end(); iter++) {
			std::string line = iter->as_UTF8();
			size_t count = line.size();
			status = sAIDataFilter->WriteDataFilter(filter, const_cast<char*>(line.c_str()), &count);
			aisdk::check_ai_error(status);
			status = sAIDataFilter->WriteDataFilter(filter, const_cast<char*>(lineEnding), &lineEndingCount);
			aisdk::check_ai_error(status);
		}
	}
	catch (...) {
	}
	if (filter) {
		sAIDataFilter->UnlinkDataFilter(filter, NULL);
	}
}

/*
*/
const SnippetRunnerLog::Log& SnippetRunnerLog::GetLog()
{
	return this->fLog;
}

/*
*/
void SnippetRunnerLog::Clear()
{
	this->fLog.clear();
}

/*
*/
void SnippetRunnerLog::NotifyLogChanged()
{
	if (gPlugin) {
		gPlugin->NotifyLogChanged();
	}
}

// End SnippetRunnerLog.cpp
