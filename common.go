package go_common

import (
	"github.com/IntelXLabs-LLC/go-common/config"
	"github.com/IntelXLabs-LLC/go-common/logger"
	"os"
	"os/signal"
	"syscall"
)

// Common is a type representing a common structure.
//
// Example usage:
//
//	func UseDefault() *Common {
//	    config.UseDefaultConfig()
//	    return &Common{}
//	}
//
//	func (c *Common) WithCustomConfig() *Common {
//	    return c
//	}
//
//	func (c *Common) WithErrors() *Common {
//	    return c
//	}
//
//	func (c *Common) Init() {
//
//	}
//
//	func (c *Common) WithSignalCheck() *Common {
//	    go func() {
//	        c := make(chan os.Signal, 1)
//	        signal.Notify(c, syscall.SIGTERM, syscall.SIGKILL)
//	    }()
//	    return c
//	}
//
//	func (c *Common) GracefulShutdown(callback func(chan os.Signal) error) {
//	    logger.Info("Started Graceful Shutdown")
//	    if err := callback(c.signalChannel); err != nil {
//	        logger.Fatal(err)
//	    }
//	    os.Exit(0)
//	}
type Common struct {
	signalChannel chan os.Signal
	isStartReady  bool
}

// UseDefault is a function that sets the default configuration by using config.UseDefaultConfig() and returns a new instance of Common.
// Example usage:
// common := UseDefault()
//
// Returns:
// - A pointer to a new instance of Common.
//
// Example:
// config.UseDefaultConfig()
// return &Common{}
func UseDefault() *Common {
	config.UseDefaultConfig()
	return &Common{}
}

func (c *Common) WithCustomConfig() *Common {

	return c
}

func (c *Common) WithErrors() *Common {

	return c
}

// UseWithConfig is a function that sets the given Config as the common configuration and returns a new instance of Common.
// Example usage:
//
//	common := UseWithConfig(commonConfig)
//
// Parameters:
//   - commonConfig: A pointer to a config.commonConfig object that contains the configuration.
//
// Returns:
//   - A pointer to a new instance of Common.
//
// Example:
//
//	config.SetCommonConfig(commonConfig)
//	return &Common{}
func UseWithConfig(commonConfig *config.CommonConfig) *Common {
	config.SetCommonConfig(commonConfig)
	return &Common{}
}

func UseAWS(awsConfig *config.AWSConfig) *Common {
	config.SetAWSConfig(awsConfig)
	return &Common{}
}

// Init is a method of the Common struct that initializes the configuration and other dependencies.
func (c *Common) Init() {

}

// WithSignalCheck is a method of the Common struct that sets up a signal handler to catch termination signals (SIGTERM, SIGKILL).
// It creates a channel to receive the signals and notifies the channel when the specified signals are received.
// This method runs as a goroutine to continuously listen for signals in the background.
// It returns the Common object itself, allowing method chaining.
// Example usage:
//
//	common := &Common{}
//	common.WithSignalCheck()
//	// ...do other operations
//	common.GracefulShutdown()
//
// more detailed examples can be seen in the Common struct docs.
func (c *Common) WithSignalCheck() *Common {
	c.signalChannel = make(chan os.Signal, 1)
	signal.Notify(c.signalChannel, syscall.SIGTERM, syscall.SIGKILL)
	return c
}

// GracefulShutdown is a method of the Common struct that waits for a signal on the signalChannel and then exits the program.
// It uses the <- operator to receive a value from the signalChannel, and then calls os.Exit(0) to exit the program with a status code of 0.
// This method is typically used to gracefully shut down a program when a termination signal is received.
//
// the function accepts a callback that
func (c *Common) GracefulShutdown(callback func() error) *Common {
	logger.Info("Started Graceful Shutdown")
	go func() {
		for {
			select {
			case <-c.signalChannel:
				logger.Info("Received termination signal. Starting graceful shutdown...")
				if err := callback(); err != nil {
					logger.Fatal(err)
				}
				os.Exit(0)
			}
		}
	}()
	c.isStartReady = true
	return c
}

// Start will start the program if start is ready, start will only be ready with GracefulShutdown being turned on.
// The program will not start without calling it
func (c *Common) Start() {
	if c.isStartReady {
		<-c.signalChannel
	} else {
		logger.Info("The Program will exit: Implement common.GracefulShutdown\nRefer Documentation")
		os.Exit(1)
	}
}
