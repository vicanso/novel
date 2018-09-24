package service

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/vicanso/novel/util"
)

var (
	json      = jsoniter.ConfigCompatibleWithStandardLibrary
	getLogger = util.GetLogger
)
