package service

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/vicanso/novel/model"
)

var (
	json               = jsoniter.ConfigCompatibleWithStandardLibrary
	getClient          = model.GetClient
	getClientByOptions = model.GetClientByOptions
)
