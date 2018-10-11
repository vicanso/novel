package validate

import (
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	jsoniter "github.com/json-iterator/go"
	"github.com/vicanso/novel/xerror"
)

var (
	paramTagRegexMap = govalidator.ParamTagRegexMap
	paramTagMap      = govalidator.ParamTagMap
	json             = jsoniter.ConfigCompatibleWithStandardLibrary
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
	AddRegex("xIntRange", "^xIntRange\\((\\d+)\\|(\\d+)\\)$", func(value string, params ...string) bool {
		return govalidator.InRangeInt(value, params[0], params[1])
	})

	AddRegex("xIntIn", `^xIntIn\((.*)\)$`, func(value string, params ...string) bool {
		if len(params) == 1 {
			rawParams := params[0]
			parsedParams := strings.Split(rawParams, "|")
			return govalidator.IsIn(value, parsedParams...)
		}
		return false
	})
}

// Do do validate
func Do(s interface{}, data interface{}) (err error) {
	if data != nil {
		switch data.(type) {
		case []byte:
			err = json.Unmarshal(data.([]byte), s)
			if err != nil {
				err = xerror.NewJSON(err.Error())
				return
			}
		default:
			buf, e := json.Marshal(data)
			if e != nil {
				err = xerror.NewJSON(e.Error())
				return
			}
			e = json.Unmarshal(buf, s)
			if e != nil {
				err = xerror.NewJSON(e.Error())
				return
			}
		}
	}
	_, err = govalidator.ValidateStruct(s)
	if err != nil {
		err = xerror.NewValidate(err.Error())
	}
	return
}

// AddRegex add a regexp validate
func AddRegex(name, reg string, fn govalidator.ParamValidator) {
	paramTagRegexMap[name] = regexp.MustCompile(reg)
	Add(name, fn)
}

// Add add validate
func Add(name string, fn govalidator.ParamValidator) {
	paramTagMap[name] = fn
}
