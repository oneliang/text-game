package logic

import (
	"fmt"
	"github.com/oneliang/util-golang/common"
	"github.com/oneliang/util-golang/logging"
	"model"
	"strings"
	"view"
)

const systemOperationLoggerTag = "SystemOperation"
const map_01 = 0x01010001

type SystemOperation struct {
	nextOperation     Operation
	logger            logging.Logger
	resourceManager   *model.ResourceManager
	playerDataManager *model.PlayerDataManager
	operationManager  *OperationManager
	operationStack    *common.Stack[Operation]
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
	nextOperationDisplayable := this.nextOperation.Operate(event)
	if this.operationManager.Size() > 1 {
		cancelKeyCode := model.EVENT_KEY_CODE_MAPPING[model.EVENT_CANCEL]
		return view.NewViewGroup(
			nextOperationDisplayable,
			view.NewTextView("-----OTHER CONTROL TIPS CONTENT BEGIN-----"),
			view.NewButtonView(model.EVENT_CANCEL, fmt.Sprintf("%s: Return to previous", strings.ToUpper(string(cancelKeyCode)))),
			view.NewTextView("-----OTHER CONTROL TIPS CONTENT END-----"),
		)
	}
	return nextOperationDisplayable
}

// GetNextOperation .
func (this *SystemOperation) GetNextOperation() Operation {
	operation := this.operationManager.NewMapOperation(map_01, 0, 1, this.resourceManager, this.playerDataManager, true)
	operation.SetOperationManager(this.operationManager)
	return operation
}

// SetOperationManager .
func (this *SystemOperation) SetOperationManager(operationManager *OperationManager) {
	this.operationManager = operationManager
}
