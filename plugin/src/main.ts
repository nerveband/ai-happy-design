import { handlePaint } from './handlers/paint';
import { handleShape } from './handlers/shape';
import { handleText } from './handlers/text';
import { handleLayout } from './handlers/layout';
import { handleNode } from './handlers/node';
import { handleLayer } from './handlers/layer';
import { handleComponent } from './handlers/component';
import { handleStyle } from './handlers/style';
import { handleVariable } from './handlers/variable';
import { handleEffect } from './handlers/effect';
import { handleBoolean } from './handlers/boolean';
import { handlePage } from './handlers/page';
import { handleDocument } from './handlers/document';
import { handleExport } from './handlers/export';
import { handleDraw } from './handlers/draw';
import { handleMotion } from './handlers/motion';
import { generateChannelKey } from './utils/channel';
import { DEFAULT_PORT, normalizeRelayUrl } from './utils/relay';

type ConnectionSettings = {
  channelKey: string;
  relayUrl: string;
  autoConnect: boolean;
};

const SETTINGS_KEY = 'connectionSettings.v1';

function sanitizeSettings(raw: any): ConnectionSettings {
  const channelKey = typeof raw?.channelKey === 'string' && raw.channelKey.trim()
    ? raw.channelKey.trim()
    : generateChannelKey();
  const legacyPort = typeof raw?.port === 'number' && Number.isFinite(raw.port) && raw.port > 0
    ? Math.floor(raw.port)
    : DEFAULT_PORT;
  const relayUrl = normalizeRelayUrl(raw?.relayUrl || `ws://localhost:${legacyPort}/ws`);
  const autoConnect = raw?.autoConnect !== false;
  return { channelKey, relayUrl, autoConnect };
}

async function loadSettings(): Promise<ConnectionSettings> {
  var raw = await figma.clientStorage.getAsync(SETTINGS_KEY);
  var settings = sanitizeSettings(raw);
  // Only persist defaults on first run (no saved settings yet)
  if (raw == null) {
    await figma.clientStorage.setAsync(SETTINGS_KEY, settings);
  }
  return settings;
}

async function saveSettings(partial: any): Promise<ConnectionSettings> {
  const current = await loadSettings();
  const combined: any = {};
  if (current && typeof current === 'object') {
    for (const key of Object.keys(current)) {
      combined[key] = (current as any)[key];
    }
  }
  if (partial && typeof partial === 'object') {
    for (const key of Object.keys(partial)) {
      combined[key] = (partial as any)[key];
    }
  }
  const merged = sanitizeSettings(combined);
  await figma.clientStorage.setAsync(SETTINGS_KEY, merged);
  return merged;
}

// Domain handler registry
const handlers: Record<string, (action: string, params: any) => Promise<any>> = {
  paint: handlePaint,
  shape: handleShape,
  text: handleText,
  layout: handleLayout,
  node: handleNode,
  layer: handleLayer,
  component: handleComponent,
  style: handleStyle,
  variable: handleVariable,
  effect: handleEffect,
  boolean: handleBoolean,
  page: handlePage,
  document: handleDocument,
  export: handleExport,
  draw: handleDraw,
  motion: handleMotion,
  fill: handlePaint,
};

function isReadOnlyCommand(domain: string, action: string): boolean {
  if (domain === 'document' || domain === 'export') {
    return action !== 'set_selection' && action !== 'select' && action !== 'set_selected';
  }
  if (domain === 'motion') {
    return action.indexOf('get_') === 0 || action.indexOf('list_') === 0;
  }
  return action.indexOf('get') === 0 || action.indexOf('list') === 0 || action.indexOf('find') === 0 || action.indexOf('scan') === 0 || action.indexOf('check') === 0 || action.indexOf('lint') === 0;
}

function commitUndoIfNeeded(domain: string, action: string, params: any): void {
  if (params && params.commitUndo === false) return;
  if (isReadOnlyCommand(domain, action)) return;
  var figmaAny = figma as any;
  if (typeof figmaAny.commitUndo === 'function') {
    figmaAny.commitUndo();
  }
}

// Show plugin UI
figma.showUI(__html__, { width: 380, height: 420, themeColors: false });

// Handle messages from UI
figma.ui.onmessage = async (msg: any) => {
  if (msg.type === 'execute-command') {
    try {
      const { domain, action, params, id } = msg;
      const handler = handlers[domain];
      if (!handler) throw new Error(`Unknown domain: ${domain}`);

      const result = await handler(action, params);
      commitUndoIfNeeded(domain, action, params);
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
  } else if (msg.type === 'request-connection-settings') {
    try {
      const settings = await loadSettings();
      figma.ui.postMessage({ type: 'connection-settings', settings });
    } catch (error: any) {
      figma.ui.postMessage({
        type: 'connection-settings-error',
        error: error?.message || String(error),
      });
    }
  } else if (msg.type === 'save-connection-settings') {
    try {
      const settings = await saveSettings(msg.settings);
      figma.ui.postMessage({ type: 'connection-settings-saved', settings });
    } catch (error: any) {
      figma.ui.postMessage({
        type: 'connection-settings-error',
        error: error?.message || String(error),
      });
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
