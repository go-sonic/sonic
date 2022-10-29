package impl

import (
	"time"

	"github.com/go-sonic/sonic/cache"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
)

const (
	oneTimeTokenPrefix = "OTT-"
	ottExpirationTime  = time.Minute * 5
)

type oneTimeTokenServiceImpl struct {
	Cache cache.Cache
}

func NewOneTimeTokenService(cache cache.Cache) service.OneTimeTokenService {
	return &oneTimeTokenServiceImpl{
		Cache: cache,
	}
}

func (o *oneTimeTokenServiceImpl) Get(oneTimeToken string) (string, bool) {
	v, ok := o.Cache.Get(oneTimeTokenPrefix + oneTimeToken)
	if !ok {
		return "", false
	}
	return v.(string), true
}

func (o *oneTimeTokenServiceImpl) Create(value string) string {
	uuid := util.GenUUIDWithOutDash()
	o.Cache.Set(oneTimeTokenPrefix+uuid, value, ottExpirationTime)
	return uuid
}
