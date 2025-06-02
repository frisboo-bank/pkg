package postgres

type PgConfig struct {
	DBName   string `mapstructure:"dbName"`
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     string `mapstructure:"port"`
	SSLMode  bool   `mapstructure:"sslMode"`
	User     string `mapstructure:"user"`
}
