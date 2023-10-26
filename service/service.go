package service

type Service struct {
}

func New() *Service {
	initData()
	return &Service{}
}

type ServiceInterface interface {
}
