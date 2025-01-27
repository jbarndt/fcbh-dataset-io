package bible_brain

//type Boundary struct {
//	TimestampID string
//	Timestamp   float64
//	Duration    float64
//	Position    int64
//	NumBytes    int64
//}

type Segment struct {
	TimestampId int64
	VerseStr    string
	Timestamp   float64
	Duration    float64
	Position    int64
	NumBytes    int64
}
