package config

type Config struct {
	MysqlConfig *MysqlConfig
	Port        string
}

type MysqlConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func (config *MysqlConfig) DSN() string {
	return config.Username + ":" + config.Password + "@tcp(" + config.Host + ":" + config.Port + ")/" + config.Database
}

func loadDefaultConfig() *Config {
	return &Config{
		MysqlConfig: &MysqlConfig{
			Host:     "localhost",
			Port:     "3306",
			Username: "root",
			Password: "secret",
			Database: "flashcard",
		},
		Port: "8080",
	}
}
