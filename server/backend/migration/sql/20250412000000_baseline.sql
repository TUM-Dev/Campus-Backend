-- Create "alarm_ban" table
CREATE TABLE `alarm_ban` (
  `ban` bigint NOT NULL AUTO_INCREMENT,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `ip` binary(16) NOT NULL,
  PRIMARY KEY (`ban`),
  UNIQUE INDEX `ip` (`ip`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "alarm_log" table
CREATE TABLE `alarm_log` (
  `alarm` bigint NOT NULL AUTO_INCREMENT,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `message` text NOT NULL,
  `send` int NOT NULL,
  `received` int NOT NULL,
  `test` enum('true','false') NOT NULL DEFAULT "false",
  `ip` binary(16) NOT NULL,
  PRIMARY KEY (`alarm`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "barrierFree_moreInfo" table
CREATE TABLE `barrierFree_moreInfo` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `title` varchar(32) CHARSET utf8mb4 NULL COLLATE utf8mb4_unicode_ci,
  `category` varchar(11) CHARSET utf8mb4 NULL COLLATE utf8mb4_unicode_ci,
  `url` varchar(128) CHARSET utf8mb4 NULL COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`id`)
) CHARSET utf8mb3 COLLATE utf8mb3_general_ci;
-- Create "barrierFree_persons" table
CREATE TABLE `barrierFree_persons` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(40) CHARSET utf8mb4 NULL COLLATE utf8mb4_unicode_ci,
  `telephone` varchar(32) CHARSET utf8mb4 NULL COLLATE utf8mb4_unicode_ci,
  `email` varchar(32) CHARSET utf8mb4 NULL COLLATE utf8mb4_unicode_ci,
  `faculty` varchar(32) CHARSET utf8mb4 NULL COLLATE utf8mb4_unicode_ci,
  `office` varchar(16) CHARSET utf8mb4 NULL COLLATE utf8mb4_unicode_ci,
  `officeHour` varchar(16) CHARSET utf8mb4 NULL COLLATE utf8mb4_unicode_ci,
  `tumID` varchar(24) CHARSET utf8mb4 NULL COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`id`)
) CHARSET utf8mb3 COLLATE utf8mb3_general_ci;
-- Create "canteen_head_counts" table
CREATE TABLE `canteen_head_counts` (
  `canteen_id` varchar(64) NOT NULL,
  `count` int NOT NULL,
  `max_count` int NOT NULL,
  `percent` float NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`canteen_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "canteen_rating_tag_options" table
CREATE TABLE `canteen_rating_tag_options` (
  `cafeteriaRatingTagOption` bigint NOT NULL AUTO_INCREMENT,
  `DE` varchar(191) NOT NULL DEFAULT (_utf8mb4'de'),
  `EN` varchar(191) NOT NULL DEFAULT (_utf8mb4'en'),
  PRIMARY KEY (`cafeteriaRatingTagOption`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "canteen_rating_tags" table
CREATE TABLE `canteen_rating_tags` (
  `CafeteriaRatingTag` bigint NOT NULL AUTO_INCREMENT,
  `correspondingRating` bigint NOT NULL,
  `tagID` bigint NOT NULL,
  `points` int NOT NULL,
  PRIMARY KEY (`CafeteriaRatingTag`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "canteens" table
CREATE TABLE `canteens` (
  `cafeteria` bigint NOT NULL AUTO_INCREMENT,
  `name` text NOT NULL,
  `address` text NOT NULL,
  `latitude` double NOT NULL,
  `longitude` double NOT NULL,
  PRIMARY KEY (`cafeteria`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "crontab" table
CREATE TABLE `crontab` (
  `cron` bigint NOT NULL AUTO_INCREMENT,
  `interval` int NOT NULL DEFAULT 7200,
  `lastRun` int NOT NULL DEFAULT 0,
  `type` enum('news','mensa','ticketsale','alarm','fileDownload','canteenHeadCount','iosNotifications','iosActivityReset','newExamResultsHook','movie','feedbackEmail','dishNameDownload','scrapeStudentClubs') NULL,
  `id` bigint NULL,
  PRIMARY KEY (`cron`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "curricula" table
CREATE TABLE `curricula` (
  `curriculum` bigint NOT NULL AUTO_INCREMENT,
  `category` enum('bachelor','master') NOT NULL DEFAULT "bachelor",
  `name` mediumtext NOT NULL,
  `url` mediumtext NOT NULL,
  PRIMARY KEY (`curriculum`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "dish_name_tag_options" table
CREATE TABLE `dish_name_tag_options` (
  `dishNameTagOption` bigint NOT NULL AUTO_INCREMENT,
  `DE` text NOT NULL,
  `EN` text NOT NULL,
  PRIMARY KEY (`dishNameTagOption`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "dish_name_tags" table
CREATE TABLE `dish_name_tags` (
  `DishNameTag` bigint NOT NULL AUTO_INCREMENT,
  `correspondingRating` bigint NOT NULL,
  `tagNameID` bigint NOT NULL,
  `points` int NOT NULL,
  PRIMARY KEY (`DishNameTag`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "excluded_dish_name_tag_options" table
CREATE TABLE `excluded_dish_name_tag_options` (
  `dishNameTagOptionExcluded` bigint NOT NULL AUTO_INCREMENT,
  `nameTagID` bigint NOT NULL,
  `expression` text NULL,
  PRIMARY KEY (`dishNameTagOptionExcluded`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "feedback" table
CREATE TABLE `feedback` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `image_count` int NOT NULL,
  `email_id` text NULL,
  `receiver` text NULL,
  `reply_to_email` text NULL,
  `reply_to_name` text NULL,
  `feedback` text NULL,
  `latitude` double NULL,
  `longitude` double NULL,
  `os_version` text NULL,
  `app_version` text NULL,
  `processed` bool NOT NULL DEFAULT 0,
  `timestamp` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `receiver_reply_to_feedback_app_version_uindex` (`receiver` (255), `reply_to_email` (100), `feedback` (255), `app_version` (100))
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "included_dish_name_tag_options" table
CREATE TABLE `included_dish_name_tag_options` (
  `dishNameTagOptionIncluded` bigint NOT NULL AUTO_INCREMENT,
  `nameTagID` bigint NOT NULL,
  `expression` text NULL,
  PRIMARY KEY (`dishNameTagOptionIncluded`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "mensaplan_mensa" table
CREATE TABLE `mensaplan_mensa` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `latitude` double NULL,
  `longitude` double NULL,
  `webid` int NULL,
  `category` varchar(50) NOT NULL,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "mensaprices" table
CREATE TABLE `mensaprices` (
  `price` bigint NOT NULL AUTO_INCREMENT,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `person` text NOT NULL,
  `type` text NOT NULL,
  `typeLong` text NOT NULL,
  `typeNumber` int NOT NULL,
  `value` decimal(10) NOT NULL,
  PRIMARY KEY (`price`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "new_exam_results_subscribers" table
CREATE TABLE `new_exam_results_subscribers` (
  `callback_url` varchar(191) NOT NULL,
  `api_key` longtext NULL,
  `created_at` datetime(3) NULL,
  `last_notified_at` datetime(3) NULL,
  PRIMARY KEY (`callback_url`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "openinghours" table
CREATE TABLE `openinghours` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `category` varchar(20) NOT NULL,
  `name` varchar(60) NOT NULL,
  `address` varchar(140) NOT NULL,
  `room` varchar(140) NULL,
  `transport_station` varchar(150) NULL,
  `opening_hours` varchar(300) NULL,
  `infos` varchar(500) NULL,
  `url` varchar(300) NOT NULL,
  `language` varchar(2) NULL DEFAULT "de",
  `reference_id` int NULL DEFAULT -1,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "published_exam_results" table
CREATE TABLE `published_exam_results` (
  `date` datetime(3) NULL,
  `exam_id` varchar(191) NOT NULL,
  `lecture_title` longtext NULL,
  `lecture_type` longtext NULL,
  `lecture_sem` longtext NULL,
  `published` bool NULL,
  PRIMARY KEY (`exam_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "sessions" table
CREATE TABLE `sessions` (
  `session` varchar(255) NOT NULL,
  `access` int unsigned NULL,
  `data` text NULL,
  PRIMARY KEY (`session`),
  UNIQUE INDEX `session` (`session`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "update_notes" table
CREATE TABLE `update_notes` (
  `version_code` bigint NOT NULL AUTO_INCREMENT,
  `version_name` text NULL,
  `message` text NULL,
  PRIMARY KEY (`version_code`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "wifi_measurement" table
CREATE TABLE `wifi_measurement` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `date` date NOT NULL,
  `SSID` varchar(32) NOT NULL,
  `BSSID` varchar(64) NOT NULL,
  `dBm` int NULL,
  `accuracyInMeters` float NOT NULL,
  `latitude` double NOT NULL,
  `longitude` double NOT NULL,
  PRIMARY KEY (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "canteen_ratings" table
CREATE TABLE `canteen_ratings` (
  `cafeteriaRating` bigint NOT NULL AUTO_INCREMENT,
  `points` int NOT NULL,
  `comment` text NULL,
  `cafeteriaID` bigint NOT NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `image` text NULL,
  PRIMARY KEY (`cafeteriaRating`),
  INDEX `cafeteria_rating_cafeteria_cafeteria_fk` (`cafeteriaID`),
  CONSTRAINT `cafeteria_rating_cafeteria_cafeteria_fk` FOREIGN KEY (`cafeteriaID`) REFERENCES `canteens` (`cafeteria`) ON UPDATE CASCADE ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "member" table
CREATE TABLE `member` (
  `member` bigint NOT NULL AUTO_INCREMENT,
  `lrz_id` varchar(7) NOT NULL,
  `name` varchar(150) NOT NULL,
  `active_day` int NULL DEFAULT 0,
  `active_day_date` date NULL,
  `student_id` text NULL,
  `employee_id` text NULL,
  `external_id` text NULL,
  PRIMARY KEY (`member`),
  UNIQUE INDEX `lrz_id` (`lrz_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci AUTO_INCREMENT 104353;
-- Create "devices" table
CREATE TABLE `devices` (
  `device` bigint NOT NULL AUTO_INCREMENT,
  `member` bigint NULL,
  `uuid` varchar(50) NOT NULL,
  `created` timestamp NULL,
  `lastAccess` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `lastApi` text NOT NULL DEFAULT (_utf8mb4''),
  `developer` enum('true','false') NOT NULL DEFAULT "false",
  `osVersion` text NOT NULL DEFAULT (_utf8mb4''),
  `appVersion` text NOT NULL DEFAULT (_utf8mb4''),
  `counter` int NOT NULL DEFAULT 0,
  `pk` text NULL,
  `pkActive` enum('true','false') NOT NULL DEFAULT "false",
  `gcmToken` text NULL,
  `gcmStatus` varchar(200) NULL,
  `confirmationKey` varchar(35) NULL,
  `keyCreated` datetime NULL,
  `keyConfirmed` datetime NULL,
  PRIMARY KEY (`device`),
  INDEX `member` (`member`),
  UNIQUE INDEX `uuid` (`uuid`),
  CONSTRAINT `devices_ibfk_1` FOREIGN KEY (`member`) REFERENCES `member` (`member`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "device2stats" table
CREATE TABLE `device2stats` (
  `device` bigint NOT NULL,
  `LecturesPersonalActivity` int NOT NULL DEFAULT 0,
  `CafeteriaActivity` int NOT NULL DEFAULT 0,
  `WizNavStartActivity` int NOT NULL DEFAULT 0,
  `NewsActivity` int NOT NULL DEFAULT 0,
  `StartupActivity` int NOT NULL DEFAULT 0,
  `MainActivity` int NOT NULL DEFAULT 0,
  `CalendarActivity` int NOT NULL DEFAULT 0,
  `WizNavCheckTokenActivity` int NOT NULL DEFAULT 0,
  `CurriculaActivity` int NOT NULL DEFAULT 0,
  `CurriculaDetailsActivity` int NOT NULL DEFAULT 0,
  `GradeChartActivity` int NOT NULL DEFAULT 0,
  `GradesActivity` int NOT NULL DEFAULT 0,
  `InformationActivity` int NOT NULL DEFAULT 0,
  `LecturesAppointmentsActivity` int NOT NULL DEFAULT 0,
  `LecturesDetailsActivity` int NOT NULL DEFAULT 0,
  `OpeningHoursDetailActivity` int NOT NULL DEFAULT 0,
  `OpeningHoursListActivity` int NOT NULL DEFAULT 0,
  `OrganisationActivity` int NOT NULL DEFAULT 0,
  `OrganisationDetailsActivity` int NOT NULL DEFAULT 0,
  `PersonsDetailsActivity` int NOT NULL DEFAULT 0,
  `PersonsSearchActivity` int NOT NULL DEFAULT 0,
  `PlansActivity` int NOT NULL DEFAULT 0,
  `PlansDetailsActivity` int NOT NULL DEFAULT 0,
  `SetupEduroamActivity` int NOT NULL DEFAULT 0,
  `TransportationActivity` int NOT NULL DEFAULT 0,
  `TransportationDetailsActivity` int NOT NULL DEFAULT 0,
  `TuitionFeesActivity` int NOT NULL DEFAULT 0,
  `UserPreferencesActivity` int NOT NULL DEFAULT 0,
  `WizNavExtrasActivity` int NOT NULL DEFAULT 0,
  `TuitionFeesCard` int NOT NULL DEFAULT 0,
  `NextLectureCard` int NOT NULL DEFAULT 0,
  `CafeteriaMenuCard` int NOT NULL DEFAULT 0,
  `NewsCard1` int NOT NULL DEFAULT 0,
  `NewsCard2` int NOT NULL DEFAULT 0,
  `NewsCard3` int NOT NULL DEFAULT 0,
  `NewsCard7` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`device`),
  CONSTRAINT `device2stats_ibfk_2` FOREIGN KEY (`device`) REFERENCES `devices` (`device`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "dishes" table
CREATE TABLE `dishes` (
  `dish` bigint NOT NULL AUTO_INCREMENT,
  `name` text NOT NULL,
  `type` text NOT NULL,
  `cafeteriaID` bigint NOT NULL,
  PRIMARY KEY (`dish`),
  INDEX `dish_cafeteria_cafeteria_fk` (`cafeteriaID`),
  UNIQUE INDEX `dish_name_cafeteriaID_uindex` (`name` (255), `cafeteriaID`),
  CONSTRAINT `dish_cafeteria_cafeteria_fk` FOREIGN KEY (`cafeteriaID`) REFERENCES `canteens` (`cafeteria`) ON UPDATE CASCADE ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "dishflags" table
CREATE TABLE `dishflags` (
  `flag` bigint NOT NULL AUTO_INCREMENT,
  `short` varchar(10) NOT NULL,
  `description` varchar(50) NOT NULL,
  PRIMARY KEY (`flag`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "dish2dishflags" table
CREATE TABLE `dish2dishflags` (
  `dish2dishflags` bigint NOT NULL AUTO_INCREMENT,
  `dish` bigint NOT NULL,
  `flag` bigint NOT NULL,
  PRIMARY KEY (`dish2dishflags`),
  UNIQUE INDEX `dish` (`dish`, `flag`),
  INDEX `flag` (`flag`),
  CONSTRAINT `dish2dishflags_dish_dish_fk` FOREIGN KEY (`dish`) REFERENCES `dishes` (`dish`) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT `dish2dishflags_dishflags_flag_fk` FOREIGN KEY (`flag`) REFERENCES `dishflags` (`flag`) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT `dish2dishflags_ibfk_1` FOREIGN KEY (`dish`) REFERENCES `dishes` (`dish`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `dish2dishflags_ibfk_2` FOREIGN KEY (`flag`) REFERENCES `dishflags` (`flag`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "mensa" table
CREATE TABLE `mensa` (
  `mensa` bigint NOT NULL AUTO_INCREMENT,
  `id` int NULL,
  `name` text NOT NULL,
  `address` text NOT NULL,
  `latitude` float NOT NULL DEFAULT 0.000000,
  `longitude` float NOT NULL DEFAULT 0.000000,
  PRIMARY KEY (`mensa`),
  UNIQUE INDEX `id` (`id`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci AUTO_INCREMENT 17;
-- Create "dish2mensa" table
CREATE TABLE `dish2mensa` (
  `dish2mensa` bigint NOT NULL AUTO_INCREMENT,
  `mensa` bigint NOT NULL,
  `dish` bigint NOT NULL,
  `date` date NOT NULL,
  `created` datetime NOT NULL,
  `modifierd` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`dish2mensa`),
  INDEX `dish` (`dish`),
  INDEX `mensa` (`mensa`),
  CONSTRAINT `dish2mensa_dish_dish_fk` FOREIGN KEY (`dish`) REFERENCES `dishes` (`dish`) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT `dish2mensa_ibfk_1` FOREIGN KEY (`mensa`) REFERENCES `mensa` (`mensa`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `dish2mensa_ibfk_2` FOREIGN KEY (`dish`) REFERENCES `dishes` (`dish`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `dish2mensa_mensa_mensa_fk` FOREIGN KEY (`mensa`) REFERENCES `mensa` (`mensa`) ON UPDATE CASCADE ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "dish_ratings" table
CREATE TABLE `dish_ratings` (
  `dishRating` bigint NOT NULL AUTO_INCREMENT,
  `points` int NOT NULL,
  `dishID` bigint NOT NULL,
  `comment` text NULL,
  `timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `image` text NULL,
  PRIMARY KEY (`dishRating`),
  INDEX `dish_rating_dish_dish_fk` (`dishID`),
  CONSTRAINT `dish_rating_dish_dish_fk` FOREIGN KEY (`dishID`) REFERENCES `dishes` (`dish`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "dish_rating_tag_options" table
CREATE TABLE `dish_rating_tag_options` (
  `dishRatingTagOption` bigint NOT NULL AUTO_INCREMENT,
  `DE` text NOT NULL DEFAULT (_utf8mb4'de'),
  `EN` text NOT NULL DEFAULT (_utf8mb4'en'),
  PRIMARY KEY (`dishRatingTagOption`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "dish_rating_tags" table
CREATE TABLE `dish_rating_tags` (
  `dishRatingTag` bigint NOT NULL AUTO_INCREMENT,
  `parentRating` bigint NOT NULL,
  `tagID` bigint NOT NULL,
  `points` int NOT NULL,
  PRIMARY KEY (`dishRatingTag`),
  INDEX `dish_rating_tag_dish_rating_dishRating_fk` (`parentRating`),
  INDEX `dish_rating_tag_dish_rating_tag_option_dishRatingTagOption_fk` (`tagID`),
  CONSTRAINT `dish_rating_tag_dish_rating_dishRating_fk` FOREIGN KEY (`parentRating`) REFERENCES `dish_ratings` (`dishRating`) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT `dish_rating_tag_dish_rating_tag_option_dishRatingTagOption_fk` FOREIGN KEY (`tagID`) REFERENCES `dish_rating_tag_options` (`dishRatingTagOption`) ON UPDATE CASCADE ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "dish_to_dish_name_tags" table
CREATE TABLE `dish_to_dish_name_tags` (
  `dishToDishNameTag` bigint NOT NULL AUTO_INCREMENT,
  `dishID` bigint NOT NULL,
  `nameTagID` bigint NOT NULL,
  PRIMARY KEY (`dishToDishNameTag`),
  INDEX `dish_to_dish_name_tag_dish_dish_fk` (`dishID`),
  CONSTRAINT `dish_to_dish_name_tag_dish_dish_fk` FOREIGN KEY (`dishID`) REFERENCES `dishes` (`dish`) ON UPDATE CASCADE ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "dishes_of_the_weeks" table
CREATE TABLE `dishes_of_the_weeks` (
  `dishesOfTheWeek` bigint NOT NULL AUTO_INCREMENT,
  `year` int NOT NULL,
  `week` int NOT NULL,
  `day` int NOT NULL,
  `dishID` bigint NOT NULL,
  PRIMARY KEY (`dishesOfTheWeek`),
  INDEX `dishes_of_the_week_dish_dish_fk` (`dishID`),
  CONSTRAINT `dishes_of_the_week_dish_dish_fk` FOREIGN KEY (`dishID`) REFERENCES `dishes` (`dish`) ON UPDATE CASCADE ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "actions" table
CREATE TABLE `actions` (
  `action` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `description` mediumtext NOT NULL,
  `color` varchar(6) NOT NULL,
  PRIMARY KEY (`action`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci AUTO_INCREMENT 19;
-- Create "users" table
CREATE TABLE `users` (
  `user` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(7) NULL,
  `firstname` varchar(100) NULL,
  `surname` varchar(100) NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted` int NOT NULL DEFAULT 0,
  `lastActive` int NOT NULL DEFAULT 0,
  `lastPage` text NULL,
  `lastLogin` datetime NULL,
  `tum_id_student` varchar(50) NULL COMMENT "OBFUSCATED_ID_ST",
  `tum_id_employee` varchar(50) NULL COMMENT "OBFUSCATED_ID_B",
  `tum_id_alumni` varchar(50) NULL COMMENT "OBFUSCATED_ID_EXT",
  `tum_id_preferred` varchar(50) NULL COMMENT "OBFUSCATED_ID_BEVORZUGT",
  PRIMARY KEY (`user`),
  UNIQUE INDEX `username` (`username`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci AUTO_INCREMENT 434;
-- Create "log" table
CREATE TABLE `log` (
  `log` bigint NOT NULL AUTO_INCREMENT,
  `time` int NOT NULL,
  `user_executed` bigint NULL,
  `user_affected` bigint NULL,
  `action` bigint NULL,
  `comment` varchar(255) NOT NULL,
  PRIMARY KEY (`log`),
  INDEX `action` (`action`),
  INDEX `user` (`user_executed`),
  INDEX `user_affected` (`user_affected`),
  CONSTRAINT `fkLog2Actions` FOREIGN KEY (`action`) REFERENCES `actions` (`action`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fkLog2UsersAf` FOREIGN KEY (`user_affected`) REFERENCES `users` (`user`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fkLog2UsersEx` FOREIGN KEY (`user_executed`) REFERENCES `users` (`user`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "rights" table
CREATE TABLE `rights` (
  `right` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NULL,
  `description` mediumtext NOT NULL,
  `category` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`right`),
  UNIQUE INDEX `Unquie` (`name`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci AUTO_INCREMENT 14;
-- Create "menu" table
CREATE TABLE `menu` (
  `menu` bigint NOT NULL AUTO_INCREMENT,
  `right` bigint NULL,
  `parent` bigint NULL,
  `name` varchar(255) NULL,
  `path` varchar(255) NULL,
  `target` enum('_blank','_self','_parent','_top') NOT NULL DEFAULT "_self",
  `icon` varchar(200) NOT NULL,
  `class` varchar(200) NOT NULL,
  `position` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`menu`),
  INDEX `parent` (`parent`),
  INDEX `right` (`right`),
  CONSTRAINT `menu_ibfk_1` FOREIGN KEY (`right`) REFERENCES `rights` (`right`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `menu_ibfk_2` FOREIGN KEY (`parent`) REFERENCES `menu` (`menu`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci AUTO_INCREMENT 25;
-- Create "modules" table
CREATE TABLE `modules` (
  `module` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NULL,
  `right` bigint NULL,
  PRIMARY KEY (`module`),
  INDEX `module_right` (`right`),
  CONSTRAINT `fkMod2Rights` FOREIGN KEY (`right`) REFERENCES `rights` (`right`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "files" table
CREATE TABLE `files` (
  `file` bigint NOT NULL AUTO_INCREMENT,
  `name` text NULL,
  `path` text NULL,
  `downloads` int NOT NULL DEFAULT 0,
  `url` varchar(191) NULL,
  `downloaded` bool NULL DEFAULT 1,
  PRIMARY KEY (`file`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "movies" table
CREATE TABLE `movies` (
  `kino` bigint NOT NULL AUTO_INCREMENT,
  `date` datetime NOT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `title` text NOT NULL,
  `year` varchar(4) NULL,
  `runtime` varchar(40) NULL,
  `genre` varchar(100) NULL,
  `director` text NULL,
  `actors` text NULL,
  `rating` varchar(4) NULL,
  `description` text NOT NULL,
  `trailer` text NULL,
  `cover` bigint NULL,
  `link` varchar(190) NOT NULL,
  `location` varchar(191) NULL,
  PRIMARY KEY (`kino`),
  INDEX `cover` (`cover`),
  UNIQUE INDEX `link` (`link`),
  CONSTRAINT `movies_ibfk_1` FOREIGN KEY (`cover`) REFERENCES `files` (`file`) ON UPDATE CASCADE ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "news_sources" table
CREATE TABLE `news_sources` (
  `source` bigint NOT NULL AUTO_INCREMENT,
  `title` text NOT NULL,
  `url` text NULL,
  `icon` bigint NOT NULL,
  `hook` enum('newspread','impulsivHook') NULL,
  PRIMARY KEY (`source`),
  INDEX `icon` (`icon`),
  CONSTRAINT `news_sources_ibfk_1` FOREIGN KEY (`icon`) REFERENCES `files` (`file`) ON UPDATE CASCADE ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "news" table
CREATE TABLE `news` (
  `news` bigint NOT NULL AUTO_INCREMENT,
  `date` datetime NOT NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `title` tinytext NOT NULL,
  `description` text NOT NULL,
  `src` bigint NOT NULL,
  `link` varchar(190) NOT NULL,
  `image` text NULL,
  `file` bigint NULL,
  PRIMARY KEY (`news`),
  INDEX `file` (`file`),
  UNIQUE INDEX `link` (`link`),
  INDEX `src` (`src`),
  CONSTRAINT `news_ibfk_1` FOREIGN KEY (`src`) REFERENCES `news_sources` (`source`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `news_ibfk_2` FOREIGN KEY (`file`) REFERENCES `files` (`file`) ON UPDATE CASCADE ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "news_alerts" table
CREATE TABLE `news_alerts` (
  `news_alert` bigint NOT NULL AUTO_INCREMENT,
  `file` bigint NOT NULL,
  `name` varchar(100) NULL,
  `link` text NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `from` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `to` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`news_alert`),
  INDEX `news_alert_files_file_fk` (`file`),
  CONSTRAINT `news_alert_files_file_fk` FOREIGN KEY (`file`) REFERENCES `files` (`file`) ON UPDATE NO ACTION ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "notification_type" table
CREATE TABLE `notification_type` (
  `type` bigint NOT NULL AUTO_INCREMENT,
  `name` text NULL,
  `confirmation` enum('true','false') NULL,
  PRIMARY KEY (`type`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "location" table
CREATE TABLE `location` (
  `location` bigint NOT NULL AUTO_INCREMENT,
  `name` text CHARSET utf8mb4 NULL COLLATE utf8mb4_unicode_ci,
  `lon` float NOT NULL,
  `lat` float NOT NULL,
  `radius` int NOT NULL DEFAULT 1000 COMMENT "in meters",
  PRIMARY KEY (`location`)
) CHARSET utf8mb3 COLLATE utf8mb3_general_ci AUTO_INCREMENT 2;
-- Create "notification" table
CREATE TABLE `notification` (
  `notification` bigint NOT NULL AUTO_INCREMENT,
  `type` bigint NOT NULL,
  `location` bigint NULL,
  `title` text NULL,
  `description` text NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `signature` text NULL,
  `silent` bool NOT NULL DEFAULT 0,
  PRIMARY KEY (`notification`),
  INDEX `location` (`location`),
  INDEX `type` (`type`),
  CONSTRAINT `notification_ibfk_1` FOREIGN KEY (`type`) REFERENCES `notification_type` (`type`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `notification_ibfk_2` FOREIGN KEY (`location`) REFERENCES `location` (`location`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "notification_confirmation" table
CREATE TABLE `notification_confirmation` (
  `notification` bigint NOT NULL,
  `device` bigint NOT NULL,
  `sent` bool NOT NULL DEFAULT 0,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `received` timestamp NULL,
  PRIMARY KEY (`notification`, `device`),
  INDEX `device` (`device`),
  CONSTRAINT `notification_confirmation_ibfk_1` FOREIGN KEY (`notification`) REFERENCES `notification` (`notification`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `notification_confirmation_ibfk_2` FOREIGN KEY (`device`) REFERENCES `devices` (`device`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "questionAnswers" table
CREATE TABLE `questionAnswers` (
  `answer` bigint NOT NULL AUTO_INCREMENT,
  `text` text NULL,
  PRIMARY KEY (`answer`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci AUTO_INCREMENT 3;
-- Create "question" table
CREATE TABLE `question` (
  `question` bigint NOT NULL AUTO_INCREMENT,
  `member` bigint NOT NULL,
  `text` text NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `end` timestamp NULL,
  PRIMARY KEY (`question`),
  INDEX `member` (`member`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci AUTO_INCREMENT 282;
-- Create "question2answer" table
CREATE TABLE `question2answer` (
  `question` bigint NOT NULL,
  `answer` bigint NOT NULL,
  `member` bigint NOT NULL,
  UNIQUE INDEX `question` (`question`, `member`),
  INDEX `question2answer_member_member_fk` (`member`),
  INDEX `question2answer_questionAnswers_answer_fk` (`answer`),
  CONSTRAINT `question2answer_member_member_fk` FOREIGN KEY (`member`) REFERENCES `member` (`member`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `question2answer_questionAnswers_answer_fk` FOREIGN KEY (`answer`) REFERENCES `questionAnswers` (`answer`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `question2answer_question_question_fk` FOREIGN KEY (`question`) REFERENCES `question` (`question`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "faculty" table
CREATE TABLE `faculty` (
  `faculty` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(150) NULL COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`faculty`),
  UNIQUE INDEX `name` (`name`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci AUTO_INCREMENT 18;
-- Create "question2faculty" table
CREATE TABLE `question2faculty` (
  `question` bigint NOT NULL,
  `faculty` bigint NOT NULL,
  PRIMARY KEY (`question`, `faculty`),
  INDEX `faculty` (`faculty`),
  CONSTRAINT `question2faculty_ibfk_1` FOREIGN KEY (`question`) REFERENCES `question` (`question`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `question2faculty_ibfk_2` FOREIGN KEY (`faculty`) REFERENCES `faculty` (`faculty`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "recover" table
CREATE TABLE `recover` (
  `recover` bigint NOT NULL AUTO_INCREMENT,
  `user` bigint NOT NULL,
  `created` int NOT NULL,
  `hash` varchar(190) NOT NULL,
  `ip` varchar(255) NOT NULL,
  PRIMARY KEY (`recover`),
  UNIQUE INDEX `hash` (`hash`),
  INDEX `user` (`user`),
  CONSTRAINT `fkRecover2User` FOREIGN KEY (`user`) REFERENCES `users` (`user`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "reports" table
CREATE TABLE `reports` (
  `report` bigint NOT NULL AUTO_INCREMENT,
  `device` bigint NULL,
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `fixed` enum('true','false') NOT NULL DEFAULT "false",
  `issue` int NULL,
  `stacktrace` mediumtext NOT NULL,
  `stacktraceGroup` text NULL,
  `log` mediumtext NOT NULL,
  `package` mediumtext NOT NULL,
  `packageVersion` mediumtext NOT NULL,
  `packageVersionCode` int NOT NULL DEFAULT -1,
  `model` mediumtext NOT NULL,
  `osVersion` mediumtext NOT NULL,
  `networkWifi` varchar(10) NOT NULL,
  `networkMobile` varchar(10) NOT NULL,
  `gps` varchar(10) NOT NULL,
  `screenWidth` varchar(100) NOT NULL,
  `screenHeight` varchar(100) NOT NULL,
  `screenOrientation` varchar(100) NOT NULL,
  `screenDpi` varchar(100) NOT NULL,
  PRIMARY KEY (`report`),
  INDEX `device` (`device`),
  CONSTRAINT `reports_ibfk_3` FOREIGN KEY (`device`) REFERENCES `devices` (`device`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "roles" table
CREATE TABLE `roles` (
  `role` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `description` mediumtext NOT NULL,
  PRIMARY KEY (`role`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci AUTO_INCREMENT 6;
-- Create "roles2rights" table
CREATE TABLE `roles2rights` (
  `role` bigint NOT NULL,
  `right` bigint NOT NULL,
  PRIMARY KEY (`role`, `right`),
  INDEX `fkRight_idx` (`right`),
  CONSTRAINT `fkRight` FOREIGN KEY (`right`) REFERENCES `rights` (`right`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fkRole` FOREIGN KEY (`role`) REFERENCES `roles` (`role`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "student_club_collections" table
CREATE TABLE `student_club_collections` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL,
  `updated_at` datetime(3) NULL,
  `deleted_at` datetime(3) NULL,
  `name` varchar(100) NULL,
  `language` enum('German','English') NULL DEFAULT "German",
  `description` longtext NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_student_club_collections_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "student_clubs" table
CREATE TABLE `student_clubs` (
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
  `student_club_collection_id` bigint unsigned NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_student_clubs_image` (`image_id`),
  INDEX `fk_student_clubs_student_club_collection` (`student_club_collection_id`),
  INDEX `idx_student_clubs_deleted_at` (`deleted_at`),
  UNIQUE INDEX `uni_student_clubs_link_url` (`link_url`, `language`),
  CONSTRAINT `fk_student_clubs_image` FOREIGN KEY (`image_id`) REFERENCES `files` (`file`) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT `fk_student_clubs_student_club_collection` FOREIGN KEY (`student_club_collection_id`) REFERENCES `student_club_collections` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "users2info" table
CREATE TABLE `users2info` (
  `user` bigint NOT NULL,
  `firstname` varchar(255) NOT NULL,
  `surname` varchar(255) NOT NULL,
  `lastPwChange` int NOT NULL,
  `pager` int NULL DEFAULT 15,
  PRIMARY KEY (`user`),
  CONSTRAINT `fkUsers` FOREIGN KEY (`user`) REFERENCES `users` (`user`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- Create "users2roles" table
CREATE TABLE `users2roles` (
  `user` bigint NOT NULL,
  `role` bigint NOT NULL,
  PRIMARY KEY (`user`, `role`),
  INDEX `fkUser2RolesRole` (`role`),
  CONSTRAINT `fkUser2RolesRole` FOREIGN KEY (`role`) REFERENCES `roles` (`role`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fkUser2RolesUser` FOREIGN KEY (`user`) REFERENCES `users` (`user`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;
