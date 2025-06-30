CREATE TABLE IF NOT EXISTS `webauthn_credential` (
    `credential_id` VARCHAR(100) PRIMARY KEY COMMENT '認証情報の識別子',
    `user_id` BIGINT NOT NULL COMMENT 'サービス側で管理しているユーザー識別子',
    `public_key` BLOB NOT NULL COMMENT '公開鍵情報',
    `attestation_type` VARCHAR(50) NOT NULL COMMENT 'サーバーから認証器へ要求する証明書レベル',
    `transport` JSON NOT NULL COMMENT '認証器がサポートしているトランスポート (USB, NFC, BLE, INTERNAL)',
    `user_present_flg` BOOLEAN NOT NULL COMMENT 'ユーザーが操作を行ったことを示すフラグ',
    `user_verified_flg` BOOLEAN NOT NULL COMMENT 'ユーザーが認証機に対して本人であることを証明したフラグ',
    `backup_eligible_flg` BOOLEAN NOT NULL COMMENT 'デバイス間同期が可能な鍵か示すフラグ',
    `backup_state_flg` BOOLEAN NOT NULL COMMENT '認証器のバックアップ状態',
    `aaguid` BLOB NOT NULL COMMENT '認証器のモデル識別子',
    `sign_count` INT NOT NULL COMMENT 'ログイン回数',
    `clone_warning` BOOLEAN NOT NULL COMMENT '認証器が複製された可能性',
    `attachment` VARCHAR(50) NOT NULL COMMENT '認証器とデバイスの接続形態',

    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    INDEX user_id(user_id)
)
