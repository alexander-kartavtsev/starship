package http

import (
	"context"
	"log"
	"net/http"

	grpcAuth "github.com/alexander-kartavtsev/starship/platform/pkg/middleware/grpc"
	authV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/auth/v1"
	commonV1 "github.com/alexander-kartavtsev/starship/shared/pkg/proto/common/v1"
)

const SessionUUIDHeader = "X-Session-Uuid"

// IamClient это алиас для сгенерированного gRPC клиента
type IamClient = authV1.AuthServiceClient

// AuthMiddleware middleware для аутентификации HTTP запросов
type AuthMiddleware struct {
	iamClient IamClient
}

// NewAuthMiddleware создает новый middleware аутентификации
func NewAuthMiddleware(iamClient IamClient) *AuthMiddleware {
	return &AuthMiddleware{
		iamClient: iamClient,
	}
}

// client (X-Session-Uuid) -> auth middleware (add session_uuid in ctx (incomming)) -> order api (outgoing)-> auth interceptor ->inventory
// Handle обрабатывает HTTP запрос с аутентификацией
func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("headers: %v\n", r.Header)

		// Извлекаем session UUID из заголовка
		sessionUUID := r.Header.Get(SessionUUIDHeader)
		if sessionUUID == "" {
			writeErrorResponse(w, http.StatusUnauthorized, "MISSING_SESSION", "Authentication required")
			return
		}
		log.Printf("sessionUUID: %v\n", sessionUUID)

		// Валидируем сессию через IAM сервис
		whoamiRes, err := m.iamClient.Whoami(r.Context(), &authV1.WhoamiRequest{
			SessionUuid: sessionUUID,
		})
		if err != nil {
			log.Printf("whoamiRes err: %v\n", err)
			writeErrorResponse(w, http.StatusUnauthorized, "INVALID_SESSION", "Authentication failed")
			return
		}

		// Добавляем пользователя и session UUID в контекст используя функции из grpc middleware
		ctx := r.Context()
		ctx = grpcAuth.AddSessionUUIDToContext(ctx, sessionUUID)
		// Также добавляем пользователя в контекст
		ctx = grpcAuth.AddUserToContext(ctx, whoamiRes.User)

		// Передаем управление следующему handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext извлекает пользователя из контекста
func GetUserFromContext(ctx context.Context) (*commonV1.User, bool) {
	return grpcAuth.GetUserFromContext(ctx)
}

// GetSessionUUIDFromContext извлекает session UUID из контекста
func GetSessionUUIDFromContext(ctx context.Context) (string, bool) {
	return grpcAuth.GetSessionUUIDFromContext(ctx)
}
