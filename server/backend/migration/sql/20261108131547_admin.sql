-- Create "admin_audits" table
CREATE TABLE `admin_audits` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) NOT NULL,
  `subject` varchar(255) NOT NULL,
  `action` varchar(32) NOT NULL,
  `source_id` bigint NOT NULL,
  `details` mediumtext NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "admin_sessions" table
CREATE TABLE `admin_sessions` (
  `token_hash` char(64) NOT NULL,
  `issuer` varchar(255) NOT NULL,
  `subject` varchar(255) NOT NULL,
  `csrf_token` char(64) NOT NULL,
  `expires_at` datetime(6) NOT NULL,
  PRIMARY KEY (`token_hash`),
  INDEX `idx_admin_sessions_expires_at` (`expires_at`)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
