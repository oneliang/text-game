package main

import (
	"fmt"
	gotty "github.com/mattn/go-tty"
	"github.com/oneliang/util-golang/base"
	"github.com/oneliang/util-golang/common"
	"log"
	"logic"
	"model"
)

const DEBUG = true

func main() {
	//fmt.Println(model.GenResourceId(4, 0, 1))
	//fmt.Println(model.GenResourceId(4, 0, 2))
	//fmt.Println(model.GenResourceId(4, 0, 3))
	////fmt.Println(model.GenResourceId(3, 0, 3))
	//return
	//read all resource
	saveData := ReadSavedData()
	basicMapList := ReadMapData()
	mapThingList := ReadMapThingData()
	itemList := ReadItemData()
	npcList := ReadNpcData()

	basicMapMap := common.ListToMap[*model.Map, uint32](basicMapList, func(index int, item *model.Map) uint32 {
		return item.Id
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
	player := model.NewPlayer(make([]*base.KeyValue[uint32, uint8], 0))

	//var playerMap *model.PlayerMap = nil
	//if saveData.PlayerMap != nil {
	//	playerMap = saveData.PlayerMap
	//} else {
	//	playerMap = model.NewPlayerMap(player, basicMapList[0], 0, 1)
	//}

	resourceManager := model.NewResourceManager(&basicMapMap, &mapThingMap, &itemMap, &npcMap)
	playerOperation := logic.NewPlayerOperation(player, resourceManager)
	mapOperation := logic.NewMapOperation(player, basicMapList[0], 0, 1, resourceManager, playerOperation)
	// new game
	game := NewGame(resourceManager)
	game.SetRootOperation(mapOperation)

	mapOperation.LoadSavedData(saveData.DataMap)

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
