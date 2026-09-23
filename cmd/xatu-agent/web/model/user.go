package model

import (
	"agentFrame/cmd/xatu-agent/utils/client"
	"database/sql"
	"log/slog"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Nickname string
	Email    string
	Status   int    `gorm:"default:1"`    //1正常 0禁用
	Role     string `gorm:"default:user"` // user / admin
}

type MySQLUser struct {
	sqlDB  *sql.DB
	gormDB *gorm.DB
}

func NewMySQLUser() (*MySQLUser, error) {

	config := client.NewMySQLClientConfig(os.Getenv("MYSQL_USER_DBNAME"))
	dsn := config.GetDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("open db failed", "error", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("get sql db failed", "error", err)
		return nil, err
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	if err := sqlDB.Ping(); err != nil {
		slog.Error("ping db failed", "error", err)
		_ = sqlDB.Close()
		return nil, err
	}

	return &MySQLUser{sqlDB: sqlDB, gormDB: db}, nil
}

func (u *MySQLUser) GetUserConnection() *gorm.DB {
	return u.gormDB
}
