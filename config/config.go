package config

// Environment represents the environment in which the software is running.
type Environment int

// ENV represents the current environment the software is running in.
// It is of type Environment.
var ENV Environment

const (
	Development Environment = iota
	Staging     Environment = iota
	UAT         Environment = iota
	Production  Environment = iota
)

// DefaultSettings represents the default settings for the software.
// It contains the environment and AWS configuration.
type DefaultSettings struct {
	Environment   `config:"ENVIRONMENT"`
	AWSConfig     `prefix:"AWS"`
	UserSecret    string `json:"user_secret" config:"USER_SECRET"`
	AdminSecret   string `json:"admin_secret" config:"ADMIN_SECRET"`
	ServiceSecret string `json:"service_secret" config:"SERVICE_SECRET"`
}

// GetAdminSecret returns the admin secret
func GetAdminSecret() string {
	return defaultSettings.AdminSecret
}

// GetUserSecret returns the user Secret
func GetUserSecret() string {
	return defaultSettings.UserSecret
}

// GetServiceSecret returns the service secret
func GetServiceSecret() string {
	return defaultSettings.ServiceSecret
}
