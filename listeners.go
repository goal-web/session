package session

import (
	"github.com/goal-web/contracts"
)

type RequestBeforeListener struct {
}

func (listener *RequestBeforeListener) Handle(event contracts.Event) {
}

type RequestAfterListener struct {
}

// Handle 如果开启了 session 那么请求结束时保存 session
func (listener *RequestAfterListener) Handle(event contracts.Event) {
	//if session.IsStarted() {
	//	session.Save()
	//}
}
