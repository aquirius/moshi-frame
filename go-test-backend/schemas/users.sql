
CREATE TABLE `users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `email` varchar(200) NOT NULL,
  `password_hash` char(64) CHARACTER SET ascii COLLATE ascii_general_ci NOT NULL,
  `registered_ts` int unsigned NOT NULL,
  `last_login_ts` int unsigned DEFAULT NULL,
  `status` enum('active','blocked') NOT NULL DEFAULT 'active',
  `display_name` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `first_name` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `last_name` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `birthday` int unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_email_IDX` (`email`) USING BTREE,
  UNIQUE KEY `users_uuid_IDX` (`uuid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci