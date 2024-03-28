package db

import (
	"fmt"
	"log"
	"testing"
)

func TestGenericQuery(t *testing.T) {
	var database = `ATIWBT_USXEDIT.db`
	conn := NewDBAdapter(database)
	query := Select{conn.DB}
	sql1 := `SELECT book_id, chapter_num, script_begin_ts FROM scripts WHERE chapter_num=?`
	results, err := query.Select(sql1, 11)
	if err != nil {
		log.Fatal(err)
	}
	for _, result := range results {
		fmt.Println(result[0].(string), result[1].(int), result[2].(float64))
	}
}
