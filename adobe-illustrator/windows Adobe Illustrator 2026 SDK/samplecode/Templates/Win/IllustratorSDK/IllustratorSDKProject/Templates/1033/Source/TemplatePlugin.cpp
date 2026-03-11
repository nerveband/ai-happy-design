#include "IllustratorSDK.h"
#include "[!output PROJECT_NAME]Plugin.h"
#include "[!output PROJECT_NAME]Suites.h"

Plugin* AllocatePlugin(SPPluginRef pluginRef)
{
	return new [!output PROJECT_NAME]Plugin(pluginRef);
}

void FixupReload(Plugin* plugin)
{
	[!output PROJECT_NAME]Plugin::FixupVTable(([!output PROJECT_NAME]Plugin*) plugin);
}

[!output PROJECT_NAME]Plugin::[!output PROJECT_NAME]Plugin(SPPluginRef pluginRef) :
	Plugin(pluginRef)
{
	strncpy(fPluginName, k[!output PROJECT_NAME]PluginName, kMaxStringLength);
}

[!output PROJECT_NAME]Plugin::~[!output PROJECT_NAME]Plugin()
{
}

ASErr [!output PROJECT_NAME]Plugin::StartupPlugin( SPInterfaceMessage *message )
{
	ASErr error = kNoErr;
	error = Plugin::StartupPlugin(message);
	sAIUser->MessageAlert(ai::UnicodeString("Hello from [!output PROJECT_NAME]!"));
	return error;
}

ASErr [!output PROJECT_NAME]Plugin::ShutdownPlugin( SPInterfaceMessage *message )
{
	ASErr error = kNoErr;
	sAIUser->MessageAlert(ai::UnicodeString("Goodbye from [!output PROJECT_NAME]!"));
	error = Plugin::ShutdownPlugin(message);
	return error;
}
