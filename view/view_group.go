package view

import (
	"fmt"
	"strings"
)

type ViewGroup struct {
	ChildList []Displayable
}

func NewViewGroup(child ...Displayable) *ViewGroup {
	childList := make([]Displayable, 0)
	if child != nil {
		childList = append(childList, child...)
	}
	return &ViewGroup{
		ChildList: childList,
	}
}
func (this *ViewGroup) Display() string {
	var stringBuilder strings.Builder

	if this.ChildList != nil {
		for index, displayable := range this.ChildList {
			if displayable != nil {
				stringBuilder.WriteString(displayable.Display())
			} else {
				fmt.Println(fmt.Sprintf("ViewGroup display child list, child is nil, viewGroup:%+v, child index:%d", this, index))
			}
		}
	}
	return stringBuilder.String()
}
