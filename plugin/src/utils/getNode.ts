/**
 * Shared node lookup utilities.
 *
 * Every lookup calls figma.currentPage.loadAsync() first so that
 * recently-created nodes are always findable by getNodeByIdAsync.
 * Without this, nodes created in one handler call are invisible
 * to subsequent handler calls within the same Figma session.
 */

/** Load the current page and look up any BaseNode by ID. */
export async function getNodeById(nodeId: string): Promise<BaseNode> {
  await figma.currentPage.loadAsync();
  const node = await figma.getNodeByIdAsync(nodeId);
  if (!node) throw new Error(`Node not found: ${nodeId}`);
  return node;
}

/** Look up a SceneNode by ID (throws if missing or not a scene node). */
export async function getSceneNodeById(nodeId: string): Promise<SceneNode> {
  const node = await getNodeById(nodeId);
  return node as SceneNode;
}

/** Look up a node and verify it's a page. */
export async function getPageNodeById(pageId: string): Promise<PageNode> {
  await figma.loadAllPagesAsync();
  const node = await figma.getNodeByIdAsync(pageId);
  if (!node || node.type !== 'PAGE') throw new Error(`Not a page: ${pageId}`);
  return node as PageNode;
}

/** Look up a node and verify it's a text node. */
export async function getTextNodeById(nodeId: string): Promise<TextNode> {
  const node = await getNodeById(nodeId);
  if (node.type !== 'TEXT') throw new Error(`Node ${nodeId} is not a text node`);
  return node as TextNode;
}

/** Look up a node and verify it supports children (frame/component/etc). */
export async function getContainerById(nodeId: string): Promise<BaseNode & ChildrenMixin> {
  const node = await getNodeById(nodeId);
  if (!('appendChild' in node)) throw new Error(`Node ${nodeId} is not a container`);
  return node as BaseNode & ChildrenMixin;
}

/** Look up a parent node for appending children. Returns null if not found or not a container (graceful). */
export async function getParentById(parentId: string): Promise<(BaseNode & ChildrenMixin) | null> {
  await figma.currentPage.loadAsync();
  const node = await figma.getNodeByIdAsync(parentId);
  if (!node || !('appendChild' in node)) return null;
  return node as BaseNode & ChildrenMixin;
}
