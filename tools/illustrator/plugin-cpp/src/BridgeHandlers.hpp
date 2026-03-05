#pragma once

#include <string>

namespace ahd {

class BridgeHandlers {
 public:
  static std::string Handle(const std::string& selector, const std::string& payload);

 private:
  static std::string HandleCapabilities();
  static std::string HandleVersion();
  static std::string HandleExec(const std::string& payload);
  static std::string HandleInspect(const std::string& payload);
  static std::string MakeOK(const std::string& request_id, const std::string& result_json);
  static std::string MakeError(const std::string& request_id, const std::string& code, const std::string& message);
};

}  // namespace ahd
