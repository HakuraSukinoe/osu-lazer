package app

import (
	"context"
	"fmt"
	"github.com/deissh/osu-lazer/server/api"
	"github.com/deissh/osu-lazer/server/middlewares/customerror"
	"github.com/deissh/osu-lazer/server/middlewares/customlogger"
	"github.com/deissh/osu-lazer/server/middlewares/permission"
	"github.com/deissh/osu-lazer/server/model"
	"github.com/deissh/osu-lazer/server/services/cache"
	"github.com/deissh/osu-lazer/server/store"
	"github.com/deissh/osu-lazer/server/store/sqlstore"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"sync/atomic"
	"time"
)

type Server struct {
	HttpServer *echo.Echo

	Store         store.Store
	CacheProvider cache.Provider

	goroutineCount      int32
	goroutineExitSignal chan struct{}
	newStore            func() store.Store
	configPath          string
}

func NewServer(options ...Option) (*Server, error) {
	s := &Server{
		goroutineExitSignal: make(chan struct{}, 1),
	}

	for _, option := range options {
		if err := option(s); err != nil {
			return nil, errors.Wrap(err, "failed to apply option")
		}
	}

	config.WithOptions(config.ParseEnv, config.Readonly)
	config.AddDriver(yaml.Driver)
	err := config.LoadFiles(s.configPath, "config.yaml")
	if err != nil {
		panic(err)
	}

	if config.Bool("debug", false) {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: false,
			},
		).With().Caller().Logger()
	}

	logCurrentVersion := fmt.Sprintf("Current version is %v (%v/%v/%v/%v)", model.CurrentVersion, model.BuildNumber, model.BuildDate, model.BuildHash, model.BuildHashEnterprise)
	log.Info().
		Str("current_version", model.CurrentVersion).
		Str("build_number", model.BuildNumber).
		Str("build_date", model.BuildDate).
		Str("build_hash", model.BuildHash).
		Str("build_hash_enterprise", model.BuildHashEnterprise).
		Msg(logCurrentVersion)

	pwd, _ := os.Getwd()
	log.Info().Msg("Printing current working " + pwd)
	log.Info().Str("path", s.configPath).Msg("Loaded config")

	// init store
	settings := model.NewSqlSettings()
	s.Store = sqlstore.NewSqlSupplier(settings)

	// create new echo server and setup middlewares
	s.HttpServer = echo.New()
	s.HttpServer.HideBanner = true
	s.HttpServer.HidePort = true
	s.HttpServer.HTTPErrorHandler = customerror.CustomHTTPErrorHandler

	s.HttpServer.Use(middleware.RequestID())
	s.HttpServer.Use(customlogger.Middleware())
	s.HttpServer.Use(permission.GlobalMiddleware(s.Store))

	api.Init(s.Store, s.HttpServer.Group(""))

	return s, nil
}

func (s *Server) Start() error {
	addr := config.String("server.host") + ":" + config.String("server.port")

	go func() {
		log.Info().Msg("Server started on " + addr)
		err := s.HttpServer.Start(addr)
		if err != nil {
			log.Error().
				Err(err).
				Send()
		}
	}()

	return nil
}

func (s *Server) Shutdown() error {
	log.Info().Msg("Stopping HttpServer...")

	s.WaitForGoroutines()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.HttpServer.Shutdown(ctx); err != nil {
		log.Error().
			Err(err).
			Send()
	}

	if s.Store != nil {
		s.Store.Close()
	}
	if s.CacheProvider != nil {
		s.CacheProvider.Close()
	}

	log.Info().Msg("HttpServer stopped")
	return nil
}

// Go creates a goroutine, but maintains a record of it to ensure that execution completes before
// the server is shutdown.
func (s *Server) Go(f func()) {
	atomic.AddInt32(&s.goroutineCount, 1)

	go func() {
		f()

		atomic.AddInt32(&s.goroutineCount, -1)
		select {
		case s.goroutineExitSignal <- struct{}{}:
		default:
		}
	}()
}

// WaitForGoroutines blocks until all goroutines created by App.Go exit.
func (s *Server) WaitForGoroutines() {
	for atomic.LoadInt32(&s.goroutineCount) != 0 {
		<-s.goroutineExitSignal
	}
}
