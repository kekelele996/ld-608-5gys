package config

import (
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"groundTurn/src/models"
)

// OpenDB 根据配置选择 MySQL（容器）或 SQLite（本地零依赖）。
func OpenDB(cfg *Config) *gorm.DB {
	var dialector gorm.Dialector
	if cfg.UseMySQL {
		dialector = mysql.Open(cfg.MySQLDSN())
	} else {
		dialector = sqlite.Open(cfg.SQLitePath)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Warn),
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatalf("open db failed: %v", err)
	}

	if !cfg.UseMySQL {
		// SQLite 单写连接，配合事务避免 database is locked。
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.SetMaxOpenConns(1)
			sqlDB.SetMaxIdleConns(1)
		}
	} else if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxOpenConns(20)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	if err := db.AutoMigrate(
		&models.FlightTurnaround{},
		&models.GroundTask{},
		&models.GroundResource{},
		&models.ResourceBooking{},
		&models.DelayEvent{},
		&models.AuditLog{},
	); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	return db
}
