package service

import (
	"context"

	"github.com/fazarmitrais/atm-simulation/lib/responseFormatter"
)

type Service struct {
}

func New() *Service {
	initData()
	return &Service{}
}

type ServiceInterface interface {
	PINValidation(c context.Context, account Account) *responseFormatter.ResponseFormatter
	Transfer(ctx, transfer Transfer) (*Account, *responseFormatter.ResponseFormatter)
	BalanceCheck(ctx context.Context, acctNbr string) (*Account, *responseFormatter.ResponseFormatter)
}
