package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/winens/enterprise-backend-template/pkg/model"
)

// birşeyleri kesinleştirelim.
// Authentication = kimlik doğrulama, sadece kimlik doğrulama işlemleri yapılır.

type AuthRepository interface {
	NewSession(ctx context.Context, userId int64, ip, userAgent string) (sessionId uuid.UUID, err error)
	FindSessionById(ctx context.Context, sessionId uuid.UUID) (*model.Session, error)
	DeleteSession(ctx context.Context, sessionId uuid.UUID) error
}
