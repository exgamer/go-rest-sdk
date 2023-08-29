package structures

import (
	"reflect"
)

type AppConfig struct {
	Name          string `mapstructure:"APP_NAME" json:"app_name"`
	HostName      string `mapstructure:"HOST_NAME" json:"host_name"`
	Version       string `mapstructure:"APP_VERSION" json:"app_version"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS" json:"server_address"`
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
