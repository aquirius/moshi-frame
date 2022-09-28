
CREATE TABLE `plants` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `pluid` varchar(36) NOT NULL,
  `created_ts` int unsigned NOT NULL,
  `planted_ts` int unsigned NOT NULL,
  `harvested_ts` int unsigned NOT NULL,
  `nutrient_id` int unsigned NOT NULL,
  `pot_id` int unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `plants_pluid_IDX` (`pluid`) USING BTREE,
  CONSTRAINT `fk_plants_nutrients` FOREIGN KEY (`nutrient_id`) REFERENCES `nutrients` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_plants_pots` FOREIGN KEY (`pot_id`) REFERENCES `pots` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci