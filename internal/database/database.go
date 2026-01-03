package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DannyAss/users/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBManager struct {
	connections map[string]*gorm.DB
}

func NewDBManager(cfg config.DBManager) *DBManager {
	conns := map[string]*gorm.DB{
		"master": InitDBConnection(&cfg.DBMaster),
	}

	return &DBManager{connections: conns}
}

func (m *DBManager) GetDB(alias ...string) *gorm.DB {
	if len(alias) > 0 {
		if db, ok := m.connections[alias[0]]; ok {
			return db
		}
	}
	return m.connections["master"]
}

func InitDBConnection(cfg *config.DBConnection) *gorm.DB {
	dialector := GetDBDialect(cfg)

	// Open connection to database
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	// Pooling configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConnection)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConnection)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBConnMaxLifetime) * time.Minute)

	// Health check to ensure the connection is alive
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil
	}

	log.Printf("Database %s connected successfully\n", cfg.DBName)

	return db
}

func GetDBDialect(cfg *config.DBConnection) gorm.Dialector {
	log.Println("cek dialect : ", cfg.DBDialect)
	switch cfg.DBDialect {
	case "mysql":
		return mysql.Open(MySQLConnectionString(cfg))
	default:
		log.Fatalf("Unsupported database dialect: %s", cfg.DBDialect)
		return nil
	}
}

func MySQLConnectionString(cfg *config.DBConnection) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBHostWrite,
		cfg.DBPort,
		cfg.DBName,
	)
}

func (m *DBManager) Close() {
	for name, conn := range m.connections {
		if sqlDB, err := conn.DB(); err == nil {
			_ = sqlDB.Close()
			log.Printf("Closed DB connection: %s\n", name)
		}
	}
}
