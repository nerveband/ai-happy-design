#include "BridgeHandlers.hpp"

#include <sstream>

namespace {

std::string EscapeJSON(const std::string& value) {
  std::ostringstream stream;
  for (char ch : value) {
    switch (ch) {
      case '\\':
        stream << "\\\\";
        break;
      case '"':
        stream << "\\\"";
        break;
      case '\n':
        stream << "\\n";
        break;
      case '\r':
        stream << "\\r";
        break;
      case '\t':
        stream << "\\t";
        break;
      default:
        stream << ch;
        break;
    }
  }
  return stream.str();
}

std::string ExtractID(const std::string& payload) {
  const std::string key = "\"id\":\"";
  const auto start = payload.find(key);
  if (start == std::string::npos) {
    return "";
  }
  const auto value_start = start + key.size();
  const auto value_end = payload.find('"', value_start);
  if (value_end == std::string::npos) {
    return "";
  }
  return payload.substr(value_start, value_end - value_start);
}

}  // namespace

namespace ahd {

std::string BridgeHandlers::Handle(const std::string& selector, const std::string& payload) {
  if (selector == "ahd.capabilities") {
    return HandleCapabilities();
  }
  if (selector == "ahd.version") {
    return HandleVersion();
  }
  if (selector == "ahd.exec") {
    return HandleExec(payload);
  }
  if (selector == "ahd.inspect") {
    return HandleInspect(payload);
  }
  return MakeError(ExtractID(payload), "UNSUPPORTED_SELECTOR", "unknown selector: " + selector);
}

std::string BridgeHandlers::HandleCapabilities() {
  return MakeOK("", R"({"selectors":["ahd.capabilities","ahd.exec","ahd.inspect","ahd.version"],"pluginMode":"sendScriptMessage"})");
}

std::string BridgeHandlers::HandleVersion() {
  return MakeOK("", R"({"plugin":"AHDIllustrator","bridgeVersion":"0.1.0"})");
}

std::string BridgeHandlers::HandleExec(const std::string& payload) {
  if (payload.empty()) {
    return MakeError("", "VALIDATION_ERROR", "exec payload cannot be empty");
  }
  return MakeOK(ExtractID(payload), std::string(R"({"accepted":true,"echo":")") + EscapeJSON(payload) + "\"}");
}

std::string BridgeHandlers::HandleInspect(const std::string& payload) {
  return MakeOK(ExtractID(payload), std::string(R"({"inspectStub":true,"echo":")") + EscapeJSON(payload) + "\"}");
}

std::string BridgeHandlers::MakeOK(const std::string& request_id, const std::string& result_json) {
  return std::string("{\"v\":\"1.0\",\"id\":\"") + EscapeJSON(request_id) + "\",\"ok\":true,\"result\":" + result_json + ",\"warnings\":[]}";
}

std::string BridgeHandlers::MakeError(const std::string& request_id, const std::string& code, const std::string& message) {
  return std::string("{\"v\":\"1.0\",\"id\":\"") + EscapeJSON(request_id) +
         "\",\"ok\":false,\"error\":\"" + EscapeJSON(code + ": " + message) + "\",\"warnings\":[]}";
}

}  // namespace ahd
