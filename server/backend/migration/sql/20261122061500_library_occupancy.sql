-- Create "library_occupancies" table
CREATE TABLE IF NOT EXISTS `library_occupancies` (
  `library_id` varchar(64) NOT NULL,
  `name` varchar(128) NOT NULL,
  `campus` varchar(64) NOT NULL,
  `count` int unsigned NOT NULL,
  `max_count` int unsigned NOT NULL,
  `percent` float NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`library_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "library_occupancy_history" table
CREATE TABLE IF NOT EXISTS `library_occupancy_history` (
  `library_id` varchar(64) NOT NULL,
  `sample_date` char(10) NOT NULL,
  `bucket` tinyint NOT NULL,
  `weekday` tinyint NOT NULL,
  `avg_count` double NOT NULL,
  `avg_percent` double NOT NULL,
  `sample_count` int NOT NULL DEFAULT 1,
  PRIMARY KEY (`library_id`, `sample_date`, `bucket`),
  INDEX `idx_library_occ_history_date` (`sample_date`),
  INDEX `idx_library_occ_history_lookup` (`library_id`, `weekday`, `bucket`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
