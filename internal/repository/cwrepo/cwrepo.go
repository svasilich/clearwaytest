package cwrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/svasilich/clearwaytest/domain/auth"
)

// Repository is struct for access to database.
type Repository struct {
	pool             *pgxpool.Pool
	connectionString string
}

// NewRepository create new instance of Repository.
func NewRepository(connectionString string) *Repository {
	return &Repository{
		connectionString: connectionString,
	}
}

// Connect init connection to db.
func (r *Repository) Connect(ctx context.Context) error {
	p, err := pgxpool.New(ctx, r.connectionString)
	if err != nil {
		return err
	}
	r.pool = p
	return nil
}

// Close connection to db.
func (r *Repository) Close() {
	r.pool.Close()
}

// Login сheck if the user exists in the database and create a new user session, or return an existing one.
func (r *Repository) Login(ctx context.Context, user string, passwordHash string) (auth.UserSession, error) {
	userID, err := r.getUserID(ctx, user, passwordHash)
	if err != nil {
		return auth.UserSession{}, err
	}

	err = r.createSession(ctx, userID)
	if err != nil {
		return auth.UserSession{}, fmt.Errorf("can't login %s: %w", user, err)
	}

	return r.GetLastSession(ctx, userID)
}

// GetLastSession returns the last active user session.
func (r *Repository) GetLastSession(ctx context.Context, userID int64) (auth.UserSession, error) {
	query := "SELECT id, created_at FROM sessions WHERE uid = @userID ORDER BY created_at DESC LIMIT 1"
	args := pgx.NamedArgs{
		"userID": userID,
	}

	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.UserSession{}, ErrNoOpenSessions
		}

		return auth.UserSession{}, fmt.Errorf("an error occurred while trying to get the user's session: %w", err)
	}
	defer rows.Close()

	session := auth.UserSession{}
	rows.Next()
	rows.Scan(&session.Token, &session.OpenedAt)
	return session, nil
}

func (r *Repository) getUserID(ctx context.Context, user string, passwordHash string) (int64, error) {
	query := "SELECT login FROM users WHERE login = @login AND password_hash = @passHash LIMIT 1"
	args := pgx.NamedArgs{
		"login":    user,
		"passHash": passwordHash,
	}

	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrUserNotExists
		}

		return 0, err
	}
	defer rows.Close()

	var userID int64
	rows.Next()
	rows.Scan(&userID)

	return userID, nil
}

func (r *Repository) createSession(ctx context.Context, userID int64) error {
	query := "INSERT INTO sessions (uid) VALUES (@userID)"
	args := pgx.NamedArgs{
		"userID": userID,
	}

	_, err := r.pool.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to create new session: %w", err)
	}

	return nil
}
