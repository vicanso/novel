package model

import (
	"github.com/lib/pq"
)

const (
	// BookStatusPadding 待审核
	BookStatusPadding = iota
	// BookStatusRejected 拒绝
	BookStatusRejected
	// BookStatusPassed 已通过
	BookStatusPassed
)

type (
	// Book book's info
	Book struct {
		BaseModel
		// Name book name
		Name string `json:"name,omitempty" gorm:"type:varchar(60);not null;unique_index:idx_books_name_author"`
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
	// Chapter chapter struct
	Chapter struct {
		BaseModel `json:"base_model,omitempty"`
		Index     int    `json:"index,omitempty" gorm:"not null;unique_index:idx_chapters_book_id_index"`
		Title     string `json:"title,omitempty" gorm:"not null;"`
		Content   string `json:"content,omitempty" gorm:"not null;"`
		WordCount int    `json:"wordCount,omitempty"`
		Book      Book   `gorm:"ForeignKey:BookID" json:"book,omitempty"`
		BookID    uint   `json:"bookId,omitempty" gorm:"not null;unique_index:idx_chapters_book_id_index"`
	}
	// BookQueryConditions book query conditions
	BookQueryConditions struct {
		Name      string
		Author    string
		Status    string
		Category  string
		UpdatedAt string
		CreatedAt string
	}
)

// IsExistsBook check the book is exists
func IsExistsBook(source string, id int) (exitst bool, err error) {
	client := GetClient()
	b := Book{}

	err = client.Where(&Book{
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

// AddBook add the book
func AddBook(b *Book) (err error) {
	err = GetClient().Create(b).Error
	return
}

// AddBookChapter add chapter for book
func AddBookChapter(name, author string, chapter *Chapter) (err error) {
	b := Book{}
	err = GetClient().Where(&Book{
		Name:   name,
		Author: author,
	}).Select("id").First(&b).Error
	if err != nil {
		return
	}
	if b.ID == 0 {
		return
	}
	chapter.BookID = b.ID
	err = GetClient().Create(chapter).Error
	return
}
