package view

import (
	"fmt"
	"github.com/oneliang/util-golang/constants"
	"model"
)

type ButtonView struct {
	event model.Event
	text  string
}

func NewButtonView(event model.Event, text string) *ButtonView {
	return &ButtonView{
		event: event,
		text:  text,
	}
}

func (this *ButtonView) Display() string {
	return fmt.Sprintf("%s", this.text) + constants.STRING_LF
}
