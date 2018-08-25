package controller

import (
	"github.com/vicanso/novel/middleware"
	"github.com/vicanso/novel/utils"
)

var (
	res            = utils.Res
	resNoContent   = utils.ResNoContent
	resCreated     = utils.ResCreated
	resJPEG        = utils.ResJPEG
	resPNG         = utils.ResPNG
	resErr         = utils.ResErr
	setCache       = utils.SetCache
	setNoCache     = utils.SetNoCache
	setNoStore     = utils.SetNoStore
	validate       = utils.Validate
	getRequestBody = utils.GetRequestBody

	newDefaultTracker = middleware.NewDefaultTracker
	newTracker        = middleware.NewTracker
)
