package match

import (
	"context"
	"dataset"
	"dataset/db"
	"dataset/fetch"
	log "dataset/logger"
	"dataset/request"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"golang.org/x/text/unicode/norm"
	"os"
	"strconv"
	"strings"
	"time"
)

type Compare struct {
	ctx         context.Context
	user        fetch.DBPUser
	baseDataset string
	dataset     string
	baseDb      db.DBAdapter
	database    db.DBAdapter
	settings    request.CompareSettings
	diffCount   int
	insertSum   int
	deleteSum   int
}

type Verse struct {
	bookId   string
	chapter  int
	verse    string
	verseEnd string
	text     string
}

func NewCompare(ctx context.Context, user fetch.DBPUser, baseDSet string, db db.DBAdapter,
	settings request.CompareSettings) Compare {
	var c Compare
	c.ctx = ctx
	c.user = user
	c.baseDataset = baseDSet
	c.dataset = strings.Split(db.Database, `.`)[0]
	c.database = db
	c.settings = settings
	return c
}

func (c *Compare) Process() (string, dataset.Status) {
	var filename string
	c.baseDb = db.NewerDBAdapter(c.ctx, false, c.user.Username, c.baseDataset)
	output, status := c.openOutput(c.baseDataset, c.dataset)
	if status.IsErr {
		return filename, status
	}
	filename = output.Name()
	for _, bookId := range db.BookNT {
		var chapInBook, _ = db.BookChapterMap[bookId] // Need to check OK, because bookId could be in error?
		var chapter = 1
		for chapter <= chapInBook {
			var lines1, lines2 []Verse
			lines1, status = c.process(c.baseDb, bookId, chapter)
			if status.IsErr {
				return filename, status
			}
			lines2, status = c.process(c.database, bookId, chapter)
			if status.IsErr {
				return filename, status
			}
			c.diff(output, lines1, lines2)
			chapter++
		}
	}
	c.baseDb.Close()
	fmt.Println("Num Diff", c.diffCount)
	_, _ = output.WriteString(`<p>Total Inserted Chars `)
	_, _ = output.WriteString(strconv.Itoa(c.insertSum))
	_, _ = output.WriteString(`, Total Deleted Chars `)
	_, _ = output.WriteString(strconv.Itoa(c.deleteSum))
	_, _ = output.WriteString("</p>\n")
	_, _ = output.WriteString(`<p>`)
	_, _ = output.WriteString("Total Difference Count: ")
	_, _ = output.WriteString(strconv.Itoa(c.diffCount))
	_, _ = output.WriteString("</p>\n")
	c.closeOutput(output)
	return filename, status
}

func (c *Compare) openOutput(database1 string, database2 string) (*os.File, dataset.Status) {
	var status dataset.Status
	output, err := os.CreateTemp(os.Getenv(`FCBH_DATASET_TMP`), c.dataset+"_*diff.html")
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
	_, _ = output.WriteString(head)
	_, _ = output.WriteString(`<h2 style="text-align:center">Compare `)
	_, _ = output.WriteString(database1)
	_, _ = output.WriteString(` to `)
	_, _ = output.WriteString(database2)
	_, _ = output.WriteString("</h2>\n")
	_, _ = output.WriteString(`<h3 style="text-align:center">`)
	_, _ = output.WriteString(time.Now().Format(`Mon Jan 2 2006 03:04:05 pm MST`))
	_, _ = output.WriteString("</h3>\n")
	_, _ = output.WriteString(`<h3 style="text-align:center">RED characters are those in `)
	_, _ = output.WriteString(database1)
	_, _ = output.WriteString(` only, while GREEN characters are in `)
	_, _ = output.WriteString(database2)
	_, _ = output.WriteString(" only</h3>\n")
	return output, status
}

func (c *Compare) closeOutput(output *os.File) {
	end := ` </body>
</html>`
	_, _ = output.WriteString(end)
	_ = output.Close()
}

func (c *Compare) process(conn db.DBAdapter, bookId string, chapterNum int) ([]Verse, dataset.Status) {
	var lines []Verse
	var status dataset.Status
	var ident db.Ident
	ident, status = conn.SelectIdent()
	if status.IsErr {
		return lines, status
	}
	scripts, status := conn.SelectScriptsByChapter(bookId, chapterNum)
	if status.IsErr {
		return lines, status
	}
	for _, script := range scripts {
		var vs Verse
		vs.bookId = script.BookId
		vs.chapter = script.ChapterNum
		vs.verse = script.VerseStr
		vs.verseEnd = script.VerseEnd
		vs.text = script.ScriptText
		lines = append(lines, vs)
	}
	if ident.TextSource == request.TextScript {
		lines = c.consolidateScript(lines)
	} else if ident.TextSource == request.TextUSXEdit {
		lines = c.consolidateUSX(lines)
	} else if ident.TextSource == request.TextPlainEdit {
		lines = c.consolidatePlainEdit(lines)
	}
	lines = c.cleanUp(lines)
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
			verse := Verse{bookId: bookId, chapter: chapter, verse: verseNum, text: part}
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
				log.Warn(c.ctx, bookId, chapter, verseNum, `Invalid char in {nn, expect n or } found `,
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
		if rec.chapter != lastChapter || rec.verse != lastVerse {
			if lastChapter != -1 {
				results = append(results, verse)
				sumOutput += len(verse.text)
			}
			lastChapter = rec.chapter
			lastVerse = rec.verse
			verse = Verse{bookId: rec.bookId, chapter: rec.chapter, verse: rec.verse, text: ``}
		}
		if !strings.HasSuffix(verse.text, ` `) && !strings.HasPrefix(rec.text, ` `) {
			verse.text += ` ` + rec.text
			sumOutput--
		} else {
			verse.text += rec.text
		}
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

func (c *Compare) consolidatePlainEdit(verses []Verse) []Verse {
	var results = make([]Verse, 0, len(verses))
	var first Verse
	for pos, rec := range verses {
		if pos == 0 {
			first = rec
		} else if rec.verse == `` {
			if !strings.HasSuffix(first.text, ` `) && !strings.HasPrefix(rec.text, ` `) {
				first.text += ` ` + rec.text
			} else {
				first.text += rec.text
			}
		} else {
			if len(first.text) > 0 {
				results = append(results, first)
				first = Verse{}
			}
			results = append(results, rec)
		}
	}
	return results
}

func (c *Compare) cleanUp(verses []Verse) []Verse {
	cfg := c.settings
	if cfg.LowerCase {
		for i, vs := range verses {
			verses[i].text = strings.ToLower(vs.text)
		}
	}
	var replace []string
	if cfg.RemovePromptChars {
		replace = append(replace, "<<", "")
		replace = append(replace, ">>", "")
		replace = append(replace, "((", "")
		replace = append(replace, "†", "") // KDLTBL
		replace = append(replace, "[", "") // KDLTBL
		replace = append(replace, "]", "") // KDLTBL
		replace = append(replace, "<", "") // AAAMLT
		replace = append(replace, ">", "") // AAAMLT
	}
	if cfg.RemovePunctuation {
		replace = append(replace, ",", "")
		replace = append(replace, ";", "")
		replace = append(replace, ":", "")
		replace = append(replace, ".", "")
		replace = append(replace, "!", "")
		replace = append(replace, "?", "")
		replace = append(replace, "~", "")
	}
	if cfg.DoubleQuotes.Normalize {
		//replace = append(replace, "\u201C", "\u0022") // left quote
		//replace = append(replace, "\u201D", "\u0022") // right quote
		//replace = append(replace, "\u2033", "\u0022") // minutes or seconds
		replace = append(replace, "\u00AB", "\u0022") // <<
		//replace = append(replace, "\u00BB", "\u0022") // >>
		//replace = append(replace, "\u201E", "\u0022") // low 9 quote
	} else if cfg.DoubleQuotes.Remove {
		replace = append(replace, "\u0022", "") // ascii
		replace = append(replace, "\u201C", "") // left quote
		replace = append(replace, "\u201D", "") // right quote
		//replace = append(replace, "\u2033", "") // minutes or seconds
		replace = append(replace, "\u00AB", "") // <<
		replace = append(replace, "\u00BB", "") // >>
		//replace = append(replace, "\u201E", "") // low 9 quote
	}
	if cfg.Apostrophe.Normalize {
		replace = append(replace, "\uA78C", "'") // ? found in script text
		replace = append(replace, "\u2018", "'") // left
		replace = append(replace, "\u2019", "'") // right
		replace = append(replace, "\u02BC", "'") // modifier letter apos
		replace = append(replace, "\u2032", "'") // prime
		replace = append(replace, "\u00B4", "'") // acute accent (not really apos)
	} else if cfg.Apostrophe.Remove {
		replace = append(replace, "'", "")
		replace = append(replace, "\uA78C", "") // ? found in script text
		replace = append(replace, "\u2018", "") // left
		replace = append(replace, "\u2019", "") // right
		//replace = append(replace, "\u02BC", "") // modifier letter apos
		//replace = append(replace, "\u2032", "") // prime
		//replace = append(replace, "\u00B4", "") // acute accent (not really apos)
	}
	if cfg.Hyphen.Normalize {
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
	} else if cfg.Hyphen.Remove {
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
	//if cfg.DiacriticalMarks.Normalize {
	if cfg.DiacriticalMarks.NormalizeNFC {
		for i, vs := range verses {
			verses[i].text = norm.NFC.String(vs.text)
		}
	} else if cfg.DiacriticalMarks.NormalizeNFD {
		for i, vs := range verses {
			verses[i].text = norm.NFD.String(vs.text)
		}
	} else if cfg.DiacriticalMarks.NormalizeNFKC {
		for i, vs := range verses {
			verses[i].text = norm.NFKC.String(vs.text)
		}
	} else if cfg.DiacriticalMarks.NormalizeNFKD {
		for i, vs := range verses {
			verses[i].text = norm.NFKD.String(vs.text)
		}
	} else if cfg.DiacriticalMarks.Remove {
		for i, vs := range verses {
			verses[i].text = norm.NFKD.String(vs.text)
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
		verse2Map[vs2.verse] = vs2
	}
	// combine the verse2 to verse1 that match
	var didMatch = make(map[string]bool)
	var pairs = make([]pair, 0, len(verses1))
	for _, vs1 := range verses1 {
		vs2, ok := verse2Map[vs1.verse]
		if ok {
			didMatch[vs1.verse] = true
		}
		p := pair{bookId: vs1.bookId, chapter: vs1.chapter, num: vs1.verse, text1: vs1.text, text2: vs2.text}
		pairs = append(pairs, p)
	}
	// pick up any verse2 that did not match verse1
	for _, vs2 := range verses2 {
		_, ok := didMatch[vs2.verse]
		if !ok {
			p := pair{bookId: vs2.bookId, chapter: vs2.chapter, num: vs2.verse, text1: ``, text2: vs2.text}
			pairs = append(pairs, p)
		}
	}
	// perform a match on pairs
	diffMatch := diffmatchpatch.New()
	for _, pair := range pairs {
		if len(pair.text1) > 0 || len(pair.text2) > 0 {
			diffs := diffMatch.DiffMain(pair.text1, pair.text2, false)
			if !c.isMatch(diffs) {
				inserts, deletes := c.measure(diffs)
				c.insertSum += inserts
				c.deleteSum += deletes
				ref := pair.bookId + " " + strconv.Itoa(pair.chapter) + ":" + pair.num + ` `
				//fmt.Println(ref, diffMatch.DiffPrettyText(diffs))
				//fmt.Println("=============")
				_, _ = output.WriteString(`<p>`)
				_, _ = output.WriteString(ref)
				_, _ = output.WriteString(` +`)
				_, _ = output.WriteString(strconv.Itoa(inserts))
				_, _ = output.WriteString(` -`)
				_, _ = output.WriteString(strconv.Itoa(deletes))
				_, _ = output.WriteString(` `)
				_, _ = output.WriteString(diffMatch.DiffPrettyHtml(diffs))
				_, _ = output.WriteString("</p>\n")
				c.diffCount++
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

func (c *Compare) measure(diffs []diffmatchpatch.Diff) (int, int) {
	var inserts = 0
	var deletes = 0
	for _, diff := range diffs {
		if diff.Type == diffmatchpatch.DiffInsert {
			inserts += len(diff.Text)
		} else if diff.Type == diffmatchpatch.DiffDelete {
			deletes += len(diff.Text)
		}
	}
	return inserts, deletes
}
