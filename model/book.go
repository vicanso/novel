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
		Category        pq.StringArray `json:"category,omitempty" gorm:"type:text[]"`
		Source          string         `json:"source,omitempty" gorm:"not null;index:idx_books_source_source_id"`
		SourceID        int            `json:"sourceId,omitempty" gorm:"not null;index:idx_books_source_source_id"`
		SourceCover     string         `json:"sourceCover,omitempty"`
		Status          int            `json:"status,omitempty" gorm:"not null;index:idx_books_status;"`
		WordCount       int            `json:"wordCount,omitempty"`
		ViewCount       int            `json:"viewCount,omitempty" gorm:"index:idx_books_view_count"`
		LatestViewCount int            `json:"latestViewCount,omitempty" gorm:"index:idx_books_latest_view_count"`
		LikeCount       int            `json:"likeCount,omitempty" gorm:"index:idx_books_like_count"`
		LatestLikeCount int            `json:"latestLikeCount,omitempty" gorm:"index:idx_books_latest_like_count"`
	}
	// Chapter chapter struct
	Chapter struct {
		BaseModel
		Index     int    `json:"index,omitempty" gorm:"not null;unique_index:idx_chapters_book_id_index"`
		Title     string `json:"title,omitempty" gorm:"not null;"`
		Content   string `json:"content,omitempty" gorm:"not null;"`
		WordCount int    `json:"wordCount,omitempty"`
		BookID    uint   `json:"bookId,omitempty" gorm:"not null;unique_index:idx_chapters_book_id_index"`
	}
	// Favorite user's favorite
	Favorite struct {
		BaseModel
		BookID          uint   `json:"bookId,omitempty" gorm:"not null;unique_index:idx_favorites_book_id_account"`
		Account         string `json:"account,omitempty" gorm:"not null;unique_index:idx_favorites_book_id_account"`
		ReaddingChapter int    `json:"readdingChapter"`
	}
)
