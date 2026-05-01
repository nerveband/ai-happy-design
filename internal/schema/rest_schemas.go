package schema

func init() {
	fileKey := Param{Name: "fileKey", Type: "string", Required: true, Desc: "Figma file key"}
	token := Param{Name: "token", Type: "string", Desc: "Optional Figma token; defaults to FIGMA_ACCESS_TOKEN or FIGMA_TOKEN"}
	baseURL := Param{Name: "baseURL", Type: "string", Desc: "Optional API base URL for tests or proxies"}
	webhookID := Param{Name: "webhookId", Type: "string", Required: true}
	for _, s := range []Schema{
		{Command: "figma.oembed", Aliases: []string{"rest.oembed"}, Description: "Fetch Figma oEmbed metadata for a file or Make URL", Safety: "read", Idempotency: "idempotent", RequiresFigma: false, Params: []Param{{Name: "url", Type: "string", Required: true}, token, baseURL}},
		{Command: "figma.file_metadata", Aliases: []string{"rest.file_metadata"}, Description: "Fetch Figma REST file metadata using least-privilege scopes", Safety: "read", Idempotency: "idempotent", RequiresFigma: false, Params: []Param{fileKey, {Name: "ids", Type: "string"}, {Name: "version", Type: "string"}, {Name: "geometry", Type: "string"}, {Name: "plugin_data", Type: "string"}, token, baseURL}},
		{Command: "figma.dev_resources_list", Aliases: []string{"dev_resource.list"}, Description: "List Dev Mode resources attached to a Figma file", Safety: "read", Idempotency: "idempotent", RequiresFigma: false, Params: []Param{fileKey, token, baseURL}},
		{Command: "figma.dev_resource_create", Aliases: []string{"dev_resource.create"}, Description: "Create a Dev Mode resource link for a node", Safety: "write", Idempotency: "non-idempotent", RequiresFigma: false, Params: []Param{fileKey, {Name: "nodeId", Type: "string", Required: true}, {Name: "name", Type: "string", Required: true}, {Name: "url", Type: "string", Required: true}, token, baseURL}},
		{Command: "figma.dev_resource_update", Aliases: []string{"dev_resource.update"}, Description: "Update a Dev Mode resource link", Safety: "write", Idempotency: "idempotent", RequiresFigma: false, Params: []Param{fileKey, {Name: "resourceId", Type: "string", Required: true}, {Name: "name", Type: "string"}, {Name: "url", Type: "string"}, token, baseURL}},
		{Command: "figma.dev_resource_delete", Aliases: []string{"dev_resource.delete"}, Description: "Delete a Dev Mode resource link", Safety: "destructive", Idempotency: "idempotent", RequiresFigma: false, Params: []Param{fileKey, {Name: "resourceId", Type: "string", Required: true}, token, baseURL}},
		{Command: "figma.webhooks_list", Aliases: []string{"webhook.list"}, Description: "List Figma Webhooks V2 subscriptions", Safety: "read", Idempotency: "idempotent", RequiresFigma: false, Params: []Param{token, baseURL}},
		{Command: "figma.webhook_create", Aliases: []string{"webhook.create"}, Description: "Create a Figma Webhooks V2 subscription, including DEV_MODE_STATUS_UPDATE", Safety: "write", Idempotency: "non-idempotent", RequiresFigma: false, Params: []Param{{Name: "eventType", Type: "string", Required: true}, {Name: "context", Type: "string", Required: true}, {Name: "contextId", Type: "string", Required: true}, {Name: "endpoint", Type: "string", Required: true}, token, baseURL}},
		{Command: "figma.webhook_get", Aliases: []string{"webhook.get"}, Description: "Get one Figma Webhooks V2 subscription", Safety: "read", Idempotency: "idempotent", RequiresFigma: false, Params: []Param{webhookID, token, baseURL}},
		{Command: "figma.webhook_update", Aliases: []string{"webhook.update"}, Description: "Update one Figma Webhooks V2 subscription", Safety: "write", Idempotency: "idempotent", RequiresFigma: false, Params: []Param{webhookID, {Name: "eventType", Type: "string"}, {Name: "context", Type: "string"}, {Name: "contextId", Type: "string"}, {Name: "endpoint", Type: "string"}, token, baseURL}},
		{Command: "figma.webhook_delete", Aliases: []string{"webhook.delete"}, Description: "Delete one Figma Webhooks V2 subscription", Safety: "destructive", Idempotency: "idempotent", RequiresFigma: false, Params: []Param{webhookID, token, baseURL}},
	} {
		Register(s)
	}
}
