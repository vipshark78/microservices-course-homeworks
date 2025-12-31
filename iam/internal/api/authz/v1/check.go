package v1

import (
	"context"
	"errors"
	"fmt"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vipshark78/microservices-course-homeworks/iam/internal/model"
)

const (
	HeaderUserUUID    = "x-user-uuid"
	HeaderUserLogin   = "x-user-login"
	HeaderUserEmail   = "x-user-email"
	HeaderSessionUUID = "X-Session-Uuid"
	HeaderContentType = "content-type"

	ContentTypeJSON = "application/json"
)

// Check реализует метод Check из envoy.service.auth.v3.Authorization
// Проверяет аутентификацию запроса и возвращает OK или Denied
func (a *api) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	// Извлекаем session UUID из заголовка session-uuid (передается Envoy через allowed_headers)
	sessionUUID, err := a.extractSessionUUID(req)
	if err != nil {
		return a.denyRequest("Missing or invalid session", 401), nil
	}

	// Вызываем Whoami для проверки сессии и получения информации о пользователе
	resp, err := a.authService.Whoami(ctx, sessionUUID)
	if err != nil {
		if errors.Is(err, model.ErrSessionNotFound) || errors.Is(err, model.ErrSessionExpired) {
			return a.denyRequest("Invalid or expired session", 401), nil
		}
		return a.denyRequest("Internal server error", 500), nil
	}

	// Разрешаем запрос и добавляем заголовки с информацией о пользователе
	return a.allowRequest(resp.Session, resp.User), nil
}

// extractSessionUUID извлекает session UUID из заголовков запроса
// Envoy передает заголовок session-uuid из входящего запроса через allowed_headers
func (a *api) extractSessionUUID(req *authv3.CheckRequest) (string, error) {
	if req.Attributes == nil || req.Attributes.Request == nil {
		return "", fmt.Errorf("no HTTP request found")
	}

	headers := req.Attributes.Request.Http.Headers

	// Проверяем заголовок session-uuid (основной способ, используется в тестах)
	if sessionUUID, ok := headers["session-uuid"]; ok && sessionUUID != "" {
		return sessionUUID, nil
	}

	return "", fmt.Errorf("session uuid not found in headers")
}

// allowRequest создает успешный ответ CheckResponse с заголовками пользователя
func (a *api) allowRequest(session model.Session, user model.User) *authv3.CheckResponse {
	headers := []*corev3.HeaderValueOption{
		{
			Header: &corev3.HeaderValue{
				Key:   HeaderUserUUID,
				Value: user.UserUUID,
			},
		},
		{
			Header: &corev3.HeaderValue{
				Key:   HeaderUserLogin,
				Value: user.UserInfo.Login,
			},
		},
		{
			Header: &corev3.HeaderValue{
				Key:   HeaderUserEmail,
				Value: user.UserInfo.Email,
			},
		},
		{
			Header: &corev3.HeaderValue{
				Key:   HeaderSessionUUID,
				Value: session.UUID(),
			},
		},
	}

	return &authv3.CheckResponse{
		Status: status.New(codes.OK, "").Proto(),
		HttpResponse: &authv3.CheckResponse_OkResponse{
			OkResponse: &authv3.OkHttpResponse{
				Headers: headers,
			},
		},
	}
}

// denyRequest создает ответ об отказе в доступе
func (a *api) denyRequest(message string, statusCode int32) *authv3.CheckResponse {
	var httpStatusCode typev3.StatusCode
	switch statusCode {
	case 401:
		httpStatusCode = typev3.StatusCode_Unauthorized
	case 403:
		httpStatusCode = typev3.StatusCode_Forbidden
	default:
		httpStatusCode = typev3.StatusCode_InternalServerError
	}

	return &authv3.CheckResponse{
		Status: status.New(codes.Unauthenticated, message).Proto(),
		HttpResponse: &authv3.CheckResponse_DeniedResponse{
			DeniedResponse: &authv3.DeniedHttpResponse{
				Status: &typev3.HttpStatus{
					Code: httpStatusCode,
				},
				Body: fmt.Sprintf(`{"error": "%s"}`, message),
				Headers: []*corev3.HeaderValueOption{
					{
						Header: &corev3.HeaderValue{
							Key:   HeaderContentType,
							Value: ContentTypeJSON,
						},
					},
				},
			},
		},
	}
}
