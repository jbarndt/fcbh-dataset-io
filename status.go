package dataset

type Status struct {
	IsErr   bool
	Message string
	Status  int
	Err     string
	Trace   string
	Request string
}

// Status implements the Error interface
func (e *Status) Error() string {
	return e.Err
}
