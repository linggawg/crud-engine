/*
 Navicat Premium Data Transfer

 Source Server         : postgresql_local
 Source Server Type    : PostgreSQL
 Source Server Version : 140003
 Source Host           : localhost:5432
 Source Catalog        : db_config
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 140003
 File Encoding         : 65001

 Date: 06/10/2022 15:54:45
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
  "port" int8,
  "username" varchar(255) COLLATE "pg_catalog"."default",
  "password" varchar(255) COLLATE "pg_catalog"."default",
  "dialect" varchar(255) COLLATE "pg_catalog"."default",
  "created_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "created_by" varchar(255) COLLATE "pg_catalog"."default",
  "modified_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "modified_by" varchar(255) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for queries
-- ----------------------------
DROP TABLE IF EXISTS "public"."queries";
CREATE TABLE "public"."queries" (
  "id" uuid NOT NULL,
  "query_definition" text COLLATE "pg_catalog"."default",
  "created_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "created_by" varchar(255) COLLATE "pg_catalog"."default",
  "modified_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "modified_by" varchar(255) COLLATE "pg_catalog"."default",
  "key" text COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for resources_mapping
-- ----------------------------
DROP TABLE IF EXISTS "public"."resources_mapping";
CREATE TABLE "public"."resources_mapping" (
  "id" uuid NOT NULL,
  "service_id" uuid,
  "source_origin" varchar(255) COLLATE "pg_catalog"."default",
  "source_alias" varchar(255) COLLATE "pg_catalog"."default",
  "created_at" timestamp(6) DEFAULT now(),
  "created_by" varchar(255) COLLATE "pg_catalog"."default",
  "modified_at" timestamp(6),
  "modified_by" varchar(255) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for roles
-- ----------------------------
DROP TABLE IF EXISTS "public"."roles";
CREATE TABLE "public"."roles" (
  "id" uuid NOT NULL,
  "name" text COLLATE "pg_catalog"."default",
  "created_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "created_by" varchar(255) COLLATE "pg_catalog"."default",
  "modified_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "modified_by" varchar(255) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for services
-- ----------------------------
DROP TABLE IF EXISTS "public"."services";
CREATE TABLE "public"."services" (
  "id" uuid NOT NULL,
  "db_id" uuid,
  "query_id" uuid,
  "method" varchar(255) COLLATE "pg_catalog"."default",
  "service_url" text COLLATE "pg_catalog"."default",
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
  "role_id" uuid NOT NULL,
  "username" varchar(255) COLLATE "pg_catalog"."default",
  "password" varchar(255) COLLATE "pg_catalog"."default",
  "created_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "created_by" varchar(255) COLLATE "pg_catalog"."default",
  "modified_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "modified_by" varchar(255) COLLATE "pg_catalog"."default"
)
;

-- ----------------------------
-- Table structure for users_services
-- ----------------------------
DROP TABLE IF EXISTS "public"."users_services";
CREATE TABLE "public"."users_services" (
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
-- Primary Key structure for table apps
-- ----------------------------
ALTER TABLE "public"."apps" ADD CONSTRAINT "pk_apps" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table dbs
-- ----------------------------
ALTER TABLE "public"."dbs" ADD CONSTRAINT "pk_dbs" PRIMARY KEY ("id");

-- ----------------------------
-- Uniques structure for table queries
-- ----------------------------
ALTER TABLE "public"."queries" ADD CONSTRAINT "unique_key" UNIQUE ("key");

-- ----------------------------
-- Primary Key structure for table queries
-- ----------------------------
ALTER TABLE "public"."queries" ADD CONSTRAINT "pk_queries" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table resources_mapping
-- ----------------------------
ALTER TABLE "public"."resources_mapping" ADD CONSTRAINT "resources_mapping_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table roles
-- ----------------------------
ALTER TABLE "public"."roles" ADD CONSTRAINT "pk_roles" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table services
-- ----------------------------
ALTER TABLE "public"."services" ADD CONSTRAINT "pk_services" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "pk_users" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table users_services
-- ----------------------------
ALTER TABLE "public"."users_services" ADD CONSTRAINT "pk_users_services" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table dbs
-- ----------------------------
ALTER TABLE "public"."dbs" ADD CONSTRAINT "fk_dbs_apps" FOREIGN KEY ("app_id") REFERENCES "public"."apps" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table resources_mapping
-- ----------------------------
ALTER TABLE "public"."resources_mapping" ADD CONSTRAINT "fk_resources_mapping_services" FOREIGN KEY ("service_id") REFERENCES "public"."services" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table services
-- ----------------------------
ALTER TABLE "public"."services" ADD CONSTRAINT "fk_services_dbs" FOREIGN KEY ("db_id") REFERENCES "public"."dbs" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."services" ADD CONSTRAINT "fk_services_queries" FOREIGN KEY ("query_id") REFERENCES "public"."queries" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "fk_users_roles" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table users_services
-- ----------------------------
ALTER TABLE "public"."users_services" ADD CONSTRAINT "fk_services_users" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "public"."users_services" ADD CONSTRAINT "fk_users_services" FOREIGN KEY ("service_id") REFERENCES "public"."services" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
