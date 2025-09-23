package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Alexey-zaliznuak/gofermart/internal/config"
	"github.com/Alexey-zaliznuak/gofermart/internal/logger"
	"github.com/Alexey-zaliznuak/gofermart/internal/model"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository/database"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository/order"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository/user"
	"github.com/Alexey-zaliznuak/gofermart/internal/utils"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

var (
	ErrOrderAlreadyAdded              = errors.New("order already added")
	ErrOrderAlreadyAddedByAnotherUser = errors.New("order already added by another user")
)

type OrderService struct {
	repository     *order.OrderRepository
	userRepository *user.UserRepository

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

func (service *OrderService) StartWorker() {
	for {
		order, err := service.repository.GetFirstNotProcessed()

		if err == database.ErrNotFound {
			logger.Log.Info("Необработанные заказы не найдены")
			time.Sleep(time.Second)
			continue
		}

		client := resty.New()

		resp, err := client.R().Get(fmt.Sprintf("%s/api/orders/%s", service.AppConfig.AccrualSystemAddress, order.Number))

		if err != nil {
			logger.Log.Error("Ошибка при попытке получить информацию по заказу")
			time.Sleep(time.Second)
			continue
		}

		if resp.StatusCode() == http.StatusTooManyRequests {
			logger.Log.Info("Превышено количество запросов к системе расчета начислений баллов лояльности")
			time.Sleep(time.Minute)
			continue
		}

		if resp.StatusCode() == http.StatusNoContent {
			logger.Log.Info("По данному заказу не найдено никакой информации")
			time.Sleep(time.Second)
			continue
		}

		if resp.StatusCode() != http.StatusOK {
			logger.Log.Info("Неожиданный статус код")
			continue
		}

		payload := &model.AccrualResponse{}

		err = json.Unmarshal(resp.Body(), payload)

		if err != nil {
			logger.Log.Error("Ошибка обработки ответа от системы расчета начислений баллов лояльности")
			time.Sleep(time.Second)
			continue
		}

		logger.Log.Info("Получен ответ от системы расчета начислений баллов лояльности", zap.Any("payload", payload))

		order.Status = payload.Status
		if payload.Accrual != nil {
			order.Accrual = payload.Accrual

			_, err := service.userRepository.AddBalance(*payload.Accrual, order.UserID)

			if err != nil {
				logger.Log.Error("Ошибка начисления баланса за заказ")
				continue
			}
		}

		err = service.repository.Update(order)

		if err != nil {
			logger.Log.Error("Ошибка сохранения результата от системы расчета начислений баллов лояльности")
			continue
		}
	}
}

func NewOrderService(repository *order.OrderRepository, config *config.AppConfig) *OrderService {
	return &OrderService{
		repository: repository,
		AppConfig:  config,
	}
}
