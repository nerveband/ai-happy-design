export async function handlePage(action: string, params: any): Promise<any> {
  switch (action) {
    case 'create': return createPage(params);
    case 'delete': return deletePage(params);
    case 'rename': return renamePage(params);
    case 'set_current': return setCurrentPage(params);
    case 'get_all': return getAllPages(params);
    case 'get_current': return getCurrentPage(params);
    case 'duplicate': return duplicatePage(params);
    default: throw new Error(`Unknown page action: ${action}`);
  }
}

async function createPage(params: any) {
  const { name } = params;
  const page = figma.createPage();
  if (name) page.name = name;
  return { id: page.id, name: page.name };
}

async function deletePage(params: any) {
  const { pageId } = params;
  const page = figma.getNodeById(pageId);
  if (!page || page.type !== 'PAGE') throw new Error(`Not a page: ${pageId}`);

  // Don't delete if it's the only page
  if (figma.root.children.length <= 1) throw new Error('Cannot delete the only page');

  const info = { id: page.id, name: page.name };
  page.remove();
  return { deleted: info };
}

async function renamePage(params: any) {
  const { pageId, name } = params;
  const page = figma.getNodeById(pageId);
  if (!page || page.type !== 'PAGE') throw new Error(`Not a page: ${pageId}`);

  page.name = name;
  return { id: page.id, name: page.name };
}

async function setCurrentPage(params: any) {
  const { pageId } = params;
  const page = figma.getNodeById(pageId);
  if (!page || page.type !== 'PAGE') throw new Error(`Not a page: ${pageId}`);

  figma.currentPage = page as PageNode;
  return { id: page.id, name: page.name };
}

async function getAllPages(_params: any) {
  const pages = figma.root.children.map(page => ({
    id: page.id,
    name: page.name,
    childCount: page.children.length,
    isCurrent: page.id === figma.currentPage.id,
  }));
  return { pages, count: pages.length };
}

async function getCurrentPage(_params: any) {
  const page = figma.currentPage;
  return {
    id: page.id,
    name: page.name,
    childCount: page.children.length,
  };
}

async function duplicatePage(params: any) {
  const { pageId, name } = params;
  const page = figma.getNodeById(pageId);
  if (!page || page.type !== 'PAGE') throw new Error(`Not a page: ${pageId}`);

  const dup = (page as PageNode).clone();
  if (name) dup.name = name;
  return { id: dup.id, name: dup.name };
}
