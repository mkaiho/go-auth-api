package infrastructure

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/mkaiho/go-auth-api/adapter/rdb"

	_ "github.com/go-sql-driver/mysql"
)

// MySQL Configuration
var _ rdb.Config = (*MySQLConfig)(nil)

func LoadMySQLConfig() (*MySQLConfig, error) {
	var c MySQLConfig
	if err := envconfig.Process("MYSQL", &c); err != nil {
		return nil, err
	}
	return &c, nil
}

type MySQLConfig struct {
	Host     string `envconfig:"HOST" required:"true"`
	User     string `envconfig:"USER" required:"true"`
	Database string `envconfig:"DATABASE" required:"true"`
	Password string `envconfig:"PASSWORD" required:"true"`
	Port     int    `envconfig:"PORT" default:"3306"`
	MaxConns int    `envconfig:"MAX_CONNS" default:"1"`
}

func (c *MySQLConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4", c.User, c.Password, c.Host, c.Port, c.Database)
}

func (c *MySQLConfig) GetMaxConns() int {
	return c.MaxConns
}
