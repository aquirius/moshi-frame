
CREATE TABLE `pots` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `puid` varchar(36) NOT NULL,
  `stack_id` int unsigned NOT NULL,
  `user_id` int unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `pots_puid_IDX` (`puid`) USING BTREE,
  CONSTRAINT `fk_pots_stacks` FOREIGN KEY (`stack_id`) REFERENCES `stacks` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_pots_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;