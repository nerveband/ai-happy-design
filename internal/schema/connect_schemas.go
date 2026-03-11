package schema

func init() {
	Register(Schema{
		Command:     "connect.join",
		Description: "Join a channel to communicate with the Figma plugin",
		Params: []Param{
			{Name: "channelKey", Type: "string", Required: true, Desc: "Channel key to join"},
		},
	})

	Register(Schema{
		Command:     "connect.status",
		Description: "Get the current connection status",
		Params:      []Param{},
	})

	Register(Schema{
		Command:     "connect.disconnect",
		Description: "Disconnect from the current channel",
		Params:      []Param{},
	})
}
