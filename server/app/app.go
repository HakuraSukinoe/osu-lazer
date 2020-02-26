package app

import (
	"context"
)

type App struct {
	srv *Server

	requestId      string
	ipAddress      string
	path           string
	userAgent      string
	acceptLanguage string

	context context.Context
}

type Option func(a *App)

func New(options ...Option) *App {
	app := &App{}

	for _, option := range options {
		option(app)
	}

	return app
}
