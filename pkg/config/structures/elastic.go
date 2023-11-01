package structures

type ElasticNoSecurityConfig struct {
	ElasticHost     string `mapstructure:"ELASTIC_HOST" json:"elastic_host"`
	ElasticShowLogs bool   `mapstructure:"ELASTIC_SHOW_LOGS" json:"elastic_show_logs"`
}
