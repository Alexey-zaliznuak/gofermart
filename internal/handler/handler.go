package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Alexey-zaliznuak/gofermart/internal/model"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository/database"
	"github.com/Alexey-zaliznuak/gofermart/internal/service"
	"github.com/gin-gonic/gin"
)

// func redirect(linksService *service.LinksService) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		shortcut := c.Param("shortcut")
// 		fullURL, err := linksService.GetFullURLFromShort(shortcut)

// 		if err != nil {
// 			if err == database.ErrObjectDeleted {
// 				c.Status(http.StatusGone)
// 				return
// 			}

// 			c.String(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		c.Redirect(http.StatusTemporaryRedirect, fullURL)
// 	}
// }

// func createLink(linksService *service.LinksService) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		body, err := c.GetRawData()

// 		if err != nil {
// 			c.String(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		link := &model.Link{FullURL: string(body)}
// 		link, created, err := linksService.CreateLink(link.ToCreateDto(), c)

// 		if err != nil {
// 			c.String(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		url, err := linksService.BuildShortURL(link.Shortcut, c)

// 		if err != nil {
// 			c.String(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		status := http.StatusCreated
// 		if !created {
// 			fmt.Printf("Duplicate: %s", link.FullURL)
// 			status = http.StatusConflict
// 		}

// 		c.String(status, url)
// 	}
// }

// func createLinkWithJSONAPI(linksService *service.LinksService) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		body, err := c.GetRawData()

// 		if err != nil {
// 			c.String(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		request := &model.CreateShortURLRequest{}
// 		err = json.Unmarshal(body, &request)

// 		if err != nil {
// 			c.String(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		l := &model.CreateLinkDto{FullURL: request.FullURL}
// 		link, created, err := linksService.CreateLink(l, c)

// 		if err != nil {
// 			c.String(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		shortURL, err := linksService.BuildShortURL(link.Shortcut, c)

// 		if err != nil {
// 			c.String(http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		status := http.StatusCreated
// 		if !created {
// 			status = http.StatusConflict
// 		}

// 		c.JSON(status, &model.CreateShortURLResponse{Result: shortURL})
// 	}
// }

// func createLinkBatch(linksService *service.LinksService) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		body, err := c.GetRawData()

// 		if err != nil {
// 			c.String(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		request := make([]*model.CreateLinkWithCorrelationIDRequestItem, 0, 100)

// 		err = json.Unmarshal(body, &request)

// 		if err != nil {
// 			c.String(http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		response, err := linksService.BulkCreateWithCorrelationID(request, c)

// 		if err != nil {
// 			c.String(http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		c.JSON(http.StatusCreated, response)
// 	}
// }

// func getUserLinks(linksService *service.LinksService, authService *service.AuthService) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		status := http.StatusOK
// 		links, err := linksService.GetUserLinks(c)

// 		if err != nil && err != http.ErrNoCookie {
// 			c.String(http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		_, err = authService.CreateAndSaveAuthorization(c)

// 		if err != nil {
// 			c.String(http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		if err == http.ErrNoCookie {
// 			c.String(http.StatusNoContent, "")

// 			_, err = authService.CreateAndSaveAuthorization(c)

// 			if err != nil {
// 				c.String(http.StatusInternalServerError, err.Error())
// 				return
// 			}
// 			return
// 		}

// 		if err == repository.ErrTokenValidation {
// 			_, err = authService.CreateAndSaveAuthorization(c)
// 			status = http.StatusNoContent
// 		}

// 		if err != nil {
// 			c.String(http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		if len(links) == 0 {
// 			status = http.StatusNoContent
// 		}

// 		c.JSON(status, links)
// 	}
// }

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

		_, err = orderService.AddOrder(string(body), claims)

		if err != nil {
			switch err {
			case service.ErrOrderAlreadyAdded:
				c.Status(http.StatusOK)

			case service.ErrOrderAlreadyAddedByAnotherUser:
				c.Status(http.StatusConflict)

			case service.ErrOrderNumberIsInvalid:
				c.Status(http.StatusUnprocessableEntity)

			default:
				c.String(http.StatusInternalServerError, err.Error())
			}
			return
		}

		c.Status(http.StatusAccepted)
	}
}

func RegisterRoutes(router *gin.Engine, userService *service.UserService, orderService *service.OrderService, authService *service.AuthService, db *sql.DB) {
	// router.GET("/:shortcut", redirect(linksService))

	// router.POST("/", createLink(linksService))
	router.POST("/api/user/register", registerUser(userService, authService))
	router.POST("/api/user/login", loginUser(userService, authService))

	router.POST("/api/user/orders", addOrder(orderService, authService))

	// router.POST("/api/shorten/batch", createLinkBatch(linksService))

	// router.GET("/api/user/urls", getUserLinks(linksService, authService))
	// router.DELETE("/api/user/urls", deleteUserLinks(linksService))
}
