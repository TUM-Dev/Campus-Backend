-- Create "dish_rating_votes" table
CREATE TABLE IF NOT EXISTS `dish_rating_votes` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `rating_id` bigint NOT NULL,
  `device_hash` varchar(64) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_dish_rating_votes_rating` (`rating_id`),
  UNIQUE INDEX `idx_dish_rating_votes_unique` (`rating_id`, `device_hash`),
  CONSTRAINT `fk_dish_rating_votes_rating` FOREIGN KEY (`rating_id`) REFERENCES `dish_ratings` (`dishRating`) ON UPDATE CASCADE ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
