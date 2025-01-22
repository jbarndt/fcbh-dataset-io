package timestamp

import (
	"context"
	"dataset"
	log "dataset/logger"
	"encoding/json"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"path/filepath"
	"strconv"
)

type ProbeData struct {
	Format ProbeFormat `json:"format"`
}

type ProbeFormat struct {
	Filename       string `json:"filename"`
	NBStreams      int    `json:"nb_streams"`
	NBProgress     int    `json:"nb_programs"`
	NBStreamGroups int    `json:"nb_stream_groups"`
	FormatName     string `json:"format_name"`
	FormatLongName string `json:"format_long_name"`
	StartTime      string `json:"start_time"`
	Duration       string `json:"duration"`
	Size           string `json:"size"`
	BitRate        string `json:"bit_rate"`
	ProbeScore     int    `json:"probe_score"`
}

func GetAudioDuration(ctx context.Context, directory string, filename string) (float64, dataset.Status) {
	var result float64
	probeData, status := GetProbeData(ctx, directory, filename)
	if status.IsErr {
		return result, status
	}
	var err error
	result, err = strconv.ParseFloat(probeData.Format.Duration, 64)
	if err != nil {
		status = log.Error(ctx, 500, err, "Data conversion error in timestamp.GetAudioDuration")
	}
	return result, status
}

func GetAudioSize(ctx context.Context, directory string, filename string) (float64, dataset.Status) {
	var result float64
	probeData, status := GetProbeData(ctx, directory, filename)
	if status.IsErr {
		return result, status
	}
	var err error
	result, err = strconv.ParseFloat(probeData.Format.Size, 64)
	if err != nil {
		status = log.Error(ctx, 500, err, "Data conversion error in timestamp.GetAudioSize")
	}
	return result, status
}

func GetAudioBitrate(ctx context.Context, directory string, filename string) (float64, dataset.Status) {
	var result float64
	var status dataset.Status
	probeData, status := GetProbeData(ctx, directory, filename)
	if status.IsErr {
		return result, status
	}
	var err error
	result, err = strconv.ParseFloat(probeData.Format.BitRate, 64)
	if err != nil {
		status = log.Error(ctx, 500, err, "Data conversion error in timestamp.GetAudioBitrate")
	}
	return result, status
}

func GetProbeData(ctx context.Context, directory string, filename string) (ProbeData, dataset.Status) {
	var result ProbeData
	var status dataset.Status
	filePath := filepath.Join(directory, filename)
	data, err := ffmpeg.Probe(filePath)
	if err != nil {
		return result, log.Error(ctx, 500, err, "Error in timestamp.GetProbeData")
	}
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		status = log.Error(ctx, 500, err, "Error in timestamp.GetProbeData")
	}
	return result, status
}
