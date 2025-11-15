package logic

import (
	"errors"
	"fmt"
	"github.com/oneliang/util-golang/base"
	"github.com/oneliang/util-golang/constants"
	"github.com/oneliang/util-golang/logging"
	"log_content"
	"model"
	"strings"
	"view"
)

const playerOperationLoggerTag = "PlayerOperation"

type PlayerOperation struct {
	player            *model.Player
	currentIndex      int
	resourceManager   *model.ResourceManager
	playerDataManager *model.PlayerDataManager
	operationManager  *OperationManager
	logger            logging.Logger
}

// NewPlayerOperation .
func NewPlayerOperation(
	resourceManager *model.ResourceManager,
	playerDataManager *model.PlayerDataManager,
) *PlayerOperation {
	return &PlayerOperation{
		resourceManager:   resourceManager,
		playerDataManager: playerDataManager,
		player:            playerDataManager.GetPlayer(),
		logger:            logging.LoggerManager.GetLogger(playerOperationLoggerTag),
	}
}

// Operate .
func (this *PlayerOperation) Operate(event model.Event) view.Displayable {
	var err error = nil
	switch event {
	case model.EVENT_UP:
		_, _, err = this.MoveUp()
		break
	case model.EVENT_DOWN:
		_, _, err = this.MoveDown()
		break
	case model.EVENT_CONFIRM:
		err = this.ClickConfirm()
		break
	case model.EVENT_CANCEL:
		return nil
	case model.EVENT_NONE:
		break
	default:
		this.logger.Debug(log_content.LogContentNormal(playerOperationLoggerTag, "not support event, %d", event))
	}
	if err != nil {
		this.logger.Debug(log_content.LogContentNormal(playerOperationLoggerTag, "err:%+v", err))
	}
	return this.getView()
}

// GetNextOperation .
func (this *PlayerOperation) GetNextOperation() Operation {
	return nil
}

// SetOperationManager .
func (this *PlayerOperation) SetOperationManager(operationManager *OperationManager) {
	this.operationManager = operationManager
}

// GetCurrentLocation .
func (this *PlayerOperation) GetCurrentLocation() int {
	return this.currentIndex
}

// SetLocationResource .
func (this *PlayerOperation) SetLocationResource(index uint32, resourceRealId uint32, count uint8) error {
	totalItemLength := uint32(len(this.player.ItemList))

	if index < totalItemLength {
		this.player.ItemList[index] = base.NewKeyValue[uint32, uint8](resourceRealId, count)
	} else {
		return errors.New(fmt.Sprintf("index not match, max item length:%d, index:%d", totalItemLength, index))
	}
	return nil
}

// GetLocationResource .
func (this *PlayerOperation) GetLocationResource(index int) (uint32, uint8, error) {
	totalItemLength := len(this.player.ItemList)
	if index > (totalItemLength-1) || index < 0 {
		return 0, 0, errors.New(fmt.Sprintf("index not match, max item length:%d, index:%d", totalItemLength, index))
	}
	item := this.player.ItemList[index]
	return item.Key, item.Value, nil
}

// GetCurrentLocationResource .
func (this *PlayerOperation) GetCurrentLocationResource() (uint32, uint8, error) {
	return this.GetLocationResource(this.currentIndex)
}

// MoveUp .
func (this *PlayerOperation) MoveUp() (uint32, uint8, error) {
	previousIndex := this.currentIndex - 1
	resourceRealId, count, err := this.GetLocationResource(previousIndex)
	if err != nil {
		return 0, 0, err
	}
	this.currentIndex = previousIndex
	return resourceRealId, count, nil
}

// MoveDown .
func (this *PlayerOperation) MoveDown() (uint32, uint8, error) {
	nextIndex := this.currentIndex + 1
	resourceRealId, count, err := this.GetLocationResource(nextIndex)
	if err != nil {
		return 0, 0, err
	}
	this.currentIndex = nextIndex
	return resourceRealId, count, nil
}

func (this *PlayerOperation) ClickConfirm() error {
	resourceRealId, count, err := this.GetLocationResource(this.currentIndex)
	if err != nil {
		return errors.New(fmt.Sprintf("location resource can not confirm, current (index:%d)", this.currentIndex))
	}
	resource, err := this.resourceManager.GetItemWithRealId(resourceRealId)
	if resource.State == model.RESOURCE_ITEM_STATE_CAN_OPEN && count > 0 {
		//add item to player item list
		this.player.DeleteItem(resourceRealId)
		//newResourceId := model.ResourceUnsetState(resourceId, model.RESOURCE_ITEM_STATE_CAN_OPEN)
		//err = this.SetLocationResource(this.currentX, this.currentY, newResourceId, 0)
		//if err != nil {
		//	return errors.New(fmt.Sprintf("location resource can not set, current (x:%d, y:%d) err:%v", this.currentX, this.currentY, err))
		//}
	} else {
		this.logger.Warning(log_content.LogContentNormal(playerOperationLoggerTag, "can not open, current (index:%d), resource id:%x", this.currentIndex, resourceRealId))
	}
	return nil
}

func (this *PlayerOperation) NearbyResourceList() []*base.Pair[int, uint32] {
	upDownDirectionList := []int{
		this.currentIndex - 1, //up
		this.currentIndex + 1, //down
	}
	resourceRealIdList := make([]*base.Pair[int, uint32], 0)
	for _, upDownDirection := range upDownDirectionList {

		currentIndex := upDownDirection
		resourceRealId, _, err := this.GetLocationResource(currentIndex)
		if err != nil {
			resourceRealIdList = append(resourceRealIdList, nil)
			continue
		}
		pairData := base.NewPair[int, uint32](currentIndex, resourceRealId)
		resourceRealIdList = append(resourceRealIdList, pairData)
	}
	return resourceRealIdList
}

func (this *PlayerOperation) getItemViewList() []view.Displayable {
	viewList := make([]view.Displayable, 0)
	for index, item := range this.player.ItemList {
		var stringBuilder strings.Builder
		resourceRealId := item.Key
		itemResource, err := this.resourceManager.GetItemWithRealId(resourceRealId)
		currentFlag := constants.STRING_BLANK
		if this.currentIndex == index {
			currentFlag = constants.EMOJI_GRINNING_FACE
		}
		if err != nil {
			this.logger.Error(log_content.LogContentNormal(playerOperationLoggerTag, constants.STRING_BLANK), err)
			stringBuilder.WriteString(fmt.Sprintf("%-24s", fmt.Sprintf("%d:ERR", index)))
		} else {

			stringBuilder.WriteString(fmt.Sprintf("%-24s", fmt.Sprintf("%s%d:%s(%d)", currentFlag, index, itemResource.Name, item.Value)))
		}

		viewList = append(viewList, view.NewTextView(stringBuilder.String()))
	}
	return viewList
}

func (this *PlayerOperation) getView() view.Displayable {
	viewList := make([]view.Displayable, 0)

	itemViewList := this.getItemViewList()
	//fmt.Println(itemViewList)
	viewList = append(viewList, itemViewList...)

	// nearby resource list
	nearByResourceList := this.NearbyResourceList()
	defaultEventList := []model.Event{
		model.EVENT_UP,
		model.EVENT_DOWN,
	}

	itemTipsViewList := make([]view.Displayable, 0)

	for index, locationResourceTuple := range nearByResourceList {
		if locationResourceTuple == nil {
			continue
		}
		currentIndex := locationResourceTuple.First
		resourceRealId := locationResourceTuple.Second
		itemResource, err := this.resourceManager.GetItemWithRealId(resourceRealId)
		//resourceRealId := this.resourceManager.GetResourceRealId(resourceId)
		//mapThing, err := this.resourceManager.GetMapThing(resourceRealId)
		if err != nil {
			fmt.Println(fmt.Sprintf("%v", err))
		}
		event := defaultEventList[index]
		keyCode := model.EVENT_KEY_CODE_MAPPING[event]
		viewList = append(viewList, view.NewButtonView(event, fmt.Sprintf("%s: %s (index:%d)", strings.ToUpper(string(keyCode)), itemResource.Name, currentIndex)))
	}

	resourceId, _, err := this.GetLocationResource(this.currentIndex)
	if err == nil {
		resource, err := this.resourceManager.GetResource(resourceId)
		if err == nil {
			if model.ResourceTypeIsItem(resourceId) && model.ResourceItemStateCanOpen(resourceId) {
				confirmKeyCode := model.EVENT_KEY_CODE_MAPPING[model.EVENT_CONFIRM]
				itemTipsViewList = append(itemTipsViewList, view.NewButtonView(model.EVENT_CONFIRM, fmt.Sprintf("%s: %s (index:%d) can open", strings.ToUpper(string(confirmKeyCode)), resource.Name, this.currentIndex)))
			}
		}
	}

	viewList = append(viewList, itemTipsViewList...)

	return view.NewViewGroup(viewList...)
}
