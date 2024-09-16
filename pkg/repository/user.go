package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/winens/enterprise-backend-template/pkg/errdefs"
	"github.com/winens/enterprise-backend-template/pkg/model"
	"github.com/winens/enterprise-backend-template/pkg/repository/interfaces"
	response "github.com/winens/enterprise-backend-template/pkg/repository/response"
	"github.com/winens/enterprise-backend-template/pkg/restapi/handler/request"
)

type DB interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type userRepo struct {
	db DB
}

func NewUserRepository(db *pgxpool.Pool) interfaces.UserRepository {
	return &userRepo{db}
}

func (r *userRepo) BeginTx(ctx context.Context, callback func(txRepo interfaces.UserRepository) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	txRepo := &userRepo{tx}

	if err := callback(txRepo); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("transaction failed (rollback) :%w", err)
	}

	return tx.Commit(ctx)
}

func (r *userRepo) FindUserById(ctx context.Context, id int64) (*model.User, error) {
	var user model.User

	err := r.db.QueryRow(ctx, `
	SELECT id, first_name, last_name, email, email_confirmed, created_at, updated_at
	FROM users WHERE id = $1`, id).
		Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.EmailConfirmed, &user.CreatedAt, &user.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errdefs.UserNotFound
	}

	return &user, err
}

func (r *userRepo) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := r.db.QueryRow(ctx, `
	SELECT id, first_name, last_name, email, email_confirmed, created_at, updated_at
	FROM users WHERE email = $1`, email).
		Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.EmailConfirmed, &user.CreatedAt, &user.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errdefs.UserNotFound
	}

	return &user, err
}

func (r *userRepo) FetchEmailPasswordLoginData(ctx context.Context, email string) (*response.UserFetchEmailPasswordLogin, error) {
	var res response.UserFetchEmailPasswordLogin

	err := r.db.QueryRow(ctx, `
	SELECT id, password_hash, email_confirmed
	FROM users WHERE email = $1`, email).
		Scan(&res.UserId, &res.PasswordHash, &res.EmailConfirmed)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errdefs.UserNotFound
	}

	return &res, err
}

func (r *userRepo) CreateUser(ctx context.Context, details request.SignUp) (userId int64, err error) {

	err = r.db.QueryRow(ctx, `
	INSERT INTO users (first_name, last_name, email, password_hash)
	VALUES ($1, $2, $3, $4)
	RETURNING id`,
		details.FirstName, details.LastName, details.Email, details.Password).Scan(&userId)

	return userId, err
}

func (r *userRepo) ConfirmEmailByUserId(ctx context.Context, userId int64) error {
	// check if already confirmed
	var emailConfirmed bool
	err := r.db.QueryRow(ctx, `
		SELECT email_confirmed FROM users WHERE id = $1`, userId).Scan(&emailConfirmed)

	if errors.Is(err, pgx.ErrNoRows) {
		return errdefs.UserNotFound
	}

	if emailConfirmed {
		return errdefs.EmailAlreadyConfirmed
	}

	// confirm email
	cmd, err := r.db.Exec(ctx, `
	UPDATE users SET email_confirmed = true, email_confirmed_at = NOW()
	WHERE id = $1`, userId)

	if err != nil {
		return fmt.Errorf("failed to confirm email: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return errdefs.UserNotFound
	}

	return nil
}
