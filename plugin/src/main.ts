import { handlePaint } from './handlers/paint';
import { handleShape } from './handlers/shape';
import { handleText } from './handlers/text';
import { handleLayout } from './handlers/layout';
import { handleNode } from './handlers/node';
import { handleLayer } from './handlers/layer';
import { handleComponent } from './handlers/component';
import { handleVariable } from './handlers/variable';
import { handleEffect } from './handlers/effect';
import { handleBoolean } from './handlers/boolean';
import { handlePage } from './handlers/page';
import { handleDocument } from './handlers/document';
import { handleExport } from './handlers/export';

// Domain handler registry
const handlers: Record<string, (action: string, params: any) => Promise<any>> = {
  paint: handlePaint,
  shape: handleShape,
  text: handleText,
  layout: handleLayout,
  node: handleNode,
  layer: handleLayer,
  component: handleComponent,
  variable: handleVariable,
  effect: handleEffect,
  boolean: handleBoolean,
  page: handlePage,
  document: handleDocument,
  export: handleExport,
};

// Show plugin UI
figma.showUI(__html__, { width: 380, height: 500 });

// Handle messages from UI
figma.ui.onmessage = async (msg: any) => {
  if (msg.type === 'execute-command') {
    try {
      const { domain, action, params, id } = msg;
      const handler = handlers[domain];
      if (!handler) throw new Error(`Unknown domain: ${domain}`);

      const result = await handler(action, params);
      figma.ui.postMessage({ type: 'command-result', id, result });
    } catch (error: any) {
      figma.ui.postMessage({
        type: 'command-error',
        id: msg.id,
        error: error.message || String(error),
      });
    }
  } else if (msg.type === 'update-settings') {
    // Handle settings updates from UI
    if (msg.width && msg.height) {
      figma.ui.resize(msg.width, msg.height);
    }
  } else if (msg.type === 'notify') {
    // Show a notification in Figma
    figma.notify(msg.message, {
      timeout: msg.timeout ?? 3000,
      error: msg.error ?? false,
    });
  }
};

// Notify that plugin is ready
figma.notify('AI Happy Design connected', { timeout: 2000 });
