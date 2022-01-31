package session

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/session/stores"
)

type ServiceProvider struct {
	app contracts.Application
}

func (this *ServiceProvider) Register(application contracts.Application) {
	this.app = application

	application.Bind("session", func(config contracts.Config, request contracts.HttpRequest, encryptor contracts.Encryptor) contracts.Session {
		if session, isSession := request.Get("session").(contracts.Session); isSession {
			return session
		}

		sessionConfig := config.Get("session").(Config)
		var store contracts.SessionStore

		switch sessionConfig.Driver {
		case "cookie":
			if sessionConfig.Encrypt {
				store = stores.CookieStore(sessionConfig.Name, sessionConfig.Lifetime, request, encryptor)
			} else {
				store = stores.CookieStore(sessionConfig.Name, sessionConfig.Lifetime, request, nil)
			}
		}

		session := New(sessionConfig, request, store)

		request.Set("session", session)
		return session
	})
}

func (this *ServiceProvider) Start() error {
	this.app.Call(func(dispatcher contracts.EventDispatcher) {
		dispatcher.Register("RESPONSE_BEFORE", &RequestAfterListener{})
	})
	return nil
}

func (this *ServiceProvider) Stop() {
}
