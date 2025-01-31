package search

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"testing"
)

func TestLanguageTree_Search(t *testing.T) {
	var tree = NewLanguageTree(context.Background())
	err := tree.Load()
	if err != nil {
		t.Error("status.IsErr:", err)
	}
	doSearch(t, tree, "eng", "whisper", 0, []string{"stan1293"})
	doSearch(t, tree, "spa", "whisper", 0, []string{"stan1288"})
}

func TestLanguageTree_AllIsoCodes(t *testing.T) {
	var isoCount int
	var notFoundCount int
	var toolCount = make(map[string]int)
	var tree = NewLanguageTree(context.Background())
	err := tree.Load()
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open("../build/data/iso-639-3.tab")
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
		var iso639 []string
		iso639 = append(iso639, record[0])
		if record[3] != "" {
			iso639 = append(iso639, record[3])
		}
		//iso639 = nil
		//iso639 = append(iso639, "zu")
		for _, isoCode := range iso639 {
			isoCount++
			_, _, err2 := tree.Search(isoCode, MMSLID)
			if err2 != nil {
				if strings.HasPrefix(err2.Error(), "") {
					fmt.Println(isoCode, "is not found")
					notFoundCount++
				} else {
					t.Fatal(err2)
				}
			} else {
				for _, tool := range []string{ESpeak, MMSASR, MMSLID, MMSTTS, Whisper} {
					langs, distance, err3 := tree.Search(isoCode, tool)
					if err3 != nil {
						t.Fatal(err3)
					} else {
						fmt.Println(isoCode, tool, "->", langs, distance)
						if len(langs) > 0 {
							toolCount[tool]++
						}
					}
				}
			}
		}
	}
	fmt.Println("isoCodeCount:", isoCount, "not FoundCount:", notFoundCount)
	var knowCount = float64(isoCount - notFoundCount)
	for tool, count := range toolCount {
		pct := math.Round(float64(count) / knowCount * 100.0)
		fmt.Println("tool:", tool, "count:", count, "percent:", pct, "%")
	}
	if isoCount != 8104 {
		t.Error("Expected iso code count", 8104, "actual", isoCount)
	}
	if notFoundCount != 143 {
		t.Error("Expected not found count", 143, "actual", notFoundCount)
	}
	if toolCount[ESpeak] != 4811 {
		t.Error("Expected tool code", ESpeak, "actual", toolCount[ESpeak])
	}
	if toolCount[MMSASR] != 6753 {
		t.Error("Expected tool code", MMSASR, "actual", toolCount[ESpeak])
	}
	if toolCount[MMSLID] != 7561 {
		t.Error("Expected tool code", MMSLID, "actual", toolCount[ESpeak])
	}
}

func TestLanguageTree_SampleData(t *testing.T) {
	// Create a sample language tree
	var langs = make([]Language, 0)
	langs = append(langs, Language{Name: "Indo-European", GlottoId: "indo1319", ParentId: "", Iso6393: "euro"})

	langs = append(langs, Language{Name: "Germanic", GlottoId: "germ1287", ParentId: "indo1319", Iso6393: "germ"})
	langs = append(langs, Language{Name: "Romance", GlottoId: "roma1334", ParentId: "indo1319", Iso6393: "roma"})
	langs = append(langs, Language{Name: "Slavic", GlottoId: "slav1255", ParentId: "indo1319", Iso6393: "slav"})

	langs = append(langs, Language{Name: "English", GlottoId: "stan1293", ParentId: "germ1287", Iso6393: "eng"})
	langs = append(langs, Language{Name: "German", GlottoId: "stan1295", ParentId: "germ1287", Iso6393: "deu"})
	langs = append(langs, Language{Name: "French", GlottoId: "stan1290", ParentId: "roma1334", Iso6393: "fra"})
	langs = append(langs, Language{Name: "Spanish", GlottoId: "stan1288", ParentId: "roma1334", Iso6393: "spa"})
	langs = append(langs, Language{Name: "Russian", GlottoId: "russ1263", ParentId: "slav1255", Iso6393: "rus"})

	langs = append(langs, Language{Name: "British", GlottoId: "british", ParentId: "stan1293", Iso6393: "brit"})
	langs = append(langs, Language{Name: "American", GlottoId: "american", ParentId: "stan1293", Iso6393: "amer"})
	langs = append(langs, Language{Name: "Australian", GlottoId: "australian", ParentId: "stan1293", Iso6393: "aust"})

	tree := NewLanguageTree(context.Background())
	tree.Table = langs
	searchType := `whisper`
	tree.validateSearch(searchType)
	tree.buildTree()
	setWhisper(tree, []string{"stan1293"})
	doSearch(t, tree, "eng", "whisper", 0, []string{"stan1293"})
	setWhisper(tree, []string{"american"})
	doSearch(t, tree, "eng", "whisper", 1, []string{"american"})
	setWhisper(tree, []string{"germ1287"})
	doSearch(t, tree, "eng", "whisper", 1, []string{"germ1287"})
	setWhisper(tree, []string{"indo1319"})
	doSearch(t, tree, "eng", "whisper", 2, []string{"indo1319"})
	setWhisper(tree, []string{"british", "american", "australian"})
	doSearch(t, tree, "eng", "whisper", 1, []string{"british", "american", "australian"})
	setWhisper(tree, []string{"germ1287", "roma1334"})
	doSearch(t, tree, "eng", "whisper", 1, []string{"germ1287"})
}

func setWhisper(tree LanguageTree, glottoIds []string) {
	for id, lang := range tree.idMap {
		lang.Whisper = ""
		tree.idMap[id] = lang
	}
	for _, id := range glottoIds {
		lang := tree.idMap[id]
		lang.Whisper = lang.Iso6393
		tree.isoMap[id] = lang
	}
}

func doSearch(t *testing.T, tree LanguageTree, iso639 string, search string, distance int, result []string) {
	langs, dist, err := tree.DetailSearch(iso639, search)
	if err != nil {
		t.Error("status.IsErr:", err)
	}
	if dist != distance {
		t.Error("Expected Depth:", distance, "Found Distance:", dist)
	}
	if len(langs) != len(result) {
		t.Error("Expected Num:", len(result), "Found Num:", len(langs))
	} else {
		var resultMap = make(map[string]bool)
		for _, lang := range result {
			resultMap[lang] = true
		}
		for _, lang := range langs {
			_, ok := resultMap[lang.GlottoId]
			if !ok {
				t.Error("Expected lang", lang, "Found lang", lang.GlottoId)
			}
		}
	}
}
