package model

import (
	"errors"
	"github.com/oneliang/util-golang/common"
)

const (
	EMPTY_LOCATION = RESOURCE_TYPE_EMPTY<<RESOURCE_TYPE_BIT_SIZE | 0x00_00_0000
)

type Map struct {
	Id           uint32           `json:"id"`
	MapResources [][]*MapResource `json:"mapResources"`
}

type MapResource struct {
	ResourceId          uint32   `json:"resourceId"`
	InnerResourceIdList []uint32 `json:"innerResourceIdList"`
}

// NewMapResource .
func NewMapResource(resourceId uint32, innerResourceIdList []uint32) *MapResource {
	if innerResourceIdList == nil {
		innerResourceIdList = make([]uint32, 0)
	}
	return &MapResource{
		ResourceId:          resourceId,
		InnerResourceIdList: innerResourceIdList,
	}
}

func NewMapWithSize(id uint32, maxXSize uint32, maxYSize uint32) *Map {
	mapResources := make([][]*MapResource, maxYSize)
	for y, _ := range mapResources {
		mapResources[y] = make([]*MapResource, maxXSize)

		for x := 0; x < int(maxXSize); x++ {
			mapResources[y][x] = NewMapResource(EMPTY_LOCATION, make([]uint32, 0))
		}
	}
	return &Map{
		Id:           id,
		MapResources: mapResources,
	}
}

func NewMapWithResourceIds(id uint32, mapResources [][]*MapResource) (*Map, error) {
	if mapResources == nil {
		return nil, errors.New("")
	}
	newMapResources := common.ListToNewList[[]*MapResource, []*MapResource](mapResources, func(index int, rowMapResources []*MapResource) []*MapResource {
		return common.ListToNewList[*MapResource, *MapResource](rowMapResources, func(index int, mapResource *MapResource) *MapResource {
			//clone MapResource
			return NewMapResource(mapResource.ResourceId, make([]uint32, 0))
		})
	})

	return &Map{
		Id:           id,
		MapResources: newMapResources,
	}, nil
}
