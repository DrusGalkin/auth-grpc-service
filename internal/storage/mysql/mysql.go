package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DrusGalkin/auth-grpc-service/internal/domain/models"
	"github.com/DrusGalkin/auth-grpc-service/internal/storage"
	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) *Storage {
	db, err := sql.Open("mysql", storagePath)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	return &Storage{
		db: db,
	}
}

func (s *Storage) SaveUser(ctx context.Context, email string, username string, hashPassword []byte) (uid int64, err error) {
	const op = "mysql.SaveUser"

	stmt, err := s.db.Prepare(`INSERT INTO users (username, email, password_hash, role) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()
	result, err := stmt.ExecContext(ctx, username, email, hashPassword, "user")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExist)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "mysql.User"

	stmt, err := s.db.Prepare(`SELECT * FROM users WHERE email = ?`)
	defer stmt.Close()

	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models.User

	if err = row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.HashPassword,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, storage.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, uid int64) (bool, error) {
	const op = "mysql.IsAdmin"

	stmt, err := s.db.Prepare(`SELECT role FROM users WHERE id = ?`)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, uid)

	var role string

	if err := row.Scan(&role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return role == "admin", nil
}
