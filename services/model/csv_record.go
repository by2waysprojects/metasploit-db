package services

// Struct to hold the data of each CSV row
type CSVRecord struct {
	Seq      string `csv:"seq"`
	Size     string `csv:"size"`
	Protocol string `csv:"protocol"`
	Request  string `csv:"request"`
	Body     string `csv:"body"`
}
