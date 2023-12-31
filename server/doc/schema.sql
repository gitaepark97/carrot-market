-- SQL dump generated using DBML (dbml-lang.org)
-- Database: PostgreSQL
-- Generated at: 2023-07-19T02:10:20.651Z

CREATE TABLE "users" (
  "user_id" serial PRIMARY KEY,
  "email" varchar(50) UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "nickname" varchar(50) UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "sessions" (
  "session_id" uuid PRIMARY KEY,
  "user_id" int NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT (false),
  "expired_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "goods" (
  "goods_id" serial PRIMARY KEY,
  "user_id" int NOT NULL,
  "title" varchar(50) NOT NULL,
  "price" int NOT NULL,
  "description" text NOT NULL,
  "default_image_url" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "categories" (
  "category_id" serial PRIMARY KEY,
  "title" varchar(50) UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "goods_categories" (
  "goods_id" int,
  "category_id" int,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  PRIMARY KEY ("goods_id", "category_id")
);

CREATE TABLE "goods_images" (
  "goods_image_id" serial PRIMARY KEY,
  "goods_id" int NOT NULL,
  "image_url" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "goods_images" ("goods_id", "image_url");

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;

ALTER TABLE "goods" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;

ALTER TABLE "goods_categories" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("category_id");

ALTER TABLE "goods_categories" ADD FOREIGN KEY ("goods_id") REFERENCES "goods" ("goods_id") ON DELETE CASCADE;

ALTER TABLE "goods_images" ADD FOREIGN KEY ("goods_id") REFERENCES "goods" ("goods_id") ON DELETE CASCADE;
