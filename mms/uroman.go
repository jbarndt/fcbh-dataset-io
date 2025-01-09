package mms

import (
	"context"
	"dataset"
	"dataset/db"
	"dataset/utility"
	"os"
	"path/filepath"
)

// uroman.go requires a pip install uroman, but it uses the uroman.pl that is included
// https://github.com/isi-nlp/uroman/tree/master

func EnsureUroman(conn db.DBAdapter, lang string) dataset.Status {
	hasUroman, status := CheckUroman(conn)
	if status.IsErr {
		return status
	}
	if !hasUroman {
		status = UpdateUroman(conn, lang)
	}
	return status
}

func CheckUroman(conn db.DBAdapter) (bool, dataset.Status) {
	var result bool
	textLen, status := conn.SelectScalarInt("SELECT sum(length(script_text)) FROM scripts")
	if status.IsErr {
		return result, status
	}
	uromanLen, status := conn.SelectScalarInt("SELECT sum(length(uroman)) FROM scripts")
	if status.IsErr {
		return result, status
	}
	return float64(uromanLen)*1.2 >= float64(textLen), status
}

func UpdateUroman(conn db.DBAdapter, lang string) dataset.Status {
	scripts, status := conn.SelectScripts()
	if status.IsErr {
		return status
	}
	scripts, status = SetUroman(conn.Ctx, scripts, lang)
	if status.IsErr {
		return status
	}
	_, status = conn.UpdateUromanText(scripts)
	return status
}

func SetUroman(ctx context.Context, lines []db.Script, lang string) ([]db.Script, dataset.Status) {
	uromanPath := filepath.Join(os.Getenv("GOPROJ"), "dataset", "mms", "uroman_stdio.py")
	uroman, status := utility.NewStdioExec(ctx, os.Getenv(`FCBH_MMS_FA_PYTHON`), uromanPath, "-l", lang)
	if status.IsErr {
		return lines, status
	}
	defer uroman.Close()
	for i := range lines {
		lines[i].URoman, status = uroman.Process(lines[i].ScriptText)
		if status.IsErr {
			return lines, status
		}
	}
	return lines, status
}
