import { WSClient, ConnectionStatus } from '../ws/client';
import { generateChannelKey } from '../utils/channel';

// ---- State ----
let wsClient: WSClient | null = null;
let channelKey = generateChannelKey();
let commandCount = 0;
let errorCount = 0;
let connectedAt: number | null = null;
let uptimeInterval: ReturnType<typeof setInterval> | null = null;

// ---- DOM Elements ----
const statusEl = document.getElementById('status')!;
const channelKeyEl = document.getElementById('channel-key')!;
const copyBtn = document.getElementById('copy-btn')!;
const portInput = document.getElementById('port-input') as HTMLInputElement;
const connectBtn = document.getElementById('connect-btn')!;
const progressContainer = document.getElementById('progress-container')!;
const progressBar = document.getElementById('progress-bar')!;
const progressText = document.getElementById('progress-text')!;
const logEl = document.getElementById('log')!;
const statCommands = document.getElementById('stat-commands')!;
const statErrors = document.getElementById('stat-errors')!;
const statUptime = document.getElementById('stat-uptime')!;

// ---- Initialize ----
channelKeyEl.textContent = channelKey;

// ---- Logging ----
function addLog(message: string, level: 'info' | 'warn' | 'error' | 'success' = 'info') {
  const entry = document.createElement('div');
  entry.className = `log-entry ${level}`;

  const time = new Date();
  const ts = `${time.getHours().toString().padStart(2, '0')}:${time.getMinutes().toString().padStart(2, '0')}:${time.getSeconds().toString().padStart(2, '0')}`;

  entry.innerHTML = `<span class="log-time">${ts}</span>${escapeHtml(message)}`;
  logEl.appendChild(entry);
  logEl.scrollTop = logEl.scrollHeight;

  // Limit log entries
  while (logEl.children.length > 200) {
    logEl.removeChild(logEl.firstChild!);
  }
}

function escapeHtml(str: string): string {
  return str.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
}

// ---- Status Updates ----
function updateStatus(status: ConnectionStatus) {
  statusEl.className = `status ${status}`;
  const labels: Record<ConnectionStatus, string> = {
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
  }
}

// ---- Uptime Timer ----
function startUptimeTimer() {
  stopUptimeTimer();
  uptimeInterval = setInterval(() => {
    if (connectedAt) {
      const seconds = Math.floor((Date.now() - connectedAt) / 1000);
      if (seconds < 60) {
        statUptime.textContent = `${seconds}s`;
      } else if (seconds < 3600) {
        statUptime.textContent = `${Math.floor(seconds / 60)}m`;
      } else {
        statUptime.textContent = `${Math.floor(seconds / 3600)}h`;
      }
    }
  }, 1000);
}

function stopUptimeTimer() {
  if (uptimeInterval) {
    clearInterval(uptimeInterval);
    uptimeInterval = null;
  }
  statUptime.textContent = '0s';
}

// ---- Progress Bar ----
function updateProgress(percent: number) {
  if (percent <= 0 || percent >= 100) {
    progressContainer.style.display = 'none';
  } else {
    progressContainer.style.display = 'flex';
    progressBar.style.width = `${percent}%`;
    progressText.textContent = `${Math.round(percent)}%`;
  }
}

// ---- WebSocket Message Handling ----
function handleWSMessage(msg: any) {
  switch (msg.type) {
    case 'command': {
      // Forward command to plugin sandbox for execution
      commandCount++;
      statCommands.textContent = String(commandCount);
      addLog(`Command: ${msg.domain}.${msg.action}`, 'info');

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
      updateProgress(msg.percent ?? 0);
      break;
    }

    case 'joined': {
      addLog(`Channel joined: ${msg.channel}`, 'success');
      break;
    }

    case 'pong': {
      // Heartbeat response - ignore silently
      break;
    }

    case 'error': {
      errorCount++;
      statErrors.textContent = String(errorCount);
      addLog(`Server error: ${msg.message}`, 'error');
      break;
    }

    default: {
      addLog(`Unknown message type: ${msg.type}`, 'warn');
    }
  }
}

// ---- Receive results from plugin sandbox ----
window.onmessage = (event: MessageEvent) => {
  const msg = event.data.pluginMessage;
  if (!msg) return;

  if (msg.type === 'command-result') {
    addLog(`Result [${msg.id}]: OK`, 'success');
    wsClient?.send({
      type: 'result',
      id: msg.id,
      result: msg.result,
    });
  } else if (msg.type === 'command-error') {
    errorCount++;
    statErrors.textContent = String(errorCount);
    addLog(`Error [${msg.id}]: ${msg.error}`, 'error');
    wsClient?.send({
      type: 'error',
      id: msg.id,
      error: msg.error,
    });
  }
};

// ---- Connect / Disconnect ----
function doConnect() {
  const port = parseInt(portInput.value, 10) || 3055;

  if (wsClient?.isConnected) {
    wsClient.disconnect();
    return;
  }

  wsClient = new WSClient({
    port,
    channelKey,
    onMessage: handleWSMessage,
    onStatusChange: updateStatus,
    onLog: addLog,
  });

  wsClient.connect();
}

// ---- Copy Channel Key ----
function doCopy() {
  // Use a textarea trick since we're in an iframe
  const textarea = document.createElement('textarea');
  textarea.value = channelKey;
  textarea.style.position = 'fixed';
  textarea.style.opacity = '0';
  document.body.appendChild(textarea);
  textarea.select();
  try {
    document.execCommand('copy');
    copyBtn.textContent = 'Copied!';
    copyBtn.classList.add('copied');
    addLog('Channel key copied to clipboard', 'success');
    setTimeout(() => {
      copyBtn.textContent = 'Copy';
      copyBtn.classList.remove('copied');
    }, 2000);
  } catch {
    addLog('Failed to copy - select and copy manually', 'warn');
  }
  document.body.removeChild(textarea);
}

// ---- Event Listeners ----
connectBtn.addEventListener('click', doConnect);
copyBtn.addEventListener('click', doCopy);

portInput.addEventListener('keydown', (e: KeyboardEvent) => {
  if (e.key === 'Enter') doConnect();
});

// ---- Initial Log ----
addLog('AI Happy Design plugin ready');
addLog(`Channel key: ${channelKey}`, 'info');
addLog('Click Connect to start', 'info');
