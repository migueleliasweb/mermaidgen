package mermaidgen

import (
	"fmt"
	"strings"
)

// A NodeStyle is used to add CSS to a Node. It renders to a classDef line.
// Retrieve an instance of NodeStyle via Flowchart's NodeStyle method, do not
// create instances directly.
type NodeStyle struct {
	id          string
	Fill        htmlColor // renders to something like fill:#f9f
	Stroke      htmlColor // renders to something like stroke:#333
	StrokeWidth uint8     // renders to something like stroke-width:2px
	StrokeDash  uint8     // renders to something like stroke-dasharray:5px
	More        string    // more styles, e.g.: stroke:#333,stroke-width:1px
}

// ID provides access to the NodeStyle's readonly field id.
func (ns *NodeStyle) ID() (id string) {
	return ns.id
}

// String renders this graph element to a classDef line.
func (ns *NodeStyle) String() (renderedElement string) {
	styles := []string{
		fmt.Sprintf(`stroke-width:%dpx`, ns.StrokeWidth),
	}
	if ns.Fill != "" {
		styles = append(styles, fmt.Sprintf(`fill:%s`, ns.Fill))
	}
	if ns.Stroke != "" {
		styles = append(styles, fmt.Sprintf(`stroke:%s`, ns.Stroke))
	}
	if ns.StrokeDash != 0 {
		styles = append(styles, fmt.Sprintf(`stroke-dasharray:%dpx`,
			ns.StrokeDash))
	}
	if ns.More != "" {
		styles = append(styles, ns.More)
	}
	return fmt.Sprintf("classDef %s %s\n", ns.id, strings.Join(styles, ","))
}
