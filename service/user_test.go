package service

import (
	"testing"

	"github.com/vicanso/novel/model"
	"github.com/vicanso/novel/util"
	"github.com/vicanso/session"
)

func TestUserSession(t *testing.T) {
	account := "vicanso"
	sess := session.Mock(session.M{
		"data": session.M{
			UserAccount: account,
		},
		"fetched": true,
	})
	us := NewUserSession(sess)
	if us.GetAccount() != account {
		t.Fatalf("get account fail")
	}
	account = "abcd"
	us.SetAccount(account)
	if us.GetAccount() != account {
		t.Fatalf("set account fail")
	}

	if us.GetUpdatedAt() == "" {
		t.Fatalf("get updated at fail")
	}

	us.SetLoginedAt(util.Now())
	if us.GetLoginedAt() == "" {
		t.Fatalf("get logined at fail")
	}

	us.SetRoles([]string{
		"su",
	})
	if us.GetRoles()[0] != "su" {
		t.Fatalf("get user roles fail")
	}

	us.SetLoginToken("a")
	if us.GetLoginToken() != "a" {
		t.Fatalf("get login token fail")
	}

	err := us.Refresh()
	if err != nil {
		t.Fatalf("refresh fail, %v", err)
	}
}

func TestUserService(t *testing.T) {
	userService := User{}
	account := util.RandomString(8)
	pwd := util.RandomString(8)
	token := "abcd"

	t.Run("register", func(t *testing.T) {
		_, err := userService.Register(account, pwd)
		if err != nil {
			t.Fatalf("register fail, %v", err)
		}
		_, err = userService.Register(account, pwd)
		if err != errUserAccountExists {
			t.Fatalf("the account exists should return error")
		}
	})

	t.Run("login", func(t *testing.T) {
		_, err := userService.Login(account, pwd, token)
		if err != errAccountOrPasswordWrong {
			t.Fatalf("the password is wrong should return error")
		}

		hash := util.Sha1(token + pwd)
		_, err = userService.Login(account, hash, token)
		if err != nil {
			t.Fatalf("login fail, %v", err)
		}
	})

	t.Run("add login record", func(t *testing.T) {
		err := userService.AddLoginRecord(&model.UserLogin{
			Account: account,
		})
		if err != nil {
			t.Fatalf("add login record fail, %v", err)
		}
	})
}
