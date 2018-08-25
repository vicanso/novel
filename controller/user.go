package controller

import (
	"encoding/base64"
	"time"

	"github.com/vicanso/cookies"
	"github.com/vicanso/session"

	"github.com/kataras/iris"
	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/middleware"
	"github.com/vicanso/novel/router"
	"github.com/vicanso/novel/utils"
)

const (
	avatar       = "/9j/4AAQSkZJRgABAQAAAQABAAD/2wCEAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDIBCQkJDAsMGA0NGDIhHCEyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMv/AABEIACgAKAMBIgACEQEDEQH/xAB6AAACAwEBAQAAAAAAAAAAAAAABQYHCAEDBBAAAQMDAgQEBAcAAAAAAAAAAQIDBAAFEQYSISJBYQcTMVEUMnGxIzNicoHB0QEAAwEAAAAAAAAAAAAAAAAAAgMEAREAAgICAwADAAAAAAAAAAAAAAECAxEhEiIxE0FR/9oADAMBAAIRAxEAPwCWeJmrLSlgwG47UmSAfxicbPoev2qiU6umxpPMlDjOeKMYOOxp3am5GrtViM8vCOZZBGcJGDX33TS2n0XNyPCbefdRkrSknaD7DNS2Tgn2Q+FUpeM43Lakxm5LaipDidyRXUSXFrwAT2FPtB6FN1tE15b5Zaj7lJj7cqBIJAz7cPvVh+E8eHM0FbbmuHHE17zPMcDYySHFAfTgBQxp5b+gZSxplUt2C+zm0Ki22Y6TxUQ2oAdhj+6K0pRTfhX6BzMv6Qix4N2aujE9tDjK1MyY76wjelXBOwn1PbtU7LUOFKVKjx0MEZJKUFWSew4npSN9iA3oidDWW0OoCXW1hsFReSQUegySTw/mlls1M1NcZZ+IdiXAgIcYeTyqV2z6Z9qmthJ9kVU2R8eiytBebcGbj5L/AMMd4CtqQN+QehzivLwoj3A6EaajzwyGJL7WwtpOCHD271Vt7kXKPdnmUzJsdBCCgNoAQvIGSDkD5sil7Ey8QElEa4T4wySQEZGT6nkUftVNeoLJPdubaNLLjX4fJcY5/e1/lFZkXr/VEFxSTcJKkpOApQUAe/MKKMXxY60wRdNalh5eUpQpLbJzzrXyAg9Nu8qz+nhxpR4gqLd+iSIzm1ZYD6ACARucWoYx7DaB7AAelFFbFdTH6TV6/wAbWHhjc5TDLEe+QG0OPoSMbkpUCVJHcbvpkjrUHtVyM5lQV+Y3jd3HQ0UUNiXE2LeRiEg9KKKKnyNP/9k="
	cookieSigned = true
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
	userCtrl struct {
	}
	userInfoResponse struct {
		Anonymous bool   `json:"anonymous"`
		Account   string `json:"account,omitempty"`
		Date      string `json:"date"`
	}
	userLoginParams struct {
		Account  string `valid:"ascii,runelength(4|10)"`
		Password string `valid:"runelength(6|20)"`
	}
)

func init() {
	users := router.NewGroup("/users", router.SessionHandler)
	ctrl := userCtrl{}
	users.Add("GET", "/v1/me/token", ctrl.getLoginToken)
	users.Add("GET", "/v1/me/avatar", ctrl.getAvatar)
	users.Add("GET", "/v1/me", ctrl.getInfo)
	users.Add("POST", "/v1/me", newDefaultTracker("login", nil), ctrl.doLogin)
	users.Add("DELETE", "/v1/me", newDefaultTracker("logout", nil), ctrl.doLogout)
	users.Add("PATCH", "/v1/me", middleware.IsLogined, ctrl.refresh)
}

// 从session中筛选用户信息
func (c *userCtrl) pickUserInfo(sess *session.Session) (userInfo *userInfoResponse) {
	userInfo = &userInfoResponse{
		Anonymous: true,
		Date:      time.Now().Format(time.RFC3339),
	}
	account := sess.GetString(utils.AccountField)

	if account != "" {
		userInfo.Account = account
		userInfo.Anonymous = false
	}
	return
}

// getUserInfo 获取用户信息
func (c *userCtrl) getInfo(ctx iris.Context) {
	userInfo := c.pickUserInfo(utils.GetSession(ctx))
	key := config.GetTrackKey()
	ck := cookies.New(ctx.Request(), ctx.ResponseWriter(), cookieOptions)
	if ck.Get(key, cookieSigned) == "" {
		cookie := ck.CreateCookie(key, utils.GenUlid())
		ck.Set(cookie, cookieSigned)
	}
	res(ctx, userInfo)
}

// getAvatar 获取用户头像
func (c *userCtrl) getAvatar(ctx iris.Context) {
	// 图像数据手工生成，不会出错，忽略出错处理
	buf, _ := base64.StdEncoding.DecodeString(avatar)
	resJPEG(ctx, buf)
}

// getLoginToken 获取登录加密使用的token
func (c *userCtrl) getLoginToken(ctx iris.Context) {
	token := utils.RandomString(8)
	sess := utils.GetSession(ctx)
	sess.Set("token", token)
	setNoStore(ctx)
	res(ctx, iris.Map{
		"token": token,
	})
}

// doLogin 登录
func (c *userCtrl) doLogin(ctx iris.Context) {
	body := utils.GetRequestBody(ctx)
	params := &userLoginParams{}
	err := validate(params, body)
	if err != nil {
		resErr(ctx, err)
		return
	}
	key := config.GetTrackKey()
	ck := cookies.New(ctx.Request(), ctx.ResponseWriter(), cookieOptions)
	// track cookie用于跟踪用户，必须保证是正确的才可以登录
	if ck.Get(key, cookieSigned) == "" {
		resErr(ctx, utils.ErrNoTrackKey)
		return
	}
	account := utils.GetAccount(ctx)
	if account != "" {
		resErr(ctx, utils.ErrLoginedAlready)
		return
	}
	// TODO 密码的校验需要增加login token来生成
	sess := utils.GetSession(ctx)
	err = sess.Set(utils.AccountField, params.Account)
	if err != nil {
		resErr(ctx, err)
		return
	}
	userInfo := c.pickUserInfo(sess)
	res(ctx, userInfo)
}

// doLogout 退出登录
func (c *userCtrl) doLogout(ctx iris.Context) {
	ctx.RemoveCookie(config.GetSessionCookie())
	sess := utils.GetSession(ctx)
	sess.Set(utils.AccountField, "")
	userInfo := c.pickUserInfo(sess)
	res(ctx, userInfo)
}

// refresh 刷新
func (c *userCtrl) refresh(ctx iris.Context) {
	sess := utils.GetSession(ctx)
	err := sess.Refresh()
	if err != nil {
		resErr(ctx, err)
		return
	}
	resNoContent(ctx)
}
