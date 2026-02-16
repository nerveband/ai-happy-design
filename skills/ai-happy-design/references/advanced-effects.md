# Advanced Effects & Creative Techniques

## Glass Morphism

### Quick Glass (auto-detects native vs simulated)
```json
{"command":"glass","params":{"nodeId":"$card","intensity":"medium","tint":"#FFFFFF"}}
```
- `light`: subtle frost, 8px blur radius
- `medium`: balanced, 12px blur radius
- `heavy`: dense frost, 16px blur radius

### Native Glass (Figma Beta)
Direct control over Figma's native GLASS effect:
```json
{"command":"effect.add_glass","params":{
  "nodeId":"$card",
  "lightIntensity": 0.5,
  "lightAngle": 45,
  "refraction": 0.5,
  "depth": 1.5,
  "dispersion": 0.1,
  "radius": 12
}}
```

### Full Glass Card Recipe
```json
[
  {"name":"bg","command":"frame","params":{"name":"BG","w":1080,"h":1080,"bg":"#0a0a1a"}},
  {"command":"gradient","params":{"nodeId":"$bg","type":"LINEAR","stops":[{"position":0,"color":"#667eea"},{"position":1,"color":"#764ba2"}]}},
  {"name":"card","command":"frame","params":{"name":"Glass Card","pid":"$bg","x":64,"y":300,"w":952,"h":400,"r":24,"stroke":"#ffffff20","sw":1,"noFill":true,"layoutMode":"VERTICAL","padding":48,"itemSpacing":16}},
  {"command":"glass","params":{"nodeId":"$card","intensity":"medium"}},
  {"command":"shadow","params":{"nodeId":"$card","offsetY":8,"radius":32,"color":"#00000033"}},
  {"command":"text","params":{"text":"Glass Card Title","pid":"$card","sz":64,"fontStyle":"Bold","color":"#ffffff","w":856}},
  {"command":"text","params":{"text":"Frosted glass effect with depth","pid":"$card","sz":36,"color":"#ffffffCC","w":856}}
]
```

## Noise & Texture

### Monotone Noise (Grain Overlay)
```json
{"command":"noise","params":{"nodeId":"$bg","noiseType":"monotone","color":"#FFFFFF","noiseSize":100,"density":0.3}}
```
- `density`: 0.1 (subtle) to 0.5 (heavy grain)
- `noiseSize`: 50 (fine) to 200 (coarse)
- Best colors: #FFFFFF (light grain on dark), #000000 (dark grain on light)

### Duotone Noise
```json
{"command":"noise","params":{"nodeId":"$bg","noiseType":"duotone","color":"#FF6B6B","secondaryColor":"#4ECDC4","density":0.2}}
```

### Texture Effect
```json
{"command":"effect.add_texture","params":{"nodeId":"$bg","noiseSize":100,"radius":0}}
```

## Shadow Recipes

### Production-Grade Presets
| Name | offsetY | radius | spread | color | Use for |
|------|---------|--------|--------|-------|---------|
| Subtle | 2 | 4 | 0 | #00000010 | Flat cards, minimal UI |
| Card | 4 | 12 | -2 | #0000001A | Content cards |
| Elevated | 8 | 24 | -4 | #00000026 | Modals, popovers |
| Floating | 16 | 48 | -8 | #00000033 | FABs, floating elements |
| Glow | 0 | 24 | 4 | accent+40 | CTAs, accent buttons |
| Inner light | -1 | 0 | 0 | #FFFFFF20 | Top bevel on buttons |
| Inner depth | 2 | 4 | 0 | #00000020 | Inset/pressed state |

### Layered Shadows (Realistic)
Stack 3 shadows for realistic depth perception:
```json
[
  {"command":"shadow","params":{"nodeId":"$card","offsetY":1,"radius":2,"color":"#0000000D"}},
  {"command":"shadow","params":{"nodeId":"$card","offsetY":4,"radius":12,"color":"#0000001A"}},
  {"command":"shadow","params":{"nodeId":"$card","offsetY":16,"radius":48,"color":"#00000012"}}
]
```

## Masking

### Circular Avatar
```json
[
  {"name":"photo","command":"image","params":{"pid":"$root","x":40,"y":40,"w":120,"h":120,"imageData":"...","scaleMode":"FILL","name":"avatar-photo"}},
  {"name":"mask_shape","command":"ellipse","params":{"pid":"$root","x":40,"y":40,"w":120,"h":120,"color":"#000000","name":"avatar-mask"}},
  {"command":"mask","params":{"nodeId":"${{steps.mask_shape.result.id}}","targetIds":["${{steps.photo.result.id}}"]}}
]
```

### Rounded Rectangle Crop
```json
[
  {"name":"img","command":"image","params":{"pid":"$root","x":0,"y":0,"w":400,"h":300,"imageData":"...","name":"hero-image"}},
  {"name":"crop","command":"rect","params":{"pid":"$root","x":0,"y":0,"w":400,"h":300,"r":24,"color":"#000000","name":"image-crop"}},
  {"command":"mask","params":{"nodeId":"${{steps.crop.result.id}}","targetIds":["${{steps.img.result.id}}"]}}
]
```

## Gradient Techniques

### Text Readability Overlay
Dark gradient at bottom of image for text:
```json
{"command":"gradient","params":{"nodeId":"$overlay","type":"LINEAR","angle":180,"stops":[
  {"position":0,"color":"#00000000"},
  {"position":0.5,"color":"#00000066"},
  {"position":1,"color":"#000000CC"}
]}}
```

### Gradient Background
```json
{"command":"gradient","params":{"nodeId":"$bg","type":"LINEAR","stops":[
  {"position":0,"color":"#667eea"},
  {"position":1,"color":"#764ba2"}
]}}
```

### Radial Spotlight
```json
{"command":"gradient","params":{"nodeId":"$bg","type":"RADIAL","stops":[
  {"position":0,"color":"#ffffff15"},
  {"position":1,"color":"#00000000"}
]}}
```

## Creative Compositions

### Minimaximalism
One oversized focal element + minimal supporting content:
- Display-tier text (200px+) in accent color as hero
- Body/caption tier for supporting text in muted gray
- Generous padding (80-120px), clean backgrounds
- Single accent color against monochrome

### Editorial Layout
Controlled asymmetry:
- Break grid with one overlapping element
- Decorative shapes at low opacity (0.05-0.10) as ambient texture
- Dramatic font size contrast (hero 152px + caption 36px)
- Rotated elements at 5-15 degrees for energy

### Textured Gradient
Grain + gradient for warmth:
```json
[
  {"command":"gradient","params":{"nodeId":"$bg","type":"LINEAR","stops":[{"position":0,"color":"#1a1a2e"},{"position":1,"color":"#16213e"}]}},
  {"command":"noise","params":{"nodeId":"$bg","noiseType":"monotone","color":"#FFFFFF","density":0.15,"noiseSize":80}}
]
```

### Dark + Neon
High contrast for social media:
- Background: #0a0a0a to #111111
- Accent: electric neon (#00ff88, #7c3aed, #f43f5e, #00d4ff)
- Glow shadow with accent color at 40% opacity
- Subtle noise grain for depth

### Warm Earth Tones
Organic, premium feel:
- Background: #1a1a1a
- Text: cream #f5f0e8
- Accents: terracotta #c4704a, sage #8fbc8f, sand #d4a574
- Noise overlay with warm tint for texture

## Color Palettes

### Social Media (Dark)
```
bg: #0a0a0a | surface: #1a1a1a | border: #ffffff12
text-primary: #ffffff | text-muted: #999999 | accent: varies
```

### Social Media (Light)
```
bg: #ffffff | surface: #f5f5f5 | border: #e5e5e5
text-primary: #1a1a1a | text-muted: #666666 | accent: varies
```

### Popular Accent Pairs
- Electric blue + coral: #3b82f6 + #f43f5e
- Purple + gold: #7c3aed + #f59e0b
- Emerald + pink: #10b981 + #ec4899
- Cyan + orange: #06b6d4 + #f97316
