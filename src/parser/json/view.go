package json

type Envelop struct {
	Records []RecordView `json:"records"`
	Errors  []ErrorView  `json:"errors"`
}

type RecordView struct {
	Date            string        `json:"date"`
	Summary         string        `json:"summary"`
	Total           string        `json:"total"`
	TotalMins       int           `json:"total_mins"`
	ShouldTotal     string        `json:"should_total"`
	ShouldTotalMins int           `json:"should_total_mins"`
	Diff            string        `json:"diff"`
	DiffMins        int           `json:"diff_mins"`
	Tags            []string      `json:"tags"`
	Entries         []interface{} `json:"entries"`
}

type EntryView struct {
	Type      string   `json:"type"`
	Summary   string   `json:"summary"`
	Tags      []string `json:"tags"`
	Total     string   `json:"total"`
	TotalMins int      `json:"total_mins"`
}

type OpenRangeView struct {
	EntryView
	Start     string `json:"start"`
	StartMins int    `json:"start_mins"`
}

type RangeView struct {
	OpenRangeView
	End     string `json:"end"`
	EndMins int    `json:"end_mins"`
}

type ErrorView struct {
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Length  int    `json:"length"`
	Title   string `json:"title"`
	Details string `json:"details"`
}
