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
	ctx    context.Context
	table  []Language
	roots  []*Language
	idMap  map[string]*Language
	isoMap map[string]*Language
}

func NewLanguageTree(ctx context.Context) LanguageTree {
	var l LanguageTree
	l.ctx = ctx
	l.idMap = make(map[string]*Language)
	l.isoMap = make(map[string]*Language)
	return l
}

func (l *LanguageTree) Load() dataset.Status {
	status := l.loadTable()
	if status.IsErr {
		return status
	}
	status = l.buildTree()
	return status
}

func (l *LanguageTree) Search(iso6393 string, search string) ([]*Language, int, dataset.Status) {
	var results []*Language
	var distance int
	var status dataset.Status
	var lang *Language
	lang, ok := l.isoMap[iso6393]
	if !ok {
		status = log.ErrorNoErr(l.ctx, 400, "iso code ", iso6393, " is not known.")
		return results, distance, status
	}
	if !l.validateSearch(search) {
		status = log.ErrorNoErr(l.ctx, 400, "Search parameter", search, "is not known")
		return results, distance, status
	}
	results, distance = l.searchRelatives(lang, search)
	return results, distance, status
}

func (l *LanguageTree) loadTable() dataset.Status {
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

func (l *LanguageTree) buildTree() dataset.Status {
	var status dataset.Status
	// Make Map of GlottoId
	for i := range l.table {
		lang := l.table[i]
		l.idMap[lang.GlottoId] = &lang
	}
	// Build Tree
	var idMapCount int
	var parentIdCount int
	for glottoId, lang := range l.idMap {
		idMapCount++
		if lang.ParentId != "" {
			parentIdCount++
			parent, ok := l.idMap[lang.ParentId]
			if !ok {
				return log.ErrorNoErr(l.ctx, 500, "Missing parent id: ", lang.ParentId)
			}
			lang.Parent = parent
			l.idMap[glottoId] = lang
			lang.Parent.Children = append(lang.Parent.Children, lang)
			l.idMap[lang.ParentId] = lang.Parent
		}
		if lang.Iso6393 != "" {
			l.isoMap[lang.Iso6393] = lang
		}
	}
	fmt.Println("count: ", idMapCount)
	fmt.Println("parent count: ", parentIdCount)
	// Build root
	for _, lang := range l.idMap {
		if lang.Parent == nil {
			l.roots = append(l.roots, lang)
		}
	}
	return status
}

func (l *LanguageTree) searchRelatives(start *Language, search string) ([]*Language, int) {
	var finalLang []*Language
	var finalDepth int
	var limit = 1000
	//for finalDepth > 0 && limit > 0 {
	for limit > 0 && start != nil {
		fmt.Println("\nSearching", start.Name, search, limit)
		results, depth := l.descendantSearch(start, search, limit)
		fmt.Println("descendantSearch Depth", depth, "num", len(results))
		for _, result := range results {
			fmt.Println("descendentSearch lang.Name", result.Name)
		}
		if len(results) > 0 {
			finalLang = results
			finalDepth = depth
			limit = depth
		}
		start = start.Parent
	}
	return finalLang, finalDepth
}

// DescendantSearch performs a breadth-first search of the LanguageTree
func (l *LanguageTree) descendantSearch(start *Language, search string, limit int) ([]*Language, int) {
	var results []*Language
	var depth int
	if start == nil {
		return results, depth
	}
	var queue LanguageQueue
	queue.Enqueue(start, 0)
	for !queue.IsEmpty() {
		item := queue.Dequeue()
		if (len(results) > 0 && item.Depth > depth) || item.Depth > limit {
			return results, depth
		}
		depth = item.Depth
		if l.isMatch(item.Lang, search) {
			results = append(results, item.Lang)
		}
		fmt.Printf("Depth: %d, Name: %s, GlottoId: %s\n", item.Depth, item.Lang.Name, item.Lang.GlottoId)

		for _, child := range item.Lang.Children {
			queue.Enqueue(child, item.Depth+1)
		}
	}
	return results, depth
}

func (l *LanguageTree) validateSearch(search string) bool {
	switch search {
	case `whisper`:
		return true
	case `mms_asr`:
		return true
	case `espeak`:
		return true
	default:
		return false
	}
}

func (l *LanguageTree) isMatch(lang *Language, search string) bool {
	switch search {
	case `whisper`:
		return lang.Whisper
	case `mms_asr`:
		return lang.MMSASR
	case `espeak`:
		return lang.ESpeak
	default:
		return false
	}
}
