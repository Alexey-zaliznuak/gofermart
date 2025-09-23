package service

import (
	"net/http"

	"github.com/Alexey-zaliznuak/gofermart/internal/config"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository"
	"github.com/gin-gonic/gin"
)

type AuthService struct {
	Repository *repository.AuthRepository
}

func (service *AuthService) GetAuthorization(c *gin.Context) (*repository.Claims, error) {
	auth, err := c.Cookie("Authorization")

	if err != nil {
		if err != http.ErrNoCookie {
			return nil, err
		}
		auth = c.GetHeader("Authorization")

		if auth == "" {
			return nil, http.ErrNoCookie
		}
	}

	return service.Repository.ParsePayload(auth)
}

func (service *AuthService) SaveAuthorization(UserID int, c *gin.Context) (string, error) {
	jwt, err := service.Repository.BuildJWTString(UserID)

	if err != nil {
		return "", err
	}

	c.SetCookie("Authorization", jwt, 86400, "/", "", false, true)

	return jwt, nil
}

func NewAuthService(config *config.AppConfig) *AuthService {
	return &AuthService{Repository: repository.NewAuthRepository(config)}
}
