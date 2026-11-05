-- Modify "crontab" table
ALTER TABLE `crontab` MODIFY COLUMN `type` enum('news','mensa','ticketsale','fileDownload','canteenHeadCount','iosNotifications','iosActivityReset','newExamResultsHook','movie','feedbackEmail','dishNameDownload','scrapeStudentClubs') NULL;
-- Create "student_council_collections" table
CREATE TABLE `student_council_collections` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `name` varchar(100) NULL,
  `language` enum('German','English') NULL DEFAULT "German",
  `description` longtext NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_student_council_collections_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "student_councils" table
CREATE TABLE `student_councils` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `name` longtext NULL,
  `language` enum('German','English') NULL DEFAULT "German",
  `description` longtext NULL,
  `link_url` varchar(190) NULL,
  `image_id` bigint NULL,
  `image_caption` longtext NULL,
  `student_council_collection_id` bigint unsigned NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_student_councils_image` (`image_id`),
  INDEX `fk_student_councils_student_council_collection` (`student_council_collection_id`),
  INDEX `idx_student_councils_deleted_at` (`deleted_at`),
  UNIQUE INDEX `uni_student_councils_link_url` (`language`, `link_url`),
  CONSTRAINT `fk_student_councils_image` FOREIGN KEY (`image_id`) REFERENCES `files` (`file`) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT `fk_student_councils_student_council_collection` FOREIGN KEY (`student_council_collection_id`) REFERENCES `student_council_collections` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
