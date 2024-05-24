package session

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/session/stores"
)

type ServiceProvider struct {
	app contracts.Application
}

func NewService() contracts.ServiceProvider {
	return &ServiceProvider{}
}

func (provider *ServiceProvider) Register(application contracts.Application) {
	provider.app = application

	application.Bind("session", func(
		config contracts.Config,
		request contracts.HttpRequest,
		encryptor contracts.Encryptor,
		redis contracts.RedisFactory,
	) contracts.Session {
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

		case "redis":
			store = stores.RedisStore(sessionConfig.Key, sessionConfig.Lifetime, redis.Connection(sessionConfig.Connection))
		}

		session := New(sessionConfig, request, store)

		request.Set("session", session)
		return session
	})
}

func (provider *ServiceProvider) Start() error {
	provider.app.Call(func(dispatcher contracts.EventDispatcher) {
		dispatcher.Register("RESPONSE_BEFORE", &RequestAfterListener{})
	})
	return nil
}

func (provider *ServiceProvider) Stop() {
}
