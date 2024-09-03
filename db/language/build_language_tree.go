package main

import (
	"dataset/db"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	languages := loadGlottoLanguoid()
	languages = loadIso6393(languages)
	languages = loadWhisper(languages)
	root := buildTree(languages)
	//outputJSON(languages)
	outputJSON(root)
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
	var count6393 = 0
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
		if lang.Iso6393 != "" {
			count6393++
		}
		lang.CountryIds = record[14]
		languages = append(languages, lang)
	}
	fmt.Println("Num iso639-3", count6393)
	return languages
}

func loadIso6393(languages []db.Language) []db.Language {
	var part1Map = make(map[string]string)
	file, err := os.Open("db/language/iso-639-3.tab")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = '\t'
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
		iso6393 := record[0]
		iso6391 := record[3]
		if iso6391 != "" {
			part1Map[iso6393] = iso6391
		}
	}
	var missing = make(map[string]int)
	for i := range languages {
		iso6393 := languages[i].Iso6393
		if iso6393 != "" {
			iso6391, ok := part1Map[iso6393]
			if ok {
				languages[i].Iso6391 = iso6391
			}
		}
	}
	for iso, cnt := range missing {
		fmt.Println("missing", iso, cnt)
	}
	fmt.Println(len(missing))
	//fmt.Println("number of mismatch iso639-3", missing)
	return languages
}

func loadWhisper(languages []db.Language) []db.Language {
	var isoMap = make(map[string]bool)
	file, err := os.Open("db/language/whisper.tab")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	var record []string
	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		isoMap[record[0]] = true
	}
	for i := range languages {
		_, ok := isoMap[languages[i].Iso6393]
		if ok {
			languages[i].Whisper = true
		}
		_, ok = isoMap[languages[i].Iso6391]
		if ok {
			languages[i].Whisper = true
		}
	}
	return languages
}

func buildTree(languages []db.Language) []db.Language {
	var root []db.Language
	var idMap = make(map[string]db.Language)
	for _, lang := range languages {
		idMap[lang.GlottoId] = lang
	}
	for _, lang := range languages {
		if lang.ParentId == "" {
			root = append(root, lang)
		} else {
			parentLang, ok := idMap[lang.ParentId]
			if !ok {
				panic(lang.ParentId + " does not exist")
			}
			parentLang.Children = append(parentLang.Children, lang)
			idMap[lang.ParentId] = parentLang
		}
	}
	for _, lang := range root {
		fmt.Println(lang.GlottoId, lang.ParentId, lang.Iso6393, lang.FamilyId, lang.Bookkeeping, lang.Level)
	}
	fmt.Println(len(root))
	return root
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
