package output

import (
	"context"
	"dataset/db"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"gonum.org/v1/gonum/floats"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestCSVWriterScript(t *testing.T) {
	ctx := context.Background()
	var conn = db.NewDBAdapter(ctx, `ENGWEB_DBPTEXT.db`)
	//prepareTimestampAndFMCCData(conn, `ENGWEB`, `ENGWEBN2DA`, t)
	structs, meta := PrepareScripts(conn, false, false)
	fmt.Println("Loaded Scripts", len(structs))
	filename := WriteCSV(structs, meta)
	mfccFile := readCSVMFCC(filename, t)
	mfccDB := readDBMFCC(conn, `script_mfcc`, t)
	compare(mfccDB, mfccFile, t)
	fmt.Println("Written CSV", filename)
}

func readCSVMFCC(filename string, t *testing.T) [][]float64 {
	var results = make([][]float64, 0, 10000)
	file, _ := os.Open(filename)
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Error(err)
	}
	header := records[0]
	var columns []int
	for i, col := range header {
		if strings.HasPrefix(col, `mfcc`) {
			columns = append(columns, i)
		}
	}
	for row, record := range records {
		if row > 0 {
			var hasData = false
			var resultRec []float64
			for _, index := range columns {
				value, err := strconv.ParseFloat(record[index], 64)
				if err != nil {
					value = 0.0
				} else {
					hasData = true
				}
				resultRec = append(resultRec, value)
			}
			if hasData {
				results = append(results, resultRec)
			}
		}
	}
	return results
}

func readDBMFCC(conn db.DBAdapter, table string, t *testing.T) [][]float64 {
	var results = make([][]float64, 0, 10000)
	var query string
	if table == `script_mfcc` {
		query = `SELECT script_id, rows, cols, mfcc_json FROM script_mfcc ORDER BY script_id`
	} else {
		query = `SELECT word_id, rows, cols, mfcc_json FROM script_mfcc ORDER BY word_id`
	}
	rows, err := conn.DB.Query(query)
	defer rows.Close()
	if err != nil {
		t.Error(err)
	}
	for rows.Next() {
		var rec db.MFCC
		var mfccJson string
		err := rows.Scan(&rec.Id, &rec.Rows, &rec.Cols, &mfccJson)
		if err != nil {
			t.Error(err)
		}
		err = json.Unmarshal([]byte(mfccJson), &rec.MFCC)
		if err != nil {
			t.Error(err)
		}
		if rec.Rows > 0 {
			for row := 0; row < rec.Rows; row++ {
				var resultRec []float64
				for _, value := range rec.MFCC[row] {
					resultRec = append(resultRec, float64(value))
				}
				results = append(results, resultRec)
			}
		}
	}
	err = rows.Err()
	if err != nil {
		t.Error(err)
	}
	return results
}

func testSumMWordEnc(conn db.DBAdapter, t *testing.T) map[string]float64 {
	var results = make(map[string]float64)
	query := `SELECT word_enc FROM words`
	rows, err := conn.DB.Query(query)
	defer rows.Close()
	if err != nil {
		t.Error(err)
	}
	for rows.Next() {
		var wordJson string
		err := rows.Scan(&wordJson)
		if err != nil {
			t.Error(err)
		}
		var wordEnc []float64
		err = json.Unmarshal([]byte(wordJson), &wordEnc)
		if err != nil {
			t.Error(err)
		}
		for col, wd := range wordEnc {
			name := `word_enc` + strconv.Itoa(col)
			total, _ := results[name]
			results[name] = wd + total
		}
	}
	err = rows.Err()
	if err != nil {
		t.Error(err)
	}
	return results
}

func compare(dbSlice [][]float64, fileSlice [][]float64, t *testing.T) {
	if len(dbSlice) != len(fileSlice) {
		t.Error(`dbSlice has length`, len(dbSlice), `fileSlice has length`, len(fileSlice))
	}
	fmt.Println(`Compare `, len(dbSlice), `records`)
	var rowDiffCount = 0
	for i, dbRec := range dbSlice {
		fileRec := fileSlice[i]
		if len(dbRec) != len(fileRec) {
			t.Error(`row`, i, `dbRec len`, len(dbRec), `fileRec len`, len(fileRec))
		} else if !floats.EqualApprox(dbRec, fileRec, 0.0000001) {
			t.Error(`row`, i, `dbRec`, dbRec, `fileRec`, fileRec)
			rowDiffCount++
		}
	}
	fmt.Println("rows", len(dbSlice), `differences`, rowDiffCount)
}
