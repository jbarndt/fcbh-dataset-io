package match

import (
	"context"
	"dataset"
	"dataset/db"
)

type FAverse struct {
	scriptId    int64
	bookId      string
	chapter     int
	verseStr    string
	verseSeq    int
	beginTS     float64
	endTS       float64
	faScore     float64
	words       []db.Audio
	critWords   []int
	critScore   float64
	questWords  []int
	questScore  float64
	startTSDiff float64
	endTSDiff   float64
	critDiff    bool
	questDiff   bool
}

type AlignErrorCalc struct {
	ctx  context.Context
	conn db.DBAdapter
}

func NewAlignErrorCalc(ctx context.Context, conn db.DBAdapter) AlignErrorCalc {
	var a AlignErrorCalc
	a.ctx = ctx
	a.conn = conn
	return a
}

func (a *AlignErrorCalc) Process() ([]FAverse, dataset.Status) {
	var verses []FAverse
	var status dataset.Status
	var words []db.Audio
	words, status = a.conn.SelectFAWordTimestamps()
	if status.IsErr {
		return verses, status
	}
	verses = a.createVerseSlice(words)
	verseMap := a.createVerseMap(verses)
	verses = a.addWordsToVerses(verseMap, words)
	a.computeCritical(verses, 0.001)
	a.computeQuestion(verses, 0.001)
	a.computeTSDiff(verses)
	a.markCriticalDiff(verses, 1.5)
	a.markQuestionDiff(verses, 0.9)
	return verses, status
}

func (a *AlignErrorCalc) createVerseSlice(words []db.Audio) []FAverse {
	var result []FAverse
	var priorScriptId int64 = -1
	for _, w := range words {
		if w.ScriptId != priorScriptId {
			priorScriptId = w.ScriptId
			var v FAverse
			v.scriptId = w.ScriptId
			v.bookId = w.BookId
			v.chapter = w.ChapterNum
			v.verseStr = w.VerseStr
			v.verseSeq = w.VerseSeq
			v.beginTS = w.ScriptBeginTS
			v.endTS = w.ScriptEndTS
			v.faScore = w.ScriptFAScore
			result = append(result, v)
		}
	}
	return result
}

func (a *AlignErrorCalc) createVerseMap(verses []FAverse) map[int64]FAverse {
	var result = make(map[int64]FAverse)
	for _, vs := range verses {
		result[vs.scriptId] = vs
	}
	return result
}

func (a *AlignErrorCalc) addWordsToVerses(verseMap map[int64]FAverse, words []db.Audio) []FAverse {
	for _, w := range words {
		vs, ok := verseMap[w.ScriptId]
		if !ok {
			panic("what??")
		}
		vs.words = append(vs.words, w)
		verseMap[w.ScriptId] = vs
	}
	var result []FAverse
	for _, vs := range verseMap {
		result = append(result, vs)
	}
	return result
}

func (a *AlignErrorCalc) computeCritical(verses []FAverse, threshold float64) {
	for i := range verses {
		verses[i].critWords, verses[i].critScore = a.findErrors(verses[i].words, threshold)
	}
}

func (a *AlignErrorCalc) computeQuestion(verses []FAverse, threshold float64) {
	for i := range verses {
		verses[i].questWords, verses[i].questScore = a.findErrors(verses[i].words, threshold)
	}
}

func (a *AlignErrorCalc) findErrors(words []db.Audio, threshold float64) ([]int, float64) {
	var results []int
	var score float64
	for i, w := range words {
		if w.FAScore <= threshold {
			results = append(results, i)
			score += 1.0 - w.FAScore
		}
	}
	return results, score
}

func (a *AlignErrorCalc) computeTSDiff(verses []FAverse) {
	for i := range verses {
		firstWord := verses[i].words[0]
		verses[i].startTSDiff = firstWord.BeginTS - verses[i].beginTS
		lastWord := verses[i].words[len(verses[i].words)-1]
		verses[i].endTSDiff = verses[i].endTS - lastWord.EndTS
	}
}

func (a *AlignErrorCalc) markCriticalDiff(verses []FAverse, threshold float64) {
	for i := range verses {
		if verses[i].startTSDiff >= threshold || verses[i].endTSDiff >= threshold {
			verses[i].critDiff = true
		}
	}
}

func (a *AlignErrorCalc) markQuestionDiff(verses []FAverse, threshold float64) {
	for i := range verses {
		if verses[i].startTSDiff >= threshold || verses[i].endTSDiff >= threshold {
			verses[i].questDiff = true
		}
	}
}
