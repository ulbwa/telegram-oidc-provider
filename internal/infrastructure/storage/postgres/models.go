package postgres

import (
	"database/sql"
	"net"
	"net/url"
	"time"
)

// UserModel represents the users table in the database
type UserModel struct {
	Id        int64          `gorm:"primaryKey;column:id"`
	FirstName string         `gorm:"column:first_name"`
	LastName  sql.NullString `gorm:"column:last_name"`
	Username  sql.NullString `gorm:"column:username"`
	PhotoUrl  sql.NullString `gorm:"column:photo_url"`
	IsPremium sql.NullBool   `gorm:"column:is_premium"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt sql.NullTime   `gorm:"column:updated_at"`
}

// TableName overrides the table name used by UserModel
func (UserModel) TableName() string {
	return "users"
}

// ToPhotoURL converts the string PhotoUrl to *url.URL
func (u *UserModel) ToPhotoURL() (*url.URL, error) {
	if !u.PhotoUrl.Valid {
		return nil, nil
	}
	return url.Parse(u.PhotoUrl.String)
}

// FromPhotoURL sets the PhotoUrl from *url.URL
func (u *UserModel) FromPhotoURL(photoUrl *url.URL) {
	if photoUrl == nil {
		u.PhotoUrl = sql.NullString{Valid: false}
		return
	}
	u.PhotoUrl = sql.NullString{String: photoUrl.String(), Valid: true}
}

// BotModel represents the bots table in the database
type BotModel struct {
	Id        int64        `gorm:"primaryKey;column:id"`
	Name      string       `gorm:"column:name"`
	ClientId  string       `gorm:"column:client_id"`
	Username  string       `gorm:"column:username"`
	Token     string       `gorm:"column:token"`
	CreatedAt time.Time    `gorm:"column:created_at"`
	UpdatedAt sql.NullTime `gorm:"column:updated_at"`
}

// TableName overrides the table name used by BotModel
func (BotModel) TableName() string {
	return "bots"
}

// UserBotLoginModel represents the user_bot_logins table in the database
// This is a junction table with additional attributes
type UserBotLoginModel struct {
	UserId      int64          `gorm:"primaryKey;column:user_id"`
	BotId       int64          `gorm:"primaryKey;column:bot_id"`
	IP          string         `gorm:"column:ip"`
	UserAgent   sql.NullString `gorm:"column:user_agent"`
	Language    sql.NullString `gorm:"column:language"`
	LastLoginAt time.Time      `gorm:"column:last_login_at"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   sql.NullTime   `gorm:"column:updated_at"`
}

// TableName overrides the table name used by UserBotLoginModel
func (UserBotLoginModel) TableName() string {
	return "user_bot_logins"
}

// ToIP converts the string IP to net.IP
func (ubl *UserBotLoginModel) ToIP() net.IP {
	return net.ParseIP(ubl.IP)
}

// FromIP sets the IP from net.IP
func (ubl *UserBotLoginModel) FromIP(ip net.IP) {
	ubl.IP = ip.String()
}
