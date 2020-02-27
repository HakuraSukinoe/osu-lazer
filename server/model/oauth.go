package model

import "time"

// Token model
type Token struct {
	ID           uint   `db:"id" json:"id"`
	UserID       uint   `db:"user_id" json:"user_id"`
	ClientID     uint   `db:"client_id" json:"client_id"`
	AccessToken  string `db:"access_token" json:"access_token"`
	RefreshToken string `db:"refresh_token" json:"refresh_token"`
	Scopes       string `db:"scopes" json:"scopes"`
	Revoked      bool   `db:"revoked" json:"revoked" json:"revoked"`
	ExpiresIn    int    `json:"expires_in"`

	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Client model
type Client struct {
	ID        uint      `db:"id" json:"id"`
	UserID    uint      `db:"user_id" json:"user_id"`
	Name      string    `db:"name" json:"name"`
	Secret    string    `db:"secret" json:"secret"`
	Redirect  string    `db:"redirect" json:"redirect"`
	Revoked   string    `db:"revoked" json:"revoked"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
