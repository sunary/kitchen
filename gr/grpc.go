package gr

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/textproto"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sunary/kitchen/e"
	"github.com/sunary/kitchen/l"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const fallbackMessage = `{"error": "failed to marshal error message", "code": 500}`

func init() {
	log.Println("Use custom runtime.HTTPError for grpc")
	runtime.HTTPError = DefaultHTTPError
}

// ValidateInterceptor returns middleware for validate
func ValidateInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		type validate interface {
			Validate() error
		}

		val := req.(validate)
		if err := val.Validate(); err != nil {
			return nil, e.Error(e.FailedPrecondition, "Invalid parameters", err)
		}

		return handler(ctx, req)
	}
}

// LogUnaryServerInterceptor returns middleware for logging with zap
func LogUnaryServerInterceptor(logger l.Logger, excepts ...error) grpc.UnaryServerInterceptor {
	m := make(map[error]struct{})
	for _, err := range excepts {
		m[err] = struct{}{}
	}
	getLogFn := func(v interface{}) (lg func(msg string, fields ...zapcore.Field)) {
		if err, ok := v.(error); ok {
			defer func() {
				if e := recover(); e != nil {
					lg = logger.Error
				}
			}()

			if _, isExcept := m[err]; isExcept {
				return logger.Warn
			}

			return logger.Error
		}

		return logger.Info
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		start := time.Now()
		defer func() {
			t := time.Now().Sub(start)

			e := recover()
			if e != nil {
				logger.Error("Panic (Recovered)", l.Error(err), l.Stack())
				err = status.Errorf(codes.Internal, "Internal Error (%v)", e)
			}

			if err == nil {
				logFn := getLogFn(resp)
				logFn(info.FullMethod, l.Duration("t", t), l.Interface("\n→", req), l.Interface("\n⇐", resp))
				return
			}

			logFn := getLogFn(err)

			// Append details
			if errorStatus, ok := status.FromError(err); ok {
				logFn(info.FullMethod, l.Duration("t", t), l.Interface("\n→", req), l.Stringer("\n⇐ERROR", errorStatus.Proto()))
			} else {
				logFn(info.FullMethod, l.Duration("t", t), l.Interface("\n→", req), l.String("\n⇐ERROR", err.Error()))
			}
		}()

		return handler(ctx, req)
	}
}

// ForwardMetadataUnaryServerInterceptor ...
func ForwardMetadataUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = ForwardMetadata(ctx)
		return handler(ctx, req)
	}
}

// ForwardMetadata ...
func ForwardMetadata(ctx context.Context) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	return metadata.NewOutgoingContext(ctx, md)
}

var jSON = runtime.JSONPb{
	OrigName:     true,
	EmitDefaults: true,
}

// NewJSONDecoder ...
func NewJSONDecoder(r io.Reader) runtime.Decoder {
	return jSON.NewDecoder(r)
}

// NewJSONEncoder ...
func NewJSONEncoder(w io.Writer) runtime.Encoder {
	return jSON.NewEncoder(w)
}

// MarshalJSON encodes JSON in compatible with GRPC
func MarshalJSON(v interface{}) ([]byte, error) {
	return jSON.Marshal(v)
}

// UnmarshalJSON decodes JSON in compatible with GRPC
func UnmarshalJSON(data []byte, v interface{}) error {
	return jSON.Unmarshal(data, v)
}

// DefaultHTTPError is extracted from grpc package
func DefaultHTTPError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	w.Header().Del("Trailer")
	w.Header().Set("Content-Type", marshaler.ContentType())

	body, buf, mErr := MarshalErrorWithDetails(marshaler, err)
	if mErr != nil {
		grpclog.Errorf("Failed to marshal error message %q: %v", body, mErr)
		w.WriteHeader(http.StatusInternalServerError)
		if _, wErr := io.WriteString(w, fallbackMessage); wErr != nil {
			grpclog.Errorf("Failed to write response: %v", wErr)
		}

		return
	}

	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		grpclog.Errorf("Failed to extract ServerMetadata from context")
	}

	handleForwardResponseServerMetadata(w, md)
	handleForwardResponseTrailerHeader(w, md)
	//w.WriteHeader(runtime.HTTPStatusFromCode(status.Code(err)))
	w.WriteHeader(httpStatusFromCode(e.Code(status.Code(err))))
	if _, err := w.Write(buf); err != nil {
		grpclog.Errorf("Failed to write response: %v", err)
	}

	handleForwardResponseTrailer(w, md)
}

func httpStatusFromCode(code e.Code) int {
	if code := e.DefaultHttpStatusFromCode(code); code != 0 {
		return code
	}

	grpclog.Infof("Unknown gRPC error code: %v", code)
	return http.StatusInternalServerError
}

// ErrorBody ...
type ErrorBody struct {
	Error   string          `protobuf:"bytes,1,name=error" json:"error"`
	Code    int32           `protobuf:"varint,2,name=code" json:"code"`
	Details []proto.Message `protobuf:"bytes,3,rep,name=details" json:"details,omitempty"`
}

// DecodeErrorWithDetails ...
func DecodeErrorWithDetails(err error) *ErrorBody {
	body := &ErrorBody{
		Error: status.Convert(err).Message(),
		Code:  int32(status.Code(err)),
	}

	// Append details
	if s, ok := status.FromError(err); ok {
		p := s.Proto()
		if len(p.Details) > 0 {
			details := make([]proto.Message, 0, len(p.Details))
			for _, d := range p.Details {
				var pmsg ptypes.DynamicAny
				e := ptypes.UnmarshalAny(d, &pmsg)
				if e != nil {
					log.Println("Error unmarshalling any: " + e.Error())
				} else {
					details = append(details, pmsg.Message)
				}
			}

			body.Details = details
		}
	}

	return body
}

// MarshalErrorWithDetails ...
func MarshalErrorWithDetails(marshaler runtime.Marshaler, err error) (body *ErrorBody, buf []byte, merr error) {
	body = DecodeErrorWithDetails(err)
	buf, merr = marshaler.Marshal(body)
	return
}

func handleForwardResponseServerMetadata(w http.ResponseWriter, md runtime.ServerMetadata) {
	for k, vs := range md.HeaderMD {
		hKey := fmt.Sprintf("%s%s", runtime.MetadataHeaderPrefix, k)
		for i := range vs {
			w.Header().Add(hKey, vs[i])
		}
	}
}

func handleForwardResponseTrailerHeader(w http.ResponseWriter, md runtime.ServerMetadata) {
	for k := range md.TrailerMD {
		tKey := textproto.CanonicalMIMEHeaderKey(fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k))
		w.Header().Add("Trailer", tKey)
	}
}

func handleForwardResponseTrailer(w http.ResponseWriter, md runtime.ServerMetadata) {
	for k, vs := range md.TrailerMD {
		tKey := fmt.Sprintf("%s%s", runtime.MetadataTrailerPrefix, k)
		for i := range vs {
			w.Header().Add(tKey, vs[i])
		}
	}
}
