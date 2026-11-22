-- Trigger an immediate movie and news cron run on deploy so newly added cinemas and news feeds populate right away.
UPDATE `crontab`
SET `lastRun` = 0
WHERE `type` IN ('movie', 'news', 'fileDownload')
  AND `interval` > 0;
