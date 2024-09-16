package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/winens/enterprise-backend-template/pkg/model"
	"github.com/winens/enterprise-backend-template/pkg/repository/interfaces"
)

type authRepo struct {
	db DB
}

func NewAuthRepository(db *pgxpool.Pool) interfaces.AuthRepository {
	return &authRepo{db}
}

// NewSession yeni bir session oluşturacak ve bu session_id'yi ChaCha20Poly1308 ile şifreleyip string olarak işleyip geri döndürecek.
// Böylece session_id brute-force saldırılarına karşı korunmuş olacak.
func (r *authRepo) NewSession(ctx context.Context, userId int64, ip, userAgent string) (sessionId uuid.UUID, err error) {
	sessionId, err = uuid.NewRandom()
	if err != nil {
		return uuid.Nil, err
	}

	_, err = r.db.Exec(ctx, `
	INSERT INTO sessions (id, user_id, ip, user_agent) VALUES ($1, $2, $3, $4)
	`, sessionId, userId, ip, userAgent)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert new session: %w", err)
	}

	return sessionId, nil

}

func (r *authRepo) FindSessionById(ctx context.Context, sessionId uuid.UUID) (*model.Session, error) {
	var sess model.Session

	err := r.db.QueryRow(ctx, `
	SELECT id, user_id, ip, user_agent, created_at, last_seen_at
	FROM sessions WHERE id = $1`, sessionId).
		Scan(&sess.Id, &sess.UserId, &sess.IP, &sess.UserAgent, &sess.CreatedAt, &sess.LastSeenAt)

	if err != nil {
		return nil, err
	}

	// USAGE: usecase(){ userRepo.FindUserById(sess.UserId) }

	return &sess, nil
}

func (r *authRepo) DeleteSession(ctx context.Context, sessionId uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM sessions WHERE id = $1`, sessionId)
	return err
}
