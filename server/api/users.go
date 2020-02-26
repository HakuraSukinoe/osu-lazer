package api

import (
	"github.com/labstack/echo/v4"
)

func (api *API) InitUser() {
	users := userHandlers{api}

	api.BaseRoutes.Users.GET("/:user/", users.getUserByToken)
}

type userHandlers struct {
	*API
}

func (api *userHandlers) getUserByToken(c echo.Context) error {
	mode := c.Param("mode")

	api.User().GetAll()

	return c.JSON(200, mode)
}
