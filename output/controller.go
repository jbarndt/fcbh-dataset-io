package output

import (
	"dataset/db"
	"fmt"
)

func PrepareScripts(conn db.DBAdapter) ([]any, []Meta) {
	var script Script
	meta := ReflectStruct(script)
	fmt.Println("meta :=", meta)
	scripts := LoadScriptStruct(conn)
	numMFCC := FindNumScriptMFCC(scripts)
	fmt.Println("numMFCC :=", numMFCC)
	SetNumMFCC(&meta, numMFCC)
	fmt.Println("meta-num :=", meta)
	scripts = NormalizeScriptMFCC(scripts, numMFCC)
	scripts = PadScriptRows(scripts, numMFCC)
	scriptAny := ConvertScriptsAny(scripts)
	fmt.Println("scriptAny :=", len(scriptAny))
	meta = FindActiveCols(scriptAny, meta)
	fmt.Println("meta-prune :=", meta)
	return scriptAny, meta
}

func PrepareWords(conn db.DBAdapter) ([]any, []Meta) {
	var word Word
	meta := ReflectStruct(word)
	fmt.Println("meta :=", meta)
	words := LoadWordStruct(conn)
	numMFCC := FindNumWordMFCC(words)
	fmt.Println("numMFCC :=", numMFCC)
	SetNumMFCC(&meta, numMFCC)
	fmt.Println("meta-num :=", meta)
	words = NormalizeWordMFCC(words, numMFCC)
	words = PadWordRows(words, numMFCC)
	wordAny := ConvertWordsAny(words)
	//fmt.Println(words)
	fmt.Println("wordAny :=", len(wordAny))
	meta = FindActiveCols(wordAny, meta)
	fmt.Println("meta-prune :=", meta)
	return wordAny, meta
}
