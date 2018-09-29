package controller

import (
	"encoding/base64"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/vicanso/cookies"

	"github.com/kataras/iris"
	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/cs"
	"github.com/vicanso/novel/middleware"
	"github.com/vicanso/novel/model"
	"github.com/vicanso/novel/router"
	"github.com/vicanso/novel/util"
)

const (
	avatar        = "/9j/4AAQSkZJRgABAQAAAQABAAD/2wCEAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDIBCQkJDAsMGA0NGDIhHCEyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMv/AABEIACgAKAMBIgACEQEDEQH/xAB6AAACAwEBAQAAAAAAAAAAAAAABQYHCAEDBBAAAQMDAgQEBAcAAAAAAAAAAQIDBAAFEQYSISJBYQcTMVEUMnGxIzNicoHB0QEAAwEAAAAAAAAAAAAAAAAAAgMEAREAAgICAwADAAAAAAAAAAAAAAECAxEhEiIxE0FR/9oADAMBAAIRAxEAPwCWeJmrLSlgwG47UmSAfxicbPoev2qiU6umxpPMlDjOeKMYOOxp3am5GrtViM8vCOZZBGcJGDX33TS2n0XNyPCbefdRkrSknaD7DNS2Tgn2Q+FUpeM43Lakxm5LaipDidyRXUSXFrwAT2FPtB6FN1tE15b5Zaj7lJj7cqBIJAz7cPvVh+E8eHM0FbbmuHHE17zPMcDYySHFAfTgBQxp5b+gZSxplUt2C+zm0Ki22Y6TxUQ2oAdhj+6K0pRTfhX6BzMv6Qix4N2aujE9tDjK1MyY76wjelXBOwn1PbtU7LUOFKVKjx0MEZJKUFWSew4npSN9iA3oidDWW0OoCXW1hsFReSQUegySTw/mlls1M1NcZZ+IdiXAgIcYeTyqV2z6Z9qmthJ9kVU2R8eiytBebcGbj5L/AMMd4CtqQN+QehzivLwoj3A6EaajzwyGJL7WwtpOCHD271Vt7kXKPdnmUzJsdBCCgNoAQvIGSDkD5sil7Ey8QElEa4T4wySQEZGT6nkUftVNeoLJPdubaNLLjX4fJcY5/e1/lFZkXr/VEFxSTcJKkpOApQUAe/MKKMXxY60wRdNalh5eUpQpLbJzzrXyAg9Nu8qz+nhxpR4gqLd+iSIzm1ZYD6ACARucWoYx7DaB7AAelFFbFdTH6TV6/wAbWHhjc5TDLEe+QG0OPoSMbkpUCVJHcbvpkjrUHtVyM5lQV+Y3jd3HQ0UUNiXE2LeRiEg9KKKKnyNP/9k="
	cookieSigned  = true
	loginTokenKey = "token"
)

var (
	cookieOptions = &cookies.Options{
		Keys:     config.GetSessionKeys(),
		MaxAge:   365 * 24 * 3600,
		Path:     config.GetCookiePath(),
		HttpOnly: true,
	}
)

type (
	// UserCtrl user controller
	UserCtrl struct{}
	// UserInfoRes the response of user info
	UserInfoRes struct {
		Anonymous bool     `json:"anonymous,omitempty"`
		Account   string   `json:"account,omitempty"`
		Date      string   `json:"date,omitempty"`
		UpdatedAt string   `json:"updatedAt,omitempty"`
		IP        string   `json:"ip,omitempty"`
		TrackID   string   `json:"trackId,omitempty"`
		LoginedAt string   `json:"loginedAt,omitempty"`
		Roles     []string `json:"roles,omitempty"`
	}
	// UserLoginParams the params of user login
	UserLoginParams struct {
		Account  string `valid:"ascii,runelength(4|10)"`
		Password string `valid:"runelength(6|64)"`
	}
	// UserRegisterParams the params of user register
	UserRegisterParams struct {
		Account  string `valid:"ascii,runelength(4|10)"`
		Password string `valid:"runelength(6|64)"`
	}
	// UserUpdateRolesParams the params of user update roles
	UserUpdateRolesParams struct {
		Role string `valid:"in(admin)"`
		Type string `valid:"in(add|remove)"`
	}
)

func init() {
	users := router.NewGroup("/users", middleware.Session)
	waitForOneSecond := middleware.WaitFor(time.Second)

	// 增加登录的限制，可以用于限制账号的登录频率
	loginLimit := middleware.NewConcurrentLimiter(middleware.ConcurrentLimiterConfig{
		Category: cs.ActionLogin,
		Keys: []string{
			"account",
		},
		// 限制3秒只能登录一次
		TTL: 3 * time.Second,
	})
	ctrl := UserCtrl{}
	users.Add("GET", "/v1/me/avatar", ctrl.getAvatar)
	users.Add("GET", "/v1/me", ctrl.getInfo)
	users.Add("PATCH", "/v1/me", middleware.IsLogined, ctrl.refresh)
	// user register
	users.Add(
		"POST",
		"/v1/me",
		waitForOneSecond,
		newDefaultTracker(cs.ActionRegister, nil),
		middleware.IsAnonymous,
		ctrl.register,
	)
	// get user login token
	users.Add(
		"GET",
		"/v1/me/login",
		middleware.IsAnonymous,
		ctrl.getLoginToken,
	)
	// user login
	users.Add(
		"POST",
		"/v1/me/login",
		waitForOneSecond,
		newDefaultTracker(cs.ActionLogin, nil),
		loginLimit,
		middleware.IsAnonymous,
		ctrl.doLogin,
	)
	// user logout
	users.Add(
		"DELETE",
		"/v1/me/logout",
		newDefaultTracker(cs.ActionLogout, nil),
		ctrl.doLogout,
	)
	// user update roles
	users.Add(
		"PATCH",
		"/v1/roles/:account",
		newDefaultTracker(cs.ActionUpdateRoles, nil),
		middleware.IsSu,
		ctrl.updateRoles,
	)
}

func addUserLoginRecord(account string, ctx iris.Context) {
	ul := &model.UserLogin{
		Account:   account,
		UserAgent: ctx.Request().UserAgent(),
		IP:        ctx.RemoteAddr(),
		TrackID:   util.GetTrackID(ctx),
		SessionID: ctx.GetCookie(config.GetSessionCookie()),
	}
	err := ul.Save()
	if err != nil {
		getContextLogger(ctx).Error("save user login record fail",
			zap.String("account", ul.Account),
			zap.Error(err),
		)
	}
}

// 从session中筛选用户信息
func (c *UserCtrl) pickUserInfo(ctx iris.Context) (userInfo *UserInfoRes) {
	sess := getSession(ctx)
	userInfo = &UserInfoRes{
		Anonymous: true,
		Date:      getNow(),
		IP:        ctx.RemoteAddr(),
		TrackID:   util.GetTrackID(ctx),
	}
	if sess == nil {
		return
	}
	account := sess.GetString(cs.SessionAccountField)
	if account != "" {
		userInfo.Account = account
		userInfo.Anonymous = false
		userInfo.UpdatedAt = sess.GetUpdatedAt()
		userInfo.LoginedAt = sess.GetString(cs.SessionLoginedAtField)
		userInfo.Roles = sess.GetStringSlice(cs.SessionRolesField)
	}
	return
}

// getUserInfo 获取用户信息
func (c *UserCtrl) getInfo(ctx iris.Context) {
	userInfo := c.pickUserInfo(ctx)
	key := config.GetTrackKey()
	ck := cookies.New(ctx.Request(), ctx.ResponseWriter(), cookieOptions)
	if ck.Get(key, cookieSigned) == "" {
		cookie := ck.CreateCookie(key, util.GenUlid())
		ck.Set(cookie, cookieSigned)
	}
	res(ctx, userInfo)
}

// getAvatar 获取用户头像
func (c *UserCtrl) getAvatar(ctx iris.Context) {
	// 图像数据手工生成，不会出错，忽略出错处理
	buf, _ := base64.StdEncoding.DecodeString(avatar)
	resJPEG(ctx, buf)
}

// register 注册
func (c *UserCtrl) register(ctx iris.Context) {
	params := &UserLoginParams{}
	err := validate(params, getRequestBody(ctx))
	if err != nil {
		resErr(ctx, err)
		return
	}
	u := model.User{
		Account:  params.Account,
		Password: params.Password,
	}
	err = u.Save()
	if err != nil {
		resErr(ctx, err)
		return
	}
	resCreated(ctx, nil)
}

// getLoginToken 获取登录加密使用的token
func (c *UserCtrl) getLoginToken(ctx iris.Context) {
	key := config.GetTrackKey()
	ck := cookies.New(ctx.Request(), ctx.ResponseWriter(), cookieOptions)
	// track cookie用于跟踪用户，必须保证是正确的才可以登录
	if ck.Get(key, cookieSigned) == "" {
		resErr(ctx, util.ErrNoTrackKey)
		return
	}
	token := util.RandomString(8)
	sess := getSession(ctx)
	sess.Set(loginTokenKey, token)
	setNoStore(ctx)
	res(ctx, iris.Map{
		loginTokenKey: token,
	})
}

// doLogin 登录
func (c *UserCtrl) doLogin(ctx iris.Context) {
	body := getRequestBody(ctx)
	params := &UserLoginParams{}
	err := validate(params, body)
	if err != nil {
		resErr(ctx, err)
		return
	}

	sess := getSession(ctx)

	token := sess.GetString(loginTokenKey)
	if token == "" {
		resErr(ctx, util.ErrLoginTokenNil)
		return
	}

	u := &model.User{
		Account: params.Account,
	}
	err = u.First()
	if err == gorm.ErrRecordNotFound || u.ID == 0 {
		resErr(ctx, util.ErrAccountOrPasswordWrong)
		return
	}

	pwd := util.Sha1(token + u.Password)
	if util.IsDevelopment() && params.Password == "tree.xie" {
		// 开发环境万能密码
		pwd = params.Password
	}
	if params.Password != pwd {
		resErr(ctx, util.ErrAccountOrPasswordWrong)
		return
	}
	roles := []string{}
	for _, v := range u.Roles {
		roles = append(roles, v)
	}
	sess.SetMap(map[string]interface{}{
		cs.SessionAccountField:   u.Account,
		cs.SessionLoginedAtField: getNow(),
		cs.SessionRolesField:     roles,
		loginTokenKey:            nil,
	})
	userInfo := c.pickUserInfo(ctx)
	addUserLoginRecord(u.Account, ctx)
	res(ctx, userInfo)
}

// doLogout 退出登录
func (c *UserCtrl) doLogout(ctx iris.Context) {
	// 删除cookie会导致所有cookies清除(仅是此次ctx)
	ctx.RemoveCookie(config.GetSessionCookie())
	util.SetSession(ctx, nil)
	userInfo := c.pickUserInfo(ctx)
	res(ctx, userInfo)
}

// refresh 刷新
func (c *UserCtrl) refresh(ctx iris.Context) {
	sess := getSession(ctx)
	sess.Refresh()
	resNoContent(ctx)
}

// updateRoles update user's roles
func (c *UserCtrl) updateRoles(ctx iris.Context) {
	body := util.GetRequestBody(ctx)
	params := &UserUpdateRolesParams{}
	err := validate(params, body)
	if err != nil {
		resErr(ctx, err)
		return
	}
	account := ctx.Params().Get("account")

	u := model.User{
		Account: account,
	}
	action := model.UserActionAddRole
	if params.Type != "add" {
		action = model.UserActionRemoveRole
	}
	err = u.UpdateRole(params.Role, action)
	if err != nil {
		resErr(ctx, err)
		return
	}
	resNoContent(ctx)
}
