package http_errors

import (
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var httpStatus = map[codes.Code]int{
	codes.AlreadyExists:      fasthttp.StatusPreconditionFailed,
	codes.FailedPrecondition: fasthttp.StatusPreconditionFailed,
	codes.Internal:           fasthttp.StatusInternalServerError,
	codes.InvalidArgument:    fasthttp.StatusBadRequest,
	codes.NotFound:           fasthttp.StatusNotFound,
	codes.PermissionDenied:   fasthttp.StatusUnauthorized,
	codes.OutOfRange:         fasthttp.StatusRequestedRangeNotSatisfiable,
	codes.Unknown:            fasthttp.StatusInternalServerError,
	codes.DeadlineExceeded:   fasthttp.StatusGatewayTimeout,
	codes.Unavailable:        fasthttp.StatusServiceUnavailable,
}

// SetInternalError ...
func SetInternalError(request *routing.Context, requestID string) error {
	request.Response.Header.Set("X-Request-ID", requestID)
	request.Response.SetStatusCode(fasthttp.StatusInternalServerError)
	return nil
}

// SetInvalidParam ...
func SetInvalidParam(request *routing.Context, requestID string, err string) error {
	request.Response.Header.Set("X-Request-ID", requestID)
	request.Response.SetStatusCode(fasthttp.StatusBadRequest)
	request.Response.SetBodyString(err)
	return nil
}

// SetMarshallingError ...
func SetMarshallingError(request *routing.Context, requestID string) error {
	request.Response.Header.Set("X-Request-ID", requestID)
	request.Response.SetStatusCode(fasthttp.StatusUnprocessableEntity)
	request.Response.SetBodyString("Marshalling/unmarshalling error")
	return nil
}

// SetRPCError ...
func SetRPCError(err error, request *routing.Context, requestID string) error {
	errorStatus, _ := status.FromError(err)
	request.Response.Header.Set("X-Request-ID", requestID)
	if errorStatus.Code() == codes.Internal {
		request.Response.SetStatusCode(fasthttp.StatusUnprocessableEntity)
	} else {
		request.Response.SetStatusCode(httpStatus[errorStatus.Code()])
		request.Response.SetBodyString(errorStatus.Message())
	}
	return nil
}
