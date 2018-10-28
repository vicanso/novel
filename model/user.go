package model

import (
	"github.com/lib/pq"
)

const (
	// UserRoleSu super user
	UserRoleSu = "su"
	// UserRoleAdmin admin user
	UserRoleAdmin = "admin"
)

type (
	// User user model
	User struct {
		BaseModel
		Account  string         `json:"account,omitempty" gorm:"type:varchar(20);not null;unique_index:idx_users_account"`
		Password string         `json:"password,omitempty" gorm:"type:varchar(128);not null;"`
		Email    string         `json:"email,omitempty" gorm:"type:varchar(128);"`
		Roles    pq.StringArray `json:"roles,omitempty" gorm:"type:text[]"`
	}
	// UserLogin user login
	UserLogin struct {
		BaseModel
		Account   string `json:"account,omitempty" gorm:"type:varchar(20);not null;index:idx_user_logins_account"`
		UserAgent string `json:"userAgent,omitempty"`
		IP        string `json:"ip,omitempty" gorm:"type:varchar(64);not null"`
		TrackID   string `json:"trackId,omitempty" gorm:"type:varchar(64);not null"`
		SessionID string `json:"sessionId,omitempty" gorm:"type:varchar(64);not null"`
	}
)
