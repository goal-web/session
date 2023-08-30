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

func (c *Cookie) LoadSession(id string) map[string]string {
	attributes := make(map[string]string)
	for _, cookie := range c.request.Cookies() {
		if strings.HasPrefix(cookie.Name, c.name) {
			value := cookie.Value
			if c.encrypt {
				decrypted, err := c.encryptor.Decrypt([]byte(cookie.Value))
				if err != nil {
					value = cookie.Value
					logs.WithError(err).Warn(fmt.Sprintf("cookie %s decryption failed", cookie.Name))
				} else {
					value = string(decrypted)
				}
			}
			attributes[strings.ReplaceAll(cookie.Name, c.name, "")] = value
		}
	}
	return attributes
}

func (c *Cookie) Save(id string, sessions map[string]string) {
	for key, value := range sessions {
		if c.encrypt {
			value = string(c.encryptor.Encrypt([]byte(value)))
		}
		c.request.SetCookie(&http.Cookie{
			Name:    c.CookieKey(key),
			Value:   value,
			Expires: time.Now().Add(time.Second * c.lifetime),
		})
	}
}

func (c *Cookie) CookieKey(key string) string {
	return c.name + key
}
