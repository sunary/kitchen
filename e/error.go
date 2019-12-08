package e

import (
	"log"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	pbstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Code ...
type Code int

const (
	// grpc error codes
	NoError            = Code(codes.OK)
	Canceled           = Code(codes.Canceled)
	Unknown            = Code(codes.Unknown)
	InvalidArgument    = Code(codes.InvalidArgument)
	DeadlineExceeded   = Code(codes.DeadlineExceeded)
	NotFound           = Code(codes.NotFound)
	AlreadyExists      = Code(codes.AlreadyExists)
	PermissionDenied   = Code(codes.PermissionDenied)
	ResourceExhausted  = Code(codes.ResourceExhausted)
	FailedPrecondition = Code(codes.FailedPrecondition)
	Aborted            = Code(codes.Aborted)
	OutOfRange         = Code(codes.OutOfRange)
	Unimplemented      = Code(codes.Unimplemented)
	Internal           = Code(codes.Internal)
	Unavailable        = Code(codes.Unavailable)
	DataLoss           = Code(codes.DataLoss)
	Unauthenticated    = Code(codes.Unauthenticated)
)

var (
	mapCodes = map[Code]string{}
)

func init() {
	mapCodes[NoError] = "ok"
	mapCodes[Canceled] = "canceled"
	mapCodes[Unknown] = "unknown"
	mapCodes[InvalidArgument] = "invalid_argument"
	mapCodes[DeadlineExceeded] = "deadline_exceeded"
	mapCodes[NotFound] = "not_found"
	mapCodes[AlreadyExists] = "already_exists"
	mapCodes[PermissionDenied] = "permission_denied"
	mapCodes[ResourceExhausted] = "resource_exhausted"
	mapCodes[FailedPrecondition] = "failed_precondition"
	mapCodes[Aborted] = "aborted"
	mapCodes[OutOfRange] = "out_of_range"
	mapCodes[Unimplemented] = "unimplemented"
	mapCodes[Internal] = "internal"
	mapCodes[Unavailable] = "unavailable"
	mapCodes[DataLoss] = "data_loss"
	mapCodes[Unauthenticated] = "unauthenticated"
}

// APIError ...
type APIError struct {
	Code    Code // A standard grpc error code.
	Err     error
	Message string
}

func Error(code Code, message string, errs ...error) *APIError {
	if message == "" {
		message = code.String()
	}

	var err error
	if len(errs) > 0 {
		err = errs[0]
	}

	return &APIError{
		Code:    code,
		Err:     err,
		Message: message,
	}
}

// Error ...
func (e *APIError) Error() string {
	return e.Message
}

// String ...
func (c Code) String() string {
	if msg, ok := mapCodes[c]; ok {
		return msg
	}

	return "Code(" + strconv.Itoa(int(c)) + ")"
}

// GRPCError ...
func GRPCError(code codes.Code, msg string, details ...proto.Message) error {
	if msg == "" {
		msg = code.String()
	}

	s := &pbstatus.Status{
		Code:    int32(code),
		Message: msg,
	}
	if len(details) > 0 {
		ds := make([]*any.Any, len(details))
		for i, d := range details {
			any, err := ptypes.MarshalAny(d)
			if err != nil {
				debug.PrintStack()
				log.Println("Unable to marshal any")
				ds[i], _ = ptypes.MarshalAny(status.New(codes.Internal, "Unable to marshal to grpc.Any").Proto())
			} else {
				ds[i] = any
			}
		}

		s.Details = ds
	}

	return status.ErrorProto(s)
}

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
