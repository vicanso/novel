package middleware

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/vicanso/novel/utils"
)

var (
	json       = jsoniter.ConfigCompatibleWithStandardLibrary
	getTrackID = utils.GetTrackID
	getAccount = utils.GetAccount
)
