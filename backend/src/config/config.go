package config

import (
	"fmt"
	"os"
)

// Config 全局配置分散从 .env.example / docker-compose.yml / 环境变量读取，新增配置需同步多处。
type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	JWTSecret  string
	UseMySQL   bool
	SQLitePath string
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Load 容器内 DB_HOST=db 时走 MySQL；本地无 DB_HOST 时使用纯 Go SQLite 文件（零外部依赖）。
func Load() *Config {
	cfg := &Config{
		Port:       getenv("PORT", "3000"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     getenv("DB_PORT", "3306"),
		DBName:     getenv("DB_NAME", "app_db"),
		DBUser:     getenv("DB_USER", "app_user"),
		DBPassword: getenv("DB_PASSWORD", "app_password"),
		JWTSecret:  getenv("JWT_SECRET", "local-dev-secret"),
		SQLitePath: getenv("SQLITE_PATH", "ground-turn.db"),
	}
	cfg.UseMySQL = cfg.DBHost != ""
	return cfg
}

// MySQLDSN GORM MySQL 连接串。
func (c *Config) MySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
