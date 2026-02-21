package json

// Envelop is the top level data structure of the JSON output.
// It contains three nodes:
// - `records`: is `null` if there are errors
// - `warnings`: only if applicable and only unless there are errors
// - `errors`: only unless there are records
type Envelop struct {
	Records  []RecordView `json:"records"`
	Warnings []string     `json:"warnings"`
	Errors   []ErrorView  `json:"errors"`
}

// RecordView is the JSON representation of a record.
// It also contains some evaluation data, such as the total time.
type RecordView struct {
	Date            string   `json:"date"`
	Summary         string   `json:"summary"`
	Total           string   `json:"total"`
	TotalMins       int      `json:"total_mins"`
	ShouldTotal     string   `json:"should_total"`
	ShouldTotalMins int      `json:"should_total_mins"`
	Diff            string   `json:"diff"`
	DiffMins        int      `json:"diff_mins"`
	Tags            []string `json:"tags"`
	Entries         []any    `json:"entries"`
}

// EntryView is the JSON representation of an entry.
type EntryView struct {
	// Type is one of `range`, `duration`, or `open_range`.
	Type    string `json:"type"`
	Summary string `json:"summary"`

	// Tags is a list of all tags that the entry summary contains.
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

// ErrorView is the JSON representation of a parsing error.
type ErrorView struct {
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Length  int    `json:"length"`
	Title   string `json:"title"`
	Details string `json:"details"`
	File    string `json:"file"`
}
