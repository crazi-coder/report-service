package controller

type Report struct {
}

type Download struct {
	ID         int64  `json:"id"`
	ReportName string `json:"report_name"`
	Status     string `json:"status"`
	Created    string `json:"created"`
	Modified   string `json:"modified"`
}
