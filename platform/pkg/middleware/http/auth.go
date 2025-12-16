package http

import (
	"context"
	"net/http"

	"github.com/vipshark78/microservices-course-homeworks/platform/pkg/middleware/session"
	authV1 "github.com/vipshark78/microservices-course-homeworks/shared/pkg/proto/auth/v1"
)

const SessionUUIDHeader = "X-Session-Uuid"

// IAMClient это алиас для сгенерированного gRPC клиента
type IAMClient = authV1.AuthServiceClient

// AuthMiddleware middleware для аутентификации HTTP запросов
type AuthMiddleware struct {
	iamClient IAMClient
}

// NewAuthMiddleware создает новый middleware аутентификации
func NewAuthMiddleware(iamClient IAMClient) *AuthMiddleware {
	return &AuthMiddleware{
		iamClient: iamClient,
	}
}

// client (X-Session-Uuid) -> auth middleware (add session_uuid in ctx (incomming)) -> order api (outgoing)-> auth interceptor ->inventory
// Handle обрабатывает HTTP запрос с аутентификацией
func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем session UUID из заголовка
		sessionUUID := r.Header.Get(SessionUUIDHeader)
		if sessionUUID == "" {
			writeErrorResponse(w, http.StatusUnauthorized, "MISSING_SESSION", "Authentication required")
			return
		}

		// Валидируем сессию через IAM сервис
		whoamiRes, err := m.iamClient.Whoami(r.Context(), &authV1.WhoamiRequest{
			SessionUuid: sessionUUID,
		})
		if err != nil {
			writeErrorResponse(w, http.StatusUnauthorized, "INVALID_SESSION", "Authentication failed")
			return
		}

		// Добавляем пользователя и session UUID в контекст используя функции из grpc middleware
		ctx := r.Context()
		ctx = AddSessionUUIDToContext(ctx, sessionUUID)
		// Также добавляем пользователя в контекст
		ctx = context.WithValue(ctx, session.UserContextKey, whoamiRes.User)

		// Передаем управление следующему handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AddSessionUUIDToContext добавляет session UUID в контекст
func AddSessionUUIDToContext(ctx context.Context, sessionUUID string) context.Context {
	return context.WithValue(ctx, session.SessionUUIDContextKey, sessionUUID)
}
