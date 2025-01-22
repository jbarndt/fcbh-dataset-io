package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
	"unicode/utf8"
)

type QueryProcessor struct {
	db *sql.DB
}

func NewQueryProcessor(db *sql.DB) *QueryProcessor {
	return &QueryProcessor{db: db}
}

// ProcessQuery executes SELECT query and prints results in table format
func (qp *QueryProcessor) ProcessQuery(query string, args ...interface{}) error {
	// Verify query starts with SELECT
	if !strings.HasPrefix(strings.TrimSpace(strings.ToUpper(query)), "SELECT") {
		return fmt.Errorf("only SELECT queries are allowed")
	}

	// Execute query
	rows, err := qp.db.Query(query, args...)
	if err != nil {
		return fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get column names: %w", err)
	}

	// Prepare value holders
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	// Collect all rows as strings
	var allRows [][]string
	columnWidths := make([]int, len(columns))

	// Initialize column widths with header lengths
	for i, col := range columns {
		columnWidths[i] = utf8.RuneCountInString(col)
	}

	// Collect rows and calculate column widths
	rowCount := 0
	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("row scan failed: %w", err)
		}

		// Convert row values to strings
		rowStrings := make([]string, len(columns))
		for i, val := range values {
			str := interfaceToString(val)
			rowStrings[i] = str
			// Update column width if this value is longer
			if width := utf8.RuneCountInString(str); width > columnWidths[i] {
				columnWidths[i] = width
			}
		}
		allRows = append(allRows, rowStrings)
		rowCount++
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("error during row iteration: %w", err)
	}

	// Print table
	printTableRow(columns, columnWidths, true) // Header
	printHorizontalLine(columnWidths)          // Separator

	// Print data rows
	for _, row := range allRows {
		printTableRow(row, columnWidths, false)
	}

	// Print row count
	fmt.Printf("\n%d rows in set\n", rowCount)
	return nil
}

// printTableRow prints a single row with proper padding
func printTableRow(values []string, widths []int, isHeader bool) {
	var parts []string
	for i, v := range values {
		// Right-pad the value with spaces to match column width
		format := fmt.Sprintf("%%-%ds", widths[i])
		parts = append(parts, fmt.Sprintf(format, v))
	}
	if isHeader {
		fmt.Printf("| %s |\n", strings.Join(parts, " | "))
	} else {
		fmt.Printf("| %s |\n", strings.Join(parts, " | "))
	}
}

// printHorizontalLine prints separator line
func printHorizontalLine(widths []int) {
	var parts []string
	for _, w := range widths {
		parts = append(parts, strings.Repeat("-", w))
	}
	fmt.Printf("+-%s-+\n", strings.Join(parts, "-+-"))
}

// interfaceToString safely converts any value to string
func interfaceToString(v interface{}) string {
	if v == nil {
		return "NULL"
	}
	switch v := v.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case bool:
		if v {
			return "1"
		}
		return "0"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func GetDBPMySqlDSN() string {
	// Format: username:password@tcp(hostname:port)/database_name
	var result string
	username := os.Getenv("DBP_MYSQL_USERNAME")
	password := os.Getenv("DBP_MYSQL_PASSWORD")
	host := os.Getenv("DBP_MYSQL_HOST")
	port := os.Getenv("DBP_MYSQL_PORT")
	database := os.Getenv("DBP_MYSQL_DATABASE")
	result = username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database
	return result
}

func main() {
	// Example usage:
	dsn := GetDBPMySqlDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	processor := NewQueryProcessor(db)

	// Example query
	query := `Select * from bible_files limit ?`

	if err := processor.ProcessQuery(query, 25); err != nil {
		panic(err)
	}

	// Output will look like:
	// +----+----------+-------------------+---------------------+
	// | id | name     | email            | created_at          |
	// +----+----------+-------------------+---------------------+
	// | 1  | John Doe | john@example.com | 2024-01-20 15:30:00|
	// | 2  | Jane Doe | jane@example.com | 2024-01-21 09:45:00|
	// +----+----------+-------------------+---------------------+
	// 2 rows in set
}
