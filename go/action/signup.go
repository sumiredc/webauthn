package action

import (
	"encoding/json"
	"log"
	"net/http"
	"wa/model"
	"wa/repository"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type SignUpRequest struct {
	Username string `json:"username"`
}

func SignUp(c echo.Context, db *gorm.DB, wa *webauthn.WebAuthn) error {
	req := &SignUpRequest{}
	if err := c.Bind(req); err != nil {
		log.Println("リクエストの解決に失敗しました", err)
		return c.JSON(http.StatusBadRequest, map[string]any{"message": "登録に失敗しました"})
	}

	if req.Username == "" {
		return c.JSON(http.StatusUnprocessableEntity, map[string]any{"message": "ユーザー名を入力してください"})
	}

	uRepo := repository.NewUserRepository(db)
	var sRepo *repository.SessionRepository
	if r, err := repository.NewSessionRepository(c); err != nil {
		log.Println("Session Repository の初期化に失敗しました", err)
		return registerOptionsErrResponse(c)
	} else {
		sRepo = r
	}

	if exists, err := uRepo.ExistsByUsername(req.Username); err != nil {
		log.Println("ユーザーの存在確認に失敗しました", err)
		return c.JSON(http.StatusInternalServerError, map[string]any{"message": "ユーザーの登録に失敗しました"})
	} else if exists {
		return c.JSON(http.StatusUnprocessableEntity, map[string]any{"message": "指定されたユーザー名は使用できません"})
	}

	user := model.User{
		Username:   req.Username,
		UserHandle: uuid.NewString(),
	}

	if userJson, err := json.Marshal(user); err != nil {
		log.Println("User の情報の生成に失敗しました", err)
		return c.JSON(http.StatusInternalServerError, map[string]any{"message": "ユーザーの登録に失敗しました"})
	} else {
		sRepo.Set(sessionUserKey, string(userJson))
	}

	if err := sRepo.Save(); err != nil {
		log.Println("Session の保存に失敗しました", err)
		return c.JSON(http.StatusInternalServerError, map[string]any{"message": "ユーザーの登録に失敗しました"})
	}

	return c.JSON(http.StatusCreated, "")
}
