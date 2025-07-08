package aws

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

// SQS object provides sqs.SQS api
type SQS struct {
	sqsSvc *sqs.SQS
	aws    *AWS
}

// NewSQSOperation creates a new SQS operation
func (a *AWS) NewSQSOperation() *SQS {
	if a.sess == nil {
		a.ConnectAws()
	}
	return &SQS{
		sqsSvc: sqs.New(a.sess),
	}
}

// Trigger will trigger any sqs url
//
//   - url: takes the URL as string
//   - data: takes []bytes
func (s *SQS) Trigger(url string, data []byte) (*string, error) {
	if s.sqsSvc == nil {
		s.sqsSvc = sqs.New(s.aws.sess)
	}
	dataString := string(data)
	message := sqs.SendMessageInput{
		MessageBody: &dataString,
		QueueUrl:    &url,
	}
	response, err := s.sqsSvc.SendMessage(&message)
	if err != nil {
		log.Println("ERROR | ", err)
		return nil, err
	}
	return response.MessageId, nil
}
