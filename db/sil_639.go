package db

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

func FindWhisperCompatibility(ctx context.Context, iso3 string) (Sil639, dataset.Status) {
	var iso639 Sil639
	var status dataset.Status
	var database = "iso_639_3.db"
	databasePath := filepath.Join(os.Getenv(`FCBH_DATASET_DB`), database)
	conn, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		status = log.Error(ctx, 500, err, `Failed to open database`, database)
		return iso639, status
	}
	defer conn.Close()
	iso639, status = SelectOnPart1(ctx, conn, iso3)
	if iso639.Lang2 != `` || status.IsErr {
		return iso639, status
	}
	iso639, status = SelectOnMacro(ctx, conn, iso3)
	return iso639, status
}

func SelectOnPart1(ctx context.Context, conn *sql.DB, iso639 string) (Sil639, dataset.Status) {
	var query = `SELECT w.Part1, l.Id, l.Ref_Name
		FROM whisper w
		JOIN languages l ON w.Part1 = l.Part1
		WHERE l.Id = ? AND l.Language_Type = 'L'`
	return genericLangQuery(ctx, conn, iso639, query)
}

func SelectOnMacro(ctx context.Context, conn *sql.DB, iso639 string) (Sil639, dataset.Status) {
	var query = `SELECT w.Part1, l.Id, l.Ref_Name
		FROM whisper w
		JOIN languages l ON w.Part1 = l.Part1
		JOIN macro m ON l.Id = m.M_Id
		WHERE m.I_Id = ? AND l.Language_Type = 'L'`
	return genericLangQuery(ctx, conn, iso639, query)
}

func genericLangQuery(ctx context.Context, conn *sql.DB, iso3 string, query string) (Sil639, dataset.Status) {
	var iso639 Sil639
	var status dataset.Status
	rows, err := conn.Query(query, iso3)
	if err != nil {
		status = log.Error(ctx, 500, err, query)
		return iso639, status
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&iso639.Lang2, &iso639.Lang3, &iso639.Name)
		if err != nil {
			status = log.Error(ctx, 500, err, query)
			return iso639, status
		}
	}
	err = rows.Err()
	if err != nil {
		log.Warn(ctx, err, query)
	}
	return iso639, status
}
