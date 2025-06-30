CREATE TABLE IF NOT EXISTS `credential` (
    `credential_id` VARCHAR(100) UNIQUE COMMENT '認証情報の識別子',
    `user_id` BIGINT NOT NULL COMMENT 'サービス側で管理しているユーザー識別子',
    `json` JSON NOT NULL COMMENT 'WebAuthn.Credential の JSON 形式',

    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX user_id(user_id)
)
