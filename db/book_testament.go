package db

import "dataset"

var BookOT = []string{`GEN`, `EXO`, `LEV`, `NUM`, `DEU`, `JOS`, `JDG`, `RUT`, `1SA`, `2SA`, `1KI`, `2KI`,
	`1CH`, `2CH`, `EZR`, `NEH`, `EST`, `JOB`, `PSA`, `PRO`, `ECC`, `SNG`, `ISA`, `JER`, `LAM`, `EZK`, `DAN`,
	`HOS`, `JOL`, `AMO`, `OBA`, `JON`, `MIC`, `NAM`, `HAB`, `ZEP`, `HAG`, `ZEC`, `MAL`}
var BookNT = []string{`MAT`, `MRK`, `LUK`, `JHN`, `ACT`, `ROM`, `1CO`, `2CO`, `GAL`, `EPH`, `PHP`, `COL`,
	`1TH`, `2TH`, `1TI`, `2TI`, `TIT`, `PHM`, `HEB`, `JAS`, `1PE`, `2PE`, `1JN`, `2JN`, `3JN`, `JUD`, `REV`}

func BookCodes(testament dataset.TestamentType) []string {
	var result []string
	switch testament {
	case dataset.OT:
		result = BookOT
	case dataset.NT:
		result = BookNT
	case dataset.C:
		result = append(result, BookOT...)
		result = append(result, BookNT...)
	}
	return result
}
