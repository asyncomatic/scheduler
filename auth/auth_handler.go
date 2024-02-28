package auth

import (
	"net/http"
	"reflect"
	"scheduler/auth/dev"
)

var AuthHandlerRegistry = map[string]AuthHandler{
	"_default": dev.NewDevAuthHandler(),
	"dev":      dev.NewDevAuthHandler(),
}

type AuthHandler interface {
	Authn(handlerFunc http.HandlerFunc) http.HandlerFunc
}

func NewAuthHandler(authType string) AuthHandler {
	return reflect.ValueOf(AuthHandlerRegistry[authType]).Interface().(AuthHandler)
}
