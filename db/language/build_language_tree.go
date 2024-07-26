package main

import (
	"dataset/db"
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"strconv"
)

func main() {
	languages := loadGlottoLanguoid()
	outputJSON(languages)
}

func loadGlottoLanguoid() []db.Language {
	var languages []db.Language
	file, err := os.Open("db/language/languoid.tab")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	first := true
	var record []string
	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if first {
			first = false
			continue
		}
		var lang db.Language
		lang.GlottoId = record[0]
		lang.FamilyId = record[1]
		lang.ParentId = record[2]
		lang.Name = record[3]
		lang.Bookkeeping, err = strconv.ParseBool(record[4])
		if err != nil {
			panic(err)
		}
		lang.Level = record[5]
		lang.Iso6393 = record[8]
		lang.CountryIds = record[14]
		languages = append(languages, lang)
	}
	return languages
}

func outputJSON(languages []db.Language) {
	bytes, err := json.MarshalIndent(languages, "", "    ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("db/language/language_tree.jason", bytes, 0644)
	if err != nil {
		panic(err)
	}
}
