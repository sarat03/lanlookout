package gdb

import (
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/glebarez/sqlite"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/sarat03/lanlookout/internal/check"
	"github.com/sarat03/lanlookout/internal/conf"
	"github.com/sarat03/lanlookout/internal/models"
)

var db *gorm.DB
var gormConf *gorm.Config

// Start working with DB
func Start() {
	var tab *gorm.DB
	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             5 * time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
	gormConf = &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			NoLowerCase: true,
			// So upper case Columns could work in both PostgreSQL and SQLite
		},
	}

	Connect()

	// Migrate the schema
	tab = db.Table("now")
	err = tab.AutoMigrate(&models.Host{})
	check.IfError(err)

	// Owner is never empty: older rows become Shared
	err = db.Table("now").Where(`"OWNER" IS NULL OR "OWNER" = ''`).Update("OWNER", models.SharedOwner).Error
	check.IfError(err)

	tab = db.Table("history")
	err = tab.AutoMigrate(&models.Host{})
	check.IfError(err)

	tab = db.Table("events")
	err = tab.AutoMigrate(&models.Event{})
	check.IfError(err)

	tab = db.Table("netstats")
	err = tab.AutoMigrate(&models.NetStat{})
	check.IfError(err)

	tab = db.Table("dnsstats")
	err = tab.AutoMigrate(&models.DNSStat{})
	check.IfError(err)
}

// Connect - choose DB and connect
func Connect() {
	var err error
	var pgFail bool

	if conf.AppConfig.UseDB == "postgres" {
		db, err = gorm.Open(postgres.Open(conf.AppConfig.PGConnect), gormConf)

		if err != nil {
			pgFail = true

			slog.Error("PostgreSQL connection error:", "err", err)
			slog.Warn("Falling back to SQLite")
		} else {
			slog.Info("Connected to DB: PostgreSQL")
		}
	}

	if pgFail || conf.AppConfig.UseDB != "postgres" {

		db, err = gorm.Open(sqlite.Open(conf.AppConfig.DBPath), gormConf)

		if !check.IfError(err) {
			slog.Info("Connected to DB: SQLite")
			db.Exec("PRAGMA journal_mode = wal;")
			db.Exec("PRAGMA busy_timeout = 5000;")
		}
	}
}
