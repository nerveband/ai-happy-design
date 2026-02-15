import { getSceneNodeById } from '../utils/getNode';

export async function handleExport(action: string, params: any): Promise<any> {
  switch (action) {
    case 'export_image':
    case 'export_png':
    case 'export_node':
    case 'export':
    case 'image': return exportImage(params);
    case 'export_svg':
    case 'svg': return exportSvg(params);
    case 'export_pdf':
    case 'pdf': return exportPdf(params);
    case 'export_json':
    case 'json': return exportJson(params);
    case 'batch':
    case 'batch_export':
    case 'export_batch':
    case 'export_all': return batchExport(params);
    default: throw new Error('Unknown export action: ' + action + '. Available: image, svg, pdf, json, batch_export');
  }
}

async function exportImage(params: any) {
  const { nodeId, format = 'PNG', scale = 2, constraint } = params;
  const node = await getSceneNodeById(nodeId);

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
  const node = await getSceneNodeById(nodeId);

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
  const node = await getSceneNodeById(nodeId);

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
  const node = await getSceneNodeById(nodeId);

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

async function batchExport(params: any) {
  const format = (params.format ?? 'PNG').toUpperCase();
  const scale = params.scale ?? 2;

  // Parse nodeIds: comma-separated string or array
  let nodeIds: string[] = [];
  if (params.nodeIds) {
    if (Array.isArray(params.nodeIds)) {
      nodeIds = params.nodeIds;
    } else {
      nodeIds = String(params.nodeIds).split(',').map((s: string) => s.trim()).filter(Boolean);
    }
  }

  // If no nodeIds, export all top-level frames on current page
  if (nodeIds.length === 0) {
    const page = figma.currentPage;
    for (const child of page.children) {
      if (child.type === 'FRAME' || child.type === 'COMPONENT' || child.type === 'COMPONENT_SET') {
        nodeIds.push(child.id);
      }
    }
  }

  const exports: any[] = [];
  for (const id of nodeIds) {
    const node = await getSceneNodeById(id);

    if (format === 'SVG') {
      const svg = await node.exportAsync({ format: 'SVG_STRING' } as ExportSettingsSVGString);
      exports.push({
        id: node.id,
        name: node.name,
        format: 'SVG',
        scale,
        size: svg.length,
        data: svg,
      });
    } else {
      const settings: ExportSettingsImage = {
        format: format as 'PNG' | 'JPG',
        constraint: { type: 'SCALE', value: scale },
      };
      const bytes = await node.exportAsync(settings);
      exports.push({
        id: node.id,
        name: node.name,
        format,
        scale,
        size: bytes.length,
        data: figma.base64Encode(bytes),
      });
    }
  }

  return { exports, count: exports.length };
}
