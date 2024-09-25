package database

import "gorm.io/gorm"

type Paging struct {
	Page     int `json:"page" query:"page" param:"page"`
	PageSize int `json:"page_size" query:"page_size" param:"page_size"`
}

func DBPaging(db *gorm.DB, paging *Paging) *gorm.DB {
	if paging.Page < 0 {
		paging.Page = 1
	}

	if paging.PageSize < 1 {
		paging.PageSize = 10
	}

	return db.Limit(paging.PageSize).Offset((paging.Page - 1) * paging.PageSize)
}
