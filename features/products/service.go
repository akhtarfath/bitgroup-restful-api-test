package products

import (
	"errors"

	"github.com/akhtarfath/domain"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("product not found")

type CreateInput struct {
	Name     string
	Price    float64
	Category *domain.Category
}

type UpdateInput struct {
	Name     string
	Price    float64
	Category *domain.Category
}

type Service struct {
	store *Store
}

func NewService(store *Store) *Service {
	return &Service{store: store}
}

func (s *Service) List() ([]domain.Product, error) {
	return s.store.FindAll()
}

func (s *Service) Get(id string) (*domain.Product, error) {
	product, err := s.store.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrNotFound
	}

	return product, nil
}

func (s *Service) Create(in CreateInput) (*domain.Product, error) {
	product := domain.Product{
		ID:       uuid.NewString(),
		Name:     in.Name,
		Price:    in.Price,
		Category: in.Category,
	}

	if err := s.store.Insert(product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *Service) Update(id string, in UpdateInput) (*domain.Product, error) {
	product := domain.Product{
		ID:       id,
		Name:     in.Name,
		Price:    in.Price,
		Category: in.Category,
	}

	updated, err := s.store.Update(id, product)
	if err != nil {
		return nil, err
	}
	if !updated {
		return nil, ErrNotFound
	}

	return &product, nil
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
