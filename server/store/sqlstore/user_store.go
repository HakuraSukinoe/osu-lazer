package sqlstore

import (
	"github.com/deissh/osu-lazer/server/model"
	"github.com/deissh/osu-lazer/server/store"
	"github.com/deissh/osu-lazer/server/utils"
)

var modes = []string{"std", "mania", "catch", "taiko"}

type SqlUserStore struct {
	SqlStore
}

func NewSqlUserStore(sqlStore SqlStore) store.UserStore {
	return &SqlUserStore{sqlStore}
}

func (s SqlUserStore) Get(id uint, mode string) (*model.UserFull, error) {
	var user model.UserFull

	if !utils.ContainsString(modes, mode) {
		mode = "std"
	}
	user.Mode = mode

	err := s.GetMaster().Get(
		&user,
		`SELECT u.*, check_online(last_visit),
				json_build_object('code', c.code, 'name', c.name) as country
			FROM users u
			INNER JOIN countries c ON c.code = u.country_code
			WHERE u.id = $1`,
		id,
	)

	return &user, err
}

func (s SqlUserStore) GetByLoginAndPass(username string, password string) (*model.UserFull, error) {
	panic("implement me")
}

func (s SqlUserStore) Create(username string, email string, password string) (*model.UserFull, error) {
	panic("implement me")
}

func (s SqlUserStore) UpdateLastVisit(userId uint) (*model.UserFull, error) {
	panic("implement me")
}

func (s SqlUserStore) GetAll() ([]*model.User, error) {
	panic("implement me")
}
