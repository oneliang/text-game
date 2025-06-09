package model

import (
	"errors"
	"fmt"
)

type ResourceManager struct {
	mapMap          *map[uint32]*Map
	mapThingMap     *map[uint32]*Resource
	itemMap         *map[uint32]*Resource
	npcMap          *map[uint32]*Resource
	resourceItemMap *map[uint32]*map[uint32]*Resource //by resource type

}

func NewResourceManager(
	mapMap *map[uint32]*Map,
	mapThingMap *map[uint32]*Resource,
	itemMap *map[uint32]*Resource,
	npcMap *map[uint32]*Resource,
) *ResourceManager {
	resourceItemMap := make(map[uint32]*map[uint32]*Resource)
	resourceItemMap[RESOURCE_TYPE_MAP_THING] = mapThingMap
	resourceItemMap[RESOURCE_TYPE_ITEM] = itemMap
	resourceItemMap[RESOURCE_TYPE_NPC] = npcMap

	resourceItemMapPointer := &resourceItemMap
	return &ResourceManager{
		mapMap:          mapMap,
		mapThingMap:     mapThingMap,
		itemMap:         itemMap,
		npcMap:          npcMap,
		resourceItemMap: resourceItemMapPointer,
	}
}

func (this *ResourceManager) GetMapThing(resourceRealId uint32) (*Resource, error) {
	mapThingMap := *this.mapThingMap
	resourceItem, ok := mapThingMap[resourceRealId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("map thing is not exist, resource real id:%d", resourceRealId))
	}
	return resourceItem, nil
}

func (this *ResourceManager) GetItem(resourceRealId uint32) (*Resource, error) {
	itemMap := *this.itemMap
	resourceItem, ok := itemMap[resourceRealId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("item is not exist, resource real id:%d", resourceRealId))
	}
	return resourceItem, nil
}

func (this *ResourceManager) GetNpc(resourceRealId uint32) (*Resource, error) {
	npcMap := *this.npcMap
	resourceItem, ok := npcMap[resourceRealId]
	if !ok {
		return nil, errors.New(fmt.Sprintf("npc is not exist, resource real id:%d", resourceRealId))
	}
	return resourceItem, nil
}

func (this *ResourceManager) GetResource(resourceId uint32) (*Resource, error) {
	resourceType := GetResourceType(resourceId)
	resourceRealId := GetResourceRealId(resourceId)
	switch resourceType {
	case RESOURCE_TYPE_MAP_THING:
		mapThingMap := *this.mapThingMap
		resourceItem, ok := mapThingMap[resourceRealId]
		if !ok {
			return nil, errors.New(fmt.Sprintf("map thing is not exist, resource real id:%d", resourceRealId))
		}
		return resourceItem, nil
	case RESOURCE_TYPE_ITEM:
		itemMap := *this.itemMap
		resourceItem, ok := itemMap[resourceRealId]
		if !ok {
			return nil, errors.New(fmt.Sprintf("item is not exist, resource real id:%d", resourceRealId))
		}
		return resourceItem, nil
	case RESOURCE_TYPE_NPC:
		npcMap := *this.npcMap
		resourceItem, ok := npcMap[resourceRealId]
		if !ok {
			return nil, errors.New(fmt.Sprintf("npc is not exist, resource real id:%d", resourceRealId))
		}
		return resourceItem, nil
	}

	return nil, errors.New(fmt.Sprintf("resource is not exist, resource id:%d", resourceId))
}
