-- Create "canteen_head_count_history" table
CREATE TABLE IF NOT EXISTS `canteen_head_count_history` (
  `canteen_id` varchar(64) NOT NULL,
  `sample_date` char(10) NOT NULL,
  `bucket` tinyint NOT NULL,
  `weekday` tinyint NOT NULL,
  `avg_count` double NOT NULL,
  `avg_percent` double NOT NULL,
  `sample_count` int NOT NULL DEFAULT 1,
  PRIMARY KEY (`canteen_id`, `sample_date`, `bucket`),
  INDEX `idx_canteen_hc_history_date` (`sample_date`),
  INDEX `idx_canteen_hc_history_lookup` (`canteen_id`, `weekday`, `bucket`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
