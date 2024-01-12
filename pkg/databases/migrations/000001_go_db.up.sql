BEGIN;

--Set timezone
SET TIME ZONE 'Asia/Bangkok';

--Install uuid extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--user_id -> U000001
--product_id -> P000001
--order_id -> O000001
--Create sequence
CREATE SEQUENCE user_id_seq START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE product_id_seq START WITH 1 INCREMENT BY 1;
CREATE SEQUENCE order_id_seq START WITH 1 INCREMENT BY 1;

--Auto update
CREATE OR REPLACE FUNCTION set_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

--Create enum
CREATE TYPE "order_status" AS ENUM (
    'waiting',
    'shipping',
    'completed',
    'canceled'
);

CREATE TABLE "users" (
  "id" VARCHAR(7) PRIMARY KEY DEFAULT CONCAT('U', LPAD(NEXTVAL('user_id_seq')::TEXT, 6, '0')),
  "username" VARCHAR UNIQUE NOT NULL,
  "password" VARCHAR NOT NULL,
  "email" VARCHAR UNIQUE NOT NULL,
  "role_id" INT NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "oauth" (
  "id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "user_id" VARCHAR NOT NULL,
  "access_token" VARCHAR NOT NULL,
  "refresh_token" VARCHAR NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "role" (
  "id" SERIAL PRIMARY KEY,
  "title" VARCHAR NOT NULL UNIQUE
);

CREATE TABLE "product" (
  "id" VARCHAR(7) PRIMARY KEY DEFAULT CONCAT('P', LPAD(NEXTVAL('product_id_seq')::TEXT, 6, '0')),
  "title" VARCHAR NOT NULL,
  "description" VARCHAR NOT NULL DEFAULT '',
  "price" FLOAT NOT NULL DEFAULT 0,
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "image" (
  "id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "filename" VARCHAR,
  "url" VARCHAR,
  "product_id" VARCHAR,
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "product_category" (
  "id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "product_id" VARCHAR NOT NULL,
  "category_id" INT NOT NULL
);

CREATE TABLE "categories" (
  "id" SERIAL PRIMARY KEY,
  "title" VARCHAR UNIQUE NOT NULL
);

CREATE TABLE "order" (
  "id" VARCHAR(7) PRIMARY KEY DEFAULT CONCAT('O', LPAD(NEXTVAL('order_id_seq')::TEXT, 6, '0')),
  "user_id" VARCHAR NOT NULL,
  "contact" VARCHAR NOT NULL,
  "address" VARCHAR NOT NULL,
  "transfer_slip" jsonb,
  "status" order_status NOT NULL,
  "created_at" TIMESTAMP NOT NULL DEFAULT now(),
  "updated_at" TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE "product_order" (
  "id" uuid NOT NULL UNIQUE PRIMARY KEY DEFAULT uuid_generate_v4(),
  "order_id" VARCHAR NOT NULL,
  "quantity" INT NOT NULL DEFAULT 1,
  "products" jsonb
);

ALTER TABLE "users" ADD FOREIGN KEY ("role_id") REFERENCES "role" ("id");
ALTER TABLE "oauth" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "image" ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id");
ALTER TABLE "product_category" ADD FOREIGN KEY ("product_id") REFERENCES "product" ("id");
ALTER TABLE "product_category" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id");
ALTER TABLE "order" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "product_order" ADD FOREIGN KEY ("order_id") REFERENCES "order" ("id");

CREATE TRIGGER set_updated_at_timestamp_users_table BEFORE UPDATE ON "users" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_oauth_table BEFORE UPDATE ON "oauth" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_product_table BEFORE UPDATE ON "product" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_image_table BEFORE UPDATE ON "image" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();
CREATE TRIGGER set_updated_at_timestamp_order_table BEFORE UPDATE ON "order" FOR EACH ROW EXECUTE PROCEDURE set_updated_at_column();


COMMIT;