package order

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Alexey-zaliznuak/gofermart/internal/config"
	"github.com/Alexey-zaliznuak/gofermart/internal/model"
	"github.com/Alexey-zaliznuak/gofermart/internal/repository/database"
)

type OrderRepository struct {
	db     *sql.DB
	config *config.AppConfig
	ctx    context.Context
}

// Получение заказа по ID
func (r *OrderRepository) GetByID(orderID string) (*model.Order, error) {
	order := &model.Order{}

	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, number, status, accrual, uploaded_at, user_id
		FROM orders
		WHERE id = $1
	`

	row := r.db.QueryRowContext(ctx, query, orderID)
	err := row.Scan(&order.ID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return order, nil
}

// Получение заказа по номеру
func (r *OrderRepository) GetByNumber(number string) (*model.Order, error) {
	order := &model.Order{}

	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, number, status, accrual, uploaded_at, user_id
		FROM orders
		WHERE number = $1
	`

	row := r.db.QueryRowContext(ctx, query, number)
	err := row.Scan(&order.ID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return order, nil
}

// Получение всех заказов пользователя по userID
func (r *OrderRepository) GetAllByUserID(userID int) ([]*model.Order, error) {
	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, number, status, accrual, uploaded_at, user_id
		FROM orders
		WHERE user_id = $1
		ORDER BY uploaded_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		order := &model.Order{}
		if err := rows.Scan(&order.ID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.UserID); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, database.ErrNotFound
	}

	return orders, nil
}


// Создание заказа
func (r *OrderRepository) CreateOrder(number string, userID int) (*model.Order, error) {
	order := &model.Order{}

	ctx, cancel := context.WithTimeout(r.ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO orders (number, status, accrual, user_id)
		VALUES ($1, $2, 0, $3)
		RETURNING id, number, status, accrual, uploaded_at, user_id
	`

	row := r.db.QueryRowContext(ctx, query, number, model.OrderStatusNew, userID)
	err := row.Scan(&order.ID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.UserID)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func NewOrderRepository(ctx context.Context, config *config.AppConfig, db *sql.DB) (*OrderRepository, error) {
	return &OrderRepository{
		db:     db,
		config: config,
		ctx:    ctx,
	}, nil
}
