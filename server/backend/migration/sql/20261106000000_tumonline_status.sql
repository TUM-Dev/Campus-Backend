-- Modify "crontab" table
ALTER TABLE `crontab` MODIFY COLUMN `type` enum('news','mensa','ticketsale','fileDownload','canteenHeadCount','iosNotifications','iosActivityReset','newExamResultsHook','movie','feedbackEmail','dishNameDownload','scrapeStudentClubs','tumonlineStatus') NULL;
-- Schedule TUMonline status synchronization every five minutes (publishing is opt-in).
INSERT INTO `crontab` (`interval`, `lastRun`, `type`) VALUES (300, 0, 'tumonlineStatus');
