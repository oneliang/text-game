package logic

import (
	"model"
	"view"
)

type Operation interface {
	// Operate . input event, output displayable
	Operate(event model.Event) view.Displayable

	// GetNextOperation .
	GetNextOperation() Operation

	// LoadSavedData .
	LoadSavedData(dataMap map[string]any)

	// GetNeedToSavedData .
	GetNeedToSavedData() map[string]any
}
