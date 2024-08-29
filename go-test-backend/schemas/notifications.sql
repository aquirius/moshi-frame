
CREATE TABLE `notifications` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `nuid` varchar(36) NOT NULL,
  `created_ts` int unsigned NOT NULL,
  `checked_ts` int unsigned NOT NULL,
  `done_ts` int unsigned NOT NULL,
  `title` varchar(36) NOT NULL,
  `message` varchar(999) NOT NULL,
  `user_id` int unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `notifications_nuid_IDX` (`nuid`) USING BTREE,
  CONSTRAINT `fk_notifications_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci