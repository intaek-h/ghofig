package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"

	_ "modernc.org/sqlite"
)

const (
	inputFile  = "reference.mdx.txt"
	outputFile = "data/ghofig.db"
)

// h2Pattern matches lines like: ## `config-name`
var h2Pattern = regexp.MustCompile("^## `(.+)`$")

type configEntry struct {
	title       string
	description string
}

func main() {
	entries, err := parseFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Parsed %d config entries\n", len(entries))

	if err := writeDatabase(outputFile, entries); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing database: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Database written to %s\n", outputFile)
}

func parseFile(filename string) ([]configEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []configEntry
	var pendingTitles []string
	var descriptionLines []string
	inDescription := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if line is an h2 header
		if match := h2Pattern.FindStringSubmatch(line); match != nil {
			// If we were building a description, flush it
			if inDescription && len(pendingTitles) > 0 {
				description := strings.TrimSpace(strings.Join(descriptionLines, "\n"))
				for _, title := range pendingTitles {
					entries = append(entries, configEntry{
						title:       title,
						description: description,
					})
				}
				pendingTitles = nil
				descriptionLines = nil
			}

			// Add this title to pending
			pendingTitles = append(pendingTitles, match[1])
			inDescription = false
			continue
		}

		// If we have pending titles and hit non-empty content, start description
		if len(pendingTitles) > 0 {
			// Skip empty lines between h2 and description start
			if !inDescription && strings.TrimSpace(line) == "" {
				continue
			}
			inDescription = true
			descriptionLines = append(descriptionLines, line)
		}
	}

	// Flush any remaining entry
	if len(pendingTitles) > 0 && len(descriptionLines) > 0 {
		description := strings.TrimSpace(strings.Join(descriptionLines, "\n"))
		for _, title := range pendingTitles {
			entries = append(entries, configEntry{
				title:       title,
				description: description,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func writeDatabase(filename string, entries []configEntry) error {
	// Remove existing database
	os.Remove(filename)

	db, err := sql.Open("sqlite", filename)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS configs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_configs_title ON configs(title);
	`)
	if err != nil {
		return err
	}

	// Insert entries
	stmt, err := db.Prepare("INSERT INTO configs (title, description) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, entry := range entries {
		_, err = stmt.Exec(entry.title, entry.description)
		if err != nil {
			return err
		}
	}

	return nil
}
