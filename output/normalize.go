package output

import "math"

func NormalizeMFCC(structs []Script, numMFCC int) []HasMFCC {
	for col := 0; col < numMFCC; col++ {
		var sum float64
		var count float64
		for _, scr := range structs {
			for _, mf := range scr.GetMFCC() {
				sum += float64(mf[col])
				count++
			}
		}
		var mean = sum / count
		var devSqr float64
		for _, scr := range structs {
			for _, mf := range scr.GetMFCC() {
				devSqr += math.Pow(float64(mf[col])-mean, 2)
			}
		}
		var stddev = math.Sqrt(devSqr / count)
		for i, scr := range structs {
			var mfccs = scr.GetMFCC()
			for j, mf := range mfccs {
				value := float64(mf[col])
				mfccs[j][col] = float32((value - mean) / stddev)
			}
			structs[i].SetMFCC(mfccs)
		}
	}
	return structs
}

func PadRows(structs []HasMFCC, numMFCC int) []HasMFCC {
	largest := 0
	for _, scr := range structs {
		if scr.Rows() > largest {
			largest = scr.Rows()
		}
	}
	var padRow = make([]float32, numMFCC)
	for i, scr := range structs {
		mfccs := scr.GetMFCC()
		needRows := largest - scr.Rows()
		for i := 0; i < needRows; i++ {
			mfccs = append(mfccs, padRow)
		}
		structs[i].SetMFCC(mfccs)
	}
	return structs
}
