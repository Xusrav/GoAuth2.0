package main

import (
	"log"

	"github.com/Xusrav/GoAuth2.0/cmd/app"
	"github.com/Xusrav/GoAuth2.0/cmd/app/handlers"
	// "github.com/gorilla/mux"
	"go.uber.org/dig"
)



func main() {
	
	deps := []interface{}{
		handlers.NewHandler,
		app.NewServer,
		// mux.NewRouter,
	}

	container := dig.New()
	for _, dep := range deps {
		err := container.Provide(dep)
		if err != nil {
			log.Fatal(err)
		}
	}

	err := container.Invoke(func(server *app.Server) {
		server.Run()
	})
	if err != nil {
		log.Fatal(err)
	}	
}
