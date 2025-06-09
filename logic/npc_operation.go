package logic

import (
	"fmt"
	"model"
	"view"
)

type NpcOperation struct {
	currentResourceId uint32
	resourceManager   *model.ResourceManager
}

// NewNpcOperation .
func NewNpcOperation(
	resourceManager *model.ResourceManager) *NpcOperation {
	return &NpcOperation{
		resourceManager: resourceManager,
	}
}

func (this *NpcOperation) SetCurrentNpcResourceId(resourceId uint32) {
	this.currentResourceId = resourceId
}

// Operate .
func (this *NpcOperation) Operate(event model.Event) view.Displayable {
	var err error = nil
	switch event {
	case model.EVENT_CONFIRM:
		//err = this.ClickConfirm()
		break
	case model.EVENT_CANCEL:
		return nil
	case model.EVENT_NONE:
		break
	default:
		fmt.Println(fmt.Sprintf("not support event, %d", event))
	}
	if err != nil {
		fmt.Println(fmt.Sprintf("%v", err))
	}
	return this.getView()
}

// GetNextOperation .
func (this *NpcOperation) GetNextOperation() Operation {
	return nil
}

func (this *NpcOperation) getView() view.Displayable {
	viewList := make([]view.Displayable, 0)
	resourceId := this.currentResourceId

	if model.GetResourceType(resourceId) != model.RESOURCE_TYPE_NPC {
		return view.NewViewGroup()
	}

	viewList = append(viewList, view.NewTextView("-----NPC CONTENT BEGIN-----"))

	resourceRealId := model.GetResourceRealId(resourceId)
	resource, err := this.resourceManager.GetNpc(resourceRealId)
	if err != nil {
		viewList = append(viewList, view.NewTextView(fmt.Sprintf("%-24s", fmt.Sprintf("ERR:%v", err))))
	}

	talkContent, ok := resource.PropertyMap[model.PROPERTY_KEY_TALK]
	if ok {
		viewList = append(viewList, view.NewTextView(talkContent))
	}

	viewList = append(viewList, view.NewTextView("-----NPC CONTENT END-----"))
	return view.NewViewGroupWithChildList(
		viewList,
	)
}

func (this *NpcOperation) LoadSavedData(dataMap map[string]any) {

}

func (this *NpcOperation) GetNeedToSavedData() map[string]any {
	return nil
}
