package session

type contextKey string

const (
	// userContextKey ключ для хранения пользователя в контексте
	UserContextKey contextKey = "user"
	// sessionUUIDContextKey ключ для хранения session UUID в контексте
	SessionUUIDContextKey contextKey = "session-uuid"
)
