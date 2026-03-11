#include "IllustratorSDK.h"
#include "[PROJECT_NAME]Suites.h"

// Suite externs
extern "C"
{
	SPBlocksSuite*			sSPBlocks = NULL;
	AIUnicodeStringSuite*	sAIUnicodeString = NULL;
}

// Import suites
ImportSuite gImportSuites[] = 
{
	kSPBlocksSuite, kSPBlocksSuiteVersion, &sSPBlocks,
	kAIUnicodeStringSuite, kAIUnicodeStringVersion, &sAIUnicodeString,
	nil, 0, nil
};
