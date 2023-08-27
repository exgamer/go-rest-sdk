package structures

type DbConfig struct {
	Host                     string `mapstructure:"DB_HOST"`
	Port                     string `mapstructure:"DB_PORT"`
	Db                       string `mapstructure:"DB_DATABASE"`
	Username                 string `mapstructure:"DB_USERNAME"`
	Password                 string `mapstructure:"DB_PASSWORD"`
	MaxPoolConnections       int    `mapstructure:"DB_MAX_POOL_CONNECTIONS"`
	MaxIdlePoolConnections   int    `mapstructure:"DB_MAX_IDLE_POOL_CONNECTIONS"`
	ConnectionTimeoutSeconds int64  `mapstructure:"DB_CONNECTION_TIMEOUT_SECONDS"`
}
