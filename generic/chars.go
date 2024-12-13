package generic

type Char struct {
	CharId int64
	WordId int64
	Seq    int
	Norm   rune
	Uroman rune
	Token  int
	Start  float64
	End    float64
	Score  float64
}
