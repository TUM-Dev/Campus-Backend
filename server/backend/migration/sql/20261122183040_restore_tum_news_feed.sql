-- 1. Restore broken or empty existing news/event sources in production.
-- Source 1: TUM News (DE)
UPDATE `news_sources`
SET `url` = 'https://www.tum.de/news.rss'
WHERE `source` = 1
  AND (`url` IS NULL OR `url` = '' OR CAST(`url` AS BINARY) = 'http://www.tum.de/nc/die-tum/aktuelles/?type=100');

-- Source 5: TUM Pressenews -> TUM News (EN)
UPDATE `news_sources`
SET `title` = 'TUM News (EN)',
    `url` = 'https://www.tum.de/en/news.rss'
WHERE `source` = 5
  AND (`url` IS NULL OR `url` = '');

-- Source 14: TUM Alumni Events -> TUM Events (DE)
UPDATE `news_sources`
SET `title` = 'TUM Events',
    `url` = 'https://www.tum.de/events.rss'
WHERE `source` = 14
  AND (`url` IS NULL OR `url` = '');

-- Source 18: TUM Community -> TUM Community Events
UPDATE `news_sources`
SET `url` = 'https://www.community.tum.de/en/events/'
WHERE `source` = 18
  AND (`url` IS NULL OR `url` = '' OR CAST(`url` AS BINARY) = 'https://www.community.tum.de/en/?event=true&eventfeed=true');

-- 2. Insert additional general TUM, SV Munich, Heilbronn (DE/EN), and School feeds
-- when a news source icon already exists in the database (production / seeded environments).
INSERT INTO `news_sources` (`title`, `url`, `icon`, `hook`)
SELECT src.`title`,
       src.`url`,
       icon_ref.`icon`,
       NULL
FROM (
  SELECT 'TUM News' AS `title`, 'https://www.tum.de/news.rss' AS `url`
  UNION ALL SELECT 'TUM News (EN)', 'https://www.tum.de/en/news.rss'
  UNION ALL SELECT 'TUM Events', 'https://www.tum.de/events.rss'
  UNION ALL SELECT 'TUM Events (EN)', 'https://www.tum.de/en/events.rss'
  UNION ALL SELECT 'TUM Community', 'https://www.community.tum.de/en/events/'
  UNION ALL SELECT 'Studentische Vertretung (SV)', 'https://www.sv.tum.de/sv/nachrichtenarchiv/rss.xml'
  UNION ALL SELECT 'Student Representation (SV EN)', 'https://www.sv.tum.de/en/sv/news-archive/rss.xml'
  UNION ALL SELECT 'TUM Campus Heilbronn Events', 'https://chn.tum.de/de/news-events/events'
  UNION ALL SELECT 'TUM Campus Heilbronn Events (EN)', 'https://chn.tum.de/news-events/events'
  UNION ALL SELECT 'TUM Campus Heilbronn News', 'https://chn.tum.de/de/news-events/newsroom'
  UNION ALL SELECT 'TUM Campus Heilbronn News (EN)', 'https://chn.tum.de/news-events/newsroom'
  UNION ALL SELECT 'TUM School of CIT', 'https://www.cit.tum.de/cit/aktuelles/rss.xml'
  UNION ALL SELECT 'TUM School of CIT (EN)', 'https://www.cit.tum.de/en/cit/news/rss.xml'
  UNION ALL SELECT 'TUM School of Natural Sciences', 'https://www.nat.tum.de/nat/aktuelles/rss.xml'
  UNION ALL SELECT 'TUM School of Life Sciences', 'https://www.ls.tum.de/ls/aktuelles/rss.xml'
) AS src
JOIN (
  SELECT COALESCE(
    (SELECT `icon` FROM `news_sources` WHERE `source` = 1 LIMIT 1),
    (SELECT `icon` FROM `news_sources` ORDER BY `source` LIMIT 1),
    (SELECT `file` FROM `files` WHERE `path` = 'news/sources/' ORDER BY `file` LIMIT 1)
  ) AS `icon`
) AS icon_ref ON icon_ref.`icon` IS NOT NULL
WHERE NOT EXISTS (
  SELECT 1 FROM `news_sources` existing WHERE CAST(existing.`url` AS BINARY) = CAST(src.`url` AS BINARY)
);

-- 3. Remove orphaned news cron jobs whose target news_source row no longer exists.
DELETE c
FROM `crontab` AS c
LEFT JOIN `news_sources` AS s ON s.`source` = c.`id`
WHERE c.`type` = 'news'
  AND (c.`id` IS NULL OR s.`source` IS NULL);

-- 4. Ensure every active news_source with a URL has a scheduled 12-hour news cron job.
INSERT INTO `crontab` (`interval`, `lastRun`, `type`, `id`)
SELECT 43200, 0, 'news', s.`source`
FROM `news_sources` AS s
WHERE s.`url` IS NOT NULL
  AND s.`url` != ''
  AND NOT EXISTS (
    SELECT 1 FROM `crontab` AS c WHERE c.`type` = 'news' AND c.`id` = s.`source`
  );

-- 5. Make newly restored/added news jobs due on the next scheduler poll.
UPDATE `crontab` AS c
JOIN `news_sources` AS s ON s.`source` = c.`id`
SET c.`lastRun` = 0
WHERE c.`type` = 'news'
  AND c.`interval` > 0
  AND s.`url` IN (
    'https://www.tum.de/news.rss',
    'https://www.tum.de/en/news.rss',
    'https://www.tum.de/events.rss',
    'https://www.tum.de/en/events.rss',
    'https://www.community.tum.de/en/events/',
    'https://www.sv.tum.de/sv/nachrichtenarchiv/rss.xml',
    'https://www.sv.tum.de/en/sv/news-archive/rss.xml',
    'https://chn.tum.de/de/news-events/events',
    'https://chn.tum.de/news-events/events',
    'https://chn.tum.de/de/news-events/newsroom',
    'https://chn.tum.de/news-events/newsroom',
    'https://www.cit.tum.de/cit/aktuelles/rss.xml',
    'https://www.cit.tum.de/en/cit/news/rss.xml',
    'https://www.nat.tum.de/nat/aktuelles/rss.xml',
    'https://www.ls.tum.de/ls/aktuelles/rss.xml'
  );
