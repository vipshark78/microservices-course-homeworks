package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/middleware/session"
	authV1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/auth/v1"
)

const (
	// SessionUUIDMetadataKey ключ для передачи UUID сессии в gRPC metadata
	SessionUUIDMetadataKey = "session-uuid"
)

// IAMClient это алиас для сгенерированного gRPC клиента
type IAMClient = authV1.AuthServiceClient

// AuthInterceptor interceptor для аутентификации gRPC запросов
type AuthInterceptor struct {
	iamClient IAMClient
}

// NewAuthInterceptor создает новый interceptor аутентификации
func NewAuthInterceptor(iamClient IAMClient) *AuthInterceptor {
	return &AuthInterceptor{
		iamClient: iamClient,
	}
}

// Unary возвращает unary server interceptor для аутентификации
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		authCtx, err := i.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		return handler(authCtx, req)
	}
}

// authenticate выполняет аутентификацию и добавляет пользователя в контекст
func (i *AuthInterceptor) authenticate(ctx context.Context) (context.Context, error) {
	// Извлекаем metadata из контекста
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	// Получаем session UUID из metadata
	sessionUUIDs := md.Get(SessionUUIDMetadataKey)
	if len(sessionUUIDs) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing session-uuid in metadata")
	}

	sessionUUID := sessionUUIDs[0]
	if sessionUUID == "" {
		return nil, status.Error(codes.Unauthenticated, "empty session-uuid")
	}

	// Валидируем сессию через IAM сервис
	whoamiRes, err := i.iamClient.Whoami(ctx, &authV1.WhoamiRequest{
		SessionUuid: sessionUUID,
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("invalid session: %v", err))
	}

	// Добавляем пользователя и session UUID в контекст
	authCtx := context.WithValue(ctx, session.UserContextKey, whoamiRes.User)
	authCtx = context.WithValue(authCtx, session.SessionUUIDContextKey, sessionUUID)
	return authCtx, nil
}

// GetSessionUUIDFromContext извлекает session UUID из контекста
func GetSessionUUIDFromContext(ctx context.Context) (string, bool) {
	sessionUUID, ok := ctx.Value(session.SessionUUIDContextKey).(string)
	return sessionUUID, ok
}

// ForwardSessionUUIDToGRPC добавляет session UUID из контекста в исходящие gRPC metadata
func ForwardSessionUUIDToGRPC(ctx context.Context) context.Context {
	sessionUUID, ok := GetSessionUUIDFromContext(ctx)
	if !ok || sessionUUID == "" {
		return ctx
	}

	return metadata.AppendToOutgoingContext(ctx, SessionUUIDMetadataKey, sessionUUID)
}
