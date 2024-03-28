package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
)

// db is a global database connection.
// Make sure to initialize it with sql.Open() before using executeQuery.

type Select struct {
	conn *sql.DB
}

// executeQuery executes a SQL query and returns the results as a slice of slice of interface{}.
// This can then be processed according to the expected types.
func (s *Select) Select(query string, args ...interface{}) ([][]interface{}, error) {
	rows, err := s.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	dtypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	fmt.Println("TYPES", dtypes)
	for _, dtype := range dtypes {
		fmt.Println("TYPE", dtype)
		fmt.Println("Nme", dtype.Name())
		fmt.Print("Decimal size ")
		fmt.Println(dtype.DecimalSize())
		fmt.Println("Length ")
		fmt.Println(dtype.Length())
		fmt.Println("type name", dtype.DatabaseTypeName())
		fmt.Print("Nullable")
		fmt.Println(dtype.Nullable())
		fmt.Println("scan type", dtype.ScanType())
	}

	// Prepare a slice of interfaces to hold each value.
	rawResult := make([][]byte, len(cols))
	//result := make([]interface{}, len(cols))
	result := make([]any, len(cols))
	for i, _ := range rawResult {
		result[i] = &rawResult[i] // Pointers to each slice cell
	}
	var results [][]any
	for rows.Next() {
		err = rows.Scan(result...)
		if err != nil {
			return nil, err
		}
		// Convert raw bytes to a new slice of interfaces that holds real types.
		// You might need type assertion here based on your schema.
		//record := make([]interface{}, len(cols))
		record := make([]any, len(cols))
		for i, raw := range rawResult {
			fmt.Println("RAW", len(raw), raw)
			if raw == nil {
				record[i] = nil
			} else {
				switch dtypes[i].ScanType().String() {
				case `string`:
					record[i] = string(raw)
				case `int64`:
					record[i], _ = strconv.Atoi(string(raw))
				case `float64`:
					record[i], _ = strconv.ParseFloat(string(raw), 64)
				default:
					os.Exit(1)
				}
			}
		}
		results = append(results, record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
