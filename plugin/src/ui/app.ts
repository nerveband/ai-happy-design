import { WSClient, ConnectionStatus } from '../ws/client';
import { generateChannelKey } from '../utils/channel';
import { DEFAULT_PORT, normalizeRelayUrl } from '../utils/relay';

type ConnectionSettings = {
  channelKey: string;
  relayUrl: string;
  autoConnect: boolean;
};

// ---- State ----
var wsClient: WSClient | null = null;
var channelKey = 'loading...';
var autoConnect = true;
var didAutoConnect = false;
var commandCount = 0;
var errorCount = 0;
var successCount = 0;
var totalCount = 0;
var connectedAt: number | null = null;
var uptimeInterval: ReturnType<typeof setInterval> | null = null;
var activeFilter = 'all';
var searchQuery = '';

// ---- DOM Elements ----
var statusEl = document.getElementById('status')!;
var channelKeyEl = document.getElementById('channel-key')!;
var copyBtn = document.getElementById('copy-btn')!;
var regenerateBtn = document.getElementById('regenerate-btn')!;
var relayUrlInput = document.getElementById('relay-url-input') as HTMLInputElement;
var connectBtn = document.getElementById('connect-btn')!;
var checkRelayBtn = document.getElementById('check-relay-btn')!;
var relayStatusEl = document.getElementById('relay-status')!;
var progressContainer = document.getElementById('progress-container')!;
var progressBar = document.getElementById('progress-bar')!;
var progressText = document.getElementById('progress-text')!;
var logEl = document.getElementById('log')!;
var statUptime = document.getElementById('stat-uptime')!;
var infoBtn = document.getElementById('info-btn')!;
var infoOverlay = document.getElementById('info-overlay')!;
var infoClose = document.getElementById('info-close')!;
var logSearchInput = document.getElementById('log-search') as HTMLInputElement;
var logClearBtn = document.getElementById('log-clear')!;
var logFilterBtns = document.querySelectorAll('.log-filter') as NodeListOf<HTMLElement>;
var countAll = document.getElementById('count-all')!;
var countInfo = document.getElementById('count-info')!;
var countError = document.getElementById('count-error')!;
var countSuccess = document.getElementById('count-success')!;
var logHeader = document.getElementById('log-header')!;
var logSection = document.getElementById('log-section')!;
var logToggleIcon = document.getElementById('log-toggle-icon')!;
var logBadge = document.getElementById('log-badge')!;
var channelInfoBtn = document.getElementById('channel-info-btn')!;

// ---- Initialize ----
channelKeyEl.textContent = channelKey;

// ---- Filter counts ----
function updateFilterCounts() {
  countAll.textContent = String(totalCount);
  countInfo.textContent = String(commandCount);
  countError.textContent = String(errorCount);
  countSuccess.textContent = String(successCount);
  logBadge.textContent = String(totalCount);
}

// ---- Logging ----
function addLog(message: string, level: 'info' | 'warn' | 'error' | 'success' = 'info') {
  var entry = document.createElement('div');
  entry.className = 'log-entry ' + level;
  entry.setAttribute('data-level', level);
  entry.setAttribute('data-text', message.toLowerCase());

  var time = new Date();
  var h = time.getHours().toString();
  var m = time.getMinutes().toString();
  var s = time.getSeconds().toString();
  if (h.length < 2) h = '0' + h;
  if (m.length < 2) m = '0' + m;
  if (s.length < 2) s = '0' + s;
  var ts = h + ':' + m + ':' + s;

  entry.innerHTML = '<span class="log-time">' + ts + '</span>' + escapeHtml(message);

  // Update counts
  totalCount++;
  if (level === 'info') commandCount++;
  if (level === 'error') errorCount++;
  if (level === 'success') successCount++;
  updateFilterCounts();

  // Apply current filter to new entry
  if (!entryMatchesFilter(entry)) {
    entry.classList.add('hidden');
  }

  logEl.appendChild(entry);
  logEl.scrollTop = logEl.scrollHeight;

  // Limit log entries
  while (logEl.children.length > 200) {
    logEl.removeChild(logEl.firstChild!);
  }
}

function entryMatchesFilter(entry: HTMLElement): boolean {
  var level = entry.getAttribute('data-level') || '';
  var text = entry.getAttribute('data-text') || '';
  if (activeFilter !== 'all' && level !== activeFilter) return false;
  if (searchQuery && text.indexOf(searchQuery) === -1) return false;
  return true;
}

function applyLogFilters() {
  var entries = logEl.querySelectorAll('.log-entry');
  for (var i = 0; i < entries.length; i++) {
    var entry = entries[i] as HTMLElement;
    if (entryMatchesFilter(entry)) {
      entry.classList.remove('hidden');
    } else {
      entry.classList.add('hidden');
    }
  }
}

function escapeHtml(str: string): string {
  return str.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
}

function normalizeSettings(raw: any): ConnectionSettings {
  var rawChannel = raw && typeof raw.channelKey === 'string' ? raw.channelKey.trim() : '';
  var sanitizedChannel = rawChannel || generateChannelKey();
  var rawPort = raw && typeof raw.port === 'number' && Number.isFinite(raw.port) && raw.port > 0
    ? Math.floor(raw.port) : DEFAULT_PORT;
  var rawRelayUrl = raw && raw.relayUrl ? raw.relayUrl : ('ws://localhost:' + rawPort + '/ws');
  var relayUrl = normalizeRelayUrl(rawRelayUrl);
  var shouldAutoConnect = !raw || raw.autoConnect !== false;
  return { channelKey: sanitizedChannel, relayUrl: relayUrl, autoConnect: shouldAutoConnect };
}

function requestConnectionSettings() {
  parent.postMessage({ pluginMessage: { type: 'request-connection-settings' } }, '*');
}

function saveConnectionSettings() {
  if (!channelKey || channelKey === 'loading...') return;
  var relayUrl = normalizeRelayUrl(relayUrlInput.value);
  relayUrlInput.value = relayUrl;
  parent.postMessage({
    pluginMessage: {
      type: 'save-connection-settings',
      settings: { channelKey: channelKey, relayUrl: relayUrl, autoConnect: autoConnect },
    },
  }, '*');
}

function applyConnectionSettings(rawSettings: any) {
  var settings = normalizeSettings(rawSettings);
  channelKey = settings.channelKey;
  autoConnect = settings.autoConnect;
  relayUrlInput.value = settings.relayUrl;
  channelKeyEl.textContent = channelKey;

  if (wsClient) {
    wsClient.updateOptions({ relayUrl: settings.relayUrl, channelKey: channelKey });
  }

  addLog('Channel: ' + channelKey, 'info');
  addLog('Relay: ' + settings.relayUrl, 'info');

  if (autoConnect && !didAutoConnect && (!wsClient || !wsClient.isConnected)) {
    didAutoConnect = true;
    doConnect();
  }
}

// ---- Status Updates ----
function updateStatus(status: ConnectionStatus) {
  statusEl.className = 'status ' + status;
  var labels: Record<ConnectionStatus, string> = {
    disconnected: 'Disconnected',
    connecting: 'Connecting...',
    connected: 'Connected',
    reconnecting: 'Reconnecting...',
  };
  statusEl.textContent = labels[status];

  if (status === 'connected') {
    connectBtn.textContent = 'Disconnect';
    connectBtn.className = 'disconnect';
    connectedAt = Date.now();
    startUptimeTimer();
  } else if (status === 'disconnected') {
    connectBtn.textContent = 'Connect';
    connectBtn.className = '';
    connectedAt = null;
    stopUptimeTimer();
  } else if (status === 'reconnecting') {
    connectBtn.textContent = 'Connect';
    connectBtn.className = '';
    connectedAt = null;
    stopUptimeTimer();
  }
}

// ---- Uptime Timer ----
function startUptimeTimer() {
  stopUptimeTimer();
  uptimeInterval = setInterval(function() {
    if (connectedAt) {
      var seconds = Math.floor((Date.now() - connectedAt) / 1000);
      if (seconds < 60) {
        statUptime.textContent = seconds + 's';
      } else if (seconds < 3600) {
        statUptime.textContent = Math.floor(seconds / 60) + 'm';
      } else {
        statUptime.textContent = Math.floor(seconds / 3600) + 'h';
      }
    }
  }, 1000);
}

function stopUptimeTimer() {
  if (uptimeInterval) {
    clearInterval(uptimeInterval);
    uptimeInterval = null;
  }
  statUptime.textContent = '';
}

// ---- Progress Bar ----
function updateProgress(percent: number) {
  if (percent <= 0 || percent >= 100) {
    progressContainer.style.display = 'none';
  } else {
    progressContainer.style.display = 'flex';
    progressBar.style.width = percent + '%';
    progressText.textContent = Math.round(percent) + '%';
  }
}

// ---- Relay Status Check ----
function getHttpUrlFromRelay(wsUrl: string): string {
  var httpUrl = wsUrl.replace(/^ws:/, 'http:').replace(/^wss:/, 'https:');
  var idx = httpUrl.lastIndexOf('/ws');
  if (idx !== -1) {
    httpUrl = httpUrl.substring(0, idx) + '/status';
  } else {
    httpUrl = httpUrl + '/status';
  }
  return httpUrl;
}

function doCheckRelay() {
  var wsUrl = normalizeRelayUrl(relayUrlInput.value);
  var statusUrl = getHttpUrlFromRelay(wsUrl);
  relayStatusEl.className = 'relay-status relay-checking';
  relayStatusEl.textContent = 'Checking...';

  var xhr = new XMLHttpRequest();
  xhr.open('GET', statusUrl, true);
  xhr.timeout = 5000;

  xhr.onload = function() {
    if (xhr.status === 200) {
      try {
        var data = JSON.parse(xhr.responseText);
        var channels = data.channels || 0;
        var clients = data.clients || 0;
        relayStatusEl.className = 'relay-status relay-ok';
        relayStatusEl.textContent = 'Relay running — ' + channels + ' ch, ' + clients + ' client(s)';
        addLog('Relay OK (' + channels + ' ch, ' + clients + ' clients)', 'success');
      } catch (e) {
        relayStatusEl.className = 'relay-status relay-ok';
        relayStatusEl.textContent = 'Relay running';
        addLog('Relay OK (status 200)', 'success');
      }
    } else {
      relayStatusEl.className = 'relay-status relay-err';
      relayStatusEl.textContent = 'Relay HTTP ' + xhr.status;
      addLog('Relay check: HTTP ' + xhr.status, 'warn');
    }
  };

  xhr.onerror = function() {
    relayStatusEl.className = 'relay-status relay-err';
    relayStatusEl.textContent = 'Relay not running — run: ai-happy-design relay start';
    addLog('Relay not reachable at ' + statusUrl, 'error');
  };

  xhr.ontimeout = function() {
    relayStatusEl.className = 'relay-status relay-err';
    relayStatusEl.textContent = 'Relay check timed out';
    addLog('Relay check timed out', 'error');
  };

  xhr.send();
}

// ---- WebSocket Message Handling ----
function handleWSMessage(msg: any) {
  switch (msg.type) {
    case 'command': {
      var label = msg.domain && msg.action
        ? msg.domain + '.' + msg.action
        : (msg.command || '<unknown>');
      addLog('CMD ' + label, 'info');

      parent.postMessage({
        pluginMessage: {
          type: 'execute-command',
          id: msg.id,
          domain: msg.domain,
          action: msg.action,
          params: msg.params || {},
        },
      }, '*');
      break;
    }

    case 'progress': {
      var pct = msg.percent != null ? msg.percent : 0;
      updateProgress(pct);
      break;
    }

    case 'joined': {
      addLog('Joined: ' + msg.channel, 'success');
      break;
    }

    case 'pong':
    case 'message': {
      break;
    }

    case 'error': {
      addLog('Server: ' + msg.message, 'error');
      break;
    }

    default: {
      addLog('Unknown type: ' + msg.type, 'warn');
    }
  }
}

// ---- Receive results from plugin sandbox ----
window.onmessage = function(event: MessageEvent) {
  var msg = event.data.pluginMessage;
  if (!msg) return;

  if (msg.type === 'command-result') {
    addLog('OK [' + msg.id + ']', 'success');
    if (wsClient) {
      wsClient.send({
        type: 'message',
        channel: channelKey,
        message: { id: msg.id, result: msg.result },
      });
    }
  } else if (msg.type === 'command-error') {
    addLog('ERR [' + msg.id + ']: ' + msg.error, 'error');
    if (wsClient) {
      wsClient.send({
        type: 'message',
        channel: channelKey,
        message: { id: msg.id, error: msg.error },
      });
    }
  } else if (msg.type === 'connection-settings') {
    applyConnectionSettings(msg.settings);
  } else if (msg.type === 'connection-settings-error') {
    addLog('Settings error: ' + msg.error, 'error');
  } else if (msg.type === 'connection-settings-saved') {
    var settings = normalizeSettings(msg.settings);
    channelKey = settings.channelKey;
    relayUrlInput.value = settings.relayUrl;
    channelKeyEl.textContent = channelKey;
  }
};

// ---- Connect / Disconnect ----
function doConnect() {
  if (!channelKey || channelKey === 'loading...') {
    addLog('Waiting for settings...', 'warn');
    requestConnectionSettings();
    return;
  }

  var relayUrl = normalizeRelayUrl(relayUrlInput.value);
  relayUrlInput.value = relayUrl;

  if (wsClient && wsClient.isConnected) {
    wsClient.disconnect();
    return;
  }

  wsClient = new WSClient({
    relayUrl: relayUrl,
    channelKey: channelKey,
    onMessage: handleWSMessage,
    onStatusChange: updateStatus,
    onLog: addLog,
  });

  wsClient.connect();
  saveConnectionSettings();
}

// ---- Copy Channel Key ----
function doCopy() {
  var textarea = document.createElement('textarea');
  textarea.value = channelKey;
  textarea.style.position = 'fixed';
  textarea.style.opacity = '0';
  document.body.appendChild(textarea);
  textarea.select();
  try {
    document.execCommand('copy');
    copyBtn.textContent = 'Copied!';
    copyBtn.classList.add('copied');
    addLog('Channel key copied', 'success');
    setTimeout(function() {
      copyBtn.textContent = 'Copy';
      copyBtn.classList.remove('copied');
    }, 2000);
  } catch (e) {
    addLog('Copy failed', 'warn');
  }
  document.body.removeChild(textarea);
}

function doRegenerateChannel() {
  if (wsClient && wsClient.isConnected) wsClient.disconnect();
  channelKey = generateChannelKey();
  channelKeyEl.textContent = channelKey;
  addLog('New channel: ' + channelKey, 'warn');
  saveConnectionSettings();
  if (autoConnect) doConnect();
}

// ---- Info Modal ----
function openInfoModal() {
  infoOverlay.classList.add('open');
}

function closeInfoModal() {
  infoOverlay.classList.remove('open');
}

// ---- Event Listeners ----
connectBtn.addEventListener('click', doConnect);
copyBtn.addEventListener('click', doCopy);
regenerateBtn.addEventListener('click', doRegenerateChannel);
checkRelayBtn.addEventListener('click', doCheckRelay);
infoBtn.addEventListener('click', openInfoModal);
infoClose.addEventListener('click', closeInfoModal);
infoOverlay.addEventListener('click', function(e: MouseEvent) {
  if (e.target === infoOverlay) closeInfoModal();
});

relayUrlInput.addEventListener('keydown', function(e: KeyboardEvent) {
  if (e.key === 'Enter') doConnect();
});
relayUrlInput.addEventListener('change', saveConnectionSettings);

// ---- Log Filter Listeners ----
for (var i = 0; i < logFilterBtns.length; i++) {
  logFilterBtns[i].addEventListener('click', function(this: HTMLElement) {
    activeFilter = this.getAttribute('data-filter') || 'all';
    for (var j = 0; j < logFilterBtns.length; j++) {
      logFilterBtns[j].classList.remove('active');
    }
    this.classList.add('active');
    applyLogFilters();
  });
}

logSearchInput.addEventListener('input', function() {
  searchQuery = logSearchInput.value.toLowerCase().trim();
  applyLogFilters();
});

logClearBtn.addEventListener('click', function() {
  logEl.innerHTML = '';
  commandCount = 0;
  errorCount = 0;
  successCount = 0;
  totalCount = 0;
  updateFilterCounts();
  addLog('Log cleared', 'info');
});

var logCopyBtn = document.getElementById('log-copy-btn')!;
var logExportBtn = document.getElementById('log-export-btn')!;

// ---- Log Copy/Export ----
function getLogText(): string {
  var entries = logEl.querySelectorAll('.log-entry');
  var lines: string[] = [];
  for (var i = 0; i < entries.length; i++) {
    lines.push((entries[i] as HTMLElement).textContent || '');
  }
  return lines.join('\n');
}

logCopyBtn.addEventListener('click', function() {
  var text = getLogText();
  var textarea = document.createElement('textarea');
  textarea.value = text;
  textarea.style.position = 'fixed';
  textarea.style.opacity = '0';
  document.body.appendChild(textarea);
  textarea.select();
  try {
    document.execCommand('copy');
    logCopyBtn.textContent = 'Copied!';
    setTimeout(function() { logCopyBtn.textContent = 'Copy'; }, 2000);
  } catch (e) {
    addLog('Copy failed', 'error');
  }
  document.body.removeChild(textarea);
});

logExportBtn.addEventListener('click', function() {
  var text = getLogText();
  var blob = new Blob([text], { type: 'text/plain' });
  var url = URL.createObjectURL(blob);
  var a = document.createElement('a');
  a.href = url;
  a.download = 'ahd-log-' + new Date().toISOString().slice(0, 19).replace(/[T:]/g, '-') + '.txt';
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
  addLog('Log exported', 'success');
});

// ---- Log Collapse Toggle ----
logHeader.addEventListener('click', function() {
  var isCollapsed = logSection.classList.contains('collapsed');
  if (isCollapsed) {
    logSection.classList.remove('collapsed');
    logToggleIcon.innerHTML = '&#9660;';
  } else {
    logSection.classList.add('collapsed');
    logToggleIcon.innerHTML = '&#9654;';
  }
});

// ---- Channel Info Button ----
channelInfoBtn.addEventListener('click', function() {
  openInfoModal();
});

// ---- Initial Log ----
addLog('AI Happy Design ready', 'info');
requestConnectionSettings();
