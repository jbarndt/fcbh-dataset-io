-- MySQL dump 10.14  Distrib 5.5.68-MariaDB, for Linux (x86_64)
--
-- Host: rds.dev.biblebrain.com    Database: dbp_TEST
-- ------------------------------------------------------
-- Server version	8.0.39

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `access_group_filesets`
--

DROP TABLE IF EXISTS `access_group_filesets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `access_group_filesets` (
  `access_group_id` int unsigned NOT NULL,
  `hash_id` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`access_group_id`,`hash_id`),
  KEY `FK_access_group_filesets__hash_id` (`hash_id`),
  CONSTRAINT `FK_access_groups_access_group_filesets` FOREIGN KEY (`access_group_id`) REFERENCES `access_groups` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_bible_filesets_access_group_filesets` FOREIGN KEY (`hash_id`) REFERENCES `bible_filesets` (`hash_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `access_group_filesets_view`
--

DROP TABLE IF EXISTS `access_group_filesets_view`;
/*!50001 DROP VIEW IF EXISTS `access_group_filesets_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `access_group_filesets_view` (
  `fileset_id` tinyint NOT NULL,
  `asset_id` tinyint NOT NULL,
  `set_type_code` tinyint NOT NULL,
  `hash_id` tinyint NOT NULL,
  `access_group_id` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `access_group_types`
--

DROP TABLE IF EXISTS `access_group_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `access_group_types` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `access_group_id` int unsigned NOT NULL,
  `access_type_id` int unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `access_group_types_access_group_id_foreign` (`access_group_id`),
  KEY `access_group_types_access_type_id_foreign` (`access_type_id`),
  CONSTRAINT `FK_access_groups_access_group_types` FOREIGN KEY (`access_group_id`) REFERENCES `access_groups` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_access_types_access_group_types` FOREIGN KEY (`access_type_id`) REFERENCES `access_types` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `access_groups`
--

DROP TABLE IF EXISTS `access_groups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `access_groups` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `lpts_fieldname` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mode_id` tinyint unsigned NOT NULL DEFAULT '1',
  `display_order` tinyint unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `FK_bible_fileset_modes` (`mode_id`),
  KEY `access_groups_display_order_index` (`display_order`),
  CONSTRAINT `FK_bible_fileset_modes` FOREIGN KEY (`mode_id`) REFERENCES `bible_fileset_modes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2006 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `access_types`
--

DROP TABLE IF EXISTS `access_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `access_types` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(24) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `country_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `continent_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `allowed` tinyint(1) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `access_types_country_id_foreign` (`country_id`),
  CONSTRAINT `FK_countries_access_types` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `all_stocknumbers`
--

DROP TABLE IF EXISTS `all_stocknumbers`;
/*!50001 DROP VIEW IF EXISTS `all_stocknumbers`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `all_stocknumbers` (
  `stocknumber` tinyint NOT NULL,
  `bible_id` tinyint NOT NULL,
  `language` tinyint NOT NULL,
  `version` tinyint NOT NULL,
  `iso` tinyint NOT NULL,
  `rolv` tinyint NOT NULL,
  `language_id` tinyint NOT NULL,
  `country_id` tinyint NOT NULL,
  `fcbh_name` tinyint NOT NULL,
  `fcbh_language_id` tinyint NOT NULL,
  `country_name` tinyint NOT NULL,
  `heart_name` tinyint NOT NULL,
  `pseudonym` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `alphabet_fonts`
--

DROP TABLE IF EXISTS `alphabet_fonts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `alphabet_fonts` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `script_id` char(4) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `font_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `font_filename` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `font_weight` int unsigned DEFAULT NULL,
  `copyright` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `url` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `italic` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `alphabet_fonts_script_id_foreign` (`script_id`),
  CONSTRAINT `FK_alphabets_alphabet_fonts` FOREIGN KEY (`script_id`) REFERENCES `alphabets` (`script`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `alphabet_language`
--

DROP TABLE IF EXISTS `alphabet_language`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `alphabet_language` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `script_id` char(4) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `language_id` int unsigned NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `alphabet_language_language_id_foreign` (`language_id`),
  KEY `alphabet_language_script_id_index` (`script_id`),
  CONSTRAINT `FK_alphabets_alphabet_language` FOREIGN KEY (`script_id`) REFERENCES `alphabets` (`script`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_languages_alphabet_language` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=328 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `alphabet_numeral_systems`
--

DROP TABLE IF EXISTS `alphabet_numeral_systems`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `alphabet_numeral_systems` (
  `numeral_system_id` char(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `script_id` char(4) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`numeral_system_id`),
  KEY `alphabet_numeral_systems_script_id_foreign` (`script_id`),
  KEY `alphabet_numeral_systems_numeral_system_id_index` (`numeral_system_id`),
  CONSTRAINT `FK_alphabets_alphabet_numeral_systems` FOREIGN KEY (`script_id`) REFERENCES `alphabets` (`script`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `alphabets`
--

DROP TABLE IF EXISTS `alphabets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `alphabets` (
  `script` char(4) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `unicode_pdf` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `family` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `white_space` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `open_type_tag` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `complex_positioning` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `requires_font` tinyint(1) NOT NULL DEFAULT '0',
  `unicode` tinyint(1) NOT NULL DEFAULT '1',
  `diacritics` tinyint(1) DEFAULT NULL,
  `contextual_forms` tinyint(1) DEFAULT NULL,
  `reordering` tinyint(1) DEFAULT NULL,
  `case` tinyint(1) DEFAULT NULL,
  `split_graphs` tinyint(1) DEFAULT NULL,
  `status` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `baseline` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ligatures` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `direction` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `direction_notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `sample` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `sample_img` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`script`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `assets`
--

DROP TABLE IF EXISTS `assets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `assets` (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `organization_id` int unsigned NOT NULL,
  `asset_type` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `hidden` tinyint(1) NOT NULL DEFAULT '0',
  `base_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `buckets_id_unique` (`id`),
  KEY `buckets_organization_id_foreign` (`organization_id`),
  CONSTRAINT `FK_organizations_assets` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `audio_storage_view`
--

DROP TABLE IF EXISTS `audio_storage_view`;
/*!50001 DROP VIEW IF EXISTS `audio_storage_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `audio_storage_view` (
  `bucket` tinyint NOT NULL,
  `type` tinyint NOT NULL,
  `set_type_code` tinyint NOT NULL,
  `bible` tinyint NOT NULL,
  `fileset` tinyint NOT NULL,
  `file_name` tinyint NOT NULL,
  `path` tinyint NOT NULL,
  `size` tinyint NOT NULL,
  `last_modified_date` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `bible_books`
--

DROP TABLE IF EXISTS `bible_books`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_books` (
  `bible_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `book_id` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name_short` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `chapters` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `book_seq` char(4) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`bible_id`,`book_id`),
  KEY `bible_books_bible_id_foreign` (`bible_id`),
  KEY `bible_books_book_id_foreign` (`book_id`),
  CONSTRAINT `FK_bibles_bible_books` FOREIGN KEY (`bible_id`) REFERENCES `bibles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_books_bible_books` FOREIGN KEY (`book_id`) REFERENCES `books` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bible_concordance`
--

DROP TABLE IF EXISTS `bible_concordance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_concordance` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `key_word` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `bible_concordance_key_word_unique` (`key_word`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bible_equivalents`
--

DROP TABLE IF EXISTS `bible_equivalents`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_equivalents` (
  `bible_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '',
  `equivalent_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `organization_id` int unsigned NOT NULL,
  `type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `site` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `suffix` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `bible_equivalents_bible_id_foreign` (`bible_id`),
  KEY `bible_equivalents_organization_id_foreign` (`organization_id`),
  CONSTRAINT `FK_bibles_bible_equivalents` FOREIGN KEY (`bible_id`) REFERENCES `bibles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_organizations_bible_equivalents` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bible_file_stream_bandwidths`
--

DROP TABLE IF EXISTS `bible_file_stream_bandwidths`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_file_stream_bandwidths` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `bible_file_id` int unsigned NOT NULL,
  `file_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `bandwidth` int unsigned NOT NULL,
  `resolution_width` int unsigned DEFAULT NULL,
  `resolution_height` int unsigned DEFAULT NULL,
  `codec` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `stream` tinyint(1) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `bible_file_video_resolutions_bible_file_id_foreign` (`bible_file_id`),
  CONSTRAINT `FK_bible_files_bible_file_stream_bandwidths` FOREIGN KEY (`bible_file_id`) REFERENCES `bible_files` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=764523 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `bible_file_stream_bandwidths_view`
--

DROP TABLE IF EXISTS `bible_file_stream_bandwidths_view`;
/*!50001 DROP VIEW IF EXISTS `bible_file_stream_bandwidths_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `bible_file_stream_bandwidths_view` (
  `fileset_id` tinyint NOT NULL,
  `asset_id` tinyint NOT NULL,
  `set_type_code` tinyint NOT NULL,
  `hash_id` tinyint NOT NULL,
  `file_id` tinyint NOT NULL,
  `file_name` tinyint NOT NULL,
  `bandwidth` tinyint NOT NULL,
  `resolution_width` tinyint NOT NULL,
  `resolution_height` tinyint NOT NULL,
  `codec` tinyint NOT NULL,
  `stream` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `bible_file_stream_bytes`
--

DROP TABLE IF EXISTS `bible_file_stream_bytes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_file_stream_bytes` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `stream_bandwidth_id` int unsigned NOT NULL,
  `runtime` double(8,2) NOT NULL,
  `bytes` int NOT NULL,
  `offset` int NOT NULL,
  `timestamp_id` int unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `FK_bible_file_bandwidth_stream_bytes` (`stream_bandwidth_id`),
  KEY `FK_bible_file_timestamp_stream_bytes` (`timestamp_id`),
  CONSTRAINT `FK_bible_file_bandwidth_stream_bytes` FOREIGN KEY (`stream_bandwidth_id`) REFERENCES `bible_file_stream_bandwidths` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_bible_file_timestamp_stream_bytes` FOREIGN KEY (`timestamp_id`) REFERENCES `bible_file_timestamps` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=7109938 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bible_file_stream_ts`
--

DROP TABLE IF EXISTS `bible_file_stream_ts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_file_stream_ts` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `stream_bandwidth_id` int unsigned NOT NULL,
  `file_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `runtime` double(8,2) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `file_name` (`file_name`),
  KEY `bible_file_video_transport_stream_video_resolution_id_foreign` (`stream_bandwidth_id`),
  CONSTRAINT `FK_stream_bandwidths_stream_ts` FOREIGN KEY (`stream_bandwidth_id`) REFERENCES `bible_file_stream_bandwidths` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=42885891 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `bible_file_stream_ts_view`
--

DROP TABLE IF EXISTS `bible_file_stream_ts_view`;
/*!50001 DROP VIEW IF EXISTS `bible_file_stream_ts_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `bible_file_stream_ts_view` (
  `fileset_id` tinyint NOT NULL,
  `asset_id` tinyint NOT NULL,
  `set_type_code` tinyint NOT NULL,
  `hash_id` tinyint NOT NULL,
  `file_id` tinyint NOT NULL,
  `stream_bandwidth_id` tinyint NOT NULL,
  `file_name` tinyint NOT NULL,
  `runtime` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `bible_file_tags`
--

DROP TABLE IF EXISTS `bible_file_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_file_tags` (
  `file_id` int unsigned NOT NULL,
  `tag` varchar(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `value` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `admin_only` tinyint(1) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`file_id`,`tag`),
  UNIQUE KEY `unique_bible_file_tag` (`file_id`,`tag`,`value`),
  KEY `index_tag` (`tag`),
  CONSTRAINT `FK_bible_files_bible_file_tags` FOREIGN KEY (`file_id`) REFERENCES `bible_files` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `bible_file_tags_view`
--

DROP TABLE IF EXISTS `bible_file_tags_view`;
/*!50001 DROP VIEW IF EXISTS `bible_file_tags_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `bible_file_tags_view` (
  `fileset_id` tinyint NOT NULL,
  `asset_id` tinyint NOT NULL,
  `set_type_code` tinyint NOT NULL,
  `hash_id` tinyint NOT NULL,
  `file_id` tinyint NOT NULL,
  `tag` tinyint NOT NULL,
  `value` tinyint NOT NULL,
  `admin_only` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `bible_file_timestamps`
--

DROP TABLE IF EXISTS `bible_file_timestamps`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_file_timestamps` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `bible_file_id` int unsigned NOT NULL,
  `verse_start` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `verse_end` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `timestamp` float NOT NULL,
  `timestamp_end` float DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `verse_sequence` tinyint unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `bible_file_timestamps_file_id_foreign` (`bible_file_id`),
  CONSTRAINT `FK_bible_files_bible_file_timestamps` FOREIGN KEY (`bible_file_id`) REFERENCES `bible_files` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=4343139 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bible_file_titles`
--

DROP TABLE IF EXISTS `bible_file_titles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_file_titles` (
  `file_id` int unsigned NOT NULL,
  `iso` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `title` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `key_words` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `bible_file_titles_file_id_foreign` (`file_id`),
  KEY `bible_file_titles_iso_foreign` (`iso`),
  CONSTRAINT `FK_bible_files_bible_file_titles` FOREIGN KEY (`file_id`) REFERENCES `bible_files` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_languages_bible_file_titles` FOREIGN KEY (`iso`) REFERENCES `languages` (`iso`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `bible_file_video_resolutions`
--

DROP TABLE IF EXISTS `bible_file_video_resolutions`;
/*!50001 DROP VIEW IF EXISTS `bible_file_video_resolutions`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `bible_file_video_resolutions` (
  `id` tinyint NOT NULL,
  `bible_file_id` tinyint NOT NULL,
  `file_name` tinyint NOT NULL,
  `bandwidth` tinyint NOT NULL,
  `resolution_width` tinyint NOT NULL,
  `resolution_height` tinyint NOT NULL,
  `codec` tinyint NOT NULL,
  `stream` tinyint NOT NULL,
  `created_at` tinyint NOT NULL,
  `updated_at` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Temporary table structure for view `bible_file_video_transport_stream`
--

DROP TABLE IF EXISTS `bible_file_video_transport_stream`;
/*!50001 DROP VIEW IF EXISTS `bible_file_video_transport_stream`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `bible_file_video_transport_stream` (
  `id` tinyint NOT NULL,
  `video_resolution_id` tinyint NOT NULL,
  `file_name` tinyint NOT NULL,
  `runtime` tinyint NOT NULL,
  `created_at` tinyint NOT NULL,
  `updated_at` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `bible_files`
--

DROP TABLE IF EXISTS `bible_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_files` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `hash_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `book_id` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `chapter_start` tinyint unsigned DEFAULT NULL,
  `chapter_end` tinyint unsigned DEFAULT NULL,
  `verse_start` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `verse_end` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `file_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `file_size` int unsigned DEFAULT NULL,
  `duration` int unsigned DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `verse_sequence` tinyint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_bible_file_by_reference` (`hash_id`,`book_id`,`chapter_start`,`verse_start`),
  KEY `bible_files_book_id_foreign` (`book_id`),
  CONSTRAINT `FK_bible_filesets_bible_files` FOREIGN KEY (`hash_id`) REFERENCES `bible_filesets` (`hash_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_books_bible_files` FOREIGN KEY (`book_id`) REFERENCES `books` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=4477870 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bible_files_secondary`
--

DROP TABLE IF EXISTS `bible_files_secondary`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_files_secondary` (
  `hash_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `file_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `file_type` varchar(12) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`hash_id`,`file_name`),
  CONSTRAINT `FK_bible_filesets_bible_files_secondary` FOREIGN KEY (`hash_id`) REFERENCES `bible_filesets` (`hash_id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `bible_files_view`
--

DROP TABLE IF EXISTS `bible_files_view`;
/*!50001 DROP VIEW IF EXISTS `bible_files_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `bible_files_view` (
  `fileset_id` tinyint NOT NULL,
  `asset_id` tinyint NOT NULL,
  `set_type_code` tinyint NOT NULL,
  `hash_id` tinyint NOT NULL,
  `book_id` tinyint NOT NULL,
  `chapter_start` tinyint NOT NULL,
  `chapter_end` tinyint NOT NULL,
  `verse_start` tinyint NOT NULL,
  `verse_end` tinyint NOT NULL,
  `file_name` tinyint NOT NULL,
  `file_size` tinyint NOT NULL,
  `duration` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Temporary table structure for view `bible_files_view2`
--

DROP TABLE IF EXISTS `bible_files_view2`;
/*!50001 DROP VIEW IF EXISTS `bible_files_view2`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `bible_files_view2` (
  `fileset_id` tinyint NOT NULL,
  `asset_id` tinyint NOT NULL,
  `set_type_code` tinyint NOT NULL,
  `hash_id` tinyint NOT NULL,
  `book_id` tinyint NOT NULL,
  `chapter_start` tinyint NOT NULL,
  `chapter_end` tinyint NOT NULL,
  `verse_start` tinyint NOT NULL,
  `verse_end` tinyint NOT NULL,
  `file_name` tinyint NOT NULL,
  `file_size` tinyint NOT NULL,
  `duration` tinyint NOT NULL,
  `created_at` tinyint NOT NULL,
  `updated_at` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `bible_fileset_connections`
--

DROP TABLE IF EXISTS `bible_fileset_connections`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_fileset_connections` (
  `hash_id` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `bible_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`hash_id`,`bible_id`),
  KEY `bible_fileset_connections_hash_id_foreign` (`hash_id`),
  KEY `bible_fileset_connections_bible_id_index` (`bible_id`),
  KEY `index_hash_id` (`hash_id`),
  CONSTRAINT `FK_bible_filesets_bible_fileset_connections` FOREIGN KEY (`hash_id`) REFERENCES `bible_filesets` (`hash_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_bibles_bible_fileset_connections` FOREIGN KEY (`bible_id`) REFERENCES `bibles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `bible_fileset_connections_view`
--

DROP TABLE IF EXISTS `bible_fileset_connections_view`;
/*!50001 DROP VIEW IF EXISTS `bible_fileset_connections_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `bible_fileset_connections_view` (
  `bible_id` tinyint NOT NULL,
  `id` tinyint NOT NULL,
  `mode` tinyint NOT NULL,
  `type` tinyint NOT NULL,
  `size` tinyint NOT NULL,
  `stockno` tinyint NOT NULL,
  `hash_id` tinyint NOT NULL,
  `content_loaded` tinyint NOT NULL,
  `archived` tinyint NOT NULL,
  `updated_at` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `bible_fileset_copyright_organizations`
--

DROP TABLE IF EXISTS `bible_fileset_copyright_organizations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_fileset_copyright_organizations` (
  `hash_id` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `organization_id` int unsigned NOT NULL,
  `organization_role` int NOT NULL DEFAULT '2',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`hash_id`,`organization_role`,`organization_id`),
  KEY `FK_org_id` (`organization_id`),
  KEY `FK_org_role` (`organization_role`),
  CONSTRAINT `FK_bible_fileset_copyright_roles_bible_fileset_copyright_organiz` FOREIGN KEY (`organization_role`) REFERENCES `bible_fileset_copyright_roles` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `FK_bible_filesets_bible_fileset_copyright_organizations` FOREIGN KEY (`hash_id`) REFERENCES `bible_filesets` (`hash_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_organizations_bible_fileset_copyright_organizations` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bible_fileset_copyright_organizations_backup`
--

DROP TABLE IF EXISTS `bible_fileset_copyright_organizations_backup`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_fileset_copyright_organizations_backup` (
  `hash_id` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `organization_id` int unsigned NOT NULL,
  `organization_role` int NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `updated_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00'
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `bible_fileset_copyright_organizations_view`
--

DROP TABLE IF EXISTS `bible_fileset_copyright_organizations_view`;
/*!50001 DROP VIEW IF EXISTS `bible_fileset_copyright_organizations_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `bible_fileset_copyright_organizations_view` (
  `fileset_id` tinyint NOT NULL,
  `asset_id` tinyint NOT NULL,
  `set_type_code` tinyint NOT NULL,
  `hash_id` tinyint NOT NULL,
  `organization_id` tinyint NOT NULL,
  `organization_role` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `bible_fileset_copyright_roles`
--

DROP TABLE IF EXISTS `bible_fileset_copyright_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_fileset_copyright_roles` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bible_fileset_copyrights`
--

DROP TABLE IF EXISTS `bible_fileset_copyrights`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_fileset_copyrights` (
  `hash_id` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `copyright_date` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `copyright` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `copyright_description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `open_access` tinyint(1) NOT NULL DEFAULT '1',
  `is_combined` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`hash_id`),
  CONSTRAINT `FK_bible_filesets_bible_fileset_copyrights` FOREIGN KEY (`hash_id`) REFERENCES `bible_filesets` (`hash_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `bible_fileset_copyrights_view`
--

DROP TABLE IF EXISTS `bible_fileset_copyrights_view`;
/*!50001 DROP VIEW IF EXISTS `bible_fileset_copyrights_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `bible_fileset_copyrights_view` (
  `fileset_id` tinyint NOT NULL,
  `asset_id` tinyint NOT NULL,
  `set_type_code` tinyint NOT NULL,
  `hash_id` tinyint NOT NULL,
  `stockno` tinyint NOT NULL,
  `copyright_date` tinyint NOT NULL,
  `copyright` tinyint NOT NULL,
  `is_combined` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `bible_fileset_fonts`
--

DROP TABLE IF EXISTS `bible_fileset_fonts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_fileset_fonts` (
  `hash_id` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `font_id` int unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `FK_bible_filesets_bible_fileset_fonts` (`hash_id`),
  KEY `FK_fonts_bible_fileset_fonts` (`font_id`),
  CONSTRAINT `FK_bible_filesets_bible_fileset_fonts` FOREIGN KEY (`hash_id`) REFERENCES `bible_filesets` (`hash_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_fonts_bible_fileset_fonts` FOREIGN KEY (`font_id`) REFERENCES `fonts` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `bible_fileset_lookup`
--

DROP TABLE IF EXISTS `bible_fileset_lookup`;
/*!50001 DROP VIEW IF EXISTS `bible_fileset_lookup`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `bible_fileset_lookup` (
  `stocknumber` tinyint NOT NULL,
  `bibleid` tinyint NOT NULL,
  `filesetid` tinyint NOT NULL,
  `type` tinyint NOT NULL,
  `mode` tinyint NOT NULL,
  `size` tinyint NOT NULL,
  `hash_id` tinyint NOT NULL,
  `language` tinyint NOT NULL,
  `version` tinyint NOT NULL,
  `licensor` tinyint NOT NULL,
  `content_loaded` tinyint NOT NULL,
  `archived` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Temporary table structure for view `bible_fileset_lookup2`
--

DROP TABLE IF EXISTS `bible_fileset_lookup2`;
/*!50001 DROP VIEW IF EXISTS `bible_fileset_lookup2`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `bible_fileset_lookup2` (
  `stocknumber` tinyint NOT NULL,
  `bibleid` tinyint NOT NULL,
  `filesetid` tinyint NOT NULL,
  `type` tinyint NOT NULL,
  `hash_id` tinyint NOT NULL,
  `language` tinyint NOT NULL,
  `version` tinyint NOT NULL,
  `licensor` tinyint NOT NULL,
  `iso` tinyint NOT NULL,
  `rolv` tinyint NOT NULL,
  `languageId` tinyint NOT NULL,
  `countryId` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `bible_fileset_modes`
--

DROP TABLE IF EXISTS `bible_fileset_modes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_fileset_modes` (
  `id` tinyint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `bible_fileset_modes_name_unique` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bible_fileset_relations`
--

DROP TABLE IF EXISTS `bible_fileset_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_fileset_relations` (
  `id` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `parent_hash_id` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `child_hash_id` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `relationship` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `bible_fileset_relations_parent_hash_id_index` (`parent_hash_id`),
  KEY `bible_fileset_relations_child_hash_id_index` (`child_hash_id`),
  CONSTRAINT `FK_bible_filesets_bible_fileset_relations_child_hash_id` FOREIGN KEY (`child_hash_id`) REFERENCES `bible_filesets` (`hash_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_bible_filesets_bible_fileset_relations_parent_hash_id` FOREIGN KEY (`parent_hash_id`) REFERENCES `bible_filesets` (`hash_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bible_fileset_sizes`
--

DROP TABLE IF EXISTS `bible_fileset_sizes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_fileset_sizes` (
  `id` tinyint unsigned NOT NULL AUTO_INCREMENT,
  `set_size_code` char(9) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `bible_fileset_sizes_set_size_code_unique` (`set_size_code`),
  UNIQUE KEY `bible_fileset_sizes_name_unique` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bible_fileset_tags`
--

DROP TABLE IF EXISTS `bible_fileset_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_fileset_tags` (
  `hash_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `admin_only` tinyint(1) NOT NULL,
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `iso` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'eng',
  `language_id` int unsigned NOT NULL DEFAULT '6414',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'UTC',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'UTC',
  PRIMARY KEY (`hash_id`,`name`,`language_id`),
  KEY `bible_fileset_tags_hash_id_index` (`hash_id`),
  KEY `bible_fileset_tags_iso_index` (`iso`),
  KEY `language_id` (`language_id`),
  KEY `description` (`description`(4)),
  KEY `hashid_name_index` (`hash_id`,`name`),
  KEY `name_index` (`name`),
  CONSTRAINT `FK_bible_filesets_bible_fileset_tags` FOREIGN KEY (`hash_id`) REFERENCES `bible_filesets` (`hash_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `bible_fileset_tags_view`
--

DROP TABLE IF EXISTS `bible_fileset_tags_view`;
/*!50001 DROP VIEW IF EXISTS `bible_fileset_tags_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `bible_fileset_tags_view` (
  `fileset_id` tinyint NOT NULL,
  `asset_id` tinyint NOT NULL,
  `set_type_code` tinyint NOT NULL,
  `hash_id` tinyint NOT NULL,
  `name` tinyint NOT NULL,
  `description` tinyint NOT NULL,
  `admin_only` tinyint NOT NULL,
  `notes` tinyint NOT NULL,
  `iso` tinyint NOT NULL,
  `language_id` tinyint NOT NULL,
  `content_loaded` tinyint NOT NULL,
  `archived` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `bible_fileset_types`
--

DROP TABLE IF EXISTS `bible_fileset_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_fileset_types` (
  `id` tinyint unsigned NOT NULL AUTO_INCREMENT,
  `set_type_code` varchar(18) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `mode_id` tinyint unsigned NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `bible_fileset_types_set_type_code_unique` (`set_type_code`),
  UNIQUE KEY `bible_fileset_types_name_unique` (`name`),
  KEY `bible_fileset_types_mode_index` (`mode_id`),
  CONSTRAINT `FK_bible_fileset_modes_bible_fileset_types` FOREIGN KEY (`mode_id`) REFERENCES `bible_fileset_modes` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bible_filesets`
--

DROP TABLE IF EXISTS `bible_filesets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_filesets` (
  `id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `hash_id` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `asset_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `set_type_code` varchar(18) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `set_size_code` char(9) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `hidden` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `content_loaded` tinyint(1) NOT NULL DEFAULT '0',
  `archived` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`hash_id`),
  UNIQUE KEY `unique_prefix_for_s3` (`id`,`asset_id`,`set_type_code`),
  KEY `bible_filesets_bucket_id_foreign` (`asset_id`),
  KEY `bible_filesets_set_type_code_foreign` (`set_type_code`),
  KEY `bible_filesets_set_size_code_foreign` (`set_size_code`),
  KEY `bible_filesets_id_index` (`id`),
  KEY `bible_filesets_hash_id_index` (`hash_id`),
  CONSTRAINT `FK_assets_bible_filesets` FOREIGN KEY (`asset_id`) REFERENCES `assets` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_bible_fileset_sizes_bible_filesets` FOREIGN KEY (`set_size_code`) REFERENCES `bible_fileset_sizes` (`set_size_code`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_bible_fileset_types_bible_filesets` FOREIGN KEY (`set_type_code`) REFERENCES `bible_fileset_types` (`set_type_code`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `bible_filesets_view`
--

DROP TABLE IF EXISTS `bible_filesets_view`;
/*!50001 DROP VIEW IF EXISTS `bible_filesets_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `bible_filesets_view` (
  `fileset_id` tinyint NOT NULL,
  `hash_id` tinyint NOT NULL,
  `asset_id` tinyint NOT NULL,
  `set_type_code` tinyint NOT NULL,
  `set_size_code` tinyint NOT NULL,
  `hidden` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `bible_links`
--

DROP TABLE IF EXISTS `bible_links`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_links` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `bible_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `url` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `title` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `provider` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `visible` tinyint(1) NOT NULL DEFAULT '1',
  `organization_id` int unsigned DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `bible_links_bible_id_foreign` (`bible_id`),
  KEY `bible_links_organization_id_foreign` (`organization_id`),
  CONSTRAINT `FK_bibles_bible_links` FOREIGN KEY (`bible_id`) REFERENCES `bibles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=52468 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bible_organizations`
--

DROP TABLE IF EXISTS `bible_organizations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_organizations` (
  `bible_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `organization_id` int unsigned DEFAULT NULL,
  `relationship_type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `bible_organizations_bible_id_foreign` (`bible_id`),
  KEY `bible_organizations_organization_id_foreign` (`organization_id`),
  CONSTRAINT `FK_bibles_bible_organizations` FOREIGN KEY (`bible_id`) REFERENCES `bibles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bible_size_translations`
--

DROP TABLE IF EXISTS `bible_size_translations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_size_translations` (
  `set_size_code` char(9) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `iso` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`set_size_code`),
  KEY `bible_size_translations_iso_index` (`iso`),
  CONSTRAINT `FK_bible_fileset_sizes_bible_size_translations` FOREIGN KEY (`set_size_code`) REFERENCES `bible_fileset_sizes` (`set_size_code`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_languages_bible_size_translations` FOREIGN KEY (`iso`) REFERENCES `languages` (`iso`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bible_translations`
--

DROP TABLE IF EXISTS `bible_translations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_translations` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `language_id` int unsigned NOT NULL,
  `bible_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `vernacular` tinyint(1) NOT NULL DEFAULT '0',
  `vernacular_trade` tinyint(1) NOT NULL DEFAULT '0',
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `background` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique` (`language_id`,`bible_id`,`vernacular`),
  KEY `bible_translations_language_id_foreign` (`language_id`),
  KEY `bible_translations_bible_id_foreign` (`bible_id`),
  FULLTEXT KEY `ft_index_bible_translations_name` (`name`),
  CONSTRAINT `FK_bibles_bible_translations` FOREIGN KEY (`bible_id`) REFERENCES `bibles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_languages_bible_translations` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=16515 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bible_translator`
--

DROP TABLE IF EXISTS `bible_translator`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_translator` (
  `bible_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `translator_id` int unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `bible_translator_translator_id_foreign` (`translator_id`),
  KEY `bible_translator_bible_id_index` (`bible_id`),
  CONSTRAINT `FK_bibles_bible_translator` FOREIGN KEY (`bible_id`) REFERENCES `bibles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_translators_bible_translator` FOREIGN KEY (`translator_id`) REFERENCES `translators` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bible_verse_concordance`
--

DROP TABLE IF EXISTS `bible_verse_concordance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_verse_concordance` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `bible_verse_id` int unsigned NOT NULL,
  `bible_concordance` int unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `bible_verse_concordance_bible_verse_id_foreign` (`bible_verse_id`),
  KEY `bible_verse_concordance_bible_concordance_foreign` (`bible_concordance`),
  CONSTRAINT `FK_bible_concordance_bible_verse_concordance` FOREIGN KEY (`bible_concordance`) REFERENCES `bible_concordance` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_bible_verses_bible_verse_concordance` FOREIGN KEY (`bible_verse_id`) REFERENCES `bible_verses` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bible_verses`
--

DROP TABLE IF EXISTS `bible_verses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bible_verses` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `hash_id` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `book_id` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `chapter` tinyint unsigned NOT NULL,
  `verse_start` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `verse_end` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `verse_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `verse_sequence` tinyint unsigned NOT NULL,
  `created_at` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_text_reference` (`hash_id`,`book_id`,`chapter`,`verse_start`),
  KEY `bible_text_book_id_foreign` (`book_id`),
  KEY `bible_text_hash_id_index` (`hash_id`),
  KEY `index_hash_id` (`hash_id`),
  FULLTEXT KEY `verse_text` (`verse_text`),
  CONSTRAINT `bible_text_hash_id_foreign` FOREIGN KEY (`hash_id`) REFERENCES `bible_filesets` (`hash_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=43025837 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `biblebrain_database_export`
--

DROP TABLE IF EXISTS `biblebrain_database_export`;
/*!50001 DROP VIEW IF EXISTS `biblebrain_database_export`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `biblebrain_database_export` (
  `iso` tinyint NOT NULL,
  `lang_name` tinyint NOT NULL,
  `bible_name` tinyint NOT NULL,
  `bible_id` tinyint NOT NULL,
  `population` tinyint NOT NULL,
  `country` tinyint NOT NULL,
  `video_date` tinyint NOT NULL,
  `audio_date` tinyint NOT NULL,
  `text_date` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `bibles`
--

DROP TABLE IF EXISTS `bibles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bibles` (
  `id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `language_id` int unsigned NOT NULL,
  `versification` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'protestant',
  `numeral_system_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `date` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `scope` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `script` char(4) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `derived` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `copyright` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `priority` tinyint unsigned NOT NULL DEFAULT '0',
  `reviewed` tinyint(1) DEFAULT '0',
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `bibles_id_unique` (`id`),
  KEY `bibles_language_id_foreign` (`language_id`),
  KEY `bibles_numeral_system_id_foreign` (`numeral_system_id`),
  KEY `bibles_script_foreign` (`script`),
  KEY `priority` (`priority`),
  CONSTRAINT `FK_alphabets_bibles` FOREIGN KEY (`script`) REFERENCES `alphabets` (`script`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_languages_bibles` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_numeral_systems_bibles` FOREIGN KEY (`numeral_system_id`) REFERENCES `numeral_systems` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `bibles_defaults`
--

DROP TABLE IF EXISTS `bibles_defaults`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bibles_defaults` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `language_code` char(6) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `bible_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_bible_default` (`language_code`,`bible_id`,`type`),
  KEY `FK_bibles_defaults` (`bible_id`),
  CONSTRAINT `FK_bibles_defaults` FOREIGN KEY (`bible_id`) REFERENCES `bibles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=87 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `bibles_view`
--

DROP TABLE IF EXISTS `bibles_view`;
/*!50001 DROP VIEW IF EXISTS `bibles_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `bibles_view` (
  `fileset_id` tinyint NOT NULL,
  `bible_id` tinyint NOT NULL,
  `language_id` tinyint NOT NULL,
  `versification` tinyint NOT NULL,
  `numeral_system_id` tinyint NOT NULL,
  `date` tinyint NOT NULL,
  `scope` tinyint NOT NULL,
  `script` tinyint NOT NULL,
  `derived` tinyint NOT NULL,
  `copyright` tinyint NOT NULL,
  `priority` tinyint NOT NULL,
  `reviewed` tinyint NOT NULL,
  `notes` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `blm`
--

DROP TABLE IF EXISTS `blm`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `blm` (
  `ISO` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Country` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `stocknos` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `blm2`
--

DROP TABLE IF EXISTS `blm2`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `blm2` (
  `ISO` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Country` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `stocknos` text CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `book_translations`
--

DROP TABLE IF EXISTS `book_translations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `book_translations` (
  `language_id` int unsigned NOT NULL,
  `book_id` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name_long` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name_short` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name_abbreviation` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`language_id`,`book_id`),
  KEY `book_translations_language_id_foreign` (`language_id`),
  KEY `book_translations_book_id_foreign` (`book_id`),
  CONSTRAINT `FK_books_book_translations` FOREIGN KEY (`book_id`) REFERENCES `books` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `FK_languages_book_translations` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `books`
--

DROP TABLE IF EXISTS `books`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `books` (
  `id` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `id_usfx` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `id_osis` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `book_testament` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `book_group` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `chapters` int unsigned DEFAULT NULL,
  `verses` int unsigned DEFAULT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `testament_order` tinyint unsigned DEFAULT NULL,
  `protestant_order` tinyint unsigned DEFAULT NULL,
  `luther_order` tinyint unsigned DEFAULT NULL,
  `synodal_order` tinyint unsigned DEFAULT NULL,
  `german_order` tinyint unsigned DEFAULT NULL,
  `kjva_order` tinyint unsigned DEFAULT NULL,
  `vulgate_order` tinyint unsigned DEFAULT NULL,
  `lxx_order` tinyint unsigned DEFAULT NULL,
  `orthodox_order` tinyint unsigned DEFAULT NULL,
  `nrsva_order` tinyint unsigned DEFAULT NULL,
  `catholic_order` tinyint unsigned DEFAULT NULL,
  `finnish_order` tinyint unsigned DEFAULT NULL,
  `messianic_order` tinyint DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `commentaries`
--

DROP TABLE IF EXISTS `commentaries`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `commentaries` (
  `id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `author` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `date` int unsigned NOT NULL,
  `features` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `publisher` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `commentary_sections`
--

DROP TABLE IF EXISTS `commentary_sections`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `commentary_sections` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `commentary_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `title` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `book_id` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `chapter_start` tinyint unsigned DEFAULT NULL,
  `chapter_end` tinyint unsigned DEFAULT NULL,
  `verse_start` tinyint unsigned DEFAULT NULL,
  `verse_end` tinyint unsigned DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `FK_commentaries_commentary_sections` (`commentary_id`),
  KEY `FK_books_commentary_sections` (`book_id`),
  CONSTRAINT `FK_books_commentary_sections` FOREIGN KEY (`book_id`) REFERENCES `books` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `FK_commentaries_commentary_sections` FOREIGN KEY (`commentary_id`) REFERENCES `commentaries` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=16108 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `commentary_translations`
--

DROP TABLE IF EXISTS `commentary_translations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `commentary_translations` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `language_id` int unsigned NOT NULL,
  `commentary_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `vernacular` tinyint(1) NOT NULL DEFAULT '0',
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `FK_languages_commentary_translations` (`language_id`),
  KEY `FK_commentaries_commentary_translations` (`commentary_id`),
  CONSTRAINT `FK_commentaries_commentary_translations` FOREIGN KEY (`commentary_id`) REFERENCES `commentaries` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_languages_commentary_translations` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `connection_translations`
--

DROP TABLE IF EXISTS `connection_translations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `connection_translations` (
  `language_id` int unsigned NOT NULL,
  `resource_id` int unsigned NOT NULL,
  `vernacular` tinyint(1) NOT NULL,
  `tag` tinyint(1) NOT NULL,
  `title` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `connection_translations_language_id_foreign` (`language_id`),
  KEY `connection_translations_resource_id_foreign` (`resource_id`),
  CONSTRAINT `FK_languages_connection_translations` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_resources_connection_translations` FOREIGN KEY (`resource_id`) REFERENCES `resources` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `connections`
--

DROP TABLE IF EXISTS `connections`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `connections` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `organization_id` int unsigned NOT NULL,
  `site_url` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `title` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `cover_thumbnail` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `date` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `connections_organization_id_foreign` (`organization_id`),
  CONSTRAINT `FK_organizations_connections` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `countries`
--

DROP TABLE IF EXISTS `countries`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `countries` (
  `id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `iso_a3` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `fips` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `wfb` tinyint(1) NOT NULL DEFAULT '0',
  `ethnologue` tinyint(1) NOT NULL DEFAULT '0',
  `continent` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `introduction` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `overview` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `countries_iso_a3_unique` (`iso_a3`),
  FULLTEXT KEY `ft_index_countries_name_iso_a3` (`name`,`iso_a3`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `country_communications`
--

DROP TABLE IF EXISTS `country_communications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `country_communications` (
  `country_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `fixed_phones_total` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `fixed_phones_subs_per_100` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `mobile_phones_total` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `mobile_phones_subs_per_100` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `phone_system_general_assessment` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `phone_system_international` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `phone_system_domestic` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `broadcast_media` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `internet_country_code` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `internet_total_users` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `internet_population_percent` decimal(4,1) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`country_id`),
  KEY `country_communications_country_id_foreign` (`country_id`),
  CONSTRAINT `FK_countries_country_communications` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `country_economy`
--

DROP TABLE IF EXISTS `country_economy`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `country_economy` (
  `country_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `overview` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `gdp_power_parity` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `gdp_real_growth` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `gdp_per_capita` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `gdp_household_consumption` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `gdp_consumption` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `gdp_investment_in_fixed_capital` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `gdp_investment_in_inventories` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `gdp_exports` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `gdp_imports` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `gdp_sector_agriculture` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `gdp_sector_industry` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `gdp_sector_services` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `agriculture_products` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `industries` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `industrial_growth_rate` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `labor_force` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `labor_force_notes` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `labor_force_services` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `labor_force_industry` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `labor_force_agriculture` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `labor_force_occupation_notes` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `unemployment_rate` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `population_below_poverty` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `household_income_lowest_10` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `household_income_highest_10` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `budget_revenues` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `taxes_revenues` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `budget_net` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `public_debt` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `external_debt` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `fiscal_year` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `inflation_rate` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `central_bank_discount_rate` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `commercial_bank_prime_lending_rate` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `stock_money_narrow` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `stock_money_broad` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `stock_domestic_credit` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `exports` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `exports_commodities` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `exports_partners` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `imports` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `imports_commodities` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `imports_partners` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `exchange_rates` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`country_id`),
  KEY `country_economy_country_id_foreign` (`country_id`),
  CONSTRAINT `FK_countries_country_economy` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `country_energy`
--

DROP TABLE IF EXISTS `country_energy`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `country_energy` (
  `country_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `electricity_production` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `electricity_consumption` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `electricity_exports` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `electricity_imports` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `electricity_generating_capacity` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `electricity_fossil_fuels` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `electricity_nuclear` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `electricity_hydroelectric` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `electricity_renewable` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `crude_oil_production` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `crude_oil_exports` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `crude_oil_imports` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `crude_oil_reserves` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `petrol_production` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `petrol_consumption` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `petrol_exports` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `petrol_imports` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `natural_gas_production` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `natural_gas_consumption` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `natural_gas_exports` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `natural_gas_imports` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `natural_gas_reserves` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `co2_output` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`country_id`),
  KEY `country_energy_country_id_foreign` (`country_id`),
  CONSTRAINT `FK_countries_country_energy` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `country_geography`
--

DROP TABLE IF EXISTS `country_geography`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `country_geography` (
  `country_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `location_description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `latitude` decimal(10,7) DEFAULT NULL,
  `longitude` decimal(10,7) DEFAULT NULL,
  `mapReferences` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `area_sqkm_total` int unsigned DEFAULT NULL,
  `area_sqkm_land` int unsigned DEFAULT NULL,
  `area_sqkm_water` int unsigned DEFAULT NULL,
  `area_km_coastline` int unsigned DEFAULT NULL,
  `area_note` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `climate` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `terrain` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `hazards` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`country_id`),
  KEY `country_geography_country_id_foreign` (`country_id`),
  CONSTRAINT `FK_countries_country_geography` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `country_government`
--

DROP TABLE IF EXISTS `country_government`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `country_government` (
  `country_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name_etymology` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `conventional_long_form` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `conventional_short_form` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `dependency_status` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `government_type` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `capital` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `capital_coordinates` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `capital_time_zone` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `administrative_divisions` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `administrative_divisions_note` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `independence` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `national_holiday` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `constitution` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `legal_system` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `citizenship` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `suffrage` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `executive_chief_of_state` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `executive_head_of_government` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `executive_cabinet` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `executive_elections` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `executive_election_results` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `legislative_description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `legislative_elections` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `legislative_election_results` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `legislative_highest_courts` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `legislative_judge_selection` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `legislative_subordinate_courts` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `political_parties` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `political_pressure` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `international_organization_participation` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `diplomatic_representation_in_usa` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `diplomatic_representation_from_usa` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `flag_description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `national_symbols` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `national_anthem` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`country_id`),
  KEY `country_government_country_id_foreign` (`country_id`),
  CONSTRAINT `FK_countries_country_government` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `country_issues`
--

DROP TABLE IF EXISTS `country_issues`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `country_issues` (
  `country_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `international_disputes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `illicit_drugs` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `refugees` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`country_id`),
  KEY `country_issues_country_id_foreign` (`country_id`),
  CONSTRAINT `FK_countries_country_issues` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `country_joshua_project`
--

DROP TABLE IF EXISTS `country_joshua_project`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `country_joshua_project` (
  `country_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `language_official_iso` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `language_official_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `population` bigint unsigned NOT NULL DEFAULT '0',
  `population_unreached` bigint unsigned NOT NULL DEFAULT '0',
  `people_groups` int unsigned NOT NULL DEFAULT '0',
  `people_groups_unreached` int unsigned NOT NULL DEFAULT '0',
  `joshua_project_scale` tinyint unsigned NOT NULL DEFAULT '0',
  `primary_religion` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `percent_christian` double(8,2) DEFAULT NULL,
  `resistant_belt` tinyint(1) NOT NULL DEFAULT '0',
  `percent_literate` double(8,2) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`country_id`),
  KEY `country_joshua_project_country_id_foreign` (`country_id`),
  KEY `country_joshua_project_language_official_iso_foreign` (`language_official_iso`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `country_language`
--

DROP TABLE IF EXISTS `country_language`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `country_language` (
  `country_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `language_id` int unsigned NOT NULL,
  `population` int NOT NULL DEFAULT '0',
  `display_order` int DEFAULT NULL,
  UNIQUE KEY `uq_country_language` (`country_id`,`language_id`),
  KEY `country_language_language_id_foreign` (`language_id`),
  CONSTRAINT `FK_countries_country_language` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_languages_country_language` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `country_maps`
--

DROP TABLE IF EXISTS `country_maps`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `country_maps` (
  `country_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `thumbnail_url` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `map_url` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `country_maps_country_id_foreign` (`country_id`),
  CONSTRAINT `FK_countries_country_maps` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `country_people`
--

DROP TABLE IF EXISTS `country_people`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `country_people` (
  `country_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `languages` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `religions` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `population` int unsigned DEFAULT NULL,
  `population_date` int unsigned DEFAULT NULL,
  `nationality_noun` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `nationality_adjective` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `age_structure_14` decimal(4,2) unsigned DEFAULT NULL,
  `age_structure_24` decimal(4,2) unsigned DEFAULT NULL,
  `age_structure_54` decimal(4,2) unsigned DEFAULT NULL,
  `age_structure_64` decimal(4,2) unsigned DEFAULT NULL,
  `age_structure_65` decimal(4,2) unsigned DEFAULT NULL,
  `dependency_total` decimal(4,2) unsigned DEFAULT NULL,
  `dependency_youth` decimal(4,2) unsigned DEFAULT NULL,
  `dependency_elder` decimal(4,2) unsigned DEFAULT NULL,
  `dependency_potential` decimal(4,2) unsigned DEFAULT NULL,
  `median_age_total` decimal(3,2) unsigned DEFAULT NULL,
  `median_age_male` decimal(3,2) unsigned DEFAULT NULL,
  `median_age_female` decimal(3,2) unsigned DEFAULT NULL,
  `population_growth_rate_percentage` decimal(3,2) DEFAULT NULL,
  `birth_rate_per_1k` decimal(8,2) unsigned DEFAULT NULL,
  `death_rate_per_1k` decimal(8,2) unsigned DEFAULT NULL,
  `net_migration_per_1k` decimal(6,2) DEFAULT NULL,
  `population_distribution` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `urban_population_percentage` decimal(4,2) unsigned DEFAULT NULL,
  `urbanization_rate` decimal(4,2) unsigned DEFAULT NULL,
  `major_urban_areas_population` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sex_ratio_birth` decimal(3,1) unsigned DEFAULT NULL,
  `sex_ratio_14` decimal(3,1) unsigned DEFAULT NULL,
  `sex_ratio_24` decimal(3,1) unsigned DEFAULT NULL,
  `sex_ratio_54` decimal(3,1) unsigned DEFAULT NULL,
  `sex_ratio_64` decimal(3,1) unsigned DEFAULT NULL,
  `sex_ratio_65` decimal(3,1) unsigned DEFAULT NULL,
  `sex_ratio_total` decimal(3,1) unsigned DEFAULT NULL,
  `mother_age_first_birth` tinyint unsigned DEFAULT NULL,
  `maternal_mortality_rate` decimal(3,1) unsigned DEFAULT NULL,
  `infant_mortality_per_1k_total` decimal(3,2) unsigned DEFAULT NULL,
  `infant_mortality_per_1k_male` decimal(3,2) unsigned DEFAULT NULL,
  `infant_mortality_per_1k_female` decimal(3,2) unsigned DEFAULT NULL,
  `life_expectancy_at_birth_total` decimal(3,1) unsigned DEFAULT NULL,
  `life_expectancy_at_birth_male` decimal(3,1) unsigned DEFAULT NULL,
  `life_expectancy_at_birth_female` decimal(3,1) unsigned DEFAULT NULL,
  `total_fertility_rate` decimal(4,2) DEFAULT NULL,
  `contraceptive_prevalence` decimal(4,2) DEFAULT NULL,
  `health_expenditures` decimal(4,2) DEFAULT NULL,
  `physicians` decimal(4,2) DEFAULT NULL,
  `hospital_beds` decimal(4,2) DEFAULT NULL,
  `drinking_water_source_urban_improved` decimal(5,2) DEFAULT NULL,
  `drinking_water_source_rural_improved` decimal(5,2) DEFAULT NULL,
  `sanitation_facility_access_urban_improved` decimal(5,2) DEFAULT NULL,
  `sanitation_facility_access_rural_improved` decimal(5,2) DEFAULT NULL,
  `hiv_infection_rate` decimal(4,2) DEFAULT NULL,
  `hiv_infected` decimal(4,2) DEFAULT NULL,
  `hiv_deaths` decimal(4,2) DEFAULT NULL,
  `obesity_rate` decimal(4,2) DEFAULT NULL,
  `underweight_children` decimal(4,2) DEFAULT NULL,
  `education_expenditures` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `literacy_definition` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `literacy_total` decimal(5,2) unsigned DEFAULT NULL,
  `literacy_male` decimal(5,2) unsigned DEFAULT NULL,
  `literacy_female` decimal(5,2) unsigned DEFAULT NULL,
  `school_years_total` tinyint unsigned DEFAULT NULL,
  `school_years_male` tinyint unsigned DEFAULT NULL,
  `school_years_female` tinyint unsigned DEFAULT NULL,
  `child_labor` int unsigned DEFAULT NULL,
  `child_labor_percentage` decimal(4,2) unsigned DEFAULT NULL,
  `unemployment_youth_total` decimal(4,2) unsigned DEFAULT NULL,
  `unemployment_youth_male` decimal(4,2) unsigned DEFAULT NULL,
  `unemployment_youth_female` decimal(4,2) unsigned DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`country_id`),
  KEY `country_people_country_id_foreign` (`country_id`),
  CONSTRAINT `FK_countries_country_people` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `country_people_ethnicities`
--

DROP TABLE IF EXISTS `country_people_ethnicities`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `country_people_ethnicities` (
  `country_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `population_percentage` decimal(5,2) unsigned NOT NULL,
  `date` tinyint unsigned DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `country_people_ethnicities_country_id_foreign` (`country_id`),
  CONSTRAINT `FK_countries_country_people_ethnicities` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `country_regions`
--

DROP TABLE IF EXISTS `country_regions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `country_regions` (
  `country_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `language_id` int unsigned NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `country_regions_country_id_foreign` (`country_id`),
  KEY `country_regions_language_id_foreign` (`language_id`),
  CONSTRAINT `FK_countries_country_regions` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE,
  CONSTRAINT `FK_languages_country_regions` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `country_religions`
--

DROP TABLE IF EXISTS `country_religions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `country_religions` (
  `country_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `population_percentage` decimal(5,2) unsigned DEFAULT NULL,
  `date` tinyint unsigned DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `country_religions_country_id_foreign` (`country_id`),
  CONSTRAINT `FK_countries_country_religions` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `country_translations`
--

DROP TABLE IF EXISTS `country_translations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `country_translations` (
  `country_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `language_id` int unsigned NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `country_translations_country_id_foreign` (`country_id`),
  KEY `country_translations_language_id_foreign` (`language_id`),
  CONSTRAINT `FK_countries_country_translations` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_languages_country_translations` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `country_transportation`
--

DROP TABLE IF EXISTS `country_transportation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `country_transportation` (
  `country_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `air_carriers` int unsigned DEFAULT NULL,
  `aircraft` int unsigned DEFAULT NULL,
  `aircraft_passengers` int unsigned DEFAULT NULL,
  `aircraft_freight` int unsigned DEFAULT NULL,
  `aircraft_code_prefix` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `airports` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `airports_paved` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `airports_info_date` tinyint unsigned DEFAULT NULL,
  `major_seaports` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `oil_terminals` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cruise_ports` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `country_transportation_country_id_foreign` (`country_id`),
  CONSTRAINT `FK_countries_country_transportation` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `fileset_permission_view`
--

DROP TABLE IF EXISTS `fileset_permission_view`;
/*!50001 DROP VIEW IF EXISTS `fileset_permission_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `fileset_permission_view` (
  `fileset_id` tinyint NOT NULL,
  `hash_id` tinyint NOT NULL,
  `license_group_id` tinyint NOT NULL,
  `license_group` tinyint NOT NULL,
  `bible_id` tinyint NOT NULL,
  `mode` tinyint NOT NULL,
  `permission_pattern_id` tinyint NOT NULL,
  `permission_pattern` tinyint NOT NULL,
  `group_concat(access_group_id)` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `fonts`
--

DROP TABLE IF EXISTS `fonts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `fonts` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `data` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=81 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `language_bibleInfo`
--

DROP TABLE IF EXISTS `language_bibleInfo`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `language_bibleInfo` (
  `language_id` int unsigned NOT NULL,
  `bible_status` tinyint DEFAULT NULL,
  `bible_translation_need` tinyint(1) DEFAULT NULL,
  `bible_year` int DEFAULT NULL,
  `bible_year_newTestament` int DEFAULT NULL,
  `bible_year_portions` int DEFAULT NULL,
  `bible_sample_text` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `bible_sample_img` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `language_bibleinfo_language_id_foreign` (`language_id`),
  CONSTRAINT `FK_languages_language_bibleInfo` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `language_classifications`
--

DROP TABLE IF EXISTS `language_classifications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `language_classifications` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `language_id` int unsigned NOT NULL,
  `classification_id` char(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `order` tinyint unsigned NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `language_classifications_language_id_foreign` (`language_id`),
  CONSTRAINT `FK_languages_language_classifications` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=45209 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `language_codes`
--

DROP TABLE IF EXISTS `language_codes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `language_codes` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `language_id` int unsigned NOT NULL,
  `source` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `code` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `language_codes_language_id_foreign` (`language_id`),
  CONSTRAINT `FK_languages_language_codes` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1752 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `language_codes_v2`
--

DROP TABLE IF EXISTS `language_codes_v2`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `language_codes_v2` (
  `id` varchar(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `language_ISO_639_3_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `family_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `english_name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `language_dialects`
--

DROP TABLE IF EXISTS `language_dialects`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `language_dialects` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `language_id` int unsigned NOT NULL,
  `dialect_id` char(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `name` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `language_dialects_language_id_foreign` (`language_id`),
  KEY `language_dialects_dialect_id_index` (`dialect_id`)
) ENGINE=InnoDB AUTO_INCREMENT=445 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `language_status`
--

DROP TABLE IF EXISTS `language_status`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `language_status` (
  `id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `title` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `language_translations`
--

DROP TABLE IF EXISTS `language_translations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `language_translations` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `language_source_id` int unsigned NOT NULL,
  `language_translation_id` int unsigned NOT NULL,
  `name` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_as_ci DEFAULT NULL,
  `priority` tinyint NOT NULL DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_language_translations` (`language_source_id`,`language_translation_id`,`name`),
  KEY `language_translations_language_source_id_foreign` (`language_source_id`),
  KEY `language_translations_language_translation_id_foreign` (`language_translation_id`),
  FULLTEXT KEY `ft_index_language_translations_name` (`name`),
  CONSTRAINT `FK_languages_language_translations_language_source_id` FOREIGN KEY (`language_source_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_languages_language_translations_language_tranlation_id` FOREIGN KEY (`language_translation_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=728938 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `language_view`
--

DROP TABLE IF EXISTS `language_view`;
/*!50001 DROP VIEW IF EXISTS `language_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `language_view` (
  `id` tinyint NOT NULL,
  `iso` tinyint NOT NULL,
  `country_id` tinyint NOT NULL,
  `country_name` tinyint NOT NULL,
  `fcbh_name` tinyint NOT NULL,
  `heart_name` tinyint NOT NULL,
  `ethnologue_name` tinyint NOT NULL,
  `alt_name` tinyint NOT NULL,
  `official_name` tinyint NOT NULL,
  `rolv_code` tinyint NOT NULL,
  `pseudonym` tinyint NOT NULL,
  `sensitivity` tinyint NOT NULL,
  `deleted_at` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `languages`
--

DROP TABLE IF EXISTS `languages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `languages` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `glotto_id` char(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `iso` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `iso2B` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `iso2T` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `iso1` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `maps` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `development` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `use` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `location` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `area` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `population` int DEFAULT NULL,
  `population_notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `typology` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `writing` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `latitude` double(11,7) DEFAULT NULL,
  `longitude` double(11,7) DEFAULT NULL,
  `status_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `country_id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `rolv_code` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sensitivity` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'Low',
  `pseudonym` tinyint unsigned NOT NULL DEFAULT '0',
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `languages_glotto_id_unique` (`glotto_id`),
  UNIQUE KEY `languages_iso2b_unique` (`iso2B`),
  UNIQUE KEY `languages_iso2t_unique` (`iso2T`),
  UNIQUE KEY `languages_iso1_unique` (`iso1`),
  KEY `languages_iso_index` (`iso`),
  KEY `language_status_foreign_key` (`status_id`),
  KEY `country_id_foreign_key` (`country_id`),
  FULLTEXT KEY `ft_index_languages_name` (`name`),
  CONSTRAINT `FK_countries_languages` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON DELETE SET NULL ON UPDATE CASCADE,
  CONSTRAINT `FK_language_status_languages` FOREIGN KEY (`status_id`) REFERENCES `language_status` (`id`) ON DELETE SET NULL ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=34311 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `lg_simple`
--

DROP TABLE IF EXISTS `lg_simple`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `lg_simple` (
  `bibleid` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `mode` varchar(18) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `testament` varchar(2) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `license_group`
--

DROP TABLE IF EXISTS `license_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `license_group` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `bible_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `mode` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `size` char(9) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mode_type` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `unique_name` (`name`),
  KEY `unique_attributes` (`bible_id`,`mode`,`size`,`mode_type`)
) ENGINE=InnoDB AUTO_INCREMENT=17871 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `license_group_fileset_view`
--

DROP TABLE IF EXISTS `license_group_fileset_view`;
/*!50001 DROP VIEW IF EXISTS `license_group_fileset_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `license_group_fileset_view` (
  `bible_id` tinyint NOT NULL,
  `mode` tinyint NOT NULL,
  `size` tinyint NOT NULL,
  `mode_type` tinyint NOT NULL,
  `lg_name` tinyint NOT NULL,
  `fileset_id` tinyint NOT NULL,
  `hash_id` tinyint NOT NULL,
  `lg_id` tinyint NOT NULL,
  `permission_pattern_id` tinyint NOT NULL,
  `updated_at` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `license_group_filesets`
--

DROP TABLE IF EXISTS `license_group_filesets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `license_group_filesets` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `hash_id` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `license_group_id` int unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_hash_license_group` (`hash_id`,`license_group_id`),
  KEY `FK_license_group_filesets_bible_filesets__hash_id` (`hash_id`),
  KEY `FK_license_group_license_group_filesets` (`license_group_id`),
  CONSTRAINT `FK_bible_filesets_license_group_filesets` FOREIGN KEY (`hash_id`) REFERENCES `bible_filesets` (`hash_id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_license_group_license_group_filesets` FOREIGN KEY (`license_group_id`) REFERENCES `license_group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=34351 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `license_group_permission_pattern`
--

DROP TABLE IF EXISTS `license_group_permission_pattern`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `license_group_permission_pattern` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `license_group_id` int unsigned NOT NULL,
  `permission_pattern_id` int unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_license_group` (`license_group_id`),
  KEY `FK_license_group_permission_pattern_permission_pattern` (`permission_pattern_id`),
  CONSTRAINT `FK_license_group_permission_pattern_license_group` FOREIGN KEY (`license_group_id`) REFERENCES `license_group` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_license_group_permission_pattern_permission_pattern` FOREIGN KEY (`permission_pattern_id`) REFERENCES `permission_pattern` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=34394 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `license_group_permission_pattern_view`
--

DROP TABLE IF EXISTS `license_group_permission_pattern_view`;
/*!50001 DROP VIEW IF EXISTS `license_group_permission_pattern_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `license_group_permission_pattern_view` (
  `license_group` tinyint NOT NULL,
  `permission_pattern` tinyint NOT NULL,
  `license_group_id` tinyint NOT NULL,
  `permission_pattern_id` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Temporary table structure for view `licensor_lookup`
--

DROP TABLE IF EXISTS `licensor_lookup`;
/*!50001 DROP VIEW IF EXISTS `licensor_lookup`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `licensor_lookup` (
  `licensor` tinyint NOT NULL,
  `bible_id` tinyint NOT NULL,
  `mode` tinyint NOT NULL,
  `stock_no` tinyint NOT NULL,
  `set_size_code` tinyint NOT NULL,
  `testament` tinyint NOT NULL,
  `copyright_date` tinyint NOT NULL,
  `copyright` tinyint NOT NULL,
  `organization_id` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `lpts_organizations`
--

DROP TABLE IF EXISTS `lpts_organizations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `lpts_organizations` (
  `lpts_organization` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `organization_role` int unsigned NOT NULL,
  `organization_id` int unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`),
  KEY `organization_id` (`organization_id`),
  CONSTRAINT `lpts_organizations_ibfk_1` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2781 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `lpts_script_codes`
--

DROP TABLE IF EXISTS `lpts_script_codes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `lpts_script_codes` (
  `lpts_name` varchar(256) NOT NULL,
  `script_id` char(4) NOT NULL,
  PRIMARY KEY (`lpts_name`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `lptslangs`
--

DROP TABLE IF EXISTS `lptslangs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `lptslangs` (
  `ISO` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `LangName` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `id` char(2) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `migrations`
--

DROP TABLE IF EXISTS `migrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `migrations` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `migration` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `batch` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `mondayexport`
--

DROP TABLE IF EXISTS `mondayexport`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `mondayexport` (
  `stocknumber` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `permission` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `audio` varchar(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `text` varchar(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `video` varchar(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `id` int NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=510814 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `numeral_system_glyphs`
--

DROP TABLE IF EXISTS `numeral_system_glyphs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `numeral_system_glyphs` (
  `numeral_system_id` char(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `value` tinyint unsigned NOT NULL,
  `glyph` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `numeral_written` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY `uq_numeral_system_glyph` (`numeral_system_id`,`value`,`glyph`),
  KEY `numeral_system_glyphs_numeral_system_id_index` (`numeral_system_id`),
  CONSTRAINT `FK_numeral_systems_numeral_system_glyphs` FOREIGN KEY (`numeral_system_id`) REFERENCES `numeral_systems` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `numeral_systems`
--

DROP TABLE IF EXISTS `numeral_systems`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `numeral_systems` (
  `id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `organization_logos`
--

DROP TABLE IF EXISTS `organization_logos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `organization_logos` (
  `organization_id` int unsigned NOT NULL,
  `language_id` int unsigned NOT NULL,
  `language_iso` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `url` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `icon` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`organization_id`,`language_id`,`icon`),
  KEY `organization_logos_organization_id_foreign` (`organization_id`),
  KEY `organization_logos_language_id_foreign` (`language_id`),
  KEY `organization_logos_language_iso_foreign` (`language_iso`),
  CONSTRAINT `FK_organizations_organization_logos` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `organization_relationships`
--

DROP TABLE IF EXISTS `organization_relationships`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `organization_relationships` (
  `organization_parent_id` int unsigned NOT NULL,
  `organization_child_id` int unsigned NOT NULL,
  `type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`organization_parent_id`,`organization_child_id`),
  KEY `organization_relationships_organization_parent_id_foreign` (`organization_parent_id`),
  KEY `organization_relationships_organization_child_id_foreign` (`organization_child_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `organization_translations`
--

DROP TABLE IF EXISTS `organization_translations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `organization_translations` (
  `language_id` int unsigned NOT NULL,
  `organization_id` int unsigned NOT NULL,
  `vernacular` tinyint(1) NOT NULL DEFAULT '0',
  `alt` tinyint(1) NOT NULL DEFAULT '0',
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `description_short` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `organization_translations_language_id_foreign` (`language_id`),
  KEY `organization_translations_organization_id_foreign` (`organization_id`),
  CONSTRAINT `FK_languages_organization_translations` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_organizations_organization_translations` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `organizations`
--

DROP TABLE IF EXISTS `organizations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `organizations` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `slug` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `abbreviation` char(6) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `primaryColor` varchar(7) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `secondaryColor` varchar(7) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `inactive` tinyint(1) DEFAULT '0',
  `url_facebook` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `url_website` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `url_donate` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `url_twitter` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `address` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `address2` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `city` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `state` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `country` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `zip` int unsigned DEFAULT NULL,
  `phone` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `email` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `email_director` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `latitude` double(11,7) DEFAULT NULL,
  `longitude` double(11,7) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `organizations_slug_unique` (`slug`),
  UNIQUE KEY `organizations_abbreviation_unique` (`abbreviation`),
  KEY `organizations_country_foreign` (`country`),
  CONSTRAINT `FK_countries_organizations` FOREIGN KEY (`country`) REFERENCES `countries` (`id`) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=4201 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `permission_pattern`
--

DROP TABLE IF EXISTS `permission_pattern`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `permission_pattern` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `mode_audio` tinyint(1) NOT NULL DEFAULT '0',
  `mode_video` tinyint(1) NOT NULL DEFAULT '0',
  `mode_text` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=138 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `permission_pattern_access_group`
--

DROP TABLE IF EXISTS `permission_pattern_access_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `permission_pattern_access_group` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `permission_pattern_id` int unsigned NOT NULL,
  `access_groups_id` int unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_permission_pattern_access_group` (`permission_pattern_id`,`access_groups_id`),
  KEY `access_groups_id` (`access_groups_id`),
  CONSTRAINT `permission_pattern_access_group_ibfk_1` FOREIGN KEY (`permission_pattern_id`) REFERENCES `permission_pattern` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `permission_pattern_access_group_ibfk_2` FOREIGN KEY (`access_groups_id`) REFERENCES `access_groups` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1525 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `pp_temp`
--

DROP TABLE IF EXISTS `pp_temp`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pp_temp` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `pp_id` int unsigned NOT NULL,
  `lg_id` int unsigned NOT NULL,
  `hash_id` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=50692 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `recording_database_export`
--

DROP TABLE IF EXISTS `recording_database_export`;
/*!50001 DROP VIEW IF EXISTS `recording_database_export`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `recording_database_export` (
  `iso` tinyint NOT NULL,
  `lang_name` tinyint NOT NULL,
  `bible_name` tinyint NOT NULL,
  `bible_id` tinyint NOT NULL,
  `population` tinyint NOT NULL,
  `country` tinyint NOT NULL,
  `video_date` tinyint NOT NULL,
  `audio_date` tinyint NOT NULL,
  `text_date` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `resource_connections`
--

DROP TABLE IF EXISTS `resource_connections`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `resource_connections` (
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `resource_links`
--

DROP TABLE IF EXISTS `resource_links`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `resource_links` (
  `resource_id` int unsigned NOT NULL,
  `title` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `size` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `url` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `resource_links_resource_id_foreign` (`resource_id`),
  CONSTRAINT `FK_resources_resource_links` FOREIGN KEY (`resource_id`) REFERENCES `resources` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `resource_translations`
--

DROP TABLE IF EXISTS `resource_translations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `resource_translations` (
  `language_id` int unsigned NOT NULL,
  `resource_id` int unsigned NOT NULL,
  `vernacular` tinyint(1) NOT NULL,
  `tag` tinyint(1) NOT NULL,
  `title` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`language_id`,`resource_id`),
  KEY `resource_translations_language_id_foreign` (`language_id`),
  KEY `resource_translations_resource_id_foreign` (`resource_id`),
  CONSTRAINT `FK_languages_resource_translations` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_resources_resource_translations` FOREIGN KEY (`resource_id`) REFERENCES `resources` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `resources`
--

DROP TABLE IF EXISTS `resources`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `resources` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `language_id` int unsigned NOT NULL,
  `slug` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `iso` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `organization_id` int unsigned NOT NULL,
  `source_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cover` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `cover_thumbnail` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `date` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `resources_language_id_foreign` (`language_id`),
  KEY `resources_organization_id_foreign` (`organization_id`),
  KEY `resources_iso_index` (`iso`),
  CONSTRAINT `FK_languages_resources_iso` FOREIGN KEY (`iso`) REFERENCES `languages` (`iso`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_languages_resources_language_id` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_organizations_resources` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=149095 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `run_batch`
--

DROP TABLE IF EXISTS `run_batch`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `run_batch` (
  `run_id` int unsigned NOT NULL,
  `batch` varchar(64) NOT NULL,
  `status` tinyint NOT NULL,
  KEY `run_id` (`run_id`),
  CONSTRAINT `run_batch_ibfk_1` FOREIGN KEY (`run_id`) REFERENCES `run_history` (`run_id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `run_history`
--

DROP TABLE IF EXISTS `run_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `run_history` (
  `run_id` int unsigned NOT NULL AUTO_INCREMENT,
  `run_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `username` varchar(64) NOT NULL,
  `location` varchar(64) NOT NULL,
  `directory` varchar(256) NOT NULL,
  `errors` text,
  PRIMARY KEY (`run_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8269 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `stocknumber_lookup`
--

DROP TABLE IF EXISTS `stocknumber_lookup`;
/*!50001 DROP VIEW IF EXISTS `stocknumber_lookup`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `stocknumber_lookup` (
  `stocknumber` tinyint NOT NULL,
  `bibleid` tinyint NOT NULL,
  `mode` tinyint NOT NULL,
  `language` tinyint NOT NULL,
  `version` tinyint NOT NULL,
  `licensor` tinyint NOT NULL,
  `iso` tinyint NOT NULL,
  `rolv` tinyint NOT NULL,
  `languageId` tinyint NOT NULL,
  `countryId` tinyint NOT NULL,
  `organization_id` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `tblataglance`
--

DROP TABLE IF EXISTS `tblataglance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tblataglance` (
  `StockID` int NOT NULL,
  `StockNo` varchar(15) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `LangName` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Pseudonym` tinyint(1) DEFAULT NULL,
  `EthName` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `AltName` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `MasterTradeLanguage` varchar(10) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `HeartName` varchar(60) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `ISO` varchar(3) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Version` varchar(110) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Country` varchar(35) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `CountryAdditional` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `Region` varchar(25) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Location` varchar(75) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `RegionalManager` varchar(25) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Population` int DEFAULT NULL,
  `RecordingStatus` varchar(25) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `TargetDate` int DEFAULT NULL,
  `Portion` varchar(15) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Selection` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `Note` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci,
  `SecurityRisk` varchar(175) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `SelectionStatus` varchar(22) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Completed Selections` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Scheduled Selections` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Audience` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL,
  `Product_Description` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `temp_hash_fix`
--

DROP TABLE IF EXISTS `temp_hash_fix`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `temp_hash_fix` (
  `id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `bad_hash` char(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `correct_hash` varchar(12) NOT NULL DEFAULT '',
  `asset_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `set_type_code` varchar(18) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `set_size_code` char(9) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `archived` tinyint(1) NOT NULL DEFAULT '0',
  `content_loaded` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `temp_language_equivalents`
--

DROP TABLE IF EXISTS `temp_language_equivalents`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `temp_language_equivalents` (
  `preferred_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `duplicate_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`duplicate_id`,`preferred_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `temp_language_equivalents_full`
--

DROP TABLE IF EXISTS `temp_language_equivalents_full`;
/*!50001 DROP VIEW IF EXISTS `temp_language_equivalents_full`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `temp_language_equivalents_full` (
  `iso` tinyint NOT NULL,
  `country_id1` tinyint NOT NULL,
  `id1` tinyint NOT NULL,
  `fcbh_name1` tinyint NOT NULL,
  `heart_name1` tinyint NOT NULL,
  `ethnologue_name1` tinyint NOT NULL,
  `official_name1` tinyint NOT NULL,
  `rolv_code1` tinyint NOT NULL,
  `pseudonym1` tinyint NOT NULL,
  `sensitivity1` tinyint NOT NULL,
  `country_id2` tinyint NOT NULL,
  `id2` tinyint NOT NULL,
  `fcbh_name2` tinyint NOT NULL,
  `heart_name2` tinyint NOT NULL,
  `ethnologue_name2` tinyint NOT NULL,
  `official_name2` tinyint NOT NULL,
  `rolv_code2` tinyint NOT NULL,
  `pseudonym2` tinyint NOT NULL,
  `sensitivity2` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `temp_stocknumber_mode_with_multiple_licensors`
--

DROP TABLE IF EXISTS `temp_stocknumber_mode_with_multiple_licensors`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `temp_stocknumber_mode_with_multiple_licensors` (
  `stock_no` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `mode` varchar(18) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `test`
--

DROP TABLE IF EXISTS `test`;
/*!50001 DROP VIEW IF EXISTS `test`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `test` (
  `iso` tinyint NOT NULL,
  `lang_name` tinyint NOT NULL,
  `bible_name` tinyint NOT NULL,
  `bible_id` tinyint NOT NULL,
  `population` tinyint NOT NULL,
  `country` tinyint NOT NULL,
  `video_date` tinyint NOT NULL,
  `audio_date` tinyint NOT NULL,
  `text_date` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `test_bibles`
--

DROP TABLE IF EXISTS `test_bibles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `test_bibles` (
  `id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `language_id` int unsigned NOT NULL,
  `versification` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'protestant',
  `numeral_system_id` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `date` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `scope` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `script` char(4) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `derived` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `copyright` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `priority` tinyint unsigned NOT NULL DEFAULT '0',
  `reviewed` tinyint(1) DEFAULT '0',
  `notes` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `language_id` (`language_id`),
  KEY `numeral_system_id` (`numeral_system_id`),
  KEY `script` (`script`),
  CONSTRAINT `test_bibles_ibfk_1` FOREIGN KEY (`script`) REFERENCES `alphabets` (`script`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `test_bibles_ibfk_3` FOREIGN KEY (`numeral_system_id`) REFERENCES `numeral_systems` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `text_storage_view`
--

DROP TABLE IF EXISTS `text_storage_view`;
/*!50001 DROP VIEW IF EXISTS `text_storage_view`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `text_storage_view` (
  `bucket` tinyint NOT NULL,
  `type` tinyint NOT NULL,
  `set_type_code` tinyint NOT NULL,
  `bible` tinyint NOT NULL,
  `fileset` tinyint NOT NULL,
  `file_name` tinyint NOT NULL,
  `path` tinyint NOT NULL,
  `size` tinyint NOT NULL,
  `last_modified_date` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `translator_relations`
--

DROP TABLE IF EXISTS `translator_relations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `translator_relations` (
  `translator_id` int unsigned NOT NULL,
  `translator_relation_id` int unsigned NOT NULL,
  `organization_id` int unsigned DEFAULT NULL,
  `type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `notes` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `translator_relations_translator_id_foreign` (`translator_id`),
  KEY `translator_relations_translator_relation_id_foreign` (`translator_relation_id`),
  KEY `translator_relations_organization_id_foreign` (`organization_id`),
  CONSTRAINT `FK_organizations_translator_relations` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `FK_translators_translator_relations` FOREIGN KEY (`translator_id`) REFERENCES `translators` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_translators_translator_relations_translator_relation_id` FOREIGN KEY (`translator_relation_id`) REFERENCES `translators` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `translators`
--

DROP TABLE IF EXISTS `translators`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `translators` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `born` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `died` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `version`
--

DROP TABLE IF EXISTS `version`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `version` (
  `id` char(3) NOT NULL,
  `name` varchar(127) NOT NULL,
  `english_name` varchar(127) NOT NULL,
  UNIQUE KEY `id` (`id`),
  FULLTEXT KEY `ft_index_version_name_english_name` (`name`,`english_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `video_sources`
--

DROP TABLE IF EXISTS `video_sources`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `video_sources` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `video_id` int unsigned DEFAULT NULL,
  `url` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `encoding` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `resolution` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `size` int NOT NULL,
  `url_type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `video_sources_video_id_foreign` (`video_id`),
  CONSTRAINT `FK_videos_video_sources` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary table structure for view `video_storage_view1`
--

DROP TABLE IF EXISTS `video_storage_view1`;
/*!50001 DROP VIEW IF EXISTS `video_storage_view1`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `video_storage_view1` (
  `bucket` tinyint NOT NULL,
  `type` tinyint NOT NULL,
  `set_type_code` tinyint NOT NULL,
  `bible` tinyint NOT NULL,
  `fileset` tinyint NOT NULL,
  `file_name` tinyint NOT NULL,
  `path` tinyint NOT NULL,
  `size` tinyint NOT NULL,
  `last_modified_date` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Temporary table structure for view `video_storage_view2`
--

DROP TABLE IF EXISTS `video_storage_view2`;
/*!50001 DROP VIEW IF EXISTS `video_storage_view2`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `video_storage_view2` (
  `bucket` tinyint NOT NULL,
  `type` tinyint NOT NULL,
  `set_type_code` tinyint NOT NULL,
  `bible` tinyint NOT NULL,
  `fileset` tinyint NOT NULL,
  `file_name` tinyint NOT NULL,
  `path` tinyint NOT NULL,
  `size` tinyint NOT NULL,
  `last_modified_date` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Temporary table structure for view `video_storage_view3`
--

DROP TABLE IF EXISTS `video_storage_view3`;
/*!50001 DROP VIEW IF EXISTS `video_storage_view3`*/;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE TABLE `video_storage_view3` (
  `bucket` tinyint NOT NULL,
  `type` tinyint NOT NULL,
  `set_type_code` tinyint NOT NULL,
  `bible` tinyint NOT NULL,
  `fileset` tinyint NOT NULL,
  `file_name` tinyint NOT NULL,
  `path` tinyint NOT NULL,
  `size` tinyint NOT NULL,
  `last_modified_date` tinyint NOT NULL
) ENGINE=MyISAM */;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `video_tags`
--

DROP TABLE IF EXISTS `video_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `video_tags` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `video_id` int unsigned DEFAULT NULL,
  `category` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `tag_type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `tag` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `language_id` int unsigned DEFAULT NULL,
  `organization_id` int unsigned DEFAULT NULL,
  `book_id` char(3) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `chapter_start` int unsigned DEFAULT NULL,
  `chapter_end` int unsigned DEFAULT NULL,
  `verse_start` int unsigned DEFAULT NULL,
  `verse_end` int unsigned DEFAULT NULL,
  `time_begin` double(8,2) unsigned DEFAULT NULL,
  `time_end` double(8,2) unsigned DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `video_tags_video_id_foreign` (`video_id`),
  KEY `video_tags_language_id_foreign` (`language_id`),
  KEY `video_tags_organization_id_foreign` (`organization_id`),
  KEY `video_tags_book_id_foreign` (`book_id`),
  CONSTRAINT `FK_books_video_tags` FOREIGN KEY (`book_id`) REFERENCES `books` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `FK_languages_video_tags` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_organizations_video_tags` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `FK_videos_video_tags` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `video_translations`
--

DROP TABLE IF EXISTS `video_translations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `video_translations` (
  `language_id` int unsigned NOT NULL AUTO_INCREMENT,
  `video_id` int unsigned NOT NULL,
  `title` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`language_id`),
  KEY `video_translations_video_id_foreign` (`video_id`),
  CONSTRAINT `FK_languages_video_translations` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_videos_video_translations` FOREIGN KEY (`video_id`) REFERENCES `videos` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `videos`
--

DROP TABLE IF EXISTS `videos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `videos` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `language_id` int unsigned DEFAULT NULL,
  `bible_id` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `series` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `episode` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `section` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `picture` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `duration` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `foo` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `videos_language_id_foreign` (`language_id`),
  KEY `videos_bible_id_foreign` (`bible_id`),
  CONSTRAINT `FK_bibles_videos` FOREIGN KEY (`bible_id`) REFERENCES `bibles` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_languages_videos` FOREIGN KEY (`language_id`) REFERENCES `languages` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Final view structure for view `access_group_filesets_view`
--

/*!50001 DROP TABLE IF EXISTS `access_group_filesets_view`*/;
/*!50001 DROP VIEW IF EXISTS `access_group_filesets_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `access_group_filesets_view` AS select `bf`.`id` AS `fileset_id`,`bf`.`asset_id` AS `asset_id`,`bf`.`set_type_code` AS `set_type_code`,`bf`.`hash_id` AS `hash_id`,`agf`.`access_group_id` AS `access_group_id` from (`bible_filesets` `bf` join `access_group_filesets` `agf`) where (`bf`.`hash_id` = `agf`.`hash_id`) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `all_stocknumbers`
--

/*!50001 DROP TABLE IF EXISTS `all_stocknumbers`*/;
/*!50001 DROP VIEW IF EXISTS `all_stocknumbers`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `all_stocknumbers` AS select distinct `bft`.`description` AS `stocknumber`,`b`.`id` AS `bible_id`,`l`.`name` AS `language`,`bt`.`name` AS `version`,`l`.`iso` AS `iso`,`l`.`rolv_code` AS `rolv`,`l`.`id` AS `language_id`,`l`.`country_id` AS `country_id`,`lt`.`name` AS `fcbh_name`,`lt`.`language_source_id` AS `fcbh_language_id`,`c`.`name` AS `country_name`,(select group_concat(distinct `lt`.`name` separator ' | ') AS `name` from `language_translations` `lt` where ((`lt`.`priority` = 9) and (`lt`.`language_source_id` = `l`.`id`) and (`lt`.`language_translation_id` = `lt`.`language_source_id`))) AS `heart_name`,`bft_pseud`.`description` AS `pseudonym` from ((((((((`bible_filesets` `bf` join `bible_fileset_connections` `bfc` on((`bfc`.`hash_id` = `bf`.`hash_id`))) join `bibles` `b` on((`b`.`id` = `bfc`.`bible_id`))) join `languages` `l` on((`l`.`id` = `b`.`language_id`))) join `language_translations` `lt` on(((`lt`.`language_source_id` = `l`.`id`) and (`lt`.`language_translation_id` = 6414) and (`lt`.`priority` = 9)))) join `countries` `c` on((`c`.`id` = `l`.`country_id`))) join `bible_fileset_tags` `bft` on((`bft`.`hash_id` = `bf`.`hash_id`))) join `bible_translations` `bt` on((`bt`.`bible_id` = `b`.`id`))) left join `bible_fileset_tags` `bft_pseud` on(((`bft_pseud`.`hash_id` = `bft`.`hash_id`) and (`bft_pseud`.`name` = 'pseudonym')))) where ((`bt`.`language_id` = 6414) and (`bft`.`name` = 'stock_no')) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `audio_storage_view`
--

/*!50001 DROP TABLE IF EXISTS `audio_storage_view`*/;
/*!50001 DROP VIEW IF EXISTS `audio_storage_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `audio_storage_view` AS select `bfs`.`asset_id` AS `bucket`,substring_index(`bfs`.`set_type_code`,'_',1) AS `type`,`bfs`.`set_type_code` AS `set_type_code`,`bfc`.`bible_id` AS `bible`,`bfs`.`id` AS `fileset`,`bf`.`file_name` AS `file_name`,concat('audio/',`bfc`.`bible_id`,'/',`bfs`.`id`,'/',`bf`.`file_name`) AS `path`,`bf`.`file_size` AS `size`,`bf`.`updated_at` AS `last_modified_date` from ((`bible_files` `bf` join `bible_filesets` `bfs` on((`bf`.`hash_id` = `bfs`.`hash_id`))) join `bible_fileset_connections` `bfc` on((`bfs`.`hash_id` = `bfc`.`hash_id`))) where ((`bfs`.`asset_id` = 'dbp-prod') and (substring_index(`bfs`.`set_type_code`,'_',1) = 'audio') and (not((`bf`.`file_name` like '%m3u8')))) order by `bfs`.`asset_id`,substring_index(`bfs`.`set_type_code`,'_',1) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `bible_file_stream_bandwidths_view`
--

/*!50001 DROP TABLE IF EXISTS `bible_file_stream_bandwidths_view`*/;
/*!50001 DROP VIEW IF EXISTS `bible_file_stream_bandwidths_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `bible_file_stream_bandwidths_view` AS select `bf`.`id` AS `fileset_id`,`bf`.`asset_id` AS `asset_id`,`bf`.`set_type_code` AS `set_type_code`,`bf`.`hash_id` AS `hash_id`,`bf2`.`id` AS `file_id`,`bfvr`.`file_name` AS `file_name`,`bfvr`.`bandwidth` AS `bandwidth`,`bfvr`.`resolution_width` AS `resolution_width`,`bfvr`.`resolution_height` AS `resolution_height`,`bfvr`.`codec` AS `codec`,`bfvr`.`stream` AS `stream` from ((`bible_filesets` `bf` join `bible_files` `bf2`) join `bible_file_stream_bandwidths` `bfvr`) where ((`bf`.`hash_id` = `bf2`.`hash_id`) and (`bf2`.`id` = `bfvr`.`bible_file_id`)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `bible_file_stream_ts_view`
--

/*!50001 DROP TABLE IF EXISTS `bible_file_stream_ts_view`*/;
/*!50001 DROP VIEW IF EXISTS `bible_file_stream_ts_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `bible_file_stream_ts_view` AS select `bf`.`id` AS `fileset_id`,`bf`.`asset_id` AS `asset_id`,`bf`.`set_type_code` AS `set_type_code`,`bf`.`hash_id` AS `hash_id`,`bf2`.`id` AS `file_id`,`bfvts`.`stream_bandwidth_id` AS `stream_bandwidth_id`,`bfvts`.`file_name` AS `file_name`,`bfvts`.`runtime` AS `runtime` from (((`bible_filesets` `bf` join `bible_files` `bf2`) join `bible_file_stream_bandwidths` `bfvr`) join `bible_file_stream_ts` `bfvts`) where ((`bf`.`hash_id` = `bf2`.`hash_id`) and (`bf2`.`id` = `bfvr`.`bible_file_id`) and (`bfvr`.`id` = `bfvts`.`stream_bandwidth_id`)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `bible_file_tags_view`
--

/*!50001 DROP TABLE IF EXISTS `bible_file_tags_view`*/;
/*!50001 DROP VIEW IF EXISTS `bible_file_tags_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `bible_file_tags_view` AS select `bf`.`id` AS `fileset_id`,`bf`.`asset_id` AS `asset_id`,`bf`.`set_type_code` AS `set_type_code`,`bf`.`hash_id` AS `hash_id`,`bft`.`file_id` AS `file_id`,`bft`.`tag` AS `tag`,`bft`.`value` AS `value`,`bft`.`admin_only` AS `admin_only` from ((`bible_filesets` `bf` join `bible_files` `bf2`) join `bible_file_tags` `bft`) where ((`bf`.`hash_id` = `bf2`.`hash_id`) and (`bft`.`file_id` = `bf2`.`id`)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `bible_file_video_resolutions`
--

/*!50001 DROP TABLE IF EXISTS `bible_file_video_resolutions`*/;
/*!50001 DROP VIEW IF EXISTS `bible_file_video_resolutions`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `bible_file_video_resolutions` AS select `bible_file_stream_bandwidths`.`id` AS `id`,`bible_file_stream_bandwidths`.`bible_file_id` AS `bible_file_id`,`bible_file_stream_bandwidths`.`file_name` AS `file_name`,`bible_file_stream_bandwidths`.`bandwidth` AS `bandwidth`,`bible_file_stream_bandwidths`.`resolution_width` AS `resolution_width`,`bible_file_stream_bandwidths`.`resolution_height` AS `resolution_height`,`bible_file_stream_bandwidths`.`codec` AS `codec`,`bible_file_stream_bandwidths`.`stream` AS `stream`,`bible_file_stream_bandwidths`.`created_at` AS `created_at`,`bible_file_stream_bandwidths`.`updated_at` AS `updated_at` from `bible_file_stream_bandwidths` */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `bible_file_video_transport_stream`
--

/*!50001 DROP TABLE IF EXISTS `bible_file_video_transport_stream`*/;
/*!50001 DROP VIEW IF EXISTS `bible_file_video_transport_stream`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `bible_file_video_transport_stream` AS select `bible_file_stream_ts`.`id` AS `id`,`bible_file_stream_ts`.`stream_bandwidth_id` AS `video_resolution_id`,`bible_file_stream_ts`.`file_name` AS `file_name`,`bible_file_stream_ts`.`runtime` AS `runtime`,`bible_file_stream_ts`.`created_at` AS `created_at`,`bible_file_stream_ts`.`updated_at` AS `updated_at` from `bible_file_stream_ts` */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `bible_files_view`
--

/*!50001 DROP TABLE IF EXISTS `bible_files_view`*/;
/*!50001 DROP VIEW IF EXISTS `bible_files_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `bible_files_view` AS select `bf`.`id` AS `fileset_id`,`bf`.`asset_id` AS `asset_id`,`bf`.`set_type_code` AS `set_type_code`,`bf`.`hash_id` AS `hash_id`,`bf2`.`book_id` AS `book_id`,`bf2`.`chapter_start` AS `chapter_start`,`bf2`.`chapter_end` AS `chapter_end`,`bf2`.`verse_start` AS `verse_start`,`bf2`.`verse_end` AS `verse_end`,`bf2`.`file_name` AS `file_name`,`bf2`.`file_size` AS `file_size`,`bf2`.`duration` AS `duration` from (`bible_filesets` `bf` join `bible_files` `bf2`) where (`bf`.`hash_id` = `bf2`.`hash_id`) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `bible_files_view2`
--

/*!50001 DROP TABLE IF EXISTS `bible_files_view2`*/;
/*!50001 DROP VIEW IF EXISTS `bible_files_view2`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `bible_files_view2` AS select `bf`.`id` AS `fileset_id`,`bf`.`asset_id` AS `asset_id`,`bf`.`set_type_code` AS `set_type_code`,`bf`.`hash_id` AS `hash_id`,`bf2`.`book_id` AS `book_id`,`bf2`.`chapter_start` AS `chapter_start`,`bf2`.`chapter_end` AS `chapter_end`,`bf2`.`verse_start` AS `verse_start`,`bf2`.`verse_end` AS `verse_end`,`bf2`.`file_name` AS `file_name`,`bf2`.`file_size` AS `file_size`,`bf2`.`duration` AS `duration`,`bf2`.`created_at` AS `created_at`,`bf2`.`updated_at` AS `updated_at` from (`bible_filesets` `bf` join `bible_files` `bf2`) where (`bf`.`hash_id` = `bf2`.`hash_id`) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `bible_fileset_connections_view`
--

/*!50001 DROP TABLE IF EXISTS `bible_fileset_connections_view`*/;
/*!50001 DROP VIEW IF EXISTS `bible_fileset_connections_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `bible_fileset_connections_view` AS select `bfc`.`bible_id` AS `bible_id`,`bfs`.`id` AS `id`,substring_index(`bfs`.`set_type_code`,'_',1) AS `mode`,`bfs`.`set_type_code` AS `type`,`bfs`.`set_size_code` AS `size`,`bft`.`description` AS `stockno`,`bfs`.`hash_id` AS `hash_id`,`bfs`.`content_loaded` AS `content_loaded`,`bfs`.`archived` AS `archived`,`bfs`.`updated_at` AS `updated_at` from ((`bible_fileset_connections` `bfc` left join `bible_filesets` `bfs` on((`bfs`.`hash_id` = `bfc`.`hash_id`))) left join `bible_fileset_tags` `bft` on(((`bft`.`hash_id` = `bfc`.`hash_id`) and (`bft`.`name` = 'stock_no')))) order by `bfc`.`bible_id`,`bfs`.`set_type_code` */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `bible_fileset_copyright_organizations_view`
--

/*!50001 DROP TABLE IF EXISTS `bible_fileset_copyright_organizations_view`*/;
/*!50001 DROP VIEW IF EXISTS `bible_fileset_copyright_organizations_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `bible_fileset_copyright_organizations_view` AS select `bf`.`id` AS `fileset_id`,`bf`.`asset_id` AS `asset_id`,`bf`.`set_type_code` AS `set_type_code`,`bf`.`hash_id` AS `hash_id`,`bfco`.`organization_id` AS `organization_id`,`bfco`.`organization_role` AS `organization_role` from (`bible_filesets` `bf` join `bible_fileset_copyright_organizations` `bfco`) where (`bf`.`hash_id` = `bfco`.`hash_id`) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `bible_fileset_copyrights_view`
--

/*!50001 DROP TABLE IF EXISTS `bible_fileset_copyrights_view`*/;
/*!50001 DROP VIEW IF EXISTS `bible_fileset_copyrights_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `bible_fileset_copyrights_view` AS select `bf`.`id` AS `fileset_id`,`bf`.`asset_id` AS `asset_id`,`bf`.`set_type_code` AS `set_type_code`,`bf`.`hash_id` AS `hash_id`,`bft`.`description` AS `stockno`,`bfc`.`copyright_date` AS `copyright_date`,`bfc`.`copyright` AS `copyright`,`bfc`.`is_combined` AS `is_combined` from ((`bible_filesets` `bf` join `bible_fileset_copyrights` `bfc`) join `bible_fileset_tags` `bft` on(((`bft`.`hash_id` = `bf`.`hash_id`) and (`bft`.`name` = 'stock_no')))) where ((`bf`.`hash_id` = `bfc`.`hash_id`) and (`bf`.`archived` = 0)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `bible_fileset_lookup`
--

/*!50001 DROP TABLE IF EXISTS `bible_fileset_lookup`*/;
/*!50001 DROP VIEW IF EXISTS `bible_fileset_lookup`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `bible_fileset_lookup` AS select distinct `bft`.`description` AS `stocknumber`,`b`.`id` AS `bibleid`,`bfs`.`id` AS `filesetid`,`bfs`.`set_type_code` AS `type`,substring_index(`bfs`.`set_type_code`,'_',1) AS `mode`,`bfs`.`set_size_code` AS `size`,`bfs`.`hash_id` AS `hash_id`,`l`.`name` AS `language`,`bt`.`name` AS `version`,`ot`.`name` AS `licensor`,`bfs`.`content_loaded` AS `content_loaded`,`bfs`.`archived` AS `archived` from ((((((((`bible_filesets` `bfs` join `bible_fileset_connections` `bfc` on((`bfc`.`hash_id` = `bfs`.`hash_id`))) join `bible_fileset_copyright_organizations` `bfco` on((`bfco`.`hash_id` = `bfs`.`hash_id`))) join `organizations` `o` on((`o`.`id` = `bfco`.`organization_id`))) join `bibles` `b` on((`b`.`id` = `bfc`.`bible_id`))) join `languages` `l` on((`l`.`id` = `b`.`language_id`))) join `bible_fileset_tags` `bft` on((`bfs`.`hash_id` = `bft`.`hash_id`))) join `organization_translations` `ot` on((`ot`.`organization_id` = `bfco`.`organization_id`))) join `bible_translations` `bt` on((`bt`.`bible_id` = `b`.`id`))) where ((`bt`.`language_id` = 6414) and (`ot`.`language_id` = 6414) and (`bfco`.`organization_role` = 2) and (`bft`.`name` = 'stock_no')) order by `filesetid` */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `bible_fileset_lookup2`
--

/*!50001 DROP TABLE IF EXISTS `bible_fileset_lookup2`*/;
/*!50001 DROP VIEW IF EXISTS `bible_fileset_lookup2`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `bible_fileset_lookup2` AS select distinct `bft`.`description` AS `stocknumber`,`b`.`id` AS `bibleid`,`bfs`.`id` AS `filesetid`,`bfs`.`set_type_code` AS `type`,`bfs`.`hash_id` AS `hash_id`,`l`.`name` AS `language`,`bt`.`name` AS `version`,`ot`.`name` AS `licensor`,`l`.`iso` AS `iso`,`l`.`rolv_code` AS `rolv`,`l`.`id` AS `languageId`,`l`.`country_id` AS `countryId` from ((((((((`bible_filesets` `bfs` join `bible_fileset_connections` `bfc` on((`bfc`.`hash_id` = `bfs`.`hash_id`))) join `bible_fileset_copyright_organizations` `bfco` on((`bfco`.`hash_id` = `bfs`.`hash_id`))) join `organizations` `o` on((`o`.`id` = `bfco`.`organization_id`))) join `bibles` `b` on((`b`.`id` = `bfc`.`bible_id`))) join `languages` `l` on((`l`.`id` = `b`.`language_id`))) join `bible_fileset_tags` `bft` on((`bfs`.`hash_id` = `bft`.`hash_id`))) join `organization_translations` `ot` on((`ot`.`organization_id` = `bfco`.`organization_id`))) join `bible_translations` `bt` on((`bt`.`bible_id` = `b`.`id`))) where ((`bt`.`language_id` = 6414) and (`ot`.`language_id` = 6414) and (`bfco`.`organization_role` = 2) and (`bft`.`name` = 'stock_no')) order by `filesetid` */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `bible_fileset_tags_view`
--

/*!50001 DROP TABLE IF EXISTS `bible_fileset_tags_view`*/;
/*!50001 DROP VIEW IF EXISTS `bible_fileset_tags_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `bible_fileset_tags_view` AS select `bf`.`id` AS `fileset_id`,`bf`.`asset_id` AS `asset_id`,`bf`.`set_type_code` AS `set_type_code`,`bf`.`hash_id` AS `hash_id`,`bft`.`name` AS `name`,`bft`.`description` AS `description`,`bft`.`admin_only` AS `admin_only`,`bft`.`notes` AS `notes`,`bft`.`iso` AS `iso`,`bft`.`language_id` AS `language_id`,`bf`.`content_loaded` AS `content_loaded`,`bf`.`archived` AS `archived` from (`bible_filesets` `bf` join `bible_fileset_tags` `bft`) where (`bf`.`hash_id` = `bft`.`hash_id`) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `bible_filesets_view`
--

/*!50001 DROP TABLE IF EXISTS `bible_filesets_view`*/;
/*!50001 DROP VIEW IF EXISTS `bible_filesets_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `bible_filesets_view` AS select `bible_filesets`.`id` AS `fileset_id`,`bible_filesets`.`hash_id` AS `hash_id`,`bible_filesets`.`asset_id` AS `asset_id`,`bible_filesets`.`set_type_code` AS `set_type_code`,`bible_filesets`.`set_size_code` AS `set_size_code`,`bible_filesets`.`hidden` AS `hidden` from `bible_filesets` */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `biblebrain_database_export`
--

/*!50001 DROP TABLE IF EXISTS `biblebrain_database_export`*/;
/*!50001 DROP VIEW IF EXISTS `biblebrain_database_export`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `biblebrain_database_export` AS select `A`.`iso` AS `iso`,`A`.`lang_name` AS `lang_name`,`A`.`bible_name` AS `bible_name`,`A`.`bible_id` AS `bible_id`,`A`.`population` AS `population`,`A`.`country` AS `country`,max(`A`.`video_date`) AS `video_date`,max(`A`.`audio_date`) AS `audio_date`,max(`A`.`text_date`) AS `text_date` from (select `l`.`iso` AS `iso`,`lt`.`name` AS `lang_name`,`l`.`population` AS `population`,`c`.`country` AS `country`,`bfc`.`bible_id` AS `bible_id`,`bf`.`set_type_code` AS `media_type`,`bt`.`name` AS `bible_name`,(case when (`bf`.`set_type_code` like 'video%') then cast(`bf`.`created_at` as date) end) AS `video_date`,(case when (`bf`.`set_type_code` like 'audio%') then cast(`bf`.`created_at` as date) end) AS `audio_date`,(case when (`bf`.`set_type_code` like 'text%') then cast(`bf`.`created_at` as date) end) AS `text_date` from ((((((((select `bible_filesets`.`set_type_code` AS `set_type_code`,`bible_filesets`.`created_at` AS `created_at`,`bible_filesets`.`hash_id` AS `hash_id` from `bible_filesets` where ((`bible_filesets`.`content_loaded` = 1) and (`bible_filesets`.`archived` = 0))) `bf` join (select `bible_fileset_connections`.`bible_id` AS `bible_id`,`bible_fileset_connections`.`hash_id` AS `hash_id` from `bible_fileset_connections`) `bfc` on((`bf`.`hash_id` = `bfc`.`hash_id`))) join (select `bibles`.`id` AS `id`,`bibles`.`language_id` AS `language_id` from `bibles`) `b` on((`bfc`.`bible_id` = `b`.`id`))) join (select `languages`.`id` AS `id`,`languages`.`name` AS `name`,`languages`.`iso` AS `iso`,`languages`.`population` AS `population`,`languages`.`country_id` AS `country_id` from `languages`) `l` on((`l`.`id` = `b`.`language_id`))) join (select `access_group_filesets`.`hash_id` AS `hash_id`,`access_group_filesets`.`access_group_id` AS `access_group_id` from `access_group_filesets` where (`access_group_filesets`.`access_group_id` in (121,123,125,191,193))) `agf` on((`bf`.`hash_id` = `agf`.`hash_id`))) join (select `bible_translations`.`name` AS `name`,`bible_translations`.`bible_id` AS `bible_id`,`bible_translations`.`vernacular` AS `vernacular` from `bible_translations` where (`bible_translations`.`vernacular` = 0)) `bt` on((`bt`.`bible_id` = `b`.`id`))) join (select `language_translations`.`name` AS `name`,`language_translations`.`priority` AS `priority`,`language_translations`.`language_translation_id` AS `language_translation_id`,`language_translations`.`language_source_id` AS `language_source_id` from `language_translations` where ((`language_translations`.`priority` = 9) and (`language_translations`.`language_translation_id` = 6414))) `lt` on((`l`.`id` = `lt`.`language_source_id`))) left join (select `countries`.`id` AS `country_id`,`countries`.`name` AS `country` from `countries`) `c` on((`l`.`country_id` = `c`.`country_id`)))) `A` group by `A`.`lang_name`,`A`.`bible_id` order by `video_date` desc */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `bibles_view`
--

/*!50001 DROP TABLE IF EXISTS `bibles_view`*/;
/*!50001 DROP VIEW IF EXISTS `bibles_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `bibles_view` AS select `bf`.`id` AS `fileset_id`,`b`.`id` AS `bible_id`,`b`.`language_id` AS `language_id`,`b`.`versification` AS `versification`,`b`.`numeral_system_id` AS `numeral_system_id`,`b`.`date` AS `date`,`b`.`scope` AS `scope`,`b`.`script` AS `script`,`b`.`derived` AS `derived`,`b`.`copyright` AS `copyright`,`b`.`priority` AS `priority`,`b`.`reviewed` AS `reviewed`,`b`.`notes` AS `notes` from ((`bibles` `b` join `bible_filesets` `bf`) join `bible_fileset_connections` `bfc`) where ((`b`.`id` = `bfc`.`bible_id`) and (`bfc`.`hash_id` = `bf`.`hash_id`)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `fileset_permission_view`
--

/*!50001 DROP TABLE IF EXISTS `fileset_permission_view`*/;
/*!50001 DROP VIEW IF EXISTS `fileset_permission_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `fileset_permission_view` AS select `bfs`.`id` AS `fileset_id`,`bfs`.`hash_id` AS `hash_id`,`lg`.`id` AS `license_group_id`,`lg`.`name` AS `license_group`,`lg`.`bible_id` AS `bible_id`,`lg`.`mode` AS `mode`,`pp`.`id` AS `permission_pattern_id`,`pp`.`name` AS `permission_pattern`,group_concat(`agf`.`access_group_id` separator ',') AS `group_concat(access_group_id)` from (((((`bible_filesets` `bfs` join `license_group_filesets` `lgf` on((`lgf`.`hash_id` = `bfs`.`hash_id`))) join `license_group` `lg` on((`lgf`.`license_group_id` = `lg`.`id`))) join `license_group_permission_pattern` `lgpp` on((`lgpp`.`license_group_id` = `lg`.`id`))) join `permission_pattern` `pp` on((`pp`.`id` = `lgpp`.`permission_pattern_id`))) join `access_group_filesets` `agf` on((`agf`.`hash_id` = `bfs`.`hash_id`))) group by `bfs`.`hash_id` order by `bfs`.`id` */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `language_view`
--

/*!50001 DROP TABLE IF EXISTS `language_view`*/;
/*!50001 DROP VIEW IF EXISTS `language_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `language_view` AS select `l`.`id` AS `id`,`l`.`iso` AS `iso`,`l`.`country_id` AS `country_id`,`c`.`name` AS `country_name`,`lt`.`name` AS `fcbh_name`,`lt2`.`name` AS `heart_name`,`lt3`.`name` AS `ethnologue_name`,(select group_concat(distinct `lt`.`name` separator ' | ') AS `name` from `language_translations` `lt` where ((`lt`.`priority` = 0) and (`lt`.`language_source_id` = `l`.`id`) and (`lt`.`language_translation_id` = 8012))) AS `alt_name`,`l`.`name` AS `official_name`,`l`.`rolv_code` AS `rolv_code`,`l`.`pseudonym` AS `pseudonym`,`l`.`sensitivity` AS `sensitivity`,`l`.`deleted_at` AS `deleted_at` from ((((`languages` `l` join `countries` `c` on((`c`.`id` = `l`.`country_id`))) join `language_translations` `lt` on(((`lt`.`language_source_id` = `l`.`id`) and (`lt`.`language_translation_id` = 6414) and (`lt`.`priority` = 9)))) left join `language_translations` `lt2` on(((`lt2`.`language_source_id` = `l`.`id`) and (`lt2`.`language_translation_id` = `lt2`.`language_source_id`) and (`lt2`.`priority` = 9)))) left join `language_translations` `lt3` on(((`lt3`.`language_source_id` = `l`.`id`) and (`lt3`.`language_translation_id` = 8012) and (`lt3`.`priority` = 5)))) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `license_group_fileset_view`
--

/*!50001 DROP TABLE IF EXISTS `license_group_fileset_view`*/;
/*!50001 DROP VIEW IF EXISTS `license_group_fileset_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `license_group_fileset_view` AS select `lg`.`bible_id` AS `bible_id`,`lg`.`mode` AS `mode`,`lg`.`size` AS `size`,`lg`.`mode_type` AS `mode_type`,`lg`.`name` AS `lg_name`,`bfs`.`id` AS `fileset_id`,`lgf`.`hash_id` AS `hash_id`,`lg`.`id` AS `lg_id`,`lgpp`.`permission_pattern_id` AS `permission_pattern_id`,`lgf`.`updated_at` AS `updated_at` from (((`license_group` `lg` left join `license_group_filesets` `lgf` on((`lgf`.`license_group_id` = `lg`.`id`))) join `bible_filesets` `bfs` on((`bfs`.`hash_id` = `lgf`.`hash_id`))) left join `license_group_permission_pattern` `lgpp` on((`lgpp`.`license_group_id` = `lg`.`id`))) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `license_group_permission_pattern_view`
--

/*!50001 DROP TABLE IF EXISTS `license_group_permission_pattern_view`*/;
/*!50001 DROP VIEW IF EXISTS `license_group_permission_pattern_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `license_group_permission_pattern_view` AS select `lg`.`name` AS `license_group`,`pp`.`name` AS `permission_pattern`,`lgpp`.`license_group_id` AS `license_group_id`,`lgpp`.`permission_pattern_id` AS `permission_pattern_id` from ((`license_group_permission_pattern` `lgpp` join `license_group` `lg` on((`lg`.`id` = `lgpp`.`license_group_id`))) join `permission_pattern` `pp` on((`pp`.`id` = `lgpp`.`permission_pattern_id`))) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `licensor_lookup`
--

/*!50001 DROP TABLE IF EXISTS `licensor_lookup`*/;
/*!50001 DROP VIEW IF EXISTS `licensor_lookup`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `licensor_lookup` AS select distinct `ot`.`name` AS `licensor`,`b`.`id` AS `bible_id`,substring_index(`bf`.`set_type_code`,'_',1) AS `mode`,`bft`.`description` AS `stock_no`,`bf`.`set_size_code` AS `set_size_code`,(case when ((`bf`.`set_size_code` = 'NT') or (`bf`.`set_size_code` = 'NTP')) then 'NT' when ((`bf`.`set_size_code` = 'OT') or (`bf`.`set_size_code` = 'OTP')) then 'OT' else '' end) AS `testament`,`bfcp`.`copyright_date` AS `copyright_date`,`bfcp`.`copyright` AS `copyright`,`bfco`.`organization_id` AS `organization_id` from ((((((`bible_filesets` `bf` join `bible_fileset_connections` `bfc` on(((`bfc`.`hash_id` = `bf`.`hash_id`) and (`bf`.`archived` = 0)))) left join `bible_fileset_copyright_organizations` `bfco` on(((`bfco`.`hash_id` = `bf`.`hash_id`) and (`bfco`.`organization_role` = 2)))) left join `organization_translations` `ot` on(((`ot`.`organization_id` = `bfco`.`organization_id`) and (`ot`.`language_id` = 6414)))) join `bible_fileset_copyrights` `bfcp` on((`bfcp`.`hash_id` = `bf`.`hash_id`))) join `bibles` `b` on((`b`.`id` = `bfc`.`bible_id`))) join `bible_fileset_tags` `bft` on(((`bft`.`hash_id` = `bf`.`hash_id`) and (`bft`.`name` = 'stock_no')))) where (`bf`.`archived` = 0) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `recording_database_export`
--

/*!50001 DROP TABLE IF EXISTS `recording_database_export`*/;
/*!50001 DROP VIEW IF EXISTS `recording_database_export`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `recording_database_export` AS select `A`.`iso` AS `iso`,`A`.`lang_name` AS `lang_name`,`A`.`bible_name` AS `bible_name`,`A`.`bible_id` AS `bible_id`,`A`.`population` AS `population`,`A`.`country` AS `country`,max(`A`.`video_date`) AS `video_date`,max(`A`.`audio_date`) AS `audio_date`,max(`A`.`text_date`) AS `text_date` from (select `l`.`iso` AS `iso`,`lt`.`name` AS `lang_name`,`l`.`population` AS `population`,`c`.`country` AS `country`,`bfc`.`bible_id` AS `bible_id`,`bf`.`set_type_code` AS `media_type`,`bt`.`name` AS `bible_name`,(case when (`bf`.`set_type_code` like 'video%') then cast(`bf`.`created_at` as date) end) AS `video_date`,(case when (`bf`.`set_type_code` like 'audio%') then cast(`bf`.`created_at` as date) end) AS `audio_date`,(case when (`bf`.`set_type_code` like 'text%') then cast(`bf`.`created_at` as date) end) AS `text_date` from ((((((((select `bible_filesets`.`set_type_code` AS `set_type_code`,`bible_filesets`.`created_at` AS `created_at`,`bible_filesets`.`hash_id` AS `hash_id` from `bible_filesets` where ((`bible_filesets`.`content_loaded` = 1) and (`bible_filesets`.`archived` = 0))) `bf` join (select `bible_fileset_connections`.`bible_id` AS `bible_id`,`bible_fileset_connections`.`hash_id` AS `hash_id` from `bible_fileset_connections`) `bfc` on((`bf`.`hash_id` = `bfc`.`hash_id`))) join (select `bibles`.`id` AS `id`,`bibles`.`language_id` AS `language_id` from `bibles`) `b` on((`bfc`.`bible_id` = `b`.`id`))) join (select `languages`.`id` AS `id`,`languages`.`name` AS `name`,`languages`.`population` AS `population`,`languages`.`country_id` AS `country_id`,`languages`.`iso` AS `iso` from `languages`) `l` on((`l`.`id` = `b`.`language_id`))) join (select `access_group_filesets`.`hash_id` AS `hash_id`,`access_group_filesets`.`access_group_id` AS `access_group_id` from `access_group_filesets`) `agf` on((`bf`.`hash_id` = `agf`.`hash_id`))) join (select `bible_translations`.`name` AS `name`,`bible_translations`.`bible_id` AS `bible_id`,`bible_translations`.`vernacular` AS `vernacular` from `bible_translations` where (`bible_translations`.`vernacular` = 0)) `bt` on((`bt`.`bible_id` = `b`.`id`))) join (select `language_translations`.`name` AS `name`,`language_translations`.`priority` AS `priority`,`language_translations`.`language_translation_id` AS `language_translation_id`,`language_translations`.`language_source_id` AS `language_source_id` from `language_translations` where ((`language_translations`.`priority` = 9) and (`language_translations`.`language_translation_id` = 6414))) `lt` on((`l`.`id` = `lt`.`language_source_id`))) left join (select `countries`.`id` AS `country_id`,`countries`.`name` AS `country` from `countries`) `c` on((`l`.`country_id` = `c`.`country_id`)))) `A` group by `A`.`lang_name`,`A`.`bible_id` order by `video_date` desc */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `stocknumber_lookup`
--

/*!50001 DROP TABLE IF EXISTS `stocknumber_lookup`*/;
/*!50001 DROP VIEW IF EXISTS `stocknumber_lookup`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `stocknumber_lookup` AS select distinct `bft`.`description` AS `stocknumber`,`b`.`id` AS `bibleid`,substring_index(`bfs`.`set_type_code`,'_',1) AS `mode`,`l`.`name` AS `language`,`bt`.`name` AS `version`,`ot`.`name` AS `licensor`,`l`.`iso` AS `iso`,`l`.`rolv_code` AS `rolv`,`l`.`id` AS `languageId`,`l`.`country_id` AS `countryId`,`bfco`.`organization_id` AS `organization_id` from ((((((((`bible_filesets` `bfs` join `bible_fileset_connections` `bfc` on((`bfc`.`hash_id` = `bfs`.`hash_id`))) join `bible_fileset_copyright_organizations` `bfco` on((`bfco`.`hash_id` = `bfs`.`hash_id`))) join `organizations` `o` on((`o`.`id` = `bfco`.`organization_id`))) join `bibles` `b` on((`b`.`id` = `bfc`.`bible_id`))) join `languages` `l` on((`l`.`id` = `b`.`language_id`))) join `bible_fileset_tags` `bft` on((`bfs`.`hash_id` = `bft`.`hash_id`))) join `organization_translations` `ot` on((`ot`.`organization_id` = `bfco`.`organization_id`))) join `bible_translations` `bt` on((`bt`.`bible_id` = `b`.`id`))) where ((`bft`.`name` = 'stock_no') and (`bt`.`language_id` = 6414) and (`ot`.`language_id` = 6414) and (`bfco`.`organization_role` = 1)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `temp_language_equivalents_full`
--

/*!50001 DROP TABLE IF EXISTS `temp_language_equivalents_full`*/;
/*!50001 DROP VIEW IF EXISTS `temp_language_equivalents_full`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `temp_language_equivalents_full` AS select `l`.`iso` AS `iso`,`l`.`country_id` AS `country_id1`,`l`.`id` AS `id1`,`l`.`fcbh_name` AS `fcbh_name1`,`l`.`heart_name` AS `heart_name1`,`l`.`ethnologue_name` AS `ethnologue_name1`,`l`.`official_name` AS `official_name1`,`l`.`rolv_code` AS `rolv_code1`,`l`.`pseudonym` AS `pseudonym1`,`l`.`sensitivity` AS `sensitivity1`,`l2`.`country_id` AS `country_id2`,`l2`.`id` AS `id2`,`l2`.`fcbh_name` AS `fcbh_name2`,`l2`.`heart_name` AS `heart_name2`,`l2`.`ethnologue_name` AS `ethnologue_name2`,`l2`.`official_name` AS `official_name2`,`l2`.`rolv_code` AS `rolv_code2`,`l2`.`pseudonym` AS `pseudonym2`,`l2`.`sensitivity` AS `sensitivity2` from ((`language_view` `l` join `language_view` `l2` on((`l2`.`iso` = `l`.`iso`))) join `temp_language_equivalents` `le` on(((`l`.`id` = `le`.`preferred_id`) and (`l2`.`id` = `le`.`duplicate_id`)))) where ((`l`.`rolv_code` <> `l2`.`rolv_code`) or (`l`.`official_name` <> `l2`.`official_name`)) order by `l`.`iso` */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `test`
--

/*!50001 DROP TABLE IF EXISTS `test`*/;
/*!50001 DROP VIEW IF EXISTS `test`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `test` AS select `A`.`iso` AS `iso`,`A`.`lang_name` AS `lang_name`,`A`.`bible_name` AS `bible_name`,`A`.`bible_id` AS `bible_id`,`A`.`population` AS `population`,`A`.`country` AS `country`,max(`A`.`video_date`) AS `video_date`,max(`A`.`audio_date`) AS `audio_date`,max(`A`.`text_date`) AS `text_date` from (select `l`.`iso` AS `iso`,`lt`.`name` AS `lang_name`,`l`.`population` AS `population`,`c`.`country` AS `country`,`bfc`.`bible_id` AS `bible_id`,`bf`.`set_type_code` AS `media_type`,`bt`.`name` AS `bible_name`,(case when (`bf`.`set_type_code` like 'video%') then cast(`bf`.`created_at` as date) end) AS `video_date`,(case when (`bf`.`set_type_code` like 'audio%') then cast(`bf`.`created_at` as date) end) AS `audio_date`,(case when (`bf`.`set_type_code` like 'text%') then cast(`bf`.`created_at` as date) end) AS `text_date` from ((((((((select `bible_filesets`.`set_type_code` AS `set_type_code`,`bible_filesets`.`created_at` AS `created_at`,`bible_filesets`.`hash_id` AS `hash_id` from `bible_filesets`) `bf` join (select `bible_fileset_connections`.`bible_id` AS `bible_id`,`bible_fileset_connections`.`hash_id` AS `hash_id` from `bible_fileset_connections`) `bfc` on((`bf`.`hash_id` = `bfc`.`hash_id`))) join (select `bibles`.`id` AS `id`,`bibles`.`language_id` AS `language_id` from `bibles`) `b` on((`bfc`.`bible_id` = `b`.`id`))) join (select `languages`.`id` AS `id`,`languages`.`name` AS `name`,`languages`.`population` AS `population`,`languages`.`country_id` AS `country_id`,`languages`.`iso` AS `iso` from `languages`) `l` on((`l`.`id` = `b`.`language_id`))) join (select `access_group_filesets`.`hash_id` AS `hash_id`,`access_group_filesets`.`access_group_id` AS `access_group_id` from `access_group_filesets`) `agf` on((`bf`.`hash_id` = `agf`.`hash_id`))) join (select `bible_translations`.`name` AS `name`,`bible_translations`.`bible_id` AS `bible_id`,`bible_translations`.`vernacular` AS `vernacular` from `bible_translations` where (`bible_translations`.`vernacular` = 0)) `bt` on((`bt`.`bible_id` = `b`.`id`))) join (select `language_translations`.`name` AS `name`,`language_translations`.`priority` AS `priority`,`language_translations`.`language_translation_id` AS `language_translation_id`,`language_translations`.`language_source_id` AS `language_source_id` from `language_translations` where ((`language_translations`.`priority` = 9) and (`language_translations`.`language_translation_id` = 6414))) `lt` on((`l`.`id` = `lt`.`language_source_id`))) left join (select `countries`.`id` AS `country_id`,`countries`.`name` AS `country` from `countries`) `c` on((`l`.`country_id` = `c`.`country_id`)))) `A` group by `A`.`lang_name`,`A`.`bible_id` order by `video_date` desc */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `text_storage_view`
--

/*!50001 DROP TABLE IF EXISTS `text_storage_view`*/;
/*!50001 DROP VIEW IF EXISTS `text_storage_view`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `text_storage_view` AS select `bfs`.`asset_id` AS `bucket`,substring_index(`bfs`.`set_type_code`,'_',1) AS `type`,`bfs`.`set_type_code` AS `set_type_code`,`bfc`.`bible_id` AS `bible`,`bfs`.`id` AS `fileset`,`bf`.`file_name` AS `file_name`,concat('text/',`bfc`.`bible_id`,'/',`bfs`.`id`,'/',`bf`.`file_name`) AS `path`,`bf`.`file_size` AS `size`,`bf`.`updated_at` AS `last_modified_date` from ((`bible_files` `bf` join `bible_filesets` `bfs` on((`bf`.`hash_id` = `bfs`.`hash_id`))) join `bible_fileset_connections` `bfc` on((`bfs`.`hash_id` = `bfc`.`hash_id`))) where ((`bfs`.`asset_id` = 'dbp-prod') and (substring_index(`bfs`.`set_type_code`,'_',1) = 'text')) order by `bfs`.`asset_id`,substring_index(`bfs`.`set_type_code`,'_',1) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `video_storage_view1`
--

/*!50001 DROP TABLE IF EXISTS `video_storage_view1`*/;
/*!50001 DROP VIEW IF EXISTS `video_storage_view1`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `video_storage_view1` AS select `bfs`.`asset_id` AS `bucket`,substring_index(`bfs`.`set_type_code`,'_',1) AS `type`,`bfs`.`set_type_code` AS `set_type_code`,`bfc`.`bible_id` AS `bible`,`bfs`.`id` AS `fileset`,`bf`.`file_name` AS `file_name`,concat('video/',`bfc`.`bible_id`,'/',`bfs`.`id`,'/',`bf`.`file_name`) AS `path`,`bf`.`file_size` AS `size`,`bf`.`updated_at` AS `last_modified_date` from ((`bible_files` `bf` join `bible_filesets` `bfs` on((`bf`.`hash_id` = `bfs`.`hash_id`))) join `bible_fileset_connections` `bfc` on((`bfs`.`hash_id` = `bfc`.`hash_id`))) where ((`bfs`.`asset_id` = 'dbp-vid') and (substring_index(`bfs`.`set_type_code`,'_',1) = 'video')) order by `bfc`.`bible_id`,`bfs`.`id` */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `video_storage_view2`
--

/*!50001 DROP TABLE IF EXISTS `video_storage_view2`*/;
/*!50001 DROP VIEW IF EXISTS `video_storage_view2`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `video_storage_view2` AS select `bfs`.`asset_id` AS `bucket`,substring_index(`bfs`.`set_type_code`,'_',1) AS `type`,`bfs`.`set_type_code` AS `set_type_code`,`bfc`.`bible_id` AS `bible`,`bfs`.`id` AS `fileset`,`bfsb`.`file_name` AS `file_name`,concat('video/',`bfc`.`bible_id`,'/',`bfs`.`id`,'/',`bfsb`.`file_name`) AS `path`,NULL AS `size`,`bfsb`.`updated_at` AS `last_modified_date` from (((`bible_files` `bf` join `bible_filesets` `bfs` on((`bf`.`hash_id` = `bfs`.`hash_id`))) join `bible_fileset_connections` `bfc` on((`bfs`.`hash_id` = `bfc`.`hash_id`))) join `bible_file_stream_bandwidths` `bfsb` on((`bfsb`.`bible_file_id` = `bf`.`id`))) where ((`bfs`.`asset_id` = 'dbp-vid') and (substring_index(`bfs`.`set_type_code`,'_',1) = 'video')) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `video_storage_view3`
--

/*!50001 DROP TABLE IF EXISTS `video_storage_view3`*/;
/*!50001 DROP VIEW IF EXISTS `video_storage_view3`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_0900_ai_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`sa`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `video_storage_view3` AS select `bfs`.`asset_id` AS `bucket`,substring_index(`bfs`.`set_type_code`,'_',1) AS `type`,`bfs`.`set_type_code` AS `set_type_code`,`bfc`.`bible_id` AS `bible`,`bfs`.`id` AS `fileset`,`bfst`.`file_name` AS `file_name`,concat('video/',`bfc`.`bible_id`,'/',`bfs`.`id`,'/',`bfst`.`file_name`) AS `path`,NULL AS `size`,`bfsb`.`updated_at` AS `last_modified_date` from ((((`bible_files` `bf` join `bible_filesets` `bfs` on((`bf`.`hash_id` = `bfs`.`hash_id`))) join `bible_fileset_connections` `bfc` on((`bfs`.`hash_id` = `bfc`.`hash_id`))) join `bible_file_stream_bandwidths` `bfsb` on((`bfsb`.`bible_file_id` = `bf`.`id`))) join `bible_file_stream_ts` `bfst` on((`bfst`.`stream_bandwidth_id` = `bfsb`.`id`))) where ((`bfs`.`asset_id` = 'dbp-vid') and (substring_index(`bfs`.`set_type_code`,'_',1) = 'video')) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-01-21 20:41:03
