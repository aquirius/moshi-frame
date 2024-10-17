
CREATE TABLE `nutrients` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `carbon` int unsigned NOT NULL,
  `hydrogen` int unsigned NOT NULL,
  `oxygen` int unsigned NOT NULL,
  `nitrogen` int unsigned NOT NULL,
  `phosphorus` int unsigned NOT NULL,
  `potassium` int unsigned NOT NULL,
  `sulfur` int unsigned NOT NULL,
  `calcium` int unsigned NOT NULL,
  `magnesium` int unsigned NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;