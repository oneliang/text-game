package logic

import (
	"fmt"
	"github.com/oneliang/util-golang/logging"
	"log_content"
	"model"
)

const operationManagerLoggerTag = "OperationManager"

type OperationManager struct {
	operationMap map[Operation]Operation
	logger       logging.Logger
}

func NewOperationManager() *OperationManager {
	return &OperationManager{
		operationMap: make(map[Operation]Operation),
		logger:       logging.LoggerManager.GetLogger(operationManagerLoggerTag),
	}
}
func (this *OperationManager) NewMapOperation(
	mapResourceId uint32,
	initializeX int,
	initializeY int,
	resourceManager *model.ResourceManager,
	playerDataManager *model.PlayerDataManager,
	isRootMap bool,
) *MapOperation {
	operation := NewMapOperation(mapResourceId, initializeX, initializeY, resourceManager, playerDataManager, isRootMap)
	this.operationMap[operation] = operation
	return operation
}

func (this *OperationManager) NewNpcOperation(
	resourceManager *model.ResourceManager,
	resourceId uint32) *NpcOperation {
	operation := NewNpcOperation(resourceManager)
	this.operationMap[operation] = operation
	operation.SetCurrentNpcResourceId(resourceId)
	return operation
}

func (this *OperationManager) NewPlayerOperation(
	resourceManager *model.ResourceManager,
	playerDataManager *model.PlayerDataManager,
) *PlayerOperation {
	operation := NewPlayerOperation(resourceManager, playerDataManager)
	this.operationMap[operation] = operation
	return operation
}

func (this *OperationManager) DestroyOperation(operation Operation) {
	operation, exist := this.operationMap[operation]
	if exist {
		delete(this.operationMap, operation)
		return
	}
	this.logger.Warning(log_content.LogContentNormal(operationManagerLoggerTag, fmt.Sprintf("operation:%+v is not exist", operation)))
}

func (this *OperationManager) Size() int {
	return len(this.operationMap)
}
