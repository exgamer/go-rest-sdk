package structures

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
