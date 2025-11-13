package main

import (
	"fmt"
	gotty "github.com/mattn/go-tty"
	"github.com/oneliang/util-golang/common"
	"github.com/oneliang/util-golang/logging"
	"log"
	"logic"
	"model"
)

const DEBUG = true

func main() {
	//fmt.Println(fmt.Sprintf("%X", 67108865))
	//fmt.Println(model.GenResourceId(1, 1, 1))
	//fmt.Println(model.GenResourceId(1, 1, 2))
	//return
	//fmt.Println(model.GenResourceId(4, 0, 2))
	//fmt.Println(model.GenResourceId(4, 0, 3))
	////fmt.Println(model.GenResourceId(3, 0, 3))
	//return
	//read all resource
	logging.LoggerManager.RegisterLoggerByPattern("*", logging.DEFAULT_LOGGER)

	basicMapList := ReadMapData()
	mapThingList := ReadMapThingData()
	itemList := ReadItemData()
	npcList := ReadNpcData()

	basicMapMap := common.ListToMap[*model.Map, uint32](basicMapList, func(index int, item *model.Map) uint32 {
		return item.RealId
	})
	mapThingMap := common.ListToMap[*model.Resource, uint32](mapThingList, func(index int, item *model.Resource) uint32 {
		return item.RealId
	})
	itemMap := common.ListToMap[*model.Resource, uint32](itemList, func(index int, item *model.Resource) uint32 {
		return item.RealId
	})
	npcMap := common.ListToMap[*model.Resource, uint32](npcList, func(index int, item *model.Resource) uint32 {
		return item.RealId
	})
	fmt.Println(fmt.Sprintf("item map:%+v", itemMap))
	fmt.Println(fmt.Sprintf("npc map:%+v", npcMap))
	// public
	resourceManager := model.NewResourceManager(&basicMapMap, &mapThingMap, &itemMap, &npcMap)
	saveData := ReadSavedData()
	playerDataManager := model.NewPlayerDataManager(resourceManager)
	playerDataManager.LoadSavedData(saveData.DataMap)
	// separate by player
	systemOperation := logic.NewSystemOperation(resourceManager, playerDataManager)
	// new game
	game := NewGame(playerDataManager, resourceManager)
	game.SetRootOperation(systemOperation)

	game.Start()
	//game.PostEvent(model.EVENT_RIGHT)

	// keyboard event
	tty, err := gotty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = tty.Close() }()

	game.PostEvent(model.EVENT_NONE)

	for {
		result, err := tty.ReadRune()
		if err != nil {
			fmt.Println(fmt.Sprintf("%v", err))
		}
		keyCode := byte(result)
		keyEvent, exist := model.KEY_CODE_EVENT_MAPPING[keyCode]
		if exist {
			game.PostEvent(keyEvent)
		}
		err = game.SaveGame("../data/save/saved_0.json")
		if err != nil {
			fmt.Println(fmt.Sprintf("%v", err))
		}
	}

}
