package match

import (
	"context"
	"dataset"
	log "dataset/logger"
	"os"
	"strconv"
	"strings"
	"time"
)

type AlignWriter struct {
	ctx         context.Context
	datasetName string
	out         *os.File
	lineNum     int
}

func NewAlignWriter(ctx context.Context) AlignWriter {
	var a AlignWriter
	a.ctx = ctx
	return a
}

func (a *AlignWriter) WriteReport(datasetName string, verses []FAverse) (string, dataset.Status) {
	var filename string
	var status dataset.Status
	var err error
	a.datasetName = datasetName
	a.out, err = os.CreateTemp(os.Getenv(`FCBH_DATASET_TMP`), datasetName+"_*.html")
	if err != nil {
		return filename, log.Error(a.ctx, 500, err, `Error creating output file for align writer`)
	}
	a.WriteHeading()
	for _, vers := range verses {
		a.WriteVerse(vers)
	}
	a.WriteEnd()
	return a.out.Name(), status
}

func (a *AlignWriter) WriteHeading() {
	head := `<!DOCTYPE html>
<html>
 <head>
  <meta charset="utf-8">
  <title>Alignment Error Report</title>
`
	_, _ = a.out.WriteString(head)
	_, _ = a.out.WriteString(`<link rel="stylesheet" type="text/css" href="https://cdn.datatables.net/1.10.21/css/jquery.dataTables.css">`)
	_, _ = a.out.WriteString("</head><body>\n")
	_, _ = a.out.WriteString(`<h2 style="text-align:center">Audio to Text Alignment Report For `)
	_, _ = a.out.WriteString(a.datasetName)
	_, _ = a.out.WriteString("</h2>\n")
	_, _ = a.out.WriteString(`<h3 style="text-align:center">`)
	_, _ = a.out.WriteString(time.Now().Format(`Mon Jan 2 2006 03:04:05 pm MST`))
	_, _ = a.out.WriteString("</h3>\n")
	//_, _ = h.out.WriteString(`<h3 style="text-align:center">RED characters are those in `)
	//_, _ = h.out.WriteString(baseDataset)
	//_, _ = h.out.WriteString(` only, while GREEN characters are in `)
	//_, _ = h.out.WriteString(h.datasetName)
	//_, _ = h.out.WriteString(" only</h3>\n")
	checkbox := `<div style="text-align: center; margin: 10px;">
		<input type="checkbox" id="hideVerse0" checked><label for="hideVerse0">Hide Headings</label>
	</div>
`
	_, _ = a.out.WriteString(checkbox)
	table := `<table id="diffTable" class="display">
    <thead>
    <tr>
        <th>Line</th>
        <th>Error</th>
		<th>End<br>Gap</th>
        <th>Ref</th>
		<th>Script</th>
    </tr>
    </thead>
    <tbody>
`
	_, _ = a.out.WriteString(table)
}

func (a *AlignWriter) WriteVerse(verse FAverse) {
	a.lineNum++
	_, _ = a.out.WriteString("<tr>\n")
	a.writeCell(strconv.Itoa(a.lineNum))
	//a.writeCell(strconv.FormatFloat(verse.critScore, 'f', 2, 64))
	a.writeCell(strconv.FormatFloat(verse.questScore, 'f', 1, 64))
	//a.writeCell(strconv.FormatFloat(verse.startTSDiff, 'f', 2, 64))
	a.writeCell(strconv.FormatFloat(verse.endTSDiff, 'f', 2, 64))
	a.writeCell(verse.bookId + ` ` + strconv.Itoa(verse.chapter) + `:` + verse.verseStr)
	var critical = a.createHighlightList(verse.critWords, len(verse.words))
	var question = a.createHighlightList(verse.questWords, len(verse.words))
	var text []string
	for i, wd := range verse.words {
		if critical[i] {
			text = append(text, `<span class="red-box">`+wd.Text+`</span>`)
		} else if question[i] {
			text = append(text, `<span class="yellow-box">`+wd.Text+`</span>`)
		} else {
			text = append(text, wd.Text)
		}
	}
	var diff int
	if verse.critDiff || verse.questDiff {
		diff = int((verse.endTSDiff - 0.9) * 10.0) // subract 1sec, and convert to 1/10 sec per char
		if diff > 0 {
			spaces := strings.Repeat("&nbsp;", diff)
			if verse.critDiff {
				text = append(text, `<span class="red-box">`+spaces+`</span>`)
			} else if verse.questDiff {
				text = append(text, `<span class="yellow-box">`+spaces+`</span>`)
			}
		}
	}
	a.writeCell(strings.Join(text, " "))
	_, _ = a.out.WriteString("</tr>\n")
}

func (a *AlignWriter) createHighlightList(indexes []int, length int) []bool {
	var list = make([]bool, length)
	for _, index := range indexes {
		list[index] = true
	}
	return list
}

func (a *AlignWriter) writeCell(content string) {
	_, _ = a.out.WriteString(`<td>`)
	_, _ = a.out.WriteString(content)
	_, _ = a.out.WriteString(`</td>`)
}

func (a *AlignWriter) WriteEnd() {
	table := `</tbody>
	</table>
`
	_, _ = a.out.WriteString(table)
	_, _ = a.out.WriteString(`<p>Lines with critical errors `)
	//_, _ = a.out.WriteString(strconv.Itoa(insertSum))
	_, _ = a.out.WriteString(`</p>`)
	_, _ = a.out.WriteString(`<p>Lines with questionable errors `)
	//_, _ = a.out.WriteString(strconv.Itoa(deleteSum))
	_, _ = a.out.WriteString("</p>\n")
	//_, _ = h.out.WriteString(`<p>`)
	//_, _ = h.out.WriteString("Total Difference Count: ")
	//_, _ = h.out.WriteString(strconv.Itoa(diffCount))
	//_, _ = h.out.WriteString("</p>\n")
	_, _ = a.out.WriteString(`<script type="text/javascript" src="https://code.jquery.com/jquery-3.5.1.js"></script>`)
	_, _ = a.out.WriteString("\n")
	_, _ = a.out.WriteString(`<script type="text/javascript" src="https://cdn.datatables.net/1.10.21/js/jquery.dataTables.js"></script>`)
	_, _ = a.out.WriteString("\n")
	style := `<style>
	.dataTables_length select {
	width: auto;
	display: inline-block;
	padding: 5px;
		margin-left: 5px;
		border-radius: 4px;
	border: 1px solid #ccc;
	}
	.dataTables_filter input {
		width: auto;
		display: inline-block;
		padding: 5px;
		border-radius: 4px;
		border: 1px solid #ccc;
	}
	.dataTables_wrapper .dataTables_length, .dataTables_wrapper .dataTables_filter {
		margin-bottom: 20px;
	}
	.red-box { 
		background-color: rgba(255, 0, 0, 0.4);
		padding: 0 3px; /* 0 top/bottom, 3px left/right */ 
		border-radius: 3px; /* rounded corners */ 
	} 
	.yellow-box { 
		background-color: rgba(255, 255, 0, 0.8);
		padding: 0 3px; 
		border-radius: 3px; 
	}
	</style>
`
	_, _ = a.out.WriteString(style)
	script := `<script>
    $(document).ready(function() {
        var table = $('#diffTable').DataTable({
            "columnDefs": [
                { "orderable": false, "targets": [3,4] }
				// { "visible": false, "targets": [8] }  
            ],
            "pageLength": 50,
            "lengthMenu": [[50, 500, -1], [50, 500, "All"]],
			"order": [[ 1, "desc" ]]
        });
    	$.fn.dataTable.ext.search.push(function(settings, data, dataIndex) {
        	var hideZeros = $('#hideVerse0').prop('checked');
        	if (!hideZeros) return true;
        	return !data[3].endsWith(":0"); 
    	});
    	$('#hideVerse0').prop('checked', true);
    	table.draw();
    	$('#hideVerse0').on('change', function() {
        	table.draw(); 
    	});
    });
</script>
</body>
</html>
`
	_, _ = a.out.WriteString(script)
	_ = a.out.Close()
}
