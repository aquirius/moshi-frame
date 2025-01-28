
CREATE TABLE `users_crops` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `ucuid` varchar(36) NOT NULL,
  `user_id` int unsigned NOT NULL,
  `crop_id` int unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_crops_cuid_IDX` (`ucuid`) USING BTREE,
  CONSTRAINT `fk_users_crops` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_crops` FOREIGN KEY (`crop_id`) REFERENCES `crops` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;