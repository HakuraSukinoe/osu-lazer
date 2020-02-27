package permission

import (
	"github.com/deissh/osu-lazer/server/model"
	"github.com/labstack/echo/v4"
)

// MustLogin check
var MustLogin = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		data := c.Get("current_user")

		_, ok := data.(*model.UserFull)
		if !ok {
			return model.NewHTTPError(401, "auth_token_required", "Need auth token in headers")
		}

		return next(c)
	}
}

// CanModerate check
// noinspection GoUnusedGlobalVariable
var CanModerate = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		data := c.Get("current_user")

		current, ok := data.(*model.UserFull)
		if !ok {
			return model.NewHTTPError(401, "auth_token_required", "Need auth token in headers")
		}

		if !current.CanModerate {
			return model.NewHTTPError(401, "user_permissions", "User can not moderate")
		}

		return next(c)
	}
}

// IsSupporter check
// noinspection GoUnusedGlobalVariable
var IsSupporter = func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		data := c.Get("current_user")

		current, ok := data.(*model.UserFull)
		if !ok {
			return model.NewHTTPError(401, "auth_token_required", "Need auth token in headers")
		}

		if !current.IsSupporter {
			return model.NewHTTPError(401, "user_permissions", "User does not have an active supporter tag")
		}

		return next(c)
	}
}
