package mtv

type day struct {
	ID    int    `json:"Id"`
	Title string `json:"Title"`
}

type daysResponse []day

type program struct {
	ID            int    `json:"Id"`
	Name          string `json:"Name"`
	ArName        string `json:"ArName"`
	Description   string `json:"Description"`
	ArDescription string `json:"ArDescription"`
}

type scheduleEntry struct {
	ID      int     `json:"Id"`
	Day     int     `json:"Day"`
	Time    string  `json:"Time"`
	Program program `json:"Program"`
	IsRerun *bool   `json:"IsRerun"`
}

type scheduleResponse []scheduleEntry
