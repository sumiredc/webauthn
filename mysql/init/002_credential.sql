CREATE TABLE IF NOT EXISTS `credential` (
    `credential_id` VARCHAR(1368) NOT NULL COMMENT '認証情報の識別子 (1023 byte / 3 * 4) = base64 max',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT 'サービス側で管理しているユーザー識別子',
    `data` JSON NOT NULL COMMENT 'webauthn.Credential の json データ',

    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    UNIQUE KEY `credential_id`(`credential_id`(255)),
    FOREIGN KEY `user_id`(`user_id`) REFERENCES `user`(`id`)
);
