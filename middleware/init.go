package middleware

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/vicanso/novel/util"
)

var (
	json       = jsoniter.ConfigCompatibleWithStandardLibrary
	getTrackID = util.GetTrackID
	getAccount = util.GetAccount
	resErr     = util.ResErr

	getLogger = util.GetLogger
)
