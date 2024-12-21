package match

import (
	"context"
	"dataset"
	"dataset/db"
	"dataset/generic"
	"dataset/mms"
	"fmt"
	"gonum.org/v1/gonum/stat"
	"math"
	"strconv"
	"strings"
)

const (
	criticalThreshold = 0.0001 // 0.001
	questionThreshold = 0.01   // 0.01
	silenceStdevs     = 4.0    // intended to make it rare
)

type ErrorType int

const (
	noError ErrorType = iota
	scoreCritical
	scoreQuestion
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

type AlignErrorCalc struct {
	ctx     context.Context
	conn    db.DBAdapter
	lang    string
	altLang string
}

func NewAlignErrorCalc(ctx context.Context, conn db.DBAdapter, lang string, altLang string) AlignErrorCalc {
	var a AlignErrorCalc
	a.ctx = ctx
	a.conn = conn
	a.lang = lang
	a.altLang = altLang
	return a
}

func (a *AlignErrorCalc) Process(audioDirectory string) ([]generic.AlignLine, string, dataset.Status) {
	var faVerse []generic.AlignLine
	var faChars []generic.AlignChar
	var status dataset.Status
	faChars, status = a.conn.SelectFACharTimestamps()
	if status.IsErr {
		return faVerse, "", status
	}
	for i := range faChars {
		if faChars[i].FAScore <= criticalThreshold {
			faChars[i].ScoreError = int(scoreCritical)
		} else if faChars[i].FAScore <= questionThreshold {
			faChars[i].ScoreError = int(scoreQuestion)
		}
	}
	for i := 0; i < len(faChars)-1; i++ {
		var curr = faChars[i]
		var next = faChars[i+1]
		faChars[i].Silence = next.BeginTS - curr.EndTS
		if curr.WordId == next.WordId {
			faChars[i].SilencePos = int(betweenChars)
		} else if curr.LineId == next.LineId {
			faChars[i].SilencePos = int(betweenWords)
		} else if curr.AudioFile == next.AudioFile {
			faChars[i].SilencePos = int(betweenVerses)
		} else {
			faChars[i].SilencePos = int(betweenChapters)
			faChars[i].Silence = next.BeginTS
			// To be correct, I should add the silence of fileDuration - curr.EndTS
		}
	}
	mean, stddev, mini, maxi := a.analyzeData(a.getSilence(faChars, betweenChars))
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
	filenameMap, status := a.generateBookChapterFilenameMap()
	asr := mms.NewMMSASR(a.ctx, a.conn, a.lang, a.altLang)
	asr.ProcessAlignSilence(audioDirectory, faVerse)
	return faVerse, filenameMap, status
}

func (a *AlignErrorCalc) getDurations(chars []generic.AlignChar) []float64 {
	var data []float64
	for _, ch := range chars {
		data = append(data, float64(ch.Duration))
	}
	return data
}

func (a *AlignErrorCalc) getSilence(chars []generic.AlignChar, pos SilencePosition) []float64 {
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
	if len(data) == 0 {
		return 0.0, 0.0, 0.0, 0.0
	}
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

func (a *AlignErrorCalc) markSilenceOutliers(chars []generic.AlignChar, mean float64, stddev float64,
	silencePos SilencePosition, errorType ErrorType) {
	var pct95 = mean + silenceStdevs*stddev
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

func (a *AlignErrorCalc) groupByVerse(chars []generic.AlignChar) []generic.AlignLine {
	var result []generic.AlignLine
	if len(chars) == 0 {
		return result
	}
	currRef := chars[0].LineRef
	start := 0
	for i, ch := range chars {
		if ch.LineRef != currRef {
			currRef = ch.LineRef
			oneLine := make([]generic.AlignChar, i-start)
			copy(oneLine, chars[start:i])
			start = i
			var errCount int
			for _, ch1 := range oneLine {
				if ch1.ScoreError > 0 || ch.SilenceLong > 0 {
					errCount++
				}
			}
			if errCount > 0 {
				var line generic.AlignLine
				line.Chars = oneLine
				result = append(result, line)
			}
		}
	}
	return result
}

func (a *AlignErrorCalc) generateBookChapterFilenameMap() (string, dataset.Status) {
	chapters, status := a.conn.SelectBookChapterFilename()
	if status.IsErr {
		return "", status
	}
	var result []string
	result = append(result, "let fileMap = {\n")
	for i, ch := range chapters {
		key := ch.BookId + strconv.Itoa(ch.ChapterNum)
		result = append(result, "\t'"+key+"': '"+ch.AudioFile+"'")
		if i < len(chapters)-1 {
			result = append(result, ",\n")
		} else {
			result = append(result, "};\n")
		}
	}
	return strings.Join(result, ""), status
}

func (a *AlignErrorCalc) countErrors(lines []generic.AlignLine) {
	var total int
	var critScoreError int
	var questScoreError int
	var count = make([]int, 8)
	for _, vs := range lines {
		for _, ch := range vs.Chars {
			total++
			if ch.ScoreError == int(scoreCritical) {
				critScoreError++
			} else if ch.ScoreError == int(scoreQuestion) {
				questScoreError++
			}
			count[ch.SilenceLong]++
		}
	}
	fmt.Println("NO Error\t", count[noError]-critScoreError-questScoreError)
	fmt.Println("ScoreCritical", critScoreError)
	fmt.Println("ScoreQuestion", questScoreError)
	fmt.Println("BetweenCharsLong", count[betweenCharsLong])
	fmt.Println("BetweenWordsLong", count[betweenWordsLong])
	fmt.Println("BetweenVersesLong", count[betweenVersesLong])
	fmt.Println("BetweenChaptersLong", count[betweenChaptersLong])
	fmt.Println("Total\t", total)
}
