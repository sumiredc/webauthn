package action

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"wa/env"
	"wa/model"
	"wa/repository"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Register(c echo.Context, db *gorm.DB, wa *webauthn.WebAuthn) error {
	var sRepo *repository.SessionRepository
	if r, err := repository.NewSessionRepository(c); err != nil {
		log.Println("Session Repository の初期化に失敗しました", err)
		return registerErrResponse(c)
	} else {
		sRepo = r
	}

	var sessionDataJson []byte
	if v := sRepo.Get(sessionWebAuthnSessionDataKey); v == "" {
		log.Printf("session.%s の取得に失敗しました\n", sessionWebAuthnSessionDataKey)
		return registerErrResponse(c)
	} else {
		sessionDataJson = []byte(v)
	}

	sessionData := &webauthn.SessionData{}
	if err := json.Unmarshal(sessionDataJson, sessionData); err != nil {
		log.Println("SessionData のパースに失敗しました", err)
		return registerErrResponse(c)
	}

	var userJson []byte
	if v := sRepo.Get(sessionUserKey); v == "" {
		log.Printf("session.%s の取得に失敗しました\n", sessionUserKey)
		return registerErrResponse(c)
	} else {
		userJson = []byte(v)
	}

	sUser := &model.User{}
	if err := json.Unmarshal(userJson, sUser); err != nil {
		log.Println("ユーザー情報の取得に失敗しました", err)
		return registerErrResponse(c)
	}

	if err := sRepo.Delete(); err != nil {
		log.Println("セッションの削除に失敗しました", err)
		return registerErrResponse(c)
	}

	var token string

	txErr := db.Transaction(func(tx *gorm.DB) error {
		uRepo := repository.NewUserRepository(tx)
		wRepo := repository.NewCredentialRepository(tx)

		var credential *webauthn.Credential
		if cr, err := wa.FinishRegistration(sUser, *sessionData, c.Request()); err != nil {
			log.Println("公開鍵の取得に失敗しました", err)
			return err
		} else {
			credential = cr
		}

		var user *model.User
		if u, err := uRepo.Create(sUser.Username, sUser.UserHandle); err != nil {
			log.Println("ユーザー情報の登録に失敗しました", err)
			return err
		} else {
			user = u
		}

		if _, err := wRepo.Create(user.ID, credential); err != nil {
			log.Println("公開鍵の登録に失敗しました", err)
			return err
		}

		claims := jwt.MapClaims{
			"sub": strconv.FormatUint(user.ID, 10),
			"exp": time.Now().Add(15 * time.Minute).Unix(),
		}

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		if t, err := jwtToken.SignedString(env.DummyJWTSecret); err != nil {
			log.Println("JWT の生成に失敗しました", err)
			return err
		} else {
			token = t
		}

		return nil
	})

	if txErr != nil {
		return registerErrResponse(c)
	}

	return c.JSON(http.StatusCreated, map[string]any{"token": token})
}

// client 側へ、一律で同じレスポンスを返却する
func registerErrResponse(c echo.Context) error {
	return c.JSON(http.StatusBadRequest, map[string]any{"message": "登録に失敗しました"})
}
