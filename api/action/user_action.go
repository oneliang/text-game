package action

import (
	"fmt"
	httpAction "github.com/oneliang/frame-golang/http/action"
	"github.com/oneliang/util-golang/constants"
	"net/http"
	"service/service"
	"time"
)

type UserAction struct {
	UserService *service.UserService
}

func NewUserAction() *UserAction {
	return &UserAction{}
}

func (this *UserAction) RouteMap() map[string]httpAction.ActionExecuteFunction {
	return map[string]httpAction.ActionExecuteFunction{
		constants.HTTP_REQUEST_METHOD_POST + constants.SYMBOL_AT + "/create": this.Create,
		constants.HTTP_REQUEST_METHOD_GET + constants.SYMBOL_AT + "/read":    this.Read,
		constants.HTTP_REQUEST_METHOD_POST + constants.SYMBOL_AT + "/update": this.Update,
		constants.HTTP_REQUEST_METHOD_POST + constants.SYMBOL_AT + "/delete": this.Delete,
		constants.HTTP_REQUEST_METHOD_POST + constants.SYMBOL_AT + "/stream": this.Stream,
	}
}

func (this *UserAction) Create(request *http.Request, writer http.ResponseWriter) (error, []byte, int) {
	fmt.Println("Create")
	this.UserService.Create()
	return nil, nil, 200
}

func (this *UserAction) Read(request *http.Request, writer http.ResponseWriter) (error, []byte, int) {
	fmt.Println("Read")
	this.UserService.Read()
	return nil, nil, 200
}

func (this *UserAction) Update(request *http.Request, writer http.ResponseWriter) (error, []byte, int) {
	fmt.Println("Update")
	this.UserService.Update()
	return nil, nil, 200
}

func (this *UserAction) Delete(request *http.Request, writer http.ResponseWriter) (error, []byte, int) {
	fmt.Println("Delete")
	this.UserService.Delete()
	return nil, nil, 200
}

func (this *UserAction) Stream(request *http.Request, writer http.ResponseWriter) (error, []byte, int) {
	fmt.Println("Stream")
	writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	writer.Header().Set("Cache-Control", "no-cache")
	writer.Header().Set("Connection", "keep-alive")
	flusher := writer.(http.Flusher)
	for i := 0; i < 10; i++ {
		_, err := writer.Write([]byte(fmt.Sprintf("stream:%d", i)))
		if err != nil {
			fmt.Println(fmt.Sprintf("error:%+v", err))
		}
		time.Sleep(1 * time.Second)
		flusher.Flush()
	}

	return nil, nil, 200
}
