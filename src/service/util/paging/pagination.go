package paging

import (
	"github.com/jinzhu/gorm"
	"math"
)

type Param struct {
	DB      *gorm.DB
	Page    int
	Limit   int
	OrderBy []string
	ShowSQL bool
}

type Paginator struct {
	TotalRecord int `json:"total_record"`
	TotalPage   int `json:"total_page"`
	Offset      int `json:"offset"`
	Limit       int `json:"limit"`
	Page        int `json:"page"`
	PrevPage    int `json:"prev_page"`
	NextPage    int `json:"next_page"`
}

func Paging(p *Param, result interface{}) (*Paginator, error) {
	db := p.DB

	if p.ShowSQL {
		db = db.Debug()
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(o)
		}
	}

	dbDone := make(chan *gorm.DB, 1)
	var paginator Paginator
	var count int
	var offset int

	go countRecords(db, result, &count, dbDone)

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}

	dbQuery := db.Limit(p.Limit).Offset(offset).Find(result)
	var dbCount *gorm.DB
	dbCount = <-dbDone

	if dbCount.Error != nil {
		return nil, dbCount.Error
	}

	if dbQuery.Error != nil {
		return nil, dbQuery.Error
	}

	paginator.TotalRecord = count
	paginator.Page = p.Page

	paginator.Offset = offset
	paginator.Limit = p.Limit
	paginator.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))

	if p.Page > 1 {
		paginator.PrevPage = p.Page - 1
	} else {
		paginator.PrevPage = p.Page
	}

	if p.Page == paginator.TotalPage {
		paginator.NextPage = p.Page
	} else {
		paginator.NextPage = p.Page + 1
	}
	return &paginator, nil
}

func countRecords(db *gorm.DB, anyType interface{}, count *int, dbDone chan *gorm.DB) {
	dbDone <- db.Model(anyType).Count(count)
}
