package model

import (
	"github.com/oneliang/util-golang/base"
	"github.com/oneliang/util-golang/common"
)

type Player struct {
	//[resourceId,count]
	ItemList     []*base.KeyValue[uint32, uint8] `json:"itemList"`
	ItemIndexMap map[uint32]uint32               `json:"itemIndexMap"`
}

func NewPlayer(itemList []*base.KeyValue[uint32, uint8]) *Player {
	itemIndexMap := common.ListToNewMap[*base.KeyValue[uint32, uint8], uint32, uint32](itemList, func(index int, item *base.KeyValue[uint32, uint8]) (uint32, uint32) {
		return item.Key, uint32(index)
	})
	return &Player{
		ItemList:     itemList,
		ItemIndexMap: itemIndexMap,
	}
}

// AddItem . //todo need to use operation lock when concurrent
func (this *Player) AddItem(resourceRealId uint32) {
	index, ok := this.ItemIndexMap[resourceRealId]
	var item *base.KeyValue[uint32, uint8] = nil
	if !ok {
		item = base.NewKeyValue[uint32, uint8](resourceRealId, 0)
		this.ItemList = append(this.ItemList, item)
		this.ItemIndexMap[resourceRealId] = uint32(len(this.ItemList) - 1)
	} else {
		item = this.ItemList[index]
	}
	item.Value = item.Value + 1
}

// AddItemList .
func (this *Player) AddItemList(resourceRealIdList []uint32) {
	for _, resourceRealId := range resourceRealIdList {
		this.AddItem(resourceRealId)
	}
}

// DeleteItem . //todo need to use operation lock when concurrent
func (this *Player) DeleteItem(resourceRealId uint32) {
	index, ok := this.ItemIndexMap[resourceRealId]
	if ok {
		item := this.ItemList[index]
		item.Value = item.Value - 1
		if item.Value == 0 {
			delete(this.ItemIndexMap, resourceRealId)
			this.ItemList = append(this.ItemList[:index], this.ItemList[index+1:]...)
		}
	}

}
