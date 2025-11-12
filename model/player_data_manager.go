package model

import (
	"encoding/json"
	"github.com/oneliang/util-golang/base"
	"github.com/oneliang/util-golang/logging"
	"log_content"
	"time"
)

const (
	DATA_KEY_PLAYER    = "PLAYER"
	DATA_KEY_MAPS      = "MAPS"
	DATA_KEY_CURRENT_X = "CURRENT_X"
	DATA_KEY_CURRENT_Y = "CURRENT_Y"
)
const playerDataManagerLoggerTag = "playerDataManagerLoggerTag"

type PlayerDataManager struct {
	resourceManager *ResourceManager
	player          *Player
	mapMap          map[uint32]*Map
	logger          logging.Logger
}

func NewPlayerDataManager(resourceManager *ResourceManager) *PlayerDataManager {
	return &PlayerDataManager{
		resourceManager: resourceManager,
		player:          NewPlayer(int32(time.Now().Unix()), make([]*base.KeyValue[uint32, uint8], 0)),
		mapMap:          make(map[uint32]*Map),
		logger:          logging.LoggerManager.GetLogger(playerDataManagerLoggerTag),
	}
}

// GetPlayer .
func (this *PlayerDataManager) GetPlayer() *Player {
	if this.player == nil {
		this.player = NewPlayer(int32(time.Now().Unix()), make([]*base.KeyValue[uint32, uint8], 0))
	}
	return this.player
}

// GetMap .
func (this *PlayerDataManager) GetMap(resourceId uint32) *Map {
	playerMap, exist := this.mapMap[resourceId]
	if exist {
		return playerMap
	}
	playerMap, err := this.resourceManager.GetMapWithResourceId(resourceId)
	if err != nil {
		this.logger.Debug(log_content.LogContentNormal(playerDataManagerLoggerTag, "err:%+v", err))
	}
	this.mapMap[resourceId] = playerMap
	return playerMap
}

// LoadSavedData .
func (this *PlayerDataManager) LoadSavedData(dataMap map[string]any) {
	inputPlayer, ok := dataMap[DATA_KEY_PLAYER]
	if ok {
		inputPlayerJson, err := json.Marshal(inputPlayer)
		if err != nil {
			this.logger.Error(log_content.LogContentNormal(playerDataManagerLoggerTag, "json.Marshal err"), err)
		}
		err = json.Unmarshal(inputPlayerJson, this.player)
		if err != nil {
			this.logger.Error(log_content.LogContentNormal(playerDataManagerLoggerTag, "json.Unmarshal err"), err)
		}
		//this.player = player

	}

	maps, ok := dataMap[DATA_KEY_MAPS]
	if ok {
		inputCurrentMapJson, err := json.Marshal(maps)
		if err != nil {
			this.logger.Error(log_content.LogContentNormal(playerDataManagerLoggerTag, "json.Marshal err"), err)
		}
		err = json.Unmarshal(inputCurrentMapJson, &this.mapMap)
		if this.mapMap == nil {
			this.mapMap = make(map[uint32]*Map)
		}
		if err != nil {
			this.logger.Error(log_content.LogContentNormal(playerDataManagerLoggerTag, "json.Unmarshal err"), err)
		}
		//this.currentMap = currentMap
	}
	//
	//currentX, ok := dataMap[MAP_OPERATION_DATA_KEY_CURRENT_X]
	//if ok {
	//	this.currentX = int(currentX.(float64))
	//}
	//
	//currentY, ok := dataMap[MAP_OPERATION_DATA_KEY_CURRENT_Y]
	//if ok {
	//	this.currentY = int(currentY.(float64))
	//}
}

func (this *PlayerDataManager) GetNeedToSavedData() map[string]any {
	dataMap := make(map[string]any, 0)
	dataMap[DATA_KEY_PLAYER] = this.player
	dataMap[DATA_KEY_MAPS] = this.mapMap
	return dataMap
}
