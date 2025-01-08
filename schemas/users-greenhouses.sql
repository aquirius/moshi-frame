CREATE TABLE `users_greenhouses` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uguid` varchar(36) NOT NULL,
  `user_id` int unsigned NOT NULL,
  `greenhouse_id` int unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_greenhouses_guid_IDX` (`uguid`) USING BTREE,
  CONSTRAINT `fk_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_greenhouses` FOREIGN KEY (`greenhouse_id`) REFERENCES `greenhouses` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;