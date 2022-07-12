package flowchart

import (
	"bytes"
	"fmt"
)

// docs: https://mermaid-js.github.io/mermaid/#/flowchart

const (
	ChartName = "flowchart"
)

var (
	OrientationTB = Orientation("TB")
	OrientationTD = Orientation("TD")
	OrientationBT = Orientation("BT")
	OrientationRL = Orientation("RL")
	OrientationLR = Orientation("LR")

	ShapeBox = Shape{
		StartChar: "[",
		EndChar:   "]",
	}
	ShapeRoundEdges = Shape{
		StartChar: "(",
		EndChar:   ")",
	}
	ShapeStadium = Shape{
		StartChar: "([",
		EndChar:   "])",
	}
	ShapeSubroutine = Shape{
		StartChar: "[[",
		EndChar:   "]]",
	}

	LinkArrow = Link{
		Definition: "-->",
	}
	LinkOpen = Link{
		Definition: "---",
	}
)

type InlineItem interface {
	OutputInlineItem() string
}

type Orientation string

type Shape struct {
	StartChar string
	EndChar   string
}

type Node struct {
	ID    string
	Text  string
	Shape *Shape
}

func (N *Node) String() string {
	if N.Shape == nil {
		N.Shape = &ShapeBox
	}

	if N.ID == "" {
		N.ID = generateID()
	}

	return fmt.Sprintf(
		"%s%s%s%s",
		N.ID,
		N.Shape.StartChar,
		N.Text,
		N.Shape.EndChar,
	)
}

func (N *Node) OutputInlineItem() string {
	return N.String()
}

type Link struct {
	Definition string
}

func (L *Link) String() string {
	return L.Definition
}

func (L *Link) WithText(txt string) *Link {
	r := L
	r.Definition = r.Definition + fmt.Sprintf("|%s|", txt)

	return r
}

func (L *Link) OutputInlineItem() string {
	return L.String()
}

type Line struct {
	InlineItems []InlineItem
}

func (L *Line) String() string {
	var b bytes.Buffer

	for _, item := range L.InlineItems {
		b.WriteString(item.OutputInlineItem())
	}

	return b.String()
}

func (L *Line) OutputInlineItem() string {
	return L.String()
}

type Flowchart struct {
	DefaultShape *Shape
	Orientation  *Orientation
	Lines        []Line
}

func (flow *Flowchart) String() string {
	var b bytes.Buffer

	b.WriteString(ChartName)
	b.WriteString(fmt.Sprintf(" %s", *flow.Orientation))
	b.WriteString("\n")

	for _, l := range flow.Lines {
		b.WriteString("    ")
		b.WriteString(l.OutputInlineItem())
		b.WriteString("\n")
	}

	return b.String()
}

func New(o Orientation, ds *Shape) *Flowchart {
	if o == "" {
		o = OrientationTB
	}

	if ds == nil {
		ds = &ShapeBox
	}

	return &Flowchart{
		DefaultShape: ds,
		Orientation:  &o,
	}
}
