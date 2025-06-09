package view

import (
	"fmt"
	"testing"
)

func TestView(t *testing.T) {
	viewGroup := NewViewGroup(
		NewTextView("1"),
		NewTextView("2"),
		NewViewGroup(
			NewTextView("2_1"),
			NewTextView("2_2"),
		),
	)
	fmt.Println(viewGroup.Display())
}
