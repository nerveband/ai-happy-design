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
import { generateChannelKey } from './utils/channel';

type ConnectionSettings = {
  channelKey: string;
  relayUrl: string;
  autoConnect: boolean;
};

const SETTINGS_KEY = 'connectionSettings.v1';
const DEFAULT_PORT = 3055;
const DEFAULT_RELAY_URL = `ws://localhost:${DEFAULT_PORT}/ws`;

function normalizeRelayUrl(raw: any): string {
  var input = typeof raw === 'string' ? raw.trim() : '';
  if (!input) return DEFAULT_RELAY_URL;

  // Bare port number: "3056" or ":3056"
  var barePortMatch = input.match(/^:?(\d+)$/);
  if (barePortMatch) {
    var p = parseInt(barePortMatch[1], 10);
    if (p < 1 || p > 65535) return DEFAULT_RELAY_URL;
    return 'ws://localhost:' + p + '/ws';
  }

  var candidate = input;
  if (!/^wss?:\/\//i.test(candidate)) {
    candidate = 'ws://' + candidate;
  }

  try {
    var url = new URL(candidate);
    if (url.protocol !== 'ws:' && url.protocol !== 'wss:') {
      return DEFAULT_RELAY_URL;
    }
    if (!url.hostname) {
      return DEFAULT_RELAY_URL;
    }
    if (url.port) {
      var portNum = parseInt(url.port, 10);
      if (portNum < 1 || portNum > 65535) return DEFAULT_RELAY_URL;
    }
    if (!url.pathname || url.pathname === '/') {
      url.pathname = '/ws';
    }
    return url.toString();
  } catch (e) {
    return DEFAULT_RELAY_URL;
  }
}

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
  const merged = sanitizeSettings({ ...current, ...(partial || {}) });
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
};

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
