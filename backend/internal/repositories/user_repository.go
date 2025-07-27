package repositories

import (
	"context"
	"fmt"

	"github.com/escuadron-404/red404/backend/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (email, password, created_at, updated_at) 
              VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRow(ctx, query, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).
		Scan(&user.ID)
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE id = $1`
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`
	user := &models.User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users ORDER BY id`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET email = $1, password = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(ctx, query, user.Email, user.Password, user.UpdatedAt, user.ID)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
