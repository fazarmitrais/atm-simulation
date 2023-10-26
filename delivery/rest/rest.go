package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fazarmitrais/atm-simulation/service"
	"github.com/gorilla/mux"
)

type Rest struct {
	service *service.Service
}

func New(svc *service.Service) *Rest {
	return &Rest{service: svc}
}

func (re *Rest) Register(m *mux.Router) {
	m = m.PathPrefix("/api/v1/account").Subrouter()
	m.HandleFunc("/validate", re.PINValidation)
}

func (re *Rest) PINValidation(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed reading request body : %s", err.Error()), http.StatusBadRequest)
		return
	}
	var acc service.Account
	err = json.Unmarshal(b, &acc)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed unmarshalling json : %s", err.Error()), http.StatusBadRequest)
		return
	}
	errl := re.service.PINValidation(r.Context(), acc)
	if errl != nil {
		http.Error(w, errl.Message, errl.StatusCode)
		return
	}
	json.NewEncoder(w).Encode("OK")
}
