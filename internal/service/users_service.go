package service

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/Alexey-zaliznuak/gofermart/internal/config"
	"github.com/Alexey-zaliznuak/gofermart/internal/model"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository/database"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository/user"
	"github.com/gin-gonic/gin"
)

type UserService struct {
	repository *user.UserRepository
	auth       *AuthService
	*config.AppConfig
}

var (
	ErrUserAlreadyRegistered = errors.New("user with same login already registered")
	ErrInvalidCredentials    = errors.New("invalid login and password pair")
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

func (service *UserService) getPasswordHash(password string) string {
	checksum := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", checksum[:])
}

// func (s *LinksService) GetFullURLFromShort(shortcut string) (string, error) {
// 	link, err := s.repository.GetByShortcut(shortcut)
// 	if err != nil {
// 		return "", err
// 	}
// 	return link.FullURL, nil
// }

// func (s *LinksService) GetUserLinks(c *gin.Context) ([]*model.GetUserLinksRequestItem, error) {
// 	claims, err := s.auth.GetAuthorization(c)

// 	if err != nil {
// 		return nil, err
// 	}

// 	links, err := s.repository.GetByUserID(claims.UserID)

// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, l := range links {
// 		l.Shortcut, err = s.BuildShortURL(l.Shortcut, c)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return links, err
// }

// func (s *LinksService) CreateLink(link *model.CreateLinkDto, c *gin.Context) (*model.Link, bool, error) {
// 	auth, err := s.auth.GetOrCreateAndSaveAuthorization(c)

// 	if err != nil {
// 		return nil, false, err
// 	}

// 	if !s.isValidURL(link.FullURL) {
// 		return link.NewLink(auth.UserID), false, fmt.Errorf("create link error: invalid URL: '%s'", link.FullURL)
// 	}

// 	if link.Shortcut == "" {
// 		var err error

// 		link.Shortcut, err = s.createUniqueShortcut()

// 		if err != nil {
// 			return link.NewLink(auth.UserID), false, err
// 		}
// 	}

// 	return s.repository.Create(link, auth.UserID, nil)
// }

// func (s *LinksService) DeleteUserLinks(shortcuts []string, c *gin.Context) error {
// 	auth, err := s.auth.GetOrCreateAndSaveAuthorization(c)

// 	if err != nil {
// 		return err
// 	}

// 	return s.repository.DeleteUserLinks(shortcuts, auth.UserID)
// }

// func (s *LinksService) BulkCreateWithCorrelationID(links []*model.CreateLinkWithCorrelationIDRequestItem, c *gin.Context) ([]*model.CreateLinkWithCorrelationIDResponseItem, error) {
// 	var result []*model.CreateLinkWithCorrelationIDResponseItem

// 	auth, err := s.auth.GetOrCreateAndSaveAuthorization(c)

// 	if err != nil {
// 		return nil, err
// 	}

// 	transactionExecuter, err := s.repository.GetTransactionExecuter(context.Background(), nil)
// 	supportTransaction := true

// 	if err != nil {
// 		if errors.Is(err, database.ErrExecuterNotSupportTransactions) {
// 			supportTransaction = false
// 		} else {
// 			return nil, err
// 		}
// 	}

// 	for index, link := range links {
// 		if !s.isValidURL(link.FullURL) {
// 			if supportTransaction {
// 				transactionExecuter.Commit()
// 			}
// 			return nil, fmt.Errorf("create link error: invalid URL: '%s'", link.FullURL)
// 		}

// 		shortcut, err := s.createUniqueShortcut()
// 		if err != nil {
// 			if supportTransaction {
// 				transactionExecuter.Commit()
// 			}
// 			return nil, err
// 		}

// 		l := &model.CreateLinkDto{FullURL: link.FullURL, Shortcut: shortcut}

// 		newLink, _, err := s.repository.Create(l, auth.UserID, transactionExecuter)

// 		if err != nil {
// 			if supportTransaction {
// 				transactionExecuter.Commit()
// 			}
// 			return nil, err
// 		}

// 		shortcut, err = s.BuildShortURL(newLink.Shortcut, c)

// 		if err != nil {
// 			if supportTransaction {
// 				transactionExecuter.Commit()
// 			}
// 			return nil, err
// 		}

// 		result = append(result, &model.CreateLinkWithCorrelationIDResponseItem{CorrelationID: link.CorrelationID, Shortcut: shortcut})

// 		if supportTransaction && (index+1%1000 == 0 || index == len(links)-1) {
// 			transactionExecuter.Commit()
// 			transactionExecuter, err = s.repository.GetTransactionExecuter(context.Background(), nil)
// 			if err != nil {
// 				return nil, err
// 			}
// 		}
// 	}

// 	if supportTransaction {
// 		transactionExecuter.Commit()
// 	}

// 	return result, nil
// }

// func (s *LinksService) createUniqueShortcut() (string, error) {
// 	maxAttempts := 5

// 	newShortcut := s.generateShortcut(s.AppConfig.Server.ShortLinksLength)

// 	for range maxAttempts {
// 		_, err := s.repository.GetByShortcut(newShortcut)
// 		if err != nil {
// 			newShortcut = s.generateShortcut(s.AppConfig.Server.ShortLinksLength)
// 			continue
// 		}
// 		break
// 	}

// 	if _, err := s.repository.GetByShortcut(newShortcut); err != database.ErrNotFound {
// 		return "", fmt.Errorf("create link error: could not generate unique shortcut after %d attempts", maxAttempts)
// 	}

// 	return newShortcut, nil
// }

// func (s *LinksService) generateShortcut(length int) string {
// 	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// 	result := make([]rune, length)

// 	for i := range result {
// 		result[i] = letters[rand.Intn(len(letters))]
// 	}
// 	return string(result)
// }

// func (s *LinksService) BuildShortURL(shortcut string, c *gin.Context) (string, error) {
// 	prefix := s.AppConfig.Server.BaseURL
// 	if prefix == "" {
// 		prefix = fmt.Sprintf("http://%s/", c.Request.Host)
// 	}
// 	return url.JoinPath(prefix, shortcut)
// }

// func (s *LinksService) isValidURL(u string) bool {
// 	parsedURL, err := url.ParseRequestURI(u)
// 	if err != nil {
// 		return false
// 	}

// 	if parsedURL.Scheme == "" || parsedURL.Host == "" {
// 		return false
// 	}

// 	return true
// }

func NewUserService(repository *user.UserRepository, config *config.AppConfig) *UserService {
	return &UserService{
		repository: repository,
		auth:       NewAuthService(config),
		AppConfig:  config,
	}
}
