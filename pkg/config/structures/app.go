package structures

type AppConfig struct {
	Name          string `mapstructure:"APP_NAME"`
	HostName      string `mapstructure:"HOST_NAME"`
	Version       string `mapstructure:"APP_VERSION"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}
