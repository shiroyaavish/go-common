package aws

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"time"
)

type S3 struct {
	ACLOperation      bool
	ACLName           string
	Metadata          map[string]*string
	MetadataOperation bool
	ContentType       string
	aws               *AWS
}

type S3Options interface {
	GetS3PermBool(s3 *S3)
}

func (a *AWS) NewS3Operation(opts ...S3Options) *S3 {
	thisS3Manager := &S3{
		aws: a,
	}
	for _, opt := range opts {
		opt.GetS3PermBool(thisS3Manager)
	}
	return thisS3Manager
}

func (s *S3) PutObjectInBucket(bucketName string, key string, object []byte) (string, error) {
	if s.aws.sess == nil {
		s.aws.ConnectAws()
	}
	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(s.aws.sess)
	// Get the bucket name from the config

	ip := &s3manager.UploadInput{
		Bucket: &bucketName,
		Key:    aws.String(key),
		Body:   bytes.NewReader(object),
	}

	if s.ACLName != "" && s.ACLOperation {
		ip.ACL = &s.ACLName
	}

	if s.Metadata != nil && s.MetadataOperation {
		ip.Metadata = s.Metadata
	}

	if s.ContentType != "" {
		ip.ContentType = &s.ContentType
	}
	// Upload the file to S3.
	_, err := uploader.Upload(ip)
	// If there is an error, return it
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	// Return the key
	return key, nil
}

func (s *S3) GetObjectFromBucket(bucketName string, key string) ([]byte, *int64, error) {
	if s.aws.sess == nil {
		s.aws.ConnectAws()
	}

	downloader := s3manager.NewDownloader(s.aws.sess)

	var buff aws.WriteAtBuffer

	data, err := downloader.Download(&buff, &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})
	if err != nil {
		return nil, nil, err
	}

	return buff.Bytes(), &data, nil
}

// GeneratePreSignedS3URL generates a pre-signed URL for accessing an S3 object.
// The function takes the key and bucket name as parameters.
// It first checks if the sess variable is nil and if so, calls the ConnectAws function to establish an AWS session.
// Then, it creates an S3 service client using the session.
// Next, it creates a request object for getting the S3 object using the specified key and bucket.
// The expiration duration for the pre-signed URL is set to 8 hours.
// The req object has the Presign method called on it, passing in the expiration duration to generate the pre-signed URL.
// If there's an error during the pre-signing process, it returns an empty string and the error.
// Otherwise, it returns the pre-signed URL and nil for the error.
func (s *S3) GeneratePreSignedS3URL(key string, bucket string) (string, error) {
	if s.aws == nil {
		s.aws.ConnectAws()
	}
	svc := s3.New(s.aws.sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	expireDuration := 8 * time.Hour
	url, err := req.Presign(expireDuration)
	if err != nil {
		return "", err
	}
	return url, nil
}
