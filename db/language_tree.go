package db

import (
	"context"
	"dataset"
	log "dataset/logger"
	"encoding/json"
	"fmt"
	"os"
)

type LanguageTree struct {
	ctx   context.Context
	table []Language
	roots []*Language
}

func NewLanguageTree(ctx context.Context) LanguageTree {
	var l LanguageTree
	l.ctx = ctx
	return l
}

func (l *LanguageTree) LoadTable() dataset.Status {
	var status dataset.Status
	// Read json file of languages
	filename := "../db/language/language_tree.jason"
	content, err := os.ReadFile(filename)
	if err != nil {
		return log.Error(l.ctx, 500, err, "Error when opening file: ", filename)
	}
	// Parse json into Language slice
	err = json.Unmarshal(content, &l.table)
	if err != nil {
		return log.Error(l.ctx, 500, err, "Error during Unmarshal(): ", filename)
	}
	return status
}

func (l *LanguageTree) BuildTree() dataset.Status {
	var status dataset.Status
	// Make Map of GlottoId
	var idMap = make(map[string]*Language)
	for i := range l.table {
		lang := l.table[i]
		idMap[lang.GlottoId] = &lang
	}
	// Build Tree
	var idMapCount int
	var parentIdCount int
	for glottoId, lang := range idMap {
		idMapCount++
		if lang.ParentId != "" {
			parentIdCount++
			parent, ok := idMap[lang.ParentId]
			if !ok {
				return log.ErrorNoErr(l.ctx, 500, "Missing parent id: ", lang.ParentId)
			}
			lang.Parent = parent
			idMap[glottoId] = lang
			lang.Parent.Children = append(lang.Parent.Children, lang)
			idMap[lang.ParentId] = lang.Parent
		}
	}
	fmt.Println("count: ", idMapCount)
	fmt.Println("parent count: ", parentIdCount)
	// Build root
	for _, lang := range idMap {
		if lang.Parent == nil {
			l.roots = append(l.roots, lang)
		}
	}
	return status
}

func (l *LanguageTree) Search(iso6393 string, search string) { // whisper, mms_asr, espeak
	// read in the json
	// load it into a tree
	// must have an iso6393 map
	// lookup in iso6393map
	// if not found 400 error
	// perform recursive or stack-based descent of children looking for a true
	// when found count the number of steps
	// continue descending children, but don't go further down.
	// Remeber the iso6393 code (or GlottoId) and minimum steps
	// Go up the tree and repeeat the process
	// Do not go up more than the minimum count.

}
