-- Create "pets" table
CREATE TABLE "pets" ("id" uuid NOT NULL, "public_id" character varying NOT NULL, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "name" character varying NOT NULL, "user_id" uuid NOT NULL, PRIMARY KEY ("id"));
-- Create index "pets_public_id_key" to table: "pets"
CREATE UNIQUE INDEX "pets_public_id_key" ON "pets" ("public_id");
-- Create "users" table
CREATE TABLE "users" ("id" uuid NOT NULL, "public_id" character varying NOT NULL, "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "name" character varying NOT NULL, "email" character varying NOT NULL, PRIMARY KEY ("id"));
-- Create index "users_public_id_key" to table: "users"
CREATE UNIQUE INDEX "users_public_id_key" ON "users" ("public_id");
