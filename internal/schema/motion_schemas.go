package schema

func init() {
	nodeID := Param{Name: "nodeId", Type: "string", Required: true, Pattern: `^[0-9]+:[0-9]+$`}
	for _, s := range []Schema{
		{Command: "motion.get_styles", Description: "List available Motion styles where supported by the Figma runtime", Params: nil, Safety: "read"},
		{Command: "motion.apply_style", Description: "Apply a Motion style to a node where supported", Params: []Param{nodeID, {Name: "styleId", Type: "string", Required: true}}},
		{Command: "motion.remove_style", Description: "Remove Motion style from a node where supported", Params: []Param{nodeID}, Safety: "destructive"},
		{Command: "motion.get_animations", Description: "Read Motion animations from a node where supported", Params: []Param{nodeID}, Safety: "read"},
		{Command: "motion.apply_keyframes", Description: "Apply Motion keyframes to a node where supported", Params: []Param{nodeID, {Name: "keyframes", Type: "array", Required: true}, {Name: "duration", Type: "number", Min: Ptr(0)}}},
		{Command: "motion.remove_keyframes", Description: "Remove Motion keyframes from a node where supported", Params: []Param{nodeID}, Safety: "destructive"},
		{Command: "motion.set_timeline_duration", Description: "Set Motion timeline duration where supported", Params: []Param{{Name: "duration", Type: "number", Required: true, Min: Ptr(0)}}},
	} {
		Register(s)
	}
}
