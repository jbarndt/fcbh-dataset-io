package match

import (
	"context"
	"dataset"
	"dataset/generic"
	log "dataset/logger"
	"math"
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
	critErrors  int
	questErrors int
	critGaps    int
	questGaps   int
}

func NewAlignWriter(ctx context.Context) AlignWriter {
	var a AlignWriter
	a.ctx = ctx
	return a
}

func (a *AlignWriter) WriteReport(datasetName string, lines []generic.AlignLine, filenameMap string) (string, dataset.Status) {
	var filename string
	var status dataset.Status
	var err error
	a.datasetName = datasetName
	a.out, err = os.CreateTemp(os.Getenv(`FCBH_DATASET_TMP`), datasetName+"_*.html")
	if err != nil {
		return filename, log.Error(a.ctx, 500, err, `Error creating output file for align writer`)
	}
	a.WriteHeading()
	for _, line := range lines {
		a.WriteLine(line.Chars)
	}
	a.WriteEnd(filenameMap)
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
	_, _ = a.out.WriteString("<audio id='validateAudio' controls></audio>\n")
	_, _ = a.out.WriteString(`<h2 style="text-align:center">Audio to Text Alignment Report For `)
	_, _ = a.out.WriteString(a.datasetName)
	_, _ = a.out.WriteString("</h2>\n")
	_, _ = a.out.WriteString(`<h3 style="text-align:center">`)
	loc, _ := time.LoadLocation("America/Denver")
	_, _ = a.out.WriteString(time.Now().In(loc).Format(`Mon Jan 2 2006 03:04:05 pm MST`))
	_, _ = a.out.WriteString("</h3>\n")
	checkbox := `<div style="text-align: center; margin: 10px;">
		<input type="checkbox" id="hideVerse0" checked><label for="hideVerse0">Hide Headings</label>
	</div>
`
	_, _ = a.out.WriteString(checkbox)
	table := `<table id="diffTable" class="display">
    <thead>
    <tr>
        <th>Line</th>
        <th>Critical<br>Error</th>
		<th>Silence<br>Long</th>
		<th>Start<br>TS</th>
		<th>Button</th>
        <th>Ref</th>
		<th>Script</th>
    </tr>
    </thead>
    <tbody>
`
	_, _ = a.out.WriteString(table)
}

func (a *AlignWriter) WriteLine(chars []generic.AlignChar) {
	var logTotal float64
	var asrChars int
	for i, char := range chars {
		if chars[i].FAScore <= criticalThreshold {
			chars[i].ScoreError = int(scoreCritical)
			logScore := -math.Log10(chars[i].FAScore)
			logTotal += logScore
		} else if chars[i].FAScore <= questionThreshold {
			chars[i].ScoreError = int(scoreQuestion)
		}
		if char.IsASR {
			asrChars++
		}
	}
	var firstChar = chars[0]
	var lastChar = chars[len(chars)-1]
	a.lineNum++
	_, _ = a.out.WriteString("<tr>\n")
	a.writeCell(strconv.Itoa(a.lineNum))
	a.writeCell(strconv.FormatFloat(logTotal, 'f', 2, 64))
	//a.writeCell(strings.Join(errors, "<br>"))
	//a.writeCell(strconv.FormatFloat(maxSilence*12.0, 'f', 0, 64))
	a.writeCell(strconv.FormatInt(int64(asrChars), 10))
	a.writeCell(a.minSecFormat(firstChar.BeginTS))
	var lineRef generic.LineRef
	ref := lineRef.ParseKey(firstChar.LineRef).(generic.LineRef)
	var params []string
	params = append(params, "'"+ref.BookId+"'")
	params = append(params, strconv.Itoa(ref.ChapterNum))
	params = append(params, strconv.FormatFloat(firstChar.BeginTS, 'f', 4, 64))
	params = append(params, strconv.FormatFloat(lastChar.EndTS, 'f', 4, 64))
	a.writeCell("<button onclick=\"playVerse(" + strings.Join(params, ",") + ")\">Play</button>")
	//a.writeCell(firstChar.BookId + ` ` + strconv.Itoa(firstChar.ChapterNum) + `:` + firstChar.VerseStr)
	a.writeCell(firstChar.LineRef)
	var text []string
	for _, ch := range chars {
		char := string(ch.CharNorm)
		if ch.SilenceLong != 0 {
			char += `<sub>` + strconv.Itoa(ch.SilencePos) + `</sub>`
		}
		if ch.CharSeq == 0 {
			text = append(text, " ")
		}
		if ch.ScoreError == int(scoreCritical) {
			text = append(text, `<span class="red-box">`+char+`</span>`)
		} else if ch.ScoreError == int(scoreQuestion) {
			text = append(text, `<span class="yellow-box">`+char+`</span>`)
		} else if ch.SilenceLong > 0 {
			text = append(text, `<span class="green-box">`+char+`</span>`)
		} else if ch.IsASR {
			text = append(text, `<span class="blue-box">`+char+`</span>`)
		} else {
			text = append(text, char)
		}
	}
	a.writeCell(strings.Join(text, ""))
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

func (a *AlignWriter) WriteEnd(filenameMap string) {
	table := `</tbody>
	</table>
`
	_, _ = a.out.WriteString(table)
	_, _ = a.out.WriteString(`<p>Lines with critical errors `)
	_, _ = a.out.WriteString(strconv.Itoa(a.critErrors))
	_, _ = a.out.WriteString(`</p>`)
	_, _ = a.out.WriteString(`<p>Lines with questionable errors `)
	_, _ = a.out.WriteString(strconv.Itoa(a.questErrors))
	_, _ = a.out.WriteString("</p>\n")
	_, _ = a.out.WriteString(`<p>Lines with large end-of-verse gaps `)
	_, _ = a.out.WriteString(strconv.Itoa(a.critGaps))
	_, _ = a.out.WriteString(`</p>`)
	_, _ = a.out.WriteString(`<p>Lines with smaller end-of-verse gaps `)
	_, _ = a.out.WriteString(strconv.Itoa(a.questGaps))
	_, _ = a.out.WriteString("</p>\n")
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
		padding: 1px 0;
	} 
	.yellow-box { 
		background-color: rgba(255, 255, 0, 0.8);
		padding: 1px 0;
	}
	.blue-box {
		background-color: rgba(0, 0, 255, 0.4);
		padding: 1px 0;
	}
	.green-box {
		background-color: rgba(0, 255, 0, 0.4);
		padding: 1px 0;
	}
	/*.blank-box {
		display: inline-block;
		height: 1em; 
		background-color: rgba(0, 0, 255, 0.4);
		padding: 2px 0;
		vertical-align: -4px;
	}*/
	</style>
`
	_, _ = a.out.WriteString(style)
	script := `<script>
    $(document).ready(function() {
        var table = $('#diffTable').DataTable({
            "columnDefs": [
                { "orderable": false, "targets": [3,4,5,6] }
				// { "visible": false, "targets": [8] }  
            ],
            "pageLength": 50,
            "lengthMenu": [[50, 500, -1], [50, 500, "All"]],
			"order": [[ 1, "desc" ]]
        });
    	$.fn.dataTable.ext.search.push(function(settings, data, dataIndex) {
        	var hideZeros = $('#hideVerse0').prop('checked');
        	if (!hideZeros) return true;
        	return !data[5].endsWith(":0"); 
    	});
    	$('#hideVerse0').prop('checked', true);
    	table.draw();
    	$('#hideVerse0').on('change', function() {
        	table.draw(); 
    	});
    });
	function playVerse(book, chapter, startTime, endTime) {
`
	_, _ = a.out.WriteString(script)
	_, _ = a.out.WriteString("\t" + filenameMap)
	script = `filename = fileMap[book + chapter]
		audioFile = './../../../../FCBH2024/download/ENGWEB/ENGWEBN2DA-mp3-64/' + filename
		console.log("audioFile", audioFile);
		const audio = document.getElementById('validateAudio');
		audio.src = audioFile;
		audio.controls = true;
		audio.playbackRate = 0.75;
		audio.currentTime = startTime;
		audio.play();
		audio.ontimeupdate = function() {
			if (audio.currentTime >= endTime) {
				audio.pause();
				audio.ontimeupdate = null;
			}
		}
	}
    </script>
</body>
</html>
`
	_, _ = a.out.WriteString("\t" + script)
	_ = a.out.Close()
}

func (a *AlignWriter) minSecFormat(duration float64) string {
	mins := int(duration / 60.0)
	secs := duration - float64(mins)*60.0
	var minStr string
	var delim string
	if int(mins) > 0 {
		minStr = strconv.FormatInt(int64(mins), 10)
		delim = ":"
	}
	secStr := strconv.FormatFloat(secs, 'f', 0, 64)
	return minStr + delim + secStr
}
