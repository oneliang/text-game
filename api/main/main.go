package main

import (
	"api/action"
	"fmt"
	"github.com/oneliang/frame-golang/context"
	httpAction "github.com/oneliang/frame-golang/http/action"
	"github.com/oneliang/frame-golang/http/server"
	"github.com/oneliang/frame-golang/ioc"
	"log"
	"net/http"
	"service/service"
)

func main() {

	//moonshot.RequestMoonshot()
	//moonshot.TestJsonDemo()
	//return
	globalContext := context.NewGlobalContext()
	actionContext := httpAction.NewActionContext()
	if putContextErr := globalContext.PutContext("ActionContext", actionContext, action.ProviderSet); putContextErr != nil {
		log.Fatalf("%+v", putContextErr)
		return
	}
	iocContext := ioc.NewIocContext()
	if putContextErr := globalContext.PutContext("IocContext", iocContext, service.ProviderSet); putContextErr != nil {
		log.Fatalf("%+v", putContextErr)
		return
	}

	globalContext.Initialize(nil)
	iocContext.AutoInject()

	a := actionContext.GetActionInstanceMap()
	for _, value := range a {
		log.Println(fmt.Sprintf("%p", value))
		userAction := value.(*action.UserAction)
		log.Println(fmt.Sprintf("%p", userAction.UserService))
		log.Println(fmt.Sprintf("%p", userAction.UserService.RoleService))
	}

	//return
	serverHandler := server.NewServerHandler(actionContext.GetActionExecuteFunctionMap())
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: serverHandler,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("server.ListenAndServe error:%v", err)
	}
}
