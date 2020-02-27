package permission

import (
	"errors"
	"github.com/deissh/osu-lazer/server/model"
	"github.com/deissh/osu-lazer/server/store"
	"github.com/labstack/echo/v4"
)

// keyFromHeader returns a `keyExtractor` that extracts key from the request header.
func keyFromHeader(header string) func(echo.Context) (string, error) {
	return func(c echo.Context) (string, error) {
		auth := c.Request().Header.Get(header)
		c.Logger().Info(auth)
		if auth == "" {
			return "", errors.New("missing key in request header")
		}
		if header == echo.HeaderAuthorization {
			l := len("Bearer")
			if len(auth) > l+1 && auth[:l] == "Bearer" {
				return auth[l+1:], nil
			}
			return "", errors.New("invalid key in the request header")
		}
		return auth, nil
	}
}

// GlobalMiddleware check access_token
func GlobalMiddleware(store store.Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			extractor := keyFromHeader(echo.HeaderAuthorization)

			// check token and write to context if user send one
			if key, err := extractor(c); err == nil {
				token, err := store.OAuth().ValidateToken(key)
				if err != nil {
					return err
				}

				current, err := store.User().UpdateLastVisit(token.UserID)
				if err != nil {
					return model.NewHTTPError(401, "auth_token_required", "Invalid token or user")
				}

				c.Set("current_user", current)
				c.Set("current_user_token", token)
			}

			return next(c)
		}
	}
}
