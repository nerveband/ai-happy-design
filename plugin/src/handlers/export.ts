export async function handleExport(action: string, params: any): Promise<any> {
  switch (action) {
    case 'image': return exportImage(params);
    case 'svg': return exportSvg(params);
    case 'pdf': return exportPdf(params);
    case 'json': return exportJson(params);
    default: throw new Error(`Unknown export action: ${action}`);
  }
}

async function exportImage(params: any) {
  const { nodeId, format = 'PNG', scale = 1, constraint } = params;
  const node = figma.getNodeById(nodeId) as SceneNode;
  if (!node) throw new Error(`Node not found: ${nodeId}`);

  const settings: ExportSettings = {
    format: format as 'PNG' | 'JPG',
  };

  if (constraint) {
    settings.constraint = constraint;
  } else {
    (settings as ExportSettingsImage).constraint = { type: 'SCALE', value: scale };
  }

  const bytes = await node.exportAsync(settings as ExportSettingsImage);
  // Convert to base64 for transport
  const base64 = figma.base64Encode(bytes);

  return {
    id: node.id,
    name: node.name,
    format,
    scale,
    size: bytes.length,
    data: base64,
  };
}

async function exportSvg(params: any) {
  const { nodeId, svgIdAttribute, svgOutlineText, svgSimplifyStroke } = params;
  const node = figma.getNodeById(nodeId) as SceneNode;
  if (!node) throw new Error(`Node not found: ${nodeId}`);

  const settings: ExportSettingsSVGString = {
    format: 'SVG_STRING',
  };

  if (svgIdAttribute !== undefined) settings.svgIdAttribute = svgIdAttribute;
  if (svgOutlineText !== undefined) settings.svgOutlineText = svgOutlineText;
  if (svgSimplifyStroke !== undefined) settings.svgSimplifyStroke = svgSimplifyStroke;

  const svg = await node.exportAsync(settings);

  return {
    id: node.id,
    name: node.name,
    format: 'SVG',
    data: svg,
    size: svg.length,
  };
}

async function exportPdf(params: any) {
  const { nodeId } = params;
  const node = figma.getNodeById(nodeId) as SceneNode;
  if (!node) throw new Error(`Node not found: ${nodeId}`);

  const settings: ExportSettingsPDF = {
    format: 'PDF',
  };

  const bytes = await node.exportAsync(settings);
  const base64 = figma.base64Encode(bytes);

  return {
    id: node.id,
    name: node.name,
    format: 'PDF',
    size: bytes.length,
    data: base64,
  };
}

async function exportJson(params: any) {
  const { nodeId, depth } = params;
  const node = figma.getNodeById(nodeId) as SceneNode;
  if (!node) throw new Error(`Node not found: ${nodeId}`);

  // Use the serialize utility for structured export
  const { serializeNode } = await import('../utils/serialize');
  const data = serializeNode(node, 0, depth ?? 10);

  return {
    id: node.id,
    name: node.name,
    format: 'JSON',
    data,
  };
}
