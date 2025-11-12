package logic

import (
	"encoding/json"
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

const (
	MAP_OPERATION_DATA_KEY_PLAYER      = "PLAYER"
	MAP_OPERATION_DATA_KEY_CURRENT_MAP = "CURRENT_MAP"
	MAP_OPERATION_DATA_KEY_CURRENT_X   = "CURRENT_X"
	MAP_OPERATION_DATA_KEY_CURRENT_Y   = "CURRENT_Y"
)

const mapOperationLoggerTag = "MapOperation"

type MapOperation struct {
	player          *model.Player
	currentMap      *model.Map
	currentX        int
	currentY        int
	resourceManager *model.ResourceManager
	playerOperation *PlayerOperation
	npcOperation    *NpcOperation
	currentEvent    model.Event
	logger          logging.Logger
	nextOperation   Operation
	isRootMap       bool
}

// NewMapOperation .
func NewMapOperation(
	player *model.Player,
	inputMap *model.Map,
	initializeX int,
	initializeY int,
	resourceManager *model.ResourceManager,
	playerOperation *PlayerOperation,
	isRootMap bool) *MapOperation {
	return &MapOperation{
		player:          player,
		currentMap:      inputMap,
		currentX:        initializeX,
		currentY:        initializeY,
		resourceManager: resourceManager,
		playerOperation: playerOperation,
		npcOperation:    NewNpcOperation(resourceManager),
		logger:          logging.LoggerManager.GetLogger(mapOperationLoggerTag),
		nextOperation:   nil,
		isRootMap:       isRootMap,
	}
}

// Operate .
func (this *MapOperation) Operate(event model.Event) view.Displayable {
	if this.nextOperation != nil {
		return this.operateNextOperation(event)
	}
	// current operation
	this.currentEvent = event
	var err error = nil
	switch event {
	case model.EVENT_UP:
		_, err = this.MoveUp()
	case model.EVENT_DOWN:
		_, err = this.MoveDown()
	case model.EVENT_LEFT:
		_, err = this.MoveLeft()
	case model.EVENT_RIGHT:
		_, err = this.MoveRight()
	case model.EVENT_CONFIRM:
		nextOperation := this.GetNextOperation()
		if nextOperation != nil {
			return this.operateNextOperation(event)
		} else {
			err = this.ClickConfirm()
		}
	case model.EVENT_CANCEL:
		if this.isRootMap {
			return this.getView()
		} else {
			return nil
		}
	case model.EVENT_MENU:
		nextOperation := this.GetNextOperation()
		if nextOperation != nil {
			return this.operateNextOperation(event)
		}
	case model.EVENT_NONE:
	default:
		this.logger.Debug(log_content.LogContentNormal(mapOperationLoggerTag, "not support event, %d", event))
	}
	if err != nil {
		this.logger.Debug(log_content.LogContentNormal(mapOperationLoggerTag, "err:%+v", err))
	}
	return this.getView()
}

// GetNextOperation .
func (this *MapOperation) GetNextOperation() Operation {
	var operation Operation = nil
	if this.currentEvent == model.EVENT_MENU {
		operation = this.playerOperation
	} else if this.currentEvent == model.EVENT_CONFIRM {
		resourceId, _ := this.GetCurrentLocationResourceId()
		if model.GetResourceType(resourceId) == model.RESOURCE_TYPE_NPC {
			this.npcOperation.SetCurrentNpcResourceId(resourceId)
			operation = this.npcOperation
		} else if model.GetResourceType(resourceId) == model.RESOURCE_TYPE_MAP {
			resourceRealId := model.GetResourceRealId(resourceId)
			basicMap, err := this.resourceManager.GetMap(resourceRealId)
			if err != nil {
				this.logger.Debug(log_content.LogContentNormal(mapOperationLoggerTag, "err:%+v", err))
			}
			operation = NewMapOperation(this.player, basicMap, 0, 1, this.resourceManager, this.playerOperation, false)
		}
	}
	this.nextOperation = operation
	return operation
}

// GetCurrentLocation .
func (this *MapOperation) GetCurrentLocation() (int, int) {
	return this.currentX, this.currentY
}

// SetLocationResource .
func (this *MapOperation) SetLocationResource(x int, y int, resourceId uint32, count uint8) error {
	totalY := len(this.currentMap.MapResources)

	if y <= totalY {
		totalX := len(this.currentMap.MapResources[y])
		if x <= totalX {
			this.currentMap.MapResources[y][x] = model.NewMapResource(resourceId, make([]uint32, 0))
		} else {
			return errors.New(fmt.Sprintf("x y not match, maxX:%d, maxY:%d, x:%d, y:%d", totalX, totalY, x, y))
		}
	} else {
		return errors.New(fmt.Sprintf("x y not match, maxY:%d, x:%d, y:%d", totalY, x, y))
	}
	return nil
}

// GetLocationMapResource .
func (this *MapOperation) GetLocationMapResource(x int, y int) (*model.MapResource, error) {
	totalYLength := len(this.currentMap.MapResources)
	if y > (totalYLength-1) || y < 0 {
		return nil, errors.New(fmt.Sprintf("x y not match, maxY:%d, x:%d, y:%d", totalYLength, x, y))
	} else {
		totalXLength := len(this.currentMap.MapResources[y])
		if x > (totalXLength-1) || x < 0 {
			return nil, errors.New(fmt.Sprintf("x y not match, maxX:%d, maxY:%d, x:%d, y:%d", totalXLength, totalYLength, x, y))
		} else {
			//continue
		}
	}
	mapResource := this.currentMap.MapResources[y][x]
	resourceId := mapResource.ResourceId
	if resourceId == model.EMPTY_LOCATION {
		return nil, errors.New(fmt.Sprintf("It is empty location, x:%d, y:%d", x, y))
	}
	return mapResource, nil
}

// GetLocationResourceId .
func (this *MapOperation) GetLocationResourceId(x int, y int) (uint32, error) {
	mapResource, err := this.GetLocationMapResource(x, y)
	if err != nil {
		return model.EMPTY_LOCATION, err
	}
	return mapResource.ResourceId, nil
}

// GetCurrentLocationMapResource .
func (this *MapOperation) GetCurrentLocationMapResource() (*model.MapResource, error) {
	return this.GetLocationMapResource(this.currentX, this.currentY)
}

// GetCurrentLocationResourceId .
func (this *MapOperation) GetCurrentLocationResourceId() (uint32, error) {
	return this.GetLocationResourceId(this.currentX, this.currentY)
}

// MoveUp .
func (this *MapOperation) MoveUp() (uint32, error) {
	newX := this.currentX
	newY := this.currentY - 1
	resourceId, err := this.GetLocationResourceId(newX, newY)
	if err != nil {
		return model.EMPTY_LOCATION, err
	}
	this.currentX = newX
	this.currentY = newY
	return resourceId, nil
}

// MoveDown .
func (this *MapOperation) MoveDown() (uint32, error) {
	newX := this.currentX
	newY := this.currentY + 1
	resourceId, err := this.GetLocationResourceId(newX, newY)
	if err != nil {
		return model.EMPTY_LOCATION, err
	}
	this.currentX = newX
	this.currentY = newY
	return resourceId, nil
}

// MoveLeft .
func (this *MapOperation) MoveLeft() (uint32, error) {
	newX := this.currentX - 1
	newY := this.currentY
	resourceId, err := this.GetLocationResourceId(newX, newY)
	if err != nil {
		return model.EMPTY_LOCATION, err
	}
	this.currentX = newX
	this.currentY = newY
	return resourceId, nil
}

// MoveRight .
func (this *MapOperation) MoveRight() (uint32, error) {
	newX := this.currentX + 1
	newY := this.currentY
	resourceId, err := this.GetLocationResourceId(newX, newY)
	if err != nil {
		return model.EMPTY_LOCATION, err
	}
	this.currentX = newX
	this.currentY = newY
	return resourceId, nil
}

func (this *MapOperation) ClickConfirm() error {
	mapResource, err := this.GetLocationMapResource(this.currentX, this.currentY)
	if err != nil {
		return errors.New(fmt.Sprintf("location resource can not confirm, current (x:%d, y:%d)", this.currentX, this.currentY))
	}
	resourceId := mapResource.ResourceId
	switch model.GetResourceType(resourceId) {
	case model.RESOURCE_TYPE_MAP:
		break
	case model.RESOURCE_TYPE_ITEM:
		this.logger.Debug(log_content.LogContentNormal(mapOperationLoggerTag, "before add the item to player, resourceId:%x", resourceId))
		if model.GetResourceState(resourceId) == model.RESOURCE_ITEM_STATE_CAN_OPEN {
			if mapResource.InnerResourceIdList != nil {
				for _, innerResourceId := range mapResource.InnerResourceIdList {
					innerResourceType := model.GetResourceType(innerResourceId)
					if innerResourceType == model.RESOURCE_TYPE_ITEM {
						innerResourceReadId := model.GetResourceRealId(innerResourceId)
						this.player.AddItem(innerResourceReadId)
					}
				}
				mapResource.InnerResourceIdList = make([]uint32, 0)
			}
			//add item to player item list
			newResourceId := model.ResourceUnsetState(resourceId, model.RESOURCE_ITEM_STATE_CAN_OPEN)
			err = this.SetLocationResource(this.currentX, this.currentY, newResourceId, 0)
			if err != nil {
				return errors.New(fmt.Sprintf("location resource can not set, current (x:%d, y:%d) err:%v", this.currentX, this.currentY, err))
			}
		} else {
			this.logger.Warning(log_content.LogContentNormal(mapOperationLoggerTag, "can not open, current (x:%d, y:%d), resource id:%x", this.currentX, this.currentY, resourceId))
		}
		break
	case model.RESOURCE_TYPE_NPC:
		return nil
	default:
		return errors.New(fmt.Sprintf("not support item type, current (x:%d, y:%d)", this.currentX, this.currentY))
	}
	return nil
}

func (this *MapOperation) NearbyResourceList() []*base.Tuple[int, int, uint32] {
	fourDirectionXYList := []*base.Pair[int, int]{
		base.NewPair[int, int](this.currentX, this.currentY-1), //up
		base.NewPair[int, int](this.currentX, this.currentY+1), //down
		base.NewPair[int, int](this.currentX-1, this.currentY), //left
		base.NewPair[int, int](this.currentX+1, this.currentY), //right
	}
	resourceIdList := make([]*base.Tuple[int, int, uint32], 0)
	for _, directionXY := range fourDirectionXYList {
		currentX := directionXY.First
		currentY := directionXY.Second
		resourceId, err := this.GetLocationResourceId(currentX, currentY)
		if err != nil {
			resourceIdList = append(resourceIdList, nil)
			continue
		}
		tupleData := base.NewTuple[int, int, uint32](currentX, currentY, resourceId)
		resourceIdList = append(resourceIdList, tupleData)
	}
	return resourceIdList
}

func (this *MapOperation) getMapViewList() []view.Displayable {
	viewList := make([]view.Displayable, 0)
	mapName := this.currentMap.Name
	viewList = append(viewList, view.NewTextView(fmt.Sprintf("-----MAP CONTENT BEGIN-----[%s]", mapName)))

	for y, rowMapResources := range this.currentMap.MapResources {
		var stringBuilder strings.Builder
		for x, mapResource := range rowMapResources {
			resourceId := mapResource.ResourceId
			resource, err := this.resourceManager.GetResource(resourceId)
			resourceState := model.GetResourceState(resourceId)

			currentFlag := constants.STRING_BLANK
			if this.currentX == x && this.currentY == y {
				currentFlag = constants.EMOJI_GRINNING_FACE
			}
			if err != nil {
				stringBuilder.WriteString(fmt.Sprintf("%-24s", fmt.Sprintf("ERR (ID:%d)(%d, %d)", mapResource.ResourceId, x, y)))
			} else {
				stringBuilder.WriteString(fmt.Sprintf("%-24s", fmt.Sprintf("%s%s(%d, %d)[%d]", currentFlag, resource.Name, x, y, resourceState)))
			}
		}

		viewList = append(viewList, view.NewTextView(stringBuilder.String()))
	}
	viewList = append(viewList, view.NewTextView(fmt.Sprintf("-----MAP CONTENT END-----[%s]", mapName)))
	return viewList
}

func (this *MapOperation) getTipsViewList() []view.Displayable {
	viewList := make([]view.Displayable, 0)
	viewList = append(viewList, view.NewTextView("-----CONTROL TIPS CONTENT BEGIN-----"))

	viewList = append(viewList, view.NewTextView(fmt.Sprintf("Current location: (x:%d, y:%d)", this.currentX, this.currentY)))

	// nearby resource list
	nearByResourceList := this.NearbyResourceList()
	defaultEventList := []model.Event{
		model.EVENT_UP,
		model.EVENT_DOWN,
		model.EVENT_LEFT,
		model.EVENT_RIGHT,
	}

	itemTipsViewList := make([]view.Displayable, 0)

	for index, locationResourceTuple := range nearByResourceList {
		if locationResourceTuple == nil {
			continue
		}
		x := locationResourceTuple.First
		y := locationResourceTuple.Second
		resourceId := locationResourceTuple.Third
		resource, err := this.resourceManager.GetResource(resourceId)
		//resourceRealId := this.resourceManager.GetResourceRealId(resourceId)
		//mapThing, err := this.resourceManager.GetMapThing(resourceRealId)
		if err != nil {
			this.logger.Debug(log_content.LogContentNormal(mapOperationLoggerTag, "err:%+v", err))
		}
		event := defaultEventList[index]
		keyCode := model.EVENT_KEY_CODE_MAPPING[event]
		viewList = append(viewList, view.NewButtonView(event, fmt.Sprintf("%s: %s (x:%d, y:%d)", strings.ToUpper(string(keyCode)), resource.Name, x, y)))
	}

	resourceId, err := this.GetLocationResourceId(this.currentX, this.currentY)
	if err == nil {
		resource, err := this.resourceManager.GetResource(resourceId)
		if err == nil {
			if model.ResourceTypeIsItem(resourceId) && model.ResourceItemStateCanOpen(resourceId) {
				confirmKeyCode := model.EVENT_KEY_CODE_MAPPING[model.EVENT_CONFIRM]
				itemTipsViewList = append(itemTipsViewList, view.NewButtonView(model.EVENT_CONFIRM, fmt.Sprintf("%s: %s (x:%d, y:%d) can open", strings.ToUpper(string(confirmKeyCode)), resource.Name, this.currentX, this.currentY)))
			} else if model.ResourceTypeIsNpc(resourceId) && model.ResourceNpcStateCanTalk(resourceId) {
				confirmKeyCode := model.EVENT_KEY_CODE_MAPPING[model.EVENT_CONFIRM]
				itemTipsViewList = append(itemTipsViewList, view.NewButtonView(model.EVENT_CONFIRM, fmt.Sprintf("%s: %s (x:%d, y:%d) can talk", strings.ToUpper(string(confirmKeyCode)), resource.Name, this.currentX, this.currentY)))
			} else if model.ResourceTypeIsMap(resourceId) && model.ResourceMapStateEnable(resourceId) {
				confirmKeyCode := model.EVENT_KEY_CODE_MAPPING[model.EVENT_CONFIRM]
				itemTipsViewList = append(itemTipsViewList, view.NewButtonView(model.EVENT_CONFIRM, fmt.Sprintf("%s: %s (x:%d, y:%d) can in", strings.ToUpper(string(confirmKeyCode)), resource.Name, this.currentX, this.currentY)))
			}
		}
	}

	itemTipsViewList = append(itemTipsViewList, view.NewButtonView(model.EVENT_MENU, fmt.Sprintf("%s: Menu", strings.ToUpper(string(model.EVENT_KEY_CODE_MAPPING[model.EVENT_MENU])))))

	viewList = append(viewList, itemTipsViewList...)

	viewList = append(viewList, view.NewTextView("-----CONTROL TIPS CONTENT END-----"))
	return viewList
}

func (this *MapOperation) getView() view.Displayable {
	viewList := make([]view.Displayable, 0)

	mapViewList := this.getMapViewList()
	tipsViewList := this.getTipsViewList()
	viewList = append(viewList, mapViewList...)
	viewList = append(viewList, tipsViewList...)

	return view.NewViewGroup(viewList...)
}

func (this *MapOperation) operateNextOperation(event model.Event) view.Displayable {
	var nextOperationDisplayable view.Displayable = nil
	if this.nextOperation != nil {
		nextOperationDisplayable = this.nextOperation.Operate(event)
	}
	if nextOperationDisplayable == nil {
		this.nextOperation = nil
		return this.getView()
	} else {
		cancelKeyCode := model.EVENT_KEY_CODE_MAPPING[model.EVENT_CANCEL]
		return view.NewViewGroup(
			nextOperationDisplayable,
			view.NewTextView("-----OTHER CONTROL TIPS CONTENT BEGIN-----"),
			view.NewButtonView(model.EVENT_CANCEL, fmt.Sprintf("%s: Return to previous", strings.ToUpper(string(cancelKeyCode)))),
			view.NewTextView("-----OTHER CONTROL TIPS CONTENT END-----"),
		)
	}
}

func (this *MapOperation) LoadSavedData(dataMap map[string]any) {
	inputPlayer, ok := dataMap[MAP_OPERATION_DATA_KEY_PLAYER]
	//player := &model.Player{}
	if ok {
		inputPlayerJson, err := json.Marshal(inputPlayer)
		if err != nil {
			this.logger.Error(log_content.LogContentNormal(mapOperationLoggerTag, "json.Marshal err"), err)
		}
		err = json.Unmarshal(inputPlayerJson, this.player)
		if err != nil {
			this.logger.Error(log_content.LogContentNormal(mapOperationLoggerTag, "json.Unmarshal err"), err)
		}
		//this.player = player

	}

	inputCurrentMap, ok := dataMap[MAP_OPERATION_DATA_KEY_CURRENT_MAP]
	//currentMap := &model.Map{}
	if ok {
		inputCurrentMapJson, err := json.Marshal(inputCurrentMap)
		if err != nil {
			this.logger.Error(log_content.LogContentNormal(mapOperationLoggerTag, "json.Marshal err"), err)
		}
		err = json.Unmarshal(inputCurrentMapJson, this.currentMap)
		if err != nil {
			this.logger.Error(log_content.LogContentNormal(mapOperationLoggerTag, "json.Unmarshal err"), err)
		}
		//this.currentMap = currentMap
	}

	currentX, ok := dataMap[MAP_OPERATION_DATA_KEY_CURRENT_X]
	if ok {
		this.currentX = int(currentX.(float64))
	}

	currentY, ok := dataMap[MAP_OPERATION_DATA_KEY_CURRENT_Y]
	if ok {
		this.currentY = int(currentY.(float64))
	}
}

func (this *MapOperation) GetNeedToSavedData() map[string]any {
	dataMap := make(map[string]any, 0)
	dataMap[MAP_OPERATION_DATA_KEY_PLAYER] = this.player
	dataMap[MAP_OPERATION_DATA_KEY_CURRENT_MAP] = this.currentMap
	dataMap[MAP_OPERATION_DATA_KEY_CURRENT_X] = this.currentX
	dataMap[MAP_OPERATION_DATA_KEY_CURRENT_Y] = this.currentY
	return dataMap
}
