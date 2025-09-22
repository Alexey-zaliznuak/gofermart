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
)

type OrderService struct {
	repository *order.OrderRepository

	*config.AppConfig
}

func (service *OrderService) AddOrder(number string, userID int) (*model.Order, error) {
	if !utils.LuhnCheck(number) {
		return nil, repository.ErrLuhnNumberIsInvalid
	}

	order, err := service.repository.GetByNumber(number)

	if err != database.ErrNotFound {
		if order.UserID == userID {
			return nil, ErrOrderAlreadyAdded
		}
		return nil, ErrOrderAlreadyAddedByAnotherUser
	}

	return service.repository.CreateOrder(number, userID)
}

func (service *OrderService) GetAll(userID int) (model.GetOrdersResponse, error) {
	return service.repository.GetAllByUserID(userID)
}

func NewOrderService(repository *order.OrderRepository, config *config.AppConfig) *OrderService {
	return &OrderService{
		repository: repository,
		AppConfig:  config,
	}
}
