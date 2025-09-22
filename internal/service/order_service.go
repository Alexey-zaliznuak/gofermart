package service

import (
	"errors"

	"github.com/Alexey-zaliznuak/gofermart/internal/config"
	"github.com/Alexey-zaliznuak/gofermart/internal/model"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository/database"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository/order"
	"github.com/Alexey-zaliznuak/gofermart/internal/utils"
)

var (
	ErrOrderAlreadyAdded              = errors.New("order already added")
	ErrOrderAlreadyAddedByAnotherUser = errors.New("order already added by another user")
	ErrOrderNumberIsInvalid           = errors.New("invalid login and password pair") // не уверен нужна ли отдельная ошибка
)

type OrderService struct {
	repository *order.OrderRepository

	*config.AppConfig
}

func (service *OrderService) AddOrder(number string, claims *repository.Claims) (*model.Order, error) {
	if !utils.LuhnCheck(number) {
		return nil, ErrOrderNumberIsInvalid
	}

	order, err := service.repository.GetByNumber(number)

	if err != database.ErrNotFound {
		if order.UserID == claims.UserID {
			return nil, ErrOrderAlreadyAdded
		}
		return nil, ErrOrderAlreadyAddedByAnotherUser
	}

	return service.repository.CreateOrder(number, claims.UserID)
}

func NewOrderService(repository *order.OrderRepository, config *config.AppConfig) *OrderService {
	return &OrderService{
		repository: repository,
		AppConfig:  config,
	}
}
