package structures

import (
	"reflect"
	"strings"
)

type DbConfig struct {
	Host                     string `mapstructure:"DB_HOST" json:"db_host"`
	Port                     string `mapstructure:"DB_PORT" json:"db_port"`
	Db                       string `mapstructure:"DB_DATABASE" json:"db_database"`
	Username                 string `mapstructure:"DB_USERNAME" json:"db_username"`
	Password                 string `mapstructure:"DB_PASSWORD" json:"db_password"`
	MaxPoolConnections       int    `mapstructure:"DB_MAX_POOL_CONNECTIONS" json:"db_max_pool_connections"`
	MaxIdlePoolConnections   int    `mapstructure:"DB_MAX_IDLE_POOL_CONNECTIONS" json:"db_max_idle_connections"`
	ConnectionTimeoutSeconds int64  `mapstructure:"DB_CONNECTION_TIMEOUT_SECONDS" json:"db_connection_timeout_seconds"`
}

func (a DbConfig) GetFieldsAsJsonTags() []string {
	result := make([]string, 0)

	val := reflect.ValueOf(a)
	t := val.Type()

	for i := 0; i < t.NumField(); i++ {
		result = append(result, t.Field(i).Tag.Get("json"))
	}

	return result
}

func (a DbConfig) GetFieldsAsUpperSnake() []string {
	result := make([]string, 0)

	for _, v := range a.GetFieldsAsJsonTags() {
		v = strings.ToUpper(v)
	}

	return result
}
