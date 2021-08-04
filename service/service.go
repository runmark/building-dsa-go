package service

import (
	"context"
	"fmt"
	"github.com/runmark/distribute-app-go/registry"
	"log"
	"net/http"
)

func Start(ctx context.Context, r registry.Registration, host, port string, registerHandlerFunc func()) (context.Context, error) {
	registerHandlerFunc()
	ctx = startService(ctx, r, host, port)

	err := registry.RegisterService(r)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func startService(ctx context.Context, r registry.Registration, host, port string) context.Context {

	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = host + ":" + port

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Printf("%v started. Press any key to stop.\n", r)
		var s string
		fmt.Scanln(&s)
		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx

}
