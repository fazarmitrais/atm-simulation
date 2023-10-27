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
}
