package diff

import (
	"context"
	"dataset/db"
	log "dataset/logger"
	"dataset/mms"
	"dataset/request"
	"github.com/sergi/go-diff/diffmatchpatch"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strconv"
	"strings"
)

type Compare struct {
	ctx         context.Context
	user        string
	baseDataset string
	dataset     string
	baseDb      db.DBAdapter
	database    db.DBAdapter
	lang        string
	baseIdent   db.Ident
	compIdent   db.Ident
	testament   request.Testament
	settings    request.CompareSettings
	replacer    *strings.Replacer
	verseRm     *regexp.Regexp
	writer      HTMLWriter
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
	uRoman   string
	beginTS  float64
	endTS    float64
}

func NewCompare(ctx context.Context, user string, baseDSet string, db db.DBAdapter,
	lang string, testament request.Testament, settings request.CompareSettings) Compare {
	var c Compare
	c.ctx = ctx
	c.user = user
	c.baseDataset = baseDSet
	c.dataset = strings.Split(db.Database, `.`)[0]
	c.database = db
	c.lang = lang
	c.testament = testament
	c.settings = settings
	c.replacer = c.cleanUpSetup()
	c.verseRm = regexp.MustCompile(`\{[0-9\-\,]+\}\s?`) // used by compareScriptLine
	return c
}

func (c *Compare) Process() (string, *log.Status) {
	var filename string
	var status *log.Status
	c.baseDb, status = db.NewerDBAdapter(c.ctx, false, c.user, c.baseDataset)
	if status != nil {
		return filename, status
	}
	status = mms.EnsureUroman(c.baseDb, c.lang)
	if status != nil {
		return filename, status
	}
	c.writer, status = NewHTMLWriter(c.ctx, c.dataset)
	if status != nil {
		return filename, status
	}
	var compHasVerse, compHasLine, baseHasVerse, baseHasLine bool
	compHasVerse, compHasLine, c.compIdent, status = c.hasVerseLine(c.database)
	if status != nil {
		return filename, status
	}
	baseHasVerse, baseHasLine, c.baseIdent, status = c.hasVerseLine(c.baseDb)
	if status != nil {
		return filename, status
	}
	if compHasVerse && baseHasVerse {
		filename, status = c.CompareVerses()
	} else if compHasLine && baseHasLine {
		filename, status = c.CompareScriptLines()
	} else {
		filename, status = c.CompareChapters()
	}
	return filename, status
}

func (c *Compare) hasVerseLine(conn db.DBAdapter) (bool, bool, db.Ident, *log.Status) {
	var hasVerse bool
	var hasLine bool
	var ident db.Ident
	var status *log.Status
	ident, status = conn.SelectIdent()
	if status != nil {
		return hasVerse, hasLine, ident, status
	}
	verseColLen, status := conn.SelectVerseLength()
	if status != nil {
		return hasVerse, hasLine, ident, status
	}
	hasVerse = verseColLen > 0 || ident.TextSource == request.TextScript
	lineColLen, status := conn.SelectScriptLineLength()
	if status != nil {
		return hasVerse, hasLine, ident, status
	}
	hasLine = lineColLen > 0 && (ident.TextSource == request.TextScript || ident.TextSource == request.TextCSV)
	return hasVerse, hasLine, ident, status
}

func (c *Compare) CompareVerses() (string, *log.Status) {
	var filename string
	var status *log.Status
	filename = c.writer.WriteHeading(c.baseDataset)
	for _, bookId := range db.RequestedBooks(c.testament) {
		var chapInBook, _ = db.BookChapterMap[bookId] // Need to check OK, because bookId could be in error?
		var chapter = 1
		for chapter <= chapInBook {
			var lines1, lines2 []Verse
			lines1, status = c.process(c.baseDb, bookId, chapter)
			if status != nil {
				return filename, status
			}
			lines2, status = c.process(c.database, bookId, chapter)
			if status != nil {
				return filename, status
			}
			c.diff(lines1, lines2)
			chapter++
		}
	}
	filenameMap, status := c.generateBookChapterFilenameMap()
	c.baseDb.Close()
	c.writer.WriteEnd(filenameMap, c.insertSum, c.deleteSum, c.diffCount)
	return filename, status
}

func (c *Compare) process(conn db.DBAdapter, bookId string, chapterNum int) ([]Verse, *log.Status) {
	var lines []Verse
	var status *log.Status
	var ident db.Ident
	ident, status = conn.SelectIdent()
	if status != nil {
		return lines, status
	}
	scripts, status := conn.SelectScriptsByChapter(bookId, chapterNum)
	if status != nil {
		return lines, status
	}
	for _, script := range scripts {
		var vs Verse
		vs.bookId = script.BookId
		vs.chapter = script.ChapterNum
		vs.verse = script.VerseStr
		vs.verseEnd = script.VerseEnd
		vs.text = script.ScriptText
		vs.uRoman = script.URoman
		vs.beginTS = script.ScriptBeginTS
		vs.endTS = script.ScriptEndTS
		lines = append(lines, vs)
	}
	if ident.TextSource == request.TextScript {
		lines = c.consolidateScript(lines)
	}
	lines = c.cleanUpVerses(lines)
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
	var verseNum = `0`
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

func (c *Compare) cleanUpSetup() *strings.Replacer {
	var replace []string
	cfg := c.settings
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
	var replacer *strings.Replacer
	if len(replace) > 0 {
		replacer = strings.NewReplacer(replace...)
	}
	return replacer
}

func (c *Compare) cleanUpVerses(verses []Verse) []Verse {
	for i, vs := range verses {
		verses[i].text = c.cleanup(vs.text)
		verses[i].uRoman = c.cleanup(vs.uRoman)
	}
	return verses
}

func (c *Compare) cleanup(text string) string {
	if c.replacer != nil {
		text = c.replacer.Replace(text)
	}
	cfg := c.settings
	if cfg.LowerCase {
		text = strings.ToLower(text)
	}
	// https://unicode.org/reports/tr15/  Normalization Doc
	//if cfg.DiacriticalMarks.Normalize {
	if cfg.DiacriticalMarks.NormalizeNFC {
		text = norm.NFC.String(text)
	} else if cfg.DiacriticalMarks.NormalizeNFD {
		text = norm.NFD.String(text)
	} else if cfg.DiacriticalMarks.NormalizeNFKC {
		text = norm.NFKC.String(text)
	} else if cfg.DiacriticalMarks.NormalizeNFKD {
		text = norm.NFKD.String(text)
	} else if cfg.DiacriticalMarks.Remove {
		text = norm.NFKD.String(text)
		var filtered = make([]rune, 0, 300)
		for _, ch := range text {
			if ch < '\u0300' || ch > '\u036F' {
				filtered = append(filtered, ch)
			}
		}
		text = string(filtered)
	}
	return text
}

type pair struct {
	bookId  string
	chapter int
	num     string
	beginTS float64
	endTS   float64
	text1   string
	text2   string
}

/* This diff method assumes one chapter at a time */
func (c *Compare) diff(verses1 []Verse, verses2 []Verse) {

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
		p := pair{bookId: vs1.bookId, chapter: vs1.chapter, num: vs1.verse, beginTS: vs1.beginTS, endTS: vs1.endTS, text1: vs1.uRoman, text2: vs2.uRoman}
		pairs = append(pairs, p)
	}
	// pick up any verse2 that did not match verse1
	for _, vs2 := range verses2 {
		_, ok := didMatch[vs2.verse]
		if !ok {
			p := pair{bookId: vs2.bookId, chapter: vs2.chapter, num: vs2.verse, beginTS: vs2.beginTS, endTS: vs2.endTS, text1: ``, text2: vs2.uRoman}
			pairs = append(pairs, p)
		}
	}
	// perform a match on pairs
	diffMatch := diffmatchpatch.New()
	for _, par := range pairs {
		if len(par.text1) > 0 || len(par.text2) > 0 {
			par.text1 = strings.TrimSpace(par.text1)
			par.text2 = strings.TrimSpace(par.text2)
			diffs := diffMatch.DiffMain(par.text1, par.text2, false)
			diffMatch.DiffCleanupMerge(diffs) // required for measure to compute largest
			if !c.isMatch(diffs) {
				inserts, deletes := c.measure(diffs)
				largest := c.largestLength(diffs)
				c.insertSum += inserts
				c.deleteSum += deletes
				avgLen := float64(len(par.text1)+len(par.text2)) / 2.0
				errPct := float64((inserts+deletes)*100) / avgLen
				c.writer.WriteVerseDiff(par, inserts, deletes, largest, errPct, diffMatch.DiffPrettyHtml(diffs))
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

// ensureClean is deprecated, we are using diff_clean_merge
func (c *Compare) ensureClean(diffs []diffmatchpatch.Diff) {
	var prior diffmatchpatch.Operation
	for _, diff := range diffs {
		if diff.Type == prior {
			log.Warn(c.ctx, "There are inserts and deletes together, must use diff_clean_merge")
		}
		prior = diff.Type
	}
}

func (c *Compare) measure(diffs []diffmatchpatch.Diff) (int, int) {
	var inserts = 0
	var deletes = 0
	for _, diff := range diffs {
		lenText := len(diff.Text)
		if diff.Type == diffmatchpatch.DiffInsert {
			inserts += lenText
		} else if diff.Type == diffmatchpatch.DiffDelete {
			deletes += lenText
		}
	}
	return inserts, deletes
}

func (c *Compare) largestLength(diffs []diffmatchpatch.Diff) int {
	var result int
	var length int
	for _, diff := range diffs {
		if diff.Type != diffmatchpatch.DiffEqual {
			length += len(diff.Text)
		} else {
			if length > result {
				result = length
			}
			length = 0
		}
	}
	if length > result {
		result = length
	}
	return result
}

func (c *Compare) generateBookChapterFilenameMap() (string, *log.Status) {
	chapters, status := c.database.SelectBookChapterFilename()
	if status != nil {
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
