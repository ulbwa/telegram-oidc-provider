package model

import (
	"database/sql"
	"time"
)

// Bot represents a Telegram bot in the database.
type Bot struct {
	Id        int64        `gorm:"column:id;primaryKey"`
	Name      string       `gorm:"column:name;type:varchar(255);not null"`
	ClientId  *string      `gorm:"column:client_id;type:varchar(255);uniqueIndex"`
	Username  string       `gorm:"column:username;type:varchar(255);not null"`
	Token     []byte       `gorm:"column:token;type:bytea;not null"`
	CreatedAt time.Time    `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt sql.NullTime `gorm:"column:updated_at"`
}

func (Bot) TableName() string { return "bots" }
