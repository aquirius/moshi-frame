
CREATE TABLE `crops` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `cuid` varchar(36) NOT NULL,
  `crop_name` varchar(36) NOT NULL,
  `air_temp_min` decimal(4,1) NOT NULL,
  `air_temp_max` decimal(4,1) NOT NULL,
  `humidity_min` decimal(4,1) NOT NULL,
  `humidity_max` decimal(4,1) NOT NULL,
  `ph_level_min` decimal(3,1) NOT NULL,
  `ph_level_max` decimal(3,1) NOT NULL,
  `orp_min` decimal(5,1) NOT NULL,
  `orp_max` decimal(5,1) NOT NULL,
  `tds_min` int unsigned NOT NULL,
  `tds_max` int unsigned NOT NULL,
  `water_temp_min` decimal(4,1) NOT NULL,
  `water_temp_max` decimal(4,1) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci