/*
 Navicat Premium Data Transfer

 Source Server         : postgresql_local
 Source Server Type    : PostgreSQL
 Source Server Version : 140003
 Source Host           : localhost:5432
 Source Catalog        : database_configuration
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 140003
 File Encoding         : 65001

 Date: 22/06/2022 15:49:55
*/


-- ----------------------------
-- Table structure for apps
-- ----------------------------
DROP TABLE IF EXISTS "public"."apps";
CREATE TABLE "public"."apps" (
  "id" uuid NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default",
  "created_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "created_by" varchar(255) COLLATE "pg_catalog"."default",
  "modified_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "modified_by" varchar(255) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for dbs
-- ----------------------------
DROP TABLE IF EXISTS "public"."dbs";
CREATE TABLE "public"."dbs" (
  "id" uuid NOT NULL,
  "app_id" uuid,
  "name" varchar(255) COLLATE "pg_catalog"."default",
  "host" varchar(255) COLLATE "pg_catalog"."default",
  "username" varchar(255) COLLATE "pg_catalog"."default",
  "password" varchar(255) COLLATE "pg_catalog"."default",
  "dialect" varchar(255) COLLATE "pg_catalog"."default",
  "created_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "created_by" varchar(255) COLLATE "pg_catalog"."default",
  "modified_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "modified_by" varchar(255) COLLATE "pg_catalog"."default",
  "port" int8
)
;

-- ----------------------------
-- Table structure for services
-- ----------------------------
DROP TABLE IF EXISTS "public"."services";
CREATE TABLE "public"."services" (
  "id" uuid NOT NULL,
  "db_id" uuid,
  "method" varchar(255) COLLATE "pg_catalog"."default",
  "service_url" text COLLATE "pg_catalog"."default",
  "service_definition" text COLLATE "pg_catalog"."default",
  "is_query" bool,
  "created_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "created_by" varchar(255) COLLATE "pg_catalog"."default",
  "modified_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "modified_by" varchar(255) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for user_service
-- ----------------------------
DROP TABLE IF EXISTS "public"."user_service";
CREATE TABLE "public"."user_service" (
  "id" uuid NOT NULL,
  "user_id" uuid,
  "service_id" uuid,
  "created_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "created_by" varchar(255) COLLATE "pg_catalog"."default",
  "modified_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "modified_by" varchar(255) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS "public"."users";
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "username" varchar(255) COLLATE "pg_catalog"."default",
  "password" varchar(255) COLLATE "pg_catalog"."default",
  "email" varchar(255) COLLATE "pg_catalog"."default",
  "created_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "created_by" varchar(255) COLLATE "pg_catalog"."default",
  "modified_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "modified_by" varchar(255) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Primary Key structure for table apps
-- ----------------------------
ALTER TABLE "public"."apps" ADD CONSTRAINT "pk_app" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table dbs
-- ----------------------------
ALTER TABLE "public"."dbs" ADD CONSTRAINT "pk_db" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table services
-- ----------------------------
ALTER TABLE "public"."services" ADD CONSTRAINT "pk_service" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table user_service
-- ----------------------------
ALTER TABLE "public"."user_service" ADD CONSTRAINT "pk_user_service" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "pk_user" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table dbs
-- ----------------------------
ALTER TABLE "public"."dbs" ADD CONSTRAINT "fk_db_app" FOREIGN KEY ("app_id") REFERENCES "public"."apps" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table services
-- ----------------------------
ALTER TABLE "public"."services" ADD CONSTRAINT "fk_service_db" FOREIGN KEY ("db_id") REFERENCES "public"."dbs" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table user_service
-- ----------------------------
ALTER TABLE "public"."user_service" ADD CONSTRAINT "fk_service_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."user_service" ADD CONSTRAINT "fk_user_service" FOREIGN KEY ("service_id") REFERENCES "public"."services" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
