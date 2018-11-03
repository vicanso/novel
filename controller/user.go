package controller

import (
	"encoding/base64"
	"time"

	"github.com/vicanso/novel/context"
	"github.com/vicanso/novel/cs"
	"github.com/vicanso/novel/model"
	"github.com/vicanso/novel/xerror"
	"go.uber.org/zap"

	"github.com/vicanso/novel/validate"

	"github.com/labstack/echo"
	"github.com/vicanso/cookies"
	"github.com/vicanso/novel/config"
	"github.com/vicanso/novel/middleware"
	"github.com/vicanso/novel/router"
	"github.com/vicanso/novel/util"
)

const (
	avatar           = "/9j/4AAQSkZJRgABAQAAAQABAAD/2wCEAAgGBgcGBQgHBwcJCQgKDBQNDAsLDBkSEw8UHRofHh0aHBwgJC4nICIsIxwcKDcpLDAxNDQ0Hyc5PTgyPC4zNDIBCQkJDAsMGA0NGDIhHCEyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMjIyMv/AABEIACgAKAMBIgACEQEDEQH/xAB6AAACAwEBAQAAAAAAAAAAAAAABQYHCAEDBBAAAQMDAgQEBAcAAAAAAAAAAQIDBAAFEQYSISJBYQcTMVEUMnGxIzNicoHB0QEAAwEAAAAAAAAAAAAAAAAAAgMEAREAAgICAwADAAAAAAAAAAAAAAECAxEhEiIxE0FR/9oADAMBAAIRAxEAPwCWeJmrLSlgwG47UmSAfxicbPoev2qiU6umxpPMlDjOeKMYOOxp3am5GrtViM8vCOZZBGcJGDX33TS2n0XNyPCbefdRkrSknaD7DNS2Tgn2Q+FUpeM43Lakxm5LaipDidyRXUSXFrwAT2FPtB6FN1tE15b5Zaj7lJj7cqBIJAz7cPvVh+E8eHM0FbbmuHHE17zPMcDYySHFAfTgBQxp5b+gZSxplUt2C+zm0Ki22Y6TxUQ2oAdhj+6K0pRTfhX6BzMv6Qix4N2aujE9tDjK1MyY76wjelXBOwn1PbtU7LUOFKVKjx0MEZJKUFWSew4npSN9iA3oidDWW0OoCXW1hsFReSQUegySTw/mlls1M1NcZZ+IdiXAgIcYeTyqV2z6Z9qmthJ9kVU2R8eiytBebcGbj5L/AMMd4CtqQN+QehzivLwoj3A6EaajzwyGJL7WwtpOCHD271Vt7kXKPdnmUzJsdBCCgNoAQvIGSDkD5sil7Ey8QElEa4T4wySQEZGT6nkUftVNeoLJPdubaNLLjX4fJcY5/e1/lFZkXr/VEFxSTcJKkpOApQUAe/MKKMXxY60wRdNalh5eUpQpLbJzzrXyAg9Nu8qz+nhxpR4gqLd+iSIzm1ZYD6ACARucWoYx7DaB7AAelFFbFdTH6TV6/wAbWHhjc5TDLEe+QG0OPoSMbkpUCVJHcbvpkjrUHtVyM5lQV+Y3jd3HQ0UUNiXE2LeRiEg9KKKKnyNP/9k="
	cookieSigned     = true
	loginTokenKey    = "token"
	loginTokenLength = 8
)

var (
	cookieOptions = &cookies.Options{
		Keys:     config.GetSessionKeys(),
		MaxAge:   365 * 24 * 3600,
		Path:     config.GetCookiePath(),
		HttpOnly: true,
	}
	errNoTrackKey    = xerror.New("track key is not found")
	errLoginTokenNil = xerror.New("login token can not be nil")
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
		Email    string `valid:"email,optional"`
	}
	// UserUpdateRolesParams the params of user update roles
	UserUpdateRolesParams struct {
		Role string `valid:"in(admin)"`
		Type string `valid:"in(add|remove)"`
	}
)

func init() {
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

	users := router.NewGroup("/users", userSession)

	ctrl := UserCtrl{}

	users.Add("GET", "/v1/me/avatar", ctrl.getAvatar)
	users.Add("GET", "/v1/me", ctrl.getInfo)

	users.Add(
		"PATCH",
		"/v1/me",
		ctrl.refresh,
		middleware.IsLogined,
	)

	// user register
	users.Add(
		"POST",
		"/v1/me",
		ctrl.register,
		createTracker(cs.ActionRegister),
		waitForOneSecond,
		middleware.IsAnonymous,
	)

	// get user login token
	users.Add(
		"GET",
		"/v1/me/login",
		ctrl.getLoginToken,
		middleware.IsAnonymous,
	)

	// user login
	users.Add(
		"POST",
		"/v1/me/login",
		ctrl.doLogin,
		createTracker(cs.ActionLogin),
		waitForOneSecond,
		loginLimit,
		middleware.IsAnonymous,
	)

	// user logout
	users.Add(
		"DELETE",
		"/v1/me/logout",
		ctrl.doLogout,
		createTracker(cs.ActionLogout),
	)
}

// getAvatar get user avatar
func (uc *UserCtrl) getAvatar(c echo.Context) (err error) {
	// 图像数据手工生成，不会出错，忽略出错处理
	buf, _ := base64.StdEncoding.DecodeString(avatar)
	setPrivateCache(c, "24h")
	setContentType(c, "image/jpeg")
	res(c, buf)
	return
}

// get user info from session
func (uc *UserCtrl) pickUserInfo(c echo.Context) (userInfo *UserInfoRes) {
	us := getUserSession(c)
	userInfo = &UserInfoRes{
		Anonymous: true,
		Date:      now(),
		IP:        c.RealIP(),
		TrackID:   getTrackID(c),
	}
	if us == nil {
		return
	}
	account := us.GetAccount()
	if account != "" {
		userInfo.Account = account
		userInfo.Anonymous = false
		userInfo.UpdatedAt = us.GetUpdatedAt()
		userInfo.LoginedAt = us.GetLoginedAt()
		userInfo.Roles = us.GetRoles()
	}
	return
}

// getUserInfo get user info ctrl
func (uc *UserCtrl) getInfo(c echo.Context) (err error) {
	key := config.GetTrackKey()
	ck := cookies.New(c.Request(), c.Response(), cookieOptions)
	if ck.Get(key, cookieSigned) == "" {
		cookie := ck.CreateCookie(key, util.GenUlid())
		ck.Set(cookie, cookieSigned)
	}
	userInfo := uc.pickUserInfo(c)
	res(c, userInfo)
	return
}

// register register user ctrl
func (uc *UserCtrl) register(c echo.Context) (err error) {
	buf := getRequestBody(c)
	params := &UserRegisterParams{}
	err = validate.Do(params, buf)
	if err != nil {
		return
	}
	user, err := userService.Register(params.Account, params.Password, params.Email)
	if err != nil {
		return
	}
	resCreated(c, user)
	return
}

// getLoginToken get the token for login function
func (uc *UserCtrl) getLoginToken(c echo.Context) (err error) {
	key := config.GetTrackKey()
	ck := cookies.New(c.Request(), c.Response(), cookieOptions)
	if ck.Get(key, cookieSigned) == "" {
		err = errNoTrackKey
		return
	}
	token := util.RandomString(loginTokenLength)
	us := getUserSession(c)
	us.SetLoginToken(token)
	us.RefreshSessionCookie()
	res(c, map[string]string{
		"token": token,
	})
	return
}

// doLogin user login
func (uc *UserCtrl) doLogin(c echo.Context) (err error) {
	buf := getRequestBody(c)
	params := &UserLoginParams{}
	err = validate.Do(params, buf)
	if err != nil {
		return
	}
	us := getUserSession(c)
	token := us.GetLoginToken()
	if token == "" {
		err = errLoginTokenNil
		return
	}
	user, err := userService.Login(params.Account, params.Password, token)
	if err != nil {
		return
	}
	roles := []string{}
	for _, v := range user.Roles {
		roles = append(roles, v)
	}
	us.SetAccount(user.Account)
	us.SetLoginedAt(now())
	us.SetRoles(roles)
	us.SetLoginToken("")
	userInfo := uc.pickUserInfo(c)

	cookie, _ := c.Cookie(config.GetSessionCookie())
	sessionID := ""
	if cookie != nil {
		sessionID = cookie.Value
	}
	ul := &model.UserLogin{
		Account:   user.Account,
		UserAgent: c.Request().UserAgent(),
		IP:        c.RealIP(),
		TrackID:   getTrackID(c),
		SessionID: sessionID,
	}

	e := userService.AddLoginRecord(ul)
	if e != nil {
		getContextLogger(c).Error("save user login record fail",
			zap.String("account", user.Account),
			zap.Error(err),
		)
	}
	res(c, userInfo)
	return
}

// doLogout user logout
func (uc *UserCtrl) doLogout(c echo.Context) (err error) {
	// 删除sig的cookie，cookie则校验不过需要重新生成
	cookie, _ := c.Cookie(config.GetSessionCookie() + ".sig")
	if cookie != nil {
		cookie.Value = ""
		cookie.Expires = time.Now()
		cookie.Path = config.GetCookiePath()
		c.SetCookie(cookie)
	}
	context.SetUserSession(c, nil)
	userInfo := uc.pickUserInfo(c)
	res(c, userInfo)
	return
}

// refresh fresh user session
func (uc *UserCtrl) refresh(c echo.Context) (err error) {
	us := getUserSession(c)
	err = us.Refresh()
	return
}
