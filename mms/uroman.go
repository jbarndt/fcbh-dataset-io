package mms

import (
	"context"
	"dataset/db"
	log "dataset/logger"
	"dataset/utility"
	"os"
	"path/filepath"
)

// uroman.go requires a pip install uroman, but it uses the uroman.pl that is included
// https://github.com/isi-nlp/uroman/tree/master

func EnsureUroman(conn db.DBAdapter, lang string) *log.Status {
	hasUroman, status := CheckUroman(conn)
	if status != nil {
		return status
	}
	if !hasUroman {
		status = UpdateUroman(conn, lang)
	}
	return status
}

func CheckUroman(conn db.DBAdapter) (bool, *log.Status) {
	var result bool
	textLen, status := conn.SelectScalarInt("SELECT sum(length(script_text)) FROM scripts")
	if status != nil {
		return result, status
	}
	uromanLen, status := conn.SelectScalarInt("SELECT sum(length(uroman)) FROM scripts")
	if status != nil {
		return result, status
	}
	return float64(uromanLen)*1.2 >= float64(textLen), status
}

func UpdateUroman(conn db.DBAdapter, lang string) *log.Status {
	scripts, status := conn.SelectScripts()
	if status != nil {
		return status
	}
	scripts, status = SetUroman(conn.Ctx, scripts, lang)
	if status != nil {
		return status
	}
	_, status = conn.UpdateUromanText(scripts)
	return status
}

func SetUroman(ctx context.Context, lines []db.Script, lang string) ([]db.Script, *log.Status) {
	uromanPath := filepath.Join(os.Getenv("GOPROJ"), "dataset", "mms", "uroman_stdio.py")
	uroman, status := utility.NewStdioExec(ctx, os.Getenv(`FCBH_MMS_FA_PYTHON`), uromanPath, "-l", lang)
	if status != nil {
		return lines, status
	}
	defer uroman.Close()
	for i := range lines {
		lines[i].URoman, status = uroman.Process(lines[i].ScriptText)
		if status != nil {
			return lines, status
		}
	}
	return lines, status
}
