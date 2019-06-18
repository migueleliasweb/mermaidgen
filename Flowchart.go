package mermaidgen

import (
	"fmt"
)

////////// ChartDirection //////////////////////////////////////////////////////

type chartDirection string

// Direction definitions for Flowcharts as described at
// https://mermaidjs.github.io/flowchart.html#graph.
// New Flowcharts get DirectionTopDown as the default.
const (
	DirectionTopDown   chartDirection = `TB`
	DirectionBottomUp  chartDirection = `BT`
	DirectionRightLeft chartDirection = `RL`
	DirectionLeftRight chartDirection = `LR`
)

////////// GraphItem ///////////////////////////////////////////////////////////

// interface to define what can be an "item" to a Flowchart/Subgraph
type graphItem interface {
	renderGraph() string
}

////////// Flowchart ///////////////////////////////////////////////////////////

// Flowchart objects are the entrypoints to this package, the whole graph is
// constructed around a Flowchart object. Create an instance of Flowchart via
// Flowchart's constructor NewFlowchart, do not create instances directly.
type Flowchart struct {
	nodeStyles map[string]*NodeStyle // internal storage for NodeStyles
	edgeStyles map[string]*EdgeStyle // internal storage for EdgeStyles
	subgraphs  map[string]*Subgraph  // internal storage for Subgraphs
	nodes      map[string]*Node      // internal storage for Nodes
	edges      []*Edge               // internal storage for Edges
	items      []graphItem           // sub-items to render
	Direction  chartDirection        // The direction used to render the graph.
}

// NewFlowchart is the constructor used to create a new Flowchart object.
// This object is the entrypoint for any further interactions with your graph.
// Always use the constructor, don't create Flowchart objects directly.
func NewFlowchart() (newFlowchart *Flowchart) {
	f := &Flowchart{}
	f.Direction = DirectionTopDown
	f.nodeStyles = make(map[string]*NodeStyle)
	f.edgeStyles = make(map[string]*EdgeStyle)
	f.subgraphs = make(map[string]*Subgraph)
	f.nodes = make(map[string]*Node)
	return f
}

// String recursively renders the whole graph to mermaid code lines.
func (fc *Flowchart) String() (renderedElement string) {
	text := fmt.Sprintf("graph %s\n", fc.Direction)
	for _, s := range fc.nodeStyles {
		text += s.String()
	}
	for _, item := range fc.items {
		text += item.renderGraph()
	}
	for _, e := range fc.edges {
		text += e.String()
	}
	return text
}

////////// add & get Styles ////////////////////////////////////////////////////

// NodeStyle is used to create new or lookup existing NodeStyles by ID.
// The returned object pointers can be assigned to any number of Nodes
// to style them using CSS.
func (fc *Flowchart) NodeStyle(id string) (style *NodeStyle) {
	s, found := fc.nodeStyles[id]
	if !found {
		s = &NodeStyle{id: id, StrokeWidth: 1}
		fc.nodeStyles[id] = s
	}
	return s
}

// EdgeStyle is used to create new or lookup existing EdgeStyles by ID.
// The returned object pointers can be assigned to any number of Edges
// to style them using CSS. Note that EdgeStyles override the shape of an Edge,
// e.g. if you color an Edge that uses EShapeDottedArrow it looses its dotted
// nature unless you define a dotted line using the EdgeStyle.
func (fc *Flowchart) EdgeStyle(id string) (style *EdgeStyle) {
	s, found := fc.edgeStyles[id]
	if !found {
		s = &EdgeStyle{id: id, StrokeWidth: 1}
		fc.edgeStyles[id] = s
	}
	return s
}

////////// add Items ///////////////////////////////////////////////////////////

// AddSubgraph is used to add a nested Subgraph to the Flowchart.
// If the provided ID already exists, no new Subgraph is created and nil is
// returned. The ID can later be used to lookup the created Subgraph using
// Flowchart's GetSubgraph method. If you want to add a Subgraph to a Subgraph,
// use that Subgraph's AddSubgraph method.
func (fc *Flowchart) AddSubgraph(id string) (newSubgraph *Subgraph) {
	_, found := fc.subgraphs[id]
	if found {
		// if already exists -> nil
		return nil
	} else {
		s := &Subgraph{id: id, flowchart: fc}
		fc.subgraphs[id] = s
		fc.items = append(fc.items, s)
		return s
	}
}

// AddNode is used to add a new Node to the Flowchart. If the provided ID
// already exists, no new Node is created and nil is returned. The ID can later
// be used to lookup the created Node using Flowchart's GetNode method.
// If you want to add a Node to a Subgraph, use that Subgraph's AddNode method.
func (fc *Flowchart) AddNode(id string) (newNode *Node) {
	_, found := fc.nodes[id]
	if found {
		// if already exists -> nil
		return nil
	} else {
		n := &Node{id: id, Shape: NShapeRect}
		fc.nodes[id] = n
		fc.items = append(fc.items, n)
		return n
	}
}

// AddEdge is used to add a new Edge to the Flowchart. Since Edges have no IDs
// this will always succeed. The (pseudo) ID is the index that defines the order
// of all Edges and is used to define linkStyles. The ID can later be used to
// lookup the created Edge using Flowchart's GetEdge method.
func (fc *Flowchart) AddEdge(from *Node, to *Node) (newEdge *Edge) {
	e := &Edge{From: from, To: to, Shape: EShapeArrow}
	fc.edges = append(fc.edges, e)
	e.id = len(fc.edges) - 1
	return e
}

////////// get Items ///////////////////////////////////////////////////////////

// GetSubgraph looks up a previously defined Subgraph by its ID.
// If this ID doesn't exist, nil is returned.
// Use Flowchart's or Subgraph's AddSubgraph to create new Subgraphs.
func (fc *Flowchart) GetSubgraph(id string) (existingSubgraph *Subgraph) {
	// if not found -> nil
	return fc.subgraphs[id]
}

// GetNode looks up a previously defined Node by its ID.
// If this ID doesn't exist, nil is returned.
// Use Flowchart's or Subgraph's AddNode to create new Nodes.
func (fc *Flowchart) GetNode(id string) (existingNode *Node) {
	// if not found -> nil
	return fc.nodes[id]
}

// GetEdge looks up a previously defined Edge by its ID (index).
// If this index doesn't exist, nil is returned.
// Use Flowchart's AddEdge to create new Edges.
func (fc *Flowchart) GetEdge(index int) (existingEdge *Edge) {
	if index < 0 || len(fc.edges) <= index {
		return nil
	}
	return fc.edges[index]
}

////////// list Items //////////////////////////////////////////////////////////

// ListSubgraphs returns a slice of all previously defined Subgraphs.
// The order is not well-defined.
func (fc *Flowchart) ListSubgraphs() (allSubgraphs []*Subgraph) {
	values := make([]*Subgraph, 0, len(fc.subgraphs))
	for _, v := range fc.subgraphs {
		values = append(values, v)
	}
	return values
}

// ListNodes returns a slice of all previously defined Nodes.
// The order is not well-defined.
func (fc *Flowchart) ListNodes() (allNodes []*Node) {
	values := make([]*Node, 0, len(fc.nodes))
	for _, v := range fc.nodes {
		values = append(values, v)
	}
	return values
}

// ListEdges returns a slice of all previously defined Edges in the order they
// were added.
func (fc *Flowchart) ListEdges() (allEdges []*Edge) {
	e := make([]*Edge, len(fc.edges))
	copy(e, fc.edges)
	return e
}
