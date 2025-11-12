package main

import (
	"fmt"
	"github.com/oneliang/util-golang/common"
	"github.com/oneliang/util-golang/constants"
	"github.com/oneliang/util-golang/logging"
	"log_content"
	"model"
)

const gameInitLoggerTag = "GameInit"

var logger = logging.LoggerManager.GetLogger(gameInitLoggerTag)

func ReadSavedData() *model.SaveData {
	saveData := &model.SaveData{}

	err := common.LoadJsonToObject("../data/save/saved_0.json", saveData)
	if err != nil {
		logger.Error(log_content.LogContentNormal(gameInitLoggerTag, constants.STRING_BLANK), err)
	}
	if DEBUG {

	}
	return saveData
}

func ReadMapData() []*model.Map {
	basicMapList := make([]*model.Map, 0)
	basicMap1 := readOneMapData("../data/map/map_1.json")
	basicMap2 := readOneMapData("../data/map/map_2.json")
	basicMapList = append(basicMapList, basicMap1)
	basicMapList = append(basicMapList, basicMap2)
	return basicMapList
}

func readOneMapData(filepath string) *model.Map {
	basicMap := &model.Map{}
	err := common.LoadJsonToObject(filepath, basicMap)
	if err != nil {
		logger.Error(log_content.LogContentNormal(gameInitLoggerTag, constants.STRING_BLANK), err)
	}
	if DEBUG {
		logger.Debug(log_content.LogContentNormal(gameInitLoggerTag, "map:%-24s", fmt.Sprintf("0x%x", model.GenResourceId(basicMap.Resource.Type, basicMap.Resource.State, basicMap.Resource.RealId))))
		for _, rowMapResource := range basicMap.MapResources {
			for _, mapResource := range rowMapResource {
				logger.Debug(log_content.LogContentNormal(gameInitLoggerTag, "map resource:%-24s", fmt.Sprintf("0x%x", mapResource.ResourceId)))
			}
			fmt.Println()
		}
	}
	return basicMap
}

func ReadMapThingData() []*model.Resource {
	mapThingList := make([]*model.Resource, 0)
	err := common.LoadJsonToObject("../data/item/map_thing.json", &mapThingList)
	if err != nil {
		logger.Error(log_content.LogContentNormal(gameInitLoggerTag, constants.STRING_BLANK), err)
	}
	if DEBUG {
		logger.Debug(log_content.LogContentNormal(gameInitLoggerTag, "map thing list:%+v", mapThingList))
	}
	return mapThingList
}

func ReadItemData() []*model.Resource {
	itemList := make([]*model.Resource, 0)
	err := common.LoadJsonToObject("../data/item/item.json", &itemList)
	if err != nil {
		logger.Error(log_content.LogContentNormal(gameInitLoggerTag, constants.STRING_BLANK), err)
	}
	if DEBUG {
		logger.Debug(log_content.LogContentNormal(gameInitLoggerTag, "items list:%+v", itemList))
	}
	return itemList
}

func ReadNpcData() []*model.Resource {
	npcList := make([]*model.Resource, 0)
	err := common.LoadJsonToObject("../data/item/npc.json", &npcList)
	if err != nil {
		logger.Error(log_content.LogContentNormal(gameInitLoggerTag, constants.STRING_BLANK), err)
	}
	if DEBUG {
		logger.Debug(log_content.LogContentNormal(gameInitLoggerTag, "npc list:%+v", npcList))
	}
	return npcList
}
