package categories

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
	if len(items) != len(seedCategories) {
		t.Errorf("expected %d seeded categories, got %d", len(seedCategories), len(items))
	}

	if _, err := os.Stat(dataFile); err != nil {
		t.Errorf("data file %q should be created: %v", dataFile, err)
	}
}

func TestStoreCRUDAndPersistence(t *testing.T) {
	chdirTemp(t)

	s := NewStore()

	// Create
	cat := domain.Category{ID: "c-new", Name: "Books"}
	if err := s.Insert(cat); err != nil {
		t.Fatalf("Insert: %v", err)
	}

	got, err := s.FindByID("c-new")
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}
	if got == nil || got.Name != "Books" {
		t.Errorf("unexpected category after insert: %+v", got)
	}

	// Update
	ok, err := s.Update("c-new", domain.Category{Name: "BookStore"})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if !ok {
		t.Fatal("Update should report the category was found")
	}

	got, _ = s.FindByID("c-new")
	if got == nil || got.Name != "BookStore" {
		t.Errorf("unexpected category after update: %+v", got)
	}

	if ok, _ = s.Update("missing", domain.Category{Name: "X"}); ok {
		t.Error("Update of missing ID should report not found")
	}

	// Delete
	if ok, err = s.Delete("c-new"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if !ok {
		t.Fatal("Delete should report the category was found")
	}
	if ok, _ = s.Delete("c-new"); ok {
		t.Error("Delete of missing ID should report not found")
	}

	if got, _ = s.FindByID("c-new"); got != nil {
		t.Errorf("category should be gone, got %+v", got)
	}

	// Persistence: a fresh store reading the same file sees the final state
	s2 := NewStore()
	items, err := s2.FindAll()
	if err != nil {
		t.Fatalf("second store FindAll: %v", err)
	}
	if len(items) != len(seedCategories) {
		t.Errorf("expected only seeds to remain after delete, got %d items", len(items))
	}
	for _, item := range items {
		if item.ID == "c-new" {
			t.Error("deleted category still present after reload")
		}
	}
}