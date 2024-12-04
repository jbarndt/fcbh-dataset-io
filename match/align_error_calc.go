package match

import (
	"context"
	"dataset"
	"dataset/db"
	"fmt"
	"gonum.org/v1/gonum/stat"
	"math"
)

const (
	criticalThreshold = 0.001
	questionThreshold = 0.01
)

type ErrorType int

const (
	noError ErrorType = iota
	scoreCritical
	scoreQuestion
	durationLong
	betweenCharsLong
	betweenWordsLong
	betweenVersesLong
	betweenChaptersLong
)

type SilencePosition int

const (
	betweenChars SilencePosition = iota + 1
	betweenWords
	betweenVerses
	betweenChapters
)

type AlignVerse struct {
	chars        []db.AlignChar
	critScore    int64
	questScore   int64
	longDuration int64
	longSilence  int64
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

func (a *AlignErrorCalc) Process() ([]AlignVerse, dataset.Status) {
	var faVerse []AlignVerse
	var faChars []db.AlignChar
	var status dataset.Status
	faChars, status = a.conn.SelectFACharTimestamps()
	if status.IsErr {
		return faVerse, status
	}
	for i := range faChars {
		if faChars[i].FAScore <= criticalThreshold {
			faChars[i].ScoreError = int(scoreCritical)
		} else if faChars[i].FAScore <= questionThreshold {
			faChars[i].ScoreError = int(scoreQuestion)
		}
	}
	for i := range faChars {
		faChars[i].Duration = faChars[i].EndTS - faChars[i].BeginTS
	}
	for i := 0; i < len(faChars)-1; i++ {
		var curr = faChars[i]
		var next = faChars[i+1]
		faChars[i].Silence = next.BeginTS - curr.EndTS
		if curr.WordId == next.WordId {
			faChars[i].SilencePos = int(betweenChars)
		} else if curr.ScriptId == next.ScriptId {
			faChars[i].SilencePos = int(betweenWords)
		} else if curr.BookId == next.BookId && curr.ChapterNum == next.ChapterNum {
			faChars[i].SilencePos = int(betweenVerses)
		} else {
			faChars[i].SilencePos = int(betweenChapters)
			faChars[i].Silence = next.BeginTS
			// To be correct, I should add the silence of fileDuration - curr.EndTS
		}
	}
	a.setCharSeq(faChars)
	mean, stddev, mini, maxi := a.analyzeData(a.getDurations(faChars))
	fmt.Println("Duration:", mean, stddev, mini, maxi)
	a.markDurationOutliers(faChars, mean, stddev)
	mean, stddev, mini, maxi = a.analyzeData(a.getSilence(faChars, betweenChars))
	fmt.Println("Between Chars:", mean, stddev, mini, maxi)
	a.markSilenceOutliers(faChars, mean, stddev, betweenChars, betweenCharsLong)
	mean, stddev, mini, maxi = a.analyzeData(a.getSilence(faChars, betweenWords))
	fmt.Println("Between Words:", mean, stddev, mini, maxi)
	a.markSilenceOutliers(faChars, mean, stddev, betweenWords, betweenWordsLong)
	mean, stddev, mini, maxi = a.analyzeData(a.getSilence(faChars, betweenVerses))
	fmt.Println("Between Verses:", mean, stddev, mini, maxi)
	a.markSilenceOutliers(faChars, mean, stddev, betweenVerses, betweenVersesLong)
	mean, stddev, mini, maxi = a.analyzeData(a.getSilence(faChars, betweenChapters))
	fmt.Println("Between Chapters:", mean, stddev, mini, maxi)
	a.markSilenceOutliers(faChars, mean, stddev, betweenChapters, betweenChaptersLong)
	faVerse = a.groupByVerse(faChars)
	a.addErrorCounts(faVerse)
	return faVerse, status
}

func (a *AlignErrorCalc) getDurations(chars []db.AlignChar) []float64 {
	var data []float64
	for _, ch := range chars {
		data = append(data, float64(ch.Duration))
	}
	return data
}

func (a *AlignErrorCalc) getSilence(chars []db.AlignChar, pos SilencePosition) []float64 {
	var data []float64
	posInt := int(pos)
	for _, ch := range chars {
		if ch.SilencePos == posInt {
			data = append(data, float64(ch.Silence))
		}
	}
	return data
}

func (a *AlignErrorCalc) analyzeData(data []float64) (mean, stddev, min, max float64) {
	mean = stat.Mean(data, nil)
	stddev = stat.StdDev(data, nil)
	min = data[0]
	max = data[0]
	for _, v := range data[1:] {
		min = math.Min(min, v)
		max = math.Max(max, v)
	}
	return mean, stddev, min, max
}

func (a *AlignErrorCalc) setCharSeq(chars []db.AlignChar) {
	var charSeq = 0
	var currWordId int64 = -1
	for i := range chars {
		if currWordId != chars[i].WordId {
			currWordId = chars[i].WordId
			charSeq = 0
		}
		chars[i].CharSeq = charSeq
		charSeq++
	}
}

func (a *AlignErrorCalc) markDurationOutliers(chars []db.AlignChar, mean float64, stddev float64) {
	var pct95 = mean + 2.2*stddev
	for i := range chars {
		if chars[i].Duration > pct95 {
			chars[i].DurationLong = int(durationLong)
		}
	}
}

func (a *AlignErrorCalc) markSilenceOutliers(chars []db.AlignChar, mean float64, stddev float64,
	silencePos SilencePosition, errorType ErrorType) {
	var pct95 = mean + 2.0*stddev
	for i := range chars {
		if chars[i].SilencePos == int(silencePos) {
			if chars[i].Silence > pct95 {
				if chars[i].SilenceLong == 0 {
					chars[i].SilenceLong = int(errorType)
				} else {
					fmt.Println("Skip Long Silence", silencePos, chars[i])
				}
			}
		}
	}
}

func (a *AlignErrorCalc) groupByVerse(chars []db.AlignChar) []AlignVerse {
	var verses []AlignVerse
	if len(chars) == 0 {
		return verses
	}
	currBookId := chars[0].BookId
	currChapter := chars[0].ChapterNum
	currVerse := chars[0].VerseStr
	start := 0
	for i, ch := range chars {
		if ch.BookId != currBookId || ch.ChapterNum != currChapter || ch.VerseStr != currVerse {
			oneVerse := chars[start:i]
			currBookId = ch.BookId
			currChapter = ch.ChapterNum
			currVerse = ch.VerseStr
			start = i
			var errCount int
			for _, ch1 := range oneVerse {
				if ch1.ScoreError > 0 || ch.DurationLong > 0 || ch.SilenceLong > 0 {
					errCount++
				}
			}
			if errCount > 0 {
				var verse AlignVerse
				verse.chars = oneVerse
				verses = append(verses, verse)
			}
		}
	}
	return verses
}

func (a *AlignErrorCalc) addErrorCounts(verses []AlignVerse) {
	for i := range verses {
		for _, ch := range verses[i].chars {
			if ch.ScoreError == int(scoreCritical) {
				verses[i].critScore++
			} else if ch.ScoreError == int(scoreQuestion) {
				verses[i].questScore++
			}
			if ch.DurationLong > 0 {
				verses[i].longDuration++
			}
			if ch.SilenceLong > 0 {
				verses[i].longSilence++
			}
		}
	}
}

func (a *AlignErrorCalc) countErrors(verses []AlignVerse) {
	var total int
	var critScoreError int
	var questScoreError int
	var durationLongCount int
	var count = make([]int, 8)
	for _, vs := range verses {
		for _, ch := range vs.chars {
			total++
			if ch.ScoreError == int(scoreCritical) {
				critScoreError++
			} else if ch.ScoreError == int(scoreQuestion) {
				questScoreError++
			}
			if ch.DurationLong > 0 {
				durationLongCount++
			}
			count[ch.SilenceLong]++
		}
	}
	fmt.Println("NO Error\t", count[noError]-critScoreError-questScoreError-durationLongCount)
	fmt.Println("ScoreCritical", critScoreError)
	fmt.Println("ScoreQuestion", questScoreError)
	fmt.Println("DurationLong", durationLongCount)
	fmt.Println("BetweenCharsLong", count[betweenCharsLong])
	fmt.Println("BetweenWordsLong", count[betweenWordsLong])
	fmt.Println("BetweenVersesLong", count[betweenVersesLong])
	fmt.Println("BetweenChaptersLong", count[betweenChaptersLong])
	fmt.Println("Total\t", total)
}
