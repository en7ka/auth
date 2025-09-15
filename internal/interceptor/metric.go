package interceptor

import (
	"context"
	"time"

	"github.com/en7ka/auth/internal/metric"
	"google.golang.org/grpc"
)

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	metric.IncRequestCounter()

	timeStart := time.Now()
	difftime := time.Since(timeStart)

	res, err := handler(ctx, req)
	if err != nil {
		metric.IncResponseCounter("error", info.FullMethod)
		metric.HistogramResponseTime("error", difftime.Seconds())
	} else {
		metric.IncResponseCounter("success", info.FullMethod)
		metric.HistogramResponseTime("success", difftime.Seconds())
	}

	return res, err
}
