package xerror

import (
	"errors"
	"net/http"
	"testing"
)

func TestHTTPError(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		message := "abc"
		err := New(message)
		he := err.(*HTTPError)
		if he.Message != message ||
			he.Category != ErrCategoryCommon ||
			he.StatusCode != http.StatusBadRequest {
			t.Fatalf("new a http error fail")
		}
	})

	t.Run("new io", func(t *testing.T) {
		err := NewIO("abc")
		he := err.(*HTTPError)
		if he.Category != ErrCategoryIO {
			t.Fatalf("new an io error fail")
		}
	})

	t.Run("new session", func(t *testing.T) {
		err := NewSession("abc")
		he := err.(*HTTPError)
		if he.Category != ErrCategorySession ||
			!he.Exception ||
			he.StatusCode != http.StatusInternalServerError {
			t.Fatalf("new a session error fail")
		}
	})

	t.Run("new json", func(t *testing.T) {
		he := NewJSON("abc").(*HTTPError)
		if he.Category != ErrCategoryJSON {
			t.Fatalf("new a json error fail")
		}
	})

	t.Run("new validate", func(t *testing.T) {
		he := NewValidate("abc").(*HTTPError)
		if he.Category != ErrCategoryValidte {
			t.Fatalf("new a validate error fail")
		}
	})

	t.Run("new user", func(t *testing.T) {
		he := NewUser("abc").(*HTTPError)
		if he.Category != ErrCategoryUser {
			t.Fatalf("new an user error fail")
		}
	})

	t.Run("get status", func(t *testing.T) {
		if GetStatusCode(errors.New("abc")) != http.StatusInternalServerError {
			t.Fatalf("get status from normal error fail")
		}

		if GetStatusCode(NewValidate("abc")) != http.StatusBadRequest {
			t.Fatalf("get status from http error fail")
		}
	})

	t.Run("get messsage", func(t *testing.T) {
		message := "abc"
		if GetMessage(errors.New(message)) != message {
			t.Fatalf("get message from normal error fail")
		}

		if GetMessage(New(message)) != message {
			t.Fatalf("get message from http error fail")
		}
	})
}

func TestReadCloser(t *testing.T) {
	r := NewErrorReadCloser(errors.New("abc"))
	_, err := r.Read([]byte(""))
	if err == nil {
		t.Fatalf("read function should return error ")
	}
	err = r.Close()
	if err != nil {
		t.Fatalf("close fail, %v", err)
	}
}
