package http

import (
	"fmt"
	"github.com/gorilla/mux"
	"invoice-service/cmd/config"
	"invoice-service/cmd/dependencies"
	"log"
	"net/http"
)

func StartServer(cfg *config.Config, dep *dependencies.Dependencies) error {
	router := mux.NewRouter()
	router = routes(*router, dep)

	err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.Server.Port), router)
	if err != nil {
		log.Fatalln("There's an error with the server", err)
	}
	return nil
}
