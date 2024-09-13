package db

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"testing"
)

var result []*Language
var count int

func TestLanguageTree_BuildTree(t *testing.T) {
	var tree = NewLanguageTree(context.Background())
	tree.Load()
	fmt.Println("count:", len(tree.table))
	recursiveDescent(tree.roots)
	sort.Slice(result, func(i, j int) bool {
		return result[i].GlottoId < result[j].GlottoId
	})
	if len(tree.table) != len(result) {
		t.Errorf("len(tree.table) = %d; actual %d", len(tree.table), len(result))
	}
	fmt.Println("count: ", count)
	outputResult(result)
}

func TestLanguageTree_Search(t *testing.T) {
	var tree = NewLanguageTree(context.Background())
	status := tree.Load()
	if status.IsErr {
		t.Error("status.IsErr:", status)
	}
	doSearch(t, tree, "eng", "whisper", 0, []string{"stan1293"})
	doSearch(t, tree, "spa", "whisper", 0, []string{"stan1288"})
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
	tree.table = langs
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
		lang.Whisper = false
		tree.idMap[id] = lang
	}
	for _, id := range glottoIds {
		lang := tree.idMap[id]
		lang.Whisper = true
		tree.isoMap[id] = lang
	}
}

func doSearch(t *testing.T, tree LanguageTree, iso639 string, search string, distance int, result []string) {
	langs, dist, status := tree.Search(iso639, search)
	if status.IsErr {
		t.Error("status.IsErr:", status)
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

func recursiveDescent(langs []*Language) {
	for _, lang := range langs {
		result = append(result, lang)
		count++
		recursiveDescent(lang.Children)
	}
}

func outputResult(results []*Language) {
	bytes, err := json.MarshalIndent(results, "", "    ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("../db/language/language_tree.jason2", bytes, 0644)
	if err != nil {
		panic(err)
	}
}
