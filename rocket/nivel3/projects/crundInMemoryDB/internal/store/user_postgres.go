package store

import (
	"database/sql"
	"errors"

	"crudMemoryDB/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserPostgresRepository struct {
	db *sqlx.DB
}

func NewUserPostgresRepository(db *sqlx.DB) *UserPostgresRepository {
	return &UserPostgresRepository{
		db: db,
	}
}

func (r *UserPostgresRepository) FindAll() ([]model.User, error) {
	query := `SELECT id , first_name, last_name, biography FROM users`

	users := make([]model.User, 0)
	if err := r.db.Select(&users, query); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserPostgresRepository) FindByID(id string) (model.User, error) {
	query := `SELECT id, first_name, last_name, biography FROM users WHERE id = $1`

	var u model.User
	if err := r.db.Get(&u, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, errors.New("user not found")
		}
	}
	return u, nil
}

func (r *UserPostgresRepository) Create(u model.User) error {
	query := `
			INSERT INTO users (id, first_name,last_name,biography)
			VALUES (:id,:first_name,:last_name,:biography)
	`
	_, err := r.db.NamedExec(query, u)
	return err
}

func (r *UserPostgresRepository) Update(u model.User) error {
	query := `
		UPDATE users
		SET first_name = :first_name, last_name = :last_name, biography = :biography
		WHERE id = :id
	`

	res, err := r.db.NamedExec(query, u)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *UserPostgresRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}
