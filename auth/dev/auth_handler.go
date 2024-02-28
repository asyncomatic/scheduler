package dev

import (
	"crypto/sha256"
	"github.com/caitlinelfring/go-env-default"
	"net/http"
	"strings"
)

type AuthHandler struct {
	token string
}

func NewDevAuthHandler() *AuthHandler {
	return &AuthHandler{
		token: env.GetDefault("SCHEDULER_AUTH_TOKEN", "")}
}

func (a *AuthHandler) Authn(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if a.token == "" {
			h(w, r)
			return
		} else {
			token := strings.TrimLeft(r.Header.Get("Authorization"), "Bearer ")
			if sha256.Sum256([]byte(token)) == sha256.Sum256([]byte(a.token)) {
				h(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Bearer", charset="UTF-8"`)
		w.WriteHeader(http.StatusUnauthorized)
	})

}
