-- Modify "canteen_ratings" table (retry-safe for interrupted MySQL DDL)
SET @canteen_ratings_alter = (
  SELECT IF(
    EXISTS(
      SELECT 1 FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'canteen_ratings' AND COLUMN_NAME = 'device_hash'
    ),
    'SELECT 1',
    'ALTER TABLE `canteen_ratings` ADD COLUMN `device_hash` varchar(64) NULL AFTER `image`, ADD COLUMN `status` enum(\'visible\',\'hidden\') NOT NULL DEFAULT "visible" AFTER `device_hash`, ADD COLUMN `updated_at` timestamp NULL, ADD INDEX `idx_canteen_ratings_device_ts` (`device_hash`, `cafeteriaID`, `timestamp`), ADD INDEX `idx_canteen_ratings_status` (`status`)'
  )
);
PREPARE canteen_ratings_stmt FROM @canteen_ratings_alter;
EXECUTE canteen_ratings_stmt;
DEALLOCATE PREPARE canteen_ratings_stmt;
-- Modify "dish_ratings" table (retry-safe for interrupted MySQL DDL)
SET @dish_ratings_alter = (
  SELECT IF(
    EXISTS(
      SELECT 1 FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'dish_ratings' AND COLUMN_NAME = 'device_hash'
    ),
    'SELECT 1',
    'ALTER TABLE `dish_ratings` ADD COLUMN `device_hash` varchar(64) NULL AFTER `image`, ADD COLUMN `status` enum(\'visible\',\'hidden\') NOT NULL DEFAULT "visible" AFTER `device_hash`, ADD COLUMN `updated_at` timestamp NULL, ADD INDEX `idx_dish_ratings_device_dish_ts` (`device_hash`, `dishID`, `timestamp`), ADD INDEX `idx_dish_ratings_status` (`status`)'
  )
);
PREPARE dish_ratings_stmt FROM @dish_ratings_alter;
EXECUTE dish_ratings_stmt;
DEALLOCATE PREPARE dish_ratings_stmt;
-- Create "dish_rating_reports" table
CREATE TABLE IF NOT EXISTS `dish_rating_reports` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `rating_id` bigint NOT NULL,
  `device_hash` varchar(64) NOT NULL,
  `reason` varchar(256) NOT NULL DEFAULT "",
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_dish_rating_reports_rating_device` (`rating_id`, `device_hash`),
  CONSTRAINT `fk_dish_rating_reports_rating` FOREIGN KEY (`rating_id`) REFERENCES `dish_ratings` (`dishRating`) ON UPDATE CASCADE ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Update rating statistic views to exclude hidden ratings and include counts/histograms
CREATE OR REPLACE VIEW `cafeteria_rating_statistics` AS
SELECT cafeteriaID,
       AVG(points) AS average,
       MIN(points) AS min,
       MAX(points) AS max,
       STD(points) AS std,
       COUNT(*) AS rating_count,
       SUM(CASE WHEN points = 1 THEN 1 ELSE 0 END) AS stars_1,
       SUM(CASE WHEN points = 2 THEN 1 ELSE 0 END) AS stars_2,
       SUM(CASE WHEN points = 3 THEN 1 ELSE 0 END) AS stars_3,
       SUM(CASE WHEN points = 4 THEN 1 ELSE 0 END) AS stars_4,
       SUM(CASE WHEN points = 5 THEN 1 ELSE 0 END) AS stars_5
FROM canteen_ratings
WHERE status = 'visible'
GROUP BY cafeteriaID
ORDER BY COUNT(*) DESC, average DESC;
CREATE OR REPLACE VIEW `dish_rating_statistics` AS
SELECT d.cafeteriaID AS cafeteriaID,
       r.dishID      AS dishID,
       AVG(r.points) AS average,
       MAX(r.points) AS max,
       MIN(r.points) AS min,
       STD(r.points) AS std,
       COUNT(*)      AS rating_count,
       SUM(CASE WHEN r.points = 1 THEN 1 ELSE 0 END) AS stars_1,
       SUM(CASE WHEN r.points = 2 THEN 1 ELSE 0 END) AS stars_2,
       SUM(CASE WHEN r.points = 3 THEN 1 ELSE 0 END) AS stars_3,
       SUM(CASE WHEN r.points = 4 THEN 1 ELSE 0 END) AS stars_4,
       SUM(CASE WHEN r.points = 5 THEN 1 ELSE 0 END) AS stars_5
FROM dishes d
JOIN dish_ratings r ON r.dishID = d.dish
WHERE r.status = 'visible'
GROUP BY d.cafeteriaID, r.dishID
ORDER BY COUNT(0) DESC, AVG(r.points) DESC;
CREATE OR REPLACE VIEW `dish_rating_tag_statistics` AS
SELECT d.cafeteriaID    AS cafeteriaID,
       mrt.tagID        AS tagID,
       mr.dishID        AS dishID,
       AVG(mrt.points)  AS average,
       MAX(mrt.points)  AS max,
       MIN(mrt.points)  AS min,
       STD(mrt.points)  AS std
FROM dishes d
JOIN dish_ratings mr ON d.dish = mr.dishID
JOIN dish_rating_tags mrt ON mr.dishRating = mrt.parentRating
WHERE mr.status = 'visible'
GROUP BY d.cafeteriaID, mrt.tagID, mr.dishID;
CREATE OR REPLACE VIEW `cafeteria_rating_tag_statistics` AS
SELECT cr.cafeteriaID  AS cafeteriaID,
       crt.tagID       AS tagID,
       AVG(crt.points) AS average,
       MAX(crt.points) AS max,
       MIN(crt.points) AS min,
       STD(crt.points) AS std
FROM canteen_ratings cr
JOIN canteen_rating_tags crt ON cr.cafeteriaRating = crt.correspondingRating
WHERE cr.status = 'visible'
GROUP BY cr.cafeteriaID, crt.tagID;
CREATE OR REPLACE VIEW `dish_name_tag_statistics` AS
SELECT d.cafeteriaID   AS cafeteriaID,
       mnt.tagNameID   AS tagID,
       AVG(mnt.points) AS average,
       MAX(mnt.points) AS max,
       MIN(mnt.points) AS min,
       STD(mnt.points) AS std
FROM dishes d
JOIN dish_ratings mr ON d.dish = mr.dishID
JOIN dish_name_tags mnt ON mr.dishRating = mnt.correspondingRating
WHERE mr.status = 'visible'
GROUP BY d.cafeteriaID, mnt.tagNameID;
