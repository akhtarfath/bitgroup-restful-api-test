package products

import (
	"os"
	"testing"

	"github.com/akhtarfath/domain"
)

func chdirTemp(t *testing.T) {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })

	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
}

func TestStoreSeedsFileOnFirstRun(t *testing.T) {
	chdirTemp(t)

	s := NewStore()

	items, err := s.FindAll()
	if err != nil {
		t.Fatalf("FindAll: %v", err)
	}
	if len(items) != len(seedProducts) {
		t.Errorf("expected %d seeded products, got %d", len(seedProducts), len(items))
	}

	if _, err := os.Stat(dataFile); err != nil {
		t.Errorf("data file %q should be created: %v", dataFile, err)
	}
}

func TestStoreCRUDAndPersistence(t *testing.T) {
	chdirTemp(t)

	s := NewStore()
	cat := &domain.Category{ID: "1", Name: "Electronics"}

	// Create
	prod := domain.Product{ID: "p-new", Name: "Laptop", Price: 12.5, Category: cat}
	if err := s.Insert(prod); err != nil {
		t.Fatalf("Insert: %v", err)
	}

	got, err := s.FindByID("p-new")
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if got == nil || got.Name != "Laptop" || got.Price != 12.5 {
		t.Errorf("unexpected product after insert: %+v", got)
	}

	// Update
	updated := domain.Product{Name: "Laptop Pro", Price: 15, Category: cat}
	ok, err := s.Update("p-new", updated)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if !ok {
		t.Fatal("Update should report the product was found")
	}

	got, _ = s.FindByID("p-new")
	if got == nil || got.Name != "Laptop Pro" || got.Price != 15 {
		t.Errorf("unexpected product after update: %+v", got)
	}

	// Update of a missing ID reports not found without error
	if ok, _ = s.Update("missing", updated); ok {
		t.Error("Update of missing ID should report not found")
	}

	// Delete
	if ok, err = s.Delete("p-new"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if !ok {
		t.Fatal("Delete should report the product was found")
	}
	if ok, _ = s.Delete("p-new"); ok {
		t.Error("Delete of missing ID should report not found")
	}

	if got, _ = s.FindByID("p-new"); got != nil {
		t.Errorf("product should be gone, got %+v", got)
	}

	// Persistence: a fresh store reading the same file sees the final state
	s2 := NewStore()
	items, err := s2.FindAll()
	if err != nil {
		t.Fatalf("second store FindAll: %v", err)
	}
	if len(items) != len(seedProducts) {
		t.Errorf("expected only seeds to remain after delete, got %d items", len(items))
	}
	for _, item := range items {
		if item.ID == "p-new" {
			t.Error("deleted product still present after reload")
		}
	}
}