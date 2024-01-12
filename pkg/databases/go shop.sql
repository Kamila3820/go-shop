CREATE TABLE "users" (
  "id" varchar PRIMARY KEY,
  "username" varchar UNIQUE,
  "password" varchar,
  "email" varchar UNIQUE,
  "role_id" int,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "oauth" (
  "id" varchar PRIMARY KEY,
  "user_id" varchar,
  "access_token" varchar,
  "refresh_token" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "role" (
  "id" int PRIMARY KEY,
  "title" varchar
);

CREATE TABLE "product" (
  "id" varchar PRIMARY KEY,
  "title" varchar,
  "description" varchar,
  "price" float,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "image" (
  "id" varchar PRIMARY KEY,
  "filename" varchar,
  "url" varchar,
  "product_id" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "product_category" (
  "id" varchar PRIMARY KEY,
  "product_id" varchar,
  "category_id" int
);

CREATE TABLE "categories" (
  "id" int PRIMARY KEY,
  "title" varchar UNIQUE
);

CREATE TABLE "order" (
  "id" varchar PRIMARY KEY,
  "user_id" varchar,
  "contact" varchar,
  "address" varchar,
  "transfer_slip" json,
  "status" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE "product_order" (
  "id" varchar PRIMARY KEY,
  "order_id" varchar,
  "quantity" int,
  "products" json
);

ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "role" ("id");

ALTER TABLE "oauth" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "image" ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id");

ALTER TABLE "product_category" ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id");

ALTER TABLE "product_category" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");

ALTER TABLE "order" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "product_order" ADD FOREIGN KEY ("order_id") REFERENCES "order" ("id");
