package categories

import (
	"errors"

	"github.com/akhtarfath/domain"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("category not found")

type CreateInput struct {
	Name string
}

type UpdateInput struct {
	Name string
}

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) List() ([]domain.Category, error) {
	return s.store.FindAll()
}

func (s *Service) Get(id string) (*domain.Category, error) {
	category, err := s.store.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrNotFound
	}

	return category, nil
}

func (s *Service) Create(in CreateInput) (*domain.Category, error) {
	category := domain.Category{
		ID:   uuid.NewString(),
		Name: in.Name,
	}

	if err := s.store.Insert(category); err != nil {
		return nil, err
	}

	return &category, nil
}

func (s *Service) Update(id string, in UpdateInput) (*domain.Category, error) {
	category := domain.Category{
		ID:   id,
		Name: in.Name,
	}

	updated, err := s.store.Update(id, category)
	if err != nil {
		return nil, err
	}
	if !updated {
		return nil, ErrNotFound
	}

	return &category, nil
}

func (s *Service) Delete(id string) error {
	deleted, err := s.store.Delete(id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}

	return nil
}
