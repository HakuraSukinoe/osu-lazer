package api

import (
	"github.com/labstack/echo/v4"
	"strconv"
)

func (api *API) InitUser() {
	users := userHandlers{api}

	api.BaseRoutes.Users.GET("/:user", users.getUserByToken)
	api.BaseRoutes.Users.GET("/:user/:mode", users.getUserByToken)
}

type userHandlers struct {
	*API
}

func (api *userHandlers) getUserByToken(c echo.Context) error {
	mode := c.Param("mode")

	userID, _ := strconv.ParseUint(c.Param("user"), 10, 32)
	user, _ := api.Store.User().Get(uint(userID), mode)

	return c.JSON(200, user)
}
