-- Create "news_cache_versions" table
CREATE TABLE `news_cache_versions` (
  `id` tinyint unsigned NOT NULL,
  `version` bigint unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
INSERT INTO `news_cache_versions` (`id`, `version`) VALUES (1, 0);
