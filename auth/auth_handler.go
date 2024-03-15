//  Copyright (c) 2024 JC Cormier
//  All rights reserved.
//  SPDX-License-Identifier: MIT
//  For full license text, see LICENSE file in the repo root or https://opensource.org/licenses/MIT

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
