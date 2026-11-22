-- Create "banned_rating_devices" table
CREATE TABLE IF NOT EXISTS `banned_rating_devices` (
  `device_hash` varchar(64) NOT NULL,
  `reason` varchar(256) NOT NULL DEFAULT "",
  `banned_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`device_hash`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
