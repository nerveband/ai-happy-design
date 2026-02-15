export const DEFAULT_PORT = 3055;
export const DEFAULT_RELAY_URL = 'ws://localhost:' + DEFAULT_PORT + '/ws';

export function normalizeRelayUrl(raw: any): string {
  var input = typeof raw === 'string' ? raw.trim() : '';
  if (!input) return DEFAULT_RELAY_URL;

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

  var parts = candidate.match(/^(wss?):\/\/(\[[^\]]+\]|[^\/\s:?#]+)(?::(\d+))?(\/[^?#\s]*)?(?:\?[^#\s]*)?(?:#\S*)?$/i);
  if (!parts) return DEFAULT_RELAY_URL;

  var scheme = parts[1].toLowerCase();
  var host = parts[2];
  var port = parts[3];
  var path = parts[4] || '/ws';

  if (!host) return DEFAULT_RELAY_URL;
  if (port) {
    var portNum = parseInt(port, 10);
    if (portNum < 1 || portNum > 65535) return DEFAULT_RELAY_URL;
  }

  if (path === '/') path = '/ws';

  return scheme + '://' + host + (port ? ':' + port : '') + path;
}
