package output

import "math"

func NormalizeMFCC(scripts []Script, numMFCC int) []Script {
	for col := 0; col < numMFCC; col++ {
		var sum float64
		var count float64
		for _, scr := range scripts {
			for _, mf := range scr.MFCC {
				value := mf[:][col]
				sum += float64(value)
				count++
			}
		}
		var mean = sum / count
		var devSqr float64
		for _, scr := range scripts {
			for _, mf := range scr.MFCC {
				value := mf[:][col]
				devSqr += math.Pow(float64(value)-mean, 2)
			}
		}
		var stddev = math.Sqrt(devSqr / count)
		for i, scr := range scripts {
			for j, mf := range scr.MFCC {
				value := float64(mf[:][col])
				scripts[i].MFCC[j][:][col] = float32((value - mean) / stddev)
			}
		}
	}
	return scripts
}

func PadRows(scripts []Script, numMFCC int) []Script {
	largest := 0
	for _, scr := range scripts {
		if scr.MFCCRows > largest {
			largest = scr.MFCCRows
		}
	}
	var padRow = make([]float32, numMFCC)
	for _, scr := range scripts {
		needRows := largest - scr.MFCCRows
		for i := 0; i < needRows; i++ {
			scr.MFCC = append(scr.MFCC, padRow)
		}
	}
	return scripts
}
