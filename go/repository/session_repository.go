package repository

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

const sessionKey = "session"

type SessionRepository struct {
	ctx     echo.Context
	session *sessions.Session
}

func NewSessionRepository(c echo.Context) (*SessionRepository, error) {
	s, err := session.Get(sessionKey, c)
	if err != nil {
		return nil, err
	}

	return &SessionRepository{
		ctx:     c,
		session: s,
	}, nil
}

func (r *SessionRepository) Set(k string, v string) {
	r.session.Values[k] = v
}

func (r *SessionRepository) Save() error {
	r.session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   300,
		HttpOnly: true,
		Secure:   true,
	}

	return r.session.Save(r.ctx.Request(), r.ctx.Response())
}

func (r *SessionRepository) Get(k string) string {
	var val any
	if v, ok := r.session.Values[k]; !ok {
		return ""
	} else {
		val = v
	}

	if v, ok := val.(string); !ok {
		return ""
	} else {
		return v
	}
}

func (r *SessionRepository) Delete() error {
	r.session.Values = make(map[any]any)
	r.session.Options.MaxAge = -1

	return r.session.Save(r.ctx.Request(), r.ctx.Response())
}
