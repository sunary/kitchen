package l

import (
	"go.uber.org/zap"
)

// Logger wraps zap.Logger
type Logger struct {
	*zap.Logger
}
