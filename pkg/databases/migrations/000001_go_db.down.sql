BEGIN;

DROP TRIGGER IF EXISTS set_updated_at_timestamp_users_table ON "users";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_oauth_table ON "oauth";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_product_table ON "product";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_image_table ON "image";
DROP TRIGGER IF EXISTS set_updated_at_timestamp_order_table ON "order";

DROP FUNCTION IF EXISTS set_updated_at_column();

DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "oauth" CASCADE;
DROP TABLE IF EXISTS "role" CASCADE;
DROP TABLE IF EXISTS "product" CASCADE;
DROP TABLE IF EXISTS "image" CASCADE;
DROP TABLE IF EXISTS "product_category" CASCADE;
DROP TABLE IF EXISTS "categories" CASCADE;
DROP TABLE IF EXISTS "order" CASCADE;
DROP TABLE IF EXISTS "product_order" CASCADE;

DROP SEQUENCE IF EXISTS user_id_seq;
DROP SEQUENCE IF EXISTS product_id_seq;
DROP SEQUENCE IF EXISTS order_id_seq;

DROP TYPE IF EXISTS "order_status";

COMMIT;