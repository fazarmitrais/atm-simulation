package service

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/fazarmitrais/atm-simulation/lib/responseFormatter"
)

/*
  --
   - Name          : John Doe
     PIN           : 012108
     Balance       : $100
     Account Number: 112233

   - Name          : Jane Doe

     PIN           : 932012

     Balance       : $30

     Account Number: 112244
- Account Number should have 6 digits length. Display message `Account Number should have 6 digits length` for invalid Account Number.
  - Account Number should only contains numbers [0-9]. Display message `Account Number should only contains numbers` for invalid Account Number.
  - PIN should have 6 digits length. Display message `PIN should have 6 digits length` for invalid PIN.
  - PIN should only contains numbers [0-9]. Display message `PIN should only contains numbers` for invalid PIN.
  - Check valid Acccount Number & PIN with ATM records. Display message `Invalid Account Number/PIN` if records is not exist.
  - Valid Account Number & PIN will take user to __Transaction Screen__.
*/

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
