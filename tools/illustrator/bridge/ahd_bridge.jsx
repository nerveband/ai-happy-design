(function () {
  function ahdQuote(value) {
    var str = String(value);
    var out = '"';
    var i;
    var ch;
    var code;
    var hex;
    for (i = 0; i < str.length; i++) {
      ch = str.charAt(i);
      code = str.charCodeAt(i);
      if (ch === "\\") {
        out += "\\\\";
      } else if (ch === "\"") {
        out += "\\\"";
      } else if (ch === "\b") {
        out += "\\b";
      } else if (ch === "\f") {
        out += "\\f";
      } else if (ch === "\n") {
        out += "\\n";
      } else if (ch === "\r") {
        out += "\\r";
      } else if (ch === "\t") {
        out += "\\t";
      } else if (code < 32) {
        hex = code.toString(16);
        out += "\\u" + ("0000" + hex).slice(-4);
      } else {
        out += ch;
      }
    }
    return out + '"';
  }

  function ahdIsArray(value) {
    return Object.prototype.toString.call(value) === "[object Array]";
  }

  function ahdStringify(value, seen) {
    var i;
    var keys;
    var parts;
    var item;
    if (!seen) {
      seen = [];
    }
    if (value === null || typeof value === "undefined") {
      return "null";
    }
    if (typeof value === "string") {
      return ahdQuote(value);
    }
    if (typeof value === "number") {
      return isFinite(value) ? String(value) : "null";
    }
    if (typeof value === "boolean") {
      return value ? "true" : "false";
    }
    if (typeof value === "function") {
      return "null";
    }
    for (i = 0; i < seen.length; i++) {
      if (seen[i] === value) {
        return ahdQuote("[Circular]");
      }
    }
    seen.push(value);
    if (ahdIsArray(value)) {
      parts = [];
      for (i = 0; i < value.length; i++) {
        parts.push(ahdStringify(value[i], seen));
      }
      seen.pop();
      return "[" + parts.join(",") + "]";
    }
    parts = [];
    keys = [];
    for (item in value) {
      if (Object.prototype.hasOwnProperty.call(value, item)) {
        keys.push(item);
      }
    }
    keys.sort();
    for (i = 0; i < keys.length; i++) {
      item = keys[i];
      if (typeof value[item] === "undefined" || typeof value[item] === "function") {
        continue;
      }
      parts.push(ahdQuote(item) + ":" + ahdStringify(value[item], seen));
    }
    seen.pop();
    return "{" + parts.join(",") + "}";
  }

  function ahdSuccess(id, result) {
    return ahdStringify({
      v: "1.0",
      id: id,
      ok: true,
      result: typeof result === "undefined" ? {} : result,
      warnings: []
    });
  }

  function ahdError(id, message) {
    return "{\"v\":\"1.0\",\"id\":" + ahdQuote(id || "") + ",\"ok\":false,\"error\":" + ahdQuote(String(message || "unknown bridge error")) + ",\"warnings\":[]}";
  }

  var payload = __AHD_PAYLOAD__;
  try {
    if (payload.mode === "script") {
      return ahdSuccess(payload.id, eval(payload.script));
    }

    var requestBody = ahdStringify(payload.request || {});
    var raw = app.sendScriptMessage("AHDIllustrator", payload.selector, requestBody);
    if (!raw || raw === "") {
      return ahdError((payload.request && payload.request.id) || payload.id, "empty plugin response");
    }
    return raw;
  } catch (err) {
    return ahdError((payload.request && payload.request.id) || payload.id, err);
  }
}());
