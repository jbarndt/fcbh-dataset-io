package timestamp

import (
	"context"
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"encoding/json"
	"fmt"
	"os"
)

type Waha struct {
	AudioFile string      `json:"audio_file"`
	TextFile  string      `json:"text_file"`
	Sections  []WahaVerse `json:"sections"`
}

type WahaVerse struct {
	VerseId    string    `json:"verse_id"`
	Timings    []float64 `json:"timings"`
	TimingsStr []string  `json:"timings_str"`
	Text       string    `json:"text"`
	Uroman     string    `json:"uroman_tokens"`
}

type WahaTimestamper struct {
	ctx context.Context
}

func NewWahaTimestamper(ctx context.Context) WahaTimestamper {
	var w WahaTimestamper
	w.ctx = ctx
	return w
}

func (w *WahaTimestamper) GetTimestamps(tsType string, mediaId string, bookId string, chapterNum int) ([]db.Audio, dataset.Status) {
	var result []db.Audio
	var status dataset.Status
	var filename = os.Getenv("HOME") + "/miniconda3/envs/torch1/waha-ai-timestamper-cli/MRK_1.json"
	bytes, err := os.ReadFile(filename)
	if err != nil {
		status = log.Error(w.ctx, 500, err, "Error reading timestamps file")
		return result, status
	}
	var response []Waha
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		status = log.Error(w.ctx, 500, err, "Error parsing timestamps json file")
		return result, status
	}
	fmt.Println("json", response)
	var chapter = response[0]
	for _, seg := range chapter.Sections {
		var aud db.Audio
		aud.Book = bookId
		aud.ChapterNum = chapterNum
		aud.AudioChapter = chapter.AudioFile
		aud.VerseStr = seg.VerseId
		aud.Text = seg.Text
		fmt.Println(aud.VerseStr, aud.Text)
		aud.Uroman = seg.Uroman
		if len(seg.Timings) > 1 {
			aud.BeginTS = seg.Timings[0]
			aud.EndTS = seg.Timings[1]
		} else {
			status = log.ErrorNoErr(w.ctx, 500, "Missing Timestamps for "+seg.VerseId)
			return result, status
		}
		result = append(result, aud)
	}
	return result, status
}

/*
{
	"audio_file": "MRK.1.mp3",
	"text_file": "MRK.1.txt",
	"sections":
} [
{"verse_id": "MRK.1.1",
	"timings": [0.7, 4.1],
	"timings_str": ["00:00:00", "00:00:04"],
	"text": "The Good News According to Mark",
	"uroman_tokens": "t h e g o o d n e w s a c c o r d i n g t o m a r k"},
{"verse_id": "MRK.1.2",
	"timings": [4.1, 9.26], "timings_str":
		["00:00:04", "00:00:09"], "text": "The beginning of the Good News of Jesus Christ, the Son of God.", "uroman_tokens": "t h e b e g i n n i n g o f t h e g o o d n e w s o f j e s u s c h r i s t t h e s o n o f g o d"},
*/
