package structures

type AppConfig struct {
	Name          string `mapstructure:"APP_NAME" json:"app_name"`
	HostName      string `mapstructure:"HOST_NAME" json:"host_name"`
	Version       string `mapstructure:"APP_VERSION" json:"app_version"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS" json:"server_address"`
}
