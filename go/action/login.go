package action

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/sumiredc/webauthn/env"
	"github.com/sumiredc/webauthn/model"
	"github.com/sumiredc/webauthn/repository"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Login(c echo.Context, db *gorm.DB, wa *webauthn.WebAuthn) error {
	uRepo := repository.NewUserRepository(db)
	cRepo := repository.NewCredentialRepository(db)
	var sRepo *repository.SessionRepository
	if r, err := repository.NewSessionRepository(c); err != nil {
		log.Println("Session Repository の初期化に失敗しました", err)
		return loginErrResponse(c)
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

	var token string
	var credential *webauthn.Credential

	if cr, err := wa.FinishDiscoverableLogin(func(rawID, userHandle []byte) (webauthn.User, error) {
		cID := base64.StdEncoding.EncodeToString(rawID)

		var user *model.User
		if u, err := uRepo.GetWithCredential(cID, string(userHandle)); err != nil {
			log.Println("ユーザー情報が見つかりませんでした", err)
			return nil, err
		} else {
			user = u
		}

		claims := jwt.MapClaims{
			"sub": strconv.FormatUint(user.ID, 10),
			"exp": time.Now().Add(15 * time.Minute).Unix(),
		}

		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		if t, err := jwtToken.SignedString(env.DummyJWTSecret); err != nil {
			log.Println("JWT の生成に失敗しました", err)
			return nil, err
		} else {
			token = t
		}

		return user, nil
	}, *sessionData, c.Request()); err != nil {
		log.Println("認証に失敗しました", err)
		return registerErrResponse(c)
	} else {
		credential = cr
	}

	if _, err := cRepo.Update(credential); err != nil {
		log.Println("認証情報の更新に失敗しました", err)
		return registerErrResponse(c)
	}

	return c.JSON(http.StatusOK, map[string]any{"token": token})
}

// client 側へ、一律で同じレスポンスを返却する
func loginErrResponse(c echo.Context) error {
	return c.JSON(http.StatusBadRequest, map[string]any{"message": "認証に失敗しました"})
}
