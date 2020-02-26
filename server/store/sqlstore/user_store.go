package sqlstore

import (
	"fmt"
	"github.com/deissh/osu-lazer/server/model"
	"github.com/deissh/osu-lazer/server/store"
)

type SqlUserStore struct {
	SqlStore
}

func NewSqlUserStore(sqlStore SqlStore) store.UserStore {
	us := &SqlUserStore{
		SqlStore: sqlStore,
	}

	return us
}

func (us SqlUserStore) Get(id string) (*model.User, *error) {
	return &model.User{}, nil
}

func (us SqlUserStore) GetAll() ([]*model.User, *error) {
	fmt.Print("here")
	return nil, nil
}
