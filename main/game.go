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

const gameLoggerTag = "Game"

type Game struct {
	rootOperation     logic.Operation
	playerDataManager *model.PlayerDataManager
	resourceManager   *model.ResourceManager
	eventExecutor     *model.EventExecutor
	operationManager  *logic.OperationManager
	logger            logging.Logger
}

func NewGame(resourceManager *model.ResourceManager, playerDataManager *model.PlayerDataManager) *Game {
	game := &Game{
		playerDataManager: playerDataManager,
		resourceManager:   resourceManager,
		operationManager:  logic.NewOperationManager(),
		logger:            logging.LoggerManager.GetLogger(gameLoggerTag),
	}
	eventExecutor, err := model.NewEventExecutor(game, game.realStopCallback)
	if err != nil {
		game.logger.Error(log_content.LogContentNormal(gameLoggerTag, constants.STRING_BLANK), err)
	}
	game.eventExecutor = eventExecutor
	systemOperation := logic.NewSystemOperation(resourceManager, playerDataManager)
	systemOperation.SetOperationManager(game.operationManager)
	game.rootOperation = systemOperation
	return game
}

//func (this *Game) SetRootOperation(rootOperation logic.Operation) {
//	this.rootOperation = rootOperation
//}

func (this *Game) Start() {

}

func (this *Game) PostEvent(event model.Event) {
	this.eventExecutor.PostEvent(event)
}

func (this *Game) Process(event model.Event) error {
	this.logger.Debug(log_content.LogContentNormal(gameLoggerTag, "event process, event:%v", event))
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

	dataMap := this.playerDataManager.GetNeedToSavedData()

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
