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

func doSearch(t *testing.T, tree LanguageTree, iso639 string, search string, depth int, result []string) {
	langs, deph, status := tree.Search(iso639, search)
	if status.IsErr {
		t.Error("status.IsErr:", status)
	}
	if deph != depth {
		t.Error("Expected Depth:", depth, "Found Depth:", deph)
	}
	if len(langs) != len(result) {
		t.Error("Expected Num:", len(result), "Found Num:", len(langs))
	}
	for i, lang := range result {
		if lang != langs[i].GlottoId {
			t.Error("Expected lang", lang, "Found lang", langs[i].GlottoId)
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
