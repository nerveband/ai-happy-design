#include "IllustratorSDK.h"
#include "[PROJECT_NAME]Plugin.h"
#include "[PROJECT_NAME]Suites.h"

Plugin* AllocatePlugin(SPPluginRef pluginRef)
{
	return new [PROJECT_NAME]Plugin(pluginRef);
}

void FixupReload(Plugin* plugin)
{
	[PROJECT_NAME]Plugin::FixupVTable(([PROJECT_NAME]Plugin*) plugin);
}

[PROJECT_NAME]Plugin::[PROJECT_NAME]Plugin(SPPluginRef pluginRef) :
	Plugin(pluginRef)
{
	strncpy(fPluginName, k[PROJECT_NAME]PluginName, kMaxStringLength);
}

[PROJECT_NAME]Plugin::~[PROJECT_NAME]Plugin()
{
}

ASErr [PROJECT_NAME]Plugin::StartupPlugin( SPInterfaceMessage *message )
{
	ASErr error = kNoErr;
	error = Plugin::StartupPlugin(message);
	sAIUser->MessageAlert(ai::UnicodeString("Hello from [PROJECT_NAME]!"));
	return error;
}

ASErr [PROJECT_NAME]Plugin::ShutdownPlugin( SPInterfaceMessage *message )
{
	ASErr error = kNoErr;
	sAIUser->MessageAlert(ai::UnicodeString("Goodbye from [PROJECT_NAME]!"));
	error = Plugin::ShutdownPlugin(message);
	return error;
}
