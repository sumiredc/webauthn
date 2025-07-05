package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/sumiredc/webauthn/action"
	"github.com/sumiredc/webauthn/env"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()

	e.Use(
		middlewareCROS(),
		middlewareSession(),
	)

	db := setUpDB()
	webAuthn := setUpWebAuthn()

	e.POST("signup", func(c echo.Context) error {
		return action.SignUp(c, db, webAuthn)
	})

	e.GET("register_options", func(c echo.Context) error {
		return action.RegisterOptions(c, db, webAuthn)
	})

	e.POST("register", func(c echo.Context) error {
		return action.Register(c, db, webAuthn)
	})

	e.GET("login_options", func(c echo.Context) error {
		return action.LoginOptions(c, db, webAuthn)
	})

	e.POST("login", func(c echo.Context) error {
		return action.Login(c, db, webAuthn)
	})

	// Error Handler
	e.Start(":8080")
}

func middlewareCROS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	})
}

func middlewareSession() echo.MiddlewareFunc {
	return session.Middleware(sessions.NewCookieStore(env.DummySessionKey))
}

func middlewareJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return echo.NewHTTPError(http.StatusUnauthorized, "認証に失敗しました")
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
			return env.DummyJWTSecret, nil
		})
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusUnauthorized, "認証に失敗しました")
		}

		var claims *jwt.RegisteredClaims
		if cl, ok := token.Claims.(*jwt.RegisteredClaims); ok {
			claims = cl
		} else {
			return echo.NewHTTPError(http.StatusUnauthorized, "認証に失敗しました")
		}

		if sub, err := claims.GetSubject(); err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusUnauthorized, "認証に失敗しました")
		} else {
			c.Set("userID", sub)
		}

		return next(c)
	}
}

func setUpWebAuthn() *webauthn.WebAuthn {
	wConf := &webauthn.Config{
		RPDisplayName: "WebAuthn Sample",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:3000"},
		Timeouts: webauthn.TimeoutsConfig{
			Registration: webauthn.TimeoutConfig{
				Enforce:    true,
				Timeout:    time.Second * 60,
				TimeoutUVD: time.Second * 60,
			},
			Login: webauthn.TimeoutConfig{
				Enforce:    true,
				Timeout:    time.Second * 60,
				TimeoutUVD: time.Second * 60,
			},
		},
	}

	webAuthn, err := webauthn.New(wConf)
	if err != nil {
		log.Fatalln("Web Authn の設定に失敗しました")
	}

	return webAuthn
}

func setUpDB() *gorm.DB {
	dsn := "user:password@tcp(mysql:3306)/sample?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("DB 接続に失敗しました")
	}

	return db
}
