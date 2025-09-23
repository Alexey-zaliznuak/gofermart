package service

import (
	"errors"

	"github.com/Alexey-zaliznuak/gofermart/internal/config"
	"github.com/Alexey-zaliznuak/gofermart/internal/model"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository/withdraw"
)

var (
	ErrWithdrawNumberIsInvalid = errors.New("invalid withdraw order") // не уверен нужна ли отдельная ошибка
)

type WithdrawService struct {
	repository *withdraw.WithdrawRepository

	*config.AppConfig
}

func (service *WithdrawService) GetAll(userID int) (model.GetWithdrawalsResponse, error) {
	return service.repository.GetAllByUserID(userID)
}

func NewWithdrawService(repository *withdraw.WithdrawRepository, config *config.AppConfig) *WithdrawService {
	return &WithdrawService{
		repository: repository,
		AppConfig:  config,
	}
}
