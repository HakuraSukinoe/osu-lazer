//go:generate go run layer_generators/main.go

package store

import (
	"github.com/deissh/osu-lazer/server/model"
)

type Store interface {
	User() UserStore
	Friend() FriendStore

	Close()
}

type UserStore interface {
	Get(id uint, mode string) (*model.UserFull, error)
	GetByLoginAndPass(username string, password string) (*model.UserFull, error)
	Create(username string, email string, password string) (*model.UserFull, error)
	UpdateLastVisit(userId uint) (*model.User, error)
	GetAll() ([]*model.User, error)
}

type FriendStore interface {
	GetFriends(userId uint) ([]*model.User, error)
	SetFriend(userId uint, targetId uint) error
	RemoveFriend(userId uint, targetId uint) error
}
