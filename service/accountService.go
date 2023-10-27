package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/fazarmitrais/atm-simulation/lib/responseFormatter"
)

type Account struct {
	Name          string  `json:"name"`
	AccountNumber string  `json:"accountNumber"`
	PIN           string  `json:"pin"`
	Balance       float64 `json:"balance"`
}

var accMap = make(map[string]*Account)

func initData() {
	accMap["112233"] = &Account{
		Name:          "John Doe",
		PIN:           "012108",
		Balance:       100,
		AccountNumber: "112233"}

	accMap["112244"] = &Account{
		Name:          "Jane Doe",
		PIN:           "932012",
		Balance:       100,
		AccountNumber: "112244"}
}

func (s *Service) PINValidation(c context.Context, account Account) *responseFormatter.ResponseFormatter {
	if strings.Trim(account.AccountNumber, " ") == "" {
		return responseFormatter.New(http.StatusBadRequest, "Account Number is required", true)
	} else if strings.Trim(account.PIN, " ") == "" {
		return responseFormatter.New(http.StatusBadRequest, "PIN is required", true)
	} else if len(account.AccountNumber) < 6 {
		return responseFormatter.New(http.StatusBadRequest, "Account Number should have 6 digits length", true)
	} else if len(account.PIN) < 6 {
		return responseFormatter.New(http.StatusBadRequest, "PIN should have 6 digits length", true)
	} else if _, err := strconv.Atoi(account.AccountNumber); err != nil {
		return responseFormatter.New(http.StatusBadRequest, "Account Number should only contains numbers", true)
	} else if _, err := strconv.Atoi(account.PIN); err != nil {
		return responseFormatter.New(http.StatusBadRequest, "PIN should only contains numbers", true)
	} else if accMap[account.AccountNumber] == nil || accMap[account.AccountNumber].PIN != account.PIN {
		return responseFormatter.New(http.StatusBadRequest, "Invalid Account Number/PIN", true)
	}
	return nil
}

func (s *Service) Withdraw(ctx context.Context, accountNumber string, withdrawAmount float64) (*Account, *responseFormatter.ResponseFormatter) {
	if accountNumber == "" {
		return nil, responseFormatter.New(http.StatusBadRequest, "Account Number is required", true)
	} else if withdrawAmount <= 0 {
		return nil, responseFormatter.New(http.StatusBadRequest, "Invalid withdraw amount", true)
	} else if withdrawAmount > 1000 {
		return nil, responseFormatter.New(http.StatusBadRequest, "Maximum amount to withdraw is $1000", true)
	}
	if accMap[accountNumber].Balance < withdrawAmount {
		return nil, responseFormatter.New(http.StatusBadRequest, fmt.Sprintf("Insufficient balance $%0.f", accMap[accountNumber].Balance), true)
	}
	accMap[accountNumber].Balance -= withdrawAmount
	return accMap[accountNumber], nil
}
