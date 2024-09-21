package xxxtobedeleted

import (
	"context"
	"database/sql"
	"dataset"
	log "dataset/logger"
	"os"
	"path/filepath"
)

type Sil639 struct {
	Lang3 string
	Lang2 string
	Name  string
}

func FindWhisperCompatibility(ctx context.Context, iso3 string) ([]Sil639, dataset.Status) {
	var iso639s []Sil639
	var status dataset.Status
	var database = "iso_639_3.db"
	databasePath := filepath.Join(os.Getenv(`FCBH_DATASET_DB`), database)
	conn, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		status = log.Error(ctx, 500, err, `Failed to open database`, database)
		return iso639s, status
	}
	defer conn.Close()
	iso639s, status = SelectOnPart1(ctx, conn, iso3)
	if status.IsErr || len(iso639s) > 0 {
		return iso639s, status
	}
	iso639s, status = SelectOnMacro(ctx, conn, iso3)
	if status.IsErr || len(iso639s) > 0 {
		return iso639s, status
	}
	iso639s, status = SelectOnCountry(ctx, conn, iso3)
	return iso639s, status
}

func SelectOnPart1(ctx context.Context, conn *sql.DB, iso639 string) ([]Sil639, dataset.Status) {
	var query = `SELECT w.Part1, l.Id, l.Ref_Name
		FROM whisper w
		JOIN languages l ON w.Part1 = l.Part1
		WHERE l.Id = ? AND l.Language_Type = 'L'`
	return genericLangQuery(ctx, conn, iso639, query)
}

func SelectOnMacro(ctx context.Context, conn *sql.DB, iso639 string) ([]Sil639, dataset.Status) {
	var query = `SELECT w.Part1, l.Id, l.Ref_Name
		FROM whisper w
		JOIN languages l ON w.Part1 = l.Part1
		JOIN macro m ON l.Id = m.M_Id
		WHERE m.I_Id = ? AND l.Language_Type = 'L'`
	return genericLangQuery(ctx, conn, iso639, query)
}

func SelectOnCountry(ctx context.Context, conn *sql.DB, iso639 string) ([]Sil639, dataset.Status) {
	var query = `SELECT w.Part1, lc1.LangId, lc1.CountryId
		FROM language_country lc1
		JOIN language_country lc2 ON lc1.CountryId = lc2.CountryId
		JOIN languages l ON lc1.LangId = l.Id
		JOIN whisper w ON l.Part1 = w.Part1
		WHERE lc2.LangId = ?
		AND l.Language_Type = 'L'`
	return genericLangQuery(ctx, conn, iso639, query)
}

func genericLangQuery(ctx context.Context, conn *sql.DB, iso3 string, query string) ([]Sil639, dataset.Status) {
	var iso639s []Sil639
	var status dataset.Status
	rows, err := conn.Query(query, iso3)
	if err != nil {
		status = log.Error(ctx, 500, err, query)
		return iso639s, status
	}
	defer rows.Close()
	for rows.Next() {
		var iso639 Sil639
		err = rows.Scan(&iso639.Lang2, &iso639.Lang3, &iso639.Name)
		if err != nil {
			status = log.Error(ctx, 500, err, query)
			return iso639s, status
		}
		iso639s = append(iso639s, iso639)
	}
	err = rows.Err()
	if err != nil {
		log.Warn(ctx, err, query)
	}
	return iso639s, status
}
