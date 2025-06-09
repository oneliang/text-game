package view

import "github.com/oneliang/util-golang/constants"

type TextView struct {
	text string
}

func NewTextView(text string) *TextView {
	return &TextView{
		text: text,
	}
}

func (this *TextView) Display() string {
	return this.text + constants.STRING_LF
}
