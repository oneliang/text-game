package test

import (
	"fmt"
	"github.com/oneliang/util-golang/common"
	"model"
	"testing"
)

func readPlayerData() *model.Player {
	return nil
}

func readMapData() []*model.Map {
	basicMapList := make([]*model.Map, 0)

	basicMap := &model.Map{}
	err := common.LoadJsonToObject("../data/map/map_1.json", basicMap)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}
	fmt.Println(fmt.Sprintf("basic map:%+v", basicMap))
	basicMapList = append(basicMapList, basicMap)
	return basicMapList
}

func readMapThingData() []*model.Resource {
	mapThingList := make([]*model.Resource, 0)
	err := common.LoadJsonToObject("../data/item/map_thing.json", &mapThingList)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}
	fmt.Println(fmt.Sprintf("map thing list:%+v", mapThingList))
	return mapThingList
}

func TestGame(t *testing.T) {
	//basicMapList := readMapData()
	//mapThingList := readMapThingData()
	//mapThingMap := common.ListToMap[*model.Resource, uint32](mapThingList, func(index int, item *model.Resource) uint32 {
	//	return item.RealId
	//})
	//basicMapMap := common.ListToMap[*model.basicMap, uint32](basicMapList, func(index int, item *model.basicMap) uint32 {
	//	return item.Id
	//})
	//
	//playerMap := model.NewPlayerMap(basicMapList[0], 0, 1)
	//resourceManager := model.NewResourceManager(&basicMapMap, &mapThingMap)
	//game := main.NewGame(nil, playerMap, resourceManager)
	//
	//game.Start()
	//game.PostEvent(model.EVENT_RIGHT)
	//game.PostEvent(model.EVENT_RIGHT)
	////game.PostEvent(model.EVENT_RIGHT)
	////game.PostEvent(model.EVENT_RIGHT)
	//time.Sleep(3 * time.Second)
}
