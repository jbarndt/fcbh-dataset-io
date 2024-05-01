package output

import (
	"context"
	"dataset/db"
	"fmt"
)

type Output struct {
	ctx context.Context
}

func NewOutput(ctx context.Context) Output {
	var o Output
	o.ctx = ctx
	return o
}

func (o *Output) PrepareScripts(conn db.DBAdapter, normalize bool, pad bool) ([]any, []Meta) {
	var script Script
	meta := o.ReflectStruct(script)
	fmt.Println("meta :=", meta)
	scripts, status := o.LoadScriptStruct(conn)
	if status.IsErr {
		panic(status.Message)
	}
	numMFCC := o.FindNumScriptMFCC(scripts)
	fmt.Println("numMFCC :=", numMFCC)
	o.SetNumMFCC(&meta, numMFCC)
	fmt.Println("meta-num :=", meta)
	if normalize {
		scripts = o.NormalizeScriptMFCC(scripts, numMFCC)
	}
	if pad {
		scripts = o.PadScriptRows(scripts, numMFCC)
	}
	scriptAny := o.ConvertScriptsAny(scripts)
	fmt.Println("scriptAny :=", len(scriptAny))
	meta = o.FindActiveCols(scriptAny, meta)
	fmt.Println("meta-prune :=", meta)
	o.SetCSVPos(&meta)
	fmt.Println("meta-pos :=", meta)
	return scriptAny, meta
}

func (o *Output) PrepareWords(conn db.DBAdapter, normalize bool, pad bool) ([]any, []Meta) {
	var word Word
	meta := o.ReflectStruct(word)
	fmt.Println("meta :=", meta)
	words, status := o.LoadWordStruct(conn)
	if status.IsErr {
		panic(status.Message)
	}
	numMFCC := o.FindNumWordMFCC(words)
	fmt.Println("numMFCC :=", numMFCC)
	o.SetNumMFCC(&meta, numMFCC)
	fmt.Println("meta-num :=", meta)
	o.FindNumWordEnc(words, &meta)
	fmt.Println("meta-enc :=", meta)
	if normalize {
		words = o.NormalizeWordMFCC(words, numMFCC)
	}
	if pad {
		words = o.PadWordRows(words, numMFCC)
	}
	wordAny := o.ConvertWordsAny(words)
	//fmt.Println(words)
	fmt.Println("wordAny :=", len(wordAny))
	meta = o.FindActiveCols(wordAny, meta)
	fmt.Println("meta-prune :=", meta)
	o.SetCSVPos(&meta)
	fmt.Println("meta-col :=", meta)
	return wordAny, meta
}
