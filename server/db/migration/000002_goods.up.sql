ALTER TABLE sessions ALTER COLUMN user_id DROP DEFAULT;

ALTER TABLE sessions ALTER COLUMN user_id TYPE int;

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

ALTER TABLE "goods" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id") ON DELETE CASCADE;

ALTER TABLE "goods_categories" ADD FOREIGN KEY ("goods_id") REFERENCES "goods" ("goods_id") ON DELETE CASCADE;

ALTER TABLE "goods_categories" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("category_id");

ALTER TABLE "goods_images" ADD FOREIGN KEY ("goods_id") REFERENCES "goods" ("goods_id") ON DELETE CASCADE;

INSERT INTO categories(
  title
) VALUES
  ('디지털 기기'),
  ('가구/인테리어'),
  ('유아동'),
  ('여성의류'),
  ('여성잡화'),
  ('남성패션/잡화'),
  ('생활가전'),
  ('생활/주방'),
  ('가공식품'),
  ('스포츠/레저'),
  ('취미/게임/음반'),
  ('뷰티/미용'),
  ('식물'),
  ('반려동물용품'),
  ('티켓/교환권'),
  ('도서'),
  ('유아도서'),
  ('기타 중고물품');