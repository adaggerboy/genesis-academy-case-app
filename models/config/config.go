package config

type DatabaseEndpointConfig struct {
	Name     string `yaml:"name"`
	Driver   string `yaml:"driver"`
	Path     string `yaml:"path"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Database string `yaml:"database"`
	Password string `yaml:"password"`
	Port     uint16 `yaml:"port"`
	Timeout  int    `yaml:"timeout"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
	Port     uint16 `yaml:"port"`
}

type OpenExchangeAPIConfig struct {
	Token         string `yaml:"token"`
	CacheDuration int    `yaml:"cache_duration"`
}

type HTTPServerConfig struct {
	Endpoints []string `yaml:"endpoints"`
}

type Config struct {
	Database     DatabaseEndpointConfig `yaml:"database"`
	SMTPConfig   SMTPConfig             `yaml:"smtp"`
	OpenExchange OpenExchangeAPIConfig  `yaml:"openexchange"`
	HTTPServer   HTTPServerConfig       `yaml:"http"`
	CronString   string                 `yaml:"cron_sender_string"`
}
