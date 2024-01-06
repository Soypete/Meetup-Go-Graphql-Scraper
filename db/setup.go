package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/marcboeker/go-duckdb"
)

type Warehouse struct {
	db *sql.DB
}

func Setup(filename string) (Warehouse, error) {
	db, err := sql.Open("duckdb", filename)
	if err != nil {
		return Warehouse{}, err
	}

	log.Println("Creating schema...")
	files, err := os.ReadDir("db/schema")
	if err != nil {
		return Warehouse{}, fmt.Errorf("could not read schema directory, %w", err)
	}
	for _, file := range files {
		fileBytes, err := os.ReadFile("db/schema/" + file.Name())
		if err != nil {
			return Warehouse{}, fmt.Errorf("could not read schema file, %w", err)
		}
		_, err = db.Exec(string(fileBytes))
		if err != nil {
			return Warehouse{}, fmt.Errorf("error runnign schema, %w", err)
		}
	}

	return Warehouse{
		db: db,
	}, nil
}
