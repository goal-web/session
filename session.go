package session

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/logs"
	"github.com/goal-web/supports/utils"
	"net/http"
	"time"
)

type Session struct {
	id         string
	name       string
	started    bool
	path       string
	domain     string
	lifetime   time.Duration
	attributes map[string]string
	request    contracts.HttpRequest
	changed    bool

	store contracts.SessionStore
}

func New(config Config, request contracts.HttpRequest, store contracts.SessionStore) contracts.Session {
	return &Session{
		id:         "",
		name:       config.Name,
		started:    false,
		request:    request,
		path:       request.Path(),
		store:      store,
		domain:     config.Domain,
		lifetime:   config.Lifetime,
		attributes: map[string]string{},
	}
}

func (session *Session) GetName() string {
	return session.name
}

func (session *Session) SetName(name string) {
	session.name = name
}

func (session *Session) GetId() string {
	return session.id
}

func (session *Session) SetId(id string) {
	session.id = id
}

func (session *Session) Start() bool {
	cookieValue, err := session.request.Cookie(session.name)
	if err != nil {
		logs.WithError(err).Debug("Failed to load cookies")
	} else {
		session.id = cookieValue
	}

	if session.id == "" {
		session.generateSessionId()
	}
	session.loadSession()
	session.started = true
	return true
}

func (session *Session) loadSession() {
	session.attributes = session.store.LoadSession(session.GetId())
}

func (session *Session) Save() {
	if session.changed {
		session.store.Save(session.GetId(), session.attributes)
	}
}

func (session *Session) All() map[string]string {
	return session.attributes
}

func (session *Session) Exists(key string) bool {
	_, exists := session.attributes[key]
	return exists
}

func (session *Session) Has(key string) bool {
	value, exists := session.attributes[key]
	return exists && value != ""
}

func (session *Session) Get(key, defaultValue string) string {
	value, exists := session.attributes[key]
	if !exists || value == "" {
		return defaultValue
	}
	return value
}

func (session *Session) Pull(key, defaultValue string) string {
	session.changed = true
	value, exists := session.attributes[key]
	if !exists || value == "" {
		return defaultValue
	}
	delete(session.attributes, key)
	return value
}

func (session *Session) Put(key, value string) {
	session.changed = true
	session.attributes[key] = value
}

func (session *Session) Token() string {
	return session.Get("_token", "")
}

func (session *Session) RegenerateToken() {
	session.id = utils.RandStr(40)
	session.request.SetCookie(&http.Cookie{
		Name:    session.name,
		Value:   session.id,
		Expires: time.Now().Add(time.Second * session.lifetime),
	})
}

func (session *Session) Remove(key string) string {
	return session.Pull(key, "")
}

func (session *Session) Forget(keys ...string) {
	session.changed = true
	for _, key := range keys {
		delete(session.attributes, key)
	}
}

func (session *Session) Flush() {
	session.changed = true
	session.attributes = make(map[string]string)
}

func (session *Session) Invalidate() bool {
	session.Flush()
	return session.Migrate(true)
}

func (session *Session) Regenerate(destroy bool) bool {
	if !session.Migrate(destroy) {
		session.RegenerateToken()
	}
	return true
}

func (session *Session) Migrate(destroy bool) bool {
	if destroy {
		// todo: $session->handler->destroy($session->getId());
	}
	session.generateSessionId()
	return true
}

func (session *Session) IsStarted() bool {
	return session.started
}

func (session *Session) generateSessionId() {
	session.id = utils.RandStr(40)
	session.request.SetCookie(&http.Cookie{
		Name:    session.name,
		Value:   session.id,
		Expires: time.Now().Add(time.Second * session.lifetime),
	})
}

func (session *Session) PreviousUrl() string {
	return session.Get("_previous.url", "")
}

func (session *Session) SetPreviousUrl(url string) {
	session.Put("_previous.url", url)
}
