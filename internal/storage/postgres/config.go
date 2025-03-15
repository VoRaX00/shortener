package postgres

type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"dbname"`
	SSLMode  string `yaml:"ssl_mode"`
}
