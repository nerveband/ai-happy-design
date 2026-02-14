package tools

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/figma"
	"github.com/nerveband/ai-happy-design/internal/ws"
)

// RegisterAll registers every domain tool with the MCP server.
func RegisterAll(s *server.MCPServer, wsServer *ws.Server) {
	commander := figma.NewCommander(wsServer)

	RegisterPaintTool(s, commander)
	RegisterShapeTool(s, commander)
	RegisterTextTool(s, commander)
	RegisterLayoutTool(s, commander)
	RegisterNodeTool(s, commander)
	RegisterLayerTool(s, commander)
	RegisterComponentTool(s, commander)
	RegisterStyleTool(s, commander)
	RegisterVariableTool(s, commander)
	RegisterEffectTool(s, commander)
	RegisterBooleanTool(s, commander)
	RegisterPageTool(s, commander)
	RegisterDocumentTool(s, commander)
	RegisterExportTool(s, commander)
	RegisterBulkTool(s, commander)
	RegisterDescribeTool(s, commander)
	RegisterConnectTool(s, commander, wsServer)
}
