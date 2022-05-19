/*
 Navicat Premium Data Transfer

 Source Server         : Local Postres DB
 Source Server Type    : PostgreSQL
 Source Server Version : 140002
 Source Host           : localhost:5432
 Source Catalog        : crud_go
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 140002
 File Encoding         : 65001

 Date: 19/05/2022 13:54:17
*/


-- ----------------------------
-- Sequence structure for category_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."category_id_seq";
CREATE SEQUENCE "public"."category_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;
ALTER SEQUENCE "public"."category_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for customer_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."customer_id_seq";
CREATE SEQUENCE "public"."customer_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;
ALTER SEQUENCE "public"."customer_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for employee_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."employee_id_seq";
CREATE SEQUENCE "public"."employee_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;
ALTER SEQUENCE "public"."employee_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for menu_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."menu_id_seq";
CREATE SEQUENCE "public"."menu_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;
ALTER SEQUENCE "public"."menu_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for notes_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."notes_id_seq";
CREATE SEQUENCE "public"."notes_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;
ALTER SEQUENCE "public"."notes_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for order_detail_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."order_detail_id_seq";
CREATE SEQUENCE "public"."order_detail_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;
ALTER SEQUENCE "public"."order_detail_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for orders_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."orders_id_seq";
CREATE SEQUENCE "public"."orders_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;
ALTER SEQUENCE "public"."orders_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for position_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."position_id_seq";
CREATE SEQUENCE "public"."position_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;
ALTER SEQUENCE "public"."position_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Sequence structure for users_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."users_id_seq";
CREATE SEQUENCE "public"."users_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;
ALTER SEQUENCE "public"."users_id_seq" OWNER TO "postgres";

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS "public"."users";
CREATE TABLE "public"."users" (
  "id" int4 NOT NULL DEFAULT nextval('users_id_seq'::regclass),
  "username" varchar(45) COLLATE "pg_catalog"."default" DEFAULT NULL::character varying,
  "email" varchar(45) COLLATE "pg_catalog"."default" DEFAULT NULL::character varying,
  "password" varchar(100) COLLATE "pg_catalog"."default" DEFAULT NULL::character varying,
  "created_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "created_by" varchar(45) COLLATE "pg_catalog"."default" DEFAULT NULL::character varying,
  "updated_at" timestamp(0) DEFAULT NULL::timestamp without time zone,
  "last_update_by" varchar(100) COLLATE "pg_catalog"."default" DEFAULT NULL::character varying,
  "is_deleted" bool
)
;
ALTER TABLE "public"."users" OWNER TO "postgres";

-- ----------------------------
-- Records of users
-- ----------------------------
BEGIN;
INSERT INTO "public"."users" VALUES (1, 'nazyli', 'evrynazyli@gmail.com', '$2a$10$BGMEQtjvlQB9/8lQXryjvez1Xug.XuZJ3P1ajVCiNFbh4qYIL91my', '2020-01-01 02:02:04', 'eb55808b848359c7566d41a69d712cc7d421dca3', '2020-06-12 06:55:37', 'eb55808b848359c7566d41a69d712cc7d421dca3', 'f');
INSERT INTO "public"."users" VALUES (12, 'evryy', 'evry@gmail.com', '$2a$10$DMhQVzpKnnaGIw/HfzFZ7ODExA5bE24YufFJgBfSAfDC0rl.Fla1C', '2022-05-19 13:51:23', '', '2022-05-19 13:51:23', '', 'f');
COMMIT;

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."category_id_seq"', 3, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."customer_id_seq"', 3, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."employee_id_seq"', 3, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."menu_id_seq"', 3, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."notes_id_seq"', 3, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."order_detail_id_seq"', 3, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."orders_id_seq"', 3, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."position_id_seq"', 3, false);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
SELECT setval('"public"."users_id_seq"', 13, true);

-- ----------------------------
-- Primary Key structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "users_pkey" PRIMARY KEY ("id");
