package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	startTime := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(startTime)

	// 응답 + 에러 로깅
	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}
	logger.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status", statusCode.String()).
		Dur("duration", duration).
		Msg("gRPC request received")

	return resp, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	StateCode int
	Body      []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StateCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: res,
			StateCode:      http.StatusOK,
		}

		handler.ServeHTTP(rec, req)
		duration := time.Since(startTime)

		logger := log.Info()
		if rec.StateCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body)
		}
		logger.Str("protocol", "http").
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Int("status_code", rec.StateCode).
			Str("status", http.StatusText(rec.StateCode)).
			Dur("duration", duration).
			Msg("gRPC request received")
	})
}
