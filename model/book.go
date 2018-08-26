package model

import (
	"github.com/lib/pq"

	"github.com/vicanso/novel-spider/novel"
)

type (
	// Book book's info
	Book struct {
		BaseModel `json:"base_model,omitempty"`
		//
		Name string `json:"name,omitempty" gorm:"type:varchar(20);not null;unique_index:idx_books_name_author"`
		// Author 作者，需要建索引
		Author string `json:"author,omitempty" gorm:"type:varchar(100);not null;unique_index:idx_books_name_author"`
		Brief  string `json:"brief,omitempty"`
		Cover  string `json:"cover,omitempty"`
		// Category 书籍分类
		Category    pq.StringArray `json:"category,omitempty" gorm:"type:text[]"`
		Source      string         `json:"source,omitempty" gorm:"index:book_source_source_id"`
		SourceID    int            `json:"sourceId,omitempty" gorm:"index:book_source_source_id"`
		SourceCover string         `json:"sourceCoveer,omitempty"`
		Status      int            `json:"status,omitempty"`
	}
)

// AddBook add book
func AddBook(info *novel.BasicInfo) (err error) {
	b := Book{}
	attr := Book{
		Brief: info.Brief,

		// Cover 在成功后保存？
		Category: []string{
			info.Category,
		},
		Source:      info.Source,
		SourceID:    info.SourceID,
		SourceCover: info.Cover,
		Status:      StatusPadding,
	}
	err = db.Where(Book{
		Name:   info.Name,
		Author: info.Author,
	}).
		Attrs(&attr).
		FirstOrCreate(&b).Error
	return
}

// IsExistsSource check the source exists
func IsExistsSource(source string, id int) (exitst bool, err error) {
	b := Book{}

	err = db.Where(&Book{
		Source:   source,
		SourceID: id,
	}).First(&b).Error
	if err != nil {
		return
	}
	if b.ID != 0 {
		exitst = true
	}
	return
}
