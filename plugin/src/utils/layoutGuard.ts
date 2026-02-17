// Check if a parent node has auto-layout enabled.
// layoutPositioning: "ABSOLUTE" is only valid inside auto-layout frames.
// Setting it on a child of a non-auto-layout frame causes a Figma error.

export function parentHasAutoLayout(parent: any): boolean {
  if (!parent) return false;
  if (!('layoutMode' in parent)) return false;
  return parent.layoutMode !== 'NONE';
}
