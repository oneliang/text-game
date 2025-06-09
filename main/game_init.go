package main

import (
	"fmt"
	"github.com/oneliang/util-golang/common"
	"model"
)

func ReadSavedData() *model.SaveData {
	saveData := &model.SaveData{}

	err := common.LoadJsonToObject("../data/save/saved_0.json", saveData)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}
	if DEBUG {

	}
	return saveData
}

func ReadMapData() []*model.Map {
	basicMapList := make([]*model.Map, 0)

	basicMap := &model.Map{}
	err := common.LoadJsonToObject("../data/map/map_1.json", basicMap)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}
	if DEBUG {
		for _, rowMapResource := range basicMap.MapResources {
			for _, mapResource := range rowMapResource {
				fmt.Printf("%-24s", fmt.Sprintf("0x%x", mapResource.ResourceId))
			}
			fmt.Println()
		}
	}
	basicMapList = append(basicMapList, basicMap)
	return basicMapList
}

func ReadMapThingData() []*model.Resource {
	mapThingList := make([]*model.Resource, 0)
	err := common.LoadJsonToObject("../data/item/map_thing.json", &mapThingList)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}
	if DEBUG {
		fmt.Println(fmt.Sprintf("map thing list:%+v", mapThingList))
	}
	return mapThingList
}

func ReadItemData() []*model.Resource {
	itemList := make([]*model.Resource, 0)
	err := common.LoadJsonToObject("../data/item/item.json", &itemList)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}
	if DEBUG {
		fmt.Println(fmt.Sprintf("items list:%+v", itemList))
	}
	return itemList
}

func ReadNpcData() []*model.Resource {
	npcList := make([]*model.Resource, 0)
	err := common.LoadJsonToObject("../data/item/npc.json", &npcList)
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}
	if DEBUG {
		fmt.Println(fmt.Sprintf("npc list:%+v", npcList))
	}
	return npcList
}
