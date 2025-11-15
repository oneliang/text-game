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

	// SetOperationManager .
	SetOperationManager(operationManager *OperationManager)
}
