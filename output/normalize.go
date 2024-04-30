package output

import "math"

func NormalizeScriptMFCC(structs []Script, numMFCC int) []Script {
	for col := 0; col < numMFCC; col++ {
		var sum float64
		var count float64
		for _, str := range structs {
			for _, mf := range str.MFCC {
				sum += mf[col]
				count++
			}
		}
		var mean = sum / count
		var devSqr float64
		for _, str := range structs {
			for _, mf := range str.MFCC {
				devSqr += math.Pow(mf[col]-mean, 2)
			}
		}
		var stddev = math.Sqrt(devSqr / count)
		for i, str := range structs {
			var mfccs = str.MFCC
			for j, mf := range mfccs {
				value := mf[col]
				mfccs[j][col] = (value - mean) / stddev
			}
			structs[i].MFCC = mfccs
		}
	}
	return structs
}

func NormalizeWordMFCC(structs []Word, numMFCC int) []Word {
	for col := 0; col < numMFCC; col++ {
		var sum float64
		var count float64
		for _, str := range structs {
			for _, mf := range str.MFCC {
				sum += float64(mf[col])
				count++
			}
		}
		var mean = sum / count
		var devSqr float64
		for _, str := range structs {
			for _, mf := range str.MFCC {
				devSqr += math.Pow(mf[col]-mean, 2)
			}
		}
		var stddev = math.Sqrt(devSqr / count)
		for i, str := range structs {
			var mfccs = str.MFCC
			for j, mf := range mfccs {
				value := mf[col]
				mfccs[j][col] = (value - mean) / stddev
			}
			structs[i].MFCC = mfccs
		}
	}
	return structs
}

func PadScriptRows(structs []Script, numMFCC int) []Script {
	largest := 0
	for _, str := range structs {
		if str.MFCCRows > largest {
			largest = str.MFCCRows
		}
	}
	var padRow = make([]float64, numMFCC)
	for i, str := range structs {
		mfccs := str.MFCC
		needRows := largest - str.MFCCRows
		for i := 0; i < needRows; i++ {
			mfccs = append(mfccs, padRow)
		}
		structs[i].MFCC = mfccs
	}
	return structs
}

func PadWordRows(structs []Word, numMFCC int) []Word {
	largest := 0
	for _, str := range structs {
		if str.MFCCRows > largest {
			largest = str.MFCCRows
		}
	}
	var padRow = make([]float64, numMFCC)
	for i, str := range structs {
		mfccs := str.MFCC
		needRows := largest - str.MFCCRows
		for i := 0; i < needRows; i++ {
			mfccs = append(mfccs, padRow)
		}
		structs[i].MFCC = mfccs
	}
	return structs
}
