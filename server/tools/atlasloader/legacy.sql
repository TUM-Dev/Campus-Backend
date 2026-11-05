-- Tables that exist in production (created by the legacy gormigrate migrations)
-- but have no GORM model. They are part of the desired schema so that Atlas does
-- not generate DROP statements for them.
--
-- To drop or change one of these tables, edit this file and run
-- `atlas migrate diff --env gorm <name>` to generate the migration.
--
-- This file is loaded after the GORM models, so it can reference model tables.

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
-- Create "curricula" table
CREATE TABLE `curricula` (
  `curriculum` bigint NOT NULL AUTO_INCREMENT,
  `category` enum('bachelor','master') NOT NULL DEFAULT "bachelor",
  `name` mediumtext NOT NULL,
  `url` mediumtext NOT NULL,
  PRIMARY KEY (`curriculum`)
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
  `latitude` float(10,6) NOT NULL DEFAULT 0.000000,
  `longitude` float(10,6) NOT NULL DEFAULT 0.000000,
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
-- Create "location" table
CREATE TABLE `location` (
  `location` bigint NOT NULL AUTO_INCREMENT,
  `name` text CHARSET utf8mb4 NULL COLLATE utf8mb4_unicode_ci,
  `lon` float(10,6) NOT NULL,
  `lat` float(10,6) NOT NULL,
  `radius` int NOT NULL DEFAULT 1000 COMMENT "in meters",
  PRIMARY KEY (`location`)
) CHARSET utf8mb3 COLLATE utf8mb3_general_ci AUTO_INCREMENT 2;
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

-- Foreign keys on model tables whose production names GORM cannot generate.
-- The corresponding GORM relationships are tagged `constraint:-`.
ALTER TABLE `canteen_ratings` ADD CONSTRAINT `cafeteria_rating_cafeteria_cafeteria_fk` FOREIGN KEY (`cafeteriaID`) REFERENCES `canteens` (`cafeteria`) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE `devices` ADD CONSTRAINT `devices_ibfk_1` FOREIGN KEY (`member`) REFERENCES `member` (`member`) ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE `dishes` ADD CONSTRAINT `dish_cafeteria_cafeteria_fk` FOREIGN KEY (`cafeteriaID`) REFERENCES `canteens` (`cafeteria`) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE `dish_ratings` ADD CONSTRAINT `dish_rating_dish_dish_fk` FOREIGN KEY (`dishID`) REFERENCES `dishes` (`dish`) ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE `dish_rating_tags` ADD CONSTRAINT `dish_rating_tag_dish_rating_dishRating_fk` FOREIGN KEY (`parentRating`) REFERENCES `dish_ratings` (`dishRating`) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE `dish_rating_tags` ADD CONSTRAINT `dish_rating_tag_dish_rating_tag_option_dishRatingTagOption_fk` FOREIGN KEY (`tagID`) REFERENCES `dish_rating_tag_options` (`dishRatingTagOption`) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE `dish_to_dish_name_tags` ADD CONSTRAINT `dish_to_dish_name_tag_dish_dish_fk` FOREIGN KEY (`dishID`) REFERENCES `dishes` (`dish`) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE `dishes_of_the_weeks` ADD CONSTRAINT `dishes_of_the_week_dish_dish_fk` FOREIGN KEY (`dishID`) REFERENCES `dishes` (`dish`) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE `movies` ADD CONSTRAINT `movies_ibfk_1` FOREIGN KEY (`cover`) REFERENCES `files` (`file`) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE `news_sources` ADD CONSTRAINT `news_sources_ibfk_1` FOREIGN KEY (`icon`) REFERENCES `files` (`file`) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE `news` ADD CONSTRAINT `news_ibfk_1` FOREIGN KEY (`src`) REFERENCES `news_sources` (`source`) ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE `news` ADD CONSTRAINT `news_ibfk_2` FOREIGN KEY (`file`) REFERENCES `files` (`file`) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE `news_alerts` ADD CONSTRAINT `news_alert_files_file_fk` FOREIGN KEY (`file`) REFERENCES `files` (`file`) ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE `notification` ADD CONSTRAINT `notification_ibfk_1` FOREIGN KEY (`type`) REFERENCES `notification_type` (`type`) ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE `notification` ADD CONSTRAINT `notification_ibfk_2` FOREIGN KEY (`location`) REFERENCES `location` (`location`) ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE `notification_confirmation` ADD CONSTRAINT `notification_confirmation_ibfk_1` FOREIGN KEY (`notification`) REFERENCES `notification` (`notification`) ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE `notification_confirmation` ADD CONSTRAINT `notification_confirmation_ibfk_2` FOREIGN KEY (`device`) REFERENCES `devices` (`device`) ON UPDATE NO ACTION ON DELETE NO ACTION;
