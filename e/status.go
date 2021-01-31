package e

import (
	"net/http"
	"strconv"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Code uint32

type detailMsg struct {
	detail string
}

func NewProtoMessage(detail string) proto.Message {
	return &detailMsg{
		detail: detail,
	}
}

func (d *detailMsg) Reset()         { *d = detailMsg{} }
func (d *detailMsg) String() string { return d.detail }
func (*detailMsg) ProtoMessage()    {}

type Status struct {
	*status.Status
}

// Error new status with code and message
func Error(code Code, message string, errs ...error) *Status {
	details := make([]proto.Message, len(errs))
	for i := range errs {
		details[i] = NewProtoMessage(errs[i].Error())
	}

	s := status.New(codes.Code(code), message)
	s, _ = s.WithDetails(details...)
	return &Status{s}
}

// Errorf new status with code and message
func Errorf(code Code, format string, args ...interface{}) *Status {
	return &Status{status.Newf(codes.Code(code), format, args...)}
}

// GRPCStatus ...
func (s Status) GRPCStatus() *status.Status {
	return s.Status
}

// Error ...
func (s Status) Error() string {
	if m := s.Message(); m != "" {
		return m
	}
	return strconv.Itoa(int(s.Code()))
}

// DefaultHttpStatusFromCode ...
func DefaultHttpStatusFromCode(code Code) int {
	switch codes.Code(code) {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	default:
		return 0
	}
}
