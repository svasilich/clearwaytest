package cwrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/svasilich/clearwaytest/domain/auth"
	"github.com/svasilich/clearwaytest/domain/common"
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

	return r.getLastSession(ctx, userID)
}

// GetUserBySession return user ID if given session is active.
func (r *Repository) GetUserBySession(ctx context.Context, token auth.Token) (int64, error) {
	id, err := r.getUserIDBySession(ctx, token)
	if err != nil {
		return 0, err
	}

	activeSession, err := r.getLastSession(ctx, id)
	if err != nil {
		return 0, err
	}

	if token != activeSession.Token {
		return 0, common.ErrNoOpenSessions
	}

	return id, nil
}

// WriteAsset store asset data to data base.
func (r *Repository) WriteAsset(ctx context.Context, asset string, uid int64, data []byte) error {
	query := "INSERT INTO assets (name, uid, data) VALUES (@asset, @uid, @data) ON CONFLICT (name, uid) DO UPDATE SET data = @data"
	args := pgx.NamedArgs{
		"asset": asset,
		"uid":   uid,
		"data":  data,
	}

	_, err := r.pool.Exec(ctx, query, args)
	return err
}

// ReadAsset read asset data from data base.
func (r *Repository) ReadAsset(ctx context.Context, asset string, uid int64) ([]byte, error) {
	query := "SELECT data FROM assets WHERE name = @asset AND uid = @uid LIMIT 1"
	args := pgx.NamedArgs{
		"asset": asset,
		"uid":   uid,
	}

	var data []byte
	err := r.pool.QueryRow(ctx, query, args).Scan(&data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []byte{}, common.ErrForbiddenAsset
		}

		return []byte{}, err
	}

	return data, nil
}

func (r *Repository) getLastSession(ctx context.Context, userID int64) (auth.UserSession, error) {
	query := "SELECT id, created_at FROM sessions WHERE uid = @userID ORDER BY created_at DESC LIMIT 1"
	args := pgx.NamedArgs{
		"userID": userID,
	}

	session := auth.UserSession{}
	err := r.pool.QueryRow(ctx, query, args).Scan(&session.Token, &session.OpenedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth.UserSession{}, common.ErrNoOpenSessions
		}

		return auth.UserSession{}, fmt.Errorf("an error occurred while trying to get the user's session: %w", err)
	}
	return session, nil
}

func (r *Repository) getUserID(ctx context.Context, user string, passwordHash string) (int64, error) {
	query := "SELECT login FROM users WHERE login = @login AND password_hash = @passHash LIMIT 1"
	args := pgx.NamedArgs{
		"login":    user,
		"passHash": passwordHash,
	}

	var userID int64
	err := r.pool.QueryRow(ctx, query, args).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, common.ErrUserNotExists
		}

		return 0, err
	}
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

func (r *Repository) getUserIDBySession(ctx context.Context, token auth.Token) (int64, error) {
	query := "SELECT uid FROM sessions WHERE id = @id LIMIT 1"
	args := pgx.NamedArgs{
		"id": string(token),
	}

	var uid int64
	err := r.pool.QueryRow(ctx, query, args).Scan(&uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, common.ErrNoOpenSessions
		}

		return 0, err
	}
	return uid, nil
}
