package read

import (
	"context"
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"encoding/xml"
	"fmt"
	"io"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Stack []string

var hasStyle = map[string]bool{
	`book`: true, `para`: true, `char`: true, `cell`: true, `ms`: true, `note`: true, `sidebar`: true, `figure`: true}

var numericPattern = regexp.MustCompile(`^\d+`)

type USXParser struct {
	ctx  context.Context
	conn db.DBAdapter
}

func NewUSXParser(conn db.DBAdapter) USXParser {
	var p USXParser
	p.ctx = conn.Ctx
	p.conn = conn
	return p
}

func (p *USXParser) ProcessFiles(inputFiles []InputFile) dataset.Status {
	var status dataset.Status
	for _, file := range inputFiles {
		filename := filepath.Join(file.Directory, file.Filename)
		fmt.Println(filename)
		var records []db.Script
		records, status = p.decode(p.ctx, filename) // Also edits out non-script elements
		if !status.IsErr {
			titleDesc := p.extractTitles(records)
			records = p.addChapterHeading(records, titleDesc)
			records = p.correctScriptNum(records)
			status = p.conn.InsertScripts(records)
			if status.IsErr {
				return status
			}
		}
	}
	count, status := p.conn.SelectScalarInt(`SELECT count(*) FROM scripts`)
	fmt.Println("Total Records", count)
	return status
}

func (p *USXParser) decode(ctx context.Context, filename string) ([]db.Script, dataset.Status) {
	var records []db.Script
	var status dataset.Status
	xmlFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer xmlFile.Close()
	var stack Stack
	var rec db.Script
	var tagName string
	var bookId string
	var chapterNum = 1
	var scriptNum = 0
	var verseNum int
	var verseStr string
	var usfmStyle string
	decoder := xml.NewDecoder(xmlFile)
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break // End of file
		}
		if err != nil {
			status = log.Error(ctx, 500, err, "Error parsing USX.")
			return records, status
		}
		switch se := token.(type) {
		//[StartElement], [EndElement], [CharData], [Comment], [ProcInst], [Directive].
		case xml.StartElement:
			tagName = se.Name.Local
			if tagName == `book` {
				bookId = p.findAttr(se, `code`)
			} else if tagName == `chapter` {
				chapterNum = p.findIntAttr(se, `number`)
			} else if tagName == `verse` {
				verseNum = p.findIntAttr(se, `number`)
				verseStr = p.findAttr(se, `number`)
			}
			if hasStyle[tagName] {
				usfmStyle = tagName + `.` + p.findAttr(se, `style`)
				if include(usfmStyle) { // This if fileters out the non-script code
					stack = stack.Push(usfmStyle)
				} else {
					err := decoder.Skip()
					if err != nil {
						status = log.Error(ctx, 500, err, "Error in USX parser")
						return records, status
					}
				}
			}
		case xml.CharData:
			text := string(se)
			if len(strings.TrimSpace(text)) > 0 {
				// This is needed because scripts use {n} as a verse number.
				if strings.Contains(text, "{") || strings.Contains(text, "}") {
					fmt.Println("Found { } changed to ( ) in:", text)
					text = strings.Replace(text, `{`, `(`, -1)
					text = strings.Replace(text, `}`, `)`, -1)
				}
				rec.ScriptTexts = append(rec.ScriptTexts, text)
			}
		case xml.EndElement:
			if hasStyle[se.Name.Local] {
				stack, usfmStyle = stack.Pop()
			}
		}
		if chapterNum != rec.ChapterNum || verseNum != rec.VerseNum || usfmStyle != rec.UsfmStyle {
			if bookId != `` && len(rec.ScriptTexts) > 0 {
				records = append(records, rec)
			}
			scriptNum++
			if chapterNum != rec.ChapterNum {
				scriptNum = 1
			}
			rec = db.Script{BookId: bookId, ChapterNum: chapterNum, ScriptNum: strconv.Itoa(scriptNum),
				VerseNum: verseNum, VerseStr: verseStr, UsfmStyle: usfmStyle}
		}
	}
	fmt.Println("num records", len(records))
	return records, status
}

type titleDesc struct {
	heading string
	title   []db.Script
	areDiff bool
}

func (p *USXParser) extractTitles(records []db.Script) titleDesc {
	var results titleDesc
	for _, rec := range records {
		if rec.UsfmStyle == `para.h` {
			results.heading = strings.Join(rec.ScriptTexts, ``)
		} else if strings.HasPrefix(rec.UsfmStyle, `para.mt`) {
			results.title = append(results.title, rec)
		}
	}
	if results.heading == `` && len(results.title) > 0 {
		results.heading = strings.Join(results.title[len(results.title)-1].ScriptTexts, ``)
	}
	return results
}

func (p *USXParser) addChapterHeading(records []db.Script, titles titleDesc) []db.Script {
	var results = make([]db.Script, 0, len(records))
	for _, rec := range titles.title {
		results = append(results, rec)
	}
	var lastChapter = -1
	for _, rec := range records {
		if rec.ChapterNum != lastChapter {
			lastChapter = rec.ChapterNum
			var rec2 = rec
			rec2.VerseNum = 0
			rec2.UsfmStyle = `para.h`
			rec2.VerseStr = ``
			rec2.ScriptTexts = []string{titles.heading + " " + strconv.Itoa(rec.ChapterNum)}
			results = append(results, rec2)
		}
		if rec.UsfmStyle != `para.h` && !strings.HasPrefix(rec.UsfmStyle, `para.mt`) {
			results = append(results, rec)
		}
	}
	return results
}

func (p *USXParser) correctScriptNum(records []db.Script) []db.Script {
	var results = make([]db.Script, 0, len(records))
	var scriptNum = 0
	var lastChapter = 0
	for _, rec := range records {
		if rec.ChapterNum != lastChapter {
			lastChapter = rec.ChapterNum
			scriptNum = 0
		}
		scriptNum += 1
		rec.ScriptNum = strconv.Itoa(scriptNum)
		results = append(results, rec)
	}
	return results
}

func (p *USXParser) findAttr(se xml.StartElement, name string) string {
	for _, attr := range se.Attr {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ``
}

func (p *USXParser) findIntAttr(se xml.StartElement, name string) int {
	val := p.findAttr(se, name)
	if val == `` {
		return 0
	} else {
		num := numericPattern.FindString(val)
		result, err := strconv.Atoi(num)
		if err != nil {
			panic(err)
		}
		return result
	}
}

func (s Stack) Push(v string) Stack {
	return append(s, v)
}

func (s Stack) Pop() (Stack, string) {
	l := len(s)
	if l < 1 {
		panic("You tried to pop an empty stack.")
	}
	return s[:l-1], s[l-1]
}
