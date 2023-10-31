package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fazarmitrais/atm-simulation/delivery/rest"
	"github.com/fazarmitrais/atm-simulation/service"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	envInit()
	svc := service.New()
	re := rest.New(svc)
	m := mux.NewRouter()
	re.Register(m)
	fmt.Println("App is running on port 8080")
	http.ListenAndServe(":8080", m)
}

func envInit() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}
}
