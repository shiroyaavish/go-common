package fiber_cache

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/shiroyaavish/go-common/logger"
)

// InitialCacheOperation checks the client cache for data based on the provided options.
// If the options indicate a "GetOrCreate" operation, it runs the cache operation and returns the data if it exists.
// If the options indicate an "Invalidation" operation, it runs the cache operation.
// If the options are nil, it returns an error.
// If there are no errors and no data is found in the cache, it returns nil.
// The function accepts a `fiber.Ctx` object and a variadic parameter of `CacheUtils` objects as options.
// The `CacheUtils` interface is used to perform cache operations.
// Each `CacheUtils` object in the `opts` parameter must implement the `RunCacheOperation`, `IsGetOrCreate`,
// `SetData`, `GetData`, and `IsInvalidation` methods.
func InitialCacheOperation(c *fiber.Ctx, opts ...CacheUtils) (interface{}, error) {
	// Check Client Cache
	if opts == nil {
		return nil, errors.New("opts is nil")
	}
	// If opts is not null
	for _, opt := range opts {
		switch opt.GetOperationType() {
		case GetOrCreateOperationType:
			return CheckGetOperation(opt)
		case InvalidateOperationType:
			return RunInvalidateOperation(opt)
		}

	}
	return nil, nil
}

func CheckGetOperation(opt CacheUtils) (interface{}, error) {
	// Run cache
	err := opt.RunCacheOperation()
	// If err is not null
	if err != nil {
		// Log the error
		logger.Error(err)
		return nil, err
	}
	// check if data exists
	data := opt.GetData()
	// If data is not null
	return data, nil

}

func RunInvalidateOperation(opt CacheUtils) (interface{}, error) {
	// Run cache
	err := opt.RunCacheOperation()
	// If err is not null
	if err != nil {
		// Log the error
		logger.Error(err)
		return nil, err
	}
	return nil, nil
}

// EndCacheOperation updates the cache with the response body if the options indicate a "GetOrCreate" operation.
// The function accepts a `fiber.Ctx` object and a variadic parameter of `CacheUtils` objects as options.
// The `CacheUtils` interface is used to perform cache operations.
// Each `CacheUtils` object in the `opts` parameter must implement the `IsGetOrCreate` and `SetData` methods.
func EndCacheOperation(c *fiber.Ctx, respBody interface{}, opts ...CacheUtils) {
	go func() {
		if opts != nil {
			for _, opt := range opts {
				switch opt.GetOperationType() {
				case GetOrCreateOperationType:
					opt.SetData(respBody)
				}
			}
		}
	}()
	return
}

type OperationType string

const (
	InvalidateOperationType  OperationType = "Invalidate"
	GetOrCreateOperationType OperationType = "getOrCreate"
)

// CacheUtils is an interface that defines a method for performing cache operations.
type CacheUtils interface {
	RunCacheOperation() error
	SetData(d interface{})
	GetData() interface{}
	GetOperationType() OperationType
}
