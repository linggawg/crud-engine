-- ----------------------------
-- Table structure for apps
-- ----------------------------
CREATE TABLE "apps" (
  "id" uuid NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default"
);
-- ----------------------------
-- Primary Key structure for table apps
-- ----------------------------
ALTER TABLE "apps" ADD CONSTRAINT "pk_app" PRIMARY KEY ("id");

-- ----------------------------
-- Table structure for dbs
-- ----------------------------
CREATE TABLE "dbs" (
   "id" uuid NOT NULL,
   "app_id" uuid,
   "name" varchar(255) COLLATE "pg_catalog"."default",
   "host" varchar(255) COLLATE "pg_catalog"."default",
   "username" varchar(255) COLLATE "pg_catalog"."default",
   "password" varchar(255) COLLATE "pg_catalog"."default",
   "dialect" varchar(255) COLLATE "pg_catalog"."default"
);
-- ----------------------------
-- Primary Key structure for table dbs
-- ----------------------------
ALTER TABLE "dbs" ADD CONSTRAINT "pk_db" PRIMARY KEY ("id");
-- ----------------------------
-- Foreign Keys structure for table dbs
-- ----------------------------
ALTER TABLE "dbs" ADD CONSTRAINT "fk_db_app" FOREIGN KEY ("app_id") REFERENCES "apps" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Table structure for services
-- ----------------------------
CREATE TABLE "services" (
   "id" uuid NOT NULL,
   "db_id" uuid,
   "method" varchar(255) COLLATE "pg_catalog"."default",
   "service_url" varchar(255) COLLATE "pg_catalog"."default",
   "service_definition" varchar(255) COLLATE "pg_catalog"."default",
   "is_query" bool
);
-- ----------------------------
-- Primary Key structure for table services
-- ----------------------------
ALTER TABLE "services" ADD CONSTRAINT "pk_service" PRIMARY KEY ("id");
-- ----------------------------
-- Foreign Keys structure for table services
-- ----------------------------
ALTER TABLE "services" ADD CONSTRAINT "fk_service_db" FOREIGN KEY ("db_id") REFERENCES "dbs" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Table structure for users
-- ----------------------------
CREATE TABLE "users" (
  "id" uuid NOT NULL,
  "username" varchar(255) COLLATE "pg_catalog"."default",
  "password" varchar(255) COLLATE "pg_catalog"."default"
);
-- ----------------------------
-- Primary Key structure for table users
-- ----------------------------
ALTER TABLE "users" ADD CONSTRAINT "pk_user" PRIMARY KEY ("id");

-- ----------------------------
-- Table structure for user_service
-- ----------------------------
CREATE TABLE "user_service" (
   "id" uuid NOT NULL,
   "user_id" uuid,
   "service_id" uuid
);
-- ----------------------------
-- Primary Key structure for table user_service
-- ----------------------------
ALTER TABLE "user_service" ADD CONSTRAINT "pk_user_service" PRIMARY KEY ("id");
-- ----------------------------
-- Foreign Keys structure for table user_service
-- ----------------------------
ALTER TABLE "user_service" ADD CONSTRAINT "fk_service_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "user_service" ADD CONSTRAINT "fk_user_service" FOREIGN KEY ("service_id") REFERENCES "services" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
