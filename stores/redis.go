package stores

import (
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/logs"
	"time"
)

type Redis struct {
	lifetime time.Duration
	redis    contracts.RedisConnection
	key      string
}

func RedisStore(key string, lifetime time.Duration, redis contracts.RedisConnection) contracts.SessionStore {
	return &Redis{
		lifetime: lifetime,
		key:      key,
		redis:    redis,
	}
}

func (this *Redis) LoadSession(id string) map[string]string {
	sessions, err := this.redis.HGetAll(fmt.Sprintf(this.key, id))
	if err != nil {
		logs.WithError(err).
			WithField("key", fmt.Sprintf(this.key, id)).
			Warn("LoadSession err")
	}
	if sessions == nil {
		return make(map[string]string)
	}
	return sessions
}

func (this *Redis) Save(id string, sessions map[string]string) {
	values := make([]interface{}, 0)
	for key, value := range sessions {
		values = append(values, key, value)
	}
	_, err := this.redis.HMSet(fmt.Sprintf(this.key, id), values...)
	if err != nil {
		logs.WithError(err).
			WithField("key", fmt.Sprintf(this.key, id)).
			Warn("session save err")
	}
}
