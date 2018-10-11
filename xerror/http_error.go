package xerror

import (
	"io"
	"net/http"
)

const (
	// ErrCategoryException exception error
	ErrCategoryException = "exception"
	// ErrCategoryCommon common error
	ErrCategoryCommon = "common"
	// ErrCategoryIO io error
	ErrCategoryIO = "io"
	// ErrCategoryPanic panic error
	ErrCategoryPanic = "panic"
	// ErrCategorySession session error
	ErrCategorySession = "session"
	// ErrCategoryJSON json error
	ErrCategoryJSON = "json"
	// ErrCategoryValidte validate error
	ErrCategoryValidte = "validate"
	// ErrCategoryRequset requset error
	ErrCategoryRequset = "request"
	// ErrCategoryUser user error
	ErrCategoryUser = "user"
)

type (
	// HTTPError http error
	HTTPError struct {
		StatusCode int                    `json:"statusCode,omitempty"`
		Exception  bool                   `json:"exception,omitempty"`
		Message    string                 `json:"message,omitempty"`
		Category   string                 `json:"category,omitempty"`
		Stack      []string               `json:"stack,omitempty"`
		Extra      map[string]interface{} `json:"extra,omitempty"`
	}
	errReadCloser struct {
		customErr error
	}
)

// Error makes it compatible with `error` interface.
func (he *HTTPError) Error() string {
	return he.Message
}

// New create a http error (common category)
func New(message string) error {
	he := &HTTPError{
		StatusCode: http.StatusBadRequest,
		Category:   ErrCategoryCommon,
		Message:    message,
	}
	return he
}

// NewIO create an io error (io categorry)
func NewIO(message string) error {
	he := &HTTPError{
		StatusCode: http.StatusBadRequest,
		Category:   ErrCategoryIO,
		Message:    message,
	}
	return he
}

// NewSession create a session error
func NewSession(message string) error {
	return &HTTPError{
		StatusCode: http.StatusInternalServerError,
		Exception:  true,
		Message:    message,
		Category:   ErrCategorySession,
	}
}

// NewJSON create a json error
func NewJSON(message string) error {
	return &HTTPError{
		StatusCode: http.StatusBadRequest,
		Exception:  false,
		Message:    message,
		Category:   ErrCategoryJSON,
	}
}

// NewValidate create a validate error
func NewValidate(message string) error {
	return &HTTPError{
		StatusCode: http.StatusBadRequest,
		Exception:  false,
		Message:    message,
		Category:   ErrCategoryValidte,
	}
}

// NewUser create a user error
func NewUser(message string) error {
	return &HTTPError{
		StatusCode: http.StatusBadRequest,
		Exception:  false,
		Message:    message,
		Category:   ErrCategoryUser,
	}
}

// GetStatusCode get status code
func GetStatusCode(err error) int {
	he, ok := err.(*HTTPError)
	if ok {
		return he.StatusCode
	}
	return http.StatusInternalServerError
}

// GetMessage get error message
func GetMessage(err error) string {
	he, ok := err.(*HTTPError)
	if ok {
		return he.Message
	}
	return err.Error()
}

// Read read function
func (er *errReadCloser) Read(p []byte) (n int, err error) {
	return 0, er.customErr
}

// Close close function
func (er *errReadCloser) Close() error {
	return nil
}

// NewErrorReadCloser create an read error
func NewErrorReadCloser(err error) io.ReadCloser {
	r := &errReadCloser{
		customErr: err,
	}
	return r
}
