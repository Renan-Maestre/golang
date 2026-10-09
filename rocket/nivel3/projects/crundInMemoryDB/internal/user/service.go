package user

import (
	"errors"

	"crudMemoryDB/internal/model"

	"github.com/google/uuid"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) FindAll() ([]model.User, error) {
	return s.repository.FindAll()
}

func (s *Service) FindByID(id string) (model.User, error) {
	return s.repository.FindByID(id)
}

func (s *Service) Create(input model.User) error {
	if input.FirstName == "" {
		return errors.New("first name is required")
	}

	if input.LastName == "" {
		return errors.New("last name is required")
	}

	if input.Biography == "" {
		return errors.New("biography is required")
	}

	id := uuid.NewString()

	user := model.User{
		ID:        id,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Biography: input.Biography,
	}

	if err := s.repository.Create(user); err != nil {
		return err
	}

	return nil
}

func (s *Service) Update(user model.User) error {
	if user.FirstName == "" {
		return errors.New("first name is required")
	}

	if user.LastName == "" {
		return errors.New("last name is required")
	}

	if user.Biography == "" {
		return errors.New("biography is required")
	}

	return s.repository.Update(user)
}

func (s *Service) Delete(id string) error {
	return s.repository.Delete(id)
}
