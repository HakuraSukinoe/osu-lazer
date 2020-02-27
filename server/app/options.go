package app

import (
	"github.com/deissh/osu-lazer/server/store"
	"github.com/pkg/errors"
)

type Option func(s *Server) error

// By default, the app will use the store specified by the configuration. This allows you to
// construct an app with a different store.
//
// The override parameter must be either a store.Store or func(App) store.Store().
func StoreOverride(override interface{}) Option {
	return func(s *Server) error {
		switch o := override.(type) {
		case store.Store:
			s.newStore = func() store.Store {
				return o
			}
			return nil

		case func(*Server) store.Store:
			s.newStore = func() store.Store {
				return o(s)
			}
			return nil

		default:
			return errors.New("invalid StoreOverride")
		}
	}
}

func SetConfig(path string) Option {
	return func(s *Server) error {
		s.configPath = path
		return nil
	}
}
