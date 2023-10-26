package service

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	errorlib "github.com/fazarmitrais/atm-simulation/errorlib"
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

func (s *Service) PINValidation(c context.Context, account Account) *errorlib.CustomError {
	if strings.Trim(account.AccountNumber, " ") == "" {
		return errorlib.New(http.StatusBadRequest, "Account Number is required")
	} else if strings.Trim(account.PIN, " ") == "" {
		return errorlib.New(http.StatusBadRequest, "PIN is required")
	} else if len(account.AccountNumber) < 6 {
		return errorlib.New(http.StatusBadRequest, "Account Number should have 6 digits length")
	} else if len(account.PIN) < 6 {
		return errorlib.New(http.StatusBadRequest, "PIN should have 6 digits length")
	} else if _, err := strconv.Atoi(account.AccountNumber); err != nil {
		return errorlib.New(http.StatusBadRequest, "Account Number should only contains numbers")
	} else if _, err := strconv.Atoi(account.PIN); err != nil {
		return errorlib.New(http.StatusBadRequest, "PIN should only contains numbers")
	} else if accMap[account.AccountNumber] == nil || accMap[account.AccountNumber].PIN != account.PIN {
		return errorlib.New(http.StatusBadRequest, "Invalid Account Number/PIN")
	}
	return nil
}
