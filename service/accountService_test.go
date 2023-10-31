package service

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/fazarmitrais/atm-simulation/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestPinValidation_AccountNbrIsRequired(t *testing.T) {
	svc := New()
	resp := svc.PINValidation(context.Background(), entity.Account{
		PIN: "456",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Account Number is required", resp.Message)
}

func TestPinValidation_PINIsRequired(t *testing.T) {
	svc := New()
	resp := svc.PINValidation(context.Background(), entity.Account{
		AccountNumber: "123",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "PIN is required", resp.Message)
}

// - Account Number should have 6 digits length. Display message `Account Number should have 6 digits length` for invalid Account Number.
func TestPinValidation_AccountNumberMustSixDigitsLength(t *testing.T) {
	svc := New()
	resp := svc.PINValidation(context.Background(), entity.Account{
		AccountNumber: "123",
		PIN:           "456",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Account Number should have 6 digits length", resp.Message)
}

//- PIN should have 6 digits length. Display message `PIN should have 6 digits length` for invalid PIN.

func TestPinValidation_PINMustSixDigitsLength(t *testing.T) {
	svc := New()
	resp := svc.PINValidation(context.Background(), entity.Account{
		AccountNumber: "123456",
		PIN:           "456",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "PIN should have 6 digits length", resp.Message)
}

// - Account Number should only contains numbers [0-9]. Display message `Account Number should only contains numbers` for invalid Account Number.
func TestPinValidation_AccountNumberOnlyContainsNumber(t *testing.T) {
	svc := New()
	resp := svc.PINValidation(context.Background(), entity.Account{
		AccountNumber: "a123456",
		PIN:           "123456",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Account Number should only contains numbers", resp.Message)
}

// - PIN should only contains numbers [0-9]. Display message `PIN should only contains numbers` for invalid PIN.
func TestPinValidation_PINOnlyContainsNumber(t *testing.T) {
	svc := New()
	resp := svc.PINValidation(context.Background(), entity.Account{
		AccountNumber: "123456",
		PIN:           "a123456",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "PIN should only contains numbers", resp.Message)
}

//- Check valid Acccount Number & PIN with ATM records. Display message `Invalid Account Number/PIN` if records is not exist.

func TestPinValidation_InvalidAccountNumber(t *testing.T) {
	svc := New()
	resp := svc.PINValidation(context.Background(), entity.Account{
		AccountNumber: "123456",
		PIN:           "1123456",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Invalid Account Number/PIN", resp.Message)
}

// - Check valid Acccount Number & PIN with ATM records. Display message `Invalid Account Number/PIN` if records is not exist.
func TestPinValidation_InvalidPIN(t *testing.T) {
	svc := New()
	resp := svc.PINValidation(context.Background(), entity.Account{
		AccountNumber: "112233",
		PIN:           "1123456",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Invalid Account Number/PIN", resp.Message)
}

func TestPinValidation_Success(t *testing.T) {
	svc := New()
	resp := svc.PINValidation(context.Background(), entity.Account{
		AccountNumber: "112233",
		PIN:           "012108",
	})
	assert.Nil(t, resp)
}

// - Maximum amount to withdraw is $1000. Display message `Maximum amount to withdraw is $1000` if withdraw amount is higher than $1000.
func TestWithdraw_MaxAmount1000(t *testing.T) {
	svc := New()
	_, resp := svc.Withdraw(context.Background(), "112233", 1001)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Maximum amount to withdraw is $1000", resp.Message)
}

// - Display message `Invalid ammount` if withdraw amount is not multiple of $10.
func TestWithdraw_AmountNotMultipleOf10(t *testing.T) {
	svc := New()
	_, resp := svc.Withdraw(context.Background(), "112233", 901)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Invalid ammount", resp.Message)
}

// - Display message `Insufficient balance $10` for insufficient balance. `$10` is the withdraw amount
func TestWithdraw_InsufficientBalance(t *testing.T) {
	svc := New()
	var withdrawAmt float64 = 200
	_, resp := svc.Withdraw(context.Background(), "112233", withdrawAmt)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, fmt.Sprintf("Insufficient balance $%0.f", withdrawAmt), resp.Message)
}

// - Display message `Invalid account` if account is not numbers
func TestTransfer_AccountMustBeNumbers(t *testing.T) {
	svc := New()
	_, resp := svc.Transfer(context.Background(), entity.Transfer{
		FromAccountNumber: "a432214213",
		ToAccountNumber:   "a432214214",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Invalid account", resp.Message)
}

// - Display message `Invalid account` if account is not found
func TestTransfer_FromAccountNumberMustBeCorrect(t *testing.T) {
	svc := New()
	_, resp := svc.Transfer(context.Background(), entity.Transfer{
		FromAccountNumber: "432214213",
		ToAccountNumber:   "112233",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Invalid account", resp.Message)
}

// - Display message `Invalid account` if account is not found
func TestTransfer_ToAccountNumberMustBeCorrect(t *testing.T) {
	svc := New()
	_, resp := svc.Transfer(context.Background(), entity.Transfer{
		FromAccountNumber: "112233",
		ToAccountNumber:   "432214214",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Invalid account", resp.Message)
}

// - Maximum amount to transfer is $1000. Display message `Maximum amount to transfer is $1000` if transfer amount is higher than $1000.
func TestTransfer_MaxTransferAmountIs1000(t *testing.T) {
	svc := New()
	_, resp := svc.Transfer(context.Background(), entity.Transfer{
		FromAccountNumber: "112233",
		ToAccountNumber:   "112244",
		Amount:            1001,
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Maximum amount to transfer is $1000", resp.Message)
}

// - Minimum amount to transfer is $1. Display message `Minimum amount to transfer is $1` if transfer amount is lower than $1.
func TestTransfer_MinTransferAmountIs1(t *testing.T) {
	svc := New()
	_, resp := svc.Transfer(context.Background(), entity.Transfer{
		FromAccountNumber: "112233",
		ToAccountNumber:   "112244",
		Amount:            0.5,
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Minimum amount to transfer is $1", resp.Message)
}

// - Display message `Insufficient balance $300` for insufficient balance. `$300` is the transfer amount
func TestTransfer_InsufficientBalance(t *testing.T) {
	svc := New()
	var trfAmount float64 = 200
	_, resp := svc.Transfer(context.Background(), entity.Transfer{
		FromAccountNumber: "112233",
		ToAccountNumber:   "112244",
		Amount:            trfAmount,
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, fmt.Sprintf("Insufficient balance $%0.f", trfAmount), resp.Message)
}

// - Display message `Invalid Reference Number` if reference number is not empty and not numbers
func TestTransfer_ReferenceNumberMustBeNumber(t *testing.T) {
	svc := New()
	_, resp := svc.Transfer(context.Background(), entity.Transfer{
		FromAccountNumber: "112233",
		ToAccountNumber:   "112244",
		Amount:            20,
		ReferenceNumber:   "Ref 213342",
	})
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Invalid Reference Number", resp.Message)
}

// - Valid amount will deduct the user balance with transfer amount and will add destination account with transfer amount. After that screen will
func TestTransfer_Success(t *testing.T) {
	svc := New()
	ctx := context.Background()
	_, resp := svc.Transfer(ctx, entity.Transfer{
		FromAccountNumber: "112233",
		ToAccountNumber:   "112244",
		Amount:            20,
		ReferenceNumber:   "213342",
	})
	fromAcct, _ := svc.BalanceCheck(ctx, "112233")
	toAcct, _ := svc.BalanceCheck(ctx, "112244")
	assert.Nil(t, resp)
	assert.Equal(t, float64(80), fromAcct.Balance)
	assert.Equal(t, float64(120), toAcct.Balance)
}
