package test

import (
	"fmt"
	"model"
	"testing"
)

func TestMap(t *testing.T) {
	var resourceId uint32 = 0x3010001
	var state uint32 = model.RESOURCE_ITEM_STATE_CAN_OPEN
	fmt.Printf("0x%x", resourceId)
	fmt.Println()
	fmt.Printf("%b", state)
	fmt.Println()
	fmt.Printf("0x%x", model.ResourceUnsetState(resourceId, state))
	return
}
