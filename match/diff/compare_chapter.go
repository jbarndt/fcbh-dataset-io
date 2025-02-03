package diff

import (
	"github.com/faithcomesbyhearing/fcbh-dataset-io/db"
	log "github.com/faithcomesbyhearing/fcbh-dataset-io/logger"
	"github.com/sergi/go-diff/diffmatchpatch"
	"strings"
)

func (c *Compare) CompareChapters() (string, *log.Status) {
	var filename string
	var status *log.Status
	filename = c.writer.WriteHeading(c.baseDataset)
	var scripts []db.Script
	scripts, status = c.database.SelectBookChapter()
	if status != nil {
		return filename, status
	}
	for _, scp := range scripts {
		var baseText, text string
		baseText, status = c.concatText(c.baseDb, scp.BookId, scp.ChapterNum)
		if status != nil {
			return ``, status
		}
		baseText = c.cleanup(baseText)
		text, status = c.concatText(c.database, scp.BookId, scp.ChapterNum)
		if status != nil {
			return ``, status
		}
		text = c.cleanup(text)
		c.chapterDiff(scp.BookId, scp.ChapterNum, baseText, text)
	}
	filenameMap, status := c.generateBookChapterFilenameMap()
	c.baseDb.Close()
	c.writer.WriteEnd(filenameMap, c.insertSum, c.deleteSum, c.diffCount)
	return filename, status
}

func (c *Compare) concatText(conn db.DBAdapter, bookId string, chapter int) (string, *log.Status) {
	var results []string
	scripts, status := conn.SelectScriptsByChapter(bookId, chapter)
	if status != nil {
		return ``, status
	}
	var priorText string
	for _, scp := range scripts {
		if !strings.HasSuffix(priorText, ` `) && !strings.HasPrefix(scp.ScriptText, ` `) {
			results = append(results, ` `)
		}
		results = append(results, scp.ScriptText)
		priorText = scp.ScriptText
	}
	return strings.Join(results, ""), status
}

func (c *Compare) chapterDiff(bookId string, chapter int, baseText string, text string) {
	diffMatch := diffmatchpatch.New()
	diffs := diffMatch.DiffMain(baseText, text, false)
	if !c.isMatch(diffs) {
		inserts, deletes := c.measure(diffs)
		c.insertSum += inserts
		c.deleteSum += deletes
		avgLen := float64(len(baseText)+len(text)) / 2.0
		errPct := float64((inserts+deletes)*100) / avgLen
		c.writer.WriteChapterDiff(bookId, chapter, inserts, deletes, errPct, diffMatch.DiffPrettyHtml(diffs))
		c.diffCount++
	}
}
