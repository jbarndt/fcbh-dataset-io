package fa_score_analysis

import (
	"dataset"
	"dataset/db"
	log "dataset/logger"
	"gonum.org/v1/gonum/stat"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func FAScoreAnalysis(conn db.DBAdapter) (string, dataset.Status) {
	var status dataset.Status
	var project = filepath.Base(conn.Database)
	var ext = filepath.Ext(project)
	project = project[:len(project)-len(ext)]
	var filename = project + "_fa_error.txt"
	file, err := os.Create(filename)
	if err != nil {
		return filename, log.Error(conn.Ctx, 500, err, "unable to create fa_error.txt")
	}
	defer file.Close()
	write(file, project, " FA Error Analysis\n")

	chars, status := conn.SelectFACharTimestamps()
	if status.IsErr {
		return filename, status
	}
	var faErrors []float64
	for _, char := range chars {
		faErrors = append(faErrors, -math.Log10(char.FAScore))
	}
	mean := stat.Mean(faErrors, nil)
	write(file, "Mean: ", strconv.FormatFloat(mean, 'f', 2, 64))
	stdDev := stat.StdDev(faErrors, nil)
	write(file, "StdDev: ", strconv.FormatFloat(stdDev, 'f', 2, 64))

	var mini = math.Inf(1)
	var maxi = 0.0
	for _, er := range faErrors {
		if er < mini {
			mini = er
		}
		if er > maxi {
			maxi = er
		}
	}
	write(file, "Minimum: ", strconv.FormatFloat(mini, 'f', 2, 64))
	write(file, "Maximum: ", strconv.FormatFloat(maxi, 'f', 2, 64))
	// Skewness (asymmetry of distribution)
	skewness := stat.Skew(faErrors, nil)
	write(file, "Skewness: ", strconv.FormatFloat(skewness, 'f', 2, 64))
	// Kurtosis (shape of distribution)
	kurtosis := stat.ExKurtosis(faErrors, nil)
	write(file, "Kurtosis: ", strconv.FormatFloat(kurtosis, 'f', 2, 64))
	// Percentile
	write(file, "\nPercentiles")
	sort.Float64s(faErrors)
	for _, percent := range []float64{0.95, 0.96, 0.97, 0.98, 0.99, 0.995, 0.996, 0.997, 0.998, 0.999} {
		percentile := stat.Quantile(percent, stat.Empirical, faErrors, nil)
		percentStr := strconv.FormatFloat((percent * 100.0), 'f', 1, 64)
		write(file, "Percentile ", percentStr, ": ", strconv.FormatFloat(percentile, 'f', 2, 64))
	}
	// Histogram
	write(file, "\nHISTOGRAM")
	var histogram = make(map[int]int)
	for _, er := range faErrors {
		histogram[int(er)]++
	}
	var keys []int
	for k := range histogram {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	numFAError := len(faErrors)
	for _, cat := range keys {
		pct := float64(histogram[cat]) / float64(numFAError) * 100.0
		write(file, "Cat: ", strconv.Itoa(cat), "-", strconv.Itoa(cat+1), " = ", strconv.FormatFloat(pct, 'f', 4, 64))
	}
	return filename, status
}

func write(file *os.File, args ...string) {
	for _, arg := range args {
		_, _ = file.WriteString(arg)
	}
	_, _ = file.WriteString("\n")
}
