package match

import (
	"context"
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"golang.org/x/text/unicode/norm"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var globalDiffCount = 0

type Compare struct {
	ctx       context.Context
	database1 string
	database2 string
	db1       db.DBAdapter
	db2       db.DBAdapter
}

type Verse struct {
	bookId  string
	chapter int
	num     string
	text    string
}

func NewCompare(ctx context.Context, db1 string, db2 string) Compare {
	var c Compare
	c.ctx = ctx
	c.database1 = db1
	c.database2 = db2
	return c
}

func (c *Compare) Process() dataset.Status {
	c.db1 = db.NewDBAdapter(c.ctx, c.database1)
	c.db2 = db.NewDBAdapter(c.ctx, c.database2)
	var numChapters, status = c.db1.ReadNumChapters()
	if status.IsErr {
		return status
	}
	var output *os.File
	output, status = c.openOutput(c.database1, c.database2)
	if status.IsErr {
		return status
	}
	for _, bookId := range db.BookNT {
		var chapInBook, _ = numChapters[bookId] // Need to check OK, because bookId could be in error?
		var chapter = 1
		for chapter <= chapInBook {
			lines1, status := c.process(c.db1, c.database1, bookId, chapter)
			if status.IsErr {
				return status
			}
			lines2, status := c.process(c.db2, c.database2, bookId, chapter)
			if status.IsErr {
				return status
			}
			c.diff(output, lines1, lines2)
			chapter++
		}
	}
	c.db1.Close()
	c.db2.Close()
	fmt.Println("Num Diff", globalDiffCount)
	output.WriteString(`<p>`)
	output.WriteString("TOTAL Difference Count: ")
	output.WriteString(strconv.Itoa(globalDiffCount))
	output.WriteString("</p>\n")
	c.closeOutput(output)
	return status
}

func (c *Compare) openOutput(database1 string, database2 string) (*os.File, dataset.Status) {
	var status dataset.Status
	bibleId := database1[:6]
	outPath := filepath.Join(os.Getenv("FCBH_DATASET_FILES"), bibleId, "diff.html")
	output, err := os.Create(outPath)
	if err != nil {
		status = log.Error(c.ctx, 500, err, `Error creating output file for diff`)
		return output, status
	}
	head := `<DOCTYPE html>
<html>
 <head>
  <meta charset="utf-8">
  <title>File Difference</title>
  <style>
p { margin: 20px 40px; }
  </style>
 </head>
 <body>`
	output.WriteString(head)
	output.WriteString(`<h2 style="text-align:center">Compare `)
	output.WriteString(database1)
	output.WriteString(` to `)
	output.WriteString(database2)
	output.WriteString("</h2>\n")
	output.WriteString(`<h3 style="text-align:center">`)
	output.WriteString(time.Now().Format(`Mon Jan 2 2006 03:04:05 pm MST`))
	output.WriteString("</h3>\n")
	output.WriteString(`<h3 style="text-align:center">RED characters are those in `)
	output.WriteString(database1[7:])
	output.WriteString(` only, while GREEN characters are in `)
	output.WriteString(database2[7:])
	output.WriteString(" only</h3>\n")
	return output, status
}

func (c *Compare) closeOutput(output *os.File) {
	end := ` </body>
</html>`
	output.WriteString(end)
	output.Close()
}

func (c *Compare) process(db db.DBAdapter, database string, bookId string, chapterNum int) ([]Verse, dataset.Status) {
	var lines []Verse
	var status dataset.Status
	scripts, status := db.ReadScriptsByChapter(bookId, chapterNum)
	for _, script := range scripts {
		var vs Verse
		vs.bookId = script.BookId
		vs.chapter = script.ChapterNum
		vs.num = script.VerseStr
		vs.text = script.ScriptText
		lines = append(lines, vs)
	}
	if status.IsErr {
		return lines, status
	}
	if strings.Contains(database, "SCRIPT") {
		lines = c.consolidateScript(lines)
	} else if strings.Contains(database, "USX") {
		lines = c.consolidateUSX(lines)
	}
	cfg := getConfig()
	lines = c.cleanUp(lines, cfg)
	return lines, status
}

func (c *Compare) consolidateScript(verses []Verse) []Verse {
	const (
		begin = iota + 1
		inNum
		endNum
	)
	//var labels = []string{``, `BEGIN`, `INNUM`, `ENDNUM`}
	var results = make([]Verse, 0, len(verses))
	var sumInput = 0
	var sumOutput = 0
	var bookId = verses[0].bookId
	var chapter = verses[0].chapter
	var parts = make([]string, 0, 100)
	for _, rec := range verses {
		parts = append(parts, rec.text)
		sumInput += len(rec.text)
	}
	text := strings.Join(parts, ``)
	var verseNum = ``
	var tmpNum []byte
	var index = 0
	var state = begin
	for index < len(text) {
		switch state {
		case begin:
			var part string
			search := text[index:]
			pos := strings.Index(search, `{`)
			if pos < 0 {
				part = search
				sumOutput += len(part)
			} else {
				part = search[:pos]
				state = inNum
				tmpNum = []byte{}
				sumOutput += len(part) + 1
			}
			verse := Verse{bookId: bookId, chapter: chapter, num: verseNum, text: part}
			results = append(results, verse)
			index += len(part) + 1
		case inNum:
			char := text[index]
			if char >= '0' && char <= '9' {
				tmpNum = append(tmpNum, char)
				index++
				sumOutput += 1
			} else if char == '}' {
				verseNum = string(tmpNum)
				state = endNum
				index++
				sumOutput += 1
			} else {
				start := max(0, index-50)
				end := min(len(text)-1, index+50)
				fmt.Println("WARNING:", bookId, chapter, verseNum, `Invalid char in {nn, expect n or } found `,
					string(char), ` in `, text[start:end])
				verseNum = string(tmpNum)
				state = begin
			}
		case endNum:
			char := text[index]
			peek := text[index+1]
			if (char == '_' || char == '-') && peek == '{' {
				state = inNum
				tmpNum = []byte(verseNum + "-")
				index += 2
				sumOutput += 2
			} else {
				state = begin
			}
		}
	}
	if sumInput != sumOutput {
		log.Warn(c.ctx, "Bug: Not all data processed by consolidateScript input:", sumInput, " output:", sumOutput)
	}
	return results
}

func (c *Compare) consolidateUSX(verses []Verse) []Verse {
	var sumInput = 0
	var sumOutput = 0
	var results = make([]Verse, 0, len(verses))
	var lastChapter = -1
	var lastVerse = ``
	var verse Verse
	for _, rec := range verses {
		sumInput += len(rec.text)
		if rec.chapter != lastChapter || rec.num != lastVerse {
			if lastChapter != -1 {
				results = append(results, verse)
				sumOutput += len(verse.text)
			}
			lastChapter = rec.chapter
			lastVerse = rec.num
			verse = Verse{bookId: rec.bookId, chapter: rec.chapter, num: rec.num, text: ``}
		}
		verse.text += rec.text
	}
	if len(verse.text) > 0 {
		results = append(results, verse)
		sumOutput += len(verse.text)
	}
	if sumInput != sumOutput {
		log.Warn(c.ctx, "Bug: Not all data processed by consolidateUSX input:", sumInput, " output:", sumOutput)
	}
	return results
}

func (c *Compare) cleanUp(verses []Verse, cfg config) []Verse {
	if cfg.lowerCase {
		for i, vs := range verses {
			verses[i].text = strings.ToLower(vs.text)
		}
	}
	var replace []string
	if cfg.removePromptChars {
		replace = append(replace, "<<", "")
		replace = append(replace, ">>", "")
		replace = append(replace, "((", "")
		replace = append(replace, "†", "") // KDLTBL
		replace = append(replace, "[", "") // KDLTBL
		replace = append(replace, "]", "") // KDLTBL
		replace = append(replace, "<", "") // AAAMLT
		replace = append(replace, ">", "") // AAAMLT
	}
	if cfg.removePunctuation {
		replace = append(replace, ",", "")
		replace = append(replace, ";", "")
		replace = append(replace, ":", "")
		replace = append(replace, ".", "")
		replace = append(replace, "!", "")
		replace = append(replace, "?", "")
		replace = append(replace, "~", "")
	}
	if cfg.doubleQuotes == normalize {
		//replace = append(replace, "\u201C", "\u0022") // left quote
		//replace = append(replace, "\u201D", "\u0022") // right quote
		//replace = append(replace, "\u2033", "\u0022") // minutes or seconds
		replace = append(replace, "\u00AB", "\u0022") // <<
		//replace = append(replace, "\u00BB", "\u0022") // >>
		//replace = append(replace, "\u201E", "\u0022") // low 9 quote
	} else if cfg.doubleQuotes == remove {
		replace = append(replace, "\u0022", "") // ascii
		replace = append(replace, "\u201C", "") // left quote
		replace = append(replace, "\u201D", "") // right quote
		//replace = append(replace, "\u2033", "") // minutes or seconds
		replace = append(replace, "\u00AB", "") // <<
		replace = append(replace, "\u00BB", "") // >>
		//replace = append(replace, "\u201E", "") // low 9 quote
	}
	if cfg.apostrophe == normalize {
		replace = append(replace, "\uA78C", "'") // ? found in script text
		replace = append(replace, "\u2018", "'") // left
		replace = append(replace, "\u2019", "'") // right
		replace = append(replace, "\u02BC", "'") // modifier letter apos
		replace = append(replace, "\u2032", "'") // prime
		replace = append(replace, "\u00B4", "'") // acute accent (not really apos)
	} else if cfg.apostrophe == remove {
		replace = append(replace, "'", "")
		replace = append(replace, "\uA78C", "") // ? found in script text
		replace = append(replace, "\u2018", "") // left
		replace = append(replace, "\u2019", "") // right
		//replace = append(replace, "\u02BC", "") // modifier letter apos
		//replace = append(replace, "\u2032", "") // prime
		//replace = append(replace, "\u00B4", "") // acute accent (not really apos)
	}
	if cfg.hyphen == normalize {
		replace = append(replace, "\u2010", "\u002D") // hypen
		replace = append(replace, "\u2011", "\u002D") // non-breaking hyphen
		replace = append(replace, "\u2012", "\u002D") // figure dash
		replace = append(replace, "\u2013", "\u002D") // en dash
		replace = append(replace, "\u2014", "\u002D") // em dash
		replace = append(replace, "\u2015", "\u002D") // horizontal bar
		replace = append(replace, "\uFE58", "\u002D") // small em dash
		replace = append(replace, "\uFE62", "\u002D") // small en dash
		replace = append(replace, "\uFE63", "\u002D") // small hyphen minus
		replace = append(replace, "\uFF0D", "\u002D") // fullwidth hypen-minus
	} else if cfg.hyphen == remove {
		replace = append(replace, "\u002D", "") // hypen
		replace = append(replace, "\u2010", "") // hypen
		replace = append(replace, "\u2011", "") // non-breaking hyphen
		replace = append(replace, "\u2012", "") // figure dash
		replace = append(replace, "\u2013", "") // en dash
		replace = append(replace, "\u2014", "") // em dash
		replace = append(replace, "\u2015", "") // horizontal bar
		replace = append(replace, "\uFE58", "") // small em dash
		replace = append(replace, "\uFE62", "") // small en dash
		replace = append(replace, "\uFE63", "") // small hyphen minus
		replace = append(replace, "\uFF0D", "") // fullwidth hypen-minus
	}
	if len(replace) > 0 {
		replacer := strings.NewReplacer(replace...)
		for i, vs := range verses {
			verses[i].text = replacer.Replace(vs.text)
		}
	}
	// https://unicode.org/reports/tr15/  Normalization Doc
	if cfg.diacritical == normalize {
		if cfg.normalizeType == norm.NFC {
			for i, vs := range verses {
				verses[i].text = norm.NFC.String(vs.text)
			}
		} else if cfg.normalizeType == norm.NFD {
			for i, vs := range verses {
				verses[i].text = norm.NFD.String(vs.text)
			}
		} else if cfg.normalizeType == norm.NFKC {
			for i, vs := range verses {
				verses[i].text = norm.NFKC.String(vs.text)
			}
		} else if cfg.normalizeType == norm.NFKD {
			for i, vs := range verses {
				verses[i].text = norm.NFKD.String(vs.text)
			}
		}
	} else if cfg.diacritical == remove {
		if cfg.normalizeType != norm.NFD && cfg.normalizeType != norm.NFKD {
			for i, vs := range verses {
				verses[i].text = norm.NFKD.String(vs.text)
			}
		}
		var filtered = make([]rune, 0, 300)
		for i, vs := range verses {
			filtered = []rune{}
			for _, ch := range vs.text {
				if ch < '\u0300' || ch > '\u036F' {
					filtered = append(filtered, ch)
				}
			}
			verses[i].text = string(filtered)
		}
	}
	return verses
}

/* This diff method assumes one chapter at a time */
func (c *Compare) diff(output *os.File, verses1 []Verse, verses2 []Verse) {
	type pair struct {
		bookId  string
		chapter int
		num     string
		text1   string
		text2   string
	}
	// Put the second data in a map
	var verse2Map = make(map[string]Verse)
	for _, vs2 := range verses2 {
		verse2Map[vs2.num] = vs2
	}
	// combine the verse2 to verse1 that match
	var didMatch = make(map[string]bool)
	var pairs = make([]pair, 0, len(verses1))
	for _, vs1 := range verses1 {
		vs2, ok := verse2Map[vs1.num]
		if ok {
			didMatch[vs1.num] = true
		}
		p := pair{bookId: vs1.bookId, chapter: vs1.chapter, num: vs1.num, text1: vs1.text, text2: vs2.text}
		pairs = append(pairs, p)
	}
	// pick up any verse2 that did not match verse1
	for _, vs2 := range verses2 {
		_, ok := didMatch[vs2.num]
		if !ok {
			p := pair{bookId: vs2.bookId, chapter: vs2.chapter, num: vs2.num, text1: ``, text2: vs2.text}
			pairs = append(pairs, p)
		}
	}
	// perform a match on pairs
	diffMatch := diffmatchpatch.New()
	for _, pair := range pairs {
		if len(pair.text1) > 0 || len(pair.text2) > 0 {
			diffs := diffMatch.DiffMain(pair.text1, pair.text2, false)
			if !c.isMatch(diffs) {
				ref := pair.bookId + " " + strconv.Itoa(pair.chapter) + ":" + pair.num + ` `
				fmt.Println(ref, diffMatch.DiffPrettyText(diffs))
				fmt.Println("=============")
				output.WriteString(`<p>`)
				output.WriteString(ref)
				output.WriteString(diffMatch.DiffPrettyHtml(diffs))
				//output.WriteString("<br/><br>")
				output.WriteString("</p>\n")
				globalDiffCount++
			}
		}
	}
}

func (c *Compare) isMatch(diffs []diffmatchpatch.Diff) bool {
	for _, diff := range diffs {
		if diff.Type == diffmatchpatch.DiffInsert || diff.Type == diffmatchpatch.DiffDelete {
			if len(strings.TrimSpace(diff.Text)) > 0 {
				return false
			}
		}
	}
	return true
}
