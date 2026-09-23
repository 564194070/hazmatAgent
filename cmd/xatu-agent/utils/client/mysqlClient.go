package client

import (
	"fmt"
	"os"
)

type MySQLClientConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
	Charset  string
}

func NewMySQLClientConfig(dbName string) *MySQLClientConfig {
	charset := os.Getenv("MYSQL_CHARSET")
	if charset == "" {
		charset = "utf8mb4"
	}
	port := os.Getenv("MYSQL_PORT")
	if port == "" {
		port = "3306"
	}

	return &MySQLClientConfig{
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     port,
		DBName:   dbName,
		Charset:  charset,
	}
}

func (c *MySQLClientConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.Charset,
	)
}
