(function () {
  function stringify(value) {
    try {
      return JSON.stringify(value);
    } catch (err) {
      return JSON.stringify({
        v: "1.0",
        id: (value && value.id) || "",
        ok: false,
        error: "failed to stringify bridge payload: " + err
      });
    }
  }

  var payload = __AHD_PAYLOAD__;
  try {
    if (payload.mode === "script") {
      var result = eval(payload.script);
      return stringify({
        v: "1.0",
        id: payload.id,
        ok: true,
        result: result || {},
        warnings: []
      });
    }

    var raw = app.sendScriptMessage("AHDIllustrator", payload.selector, stringify(payload.request));
    if (!raw || raw === "") {
      return stringify({
        v: "1.0",
        id: payload.request.id,
        ok: false,
        error: "empty plugin response",
        warnings: []
      });
    }
    return raw;
  } catch (err) {
    return stringify({
      v: "1.0",
      id: (payload.request && payload.request.id) || payload.id,
      ok: false,
      error: String(err),
      warnings: []
    });
  }
}());
