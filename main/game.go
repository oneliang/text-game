package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oneliang/util-golang/constants"
	"github.com/oneliang/util-golang/file"
	"github.com/oneliang/util-golang/logging"
	"log_content"
	"logic"
	"model"
	"os"
	"path/filepath"
	"view"
)

const loggerTag = "Game"

type Game struct {
	rootOperation   logic.Operation
	resourceManager *model.ResourceManager
	eventExecutor   *model.EventExecutor
	logger          logging.Logger
}

func NewGame(resourceManager *model.ResourceManager) *Game {
	logger := logging.LoggerManager.GetLogger(loggerTag)
	game := &Game{
		resourceManager: resourceManager,
		logger:          logger,
	}
	eventExecutor, err := model.NewEventExecutor(game, game.realStopCallback)
	if err != nil {
		logger.Error(log_content.LogContentNormal(loggerTag, constants.STRING_BLANK), err)
	}
	game.eventExecutor = eventExecutor
	return game
}

func (this *Game) SetRootOperation(rootOperation logic.Operation) {
	this.rootOperation = rootOperation
}

func (this *Game) Start() {

}

func (this *Game) PostEvent(event model.Event) {
	this.eventExecutor.PostEvent(event)
}

func (this *Game) Process(event model.Event) error {
	this.logger.Debug(log_content.LogContentNormal(loggerTag, "event process, event:%v", event))
	if this.rootOperation == nil {
		return errors.New("need to invoke SetOperation first, the operation is nil")
	}
	displayable := this.rootOperation.Operate(event)

	this.printView(displayable)

	return nil
}

func (this *Game) realStopCallback() {

}

func (this *Game) printView(displayable view.Displayable) {
	fmt.Println("----------GAME CONTENT BEGIN----------")
	if displayable != nil {
		fmt.Print(displayable.Display())
	}
	fmt.Println("----------GAME CONTENT END----------")
}

func (this *Game) SaveGame(outputFullFilename string) error {
	if !file.Exists(outputFullFilename) {
		outputDirectory := filepath.Dir(outputFullFilename)
		err := os.MkdirAll(outputDirectory, 0755)
		if err != nil {
			return err
		}
	}
	outputFile, err := os.Create(outputFullFilename)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(outputFile)

	dataMap := this.rootOperation.GetNeedToSavedData()

	saveData := &model.SaveData{
		DataMap: dataMap,
	}
	dataBytes, err := json.Marshal(saveData)
	if err != nil {
		return err
	}
	_, err = writer.Write(dataBytes)
	if err != nil {
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
