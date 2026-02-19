package extract

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestParseStyleBlocks(t *testing.T) {
	html := `<html><head><style>
		.bg-dark { background: linear-gradient(150deg, #0C1E2C 0%, #14344A 100%); }
		.bg-navy { background: #14344A; }
		.h { font-family: 'Poppins', sans-serif; font-weight: 700; color: #FAFCFB; }
		.h-xl { font-size: 1.6rem; }
		.ey { font-family: 'Poppins'; text-transform: uppercase; letter-spacing: .08em; }
	</style></head><body></body></html>`

	styles := parseStyleBlocks(html)

	// Check .bg-dark
	bgDark, ok := styles[".bg-dark"]
	if !ok {
		t.Fatal(".bg-dark not found in styles")
	}
	if !strings.Contains(bgDark["background"], "linear-gradient") {
		t.Errorf(".bg-dark background = %q, want linear-gradient", bgDark["background"])
	}

	// Check .bg-navy
	bgNavy, ok := styles[".bg-navy"]
	if !ok {
		t.Fatal(".bg-navy not found in styles")
	}
	if bgNavy["background"] != "#14344A" {
		t.Errorf(".bg-navy background = %q, want #14344A", bgNavy["background"])
	}

	// Check .h
	h, ok := styles[".h"]
	if !ok {
		t.Fatal(".h not found in styles")
	}
	if h["font-weight"] != "700" {
		t.Errorf(".h font-weight = %q, want 700", h["font-weight"])
	}
	if h["color"] != "#FAFCFB" {
		t.Errorf(".h color = %q, want #FAFCFB", h["color"])
	}

	// Check .ey
	ey, ok := styles[".ey"]
	if !ok {
		t.Fatal(".ey not found in styles")
	}
	if ey["text-transform"] != "uppercase" {
		t.Errorf(".ey text-transform = %q, want uppercase", ey["text-transform"])
	}
}

func TestParseInlineStyle(t *testing.T) {
	s := parseInlineStyle("color: red; font-size: 16px; font-weight: bold")
	if s["color"] != "red" {
		t.Errorf("color = %q, want red", s["color"])
	}
	if s["font-size"] != "16px" {
		t.Errorf("font-size = %q, want 16px", s["font-size"])
	}
	if s["font-weight"] != "bold" {
		t.Errorf("font-weight = %q, want bold", s["font-weight"])
	}
}

func TestParseLinearGradient(t *testing.T) {
	tests := []struct {
		input     string
		wantAngle float64
		wantStops int
		wantErr   bool
	}{
		{
			input:     "linear-gradient(150deg, #0C1E2C 0%, #14344A 100%)",
			wantAngle: 150,
			wantStops: 2,
		},
		{
			input:     "linear-gradient(145deg, #14344A 0%, #1A5F7A 45%, #1B9B6E 100%)",
			wantAngle: 145,
			wantStops: 3,
		},
		{
			input:     "linear-gradient(to right, red 0%, blue 100%)",
			wantAngle: 90,
			wantStops: 2,
		},
		{
			input:   "not-a-gradient",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		angle, stops, err := parseLinearGradient(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("parseLinearGradient(%q): expected error", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseLinearGradient(%q): %v", tt.input, err)
			continue
		}
		if angle != tt.wantAngle {
			t.Errorf("parseLinearGradient(%q) angle = %v, want %v", tt.input, angle, tt.wantAngle)
		}
		if len(stops) != tt.wantStops {
			t.Errorf("parseLinearGradient(%q) stops = %d, want %d", tt.input, len(stops), tt.wantStops)
		}
	}

	// Verify stop details for the 3-stop gradient
	_, stops, _ := parseLinearGradient("linear-gradient(145deg, #14344A 0%, #1A5F7A 45%, #1B9B6E 100%)")
	if stops[0].Color != "#14344A" {
		t.Errorf("stop[0].Color = %q, want #14344A", stops[0].Color)
	}
	if stops[0].Position != 0.0 {
		t.Errorf("stop[0].Position = %v, want 0.0", stops[0].Position)
	}
	if stops[1].Color != "#1A5F7A" {
		t.Errorf("stop[1].Color = %q, want #1A5F7A", stops[1].Color)
	}
	if stops[1].Position != 0.45 {
		t.Errorf("stop[1].Position = %v, want 0.45", stops[1].Position)
	}
	if stops[2].Color != "#1B9B6E" {
		t.Errorf("stop[2].Color = %q, want #1B9B6E", stops[2].Color)
	}
	if stops[2].Position != 1.0 {
		t.Errorf("stop[2].Position = %v, want 1.0", stops[2].Position)
	}
}

func TestCSSColorToHex(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"#FF0000", "#FF0000"},
		{"#ff0000", "#FF0000"},
		{"#f00", "#FF0000"},
		{"#FF000080", "#FF0000"}, // 8-digit hex, alpha dropped
		{"rgb(255, 0, 0)", "#FF0000"},
		{"rgba(255, 0, 0, 0.5)", "#FF0000"},
		{"rgb(0, 128, 0)", "#008000"},
		{"black", "#000000"},
		{"white", "#FFFFFF"},
		{"red", "#FF0000"},
		{"transparent", ""},
		{"", ""},
		{"not-a-color", ""},
	}

	for _, tt := range tests {
		got := cssColorToHex(tt.input)
		if got != tt.want {
			t.Errorf("cssColorToHex(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestRemToPx(t *testing.T) {
	tests := []struct {
		input    string
		base     float64
		expected float64
	}{
		{"1.6rem", 16, 25.6},
		{"2rem", 16, 32},
		{"1em", 16, 16},
		{"24px", 16, 24},
		{"100", 16, 100},
	}

	for _, tt := range tests {
		got := remToPx(tt.input, tt.base)
		if got != tt.expected {
			t.Errorf("remToPx(%q, %v) = %v, want %v", tt.input, tt.base, got, tt.expected)
		}
	}
}

func TestCSSFontWeight(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"400", "Regular"},
		{"700", "Bold"},
		{"bold", "Bold"},
		{"normal", "Regular"},
		{"600", "SemiBold"},
		{"300", "Light"},
	}
	for _, tt := range tests {
		got := cssFontWeight(tt.input)
		if got != tt.want {
			t.Errorf("cssFontWeight(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFromHTMLMiniSlide(t *testing.T) {
	htmlDoc := `<!DOCTYPE html>
<html>
<head>
<style>
  .slide { width: 280px; height: 350px; }
  .bg-dark { background: linear-gradient(150deg, #0C1E2C 0%, #14344A 100%); }
  .h { font-family: 'Poppins', sans-serif; font-weight: 700; color: #FAFCFB; }
  .h-xl { font-size: 1.6rem; }
  .ey { font-family: 'Poppins'; text-transform: uppercase; letter-spacing: .08em; color: #87CEEB; }
  .bar-grn { background: #22C55E; }
</style>
</head>
<body>
<div class="slides-row">
  <div class="slide bg-dark">
    <div class="s-num">1 / 7</div>
    <div class="s-c">
      <div class="ey">Ramadan 2026</div>
      <div class="h h-xl">The Care<br>Behind<br>the Care</div>
      <div class="bar bar-grn"></div>
      <div class="s-url">muslimchaplains.net</div>
    </div>
  </div>
</div>
</body>
</html>`

	ops, err := FromHTML(strings.NewReader(htmlDoc), Options{
		CanvasWidth:  1080,
		CanvasHeight: 1350,
	})
	if err != nil {
		t.Fatalf("FromHTML failed: %v", err)
	}

	if len(ops) != 1 {
		t.Fatalf("expected 1 slide op, got %d", len(ops))
	}

	op := ops[0]
	if op["command"] != "slide" {
		t.Errorf("command = %v, want slide", op["command"])
	}
	if op["name"] != "s1" {
		t.Errorf("name = %v, want s1", op["name"])
	}

	params, ok := op["params"].(map[string]interface{})
	if !ok {
		t.Fatal("params not a map")
	}

	// Should have canvas size
	if params["canvas"] != "1080x1350" {
		t.Errorf("canvas = %v, want 1080x1350", params["canvas"])
	}

	// Should have gradient
	grad, ok := params["gradient"].(map[string]interface{})
	if !ok {
		t.Fatal("gradient not present")
	}
	if grad["type"] != "LINEAR" {
		t.Errorf("gradient type = %v, want LINEAR", grad["type"])
	}
	if grad["angle"].(float64) != 150 {
		t.Errorf("gradient angle = %v, want 150", grad["angle"])
	}

	// Check elements
	elements, ok := params["elements"].([]interface{})
	if !ok {
		t.Fatal("elements not present")
	}

	// Expect: counter, eyebrow, headline, bar, url
	if len(elements) < 4 {
		t.Fatalf("expected at least 4 elements, got %d", len(elements))
	}

	// First element should be counter
	counter, ok := elements[0].(map[string]interface{})
	if !ok {
		t.Fatal("element 0 not a map")
	}
	if counter["type"] != "counter" {
		t.Errorf("element 0 type = %v, want counter", counter["type"])
	}
	if counter["current"] != "1" {
		t.Errorf("counter current = %v, want 1", counter["current"])
	}
	if counter["total"] != "7" {
		t.Errorf("counter total = %v, want 7", counter["total"])
	}

	// Second element should be eyebrow
	eyebrow, ok := elements[1].(map[string]interface{})
	if !ok {
		t.Fatal("element 1 not a map")
	}
	if eyebrow["type"] != "eyebrow" {
		t.Errorf("element 1 type = %v, want eyebrow", eyebrow["type"])
	}
	if eyebrow["text"] != "Ramadan 2026" {
		t.Errorf("eyebrow text = %v, want 'Ramadan 2026'", eyebrow["text"])
	}
	if eyebrow["fontFamily"] != "Poppins" {
		t.Errorf("eyebrow fontFamily = %v, want Poppins", eyebrow["fontFamily"])
	}

	// Third element should be headline
	headline, ok := elements[2].(map[string]interface{})
	if !ok {
		t.Fatal("element 2 not a map")
	}
	if headline["type"] != "headline" {
		t.Errorf("element 2 type = %v, want headline", headline["type"])
	}
	expectedText := "The Care\nBehind\nthe Care"
	if headline["text"] != expectedText {
		t.Errorf("headline text = %q, want %q", headline["text"], expectedText)
	}
	if headline["tier"] != "hero" {
		t.Errorf("headline tier = %v, want hero (from h-xl class)", headline["tier"])
	}
	if headline["fontStyle"] != "Bold" {
		t.Errorf("headline fontStyle = %v, want Bold", headline["fontStyle"])
	}

	// Fourth element should be bar
	bar, ok := elements[3].(map[string]interface{})
	if !ok {
		t.Fatal("element 3 not a map")
	}
	if bar["type"] != "bar" {
		t.Errorf("element 3 type = %v, want bar", bar["type"])
	}
	if bar["color"] != "#22C55E" {
		t.Errorf("bar color = %v, want #22C55E", bar["color"])
	}

	// Fifth element should be URL
	if len(elements) >= 5 {
		url, ok := elements[4].(map[string]interface{})
		if !ok {
			t.Fatal("element 4 not a map")
		}
		if url["type"] != "url" {
			t.Errorf("element 4 type = %v, want url", url["type"])
		}
		if url["text"] != "muslimchaplains.net" {
			t.Errorf("url text = %v, want muslimchaplains.net", url["text"])
		}
	}
}

func TestFromHTMLMultipleSlides(t *testing.T) {
	htmlDoc := `<html>
<head><style>
  .slide { width: 280px; height: 350px; }
  .bg-dark { background: #0C1E2C; }
  .bg-hero { background: linear-gradient(145deg, #14344A 0%, #1B9B6E 100%); }
  .h { font-weight: 700; color: #FFF; }
</style></head>
<body>
<div class="slides-row">
  <div class="slide bg-dark">
    <div class="s-c"><div class="h h-xl">Slide One</div></div>
  </div>
  <div class="slide bg-hero">
    <div class="s-c"><div class="h h-md">Slide Two</div></div>
  </div>
</div>
</body>
</html>`

	ops, err := FromHTML(strings.NewReader(htmlDoc), Options{
		CanvasWidth:  1080,
		CanvasHeight: 1350,
	})
	if err != nil {
		t.Fatalf("FromHTML failed: %v", err)
	}

	if len(ops) != 2 {
		t.Fatalf("expected 2 slide ops, got %d", len(ops))
	}

	if ops[0]["name"] != "s1" {
		t.Errorf("first slide name = %v, want s1", ops[0]["name"])
	}
	if ops[1]["name"] != "s2" {
		t.Errorf("second slide name = %v, want s2", ops[1]["name"])
	}

	// First slide: solid bg
	p1 := ops[0]["params"].(map[string]interface{})
	if p1["color"] != "#0C1E2C" {
		t.Errorf("slide 1 color = %v, want #0C1E2C", p1["color"])
	}
	if _, hasGrad := p1["gradient"]; hasGrad {
		t.Error("slide 1 should not have gradient")
	}

	// Second slide: gradient bg
	p2 := ops[1]["params"].(map[string]interface{})
	if _, hasGrad := p2["gradient"]; !hasGrad {
		t.Error("slide 2 should have gradient")
	}

	// Check headline tiers
	e1 := p1["elements"].([]interface{})
	h1 := e1[0].(map[string]interface{})
	if h1["tier"] != "hero" {
		t.Errorf("slide 1 headline tier = %v, want hero (h-xl)", h1["tier"])
	}

	e2 := p2["elements"].([]interface{})
	h2 := e2[0].(map[string]interface{})
	if h2["tier"] != "heading" {
		t.Errorf("slide 2 headline tier = %v, want heading (h-md)", h2["tier"])
	}
}

func TestFromHTMLEmailBanner(t *testing.T) {
	htmlDoc := `<html>
<head><style>
  .eb { width: 600px; height: 200px; }
  .bg-dark { background: linear-gradient(150deg, #0C1E2C 0%, #14344A 100%); }
</style></head>
<body>
<div class="email-banners">
  <div class="eb bg-dark">
    <div class="eb-text">
      <div class="eb-h eb-h-xl">The Care Behind the Care</div>
      <div class="eb-sub">AMC Ramadan 2026</div>
    </div>
  </div>
</div>
</body>
</html>`

	ops, err := FromHTML(strings.NewReader(htmlDoc), Options{
		CanvasWidth:  1200,
		CanvasHeight: 400,
	})
	if err != nil {
		t.Fatalf("FromHTML failed: %v", err)
	}

	if len(ops) != 1 {
		t.Fatalf("expected 1 banner op, got %d", len(ops))
	}

	op := ops[0]
	if op["command"] != "banner" {
		t.Errorf("command = %v, want banner", op["command"])
	}
	if op["name"] != "b1" {
		t.Errorf("name = %v, want b1", op["name"])
	}

	params := op["params"].(map[string]interface{})
	if params["canvas"] != "1200x400" {
		t.Errorf("canvas = %v, want 1200x400", params["canvas"])
	}

	// Should have gradient
	if _, ok := params["gradient"]; !ok {
		t.Error("banner should have gradient")
	}

	// Check elements
	elements := params["elements"].([]interface{})
	if len(elements) != 2 {
		t.Fatalf("expected 2 elements (headline + subtitle), got %d", len(elements))
	}

	hl := elements[0].(map[string]interface{})
	if hl["type"] != "headline" {
		t.Errorf("element 0 type = %v, want headline", hl["type"])
	}
	if hl["text"] != "The Care Behind the Care" {
		t.Errorf("headline text = %v", hl["text"])
	}

	sub := elements[1].(map[string]interface{})
	if sub["type"] != "subtitle" {
		t.Errorf("element 1 type = %v, want subtitle", sub["type"])
	}
	if sub["text"] != "AMC Ramadan 2026" {
		t.Errorf("subtitle text = %v", sub["text"])
	}
}

func TestFromHTMLBannerIgnoresSlideCanvasDefaults(t *testing.T) {
	htmlDoc := `<html>
<head><style>
  .eb { width: 600px; height: 200px; }
</style></head>
<body>
<div class="email-banners">
  <div class="eb">
    <div class="eb-text">
      <div class="eb-h eb-h-xl">Banner Headline</div>
      <div class="eb-sub">Banner Subtitle</div>
    </div>
  </div>
</div>
</body>
</html>`

	ops, err := FromHTML(strings.NewReader(htmlDoc), Options{
		CanvasWidth:  1080,
		CanvasHeight: 1350,
	})
	if err != nil {
		t.Fatalf("FromHTML failed: %v", err)
	}

	if len(ops) != 1 {
		t.Fatalf("expected 1 banner op, got %d", len(ops))
	}

	params, ok := ops[0]["params"].(map[string]interface{})
	if !ok {
		t.Fatalf("banner params missing")
	}
	if params["canvas"] != "1200x400" {
		t.Fatalf("banner canvas = %v, want 1200x400", params["canvas"])
	}
}

func TestFromHTMLSolidBackground(t *testing.T) {
	htmlDoc := `<html>
<head><style>
  .slide { width: 280px; height: 350px; }
  .bg-navy { background: #14344A; }
  .h { font-weight: 700; color: #FFF; }
</style></head>
<body>
<div class="slides-row">
  <div class="slide bg-navy">
    <div class="s-c"><div class="h">Hello World</div></div>
  </div>
</div>
</body>
</html>`

	ops, err := FromHTML(strings.NewReader(htmlDoc), Options{CanvasWidth: 1080, CanvasHeight: 1350})
	if err != nil {
		t.Fatalf("FromHTML failed: %v", err)
	}

	if len(ops) != 1 {
		t.Fatalf("expected 1 op, got %d", len(ops))
	}

	params := ops[0]["params"].(map[string]interface{})
	if params["color"] != "#14344A" {
		t.Errorf("color = %v, want #14344A", params["color"])
	}
	if _, hasGrad := params["gradient"]; hasGrad {
		t.Error("solid bg should not have gradient")
	}
}

func TestFromHTMLEmptyInput(t *testing.T) {
	ops, err := FromHTML(strings.NewReader(""), Options{})
	if err != nil {
		t.Fatalf("FromHTML failed on empty input: %v", err)
	}
	if len(ops) != 0 {
		t.Errorf("expected 0 ops for empty input, got %d", len(ops))
	}
}

func TestFromHTMLDefaultCanvasSize(t *testing.T) {
	htmlDoc := `<html><body>
<div class="slides-row">
  <div class="slide">
    <div class="s-c"><div class="h">Test</div></div>
  </div>
</div>
</body></html>`

	ops, err := FromHTML(strings.NewReader(htmlDoc), Options{})
	if err != nil {
		t.Fatalf("FromHTML failed: %v", err)
	}

	if len(ops) != 1 {
		t.Fatalf("expected 1 op, got %d", len(ops))
	}

	params := ops[0]["params"].(map[string]interface{})
	// Defaults to 1080x1350
	if params["canvas"] != "1080x1350" {
		t.Errorf("default canvas = %v, want 1080x1350", params["canvas"])
	}
}

func TestResolveBackgroundGradient(t *testing.T) {
	props := map[string]string{
		"background": "linear-gradient(150deg, #0C1E2C 0%, #14344A 100%)",
	}
	color, grad := resolveBackground(props)
	if color != "#0C1E2C" {
		t.Errorf("base color = %v, want #0C1E2C", color)
	}
	if grad == nil {
		t.Fatal("gradient should not be nil")
	}
	if grad["type"] != "LINEAR" {
		t.Errorf("gradient type = %v", grad["type"])
	}
}

func TestResolveBackgroundSolid(t *testing.T) {
	props := map[string]string{
		"background": "#14344A",
	}
	color, grad := resolveBackground(props)
	if color != "#14344A" {
		t.Errorf("color = %v, want #14344A", color)
	}
	if grad != nil {
		t.Error("solid bg should not have gradient")
	}
}

func TestTextContentWithBreaks(t *testing.T) {
	htmlStr := `<div>Hello<br>World<br>Test</div>`
	doc, _ := html.Parse(strings.NewReader(htmlStr))
	// Find the div
	var div *html.Node
	walkDOM(doc, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "div" {
			div = n
			return false
		}
		return true
	})
	if div == nil {
		t.Fatal("div not found")
	}
	got := textContentWithBreaks(div)
	if got != "Hello\nWorld\nTest" {
		t.Errorf("textContentWithBreaks = %q, want %q", got, "Hello\nWorld\nTest")
	}
}

func TestHasClass(t *testing.T) {
	if !hasClass("slide bg-dark", "slide") {
		t.Error("should find 'slide'")
	}
	if !hasClass("slide bg-dark", "bg-dark") {
		t.Error("should find 'bg-dark'")
	}
	if hasClass("slide bg-dark", "bg") {
		t.Error("should not match partial class 'bg'")
	}
	if hasClass("", "anything") {
		t.Error("empty classes should not match")
	}
}

func TestCleanFontFamily(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"'Poppins', sans-serif", "Poppins"},
		{"\"Inter\", Helvetica, sans-serif", "Inter"},
		{"Arial", "Arial"},
		{"'DM Sans'", "DM Sans"},
	}
	for _, tt := range tests {
		got := cleanFontFamily(tt.input)
		if got != tt.want {
			t.Errorf("cleanFontFamily(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestMergeClassStyles(t *testing.T) {
	styles := map[string]map[string]string{
		".h":    {"font-weight": "700", "color": "#FFF"},
		".h-xl": {"font-size": "1.6rem"},
	}

	merged := mergeClassStyles("h h-xl", styles)
	if merged["font-weight"] != "700" {
		t.Errorf("font-weight = %q, want 700", merged["font-weight"])
	}
	if merged["font-size"] != "1.6rem" {
		t.Errorf("font-size = %q, want 1.6rem", merged["font-size"])
	}
	if merged["color"] != "#FFF" {
		t.Errorf("color = %q, want #FFF", merged["color"])
	}
}

func TestDistributeStopPositions(t *testing.T) {
	stops := []GradientStop{
		{Color: "#000", Position: -1},
		{Color: "#555", Position: -1},
		{Color: "#FFF", Position: -1},
	}
	distributeStopPositions(stops)
	if stops[0].Position != 0 {
		t.Errorf("stop[0] = %v, want 0", stops[0].Position)
	}
	if stops[1].Position != 0.5 {
		t.Errorf("stop[1] = %v, want 0.5", stops[1].Position)
	}
	if stops[2].Position != 1.0 {
		t.Errorf("stop[2] = %v, want 1.0", stops[2].Position)
	}
}

func TestFromHTMLSlidePhotoAndOverlay(t *testing.T) {
	baseDir := "/tmp/amc"
	htmlDoc := `<!DOCTYPE html>
<html>
<head>
<style>
  .slide { width: 280px; height: 350px; }
  .s-ov { background: linear-gradient(180deg, rgba(12,26,40,.2) 0%, rgba(12,26,40,.97) 100%); }
</style>
</head>
<body>
<div class="slides-row">
  <div class="slide">
    <div class="s-img" style="background-image:url('assets/photos/hero.jpg');"></div>
    <div class="s-ov"></div>
    <div class="s-c"><div class="h h-sm">Headline</div></div>
  </div>
</div>
</body>
</html>`

	ops, err := FromHTML(strings.NewReader(htmlDoc), Options{
		CanvasWidth:  1080,
		CanvasHeight: 1350,
		BaseDir:      baseDir,
	})
	if err != nil {
		t.Fatalf("FromHTML failed: %v", err)
	}
	if len(ops) != 1 {
		t.Fatalf("expected 1 op, got %d", len(ops))
	}

	params := ops[0]["params"].(map[string]interface{})
	if params["bgImage"] != "/tmp/amc/assets/photos/hero.jpg" {
		t.Fatalf("bgImage = %v, want normalized absolute path", params["bgImage"])
	}
	if _, ok := params["overlayGradient"]; !ok {
		t.Fatalf("overlayGradient missing")
	}
}

func TestTextContentPreservesWordSpacingAcrossInlineNodes(t *testing.T) {
	htmlStr := `<div>In <strong>hospitals</strong> during care.</div>`
	doc, _ := html.Parse(strings.NewReader(htmlStr))
	var div *html.Node
	walkDOM(doc, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "div" {
			div = n
			return false
		}
		return true
	})
	if div == nil {
		t.Fatal("div not found")
	}
	got := textContentWithBreaks(div)
	want := "In hospitals during care."
	if got != want {
		t.Fatalf("textContentWithBreaks = %q, want %q", got, want)
	}
}

func TestExtractSlideElements_AMCBodyCtaAndStatsClasses(t *testing.T) {
	htmlDoc := `<html><head><style>
  .b-lg { color: rgba(250,252,251,.75); }
  .cta-g { background: #029056; color: #ffffff; }
  .sv { color: #B3D9E8; }
  .sl { color: rgba(250,252,251,.55); }
</style></head><body>
<div class="slide">
  <div class="s-c">
    <div class="b-lg">Body copy</div>
    <div class="cta-g">Donate Now</div>
    <div class="stat-g2">
      <div><div class="sv">250+</div><div class="sl">Muslim Chaplains</div></div>
      <div><div class="sv">5</div><div class="sl">Settings</div></div>
    </div>
  </div>
</div>
</body></html>`

	doc, err := html.Parse(strings.NewReader(htmlDoc))
	if err != nil {
		t.Fatalf("parse html: %v", err)
	}
	styles := parseStyleBlocks(htmlDoc)

	var slide *html.Node
	walkDOM(doc, func(n *html.Node) bool {
		if n.Type == html.ElementNode && hasClass(getAttr(n, "class"), "slide") {
			slide = n
			return false
		}
		return true
	})
	if slide == nil {
		t.Fatal("slide not found")
	}

	elements := extractSlideElements(slide, styles, 1)
	if len(elements) < 3 {
		t.Fatalf("expected at least 3 extracted elements, got %d", len(elements))
	}

	foundBody := false
	foundCTA := false
	foundStats := false
	for _, e := range elements {
		em, _ := e.(map[string]interface{})
		switch em["type"] {
		case "body":
			foundBody = true
		case "cta":
			foundCTA = true
		case "stats":
			foundStats = true
		}
	}
	if !foundBody {
		t.Fatal("expected body element from .b-lg")
	}
	if !foundCTA {
		t.Fatal("expected cta element from .cta-g")
	}
	if !foundStats {
		t.Fatal("expected stats element from .stat-g2")
	}
}
