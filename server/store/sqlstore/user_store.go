package sqlstore

import (
	"github.com/deissh/osu-lazer/server/einterfaces"
	"github.com/deissh/osu-lazer/server/model"
	"github.com/deissh/osu-lazer/server/store"
)

type SqlUserStore struct {
	SqlStore

	metrics einterfaces.MetricsInterface
}

func NewSqlUserStore(sqlStore SqlStore, metrics einterfaces.MetricsInterface) store.UserStore {
	us := &SqlUserStore{
		SqlStore: sqlStore,
		metrics:  metrics,
	}

	return us
}

func (us SqlUserStore) Get(id string) (*model.User, *error) {
	return &model.User{}, nil
}

func (us SqlUserStore) GetAll() ([]*model.User, *error) {
	return nil, nil
}
