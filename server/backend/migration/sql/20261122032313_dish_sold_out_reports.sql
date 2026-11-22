-- Create "dish_sold_out_reports" table
CREATE TABLE IF NOT EXISTS `dish_sold_out_reports` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `dish_id` bigint NOT NULL,
  `device_hash` varchar(64) NOT NULL,
  `report_date` char(10) NOT NULL,
  `reported_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_dish_sold_out_lookup` (`dish_id`, `report_date`),
  UNIQUE INDEX `idx_dish_sold_out_unique` (`dish_id`, `device_hash`, `report_date`),
  CONSTRAINT `fk_dish_sold_out_reports_dish` FOREIGN KEY (`dish_id`) REFERENCES `dishes` (`dish`) ON UPDATE CASCADE ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
