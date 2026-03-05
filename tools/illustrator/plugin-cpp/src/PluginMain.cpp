#include "BridgeHandlers.hpp"

#include <string>

namespace {

std::string g_last_response;

}  // namespace

extern "C" const char* AHD_HandleScriptMessage(const char* selector, const char* payload) {
  const std::string safe_selector = selector == nullptr ? "" : selector;
  const std::string safe_payload = payload == nullptr ? "" : payload;
  g_last_response = ahd::BridgeHandlers::Handle(safe_selector, safe_payload);
  return g_last_response.c_str();
}

extern "C" const char* AHD_PluginVersion() {
  static const std::string version = "0.1.0";
  return version.c_str();
}

/*
Illustrator SDK integration note:

The production plugin entrypoint should forward sendScriptMessage selectors into
AHD_HandleScriptMessage. The exported function keeps the bridge logic isolated
from SDK-specific bootstrapping so the JSON protocol can be tested separately.
*/
