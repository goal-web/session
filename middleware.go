package session

import (
	"github.com/goal-web/contracts"
)

func StartSession(session contracts.Session, request contracts.HttpRequest, next contracts.Pipe) any {
	session.Start()
	return next(request)
}
