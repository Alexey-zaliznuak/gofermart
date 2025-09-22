package service

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/Alexey-zaliznuak/gofermart/internal/config"
	"github.com/Alexey-zaliznuak/gofermart/internal/model"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository/database"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository/user"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository/withdraw"
	"github.com/Alexey-zaliznuak/gofermart/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserService struct {
	repository         *user.UserRepository
	withdrawRepository *withdraw.WithdrawRepository
	auth               *AuthService
	*config.AppConfig
}

var (
	ErrUserAlreadyRegistered           = errors.New("user with same login already registered")
	ErrInvalidCredentials              = errors.New("invalid login and password pair")
	ErrUserAlreadyUsedOrderForWithdraw = errors.New("order for withdraw already used")
)

func (service *UserService) RegisterUser(request *model.RegisterUserRequest, ginCtx *gin.Context) (*model.User, error) {
	_, err := service.repository.GetByUsername(request.Login)

	if err != nil && err != database.ErrNotFound {
		return nil, err
	}

	if err != database.ErrNotFound {
		return nil, ErrUserAlreadyRegistered
	}

	return service.repository.CreateUser(request.Login, service.getPasswordHash(request.Password))
}

func (service *UserService) LoginUser(request *model.LoginUserRequest, ginCtx *gin.Context) (string, error) {
	user, err := service.repository.GetByUsername(request.Login)

	if err != nil {
		return "", err
	}

	if user.PasswordHash != service.getPasswordHash(request.Password) {
		return "", ErrInvalidCredentials
	}

	return service.auth.SaveAuthorization(user.ID, ginCtx)
}

func (service *UserService) GetUserBalanceInRubs(userID int) (*model.GetUserBalanceResponse, error) {
	user, err := service.repository.GetByID(userID)

	if err != nil {
		return nil, err
	}

	return &model.GetUserBalanceResponse{Current: float64(user.Balance) / 100, Withdrawn: float64(user.Withdraw) / 100}, nil
}

func (service *UserService) AddWithdraw(request model.AddWithdrawRequest, userID int) (*model.Withdraw, error) {
	if !utils.LuhnCheck(request.Order) {
		return nil, ErrWithdrawNumberIsInvalid
	}

	if _, err := service.withdrawRepository.GetByNumber(request.Order); err != database.ErrNotFound {
		return nil, ErrUserAlreadyUsedOrderForWithdraw
	}

	_, err := service.repository.Withdraw(int64(request.Sum*100), userID)

	if err != nil {
		return nil, err
	}

	return service.withdrawRepository.CreateWithdraw(request.Order, int64(request.Sum*100), userID)
}

func (service *UserService) getPasswordHash(password string) string {
	checksum := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", checksum[:])
}

func NewUserService(repository *user.UserRepository, withdrawRepository *withdraw.WithdrawRepository, config *config.AppConfig) *UserService {
	return &UserService{
		repository:         repository,
		withdrawRepository: withdrawRepository,
		auth:               NewAuthService(config),
		AppConfig:          config,
	}
}
