package fiber_cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/shiroyaavish/go-common/logger"
	"time"
)

/*Operation Type: CreateOrGetOperation*/

// CreateOrGetOperation is a struct that represents a cache create operation with a key and data or get data
type CreateOrGetOperation struct {
	// key
	key string
	// data
	data interface{}
	// dataChan
	dataChan chan []byte
}

// RunCacheOperation executes the cache create operation for a specific cache entry.
// It sets the cache entry using the key and data stored in the CreateOperation instance.
// If there is an error during the operation, it returns the error.
// This method returns the error as nil if the operation is successful.
func (c *CreateOrGetOperation) RunCacheOperation() error {
	// If c.key exists
	response := GetClient().Get(context.Background(), c.key)
	// Get Err
	err := response.Err()
	// if err is not null
	if err != nil {
		// If errors Is Redis Null
		if errors.Is(err, redis.Nil) {
			// GetClient and set
			go func() {
				c.data = <-c.dataChan
				logger.Info("Received Data")
				err = GetClient().Set(context.Background(), c.key, c.data, time.Hour*24).Err()
				if err != nil {
					logger.Error(err)
				}
				logger.Info("cache operation finished: %s", c.key)
			}()
		} else {
			// Logger is an error
			logger.Error(err)
		}
	}
	// Get Bytes
	data, err := response.Bytes()
	// If err is not null
	if err != nil {
		// Return err
		return err
	}
	// Var i interface
	var i interface{}
	// Unmarshal data
	err = json.Unmarshal(data, &i)
	// If err is not null
	if err != nil {
		// logger.Error
		logger.Error(err)
	}
	c.data = i
	return nil
}

// WithCreateOrGetOperation returns a new instance of CreateOperation with the specified key and data
// that represents a cache create operation.
func WithCreateOrGetOperation(key string, data interface{}) *CreateOrGetOperation {
	return &CreateOrGetOperation{key: key, data: data, dataChan: make(chan []byte, 1)}
}

// SetData will set the data
func (c *CreateOrGetOperation) SetData(d interface{}) {
	dataToSet, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	logger.Info("cache operation data sent: %s", c.key)
	c.dataChan <- dataToSet
}

func (c *CreateOrGetOperation) GetData() interface{} {
	return c.data
}

func (c *CreateOrGetOperation) GetOperationType() OperationType {
	return GetOrCreateOperationType
}

/*Operation Type: Invalidate*/

// InvalidateOperation is a type that represents an operation to invalidate a cache entry.
type InvalidateOperation struct {
	key string
}

// WithInvalidate returns a new instance of InvalidateOperation with the specified key
// that represents an operation to invalidate a cache entry.
func WithInvalidate(key string) *InvalidateOperation {
	return &InvalidateOperation{key: key}
}

// RunCacheOperation executes the cache invalidation operation for a specific cache entry.
// It deletes the cache entry using the key stored in the InvalidateOperation instance.
// If the cache entry doesn't exist, it returns nil (redis.Nil error).
// If there is an error during the operation, it returns the error.
// This method returns the error as nil if the operation is successful.
func (i *InvalidateOperation) RunCacheOperation() error {
	err := GetClient().Del(context.Background(), i.key).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}
	return nil
}

// SetData sets the data
func (i *InvalidateOperation) SetData(d interface{}) {
	return
}

func (i *InvalidateOperation) GetData() interface{} {
	return nil
}

func (i *InvalidateOperation) GetOperationType() OperationType {
	return InvalidateOperationType
}
