package db

import (
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	// Read the actual embedded DB file for testing
	embeddedDB, err := os.ReadFile("../../data/ghofig.db")
	if err != nil {
		t.Fatalf("Failed to read test db: %v", err)
	}

	if err := Init(embeddedDB); err != nil {
		t.Fatalf("Failed to init db: %v", err)
	}
	defer Close()

	// Test search with results
	results, err := Search("font")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) == 0 {
		t.Error("Expected results for 'font' query")
	}

	// Verify title matches come first
	foundFontFamily := false
	for i, r := range results {
		if r.Title == "font-family" {
			foundFontFamily = true
			if i > 10 {
				t.Errorf("font-family should be near top of results, got index %d", i)
			}
			break
		}
	}
	if !foundFontFamily {
		t.Error("Expected font-family in results")
	}

	// Test GetByID
	config, err := GetByID(results[0].ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if config.Title != results[0].Title {
		t.Errorf("GetByID returned wrong config: got %s, want %s", config.Title, results[0].Title)
	}
}
