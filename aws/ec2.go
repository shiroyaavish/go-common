package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/shiroyaavish/go-common/errors"
	"github.com/shiroyaavish/go-common/logger"
	"time"
)

type EC2 struct {
	listen          bool
	listenerChannel chan bool
	listenerHandler ListenerCallback
	ec2Svc          *ec2.EC2
}

// NewEC2Operation initializes and returns a new instance of EC2.
// If ec2Svc is nil, it calls CreateEC2Session to create_svc a new EC2 session.
func (a *AWS) NewEC2Operation() *EC2 {
	if a.sess == nil {
		a.ConnectAws()
	}

	return &EC2{ec2Svc: ec2.New(a.sess)}
}

type Opts interface {
	getBool(ec2 *EC2)
}

type withStatusListener struct {
	d       chan bool
	handler ListenerCallback
}

func (w *withStatusListener) getBool(ec2 *EC2) {
	ec2.listen = true
	ec2.listenerChannel = w.d
	ec2.listenerHandler = w.handler
}

type ListenerCallback func() error

func WithListener(channel chan bool, callback ListenerCallback) *withStatusListener {
	return &withStatusListener{
		d:       channel,
		handler: callback,
	}
}

func (s *EC2) Raw() *ec2.EC2 {
	return s.ec2Svc
}

// EC2Builder is a type used to build EC2 instances with various configurations.
//
// EC2Builder is a struct that contains a `data` field of type `*ec2.RunInstancesInput`
// and an `isNetworkSet` field of type `bool`.
//
// The `data` field is a pointer to `ec2.RunInstancesInput`, which holds the input parameters
// for running EC2 instances. The `isNetworkSet` field is a boolean value indicating whether
// the network configuration has been set for the EC2 instance.
//
// Usage example:
//
//	func (s *EC2) NewInstance() *EC2Builder {
//	   return &EC2Builder{}
//	}
type EC2Builder struct {
	data         *ec2.RunInstancesInput
	isNetworkSet bool
	ec2Builder   *EC2
}

// NewInstance runs new instances
func (s *EC2) NewInstance() *EC2Builder {
	return &EC2Builder{
		ec2Builder: s,
	}
}

// AddAMI adds the AMI (Amazon Machine Image) ID to the EC2Builder instance.
// If the imageId is empty, it will panic with a custom error.
// It returns the EC2Builder instance to allow method chaining.
func (e *EC2Builder) AddAMI(imageId string) *EC2Builder {
	if imageId == "" {
		panic(errors.NewError(404, "image id cannot be empty"))
	}
	e.data.ImageId = &imageId
	return e
}

// AddInstanceType sets the instance type for the EC2Builder.
// It panics if the instanceType provided is empty.
// Returns the updated EC2Builder.
func (e *EC2Builder) AddInstanceType(instanceType string) *EC2Builder {
	if instanceType == "" {
		panic(errors.NewError(404, "instance type cannot be empty"))
	}
	e.data.InstanceType = &instanceType
	return e
}

// AddMinMax sets the minimum and maximum count of instances for the EC2 builder.
// If either min or max is negative, it will panic with an error.
// Returns the EC2 builder itself.
func (e *EC2Builder) AddMinMax(min, max int) *EC2Builder {
	if min < 0 || max < 0 {
		panic(errors.NewError(404, "min and max values cannot be negative"))
	}
	e.data.MinCount = aws.Int64(int64(min))
	e.data.MaxCount = aws.Int64(int64(max))

	return e
}

// AddEC2InstanceKey sets the key name for the EC2 instance.
// It panics with an error if the key name is empty.
// Returns the EC2Builder object.
func (e *EC2Builder) AddEC2InstanceKey(key string) *EC2Builder {
	if key == "" {
		panic(errors.NewError(404, "key name cannot be empty"))
	}
	e.data.KeyName = &key
	return e
}

// AddCloudInitScript adds a cloud-init script to the EC2Builder instance.
// The script parameter should not be empty. If it is empty, the function will panic with an error.
// The cloud-init script will be set to the UserData field of the EC2Builder instance.
// The method returns the updated EC2Builder instance, allowing for method chaining.
func (e *EC2Builder) AddCloudInitScript(script string) *EC2Builder {
	if script == "" {
		panic(errors.NewError(404, "cloud init script, if used should not be empty"))
	}
	e.data.UserData = &script
	return e
}

// NewEBS creates a new EBS (Elastic Block Store) with the specified volume size.
// If the volume size is negative, it panics with an error indicating that the volume size cannot be negative.
// The created EBS is added to the EC2Builder's block device mappings.
// Returns the EC2Builder instance.
func (e *EC2Builder) NewEBS(volume int64) *EC2Builder {
	if volume < 0 {
		panic(errors.NewError(404, "volume size cannot be negative"))
	}
	e.data.BlockDeviceMappings = []*ec2.BlockDeviceMapping{
		{
			DeviceName: aws.String("/dev/sda1"),
			Ebs: &ec2.EbsBlockDevice{
				VolumeSize: aws.Int64(volume),
			},
		},
	}
	return e
}

// AddNetworkInstances adds network interfaces to the EC2Builder struct
// It takes four parameters:
// - index: an int representing the device index
// - isPublicIP: a bool indicating whether to associate a public IP address
// - subnetID: a string representing the subnet ID
// - groups: a string representing the security groups
// It sets the NetworkInterfaces field of the EC2Builder's data with the provided values
// It also sets the isNetworkSet field to true
// It returns the EC2Builder instance
func (e *EC2Builder) AddNetworkInstances(index int, isPublicIP bool, subnetID string, groups string) *EC2Builder {
	if index < 0 {
		panic(errors.NewError(404, "index cannot be negative"))
	}
	if subnetID == "" {
		panic(errors.NewError(404, "subnet id cannot be empty"))
	}
	if groups == "" {
		panic(errors.NewError(404, "groups cannot be empty"))
	}
	e.data.NetworkInterfaces = []*ec2.InstanceNetworkInterfaceSpecification{
		{
			DeviceIndex:              aws.Int64(int64(index)),
			AssociatePublicIpAddress: aws.Bool(isPublicIP),
			SubnetId:                 aws.String(subnetID),
			Groups: []*string{
				aws.String(groups),
			},
		},
	}
	e.isNetworkSet = true
	return e
}

// SetInstanceMarketType sets the market type for the instance
func (e *EC2Builder) SetInstanceMarketType(s string) *EC2Builder {
	if s == "" {
		panic(errors.NewError(404, "market type cannot be empty"))
	}
	e.data.InstanceMarketOptions = &ec2.InstanceMarketOptionsRequest{
		MarketType: aws.String(s),
	}
	return e
}

// SetInstanceIAMRole sets the IAM role for the instance
func (e *EC2Builder) SetInstanceIAMRole(s string) *EC2Builder {
	if s == "" {
		panic(errors.NewError(404, "IAM role cannot be empty"))
	}
	e.data.IamInstanceProfile = &ec2.IamInstanceProfileSpecification{
		Name: aws.String(s),
	}
	return e
}

// SetInstanceTag sets the tag for the EC2 instance.
// If the tag is empty, it panics with a custom error.
// The tag is specified using the "Name" key.
// It returns the EC2Builder instance for method chaining.
func (e *EC2Builder) SetInstanceTag(t string) *EC2Builder {
	if t == "" {
		panic(errors.NewError(404, "tag cannot be empty"))
	}
	e.data.TagSpecifications = []*ec2.TagSpecification{
		{
			ResourceType: aws.String("instance"),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String(t),
				},
			},
		},
	}
	return e
}

// Run runs the EC2 instances and returns the instance ID of the first instance created and any error encountered during the process.
func (e *EC2Builder) Run(opts ...Opts) (string, error) {
	for _, el := range opts {
		el.getBool(e.ec2Builder)
	}
	if !e.isNetworkSet {
		return "", errors.NewError(500, "network is not set")
	}
	result, err := e.ec2Builder.ec2Svc.RunInstances(e.data)
	if err != nil {
		return "", errors.NewErrorWrapper(err, 500, "failed to run instances")
	}

	if e.ec2Builder.listen {
		e.ec2Builder.monitorStatus(e.ec2Builder.listenerChannel, *result.Instances[0].InstanceId, "running")
	}
	return *result.Instances[0].InstanceId, nil
}

type EC2Stopper struct {
	instanceId string
	EC2        *EC2
}

func (s *EC2) StopInstance() *EC2Stopper {
	return &EC2Stopper{}
}

func (e *EC2Stopper) AddInstanceId(id string) *EC2Stopper {
	e.instanceId = id
	return e
}

// monitorStatus monitors the EC2 instance status every 20 seconds, and empties the channel if the instance is stopped.
func (s *EC2) monitorStatus(ch chan bool, instanceId string, status string) {
	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			out, _ := s.ec2Svc.DescribeInstances(&ec2.DescribeInstancesInput{
				InstanceIds: []*string{aws.String(instanceId)},
			})

			statusCheck := false

			for _, reservation := range out.Reservations {
				for _, instance := range reservation.Instances {
					if *instance.State.Name == status {
						for len(ch) > 0 {
							statusCheck = true
						}
					}
				}
			}
			if statusCheck {
				<-ch
				err := s.listenerHandler()
				if err != nil {
					logger.Error(err)
				}
			}
		}
	}()
}

func (e *EC2Stopper) Stop(opts ...Opts) error {
	for _, el := range opts {
		el.getBool(e.EC2)
	}
	if e.EC2.listen {
		e.EC2.monitorStatus(e.EC2.listenerChannel, e.instanceId, "stopped")
	}
	out, err := e.EC2.ec2Svc.StopInstances(&ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(e.instanceId),
		},
	})
	if err != nil {
		return err
	}
	logger.Info(out.String())
	return nil
}

func NoCallback() ListenerCallback {
	return func() error {
		return nil
	}
}
