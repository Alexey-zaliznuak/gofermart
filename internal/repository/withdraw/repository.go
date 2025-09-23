package withdraw

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Alexey-zaliznuak/gofermart/internal/config"
	"github.com/Alexey-zaliznuak/gofermart/internal/model"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository/database"
)

type WithdrawRepository struct {
	db     *sql.DB
	config *config.AppConfig
	ctx    context.Context
}

// Получение всех выводов пользователя по userID
func (r *WithdrawRepository) GetAllByUserID(userID int) ([]*model.Withdraw, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, number, sum, processed_at, user_id
		FROM withdrawals
		WHERE user_id = $1
		ORDER BY processed_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []*model.Withdraw
	for rows.Next() {
		w := &model.Withdraw{}
		if err := rows.Scan(&w.ID, &w.Number, &w.Sum, &w.ProcessedAt, &w.UserID); err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(withdrawals) == 0 {
		return nil, database.ErrNotFound
	}

	return withdrawals, nil
}

// Получение вывода по номеру
func (r *WithdrawRepository) GetByNumber(number string) (*model.Withdraw, error) {
	w := &model.Withdraw{}

	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, number, sum, processed_at, user_id
		FROM withdrawals
		WHERE number = $1
	`

	row := r.db.QueryRowContext(ctx, query, number)
	err := row.Scan(&w.ID, &w.Number, &w.Sum, &w.ProcessedAt, &w.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return w, nil
}

// Создание вывода
func (r *WithdrawRepository) CreateWithdraw(number string, sum int64, userID int) (*model.Withdraw, error) {
	w := &model.Withdraw{}

	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO withdrawals (number, sum, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, number, sum, processed_at, user_id
	`

	row := r.db.QueryRowContext(ctx, query, number, sum, userID)
	err := row.Scan(&w.ID, &w.Number, &w.Sum, &w.ProcessedAt, &w.UserID)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func NewWithdrawRepository(ctx context.Context, config *config.AppConfig, db *sql.DB) (*WithdrawRepository, error) {
	return &WithdrawRepository{
		db:     db,
		config: config,
		ctx:    ctx,
	}, nil
}
