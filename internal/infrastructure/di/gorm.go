package di

import (
	"context"

	gormzerolog "github.com/mpalmer/gorm-zerolog"
	"github.com/samber/do/v2"
	"github.com/ulbwa/telegram-oidc-provider/internal/common"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// postgresConnection wraps gorm.DB and provides graceful shutdown
type postgresConnection struct {
	db *gorm.DB
}

// Shutdown implements graceful shutdown for database connection
func (p *postgresConnection) Shutdown(ctx context.Context) error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func provideGorm(injector do.Injector) {
	// Provide dialector
	do.Provide(injector, func(i do.Injector) (gorm.Dialector, error) {
		cfg, err := do.Invoke[*common.Config](i)
		if err != nil {
			return nil, err
		}

		return postgres.Open(cfg.Database.DSN), nil
	})

	// Provide wrapped connection with graceful shutdown support
	do.Provide(injector, func(i do.Injector) (*postgresConnection, error) {
		dialector, err := do.Invoke[gorm.Dialector](i)
		if err != nil {
			return nil, err
		}

		// Use gorm-zerolog which reads logger from context via zerolog.Ctx()
		db, err := gorm.Open(dialector, &gorm.Config{
			Logger: gormzerolog.Logger{},
		})
		if err != nil {
			return nil, err
		}

		return &postgresConnection{db: db}, nil
	})

	// Provide unwrapped *gorm.DB for convenience
	do.Provide(injector, func(i do.Injector) (*gorm.DB, error) {
		conn, err := do.Invoke[*postgresConnection](i)
		if err != nil {
			return nil, err
		}
		return conn.db, nil
	})
}
