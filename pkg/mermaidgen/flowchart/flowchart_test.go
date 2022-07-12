package flowchart

import (
	"fmt"
	"testing"
)

func TestFlowchartWithID(t *testing.T) {
	fc := New("", nil)
	fc.Lines = []Line{
		{
			InlineItems: []InlineItem{
				&Node{
					ID:   "A",
					Text: "Christmas",
				},
				&LinkArrow,
				&Node{
					ID:   "B",
					Text: "Go shopping",
				},
			},
		},
	}

	fmt.Println(fc)
}

func TestFlowchartWithText(t *testing.T) {
	fc := New(OrientationLR, nil)
	fc.Lines = []Line{
		{
			InlineItems: []InlineItem{
				&Node{
					ID:   "A",
					Text: "Christmas",
				},
				LinkArrow.WithText("inline text"),
				&Node{
					ID:   "B",
					Text: "Go shopping",
				},
			},
		},
	}

	fmt.Println(fc)
}
