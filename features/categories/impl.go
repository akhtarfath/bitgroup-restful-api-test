package categories

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/akhtarfath/domain"
	"github.com/spf13/viper"
)

const dataFile = "data/categories.json"

var seedCategories = []domain.Category{}

type data struct {
	Categories []domain.Category `json:"categories"`
}

type Store struct {
	mu sync.Mutex
}

func NewStore() *Store {
	s := &Store{}
	s.ensureFile()
	return s
}

func (s *Store) ensureFile() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := os.Stat(dataFile); err == nil {
		return
	}

	_ = s.writeLocked(seedCategories)
}

func (s *Store) readLocked() ([]domain.Category, error) {
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return []domain.Category{}, nil
	}

	v := viper.New()
	v.SetConfigFile(dataFile)
	v.SetConfigType("json")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var d data
	if err := v.Unmarshal(&d); err != nil {
		return nil, err
	}

	return d.Categories, nil
}

func (s *Store) writeLocked(items []domain.Category) error {
	b, err := json.MarshalIndent(data{Categories: items}, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dataFile), 0o755); err != nil {
		return err
	}

	return os.WriteFile(dataFile, b, 0o644)
}

func (s *Store) FindAll() ([]domain.Category, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.readLocked()
}

func (s *Store) FindByID(id string) (*domain.Category, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	items, err := s.readLocked()
	if err != nil {
		return nil, err
	}

	for i := range items {
		if items[i].ID == id {
			return &items[i], nil
		}
	}

	return nil, nil
}

func (s *Store) Insert(category domain.Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	items, err := s.readLocked()
	if err != nil {
		return err
	}

	items = append(items, category)
	return s.writeLocked(items)
}

func (s *Store) Update(id string, category domain.Category) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	items, err := s.readLocked()
	if err != nil {
		return false, err
	}

	updated := false
	for i := range items {
		if items[i].ID == id {
			category.ID = id
			items[i] = category
			updated = true
			break
		}
	}

	if !updated {
		return false, nil
	}

	return true, s.writeLocked(items)
}

func (s *Store) Delete(id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	items, err := s.readLocked()
	if err != nil {
		return false, err
	}

	kept := make([]domain.Category, 0, len(items))
	deleted := false
	for _, item := range items {
		if item.ID == id {
			deleted = true
			continue
		}
		kept = append(kept, item)
	}

	if !deleted {
		return false, nil
	}

	return true, s.writeLocked(kept)
}
