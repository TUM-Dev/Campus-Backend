-- Modify "dishes" table (retry-safe for interrupted MySQL DDL)
SET @dishes_metadata_alter = (
  SELECT IF(
    EXISTS(
      SELECT 1 FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'dishes' AND COLUMN_NAME = 'labels'
    ),
    'SELECT 1',
    'ALTER TABLE `dishes` ADD COLUMN `labels` varchar(255) NULL AFTER `cafeteriaID`, ADD COLUMN `price_student_cents` bigint NULL'
  )
);
PREPARE dishes_metadata_stmt FROM @dishes_metadata_alter;
EXECUTE dishes_metadata_stmt;
DEALLOCATE PREPARE dishes_metadata_stmt;
-- Create "push_delivery_logs" table
CREATE TABLE IF NOT EXISTS `push_delivery_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `device_id` varchar(191) NOT NULL,
  `dedup_key` varchar(191) NOT NULL,
  `sent_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_push_delivery_dedup` (`device_id`, `dedup_key`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "push_devices" table
CREATE TABLE IF NOT EXISTS `push_devices` (
  `device_id` varchar(191) NOT NULL,
  `device_type` enum('ios','android','windows') NOT NULL,
  `fcm_token` varchar(512) NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`device_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "service_status_states" table
CREATE TABLE IF NOT EXISTS `service_status_states` (
  `service_key` varchar(32) NOT NULL,
  `status` enum('ok','degraded','maintenance') NOT NULL DEFAULT "ok",
  `message_de` varchar(512) NOT NULL DEFAULT "",
  `message_en` varchar(512) NOT NULL DEFAULT "",
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`service_key`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "push_preferences" table
CREATE TABLE IF NOT EXISTS `push_preferences` (
  `device_id` varchar(191) NOT NULL,
  `canteen_ids` text NOT NULL,
  `favorite_dishes` text NOT NULL,
  `dietary_labels` text NOT NULL,
  `notify_favorite_dishes` bool NOT NULL DEFAULT 0,
  `notify_it_status` bool NOT NULL DEFAULT 0,
  `notify_tu_film` bool NOT NULL DEFAULT 0,
  `language` enum('de','en') NOT NULL DEFAULT "en",
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`device_id`),
  CONSTRAINT `fk_push_devices_preference` FOREIGN KEY (`device_id`) REFERENCES `push_devices` (`device_id`) ON UPDATE CASCADE ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
