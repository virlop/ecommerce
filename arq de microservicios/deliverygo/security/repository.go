package security

import (
	"encoding/json"
	"net/http"
	"time"

	"deliverygo/tools/env"
	"deliverygo/tools/errs"
	"deliverygo/tools/log"

	"github.com/go-playground/validator/v10"
	gocache "github.com/patrickmn/go-cache"
)

var cache = gocache.New(60*time.Minute, 10*time.Minute)

func getRemoteToken(token string, deps ...interface{}) (*User, error) {
	// Buscamos el usuario remoto
	req, err := http.NewRequest("GET", env.Get().SecurityServerURL+"/users/current", nil)
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, errs.Unauthorized
	}
	req.Header.Add("Authorization", "Bearer "+token)
	if corrId, ok := log.Get(deps...).Data[log.LOG_FIELD_CORRELATION_ID].(string); ok {
		req.Header.Add(log.LOG_FIELD_CORRELATION_ID, corrId)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Get(deps...).Error(err)
		return nil, errs.Unauthorized
	}
	defer resp.Body.Close()

	user := &User{}
	err = json.NewDecoder(resp.Body).Decode(user)
	if err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		log.Get(deps...).Error(err)
		return nil, err
	}
	return user, nil
}
