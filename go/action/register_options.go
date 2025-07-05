package action

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sumiredc/webauthn/model"
	"github.com/sumiredc/webauthn/repository"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterOptions(c echo.Context, db *gorm.DB, wa *webauthn.WebAuthn) error {
	var sRepo *repository.SessionRepository
	if r, err := repository.NewSessionRepository(c); err != nil {
		log.Println("Session Repository の初期化に失敗しました", err)
		return registerOptionsErrResponse(c)
	} else {
		sRepo = r
	}

	var userJson []byte
	if v := sRepo.Get(sessionUserKey); v == "" {
		log.Printf("session.%s の取得に失敗しました\n", sessionUserKey)
		return registerOptionsErrResponse(c)
	} else {
		userJson = []byte(v)
	}

	user := &model.User{}
	if err := json.Unmarshal(userJson, user); err != nil {
		log.Println("ユーザー情報の取得に失敗しました", err)
		return registerOptionsErrResponse(c)
	}

	// 認証器に対して指定するオプションを構築
	authSelect := protocol.AuthenticatorSelection{
		AuthenticatorAttachment: protocol.AuthenticatorAttachment("platform"),
		RequireResidentKey:      protocol.ResidentKeyRequired(),
		UserVerification:        protocol.VerificationPreferred,
	}
	attestation := protocol.PreferNoAttestation

	// Challenge の生成
	options, sessionData, err := wa.BeginRegistration(
		user,
		webauthn.WithAuthenticatorSelection(authSelect),
		webauthn.WithConveyancePreference(attestation),
	)

	if err != nil {
		log.Println("Challenge の生成に失敗しました", err)
		return registerOptionsErrResponse(c)
	}

	if j, err := json.Marshal(sessionData); err != nil {
		log.Println("webauthn.SessionData の json 変換に失敗しました", err)
		return registerOptionsErrResponse(c)
	} else {
		sRepo.Set(sessionWebAuthnSessionDataKey, string(j))
	}

	if err := sRepo.Save(); err != nil {
		log.Println("セッションへの書き込みに失敗しました", err)
		return registerOptionsErrResponse(c)
	}

	return c.JSON(http.StatusOK, map[string]any{"options": options})
}

// client 側へ、一律で同じレスポンスを返却する
func registerOptionsErrResponse(c echo.Context) error {
	return c.JSON(http.StatusBadRequest, map[string]any{"message": "登録オプションの取得に失敗しました"})
}
