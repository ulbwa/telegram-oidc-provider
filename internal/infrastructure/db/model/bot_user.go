package model

import (
	"database/sql"
	"net/netip"
	"time"
)

// BotUser represents a relationship between a bot and a Telegram user.
type BotUser struct {
	BotId       int64          `gorm:"column:bot_id;primaryKey;not null"`
	UserId      int64          `gorm:"column:user_id;primaryKey;not null"`
	FirstName   string         `gorm:"column:first_name;type:varchar(255);not null"`
	LastName    sql.NullString `gorm:"column:last_name;type:varchar(255)"`
	Username    sql.NullString `gorm:"column:username;type:varchar(255)"`
	PhotoUrl    sql.NullString `gorm:"column:photo_url;type:text"`
	IsPremium   sql.NullBool   `gorm:"column:is_premium"`
	IP          netip.Addr     `gorm:"column:ip;type:inet;not null"`
	UserAgent   sql.NullString `gorm:"column:user_agent;type:text"`
	Language    sql.NullString `gorm:"column:language;type:varchar(10)"`
	LastLoginAt time.Time      `gorm:"column:last_login_at;not null;default:CURRENT_TIMESTAMP"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   sql.NullTime   `gorm:"column:updated_at"`
}

func (BotUser) TableName() string { return "bot_users" }
