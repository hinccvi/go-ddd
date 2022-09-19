package tools

import (
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)


func Validator[I any](c echo.Context, i *I) error {
	if err := c.Bind(i); err != nil {
		return err
	}

	if err := c.Validate(i); err != nil {
		return err
	}

	return nil
}

func GetSession(c echo.Context, path, env string) *sessions.Session {
	sess, _ := session.Get("session", c)

	sess.Options = &sessions.Options{
		Path:     path,
		MaxAge:   int((15 * time.Minute).Seconds()),
		HttpOnly: true,
	}

	if env != "local" {
		sess.Options.Secure = true
	}

	return sess
}

func SaveSession(c echo.Context, sess *sessions.Session) error {
	return sess.Save(c.Request(), c.Response())
}
