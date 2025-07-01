package action

import (
	"encoding/json"
	"log"
	"net/http"
	"wa/repository"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func LoginOptions(c echo.Context, db *gorm.DB, wa *webauthn.WebAuthn) error {
	var sRepo *repository.SessionRepository
	if r, err := repository.NewSessionRepository(c); err != nil {
		log.Println("Session Repository の初期化に失敗しました", err)
		return loginOptionsErrResponse(c)
	} else {
		sRepo = r
	}

	mediation := protocol.CredentialMediationRequirement("required")
	userVerification := protocol.UserVerificationRequirement("preferred")

	var options *protocol.CredentialAssertion
	var sessionData *webauthn.SessionData
	if o, s, err := wa.BeginDiscoverableMediatedLogin(
		mediation,
		webauthn.WithUserVerification(userVerification),
	); err != nil {
		log.Println("Challenge の生成に失敗しました", err)
		return loginOptionsErrResponse(c)
	} else {
		options = o
		sessionData = s
	}

	if j, err := json.Marshal(sessionData); err != nil {
		log.Println("webauthn.SessionData の json 変換に失敗しました", err)
		return registerOptionsErrResponse(c)
	} else {
		sRepo.Set(sessionWebAuthnSessionDataKey, string(j))
	}

	if err := sRepo.Save(); err != nil {
		log.Println("session の保存に失敗しました", err)
		return registerOptionsErrResponse(c)
	}

	return c.JSON(http.StatusOK, map[string]any{"options": options})
}

// client 側へ、一律で同じレスポンスを返却する
func loginOptionsErrResponse(c echo.Context) error {
	return c.JSON(http.StatusBadRequest, map[string]any{"message": "認証オプションの取得に失敗しました"})
}
