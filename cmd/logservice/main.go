package main

import (
	"context"
	"fmt"
	"github.com/runmark/distribute-app-go/log"
	"github.com/runmark/distribute-app-go/registry"
	"github.com/runmark/distribute-app-go/service"

	stdlog "log"
)

func main() {

	log.Run("./app.log")

	ctx := context.Background()

	Host, Port := "localhost", "50011"

	var reg registry.Registration
	reg.ServiceName = registry.LogService
	reg.ServiceUrl = fmt.Sprintf("http://%v:%v", Host, Port)
	reg.RequiredServices = make([]registry.ServiceName, 0)
	reg.UpdateServiceURL = reg.ServiceUrl + "/services"
	reg.HearbeatURL = reg.ServiceUrl + "/heartbeat"

	ctx, err := service.Start(ctx, reg, Host, Port, log.RegisterHandlers)
	if err != nil {
		stdlog.Fatal(err)
	}

	<-ctx.Done()
	fmt.Println("Shutting down log service")
}
