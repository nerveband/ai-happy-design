/**
 * Font loading helpers.
 * Wraps figma.loadFontAsync with error handling and caching.
 */

const loadedFonts = new Set<string>();

function fontKey(family: string, style: string): string {
  return `${family}::${style}`;
}

/**
 * Load a font, caching the result to avoid redundant loads.
 */
export async function loadFont(family: string, style: string = 'Regular'): Promise<void> {
  const key = fontKey(family, style);
  if (loadedFonts.has(key)) return;

  try {
    await figma.loadFontAsync({ family, style });
    loadedFonts.add(key);
  } catch (err: any) {
    throw new Error(`Failed to load font "${family} ${style}": ${err.message}`);
  }
}

/**
 * Load all fonts used in a text node, handling mixed fonts.
 */
export async function loadNodeFonts(node: TextNode): Promise<void> {
  const fontName = node.fontName;

  if (fontName === figma.mixed) {
    // For mixed fonts, load each segment's font
    const len = node.characters.length;
    const seen = new Set<string>();

    for (let i = 0; i < len; i++) {
      const fn = node.getRangeFontName(i, i + 1) as FontName;
      const key = fontKey(fn.family, fn.style);
      if (!seen.has(key)) {
        seen.add(key);
        await loadFont(fn.family, fn.style);
      }
    }
  } else {
    await loadFont(fontName.family, fontName.style);
  }
}

/**
 * Common font families with fallbacks.
 */
export const COMMON_FONTS: Record<string, string> = {
  'inter': 'Inter',
  'roboto': 'Roboto',
  'arial': 'Arial',
  'helvetica': 'Helvetica',
  'times': 'Times New Roman',
  'courier': 'Courier New',
  'georgia': 'Georgia',
  'verdana': 'Verdana',
  'sf pro': 'SF Pro Display',
  'sf mono': 'SF Mono',
};

/**
 * Resolve a font family name, trying common aliases.
 */
export function resolveFontFamily(name: string): string {
  const lower = name.toLowerCase();
  return COMMON_FONTS[lower] || name;
}
