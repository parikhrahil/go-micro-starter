package sqldb

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/parikhrahil/go-micro-starter/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SQLDB interface {
	DB() *gorm.DB
	Close() error
}

type database struct {
	db *gorm.DB
}

type Opts struct {
	DSN                 string        `validate:"required"`
	Driver              string        `validate:"required,oneof=postgres"`
	PoolMaxIdleConns    int           `validate:"gte=0"`
	PoolMaxOpenConns    int           `validate:"gte=0"`
	PoolConnMaxLifetime time.Duration `validate:"gte=0"`
	EnableDebug         bool
	Log                 logger.Logger
}

func New(opts *Opts) (SQLDB, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(opts); err != nil {
		return nil, err
	}

	dsn := opts.DSN
	driver := opts.Driver
	log := opts.Log

	var db *gorm.DB
	var err error

	switch driver {
	case "postgres":
		if dsn == "" {
			return nil, fmt.Errorf("invalid DSN: postgres requires a non empty DSN")
		}
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("Failed to connect to %s: %v", driver, err)
		}

		if opts.EnableDebug {
			db.Debug()
		}

	default:
		return nil, fmt.Errorf("Unsupported database driver %s", driver)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("Failed to get DB instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(opts.PoolMaxIdleConns)
	sqlDB.SetMaxOpenConns(opts.PoolMaxOpenConns)
	sqlDB.SetConnMaxLifetime(opts.PoolConnMaxLifetime)

	log.Info(fmt.Sprintf("%s connected at %s", driver, dsn))
	return &database{db: db}, nil
}

func (db *database) DB() *gorm.DB {
	return db.db
}

func (db *database) Close() error {
	if db == nil || db.db == nil {
		return fmt.Errorf("Cannot close: Repository not initialized")
	}
	sqlDB, err := db.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
