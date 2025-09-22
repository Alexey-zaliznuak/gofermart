package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Alexey-zaliznuak/gofermart/internal/model"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository/database"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository/user"
	"github.com/Alexey-zaliznuak/gofermart/internal/service"
	"github.com/gin-gonic/gin"
)

func registerUser(userService *service.UserService, authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := c.GetRawData()

		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		request := &model.RegisterUserRequest{}

		err = json.Unmarshal(body, &request)

		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		user, err := userService.RegisterUser(request, c)

		if err != nil {
			if err == service.ErrUserAlreadyRegistered {
				c.Status(http.StatusConflict)
				return
			}

			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		authService.SaveAuthorization(user.ID, c)
		c.Status(http.StatusOK)
	}
}

func loginUser(userService *service.UserService, authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := c.GetRawData()

		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		request := &model.LoginUserRequest{}

		err = json.Unmarshal(body, &request)

		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		jwt, err := userService.LoginUser(request, c)

		if err != nil {
			switch err {
			case database.ErrNotFound, service.ErrInvalidCredentials:
				c.Status(http.StatusUnauthorized)
				return
			default:
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		}

		c.String(http.StatusOK, jwt)
	}
}

func addOrder(orderService *service.OrderService, authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := authService.GetAuthorization(c)

		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		body, err := c.GetRawData()

		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		_, err = orderService.AddOrder(string(body), claims.UserID)

		if err != nil {
			switch err {
			case service.ErrOrderAlreadyAdded:
				c.Status(http.StatusOK)

			case service.ErrOrderAlreadyAddedByAnotherUser:
				c.Status(http.StatusConflict)

			case repository.ErrLuhnNumberIsInvalid:
				c.Status(http.StatusUnprocessableEntity)

			default:
				c.String(http.StatusInternalServerError, err.Error())
			}
			return
		}

		c.Status(http.StatusAccepted)
	}
}

func getUserOrders(orderService *service.OrderService, authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := authService.GetAuthorization(c)

		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		orders, err := orderService.GetAll(claims.UserID)

		if err != nil {
			if err == database.ErrNotFound {
				c.Status(http.StatusNoContent)
				return
			}
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, orders)
	}
}

func getUserBalance(userService *service.UserService, authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := authService.GetAuthorization(c)

		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		balance, err := userService.GetUserBalanceInRubs(claims.UserID)

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, balance)
	}
}

func addWithdraw(userService *service.UserService, authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := authService.GetAuthorization(c)

		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		body, err := c.GetRawData()

		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		request := &model.AddWithdrawRequest{}

		err = json.Unmarshal(body, &request)

		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		_, err = userService.AddWithdraw(*request, claims.UserID)

		if err != nil {
			switch err {
			case service.ErrUserAlreadyUsedOrderForWithdraw:
				c.Status(http.StatusConflict)

			case repository.ErrLuhnNumberIsInvalid:
				c.Status(http.StatusUnprocessableEntity)

			case user.ErrUserInsufficientFunds:
				c.Status(http.StatusPaymentRequired)

			default:
				c.String(http.StatusInternalServerError, err.Error())
			}
			return
		}

		c.Status(http.StatusOK)
	}
}

func getUserWithdrawals(withdrawService *service.WithdrawService, authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := authService.GetAuthorization(c)

		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		withdrawals, err := withdrawService.GetAll(claims.UserID)

		if err != nil {
			if err == database.ErrNotFound {
				c.Status(http.StatusNoContent)
				return
			}
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, withdrawals)
	}
}

func RegisterRoutes(router *gin.Engine, userService *service.UserService, orderService *service.OrderService, withdrawService *service.WithdrawService, authService *service.AuthService, db *sql.DB) {
	router.POST("/api/user/register", registerUser(userService, authService))
	router.POST("/api/user/login", loginUser(userService, authService))
	router.GET("/api/user/balance", getUserBalance(userService, authService))

	router.POST("/api/user/orders", addOrder(orderService, authService))
	router.GET("/api/user/orders", getUserOrders(orderService, authService))

	router.POST("/api/user/balance/withdraw", addWithdraw(userService, authService))
	router.GET("/api/user/balance/withdrawals", getUserWithdrawals(withdrawService, authService))
}
