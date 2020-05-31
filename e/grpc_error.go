package e

import (
	"log"
	"runtime/debug"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	pbStatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

// GRPCError ...
func GRPCError(code codes.Code, message string, details ...proto.Message) error {
	if message == "" {
		message = code.String()
	}

	s := &pbStatus.Status{
		Code:    int32(code),
		Message: message,
	}

	if len(details) > 0 {
		ds := make([]*any.Any, len(details))
		for i, d := range details {
			mAny, err := ptypes.MarshalAny(d)
			if err != nil {
				debug.PrintStack()
				log.Println("Unable to marshal any")
				ds[i], _ = ptypes.MarshalAny(status.New(codes.Internal, "Unable to marshal to grpc.Any").Proto())
			} else {
				ds[i] = mAny
			}
		}

		s.Details = ds
	}

	return status.ErrorProto(s)
}

// ToGRPCError ...
func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	if eErr, ok := err.(*APIError); ok {
		return GRPCError(codes.Code(eErr.Code), eErr.Message)
	}

	code := status.Code(err)
	if code == codes.Unknown {
		return GRPCError(codes.Internal, "Internal Server Error")
	}

	return err
}
