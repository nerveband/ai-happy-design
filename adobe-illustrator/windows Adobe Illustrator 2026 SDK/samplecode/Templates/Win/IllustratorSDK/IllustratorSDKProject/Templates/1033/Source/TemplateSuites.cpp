#include "IllustratorSDK.h"
#include "[!output PROJECT_NAME]Suites.h"

// Suite externs
extern "C"
{
	SPBlocksSuite*			sSPBlocks = nullptr;
	AIUnicodeStringSuite*	sAIUnicodeString = nullptr;
}

// Import suites
ImportSuite gImportSuites[] = 
{
	kSPBlocksSuite, kSPBlocksSuiteVersion, &sSPBlocks,
	kAIUnicodeStringSuite, kAIUnicodeStringVersion, &sAIUnicodeString,
	nullptr, 0, nullptr
};
