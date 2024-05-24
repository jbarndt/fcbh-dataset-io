package match

import (
	"dataset"
	"dataset/db"
	"github.com/sergi/go-diff/diffmatchpatch"
	"os"
	"strconv"
	"strings"
)

func (c *Compare) CompareChapters() (string, dataset.Status) {
	var filename string
	var status dataset.Status
	c.baseDb, status = db.NewerDBAdapter(c.ctx, false, c.user.Username, c.baseDataset)
	if status.IsErr {
		return filename, status
	}
	output, status := c.openOutput(c.baseDataset, c.dataset)
	if status.IsErr {
		return filename, status
	}
	filename = output.Name()
	for _, bookId := range db.RequestedBooks(c.testament) {
		var chapInBook, _ = db.BookChapterMap[bookId]
		var chapter = 1
		for chapter <= chapInBook {
			var baseText, text string
			baseText, status = c.concatText(c.baseDb, bookId, chapter)
			if status.IsErr {
				return ``, status
			}
			baseText = c.cleanup(baseText)
			text, status = c.concatText(c.database, bookId, chapter)
			if status.IsErr {
				return ``, status
			}
			text = c.cleanup(text)
			c.chapterDiff(output, bookId, chapter, baseText, text)
			chapter++
		}
	}
	c.baseDb.Close()
	c.endReport(output)
	return filename, status
}

func (c *Compare) concatText(conn db.DBAdapter, bookId string, chapter int) (string, dataset.Status) {
	var results []string
	scripts, status := conn.SelectScriptsByChapter(bookId, chapter)
	if status.IsErr {
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

func (c *Compare) chapterDiff(output *os.File, bookId string, chapter int, baseText string, text string) {
	diffMatch := diffmatchpatch.New()
	diffs := diffMatch.DiffMain(baseText, text, false)
	if !c.isMatch(diffs) {
		inserts, deletes := c.measure(diffs)
		c.insertSum += inserts
		c.deleteSum += deletes
		c.writer.WriteChapterDiff(bookId, chapter, inserts, deletes, diffMatch.DiffPrettyHtml(diffs))
		_, _ = output.WriteString(`<h3 style="padding-left:50px;">`)
		_, _ = output.WriteString(bookId)
		_, _ = output.WriteString(" ")
		_, _ = output.WriteString(strconv.Itoa(chapter))
		_, _ = output.WriteString(", Inserts: ")
		_, _ = output.WriteString(strconv.Itoa(inserts))
		_, _ = output.WriteString(", Deletes: ")
		_, _ = output.WriteString(strconv.Itoa(deletes))
		_, _ = output.WriteString("</h3>\n")
		_, _ = output.WriteString(`<p>`)
		_, _ = output.WriteString(diffMatch.DiffPrettyHtml(diffs))
		_, _ = output.WriteString("</p>\n")
		c.diffCount++
	}
}
