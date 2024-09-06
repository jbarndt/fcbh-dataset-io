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
	tree.LoadTable()
	fmt.Println("count:", len(tree.table))
	tree.BuildTree()
	recursiveDescent(tree.roots)
	sort.Slice(result, func(i, j int) bool {
		return result[i].GlottoId < result[j].GlottoId
	})
	if len(tree.table) != len(result) {
		t.Errorf("len(tree.table) = %d; actual %d", len(tree.table), len(result))
	}
	fmt.Println("count: ", count)
	outputResult()
	// compare both with diff-patch-match
}

func recursiveDescent(langs []*Language) {
	for _, lang := range langs {
		result = append(result, lang)
		count++
		recursiveDescent(lang.Children)
	}
}

func outputResult() {
	bytes, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("../db/language/language_tree.jason2", bytes, 0644)
	if err != nil {
		panic(err)
	}
}
