package session

import (
	"github.com/goal-web/contracts"
)

func StartSession(session contracts.Session, request contracts.HttpRequest, next contracts.Pipe) interface{} {
	session.Start()
	return next(request)
}
