//go:generate go run layer_generators/main.go

package store

import "github.com/deissh/osu-lazer/server/model"

type Store interface {
	User() UserStore

	Close()
}

type UserStore interface {
	Get(id string) (*model.User, *error)
	GetAll() ([]*model.User, *error)
}
