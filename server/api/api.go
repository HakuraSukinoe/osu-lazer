package api

import (
	"github.com/deissh/osu-lazer/server/store"
	"github.com/labstack/echo/v4"
)

type Routes struct {
	Root    *echo.Group // ''
	ApiRoot *echo.Group // '/api/v2'

	Users *echo.Group
}

type API struct {
	Store      store.Store
	BaseRoutes *Routes
}

func Init(store store.Store, root *echo.Group) {
	api := &API{
		Store:      store,
		BaseRoutes: &Routes{},
	}

	api.BaseRoutes.Root = root
	api.BaseRoutes.ApiRoot = api.BaseRoutes.Root.Group("/api/v2")

	api.BaseRoutes.Users = api.BaseRoutes.ApiRoot.Group("/users")

	api.InitUser()
}
