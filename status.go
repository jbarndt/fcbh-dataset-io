package dataset

type Status struct {
	IsErr   bool
	Message string
	Status  int
	Err     string
	Trace   string
	Request string
}
