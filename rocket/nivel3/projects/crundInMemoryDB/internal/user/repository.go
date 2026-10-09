package user

import "crudMemoryDB/internal/model"

type Repository interface {
	FindAll() ([]model.User, error)
	FindByID(id string) (model.User, error)
	Create(user model.User) error
	Update(user model.User) error
	Delete(id string) error
}
