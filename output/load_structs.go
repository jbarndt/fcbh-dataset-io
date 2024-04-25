package output

import (
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"encoding/json"
	"math"
	"reflect"
	"strings"
)

func LoadScriptStruct(d db.DBAdapter) []Script {
	var results []Script
	var status dataset.Status
	query := `SELECT scripts.script_id, book_id, chapter_num, chapter_end, audio_file, script_num, 
		usfm_style, person, actor, verse_str, verse_end, script_text, 
		script_begin_ts, script_end_ts, rows, cols, mfcc_json
		FROM scripts LEFT OUTER JOIN script_mfcc ON script_mfcc.script_id = scripts.script_id
		WHERE book_id = 'MRK' AND chapter_num = 1 ORDER BY scripts.script_id`
	rows, err := d.DB.Query(query)
	if err != nil {
		status = log.Error(d.Ctx, 500, err, "Error during select scripts")
		panic(status.Message)
	}
	defer rows.Close()
	for rows.Next() {
		var sc Script
		var mfccJson string
		err := rows.Scan(&sc.ScriptId, &sc.BookId, &sc.ChapterNum, &sc.ChapterEnd, &sc.AudioFile,
			&sc.ScriptNum, &sc.UsfmStyle, &sc.Person, &sc.Actor, &sc.VerseStr, &sc.VerseEnd,
			&sc.ScriptText, &sc.ScriptBeginTS, &sc.ScriptEndTS, &sc.MFCCRows, &sc.MFCCCols, &mfccJson)
		if err != nil {
			status = log.Error(d.Ctx, 500, err, "Error in SelectScripts.")
			//return results, status
			panic(status.Message)
		}
		err = json.Unmarshal([]byte(mfccJson), &sc.MFCC)
		if err != nil {
			panic(err)
		}
		results = append(results, sc)
	}
	err = rows.Err()
	if err != nil {
		log.Warn(d.Ctx, err, query)
	}
	return results
}

type Meta struct {
	Index int
	Name  string
	Tag   string
	Dtype string
	Cols  int // I don't
}

func ReflectStruct(aStruct any) []Meta {
	var results []Meta
	sc := reflect.TypeOf(aStruct)
	for i := 0; i < sc.NumField(); i++ {
		var meta Meta
		meta.Index = i
		field := sc.Field(i)
		meta.Name = field.Name
		parts := strings.Split(field.Tag.Get("name"), ",")
		if len(parts) > 1 {
			meta.Tag = parts[0]
			meta.Dtype = parts[1]
			meta.Cols = 1
			//meta.Rows = false
			//if field.Type.Kind() == reflect.Slice {
			//	var scr, ok = aStruct.(Script)
			//	if ok {
			//		meta.Cols = scr.MFCCCols
			//		//		meta.Rows = true
			//	}
			//}
			results = append(results, meta)
		}
	}
	return results
}

func FindActiveScriptCols(structs []Script, meta []Meta) []Meta {
	var results []Meta
	for _, mt := range meta {
		for _, scr := range structs {
			data := reflect.ValueOf(scr).Field(mt.Index)
			if !data.IsZero() {
				results = append(results, mt)
				break
			}
		}
	}
	return results
}

/*
	func FindActiveWordCols(structs []Word, meta []MetaStruct) []MetaStruct {
		var results []MetaStruct
		for _, mt := range meta {
			for _, scr := range structs {
				data := reflect.ValueOf(scr).Field(mt.Index)
				if !data.IsZero() {
					results = append(results, mt)
					break
				}
			}
		}
		return results
	}
*/
func FindNumMFCC(scripts []Script) int {
	var result int
	for _, scr := range scripts {
		cols := scr.MFCCCols
		if cols > 0 {
			result = cols
			break
		}
	}
	return result
}

func NormalizeMFCC(scripts []Script, numMFCC int) []Script {
	for col := 0; col < numMFCC; col++ {
		var sum float64
		var count float64
		for _, scr := range scripts {
			for _, mf := range scr.MFCC {
				value := mf[:][col]
				sum += float64(value)
				count++
			}
		}
		var mean = sum / count
		var devSqr float64
		for _, scr := range scripts {
			for _, mf := range scr.MFCC {
				value := mf[:][col]
				devSqr += math.Pow(float64(value)-mean, 2)
			}
		}
		var stddev = math.Sqrt(devSqr / count)
		for i, scr := range scripts {
			for j, mf := range scr.MFCC {
				value := float64(mf[:][col])
				scripts[i].MFCC[j][:][col] = float32((value - mean) / stddev)
			}
		}
	}
	return scripts
}

func PadRows(scripts []Script, numMFCC int) []Script {
	largest := 0
	for _, scr := range scripts {
		if scr.MFCCRows > largest {
			largest = scr.MFCCRows
		}
	}
	var padRow = make([]float32, numMFCC)
	for _, scr := range scripts {
		needRows := largest - scr.MFCCRows
		for i := 0; i < needRows; i++ {
			scr.MFCC = append(scr.MFCC, padRow)
		}
	}
	return scripts
}
