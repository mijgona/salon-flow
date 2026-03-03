package cmd

// Config holds the application configuration.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Outbox   OutboxConfig   `yaml:"outbox"`
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Port int `yaml:"port"`
}

// DatabaseConfig holds database connection configuration.
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslmode"`
}

// OutboxConfig holds outbox job configuration.
type OutboxConfig struct {
	Interval  string `yaml:"interval"`
	BatchSize int    `yaml:"batch_size"`
}
