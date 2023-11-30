package repository

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UserRepositoryImpl struct {
	db        *sql.DB
	timeoutMs int
}

type UserRepository interface {
	SaveUser(ctx context.Context, user User) error
}

func NewUserRepository(db *sql.DB, timeoutMs int) UserRepository {
	return &UserRepositoryImpl{
		db:        db,
		timeoutMs: timeoutMs,
	}
}

// SaveUser implements UserRepository.
func (r *UserRepositoryImpl) SaveUser(ctx context.Context, user User) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(r.timeoutMs)*time.Millisecond)
	defer cancel()

	statement := `INSERT INTO user(name, age) VALUES(?,?)`
	_, err := r.db.ExecContext(ctxTimeout, statement, user.Name, user.Age)
	if err != nil {
		return err
	}

	return nil
}
