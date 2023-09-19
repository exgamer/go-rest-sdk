package structures

import (
	"reflect"
	"strings"
)

type AppConfig struct {
	Name          string `mapstructure:"APP_NAME" json:"app_name"`
	HostName      string `mapstructure:"HOST_NAME" json:"host_name"`
	Version       string `mapstructure:"APP_VERSION" json:"app_version"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS" json:"server_address"`
	SentryDsn     string `mapstructure:"SENTRY_DSN"    json:"sentry_dsn"`
	AppEnv        string `mapstructure:"APP_ENV"    json:"app_env"`
}

func (a AppConfig) GetFieldsAsJsonTags() []string {
	result := make([]string, 0)

	val := reflect.ValueOf(a)
	t := val.Type()

	for i := 0; i < t.NumField(); i++ {
		result = append(result, t.Field(i).Tag.Get("json"))
	}

	return result
}

func (a AppConfig) GetFieldsAsUpperSnake() []string {
	result := make([]string, 0)

	for _, v := range a.GetFieldsAsJsonTags() {
		result = append(result, strings.ToUpper(v))
	}

	return result
}
