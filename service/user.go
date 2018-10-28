package service

import (
	"github.com/jinzhu/gorm"
	"github.com/vicanso/novel/model"
	"github.com/vicanso/novel/util"
	"github.com/vicanso/novel/xerror"
	"github.com/vicanso/session"
)

const (
	// UserAccount user account field
	UserAccount = "account"
	// UserLoginedAt user logined at
	UserLoginedAt = "loginedAt"
	// UserRoles user roles
	UserRoles = "roles"
	// UserLoginToken user login token
	UserLoginToken = "loginToken"
)

var (
	errUserAccountExists      = xerror.NewUser("account already exists")
	errAccountOrPasswordWrong = xerror.NewUser("account or password is wrong")
)

type (
	// UserSession user session struct
	UserSession struct {
		Sess *session.Session
	}
	// User user
	User struct{}
)

// GetAccount get the account
func (u *UserSession) GetAccount() string {
	if u.Sess == nil {
		return ""
	}
	return u.Sess.GetString(UserAccount)
}

// SetAccount set the account
func (u *UserSession) SetAccount(account string) error {
	return u.Sess.Set(UserAccount, account)
}

// GetUpdatedAt get updated at
func (u *UserSession) GetUpdatedAt() string {
	return u.Sess.GetUpdatedAt()
}

// SetLoginedAt set the logined at
func (u *UserSession) SetLoginedAt(date string) error {
	return u.Sess.Set(UserLoginedAt, date)
}

// GetLoginedAt get logined at
func (u *UserSession) GetLoginedAt() string {
	return u.Sess.GetString(UserLoginedAt)
}

// SetRoles set the roles
func (u *UserSession) SetRoles(roles []string) error {
	return u.Sess.Set(UserRoles, roles)
}

// GetRoles get user roles
func (u *UserSession) GetRoles() []string {
	return u.Sess.GetStringSlice(UserRoles)
}

// SetLoginToken get user login token
func (u *UserSession) SetLoginToken(token string) error {
	return u.Sess.Set(UserLoginToken, token)
}

// GetLoginToken get user login token
func (u *UserSession) GetLoginToken() string {
	return u.Sess.GetString(UserLoginToken)
}

// Destroy destroy user session
func (u *UserSession) Destroy() error {
	return u.Sess.Destroy()
}

// Refresh refresh user sesion
func (u *UserSession) Refresh() error {
	return u.Sess.Refresh()
}

// NewUserSession create a new user session
func NewUserSession(sess *session.Session) *UserSession {
	return &UserSession{
		Sess: sess,
	}
}

// Register register
func (u *User) Register(account, password, email string) (user *model.User, err error) {
	user = &model.User{}
	err = getClient().Where(&model.User{
		Account: account,
	}).First(user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	err = nil
	if user.ID != 0 {
		err = errUserAccountExists
		return
	}
	user = &model.User{
		Account:  account,
		Password: password,
		Email:    email,
	}
	err = getClient().Create(user).Error
	// TODO 对于第一个注册用户增加su权限
	if user.ID == 1 {
		go getClient().Model(user).Update(&model.User{
			Roles: []string{
				model.UserRoleSu,
			},
		})
	}
	user.Password = ""
	return
}

// Login user login
func (u *User) Login(account, password, token string) (user *model.User, err error) {
	user = &model.User{}
	err = getClient().Where(&model.User{
		Account: account,
	}).First(user).Error
	if err == gorm.ErrRecordNotFound || user.ID == 0 {
		err = errAccountOrPasswordWrong
		return
	}
	pwd := util.Sha256(token + user.Password)
	if util.IsDevelopment() && password == "tree.xie" {
		// 开发环境万能密码
		pwd = password
	}
	if password != pwd {
		err = errAccountOrPasswordWrong
		return
	}
	user.Password = ""
	return
}

// AddLoginRecord add login record
func (u *User) AddLoginRecord(ul *model.UserLogin) (err error) {
	err = getClient().Create(ul).Error
	return
}
