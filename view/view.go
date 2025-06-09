package view

import "github.com/oneliang/util-golang/constants"

type Displayable interface {
	Display() string
}

type View struct {
}

func (this *View) Display() string {
	return constants.STRING_BLANK + constants.STRING_LF
}
