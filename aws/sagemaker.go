package aws

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/sagemakerruntime"
)

// Sagemaker is a type that represents an interface for interacting with the Amazon SageMaker runtime service.
// It encapsulates the SageMakerRuntime client from the sagemakerruntime package.
type Sagemaker struct {
	sageMakerSvc *sagemakerruntime.SageMakerRuntime
}

// NewSagemakerOperation creates a new instance of the Sagemaker struct which allows invoking an endpoint using the SagemakerRuntime service. It ensures that the AWS session is connected
func (a *AWS) NewSagemakerOperation() *Sagemaker {
	if a.sess == nil {
		a.ConnectAws()
	}
	return &Sagemaker{
		sageMakerSvc: sagemakerruntime.New(a.sess),
	}
}

// Invoke sends a request to invoke an endpoint in Sagemaker, passing the given `endpointName`, `body`, and `respTarget`.
// It marshals the `body` to JSON, creates an `InvokeEndpointInput` object with the specified `Accept`, `Body`, `ContentType`, and `EndpointName`.
// Then, it calls the `InvokeEndpoint` method of the `SageMakerRuntime` service in AWS SageMaker.
// If there's an error in marshaling the `body`, it returns the error.
// If there's an error in calling the `InvokeEndpoint` method, it returns the error.
// Finally, it unmarshals the response body into the `respTarget`.
// It returns any error that occurred during the process.
func (s *Sagemaker) Invoke(endpointName string, body interface{}, respTarget interface{}) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req := &sagemakerruntime.InvokeEndpointInput{
		Accept:       pointer("application/json"),
		Body:         bodyBytes,
		ContentType:  pointer("application/json"),
		EndpointName: pointer(endpointName),
	}

	response, err := s.sageMakerSvc.InvokeEndpoint(req)
	if err != nil {
		return err
	}

	return json.Unmarshal(response.Body, &respTarget)
}
