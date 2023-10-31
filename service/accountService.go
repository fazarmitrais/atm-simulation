package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/fazarmitrais/atm-simulation/domain/entity"
	"github.com/fazarmitrais/atm-simulation/lib/responseFormatter"
)

var accMap = make(map[string]*entity.Account)

func initData() {
	accMap["112233"] = &entity.Account{
		Name:          "John Doe",
		PIN:           "012108",
		Balance:       100,
		AccountNumber: "112233"}

	accMap["112244"] = &entity.Account{
		Name:          "Jane Doe",
		PIN:           "932012",
		Balance:       100,
		AccountNumber: "112244"}
}

func (s *Service) PINValidation(c context.Context, account entity.Account) *responseFormatter.ResponseFormatter {
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

func (s *Service) Withdraw(ctx context.Context, accountNumber string, withdrawAmount float64) (*entity.AccountResponse, *responseFormatter.ResponseFormatter) {
	if accountNumber == "" {
		return nil, responseFormatter.New(http.StatusBadRequest, "Account Number is required", true)
	} else if withdrawAmount <= 0 {
		return nil, responseFormatter.New(http.StatusBadRequest, "Invalid withdraw amount", true)
	} else if withdrawAmount > 1000 {
		return nil, responseFormatter.New(http.StatusBadRequest, "Maximum amount to withdraw is $1000", true)
	} else if int(withdrawAmount)%10 != 0 {
		return nil, responseFormatter.New(http.StatusBadRequest, "Invalid ammount", true)
	} else if accMap[accountNumber].Balance < withdrawAmount {
		return nil, responseFormatter.New(http.StatusBadRequest, fmt.Sprintf("Insufficient balance $%0.f", withdrawAmount), true)
	}
	accMap[accountNumber].Balance -= withdrawAmount
	return accMap[accountNumber].ToAccountResponse(), nil
}

func (s *Service) BalanceCheck(ctx context.Context, acctNbr string) (*entity.AccountResponse, *responseFormatter.ResponseFormatter) {
	if strings.Trim(acctNbr, " ") == "" {
		return nil, responseFormatter.New(http.StatusBadRequest, "Account Number is required", true)
	} else if len(acctNbr) < 6 {
		return nil, responseFormatter.New(http.StatusBadRequest, "Account Number should have 6 digits length", true)
	} else if _, err := strconv.Atoi(acctNbr); err != nil {
		return nil, responseFormatter.New(http.StatusBadRequest, "Account Number should only contains numbers", true)
	} else if accMap[acctNbr] == nil {
		return nil, responseFormatter.New(http.StatusBadRequest, "Invalid Account Number/PIN", true)
	}
	return accMap[acctNbr].ToAccountResponse(), nil
}

func (s *Service) Transfer(ctx context.Context, transfer entity.Transfer) (*entity.AccountResponse, *responseFormatter.ResponseFormatter) {
	if transfer.FromAccountNumber == "" || transfer.ToAccountNumber == "" {
		return nil, responseFormatter.New(http.StatusBadRequest, "Account Number is required", true)
	} else if transfer.FromAccountNumber == transfer.ToAccountNumber {
		return nil, responseFormatter.New(http.StatusBadRequest, "From and Destination account number cannot be the same", true)
	} else if _, err := strconv.Atoi(transfer.FromAccountNumber); err != nil {
		return nil, responseFormatter.New(http.StatusBadRequest, "Invalid account", true)
	} else if accMap[transfer.FromAccountNumber] == nil {
		return nil, responseFormatter.New(http.StatusBadRequest, "Invalid account", true)
	} else if accMap[transfer.ToAccountNumber] == nil {
		return nil, responseFormatter.New(http.StatusBadRequest, "Invalid account", true)
	} else if transfer.Amount <= 0 {
		return nil, responseFormatter.New(http.StatusBadRequest, "Invalid transfer amount", true)
	} else if transfer.Amount > 1000 {
		return nil, responseFormatter.New(http.StatusBadRequest, "Maximum amount to transfer is $1000", true)
	} else if transfer.Amount < 1 {
		return nil, responseFormatter.New(http.StatusBadRequest, "Minimum amount to transfer is $1", true)
	} else if accMap[transfer.FromAccountNumber].Balance < transfer.Amount {
		return nil, responseFormatter.New(http.StatusBadRequest, fmt.Sprintf("Insufficient balance $%0.f", transfer.Amount), true)
	} else if strings.Trim(transfer.ReferenceNumber, " ") != "" {
		if _, err = strconv.Atoi(transfer.ReferenceNumber); err != nil {
			return nil, responseFormatter.New(http.StatusBadRequest, "Invalid Reference Number", true)
		}
	}
	accMap[transfer.FromAccountNumber].Balance -= transfer.Amount
	accMap[transfer.ToAccountNumber].Balance += transfer.Amount
	return accMap[transfer.FromAccountNumber].ToAccountResponse(), nil
}
