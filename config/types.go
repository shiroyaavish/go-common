package config

// DatabaseConfig represents the configuration for connecting to a database.
// It contains the following fields:
// - Host: the hostname or IP address of the database server
// - Port: the port number of the database server
// - User: the username for authenticating to the database
// - Password: the password for authenticating to the database
// - DBName: the name of the database to connect to
// - SSLMode: the SSL mode for the database connection
type DatabaseConfig struct {
	Host     string `json:"host" validate:"required" config:"HOST"`
	Port     string `json:"port" validate:"required" config:"PORT"`
	User     string `json:"user" validate:"required" config:"USER"`
	Password string `json:"password" validate:"required" config:"PASSWORD"`
	DBName   string `json:"db_name" validate:"required" config:"NAME"`
	SSLMode  string `json:"ssl_mode" validate:"required" config:"SSL_MODE"`
}

// GRPCConfig represents the configuration for a GRPC server.
// It contains the following field:
// - Port: the port number on which the GRPC server will listen for incoming connections.
type GRPCConfig struct {
	Port int `json:"port" validate:"required" config:"PORT"`
}

// RedisConfig represents the configuration for connecting to a Redis server.
// It contains the following fields:
// - Host: the hostname or IP address of the Redis server
// - Password: the password for authenticating to the Redis server
// - DBNumber: the database number to use for the Redis connection
type RedisConfig struct {
	Host     string `json:"host" validate:"required" config:"HOST"`
	Password string `json:"password" validate:"required" config:"PASSWORD"`
	DBNumber int    `json:"db_number" validate:"required" config:"DB_NUMBER"`
}

// StripeConfig represents the configuration for interacting with the Stripe API.
// It contains the following field:
// - SecretKey: the secret key for authenticating with the Stripe API
type StripeConfig struct {
	SecretKey string `json:"secret_key" config:"SECRET_KEY" validate:"required"`
}

// WrapperConfig contains the config for a GRPC wrapper
type WrapperConfig struct {
	TimeoutSec int    `config:"TIMEOUT_SEC"`
	GrpcUrl    string `config:"GRPC_URL"`
}

// ServerConfig is to run fiber.App, all the other config variables required will be present in data.
type ServerConfig struct {
	Port        string `json:"port" validate:"required" config:"PORT"`
	Version     string `json:"version" validate:"required" config:"VERSION"`
	ProjectName string `json:"project_name" validate:"required" config:"PROJECT_NAME"`
}

// CashFreeConfig is the overall configuration for the application
type CashFreeConfig struct {
	APPId  string `json:"app_id" validate:"required" config:"APP_ID"`
	APIKey string `json:"api_key" validate:"required" config:"API_KEY"`
}

type RabbitMQConfig struct {
	URL      string `json:"url" validate:"required" config:"URL"`
	Username string `json:"username" validate:"required" config:"USERNAME"`
	Password string `json:"password" validate:"required" config:"PASSWORD"`
	Host     string `json:"host" validate:"required" config:"HOST"`
	Port     string `json:"port" validate:"required" config:"PORT"`
}
