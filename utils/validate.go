package utils

import (
	"fmt"
	"regexp"

	"github.com/asaskevich/govalidator"
)

var (
	paramTagRegexMap = govalidator.ParamTagRegexMap
	paramTagMap      = govalidator.ParamTagMap
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
	AddRegexValidate("xIntRange", "^xIntRange\\((\\d+)\\|(\\d+)\\)$", func(value string, params ...string) bool {
		return govalidator.InRangeInt(value, params[0], params[1])
	})
}

// Validate 校验数据
func Validate(s interface{}, data interface{}) (err error) {
	// TODO 了解自定义的校验方式
	if data != nil {
		switch data.(type) {
		case []byte:
			err = json.Unmarshal(data.([]byte), s)
			if err != nil {
				fmt.Println(err)
				err = NewJSONParseError(err)
				return
			}
		default:
			buf, e := json.Marshal(data)
			if e != nil {
				err = NewJSONParseError(e)
				return
			}
			e = json.Unmarshal(buf, s)
			if e != nil {
				err = NewJSONParseError(e)
				return
			}
		}
	}
	_, err = govalidator.ValidateStruct(s)
	if err != nil {
		err = NewValidateError(err)
	}
	return
}

// AddRegexValidate add a regexp validate
func AddRegexValidate(name, reg string, fn govalidator.ParamValidator) {
	paramTagRegexMap[name] = regexp.MustCompile(reg)
	AddValidate(name, fn)
}

// AddValidate add validate
func AddValidate(name string, fn govalidator.ParamValidator) {
	paramTagMap[name] = fn
}
