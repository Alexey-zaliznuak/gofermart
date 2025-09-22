package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Alexey-zaliznuak/gofermart/internal/config"
	"github.com/Alexey-zaliznuak/gofermart/internal/model"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository/database"
)

type UserRepository struct {
	db     *sql.DB
	config *config.AppConfig
	ctx    context.Context
}

// Получение пользователя по ID
func (r *UserRepository) GetByID(userID string) (*model.User, error) {
	user := &model.User{}

	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, username, password_hash, balance, withdraw
		FROM users
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, userID)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Balance, &user.Withdraw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

// Получение пользователя по username
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	user := &model.User{}

	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, username, password_hash, balance, withdraw
		FROM users
		WHERE username = $1
	`

	row := r.db.QueryRowContext(ctx, query, username)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Balance, &user.Withdraw)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) CreateUser(username, passwordHash string) (*model.User, error) {
	user := &model.User{}

	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (username, password_hash, balance, withdraw)
		VALUES ($1, $2, 0, 0)
		RETURNING id, username, password_hash, balance, withdraw
	`

	row := r.db.QueryRowContext(ctx, query, username, passwordHash)
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Balance, &user.Withdraw)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewUserRepository(ctx context.Context, config *config.AppConfig, db *sql.DB) (*UserRepository, error) {
	return &UserRepository{
		db:     db,
		config: config,
		ctx:    ctx,
	}, nil
}
