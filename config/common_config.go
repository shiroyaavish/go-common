package config

// CommonConfig represents the configuration
type CommonConfig struct{}

var commonConfig CommonConfig

// UseDefaultConfig sets the default values for the common configuration.
func UseDefaultConfig() {
	commonConfig = CommonConfig{}

}

// GetCommonConfig returns a pointer to the `commonConfig` variable.
func GetCommonConfig() *CommonConfig {
	return &commonConfig
}

// SetCommonConfig sets the values of the commonConfig variable based on the provided config parameter.
// It assigns the following values to the commonConfig variable:
func SetCommonConfig(config *CommonConfig) {
	commonConfig = *config
}

// AWSConfig represents the configuration for AWS service
type AWSConfig struct {
	SecretAccessKey string `json:"secret_access_key" config:"SECRET_ACCESS_KEY"`
	AccessKeyID     string `json:"access_key_id" config:"ACCESS_KEY_ID"`
	Region          string `json:"region" config:"REGION"`
}

// GetAWSConfig returns a pointer to the default AWSConfig object.
// AWSConfig is a struct that contains the SecretAccessKey, AccessKeyID, and Region.
// It is accessed from the defaultSettings package variable.
// Example usage:
//
//	cfg := &aws.Config{
//	    Region: aws.String("us-west-2"),
//	}
//	if config.GetAWSConfig() != nil {
//	    cfg.Credentials = credentials.NewStaticCredentials(config.GetAWSConfig().AccessKeyID, config.GetAWSConfig().SecretAccessKey, "")
//	}
//	sess, err = session.NewSession(cfg)
//	if err != nil {
//	    panic(err)
//	}
//	a.sess = sess
//	return a
func GetAWSConfig() *AWSConfig {
	return &defaultSettings.AWSConfig
}

// SetAWSConfig sets the AWS configuration by replacing the default settings with the provided configuration.
// The `config` parameter is a pointer to an `AWSConfig` struct which contains the following properties:
// - `SecretAccessKey`: The secret access key for AWS authentication.
// - `AccessKeyID`: The access key ID for AWS authentication.
// - `Region`: The AWS region.
func SetAWSConfig(config *AWSConfig) {
	defaultSettings.AWSConfig = *config
}
