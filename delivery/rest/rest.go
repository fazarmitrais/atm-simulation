package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/fazarmitrais/atm-simulation/cookie"
	"github.com/fazarmitrais/atm-simulation/lib/responseFormatter"
	middleware "github.com/fazarmitrais/atm-simulation/middleware"
	"github.com/fazarmitrais/atm-simulation/service"
	"github.com/gorilla/mux"
)

type Rest struct {
	service *service.Service
	cookie  *cookie.Cookie
}

type ResponseFormatter struct {
	IsError bool   `json:"isError"`
	Mesage  string `json:"message"`
}

func New(svc *service.Service) *Rest {
	c := cookie.New()
	return &Rest{service: svc, cookie: c}
}

func (re *Rest) Register(m *mux.Router) {
	m = m.PathPrefix("/api/v1/account").Subrouter()
	m.HandleFunc("/validate", re.PINValidation).Methods(http.MethodPost)
	m.HandleFunc("/withdraw", middleware.Chain(re.Withdraw, middleware.Required(re.cookie))).Methods(http.MethodGet)
	m.HandleFunc("/exit", re.Exit).Methods(http.MethodGet)
}

func (re *Rest) Exit(w http.ResponseWriter, r *http.Request) {
	cookieStore, err := re.cookie.Store.Get(r, "cookie-store")
	if err != nil {
		responseFormatter.New(http.StatusBadRequest,
			fmt.Sprintf("Error getting cookie store : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	cookieStore.Values["authenticated"] = false
	cookieStore.Values["acctNbr"] = nil
	cookieStore.Save(r, w)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(responseFormatter.New(http.StatusOK, "Logout success", false))
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
	cookieStore, err := re.cookie.Store.Get(r, "cookie-store")
	if err != nil {
		responseFormatter.New(http.StatusInternalServerError,
			fmt.Sprintf("Error getting cookie store : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	cookieStore.Values["authenticated"] = true
	cookieStore.Values["acctNbr"] = acc.AccountNumber
	cookieStore.Save(r, w)
	errl.ReturnAsJson(w)
}

func (re *Rest) Withdraw(w http.ResponseWriter, r *http.Request) {
	cookieStore, err := re.cookie.Store.Get(r, "cookie-store")
	if err != nil {
		responseFormatter.New(http.StatusInternalServerError,
			fmt.Sprintf("Error getting cookie store : %s", err.Error()), true).
			ReturnAsJson(w)
		return
	}
	amountStr := r.URL.Query().Get("amount")
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		responseFormatter.New(http.StatusBadRequest, "Invalid ammount", true).
			ReturnAsJson(w)
		return
	}
	acc, resp := re.service.Withdraw(r.Context(), fmt.Sprintf("%v", cookieStore.Values["acctNbr"]), float64(amount))
	if resp != nil {
		resp.ReturnAsJson(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(acc)
}
