package rabbit

import "github.com/shiroyaavish/go-common/errors"

var (
	ErrInvalidConfig = errors.NewError(50001, "invalid config")
)
