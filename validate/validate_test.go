package validate

import (
	"testing"

	"github.com/vicanso/novel/xerror"

	"github.com/asaskevich/govalidator"
)

type (
	customValidate struct {
		Age int `json:"age,omitempty" valid:"xMyValidate(0|10)"`
	}
	validateStruct struct {
		Age  int `json:"age,omitempty" valid:"xIntRange(0|100)"`
		Type int `json:"type,omitempty" valid:"xIntIn(1|5|10)"`
	}
	params struct {
		// Account account
		Account string `json:"account" valid:"ascii,runelength(4|10)"`
	}
)

func TestValidate(t *testing.T) {
	t.Run("default custom valid", func(t *testing.T) {
		buf := []byte(`{
			"age": 10,
			"type": 1
		}`)
		s := &validateStruct{}
		err := Do(s, buf)
		if err != nil {
			t.Fatalf("default custom valid fail, %v", err)
		}
	})

	t.Run("validate fail", func(t *testing.T) {
		p := &params{}
		buf := []byte(`{"account":"abd"}`)
		err := Do(p, buf)
		if err == nil {
			t.Fatalf("validate should be fail")
		}
	})
	t.Run("validate fail with not json buffer", func(t *testing.T) {
		p := &params{}
		buf := []byte(`{"account":"vicanso}`)
		err := Do(p, buf)
		he := err.(*xerror.HTTPError)
		if he.Category != xerror.ErrCategoryJSON {
			t.Fatalf("validate should be json fail")
		}
	})

	t.Run("validate success", func(t *testing.T) {
		p := &params{}
		account := "vicanso"
		buf := []byte(`{"account":"vicanso"}`)
		err := Do(p, buf)
		if err != nil || p.Account != account {
			t.Fatalf("validate fail, %v", err)
		}
		tmp := &params{}
		err = Do(tmp, p)
		if err != nil || tmp.Account != account {
			t.Fatalf("validate fail, %v", err)
		}
	})

	t.Run("custom validate", func(t *testing.T) {

		AddRegex("xMyValidate", "^xMyValidate\\((\\d+)\\|(\\d+)\\)$", func(value string, params ...string) bool {
			return govalidator.InRangeInt(value, params[0], params[1])
		})
		s := &customValidate{}
		err := Do(s, []byte(`{
			"age": 10
		}`))
		if err != nil {
			t.Fatalf("add regexp validate fail, %v", err)
		}
		err = Do(s, []byte(`{
			"age": 11
		}`))
		if err == nil {
			t.Fatalf("the age over the limit should return error")
		}
	})
}
