
CREATE TABLE `plans` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `pluid` varchar(36) NOT NULL,
  `display_name` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `settings` varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `start_ts` int unsigned NOT NULL,
  `end_ts` int unsigned NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;