package app

import (
	"context"
	"fmt"
	"github.com/deissh/osu-lazer/server/api"
	"github.com/deissh/osu-lazer/server/mlog"
	"github.com/deissh/osu-lazer/server/model"
	"github.com/deissh/osu-lazer/server/store"
	"github.com/deissh/osu-lazer/server/store/sqlstore"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"sync/atomic"
	"time"
)

type Server struct {
	Store store.Store

	HttpServer *echo.Echo

	didFinishListen chan struct{}

	goroutineCount      int32
	goroutineExitSignal chan struct{}

	newStore func() store.Store

	Log *mlog.Logger
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

	if s.Log == nil {
		s.Log = mlog.NewLogger(&mlog.LoggerConfiguration{EnableConsole: true})
	}

	// Redirect default golang logger to this logger
	mlog.RedirectStdLog(s.Log)

	// Use this app logger as the global logger (eventually remove all instances of global logging)
	mlog.InitGlobalLogger(s.Log)

	logCurrentVersion := fmt.Sprintf("Current version is %v (%v/%v/%v/%v)", model.CurrentVersion, model.BuildNumber, model.BuildDate, model.BuildHash, model.BuildHashEnterprise)
	mlog.Info(
		logCurrentVersion,
		mlog.String("current_version", model.CurrentVersion),
		mlog.String("build_number", model.BuildNumber),
		mlog.String("build_date", model.BuildDate),
		mlog.String("build_hash", model.BuildHash),
		mlog.String("build_hash_enterprise", model.BuildHashEnterprise),
	)

	s.HttpServer = echo.New()
	s.HttpServer.HideBanner = true
	s.HttpServer.HidePort = true

	s.Store = sqlstore.NewSqlSupplier()

	return s, nil
}

func (s *Server) Start() error {
	api.Init(s.Store, s.HttpServer.Group(""))

	go func() {
		err := s.HttpServer.Start("127.0.0.1:2100")
		if err != nil {
			mlog.Err(err)
		}
	}()

	return nil
}

func (s *Server) Shutdown() error {
	mlog.Info("Stopping HttpServer...")

	s.WaitForGoroutines()

	if s.Store != nil {
		s.Store.Close()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.HttpServer.Shutdown(ctx); err != nil {
		mlog.Err(err)
		return err
	}

	mlog.Info("HttpServer stopped")
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
