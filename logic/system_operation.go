package logic

import (
	"github.com/oneliang/util-golang/logging"
	"model"
	"view"
)

const systemOperationLoggerTag = "SystemOperation"
const map_01 = 0x01010001

type SystemOperation struct {
	nextOperation     Operation
	logger            logging.Logger
	resourceManager   *model.ResourceManager
	playerDataManager *model.PlayerDataManager
}

func NewSystemOperation(
	resourceManager *model.ResourceManager,
	playerDataManager *model.PlayerDataManager,
) *SystemOperation {
	return &SystemOperation{
		logger:            logging.LoggerManager.GetLogger(systemOperationLoggerTag),
		resourceManager:   resourceManager,
		playerDataManager: playerDataManager,
	}
}

// Operate . input event, output displayable
func (this *SystemOperation) Operate(event model.Event) view.Displayable {
	if this.nextOperation == nil {
		this.nextOperation = this.GetNextOperation()
	}
	return this.nextOperation.Operate(event)
}

// GetNextOperation .
func (this *SystemOperation) GetNextOperation() Operation {
	return NewMapOperation(map_01, 0, 1, this.resourceManager, this.playerDataManager, true)
}
