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

type ToolName string

const (
	ESpeak  ToolName = "espeak"
	MMSASR  ToolName = "mms_asr"
	MMSLID  ToolName = "mms_lid"
	MMSTTS  ToolName = "mms_tts"
	Whisper ToolName = "whisper"
)

func main() {
	var numLanguages = 26879
	languages := loadGlottoLanguoid()
	languages = loadIso6393(languages)
	languages = loadAIToolCompatibility(languages, "db/language/espeak.tab", ESpeak, 1)
	languages = loadAIToolCompatibility(languages, "db/language/mms_asr.tab", MMSASR, 0)
	languages = loadAIToolCompatibility(languages, "db/language/mms_lid.tab", MMSLID, 0)
	languages = loadAIToolCompatibility(languages, "db/language/mms_tts.tab", MMSTTS, 0)
	languages = loadAIToolCompatibility(languages, "db/language/whisper.tab", Whisper, 0)
	if len(languages) != numLanguages {
		fmt.Println("Load Iso6393: Expected ", numLanguages, " got ", len(languages))
		os.Exit(1)
	}
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
	var isoMap = make(map[string]string)
	var inGlotto = make(map[string]bool)
	var notInIso = make(map[string]bool)
	var notInGlotto = make(map[string]bool)
	fmt.Println("Num Glotto Records", len(languages))
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
		isoMap[iso6393] = iso6391
	}
	fmt.Println("Num iso639-3 records", len(isoMap))
	for i := range languages {
		iso6393 := languages[i].Iso6393
		if iso6393 != "" {
			inGlotto[iso6393] = languages[i].Bookkeeping
			iso6391, ok := isoMap[iso6393]
			if ok {
				if iso6391 != "" {
					languages[i].Iso6391 = iso6391
				}
			} else {
				notInIso[iso6393] = true
				fmt.Println("Glotto ISO639-3 Not In ISO List", iso6393)
			}
		}
	}
	fmt.Println("Num iso639-3 values in glotto", len(inGlotto))
	fmt.Println("Num Glotto Not In ISO", len(notInIso))
	for iso6393 := range isoMap {
		bookkeeping, ok := inGlotto[iso6393]
		if !ok && !bookkeeping {
			notInGlotto[iso6393] = true
			fmt.Println("ISO639-3 Not In Glotto List", iso6393)
		}
	}
	fmt.Println("Num iso639-3 records not in glotto", len(notInGlotto))
	return languages
}

func loadAIToolCompatibility(languages []db.Language, filePath string, toolName ToolName, iso3Col int) []db.Language {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	var toolMap = make(map[string]bool)
	var record []string
	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		toolMap[record[iso3Col]] = true
	}
	var usedMap = make(map[string]bool)
	for i := range languages {
		_, ok := toolMap[languages[i].Iso6393]
		if ok {
			languages[i] = setLanguage(languages[i], toolName)
			usedMap[languages[i].Iso6393] = true
		} else {
			_, ok = toolMap[languages[i].Iso6391]
			if ok {
				//languages[i].Whisper = true
				languages[i] = setLanguage(languages[i], toolName)
				usedMap[languages[i].Iso6391] = true
			}
		}
	}
	for iso := range toolMap {
		_, ok := usedMap[iso]
		if !ok {
			fmt.Println("AI Tool", toolName, "Has iso code", iso, " but it has no match in table")
		}
	}
	return languages
}

func setLanguage(language db.Language, toolName ToolName) db.Language {
	switch toolName {
	case ESpeak:
		language.ESpeak = true
	case MMSASR:
		language.MMSASR = true
	case MMSLID:
		language.MMSLID = true
	case MMSTTS:
		language.MMSTTS = true
	case Whisper:
		language.Whisper = true
	default:
		panic("Unknown tool name: " + toolName)
	}
	return language
}

func outputJSON(languages []db.Language) {
	bytes, err := json.MarshalIndent(languages, "", "    ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("db/language/language_tree.jason_new", bytes, 0644)
	if err != nil {
		panic(err)
	}
}
