package model

type SaveData struct {
	DataMap map[string]any `json:"dataMap"`
}

func NewSaveData(
	dataMap map[string]any,
) *SaveData {
	return &SaveData{
		DataMap: dataMap,
	}
}
