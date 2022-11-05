CREATE TABLE `messages` (
  `id` INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
	`avatar` VARCHAR(32) NOT NULL,
  `date` BIGINT NOT NULL,
	`name` VARCHAR(64) NOT NULL,
	`content` VARCHAR(256) NOT NULL,
	`site` VARCHAR(64) NOT NULL,
	`reply` INT NOT NULL,
	`email` VARCHAR(64) NOT NULL,
	`mail_notice` TINYINT(1) NOT NULL,
	`owner` TINYINT(1) NOT NULL
) DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
