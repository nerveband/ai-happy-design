/**
 * WebSocket client with auto-reconnection and message queuing.
 * Runs in the UI iframe context (has access to WebSocket API).
 */

export type ConnectionStatus = 'disconnected' | 'connecting' | 'connected' | 'reconnecting';

export interface WSClientOptions {
  relayUrl: string;
  channelKey: string;
  onMessage: (msg: any) => void;
  onStatusChange: (status: ConnectionStatus) => void;
  onLog: (message: string, level?: 'info' | 'warn' | 'error' | 'success') => void;
  onConnectionError?: (code: number, reason: string) => void;
}

export class WSClient {
  private ws: WebSocket | null = null;
  private options: WSClientOptions;
  private reconnectAttempts = 0;
  // 0 = retry forever (turnkey startup if relay starts later).
  private maxReconnectAttempts = 0;
  private hasLoggedTroubleshooting = false;
  private baseDelay = 1000;
  private maxDelay = 30000;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private messageQueue: any[] = [];
  private intentionalClose = false;
  private pingInterval: ReturnType<typeof setInterval> | null = null;
  private awaitingPong = false;

  constructor(options: WSClientOptions) {
    this.options = options;
  }

  connect(): void {
    if (this.ws && (this.ws.readyState === WebSocket.CONNECTING || this.ws.readyState === WebSocket.OPEN)) {
      return;
    }

    this.intentionalClose = false;
    this.options.onStatusChange('connecting');
    this.options.onLog(`Connecting to ${this.options.relayUrl}...`);

    try {
      this.ws = new WebSocket(this.options.relayUrl);
    } catch (err: any) {
      this.options.onLog(`Failed to create WebSocket: ${err.message}`, 'error');
      this.scheduleReconnect();
      return;
    }

    this.ws.onopen = () => {
      this.reconnectAttempts = 0;
      this.options.onStatusChange('connected');
      this.options.onLog(`Connected! Joining channel: ${this.options.channelKey}`);

      // Join channel
      this.sendRaw({
        type: 'join',
        channel: this.options.channelKey,
        role: 'plugin',
      });

      // Flush queued messages
      this.flushQueue();

      // Start ping
      this.startPing();
    };

    this.ws.onmessage = (event: MessageEvent) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === 'pong') {
          this.awaitingPong = false;
        }
        this.options.onMessage(data);
      } catch (err: any) {
        this.options.onLog(`Invalid message: ${err.message}`, 'warn');
      }
    };

    this.ws.onclose = (event: CloseEvent) => {
      this.stopPing();
      if (this.intentionalClose) {
        this.options.onStatusChange('disconnected');
        this.options.onLog('Disconnected');
        return;
      }

      // Actionable error messages based on close code
      if (event.code === 1006) {
        this.options.onLog('Relay not running or unreachable.', 'error');
        this.options.onLog('Start it with: ai-happy-design relay start', 'warn');
      } else if (event.code === 1008 || event.code === 1003) {
        this.options.onLog('Relay rejected the connection (code ' + event.code + '). Check your channel key.', 'error');
      } else {
        this.options.onLog('Connection lost (code: ' + event.code + '). Reconnecting...', 'warn');
      }

      if (this.options.onConnectionError) {
        this.options.onConnectionError(event.code, event.reason || '');
      }

      // Log troubleshooting steps on first failure only
      if (!this.hasLoggedTroubleshooting) {
        this.hasLoggedTroubleshooting = true;
        this.options.onLog('--- Troubleshooting ---', 'info');
        this.options.onLog('1. Run: ai-happy-design relay start', 'info');
        this.options.onLog('2. Check: ai-happy-design relay status', 'info');
        this.options.onLog('3. Verify the relay URL matches the port', 'info');
        this.options.onLog('-----------------------', 'info');
      }

      this.scheduleReconnect();
    };

    this.ws.onerror = () => {
      this.options.onLog('WebSocket error — is the relay running?', 'error');
    };
  }

  disconnect(): void {
    this.intentionalClose = true;
    this.cancelReconnect();
    this.stopPing();
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    this.options.onStatusChange('disconnected');
  }

  send(message: any): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    } else {
      // Queue message for later delivery
      this.messageQueue.push(message);
      if (this.messageQueue.length > 100) {
        this.messageQueue.shift(); // Drop oldest if queue too large
      }
    }
  }

  updateOptions(options: Partial<WSClientOptions>): void {
    Object.assign(this.options, options);
  }

  get isConnected(): boolean {
    return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
  }

  private sendRaw(message: any): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    }
  }

  private flushQueue(): void {
    while (this.messageQueue.length > 0) {
      const msg = this.messageQueue.shift();
      this.send(msg);
    }
  }

  private scheduleReconnect(): void {
    if (this.maxReconnectAttempts > 0 && this.reconnectAttempts >= this.maxReconnectAttempts) {
      this.options.onStatusChange('disconnected');
      this.options.onLog('Max reconnect attempts reached. Please reconnect manually.', 'error');
      return;
    }

    this.options.onStatusChange('reconnecting');
    this.reconnectAttempts++;

    // Exponential backoff with jitter
    const delay = Math.min(
      this.baseDelay * Math.pow(2, this.reconnectAttempts - 1) + Math.random() * 1000,
      this.maxDelay
    );

    const attemptLabel = this.maxReconnectAttempts > 0
      ? `${this.reconnectAttempts}/${this.maxReconnectAttempts}`
      : `${this.reconnectAttempts}`;
    this.options.onLog(`Reconnecting in ${(delay / 1000).toFixed(1)}s (attempt ${attemptLabel})...`);

    this.reconnectTimer = setTimeout(() => {
      this.connect();
    }, delay);
  }

  private cancelReconnect(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    this.reconnectAttempts = 0;
  }

  private startPing(): void {
    this.stopPing();
    this.awaitingPong = false;
    this.pingInterval = setInterval(() => {
      if (this.awaitingPong) {
        // No pong received since last ping — connection is dead
        this.options.onLog('No response from relay — connection lost', 'error');
        this.awaitingPong = false;
        if (this.ws) {
          this.ws.close();
        }
        return;
      }
      this.awaitingPong = true;
      this.sendRaw({ type: 'ping' });
    }, 30000);
  }

  private stopPing(): void {
    if (this.pingInterval) {
      clearInterval(this.pingInterval);
      this.pingInterval = null;
    }
  }
}
