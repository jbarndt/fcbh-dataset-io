package db

import (
	"database/sql"
	"fmt"
	"log"
)

// db is a global database connection.
// Make sure to initialize it with sql.Open() before using executeQuery.
var db *sql.DB

// executeQuery executes a SQL query and returns the results as a slice of slice of interface{}.
// This can then be processed according to the expected types.
func executeQuery(query string, args ...interface{}) ([][]interface{}, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Prepare a slice of interfaces to hold each value.
	rawResult := make([][]byte, len(cols))
	result := make([]interface{}, len(cols))
	for i, _ := range rawResult {
		result[i] = &rawResult[i] // Pointers to each slice cell
	}

	var results [][]interface{}
	for rows.Next() {
		err = rows.Scan(result...)
		if err != nil {
			return nil, err
		}
		// Convert raw bytes to a new slice of interfaces that holds real types.
		// You might need type assertion here based on your schema.
		record := make([]interface{}, len(cols))
		for i, raw := range rawResult {
			if raw == nil {
				record[i] = nil
			} else {
				record[i] = string(raw) // Simple conversion; customize as needed.
			}
		}
		results = append(results, record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func main() {
	// Example usage
	results, err := executeQuery("SELECT id, name FROM users WHERE id > ?", 1)
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range results {
		// Type assertions and error checks should be handled properly in production code.
		id := row[0].(string) // Assuming the ID is a string
		name := row[1].(string)
		fmt.Println("ID:", id, "Name:", name)
	}
}
