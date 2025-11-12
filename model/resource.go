package model

const (
	RESOURCE_TYPE_BIT_SIZE  uint32 = 24
	RESOURCE_TYPE_EMPTY     uint32 = 0x00 // << RESOURCE_TYPE_BIT_SIZE
	RESOURCE_TYPE_MAP       uint32 = 0x01 // << RESOURCE_TYPE_BIT_SIZE
	RESOURCE_TYPE_MAP_THING uint32 = 0x02 // << RESOURCE_TYPE_BIT_SIZE
	RESOURCE_TYPE_ITEM      uint32 = 0x03 // << RESOURCE_TYPE_BIT_SIZE
	RESOURCE_TYPE_NPC       uint32 = 0x04 // << RESOURCE_TYPE_BIT_SIZE

	RESOURCE_STATE_BIT_SIZE uint32 = 16
	RESOURCE_STATE_EMPTY    uint32 = 0x00 // << RESOURCE_STATE_BIT_SIZE

	RESOURCE_TYPE_MASK    uint32 = 0xFF_00_0000
	RESOURCE_STATE_MASK   uint32 = 0x00_FF_0000
	RESOURCE_REAL_ID_MASK uint32 = 0x00_00_FFFF

	// resource item state use bit set to store item state
	RESOURCE_ITEM_STATE_NONE     uint32 = 0b0000_0000 // << 0 << RESOURCE_STATE_BIT_SIZE
	RESOURCE_ITEM_STATE_CAN_OPEN uint32 = 0b0000_0001 // << 0 << RESOURCE_STATE_BIT_SIZE
	RESOURCE_NPC_STATE_CAN_TALK  uint32 = 0b0000_0001 // << 0 << RESOURCE_STATE_BIT_SIZE
	RESOURCE_MAP_STATE_ENABLE    uint32 = 0b0000_0001 // << 0 << RESOURCE_STATE_BIT_SIZE
)

type Resource struct {
	Name        string            `json:"name"`
	Type        uint32            `json:"type"`
	State       uint32            `json:"state"`
	RealId      uint32            `json:"realId"`
	PropertyMap map[string]string `json:"propertyMap"`
}

func NewResource(name string, resourceType uint32, state uint32, readId uint32) *Resource {
	return &Resource{
		Name:        name,
		Type:        resourceType,
		State:       state,
		RealId:      readId,
		PropertyMap: make(map[string]string),
	}
}

func GenResourceId(resourceType uint32, state uint32, realId uint32) uint32 {
	return resourceType<<RESOURCE_TYPE_BIT_SIZE | state<<RESOURCE_STATE_BIT_SIZE | realId
}
func GetResourceType(resourceId uint32) uint32 {
	resourceType := resourceId & RESOURCE_TYPE_MASK >> RESOURCE_TYPE_BIT_SIZE
	return resourceType
}

func GetResourceState(resourceId uint32) uint32 {
	resourceState := resourceId & RESOURCE_STATE_MASK >> RESOURCE_STATE_BIT_SIZE
	return resourceState
}

func GetResourceRealId(resourceId uint32) uint32 {
	resourceRealId := resourceId & RESOURCE_REAL_ID_MASK
	return resourceRealId
}

func ResourceTypeIsMap(resourceId uint32) bool {
	return GetResourceType(resourceId) == RESOURCE_TYPE_MAP
}

func ResourceTypeIsMapThing(resourceId uint32) bool {
	return GetResourceType(resourceId) == RESOURCE_TYPE_MAP_THING
}

func ResourceTypeIsItem(resourceId uint32) bool {
	return GetResourceType(resourceId) == RESOURCE_TYPE_ITEM
}

func ResourceTypeIsNpc(resourceId uint32) bool {
	return GetResourceType(resourceId) == RESOURCE_TYPE_NPC
}

func ResourceMapStateEnable(resourceId uint32) bool {
	return GetResourceState(resourceId) == RESOURCE_MAP_STATE_ENABLE
}

func ResourceItemStateCanOpen(resourceId uint32) bool {
	return GetResourceState(resourceId) == RESOURCE_ITEM_STATE_CAN_OPEN
}

func ResourceNpcStateCanTalk(resourceId uint32) bool {
	return GetResourceState(resourceId) == RESOURCE_NPC_STATE_CAN_TALK
}

func ResourceUnsetState(resourceId uint32, resourceState uint32) uint32 {
	return ^(resourceState << RESOURCE_STATE_BIT_SIZE) & resourceId
}
