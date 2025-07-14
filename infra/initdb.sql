CREATE DATABASE user_db;
\c user_db

CREATE TABLE "public"."users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "business_id" uuid,
  "name" text,
  "role" text NOT NULL DEFAULT 'user',
  "email" text,
  CONSTRAINT "users_pkey" PRIMARY KEY ("id"),
  CONSTRAINT "users_email_unique" UNIQUE ("email", "business_id")
);