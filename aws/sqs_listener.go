package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"time"
)

// SQSListener represents a listener for receiving messages from an SQS queue.
// It is associated with a specific queue URL and uses the provided SQS client for communicating with SQS service.
// SQSListener can be created using the NewListener function.
type SQSListener struct {
	sqs       *sqs.SQS
	QueueUrl  string
	closeChan chan int
}

// NewListener returns a new instance of SQSListener,
// initialized with the provided queue URL.
// It creates a new AWS session in the us-west-2 region and
// configures the SQS service with that session.
// Returns a pointer to the created SQSListener and any error encountered.
func NewListener(queueUrl string) (*SQSListener, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	if err != nil {
		return nil, err
	}

	return &SQSListener{
		sqs:      sqs.New(sess),
		QueueUrl: queueUrl,
	}, nil
}

// ListenForMessages continuously listens for messages from the SQS queue
// specified in the SQSListener struct. It takes a callback function as a parameter,
// which is executed for each received message.
// The callback function should accept a pointer to an SQS message as a parameter
// and return an error if any. If the callback function returns an error, ListenForMessages
// stops listening and returns that error.
// If there is an error while receiving messages from the SQS queue,
// ListenForMessages returns that error.
//
// Example:
//
//	err := listener.ListenForMessages(func(message *sqs.Message) error {
//	   // Process the received message here
//	   return nil
//	})
//
//	if err != nil {
//	    // Handle the error
//	}
func (s *SQSListener) ListenForMessages(callback func(*sqs.Message) error) {
	go func() {
		for {
			if len(s.closeChan) == 1 {
				return
			}
			result, err := s.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(s.QueueUrl),
				MaxNumberOfMessages: aws.Int64(1),
			})

			if err != nil {
				panic(err)
			}

			if len(result.Messages) == 0 {
				time.Sleep(time.Second * 2)
				continue
			}

			for _, message := range result.Messages {
				if err := callback(message); err != nil {
					panic(err)
				}
			}
		}
	}()
}

// Close the listener
func (s *SQSListener) Close() error {
	s.closeChan <- 1
	return nil
}
