package store

import (
	"errors"

	"crudMemoryDB/internal/model"
)

type UserMemoryRepository struct {
	db map[string]model.User
}

func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		db: make(map[string]model.User),
	}
}

func (r *UserMemoryRepository) FindAll() ([]model.User, error) {
	users := make([]model.User, 0, len(r.db))

	for _, user := range r.db {
		users = append(users, user)
	}

	return users, nil
}

func (r *UserMemoryRepository) FindByID(id string) (model.User, error) {
	user, ok := r.db[id]

	if !ok {
		return model.User{}, errors.New("user not found")
	}

	return user, nil
}

func (r *UserMemoryRepository) Create(user model.User) error {
	if _, exists := r.db[user.ID]; exists {
		return errors.New("user already exists")
	}

	r.db[user.ID] = user

	return nil
}

func (r *UserMemoryRepository) Update(user model.User) error {
	if _, existes := r.db[user.ID]; !existes {
		return errors.New("user not found")
	}
	r.db[user.ID] = user
	return nil
}

func (r *UserMemoryRepository) Delete(id string) error {
	if _, existes := r.db[id]; !existes {
		return errors.New("user not found")
	}
	delete(r.db, id)
	return nil
}
