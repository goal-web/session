package stores

import (
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/logs"
	"net/http"
	"strings"
	"time"
)

type Cookie struct {
	name      string
	encrypt   bool
	lifetime  time.Duration
	encryptor contracts.Encryptor
	request   contracts.HttpRequest
}

func CookieStore(name string, lifetime time.Duration, request contracts.HttpRequest, encryptor contracts.Encryptor) contracts.SessionStore {
	return &Cookie{
		name:      name,
		lifetime:  lifetime,
		request:   request,
		encrypt:   encryptor != nil,
		encryptor: encryptor,
	}
}

func (this *Cookie) LoadSession(id string) map[string]string {
	attributes := make(map[string]string, 0)
	for _, cookie := range this.request.Cookies() {
		if strings.HasPrefix(cookie.Name, this.name) {
			value := cookie.Value
			if this.encrypt {
				decrypted, err := this.encryptor.Decode(cookie.Value)
				if err != nil {
					value = cookie.Value
					logs.WithError(err).Warn(fmt.Sprintf("cookie %s decryption failed", cookie.Name))
				} else {
					value = decrypted
				}
			}
			attributes[strings.ReplaceAll(cookie.Name, this.name, "")] = value
		}
	}
	return attributes
}

func (this *Cookie) Save(id string, sessions map[string]string) {
	for key, value := range sessions {
		if this.encrypt {
			value = this.encryptor.Encode(value)
		}
		this.request.SetCookie(&http.Cookie{
			Name:    this.CookieKey(key),
			Value:   value,
			Expires: time.Now().Add(time.Second * this.lifetime),
		})
	}
}

func (this *Cookie) CookieKey(key string) string {
	return this.name + key
}
