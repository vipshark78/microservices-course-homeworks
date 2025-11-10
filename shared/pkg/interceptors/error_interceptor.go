package interceptors

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vipshark78/microservices-course-homeworks/shared/pkg/errors"
)

// UnaryErrorInterceptor возвращает middleware для обработки ошибок при вызове методов сервиса
func UnaryErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, convertError(err, info.FullMethod)
		}
		return resp, nil
	}
}

// convertError конвертирует ошибку в gRPC статус ошибки
func convertError(err error, method string) error {
	// Check if it's a businessError
	if businessErr := errors.GetBusinessError(err); businessErr != nil {
		grpcStatus := errors.BusinessErrorToGRPCStatus(businessErr)
		log.Printf("BusinessError in method %s: code=%d, message=%s",
			method, businessErr.Code(), businessErr.Error())
		return grpcStatus.Err()
	}

	// Проверяем на gRPC статус ошибку
	if _, ok := status.FromError(err); ok {
		return err
	}

	// если не бизнес-ошибка и не gRPC статус ошибки - то это внутренняя ошибка сервера
	log.Printf("Unknown error in method %s: %v", method, err)
	return status.Error(codes.Internal, "internal server error")
}
