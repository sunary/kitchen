package e

import (
	"net/http"
	"strconv"

	"google.golang.org/grpc/status"
)

// Code ...
type Code int32

var (
	mapCodes = map[Code]string{
		NoError:            "ok",
		Canceled:           "canceled",
		Unknown:            "unknown",
		InvalidArgument:    "invalid_argument",
		DeadlineExceeded:   "deadline_exceeded",
		NotFound:           "not_found",
		AlreadyExists:      "already_exists",
		PermissionDenied:   "permission_denied",
		ResourceExhausted:  "resource_exhausted",
		FailedPrecondition: "failed_precondition",
		Aborted:            "aborted",
		OutOfRange:         "out_of_range",
		Unimplemented:      "unimplemented",
		Internal:           "internal",
		Unavailable:        "unavailable",
		DataLoss:           "data_loss",
		Unauthenticated:    "unauthenticated",
	}
)

// String ...
func (c Code) String() string {
	if msg, ok := mapCodes[c]; ok {
		return msg
	}

	return "Code(" + strconv.Itoa(int(c)) + ")"
}

// ErrorCode ...
func ErrorCode(err error) Code {
	if err == nil {
		return NoError
	}

	if err, ok := err.(*APIError); ok {
		return err.Code
	}

	return Unknown
}

// ErrorMessage ...
func ErrorMessage(err error) string {
	return status.Convert(err).Message()
}

// DefaultHttpStatusFromCode ...
func DefaultHttpStatusFromCode(code Code) int {
	switch code {
	case NoError:
		return http.StatusOK
	case Canceled:
		return http.StatusRequestTimeout
	case Unknown:
		return http.StatusInternalServerError
	case InvalidArgument:
		return http.StatusBadRequest
	case DeadlineExceeded:
		return http.StatusGatewayTimeout
	case NotFound:
		return http.StatusNotFound
	case AlreadyExists:
		return http.StatusConflict
	case PermissionDenied:
		return http.StatusForbidden
	case Unauthenticated:
		return http.StatusUnauthorized
	case ResourceExhausted:
		return http.StatusTooManyRequests
	case FailedPrecondition:
		return http.StatusPreconditionFailed
	case Aborted:
		return http.StatusConflict
	case OutOfRange:
		return http.StatusBadRequest
	case Unimplemented:
		return http.StatusNotImplemented
	case Internal:
		return http.StatusInternalServerError
	case Unavailable:
		return http.StatusServiceUnavailable
	case DataLoss:
		return http.StatusInternalServerError
	default:
		return 0
	}
}
