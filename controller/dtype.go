package controller

import "time"

type Request struct {
	SessionID               []string  `json:"session_id"`
	Category                []int     `json:"category"`
	Store                   []int     `json:"store"`
	StoreBrand              []int     `json:"store_brand"`
	StoreChannel            []int     `json:"store_channel"`
	PhotoTakenBy            []int     `json:"photo_taken_by"`
	VisitedFrom             time.Time `json:"visited_from"`
	VisitedTo               time.Time `json:"visited_to"`
	SessionProcessingStatus string    `json:"session_processing_status"`
	EvidenceProgressStatus  string    `json:"evidence_progress_status"`
	QualityProcessionStatus string    `json:"quality_processing_status"`
}

type Download struct {
	ID         int64  `json:"id"`
	ReportName string `json:"report_name"`
	Status     string `json:"status"`
	Created    string `json:"created"`
	Modified   string `json:"modified"`
}
type Store struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type StoreBrand struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type StoreChannel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type PhotoType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PhotoSession struct {
	ID                      string    `json:"session_id"`
	VisitedOn               string    `json:"visited_on"`
	CreatedAt               string    `json:"created_at"`
	Category                Category  `json:"category"`
	Store                   Store     `json:"store"`
	PhotoTakenBy            User      `json:"photo_taken_by"`
	PhotoType               PhotoType `json:"photo_type"`
	PhotoCount              int       `json:"photo_count"`
	SessionProcessingStatus string    `json:"session_processing_status"`
	EvidenceProgressStatus  string    `json:"evidence_progress_status"`
	QualityProcessionStatus string    `json:"quality_processing_status"`
}
