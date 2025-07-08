package aws

import (
	"github.com/IntelXLabs-LLC/go-common/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// AWS is the main structure that keeps the session.Session from aws-sdk-go
//
// This is used throughout.
//
// Usage: aws.New().NewS3Operation(opts ...S3Options)
type AWS struct {
	sess   *session.Session
	region string
}

// New will create a new AWS packages
func New(region string) *AWS {
	a := new(AWS)
	if region != "" {
		a.region = region
		return a.ConnectAws()
	}
	a.region = "us-east-1"
	return a.ConnectAws()
}

// Raw returns the raw session from aws.
func (a *AWS) Raw() *session.Session {
	return a.sess
}

// ConnectAws establishes an AWS session by creating a new session with the given AWS region.
// It uses the credentials from the GetAWSConfig function if it is available.
// If there is an error during session creation, a panic is raised.
// It returns the created AWS session.
func (a *AWS) ConnectAws() *AWS {
	var err error
	cfg := &aws.Config{
		Region: aws.String(a.region),
	}
	if config.GetAWSConfig() != nil {
		cfg.Credentials = credentials.NewStaticCredentials(config.GetAWSConfig().AccessKeyID, config.GetAWSConfig().SecretAccessKey, "")
	}
	a.sess, err = session.NewSession(cfg)
	if err != nil {
		panic(err)
	}
	return a
}
