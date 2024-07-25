CREATE TABLE sprouts (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `sproutuid` varchar(36) NOT NULL,
  `pH` FLOAT NOT NULL CHECK (pH >= 4 AND pH <= 8),
  `TDS` INT NOT NULL CHECK (TDS >= 300 AND TDS <= 3500),
  `ORP` INT NOT NULL CHECK (ORP >= 300 AND ORP <= 700),
  `h2oTemp` FLOAT NOT NULL CHECK (water_temperature >= 10 AND water_temperature <= 30),
  `airTemp` FLOAT NOT NULL CHECK (air_temperature >= 20 AND air_temperature <= 35),
  `humidity` INT NOT NULL CHECK (humidity >= 50 AND humidity <= 100),
  `stack_id` int unsigned NOT NULL,
  UNIQUE KEY `sprout_sproutuid_IDX` (`sproutuid`) USING BTREE,
  CONSTRAINT `fk_sprout_stacks` FOREIGN KEY (`stack_id`) REFERENCES `stacks` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci