package controller

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/crazi-coder/report-service/core/utils/helpers"
)

type Request struct {
	SessionID               []string  `json:"session_id"`
	Category                []int     `json:"category_list"`
	Store                   []int     `json:"store_list"`
	StoreBrand              []int     `json:"store_brand_list"`
	StoreChannel            []int     `json:"store_channel_list"`
	PhotoTakenBy            []int     `json:"photo_taken_by_list"`
	PhotoType               []int     `json:"photo_type"`
	VisitedFrom             time.Time `json:"visited_from"`
	VisitedTo               time.Time `json:"visited_to"`
	SessionProcessingStatus string    `json:"session_processing_status"`
	EvidenceProgressStatus  string    `json:"evidence_progress_status"`
	QualityProcessionStatus string    `json:"quality_processing_status"`
	PageSize                uint      `json:"page_size"`
	PageNumber              uint      `json:"page_number"`
}

func (r *Request) SetPageSize(pageSize string) {
	i, err := strconv.ParseUint(pageSize, 10, 64)
	if err == nil {
		r.PageSize = uint(i)
	}
}

func (r *Request) SetPageNumber(pageNumber string) {
	i, err := strconv.ParseUint(pageNumber, 10, 64)
	if err == nil {
		r.PageNumber = uint(i)
	}
}

func (r *Request) StoreList(storeList string) {
	storeStrList := strings.Split(storeList, ",")
	for _, e := range storeStrList {
		i, err := strconv.ParseInt(e, 10, 64)
		if err == nil {
			r.Store = append(r.Store, int(i))
		}
	}
}

func (r *Request) StoreChannelList(storeChannelList string) {
	storeChannelStrList := strings.Split(storeChannelList, ",")
	for _, e := range storeChannelStrList {
		i, err := strconv.ParseInt(e, 10, 64)
		if err == nil {
			r.StoreChannel = append(r.StoreChannel, int(i))
		}
	}
}

func (r *Request) StoreBrandList(storeBrandList string) {
	storeBrandStrList := strings.Split(storeBrandList, ",")
	for _, e := range storeBrandStrList {
		i, err := strconv.ParseInt(e, 10, 64)
		if err == nil {
			r.StoreBrand = append(r.StoreBrand, int(i))
		}
	}
}

func (r *Request) CategoryList(categoryList string) {
	categoryStrList := strings.Split(categoryList, ",")
	for _, e := range categoryStrList {
		i, err := strconv.ParseInt(e, 10, 64)
		if err == nil {
			r.Category = append(r.Category, int(i))
		}
	}
}

func (r *Request) SessionIDList(storeChannelList string) {
	r.SessionID = strings.Split(storeChannelList, ",")
}

func (r *Request) PhotoTypeList(photoTypeStr string) {
	photoTypeStrList := strings.Split(photoTypeStr, ",")
	for _, e := range photoTypeStrList {
		i, err := strconv.ParseInt(e, 10, 64)
		if err == nil {
			r.PhotoType = append(r.PhotoType, int(i))
		}
	}

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
	Category                Category  `json:"Category"`
	Store                   Store     `json:"Store"`
	PhotoTakenBy            User      `json:"PhotoTakenBy"`
	PhotoType               PhotoType `json:"PhotoType"`
	PhotoCount              int       `json:"photo_count"`
	SessionProcessingStatus string    `json:"session_processing_status"`
	EvidenceProgressStatus  string    `json:"evidence_progress_status"`
	QualityProcessionStatus string    `json:"quality_processing_status"`
}

// Paginator is a  Generic Type used for pagination.
type Paginator struct {
	Next string `json:"next"`
	Prev string `json:"prev"`
}

// Pagination implementation for pagination.
func (p *Paginator) Pagination(requestURL string, requestedPageNumber uint,
	itemPerPage uint, totalItem uint) (*Paginator, error) {
	if (totalItem/itemPerPage+1 <= requestedPageNumber) {
		return nil, helpers.ErrPageLimitExceeded
	}
	var nextPageNumber uint
	var prevPageNumber uint
	if requestedPageNumber == 0 {
		requestedPageNumber = 1
	}
	if requestedPageNumber*itemPerPage < totalItem {
		nextPageNumber = requestedPageNumber + 1
	}
	if requestedPageNumber > 1 {
		prevPageNumber = requestedPageNumber - 1
	}

	urlNext, err := url.Parse(requestURL)

	if err != nil {
		return nil, err
	}
	urlPrev, err := url.Parse(requestURL)
	if err != nil {
		return nil, err
	}
	nextValues := urlNext.Query()
	prevValues := urlPrev.Query()
	nextValues.Set("page", strconv.Itoa(int(nextPageNumber)))
	prevValues.Set("page", strconv.Itoa(int(prevPageNumber)))
	urlNext.RawQuery = nextValues.Encode()
	urlPrev.RawQuery = prevValues.Encode()

	var next, prev string

	if nextPageNumber == 0 {
		next = ""
	} else {
		next = urlNext.String()
	}

	if prevPageNumber == 0 {
		prev = ""
	} else {
		prev = urlPrev.String()
	}

	paginator := Paginator{
		Next: next,
		Prev: prev,
	}

	return &paginator, nil
}

// PaginatedResult stores the search results after the search run.
type PaginatedResult struct {
	Count     uint         `json:"count"`
	Paginator Paginator   `json:"pagination"`
	Result    interface{} `json:"results"`
}
