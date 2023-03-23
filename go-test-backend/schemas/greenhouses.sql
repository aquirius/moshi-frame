
CREATE TABLE `greenhouses` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `guid` varchar(36) NOT NULL,
  `display_name` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `address` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `zip` varchar(100) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `status` enum('active','blocked') NOT NULL DEFAULT 'active',
  `type` enum('indoor','outdoor','greenhouse') NOT NULL DEFAULT 'indoor',
  PRIMARY KEY (`id`),
  UNIQUE KEY `greenhouses_guid_IDX` (`guid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci