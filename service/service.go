package service

import (
	"context"

	"github.com/fazarmitrais/atm-simulation/domain/entity"
	"github.com/fazarmitrais/atm-simulation/lib/responseFormatter"
)

type Service struct {
}

func New() *Service {
	initData()
	return &Service{}
}

type ServiceInterface interface {
	PINValidation(c context.Context, account entity.Account) *responseFormatter.ResponseFormatter
	Transfer(ctx, transfer entity.Transfer) (*entity.Account, *responseFormatter.ResponseFormatter)
	BalanceCheck(ctx context.Context, acctNbr string) (*entity.Account, *responseFormatter.ResponseFormatter)
}
