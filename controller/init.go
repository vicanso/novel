package controller

import (
	"github.com/vicanso/novel/middleware"
	"github.com/vicanso/novel/util"
)

var (
	res             = util.Res
	resNoContent    = util.ResNoContent
	resCreated      = util.ResCreated
	resJPEG         = util.ResJPEG
	resPNG          = util.ResPNG
	resErr          = util.ResErr
	setCache        = util.SetCache
	setNoCache      = util.SetNoCache
	setNoStore      = util.SetNoStore
	validate        = util.Validate
	getSession      = util.GetSession
	getRequestBody  = util.GetRequestBody
	getRequestQuery = util.GetRequestQuery
	getUTCNow       = util.GetUTCNow
	getNow          = util.GetNow

	newDefaultTracker = middleware.NewDefaultTracker
	newTracker        = middleware.NewTracker

	getLogger        = util.GetLogger
	getUserLogger    = util.CreateUserLogger
	getContextLogger = util.GetContextLogger
)
