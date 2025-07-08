package rabbit

import "github.com/IntelXLabs-LLC/go-common/errors"

var (
	ErrInvalidConfig = errors.NewError(50001, "invalid config")
)
