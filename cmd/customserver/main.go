package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {

	var srv http.Server
	srv.Addr = "localhost:8080"

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {

		fmt.Printf("customserver has been started. Press any key to abort!")
		var input string
		fmt.Scanln(&input)
		srv.Shutdown(ctx)
		cancel()
	}()

	<- ctx.Done()

	fmt.Printf("customserver has been stopped.")
}
