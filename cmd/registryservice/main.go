package main

import (
	"context"
	"fmt"
	"github.com/runmark/distribute-app-go/registry"
	"log"
	"net/http"
)

func main() {

	registry.SetupRegistryService()
	http.Handle("/services", &registry.RegistryService{})

	var srv http.Server
	srv.Addr = registry.ServerPort

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		log.Println("service registry started. Press any key to stop.")
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel()
	}()

	<-ctx.Done()

	log.Println("service registry stopped.")
}
