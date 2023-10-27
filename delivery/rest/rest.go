package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fazarmitrais/atm-simulation/lib/responseFormatter"
	"github.com/fazarmitrais/atm-simulation/service"
	"github.com/gorilla/mux"
)

type Rest struct {
	service *service.Service
}

type ResponseFormatter struct {
	IsError bool   `json:"isError"`
	Mesage  string `json:"message"`
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
		responseFormatter.New(http.StatusBadRequest,
			fmt.Sprintf("Failed unmarshalling json : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	var acc service.Account
	err = json.Unmarshal(b, &acc)
	if err != nil {
		responseFormatter.New(http.StatusBadRequest,
			fmt.Sprintf("Failed unmarshalling json : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	errl := re.service.PINValidation(r.Context(), acc)
	if errl != nil {
		errl.ReturnAsJson(w)
		return
	}
	errl.ReturnAsJson(w)
}
