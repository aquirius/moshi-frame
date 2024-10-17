CREATE TABLE sprouts (
	`id` INT unsigned primary key NOT NULL AUTO_INCREMENT,
  `sproutuid` varchar(36) NOT NULL,
  `pH` FLOAT NOT NULL CHECK (pH >= 4 AND pH <= 8),
  `TDS` INT NOT NULL CHECK (TDS >= 300 AND TDS <= 3500),
  `ORP` INT NOT NULL CHECK (ORP >= 300 AND ORP <= 700),
  `h2oTemp` FLOAT NOT NULL CHECK (h2oTemp >= 10 AND h2oTemp <= 30),
  `airTemp` FLOAT NOT NULL CHECK (airTemp >= 20 AND airTemp <= 35),
  `humidity` INT NOT NULL CHECK (humidity >= 50 AND humidity <= 100),
  `stack_id` int unsigned NOT NULL,
  UNIQUE KEY `sprout_sproutuid_IDX` (`sproutuid`) USING BTREE,
  CONSTRAINT `fk_sprout_stacks` FOREIGN KEY (`stack_id`) REFERENCES `stacks` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;