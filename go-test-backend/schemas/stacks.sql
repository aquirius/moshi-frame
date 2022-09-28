
CREATE TABLE `stacks` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `suid` varchar(36) NOT NULL,
  `greenhouse_id` int unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `stacks_suid_IDX` (`suid`) USING BTREE,
  CONSTRAINT `fk_stacks_greenhouses` FOREIGN KEY (`greenhouse_id`) REFERENCES `greenhouses` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci