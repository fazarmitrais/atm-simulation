package main

import (
	"fmt"
	"net/http"

	"github.com/fazarmitrais/atm-simulation/delivery/rest"
	"github.com/fazarmitrais/atm-simulation/service"
	"github.com/gorilla/mux"
)

func main() {
	svc := service.New()
	re := rest.New(svc)
	m := mux.NewRouter()
	re.Register(m)
	fmt.Println("App is running on port 8080")
	http.ListenAndServe(":8080", m)
}
