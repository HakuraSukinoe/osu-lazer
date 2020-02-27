package sqlstore

import (
	"github.com/deissh/osu-lazer/server/model"
	"github.com/deissh/osu-lazer/server/store"
)

type SqlOauthStore struct {
	SqlStore
}

func (s SqlOauthStore) CreateToken(userID uint, clientID uint, clientSecret string, scopes string) (*model.Token, error) {
	panic("implement me")
}

func (s SqlOauthStore) ValidateToken(accessToken string) (*model.Token, error) {
	panic("implement me")
}

func (s SqlOauthStore) RefreshOAuthToken(refreshToken string, clientID uint, clientSecret string, scopes string) (*model.Token, error) {
	panic("implement me")
}

func (s SqlOauthStore) CreateClient(userID uint, name string, redirect string) (*model.Client, error) {
	panic("implement me")
}

func (s SqlOauthStore) FindClient(clientID uint, clientSecret string) (*model.Client, error) {
	panic("implement me")
}

func NewSqlOAuthStore(sqlStore SqlStore) store.OAuthStore {
	return &SqlOauthStore{sqlStore}
}
