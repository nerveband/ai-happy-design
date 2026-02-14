// Figma assigns transient session IDs to newly created nodes. These IDs
// expire after the current plugin handler completes. Reading from the
// parent's children list returns the committed/persisted ID that works
// across separate commands (both CLI and MCP).
//
// We use a unique temporary name marker to find our node, with retries
// and increasing delays to ensure Figma has committed the stable ID.
export async function resolveStableId(node: SceneNode, parent: BaseNode & ChildrenMixin): Promise<string> {
  var originalName = node.name;
  var marker = '__ahd_' + Math.random().toString(36).slice(2, 10);
  node.name = marker;

  var transientId = node.id;

  // Retry with increasing delays until we get a committed (different) ID
  for (var attempt = 0; attempt < 4; attempt++) {
    await new Promise(resolve => setTimeout(resolve, attempt * 15));

    for (var i = parent.children.length - 1; i >= 0; i--) {
      if (parent.children[i].name === marker) {
        var candidateId = parent.children[i].id;
        if (candidateId !== transientId) {
          // Committed ID found — different from transient
          node.name = originalName;
          return candidateId;
        }
        break;
      }
    }
  }

  // Fallback: return whatever we found (may still be transient)
  node.name = originalName;
  return transientId;
}
