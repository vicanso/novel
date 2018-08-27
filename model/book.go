package model

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"

	"github.com/vicanso/novel-spider/novel"
)

type (
	// Book book's info
	Book struct {
		BaseModel
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

// AddBook add book
func AddBook(info *novel.BasicInfo) (err error) {
	client := GetClient()
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
	err = client.Where(Book{
		Name:   info.Name,
		Author: info.Author,
	}).
		Attrs(&attr).
		FirstOrCreate(&b).Error
	return
}

// IsExistsSource check the source exists
func IsExistsSource(source string, id int) (exitst bool, err error) {
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

func getClientByConditions(params *BookQueryConditions, client *gorm.DB) *gorm.DB {
	if client == nil {
		client = GetClient()
	}
	client = enhanceWhere(client, "name", params.Name)
	client = enhanceWhere(client, "author", params.Author)
	client = enhanceWhere(client, "status", params.Status)
	if params.Category != "" {
		client = client.Where("? = ANY (category)", params.Category)
	}
	client = enhanceRangeWhere(client, "updated_at", params.UpdatedAt)
	client = enhanceRangeWhere(client, "created_at", params.CreatedAt)
	return client
}

// ListBook list book
func ListBook(params *BookQueryConditions, options *QueryOptions) (books []*Book, err error) {
	client := getClientByOptions(options, nil)
	client = getClientByConditions(params, client)
	books = make([]*Book, 0)
	err = client.Find(&books).Error
	return
}

// CountBook count the book
func CountBook(params *BookQueryConditions) (count int, err error) {
	client := GetClient()
	client = getClientByConditions(params, client.Model(&Book{}))
	err = client.Count(&count).Error
	return
}

// UpdateBookByID update book by id
func UpdateBookByID(id uint, data *Book) (err error) {
	client := GetClient()
	b := &Book{}
	b.ID = id
	err = client.Model(b).Updates(data).Error
	return
}
